package algebra

// Seven table-driven concepts turned into computed questions.
//
// Each of these used to be a fixed question with a fixed answer. The answer is
// a number that follows from the parameters the question shows, so the
// generator now chooses the parameters, computes the answer and publishes both.
// That is what makes an authored solution schema possible: the schema
// interpolates the same facts the test recomputes from.
//
// The remaining fixed questions in the domain stay exempt. The Dirichlet
// function's continuity, the parity of a named function and the sigmoid's
// limit are statements rather than computations, and a schema that interpolated
// nothing would be a constant string dressed as a solution.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func convInt(n int) string { return strconv.Itoa(n) }

// ---------------------------------------------------------------------------
// The domain of a square root or a reciprocal
// ---------------------------------------------------------------------------

type domainConvertedGen struct{}

func (g *domainConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The shape is fixed by the authored walkthrough in data/lessons/solutions,
	// which describes \sqrt{x - k} and names the facts k and answer, so variety
	// comes from the shift rather than from a second family of functions.
	k := 1 + rand.Intn(60)
	return generator.Problem{
		Question: fmt.Sprintf("The function \\(f(x) = \\sqrt{x - %d}\\) is defined over the reals only where the radicand is not negative. What is the smallest integer in its domain?", k),
		Answer:   convInt(k),
		Explanation: fmt.Sprintf("The radicand \\(x - %d\\) must be at least zero, so \\(x \\ge %d\\). The smallest integer that satisfies that is %d, and every smaller integer makes the radicand negative.",
			k, k, k),
		Facts: map[string]string{"k": convInt(k), "answer": convInt(k)},
	}
}

// ---------------------------------------------------------------------------
// Square roots of perfect squares and cube roots of perfect cubes
// ---------------------------------------------------------------------------

type irrationalConvertedGen struct{}

func (g *irrationalConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	pickSquare := rand.Intn(2) == 0
	if ctx.KPOrdinal == 1 {
		pickSquare = true
	} else if ctx.KPOrdinal >= 2 {
		pickSquare = false
	}
	if pickSquare {
		n := 4 + rand.Intn(22)
		return generator.Problem{
			Question: fmt.Sprintf("What is \\(\\sqrt{%d}\\)?", n*n),
			Answer:   convInt(n),
			Explanation: fmt.Sprintf("The radicand is the square of %d, because \\(%d \\times %d = %d\\). A square root of a perfect square is an integer, so the value is %d.",
				n, n, n, n*n, n),
			Facts: map[string]string{"radicand": convInt(n * n), "root": "square", "answer": convInt(n)},
		}
	}
	n := 2 + rand.Intn(9)
	return generator.Problem{
		Question: fmt.Sprintf("What is \\(\\sqrt[3]{%d}\\)?", n*n*n),
		Answer:   convInt(n),
		Explanation: fmt.Sprintf("The radicand is the cube of %d, because \\(%d \\times %d \\times %d = %d\\). A cube root of a perfect cube is an integer, so the value is %d.",
			n, n, n, n, n*n*n, n),
		Facts: map[string]string{"radicand": convInt(n * n * n), "root": "cube", "answer": convInt(n)},
	}
}

// ---------------------------------------------------------------------------
// Powers with a positive integer exponent
// ---------------------------------------------------------------------------

type powerConvertedGen struct{}

func (g *powerConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	base := 2 + rand.Intn(9)
	exponent := 2 + rand.Intn(4)
	value := 1
	for i := 0; i < exponent; i++ {
		value *= base
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is \\(%d^{%d}\\)?", base, exponent),
		Answer:   convInt(value),
		Explanation: fmt.Sprintf("A power is repeated multiplication: %d used as a factor %d times. That product is %d.",
			base, exponent, value),
		Facts: map[string]string{
			"base": convInt(base), "exponent": convInt(exponent), "answer": convInt(value),
		},
	}
}

// ---------------------------------------------------------------------------
// Reading a point off a line
// ---------------------------------------------------------------------------

type graphAnalysisConvertedGen struct{}

func (g *graphAnalysisConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	slope := 2 + rand.Intn(5)
	intercept := slope * (rand.Intn(9) - 4)
	xIntercept := -intercept / slope
	return generator.Problem{
		Question: fmt.Sprintf("For \\(f(x) = %dx + (%d)\\), at which value of \\(x\\) does the graph cross the horizontal axis?", slope, intercept),
		Answer:   convInt(xIntercept),
		Explanation: fmt.Sprintf("The horizontal axis is where \\(f(x) = 0\\). Solving \\(%dx + (%d) = 0\\) gives \\(x = %d\\), and the slope is not zero so the crossing is unique.",
			slope, intercept, xIntercept),
		Facts: map[string]string{"slope": convInt(slope), "intercept": convInt(intercept), "answer": convInt(xIntercept)},
	}
}

// ---------------------------------------------------------------------------
// The extraneous candidate from squaring a radical equation
// ---------------------------------------------------------------------------

type extraneousRootsConvertedGen struct{}

func (g *extraneousRootsConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := 3 + rand.Intn(8)
	// p is the candidate below b, so squaring produced it but the original
	// equation cannot hold there; q is the candidate at or above b.
	maxP := 3
	if b-1 < maxP {
		maxP = b - 1
	}
	p := 1 + rand.Intn(maxP)
	q := 2*b + 1 - p
	a := b*b - p*q
	if a < 1 {
		// Fall back to the smallest p, which always yields a positive a.
		p = 1
		q = 2*b + 1 - p
		a = b*b - p*q
	}
	return generator.Problem{
		Question: fmt.Sprintf("Squaring both sides of \\(\\sqrt{x + %d} = x - %d\\) produces the candidates \\(x = %d\\) and \\(x = %d\\). Which candidate satisfies the original equation?",
			a, b, p, q),
		Answer: convInt(q),
		Explanation: fmt.Sprintf("The original equation needs \\(x - %d \\ge 0\\), so \\(x \\ge %d\\). The candidate %d is below that bound, and substituting it would equate a square root, which is never negative, with a negative number. The candidate %d satisfies the bound and the equation.",
			b, b, p, q),
		Facts: map[string]string{
			"a": convInt(a), "b": convInt(b),
			"candidates": fmt.Sprintf("%d, %d", p, q),
			"answer":     convInt(q),
		},
	}
}

// ---------------------------------------------------------------------------
// Solving an easy linear system
// ---------------------------------------------------------------------------

type systemsConvertedGen struct{}

func (g *systemsConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x := 2 + rand.Intn(19)
	y := 1 + rand.Intn(15)
	s := x + y
	d := x - y
	return generator.Problem{
		Question: fmt.Sprintf("Solve the system \\(x + y = %d\\) and \\(x - y = %d\\). What is the value of \\(x\\)?", s, d),
		Answer:   convInt(x),
		Explanation: fmt.Sprintf("Adding the two equations cancels \\(y\\): \\(2x = %d\\), so \\(x = %d\\). Subtracting them instead gives \\(2y = %d\\), so \\(y = %d\\), and the pair satisfies both equations.",
			s+d, x, s-d, y),
		Facts: map[string]string{"s": convInt(s), "d": convInt(d), "answer": convInt(x)},
	}
}

// ---------------------------------------------------------------------------
// Counting injective functions
// ---------------------------------------------------------------------------

type injectivityConvertedGen struct{}

type injectivityClaim struct {
	body      string
	injective bool
}

var injectivityClaims = []injectivityClaim{
	{"2x + 3", true},
	{"-4x + 1", true},
	{"5x - 2", true},
	{"x^3", true},
	{"x^2", false},
	{"|x|", false},
	{"x^2 + 1", false},
	{"x^4", false},
}

func (g *injectivityConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	order := rand.Perm(len(injectivityClaims))
	picked := []injectivityClaim{
		injectivityClaims[order[0]],
		injectivityClaims[order[1]],
		injectivityClaims[order[2]],
	}
	names := []string{"f", "g", "h"}
	parts := make([]string, 3)
	count := 0
	for i, c := range picked {
		parts[i] = fmt.Sprintf("%s(x) = %s", names[i], c.body)
		if c.injective {
			count++
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("Consider \\(%s\\), \\(%s\\) and \\(%s\\) on the real line. How many of the three are injective?",
			parts[0], parts[1], parts[2]),
		Answer: convInt(count),
		Explanation: fmt.Sprintf("A function is injective when no two different inputs share an output. A non-zero linear function and an odd power are injective on the real line, while \\(x^2\\), \\(|x|\\), \\(x^2 + 1\\) and \\(x^4\\) are not, because opposite inputs share an output. Of the three listed, %d are injective.",
			count),
		Facts: map[string]string{
			"functions": fmt.Sprintf("%s; %s; %s", parts[0], parts[1], parts[2]),
			"answer":    convInt(count),
		},
	}
}

// ---------------------------------------------------------------------------
// The factor theorem (alg.poly.factor_theorem)
// ---------------------------------------------------------------------------

type factorTheoremConvertedGen struct{}

func (g *factorTheoremConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	roots := []int{-3, -2, -1, 1, 2, 3}
	r := roots[rand.Intn(len(roots))]
	a := rand.Intn(9) - 4
	k := rand.Intn(17) - 8
	c := -(r*r*r + a*r*r + k*r)
	if c == 0 {
		k += 2
		c = -(r*r*r + a*r*r + k*r)
	}

	factorStr := fmt.Sprintf("x - %d", r)
	if r < 0 {
		factorStr = fmt.Sprintf("x + %d", -r)
	}

	constTerm := r*r*r + a*r*r + c
	rhs := -constTerm

	q := fmt.Sprintf("The linear polynomial \\(%s\\) is a factor of \\(p(x) = x^3 + (%d)x^2 + kx + (%d)\\). By the factor theorem, what is the value of \\(k\\)?",
		factorStr, a, c)
	exp := fmt.Sprintf("By the factor theorem, \\(%s\\) divides \\(p(x)\\) exactly when \\(p(%d) = 0\\). Evaluating \\(p(%d) = (%d)^3 + (%d)(%d)^2 + k(%d) + (%d) = 0\\) simplifies to \\(%d + (%d)k = 0\\), so \\(%dk = %d\\) and \\(k = %d\\).",
		factorStr, r, r, r, a, r, r, c, constTerm, r, r, rhs, k)

	return generator.Problem{
		Question:    q,
		Answer:      convInt(k),
		Explanation: exp,
		Facts: map[string]string{
			"factor": factorStr,
			"root":   convInt(r),
			"a":      convInt(a),
			"c":      convInt(c),
			"rhs":    convInt(rhs),
			"answer": convInt(k),
		},
	}
}

// registerConversions registers the concepts of this batch that used to be
// table-driven. Their old registrations are removed from Register, because a
// concept id may be registered only once.
func registerConversions(reg *generator.Registry) {
	reg.Register("alg.func.domain", &domainConvertedGen{})
	reg.Register("alg.func.irrational", &irrationalConvertedGen{})
	reg.Register("alg.func.power", &powerConvertedGen{})
	reg.Register("alg.func.graph_analysis", &graphAnalysisConvertedGen{})
	reg.Register("alg.eq.extraneous_roots", &extraneousRootsConvertedGen{})
	reg.Register("alg.systems.concept", &systemsConvertedGen{})
	reg.Register("alg.func.injectivity", &injectivityConvertedGen{})
	reg.Register("alg.poly.factor_theorem", &factorTheoremConvertedGen{})
}
