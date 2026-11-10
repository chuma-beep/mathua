package prealgebra

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10TestTriple(t *testing.T, s string) (int, int, int) {
	t.Helper()
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		t.Fatalf("bad triple %q", s)
	}
	a, _ := strconv.Atoi(parts[0])
	b, _ := strconv.Atoi(parts[1])
	c, _ := strconv.Atoi(parts[2])
	return a, b, c
}

func TestWave10PrealgebraAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 60; draw++ {
			p := (&w10ThreeTermRatioGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			r1, r2, r3 := w10TestInt(t, p, "r1"), w10TestInt(t, p, "r2"), w10TestInt(t, p, "r3")
			scale := w10TestInt(t, p, "scale")
			if kp == 3 {
				largest := r1
				if r2 > largest {
					largest = r2
				}
				if r3 > largest {
					largest = r3
				}
				want := strconv.Itoa(largest * scale)
				if p.Answer != want {
					t.Fatalf("three-term kp3: answer = %s, want %s\nfacts: %v", p.Answer, want, p.Facts)
				}
			} else {
				a, b, c := w10TestTriple(t, p.Answer)
				if a != r1 || b != r2 || c != r3 {
					t.Fatalf("three-term kp%d: answer = %s, want %d:%d:%d\nfacts: %v", kp, p.Answer, r1, r2, r3, p.Facts)
				}
				if w10GCD(w10GCD(r1, r2), r3) != 1 {
					t.Fatalf("three-term kp%d: ratio is not reduced: %s", kp, p.Answer)
				}
			}

			p = (&w10PartWholeGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			rp, rw := w10TestInt(t, p, "ratio_part"), w10TestInt(t, p, "ratio_whole")
			part, whole := w10TestInt(t, p, "part"), w10TestInt(t, p, "whole")
			switch kp {
			case 1:
				if p.Answer != fmt.Sprintf("%d:%d", rp, rw) {
					t.Fatalf("part-whole kp1: answer = %s, want %d:%d\nfacts: %v", p.Answer, rp, rw, p.Facts)
				}
			case 2:
				if p.Answer != strconv.Itoa(whole) {
					t.Fatalf("part-whole kp2: answer = %s, want %d\nfacts: %v", p.Answer, whole, p.Facts)
				}
			default:
				if p.Answer != strconv.Itoa(part) {
					t.Fatalf("part-whole kp3: answer = %s, want %d\nfacts: %v", p.Answer, part, p.Facts)
				}
			}

			p = (&w10OriginalAmountGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			original := w10TestInt(t, p, "original")
			change := w10TestInt(t, p, "change")
			want := original
			if kp == 3 {
				want = change
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("original kp%d: answer = %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}

			p = (&w10LiteralEquationGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			v1, v2, v3 := w10TestInt(t, p, "value1"), w10TestInt(t, p, "value2"), w10TestInt(t, p, "value3")
			want = 0
			switch p.Facts["target"] {
			case "w":
				want = (v1 - 2*v2) / 2
			case "h":
				if strings.HasPrefix(p.Facts["formula"], "V") {
					want = v1 / (v2 * v3)
				} else {
					want = 2 * v1 / (v2 + v3)
				}
			}
			if p.Answer != strconv.Itoa(want) {
				t.Fatalf("literal kp%d: answer = %s, want %d\nfacts: %v", kp, p.Answer, want, p.Facts)
			}
		}
	}
}

func TestWave10PrealgebraQuestionDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"three-term", &w10ThreeTermRatioGen{}},
		{"part-whole", &w10PartWholeGen{}},
		{"original", &w10OriginalAmountGen{}},
		{"literal", &w10LiteralEquationGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 10 {
				t.Fatalf("%s kp%d: only %d distinct questions in 200 draws", tc.name, kp, len(seen))
			}
			short := map[string]bool{}
			for i := 0; i < 12; i++ {
				p := tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				short[p.Question] = true
			}
			if len(short) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", tc.name, kp, len(short))
			}
		}
	}
}
