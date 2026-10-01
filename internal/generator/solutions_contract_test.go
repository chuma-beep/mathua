package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/solutions"
)

// solutionsDir is the shipped corpus, relative to this package.
const solutionsDir = "../../data/lessons/solutions"

// TestShippedSchemasResolveAgainstTheirGenerator is the contract that makes a
// schema safe to ship: every placeholder in every authored schema must have a
// value behind it, on a problem the generator actually produces.
//
// A schema and a generator that disagree is a silent content bug — the
// assembler falls back to Problem.Explanation and the authored prose silently
// stops being used, so nothing else would ever notice. This test is what
// notices.
func TestShippedSchemasResolveAgainstTheirGenerator(t *testing.T) {
	l, err := solutions.Load(solutionsDir)
	if err != nil {
		t.Fatalf("load shipped schemas: %v", err)
	}
	ids := l.Concepts()
	if len(ids) == 0 {
		t.Skip("no solution schemas authored yet")
	}
	reg := generator.NewRegistry()
	reg.SetSolutions(l)
	registerAllDomains(reg)

	// Several instances per concept: a generator can publish a fact only on
	// some difficulty paths, and one sample would miss that.
	const samples = 25
	for _, id := range ids {
		schema, ok := l.Schema(id)
		if !ok {
			t.Fatalf("Concepts() listed %q but Schema() has none", id)
		}
		var unresolved []string
		for i := range samples {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i) / float64(samples-1),
				Seed:       int64(1000 + i),
			})
			if err != nil {
				t.Errorf("%s: generate: %v", id, err)
				break
			}
			facts := p.Facts
			if facts == nil {
				facts = map[string]string{}
			}
			facts["answer"] = p.Answer
			for _, name := range solutions.Missing(schema, facts) {
				if !contains(unresolved, name) {
					unresolved = append(unresolved, name)
				}
			}
		}
		if len(unresolved) > 0 {
			t.Errorf("%s: schema references facts the generator never publishes: %s\n"+
				"    add them to the Problem{Facts: ...} literal, or the schema is skipped at runtime",
				id, strings.Join(unresolved, ", "))
		}
	}
}

// TestSchemasNeverLeakPlaceholders guards the last hop: whatever the assembler
// produces is what a learner reads, so a placeholder must never survive into
// a served explanation.
func TestSchemasNeverLeakPlaceholders(t *testing.T) {
	l, err := solutions.Load(solutionsDir)
	if err != nil {
		t.Fatalf("load shipped schemas: %v", err)
	}
	ids := l.Concepts()
	if len(ids) == 0 {
		t.Skip("no solution schemas authored yet")
	}
	reg := generator.NewRegistry()
	reg.SetSolutions(l)
	registerAllDomains(reg)
	for _, id := range ids {
		for i := range 8 {
			p, err := reg.GenerateContext(id, generator.GeneratorContext{
				Difficulty: float64(i) / 8,
				Seed:       int64(2000 + i),
			})
			if err != nil {
				t.Errorf("%s: generate: %v", id, err)
				break
			}
			if strings.Contains(p.Explanation, "[[") {
				t.Errorf("%s: served explanation contains a raw placeholder:\n%s", id, p.Explanation)
			}
			if p.Explanation == "" {
				t.Errorf("%s: served explanation is empty", id)
			}
		}
	}
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}

// TestSchemaReachesTheServedExplanation exercises the whole path end to end with
// a schema that needs nothing but the answer, which every problem has. Without
// this the mechanism would have no coverage until a generator published its
// first fact, and a break in the wiring would only surface mid content wave.
func TestSchemaReachesTheServedExplanation(t *testing.T) {
	dir := t.TempDir()
	schema := `{
	  "concept": "arith.add.single",
	  "setup": "You are adding two numbers.",
	  "steps": [{"fact": "answer", "say": "Add the two numbers together."}],
	  "answer": "The sum is [[answer]]."
	}`
	if err := os.WriteFile(filepath.Join(dir, "arith.add.single.json"), []byte(schema), 0o600); err != nil {
		t.Fatal(err)
	}
	l, err := solutions.Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	reg := generator.NewRegistry()
	registerAllDomains(reg)
	reg.SetSolutions(l)

	p, err := reg.GenerateContext("arith.add.single", generator.GeneratorContext{Difficulty: 0.5, Seed: 7})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if !strings.Contains(p.Explanation, "You are adding two numbers.") {
		t.Errorf("schema prose missing from the served explanation:\n%s", p.Explanation)
	}
	if !strings.Contains(p.Explanation, "Add the two numbers together.") {
		t.Errorf("schema step missing from the served explanation:\n%s", p.Explanation)
	}
	// The answer is interpolated, and the explanation still closes on it.
	if !strings.Contains(p.Explanation, "The sum is "+p.Answer+".") {
		t.Errorf("schema answer step did not interpolate %q:\n%s", p.Answer, p.Explanation)
	}
	// The generator's own text is not what got served.
	if !strings.Contains(p.Explanation, "Add the two numbers together.") || strings.Contains(p.Explanation, "Step 1") {
		t.Errorf("expected the schema, not the generator template:\n%s", p.Explanation)
	}
}

// A schema whose placeholders cannot all resolve must leave the generator's own
// explanation in place, not a half-rendered one.
func TestSchemaWithUnresolvableFactFallsBackToGenerator(t *testing.T) {
	dir := t.TempDir()
	schema := `{
	  "concept": "arith.add.single",
	  "steps": [{"fact": "nonexistent", "say": "Using [[nonexistent]] here."}],
	  "answer": "The sum is [[answer]]."
	}`
	if err := os.WriteFile(filepath.Join(dir, "arith.add.single.json"), []byte(schema), 0o600); err != nil {
		t.Fatal(err)
	}
	l, err := solutions.Load(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	reg := generator.NewRegistry()
	registerAllDomains(reg)
	reg.SetSolutions(l)

	p, err := reg.GenerateContext("arith.add.single", generator.GeneratorContext{Difficulty: 0.5, Seed: 7})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if strings.Contains(p.Explanation, "[[") {
		t.Errorf("an unresolved schema leaked a placeholder:\n%s", p.Explanation)
	}
	if strings.Contains(p.Explanation, "Using ") {
		t.Errorf("an unresolved schema was served instead of the fallback:\n%s", p.Explanation)
	}
}
