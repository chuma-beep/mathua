package probability

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6RemTestInts(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' })
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		v, _ := strconv.Atoi(part)
		out = append(out, v)
	}
	return out
}

func w6RemTestMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w6RemTestMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6RemTestLog2(n int) int {
	c := 0
	for n > 1 {
		n /= 2
		c++
	}
	return c
}

func TestWave6RemainingProbabilityAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"renewal", "queue", "ito", "sde", "gaussian", "hmm", "kalman",
		"hoeffding", "bernstein", "mcdiarmid", "lindeberg", "berry", "local_clt",
		"conjugate", "predictive", "bayes_risk", "minimax", "entropy", "kl",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6ProbRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w6RemTestProbAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6RemTestProbAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "renewal":
		mu, tt := w6RemTestInt(t, p, "mu"), w6RemTestInt(t, p, "t")
		if kp == 1 {
			return tt / mu
		}
		if kp == 2 {
			return mu
		}
		return 2 * tt / mu
	case "queue":
		lambda, mu := w6RemTestInt(t, p, "lambda"), w6RemTestInt(t, p, "mu")
		rho := 100 * lambda / mu
		if kp == 1 {
			return rho
		}
		if kp == 2 {
			return 100 - rho
		}
		return lambda / (mu - lambda)
	case "ito":
		tt, c := w6RemTestInt(t, p, "t"), w6RemTestInt(t, p, "c")
		if kp == 1 {
			return 0
		}
		if kp == 2 {
			return c * c * tt
		}
		return 2 * c
	case "sde":
		a, b, h, x0 := w6RemTestInt(t, p, "a"), w6RemTestInt(t, p, "b"), w6RemTestInt(t, p, "h"), w6RemTestInt(t, p, "x0")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return b
		}
		return x0 + a*h
	case "gaussian":
		s, tt, n := w6RemTestInt(t, p, "s"), w6RemTestInt(t, p, "t"), w6RemTestInt(t, p, "n")
		_ = tt
		if kp == 2 {
			return n
		}
		return s
	case "hmm":
		if kp == 1 {
			k := w6RemTestInt(t, p, "states")
			return k * (k - 1)
		}
		if kp == 2 {
			return w6RemTestInt(t, p, "emit00")
		}
		return w6RemTestInt(t, p, "emit01")
	case "kalman":
		P, Q, R := w6RemTestInt(t, p, "P"), w6RemTestInt(t, p, "Q"), w6RemTestInt(t, p, "R")
		prior, z := w6RemTestInt(t, p, "prior"), w6RemTestInt(t, p, "z")
		if kp == 1 {
			return P + Q
		}
		if kp == 2 {
			return P
		}
		return R*prior + P*z
	case "hoeffding":
		n := w6RemTestInt(t, p, "n")
		if kp == 1 {
			return n / 2
		}
		if kp == 2 {
			return 2 * n
		}
		return 2 * w6RemTestInt(t, p, "target")
	case "bernstein":
		v := w6RemTestInt(t, p, "v_num")
		if kp == 1 {
			return 6*v + 4
		}
		if kp == 2 {
			return 3*v + 2
		}
		return 6*v + 2
	case "mcdiarmid":
		n := w6RemTestInt(t, p, "n")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return n / 2
		}
		return 4
	case "lindeberg":
		total, mx := w6RemTestInt(t, p, "total"), w6RemTestInt(t, p, "max")
		if kp == 1 {
			if mx*10 < total {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return total / mx
		}
		if mx*10 >= total {
			return 1
		}
		return 0
	case "berry":
		sigma, n := w6RemTestInt(t, p, "sigma"), w6RemTestInt(t, p, "n")
		if kp == 1 {
			return w6RemTestInt(t, p, "k")
		}
		if kp == 2 {
			return sigma * sigma * sigma
		}
		return n * sigma * sigma * sigma
	case "local_clt":
		v, n := w6RemTestInt(t, p, "v"), w6RemTestInt(t, p, "n")
		if kp == 1 {
			return 2
		}
		if kp == 2 {
			return n * v
		}
		return 1
	case "conjugate":
		a, b, s, fail := w6RemTestInt(t, p, "a"), w6RemTestInt(t, p, "b"), w6RemTestInt(t, p, "successes"), w6RemTestInt(t, p, "failures")
		if kp == 1 {
			return a + s
		}
		if kp == 2 {
			return b + fail
		}
		return a + b + s + fail
	case "predictive":
		a, b, s, fail := w6RemTestInt(t, p, "a"), w6RemTestInt(t, p, "b"), w6RemTestInt(t, p, "successes"), w6RemTestInt(t, p, "failures")
		total := a + b + s + fail
		if kp == 1 {
			return a + s
		}
		if kp == 2 {
			return total
		}
		return b + fail
	case "bayes_risk":
		prior, l1, l2 := w6RemTestInt(t, p, "prior"), w6RemTestInt(t, p, "loss1"), w6RemTestInt(t, p, "loss2")
		r1, r2 := prior*l1, (10-prior)*l2
		if kp == 1 {
			return r1
		}
		if kp == 2 {
			return r2
		}
		if r2 < r1 {
			return 2
		}
		return 1
	case "minimax":
		a, b, c, d := w6RemTestInt(t, p, "a"), w6RemTestInt(t, p, "b"), w6RemTestInt(t, p, "c"), w6RemTestInt(t, p, "d")
		r1, r2 := w6RemTestMax(a, b), w6RemTestMax(c, d)
		if kp == 1 {
			return r1
		}
		if kp == 2 {
			return r2
		}
		return w6RemTestMin(r1, r2)
	case "entropy":
		weights := w6RemTestInts(p.Facts["weights"])
		total := 0
		for _, w := range weights {
			total += w
		}
		if kp == 1 {
			h := 0
			for _, w := range weights {
				h += w * (w6RemTestLog2(total) - w6RemTestLog2(w))
			}
			return 2 * h / total
		}
		if kp == 2 {
			n := w6RemTestInt(t, p, "n")
			return 2 * w6RemTestLog2(n*n)
		}
		return w6RemTestLog2(w6RemTestInt(t, p, "outcomes"))
	case "kl":
		if kp == 1 {
			return w6RemTestLog2(w6RemTestInt(t, p, "n"))
		}
		if kp == 2 {
			return 0
		}
		return w6RemTestInt(t, p, "a") + w6RemTestInt(t, p, "b")
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6RemainingProbabilityDiversity(t *testing.T) {
	kinds := []string{
		"renewal", "queue", "ito", "sde", "gaussian", "hmm", "kalman",
		"hoeffding", "bernstein", "mcdiarmid", "lindeberg", "berry", "local_clt",
		"conjugate", "predictive", "bayes_risk", "minimax", "entropy", "kl",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6ProbRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Errorf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
