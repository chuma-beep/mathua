package discrete

// The 2026-10-17 batch: sixteen concepts that take discrete_math from 39 to the domain target of 55.
//
// Every answer here is computed from the instance, never tabulated: the counterexample concept scans
// for the first composite value, the matching concept searches for a maximum matching, and the
// recurrence concept derives its closed form from the roots and coefficients it chose. Facts are
// published for the authored solution schemas, and each concept's grading_type is one the router
// already supports.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("discrete.logic.inference", &inferenceGen{})
	reg.Register("discrete.logic.predicates", &predicatesGen{})
	reg.Register("discrete.logic.satisfiability", &satisfiabilityGen{})
	reg.Register("discrete.proof.contrapositive", &contrapositiveGen{})
	reg.Register("discrete.proof.counterexample", &counterexampleGen{})
	reg.Register("discrete.proof.structural", &structuralGen{})
	reg.Register("discrete.sets.powerset", &powersetGen{})
	reg.Register("discrete.sets.partial_order", &partialOrderGen{})
	reg.Register("discrete.sets.countable", &countableGen{})
	reg.Register("discrete.combinatorics.complement", &complementGen{})
	reg.Register("discrete.graphs.bipartite", &bipartiteGen{})
	reg.Register("discrete.graphs.planar", &planarGen{})
	reg.Register("discrete.graphs.cayley", &cayleyGen{})
	reg.Register("discrete.graphs.degree_sequence", &degreeSequenceGen{})
	reg.Register("discrete.sequences.linear_recurrence", &linearRecurrenceGen{})
	reg.Register("discrete.sequences.growth_rates", &growthRatesGen{})
}

// ----- helpers -----

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func divisors(n int) []int {
	var out []int
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			out = append(out, d)
			if d != n/d {
				out = append(out, n/d)
			}
		}
	}
	sort.Ints(out)
	return out
}

func intsString(xs []int, sep string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, sep)
}

// maxMatching searches for a maximum matching in a bipartite graph. Sizes here are at most 4x4, so
// exhaustive search is the honest implementation: it cannot be wrong about the answer.
func maxMatching(leftCount int, adj [][]bool) int {
	best := 0
	var walk func(i int, used map[int]bool, size int)
	walk = func(i int, used map[int]bool, size int) {
		if i == leftCount {
			if size > best {
				best = size
			}
			return
		}
		walk(i+1, used, size)
		for r, ok := range adj[i] {
			if ok && !used[r] {
				used[r] = true
				walk(i+1, used, size+1)
				delete(used, r)
			}
		}
	}
	walk(0, map[int]bool{}, 0)
	return best
}

// ----- 1. inference -----

type inferenceGen struct{}

func (g *inferenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return inferenceProblem()
}

// ----- 2. predicates -----

type predicatesGen struct{}

func (g *predicatesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return predicatesProblem()
}

// ----- 3. satisfiability -----

type satisfiabilityGen struct{}

func (g *satisfiabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return satisfiabilityProblem()
}

// ----- 4. contrapositive -----

type contrapositiveGen struct{}

func (g *contrapositiveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return contrapositiveProblem()
}

// ----- 5. counterexample -----

type counterexampleGen struct{}

func (g *counterexampleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The constants are found by scanning rather than chosen by hand: any c whose
	// polynomial n^2 + n + c stays prime for at least two steps makes the question
	// meaningful, and the run length is re-derived for the constant that was drawn, so a
	// wrong constant cannot ship a wrong answer.
	candidates := make([]int, 0, 16)
	for c := 2; c <= 200; c++ {
		if runLength(c) >= 2 {
			candidates = append(candidates, c)
		}
	}
	c := candidates[rand.Intn(len(candidates))]
	n := runLength(c)
	value := n*n + n + c
	f1 := 0
	for d := 2; d*d <= value; d++ {
		if value%d == 0 {
			f1 = d
			break
		}
	}
	p := generator.Problem{
		Facts: map[string]string{
			"c": strconv.Itoa(c), "n": strconv.Itoa(n), "value": strconv.Itoa(value),
			"factor1": strconv.Itoa(f1), "factor2": strconv.Itoa(value / f1),
			"last_prime": strconv.Itoa(n - 1),
		},
	}
	switch rand.Intn(3) {
	case 0:
		p.Question = fmt.Sprintf("Let f(n) = n\u00b2 + n + %d. Find the smallest positive integer n for which f(n) is not prime.", c)
		p.Answer = strconv.Itoa(n)
	case 1:
		p.Question = fmt.Sprintf("For how many consecutive values of n, starting at n = 0, is f(n) = n\u00b2 + n + %d prime?", c)
		p.Answer = strconv.Itoa(n)
	default:
		p.Question = fmt.Sprintf("Let f(n) = n\u00b2 + n + %d. What is the value of f(n) at the smallest n for which f(n) is not prime?", c)
		p.Answer = strconv.Itoa(value)
	}
	p.Explanation = fmt.Sprintf("f(%d) = %d = %d \u00b7 %d, which is composite, and f(n) is prime for every n from 0 to %d.",
		n, value, f1, value/f1, n-1)
	return p
}

// runLength is the number of consecutive primes produced by f(n) = n^2 + n + c
// starting at n = 0. It is the smallest n at which f(n) is composite.
func runLength(c int) int {
	n := 0
	for isPrime(n*n + n + c) {
		n++
	}
	return n
}

// ----- 6. structural induction -----

type structuralGen struct{}

func (g *structuralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return structuralProblem()
}

// ----- 7. power set -----

type powersetGen struct{}

func (g *powersetGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(5) // 4..8
	k := 1 + rand.Intn(n-1)
	choose := 1
	for i := 1; i <= k; i++ {
		choose = choose * (n - k + i) / i
	}
	total := 1 << n
	q := fmt.Sprintf("A set has %d elements. How many of its subsets have exactly %d elements?", n, k)
	p := generator.Problem{
		Question: q,
		Answer:   strconv.Itoa(choose),
		Facts: map[string]string{
			"n": strconv.Itoa(n), "k": strconv.Itoa(k), "value": strconv.Itoa(choose),
			"total": strconv.Itoa(total),
		},
	}
	p.Explanation = fmt.Sprintf("Choosing a %d-element subset means choosing %d of the %d elements, which is C(%d, %d) = %d. The whole power set has 2^%d = %d members.",
		k, k, n, n, k, choose, n, total)
	return p
}

// ----- 8. partial order -----

type partialOrderGen struct{}

func (g *partialOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The pair count is computed from the divisors of N, so any N is a valid
	// draw. Twelve values replace the original four, and they differ in the
	// number of divisors rather than only in the numeral.
	ns := []int{12, 18, 20, 24, 28, 30, 36, 40, 42, 45, 48, 60}
	N := ns[rand.Intn(len(ns))]
	ds := divisors(N)
	count := 0
	for i := range ds {
		for j := range ds {
			if ds[i] < ds[j] && ds[j]%ds[i] == 0 {
				count++
			}
		}
	}
	q := fmt.Sprintf("The positive divisors of %d are ordered by divisibility. How many pairs (a, b) with a < b and a | b are there?", N)
	p := generator.Problem{
		Question: q,
		Answer:   strconv.Itoa(count),
		Facts: map[string]string{
			"N": strconv.Itoa(N), "set": "\\{" + intsString(ds, ", ") + "\\}", "value": strconv.Itoa(count),
		},
	}
	p.Explanation = fmt.Sprintf("On %s, a pair counts when the smaller divides the larger; there are %d such pairs.", "{"+intsString(ds, ", ")+"}", count)
	return p
}

// ----- 9. countable -----

type countableGen struct{}

func (g *countableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return countableProblem()
}

// ----- 10. counting by complement -----

type complementGen struct{}

func (g *complementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 2 + rand.Intn(3) // 2..4 letters, one of them a
	L := 3 + rand.Intn(3) // 3..5
	total := 1
	for i := 0; i < L; i++ {
		total *= m
	}
	none := 1
	for i := 0; i < L; i++ {
		none *= m - 1
	}
	answer := total - none
	q := fmt.Sprintf("How many strings of length %d over an alphabet of %d letters contain at least one \u2018a\u2019?", L, m)
	p := generator.Problem{
		Question: q,
		Answer:   strconv.Itoa(answer),
		Facts: map[string]string{
			"letters": strconv.Itoa(m), "length": strconv.Itoa(L),
			"total": strconv.Itoa(total), "none": strconv.Itoa(none), "value": strconv.Itoa(answer),
		},
	}
	p.Explanation = fmt.Sprintf("Count the complement instead: %d^%d = %d strings in all, and %d^%d = %d contain no \u2018a\u2019. The difference is %d.",
		m, L, total, m-1, L, none, answer)
	return p
}

// ----- 11. bipartite matching -----

type bipartiteGen struct{}

func (g *bipartiteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	leftCount := 3 + rand.Intn(2)  // 3..4
	rightCount := 3 + rand.Intn(2) // 3..4
	adj := make([][]bool, leftCount)
	for i := range adj {
		adj[i] = make([]bool, rightCount)
	}
	// Give every left vertex at least one edge, then add a few more: a graph with an isolated
	// vertex would make the matching question trivial in a way the learner cannot see.
	for i := 0; i < leftCount; i++ {
		adj[i][rand.Intn(rightCount)] = true
	}
	for k := 0; k < 2+rand.Intn(2); k++ {
		adj[rand.Intn(leftCount)][rand.Intn(rightCount)] = true
	}
	size := maxMatching(leftCount, adj)
	var edges []string
	for i := 0; i < leftCount; i++ {
		for j := 0; j < rightCount; j++ {
			if adj[i][j] {
				edges = append(edges, fmt.Sprintf("%c%d", 'a'+i, j+1))
			}
		}
	}
	q := fmt.Sprintf("A bipartite graph has left vertices %s and right vertices %s, with edges %s. What is the size of a maximum matching?",
		commas(strings.Split(letters(leftCount), "")), commas(strings.Split(digits(rightCount), "")), strings.Join(edges, ", "))
	p := generator.Problem{
		Question: q,
		Answer:   strconv.Itoa(size),
		Facts: map[string]string{
			"left": commas(strings.Split(letters(leftCount), "")), "right": commas(strings.Split(digits(rightCount), "")),
			"edges": strings.Join(edges, ", "), "value": strconv.Itoa(size),
		},
	}
	p.Explanation = fmt.Sprintf("Pair as many left and right vertices as possible using distinct vertices: %d pairs is the maximum here, so no matching of size %d exists.", size, size+1)
	return p
}

func letters(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(byte('a' + i))
	}
	return b.String()
}

// commas joins single characters for display: "abc" becomes "a, b, c".
func commas(parts []string) string { return strings.Join(parts, ", ") }

func digits(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}

// ----- 12. planar graphs -----

type planarGen struct{}

func (g *planarGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	V := 6 + rand.Intn(5) // 6..10
	maxE := 3*V - 6
	E := V + rand.Intn(maxE-V+1) // V..3V-6, so at least one cycle
	F := E - V + 2
	q := fmt.Sprintf("A connected planar graph has %d vertices and %d edges. How many faces does it have, counting the outer face?", V, E)
	p := generator.Problem{
		Question: q,
		Answer:   strconv.Itoa(F),
		Facts: map[string]string{
			"v": strconv.Itoa(V), "e": strconv.Itoa(E), "value": strconv.Itoa(F),
		},
	}
	p.Explanation = fmt.Sprintf("Euler's formula for a connected planar graph is V - E + F = 2, so F = %d - %d + 2 = %d.", E, V, F)
	return p
}

// ----- 13. Cayley's formula -----

type cayleyGen struct{}

func (g *cayleyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 3 + rand.Intn(8)
	value := 1
	for i := 0; i < n-2; i++ {
		value *= n
	}
	answer := strconv.Itoa(value)
	facts := map[string]string{"n": strconv.Itoa(n), "exponent": strconv.Itoa(n - 2), "value": answer}
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("How many labelled trees are there on %d given vertices?", n),
			Answer:      answer,
			Explanation: fmt.Sprintf("Cayley's formula gives n^(n-2) labelled trees on n vertices: %d^%d = %d.", n, n-2, value),
			Facts:       facts,
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("How many distinct Pruefer sequences of length %d are there over an alphabet of %d symbols?", n-2, n),
			Answer:      answer,
			Explanation: fmt.Sprintf("Each position holds one of %d symbols, so the count is %d^%d = %d. Pruefer sequences are in bijection with labelled trees, which is why the two counts agree.", n, n, n-2, value),
			Facts:       facts,
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Cayley's formula says the number of labelled trees on n vertices is n^(n-2). For n = %d, what exponent appears?", n),
			Answer:      strconv.Itoa(n - 2),
			Explanation: fmt.Sprintf("Substituting n = %d gives the exponent %d - 2 = %d, so there are %d labelled trees.", n, n, n-2, value),
			Facts:       facts,
		}
	}
}

// ----- 14. degree sequence -----

type degreeSequenceGen struct{}

func (g *degreeSequenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 + rand.Intn(3) // left part size, 2..4
	b := 2 + rand.Intn(3) // right part size, 2..4
	var degrees []int
	for i := 0; i < a; i++ {
		degrees = append(degrees, b)
	}
	for i := 0; i < b; i++ {
		degrees = append(degrees, a)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(degrees)))
	sum := 0
	for _, d := range degrees {
		sum += d
	}
	answer := intsString(degrees, ",")
	q := fmt.Sprintf("Give the degree sequence of the complete bipartite graph K(%d, %d) in non-increasing order, separated by commas (for example 3,3,2).", a, b)
	p := generator.Problem{
		Question: q,
		Answer:   answer,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "sequence": answer,
			"sum": strconv.Itoa(sum), "edges": strconv.Itoa(a * b),
		},
	}
	p.Explanation = fmt.Sprintf("Each of the %d vertices in the first part has degree %d and each of the %d in the second has degree %d, so the sequence is %s. Its sum is %d = 2 \u00b7 %d, as the handshake lemma requires.",
		a, b, b, a, answer, sum, a*b)
	return p
}

// ----- 15. linear recurrence -----

type linearRecurrenceGen struct{}

func (g *linearRecurrenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	bases := []int{2, 3, 4, 5}
	rand.Shuffle(len(bases), func(i, j int) { bases[i], bases[j] = bases[j], bases[i] })
	r1, r2 := bases[0], bases[1]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	a := 1 + rand.Intn(2) // 1..2
	b := 1 + rand.Intn(2)
	a0 := a + b
	a1 := a*r1 + b*r2
	c1 := r1 + r2
	c2 := r1 * r2
	term1 := powerTerm(a, r1)
	term2 := powerTerm(b, r2)
	answer := term1 + " + " + term2
	q := fmt.Sprintf("A sequence satisfies a(n) = %d\u00b7a(n-1) - %d\u00b7a(n-2), with a(0) = %d and a(1) = %d. Write a closed form for a(n) (use ^ for powers and * for multiplication).",
		c1, c2, a0, a1)
	p := generator.Problem{
		Question: q,
		Answer:   answer,
		Facts: map[string]string{
			"r1": strconv.Itoa(r1), "r2": strconv.Itoa(r2),
			"a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"c1": strconv.Itoa(c1), "c2": strconv.Itoa(c2),
			"a0": strconv.Itoa(a0), "a1": strconv.Itoa(a1),
			"term1": term1, "term2": term2, "value": answer,
		},
	}
	p.Explanation = fmt.Sprintf("The characteristic equation is x\u00b2 - %d\u00b7x + %d = 0, with roots %d and %d, so a(n) = A\u00b7%d^n + B\u00b7%d^n. The initial conditions give A + B = %d and %dA + %dB = %d, hence A = %d and B = %d, and a(n) = %s.",
		c1, c2, r1, r2, r1, r2, a0, r1, r2, a1, a, b, answer)
	return p
}

func powerTerm(coef, base int) string {
	if coef == 1 {
		return fmt.Sprintf("%d^n", base)
	}
	return fmt.Sprintf("%d*%d^n", coef, base)
}

// ----- 16. growth rates -----

type growthRatesGen struct{}

func (g *growthRatesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return growthRatesProblem()
}
