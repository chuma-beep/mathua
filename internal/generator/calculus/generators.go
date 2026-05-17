package calculus

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("calc.limit.concept", &limitConceptGen{})
	reg.Register("calc.limit.numeric", &limitNumericGen{})
	reg.Register("calc.limit.properties", &limitPropertiesGen{})
	reg.Register("calc.limit.infinity", &limitInfinityGen{})
	reg.Register("calc.limit.continuity", &limitContinuityGen{})

	reg.Register("calc.deriv.concept", &derivConceptGen{})
	reg.Register("calc.deriv.power_rule", &derivPowerRuleGen{})
	reg.Register("calc.deriv.sum_rule", &derivSumRuleGen{})
	reg.Register("calc.deriv.product_rule", &derivProductRuleGen{})
	reg.Register("calc.deriv.quotient_rule", &derivQuotientRuleGen{})
	reg.Register("calc.deriv.chain_rule", &derivChainRuleGen{})
	reg.Register("calc.deriv.trig", &derivTrigGen{})
	reg.Register("calc.deriv.exp_log", &derivExpLogGen{})
	reg.Register("calc.deriv.applications", &derivApplicationsGen{})
	reg.Register("calc.deriv.optimization", &derivOptimizationGen{})

	reg.Register("calc.integral.indefinite", &integralIndefiniteGen{})
	reg.Register("calc.integral.power_rule", &integralPowerRuleGen{})
	reg.Register("calc.integral.substitution", &integralSubstitutionGen{})
	reg.Register("calc.integral.definite", &integralDefiniteGen{})
	reg.Register("calc.integral.ftc", &integralFTCGen{})
	reg.Register("calc.integral.area_between", &integralAreaBetweenGen{})
	reg.Register("calc.integral.volume", &integralVolumeGen{})

	reg.Register("calc.deriv.implicit", &derivImplicitGen{})
	reg.Register("calc.deriv.related_rates", &derivRelatedRatesGen{})
	reg.Register("calc.integral.parts", &integralPartsGen{})
	reg.Register("calc.integral.partial_fractions", &integralPartialFractionsGen{})
}

type limitConceptGen struct{}

func (g *limitConceptGen) Generate(difficulty float64) generator.Problem {
	m := rand.Intn(4) + 1
	b := rand.Intn(10) + 1
	x0 := rand.Intn(5) + 1
	ans := m*x0 + b
	return generator.Problem{
		Question:    fmt.Sprintf("As x approaches %d, what value does f(x)=%dx+%d approach?", x0, m, b),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("Since f(x)=%dx+%d is continuous, the limit as x→%d equals f(%d)=%d(%d)+%d=%d.", m, b, x0, x0, m, x0, b, ans),
	}
}

type limitNumericGen struct{}

func (g *limitNumericGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 2
	ans := 2 * a
	return generator.Problem{
		Question:    fmt.Sprintf("Estimate lim x→%d of (x²-%d)/(x-%d) by evaluating near x=%d.", a, a*a, a, a),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("Factor: (x²-%d)/(x-%d) = (x-%d)(x+%d)/(x-%d) = x+%d for x≠%d. The limit as x→%d is %d+%d=%d.", a*a, a, a, a, a, a, a, a, a, a, ans),
	}
}

type limitPropertiesGen struct{}

func (g *limitPropertiesGen) Generate(difficulty float64) generator.Problem {
	fLim := rand.Intn(5) + 1
	gLim := rand.Intn(5) + 1
	op := rand.Intn(5)
	var ans int
	var opStr, desc string
	switch op {
	case 0:
		ans = fLim + gLim
		opStr = "f(x)+g(x)"
		desc = fmt.Sprintf("lim(f+g) = lim f + lim g = %d+%d=%d.", fLim, gLim, ans)
	case 1:
		ans = fLim - gLim
		opStr = "f(x)-g(x)"
		desc = fmt.Sprintf("lim(f-g) = lim f - lim g = %d-%d=%d.", fLim, gLim, ans)
	case 2:
		ans = fLim * gLim
		opStr = "f(x)·g(x)"
		desc = fmt.Sprintf("lim(f·g) = lim f · lim g = %d·%d=%d.", fLim, gLim, ans)
	case 3:
		c := rand.Intn(4) + 2
		ans = c * fLim
		opStr = fmt.Sprintf("%d·f(x)", c)
		desc = fmt.Sprintf("lim(%d·f) = %d·lim f = %d·%d=%d.", c, c, c, fLim, ans)
	case 4:
		if gLim == 0 {
			gLim = 1
		}
		ans = fLim / gLim
		opStr = "f(x)/g(x)"
		desc = fmt.Sprintf("lim(f/g) = lim f / lim g = %d/%d=%d.", fLim, gLim, ans)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("If lim f(x)=%d and lim g(x)=%d, what is lim (%s)?", fLim, gLim, opStr),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: desc,
	}
}

type limitInfinityGen struct{}

func (g *limitInfinityGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"lim x→∞ of 1/x", "0", "As x→∞, 1/x → 0."},
		{"lim x→∞ of 1/x²", "0", "As x→∞, 1/x² → 0."},
		{"lim x→∞ of x/(x+1)", "1", "Divide numerator and denominator by x: 1/(1+1/x) → 1."},
		{"lim x→∞ of 2x/(x+1)", "2", "Divide numerator and denominator by x: 2/(1+1/x) → 2."},
		{"lim x→∞ of (x²+1)/x²", "1", "(x²+1)/x² = 1+1/x² → 1."},
		{"lim x→∞ of 3", "3", "The limit of a constant is the constant itself."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is %s?", e.q),
		Answer:      e.a,
		Explanation: e.e,
	}
}

type limitContinuityGen struct{}

func (g *limitContinuityGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is f(x)=x² continuous at x=2?", "yes", "f(x)=x² is a polynomial, continuous everywhere. f(2)=4 and lim x→2 f(x)=4."},
		{"Is f(x)=sin(x) continuous at x=0?", "yes", "sin(x) is continuous everywhere. sin(0)=0 and lim x→0 sin(x)=0."},
		{"Is f(x)=e^x continuous at x=0?", "yes", "e^x is continuous everywhere. e^0=1 and lim x→0 e^x=1."},
		{"Is f(x)=|x| continuous at x=0?", "yes", "|x| is continuous at 0. |0|=0 and lim x→0 |x|=0."},
		{"Is f(x)=1/x continuous at x=0?", "no", "1/x has an infinite discontinuity at x=0 (vertical asymptote)."},
		{"Is f(x)=1/x² continuous at x=0?", "no", "1/x² has an infinite discontinuity at x=0."},
		{"Is f(x)=tan(x) continuous at x=π/2?", "no", "tan(x) has an infinite discontinuity at x=π/2."},
		{"Is f(x)=|x|/x continuous at x=0?", "no", "|x|/x has a jump discontinuity at x=0 (left limit=-1, right limit=1)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
	}
}

type derivConceptGen struct{}

func (g *derivConceptGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		q string
		a int
		e string
	}
	table := []entry{
		{"What is the slope of f(x)=x² at x=3?", 6, "f'(x)=2x, so f'(3)=2·3=6."},
		{"What is the slope of f(x)=x² at x=2?", 4, "f'(x)=2x, so f'(2)=2·2=4."},
		{"What is the slope of f(x)=x² at x=1?", 2, "f'(x)=2x, so f'(1)=2·1=2."},
		{"What is the slope of f(x)=2x² at x=2?", 8, "f'(x)=4x, so f'(2)=4·2=8."},
		{"What is the slope of f(x)=x³ at x=2?", 12, "f'(x)=3x², so f'(2)=3·4=12."},
		{"What is the slope of f(x)=x³ at x=1?", 3, "f'(x)=3x², so f'(1)=3·1=3."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      fmt.Sprintf("%d", e.a),
		Explanation: e.e,
	}
}

type derivPowerRuleGen struct{}

func (g *derivPowerRuleGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(5) + 1
	n := rand.Intn(4) + 2
	coef := a * n
	exp := n - 1
	return generator.Problem{
		Question:    fmt.Sprintf("Find f'(x) if f(x)=%dx^%d.", a, n),
		Answer:      singleTerm(coef, exp),
		Explanation: fmt.Sprintf("Power rule: f'(x)=%d·%d·x^(%d-1)=%s.", a, n, n, singleTerm(coef, exp)),
	}
}

type derivSumRuleGen struct{}

func (g *derivSumRuleGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 2
	n := rand.Intn(3) + 2
	b := rand.Intn(4) + 2
	m := rand.Intn(3) + 2
	for m == n {
		m = rand.Intn(3) + 2
	}
	coef1 := a * n
	exp1 := n - 1
	coef2 := b * m
	exp2 := m - 1
	answer := formatPoly([]term{{coef1, exp1}, {coef2, exp2}})
	return generator.Problem{
		Question:    fmt.Sprintf("Find f'(x) if f(x)=%dx^%d+%dx^%d.", a, n, b, m),
		Answer:      answer,
		Explanation: fmt.Sprintf("f'(x) = %d·%dx^%d + %d·%dx^%d = %s.", a, n, n-1, b, m, m-1, answer),
	}
}

type derivProductRuleGen struct{}

func (g *derivProductRuleGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(3) + 1
	m := rand.Intn(2) + 1
	a := rand.Intn(3) + 1
	// f(x) = x^n (x^m + a) = x^(n+m) + a·x^n
	// f'(x) = (n+m)x^(n+m-1) + a·n·x^(n-1)
	coef1 := n + m
	exp1 := n + m - 1
	coef2 := a * n
	exp2 := n - 1
	answer := formatPoly([]term{{coef1, exp1}, {coef2, exp2}})
	return generator.Problem{
		Question:    fmt.Sprintf("Find f'(x) if f(x)=x^%d(x^%d+%d).", n, m, a),
		Answer:      answer,
		Explanation: fmt.Sprintf("Product rule: f'(x)=%dx^%d·(x^%d+%d)+x^%d·%dx^%d = %s.", n, n-1, m, a, n, m, m-1, answer),
	}
}

type derivQuotientRuleGen struct{}

func (g *derivQuotientRuleGen) Generate(difficulty float64) generator.Problem {
	// f(x) = (x+1)/(x-1), f'(x) = -2/(x-1)², ask for f'(2) = -2
	return generator.Problem{
		Question:    "Find f'(2) if f(x)=(x+1)/(x-1).",
		Answer:      "-2",
		Explanation: "Quotient rule: f'(x)=((x-1)·1-(x+1)·1)/(x-1)² = -2/(x-1)². f'(2) = -2/(1)² = -2.",
	}
}

type derivChainRuleGen struct{}

func (g *derivChainRuleGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 1
	b := rand.Intn(4) + 1
	n := rand.Intn(2) + 2
	// f(x) = (ax+b)^n
	// f'(x) = n·a·(ax+b)^(n-1)
	coef := n * a
	inExp := n - 1
	inner := fmt.Sprintf("%dx+%d", a, b)
	answer := fmt.Sprintf("%d(%s)^%d", coef, inner, inExp)
	return generator.Problem{
		Question:    fmt.Sprintf("Find f'(x) if f(x)=%s^%d.", inner, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Chain rule: f'(x)=%d·(%s)^%d·%d = %s.", n, inner, inExp, a, answer),
	}
}

type derivTrigGen struct{}

func (g *derivTrigGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find f'(x) if f(x)=sin(x).", "cos(x)", "d/dx sin(x) = cos(x)."},
		{"Find f'(x) if f(x)=cos(x).", "-sin(x)", "d/dx cos(x) = -sin(x)."},
		{"What is d/dx sin(x)?", "cos(x)", "The derivative of sine is cosine."},
		{"What is d/dx cos(x)?", "-sin(x)", "The derivative of cosine is -sine."},
		{"Find f'(x) if f(x)=sin(2x).", "2cos(2x)", "Chain rule: d/dx sin(2x) = cos(2x)·2 = 2cos(2x)."},
		{"Find f'(x) if f(x)=cos(3x).", "-3sin(3x)", "Chain rule: d/dx cos(3x) = -sin(3x)·3 = -3sin(3x)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
	}
}

type derivExpLogGen struct{}

func (g *derivExpLogGen) Generate(difficulty float64) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find f'(x) if f(x)=e^x.", "e^x", "d/dx e^x = e^x."},
		{"Find f'(x) if f(x)=ln(x).", "1/x", "d/dx ln(x) = 1/x."},
		{"Find f'(x) if f(x)=e^(3x).", "3e^(3x)", "Chain rule: d/dx e^(3x) = e^(3x)·3 = 3e^(3x)."},
		{"Find f'(x) if f(x)=e^(2x).", "2e^(2x)", "Chain rule: d/dx e^(2x) = e^(2x)·2 = 2e^(2x)."},
		{"Find f'(x) if f(x)=ln(2x).", "1/x", "Chain rule: d/dx ln(2x) = (1/(2x))·2 = 1/x."},
		{"What is d/dx e^x?", "e^x", "The derivative of e^x is e^x."},
		{"What is d/dx ln(x)?", "1/x", "The derivative of ln(x) is 1/x."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
	}
}

type derivApplicationsGen struct{}

func (g *derivApplicationsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(3) + 1
	b := rand.Intn(10) + 1
	c := rand.Intn(10) + 1
	t0 := rand.Intn(4) + 1
	vel := 2*a*t0 + b
	return generator.Problem{
		Question:    fmt.Sprintf("If s(t)=%dt²+%dt+%d, what is velocity at t=%d?", a, b, c, t0),
		Answer:      fmt.Sprintf("%d", vel),
		Explanation: fmt.Sprintf("v(t)=s'(t)=%dt+%d. v(%d)=%d(%d)+%d=%d.", 2*a, b, t0, 2*a, t0, b, vel),
	}
}

type derivOptimizationGen struct{}

func (g *derivOptimizationGen) Generate(difficulty float64) generator.Problem {
	perim := (rand.Intn(4) + 3) * 4 // 12, 16, 20, 24
	width := perim / 4
	area := width * width
	return generator.Problem{
		Question:    fmt.Sprintf("A rectangle has perimeter %d. What width (in same units) maximizes area?", perim),
		Answer:      fmt.Sprintf("%d", width),
		Explanation: fmt.Sprintf("Let width=w, height=%d-w. Area A(w)=w(%d-w)=%dw-w². A'(w)=%d-2w=0 → w=%d. Maximum area = %d.", perim/2, perim/2, perim/2, perim/2, width, area),
	}
}

type integralIndefiniteGen struct{}

func (g *integralIndefiniteGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	n := rand.Intn(4) + 1
	num := a
	den := n + 1
	exp := n + 1
	var answer string
	if num%den == 0 {
		coef := num / den
		answer = fmt.Sprintf("%dx^%d+C", coef, exp)
	} else {
		answer = fmt.Sprintf("(%d/%d)x^%d+C", num, den, exp)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find ∫%dx^%d dx.", a, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Power rule for integration: ∫%dx^%d dx = %d/(%d)·x^(%d+1)+C = %s.", a, n, a, n, n, answer),
	}
}

type integralPowerRuleGen struct{}

func (g *integralPowerRuleGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(4) + 1
	a := n + 1
	exp := n + 1
	coef := a / (n + 1)
	var answer string
	if coef == 1 {
		answer = fmt.Sprintf("x^%d", exp)
	} else {
		answer = fmt.Sprintf("%dx^%d", coef, exp)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find ∫%dx^%d dx (use power rule for integration, omit +C).", a, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("∫%dx^%d dx = %d/(%d+1)·x^(%d+1)+C = %s+C.", a, n, a, n, n, answer),
	}
}

type integralSubstitutionGen struct{}

func (g *integralSubstitutionGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(3) + 2
	coef := n
	exp := n - 1
	answer := fmt.Sprintf("e^(x^%d)+C", n)
	return generator.Problem{
		Question:    fmt.Sprintf("Find ∫%dx^%d·e^(x^%d) dx.", coef, exp, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Let u=x^%d, du=%dx^%d dx. Integral becomes ∫e^u du = e^u+C = e^(x^%d)+C.", n, n, exp, n),
	}
}

type integralDefiniteGen struct{}

func (g *integralDefiniteGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(4) + 1
	n := rand.Intn(3) + 1
	lower := rand.Intn(2)
	upper := rand.Intn(3) + 1
	for upper <= lower {
		upper = rand.Intn(3) + 1
	}
	f := func(x int) int {
		return mathutil.IntPow(x, n+1) * a / (n + 1)
	}
	ans := f(upper) - f(lower)
	return generator.Problem{
		Question:    fmt.Sprintf("Find ∫₀^%d %dx^%d dx.", upper, a, n),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("∫%dx^%d dx = %d/(%d)·x^(%d). From 0 to %d: %d·%d^(%d)/%d - 0 = %d.", a, n, a, n+1, n+1, upper, a, upper, n+1, n+1, ans),
	}
}

type integralFTCGen struct{}

func (g *integralFTCGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(3) + 1
	answer := fmt.Sprintf("x^%d", n)
	return generator.Problem{
		Question:    fmt.Sprintf("If F(x)=∫₀ˣ t^%d dt, what is F'(x)?", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("By the Fundamental Theorem of Calculus, F'(x)=x^%d.", n),
	}
}

type integralAreaBetweenGen struct{}

func (g *integralAreaBetweenGen) Generate(difficulty float64) generator.Problem {
	// Area between y=x and y=x² from 0 to 1 = 1/6 ≈ 0.1667
	return generator.Problem{
		Question:    "Find area between y=x and y=x² from x=0 to x=1.",
		Answer:      "0.1667",
		Explanation: "∫(x-x²)dx from 0 to 1 = [x²/2-x³/3]₀¹ = 1/2-1/3 = 1/6 ≈ 0.1667.",
	}
}

type integralVolumeGen struct{}

func (g *integralVolumeGen) Generate(difficulty float64) generator.Problem {
	// V = π∫(√x)²dx from 0 to 4 = π∫x dx = π·x²/2 from 0 to 4 = π·16/2 = 8π
	return generator.Problem{
		Question:    "Find the volume when y=√x from x=0 to 4 is revolved around the x-axis.",
		Answer:      "8π",
		Explanation: "V = π∫(√x)²dx = π∫x dx = π[x²/2]₀⁴ = π(16/2-0) = 8π.",
	}
}

type derivImplicitGen struct{}

func (g *derivImplicitGen) Generate(difficulty float64) generator.Problem {
	r := rand.Intn(3) + 3
	a := rand.Intn(3) + 1
	b := r*r - a*a
	// Ensure integer b
	for b <= 0 {
		a = rand.Intn(r) + 1
		b = r*r - a*a
	}
	bSqrt := 1
	for i := 2; i*i <= b; i++ {
		if i*i == b {
			bSqrt = i
			break
		}
	}
	if bSqrt == 1 {
		bSqrt = b
	}
	// dy/dx = -a/bSqrt
	gcd := mathutil.GCD(mathutil.Abs(a), bSqrt)
	num := -a / gcd
	den := bSqrt / gcd
	answer := fmt.Sprintf("%d/%d", num, den)
	if den == 1 {
		answer = fmt.Sprintf("%d", num)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find dy/dx at (%d,%d) for x²+y²=%d.", a, bSqrt, r*r),
		Answer:      answer,
		Explanation: fmt.Sprintf("2x+2y·dy/dx=0 → dy/dx=-x/y. At (%d,%d): dy/dx=-%d/%d=%s.", a, bSqrt, a, bSqrt, answer),
	}
}

type derivRelatedRatesGen struct{}

func (g *derivRelatedRatesGen) Generate(difficulty float64) generator.Problem {
	dr := rand.Intn(3) + 1
	r := rand.Intn(3) + 3
	ans := 2 * 3 * r * dr
	answer := fmt.Sprintf("%dπ", ans)
	return generator.Problem{
		Question:    fmt.Sprintf("A circle's radius grows at %d cm/s. How fast is area growing when r=%d?", dr, r),
		Answer:      answer,
		Explanation: fmt.Sprintf("A=πr², dA/dt=2πr·dr/dt=2π(%d)(%d)=%dπ cm²/s.", r, dr, ans),
	}
}

type integralPartsGen struct{}

func (g *integralPartsGen) Generate(difficulty float64) generator.Problem {
	return generator.Problem{
		Question:    "Find ∫x·e^x dx.",
		Answer:      "xe^x-e^x+C",
		Explanation: "Integration by parts: let u=x, dv=e^x dx → du=dx, v=e^x. ∫x·e^x dx = x·e^x - ∫e^x dx = x·e^x - e^x + C.",
	}
}

type integralPartialFractionsGen struct{}

func (g *integralPartialFractionsGen) Generate(difficulty float64) generator.Problem {
	return generator.Problem{
		Question:    "Find ∫1/(x²-1) dx.",
		Answer:      "(1/2)ln|x-1|-(1/2)ln|x+1|+C",
		Explanation: "Partial fractions: 1/(x²-1) = 1/2·(1/(x-1) - 1/(x+1)). Integrate: (1/2)ln|x-1| - (1/2)ln|x+1| + C.",
	}
}

func singleTerm(coef, exp int) string {
	if coef == 0 {
		return "0"
	}
	if exp == 0 {
		return fmt.Sprintf("%d", coef)
	}
	if exp == 1 {
		if coef == 1 {
			return "x"
		}
		if coef == -1 {
			return "-x"
		}
		return fmt.Sprintf("%dx", coef)
	}
	if coef == 1 {
		return fmt.Sprintf("x^%d", exp)
	}
	if coef == -1 {
		return fmt.Sprintf("-x^%d", exp)
	}
	return fmt.Sprintf("%dx^%d", coef, exp)
}

type term struct {
	coef int
	exp  int
}

func formatPoly(terms []term) string {
	var parts []string
	for _, t := range terms {
		if t.coef == 0 {
			continue
		}
		parts = append(parts, singleTerm(t.coef, t.exp))
	}
	if len(parts) == 0 {
		return "0"
	}
	// Build the string, handling signs
	var sb strings.Builder
	sb.WriteString(parts[0])
	for _, p := range parts[1:] {
		if strings.HasPrefix(p, "-") {
			sb.WriteString(p)
		} else {
			sb.WriteString("+")
			sb.WriteString(p)
		}
	}
	return sb.String()
}
