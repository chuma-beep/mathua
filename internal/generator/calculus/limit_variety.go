package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The limit concepts mix questions with one exact answer and questions that
// classify a theorem. A classification cannot be typed into a maths field, so
// those shapes are multiple choice and every pool contains both answers.

// limitChoice builds a multiple choice problem with an explanation, since a
// classification answer cannot be typed into a maths field.
func limitChoice(question, answer, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, answer, distractors...)
	p.Explanation = explanation
	return p
}

type limitInfinityGen struct{}

func (g *limitInfinityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(3) {
	case 0:
		k := randRange(ctx.Difficulty, 2, 9)
		c := randRange(ctx.Difficulty, 2, 5)
		a := k * c
		b := randRange(ctx.Difficulty, 1, 9)
		d := randRange(ctx.Difficulty, 1, 9)
		p := randRange(ctx.Difficulty, 1, 3)
		return generator.Problem{
			Question: fmt.Sprintf("What is \\(\\lim_{x\\to\\infty} \\frac{%dx^{%d}%+d}{%dx^{%d}%+d}\\)?", a, p, b, c, p, d),
			Answer:   strconv.Itoa(k),
			Explanation: fmt.Sprintf(
				"Both polynomials have degree %d, so the limit is the ratio of the leading coefficients: %d divided by %d is %d. The lower-order terms cannot change the end behaviour.",
				p, a, c, k),
		}
	case 1:
		p := randRange(ctx.Difficulty, 1, 3)
		c := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{x\\to\\infty} \\frac{%d}{x^{%d}}\\)?", c, p),
			Answer:      "0",
			Explanation: fmt.Sprintf("As \\(x\\) grows, \\(x^{%d}\\) grows without bound, so %d divided by it shrinks to 0.", p, c),
		}
	default:
		c := randRange(ctx.Difficulty, 2, 9)
		k := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{x\\to\\infty} \\left(%d + \\frac{%d}{x}\\right)\\)?", c, k),
			Answer:      strconv.Itoa(c),
			Explanation: fmt.Sprintf("The fraction \\(%d/x\\) tends to 0, so the expression tends to the constant %d.", k, c),
		}
	}
}

type limitAlgebraGen struct{}

func (g *limitAlgebraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 9)
	b := randRange(ctx.Difficulty, 2, 6)
	c := randRange(ctx.Difficulty, 2, 5)
	switch rand.Intn(5) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(\\lim f(x)=%d\\) and \\(\\lim g(x)=%d\\), what is \\(\\lim (f(x)+g(x))\\)?", a, b),
			Answer:      strconv.Itoa(a + b),
			Explanation: "The limit of a sum is the sum of the limits, provided both limits exist: " + strconv.Itoa(a) + " plus " + strconv.Itoa(b) + ".",
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(\\lim f(x)=%d\\) and \\(\\lim g(x)=%d\\), what is \\(\\lim f(x)g(x)\\)?", a, b),
			Answer:      strconv.Itoa(a * b),
			Explanation: "The limit of a product is the product of the limits: " + strconv.Itoa(a) + " times " + strconv.Itoa(b) + ".",
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(\\lim f(x)=%d\\), what is \\(\\lim %df(x)\\)?", a, c),
			Answer:      strconv.Itoa(c * a),
			Explanation: fmt.Sprintf("A constant multiple factors out of a limit, so the answer is %d times %d.", c, a),
		}
	case 3:
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(\\lim f(x)=%d\\) and \\(\\lim g(x)=%d\\), what is \\(\\lim \\frac{f(x)}{g(x)}\\)?", b*a, b),
			Answer:      strconv.Itoa(a),
			Explanation: fmt.Sprintf("The denominator's limit is %d, which is nonzero, so the quotient rule applies: %d divided by %d is %d.", b, b*a, b, a),
		}
	default:
		return limitChoice(
			"If \\(\\lim g(x)=0\\), is it always true that \\(\\lim \\frac{f(x)}{g(x)} = \\frac{\\lim f(x)}{\\lim g(x)}\\)?",
			"no",
			"Dividing by zero is undefined, and the quotient rule requires the denominator's limit to be nonzero.",
			"yes",
		)
	}
}

type asymptotesGen struct{}

func (g *asymptotesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(6) {
	case 0:
		a := randRange(ctx.Difficulty, 2, 9)
		c := randRange(ctx.Difficulty, 1, 5)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the vertical asymptote of \\(f(x)=\\frac{%d}{x%+d}\\)?", c, -a),
			Answer:      fmt.Sprintf("x=%d", a),
			Explanation: fmt.Sprintf("The denominator is zero at \\(x=%d\\) and the numerator is not, so the graph runs to infinity there.", a),
		}
	case 1:
		k := randRange(ctx.Difficulty, 2, 6)
		c := randRange(ctx.Difficulty, 2, 5)
		a := k * c
		b := randRange(ctx.Difficulty, 1, 9)
		d := randRange(ctx.Difficulty, 1, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the horizontal asymptote of \\(f(x)=\\frac{%dx%+d}{%dx%+d}\\)?", a, b, c, d),
			Answer:      fmt.Sprintf("y=%d", k),
			Explanation: fmt.Sprintf("Both numerator and denominator are linear, so the end behaviour is the ratio of the leading coefficients: %d over %d is %d.", a, c, k),
		}
	case 2:
		p := randRange(ctx.Difficulty, 1, 3)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the horizontal asymptote of \\(f(x)=\\frac{1}{x^{%d}}\\)?", p),
			Answer:      "y=0",
			Explanation: "The denominator grows without bound, so the values approach zero from above: the horizontal asymptote is the x-axis.",
		}
	case 3:
		return limitChoice(
			"Can the graph of a function cross its horizontal asymptote?",
			"yes",
			"An asymptote describes end behaviour, not a barrier. A function may cross the line and still approach it as x grows.",
			"no",
		)
	case 4:
		return limitChoice(
			"Does \\(f(x)=\\frac{x^{2}+1}{x-3}\\) have a horizontal asymptote?",
			"no",
			"The numerator has higher degree than the denominator, so the values grow without bound in both directions.",
			"yes",
		)
	default:
		a := randRange(ctx.Difficulty, 2, 7)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the vertical asymptote of \\(f(x)=\\frac{x+1}{x%+d}\\)?", -a),
			Answer:      fmt.Sprintf("x=%d", a),
			Explanation: fmt.Sprintf("The denominator vanishes at \\(x=%d\\), and the numerator there is %d, which is not zero, so the line is a vertical asymptote.", a, 1-a),
		}
	}
}

type lhopitalGen struct{}

func (g *lhopitalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 2, 7)
	switch rand.Intn(6) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} \\frac{\\sin(%dx)}{x}\\) using L'Hôpital's rule.", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("Both numerator and denominator tend to 0. Differentiating gives \\(%d\\cos(%dx)\\) over 1, which tends to %d.", k, k, k),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} \\frac{e^{%dx}-1}{x}\\) using L'Hôpital's rule.", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("This is 0 over 0. Differentiating gives \\(%de^{%dx}\\) over 1, which tends to %d.", k, k, k),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} \\frac{\\ln(1+%dx)}{x}\\) using L'Hôpital's rule.", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("This is 0 over 0. The derivative of the numerator is \\(%d/(1+%dx)\\), which tends to %d over 1.", k, k, k),
		}
	case 3:
		p := randRange(ctx.Difficulty, 1, 3)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to\\infty} \\frac{x^{%d}}{e^{%dx}}\\) using L'Hôpital's rule.", p, k),
			Answer:      "0",
			Explanation: "The exponential grows faster than any power, so after finitely many applications the numerator becomes a constant and the denominator keeps its exponential factor: the limit is 0.",
		}
	case 4:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} \\frac{1-\\cos(%dx)}{x^{2}}\\) using L'Hôpital's rule.", k),
			Answer:      fmt.Sprintf("%d/2", k*k),
			Explanation: fmt.Sprintf("Two applications: the first gives \\(%d\\sin(%dx)/(2x)\\) after rewriting, the second leaves \\(%d\\cos(%dx)/2\\), which is %d/2 at zero.", k, k, k*k, k, k*k),
		}
	default:
		return limitChoice(
			"Does L'Hôpital's rule apply directly to a limit of the form \\(1/0\\)?",
			"no",
			"The rule requires an indeterminate form such as 0/0 or infinity over infinity. A nonzero numerator over a vanishing denominator diverges.",
			"yes",
		)
	}
}

type squeezeGen struct{}

func (g *squeezeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		p := randRange(ctx.Difficulty, 1, 3)
		c := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} %dx^{%d}\\sin(1/x)\\).", c, p),
			Answer:      "0",
			Explanation: fmt.Sprintf("The sine factor is bounded between -1 and 1, so the expression lies between \\(-%dx^{%d}\\) and \\(%dx^{%d}\\). Both bounds tend to 0, so the limit is 0.", c, p, c, p),
		}
	case 1:
		k := randRange(ctx.Difficulty, 2, 7)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} x\\cos(%d/x)\\).", k),
			Answer:      "0",
			Explanation: fmt.Sprintf("The cosine factor stays between -1 and 1, so the expression is trapped between \\(-x\\) and \\(x\\), and both tend to 0 as \\(x\\to 0\\)."),
		}
	case 2:
		k := randRange(ctx.Difficulty, 2, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0} \\frac{\\sin(%dx)}{x}\\) by squeezing between cosines.", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("For small \\(x\\), \\(\\cos(%dx)\\) traps \\(\\sin(%dx)/x\\) between two cosines, and both tend to %d as \\(x\\to 0\\).", k, k, k),
		}
	case 3:
		return limitChoice(
			"If \\(g(x)\\le f(x)\\le h(x)\\) and \\(\\lim g=\\lim h=L\\), must \\(\\lim f\\) exist?",
			"yes",
			"The two bounds force the value: for any tolerance, f is trapped within it once both bounds are close enough to L.",
			"no",
		)
	default:
		return limitChoice(
			"If \\(\\lim g=1\\) and \\(\\lim h=2\\), can the squeeze theorem conclude that \\(\\lim f\\) exists?",
			"no",
			"The bounds must agree. A growing gap to cover tells you nothing about where f settles.",
			"yes",
		)
	}
}

type supremumGen struct{}

func (g *supremumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(6) {
	case 0:
		c := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the supremum of the set \\(\\{%d-%d/n : n \\in \\mathbb{N}\\}\\)?", c, c),
			Answer:      strconv.Itoa(c),
			Explanation: fmt.Sprintf("Every element is below %d and the elements approach it from below, so %d is the least upper bound. It is not in the set, because no finite \\(n\\) attains it.", c, c),
		}
	case 1:
		c := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the infimum of the set \\(\\{%d/n : n \\in \\mathbb{N}\\}\\)?", c),
			Answer:      "0",
			Explanation: "The elements are positive and shrink towards 0, so 0 is the greatest lower bound. It is not attained.",
		}
	case 2:
		a := randRange(ctx.Difficulty, 2, 12)
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\sup\\{x \\in \\mathbb{R} : x^{2}<%d\\}\\)? (enter as sqrt(...))", a),
			Answer:      fmt.Sprintf("sqrt(%d)", a),
			Explanation: fmt.Sprintf("The set is the open interval between \\(-\\sqrt{%d}\\) and \\(\\sqrt{%d}\\). The least upper bound is \\(\\sqrt{%d}\\), which the strict inequality excludes from the set.", a, a, a),
		}
	case 3:
		a := randRange(ctx.Difficulty, 1, 8)
		b := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is the infimum of the set \\(\\{%d+%d/n : n \\in \\mathbb{N}\\}\\)?", a, b),
			Answer:      strconv.Itoa(a),
			Explanation: fmt.Sprintf("Every element exceeds %d, and the elements approach it as \\(n\\) grows, so %d is the greatest lower bound.", a, a),
		}
	case 4:
		return limitChoice(
			"Is the supremum of a set always an element of the set?",
			"no",
			"An open interval has a supremum at its endpoint, and that endpoint is not in the set.",
			"yes",
		)
	default:
		return limitChoice(
			"If a set has a maximum, is that maximum equal to its supremum?",
			"yes",
			"A maximum is an upper bound that belongs to the set, so no smaller bound can work.",
			"no",
		)
	}
}
