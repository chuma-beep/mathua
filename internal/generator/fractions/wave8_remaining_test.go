package fractions

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w8RFInt(t *testing.T, p generator.Problem, k string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[k])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer: %v", p.Question, k, p.Facts[k], p.Facts)
	}
	return v
}

func w8RFGCDTest(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func TestWave8RemainingFractions(t *testing.T) {
	kinds := []string{"improper_number_line", "order_three", "add_number_line", "sub_number_line", "mult_area_model", "missing_factor", "mixed_decimal", "percent_convert", "remaining_two_step"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8FracRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8RemFracAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8RemFracAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "improper_number_line":
		if kp == 1 || kp == 3 {
			return w8RFInt(t, p, "whole")
		}
		return w8RFInt(t, p, "remainder")
	case "order_three":
		a1, b1 := w8RFInt(t, p, "a1"), w8RFInt(t, p, "b1")
		a2, b2 := w8RFInt(t, p, "a2"), w8RFInt(t, p, "b2")
		a3, b3 := w8RFInt(t, p, "a3"), w8RFInt(t, p, "b3")
		v1, v2, v3 := a1*b2*b3, a2*b1*b3, a3*b1*b2
		if kp == 1 {
			small := 1
			if v2 < v1 && v2 < v3 {
				small = 2
			} else if v3 < v1 && v3 < v2 {
				small = 3
			}
			return small
		}
		if kp == 2 {
			large := 1
			if v2 > v1 && v2 > v3 {
				large = 2
			} else if v3 > v1 && v3 > v2 {
				large = 3
			}
			return large
		}
		return w8RFInt(t, p, "lcm")
	case "add_number_line":
		a, b, c := w8RFInt(t, p, "a"), w8RFInt(t, p, "b"), w8RFInt(t, p, "c")
		if kp == 1 {
			return a + c
		}
		if kp == 2 {
			return b
		}
		g := w8RFGCDTest(a+c, b)
		return (a + c) / g
	case "sub_number_line":
		a, b, c := w8RFInt(t, p, "a"), w8RFInt(t, p, "b"), w8RFInt(t, p, "c")
		if kp == 1 {
			return a - c
		}
		if kp == 2 {
			return b
		}
		g := w8RFGCDTest(a-c, b)
		return (a - c) / g
	case "mult_area_model":
		a, b, c, d := w8RFInt(t, p, "a"), w8RFInt(t, p, "b"), w8RFInt(t, p, "c"), w8RFInt(t, p, "d")
		if kp == 1 {
			return a * c
		}
		if kp == 2 {
			return b * d
		}
		g := w8RFGCDTest(a*c, b*d)
		return (a * c) / g
	case "missing_factor":
		a, b, c, d := w8RFInt(t, p, "a"), w8RFInt(t, p, "b"), w8RFInt(t, p, "c"), w8RFInt(t, p, "d")
		num, den := a*d, b*c
		g := w8RFGCDTest(num, den)
		if kp == 1 {
			return num / g
		}
		if kp == 2 {
			return den / g
		}
		return num
	case "mixed_decimal":
		w, n := w8RFInt(t, p, "whole"), w8RFInt(t, p, "numerator")
		if kp == 1 {
			return w*100 + n*10
		}
		if kp == 2 {
			return w*10 + n
		}
		return 1
	case "percent_convert":
		a, b := w8RFInt(t, p, "a"), w8RFInt(t, p, "b")
		pct := a * 100 / b
		if kp == 1 {
			return pct
		}
		if kp == 2 {
			return a * 100 / b
		}
		if pct > 50 {
			return 1
		}
		return 0
	case "remaining_two_step":
		n := w8RFInt(t, p, "n")
		n1, d1 := w8RFInt(t, p, "num1"), w8RFInt(t, p, "den1")
		n2, d2 := w8RFInt(t, p, "num2"), w8RFInt(t, p, "den2")
		afterFirst := n * (d1 - n1) / d1
		afterSecond := afterFirst * (d2 - n2) / d2
		if kp == 1 {
			return afterFirst
		}
		if kp == 2 {
			return afterSecond
		}
		g := w8RFGCDTest(afterSecond, n)
		return afterSecond / g
	}
	t.Fatalf("unknown kind %s", kind)
	return 0
}

func TestWave8RemainingFractionsDiversity(t *testing.T) {
	kinds := []string{"improper_number_line", "order_three", "add_number_line", "sub_number_line", "mult_area_model", "missing_factor", "mixed_decimal", "percent_convert", "remaining_two_step"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8FracRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
