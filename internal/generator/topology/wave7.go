package topology

// Wave 7 topology, first batch: nets, sequential compactness, deck
// transformations and the cup product. Answers are exact integers computed from
// published facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerWave7(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"topo.nets":               "nets",
		"topo.sequential_compact": "sequential_compact",
		"topo.deck":               "deck",
		"topo.cup_product":        "cup_product",
	} {
		reg.Register(id, &w7TopologyGen{kind: kind})
	}
}

type w7TopologyGen struct{ kind string }

func (g *w7TopologyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "nets":
		return w7Nets(kp)
	case "sequential_compact":
		return w7SequentialCompact(kp)
	case "deck":
		return w7Deck(kp)
	case "cup_product":
		return w7CupProduct(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w7T(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w7TI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w7TR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w7TPhi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if mathutil.GCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w7Nets(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w7TR(3, 15)
		w7TI(f, "n", n)
		return w7T(kp, fmt.Sprintf("Let \\(D\\) be the finite subsets of \\(\\{1,\\dots,%d\\}\\) ordered by inclusion and let \\(x_S=|S|\\). What is the limit of the net \\((x_S)\\)?", n), strconv.Itoa(n),
			fmt.Sprintf("The net is eventually %d because the full set is an upper bound in the directed set.", n), f)
	case 2:
		a := w7TR(-3, 3)
		b := w7TR(1, 4)
		f["a"] = strconv.Itoa(a)
		f["b"] = strconv.Itoa(b)
		return w7T(kp, fmt.Sprintf("A net indexed by \\(\\mathbb{N}\\) has values \\(a_k=%d+(-1)^k\\cdot%d\\). How many distinct cluster points does it have?", a, b), "2",
			"The even and odd subsequences converge to \\(a+b\\) and \\(a-b\\), which are distinct.", f)
	default:
		c := w7TR(-5, 5)
		k := w7TR(2, 6)
		f["c"] = strconv.Itoa(c)
		f["k"] = strconv.Itoa(k)
		return w7T(kp, fmt.Sprintf("A net is eventually constant with value %d after index %d. What is its limit?", c, k), strconv.Itoa(c),
			"An eventually constant net converges to its eventual value.", f)
	}
}

func w7SequentialCompact(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		m := w7TR(2, 12)
		w7TI(f, "m", m)
		return w7T(kp, fmt.Sprintf("The sequence \\(a_n=n\\bmod %d\\) takes values in \\(\\{0,\\dots,%d\\}\\). How many distinct subsequential limits does it have?", m, m-1), strconv.Itoa(m),
			"Every residue repeats infinitely often, so each residue is a subsequential limit.", f)
	case 2:
		c := w7TR(1, 12)
		f["c"] = strconv.Itoa(c)
		return w7T(kp, fmt.Sprintf("The sequence repeats the four values \\(%d,0,-%d,0\\). How many distinct subsequential limits does it have?", c, c), "3",
			"The constant subsequences give the three values c, 0 and -c.", f)
	default:
		lo := w7TR(0, 4)
		hi := lo + w7TR(1, 5)
		closed := rand.Intn(2) == 0
		f["lo"] = strconv.Itoa(lo)
		f["hi"] = strconv.Itoa(hi)
		f["closed"] = strconv.FormatBool(closed)
		left, right := "(", ")"
		if closed {
			left, right = "[", "]"
		}
		ans := 0
		if closed {
			ans = 1
		}
		return w7T(kp, fmt.Sprintf("Is the interval \\(%s%d,%d%s\\) sequentially compact as a subset of \\(\\mathbb{R}\\)? Enter 1 for yes and 0 for no.", left, lo, hi, right), strconv.Itoa(ans),
			"A closed bounded interval in \\(\\mathbb{R}\\) is compact and hence sequentially compact; an open interval is not.", f)
	}
}

func w7Deck(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w7TR(2, 12)
		w7TI(f, "n", n)
		return w7T(kp, fmt.Sprintf("The covering \\(S^1\\to S^1\\), \\(z\\mapsto z^{%d}\\), has how many deck transformations?", n), strconv.Itoa(n),
			"The deck group is cyclic of order equal to the number of sheets.", f)
	case 2:
		n := w7TR(3, 15)
		w7TI(f, "n", n)
		return w7T(kp, fmt.Sprintf("A regular covering has cyclic deck group of order %d. How many generators does the deck group have?", n), strconv.Itoa(w7TPhi(n)),
			fmt.Sprintf("A cyclic group of order %d has \\(\\varphi(%d)=%d\\) generators.", n, n, w7TPhi(n)), f)
	default:
		m := w7TR(2, 12)
		w7TI(f, "m", m)
		return w7T(kp, fmt.Sprintf("The universal cover of a wedge of %d circles has deck group the free group of rank what?", m), strconv.Itoa(m),
			"The fundamental group of the wedge is free of rank m, and the deck group is isomorphic to it.", f)
	}
}

func w7CupProduct(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		p := w7TR(1, 4)
		q := w7TR(1, 4)
		f["p"] = strconv.Itoa(p)
		f["q"] = strconv.Itoa(q)
		return w7T(kp, fmt.Sprintf("If \\(\\alpha\\) has degree %d and \\(\\beta\\) has degree %d, what is the degree of \\(\\alpha\\cup\\beta\\)?", p, q), strconv.Itoa(p+q),
			"Cup products add degrees.", f)
	case 2:
		p := w7TR(1, 4)
		q := w7TR(1, 4)
		f["p"] = strconv.Itoa(p)
		f["q"] = strconv.Itoa(q)
		sign := 1
		if (p*q)%2 == 1 {
			sign = -1
		}
		return w7T(kp, fmt.Sprintf("For classes of degrees %d and %d, what is the sign in \\(\\alpha\\cup\\beta=(-1)^{pq}\\beta\\cup\\alpha\\)?", p, q), strconv.Itoa(sign),
			fmt.Sprintf("Evaluate \\((-1)^{%d\\cdot%d}\\).", p, q), f)
	default:
		n := w7TR(1, 12)
		w7TI(f, "n", n)
		return w7T(kp, fmt.Sprintf("In \\(H^*(\\mathbb{CP}^{%d})\\), the generator \\(\\alpha\\in H^2\\) satisfies \\(\\alpha^k=0\\) for the first time at which \\(k\\)?", n), strconv.Itoa(n+1),
			"The cohomology ring is truncated above degree \\(2n\\), so the first vanishing power is \\(n+1\\).", f)
	}
}
