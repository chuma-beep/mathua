package machinelearning

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w10bF(t *testing.T, p generator.Problem, key string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(p.Facts[key], 64)
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not numeric\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w10bNear(got, want float64) bool { return math.Abs(got-want) < 1e-3 }

func w10bVals(t *testing.T, s string) []int {
	t.Helper()
	var out []int
	for _, part := range strings.Split(s, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("bad Fiedler vector %q", s)
		}
		out = append(out, v)
	}
	return out
}

func TestWave10RestMLAnswerRecomputation(t *testing.T) {
	for kp := 1; kp <= 3; kp++ {
		for draw := 0; draw < 40; draw++ {
			p := (&w10gRMSPropGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			g1, g2, eta := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			v1 := 0.5 * float64(g1*g1)
			v2 := 0.5*float64(g2*g2) + 0.5*v1
			u1 := float64(eta*g1) / math.Sqrt(v1)
			u2 := float64(eta*g2) / math.Sqrt(v2)
			wantF := u1
			if kp == 2 {
				wantF = v2
			} else if kp == 3 {
				wantF = u1 - u2
			}
			if !w10bNear(w10bF(t, p, "answer"), wantF) || !w10bNear(w10bF(t, p, "v1"), v1) || !w10bNear(w10bF(t, p, "v2"), v2) {
				t.Fatalf("rmsprop kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}

			p = (&w10gAdamWGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			w0, lr, lambda, grad := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c"), w10TestInt(t, p, "grad")
			decay := w0 * lr * lambda / 100
			adaptive := lr * grad
			newW := w0 - adaptive - decay
			wantI := decay
			if kp == 2 {
				wantI = adaptive
			} else if kp == 3 {
				wantI = newW
			}
			if p.Answer != strconv.Itoa(wantI) || w10TestInt(t, p, "new") != newW {
				t.Fatalf("adamw kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, wantI, p.Facts)
			}

			p = (&w10gStackingGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			p1, p2, meta := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			w1, w2 := w10TestInt(t, p, "w1"), w10TestInt(t, p, "w2")
			bias, target := w10TestInt(t, p, "bias"), w10TestInt(t, p, "target")
			if meta != w1*p1+w2*p2+bias {
				t.Fatalf("stacking kp%d: meta fact does not match weights\nfacts: %v", kp, p.Facts)
			}
			wantI = meta
			if kp == 2 {
				wantI = meta - p1
			} else if kp == 3 {
				e1 := int(math.Abs(float64(p1 - target)))
				e2 := int(math.Abs(float64(p2 - target)))
				wantI = 1
				if e2 < e1 {
					wantI = 2
				}
			}
			if p.Answer != strconv.Itoa(wantI) {
				t.Fatalf("stacking kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, wantI, p.Facts)
			}

			p = (&w10gVotingGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			ones, zeros, margin := w10TestInt(t, p, "a"), w10TestInt(t, p, "b"), w10TestInt(t, p, "c")
			wantI = ones
			if kp == 2 {
				wantI = 0
				if ones > zeros {
					wantI = 1
				}
			} else if kp == 3 {
				wantI = ones - zeros
				if wantI < 0 {
					wantI = -wantI
				}
			}
			if p.Answer != strconv.Itoa(wantI) || margin != w10bAbsInt(ones-zeros) {
				t.Fatalf("voting kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, wantI, p.Facts)
			}

			p = (&w10gSpectralGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			vals := w10bVals(t, p.Facts["vals"])
			node := w10TestInt(t, p, "node")
			neg, cross := 0, 0
			for i, v := range vals {
				if v < 0 {
					neg++
				}
				if i > 0 && (vals[i] < 0) != (vals[i-1] < 0) {
					cross++
				}
			}
			wantI = neg
			if kp == 2 {
				wantI = 1
				if vals[node] > 0 {
					wantI = 2
				}
			} else if kp == 3 {
				wantI = cross
			}
			if p.Answer != strconv.Itoa(wantI) || (kp == 2 && w10TestInt(t, p, "c") != wantI) {
				t.Fatalf("spectral kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, wantI, p.Facts)
			}

			p = (&w10gHingeGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			y, f := w10TestInt(t, p, "a"), w10TestInt(t, p, "b")
			margin = y * f
			loss := 1 - margin
			if loss < 0 {
				loss = 0
			}
			sub := 0
			if margin < 1 {
				sub = -y
			}
			wantI = loss
			if kp == 2 {
				wantI = margin
			} else if kp == 3 {
				wantI = sub
			}
			if p.Answer != strconv.Itoa(wantI) {
				t.Fatalf("hinge kp%d: answer %s, want %d\nfacts: %v", kp, p.Answer, wantI, p.Facts)
			}

			p = (&w10gHuberGen{}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			delta, r := w10TestInt(t, p, "a"), w10TestInt(t, p, "b")
			abs := r
			if abs < 0 {
				abs = -abs
			}
			lossF := float64(r*r) / 2
			if abs > delta {
				lossF = float64(delta) * (float64(abs) - float64(delta)/2)
			}
			mse := float64(r*r) / 2
			wantF = lossF
			if kp == 3 {
				wantF = mse - lossF
			}
			if !w10bNear(w10bF(t, p, "answer"), wantF) || !w10bNear(w10bF(t, p, "c"), lossF) {
				t.Fatalf("huber kp%d: answer %s, want %v\nfacts: %v", kp, p.Answer, wantF, p.Facts)
			}
		}
	}
}

func w10bAbsInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func TestWave10RestMLDiversity(t *testing.T) {
	gens := []struct {
		name string
		gen  generator.Generator
	}{
		{"rmsprop", &w10gRMSPropGen{}}, {"adamw", &w10gAdamWGen{}}, {"stacking", &w10gStackingGen{}},
		{"voting", &w10gVotingGen{}}, {"spectral", &w10gSpectralGen{}}, {"hinge", &w10gHingeGen{}},
		{"huber", &w10gHuberGen{}},
	}
	for _, tc := range gens {
		for kp := 1; kp <= 3; kp++ {
			rand.Seed(int64(20261110 + kp))
			seen := map[string]bool{}
			for i := 0; i < 200; i++ {
				seen[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(seen) < 10 {
				t.Fatalf("%s kp%d: only %d distinct questions in 200 draws", tc.name, kp, len(seen))
			}
			short := map[string]bool{}
			for i := 0; i < 12; i++ {
				short[tc.gen.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp}).Question] = true
			}
			if len(short) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 12 draws", tc.name, kp, len(short))
			}
		}
	}
}
