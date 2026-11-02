package calculus

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// calc1BridgeGen serves the single-variable calculus concepts added by the
// content expansion. The kind selects the concept and KPOrdinal selects the
// lesson section.
type calc1BridgeGen struct{ kind string }

func (g *calc1BridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch g.kind {
	case "integral_properties":
		return integralPropertiesProblem(o)
	case "deriv_definition":
		return derivDefinitionProblem(o)
	case "initial_value":
		return initialValueProblem(o)
	case "deriv_rate":
		return derivativeRateProblem(o)
	case "graph_reading":
		return graphReadingProblem(o)
	case "integral_symmetry":
		return integralSymmetryProblem(o)
	case "limit_trig":
		return trigLimitProblem(o)
	}
	return calcBridgeProblem(o, "Evaluate the expression. (enter a number)", "0", "The expression evaluates to zero.", nil)
}

func integralPropertiesProblem(o int) generator.Problem {
	a := 1 + rand.Intn(4)
	b := a + 2 + rand.Intn(5)
	switch o {
	case 1:
		c := 2 + rand.Intn(4)
		d := 2 + rand.Intn(4)
		fv := 2 + rand.Intn(8)
		gv := 2 + rand.Intn(8)
		ans := c*fv + d*gv
		return calcBridgeProblem(o,
			fmt.Sprintf("For a function f, the integral from %d to %d is %d. For a function g over the same interval the integral is %d. What is the integral of %df + %dg over that interval? (enter a number)", a, b, fv, gv, c, d),
			strconv.Itoa(ans),
			fmt.Sprintf("Linearity gives %d times %d plus %d times %d, which is %d.", c, fv, d, gv, ans),
			map[string]string{"kind": "linearity"})
	case 2:
		c := b + 2 + rand.Intn(4)
		part := 3 + rand.Intn(15)
		whole := part + 3 + rand.Intn(15)
		return calcBridgeProblem(o,
			fmt.Sprintf("The integral of f from %d to %d is %d and the integral from %d to %d is %d. What is the integral from %d to %d? (enter a number)", a, b, part, b, c, whole, a, c),
			strconv.Itoa(part+whole),
			fmt.Sprintf("Interval additivity gives %d + %d = %d.", part, whole, part+whole),
			map[string]string{"kind": "additivity"})
	default:
		even := 2 * (1 + rand.Intn(9))
		if rand.Intn(2) == 0 {
			return calcBridgeProblem(o,
				fmt.Sprintf("A function f is even and its integral from -%d to %d is %d. What is its integral from 0 to %d? (enter a number)", a, a, even, a),
				strconv.Itoa(even/2),
				fmt.Sprintf("The right half carries half of the symmetric integral, so the answer is %d.", even/2),
				map[string]string{"kind": "even"})
		}
		c := 2 + rand.Intn(5)
		return calcBridgeProblem(o,
			fmt.Sprintf("A function f is even with integral from 0 to %d equal to %d. What is the integral of %df from -%d to %d? (enter a number)", a, even/2, c, a, a),
			strconv.Itoa(c*even),
			fmt.Sprintf("The symmetric integral of f is %d, and the constant multiplies it by %d to give %d.", even, c, c*even),
			map[string]string{"kind": "even-scaled"})
	}
}

func derivDefinitionProblem(o int) generator.Problem {
	a := 1 + rand.Intn(6)
	switch o {
	case 1:
		return calcBridgeProblem(o,
			fmt.Sprintf("Use the limit definition to find the derivative of f(x)=x^2 at x=%d. (enter a number)", a),
			strconv.Itoa(2*a),
			fmt.Sprintf("The difference quotient simplifies to 2x+h, and as h tends to zero it becomes 2x. At x=%d this is %d.", a, 2*a),
			map[string]string{"x": strconv.Itoa(a)})
	case 2:
		return calcBridgeProblem(o,
			fmt.Sprintf("Use the limit definition to find the derivative of f(x)=x^3 at x=%d. (enter a number)", a),
			strconv.Itoa(3*a*a),
			fmt.Sprintf("The difference quotient simplifies to 3x^2+3xh+h^2, whose limit is 3x^2. At x=%d this is %d.", a, 3*a*a),
			map[string]string{"x": strconv.Itoa(a)})
	default:
		k := 2 + rand.Intn(4)
		return calcBridgeProblem(o,
			fmt.Sprintf("Use the limit definition to find the derivative of f(x)=%dx^2 at x=%d. (enter a number)", k, a),
			strconv.Itoa(2*k*a),
			fmt.Sprintf("The derivative is %dx, and at x=%d it equals %d.", 2*k, a, 2*k*a),
			map[string]string{"x": strconv.Itoa(a), "k": strconv.Itoa(k)})
	}
}

func initialValueProblem(o int) generator.Problem {
	c := rand.Intn(7) - 3
	switch o {
	case 1:
		k := 2 + rand.Intn(6)
		n := 1 + rand.Intn(6)
		ans := k*n + c
		return calcBridgeProblem(o,
			fmt.Sprintf("A function satisfies f'(x)=%d and f(0)=%d. What is f(%d)? (enter a number)", k, c, n),
			strconv.Itoa(ans),
			fmt.Sprintf("The antiderivative is %dx plus the constant. Using f(0)=%d gives f(x)=%dx%+d, so f(%d)=%d.", k, c, k, c, n, ans),
			map[string]string{"kind": "constant-rate"})
	case 2:
		n := 1 + rand.Intn(4)
		q := 1 + rand.Intn(4)
		a := q * (n + 1)
		ans := q + c
		return calcBridgeProblem(o,
			fmt.Sprintf("A function satisfies f'(x)=%dx^%d and f(0)=%d. What is f(1)? (enter a number)", a, n, c),
			strconv.Itoa(ans),
			fmt.Sprintf("The antiderivative is %dx^%d plus the constant. With f(0)=%d, f(1)=%d.", q, n+1, c, ans),
			map[string]string{"kind": "power"})
	default:
		k := 2 * (1 + rand.Intn(3))
		p := rand.Intn(7) - 3
		t := 1 + rand.Intn(4)
		ans := k*t*t/2 + p*t + c
		return calcBridgeProblem(o,
			fmt.Sprintf("A function satisfies f''(x)=%d, f'(0)=%d and f(0)=%d. What is f(%d)? (enter a number)", k, p, c, t),
			strconv.Itoa(ans),
			fmt.Sprintf("Integrating twice gives f(x)=%dx^2/2%+dx%+d. At x=%d the value is %d.", k, p, c, t, ans),
			map[string]string{"kind": "second-order"})
	}
}

func derivativeRateProblem(o int) generator.Problem {
	a := 1 + rand.Intn(5)
	b := a + 1 + rand.Intn(5)
	switch o {
	case 1:
		return calcBridgeProblem(o,
			fmt.Sprintf("What is the average rate of change of f(x)=x^2 on the interval [%d,%d]? (enter a number)", a, b),
			strconv.Itoa(a+b),
			fmt.Sprintf("The average rate is (b^2-a^2)/(b-a)=a+b, which is %d.", a+b),
			map[string]string{"kind": "average"})
	case 2:
		return calcBridgeProblem(o,
			fmt.Sprintf("What is the instantaneous rate of change of f(x)=x^2 at x=%d? (enter a number)", a),
			strconv.Itoa(2*a),
			fmt.Sprintf("The derivative is 2x, so at x=%d the rate is %d.", a, 2*a),
			map[string]string{"kind": "instantaneous"})
	default:
		k := 1 + rand.Intn(4)
		v := 2 + rand.Intn(6)
		t := 1 + rand.Intn(4)
		ans := 2*k*t + v
		return calcBridgeProblem(o,
			fmt.Sprintf("A particle has position s(t)=%dt^2+%dt. What is its velocity at t=%d? (enter a number)", k, v, t),
			strconv.Itoa(ans),
			fmt.Sprintf("Velocity is the derivative %dt+%d. At t=%d it is %d.", 2*k, v, t, ans),
			map[string]string{"kind": "velocity"})
	}
}

func graphReadingProblem(o int) generator.Problem {
	switch o {
	case 1:
		b := rand.Intn(9) - 4
		c := rand.Intn(9) - 4
		return calcBridgeProblem(o,
			fmt.Sprintf("The graph of f(x)=x^2%+dx%+d has a tangent at x=0. What is the slope of that tangent? (enter a number)", b, c),
			strconv.Itoa(b),
			fmt.Sprintf("The derivative is 2x%+d, and at x=0 the slope is %d.", b, b),
			map[string]string{"kind": "tangent-at-zero"})
	case 2:
		k := 1 + rand.Intn(5)
		return calcBridgeProblem(o,
			fmt.Sprintf("The graph of f' crosses the x-axis at %d distinct points and changes sign at every one. How many turning points does the graph of f have? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("Each sign change of f' produces one turning point of f, so there are %d.", k),
			map[string]string{"kind": "turning-points"})
	default:
		n := 5
		signs := make([]string, n)
		for i := range signs {
			if rand.Intn(2) == 0 {
				signs[i] = "+"
			} else {
				signs[i] = "-"
			}
		}
		maxima := 0
		for i := 1; i < n; i++ {
			if signs[i-1] == "+" && signs[i] == "-" {
				maxima++
			}
		}
		return calcBridgeProblem(o,
			fmt.Sprintf("A derivative changes sign across five consecutive intervals in the order %s. How many local maxima does the function have? (enter a number)", strings.Join(signs, " ")),
			strconv.Itoa(maxima),
			fmt.Sprintf("A local maximum needs a change from positive to negative. The given order has %d such changes.", maxima),
			map[string]string{"kind": "sign-chart"})
	}
}

func integralSymmetryProblem(o int) generator.Problem {
	a := 1 + rand.Intn(5)
	switch o {
	case 1:
		k := 1 + rand.Intn(9)
		return calcBridgeProblem(o,
			fmt.Sprintf("An even function f has integral from -%d to %d equal to %d. What is its integral from 0 to %d? (enter a number)", a, a, 2*k, a),
			strconv.Itoa(k),
			fmt.Sprintf("For an even function the symmetric integral is twice the half integral, so the answer is %d.", k),
			map[string]string{"kind": "even"})
	case 2:
		names := []string{"x^3", "x^5", "sin(x)", "tan(x)", "x cos(x)"}
		name := names[rand.Intn(len(names))]
		return calcBridgeProblem(o,
			fmt.Sprintf("What is the integral of the odd function %s from -%d to %d? (enter a number)", name, a, a),
			"0",
			"The signed areas on the two sides cancel because the function is odd.",
			map[string]string{"kind": "odd"})
	default:
		c := 2 + rand.Intn(5)
		e := 1 + rand.Intn(8)
		return calcBridgeProblem(o,
			fmt.Sprintf("An even function f has integral from 0 to %d equal to %d. What is the integral of %df from -%d to %d? (enter a number)", a, e, c, a, a),
			strconv.Itoa(2*c*e),
			fmt.Sprintf("The symmetric integral of f is %d, and multiplying by %d gives %d.", 2*e, c, 2*c*e),
			map[string]string{"kind": "even-scaled"})
	}
}

func trigLimitProblem(o int) generator.Problem {
	k := 2 + rand.Intn(8)
	switch o {
	case 1:
		return calcBridgeProblem(o,
			fmt.Sprintf("Evaluate the limit of sin(%dx)/x as x approaches 0. (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("Sine of %dx divided by x tends to %d.", k, k),
			map[string]string{"kind": "sine"})
	case 2:
		m := 2 + rand.Intn(4)
		ans := calcBridgeDecimal(float64(k) / float64(m))
		return calcBridgeProblem(o,
			fmt.Sprintf("Evaluate the limit of sin(%dx)/(%dx) as x approaches 0. (enter a decimal)", k, m),
			ans,
			fmt.Sprintf("The limit is the ratio of the coefficients, %d/%d = %s.", k, m, ans),
			map[string]string{"kind": "scaled-sine"})
	default:
		if rand.Intn(2) == 0 {
			return calcBridgeProblem(o,
				fmt.Sprintf("Evaluate the limit of tan(%dx)/x as x approaches 0. (enter a number)", k),
				strconv.Itoa(k),
				fmt.Sprintf("Tangent of %dx divided by x has the same limit as sine, namely %d.", k, k),
				map[string]string{"kind": "tangent"})
		}
		ans := calcBridgeDecimal(float64(k*k) / 2)
		return calcBridgeProblem(o,
			fmt.Sprintf("Evaluate the limit of (1-cos(%dx))/x^2 as x approaches 0. (enter a decimal)", k),
			ans,
			fmt.Sprintf("The limit is %d squared divided by two, which is %s.", k, ans),
			map[string]string{"kind": "cosine"})
	}
}

func (g *calc1BridgeGen) String() string { return "calc1BridgeGen(" + g.kind + ")" }
