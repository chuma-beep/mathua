package machinelearning

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w7RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w7RemTestInts(s string) []int {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	if s == "" {
		return nil
	}
	fields := strings.Fields(s)
	out := make([]int, len(fields))
	for i, f := range fields {
		out[i], _ = strconv.Atoi(f)
	}
	return out
}

func TestWave7RemainingMLAnswerRecomputation(t *testing.T) {
	kinds := []string{"elastic_net", "polynomial", "soft_margin", "dbscan", "hierarchical", "autoencoder", "sgd", "leakage"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7MLRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7RemTestMLAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7RemTestMLAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "elastic_net":
		l1, l2 := w7RemTestInt(t, p, "l1"), w7RemTestInt(t, p, "l2")
		w := w7RemTestInt(t, p, "w")
		if kp == 1 {
			return l1*w7Abs(w) + l2*w*w
		}
		if kp == 2 {
			return l1 + l2
		}
		c := w7RemTestInt(t, p, "c")
		lambda := l1 + l2
		if w7Abs(c) <= lambda {
			return 0
		}
		if c > 0 {
			return c - lambda
		}
		return c + lambda
	case "polynomial":
		x, d, w := w7RemTestInt(t, p, "x"), w7RemTestInt(t, p, "d"), w7RemTestInt(t, p, "w")
		pow := 1
		for i := 0; i < d; i++ {
			pow *= x
		}
		if kp == 1 {
			return pow
		}
		if kp == 2 {
			return d
		}
		return w * pow
	case "soft_margin":
		if kp == 1 {
			return w7RemTestInt(t, p, "sum")
		}
		if kp == 2 {
			return 20 / w7RemTestInt(t, p, "norm")
		}
		return 2
	case "dbscan":
		if kp == 1 {
			if w7RemTestInt(t, p, "neighbors") >= w7RemTestInt(t, p, "minPts") {
				return 1
			}
			return 0
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "clusters")
		}
		return w7RemTestInt(t, p, "noise")
	case "hierarchical":
		if kp == 1 {
			return w7RemTestInt(t, p, "d1")
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "d3")
		}
		return 2
	case "autoencoder":
		if kp == 1 {
			return w7RemTestInt(t, p, "b")
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "n") * w7RemTestInt(t, p, "h")
		}
		sum := 0
		for _, e := range w7RemTestInts(p.Facts["errs"]) {
			sum += e * e
		}
		return sum
	case "sgd":
		w, g, k := w7RemTestInt(t, p, "w"), w7RemTestInt(t, p, "g"), w7RemTestInt(t, p, "k")
		if kp == 1 {
			return w - g/k
		}
		if kp == 2 {
			return (w7RemTestInt(t, p, "g1") + w7RemTestInt(t, p, "g2")) / 2
		}
		return k
	case "leakage":
		if kp == 1 {
			return w7RemTestInt(t, p, "leak")
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "bad")
		}
		return 1 - w7RemTestInt(t, p, "bad")
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7RemainingMLDiversity(t *testing.T) {
	kinds := []string{"elastic_net", "polynomial", "soft_margin", "dbscan", "hierarchical", "autoencoder", "sgd", "leakage"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7MLRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
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
