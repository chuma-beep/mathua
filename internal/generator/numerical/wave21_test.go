package numerical

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
	for kp := 1; kp <= 3; kp++ {
		for i := 0; i < 8; i++ {
			p := wave21Draw(t, reg, "numa.optimization.anderson", kp)
			want := 0
			switch kp {
			case 1:
				want = wave21Fact(t, p, "m")
			case 2:
				want = 100 - wave21Fact(t, p, "beta")
			case 3:
				n, s := wave21Fact(t, p, "n"), wave21Fact(t, p, "s")
				want = (n*s + 99) / 100
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("kp%d: answer %q want %d; facts %v", kp, p.Answer, want, p.Facts)
			}
			if p.Facts["answer"] != p.Answer || p.Explanation == "" {
				t.Fatalf("kp%d: incomplete problem %+v", kp, p)
			}
		}
	}
}

func TestWave21QuestionVariety(t *testing.T) {
	reg := wave21Registry(t)
	for kp := 1; kp <= 3; kp++ {
		seen := map[string]bool{}
		for i := 0; i < 12; i++ {
			seen[wave21Draw(t, reg, "numa.optimization.anderson", kp).Question] = true
		}
		if len(seen) < 5 {
			t.Errorf("kp%d: only %d distinct questions in 12 draws", kp, len(seen))
		}
	}
}
