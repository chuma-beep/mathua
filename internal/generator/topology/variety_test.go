package topology

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

// TestElevatedTopologyVariety checks the five elevated topology concepts that
// now ask for counts. They used to serve multiple choice, so this test watched
// their option sets; now it watches that the same generators keep producing a
// spread of counts and stems.
func TestElevatedTopologyVariety(t *testing.T) {
	rand.Seed(42)
	reg := generator.NewRegistry()
	Register(reg)

	conceptIDs := []string{
		"topo.homotopy",
		"topo.manifold",
		"topo.connected.local_path",
		"topo.compact.one_point",
		"topo.separation.completely_regular",
	}

	for _, id := range conceptIDs {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("concept %s not registered: %v", id, err)
		}

		answers := make(map[string]bool)
		stems := make(map[string]bool)
		for i := 0; i < 100; i++ {
			ctx := generator.GeneratorContext{Difficulty: 0.5}
			prob := gen.Generate(ctx)

			if strings.TrimSpace(prob.Question) == "" {
				t.Errorf("%s: empty question at draw %d", id, i)
			}
			if strings.TrimSpace(prob.Answer) == "" {
				t.Errorf("%s: empty answer at draw %d", id, i)
			}
			if strings.TrimSpace(prob.Explanation) == "" {
				t.Errorf("%s: empty explanation at draw %d", id, i)
			}
			answers[prob.Answer] = true
			stems[prob.Question] = true
		}

		if len(answers) < 5 {
			t.Errorf("%s: only %d distinct answers in 100 draws", id, len(answers))
		}
		if len(stems) < 10 {
			t.Errorf("%s: only %d distinct questions in 100 draws", id, len(stems))
		}
	}
}
