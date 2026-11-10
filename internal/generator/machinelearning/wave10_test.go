package machinelearning

import (
	"strconv"
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

func TestWave10AdaGradAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 60; draw++ {
			p := (&w10AdaGradGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			g1, g2 := w10TestInt(t, p, "g1"), w10TestInt(t, p, "g2")
			eta := w10TestInt(t, p, "eta")
			G1, G2 := w10TestInt(t, p, "G1"), w10TestInt(t, p, "G2")
			u1, u2 := w10TestInt(t, p, "u1"), w10TestInt(t, p, "u2")
			w0, w1 := w10TestInt(t, p, "w0"), w10TestInt(t, p, "w1")
			if G1 != g1*g1 {
				t.Fatalf("adagrad kp%d: G1 = %d, want %d", kp, G1, g1*g1)
			}
			if G2 != G1+g2*g2 {
				t.Fatalf("adagrad kp%d: G2 = %d, want %d", kp, G2, G1+g2*g2)
			}
			if u1 != eta {
				t.Fatalf("adagrad kp%d: u1 = %d, want %d", kp, u1, eta)
			}
			root2 := 0
			for root2*root2 < G2 {
				root2++
			}
			if u2 != eta*g2/root2 {
				t.Fatalf("adagrad kp%d: u2 = %d, want %d", kp, u2, eta*g2/root2)
			}
			if w1 != w0-u1 {
				t.Fatalf("adagrad kp%d: w1 = %d, want %d", kp, w1, w0-u1)
			}
			switch kp {
			case 1:
				if p.Answer != strconv.Itoa(w1) {
					t.Fatalf("adagrad kp1: answer = %s, want %d", p.Answer, w1)
				}
			case 2:
				if p.Answer != strconv.Itoa(G2) {
					t.Fatalf("adagrad kp2: answer = %s, want %d", p.Answer, G2)
				}
			default:
				if p.Answer != strconv.Itoa(u1-u2) {
					t.Fatalf("adagrad kp3: answer = %s, want %d", p.Answer, u1-u2)
				}
			}
		}
	}
}

func TestWave10AdaGradQuestionDiversity(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		seen := map[string]bool{}
		for i := 0; i < 200; i++ {
			seen[(&w10AdaGradGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
		}
		if len(seen) < 10 {
			t.Fatalf("adagrad kp%d: only %d distinct questions in 200 draws", kp, len(seen))
		}
		short := map[string]bool{}
		for i := 0; i < 12; i++ {
			short[(&w10AdaGradGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
		}
		if len(short) < 5 {
			t.Fatalf("adagrad kp%d: only %d distinct questions in 12 draws", kp, len(short))
		}
	}
}
