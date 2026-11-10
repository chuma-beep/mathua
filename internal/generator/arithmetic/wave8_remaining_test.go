package arithmetic

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w8RTInt(t *testing.T, p generator.Problem, k string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[k])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer: %v", p.Question, k, p.Facts[k], p.Facts)
	}
	return v
}

func w8RFact(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func w8RGCDTest(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func TestWave8RemainingArithmetic(t *testing.T) {
	kinds := []string{"base_five", "roman", "tally", "doubling_halving", "compensation", "near_doubles", "partial_quotients", "area_model", "lattice", "square_numbers", "triangular_numbers", "factorial", "gcf_lcm", "elapsed", "change"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8ArithRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8RemArithAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8RemArithAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "base_five":
		n := w8RTInt(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return w8RTInt(t, p, "power")
		}
		s := 0
		for _, c := range p.Facts["numeral"] {
			s += int(c - '0')
		}
		return s
	case "roman":
		if kp == 1 {
			return w8RTInt(t, p, "n")
		}
		if kp == 2 {
			return strings.Count(p.Facts["roman"], "X")
		}
		for _, pair := range []string{"IV", "IX", "XL", "XC", "CD", "CM"} {
			if strings.Contains(p.Facts["roman"], pair) {
				return 1
			}
		}
		return 0
	case "tally":
		if kp == 1 {
			return w8RTInt(t, p, "n")
		}
		if kp == 2 {
			return w8RTInt(t, p, "groups")
		}
		return w8RTInt(t, p, "remainder")
	case "doubling_halving":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return 2 * a
		}
		if kp == 2 {
			return a
		}
		return a * b
	case "compensation":
		a, d := w8RTInt(t, p, "a"), w8RTInt(t, p, "delta")
		if kp == 1 {
			return 10 - d
		}
		if kp == 2 {
			return a + d
		}
		return a + 10
	case "near_doubles":
		a := w8RTInt(t, p, "a")
		if kp == 1 {
			return 2 * a
		}
		if kp == 2 {
			return 2*a + 1
		}
		return 1
	case "partial_quotients":
		d, q, r := w8RTInt(t, p, "divisor"), w8RTInt(t, p, "quotient"), w8RTInt(t, p, "remainder")
		if kp == 1 {
			return q
		}
		if kp == 2 {
			return r
		}
		return d * q
	case "area_model":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return a * b
		}
		return b
	case "lattice":
		if kp == 1 {
			return w8RTInt(t, p, "cross")
		}
		if kp == 2 {
			return w8RTInt(t, p, "units")
		}
		return w8RTInt(t, p, "a") * w8RTInt(t, p, "b")
	case "square_numbers":
		n := w8RTInt(t, p, "n")
		if kp == 1 {
			return n * n
		}
		if kp == 2 {
			return n
		}
		return 2*n + 1
	case "triangular_numbers":
		n := w8RTInt(t, p, "n")
		if kp == 2 {
			return n
		}
		return n * (n + 1) / 2
	case "factorial":
		n := w8RTInt(t, p, "n")
		if kp == 1 {
			return w8RFact(n)
		}
		if kp == 2 {
			return n
		}
		return n * (n - 1)
	case "gcf_lcm":
		a, b := w8RTInt(t, p, "a"), w8RTInt(t, p, "b")
		if kp == 1 {
			return w8RGCDTest(a, b)
		}
		if kp == 2 {
			return w8RTInt(t, p, "lcm")
		}
		return a * b
	case "elapsed":
		if kp == 1 {
			return w8RTInt(t, p, "duration")
		}
		if kp == 2 {
			return w8RTInt(t, p, "end")
		}
		return w8RTInt(t, p, "start")
	case "change":
		if kp == 1 {
			return w8RTInt(t, p, "change")
		}
		if kp == 2 {
			return w8RTInt(t, p, "change") / 25
		}
		return w8RTInt(t, p, "paid")
	}
	t.Fatalf("unknown kind %s", kind)
	return 0
}

func TestWave8RemainingArithmeticDiversity(t *testing.T) {
	kinds := []string{"base_five", "roman", "tally", "doubling_halving", "compensation", "near_doubles", "partial_quotients", "area_model", "lattice", "square_numbers", "triangular_numbers", "factorial", "gcf_lcm", "elapsed", "change"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8ArithRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d only %d distinct questions", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("below target: %s kp%d = %d distinct in 60", kind, kp, len(seen))
			}
		}
	}
}
