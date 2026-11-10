package arithmetic

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w8TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w8TestHighestPower(n int) int {
	p := 1
	for p<<1 <= n {
		p <<= 1
	}
	return p
}

func w8TestPopcount(n int) int {
	c := 0
	for n > 0 {
		c += n & 1
		n >>= 1
	}
	return c
}

func TestWave8ArithAnswers(t *testing.T) {
	for _, kind := range []string{"binary", "making_ten", "partial_products"} {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w8ArithGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w8TestArithAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w8TestArithAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "binary":
		n := w8TestInt(t, p, "n")
		switch kp {
		case 1:
			return n
		case 2:
			return w8TestHighestPower(n)
		default:
			return w8TestPopcount(n)
		}
	case "making_ten":
		a, b := w8TestInt(t, p, "a"), w8TestInt(t, p, "b")
		switch kp {
		case 1:
			return 10 - a
		case 2:
			return a + b
		default:
			return a + b - 10
		}
	case "partial_products":
		a, b := w8TestInt(t, p, "a"), w8TestInt(t, p, "b")
		tens, u, v, w := w8TestInt(t, p, "t"), w8TestInt(t, p, "u"), w8TestInt(t, p, "v"), w8TestInt(t, p, "w")
		if a != 10*tens+u || b != 10*v+w {
			t.Fatalf("partial-product facts do not match a and b: %v", p.Facts)
		}
		switch kp {
		case 1:
			return 10 * tens * b
		case 2:
			return u * b
		default:
			return a * b
		}
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave8ArithDiversity(t *testing.T) {
	for _, kind := range []string{"binary", "making_ten", "partial_products"} {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w8ArithGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
