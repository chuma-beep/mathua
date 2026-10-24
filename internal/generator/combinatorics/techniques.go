package combinatorics

// The second combinatorics block: the counting techniques that are not a single formula.
//
// The first block was the four standard counts -- arrangements, multisets, circles, derangements --
// plus the Catalan numbers. These are the arguments instead: counting one set in two ways
// (double_counting), evaluating a convolution of binomial coefficients (vandermonde), counting by
// subtracting the cases that miss (surjections, and the Stirling number that divides them by the
// labelling), and the three classic sequences that a counting course ends on: partitions of a set
// (Bell), partitions of an integer, and shortest paths on a grid.
//
// Bell numbers and integer partitions are tabulated rather than computed, because a closed form does
// not exist for either; the tables stop where the values still fit in the answer field, and the
// generator never indexes past them.
import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerCountingTechniques(reg *generator.Registry) {
	reg.Register("discrete.combinatorics.vandermonde", &vandermondeGen{})
	reg.Register("discrete.combinatorics.double_counting", &doubleCountingGen{})
	reg.Register("discrete.combinatorics.surjections", &surjectionGen{})
	reg.Register("discrete.combinatorics.stirling_second", &stirlingSecondGen{})
	reg.Register("discrete.combinatorics.bell", &bellGen{})
	reg.Register("discrete.combinatorics.integer_partitions", &integerPartitionGen{})
	reg.Register("discrete.combinatorics.lattice_paths", &latticePathGen{})
	reg.Register("discrete.combinatorics.compositions", &compositionGen{})
}

// vandermondeGen asks for the convolution sum of two binomial rows, which is one binomial
// coefficient: sum over i of C(a,i)C(b,r-i) = C(a+b,r).
type vandermondeGen struct{}

func (g *vandermondeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 2 + pick(2+int(ctx.Difficulty*3))
	b := 2 + pick(2+int(ctx.Difficulty*3))
	r := 1 + pick(a+b)
	v := choose(a+b, r)
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate \\(\\sum_{i=0}^{%d} \\binom{%d}{i}\\binom{%d}{%d-i}\\), taking a binomial coefficient to be zero when its lower index falls outside \\(0 \\le i \\le n\\).", r, a, b, r),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("Vandermonde's identity evaluates the convolution of two binomial rows as a single one: the sum is C(%d+%d, %d) = C(%d, %d) = %d.",
			a, b, r, a+b, r, v),
		Facts: map[string]string{"value": fmt.Sprintf("%d", v)},
	}
}

// doubleCountingGen asks for the number of subsets containing a fixed element, which is 2^(n-1) by
// the pairing argument.
type doubleCountingGen struct{}

func (g *doubleCountingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + pick(3+int(ctx.Difficulty*5))
	return generator.Problem{
		Question: fmt.Sprintf("How many of the subsets of a set with %d elements contain a particular one of those elements?", n),
		Answer:   fmt.Sprintf("%d", 1<<uint(n-1)),
		Explanation: fmt.Sprintf("Count the pairs (element, subset containing it) two ways. Each of the %d elements lies in 2^(%d-1) subsets, so there are %d*%d pairs; each subset contains as many elements as its size, so the total is also the sum of the subset sizes. The direct way is simpler: every subset of the other %d elements can be extended by that element, giving 2^%d = %d.",
			n, n, n, 1<<uint(n-1), n-1, n-1, 1<<uint(n-1)),
		Facts: map[string]string{"n": fmt.Sprintf("%d", n), "subsets": fmt.Sprintf("%d", 1<<uint(n-1))},
	}
}

// surjectionGen counts onto functions by inclusion and exclusion over the images that are missed.
type surjectionGen struct{}

func (g *surjectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + pick(2) // the codomain, kept small so the count stays small
	m := n + 1 + pick(2+int(ctx.Difficulty*2))
	total := 0
	for i := 0; i <= n; i++ {
		term := choose(n, i) * pow(n-i, m)
		if i%2 == 0 {
			total += term
		} else {
			total -= term
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many functions from a set with %d elements onto a set with %d elements are there?", m, n),
		Answer:   fmt.Sprintf("%d", total),
		Explanation: fmt.Sprintf("All %d^%d functions, minus the ones missing at least one image: inclusion and exclusion over the %d images gives %d onto functions.",
			n, m, n, total),
		Facts: map[string]string{"onto": fmt.Sprintf("%d", total)},
	}
}

// stirlingSecondGen counts the partitions of a labelled set into k unlabelled blocks, which is the
// number of onto functions divided by the ways to label the blocks.
type stirlingSecondGen struct{}

func (g *stirlingSecondGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 2 + pick(2) // 2 or 3 blocks
	n := k + 1 + pick(2+int(ctx.Difficulty*4))
	return generator.Problem{
		Question: fmt.Sprintf("In how many ways can a set of %d distinct items be split into exactly %d non-empty groups, where the groups themselves are not labelled?", n, k),
		Answer:   fmt.Sprintf("%d", stirling2(n, k)),
		Explanation: fmt.Sprintf("That count is the Stirling number of the second kind S(%d, %d) = %d. One way to see it: there are %d onto assignments to %d labelled groups, and dividing by the %d! ways to label the groups removes the labelling.",
			n, k, stirling2(n, k), ontoCount(n, k), k, k),
		Facts: map[string]string{"S": fmt.Sprintf("%d", stirling2(n, k))},
	}
}

// bellGen asks how many partitions a set has in total, which is the sum of its Stirling numbers.
type bellGen struct{}

func (g *bellGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Bell numbers up to B(15); beyond that the value stops reading as a count
	// a learner could check, but every result here is computed, not looked up.
	n := 1 + pick(2+int(ctx.Difficulty*13))
	if n > 15 {
		n = 15
	}
	b := bellNumber(n)
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("In how many ways can a set of %d distinct items be split into non-empty groups, where the number of groups is not fixed and the groups are not labelled?", n),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("That is the Bell number B(%d) = %d: sum the Stirling numbers S(%d, k) over every k from 1 to %d, which counts the same partitions split by how many groups they have.", n, b, n, n),
			Facts:       map[string]string{"B": strconv.Itoa(b), "n": strconv.Itoa(n)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the Bell number B(%d), the number of partitions of a set with %d elements?", n, n),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("B(%d) = %d. The recurrence B(n+1) = sum over k of C(n,k) B(k) adds the block containing the (n+1)-th element, so the count grows faster than any polynomial.", n, b),
			Facts:       map[string]string{"B": strconv.Itoa(b), "n": strconv.Itoa(n)},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("A poem has %d lines. How many rhyme schemes are possible, counting schemes that match in pattern as the same?", n),
			Answer:      strconv.Itoa(b),
			Explanation: fmt.Sprintf("Each rhyme scheme is a partition of the %d line positions into rhyme classes, so the count is B(%d) = %d.", n, n, b),
			Facts:       map[string]string{"B": strconv.Itoa(b), "n": strconv.Itoa(n)},
		}
	}
}

// bellNumber returns the Bell number B(n) through the Stirling recurrence, so
// the generator owns no table of values.
func bellNumber(n int) int {
	stirling := make([][]int, n+1)
	for i := range stirling {
		stirling[i] = make([]int, n+1)
	}
	stirling[0][0] = 1
	for i := 1; i <= n; i++ {
		for k := 1; k <= i; k++ {
			stirling[i][k] = stirling[i-1][k-1] + k*stirling[i-1][k]
		}
	}
	total := 0
	for k := 0; k <= n; k++ {
		total += stirling[n][k]
	}
	return total
}

// integerPartitionGen asks for the number of ways to write n as a sum of positive integers where
// order does not matter. There is no closed form, so the values are tabulated.
type integerPartitionGen struct{}

func (g *integerPartitionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := []int{1, 2, 3, 5, 7, 11, 15, 22, 30, 42, 56, 77}
	n := 1 + pick(2+int(ctx.Difficulty*6))
	if n > len(p) {
		n = len(p)
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many ways can %d be written as a sum of positive integers, if two sums that differ only in the order of their terms are the same?", n),
		Answer:   fmt.Sprintf("%d", p[n-1]),
		Explanation: fmt.Sprintf("That is the partition number p(%d) = %d. Unlike compositions, order is ignored, so 1+2 and 2+1 are one partition -- which is why no closed form exists and the values are tabulated.",
			n, p[n-1]),
		Facts: map[string]string{"p": fmt.Sprintf("%d", p[n-1])},
	}
}

// latticePathGen counts monotone paths across a grid, which is a binomial coefficient.
type latticePathGen struct{}

func (g *latticePathGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 1 + pick(2+int(ctx.Difficulty*3))
	n := 1 + pick(2+int(ctx.Difficulty*3))
	return generator.Problem{
		Question: fmt.Sprintf("How many shortest paths run from \\((0,0)\\) to \\((%d,%d)\\) if each step moves one unit right or one unit up?", m, n),
		Answer:   fmt.Sprintf("%d", choose(m+n, m)),
		Explanation: fmt.Sprintf("A shortest path is a word of %d right steps and %d up steps, so choosing which of the %d positions hold the right steps counts them: C(%d, %d) = %d.",
			m, n, m+n, m+n, m, choose(m+n, m)),
		Facts: map[string]string{"paths": fmt.Sprintf("%d", choose(m+n, m))},
	}
}

// compositionGen counts ordered sums into a fixed number of positive parts, which is stars and bars
// with the bars placed between items rather than among them.
type compositionGen struct{}

func (g *compositionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := 2 + pick(1+int(ctx.Difficulty*2))
	n := k + 1 + pick(3+int(ctx.Difficulty*4))
	v := choose(n-1, k-1)
	return generator.Problem{
		Question: fmt.Sprintf("How many ways can %d be written as a sum of exactly %d positive integers, if the order of the terms matters?", n, k),
		Answer:   fmt.Sprintf("%d", v),
		Explanation: fmt.Sprintf("Lay %d ones in a row; choosing %d of the %d gaps between them as the places to cut gives C(%d, %d) = %d ordered sums.",
			n, k-1, n-1, n-1, k-1, v),
		Facts: map[string]string{"compositions": fmt.Sprintf("%d", v)},
	}
}

func pow(base, exp int) int {
	r := 1
	for i := 0; i < exp; i++ {
		r *= base
	}
	return r
}

// ontoCount is the number of surjections from m items onto n labelled images.
func ontoCount(m, n int) int {
	total := 0
	for i := 0; i <= n; i++ {
		term := choose(n, i) * pow(n-i, m)
		if i%2 == 0 {
			total += term
		} else {
			total -= term
		}
	}
	return total
}

// stirling2 is the number of partitions of an m-element set into exactly n non-empty blocks,
// computed by the recurrence S(m,n) = n*S(m-1,n) + S(m-1,n-1).
func stirling2(m, n int) int {
	if n == 0 || m < n {
		return 0
	}
	if n == 1 || m == n {
		return 1
	}
	return n*stirling2(m-1, n) + stirling2(m-1, n-1)
}
