package generator_test

// Pool discrimination ratchet: every concept's easy pool must discriminate
// (more than one distinct answer across samples) and must not be yes/no-only
// (recognition without production). Known violations live in
// scripts/audit_baseline.json and shrink to empty as content waves land;
// any NEW violation fails the build.
import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/abstract"
	"github.com/chuma-beep/mathua/internal/generator/algebra"
	"github.com/chuma-beep/mathua/internal/generator/arithmetic"
	"github.com/chuma-beep/mathua/internal/generator/calculus"
	"github.com/chuma-beep/mathua/internal/generator/complex"
	"github.com/chuma-beep/mathua/internal/generator/discrete"
	"github.com/chuma-beep/mathua/internal/generator/fractions"
	"github.com/chuma-beep/mathua/internal/generator/geometry"
	"github.com/chuma-beep/mathua/internal/generator/linalg"
	"github.com/chuma-beep/mathua/internal/generator/machinelearning"
	"github.com/chuma-beep/mathua/internal/generator/numtheory"
	"github.com/chuma-beep/mathua/internal/generator/odes"
	"github.com/chuma-beep/mathua/internal/generator/prealgebra"
	"github.com/chuma-beep/mathua/internal/generator/precalculus"
	"github.com/chuma-beep/mathua/internal/generator/statistics"
	"github.com/chuma-beep/mathua/internal/generator/topology"
	"github.com/chuma-beep/mathua/internal/generator/trigonometry"
)

func registerAllDomains(reg *generator.Registry) {
	abstract.Register(reg)
	algebra.Register(reg)
	arithmetic.Register(reg)
	calculus.Register(reg)
	complex.Register(reg)
	discrete.Register(reg)
	fractions.Register(reg)
	geometry.Register(reg)
	linalg.Register(reg)
	machinelearning.Register(reg)
	numtheory.Register(reg)
	odes.Register(reg)
	prealgebra.Register(reg)
	precalculus.Register(reg)
	statistics.Register(reg)
	topology.Register(reg)
	trigonometry.Register(reg)
}

func normAnswer(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func TestEasyPoolsDiscriminate(t *testing.T) {
	rand.Seed(20260923) // deterministic sampling; pools are small
	const samples = 150

	reg := generator.NewRegistry()
	registerAllDomains(reg)

	var singleAnswer, yesNoOnly []string
	for _, id := range reg.Concepts() {
		gen, err := reg.Get(id)
		if err != nil {
			t.Fatalf("registered concept %q not retrievable: %v", id, err)
		}
		answers := map[string]bool{}
		allYesNo := true
		for i := 0; i < samples; i++ {
			p := gen.Generate(generator.GeneratorContext{Difficulty: 0.12})
			answers[normAnswer(p.Answer)] = true
			if !strings.Contains(p.Question, "(yes/no)") {
				allYesNo = false
			}
		}
		if len(answers) < 2 {
			singleAnswer = append(singleAnswer, id)
		}
		if allYesNo {
			yesNoOnly = append(yesNoOnly, id)
		}
	}
	sort.Strings(singleAnswer)
	sort.Strings(yesNoOnly)

	raw, err := os.ReadFile(filepath.Join("..", "..", "scripts", "audit_baseline.json"))
	if err != nil {
		t.Fatalf("read audit baseline: %v", err)
	}
	var baseline struct {
		ThinPoolsSingleAnswer []string `json:"thin_pools_single_answer"`
		ThinPoolsYesNoOnly    []string `json:"thin_pools_yesno_only"`
	}
	if err := json.Unmarshal(raw, &baseline); err != nil {
		t.Fatalf("parse audit baseline: %v", err)
	}
	checkRatchet := func(name string, current, allowed []string) {
		allow := map[string]bool{}
		for _, id := range allowed {
			allow[id] = true
		}
		var fresh []string
		for _, id := range current {
			if !allow[id] {
				fresh = append(fresh, id)
			}
		}
		t.Logf("%s: %d current (%d baselined), new: %v", name, len(current), len(allowed), fresh)
		if len(fresh) > 0 {
			t.Errorf("%s: %d NEW violations not in baseline (fix the pools, do not extend the baseline): %v",
				name, len(fresh), fresh)
		}
	}
	checkRatchet("single-answer easy pools", singleAnswer, baseline.ThinPoolsSingleAnswer)
	checkRatchet("yes/no-only easy pools", yesNoOnly, baseline.ThinPoolsYesNoOnly)
}
