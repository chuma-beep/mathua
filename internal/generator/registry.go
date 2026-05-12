package generator

import (
	"fmt"
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
	r.gens[conceptID] = gen
	return nil
}

func (r *Registry) Generate(conceptID string, difficulty float64) (Problem, error) {
	if difficulty < 0 {
		difficulty = 0
	}
	if difficulty > 1 {
		difficulty = 1
	}
	r.mu.RLock()
	gen, exists := r.gens[conceptID]
	r.mu.RUnlock()
	if !exists {
		return Problem{}, fmt.Errorf("no generator registered for concept %q", conceptID)
	}
	return gen.Generate(difficulty), nil
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
