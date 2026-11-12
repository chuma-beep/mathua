package calculus

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

func wave21IDs() []string {
	return []string{
		"calc.deriv.reciprocal", "calc.deriv.leibniz_rule",
		"calc.deriv.linear_error", "calc.limit.stolz",
		"calc.seq.recursive", "calc.integral.periodic",
	}
}

func TestWave21AnswersRecompute(t *testing.T) {
	reg := wave21Registry(t)
	for _, id := range wave21IDs() {
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < 8; i++ {
				p := wave21Draw(t, reg, id, kp)
				want := 0
				switch id {
				case "calc.deriv.reciprocal":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "n")
					case 2:
						a := wave21Fact(t, p, "a")
						want = a * a
					case 3:
						want = wave21Fact(t, p, "q")
					}
				case "calc.deriv.leibniz_rule":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "n") + 1
					case 2:
						n, k := wave21Fact(t, p, "n"), wave21Fact(t, p, "k")
						want = w21Binom(n, k)
					case 3:
						want = wave21Fact(t, p, "d") + 1
					}
				case "calc.deriv.linear_error":
					switch kp {
					case 1:
						m, d := wave21Fact(t, p, "m"), wave21Fact(t, p, "d")
						want = m * d * d
					case 2:
						bigM, d := wave21Fact(t, p, "m"), wave21Fact(t, p, "d")
						want = bigM * d * d / 2
					case 3:
						k := wave21Fact(t, p, "k")
						want = k * k
					}
				case "calc.limit.stolz":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "p") + 1
					case 2:
						want = wave21Fact(t, p, "k") + 1
					case 3:
						want = wave21Fact(t, p, "q")
					}
				case "calc.seq.recursive":
					switch kp {
					case 1:
						steps := 0
						for v := wave21Fact(t, p, "k"); v >= 1; v /= 2 {
							steps++
						}
						want = steps
					case 2:
						want = wave21Fact(t, p, "a")
					case 3:
						want = 2 * wave21Fact(t, p, "c")
					}
				case "calc.integral.periodic":
					switch kp {
					case 1:
						want = wave21Fact(t, p, "i") * wave21Fact(t, p, "k")
					case 2:
						want = wave21Fact(t, p, "i") - wave21Fact(t, p, "j")
					case 3:
						want = wave21Fact(t, p, "i") / wave21Fact(t, p, "t")
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
	for _, id := range wave21IDs() {
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
