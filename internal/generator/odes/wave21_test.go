package odes

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func wave21Registry(t *testing.T) *generator.Registry {
	t.Helper()
	reg := generator.NewRegistry()
	Register(reg)
	return reg
}

func wave21Draw(t *testing.T, reg *generator.Registry, id string, kp int) generator.Problem {
	t.Helper()
	g, err := reg.Get(id)
	if err != nil {
		t.Fatalf("get %s: %v", id, err)
	}
	return g.Generate(generator.GeneratorContext{KPOrdinal: kp})
}

func wave21Fact(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer", p.Question, key, p.Facts[key])
	}
	return v
}

func TestWave21AnswersRecompute(t *testing.T) {
	reg := wave21Registry(t)
	for _, id := range []string{
		"ode.systems.coupled", "ode.boundary.eigenfunction",
		"ode.integral.resolvent", "ode.nonlinear.hopf",
	} {
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 8; i++ {
				p := wave21Draw(t, reg, id, kp)
				want := 0
				switch id {
				case "ode.systems.coupled":
					switch kp {
					case 1:
						want = 2 * wave21Fact(t, p, "n")
					case 2:
						n := wave21Fact(t, p, "n")
						want = n * n
					case 3:
						want = 2 * wave21Fact(t, p, "k")
					}
				case "ode.boundary.eigenfunction":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "n") - 1
					case 2:
						want = wave21Fact(t, p, "terms")
					case 3:
						want = 2 * wave21Fact(t, p, "val")
					}
				case "ode.integral.resolvent":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "n")
					case 2:
						want = wave21Fact(t, p, "terms")
					case 3:
						want = wave21Fact(t, p, "n") - 1
					}
				case "ode.nonlinear.hopf":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "k")
					case 2:
						want = wave21Fact(t, p, "sq")
					case 3:
						want = wave21Fact(t, p, "d")
					}
				}
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %q want %d; facts %v", id, kp, p.Answer, want, p.Facts)
				}
				if p.Facts["answer"] != p.Answer || p.Explanation == "" {
					t.Fatalf("%s kp%d: incomplete problem %+v", id, kp, p)
				}
			}
		}
	}
}

func TestWave21QuestionVariety(t *testing.T) {
	reg := wave21Registry(t)
	for _, id := range []string{
		"ode.systems.coupled", "ode.boundary.eigenfunction",
		"ode.integral.resolvent", "ode.nonlinear.hopf",
	} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 12; i++ {
				seen[wave21Draw(t, reg, id, kp).Question] = true
			}
			if len(seen) < 5 {
				t.Errorf("%s kp%d: only %d distinct questions in 12 draws", id, kp, len(seen))
			}
		}
	}
}
