package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Five concepts that were short literal tables. Each draws a parameter and
// derives its answer, and the classification shapes are multiple choice
// because a yes or no cannot be typed into a maths field.

type trigSubstitutionGen struct{}

func (g *trigSubstitutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 2, 9)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int \\frac{1}{\\sqrt{%d-x^{2}}} \\, dx\\).", a*a),
			Answer:      fmt.Sprintf("arcsin(x/%d)+C", a),
			Explanation: fmt.Sprintf("With \\(x=%d\\sin(\\theta)\\) the radical becomes \\(%d\\cos(\\theta)\\) and \\(dx\\) contributes the same factor, so the integral is \\(\\theta+C\\) with \\(\\theta=\\arcsin(x/%d)\\).", a, a, a),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int \\frac{1}{%d+x^{2}} \\, dx\\).", a*a),
			Answer:      fmt.Sprintf("(1/%d)arctan(x/%d)+C", a, a),
			Explanation: fmt.Sprintf("With \\(x=%d\\tan(\\theta)\\) the denominator becomes \\(%d\\sec^{2}(\\theta)\\), which cancels \\(dx\\), leaving \\((1/%d)\\theta+C\\).", a, a*a, a),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int \\frac{1}{x\\sqrt{x^{2}-%d}} \\, dx\\) for \\(x>%d\\).", a*a, a),
			Answer:      fmt.Sprintf("(1/%d)arcsec(x/%d)+C", a, a),
			Explanation: fmt.Sprintf("With \\(x=%d\\sec(\\theta)\\) the radical becomes \\(%d\\tan(\\theta)\\) and the integrand reduces to \\((1/%d)\\,d\\theta\\).", a, a, a),
		}
	default:
		return limitChoice(
			fmt.Sprintf("Which substitution removes the radical from \\(\\sqrt{%d-x^{2}}\\)?", a*a),
			fmt.Sprintf("x = %d\\sin(\\theta)", a),
			fmt.Sprintf("The identity \\(1-\\sin^{2}(\\theta)=\\cos^{2}(\\theta)\\) turns \\(%d-x^{2}\\) into a perfect square when \\(x=%d\\sin(\\theta)\\).", a*a, a),
			fmt.Sprintf("x = %d\\tan(\\theta)", a),
			fmt.Sprintf("x = %d\\sec(\\theta)", a),
		)
	}
}

type eulerSeqGen struct{}

func (g *eulerSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 2, 12)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{n\\to\\infty} \\left(1+\\frac{%d}{n}\\right)^{n}\\)?", k),
			Answer:      fmt.Sprintf("e^%d", k),
			Explanation: fmt.Sprintf("The standard limit \\((1+1/n)^{n}\\to e\\) generalises: \\((1+%d/n)^{n} = \\left((1+%d/n)^{n/%d}\\right)^{%d} \\to e^{%d}\\).", k, k, k, k, k),
		}
	case 1:
		d := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{n\\to\\infty} \\left(1+\\frac{1}{%dn}\\right)^{n}\\)?", d),
			Answer:      fmt.Sprintf("e^(1/%d)", d),
			Explanation: fmt.Sprintf("Write the exponent as \\(n = (%dn)/%d\\) so the expression becomes \\((1+1/(%dn))^{%dn/%d}\\), which tends to \\(e^{1/%d}\\).", d, d, d, d, d, d),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("What is \\(\\lim_{n\\to\\infty} \\left(1+\\frac{1}{n}\\right)^{%dn}\\)?", k),
			Answer:      fmt.Sprintf("e^%d", k),
			Explanation: fmt.Sprintf("The base tends to \\(e\\) and the extra exponent %d is constant relative to \\(n\\), so the limit is \\(e^{%d}\\).", k, k),
		}
	default:
		return limitChoice(
			"Is Euler's number \\(e\\) rational?",
			"no",
			"e is irrational, and transcendental: no ratio of integers equals it.",
			"yes",
		)
	}
}

type integralImproperCompareGen struct{}

func (g *integralImproperCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 2, 9)
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Evaluate \\(\\int_{1}^{\\infty} \\frac{1}{x^{%d}} \\, dx\\).", p),
			Answer:      fmt.Sprintf("1/%d", p-1),
			Explanation: fmt.Sprintf("The antiderivative is \\(x^{%d}/(%d)\\) with a negative exponent, and its value at infinity is zero, leaving \\(1/%d\\).", 1-p, 1-p, p-1),
		}
	case 1:
		return limitChoice(
			fmt.Sprintf("Does \\(\\int_{1}^{\\infty} \\frac{1}{x^{%d}} \\, dx\\) converge?", p),
			"yes",
			fmt.Sprintf("The p-integral converges exactly when p exceeds 1, and here p is %d.", p),
			"no",
		)
	default:
		return limitChoice(
			fmt.Sprintf("Does \\(\\int_{1}^{\\infty} \\frac{1}{x^{1/%d}} \\, dx\\) converge?", p),
			"no",
			fmt.Sprintf("The exponent 1/%d is at most 1, so the integral diverges like the harmonic case. Only exponents above 1 converge.", p),
			"yes",
		)
	}
}

type curvatureGen struct{}

func (g *curvatureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 2, 12)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find the curvature of \\(y = %dx^{2}\\) at \\(x = 0\\).", a),
			Answer:      strconv.Itoa(2 * a),
			Explanation: fmt.Sprintf("Curvature is \\(|y''|/(1+(y')^{2})^{3/2}\\). Here \\(y'=%dx\\) vanishes at zero, so only \\(y''=%d\\) survives.", 2*a, 2*a),
			Facts:       map[string]string{"a": strconv.Itoa(a), "curvature": strconv.Itoa(2 * a)},
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find the curvature of \\(y = %dx^{3}\\) at \\(x = 0\\).", a),
			Answer:      "0",
			Explanation: fmt.Sprintf("Here \\(y'' = %dx\\), which is zero at the origin. The curve is straight there, so the curvature is 0.", 6*a),
			Facts:       map[string]string{"a": strconv.Itoa(a), "curvature": "0"},
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find the curvature of the line \\(y = %dx + %d\\).", a, a+1),
			Answer:      "0",
			Explanation: "A line has zero second derivative everywhere, so its curvature is 0 at every point.",
			Facts:       map[string]string{"a": strconv.Itoa(a), "curvature": "0"},
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find the curvature of \\(y = -%dx^{2}\\) at \\(x = 0\\) (enter a non-negative number).", a),
			Answer:      strconv.Itoa(2 * a),
			Explanation: fmt.Sprintf("Curvature uses the absolute value of the second derivative, so the downward bowl has \\(|y''|=%d\\) at the vertex and the same curvature as the upward one.", 2*a),
			Facts:       map[string]string{"a": strconv.Itoa(a), "curvature": strconv.Itoa(2 * a)},
		}
	}
}

type parametricDerivativeGen struct{}

func (g *parametricDerivativeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(4) {
	case 0:
		k := randRange(ctx.Difficulty, 2, 8)
		t := 2 * k
		return generator.Problem{
			Question:    fmt.Sprintf("For the curve \\(x = t^{2},\\ y = t^{3}\\), find \\(\\frac{dy}{dx}\\) at \\(t = %d\\).", t),
			Answer:      strconv.Itoa(3 * k),
			Explanation: fmt.Sprintf("Divide the derivatives: \\(dy/dx = 3t^{2}/(2t) = 3t/2\\). At \\(t=%d\\) that is %d.", t, 3*k),
			Facts:       map[string]string{"t": strconv.Itoa(t), "slope": strconv.Itoa(3 * k)},
		}
	case 1:
		t := randRange(ctx.Difficulty, 2, 12)
		return generator.Problem{
			Question:    fmt.Sprintf("For the curve \\(x = t^{2},\\ y = t^{4}\\), find \\(\\frac{dy}{dx}\\) at \\(t = %d\\).", t),
			Answer:      strconv.Itoa(2 * t * t),
			Explanation: fmt.Sprintf("The derivatives are \\(dx/dt = 2t\\) and \\(dy/dt = 4t^{3}\\), so \\(dy/dx = 2t^{2}\\). At \\(t=%d\\) that is %d.", t, 2*t*t),
			Facts:       map[string]string{"t": strconv.Itoa(t), "slope": strconv.Itoa(2 * t * t)},
		}
	case 2:
		m := randRange(ctx.Difficulty, 1, 6)
		t := 3 * m
		return generator.Problem{
			Question:    fmt.Sprintf("For the curve \\(x = t^{3},\\ y = t^{5}\\), find \\(\\frac{dy}{dx}\\) at \\(t = %d\\).", t),
			Answer:      strconv.Itoa(15 * m * m),
			Explanation: fmt.Sprintf("\\(dy/dx = 5t^{4}/(3t^{2}) = (5/3)t^{2}\\). At \\(t=%d\\) that is \\((5/3)\\cdot%d = %d\\).", t, t*t, 15*m*m),
			Facts:       map[string]string{"t": strconv.Itoa(t), "slope": strconv.Itoa(15 * m * m)},
		}
	default:
		t := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("For the curve \\(x = t^{2},\\ y = t^{6}\\), find \\(\\frac{dy}{dx}\\) at \\(t = %d\\).", t),
			Answer:      strconv.Itoa(3 * intPow(t, 4)),
			Explanation: fmt.Sprintf("\\(dy/dx = 6t^{5}/(2t) = 3t^{4}\\). At \\(t=%d\\) that is %d.", t, 3*intPow(t, 4)),
			Facts:       map[string]string{"t": strconv.Itoa(t), "slope": strconv.Itoa(3 * intPow(t, 4))},
		}
	}
}
