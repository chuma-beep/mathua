package generator

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu   sync.RWMutex
	gens map[string]Generator
}

func NewRegistry() *Registry {
	return &Registry{
		gens: make(map[string]Generator),
	}
}

func (r *Registry) Register(conceptID string, gen Generator) error {
	if conceptID == "" {
		return fmt.Errorf("concept ID must not be empty")
	}
	if gen == nil {
		return fmt.Errorf("generator must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.gens[conceptID]; exists {
		return fmt.Errorf("generator already registered for concept %q", conceptID)
	}
	// Self-check: if this is a GradedGenerator, make sure it can grade its own output.
	if gg, ok := gen.(GradedGenerator); ok {
		p := gg.Generate(GeneratorContext{Difficulty: 0.5})
		result := gg.Grade(p.Answer, p.Answer)
		if !result.Correct {
			return fmt.Errorf("generator for %q cannot grade its own answer (expected=%q): %s",
				conceptID, p.Answer, result.Feedback)
		}
		// Also verify it rejects a clearly wrong answer
		wrongResult := gg.Grade(p.Answer, "")
		if wrongResult.Correct {
			return fmt.Errorf("generator for %q accepted empty answer as correct", conceptID)
		}
	}
	r.gens[conceptID] = gen
	return nil
}


func (r *Registry) Get(conceptID string) (Generator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gen, exists := r.gens[conceptID]
	if !exists {
		return nil, fmt.Errorf("no generator registered for concept %q", conceptID)
	}
	return gen, nil
}

func (r *Registry) Generate(conceptID string, difficulty float64) (Problem, error) {
	if difficulty < 0 {
		difficulty = 0
	}
	if difficulty > 1 {
		difficulty = 1
	}
	return r.GenerateContext(conceptID, GeneratorContext{Difficulty: difficulty})
}

func (r *Registry) GenerateContext(conceptID string, ctx GeneratorContext) (Problem, error) {
	r.mu.RLock()
	gen, exists := r.gens[conceptID]
	r.mu.RUnlock()
	if !exists {
		return Problem{}, fmt.Errorf("no generator registered for concept %q", conceptID)
	}
	return gen.Generate(ctx), nil
}

func (r *Registry) BatchGenerate(conceptID string, count int, difficulty float64) ([]Problem, error) {
	return r.BatchGenerateContext(conceptID, count, GeneratorContext{Difficulty: difficulty})
}

func (r *Registry) BatchGenerateContext(conceptID string, count int, ctx GeneratorContext) ([]Problem, error) {
	if count < 1 {
		count = 1
	}
	if count > 20 {
		count = 20
	}
	r.mu.RLock()
	gen, exists := r.gens[conceptID]
	r.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("no generator registered for concept %q", conceptID)
	}
	problems := make([]Problem, 0, count)
	seen := make(map[string]bool)
	for i := 0; i < count*3 && len(problems) < count; i++ {
		p := gen.Generate(ctx)
		if !seen[p.Question] {
			seen[p.Question] = true
			problems = append(problems, p)
		}
	}
	return problems, nil
}

func (r *Registry) Has(conceptID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.gens[conceptID]
	return exists
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.gens)
}

func (r *Registry) Concepts() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.gens))
	for id := range r.gens {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
