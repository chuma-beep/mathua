package realanalysis

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6bTestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w6bTestMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w6bTestMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6bTestPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func w6bTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("missing fact %q\nquestion: %s\nfacts: %v", key, p.Question, p.Facts)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("fact %q is not an integer: %q", key, raw)
	}
	return v
}

func w6bTestInts(t *testing.T, p generator.Problem, key string) []int {
	t.Helper()
	raw, ok := p.Facts[key]
	if !ok {
		t.Fatalf("missing fact %q\nquestion: %s\nfacts: %v", key, p.Question, p.Facts)
	}
	var out []int
	for _, part := range strings.Split(strings.TrimSpace(raw), ",") {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			t.Fatalf("fact %q holds a non-integer entry: %q", key, raw)
		}
		out = append(out, v)
	}
	return out
}

var w6bKinds = []string{
	"cantor_function", "outer_measure", "caratheodory", "radon_nikodym", "signed",
	"convergence_in_measure", "egorov", "lusin", "cauchy_schwarz",
	"orthogonal_projection", "uniform_boundedness", "open_mapping",
	"closed_graph", "hahn_banach", "weak_convergence",
}

func TestWave6BRealAnswers(t *testing.T) {
	for _, kind := range w6bKinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6bGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				if !strings.HasSuffix(p.Question, "(enter a number)") {
					t.Fatalf("%s kp%d: question does not end with the numeric prompt: %s", kind, kp, p.Question)
				}
				if p.KPOrdinal != kp {
					t.Fatalf("%s kp%d: problem carries KPOrdinal %d", kind, kp, p.KPOrdinal)
				}
				if p.Facts["answer"] != p.Answer {
					t.Fatalf("%s kp%d: facts answer %q disagrees with answer %q", kind, kp, p.Facts["answer"], p.Answer)
				}
				want := w6bTestReal(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6bTestReal(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "cantor_function":
		a := 0
		for _, d := range w6bTestInts(t, p, "ternary_digits") {
			a = a*2 + d/2
		}
		return a
	case "outer_measure":
		if kp == 3 {
			return w6bTestInt(t, p, "d") - w6bTestInt(t, p, "a")
		}
		return w6bTestInt(t, p, "b") - w6bTestInt(t, p, "a")
	case "caratheodory":
		if kp == 1 {
			return w6bTestInt(t, p, "b") - w6bTestInt(t, p, "c")
		}
		if kp == 2 {
			a, b := w6bTestInt(t, p, "a"), w6bTestInt(t, p, "b")
			c, d := w6bTestInt(t, p, "c"), w6bTestInt(t, p, "d")
			return (b - a) + (d - c)
		}
		l1, x1 := w6bTestInt(t, p, "l1"), w6bTestInt(t, p, "x1")
		l2, x2 := w6bTestInt(t, p, "l2"), w6bTestInt(t, p, "x2")
		return (l1 - x1) + (l2 - x2)
	case "radon_nikodym":
		if kp == 1 {
			return w6bTestInt(t, p, "density") * w6bTestInt(t, p, "length")
		}
		if kp == 2 {
			return w6bTestInt(t, p, "density1")*w6bTestInt(t, p, "length1") +
				w6bTestInt(t, p, "density2")*w6bTestInt(t, p, "length2")
		}
		return w6bTestInt(t, p, "density1")*w6bTestInt(t, p, "length1") +
			w6bTestInt(t, p, "density2")*w6bTestInt(t, p, "length2") +
			w6bTestInt(t, p, "density3")*w6bTestInt(t, p, "length3")
	case "signed":
		sum := 0
		for _, m := range w6bTestInts(t, p, "masses") {
			sum += w6bTestAbs(m)
		}
		return sum
	case "convergence_in_measure":
		c := w6bTestInt(t, p, "c")
		n := w6bTestInt(t, p, "n")
		e := w6bTestInt(t, p, "epsilon")
		q := (e * n) / c
		if kp == 2 {
			return w6bTestInt(t, p, "N") - q
		}
		if kp == 3 {
			return q
		}
		return w6bTestInt(t, p, "L") - q
	case "egorov":
		if kp == 1 {
			return w6bTestInt(t, p, "total") - w6bTestInt(t, p, "exceptional")
		}
		if kp == 2 {
			return (w6bTestInt(t, p, "m1") - w6bTestInt(t, p, "exceptional1")) + w6bTestInt(t, p, "m2")
		}
		return w6bTestInt(t, p, "total") - w6bTestInt(t, p, "epsilon")
	case "lusin":
		total := w6bTestInt(t, p, "total")
		if kp == 1 {
			return total - w6bTestInt(t, p, "removed")
		}
		if kp == 2 {
			return total - w6bTestInt(t, p, "intervals")*w6bTestInt(t, p, "each")
		}
		return total - w6bTestInt(t, p, "removed1") - w6bTestInt(t, p, "removed2") - w6bTestInt(t, p, "removed3")
	case "cauchy_schwarz":
		a, b := w6bTestInt(t, p, "a"), w6bTestInt(t, p, "b")
		c, d := w6bTestInt(t, p, "c"), w6bTestInt(t, p, "d")
		if kp == 1 {
			return a*c + b*d
		}
		if kp == 2 {
			return (a+c)*(a+c) + (b+d)*(b+d)
		}
		return (a*d - b*c) * (a*d - b*c)
	case "orthogonal_projection":
		a, b := w6bTestInt(t, p, "a"), w6bTestInt(t, p, "b")
		if kp == 1 {
			return w6bTestInt(t, p, "k")
		}
		if kp == 2 {
			return a*w6bTestInt(t, p, "c") + b*w6bTestInt(t, p, "d")
		}
		k := w6bTestInt(t, p, "k")
		return k * k * (a*a + b*b)
	case "uniform_boundedness":
		if kp == 1 {
			return w6bTestInt(t, p, "bound1")
		}
		if kp == 2 {
			return w6bTestMax(w6bTestInt(t, p, "bound1"), w6bTestInt(t, p, "bound2"))
		}
		return w6bTestMax(w6bTestMax(w6bTestInt(t, p, "bound1"), w6bTestInt(t, p, "bound2")), w6bTestInt(t, p, "bound3"))
	case "open_mapping":
		if kp == 1 {
			return w6bTestAbs(w6bTestInt(t, p, "a")) * w6bTestInt(t, p, "radius")
		}
		if kp == 2 {
			return w6bTestMin(w6bTestAbs(w6bTestInt(t, p, "a")), w6bTestAbs(w6bTestInt(t, p, "b")))
		}
		return w6bTestAbs(w6bTestInt(t, p, "a")*w6bTestInt(t, p, "b")) * w6bTestInt(t, p, "radius")
	case "closed_graph":
		a, b := w6bTestInt(t, p, "a"), w6bTestInt(t, p, "b")
		if kp == 1 {
			return w6bTestMax(w6bTestAbs(a), w6bTestAbs(b))
		}
		if kp == 2 {
			pp, qq := w6bTestInt(t, p, "p"), w6bTestInt(t, p, "q")
			return (pp + qq) + (a*pp + b*qq)
		}
		return w6bTestAbs(a * b)
	case "hahn_banach":
		a, b := w6bTestInt(t, p, "a"), w6bTestInt(t, p, "b")
		if kp == 1 {
			return w6bTestAbs(a)
		}
		if kp == 2 {
			return w6bTestAbs(w6bTestInt(t, p, "p"))
		}
		return w6bTestAbs(a) + w6bTestAbs(b)
	case "weak_convergence":
		if kp == 1 {
			coords := w6bTestInts(t, p, "coords")
			return coords[w6bTestInt(t, p, "N")-1]
		}
		if kp == 2 {
			return w6bTestInt(t, p, "c") * w6bTestPow(w6bTestInt(t, p, "ratio"), w6bTestInt(t, p, "N")-1)
		}
		return w6bTestInt(t, p, "a1") + w6bTestInt(t, p, "a2") + w6bTestInt(t, p, "a3")
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6BDiversity(t *testing.T) {
	for _, kind := range w6bKinds {
		seen := map[string]bool{}
		perKP := [4]map[string]bool{{}, {}, {}, {}}
		for i := 0; i < 200; i++ {
			kp := i%3 + 1
			p := (&w6bGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
			seen[p.Answer] = true
			perKP[kp][p.Answer] = true
		}
		if len(seen) < 8 {
			t.Fatalf("%s: only %d distinct answers in 200 draws: %v", kind, len(seen), seen)
		}
		for kp := 1; kp <= 3; kp++ {
			if len(perKP[kp]) < 8 {
				t.Fatalf("%s kp%d: only %d distinct answers in 200 draws", kind, kp, len(perKP[kp]))
			}
		}
	}
}
