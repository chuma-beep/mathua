package calculus

// Wave 9 remaining calculus: eleven Calculus I concepts and six Calculus II
// concepts. Every family is exact, integer-valued and publishes its operands.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

var w9rPureNumber = regexp.MustCompile(`\\\((-?\d+)\\\)`)
var w9rUnitCoefficient = regexp.MustCompile(`(^|[^0-9])1 ?([xyzuvtr])([^A-Za-z]|$)`)
var w9rUnitParen = regexp.MustCompile(`(^|[^0-9])1\(`)
var w9rSpaces = regexp.MustCompile(`[ \t]+`)

// w9rClean removes rendering defects that a learner should not have to read.
func w9rClean(s string) string {
	s = w9rPureNumber.ReplaceAllString(s, "$1")
	s = strings.ReplaceAll(s, "+-", "-")
	s = strings.ReplaceAll(s, "--", "+")
	s = w9rSpaces.ReplaceAllString(s, " ")
	s = w9rUnitCoefficient.ReplaceAllString(s, "${1}${2}${3}")
	s = w9rUnitParen.ReplaceAllString(s, "${1}(")
	return s
}

func registerWave9Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.limit.continuity_composition": "continuity_composition",
		"calc.limit.continuity_inverse":     "continuity_inverse",
		"calc.limit.continuity_extension":   "continuity_extension",
		"calc.deriv.absolute":               "deriv_absolute",
		"calc.deriv.floor":                  "deriv_floor",
		"calc.deriv.sign":                   "deriv_sign",
		"calc.deriv.related_rates_shadow":   "shadow",
		"calc.deriv.related_rates_angle":    "angle",
		"calc.deriv.optimization_distance":  "optimization_distance",
		"calc.deriv.optimization_cylinder":  "optimization_cylinder",
		"calc.integral.floor":               "integral_floor",
		"calc.series.abel_theorem":          "abel_theorem",
		"calc.series.asymptotic":            "asymptotic_series",
		"calc.series.euler_maclaurin":       "euler_maclaurin",
		"calc.series.borel":                 "borel",
		"calc.series.taubian":               "taubian",
		"calc.series.lambert":               "lambert",
	} {
		reg.Register(id, &w9rCalcGen{kind: kind})
	}
}

type w9rCalcGen struct{ kind string }

func (g *w9rCalcGen) Generate(ctx generator.GeneratorContext) (prob generator.Problem) {
	defer func() {
		prob.Question = w9rClean(prob.Question)
		prob.Explanation = w9rClean(prob.Explanation)
	}()
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "continuity_composition":
		return w9rContinuityComposition(kp)
	case "continuity_inverse":
		return w9rContinuityInverse(kp)
	case "continuity_extension":
		return w9rContinuityExtension(kp)
	case "deriv_absolute":
		return w9rDerivAbsolute(kp)
	case "deriv_floor":
		return w9rDerivFloor(kp)
	case "deriv_sign":
		return w9rDerivSign(kp)
	case "shadow":
		return w9rShadow(kp)
	case "angle":
		return w9rAngle(kp)
	case "optimization_distance":
		return w9rOptimizationDistance(kp)
	case "optimization_cylinder":
		return w9rOptimizationCylinder(kp)
	case "integral_floor":
		return w9rIntegralFloor(kp)
	case "abel_theorem":
		return w9rAbel(kp)
	case "asymptotic_series":
		return w9rAsymptoticSeries(kp)
	case "euler_maclaurin":
		return w9rEulerMaclaurin(kp)
	case "borel":
		return w9rBorel(kp)
	case "taubian":
		return w9rTaubian(kp)
	case "lambert":
		return w9rLambert(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w9rC(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w9rCI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w9rCR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w9rCIPow(b, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= b
	}
	return r
}

func w9rContinuityComposition(kp int) generator.Problem {
	a := w9rCR(-4, 4)
	x0 := w9rCR(-3, 3)
	gx := x0 + a
	f := map[string]string{}
	w9rCI(f, "a", a)
	w9rCI(f, "x0", x0)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(g(x)=x+%d\\) and \\(f(x)=x^2\\), what is \\(\\lim_{x\\to %d} f(g(x))\\)? (enter a number)", a, x0), strconv.Itoa(gx*gx),
			fmt.Sprintf("The inner limit is \\(g(%d)=%d\\), and \\(f(%d)=%d\\).", x0, gx, gx, gx*gx), f)
	case 2:
		inner := x0*x0 + a
		return w9rC(kp, fmt.Sprintf("For \\(g(x)=x+%d\\) and \\(f(x)=x^2\\), what is \\(g(f(%d))\\)? (enter a number)", a, x0), strconv.Itoa(inner),
			fmt.Sprintf("\\(f(%d)=%d\\), so \\(g(%d)=%d+%d=%d\\).", x0, x0*x0, x0*x0, x0*x0, a, inner), f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\(g(x)=x+%d\\) and \\(f(x)=x^2\\), what is \\(g(%d)\\)? (enter a number)", a, x0), strconv.Itoa(gx),
			fmt.Sprintf("\\(g(%d)=%d\\).", x0, gx), f)
	}
}

func w9rContinuityInverse(kp int) generator.Problem {
	r := w9rCR(-4, 4)
	a := w9rCR(-4, 4)
	y := r*r*r + a
	f := map[string]string{}
	w9rCI(f, "a", a)
	w9rCI(f, "y", y)
	w9rCI(f, "r", r)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=x^3+%d\\), what is \\(f^{-1}(%d)\\)? (enter a number)", a, y), strconv.Itoa(r),
			fmt.Sprintf("Solve \\(x^3+%d=%d\\), so \\(x^3=%d\\) and \\(x=%d\\).", a, y, r*r*r, r), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=x^3+%d\\) and \\(y=%d\\), what is \\(f^{-1}(y)+1\\)? (enter a number)", a, y), strconv.Itoa(r+1),
			fmt.Sprintf("The inverse value is \\(%d\\), so adding 1 gives %d.", r, r+1), f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=x^3+%d\\), what is \\(f'(%d)\\)? (enter a number)", a, r), strconv.Itoa(3*r*r),
			fmt.Sprintf("\\(f'(x)=3x^2\\), so \\(f'(%d)=3(%d)^2=%d\\).", r, r, 3*r*r), f)
	}
}

func w9rContinuityExtension(kp int) generator.Problem {
	a := w9rCR(1, 8)
	f := map[string]string{}
	w9rCI(f, "a", a)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=(x^2-%d)/(x-%d)\\), what is \\(\\lim_{x\\to %d}f(x)\\)? (enter a number)", a*a, a, a), strconv.Itoa(2*a),
			fmt.Sprintf("Factor the numerator as \\((x-%d)(x+%d)\\), so the limit is \\(2\\cdot%d=%d\\).", a, a, a, 2*a), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("What value at \\(x=%d\\) makes \\(f(x)=(x^2-%d)/(x-%d)\\) continuous? (enter a number)", a, a*a, a), strconv.Itoa(2*a),
			fmt.Sprintf("The continuous extension value is the limit \\(2\\cdot%d=%d\\).", a, 2*a), f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=(x^2-%d)/(x-%d)\\), what is \\(f(%d)\\) after the removable discontinuity is filled?", a*a, a, a+1), strconv.Itoa(2*a+1),
			fmt.Sprintf("After cancelling the common factor the function is \\(x+%d\\), so at \\(x=%d\\) the value is %d.", a, a+1, 2*a+1), f)
	}
}

func w9rDerivAbsolute(kp int) generator.Problem {
	a := w9rCR(-5, 5)
	f := map[string]string{}
	w9rCI(f, "a", a)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=|x-%d|\\), what is the left derivative at \\(x=%d\\)? (enter a number)", a, a), "-1",
			"To the left of the corner the absolute value is the negative of its linear argument, so its derivative is -1.", f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=|x-%d|\\), what is the right derivative at \\(x=%d\\)? (enter a number)", a, a), "1",
			"To the right of the corner the absolute value is its linear argument, so its derivative is 1.", f)
	default:
		return w9rC(kp, fmt.Sprintf("Is \\(f(x)=|x-%d|\\) differentiable at \\(x=%d\\)? Enter 1 for yes and 0 for no.", a, a), "0",
			"The left and right derivatives differ, so the function is not differentiable at the corner.", f)
	}
}

func w9rDerivFloor(kp int) generator.Problem {
	a := w9rCR(-5, 5)
	f := map[string]string{}
	w9rCI(f, "a", a)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=\\lfloor x\\rfloor\\), what is \\(f'(x)\\) on the interval \\((%d,%d)\\)? (enter a number)", a, a+1), "0",
			"The floor function is constant between consecutive integers.", f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=\\lfloor x\\rfloor\\), what is \\(f(%d)\\)? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("Because %d is an integer, \\(\\lfloor %d\\rfloor=%d\\).", a, a, a), f)
	default:
		return w9rC(kp, fmt.Sprintf("Is \\(f(x)=\\lfloor x\\rfloor\\) differentiable at \\(x=%d\\)? Enter 1 for yes and 0 for no.", a), "0",
			"The floor function has a jump discontinuity at every integer, so it is not differentiable there.", f)
	}
}

func w9rDerivSign(kp int) generator.Problem {
	a := w9rCR(-5, 5)
	f := map[string]string{}
	w9rCI(f, "a", a)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=\\operatorname{sgn}(x-%d)\\), what is \\(f'(x)\\) for \\(x>%d\\)? (enter a number)", a, a), "0",
			"The sign function is constant on each side of the jump.", f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=\\operatorname{sgn}(x-%d)\\), what is \\(f(%d)\\)? (enter a number)", a, a+1), "1",
			"To the right of the jump the sign is 1.", f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\(g(x)=(x-%d)\\operatorname{sgn}(x-%d)=|x-%d|\\), what is \\(g'(x)\\) for \\(x>%d\\)? (enter a number)", a, a, a, a), "1",
			"On the right side the sign is 1, so the derivative of the product is 1.", f)
	}
}

func w9rShadow(kp int) generator.Problem {
	h := w9rCR(2, 4)
	H := 3 * h
	x := 2 * w9rCR(2, 8)
	rate := 2 * w9rCR(1, 5)
	f := map[string]string{}
	for key, v := range map[string]int{"h": h, "H": H, "x": x, "rate": rate} {
		w9rCI(f, key, v)
	}
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("A person %d units tall walks away from a %d-unit pole at %d units per second. How fast is the shadow length increasing? (enter a number)", h, H, rate), strconv.Itoa(h*rate/(H-h)),
			fmt.Sprintf("The shadow rate is \\(h\\,dx/dt/(H-h)=%d\\cdot%d/%d=%d\\).", h, rate, H-h, h*rate/(H-h)), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("With a person %d units tall and a pole %d units tall, the person is %d units from the pole. What is the shadow length? (enter a number)", h, H, x), strconv.Itoa(h*x/(H-h)),
			fmt.Sprintf("Similar triangles give \\(s=hx/(H-h)=%d\\cdot%d/%d=%d\\).", h, x, H-h, h*x/(H-h)), f)
	default:
		return w9rC(kp, fmt.Sprintf("For a person %d units tall and a pole %d units tall, with the person %d units from the pole, what is the distance from the person to the shadow tip? (enter a number)", h, H, x), strconv.Itoa(3*x/2),
			fmt.Sprintf("The shadow length is %d, so the tip is \\(%d+%d=%d\\) units from the person.", x/2, x, x/2, 3*x/2), f)
	}
}

func w9rAngle(kp int) generator.Problem {
	h := w9rCR(2, 8)
	x := w9rCR(3, 10)
	rate := w9rCR(1, 5)
	f := map[string]string{}
	w9rCI(f, "h", h)
	w9rCI(f, "x", x)
	w9rCI(f, "rate", rate)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("An observer is %d units from a tower %d units tall. What is \\(x+h\\)? (enter a number)", x, h), strconv.Itoa(x+h),
			fmt.Sprintf("\\(x+h=%d+%d=%d\\).", x, h, x+h), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("An observer is %d units from a tower %d units tall. What is \\(x^2+h^2\\)? (enter a number)", x, h), strconv.Itoa(x*x+h*h),
			fmt.Sprintf("\\(%d^2+%d^2=%d\\).", x, h, x*x+h*h), f)
	default:
		return w9rC(kp, fmt.Sprintf("For the same geometry with \\(dx/dt=%d\\), what is the numerator of the angle rate \\(|d\\theta/dt|\\) before dividing by \\(x^2+h^2\\)? (enter a number)", rate), strconv.Itoa(h*rate),
			fmt.Sprintf("\\(|d\\theta/dt|=h\\,dx/dt/(x^2+h^2)\\), so the numerator is \\(%d\\cdot%d=%d\\).", h, rate, h*rate), f)
	}
}

func w9rOptimizationDistance(kp int) generator.Problem {
	a := 2 * w9rCR(1, 8)
	f := map[string]string{}
	w9rCI(f, "a", a)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("What is the x-coordinate of the point on \\(y=x\\) closest to \\((%d,0)\\)? (enter a number)", a), strconv.Itoa(a/2),
			fmt.Sprintf("Minimise \\((x-%d)^2+x^2\\). The derivative gives \\(x=%d/2=%d\\).", a, a, a/2), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\((%d,0)\\), what is the minimum squared distance to the line \\(y=x\\)? (enter a number)", a), strconv.Itoa(a*a/2),
			fmt.Sprintf("At \\(x=%d/2\\), the squared distance is \\((a/2)^2+(a/2)^2=%d\\).", a, a*a/2), f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\((%d,0)\\) and the line \\(y=x\\), what is \\(x+y\\) at the closest point? (enter a number)", a), strconv.Itoa(a),
			fmt.Sprintf("The closest point is \\((a/2,a/2)\\), so \\(x+y=a=%d\\).", a), f)
	}
}

func w9rOptimizationCylinder(kp int) generator.Problem {
	r := w9rCR(1, 6)
	f := map[string]string{}
	w9rCI(f, "r", r)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("A cylinder with volume \\(%d\\pi\\) is designed to minimise surface area. What is its radius? (enter a number)", 2*r*r*r), strconv.Itoa(r),
			fmt.Sprintf("The optimal condition is \\(h=2r\\) with \\(V=2\\pi r^3=%d\\pi\\), so \\(r=%d\\).", 2*r*r*r, r), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For the optimal cylinder with radius %d, what is its height? (enter a number)", r), strconv.Itoa(2*r),
			fmt.Sprintf("Surface area is minimised when \\(h=2r=2\\cdot%d=%d\\).", r, 2*r), f)
	default:
		return w9rC(kp, fmt.Sprintf("For the optimal cylinder with radius %d, what is \\(V/\\pi\\)? (enter a number)", r), strconv.Itoa(2*r*r*r),
			fmt.Sprintf("\\(V/\\pi=2r^3=2\\cdot%d=%d\\).", r*r*r, 2*r*r*r), f)
	}
}

func w9rIntegralFloor(kp int) generator.Problem {
	n := w9rCR(2, 10)
	f := map[string]string{}
	w9rCI(f, "n", n)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("What is \\(\\int_0^{%d}\\lfloor x\\rfloor\\,dx\\)? (enter a number)", n), strconv.Itoa(n*(n-1)/2),
			fmt.Sprintf("The integral is \\(0+1+\\cdots+(%d)=%d\\).", n-1, n*(n-1)/2), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("What is \\(\\sum_{k=0}^{%d}\\lfloor k\\rfloor\\)? (enter a number)", n-1), strconv.Itoa(n*(n-1)/2),
			fmt.Sprintf("The sum is \\(0+1+\\cdots+(%d)=%d\\).", n-1, n*(n-1)/2), f)
	default:
		return w9rC(kp, fmt.Sprintf("How many jump discontinuities does \\(\\lfloor x\\rfloor\\) have on \\((0,%d)\\)? (enter a number)", n), strconv.Itoa(n-1),
			fmt.Sprintf("The floor jumps at \\(1,2,\\ldots,%d\\), so there are %d jumps.", n-1, n-1), f)
	}
}

func w9rAbel(kp int) generator.Problem {
	k := w9rCR(2, 15)
	f := map[string]string{}
	w9rCI(f, "k", k)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(\\sum (1/%d)^n x^n\\), what is the denominator of the Abel boundary sum at \\(x=%d\\)? (enter a number)", k, k), strconv.Itoa(k-1),
			fmt.Sprintf("The boundary sum is \\(1/(1-1/%d)=%d/%d\\), whose denominator is %d.", k, k, k-1, k-1), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(\\sum (1/%d)^n x^n\\), what is the numerator of the Abel boundary sum at \\(x=%d\\)? (enter a number)", k, k), strconv.Itoa(k),
			fmt.Sprintf("The boundary sum is \\(%d/%d\\), so the numerator is %d.", k, k-1, k), f)
	default:
		return w9rC(kp, fmt.Sprintf("For the power series \\(\\sum (1/%d)^n x^n\\), what is its sum at \\(x=%d/2\\)? (enter a number)", k, k), "2",
			fmt.Sprintf("The ratio is \\(1/2\\), so the sum is \\(1/(1-1/2)=2\\)."), f)
	}
}

func w9rAsymptoticSeries(kp int) generator.Problem {
	n := w9rCR(2, 8)
	c := w9rCR(1, 5)
	f := map[string]string{}
	w9rCI(f, "n", n)
	w9rCI(f, "c", c)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("An asymptotic series begins \\(%d-%d!/x+\\cdots\\). What is its leading coefficient? (enter a number)", c, c), strconv.Itoa(c),
			fmt.Sprintf("As \\(x\\to\\infty\\), the constant term %d dominates.", c), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("An asymptotic series is truncated after %d factorial terms. Including the leading constant, how many terms are kept? (enter a number)", n), strconv.Itoa(n+1),
			fmt.Sprintf("There are %d factorial terms plus the constant term, so %d terms.", n, n+1), f)
	default:
		return w9rC(kp, fmt.Sprintf("If the asymptotic series is truncated after %d factorial terms, what is the order of the first omitted term in powers of \\(1/x\\)? (enter a number)", n), strconv.Itoa(n+1),
			fmt.Sprintf("The first omitted term is proportional to \\(x^{-(%d)}\\).", n+1), f)
	}
}

func w9rEulerMaclaurin(kp int) generator.Problem {
	n := w9rCR(2, 12)
	f := map[string]string{}
	w9rCI(f, "n", n)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=2x\\) and \\(n=%d\\), what is the Euler-Maclaurin endpoint correction \\((f(1)+f(n))/2\\)? (enter a number)", n), strconv.Itoa(n+1),
			fmt.Sprintf("The correction is \\((2+2n)/2=n+1=%d\\).", n+1), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=2x\\) and \\(n=%d\\), what is \\(\\sum_{k=1}^n f(k)-\\int_1^n f(x)\\,dx\\)? (enter a number)", n), strconv.Itoa(n+1),
			fmt.Sprintf("The sum is \\(n(n+1)\\) and the integral is \\(n^2-1\\); the difference is \\(n+1=%d\\).", n+1), f)
	default:
		return w9rC(kp, fmt.Sprintf("For \\(f(x)=2x\\) and \\(n=%d\\), what is \\(\\int_1^n 2x\\,dx\\)? (enter a number)", n), strconv.Itoa(n*n-1),
			fmt.Sprintf("\\(\\int_1^n 2x\\,dx=n^2-1=%d\\).", n*n-1), f)
	}
}

func w9rBorel(kp int) generator.Problem {
	n := w9rCR(2, 6)
	k := w9rCR(2, 12)
	f := map[string]string{}
	w9rCI(f, "n", n)
	w9rCI(f, "k", k)
	switch kp {
	case 1:
		fact := 1
		for i := 2; i <= n; i++ {
			fact *= i
		}
		return w9rC(kp, fmt.Sprintf("What is the coefficient \\(n!\\) for \\(n=%d\\) in the Borel transform of \\(\\sum n!x^n\\)? (enter a number)", n), strconv.Itoa(fact),
			fmt.Sprintf("\\(%d!=%d\\).", n, fact), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(\\sum n!x^n\\), the Borel sum at \\(x=1/%d\\) has denominator \\(%d-1\\). What is that denominator? (enter a number)", k, k), strconv.Itoa(k-1),
			fmt.Sprintf("The denominator is \\(%d-1=%d\\).", k, k-1), f)
	default:
		return w9rC(kp, fmt.Sprintf("Does the power series \\(\\sum n!x^n\\) converge at \\(x=1/%d\\)? Enter 1 for yes and 0 for no.", k), "0",
			"The factorial coefficients grow faster than any fixed power, so the radius of convergence is zero.", f)
	}
}

func w9rTaubian(kp int) generator.Problem {
	c := w9rCR(1, 9)
	n := w9rCR(2, 9)
	p := 1 + rand.Intn(2)
	f := map[string]string{}
	w9rCI(f, "c", c)
	w9rCI(f, "n", n)
	w9rCI(f, "p", p)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("For \\(a_n=%d/%d\\), what is \\(n a_n\\)? (enter a number)", c, n), strconv.Itoa(c),
			fmt.Sprintf("\\(n a_n=%d\\).", c), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("For \\(a_n=%d/%d\\), does \\(n a_n\\to0\\) as \\(n\\to\\infty\\)? Enter 1 for yes and 0 for no.", c, n), "0",
			fmt.Sprintf("\\(n a_n=%d\\), which does not tend to zero.", c), f)
	default:
		code := 0
		if p == 2 {
			code = 1
		}
		return w9rC(kp, fmt.Sprintf("For \\(a_n=%d/n^%d\\), does \\(\\sum a_n\\) converge? Enter 1 for yes and 0 for no.", c, p), strconv.Itoa(code),
			fmt.Sprintf("The series is a constant multiple of the p-series with p=%d, so it converges exactly when p>1.", p), f)
	}
}

func w9rLambert(kp int) generator.Problem {
	n := w9rCR(2, 12)
	tau, sigma, sigma2 := 0, 0, 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			tau++
			sigma += d
			sigma2 += d * d
		}
	}
	f := map[string]string{}
	w9rCI(f, "n", n)
	switch kp {
	case 1:
		return w9rC(kp, fmt.Sprintf("In the Lambert series with \\(a_k=1\\), what is the coefficient of \\(x^{%d}\\)? (enter a number)", n), strconv.Itoa(tau),
			fmt.Sprintf("The coefficient is the number of divisors of %d, which is %d.", n, tau), f)
	case 2:
		return w9rC(kp, fmt.Sprintf("In the Lambert series with \\(a_k=k\\), what is the coefficient of \\(x^{%d}\\)? (enter a number)", n), strconv.Itoa(sigma),
			fmt.Sprintf("The coefficient is the sum of the divisors of %d, which is %d.", n, sigma), f)
	default:
		return w9rC(kp, fmt.Sprintf("In the Lambert series with \\(a_k=k^2\\), what is the coefficient of \\(x^{%d}\\)? (enter a number)", n), strconv.Itoa(sigma2),
			fmt.Sprintf("The coefficient is the sum of the squares of the divisors of %d, which is %d.", n, sigma2), f)
	}
}
