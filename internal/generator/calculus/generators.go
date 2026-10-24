package calculus

import (
	"fmt"
	"math/rand"
	"strconv"
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

	reg.Register("calc.deriv.cauchy_mvt", &cauchyMVTGen{})
	reg.Register("calc.deriv.convexity", &convexityGen{})
	reg.Register("calc.deriv.difference_quotient", &differenceQuotientGen{})
	reg.Register("calc.deriv.differential", &differentialGen{})
	reg.Register("calc.deriv.fermat", &fermatGen{})
	reg.Register("calc.deriv.mvt", &mvtGen{})
	reg.Register("calc.deriv.non_differentiability", &nonDiffGen{})
	reg.Register("calc.deriv.partial", &partialDerivGen{})
	reg.Register("calc.deriv.rolle", &rolleGen{})
	reg.Register("calc.integral.arc_length", &arcLengthGen{})
	reg.Register("calc.integral.exp_integral", &expIntegralGen{})
	reg.Register("calc.integral.improper", &improperIntegralGen{})
	reg.Register("calc.integral.numerical", &numericalIntegralGen{})
	reg.Register("calc.integral.riemann_criteria", &riemannCriteriaGen{})
	reg.Register("calc.integral.trig_integrals", &trigIntegralsGen{})
	reg.Register("calc.integral.trig_substitution", &trigSubstitutionGen{})
	reg.Register("calc.integral.weierstrass_sub", &weierstrassSubGen{})
	reg.Register("calc.limit.algebra", &limitAlgebraGen{})
	reg.Register("calc.limit.asymptotes", &asymptotesGen{})
	reg.Register("calc.limit.big_o", &bigOGen{})
	reg.Register("calc.limit.discontinuity", &discontinuityGen{})
	reg.Register("calc.limit.indeterminate", &indeterminateGen{})
	reg.Register("calc.limit.lhopital", &lhopitalGen{})
	reg.Register("calc.limit.little_o", &littleOGen{})
	reg.Register("calc.limit.squeeze", &squeezeGen{})
	reg.Register("calc.limit.supremum", &supremumGen{})
	reg.Register("calc.limit.uniform_continuity", &uniformContinuityGen{})
	reg.Register("calc.limit.weierstrass", &weierstrassLimitGen{})
	reg.Register("calc.seq.cauchy", &cauchySeqGen{})
	reg.Register("calc.seq.concept", &seqConceptGen{})
	reg.Register("calc.seq.convergence", &seqConvergenceGen{})
	reg.Register("calc.seq.euler", &eulerSeqGen{})
	reg.Register("calc.seq.function_sequences", &functionSequencesGen{})
	reg.Register("calc.seq.monotone", &monotoneSeqGen{})
	reg.Register("calc.series.alternating", &alternatingSeriesGen{})
	reg.Register("calc.series.cauchy_criterion", &cauchyCriterionSeriesGen{})
	reg.Register("calc.series.concept", &seriesConceptGen{})
	reg.Register("calc.series.fourier", &fourierSeriesGen{})
	reg.Register("calc.series.function_series", &functionSeriesGen{})
	reg.Register("calc.series.harmonic", &harmonicSeriesGen{})
	reg.Register("calc.series.integral_test", &integralTestGen{})
	reg.Register("calc.series.positive_terms", &positiveTermsGen{})
	reg.Register("calc.series.power", &powerSeriesGen{})
	reg.Register("calc.series.root_test", &rootTestGen{})
	reg.Register("calc.series.taylor", &taylorSeriesGen{})
	reg.Register("calc.deriv.trig_other", &derivTrigOtherGen{})
	reg.Register("calc.deriv.inverse_trig", &derivInverseTrigGen{})
	reg.Register("calc.deriv.general_exp", &derivGeneralExpGen{})
	reg.Register("calc.deriv.log_diff", &derivLogDiffGen{})
	reg.Register("calc.deriv.higher_order", &derivHigherOrderGen{})
	reg.Register("calc.deriv.inverse_func", &derivInverseFuncGen{})
	reg.Register("calc.deriv.critical_points", &derivCriticalPointsGen{})
	reg.Register("calc.deriv.second_test", &derivSecondTestGen{})
	reg.Register("calc.deriv.inflection", &derivInflectionGen{})
	reg.Register("calc.deriv.global_extrema", &derivGlobalExtremaGen{})
	reg.Register("calc.deriv.curve_sketch", &derivCurveSketchGen{})
	reg.Register("calc.deriv.tangent_line", &derivTangentLineGen{})
	reg.Register("calc.deriv.linear_approx", &derivLinearApproxGen{})
	reg.Register("calc.deriv.newton", &derivNewtonGen{})
	reg.Register("calc.deriv.error_prop", &derivErrorPropGen{})
	reg.Register("calc.deriv.implicit_second", &derivImplicitSecondGen{})
	reg.Register("calc.deriv.chain_power", &derivChainPowerGen{})
	reg.Register("calc.deriv.chain_multi", &derivChainMultiGen{})
	reg.Register("calc.deriv.gradient", &derivGradientGen{})
	reg.Register("calc.deriv.directional", &derivDirectionalGen{})
	reg.Register("calc.deriv.partial_higher", &derivPartialHigherGen{})
	reg.Register("calc.deriv.clairaut", &derivClairautGen{})
	reg.Register("calc.deriv.total_diff", &derivTotalDiffGen{})
	reg.Register("calc.deriv.diff_continuity", &derivDiffContinuityGen{})
	reg.Register("calc.deriv.jacobian", &derivJacobianGen{})
	reg.Register("calc.integral.net_change", &integralNetChangeGen{})
	reg.Register("calc.integral.variable_limits", &integralVariableLimitsGen{})
	reg.Register("calc.integral.mvt", &integralMVTGen{})
	reg.Register("calc.integral.work", &integralWorkGen{})
	reg.Register("calc.integral.center_mass", &integralCenterMassGen{})
	reg.Register("calc.integral.surface_area", &integralSurfaceAreaGen{})
	reg.Register("calc.integral.disc", &discMethodGen{})
	reg.Register("calc.integral.washer", &washerMethodGen{})
	reg.Register("calc.integral.shell", &shellMethodGen{})
	reg.Register("calc.integral.trapezoid", &integralTrapezoidGen{})
	reg.Register("calc.integral.simpson", &integralSimpsonGen{})
	reg.Register("calc.integral.improper_compare", &integralImproperCompareGen{})
	reg.Register("calc.integral.p_test", &integralPTestGen{})
	reg.Register("calc.series.ratio_test", &ratioTestGen{})
	reg.Register("calc.series.geometric", &geometricSeriesGen{})
	reg.Register("calc.series.nth_term", &nthTermGen{})
	reg.Register("calc.series.telescoping", &telescopingSeriesGen{})
	reg.Register("calc.series.direct_compare", &directCompareGen{})
	reg.Register("calc.series.limit_compare", &limitCompareGen{})
	reg.Register("calc.series.absolute", &absoluteConvergenceGen{})
	reg.Register("calc.series.remainder", &alternatingRemainderGen{})
	reg.Register("calc.series.radius", &convergenceRadiusGen{})
	reg.Register("calc.limit.epsilon_delta", &epsilonDeltaGen{})
	reg.Register("calc.limit.one_sided", &oneSidedLimitGen{})
	reg.Register("calc.limit.infinite", &infiniteLimitGen{})
	reg.Register("calc.limit.sine_limit", &sineLimitGen{})
	reg.Register("calc.limit.exp_limit", &expLimitGen{})
	reg.Register("calc.seq.bounded", &boundedSeqGen{})
	reg.Register("calc.integral.convergence", &convergenceGen{})
	reg.Register("calc.integral.irrational", &irrationalIntegralGen{})
	reg.Register("calc.integral.strategy", &strategyGen{})
	reg.Register("calc.integral.reduction", &reductionGen{})
	registerMultivariable(reg)
	registerVector(reg)
	registerVector3D(reg)
	registerSurface(reg)
	registerRemaining(reg)
	registerNumeric(reg)
	registerBatch(reg)
	registerRaise(reg)
}

type limitConceptGen struct{}

func (g *limitConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	m := rand.Intn(scale*2) + 1
	b := rand.Intn(scale*5) + 1
	x0 := rand.Intn(scale*2) + 1
	ans := m*x0 + b
	return generator.Problem{
		Question:    fmt.Sprintf("As \\(x\\) approaches %d, what value does \\(f(x)=%dx+%d\\) approach?", x0, m, b),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("Since f(x)=%dx+%d is continuous, the limit as x→%d equals f(%d)=%d(%d)+%d=%d.", m, b, x0, x0, m, x0, b, ans),
	}
}

type limitNumericGen struct{}

func (g *limitNumericGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) + 2
	ans := 2 * a
	return generator.Problem{
		Question:    fmt.Sprintf("Estimate \\(\\lim_{x \\to %d}\\) of \\(\\frac{x^{2}-%d}{x-%d}\\) by evaluating near \\(x=%d\\).", a, a*a, a, a),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("Factor: (x²-%d)/(x-%d) = (x-%d)(x+%d)/(x-%d) = x+%d for x≠%d. The limit as x→%d is %d+%d=%d.", a*a, a, a, a, a, a, a, a, a, a, ans),
	}
}

// selectMode selects template mode 0, 1 or 2 when ctx.KPOrdinal is 1, 2 or 3.
// When ctx.KPOrdinal is 0, it selects randomly from totalModes.
func selectMode(ctx generator.GeneratorContext, totalModes int) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal - 1
	}
	return rand.Intn(totalModes)
}

type limitPropertiesGen struct{}

func (g *limitPropertiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	fLim := rand.Intn(scale*2) + 1
	gLim := rand.Intn(scale*2) + 1
	op := selectMode(ctx, 5)
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
		Question:    fmt.Sprintf("If \\(\\lim f(x)=%d\\) and \\(\\lim g(x)=%d\\), what is \\(\\lim (%s)\\)?", fLim, gLim, opStr),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: desc,
	}
}

type limitContinuityGen struct{}

func (g *limitContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is \\(f(x)=x^{2}\\) continuous at \\(x=2\\)?", "yes", "\\(f(x)=x^{2}\\) is a polynomial, continuous everywhere. \\(f(2)=4\\) and \\(\\lim_{x \\to 2} f(x)=4\\)."},
		{"Is \\(f(x)=\\sin(x)\\) continuous at \\(x=0\\)?", "yes", "\\(\\sin(x)\\) is continuous everywhere. \\(\\sin(0)=0\\) and \\(\\lim_{x \\to 0} \\sin(x)=0\\)."},
		{"Is \\(f(x)=e^{x}\\) continuous at \\(x=0\\)?", "yes", "\\(e^{x}\\) is continuous everywhere. \\(e^{0}=1\\) and \\(\\lim_{x \\to 0} e^{x}=1\\)."},
		{"Is \\(f(x)=|x|\\) continuous at \\(x=0\\)?", "yes", "\\(|x|\\) is continuous at 0. \\(|0|=0\\) and \\(\\lim_{x \\to 0} |x|=0\\)."},
		{"Is \\(f(x)=1/x\\) continuous at \\(x=0\\)?", "no", "\\(1/x\\) has an infinite discontinuity at x=0 (vertical asymptote)."},
		{"Is \\(f(x)=1/x^{2}\\) continuous at \\(x=0\\)?", "no", "\\(1/x^{2}\\) has an infinite discontinuity at \\(x=0\\)."},
		{"Is \\(f(x)=\\tan(x)\\) continuous at \\(x=\\pi/2\\)?", "no", "\\(\\tan(x)\\) has an infinite discontinuity at \\(x=\\pi/2\\)."},
		{"Is \\(f(x)=|x|/x\\) continuous at \\(x=0\\)?", "no", "\\(|x|/x\\) has a jump discontinuity at \\(x=0\\) (left limit \\(=-1\\), right limit \\(=1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
	}
}

type derivConceptGen struct{}

func (g *derivConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q string
		a int
		e string
	}
	a := randRange(ctx.Difficulty, 1, 5)
	n := randRange(ctx.Difficulty, 2, 4)
	b := randRange(ctx.Difficulty, 1, 6)
	deriv := scaledXPow(a*n, n-1)
	return generator.Problem{
		Question:    fmt.Sprintf("What is the slope of \\(f(x)=%s\\) at \\(x=%d\\)?", scaledXPow(a, n), b),
		Answer:      strconv.Itoa(a * n * intPow(b, n-1)),
		Explanation: fmt.Sprintf("\\(f'(x)=%s\\), so \\(f'(%d)=%d\\).", deriv, b, a*n*intPow(b, n-1)),
	}
}

type derivPowerRuleGen struct{}

func (g *derivPowerRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) + 1
	n := rand.Intn(scale) + 2
	coef := a * n
	exp := n - 1
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=%dx^{%d}\\).", a, n),
		Answer:      singleTerm(coef, exp),
		Explanation: fmt.Sprintf("Power rule: f'(x)=%d·%d·x^(%d-1)=%s.", a, n, n, singleTerm(coef, exp)),
	}
}

type derivSumRuleGen struct{}

func (g *derivSumRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 2
	n := rand.Intn(max(2, scale)) + 2
	b := rand.Intn(max(1, scale*2)) + 2
	m := rand.Intn(max(2, scale)) + 2
	for m == n {
		m = rand.Intn(max(2, scale)) + 2
	}
	coef1 := a * n
	exp1 := n - 1
	coef2 := b * m
	exp2 := m - 1
	answer := formatPoly([]term{{coef1, exp1}, {coef2, exp2}})
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=%dx^{%d}+%dx^{%d}\\).", a, n, b, m),
		Answer:      answer,
		Explanation: fmt.Sprintf("f'(x) = %d·%dx^%d + %d·%dx^%d = %s.", a, n, n-1, b, m, m-1, answer),
	}
}

type derivProductRuleGen struct{}

func (g *derivProductRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(1, scale)) + 1
	m := rand.Intn(max(1, scale)) + 1
	a := rand.Intn(max(1, scale*2)) + 1
	// f(x) = x^n (x^m + a) = x^(n+m) + a·x^n
	// f'(x) = (n+m)x^(n+m-1) + a·n·x^(n-1)
	coef1 := n + m
	exp1 := n + m - 1
	coef2 := a * n
	exp2 := n - 1
	answer := formatPoly([]term{{coef1, exp1}, {coef2, exp2}})
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=x^{%d}(x^{%d}+%d)\\).", n, m, a),
		Answer:      answer,
		Explanation: fmt.Sprintf("Product rule: f'(x)=%dx^%d·(x^%d+%d)+x^%d·%dx^%d = %s.", n, n-1, m, a, n, m, m-1, answer),
	}
}

type derivQuotientRuleGen struct{}

func (g *derivQuotientRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 6)
	b := randRange(ctx.Difficulty, 1, 6)
	if a == b {
		b++
	}
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\frac{x+%d}{x+%d}\\).", a, b),
			Answer:      fmt.Sprintf("%d/(x+%d)^2", b-a, b),
			Explanation: fmt.Sprintf("Quotient rule: the numerator is \\((x+%d)-(x+%d)=%d\\), so \\(f'(x)=%d/(x+%d)^2\\).", b, a, b-a, b-a, b),
		}
	}
	c := randRange(ctx.Difficulty, 1, 6)
	d := randRange(ctx.Difficulty, 1, 6)
	num := fmt.Sprintf("%dx+%d", a, b)
	den := fmt.Sprintf("%dx+%d", c, d)
	if a == 1 {
		num = fmt.Sprintf("x+%d", b)
	}
	if c == 1 {
		den = fmt.Sprintf("x+%d", d)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\frac{%s}{%s}\\).", num, den),
		Answer:      fmt.Sprintf("%d/(%s)^2", a*d-b*c, den),
		Explanation: fmt.Sprintf("Quotient rule: \\((ad-bc)\\) with \\(a=%d\\), \\(d=%d\\), \\(b=%d\\), \\(c=%d\\) gives \\(%d/(%s)^2\\).", a, d, b, c, a*d-b*c, den),
	}
}

type derivChainRuleGen struct{}

func (g *derivChainRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	n := rand.Intn(max(1, scale)) + 2
	// f(x) = (ax+b)^n
	// f'(x) = n·a·(ax+b)^(n-1)
	coef := n * a
	inExp := n - 1
	inner := fmt.Sprintf("%dx+%d", a, b)
	answer := fmt.Sprintf("%d(%s)^%d", coef, inner, inExp)
	return generator.Problem{
		// The exponent used to bind only to the constant: f(x)=6x+1^4 read as
		// 6x + 1^4, whose derivative is a constant, and three reports said so.
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=(%s)^{%d}\\).", inner, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Chain rule: f'(x)=%d·(%s)^%d·%d = %s.", n, inner, inExp, a, answer),
	}
}

type derivTrigGen struct{}

func (g *derivTrigGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	a := randRange(ctx.Difficulty, 1, 9)
	b := randRange(ctx.Difficulty, 1, 9)
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\sin(%s)\\).", argX(a)),
			Answer:      coefTerm(a, fmt.Sprintf("cos(%s)", argX(a))),
			Explanation: fmt.Sprintf("Chain rule: \\(\\frac{d}{dx}\\sin(%s)=\\cos(%s)\\cdot %d\\).", argX(a), argX(a), a),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\cos(%s)\\).", argX(a)),
			Answer:      "-" + coefTerm(a, fmt.Sprintf("sin(%s)", argX(a))),
			Explanation: fmt.Sprintf("Chain rule: \\(\\frac{d}{dx}\\cos(%s)=-\\sin(%s)\\cdot %d\\).", argX(a), argX(a), a),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\sin(%s)+\\cos(%s)\\).", argX(a), argX(b)),
			Answer:      fmt.Sprintf("%s-%s", coefTerm(a, fmt.Sprintf("cos(%s)", argX(a))), coefTerm(b, fmt.Sprintf("sin(%s)", argX(b)))),
			Explanation: fmt.Sprintf("Differentiate each term: \\(%s\\) then \\(-%s\\).", coefTerm(a, fmt.Sprintf("cos(%s)", argX(a))), coefTerm(b, fmt.Sprintf("sin(%s)", argX(b)))),
		}
	}
}

type derivExpLogGen struct{}

func (g *derivExpLogGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	a := randRange(ctx.Difficulty, 1, 9)
	b := randRange(ctx.Difficulty, 1, 5)
	switch rand.Intn(3) {
	case 0:
		inner := "x"
		if a > 1 {
			inner = fmt.Sprintf("%dx", a)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=e^{%s}\\).", inner),
			Answer:      coefTerm(a, fmt.Sprintf("e^(%s)", inner)),
			Explanation: fmt.Sprintf("Chain rule: \\(\\frac{d}{dx}e^{%s}=e^{%s}\\cdot %d\\).", inner, inner, a),
		}
	case 1:
		arg := fmt.Sprintf("%dx+%d", a, b)
		if a == 1 {
			arg = fmt.Sprintf("x+%d", b)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\ln(%s)\\).", arg),
			Answer:      fmt.Sprintf("%d/(%s)", a, arg),
			Explanation: fmt.Sprintf("Chain rule: \\(\\frac{d}{dx}\\ln(%s)=\\frac{%d}{%s}\\).", arg, a, arg),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\ln(x^{%d})\\).", a),
			Answer:      fmt.Sprintf("%d/x", a),
			Explanation: fmt.Sprintf("\\(\\ln(x^{%d})=%d\\ln x\\), so the derivative is \\(%d/x\\).", a, a, a),
		}
	}
}

type derivApplicationsGen struct{}

func (g *derivApplicationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	c := rand.Intn(max(1, scale*2)) + 1
	t0 := rand.Intn(max(1, scale*2)) + 1
	vel := 2*a*t0 + b
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(s(t)=%dt^{2}+%dt+%d\\), what is velocity at \\(t=%d\\)?", a, b, c, t0),
		Answer:      fmt.Sprintf("%d", vel),
		Explanation: fmt.Sprintf("v(t)=s'(t)=%dt+%d. v(%d)=%d(%d)+%d=%d.", 2*a, b, t0, 2*a, t0, b, vel),
	}
}

type derivOptimizationGen struct{}

func (g *derivOptimizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	perim := (rand.Intn(max(1, scale*2)) + 3) * 4 // 12, 16, 20, 24
	width := perim / 4
	area := width * width
	return generator.Problem{
		Question:    fmt.Sprintf("A rectangle has perimeter %d. What width (in same units) maximizes area?", perim),
		Answer:      fmt.Sprintf("%d", width),
		Explanation: fmt.Sprintf("Let width=w, height=%d-w. Area A(w)=w(%d-w)=%dw-w². A'(w)=%d-2w=0 → w=%d. Maximum area = %d.", perim/2, perim/2, perim/2, perim/2, width, area),
	}
}

type integralIndefiniteGen struct{}

func (g *integralIndefiniteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	n := rand.Intn(max(1, scale)) + 1
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
		Question:    fmt.Sprintf("Find \\(\\int %d x^{%d} \\, dx\\).", a, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Power rule for integration: ∫%dx^%d dx = %d/(%d)·x^(%d+1)+C = %s.", a, n, a, n, n, answer),
	}
}

type integralPowerRuleGen struct{}

func (g *integralPowerRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := randRange(ctx.Difficulty, 1, 2+scale)
	// The coefficient is a multiple of n+1 so that the antiderivative stays an
	// integer multiple of x^(n+1). An arbitrary coefficient would divide to a
	// fraction, and the concept is graded as an expression.
	k := randRange(ctx.Difficulty, 1, 4)
	c := k * (n + 1)
	answer := coefStr(k) + varPow("x", n+1)
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\int %d x^{%d} \\, dx\\) (use power rule for integration, omit \\(+C\\)).", c, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("∫%dx^%d dx = %d/(%d+1)·x^(%d+1)+C = %s+C.", c, n, c, n, n, answer),
	}
}

type integralSubstitutionGen struct{}

// sinPow renders the power of sine used in a substitution stem. A first power
// is the bare function, because \\sin(x)^{1} reads as machine output.
func sinPow(p int) string {
	if p == 1 {
		return "\\sin(x)"
	}
	return fmt.Sprintf("\\sin^{%d}(x)", p)
}

// coefTerm joins a coefficient to a rendered term, dropping the coefficient when
// it is one so the stem does not read as 1 t^{2}.
func coefTerm(c int, term string) string {
	if c == 1 {
		return term
	}
	return fmt.Sprintf("%d %s", c, term)
}

func (g *integralSubstitutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// One shape is not enough here: a single substitution family varies only in the
	// exponents, which is why this generator produced five stems. Each case below is
	// a different substitution with its answer derived from the parameters, not chosen.
	p := randRange(ctx.Difficulty, 1, 4)
	switch rand.Intn(5) {
	case 0:
		k := randRange(ctx.Difficulty, 1, 4)
		c := k * (p + 1)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int %d %s e^{%s} \\, dx\\).", c, texPow("x", p), texPow("x", p+1)),
			Answer:      fmt.Sprintf("%se^(x^%d)+C", coefStr(k), p+1),
			Explanation: fmt.Sprintf("Let u=x^%d, so du=%d x^%d dx. The integral becomes %d times the integral of e^u du, which is %de^u+C, i.e. %se^(x^%d)+C.", p+1, p+1, p, k, k, coefStr(k), p+1),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int %s \\sin(%s) \\, dx\\).", texPow("x", p), texPow("x", p+1)),
			Answer:      fmt.Sprintf("-cos(x^%d)/%d+C", p+1, p+1),
			Explanation: fmt.Sprintf("Let u=x^%d, so du=%d x^%d dx. The integral is the integral of sin(u) du over %d, which is -cos(u)/%d+C.", p+1, p+1, p, p+1, p+1),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int \\cos(x) %s \\, dx\\).", sinPow(p)),
			Answer:      fmt.Sprintf("sin^%d(x)/%d+C", p+1, p+1),
			Explanation: fmt.Sprintf("Let u=sin(x), so du=cos(x) dx. The integral becomes the integral of u^%d du, which is u^%d/%d+C.", p, p+1, p+1),
		}
	case 3:
		b := randRange(ctx.Difficulty, 1, 4)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int (x+%d)^{%d} \\, dx\\).", b, p+1),
			Answer:      fmt.Sprintf("(x+%d)^%d/%d+C", b, p+2, p+2),
			Explanation: fmt.Sprintf("Let u=x+%d, so du=dx. The integral becomes the integral of u^%d du, which is u^%d/%d+C.", b, p+1, p+2, p+2),
		}
	default:
		c := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\int \\frac{x}{x^{2}+%d} \\, dx\\).", c),
			Answer:      fmt.Sprintf("(1/2)ln(x^2+%d)+C", c),
			Explanation: fmt.Sprintf("Let u=x^2+%d, so du=2x dx and x dx = du/2. The integral becomes half the integral of du/u, which is (1/2)ln(x^2+%d)+C.", c, c),
		}
	}
}

type integralDefiniteGen struct{}

func (g *integralDefiniteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	n := rand.Intn(max(1, scale)) + 1
	lower := rand.Intn(2)
	upper := rand.Intn(max(1, scale*2)) + 1
	for upper <= lower {
		upper = rand.Intn(max(1, scale*2)) + 1
	}
	f := func(x int) int {
		return mathutil.IntPow(x, n+1) * a / (n + 1)
	}
	ans := f(upper) - f(lower)
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\int_{0}^{%d} %d x^{%d} \\, dx\\).", upper, a, n),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("∫%dx^%d dx = %d/(%d)·x^(%d). From 0 to %d: %d·%d^(%d)/%d - 0 = %d.", a, n, a, n+1, n+1, upper, a, upper, n+1, n+1, ans),
	}
}

type integralFTCGen struct{}

func (g *integralFTCGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The original case only varied the exponent of one integrand. The cases below
	// ask about the derivative of an accumulation function whose upper limit, lower
	// limit and integrand differ, so the answer has to be derived through the rule.
	p := randRange(ctx.Difficulty, 1, 5)
	switch rand.Intn(4) {
	case 0:
		c := randRange(ctx.Difficulty, 2, 5)
		b := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(F(x)=\\int_{%d}^{x} (%d %s + %d)\\, dt\\), what is \\(F'(x)\\)?", c, c, texPow("t", p), b),
			Answer:      fmt.Sprintf("%s+%d", coefStr(c)+varPow("x", p), b),
			Explanation: fmt.Sprintf("By the Fundamental Theorem of Calculus the lower limit is a constant, so it contributes nothing to the derivative. Differentiating the integrand at t=x gives %s+%d.", coefStr(c)+varPow("x", p), b),
		}
	case 1:
		q := randRange(ctx.Difficulty, 2, 3)
		// F(x) = integral from 0 to x^q of t^p dt, so F'(x) = q x^(q-1) * (x^q)^p.
		exp := q*p + q - 1
		coef := q
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(F(x)=\\int_{0}^{%s} %s \\, dt\\), what is \\(F'(x)\\)?", texPow("x", q), texPow("t", p)),
			Answer:      coefStr(coef) + varPow("x", exp),
			Explanation: fmt.Sprintf("Let u=x^%d, so du=%d x^%d dx. Then F(u)=u^%d/%d and F'(x)=(x^%d)^%d times %d x^%d, i.e. %s.", q, q, q-1, p+1, p+1, q, p, q, q-1, coefStr(coef)+varPow("x", exp)),
		}
	case 2:
		c := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(F(x)=\\int_{1}^{x} \\frac{%d}{t} \\, dt\\), what is \\(F'(x)\\)?", c),
			Answer:      fmt.Sprintf("%d/x", c),
			Explanation: fmt.Sprintf("The integrand evaluated at t=x is %d/x, and the lower limit is constant, so F'(x)=%d/x.", c, c),
		}
	default:
		q := 1 + rand.Intn(3)
		for q == p {
			q = 1 + rand.Intn(3)
		}
		c := randRange(ctx.Difficulty, 1, 4)
		// F(x) = integral of t^p + c t^q, differentiated term by term.
		answer := fmt.Sprintf("%s+%s", powTerm(p, p-1), powTerm(c*q, q-1))
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(F(x)=\\int_{0}^{x} (%s + %s)\\, dt\\), what is \\(F'(x)\\)?", texPow("t", p), coefTerm(c, texPow("t", q))),
			Answer:      answer,
			Explanation: fmt.Sprintf("Differentiating the integrand at t=x term by term gives %s+%s.", powTerm(p, p-1), powTerm(c*q, q-1)),
		}
	}
}

type integralAreaBetweenGen struct{}

func (g *integralAreaBetweenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Area between y = mx and y = 0 from x = 0 to x = L with m even
		m := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		L := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := (m * L * L) / 2
		return generator.Problem{
			Question: fmt.Sprintf("Find the area enclosed between \\(y = %d x\\) and the x-axis \\(y = 0\\) from \\(x = 0\\) to \\(x = %d\\).", m, L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for the area between curves: \\(A = \\int_{a}^{b} (f(x) - g(x))\\,dx\\).\\n"+
				"Subgoal 2: Substitute parameters \\(f(x) = %d x\\), \\(g(x) = 0\\), \\(a = 0\\) and \\(b = %d\\).\\n"+
				"Subgoal 3: Evaluate the definite integral: \\(\\int_{0}^{%d} %d x\\,dx = [%d x^2 / 2]_{0}^{%d} = %d\\).", m, L, L, m, m, L, ans),
		}
	case 1:
		// Area between parabola y = r^2 - x^2 and y = 0 on [-r, r] with r = 3k
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		r := 3 * k
		H := r * r
		ans := (4 * r * r * r) / 3
		return generator.Problem{
			Question: fmt.Sprintf("Find the area bounded by the parabola \\(y = %d - x^2\\) and the x-axis \\(y = 0\\).", H),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: find roots \\(\\pm %d\\) and integrate \\(A = \\int_{-%d}^{%d} (%d - x^2)\\,dx\\).\\n"+
				"Subgoal 2: Use symmetry: \\(A = 2 \\int_{0}^{%d} (%d - x^2)\\,dx\\).\\n"+
				"Subgoal 3: Evaluate and simplify: \\(2 [%d x - x^3 / 3]_{0}^{%d} = 2 (%d - %d) = %d\\).", r, r, r, H, r, H, H, r, H*r, (r*r*r)/3, ans),
		}
	case 2:
		// Area between y = kx and y = x^2 from x = 0 to x = k with k = 6m
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 2)
		k := 6 * m
		ans := (k * k * k) / 6
		return generator.Problem{
			Question: fmt.Sprintf("Find the area of the region bounded by \\(y = %d x\\) and \\(y = x^2\\).", k),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: intersections occur at \\(x = 0\\) and \\(x = %d\\), so \\(A = \\int_{0}^{%d} (%d x - x^2)\\,dx\\).\\n"+
				"Subgoal 2: Find antiderivative: \\([%d x^2 / 2 - x^3 / 3]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Evaluate and simplify: \\(\\frac{%d \\cdot %d}{2} - \\frac{%d}{3} = %d\\).", k, k, k, k, k, k, k*k, k*k*k, ans),
		}
	case 3:
		// Area between y = 3kx^2 and y = 0 from x = 0 to x = L
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		L := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k * L * L * L
		coeff := 3 * k
		return generator.Problem{
			Question: fmt.Sprintf("Find the area under the curve \\(y = %d x^2\\) from \\(x = 0\\) to \\(x = %d\\).", coeff, L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(A = \\int_{0}^{%d} %d x^2\\,dx\\).\\n"+
				"Subgoal 2: Find the antiderivative: \\([%d x^3]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Evaluate and simplify: \\(%d \\cdot %d^3 = %d\\).", L, coeff, k, L, k, L, ans),
		}
	default:
		// Area between two parallel lines y = mx + b and y = mx from x = 0 to x = L
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		L := mathutil.ScaleInt(ctx.Difficulty, 1, 7)
		ans := b * L
		return generator.Problem{
			Question: fmt.Sprintf("Find the area between the lines \\(y = %d x + %d\\) and \\(y = %d x\\) over the interval [%d, %d].", m, b, m, 0, L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(A = \\int_{a}^{b} (f(x) - g(x))\\,dx\\).\\n"+
				"Subgoal 2: Substitute \\(f(x) - g(x) = (%d x + %d) - %d x = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(\\int_{0}^{%d} %d\\,dx = %d \\cdot %d = %d\\).", m, b, m, b, L, b, b, L, ans),
		}
	}
}

type integralVolumeGen struct{}

func (g *integralVolumeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The disc method is one formula, so the variety has to come from the curve, the
	// limits and the axis of revolution rather than from picking a row. Each case
	// derives V = pi * integral of radius squared from the parameters it drew.
	reduce := func(num, den int) string {
		d := mathutil.GCD(num, den)
		num, den = num/d, den/d
		if den == 1 {
			return fmt.Sprintf("%dpi", num)
		}
		return fmt.Sprintf("%dpi/%d", num, den)
	}
	switch rand.Intn(4) {
	case 0:
		p := randRange(ctx.Difficulty, 1, 3)
		c := randRange(ctx.Difficulty, 1, 3)
		b := randRange(ctx.Difficulty, 1, 5)
		num := c * c * intPow(b, 2*p+1)
		den := 2*p + 1
		return generator.Problem{
			Question:    fmt.Sprintf("Find the volume when \\(y=%s\\) from \\(x=0\\) to %d is revolved around the \\(x\\)-axis.", coefStr(c)+varPow("x", p), b),
			Answer:      reduce(num, den),
			Explanation: fmt.Sprintf("The disc method gives V = pi times the integral of (%s)^2 from 0 to %d, i.e. %d pi x^(%d) integrated, which is %d pi b^(%d)/(%d) with b = %d.", coefStr(c)+varPow("x", p), b, c*c, 2*p, c*c, 2*p+1, 2*p+1, b),
		}
	case 1:
		r := randRange(ctx.Difficulty, 1, 5)
		l := randRange(ctx.Difficulty, 1, 6)
		return generator.Problem{
			Question:    fmt.Sprintf("Find the volume when \\(y=%d\\) from \\(x=0\\) to %d is revolved around the \\(x\\)-axis.", r, l),
			Answer:      fmt.Sprintf("%dpi", r*r*l),
			Explanation: fmt.Sprintf("A constant radius revolving about the x-axis sweeps a cylinder of radius %d and length %d, so V = pi r^2 L = %dpi.", r, l, r*r*l),
		}
	case 2:
		c := randRange(ctx.Difficulty, 1, 4)
		b := randRange(ctx.Difficulty, 1, 6)
		num := c * c * b * b
		curve := "\\sqrt{x}"
		if c != 1 {
			curve = fmt.Sprintf("%d\\sqrt{x}", c)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Find the volume when \\(y=%s\\) from \\(x=0\\) to %d is revolved around the \\(x\\)-axis.", curve, b),
			Answer:      reduce(num, 2),
			Explanation: fmt.Sprintf("Squaring the radius gives %dx, so V = pi times %d times b^2/2 with b = %d, i.e. %d pi/2 reduced.", c*c, c*c, b, num),
		}
	default:
		p := randRange(ctx.Difficulty, 1, 2)
		c := randRange(ctx.Difficulty, 1, 2)
		a := randRange(ctx.Difficulty, 1, 2)
		b := a + randRange(ctx.Difficulty, 1, 4)
		num := c * c * (intPow(b, 2*p+1) - intPow(a, 2*p+1))
		den := 2*p + 1
		return generator.Problem{
			Question:    fmt.Sprintf("Find the volume when \\(y=%s\\) from \\(x=%d\\) to %d is revolved around the \\(x\\)-axis.", coefStr(c)+varPow("x", p), a, b),
			Answer:      reduce(num, den),
			Explanation: fmt.Sprintf("V = pi times the integral of (%s)^2 from %d to %d, which is %d pi (b^(%d) - a^(%d))/(%d) with a = %d and b = %d.", coefStr(c)+varPow("x", p), a, b, c*c, 2*p+1, 2*p+1, 2*p+1, a, b),
		}
	}
}

type derivImplicitGen struct{}

func (g *derivImplicitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	r := rand.Intn(max(1, scale*2)) + 3
	a := rand.Intn(max(1, scale*2)) + 1
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
		Question:    fmt.Sprintf("Find \\(dy/dx\\) at \\((%d,%d)\\) for \\(x^{2}+y^{2}=%d\\).", a, bSqrt, r*r),
		Answer:      answer,
		Explanation: fmt.Sprintf("2x+2y·dy/dx=0 → dy/dx=-x/y. At (%d,%d): dy/dx=-%d/%d=%s.", a, bSqrt, a, bSqrt, answer),
	}
}

type derivRelatedRatesGen struct{}

func (g *derivRelatedRatesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	dr := rand.Intn(max(1, scale*2)) + 1
	r := rand.Intn(max(1, scale*2)) + 3
	ans := 2 * 3 * r * dr
	answer := fmt.Sprintf("%dpi", ans)
	return generator.Problem{
		Question:    fmt.Sprintf("A circle's radius grows at %d cm/s. How fast is area growing when \\(r=%d\\)?", dr, r),
		Answer:      answer,
		Explanation: fmt.Sprintf("A=πr², dA/dt=2πr·dr/dt=2π(%d)(%d)=%dpi cm²/s.", r, dr, ans),
	}
}

type integralPartsGen struct{}

func (g *integralPartsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Find \\(\\int x e^{x} \\, dx\\).", "xe^x-e^x+C", "Integration by parts: let u=x, dv=e^x dx → du=dx, v=e^x. ∫x·e^x dx = x·e^x - ∫e^x dx = x·e^x - e^x + C."},
		{"Find \\(\\int 2x e^{x} \\, dx\\).", "2xe^x-2e^x+C", "Let u=2x, dv=e^x dx → du=2dx, v=e^x. ∫2x·e^x dx = 2x·e^x - ∫2e^x dx = 2x·e^x - 2e^x + C."},
		{"For \\(\\int x e^{x} \\, dx\\) with \\(u=x\\) and \\(dv=e^{x}dx\\), what is \\(du\\)?", "dx", "Differentiate u=x to get du=dx; then v=e^x, and the parts formula gives x·e^x - e^x + C."},
		{"Find \\(\\int x e^{2x} \\, dx\\).", "(2x-1)e^(2x)/4+C", "Let u=x, dv=e^(2x)dx → du=dx, v=e^(2x)/2. ∫x·e^(2x)dx = x·e^(2x)/2 - ∫e^(2x)/2 dx = (2x-1)e^(2x)/4 + C."},
	}
	tableHard := []entry{
		{"Find \\(\\int x^{2} e^{x} \\, dx\\).", "x^2e^x-2xe^x+2e^x+C", "Apply parts twice: ∫x²e^x dx = x²e^x - ∫2xe^x dx = x²e^x - 2xe^x + 2e^x + C."},
		{"Find \\(\\int x \\ln(x) \\, dx\\).", "x^2ln(x)/2-x^2/4+C", "Let u=ln(x), dv=x dx → du=dx/x, v=x²/2. ∫x·ln(x)dx = x²ln(x)/2 - ∫x/2 dx = x²ln(x)/2 - x²/4 + C."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralPartialFractionsGen struct{}

func (g *integralPartialFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Find \\(\\int \\frac{1}{x^{2}-1} \\, dx\\).", "(1/2)ln|x-1|-(1/2)ln|x+1|+C", "Partial fractions: 1/(x²-1) = 1/2·(1/(x-1) - 1/(x+1)). Integrate: (1/2)ln|x-1| - (1/2)ln|x+1| + C."},
		{"Find \\(\\int \\frac{1}{x(x+1)} \\, dx\\).", "ln|x|-ln|x+1|+C", "Partial fractions: 1/(x(x+1)) = 1/x - 1/(x+1). Integrate: ln|x| - ln|x+1| + C."},
		{"Find \\(\\int \\frac{1}{x(x-1)} \\, dx\\).", "ln|x-1|-ln|x|+C", "Partial fractions: 1/(x(x-1)) = 1/(x-1) - 1/x. Integrate: ln|x-1| - ln|x| + C."},
		{"In \\(\\frac{1}{x^{2}-1} = \\frac{A}{x-1}+\\frac{B}{x+1}\\), what is \\(A\\)?", "1/2", "Write 1 = A(x+1)+B(x-1). Setting x=1 gives 1 = 2A, so A = 1/2."},
	}
	tableHard := []entry{
		{"Find \\(\\int \\frac{2}{x^{2}-4} \\, dx\\).", "(1/2)ln|x-2|-(1/2)ln|x+2|+C", "Partial fractions: 2/(x²-4) = 1/2·(1/(x-2) - 1/(x+2)). Integrate: (1/2)ln|x-2| - (1/2)ln|x+2| + C."},
		{"In \\(\\frac{2}{x^{2}-4} = \\frac{A}{x-2}+\\frac{B}{x+2}\\), what is \\(B\\)?", "-1/2", "Write 2 = A(x+2)+B(x-2). Setting x=-2 gives 2 = -4B, so B = -1/2."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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

// --- Derivative generators ---

type rolleGen struct{}

func (g *rolleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Quadratic f(x) = x^2 - (r1+r2)x + r1*r2 on [r1, r2] with r1, r2 same parity
		r1 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		span := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		r2 := r1 + span
		b := r1 + r2
		cVal := r1 * r2
		cMid := b / 2
		return generator.Problem{
			Question: fmt.Sprintf("For the function \\(f(x) = x^2 - %d x + %d\\) on the interval \\([%d, %d]\\), find the value \\(c \\in (%d, %d)\\) guaranteed by Rolle's Theorem such that \\(f'(c) = 0\\).", b, cVal, r1, r2, r1, r2),
			Answer:   fmt.Sprintf("%d", cMid),
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: if \\(f\\) is continuous on \\([a, b]\\), differentiable on \\((a, b)\\) and \\(f(a) = f(b)\\), there exists \\(c \\in (a, b)\\) with \\(f'(c) = 0\\).\\n"+
				"Subgoal 2: Compute the derivative \\(f'(x) = 2x - %d\\) and set it to zero.\\n"+
				"Subgoal 3: Solve for \\(c\\): \\(2c - %d = 0 \\implies c = %d\\).", b, b, cMid),
		}
	case 1:
		// Inverted quadratic f(x) = -x^2 + 2kx on [0, 2k]
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		b := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = -x^2 + %d x\\) on \\([0, %d]\\), find the critical point \\(c\\) guaranteed by Rolle's Theorem.", b, b),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Subgoal 1: Confirm the boundary condition \\(f(0) = 0\\) and \\(f(%d) = -(%d)^2 + %d(%d) = 0\\).\\n"+
				"Subgoal 2: Compute \\(f'(x) = -2x + %d\\) and equate to zero.\\n"+
				"Subgoal 3: Solve \\(-2c + %d = 0 \\implies c = %d\\).", b, b, b, b, b, b, k),
		}
	case 2:
		// Shifted quadratic f(x) = -(x - c0)^2 + d^2 on [c0-d, c0+d]
		c0 := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		d := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		aEnd := c0 - d
		bEnd := c0 + d
		lin := 2 * c0
		constTerm := d*d - c0*c0
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = -x^2 + %d x + (%d)\\) on \\([%d, %d]\\), find the point \\(c\\) guaranteed by Rolle's Theorem.", lin, constTerm, aEnd, bEnd),
			Answer:   fmt.Sprintf("%d", c0),
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: since \\(f(%d) = f(%d) = 0\\), a point \\(c \\in (%d, %d)\\) satisfies \\(f'(c) = 0\\).\\n"+
				"Subgoal 2: Differentiate: \\(f'(x) = -2x + %d\\).\\n"+
				"Subgoal 3: Set \\(-2c + %d = 0\\) and solve for \\(c = %d\\).", aEnd, bEnd, aEnd, bEnd, lin, lin, c0),
		}
	case 3:
		// Value of derivative guaranteed by Rolle
		x1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x2 := x1 + mathutil.ScaleInt(ctx.Difficulty, 2, 6)
		return generator.Problem{
			Question: fmt.Sprintf("If a differentiable function \\(f(x)\\) satisfies \\(f(%d) = f(%d)\\) on \\([%d, %d]\\), what is the value of \\(f'(c)\\) at the point \\(c\\) guaranteed by Rolle's Theorem?", x1, x2, x1, x2),
			Answer:   "0",
			Explanation: "Subgoal 1: State the conclusion of Rolle's Theorem: the derivative must vanish at some interior point \\(c \\in (a, b)\\).\\n" +
				"Subgoal 2: Identify the critical condition: \\(f'(c) = 0\\).\\n" +
				"Subgoal 3: Conclude that the derivative value is 0.",
		}
	default:
		// Scaled quadratic f(x) = 2x^2 - 4cx + K on [c-d, c+d]
		c := mathutil.ScaleInt(ctx.Difficulty, 2, 7)
		d := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		x1 := c - d
		x2 := c + d
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = 2x^2 - %d x\\) on \\([%d, %d]\\), find the value \\(c\\) where \\(f'(c) = 0\\) as guaranteed by Rolle's Theorem.", 4*c, x1, x2),
			Answer:   fmt.Sprintf("%d", c),
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: since \\(f(%d) = f(%d)\\), there exists \\(c \\in (%d, %d)\\) with \\(f'(c) = 0\\).\\n"+
				"Subgoal 2: Differentiate \\(f(x)\\): \\(f'(x) = 4x - %d\\).\\n"+
				"Subgoal 3: Solve \\(4c - %d = 0 \\implies c = %d\\).", x1, x2, x1, x2, 4*c, 4*c, c),
		}
	}
}

type mvtGen struct{}

func (g *mvtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Quadratic f(x) = x^2 on [a, b] with same parity
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		span := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := a + span
		cMid := (a + b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^2\\) on the interval \\([%d, %d]\\), find the value \\(c \\in (%d, %d)\\) satisfying the Mean Value Theorem.", a, b, a, b),
			Answer:   fmt.Sprintf("%d", cMid),
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem formula: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\\n"+
				"Subgoal 2: Compute the average rate of change: \\(\\frac{(%d)^2 - (%d)^2}{%d - %d} = \\frac{%d - %d}{%d} = %d\\).\\n"+
				"Subgoal 3: Equate to \\(f'(c) = 2c = %d\\) and solve for \\(c = %d\\).", b, a, b, a, b*b, a*a, b-a, a+b, a+b, cMid),
		}
	case 1:
		// General quadratic f(x) = kx^2 + mx on [a, b]
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		a := mathutil.ScaleInt(ctx.Difficulty, 0, 5)
		span := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		b := a + span
		cMid := (a + b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = %d x^2 + %d x\\) on \\([%d, %d]\\), find the point \\(c\\) guaranteed by the Mean Value Theorem.", k, m, a, b),
			Answer:   fmt.Sprintf("%d", cMid),
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\\n"+
				"Subgoal 2: For any quadratic function, the instantaneous rate of change equals the average rate of change at the midpoint \\(c = \\frac{a + b}{2}\\).\\n"+
				"Subgoal 3: Calculate the midpoint: \\(c = \\frac{%d + %d}{2} = %d\\).", a, b, cMid),
		}
	case 2:
		// Average rate of change on [0, b]
		b := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		ans := b + k
		return generator.Problem{
			Question: fmt.Sprintf("Find the average rate of change \\(\\frac{f(b) - f(a)}{b - a}\\) of \\(f(x) = x^2 + %d x\\) on \\([0, %d]\\).", k, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the average rate of change formula: \\(\\frac{f(b) - f(a)}{b - a}\\).\\n"+
				"Subgoal 2: Evaluate \\(f(%d) = (%d)^2 + %d(%d) = %d\\) and \\(f(0) = 0\\).\\n"+
				"Subgoal 3: Divide by interval width: \\(\\frac{%d - 0}{%d} = %d\\).", b, b, k, b, b*b+k*b, b*b+k*b, b, ans),
		}
	case 3:
		// Kinematics velocity interpretation
		hours := 2
		if ctx.Difficulty > 0.5 {
			hours = 3
		}
		speed := mathutil.ScaleInt(ctx.Difficulty, 45, 80)
		miles := speed * hours
		return generator.Problem{
			Question: fmt.Sprintf("A car travels %d miles in %d hours. By the Mean Value Theorem, at what instantaneous speed (in mph) was the car traveling at some moment during the journey?", miles, hours),
			Answer:   fmt.Sprintf("%d", speed),
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem in kinematics: \\(v_{\\text{avg}} = \\frac{s(t_2) - s(t_1)}{t_2 - t_1} = s'(c)\\).\\n"+
				"Subgoal 2: Substitute distance \\(\\Delta s = %d\\) and time \\(\\Delta t = %d\\): \\(\\frac{%d}{%d}\\).\\n"+
				"Subgoal 3: Divide to obtain the instantaneous speed %d mph.", miles, hours, miles, hours, speed),
		}
	default:
		// Cubic f(x) = x^3 on [0, 3k] derivative value f'(c)
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 3 * k
		derivVal := b * b
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^3\\) on \\([0, %d]\\), what is the value of the derivative \\(f'(c)\\) guaranteed by the Mean Value Theorem?", b),
			Answer:   fmt.Sprintf("%d", derivVal),
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\\n"+
				"Subgoal 2: Evaluate \\(f(%d) = (%d)^3 = %d\\) and \\(f(0) = 0\\).\\n"+
				"Subgoal 3: Calculate \\(f'(c) = \\frac{%d}{%d} = %d\\).", b, b, b*b*b, b*b*b, b, derivVal),
		}
	}
}

type cauchyMVTGen struct{}

func (g *cauchyMVTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// f(x) = x^2, g(x) = x on [a, b]
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		span := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		b := a + span
		cMid := (a + b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^2\\) and \\(g(x) = x\\) on \\([%d, %d]\\), find the value \\(c\\) satisfying Cauchy's Mean Value Theorem: \\(\\frac{f(b)-f(a)}{g(b)-g(a)} = \\frac{f'(c)}{g'(c)}\\).", a, b),
			Answer:   fmt.Sprintf("%d", cMid),
			Explanation: fmt.Sprintf("Subgoal 1: State Cauchy's Mean Value Theorem: \\(\\frac{f(b) - f(a)}{g(b) - g(a)} = \\frac{f'(c)}{g'(c)}\\).\\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 2x\\) and \\(g'(x) = 1\\), giving quotient \\(\\frac{2c}{1} = \\frac{(%d)^2 - (%d)^2}{%d - %d} = %d\\).\\n"+
				"Subgoal 3: Solve for \\(c\\): \\(2c = %d \\implies c = %d\\).", a, b, a, b, a+b, a+b, cMid),
		}
	case 1:
		// f(x) = x^3, g(x) = x^2 on [0, 3k] -> c = 2k
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := 3 * k
		cVal := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^3\\) and \\(g(x) = x^2\\) on \\([0, %d]\\), find \\(c\\) satisfying Cauchy's Mean Value Theorem.", b),
			Answer:   fmt.Sprintf("%d", cVal),
			Explanation: fmt.Sprintf("Subgoal 1: State Cauchy's Mean Value Theorem: \\(\\frac{f(b) - f(a)}{g(b) - g(a)} = \\frac{f'(c)}{g'(c)}\\).\\n"+
				"Subgoal 2: Substitute derivatives \\(f'(c) = 3c^2\\) and \\(g'(c) = 2c\\), so \\(\\frac{3c^2}{2c} = \\frac{3}{2}c\\).\\n"+
				"Subgoal 3: Compute quotient \\(\\frac{(%d)^3 - 0}{(%d)^2 - 0} = %d\\) and solve \\(\\frac{3}{2}c = %d \\implies c = %d\\).", b, b, b, b, cVal),
		}
	case 2:
		// Cauchy quotient evaluation for f(x) = kx^2 and g(x) = x on [0, b]
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := mathutil.ScaleInt(ctx.Difficulty, 2, 7)
		ans := k * b
		return generator.Problem{
			Question: fmt.Sprintf("Compute the Cauchy quotient \\(\\frac{f(b)-f(a)}{g(b)-g(a)}\\) for \\(f(x) = %d x^2\\) and \\(g(x) = x\\) on \\([0, %d]\\).", k, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Cauchy difference quotient: \\(\\frac{f(b) - f(a)}{g(b) - g(a)}\\).\\n"+
				"Subgoal 2: Evaluate \\(f(%d) = %d(%d)^2 = %d\\) and \\(g(%d) = %d\\).\\n"+
				"Subgoal 3: Divide: \\(\\frac{%d - 0}{%d - 0} = %d\\).", b, k, b, k*b*b, b, b, k*b*b, b, ans),
		}
	case 3:
		// Cauchy ratio for f(x) = a*x, g(x) = b*x
		m2 := mathutil.ScaleInt(ctx.Difficulty, 2, 4)
		q := mathutil.ScaleInt(ctx.Difficulty, 2, 5)
		m1 := m2 * q
		return generator.Problem{
			Question: fmt.Sprintf("For linear functions \\(f(x) = %d x\\) and \\(g(x) = %d x\\) on any interval \\([a, b]\\), what is the value of \\(\\frac{f'(c)}{g'(c)}\\)?", m1, m2),
			Answer:   fmt.Sprintf("%d", q),
			Explanation: fmt.Sprintf("Subgoal 1: Differentiate the linear functions: \\(f'(x) = %d\\) and \\(g'(x) = %d\\).\\n"+
				"Subgoal 2: Form the derivative ratio: \\(\\frac{f'(c)}{g'(c)} = \\frac{%d}{%d}\\).\\n"+
				"Subgoal 3: Divide to obtain %d.", m1, m2, m1, m2, q),
		}
	default:
		// Denominator reduction when g(x) = x
		x1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x2 := x1 + mathutil.ScaleInt(ctx.Difficulty, 2, 7)
		width := x2 - x1
		return generator.Problem{
			Question: fmt.Sprintf("In Cauchy's Mean Value Theorem on \\([%d, %d]\\) with \\(g(x) = x\\), what is the denominator value \\(g(b) - g(a)\\)?", x1, x2),
			Answer:   fmt.Sprintf("%d", width),
			Explanation: fmt.Sprintf("Subgoal 1: Substitute \\(g(x) = x\\) into the denominator \\(g(b) - g(a)\\).\\n"+
				"Subgoal 2: Compute \\(b - a = %d - %d\\).\\n"+
				"Subgoal 3: Evaluate the difference: %d.", x2, x1, width),
		}
	}
}

type fermatGen struct{}

func (g *fermatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Parabola vertex local extremum c
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		b := 2 * a * h
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 10)
		cConst := a*h*h + k
		return generator.Problem{
			Question: fmt.Sprintf("Find the critical point \\(c\\) where \\(f(x) = %d x^2 - %d x + %d\\) attains its local extremum according to Fermat's Theorem.", a, b, cConst),
			Answer:   fmt.Sprintf("%d", h),
			Explanation: fmt.Sprintf("Subgoal 1: State Fermat's Theorem: differentiable local extrema satisfy \\(f'(c) = 0\\).\\n"+
				"Subgoal 2: Differentiate: \\(f'(x) = 2(%d)x - %d = %d x - %d\\).\\n"+
				"Subgoal 3: Set \\(%d c - %d = 0\\) and solve for \\(c = %d\\).", a, b, 2*a, b, 2*a, b, h),
		}
	case 1:
		// Value of derivative at local extremum
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 12)
		return generator.Problem{
			Question: fmt.Sprintf("If a differentiable function \\(f(x)\\) attains a local extremum at \\(x = %d\\), what is the value of \\(f'(%d)\\) by Fermat's Theorem?", x0, x0),
			Answer:   "0",
			Explanation: "Subgoal 1: State Fermat's Theorem on stationary points: if \\(f\\) has a local extremum at \\(c\\) and \\(f'(c)\\) exists, then \\(f'(c) = 0\\).\\n" +
				"Subgoal 2: Identify that \\(f\\) is differentiable at the local extremum.\\n" +
				"Subgoal 3: Conclude that \\(f'(c) = 0\\).",
		}
	case 2:
		// Smaller critical point of cubic
		h1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		h2 := h1 + mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 3 * (h1 + h2)
		cVal := 6 * h1 * h2
		return generator.Problem{
			Question: fmt.Sprintf("Find the smaller critical point of \\(f(x) = 2x^3 - %d x^2 + %d x\\) where a local extremum may occur.", b, cVal),
			Answer:   fmt.Sprintf("%d", h1),
			Explanation: fmt.Sprintf("Subgoal 1: Differentiate \\(f(x)\\): \\(f'(x) = 6x^2 - %d x + %d = 6(x - %d)(x - %d)\\).\\n"+
				"Subgoal 2: Set \\(f'(x) = 0\\) to find stationary points \\(x = %d\\) and \\(x = %d\\).\\n"+
				"Subgoal 3: Select the smaller critical point: %d.", 2*b, cVal, h1, h2, h1, h2, h1),
		}
	case 3:
		// Maximum value at stationary point for f(x) = -(x - h)^2 + K
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 7)
		k := mathutil.ScaleInt(ctx.Difficulty, 5, 25)
		return generator.Problem{
			Question: fmt.Sprintf("The function \\(f(x) = -(x - %d)^2 + %d\\) attains its local maximum at its stationary point. What is this maximum value?", h, k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Subgoal 1: Apply Fermat's Theorem: the stationary point occurs where \\(f'(x) = -2(x - %d) = 0 \\implies x = %d\\).\\n"+
				"Subgoal 2: Substitute \\(x = %d\\) into \\(f(x)\\): \\(f(%d) = -(%d - %d)^2 + %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(0 + %d = %d\\).", h, h, h, h, h, h, k, k, k),
		}
	default:
		// Count of stationary points for non-degenerate quadratic
		return generator.Problem{
			Question: "How many stationary points (where \\(f'(c) = 0\\)) does a non-degenerate quadratic polynomial \\(f(x) = ax^2 + bx + c\\) possess?",
			Answer:   "1",
			Explanation: "Subgoal 1: Compute the derivative of a quadratic: \\(f'(x) = 2ax + b\\).\\n" +
				"Subgoal 2: Set \\(2ax + b = 0\\).\\n" +
				"Subgoal 3: Since \\(a \\ne 0\\), there is exactly one unique solution \\(x = -\\frac{b}{2a}\\), so the count is 1.",
		}
	}
}

type differenceQuotientGen struct{}

func (g *differenceQuotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// f(x) = a x^2 + b x
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		fx0 := a*x0*x0 + b*x0
		x1 := x0 + h
		fx1 := a*x1*x1 + b*x1
		ans := (fx1 - fx0) / h
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = %d x^2 + %d x\\), compute the difference quotient \\(\\frac{f(%d + %d) - f(%d)}{%d}\\).", a, b, x0, h, x0, h),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the difference quotient formula: \\(\\frac{f(x_0 + h) - f(x_0)}{h}\\).\\n"+
				"Subgoal 2: Evaluate \\(f(%d) = %d\\) and \\(f(%d) = %d\\), so \\(f(%d) - f(%d) = %d - %d = %d\\).\\n"+
				"Subgoal 3: Divide by \\(h = %d\\): \\(\\frac{%d}{%d} = %d\\).", x1, fx1, x0, fx0, x1, x0, fx1, fx0, fx1-fx0, h, fx1-fx0, h, ans),
		}
	case 1:
		// Linear function f(x) = m x + b
		m := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 10)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		return generator.Problem{
			Question: fmt.Sprintf("For linear function \\(f(x) = %d x + %d\\), compute the difference quotient \\(\\frac{f(%d + %d) - f(%d)}{%d}\\).", m, b, x0, h, x0, h),
			Answer:   fmt.Sprintf("%d", m),
			Explanation: fmt.Sprintf("Subgoal 1: State the difference quotient for a linear function: \\(\\frac{f(x + h) - f(x)}{h}\\).\\n"+
				"Subgoal 2: Substitute \\(f(x + h) - f(x) = [%d(x + h) + %d] - [%d x + %d] = %d h\\).\\n"+
				"Subgoal 3: Divide by \\(h\\): \\(\\frac{%d h}{h} = %d\\).", m, b, m, b, m, m, m),
		}
	case 2:
		// f(x) = x^3 at x0 = 1
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		ans := 3 + 3*h + h*h
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^3\\), compute the difference quotient \\(\\frac{f(1 + %d) - f(1)}{%d}\\).", h, h),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Expand \\(f(1 + h) = (1 + h)^3 = 1 + 3h + 3h^2 + h^3\\).\\n"+
				"Subgoal 2: Subtract \\(f(1) = 1\\) to get \\(3h + 3h^2 + h^3\\), and divide by \\(h\\) to obtain \\(3 + 3h + h^2\\).\\n"+
				"Subgoal 3: Substitute \\(h = %d\\): \\(3 + 3(%d) + (%d)^2 = 3 + %d + %d = %d\\).", h, h, h, 3*h, h*h, ans),
		}
	case 3:
		// Limit of difference quotient (instantaneous derivative)
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := 2 * a * x0
		return generator.Problem{
			Question: fmt.Sprintf("What is the limit of the difference quotient \\(\\lim_{h \\to 0} \\frac{f(%d + h) - f(%d)}{h}\\) for \\(f(x) = %d x^2\\)?", x0, x0, a),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Recognize that the limit of the difference quotient defines the derivative: \\(f'(x_0) = \\lim_{h \\to 0} \\frac{f(x_0 + h) - f(x_0)}{h}\\).\\n"+
				"Subgoal 2: Differentiate \\(f(x) = %d x^2\\): \\(f'(x) = 2(%d)x = %d x\\).\\n"+
				"Subgoal 3: Evaluate at \\(x = %d\\): \\(f'(%d) = %d(%d) = %d\\).", a, a, 2*a, x0, x0, 2*a, x0, ans),
		}
	default:
		// Quadratic difference quotient at x0 = 0
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := a*h + b
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = %d x^2 + %d x\\), compute \\(\\frac{f(%d) - f(0)}{%d}\\).", a, b, h, h),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Substitute \\(x = %d\\) into \\(f(x)\\): \\(f(%d) = %d(%d)^2 + %d(%d) = %d\\).\\n"+
				"Subgoal 2: Subtract \\(f(0) = 0\\), giving numerator %d.\\n"+
				"Subgoal 3: Divide by %d: \\(\\frac{%d}{%d} = %d\\).", h, h, a, h, b, h, a*h*h+b*h, a*h*h+b*h, h, a*h*h+b*h, h, ans),
		}
	}
}

type differentialGen struct{}

func (g *differentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Differential dy = f'(x) dx for y = a x^2
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		dx := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		dy := 2 * a * x0 * dx
		return generator.Problem{
			Question: fmt.Sprintf("If \\(y = %d x^2\\) at \\(x = %d\\), compute the differential \\(dy\\) for \\(dx = %d\\).", a, x0, dx),
			Answer:   fmt.Sprintf("%d", dy),
			Explanation: fmt.Sprintf("Subgoal 1: State the differential definition: \\(dy = f'(x) dx\\).\\n"+
				"Subgoal 2: Compute derivative \\(f'(x) = 2(%d)x = %d x\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d\\).\\n"+
				"Subgoal 3: Multiply by \\(dx = %d\\): \\(dy = (%d)(%d) = %d\\).", a, 2*a, x0, x0, 2*a*x0, dx, 2*a*x0, dx, dy),
		}
	case 1:
		// Differential for y = a x^3
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		dx := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		dy := 3 * a * x0 * x0 * dx
		return generator.Problem{
			Question: fmt.Sprintf("If \\(y = %d x^3\\) at \\(x = %d\\), compute the differential \\(dy\\) for \\(dx = %d\\).", a, x0, dx),
			Answer:   fmt.Sprintf("%d", dy),
			Explanation: fmt.Sprintf("Subgoal 1: State the differential formula: \\(dy = f'(x) dx\\).\\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 3(%d)x^2 = %d x^2\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d(%d)^2 = %d\\).\\n"+
				"Subgoal 3: Multiply by \\(dx = %d\\): \\(dy = (%d)(%d) = %d\\).", a, 3*a, x0, x0, 3*a, x0, 3*a*x0*x0, dx, 3*a*x0*x0, dx, dy),
		}
	case 2:
		// Linear approximation change Delta y approx f'(x) dx for f(x) = a x^2 + b x
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		dx := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		dy := (2*a*x0 + b) * dx
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = %d x^2 + %d x\\) at \\(x = %d\\), compute the differential approximation \\(dy\\) for \\(dx = %d\\).", a, b, x0, dx),
			Answer:   fmt.Sprintf("%d", dy),
			Explanation: fmt.Sprintf("Subgoal 1: State the linear approximation differential: \\(dy = f'(x) dx\\).\\n"+
				"Subgoal 2: Compute \\(f'(x) = 2(%d)x + %d = %d x + %d\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d\\).\\n"+
				"Subgoal 3: Multiply by \\(dx = %d\\): \\(dy = (%d)(%d) = %d\\).", a, b, 2*a, b, x0, x0, 2*a*x0+b, dx, 2*a*x0+b, dx, dy),
		}
	case 3:
		// Ratio of differentials dy / dx = f'(x)
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := 2*a*x0 + b
		return generator.Problem{
			Question: fmt.Sprintf("What is the ratio of differentials \\(\\frac{dy}{dx}\\) for \\(y = %d x^2 + %d x\\) evaluated at \\(x = %d\\)?", a, b, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Recall that the ratio of differentials equals the derivative: \\(\\frac{dy}{dx} = f'(x)\\).\\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 2(%d)x + %d = %d x + %d\\).\\n"+
				"Subgoal 3: Evaluate at \\(x = %d\\): \\(f'(%d) = %d(%d) + %d = %d\\).", a, b, 2*a, b, x0, x0, 2*a, x0, b, ans),
		}
	default:
		// Differential for affine function y = m x + c
		m := mathutil.ScaleInt(ctx.Difficulty, 2, 9)
		dx := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		dy := m * dx
		return generator.Problem{
			Question: fmt.Sprintf("If \\(y = %d x + 7\\), compute \\(dy\\) for \\(dx = %d\\).", m, dx),
			Answer:   fmt.Sprintf("%d", dy),
			Explanation: fmt.Sprintf("Subgoal 1: State the differential rule: \\(dy = f'(x) dx\\).\\n"+
				"Subgoal 2: Since \\(f(x) = %d x + 7\\), its derivative is constant: \\(f'(x) = %d\\).\\n"+
				"Subgoal 3: Multiply by \\(dx = %d\\): \\(dy = %d \\times %d = %d\\).", m, m, dx, m, dx, dy),
		}
	}
}

type convexityGen struct{}

func (g *convexityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The answer is a classification rather than an expression, so the question
	// is multiple choice. A symbolic concept renders a maths field, and a learner
	// who types "convex" into it is marked wrong (report 118).
	const convex, concave, neither = "convex", "concave", "neither"
	var question, answer, explanation string
	switch rand.Intn(3) {
	case 0:
		a := randRange(ctx.Difficulty, 1, 6)
		if rand.Intn(2) == 0 {
			a = -a
		}
		b := randRange(ctx.Difficulty, 0, 9)
		f := fmt.Sprintf("%sx^{2}%s", coefStr(a), signedTerm(b, "x"))
		question = fmt.Sprintf("Is \\(f(x)=%s\\) convex or concave?", f)
		if a > 0 {
			answer = convex
		} else {
			answer = concave
		}
		explanation = fmt.Sprintf("\\(f'(x)=2\\cdot(%d)x%+d\\) and \\(f''(x)=%d\\), which is %s zero everywhere, so \\(f\\) is %s everywhere.", a, b, 2*a, map[bool]string{true: "above", false: "below"}[a > 0], answer)
	case 1:
		a := randRange(ctx.Difficulty, 1, 3)
		if rand.Intn(2) == 0 {
			a = -a
		}
		p := randRange(ctx.Difficulty, 1, 4)
		if rand.Intn(2) == 0 {
			p = -p
		}
		b := randRange(ctx.Difficulty, 0, 6)
		f := fmt.Sprintf("%sx^{3}%s", coefStr(a), signedTerm(b, "x"))
		question = fmt.Sprintf("Is \\(f(x)=%s\\) convex or concave at \\(x=%d\\)?", f, p)
		second := 6 * a * p
		if second > 0 {
			answer = convex
		} else {
			answer = concave
		}
		explanation = fmt.Sprintf("\\(f''(x)=%dx\\), so \\(f''(%d)=%d\\), which is %s zero. The curve is %s there.", 6*a, p, second, map[bool]string{true: "above", false: "below"}[second > 0], answer)
	default:
		n := randRange(ctx.Difficulty, 2, 6)
		question = fmt.Sprintf("Is \\(f(x)=x^{%d}\\) convex or concave at \\(x=1\\)?", n)
		answer = convex
		explanation = fmt.Sprintf("\\(f''(x)=%s\\), and at \\(x=1\\) that is %d, which is positive. So \\(f\\) is convex at \\(x=1\\).", powTerm(n*(n-1), n-2), n*(n-1))
	}
	distractors := []string{}
	for _, o := range []string{convex, concave, neither} {
		if o != answer {
			distractors = append(distractors, o)
		}
	}
	p := generator.MultipleChoice(question, answer, distractors...)
	p.Explanation = explanation
	return p
}

type nonDiffGen struct{}

func (g *nonDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Same reason as the continuity generator: the answer is a verdict, not an
	// expression, so it is offered as choices instead of a maths field.
	a := randRange(ctx.Difficulty, 1, 5)
	if rand.Intn(2) == 0 {
		a = -a
	}
	c := randRange(ctx.Difficulty, 1, 4)
	var question, answer, explanation string
	switch rand.Intn(5) {
	case 0:
		question = fmt.Sprintf("Is \\(f(x)=|x%+d|\\) differentiable at \\(x=%d\\)?", -a, a)
		answer = "no"
		explanation = fmt.Sprintf("The left slope at \\(x=%d\\) is -1 and the right slope is 1. A corner has no single tangent, so the derivative does not exist even though the function is continuous there.", a)
	case 1:
		question = fmt.Sprintf("Is \\(f(x)=x^{2}%+d\\) differentiable at \\(x=%d\\)?", c, a)
		answer = "yes"
		explanation = fmt.Sprintf("A polynomial is differentiable everywhere. \\(f'(x)=2x\\), so the slope at \\(x=%d\\) is %d.", a, 2*a)
	case 2:
		question = fmt.Sprintf("Is \\(f(x)=\\sqrt[3]{x%+d}\\) differentiable at \\(x=%d\\)?", -a, a)
		answer = "no"
		explanation = fmt.Sprintf("The function is continuous at \\(x=%d\\), but its slope grows without bound as \\(x\\) approaches %d. A vertical tangent means the derivative does not exist.", a, a)
	case 3:
		question = fmt.Sprintf("Is \\(f(x)=\\frac{%d}{x%+d}\\) differentiable at \\(x=%d\\)?", c, -a, a)
		answer = "no"
		explanation = fmt.Sprintf("The function is not even defined at \\(x=%d\\), because the denominator is zero there. Differentiability requires the function to exist at the point, so it fails earlier than the derivative.", a)
	default:
		question = fmt.Sprintf("Is \\(f(x)=(x%+d)|x%+d|\\) differentiable at \\(x=%d\\)?", -a, -a, a)
		answer = "yes"
		explanation = fmt.Sprintf("Both one-sided difference quotients at \\(x=%d\\) tend to 0, so the derivative exists and equals 0. The absolute value has a corner on its own, but multiplying by \\((x%+d)\\) flattens that corner into a horizontal tangent.", a, -a)
	}
	p := generator.MultipleChoice(question, answer, "yes", "no")
	p.Explanation = explanation
	return p
}

type partialDerivGen struct{}

func (g *partialDerivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	// f(x,y) = a·x²·y + b·xy²
	// ∂f/∂x = 2a·xy + b·y²
	// ∂f/∂y = a·x² + 2b·xy
	x0 := rand.Intn(max(1, scale*2)) + 1
	y0 := rand.Intn(max(1, scale*2)) + 1
	which := rand.Intn(2)
	var answerStr, explainStr string
	if which == 0 {
		val := 2*a*x0*y0 + b*y0*y0
		answerStr = fmt.Sprintf("%d", val)
		explainStr = fmt.Sprintf("df/dx = 2axy + by^2 = 2(%d)(%d)(%d)+%d(%d)^2 = %d.", a, x0, y0, b, y0, val)
	} else {
		val := a*x0*x0 + 2*b*x0*y0
		answerStr = fmt.Sprintf("%d", val)
		explainStr = fmt.Sprintf("df/dy = ax^2 + 2bxy = (%d)(%d)^2+2(%d)(%d)(%d) = %d.", a, x0, b, x0, y0, val)
	}
	dir := []string{"x", "y"}[which]
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\partial f/\\partial %s\\) at \\((%d,%d)\\) for \\(f(x,y)=%dx^{2}y+%dxy^{2}\\). (enter a number)", dir, x0, y0, a, b),
		Answer:      answerStr,
		Explanation: explainStr,
	}
}

// --- Integral generators ---

type arcLengthGen struct{}

func (g *arcLengthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Arc length integrates speed sqrt(1+(f')^2). For \\(y=3x\\), what is \\((f')^2\\)? (enter a number)", "9", "f'=3, so (f')^2=9 and the integrand is sqrt(1+9)=sqrt(10). Arc length is the integral of speed: ∫ₐᵇ √(1+(f'(x))²) dx."},
		{"What is the arc length of \\(y=2x\\) from \\(x=0\\) to 1?", "sqrt(5)", "f'(x)=2, so L = ∫₀¹√(1+4)dx = √5."},
		{"What is the arc length of \\(y=x\\) from \\(x=0\\) to 3?", "3sqrt(2)", "f'(x)=1, so L = ∫₀³√2 dx = 3√2."},
		{"In the arc-length integrand √(1+(f')²), what is f' for \\(y=3x\\)?", "3", "The derivative of 3x is the constant 3."},
	}
	tableHard := []entry{
		{"What is the arc length of \\(y=4x\\) from \\(x=0\\) to 2?", "2sqrt(17)", "f'(x)=4, so L = ∫₀²√17 dx = 2√17."},
		{"What integrand replaces √(1+(f')²) for a parametric curve \\((x(t),y(t))\\)?", "sqrt((dx/dt)^2+(dy/dt)^2)", "Parametric arc length uses ∫√((dx/dt)²+(dy/dt)²)dt."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type expIntegralGen struct{}

func (g *expIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\int %d e^{%d x} \\, dx\\).", a, b),
		Answer:      fmt.Sprintf("(%d/%d)e^(%dx)+C", a, b, b),
		Explanation: fmt.Sprintf("∫%d·e^(%dx) dx = %d·(1/%d)·e^(%dx)+C = (%d/%d)e^(%dx)+C.", a, b, a, b, b, a, b, b),
	}
}

type improperIntegralGen struct{}

func (g *improperIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Determine if \\(\\int_{1}^{\\infty} \\frac{1}{x^{2}} \\, dx\\) converges or diverges.", "converges (to 1)", "∫₁^∞ 1/x² dx = lim_{b→∞} [-1/x]₁ᵇ = lim_{b→∞} (-1/b+1) = 1. The integral converges to 1."},
		{"Determine if \\(\\int_{1}^{\\infty} \\frac{1}{x^{3}} \\, dx\\) converges or diverges.", "converges (to 1/2)", "∫₁^∞ x⁻³dx = lim_{b→∞} [-1/(2x²)]₁ᵇ = 1/2. The integral converges to 1/2."},
		{"Determine if \\(\\int_{1}^{\\infty} \\frac{1}{x} \\, dx\\) converges or diverges.", "diverges", "∫₁^∞ 1/x dx = lim_{b→∞} [ln x]₁ᵇ = ∞. The integral diverges (logarithmic growth)."},
		{"What is the smallest integer p for which \\(\\int_{1}^{\\infty} dx/x^{p}\\) converges? (enter a number)", "2", "The p-integral converges iff p>1 (the cases p=2 and p=3 above converge; p=1 diverges), so the smallest integer is 2."},
	}
	tableHard := []entry{
		{"Determine if \\(\\int_{0}^{1} \\frac{1}{\\sqrt{x}} \\, dx\\) converges or diverges.", "converges (to 2)", "∫₀¹ x^{-1/2}dx = lim_{a→0+} [2√x]ₐ¹ = 2. The integral converges to 2."},
		{"Determine if \\(\\int_{1}^{\\infty} \\frac{1}{\\sqrt{x}} \\, dx\\) converges or diverges.", "diverges", "Here p=1/2 ≤ 1, so the p-integral diverges."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type numericalIntegralGen struct{}

func (g *numericalIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	switch rand.Intn(8) {
	case 0:
		return generator.Problem{
			Question:    "What method approximates the integral using rectangles? (enter trapezoids, parabolas, or rectangles)",
			Answer:      "rectangles",
			Explanation: "Riemann sums partition [a,b] into subintervals and use rectangles to approximate the area.",
		}
	case 1:
		return generator.Problem{
			Question:    "What method approximates the integral using trapezoids? (enter trapezoids, parabolas, or rectangles)",
			Answer:      "trapezoids",
			Explanation: "The trapezoidal rule approximates the integral with (Δx/2)(f(x₀)+2f(x₁)+...+2f(xₙ₋₁)+f(xₙ)).",
		}
	case 2:
		return generator.Problem{
			Question:    "What method approximates the integral using parabolas? (enter trapezoids, parabolas, or rectangles)",
			Answer:      "parabolas",
			Explanation: "Simpson's rule uses quadratic polynomials and is more accurate than trapezoidal for smooth functions.",
		}
	case 3:
		return generator.Problem{
			Question:    "What power of h gives the accuracy order of the trapezoidal rule? (enter a number)",
			Answer:      "2",
			Explanation: "The trapezoidal rule error is proportional to h^2 times f'' at some point, making it second-order accurate: O(h^2).",
		}
	case 4:
		return generator.Problem{
			Question:    "What power of h gives the accuracy order of Simpson's rule? (enter a number)",
			Answer:      "4",
			Explanation: "Simpson's rule error is proportional to h^4 times the fourth derivative, so it is O(h^4) and beats the trapezoidal rule on smooth functions.",
		}
	case 5:
		return generator.Problem{
			Question:    "What power of h gives the accuracy order of the midpoint rule? (enter a number)",
			Answer:      "2",
			Explanation: "The midpoint rule error is proportional to h^2 times f'' at some point, so it is second-order accurate like the trapezoidal rule.",
		}
	case 6:
		return generator.Problem{
			Question:    "What power of h gives the accuracy order of the left-endpoint Riemann sum? (enter a number)",
			Answer:      "1",
			Explanation: "A one-point Riemann sum error is proportional to h times f' at some point, so it is only first-order accurate: O(h).",
		}
	}
	// The step size for n equal subintervals of [a,b] is h = (b-a)/n.
	n := randRange(ctx.Difficulty, 2, 4+scale)
	a := randRange(ctx.Difficulty, 0, 3)
	b := a + n*randRange(ctx.Difficulty, 1, 3)
	h := decimal(float64(b-a)/float64(n), 4)
	return generator.Problem{
		Question: fmt.Sprintf("On \\([%d,%d]\\) with \\(n=%d\\) equal subintervals, what is the step size \\(h\\)? (enter a number)", a, b, n),
		Answer:   h,
		Explanation: fmt.Sprintf("\\(h = (b-a)/n = (%d-%d)/%d = %s\\). Every quadrature rule on this panel uses that step size.",
			b, a, n, h),
	}
}

type riemannCriteriaGen struct{}

func (g *riemannCriteriaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "Is f(x)=x² Riemann integrable on [0,1]?", correct: "yes", reason: "x² is continuous on [0,1], and continuous functions on closed intervals are Riemann integrable.", distractors: []string{"no"}},
		{question: "Is f(x)={0 for x rational, 1 for x irrational} Riemann integrable on [0,1]?", correct: "no", reason: "The upper sum = 1 and lower sum = 0 for any partition, so the integral does not exist (Dirichlet function).", distractors: []string{"yes"}},
		{question: "Is f(x)=1/x Riemann integrable on [0,1]?", correct: "no", reason: "1/x is unbounded near 0, so it is not Riemann integrable on [0,1] (though the improper integral exists).", distractors: []string{"yes"}},
		{question: "Is a function with finitely many discontinuities Riemann integrable?", correct: "yes", reason: "A bounded function with finitely many discontinuities on [a,b] is Riemann integrable.", distractors: []string{"no"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.integral.riemann_criteria")...))
}

type trigIntegralsGen struct{}

func (g *trigIntegralsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\int \\sin(x) \\, dx\\).", "-cos(x)+C", "\\(\\int \\sin(x) \\, dx = -\\cos(x)+C\\)."},
		{"Find \\(\\int \\cos(x) \\, dx\\).", "sin(x)+C", "\\(\\int \\cos(x) \\, dx = \\sin(x)+C\\)."},
		{"Find \\(\\int \\sec^{2}(x) \\, dx\\).", "tan(x)+C", "\\(\\int \\sec^{2}(x) \\, dx = \\tan(x)+C\\)."},
		{"Find \\(\\int \\csc^{2}(x) \\, dx\\).", "-cot(x)+C", "\\(\\int \\csc^{2}(x) \\, dx = -\\cot(x)+C\\)."},
		{"Find \\(\\int \\sec(x)\\tan(x) \\, dx\\).", "sec(x)+C", "\\(\\int \\sec(x)\\tan(x) \\, dx = \\sec(x)+C\\)."},
		{"Find \\(\\int \\csc(x)\\cot(x) \\, dx\\).", "-csc(x)+C", "\\(\\int \\csc(x)\\cot(x) \\, dx = -\\csc(x)+C\\)."},
		{"Find \\(\\int \\tan(x) \\, dx\\).", "ln|sec(x)|+C", "\\(\\int \\tan(x) \\, dx = \\ln|\\sec(x)|+C\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type weierstrassSubGen struct{}

func (g *weierstrassSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(6) {
	case 0:
		return generator.Problem{
			Question:    "What is the Weierstrass substitution? (enter as t = ...)",
			Answer:      "t = tan(x/2)",
			Explanation: "Let \\(t=\\tan(x/2)\\). Then \\(\\sin(x)=2t/(1+t^{2})\\), \\(\\cos(x)=(1-t^{2})/(1+t^{2})\\), \\(dx=2/(1+t^{2})dt\\).",
		}
	case 1:
		return generator.Problem{
			Question:    "What is \\(\\sin(x)\\) in terms of \\(t=\\tan(x/2)\\)?",
			Answer:      "2t/(1+t^2)",
			Explanation: "\\(\\sin(x) = 2t/(1+t^{2})\\) where \\(t=\\tan(x/2)\\).",
		}
	case 2:
		return generator.Problem{
			Question:    "What is \\(\\cos(x)\\) in terms of \\(t=\\tan(x/2)\\)?",
			Answer:      "(1-t^2)/(1+t^2)",
			Explanation: "\\(\\cos(x) = (1-t^{2})/(1+t^{2})\\) where \\(t=\\tan(x/2)\\).",
		}
	case 3:
		return generator.Problem{
			Question:    "What is \\(dx\\) in terms of \\(t=\\tan(x/2)\\)?",
			Answer:      "2/(1+t^2) dt",
			Explanation: "\\(dx = 2/(1+t^{2}) dt\\) where \\(t=\\tan(x/2)\\).",
		}
	case 4:
		return generator.Problem{
			Question:    "What is \\(\\tan(x)\\) in terms of \\(t=\\tan(x/2)\\)?",
			Answer:      "2t/(1-t^2)",
			Explanation: "\\(\\tan(x) = \\sin(x)/\\cos(x) = \\frac{2t}{1-t^{2}}\\) where \\(t=\\tan(x/2)\\), valid where \\(t \\ne \\pm 1\\).",
		}
	}
	// A concrete value keeps the substitution honest: the learner cannot answer by
	// pattern-matching the general identity.
	t := randRange(ctx.Difficulty, 1, 4)
	sinValue := decimal(2*float64(t)/float64(1+t*t), 4)
	cosValue := decimal(float64(1-t*t)/float64(1+t*t), 4)
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(t=\\tan(x/2)=%d\\), what is \\(\\sin(x)\\)? (enter a decimal)", t),
			Answer:      sinValue,
			Explanation: fmt.Sprintf("\\(\\sin(x) = 2t/(1+t^{2}) = %d/(1+%d) = %s\\).", 2*t, t*t, sinValue),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(t=\\tan(x/2)=%d\\), what is \\(\\cos(x)\\)? (enter a decimal)", t),
		Answer:      cosValue,
		Explanation: fmt.Sprintf("\\(\\cos(x) = (1-t^{2})/(1+t^{2}) = (1-%d)/(1+%d) = %s\\).", t*t, t*t, cosValue),
	}
}

// --- Limit generators ---

type discontinuityGen struct{}

func (g *discontinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The four kinds are fixed mathematics; the function, the coefficient and the
	// singular point are drawn, so the answer word is derived from the shape that
	// was actually built rather than selected from a list of answers.
	a := randRange(ctx.Difficulty, 1, 5)
	if rand.Intn(2) == 0 {
		a = -a
	}
	k := randRange(ctx.Difficulty, 1, 3)
	b := randRange(ctx.Difficulty, 1, 4)
	switch rand.Intn(4) {
	case 0:
		f := fmt.Sprintf("\\frac{%d}{x-%d}", k, a)
		if rand.Intn(2) == 0 {
			f = fmt.Sprintf("\\frac{%d}{(x-%d)^{2}}", k, a)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What kind of discontinuity does \\(f(x)=%s\\) have at \\(x=%d\\)?", f, a),
			Answer:      "infinite",
			Explanation: fmt.Sprintf("As \\(x \\to %d\\) the denominator tends to zero while the numerator stays at %d, so \\(f\\) grows without bound. The line \\(x=%d\\) is a vertical asymptote, which makes the discontinuity infinite.", a, k, a),
		}
	case 1:
		f := fmt.Sprintf("\\frac{|x-%d|}{x-%d}", a, a)
		if rand.Intn(2) == 0 {
			f = fmt.Sprintf("\\frac{x-%d}{|x-%d|}", a, a)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What kind of discontinuity does \\(f(x)=%s\\) have at \\(x=%d\\)?", f, a),
			Answer:      "jump",
			Explanation: fmt.Sprintf("The one-sided limits differ: from the left the ratio tends to -1 and from the right to 1. Both are finite but unequal, so the discontinuity at \\(x=%d\\) is a jump.", a),
		}
	case 2:
		f := fmt.Sprintf("\\frac{x^{2}-%d}{x-%d}", a*a, a)
		if rand.Intn(2) == 0 {
			f = fmt.Sprintf("\\frac{(x-%d)(x+%d)}{x-%d}", a, b, a)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What kind of discontinuity does \\(f(x)=%s\\) have at \\(x=%d\\)?", f, a),
			Answer:      "removable",
			Explanation: fmt.Sprintf("The numerator shares the factor \\(x-%d\\) with the denominator, so the quotient simplifies to a polynomial for \\(x \\neq %d\\). The limit exists and equals the value that simplification gives, so the single missing point is removable.", a, a),
		}
	default:
		f := fmt.Sprintf("\\sin\\left(\\frac{%d}{x-%d}\\right)", k, a)
		if rand.Intn(2) == 0 {
			f = fmt.Sprintf("\\cos\\left(\\frac{%d}{x-%d}\\right)", k, a)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What kind of discontinuity does \\(f(x)=%s\\) have at \\(x=%d\\)?", f, a),
			Answer:      "essential",
			Explanation: fmt.Sprintf("As \\(x \\to %d\\) the argument grows without bound, so the function oscillates through its whole range infinitely often. No limit exists and the value never settles, which makes the discontinuity essential.", a),
		}
	}
}

type indeterminateGen struct{}

func (g *indeterminateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is \\(0/0\\) an indeterminate form?", "yes", "\\(0/0\\) is indeterminate, the limit could be any real number depending on the functions."},
		{"Is \\(\\infty/\\infty\\) an indeterminate form?", "yes", "\\(\\infty/\\infty\\) is indeterminate, the ratio of functions both approaching infinity could converge to anything."},
		{"Is \\(0 \\cdot \\infty\\) an indeterminate form?", "yes", "\\(0 \\cdot \\infty\\) is indeterminate, it can be rewritten as \\(0/0\\) or \\(\\infty/\\infty\\)."},
		{"Is \\(\\infty - \\infty\\) an indeterminate form?", "yes", "\\(\\infty - \\infty\\) is indeterminate, the difference of two quantities approaching infinity could be anything."},
		{"Is \\(0^{\\infty}\\) an indeterminate form?", "no", "\\(0^{\\infty} = 0\\). If the base approaches 0 and the exponent approaches \\(\\infty\\), the result is 0 (not indeterminate)."},
		{"Is \\(1^{\\infty}\\) an indeterminate form?", "yes", "\\(1^{\\infty}\\) is indeterminate, it often arises with exponential limits and can approach \\(e\\), 1, or other values."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type bigOGen struct{}

func (g *bigOGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does f(x)=O(g(x)) as x→∞ mean? (yes/no: |f| bounded by M|g| for large x)", "yes", "f=O(g) means there exists M>0 and x₀ such that |f(x)|≤M·|g(x)| for all x>x₀."},
		{"Is x²=O(x³) as x→∞?", "yes", "x² ≤ 1·x³ for x≥1, so x²=O(x³)."},
		{"Is x³=O(x²) as x→∞?", "no", "x³/x² = x → ∞, so no constant M can bound x³ by x² for all large x."},
		{"Is sin(x)=O(1) as x→∞?", "yes", "|sin(x)| ≤ 1, so sin(x)=O(1). In fact, any bounded function is O(1)."},
		{"What is the big-O notation primarily used for?", "describing asymptotic growth rates", "Big-O notation describes the limiting behavior of functions, which function dominates as the argument grows."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type littleOGen struct{}

func (g *littleOGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "What does f(x)=o(g(x)) as x→∞ mean? (yes/no: f(x)/g(x) tends to 0)", correct: "yes", reason: "f=o(g) means that f(x)/g(x) → 0 as x → ∞ (f grows strictly slower than g).", distractors: []string{"no"}},
		{question: "Is x²=o(x³) as x→∞?", correct: "yes", reason: "x²/x³ = 1/x → 0, so x²=o(x³).", distractors: []string{"no"}},
		{question: "Is x³=o(x³) as x→∞?", correct: "no", reason: "x³/x³ = 1 → 1, not 0, so x³≠o(x³). However, x³=O(x³).", distractors: []string{"yes"}},
		{question: "What is the difference between O and o? (yes/no: o needs ratio → 0, O allows any bounded ratio)", correct: "yes", reason: "f=O(g) means |f|≤M·|g| for large x. f=o(g) means f/g → 0, which is a stricter requirement: o requires ratio → 0; O allows any bounded ratio.", distractors: []string{"no"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.limit.little_o")...))
}

type uniformContinuityGen struct{}

func (g *uniformContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is uniform continuity? (yes/no: delta depends only on epsilon)", "yes", "A function is uniformly continuous if for any ε>0, there exists δ>0 (depending only on ε, not on the point) such that |x-y|<δ implies |f(x)-f(y)|<ε for ALL x,y."},
		{"Is f(x)=x² uniformly continuous on ℝ?", "no", "x² is not uniformly continuous on ℝ because the slope grows unbounded. However, it IS uniformly continuous on bounded intervals."},
		{"Is f(x)=x uniformly continuous on ℝ?", "yes", "|x-y|<δ implies |f(x)-f(y)|=|x-y|<δ. Just take δ=ε. f(x)=x is uniformly continuous on ℝ."},
		{"Is f(x)=√x uniformly continuous on [0,∞)?", "yes", "√x is continuous on [0,∞) and continuous functions on closed bounded intervals are uniformly continuous. The extension to [0,∞) works due to decreasing slope."},
		{"What theorem guarantees uniform continuity on [a,b]?", "Heine-Cantor theorem", "The Heine-Cantor theorem states that a continuous function on a closed bounded interval [a,b] is uniformly continuous."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type weierstrassLimitGen struct{}

func (g *weierstrassLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "What does the extreme value theorem (Weierstrass) state? (yes/no: continuous f on [a,b] attains max and min)", correct: "yes", reason: "If f is continuous on [a,b], then there exist c,d∈[a,b] such that f(c)≤f(x)≤f(d) for all x∈[a,b]: a continuous function on [a,b] attains its max and min.", distractors: []string{"no"}},
		{question: "Does the extreme value theorem apply to f(x)=1/x on (0,1]?", correct: "no", reason: "The interval must be closed. (0,1] is not closed, and in fact 1/x → ∞ as x→0+.", distractors: []string{"yes"}},
		{question: "Does the extreme value theorem apply to f(x)=x on [0,1]?", correct: "yes", reason: "f(x)=x is continuous on the closed bounded interval [0,1], so it attains max=1 at x=1 and min=0 at x=0.", distractors: []string{"no"}},
		{question: "Does the extreme value theorem apply to f(x)={x for x<1, 2 for x=1} on [0,1]?", correct: "no", reason: "f is not continuous on [0,1] (jump at x=1), so the EVT does not apply.", distractors: []string{"yes"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.limit.weierstrass")...))
}

// --- Sequence generators ---

type seqConceptGen struct{}

func (g *seqConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	switch rand.Intn(3) {
	case 0:
		// a_n = n/(n+k): every term is a fresh fraction.
		k := randRange(ctx.Difficulty, 1, 2+scale)
		terms := make([]string, 3)
		for i := 1; i <= 3; i++ {
			terms[i-1] = fmt.Sprintf("%d/%d", i, i+k)
		}
		answer := strings.Join(terms, ", ")
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(a_{n} = n/(n+%d)\\), write the first 3 terms. (separate with commas)", k),
			Answer:      answer,
			Explanation: fmt.Sprintf("Substituting \\(n=1,2,3\\) gives \\(%s\\). The terms increase toward 1 as \\(n\\to\\infty\\).", answer),
		}
	case 1:
		// An arithmetic sequence: a_n = a_1 + (n-1)d.
		a := randRange(ctx.Difficulty, 1, 9)
		d := randRange(ctx.Difficulty, 1, 5)
		terms := make([]string, 4)
		for i := 0; i < 4; i++ {
			terms[i] = strconv.Itoa(a + i*d)
		}
		answer := strings.Join(terms, ", ")
		return generator.Problem{
			Question:    fmt.Sprintf("For the arithmetic sequence with \\(a_{1}=%d\\) and common difference \\(d=%d\\), write the first 4 terms. (separate with commas)", a, d),
			Answer:      answer,
			Explanation: fmt.Sprintf("Each term adds \\(d=%d\\): \\(%s\\).", d, answer),
		}
	}
	// A geometric sequence: a_n = a r^(n-1).
	a := randRange(ctx.Difficulty, 1, 4)
	r := randRange(ctx.Difficulty, 2, 3)
	terms := make([]string, 3)
	v := a
	for i := 0; i < 3; i++ {
		terms[i] = strconv.Itoa(v)
		v *= r
	}
	answer := strings.Join(terms, ", ")
	return generator.Problem{
		Question:    fmt.Sprintf("For the geometric sequence with \\(a_{1}=%d\\) and ratio \\(r=%d\\), write the first 3 terms. (separate with commas)", a, r),
		Answer:      answer,
		Explanation: fmt.Sprintf("Multiply by \\(r=%d\\) each time: \\(%s\\).", r, answer),
	}
}

type seqConvergenceGen struct{}

func (g *seqConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The answers used to be prose ("no (diverges to \\(\\infty\\))") on a concept
	// graded as a tuple, so a learner who answered "False" was compared against
	// text and marked wrong, and nobody could type the parenthetical anyway
	// (report 178). The question is now answerable by tapping yes or no, and the
	// limit it converges to lives in the explanation where it belongs.
	type entry struct {
		q      string
		answer string
		reason string
	}
	table := []entry{
		{"Does \\(a_{n} = 1/n\\) converge?", "yes", "As \\(n\\to\\infty\\), \\(1/n \\to 0\\), so the sequence converges to 0."},
		{"Does \\(a_{n} = n/(n+1)\\) converge?", "yes", "\\(n/(n+1) = 1/(1+1/n) \\to 1/(1+0) = 1\\), so it converges to 1."},
		{"Does \\(a_{n} = (-1)^{n}\\) converge?", "no", "\\((-1)^{n}\\) oscillates between -1 and 1 without approaching a single limit."},
		{"Does \\(a_{n} = n^{2}\\) converge?", "no", "\\(n^{2} \\to \\infty\\) as \\(n\\to\\infty\\), so the sequence diverges."},
		{"Does \\(a_{n} = 2^{n}/n!\\) converge?", "yes", "The factorial grows faster than the exponential, so \\(2^{n}/n! \\to 0\\)."},
	}
	e := table[rand.Intn(len(table))]
	p := generator.MultipleChoice(e.q, e.answer, "yes", "no")
	p.Explanation = e.reason
	return p
}

type cauchySeqGen struct{}

func (g *cauchySeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "Does Cauchy mean terms get arbitrarily close as n increases? (yes/no)", correct: "yes", reason: "A sequence is Cauchy where for any ε>0, there exists N such that |aₘ-aₙ|<ε for all m,n>N; terms get arbitrarily close as n increases.", distractors: []string{"no"}},
		{question: "In ℝ, do all Cauchy sequences converge?", correct: "yes", reason: "ℝ is complete, every Cauchy sequence in ℝ converges to a real number.", distractors: []string{"no"}},
		{question: "In ℚ, do all Cauchy sequences converge?", correct: "no", reason: "ℚ is not complete. E.g., a sequence of rationals converging to √2 is Cauchy in ℚ but does not converge in ℚ.", distractors: []string{"yes"}},
		{question: "Is every convergent sequence Cauchy?", correct: "yes", reason: "Every convergent sequence is Cauchy. In ℝ, the converse also holds: Cauchy ⟺ convergent.", distractors: []string{"no"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.seq.cauchy")...))
}

type monotoneSeqGen struct{}

func (g *monotoneSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does monotone increasing mean each term is at least the previous one? (yes/no)", "yes", "Each term satisfies a(n+1) >= a(n). E.g., 1, 2, 3, 4, ..."},
		{"Does a bounded monotone sequence converge? (yes/no)", "yes", "Monotone convergence theorem: if a sequence is monotone (increasing or decreasing) and bounded, then it converges; a bounded monotone sequence converges."},
		{"Does a(n) = 1-1/n converge? (yes/no)", "yes", "a(n) increases: 1-1/(n+1) > 1-1/n. It is bounded above by 1. By the monotone convergence theorem, it converges to 1; it is bounded and increasing."},
		{"Does a(n) = n converge? (yes/no)", "no", "a(n) is increasing but not bounded above (unbounded), so it diverges."},
		{"If a sequence is not monotone, can it still converge?", "yes", "E.g., aₙ = (-1)ⁿ/n is not monotone but converges to 0."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type functionSequencesGen struct{}

func (g *functionSequencesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is a function sequence one where each term is a function? (yes/no)", "yes", "A function sequence {fₙ} has each fₙ: D → ℝ defined on a common domain D; a sequence where each term is a function."},
		{"Does pointwise convergence mean f_n(x) converges for each fixed x? (yes/no)", "yes", "fₙ → f pointwise if for each x in D, the numerical sequence fₙ(x) → f(x); sum f_n(x) converges for each fixed x."},
		{"Does uniform convergence imply pointwise convergence? (yes/no)", "yes", "Uniform means sup|fₙ(x)-f(x)| → 0; uniform convergence implies pointwise convergence, but the converse is false. E.g., fₙ(x)=xⁿ on [0,1] converges pointwise but not uniformly."},
		{"Does pointwise convergence imply uniform convergence? (yes/no)", "no", "fₙ(x)=xⁿ on [0,1] converges pointwise but not uniformly: sup norm stays 1."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    "f_n(x)=x^n on [0,1] converges pointwise but not uniformly. What is lim f_n(1)? (enter a number)",
			Answer:      "1",
			Explanation: "At x=1, x^n=1 for all n, so the pointwise limit there is 1 (elsewhere 0).",
		}
	}
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// --- Series generators ---

type seriesConceptGen struct{}

type harmonicSeriesGen struct{}

type positiveTermsGen struct{}

func (g *positiveTermsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 2, 5)
	c := randRange(ctx.Difficulty, 1, 9)
	// Classification answers are multiple choice: this concept is graded
	// symbolic, and "yes" cannot be typed into a maths field. Two of the four
	// shapes answer "no", because an easy pool where every answer is "yes"
	// lets a learner score without reading the question.
	switch rand.Intn(4) {
	case 0:
		q := fmt.Sprintf("Does every term of \\(\\sum_{n=1}^{\\infty}\\frac{1}{n^{%d}+%d}\\) satisfy \\(a_n\\ge 0\\)?", p, c)
		r := generator.MultipleChoice(q, "yes", "no")
		r.Explanation = "For every \\(n\\ge 1\\), the denominator is positive, so each term is positive. That positivity is what makes the comparison tests usable at all."
		return r
	case 1:
		q := fmt.Sprintf("Does every term of \\(\\sum_{n=1}^{\\infty}\\frac{(-1)^{n}}{n+%d}\\) satisfy \\(a_n\\ge 0\\)?", c)
		r := generator.MultipleChoice(q, "no", "yes")
		r.Explanation = "The signs alternate: odd \\(n\\) gives a negative term, even \\(n\\) a positive one. The comparison tests need non-negative terms, so this series must be handled as an alternating series instead."
		return r
	case 2:
		q := "If \\(0\\le a_n\\le b_n\\) and \\(\\sum b_n\\) converges, must \\(\\sum a_n\\) converge?"
		r := generator.MultipleChoice(q, "yes", "no")
		r.Explanation = "Direct comparison: the partial sums of \\(\\sum a_n\\) increase and are bounded above by the partial sums of the convergent \\(\\sum b_n\\), so they converge."
		return r
	}
	q := "If \\(0\\le a_n\\le b_n\\) and \\(\\sum b_n\\) diverges, must \\(\\sum a_n\\) diverge?"
	r := generator.MultipleChoice(q, "no", "yes")
	r.Explanation = "No. A smaller series can converge while the larger one diverges, for example \\(a_n=1/n^{2}\\) under \\(b_n=1/n\\). Comparison only transfers convergence downwards and divergence upwards."
	return r
}

type alternatingSeriesGen struct{}

type integralTestGen struct{}

type rootTestEntry struct {
	q, a, e string
}

// rootTestTable is package-level so the completeness test can assert the entry
// count and question uniqueness directly. A sampled test cannot tell an added
// duplicate entry from the eight intended items.
var rootTestTable = []rootTestEntry{
	{"Does the root test use the nth root of |a_n|? (yes/no)", "yes", "Root test evaluates limsup |a(n)|^(1/n); if L < 1, sum a(n) converges absolutely. If L > 1, diverges. Inconclusive if L=1."},
	{"Does sum (n/(n+1))^(n^2) converge? (yes/no)", "yes", "|a(n)|^(1/n) = (n/(n+1))^n = 1/(1+1/n)^n → 1/e < 1. The series converges by the root test."},
	{"Does sum 1/2^n converge? (yes/no)", "yes", "|a(n)|^(1/n) = 1/2 < 1. The series sum 1/2^n converges by root test (it is a geometric series)."},
	{"Is the root test useful when a_n involves nth powers? (yes/no)", "yes", "The root test is especially useful when a(n) contains expressions like (something)^n or nth powers."},
	{"For a_n = (2n/(3n+1))^n, what is the root-test limit L = lim |a_n|^(1/n)? (enter an exact fraction)", "2/3", "|a(n)|^(1/n) = 2n/(3n+1) → 2/3 < 1, so the series converges absolutely by the root test."},
	{"For a_n = ((3n+1)/n)^n, what is the root-test limit L = lim |a_n|^(1/n)? (enter a number)", "3", "|a(n)|^(1/n) = (3n+1)/n → 3 > 1, so the series diverges by the root test."},
	{"For a_n = ((n+1)/n)^n, what is the root-test limit L = lim |a_n|^(1/n)? (enter a number)", "1", "|a(n)|^(1/n) = (n+1)/n → 1, so the root test is inconclusive. The series still diverges: a(n) → e ≠ 0, so the nth-term test applies."},
	{"A series has root-test limit L = 1. Does the root test decide convergence? (yes/no)", "no", "When L = 1 the root test gives no conclusion, so another test is needed."},
}

type rootTestGen struct{}

func (g *rootTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	e := rootTestTable[rand.Intn(len(rootTestTable))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the coefficient of x^2 in the Maclaurin series for cos? (enter a fraction)", "-1/2", "General form sum c(n)(x-a)^n; cos(x) = 1 - x^2/2! + x^4/4! - ...; coefficient of x^2 is -1/2."},
		{"Does a power series converge for |x-a|<R? (yes/no)", "yes", "Radius R: the series converges for |x-a|<R, diverges for |x-a|>R; converges absolutely inside."},
		{"Does 1/(1-x) = sum x^n for |x|<1? (yes/no)", "yes", "1/(1-x) = 1 + x + x^2 + x^3 + ... = sum x^n, which converges for |x|<1."},
		{"What is the coefficient of x^2 in the Maclaurin series for e^x? (enter a fraction)", "1/2", "e^x = 1 + x + x^2/2! + x^3/3! + ... = sum x^n/n!, which converges for all real x; coefficient of x^2 is 1/2."},
		{"What is the coefficient of x^3 in the Maclaurin series for sin? (enter a fraction)", "-1/6", "sin(x) = x - x^3/3! + x^5/5! - ..., converging for all real x; coefficient of x^3 is -1/6."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type taylorSeriesGen struct{}

func (g *taylorSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is a Taylor series built from derivatives at a? (yes/no)", "yes", "Taylor: f(x) = sum f^(n)(a)(x-a)^n/n!; centered at a with radius R."},
		{"What is the coefficient of x^2 in the Maclaurin series for e^x? (enter a fraction)", "1/2", "e^x = sum x^n/n! with radius R=infinity; coefficient of x^2 is 1/2."},
		{"What is the coefficient of x^3 in the Maclaurin series for sin? (enter a fraction)", "-1/6", "sin(x) = x - x^3/3! + x^5/5! - ... with R=infinity; coefficient of x^3 is -1/6."},
		{"What is the coefficient of x^2 in the Maclaurin series for cos? (enter a fraction)", "-1/2", "cos(x) = 1 - x^2/2! + x^4/4! - ... with R=infinity; coefficient of x^2 is -1/2."},
		{"Does ln(1+x) = x - x^2/2 + x^3/3 - ... for |x|<1? (yes/no)", "yes", "ln(1+x) = x - x^2/2 + x^3/3 - x^4/4 + ... with R=1."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cauchyCriterionSeriesGen struct{}

func (g *cauchyCriterionSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Do convergent series have partial sums forming a Cauchy sequence? (yes/no)", "yes", "Sum a(n) converges iff for every e>0, there exists N such that |sum_{k=m}^n a(k)|<e for all n>=m>N; partial sums form a Cauchy sequence."},
		{"If sum a_n converges, must a_n tend to 0? (yes/no)", "yes", "If sum a(n) converges, then by Cauchy criterion with n=m, |a(m)|<e for m>N, so a(n) → 0."},
		{"Does condensation say sum f(n) and sum 2^n f(2^n) share fate? (yes/no)", "yes", "Cauchy condensation: for decreasing positive f, the condensed series sum 2^n f(2^n) has the same convergence as sum f(n); sum f(n) converges iff sum 2^n f(2^n) converges."},
		{"Does a_n → 0 guarantee sum a_n converges? (yes/no)", "no", "Harmonic series: terms tend to 0 but the sum diverges. Necessary, not sufficient."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    "Cauchy condensation replaces sum f(n) with sum 2^n f(2^n). What is the multiplier? (enter like 2^n)",
			Answer:      "2^n",
			Explanation: "Condensed term: 2^n f(2^n).",
		}
	}
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type fourierSeriesGen struct{}

func (g *fourierSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does a Fourier series use sines and cosines? (yes/no)", "yes", "A Fourier series sum (a(n)cos(nx) + b(n)sin(nx)) represents a periodic function as a sum of sines and cosines."},
		{"Is a_0 related to the average value? (yes/no)", "yes", "a(0) = (1/pi) int f(x) dx over one period, giving twice the average value; the average value over one period."},
		{"Does a_n involve cos(nx)? (yes/no)", "yes", "Fourier coefficient a(n) = (1/pi) int f(x)cos(nx) dx for n>=0."},
		{"Does b_n involve sin(nx)? (yes/no)", "yes", "Fourier coefficient b(n) = (1/pi) int f(x)sin(nx) dx for n>=1."},
		{"Which terms remain for an odd function? (enter sine)", "sine", "Odd functions have Fourier series with only sine terms: all a(n)=0."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type functionSeriesGen struct{}

func (g *functionSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "Is a function series one where each term is a function? (yes/no)", correct: "yes", reason: "A function series sum f(n)(x) is the sum of a sequence of functions; sum f_n(x) where each term is a function.", distractors: []string{"no"}},
		{question: "Does pointwise convergence mean convergence for each fixed x? (yes/no)", correct: "yes", reason: "The series converges pointwise if for each x in the domain, the numerical series converges; sum f_n(x) converges for each fixed x.", distractors: []string{"no"}},
		{question: "Does uniform convergence mean partial sums converge uniformly? (yes/no)", correct: "yes", reason: "The series converges uniformly if the sequence of partial sums converges uniformly; partial sums converge uniformly.", distractors: []string{"no"}},
		{question: "Which convergence does the Weierstrass M-test give? (enter uniform)", correct: "uniform", reason: "Weierstrass M-test for testing uniform convergence: if |f(n)(x)| ≤ M(n) and sum M(n) converges, then sum f(n)(x) converges uniformly and absolutely.", distractors: []string{"pointwise", "absolutely", "almost everywhere"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.series.function_series")...))
}

type derivTrigOtherGen struct{}

func (g *derivTrigOtherGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	a := randRange(ctx.Difficulty, 1, 9)
	inner := argX(a)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\tan(%s)\\).", inner),
			Answer:      coefTerm(a, fmt.Sprintf("sec(%s)^2", inner)),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\tan(%s)=\\sec^{2}(%s)\\cdot %d\\).", inner, inner, a),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\sec(%s)\\).", inner),
			Answer:      coefTerm(a, fmt.Sprintf("sec(%s)tan(%s)", inner, inner)),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\sec(%s)=\\sec(%s)\\tan(%s)\\cdot %d\\).", inner, inner, inner, a),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\csc(%s)\\).", inner),
			Answer:      "-" + coefTerm(a, fmt.Sprintf("csc(%s)cot(%s)", inner, inner)),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\csc(%s)=-\\csc(%s)\\cot(%s)\\cdot %d\\).", inner, inner, inner, a),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\cot(%s)\\).", inner),
			Answer:      "-" + coefTerm(a, fmt.Sprintf("csc(%s)^2", inner)),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\cot(%s)=-\\csc^{2}(%s)\\cdot %d\\).", inner, inner, a),
		}
	}
}

type derivInverseTrigGen struct{}

func (g *derivInverseTrigGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 2, 12)
	a2 := a * a
	ax := scaledXPow(a, 1)
	a2x := scaledXPow(a2, 2)
	switch rand.Intn(5) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\arcsin(%s)\\).", ax),
			Answer:      fmt.Sprintf("%d/sqrt(1-%s)", a, a2x),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\arcsin(%s) = \\frac{1}{\\sqrt{1-(%s)^{2}}}\\cdot %d = \\frac{%d}{\\sqrt{1-%s}}\\).", ax, ax, a, a, a2x),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\arccos(%s)\\).", ax),
			Answer:      fmt.Sprintf("-%d/sqrt(1-%s)", a, a2x),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\arccos(%s) = -\\frac{1}{\\sqrt{1-(%s)^{2}}}\\cdot %d = -\\frac{%d}{\\sqrt{1-%s}}\\).", ax, ax, a, a, a2x),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\arctan(%s)\\).", ax),
			Answer:      fmt.Sprintf("%d/(1+%s)", a, a2x),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\arctan(%s) = \\frac{1}{1+(%s)^{2}}\\cdot %d = \\frac{%d}{1+%s}\\).", ax, ax, a, a, a2x),
		}
	case 3:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\text{arccot}(%s)\\).", ax),
			Answer:      fmt.Sprintf("-%d/(1+%s)", a, a2x),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\text{arccot}(%s) = -\\frac{1}{1+(%s)^{2}}\\cdot %d = -\\frac{%d}{1+%s}\\).", ax, ax, a, a, a2x),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=\\arcsin(x/%d)\\).", a),
		Answer:      fmt.Sprintf("1/sqrt(%d-x^2)", a2),
		Explanation: fmt.Sprintf("\\(\\frac{d}{dx}\\arcsin(x/%d) = \\frac{1}{\\sqrt{1-(x/%d)^{2}}}\\cdot\\frac{1}{%d} = \\frac{1}{\\sqrt{%d-x^{2}}}\\).", a, a, a, a2),
	}
}

type derivGeneralExpGen struct{}

func (g *derivGeneralExpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	bases := []int{2, 3, 5, 6, 7, 10}
	a := bases[rand.Intn(len(bases))]
	b := randRange(ctx.Difficulty, 1, 8)
	if b == 1 {
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=%d^{x}\\).", a),
			Answer:      fmt.Sprintf("%d^x*ln(%d)", a, a),
			Explanation: fmt.Sprintf("\\(\\frac{d}{dx}%d^{x} = %d^{x}\\ln %d\\).", a, a, a),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=%d^{%dx}\\).", a, b),
		Answer:      fmt.Sprintf("%d*%d^(%dx)*ln(%d)", b, a, b, a),
		Explanation: fmt.Sprintf("Chain rule: \\(\\frac{d}{dx}%d^{%dx} = %d^{%dx}\\ln %d \\cdot %d\\).", a, b, a, b, a, b),
	}
}

type derivLogDiffGen struct{}

func (g *derivLogDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(3) {
	case 0:
		k := randRange(ctx.Difficulty, 2, 12)
		return generator.Problem{
			Question:    fmt.Sprintf("Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=x^{%dx}\\).", k),
			Answer:      fmt.Sprintf("x^(%dx)*(%dln(x)+%d)", k, k, k),
			Explanation: fmt.Sprintf("Take logs: \\(\\ln f = %dx\\ln x\\), so \\(f'/f = %d\\ln x + %d\\) and \\(f' = x^{%dx}(%d\\ln x + %d)\\).", k, k, k, k, k, k),
		}
	case 1:
		a := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=(%dx)^{x}\\).", a),
			Answer:      fmt.Sprintf("(%dx)^x*(ln(%dx)+1)", a, a),
			Explanation: fmt.Sprintf("Take logs: \\(\\ln f = x\\ln(%dx)\\), so \\(f'/f = \\ln(%dx) + 1\\) and \\(f' = (%dx)^{x}(\\ln(%dx) + 1)\\).", a, a, a, a),
		}
	}
	n := randRange(ctx.Difficulty, 2, 8)
	return generator.Problem{
		Question:    fmt.Sprintf("Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=x^{%d}e^{x}\\).", n),
		Answer:      fmt.Sprintf("x^%d*e^x*(%d/x+1)", n, n),
		Explanation: fmt.Sprintf("Take logs: \\(\\ln f = %d\\ln x + x\\), so \\(f'/f = %d/x + 1\\) and \\(f' = x^{%d}e^{x}(%d/x + 1)\\).", n, n, n, n),
	}
}

type derivHigherOrderGen struct{}

func (g *derivHigherOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(4) {
	case 0:
		n := randRange(ctx.Difficulty, 3, 12)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f''(x)\\) if \\(f(x)=x^{%d}\\).", n),
			Answer:      powTerm(n*(n-1), n-2),
			Explanation: fmt.Sprintf("\\(f'(x) = %s\\), so \\(f''(x) = %s\\).", powTerm(n, n-1), powTerm(n*(n-1), n-2)),
		}
	case 1:
		a := randRange(ctx.Difficulty, 1, 8)
		if a == 1 {
			return generator.Problem{
				Question:    "Find \\(f''(x)\\) if \\(f(x)=\\sin(x)\\).",
				Answer:      "-sin(x)",
				Explanation: "\\(f'(x) = \\cos(x)\\), so \\(f''(x) = -\\sin(x)\\).",
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f''(x)\\) if \\(f(x)=\\sin(%dx)\\).", a),
			Answer:      fmt.Sprintf("-%d*sin(%dx)", a*a, a),
			Explanation: fmt.Sprintf("\\(f'(x) = %d\\cos(%dx)\\), so \\(f''(x) = -%d\\sin(%dx)\\).", a, a, a*a, a),
		}
	case 2:
		a := randRange(ctx.Difficulty, 1, 8)
		if a == 1 {
			return generator.Problem{
				Question:    "Find \\(f''(x)\\) if \\(f(x)=e^{x}\\).",
				Answer:      "e^x",
				Explanation: "\\(f'(x) = e^{x}\\), so \\(f''(x) = e^{x}\\).",
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(f''(x)\\) if \\(f(x)=e^{%dx}\\).", a),
			Answer:      fmt.Sprintf("%de^(%dx)", a*a, a),
			Explanation: fmt.Sprintf("\\(f'(x) = %de^{%dx}\\), so \\(f''(x) = %de^{%dx}\\).", a, a, a*a, a),
		}
	}
	n := randRange(ctx.Difficulty, 4, 13)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the third derivative of \\(f(x)=x^{%d}\\).", n),
		Answer:      powTerm(n*(n-1)*(n-2), n-3),
		Explanation: fmt.Sprintf("\\(f' = %s\\), \\(f'' = %s\\), \\(f''' = %s\\).", powTerm(n, n-1), powTerm(n*(n-1), n-2), powTerm(n*(n-1)*(n-2), n-3)),
	}
}

type derivInverseFuncGen struct{}

func (g *derivInverseFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		a := randRange(ctx.Difficulty, 1, 5)
		b := randRange(ctx.Difficulty, 2, 9)
		return generator.Problem{
			Question:    fmt.Sprintf("If \\(f\\) is invertible and \\(f'(%d)=%d\\), what is \\((f^{-1})'(f(%d))\\)?", a, b, a),
			Answer:      fmt.Sprintf("1/%d", b),
			Explanation: fmt.Sprintf("The inverse function rule gives \\((f^{-1})'(f(%d)) = 1/f'(%d) = 1/%d\\): the slopes are reciprocals at matching points.", a, a, b),
		}
	}
	n := randRange(ctx.Difficulty, 2, 9)
	a := randRange(ctx.Difficulty, 2, 6)
	y := intPow(a, n)
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(f(x)=x^{%d}\\) and \\(g\\) is its inverse, what is \\(g'(%d)\\)?", n, y),
		Answer:      fmt.Sprintf("1/%d", n*intPow(a, n-1)),
		Explanation: fmt.Sprintf("\\(f'(x) = %s\\); at \\(x=%d\\), \\(f(%d)=%d\\) and \\(f'(%d)=%d\\), so \\(g'(%d)=1/%d\\).", powTerm(n, n-1), a, a, y, a, n*intPow(a, n-1), y, n*intPow(a, n-1)),
	}
}

type derivCriticalPointsGen struct{}

func (g *derivCriticalPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	b := 2 * randRange(ctx.Difficulty, 1, 8)
	c := randRange(ctx.Difficulty, 0, 9)
	if rand.Intn(2) == 0 {
		b = -b
	}
	f := fmt.Sprintf("x^{2}%s", signedTerm(b, "x"))
	if c != 0 {
		f += fmt.Sprintf("+%d", c)
	}
	root := -b / 2
	return generator.Problem{
		Question:    fmt.Sprintf("Find the critical point of \\(f(x)=%s\\).", f),
		Answer:      fmt.Sprintf("x=%d", root),
		Explanation: fmt.Sprintf("\\(f'(x)=2x%+d=0\\) gives \\(x=%d\\). A quadratic has exactly one critical point, and it is the vertex.", b, root),
	}
}

type derivSecondTestGen struct{}

func (g *derivSecondTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is \\(f''(0)\\) for \\(f(x)=x^{2}\\)? (enter a number)", "2", "\\(f''(x) = 2\\), so \\(f''(0) = 2 > 0\\): \\(x=0\\) is a local minimum."},
		{"What is \\(f''(0)\\) for \\(f(x)=-x^{2}\\)? (enter a number)", "-2", "\\(f''(x) = -2\\), so \\(f''(0) = -2 < 0\\): \\(x=0\\) is a local maximum."},
		{"If \\(f'(c)=0\\) and \\(f''(c)>0\\), is \\(c\\) a local minimum? (yes/no)", "yes", "Second derivative test: f''(c)>0 gives a local minimum, f''(c)<0 a local maximum, f''(c)=0 is inconclusive."},
		{"Is \\(x=0\\) a local minimum of \\(f(x)=x^{3}-3x^{2}\\) by the second derivative test? (yes/no)", "no", "\\(f''(x) = 6x-6\\); \\(f''(0) = -6 < 0\\), so \\(x=0\\) is a local maximum, not a minimum."},
		{"Is \\(x=2\\) a local minimum of \\(f(x)=x^{3}-3x^{2}\\) by the second derivative test? (yes/no)", "yes", "\\(f''(x) = 6x-6\\); \\(f''(2) = 6 > 0\\), so \\(x=2\\) is a local minimum."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivInflectionGen struct{}

func (g *derivInflectionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	b := 3 * randRange(ctx.Difficulty, 1, 12)
	if rand.Intn(2) == 0 {
		b = -b
	}
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("How many inflection points does \\(f(x)=x^{4}%s\\) have? (enter a number)", signedTerm(2*3*b, "x^{2}")),
			Answer:      strconv.Itoa(0),
			Explanation: fmt.Sprintf("\\(f''(x)=12x^{2}%+d\\) is positive everywhere, so the concavity never changes: no inflection point.", 12*3*b),
		}
	}
	root := -b / 3
	return generator.Problem{
		Question:    fmt.Sprintf("Find the inflection point of \\(f(x)=x^{3}%s\\).", signedTerm(b, "x^{2}")),
		Answer:      fmt.Sprintf("x=%d", root),
		Explanation: fmt.Sprintf("\\(f''(x)=6x%+d=0\\) gives \\(x=%d\\), and \\(f''\\) changes sign there, so it is an inflection point.", 2*b, root),
	}
}

type derivGlobalExtremaGen struct{}

func (g *derivGlobalExtremaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	b := 2 * randRange(ctx.Difficulty, 1, 10)
	m := b + 2
	vertex := b / 2
	if rand.Intn(2) == 0 {
		f := fmt.Sprintf("x^{2}-%dx", b)
		if rand.Intn(2) == 0 {
			return generator.Problem{
				Question:    fmt.Sprintf("What is the global minimum VALUE of \\(f(x)=%s\\) on \\([0,%d]\\)? (enter a number)", f, m),
				Answer:      strconv.Itoa(-(b * b) / 4),
				Explanation: fmt.Sprintf("Critical point \\(x=%d\\): \\(f(%d)=%d\\). Endpoints: \\(f(0)=0\\), \\(f(%d)=%d\\). The global minimum value is %d at \\(x=%d\\).", vertex, vertex, -(b*b)/4, m, m*(m-b), -(b*b)/4, vertex),
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What is the global maximum VALUE of \\(f(x)=%s\\) on \\([0,%d]\\)? (enter a number)", f, m),
			Answer:      strconv.Itoa(m * (m - b)),
			Explanation: fmt.Sprintf("Critical point \\(x=%d\\): \\(f(%d)=%d\\). Endpoints: \\(f(0)=0\\), \\(f(%d)=%d\\). The global maximum value is %d at \\(x=%d\\).", vertex, vertex, -(b*b)/4, m, m*(m-b), m*(m-b), m),
		}
	}
	f := fmt.Sprintf("-x^{2}+%dx", b)
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("What is the global maximum VALUE of \\(f(x)=%s\\) on \\([0,%d]\\)? (enter a number)", f, m),
			Answer:      strconv.Itoa((b * b) / 4),
			Explanation: fmt.Sprintf("Critical point \\(x=%d\\): \\(f(%d)=%d\\). Endpoints: \\(f(0)=0\\), \\(f(%d)=%d\\). The global maximum value is %d at \\(x=%d\\).", vertex, vertex, (b*b)/4, m, -m*(m-b), (b*b)/4, vertex),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the global minimum VALUE of \\(f(x)=%s\\) on \\([0,%d]\\)? (enter a number)", f, m),
		Answer:      strconv.Itoa(-m * (m - b)),
		Explanation: fmt.Sprintf("Critical point \\(x=%d\\): \\(f(%d)=%d\\). Endpoints: \\(f(0)=0\\), \\(f(%d)=%d\\). The global minimum value is %d at \\(x=%d\\).", vertex, vertex, (b*b)/4, m, -m*(m-b), -m*(m-b), m),
	}
}

type derivCurveSketchGen struct{}

// Every entry is a choice now. Three of the four used to answer "yes"/"no" on a
// concept graded symbolic, so the client showed a maths field: a learner typed
// y for yes and was marked wrong (report 118), and no form of the answer was
// typeable in a maths editor.
func (g *derivCurveSketchGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Every entry is a choice. Three of the four used to answer "yes"/"no" on a
	// concept graded symbolic, so the client showed a maths field: a learner typed
	// y for yes and was marked wrong (report 118).
	type row struct {
		q       string
		a       string
		e       string
		options []string
	}
	table := []row{
		{
			q:       "Where is the inflection point of \\(f(x)=x^{3}-3x\\)?",
			a:       "0",
			e:       "\\(f' = 3x^{2}-3\\): critical points at -1 and 1. \\(f'' = 6x\\): inflection at \\(x=0\\).",
			options: []string{"-1", "0", "1", "2"},
		},
		{
			q:       "Does curve sketching use the first derivative?",
			a:       "yes",
			e:       "Sketching combines intercepts, asymptotes (limits at infinity), critical points from the first derivative, and concavity/inflection from the second derivative.",
			options: []string{"yes", "no"},
		},
		{
			q:       "Is \\(f(x)=x^{2}\\) increasing on \\((0, \\infty)\\)?",
			a:       "yes",
			e:       "\\(f'(x) = 2x > 0\\) for \\(x>0\\), so f is increasing on \\((0,\\infty)\\).",
			options: []string{"yes", "no"},
		},
		{
			q:       "Does \\(f(x)=1/x\\) have a vertical asymptote at \\(x=0\\)?",
			a:       "yes",
			e:       "\\(\\lim_{x\\to 0} 1/x = \\pm\\infty\\) gives vertical asymptote \\(x=0\\); \\(\\lim_{x\\to\\pm\\infty} 1/x = 0\\) gives horizontal asymptote \\(y=0\\).",
			options: []string{"yes", "no"},
		},
	}
	claims := make([]conceptClaim, 0, len(table))
	for _, e := range table {
		distractors := make([]string, 0, len(e.options))
		for _, o := range e.options {
			if o != e.a {
				distractors = append(distractors, o)
			}
		}
		claims = append(claims, conceptClaim{
			question:    e.q,
			correct:     e.a,
			distractors: distractors,
			reason:      e.e,
		})
	}
	return conceptPickMC(append(claims, computedCurveSketchClaims()...))
}

type derivTangentLineGen struct{}

func (g *derivTangentLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// f(x)=x^n at x=a: the tangent is y = a^n + n a^(n-1) (x-a).
	n := randRange(ctx.Difficulty, 2, 4)
	a := randRange(ctx.Difficulty, 1, 4)
	fv := intPow(a, n)
	slope := n * intPow(a, n-1)
	return generator.Problem{
		Question: fmt.Sprintf("Find the tangent line to \\(f(x)=x^{%d}\\) at \\(x=%d\\).", n, a),
		Answer:   fmt.Sprintf("%d+%d(x-%d)", fv, slope, a),
		Explanation: fmt.Sprintf(
			"Subgoal 1: Differentiate: \\(f'(x)=%d%s\\).\\n"+
				"Subgoal 2: The slope at \\(x=%d\\) is \\(f'(%d)=%d\\) and the point of tangency is \\((%d,%d)\\).\\n"+
				"Subgoal 3: Point-slope form gives \\(y-%d=%d(x-%d)\\), that is \\(y=%d+%d(x-%d)\\).",
			n, varPow("x", n-1), a, a, slope, a, fv, fv, slope, a, fv, slope, a),
	}
}

type derivLinearApproxGen struct{}

func (g *derivLinearApproxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Linearize f(x)=x^n at x=a and evaluate the tangent at a+h.
	n := randRange(ctx.Difficulty, 2, 3)
	a := randRange(ctx.Difficulty, 2, 6)
	h := []float64{0.1, 0.2}[rand.Intn(2)]
	fv := intPow(a, n)
	slope := n * intPow(a, n-1)
	estimate := float64(fv) + float64(slope)*h
	return generator.Problem{
		Question: fmt.Sprintf("Use the linearization of \\(f(x)=x^{%d}\\) at \\(x=%d\\) to estimate \\((%.1f)^{%d}\\). (enter a decimal)", n, a, float64(a)+h, n),
		Answer:   decimal(estimate, 4),
		Explanation: fmt.Sprintf(
			"Subgoal 1: \\(f(%d)=%d\\) and \\(f'(x)=%d%s\\), so \\(f'(%d)=%d\\).\\n"+
				"Subgoal 2: The linearization is \\(L(x)=%d+%d(x-%d)\\).\\n"+
				"Subgoal 3: Evaluate at \\(x=%.1f\\): \\(L(%.1f)=%d+%d(%.1f)=%s\\).",
			a, fv, n, varPow("x", n-1), a, slope, fv, slope, a, float64(a)+h, float64(a)+h, fv, slope, h, decimal(estimate, 4)),
	}
}

type derivNewtonGen struct{}

func (g *derivNewtonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// One Newton step for f(x)=x^2-c from x0: x1 = x0 - (x0^2-c)/(2 x0).
	c := randRange(ctx.Difficulty, 2, 9)
	x0 := randRange(ctx.Difficulty, 1, 4)
	if x0*x0 == c {
		// x0 would already be the root, which makes the step a no-op.
		x0 = x0%4 + 1
	}
	x1 := (float64(x0*x0) + float64(c)) / float64(2*x0)
	return generator.Problem{
		Question: fmt.Sprintf("One Newton step for \\(f(x)=x^{2}-%d\\) from \\(x_{0}=%d\\). (enter a number)", c, x0),
		Answer:   decimal(x1, 4),
		Explanation: fmt.Sprintf(
			"Subgoal 1: Newton's method uses \\(x_{1}=x_{0}-f(x_{0})/f'(x_{0})\\) with \\(f'(x)=2x\\).\\n"+
				"Subgoal 2: \\(f(%d)=%d-%d=%d\\) and \\(f'(%d)=%d\\).\\n"+
				"Subgoal 3: \\(x_{1}=%d-%d/%d=%s\\).",
			x0, x0*x0, c, x0*x0-c, x0, 2*x0, x0, x0*x0-c, 2*x0, decimal(x1, 4)),
	}
}

type derivErrorPropGen struct{}

func (g *derivErrorPropGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// y=x^n with x measured as x0±dx gives dy = n x0^(n-1) dx.
	n := randRange(ctx.Difficulty, 2, 3)
	x0 := randRange(ctx.Difficulty, 2, 5)
	dx := []float64{0.1, 0.2}[rand.Intn(2)]
	dy := float64(n*intPow(x0, n-1)) * dx
	return generator.Problem{
		Question: fmt.Sprintf("If \\(y=x^{%d}\\) and \\(x=%d\\pm %.1f\\), estimate \\(\\Delta y\\). (enter a decimal)", n, x0, dx),
		Answer:   decimal(dy, 4),
		Explanation: fmt.Sprintf(
			"Subgoal 1: For small changes \\(\\Delta y \\approx dy = f'(x)\\,dx\\).\\n"+
				"Subgoal 2: \\(f'(x)=%d%s\\), so \\(f'(%d)=%d\\).\\n"+
				"Subgoal 3: \\(\\Delta y \\approx %d \\cdot %.1f = %s\\).",
			n, varPow("x", n-1), x0, n*intPow(x0, n-1), n*intPow(x0, n-1), dx, decimal(dy, 4)),
	}
}

type derivImplicitSecondGen struct{}

func (g *derivImplicitSecondGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// x^2+y^2=r^2 gives y'=-x/y and y''=-r^2/y^3.
	r := randRange(ctx.Difficulty, 3, 6)
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question: fmt.Sprintf("For \\(x^{2}+y^{2}=%d\\), what is \\(dy/dx\\)?", r*r),
			Answer:   "-x/y",
			Explanation: fmt.Sprintf(
				"Subgoal 1: Differentiate both sides implicitly: \\(2x+2y\\,y'=0\\).\\n"+
					"Subgoal 2: Solve for \\(y'\\): \\(y'=-x/y\\).\\n"+
					"Subgoal 3: This holds on the circle \\(x^{2}+y^{2}=%d\\) wherever \\(y\\neq 0\\).", r*r),
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("For \\(x^{2}+y^{2}=%d\\), what is \\(d^{2}y/dx^{2}\\) in terms of \\(x,y\\)?", r*r),
		Answer:   fmt.Sprintf("-%d/y^3", r*r),
		Explanation: fmt.Sprintf(
			"Subgoal 1: From \\(y'=-x/y\\), differentiate again: \\(y''=-(y-xy')/y^{2}\\).\\n"+
				"Subgoal 2: Substitute \\(y'=-x/y\\): \\(y''=-(y+x^{2}/y)/y^{2}=-(y^{2}+x^{2})/y^{3}\\).\\n"+
				"Subgoal 3: On the circle \\(x^{2}+y^{2}=%d\\), so \\(y''=-%d/y^{3}\\).", r*r, r*r),
	}
}

type derivChainPowerGen struct{}

func (g *derivChainPowerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// (ax+b)^n differentiates to n a (ax+b)^(n-1).
	a := randRange(ctx.Difficulty, 1, 4)
	b := randRange(ctx.Difficulty, 1, 5)
	n := randRange(ctx.Difficulty, 2, 4)
	coef := n * a
	inner := innerLinear(a, b)
	answer := coefStr(coef) + "(" + inner + ")" + powSuffix(n-1)
	return generator.Problem{
		Question: fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=(%s)^{%d}\\).", inner, n),
		Answer:   answer,
		Explanation: fmt.Sprintf(
			"Subgoal 1: The outer power rule gives \\(%d(%s)^{%d}\\).\\n"+
				"Subgoal 2: The inner derivative is %d.\\n"+
				"Subgoal 3: Multiply: \\(f'(x)=%d(%s)^{%d}=%s\\).",
			n, inner, n-1, a, coef, inner, n-1, answer),
	}
}

type derivChainMultiGen struct{}

func (g *derivChainMultiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// z=ax+by with x=t, y=t^2 gives dz/dt = a + 2 b t.
	a := randRange(ctx.Difficulty, 1, 5)
	b := randRange(ctx.Difficulty, 1, 5)
	t0 := randRange(ctx.Difficulty, 1, 4)
	answer := a + 2*b*t0
	return generator.Problem{
		Question: fmt.Sprintf("For \\(z=%s\\), \\(x=t\\), \\(y=t^{2}\\), find \\(dz/dt\\) at \\(t=%d\\). (enter a number)", linearIn(a, b), t0),
		Answer:   fmt.Sprintf("%d", answer),
		Explanation: fmt.Sprintf(
			"Subgoal 1: The multivariable chain rule is \\(dz/dt = z_{x}\\,x'(t) + z_{y}\\,y'(t)\\).\\n"+
				"Subgoal 2: \\(z_{x}=%d\\), \\(z_{y}=%d\\), \\(x'(t)=1\\), \\(y'(t)=2t\\).\\n"+
				"Subgoal 3: \\(dz/dt = %d + %d t\\), so at \\(t=%d\\) it equals %d.",
			a, b, a, 2*b, t0, answer),
	}
}

type derivGradientGen struct{}

func (g *derivGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// f = c x^p y^q gives grad f = (c p x^(p-1) y^q, c q x^p y^(q-1)).
	c := randRange(ctx.Difficulty, 1, 4)
	p := randRange(ctx.Difficulty, 2, 3)
	q := randRange(ctx.Difficulty, 2, 3)
	f := monoTerm(c, p, q)
	fx := monoTerm(c*p, p-1, q)
	fy := monoTerm(c*q, p, q-1)
	return generator.Problem{
		Question: fmt.Sprintf("Find \\(\\nabla f\\) for \\(f(x,y)=%s\\).", f),
		Answer:   fmt.Sprintf("(%s, %s)", fx, fy),
		Explanation: fmt.Sprintf(
			"Subgoal 1: \\(\\nabla f=(f_{x},f_{y})\\).\\n"+
				"Subgoal 2: Differentiate in \\(x\\) holding \\(y\\) fixed: \\(f_{x}=%s\\).\\n"+
				"Subgoal 3: Differentiate in \\(y\\) holding \\(x\\) fixed: \\(f_{y}=%s\\). So \\(\\nabla f=(%s, %s)\\).",
			fx, fy, fx, fy),
	}
}

type derivDirectionalGen struct{}

func (g *derivDirectionalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// f=a x^2+b y^2, so D_u f at (p,q) is 2ap or 2bq with the sign of u.
	a := randRange(ctx.Difficulty, 1, 4)
	b := randRange(ctx.Difficulty, 1, 4)
	p := randRange(ctx.Difficulty, 1, 4)
	q := randRange(ctx.Difficulty, 1, 4)
	dirs := []struct {
		label string
		x, y  int
	}{
		{"(1,0)", 1, 0},
		{"(0,1)", 0, 1},
		{"(-1,0)", -1, 0},
		{"(0,-1)", 0, -1},
	}
	d := dirs[rand.Intn(len(dirs))]
	answer := d.x*2*a*p + d.y*2*b*q
	return generator.Problem{
		Question: fmt.Sprintf("For \\(f(x,y)=%s\\), find the directional derivative at \\((%d,%d)\\) in the direction \\(u=%s\\). (enter a number)", quadForm(a, b), p, q, d.label),
		Answer:   fmt.Sprintf("%d", answer),
		Explanation: fmt.Sprintf(
			"Subgoal 1: \\(D_{u}f=\\nabla f\\cdot u\\) for a unit direction \\(u\\).\\n"+
				"Subgoal 2: \\(\\nabla f=(%s, %s)\\), so at \\((%d,%d)\\) it is \\((%d, %d)\\).\\n"+
				"Subgoal 3: Dot with \\(u=%s\\): %d.",
			monoTerm(2*a, 1, 0), monoTerm(2*b, 0, 1), p, q, 2*a*p, 2*b*q, d.label, answer),
	}
}

type derivPartialHigherGen struct{}

func (g *derivPartialHigherGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// f = c x^p y^q; ask for one of the three second partials.
	c := randRange(ctx.Difficulty, 1, 4)
	p := randRange(ctx.Difficulty, 2, 4)
	q := randRange(ctx.Difficulty, 2, 4)
	f := monoTerm(c, p, q)
	which := []struct {
		label string
		coef  int
		px    int
		py    int
	}{
		{"f_{xx}", c * p * (p - 1), p - 2, q},
		{"f_{yy}", c * q * (q - 1), p, q - 2},
		{"f_{xy}", c * p * q, p - 1, q - 1},
	}[rand.Intn(3)]
	answer := monoTerm(which.coef, which.px, which.py)
	return generator.Problem{
		Question: fmt.Sprintf("For \\(f(x,y)=%s\\), find \\(%s\\).", f, which.label),
		Answer:   answer,
		Explanation: fmt.Sprintf(
			"Subgoal 1: Differentiate once in each variable the label names, holding the other fixed.\\n"+
				"Subgoal 2: Powers drop by one per derivative, so the surviving exponents are \\(x^{%d}y^{%d}\\).\\n"+
				"Subgoal 3: The coefficients multiply to %d, giving \\(%s\\).",
			which.px, which.py, which.coef, answer),
	}
}

type derivClairautGen struct{}

func (g *derivClairautGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// f = c x^m y^n, so f_xy = f_yx = c m n x^(m-1) y^(n-1) for polynomials.
	c := randRange(ctx.Difficulty, 1, 4)
	m := randRange(ctx.Difficulty, 2, 4)
	n := randRange(ctx.Difficulty, 2, 4)
	f := monoTerm(c, m, n)
	answer := monoTerm(c*m*n, m-1, n-1)
	return generator.Problem{
		Question: fmt.Sprintf("For \\(f(x,y)=%s\\), compute \\(f_{xy}\\) (Clairaut's theorem says this equals \\(f_{yx}\\)).", f),
		Answer:   answer,
		Explanation: fmt.Sprintf(
			"Subgoal 1: Differentiate in \\(x\\) first: \\(f_{x}=%s\\).\\n"+
				"Subgoal 2: Then differentiate that in \\(y\\): \\(f_{xy}=%s\\).\\n"+
				"Subgoal 3: The other order gives the same term because the mixed partials of a polynomial are continuous, which is exactly Clairaut's hypothesis.",
			monoTerm(c*m, m-1, n), answer),
	}
}

type derivTotalDiffGen struct{}

func (g *derivTotalDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// z=c x^p y^q, so dz = (p c x^(p-1) y^q) dx + (q c x^p y^(q-1)) dy.
	c := randRange(ctx.Difficulty, 1, 3)
	p := randRange(ctx.Difficulty, 2, 3)
	q := randRange(ctx.Difficulty, 2, 3)
	x0 := randRange(ctx.Difficulty, 1, 3)
	y0 := randRange(ctx.Difficulty, 1, 3)
	dx, dy := 0.1, 0.1
	fx := c * p * intPow(x0, p-1) * intPow(y0, q)
	fy := c * q * intPow(x0, p) * intPow(y0, q-1)
	dz := (float64(fx) + float64(fy)) * dx
	return generator.Problem{
		Question: fmt.Sprintf("For \\(z=%s\\), estimate \\(\\Delta z\\) when \\(x\\) grows from %d to %.1f and \\(y\\) grows from %d to %.1f. (enter a decimal)", monoTerm(c, p, q), x0, float64(x0)+dx, y0, float64(y0)+dy),
		Answer:   decimal(dz, 4),
		Explanation: fmt.Sprintf(
			"Subgoal 1: The total differential is \\(dz=z_{x}\\,dx+z_{y}\\,dy\\).\\n"+
				"Subgoal 2: At \\((%d,%d)\\), \\(z_{x}=%d\\) and \\(z_{y}=%d\\).\\n"+
				"Subgoal 3: \\(dz\\approx %d(%.1f)+%d(%.1f)=%s\\).",
			x0, y0, fx, fy, fx, dx, fy, dy, decimal(dz, 4)),
	}
}

type derivDiffContinuityGen struct{}

func (g *derivDiffContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// A classification answer cannot be typed into a maths field, so every case
	// is multiple choice. "differentiable but not continuous" is impossible by
	// the theorem this concept teaches, and is the distractor that catches a
	// learner who has the implication backwards.
	const both = "both"
	const contNotDiff = "continuous but not differentiable"
	const diffNotCont = "differentiable but not continuous"
	const neither = "neither"
	a := randRange(ctx.Difficulty, 1, 5)
	if rand.Intn(2) == 0 {
		a = -a
	}
	c := randRange(ctx.Difficulty, 1, 4)
	var question, answer, explanation string
	// Six shapes: the removable case at the end is the sixth, so the bound must
	// be six or the default branch is unreachable.
	switch rand.Intn(6) {
	case 0:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=|x%+d|\\) continuous, differentiable, both or neither?", a, -a)
		answer = contNotDiff
		explanation = fmt.Sprintf("The two one-sided limits agree, so \\(f\\) is continuous at \\(x=%d\\). The left slope is -1 and the right slope is 1, so the derivative does not exist there: a corner.", a)
	case 1:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=x^{2}%+d\\) continuous, differentiable, both or neither?", a, c)
		answer = both
		explanation = fmt.Sprintf("A polynomial is continuous and differentiable everywhere, so at \\(x=%d\\) it is both. \\(f'(x)=2x\\), so the slope there is %d.", a, 2*a)
	case 2:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=\\frac{%d}{x%+d}\\) continuous, differentiable, both or neither?", a, c, -a)
		answer = neither
		explanation = fmt.Sprintf("The denominator is zero at \\(x=%d\\), so \\(f\\) is not defined there and cannot be continuous. Differentiability requires continuity, so it is neither.", a)
	case 3:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=|x%+d|\\) continuous, differentiable, both or neither?", a+1, -a)
		answer = both
		explanation = fmt.Sprintf("The corner of \\(|x%+d|\\) sits at \\(x=%d\\), not at \\(x=%d\\). At \\(x=%d\\) the function is smooth, so it is continuous and differentiable: at a point away from a corner, \\(|x|\\) behaves like a line.", -a, a, a+1, a+1)
	case 4:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=\\sqrt[3]{x%+d}\\) continuous, differentiable, both or neither?", a, -a)
		answer = contNotDiff
		explanation = fmt.Sprintf("A cube root is continuous at \\(x=%d\\), since the limit equals the value. Its derivative behaves like \\((x%+d)^{-2/3}\\), which grows without bound near \\(x=%d\\), so the tangent is vertical and the derivative does not exist.", a, -a, a)
	default:
		question = fmt.Sprintf("At \\(x=%d\\), is \\(f(x)=\\frac{x^{2}%+d}{x%+d}\\) continuous, differentiable, both or neither?", a, -a*a, -a)
		answer = neither
		explanation = fmt.Sprintf("The numerator factors as \\((x%+d)(x%+d)\\), so the limit at \\(x=%d\\) exists and equals %d. The function itself is undefined there, because the expression still divides by zero: the discontinuity is removable, but at \\(x=%d\\) it is neither continuous nor differentiable.", -a, a, a, 2*a, a)
	}
	distractors := []string{}
	for _, o := range []string{both, contNotDiff, diffNotCont, neither} {
		if o != answer {
			distractors = append(distractors, o)
		}
	}
	p := generator.MultipleChoice(question, answer, distractors...)
	p.Explanation = explanation
	return p
}

type derivJacobianGen struct{}

func (g *derivJacobianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Two shapes: a power map, whose Jacobian is diagonal, and a linear map,
	// whose Jacobian is the constant coefficient matrix.
	if rand.Intn(2) == 0 {
		p := randRange(ctx.Difficulty, 2, 3)
		q := randRange(ctx.Difficulty, 2, 3)
		fx := monoTerm(p, p-1, 0)
		fy := monoTerm(q, 0, q-1)
		return generator.Problem{
			Question: fmt.Sprintf("Write the Jacobian of \\(T(x,y)=(x^{%d},y^{%d})\\).", p, q),
			Answer:   fmt.Sprintf("[[%s,0],[0,%s]]", fx, fy),
			Explanation: fmt.Sprintf(
				"Subgoal 1: Row 1 holds the partials of the first component \\(x^{%d}\\): \\(%s\\) in \\(x\\) and 0 in \\(y\\).\\n"+
					"Subgoal 2: Row 2 holds the partials of \\(y^{%d}\\): 0 in \\(x\\) and \\(%s\\) in \\(y\\).\\n"+
					"Subgoal 3: Assemble the matrix \\(\\begin{pmatrix}%s&0\\\\0&%s\\end{pmatrix}\\).",
				p, fx, q, fy, fx, fy),
		}
	}
	a := randRange(ctx.Difficulty, 1, 4)
	b := randRange(ctx.Difficulty, 1, 4)
	c := randRange(ctx.Difficulty, 1, 4)
	d := randRange(ctx.Difficulty, 1, 4)
	return generator.Problem{
		Question: fmt.Sprintf("Write the Jacobian of \\(T(x,y)=(%s,%s)\\).", linearIn(a, b), linearIn(c, d)),
		Answer:   fmt.Sprintf("[[%d,%d],[%d,%d]]", a, b, c, d),
		Explanation: fmt.Sprintf(
			"Subgoal 1: The first row holds \\(\\partial T_{1}/\\partial x=%d\\) and \\(\\partial T_{1}/\\partial y=%d\\).\\n"+
				"Subgoal 2: The second row holds \\(\\partial T_{2}/\\partial x=%d\\) and \\(\\partial T_{2}/\\partial y=%d\\).\\n"+
				"Subgoal 3: A linear map has a constant Jacobian, so the matrix is the coefficient matrix itself.",
			a, b, c, d),
	}
}

type integralNetChangeGen struct{}

func (g *integralNetChangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// v(t) = 3 k t^2 from t = 0 to t = T -> displacement = k T^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		T := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k * T * T * T
		coeff := 3 * k
		return generator.Problem{
			Question: fmt.Sprintf("A particle moves along a line with velocity \\(v(t) = %d t^2\\) m/s. Find the total displacement of the particle from \\(t = 0\\) to \\(t = %d\\) seconds.", coeff, T),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Net Change Theorem: displacement equals \\(\\int_{t_1}^{t_2} v(t)\\,dt\\).\\n"+
				"Subgoal 2: Substitute \\(v(t) = %d t^2\\), \\(t_1 = 0\\) and \\(t_2 = %d\\).\\n"+
				"Subgoal 3: Evaluate the definite integral: \\(\\int_{0}^{%d} %d t^2\\,dt = [%d t^3]_{0}^{%d} = %d\\).", coeff, T, T, coeff, k, T, ans),
		}
	case 1:
		// r(t) = 2 k t + b from t = 0 to t = T -> net change = k T^2 + b T
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		T := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k*T*T + b*T
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Water flows into a tank at rate \\(r(t) = %d t + %d\\) L/min. How many liters of water enter the tank between \\(t = 0\\) and \\(t = %d\\) minutes?", coeff, b, T),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: net volume change is \\(\\int_{a}^{b} r(t)\\,dt\\).\\n"+
				"Subgoal 2: Substitute integrand \\(r(t) = %d t + %d\\) with limits from 0 to %d.\\n"+
				"Subgoal 3: Evaluate antiderivative: \\([%d t^2 + %d t]_{0}^{%d} = %d \\cdot %d + %d \\cdot %d = %d\\).", coeff, b, T, k, b, T, k, T*T, b, T, ans),
		}
	case 2:
		// f'(x) = 2 k x and f(0) = c, find f(b) -> f(b) = c + k b^2
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		c := mathutil.ScaleInt(ctx.Difficulty, 1, 10)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := c + k*b*b
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Given \\(f'(x) = %d x\\) and \\(f(0) = %d\\), determine the value of \\(f(%d)\\).", coeff, c, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Apply the Net Change Theorem: \\(f(b) = f(a) + \\int_{a}^{b} f'(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute parameters \\(f(0) = %d\\) and \\(\\int_{0}^{%d} %d x\\,dx\\).\\n"+
				"Subgoal 3: Evaluate and sum: \\(%d + [%d x^2]_{0}^{%d} = %d + %d = %d\\).", c, b, coeff, c, k, b, c, k*b*b, ans),
		}
	case 3:
		// Inflow rate r(t) = 6 k t^2 from t = 0 to t = T -> volume = 2 k T^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		T := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		ans := 2 * k * T * T * T
		coeff := 6 * k
		return generator.Problem{
			Question: fmt.Sprintf("A chemical solution is added at rate \\(r(t) = %d t^2\\) mL/s. Compute the total amount added over the interval [0, %d].", coeff, T),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: net accumulation is \\(\\Delta V = \\int_{0}^{T} r(t)\\,dt\\).\\n"+
				"Subgoal 2: Integrate \\(\\int_{0}^{%d} %d t^2\\,dt = [%d t^3]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Evaluate: \\(%d \\cdot %d = %d\\).", T, coeff, 2*k, T, 2*k, T*T*T, ans),
		}
	default:
		// Acceleration a(t) = 2 k with v(0) = v0, find v(T) -> v(T) = v0 + 2 k T
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		v0 := mathutil.ScaleInt(ctx.Difficulty, 1, 15)
		T := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := v0 + 2*k*T
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("An accelerating object has constant acceleration \\(a(t) = %d\\) m/s\\(^2\\) and initial velocity \\(v(0) = %d\\) m/s. Find \\(v(%d)\\).", coeff, v0, T),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Net Change Theorem for acceleration: \\(v(T) = v(0) + \\int_{0}^{T} a(t)\\,dt\\).\\n"+
				"Subgoal 2: Substitute \\(v(0) = %d\\) and \\(\\int_{0}^{%d} %d\\,dt\\).\\n"+
				"Subgoal 3: Evaluate and compute: \\(%d + (%d \\cdot %d) = %d\\).", v0, T, coeff, v0, coeff, T, ans),
		}
	}
}

type integralVariableLimitsGen struct{}

func (g *integralVariableLimitsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	p := randRange(ctx.Difficulty, 1, 2+scale)
	// q is a multiple of p+1 so the chain rule cancels the antiderivative's
	// denominator and the answer stays an integer multiple of a power of x.
	k := randRange(ctx.Difficulty, 1, 3)
	q := k * (p + 1)
	a := randRange(ctx.Difficulty, 0, 3)
	exp := q*p + q - 1
	body := coefStr(k) + varPow("x", exp)
	upper := fmt.Sprintf("x^{%d}", q)
	lower := fmt.Sprintf("%d", a)
	flipped := rand.Intn(3) == 0
	if flipped {
		upper, lower = lower, upper
		body = "-" + body
	}
	dir := ""
	if flipped {
		dir = "Flipping the limits introduces a minus sign. "
	}
	return generator.Problem{
		Question: fmt.Sprintf("Find \\(\\frac{d}{dx}\\int_{%s}^{%s} t^{%d}\\,dt\\).", lower, upper, p),
		Answer:   body,
		Explanation: fmt.Sprintf("%sFTC part 1 with the chain rule: the integrand at the upper limit times the derivative of the upper limit, "+
			"\\(%s \\cdot %d x^{%d} = %s\\).", dir, upper, q, q-1, body),
	}
}

type integralMVTGen struct{}

func (g *integralMVTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Average value of f(x) = 2 k x on [0, b] -> f_avg = k b
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := k * b
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Find the average value of \\(f(x) = %d x\\) on the interval [0, %d].", coeff, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute parameters \\(a = 0\\), \\(b = %d\\) and \\(f(x) = %d x\\).\\n"+
				"Subgoal 3: Evaluate: \\(\\frac{1}{%d}\\int_{0}^{%d} %d x\\,dx = \\frac{1}{%d} [%d x^2]_{0}^{%d} = \\frac{%d}{%d} = %d\\).", b, coeff, b, b, coeff, b, k, b, k*b*b, b, ans),
		}
	case 1:
		// Average value of f(x) = 3 k x^2 on [0, b] -> f_avg = k b^2
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k * b * b
		coeff := 3 * k
		return generator.Problem{
			Question: fmt.Sprintf("Compute the average value of \\(f(x) = %d x^2\\) on [0, %d].", coeff, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the average value formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute limits from 0 to %d with integrand \\(%d x^2\\).\\n"+
				"Subgoal 3: Evaluate antiderivative: \\(\\frac{1}{%d} [%d x^3]_{0}^{%d} = \\frac{%d \\cdot %d}{%d} = %d\\).", b, coeff, b, k, b, k, b*b*b, b, ans),
		}
	case 2:
		// Find c in (0, 2k) guaranteed by MVT for f(x) = 2x on [0, 2k] -> c = k
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		b := 2 * k
		ans := k
		return generator.Problem{
			Question: fmt.Sprintf("Find the value \\(c\\) guaranteed by the Mean Value Theorem for Integrals for \\(f(x) = 2x\\) on [0, %d].", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the theorem: there exists \\(c \\in [a, b]\\) such that \\(f(c) = f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\\n"+
				"Subgoal 2: Compute \\(f_{\\text{avg}} = \\frac{1}{%d} \\int_{0}^{%d} 2x\\,dx = \\frac{1}{%d} [x^2]_{0}^{%d} = %d\\).\\n"+
				"Subgoal 3: Set \\(f(c) = f_{\\text{avg}} \\implies 2c = %d \\implies c = %d\\).", b, b, b, b, b, b, ans),
		}
	case 3:
		// Average value of linear f(x) = 2 k x + d on [0, b] -> f_avg = k b + d
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		d := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 7)
		ans := k*b + d
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Find the average value of \\(f(x) = %d x + %d\\) over [0, %d].", coeff, d, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\\n"+
				"Subgoal 2: Compute \\(\\int_{0}^{%d} (%d x + %d)\\,dx = [%d x^2 + %d x]_{0}^{%d} = %d\\).\\n"+
				"Subgoal 3: Divide by interval length %d: \\(\\frac{%d}{%d} = %d\\).", b, coeff, d, k, d, b, k*b*b+d*b, b, k*b*b+d*b, b, ans),
		}
	default:
		// Average value of f(x) = 4 k x^3 on [0, b] -> f_avg = k b^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		ans := k * b * b * b
		coeff := 4 * k
		return generator.Problem{
			Question: fmt.Sprintf("Compute the average value of \\(f(x) = %d x^3\\) on [0, %d].", coeff, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute integrand \\(%d x^3\\) on interval [0, %d].\\n"+
				"Subgoal 3: Evaluate: \\(\\frac{1}{%d} [k x^4]_{0}^{%d} = \\frac{%d \\cdot %d}{%d} = %d\\).", coeff, b, b, b, k, b*b*b*b, b, ans),
		}
	}
}

type integralWorkGen struct{}

func (g *integralWorkGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Spring force F(x) = 2 k x stretched from 0 to D -> W = k D^2
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		D := mathutil.ScaleInt(ctx.Difficulty, 1, 7)
		ans := k * D * D
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("A spring requires a force of \\(F(x) = %d x\\) N to stretch a distance of \\(x\\) meters. Calculate the work done in stretching the spring from \\(x = 0\\) to \\(x = %d\\) meters.", coeff, D),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for work: \\(W = \\int_{a}^{b} F(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute \\(F(x) = %d x\\) with limits \\(a = 0\\) and \\(b = %d\\).\\n"+
				"Subgoal 3: Evaluate the integral: \\(\\int_{0}^{%d} %d x\\,dx = [%d x^2]_{0}^{%d} = %d \\cdot %d = %d\\) J.", coeff, D, D, coeff, k, D, k, D*D, ans),
		}
	case 1:
		// Force F(x) = 3 k x^2 moving an object from 0 to D -> W = k D^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		D := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k * D * D * D
		coeff := 3 * k
		return generator.Problem{
			Question: fmt.Sprintf("A variable force \\(F(x) = %d x^2\\) N acts on an object moving along the x-axis. Find the work done in moving the object from \\(x = 0\\) to \\(x = %d\\) meters.", coeff, D),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{0}^{D} F(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute force \\(F(x) = %d x^2\\) and distance \\(D = %d\\).\\n"+
				"Subgoal 3: Evaluate antiderivative: \\([%d x^3]_{0}^{%d} = %d \\cdot %d = %d\\) J.", coeff, D, k, D, k, D*D*D, ans),
		}
	case 2:
		// Spring stretched from x1 to x2 with F(x) = 2 k x -> W = k(x2^2 - x1^2)
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x1 := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		span := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		x2 := x1 + span
		ans := k * (x2*x2 - x1*x1)
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Find the work done (in Joules) by a force \\(F(x) = %d x\\) N in moving an object from \\(x = %d\\) to \\(x = %d\\) meters.", coeff, x1, x2),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{x_1}^{x_2} F(x)\\,dx\\).\\n"+
				"Subgoal 2: Substitute \\(F(x) = %d x\\), \\(x_1 = %d\\) and \\(x_2 = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\([%d x^2]_{%d}^{%d} = %d (%d - %d) = %d\\) J.", coeff, x1, x2, k, x1, x2, k, x2*x2, x1*x1, ans),
		}
	case 3:
		// Linear force F(x) = 2 k x + C from 0 to D -> W = k D^2 + C D
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		C := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		D := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		ans := k*D*D + C*D
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Calculate the work done by a force \\(F(x) = %d x + %d\\) N in displacing a particle from \\(x = 0\\) to \\(x = %d\\) meters.", coeff, C, D),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{0}^{D} (mx + C)\\,dx\\).\\n"+
				"Subgoal 2: Compute antiderivative \\([%d x^2 + %d x]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Evaluate: \\(%d \\cdot %d + %d \\cdot %d = %d\\) J.", k, C, D, k, D*D, C, D, ans),
		}
	default:
		// Constant force F0 over distance D -> W = F0 * D
		F0 := mathutil.ScaleInt(ctx.Difficulty, 5, 25)
		D := mathutil.ScaleInt(ctx.Difficulty, 2, 10)
		ans := F0 * D
		return generator.Problem{
			Question: fmt.Sprintf("A constant force of %d N pushes a box across a frictionless floor for a distance of %d meters. Compute the total work done.", F0, D),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = F \\cdot d\\).\\n"+
				"Subgoal 2: Substitute constant force \\(F = %d\\) N and displacement \\(d = %d\\) m.\\n"+
				"Subgoal 3: Multiply: \\(%d \\cdot %d = %d\\) J.", F0, D, F0, D, ans),
		}
	}
}

type integralCenterMassGen struct{}

func (g *integralCenterMassGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Uniform rod on [0, 2k] -> xbar = k
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 9)
		L := 2 * k
		ans := k
		return generator.Problem{
			Question: fmt.Sprintf("A thin rod of uniform density occupies the interval [0, %d] on the x-axis. Find the x-coordinate of its center of mass.", L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for a uniform rod: \\(\\bar{x} = \\frac{a + b}{2}\\).\\n"+
				"Subgoal 2: Substitute endpoints \\(a = 0\\) and \\(b = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(\\bar{x} = \\frac{0 + %d}{2} = %d\\).", L, L, ans),
		}
	case 1:
		// Non-uniform rod on [0, L] with rho(x) = k x, where L = 3 m -> xbar = 2 m
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		L := 3 * m
		ans := 2 * m
		return generator.Problem{
			Question: fmt.Sprintf("A rod of length %d on [0, %d] has linear density \\(\\rho(x) = x\\). Find its center of mass \\(\\bar{x}\\).", L, L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{M_y}{M} = \\frac{\\int_{0}^{L} x \\rho(x)\\,dx}{\\int_{0}^{L} \\rho(x)\\,dx}\\).\\n"+
				"Subgoal 2: Compute total mass \\(M = \\int_{0}^{%d} x\\,dx = \\frac{%d^2}{2}\\) and moment \\(M_y = \\int_{0}^{%d} x^2\\,dx = \\frac{%d^3}{3}\\).\\n"+
				"Subgoal 3: Simplify the ratio: \\(\\bar{x} = \\frac{2}{3} L = \\frac{2}{3} (%d) = %d\\).", L, L, L, L, L, ans),
		}
	case 2:
		// Rod on [0, L] with rho(x) = x^2, where L = 4 m -> xbar = 3 m
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		L := 4 * m
		ans := 3 * m
		return generator.Problem{
			Question: fmt.Sprintf("A rod occupying [0, %d] has density \\(\\rho(x) = x^2\\). Determine the center of mass \\(\\bar{x}\\).", L),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{\\int_{0}^{L} x^3\\,dx}{\\int_{0}^{L} x^2\\,dx}\\).\\n"+
				"Subgoal 2: Evaluate integrals: numerator is \\(L^4 / 4\\) and denominator is \\(L^3 / 3\\).\\n"+
				"Subgoal 3: Compute \\(\\bar{x} = \\frac{3}{4} L = \\frac{3}{4} (%d) = %d\\).", L, ans),
		}
	case 3:
		// Uniform rod on [a, b] with same parity -> xbar = (a + b) / 2
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		span := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := a + span
		ans := (a + b) / 2
		return generator.Problem{
			Question: fmt.Sprintf("A uniform rod extends along the x-axis from \\(x = %d\\) to \\(x = %d\\). Find the center of mass \\(\\bar{x}\\).", a, b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: by symmetry, \\(\\bar{x} = \\frac{a + b}{2}\\).\\n"+
				"Subgoal 2: Substitute endpoints \\(a = %d\\) and \\(b = %d\\).\\n"+
				"Subgoal 3: Calculate: \\(\\bar{x} = \\frac{%d + %d}{2} = %d\\).", a, b, a, b, ans),
		}
	default:
		// Two equal point masses m at x = 0 and x = 2d -> xbar = d
		d := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := d
		return generator.Problem{
			Question: fmt.Sprintf("Two equal point masses \\(m\\) are placed at positions \\(x = 0\\) and \\(x = %d\\). Find the center of mass of the system.", 2*d),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{\\sum m_i x_i}{\\sum m_i}\\).\\n"+
				"Subgoal 2: Substitute \\(m_1 = m_2 = m\\), \\(x_1 = 0\\) and \\(x_2 = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(\\bar{x} = \\frac{m(0) + m(%d)}{2m} = \\frac{%d}{2} = %d\\).", 2*d, 2*d, 2*d, ans),
		}
	}
}

type integralSurfaceAreaGen struct{}

func (g *integralSurfaceAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	if rand.Intn(2) == 0 {
		// Revolving a semicircle of radius r gives S = 4 pi r^2.
		r := randRange(ctx.Difficulty, 1, 2+scale)
		answer := fmt.Sprintf("%dpi", 4*r*r)
		return generator.Problem{
			Question: fmt.Sprintf("What is the surface area of a sphere of radius %d? (enter in terms of pi)", r),
			Answer:   answer,
			Explanation: fmt.Sprintf("Revolving a semicircle gives \\(S = 4\\pi r^{2}\\); with \\(r=%d\\): \\(4\\pi \\cdot %d = %s\\). "+
				"The general band formula is \\(S = 2\\pi\\int f(x)\\sqrt{1+(f'(x))^{2}}\\,dx\\).", r, r*r, answer),
		}
	}
	// The slant factor of a line y = m x: 1 + (f')^2 = 1 + m^2.
	m := randRange(ctx.Difficulty, 1, 2+scale)
	answer := fmt.Sprintf("%d", 1+m*m)
	return generator.Problem{
		Question: fmt.Sprintf("For \\(y=%dx\\), what is \\(1+(f'(x))^{2}\\)? (enter a number)", m),
		Answer:   answer,
		Explanation: fmt.Sprintf("\\(f'=%d\\), so \\(1+(f')^{2} = 1+%d = %s\\). "+
			"The arc-length element \\(\\sqrt{1+(f')^{2}}\\,dx\\) measures slant length in the surface formula.", m, m*m, answer),
	}
}

type integralTrapezoidGen struct{}

func (g *integralTrapezoidGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Trapezoidal approximation for f(x) = x^2 on [0, 2k] with n = 1 -> T_1 = 4 k^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 2 * k
		ans := 4 * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf("Approximate \\(\\int_{0}^{%d} x^2\\,dx\\) using the Trapezoidal Rule with \\(n = 1\\) subinterval.", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule for \\(n = 1\\): \\(T_1 = \\frac{b-a}{2} (f(a) + f(b))\\).\\n"+
				"Subgoal 2: Substitute \\(a = 0\\), \\(b = %d\\), \\(f(0) = 0\\) and \\(f(%d) = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(T_1 = \\frac{%d}{2} (0 + %d) = %d \\cdot %d = %d\\).", b, b, b*b, b, b*b, k, b*b, ans),
		}
	case 1:
		// Trapezoidal rule for linear function f(x) = 2mx + b0 on [0, L] with n = 2 -> exact integral m L^2 + b0 L
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b0 := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		L := 2 * mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		ans := m*L*L + b0*L
		coeff := 2 * m
		return generator.Problem{
			Question: fmt.Sprintf("Approximate \\(\\int_{0}^{%d} (%d x + %d)\\,dx\\) using the Trapezoidal Rule with \\(n = 2\\).", L, coeff, b0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule: \\(T_n = \\frac{\\Delta x}{2} [f(x_0) + 2f(x_1) + f(x_2)]\\).\\n"+
				"Subgoal 2: Note that the Trapezoidal Rule is exact for linear polynomials, so \\(T_2 = \\int_{0}^{%d} (%d x + %d)\\,dx\\).\\n"+
				"Subgoal 3: Evaluate the exact integral: \\([%d x^2 + %d x]_{0}^{%d} = %d\\).", L, coeff, b0, m, b0, L, ans),
		}
	case 2:
		// Trapezoidal approximation for f(x) = x^2 on [0, 2k] with n = 2 -> T_2 = 3 k^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 2 * k
		ans := 3 * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf("Use the Trapezoidal Rule with \\(n = 2\\) subintervals to approximate \\(\\int_{0}^{%d} x^2\\,dx\\).", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\Delta x = \\frac{b-a}{n} = %d\\), \\(T_2 = \\frac{%d}{2} [f(0) + 2f(%d) + f(%d)]\\).\\n"+
				"Subgoal 2: Substitute function values \\(f(0) = 0\\), \\(f(%d) = %d\\) and \\(f(%d) = %d\\).\\n"+
				"Subgoal 3: Calculate: \\(\\frac{%d}{2} [0 + 2(%d) + %d] = \\frac{%d}{2} [%d] = %d\\).", k, k, k, b, k, k*k, b, b*b, k, k*k, b*b, k, 2*k*k+b*b, ans),
		}
	case 3:
		// Subinterval count n given [a, b] and h -> n = (b - a) / h
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		n := mathutil.ScaleInt(ctx.Difficulty, 4, 12)
		b := a + n*h
		ans := n
		return generator.Problem{
			Question: fmt.Sprintf("If the interval \\([%d, %d]\\) is divided into subintervals of uniform width \\(\\Delta x = %d\\), how many subintervals \\(n\\) are created?", a, b, h),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\Delta x = \\frac{b - a}{n} \\implies n = \\frac{b - a}{\\Delta x}\\).\\n"+
				"Subgoal 2: Substitute parameters \\(a = %d\\), \\(b = %d\\) and \\(\\Delta x = %d\\).\\n"+
				"Subgoal 3: Compute: \\(n = \\frac{%d - %d}{%d} = \\frac{%d}{%d} = %d\\).", a, b, h, b, a, h, b-a, h, ans),
		}
	default:
		// Trapezoidal rule for f(x) = 2kx^2 on [0, 2] with n = 1 -> T_1 = 8k
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := 8 * k
		coeff := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("Approximate \\(\\int_{0}^{2} %d x^2\\,dx\\) using the Trapezoidal Rule with \\(n = 1\\).", coeff),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule: \\(T_1 = \\frac{b-a}{2} (f(a) + f(b))\\).\\n"+
				"Subgoal 2: For \\([0, 2]\\), \\(\\frac{b-a}{2} = 1\\), \\(f(0) = 0\\) and \\(f(2) = %d(4) = %d\\).\\n"+
				"Subgoal 3: Evaluate: \\(T_1 = 1 \\cdot (0 + %d) = %d\\).", coeff, 4*coeff, 4*coeff, ans),
		}
	}
}

type integralSimpsonGen struct{}

func (g *integralSimpsonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := selectMode(ctx, 5)
	switch mode {
	case 0:
		// Simpson's rule for cubic f(x) = x^3 on [0, 2k] with n = 2 -> exact integral (2k)^4 / 4 = 4 k^4
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 2 * k
		ans := 4 * k * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf("Approximate \\(\\int_{0}^{%d} x^3\\,dx\\) using Simpson's Rule with \\(n = 2\\).", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for polynomials of degree up to 3.\\n"+
				"Subgoal 2: Recognize that \\(S_2 = \\int_{0}^{%d} x^3\\,dx\\).\\n"+
				"Subgoal 3: Evaluate the exact integral: \\([\\frac{x^4}{4}]_{0}^{%d} = \\frac{%d^4}{4} = %d\\).", b, b, b, ans),
		}
	case 1:
		// Simpson's rule for quadratic f(x) = 3 x^2 on [0, 2k] with n = 2 -> exact integral (2k)^3 = 8 k^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		b := 2 * k
		ans := 8 * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf("Use Simpson's Rule with \\(n = 2\\) subintervals to approximate \\(\\int_{0}^{%d} 3x^2\\,dx\\).", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for quadratic polynomials.\\n"+
				"Subgoal 2: Evaluate \\(S_2 = \\int_{0}^{%d} 3x^2\\,dx = [x^3]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Compute: \\(%d^3 = %d\\).", b, b, b, ans),
		}
	case 2:
		// What is the smallest even number of subintervals n usable in Simpson's rule?
		return generator.Problem{
			Question: "What is the smallest number of subintervals \\(n\\) required to apply Simpson's Rule?",
			Answer:   "2",
			Explanation: "Subgoal 1: State the requirement for Simpson's Rule: parabolas are fit across pairs of adjacent subintervals.\\n" +
				"Subgoal 2: Deduce that \\(n\\) must be an even positive integer.\\n" +
				"Subgoal 3: Identify the smallest positive even integer: \\(n = 2\\).",
		}
	case 3:
		// Simpson's rule for linear f(x) = 2mx on [0, 2k] with n = 2 -> exact integral 4 m k^2
		m := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		b := 2 * k
		ans := 4 * m * k * k
		coeff := 2 * m
		return generator.Problem{
			Question: fmt.Sprintf("Approximate \\(\\int_{0}^{%d} %d x\\,dx\\) using Simpson's Rule with \\(n = 2\\).", b, coeff),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule \\(S_2 = \\frac{\\Delta x}{3} [f(x_0) + 4f(x_1) + f(x_2)]\\).\\n"+
				"Subgoal 2: Because the integrand is linear, Simpson's Rule gives the exact integral value \\(\\int_{0}^{%d} %d x\\,dx\\).\\n"+
				"Subgoal 3: Evaluate: \\([%d x^2]_{0}^{%d} = %d \\cdot %d = %d\\).", b, coeff, m, b, m, b*b, ans),
		}
	default:
		// Simpson's rule for f(x) = x^2 on [0, 6k] with n = 2 -> exact value 72 k^3
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 3)
		b := 6 * k
		ans := 72 * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf("Apply Simpson's Rule with \\(n = 2\\) subintervals to approximate \\(\\int_{0}^{%d} x^2\\,dx\\).", b),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for quadratic polynomials.\\n"+
				"Subgoal 2: Compute exact integral \\(\\int_{0}^{%d} x^2\\,dx = [\\frac{x^3}{3}]_{0}^{%d}\\).\\n"+
				"Subgoal 3: Calculate: \\(\\frac{%d^3}{3} = \\frac{%d}{3} = %d\\).", b, b, b, b*b*b, ans),
		}
	}
}

type integralPTestGen struct{}

func (g *integralPTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "What is the smallest integer p for which \\(\\int_{1}^{\\infty}dx/x^{p}\\) converges? (enter a number)", correct: "2", reason: "The \\(p\\)-integral converges iff \\(p>1\\): smallest integer 2.", distractors: []string{"1", "3", "0"}},
		{question: "Does \\(\\int_{1}^{\\infty}dx/x^{3}\\) converge? (yes/no)", correct: "yes", reason: "Since \\(p=3>1\\), the integral converges.", distractors: []string{"no"}},
		{question: "Does \\(\\int_{1}^{\\infty}dx/x\\) converge? (yes/no)", correct: "no", reason: "The harmonic case \\(p=1\\) diverges (logarithmic growth).", distractors: []string{"yes"}},
		{question: "Does \\(\\int_{0}^{1}dx/x^{p}\\) converge for \\(p<1\\)? (yes/no)", correct: "yes", reason: "At 0, \\(\\int_{0}^{1}dx/x^{p}\\) converges iff \\(p<1\\) (flipped condition).", distractors: []string{"no"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.integral.p_test")...))
}

type ratioTestGen struct{}

func (g *ratioTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"In the ratio test, does L<1 imply convergence? (yes/no)", "yes", "With L=lim|a(n+1)/a(n)|: L<1 converges, L>1 diverges, L=1 says nothing; L<1 converges, L>1 diverges, L=1 inconclusive."},
		{"Does sum 1/n! converge? (yes/no)", "yes", "|a(n+1)/a(n)| = 1/(n+1) → 0 < 1, so it converges (L=0)."},
		{"Does sum n!/n^n converge? (yes/no)", "yes", "The ratio tends to 1/e < 1, so it converges (L=1/e)."},
		{"Does L=1 decide convergence in the ratio test? (yes/no)", "no", "L=1 is inconclusive and says nothing: try comparison, integral, or root tests."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    "For sum 1/n!, the ratio limit L is 0. What is L? (enter a number)",
			Answer:      "0",
			Explanation: "|a(n+1)/a(n)| = 1/(n+1) → 0.",
		}
	}
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type geometricSeriesGen struct{}

func (g *geometricSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 6)
	q := randRange(ctx.Difficulty, 2, 6)
	if rand.Intn(3) == 0 {
		ratio := fmt.Sprintf("%d", q)
		if q-1 != 1 {
			ratio = fmt.Sprintf("%d/%d", q, q-1)
		}
		r := fmt.Sprintf("Does \\(\\sum_{n=0}^{\\infty}\\left(%s\\right)^{n}\\) converge?", ratio)
		res := generator.MultipleChoice(r, "no", "yes")
		res.Explanation = fmt.Sprintf("The ratio is \\(%s>1\\), so the terms grow without bound and the series diverges. A geometric series converges only when \\(|r|<1\\).", ratio)
		return res
	}
	// Sum from n = 0: a/(1 - 1/q) = aq/(q-1), reduced to lowest terms.
	num, den := a*q, q-1
	d := mathutil.GCD(num, den)
	num, den = num/d, den/d
	answer := strconv.Itoa(num)
	if den != 1 {
		answer = fmt.Sprintf("%d/%d", num, den)
	}
	stem := fmt.Sprintf("\\left(\\frac{1}{%d}\\right)^{n}", q)
	if a != 1 {
		stem = fmt.Sprintf("%d\\left(\\frac{1}{%d}\\right)^{n}", a, q)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the sum of \\(\\sum_{n=0}^{\\infty}%s\\).", stem),
		Answer:      answer,
		Explanation: fmt.Sprintf("Geometric with first term %d and ratio 1/%d: the sum is \\(\\frac{%d}{1-1/%d}=\\frac{%d}{%d}\\).", a, q, a, q, num, den),
		Facts:       map[string]string{"a": strconv.Itoa(a), "q": strconv.Itoa(q)},
	}
}

type nthTermGen struct{}

type telescopingSeriesGen struct{}

func (g *telescopingSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 1, 9)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\sum_{n=1}^{\\infty}\\left(\\frac{%d}{n}-\\frac{%d}{n+1}\\right)\\).", k, k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("Every partial sum collapses to \\(%d-\\frac{%d}{N+1}\\), which tends to %d.", k, k, k),
		}
	case 1:
		a := randRange(ctx.Difficulty, 1, 5)
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\sum_{n=1}^{\\infty}\\left(\\frac{1}{n+%d}-\\frac{1}{n+%d}\\right)\\).", a, a+1),
			Answer:      fmt.Sprintf("1/%d", a+1),
			Explanation: fmt.Sprintf("The partial sums collapse to \\(\\frac{1}{%d}-\\frac{1}{N+%d}\\), so the series has value \\(1/%d\\).", a+1, a+1, a+1),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\sum_{n=1}^{\\infty}\\frac{%d}{n(n+1)}\\) via partial fractions.", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(\\frac{%d}{n(n+1)}=\\frac{%d}{n}-\\frac{%d}{n+1}\\), and the partial sums telescope to %d.", k, k, k, k),
		}
	default:
		return generator.Problem{
			Question:    "Do the interior terms cancel in a telescoping partial sum? (yes/no)",
			Answer:      "yes",
			Explanation: "Writing each term as a difference makes consecutive terms cancel, so only the first and last survive.",
		}
	}
}

type directCompareGen struct{}

func (g *directCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 2, 5)
	c := randRange(ctx.Difficulty, 1, 9)
	if rand.Intn(4) == 0 {
		res := generator.MultipleChoice(
			"Which test proves convergence by bounding each term above by a larger convergent series?",
			"direct comparison", "ratio test", "root test", "integral test")
		res.Explanation = "Direct comparison needs a second series you already know: \\(0\\le a_n\\le b_n\\) with \\(\\sum b_n\\) convergent forces \\(\\sum a_n\\) to converge."
		return res
	}
	if rand.Intn(3) == 0 {
		// Divergent comparison: the harmonic-like series is the lower bound.
		var q string
		if rand.Intn(2) == 0 {
			q = fmt.Sprintf("Does \\(\\sum_{n=1}^{\\infty}\\frac{1}{n+%d}\\) converge?", c)
		} else {
			q = fmt.Sprintf("Does \\(\\sum_{n=1}^{\\infty}\\frac{1}{\\sqrt{n}+%d}\\) converge?", c)
		}
		res := generator.MultipleChoice(q, "no", "yes")
		res.Explanation = fmt.Sprintf("The denominator grows like \\(n\\) (or slower), so the terms dominate a shifted harmonic series, which diverges. Adding %d does not change that.", c)
		return res
	}
	q := fmt.Sprintf("Does \\(\\sum_{n=1}^{\\infty}\\frac{1}{n^{%d}+%d}\\) converge?", p, c)
	res := generator.MultipleChoice(q, "yes", "no")
	res.Explanation = fmt.Sprintf("For \\(n\\ge 1\\) we have \\(\\frac{1}{n^{%d}+%d}<\\frac{1}{n^{%d}}\\), and \\(\\sum 1/n^{%d}\\) converges because \\(%d>1\\). Direct comparison gives convergence.", p, c, p, p, p)
	return res
}

type limitCompareGen struct{}

func (g *limitCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 1, 3)
	d := randRange(ctx.Difficulty, 1, 3)
	m := k + d
	c := randRange(ctx.Difficulty, 1, 9)
	if rand.Intn(3) == 0 {
		// Ask for the comparison exponent directly; a small integer is typeable.
		return generator.Problem{
			Question:    fmt.Sprintf("For \\(\\sum_{n=1}^{\\infty}\\frac{n^{%d}}{n^{%d}+%d}\\), limit comparison with \\(\\sum 1/n^{p}\\) uses which \\(p\\)?", k, m, c),
			Answer:      strconv.Itoa(d),
			Explanation: fmt.Sprintf("The numerator has degree %d and the denominator degree %d, so the terms behave like \\(1/n^{%d}\\): \\(p=%d\\).", k, m, d, d),
			Facts:       map[string]string{"k": strconv.Itoa(k), "m": strconv.Itoa(m), "p": strconv.Itoa(d)},
		}
	}
	ans := "no"
	why := fmt.Sprintf("The terms behave like \\(1/n^{%d}\\), and a p-series diverges when \\(p\\le 1\\).", d)
	if d > 1 {
		ans = "yes"
		why = fmt.Sprintf("The terms behave like \\(1/n^{%d}\\), and \\(\\sum 1/n^{%d}\\) converges because \\(%d>1\\), so the limit comparison test gives convergence.", d, d, d)
	}
	q := fmt.Sprintf("Does \\(\\sum_{n=1}^{\\infty}\\frac{n^{%d}}{n^{%d}+%d}\\) converge?", k, m, c)
	res := generator.MultipleChoice(q, ans, "yes", "no")
	res.Explanation = why
	return res
}

type absoluteConvergenceGen struct{}

type alternatingRemainderGen struct{}

func (g *alternatingRemainderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := randRange(ctx.Difficulty, 1, 3)
	n := randRange(ctx.Difficulty, 4, 20)
	first := n + 1
	den := intPow(first, p)
	if p == 1 {
		return generator.Problem{
			Question:    fmt.Sprintf("Bound the error from stopping \\(\\sum (-1)^{n}/n\\) after \\(n=%d\\).", n),
			Answer:      fmt.Sprintf("1/%d", first),
			Explanation: fmt.Sprintf("For an alternating series the error is at most the first omitted term: \\(a_{%d}=1/%d\\).", first, first),
		}
	}
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Bound the error from stopping \\(\\sum (-1)^{n}/n^{%d}\\) after \\(n=%d\\).", p, n),
			Answer:      fmt.Sprintf("1/%d", den),
			Explanation: fmt.Sprintf("The first omitted term bounds the error: \\(1/%d^{%d}=1/%d\\).", first, p, den),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("An alternating series is stopped after the term with \\(n=%d\\). Which term number bounds the truncation error? (enter a number)", n),
		Answer:      strconv.Itoa(first),
		Explanation: fmt.Sprintf("The error is at most the first dropped term, which is \\(a_{%d}\\).", first),
	}
}

type convergenceRadiusGen struct{}

func (g *convergenceRadiusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	c := randRange(ctx.Difficulty, 2, 9)
	switch rand.Intn(6) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the radius of convergence of \\(\\sum x^{n}/%d^{n}\\)?", c),
			Answer:      strconv.Itoa(c),
			Explanation: fmt.Sprintf("The ratio test gives \\(|x|/%d<1\\), so \\(R=%d\\).", c, c),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the radius of convergence of \\(\\sum %d^{n}x^{n}\\)?", c),
			Answer:      fmt.Sprintf("1/%d", c),
			Explanation: fmt.Sprintf("The ratio test gives \\(%d|x|<1\\), so \\(R=1/%d\\).", c, c),
		}
	case 2:
		return generator.Problem{
			Question:    "What is the radius of convergence of \\(\\sum x^{n}\\)?",
			Answer:      "1",
			Explanation: "The ratio test gives \\(|x|<1\\), so \\(R=1\\).",
		}
	case 3:
		return generator.Problem{
			Question:    "What is the radius of convergence of \\(\\sum x^{n}/n!\\)?",
			Answer:      "infinity",
			Explanation: "The ratio test gives \\(L=0\\) for every \\(x\\), so \\(R=\\infty\\).",
		}
	case 4:
		return generator.Problem{
			Question:    "What is the radius of convergence of \\(\\sum n!x^{n}\\)?",
			Answer:      "0",
			Explanation: "The ratio test diverges for every \\(x\\neq 0\\), so \\(R=0\\).",
		}
	default:
		return generator.Problem{
			Question:    "What does the interval of convergence add to the radius?",
			Answer:      "endpoint checks",
			Explanation: "The radius gives the open interval; each endpoint must be tested separately.",
		}
	}
}

type epsilonDeltaGen struct{}

func (g *epsilonDeltaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 2, 9)
	a := randRange(ctx.Difficulty, 1, 5)
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Prove \\(\\lim_{x\\to %d}%dx = %d\\): what number \\(N\\) satisfies \\(\\delta=\\varepsilon/N\\)? (enter a number)", a, k, k*a),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(|%dx-%d| = %d|x-%d| < \\varepsilon\\) when \\(|x-%d|<\\varepsilon/%d\\), so \\(N=%d\\).", k, k*a, k, a, a, k, k),
		}
	case 1:
		return generator.Problem{
			Question:    "Does the epsilon-delta definition require a delta for every epsilon? (yes/no)",
			Answer:      "yes",
			Explanation: "For every tolerance there is a distance keeping the output within it: given \\(\\varepsilon>0\\) there is \\(\\delta>0\\) with \\(0<|x-x_0|<\\delta\\) implying \\(|f(x)-L|<\\varepsilon\\).",
		}
	default:
		return generator.Problem{
			Question:    "In epsilon-delta, what does delta control?",
			Answer:      "how close x must be to x_0",
			Explanation: "\\(\\delta\\) bounds the input distance while \\(\\varepsilon\\) bounds the output error.",
		}
	}
}

type oneSidedLimitGen struct{}

func (g *oneSidedLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 6)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d^{+}}\\frac{|x-%d|}{x-%d}\\).", a, a, a),
			Answer:      "1",
			Explanation: fmt.Sprintf("For \\(x>%d\\) the numerator equals the denominator, so the quotient is 1.", a),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d^{-}}\\frac{|x-%d|}{x-%d}\\).", a, a, a),
			Answer:      "-1",
			Explanation: fmt.Sprintf("For \\(x<%d\\) the numerator is the negative of the denominator, so the quotient is -1.", a),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d}\\frac{x^{2}-%d}{x-%d}\\).", a, a*a, a),
			Answer:      strconv.Itoa(2 * a),
			Explanation: fmt.Sprintf("Factor: \\(\\frac{(x-%d)(x+%d)}{x-%d}=x+%d\\), which tends to %d.", a, a, a, a, 2*a),
		}
	default:
		return generator.Problem{
			Question:    "How many one-sided limits must agree for a two-sided limit to exist? (enter a number)",
			Answer:      "2",
			Explanation: "\\(\\lim_{x\\to a}f\\) exists exactly when both one-sided limits exist and are equal.",
		}
	}
}

type infiniteLimitGen struct{}

func (g *infiniteLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := randRange(ctx.Difficulty, 1, 6)
	c := randRange(ctx.Difficulty, 1, 5)
	switch rand.Intn(5) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d}\\frac{%d}{(x-%d)^{2}}\\).", a, c, a),
			Answer:      "infinity",
			Explanation: fmt.Sprintf("\\(%d/(x-%d)^{2}\\) grows without bound from both sides, so the limit is \\(+\\infty\\).", c, a),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d^{+}}\\frac{%d}{x-%d}\\).", a, c, a),
			Answer:      "infinity",
			Explanation: fmt.Sprintf("For \\(x>%d\\) the denominator is positive and tiny, so the quotient tends to \\(+\\infty\\).", a),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to %d^{-}}\\frac{%d}{x-%d}\\).", a, c, a),
			Answer:      "-infinity",
			Explanation: fmt.Sprintf("For \\(x<%d\\) the denominator is negative and tiny, so the quotient tends to \\(-\\infty\\).", a),
		}
	case 3:
		return generator.Problem{
			Question:    fmt.Sprintf("What asymptote does \\(\\lim_{x\\to %d}\\frac{%d}{x-%d}=\\pm\\infty\\) signal? (enter vertical or horizontal)", a, c, a),
			Answer:      "vertical",
			Explanation: fmt.Sprintf("An infinite limit at a finite input means \\(x=%d\\) is a vertical asymptote.", a),
		}
	default:
		return generator.Problem{
			Question:    "Find \\(\\lim_{x\\to 0}\\frac{1}{x^{2}}\\).",
			Answer:      "infinity",
			Explanation: "\\(1/x^{2}\\) tends to \\(+\\infty\\) from both sides.",
		}
	}
}

type sineLimitGen struct{}

func (g *sineLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 1, 9)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{\\sin(%dx)}{x}\\).", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(\\frac{\\sin(%dx)}{x} = %d\\cdot\\frac{\\sin(%dx)}{%dx} \\to %d\\).", k, k, k, k, k),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{\\tan(%dx)}{x}\\).", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(\\frac{\\tan(%dx)}{x} = %d\\cdot\\frac{\\sin(%dx)}{%dx}\\cdot\\frac{1}{\\cos(%dx)} \\to %d\\).", k, k, k, k, k, k),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{1-\\cos(%dx)}{x^{2}}\\).", k),
			Answer:      decimal(float64(k*k)/2, 1),
			Explanation: fmt.Sprintf("\\(1-\\cos(%dx) \\approx \\frac{(%dx)^{2}}{2}\\), so the quotient tends to \\(%d/2\\).", k, k, k*k),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{x}{\\sin(%dx)}\\).", k),
			Answer:      fmt.Sprintf("1/%d", k),
			Explanation: fmt.Sprintf("\\(\\frac{x}{\\sin(%dx)} = \\frac{1}{%d}\\cdot\\frac{%dx}{\\sin(%dx)} \\to 1/%d\\).", k, k, k, k, k),
		}
	}
}

type expLimitGen struct{}

func (g *expLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	k := randRange(ctx.Difficulty, 1, 9)
	switch rand.Intn(4) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{e^{%dx}-1}{x}\\).", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(\\frac{e^{%dx}-1}{x} = %d\\cdot\\frac{e^{%dx}-1}{%dx} \\to %d\\).", k, k, k, k, k),
		}
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{x\\to 0}\\frac{\\ln(1+%dx)}{x}\\).", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("\\(\\frac{\\ln(1+%dx)}{x} = %d\\cdot\\frac{\\ln(1+%dx)}{%dx} \\to %d\\).", k, k, k, k, k),
		}
	case 2:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{n\\to\\infty}\\left(1+\\frac{%d}{n}\\right)^{n}\\).", k),
			Answer:      fmt.Sprintf("e^%d", k),
			Explanation: fmt.Sprintf("\\(\\left(1+\\frac{%d}{n}\\right)^{n} \\to e^{%d}\\).", k, k),
		}
	default:
		return generator.Problem{
			Question:    fmt.Sprintf("Find \\(\\lim_{n\\to\\infty}\\left(1+\\frac{1}{n}\\right)^{%dn}\\).", k),
			Answer:      fmt.Sprintf("e^%d", k),
			Explanation: fmt.Sprintf("The inner bracket tends to \\(e\\), so the %dth power tends to \\(e^{%d}\\).", k, k),
		}
	}
}

type boundedSeqGen struct{}

func (g *boundedSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claims := []conceptClaim{
		{question: "Is \\(a_{n}=(-1)^{n}\\) bounded? (yes/no)", correct: "yes", reason: "It oscillates between -1 and 1: bounded (|a_n|≤1) but divergent.", distractors: []string{"no"}},
		{question: "Does every convergent sequence have to be bounded?", correct: "yes", reason: "Convergence forces boundedness (terms eventually stay near the limit).", distractors: []string{"no"}},
		{question: "Is \\(a_{n}=n^{2}\\) bounded?", correct: "no", reason: "It grows without bound.", distractors: []string{"yes"}},
		{question: "Is a bounded sequence always convergent?", correct: "no", reason: "\\((-1)^{n}\\) is bounded yet oscillates: boundedness is necessary, not sufficient.", distractors: []string{"yes"}},
	}
	return conceptPickMC(append(claims, computedConceptClaims("calc.seq.bounded")...))
}

type convergenceGen struct{}

func (g *convergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"The integral \\(\\int_{1}^{\\infty} \\frac{\\sin x}{x^{2}} \\, dx\\) converges absolutely, conditionally, or diverges? (enter absolute, conditional, or diverges)", "absolute", "For \\(p=2>1\\), \\(|\\sin x|/x^{2} \\le 1/x^{2}\\) and the p-integral converges: absolute convergence."},
		{"The integral \\(\\int_{1}^{\\infty} \\frac{\\sin x}{\\sqrt{x}} \\, dx\\) converges absolutely, conditionally, or diverges? (enter absolute, conditional, or diverges)", "conditional", "Dirichlet: \\(\\int_{1}^{t}\\sin x\\,dx\\) is bounded by 2 and \\(x^{-1/2}\\) decreases to 0, so it converges; but \\(p=1/2<1\\) so the absolute integral diverges."},
		{"The integral \\(\\int_{1}^{\\infty} \\sin x \\, dx\\) converges absolutely, conditionally, or diverges? (enter absolute, conditional, or diverges)", "diverges", "Here \\(p=0\\le 0\\): the Cauchy criterion fails on the humps \\(I_{k}=[2k\\pi+\\pi/6,2k\\pi+5\\pi/6]\\), each contributing at least \\(\\pi/3\\)."},
		{"What bounds \\(\\left|\\int_{1}^{t} \\sin x \\, dx\\right|\\) for all \\(t\\)? (enter a number)", "2", "\\(\\left|\\int_{1}^{t}\\sin x\\,dx\\right| = |\\cos 1 - \\cos t| \\le 2\\): the bounded accumulation Dirichlet's test needs."},
		{"The Cauchy principal value of \\(\\int_{-\\infty}^{\\infty} \\frac{2x}{1+x^{2}} \\, dx\\) is what? (enter a number)", "0", "Symmetric bounds give \\(\\int_{-R}^{R} = 0\\) for this odd function, so the principal value is 0 (the improper integral itself diverges)."},
	}
	tableHard := []entry{
		{"The integral of \\(|\\sin x|/x^{p}\\) over \\([1,\\infty)\\) converges for \\(p\\) above what threshold? (enter a number)", "1", "Absolute convergence needs \\(p>1\\) by comparison with \\(1/x^{p}\\); at \\(p\\le 1\\) the sine-squared split shows divergence."},
		{"In Abel's test, \\(g\\) must be monotone and what? (enter bounded or vanishing)", "bounded", "Abel drops Dirichlet's \\(g\\to 0\\) but keeps monotonicity and adds boundedness, writing \\(g = h + L\\) with \\(h\\to 0\\)."},
		{"The Fresnel integral \\(\\int_{0}^{\\infty} \\sin(x^{2}) \\, dx\\) converges absolutely or conditionally? (enter absolute or conditional)", "conditional", "With \\(t=x^{2}\\) it becomes half of the \\(p=1/2\\) sine family: convergent, but the absolute integral diverges."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type irrationalIntegralGen struct{}

func (g *irrationalIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"What substitution rationalizes \\(\\int \\frac{dx}{\\sqrt{x}+\\sqrt[3]{x}}\\)? (enter as x = ...)", "x = t^6", "Indices 2 and 3 have lcm 6: with \\(x=t^{6}\\), \\(\\sqrt{x}=t^{3}\\) and \\(\\sqrt[3]{x}=t^{2}\\), and \\(dx=6t^{5}dt\\)."},
		{"With \\(x=t^{6}\\), \\(dx\\) equals what? (enter in terms of t, with dt)", "6t^5 dt", "Differentiating \\(x=t^{6}\\) gives \\(dx = 6t^{5}dt\\)."},
		{"What exponent \\(m\\) rationalizes \\(\\sqrt{ax+b}\\) and \\(\\sqrt[3]{ax+b}\\) together? (enter a number)", "6", "Take \\(m\\) as the lcm of the indices: \\(m=6\\) makes both roots integer powers of \\(t\\)."},
		{"Which Euler substitution applies when \\(a>0\\)? (enter first, second, or third)", "first", "For \\(a>0\\), set \\(\\sqrt{ax^{2}+bx+c} = t - x\\sqrt{a}\\): squaring cancels \\(ax^{2}\\) and leaves \\(x\\) rational in \\(t\\)."},
	}
	tableHard := []entry{
		{"Which Euler substitution needs \\(c>0\\)? (enter first, second, or third)", "second", "For \\(c>0\\), set \\(\\sqrt{ax^{2}+bx+c} = xt+\\sqrt{c}\\): squaring cancels the constant term."},
		{"Which Euler substitution needs two distinct real roots? (enter first, second, or third)", "third", "With roots \\(\\alpha,\\beta\\), set \\(\\sqrt{ax^{2}+bx+c} = t(x-\\alpha)\\): dividing by \\(x-\\alpha\\) leaves a first-degree equation."},
		{"What is the antiderivative of \\(\\frac{1}{\\sqrt{x^{2}+1}}\\)? (enter it)", "ln(x+sqrt(x^2+1))", "The \\(a>0\\) completing-the-square formula gives \\(\\ln|x+\\sqrt{x^{2}+1}|\\)."},
		{"For \\(\\int \\frac{x^{2}}{\\sqrt{x^{2}+1}} \\, dx\\), the undetermined-coefficient polynomial \\(S_{1}(x)\\) has what degree? (enter a number)", "1", "The ansatz is \\(S_{n-1}\\sqrt{Q} + \\lambda\\int dx/\\sqrt{Q}\\): for numerator degree 2, \\(S_{1}(x)=Ax+B\\) has degree 1."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type strategyGen struct{}

func (g *strategyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Which technique for \\(\\int x e^{x^{2}} \\, dx\\)? (enter substitution, parts, or partial-fractions)", "substitution", "The integrand is \\(f(g)g'\\): \\(u=x^{2}\\) gives \\(\\int e^{u}du/2\\). An inner function with its derivative means substitution."},
		{"Which technique for \\(\\int x e^{x} \\, dx\\)? (enter substitution, parts, or partial-fractions)", "parts", "A product of unrelated families (algebraic times exponential) with no composite derivative: LIATE picks \\(u=x\\)."},
		{"Which technique for \\(\\int \\frac{1}{x(x+1)} \\, dx\\)? (enter substitution, parts, or partial-fractions)", "partial-fractions", "A proper rational function: split \\(1/(x(x+1)) = 1/x - 1/(x+1)\\) and integrate term by term."},
		{"Which technique for \\(\\int \\sin^{3}x \\cos x \\, dx\\)? (enter substitution, parts, or partial-fractions)", "substitution", "Composite power form: \\(u=\\sin x\\) gives \\(\\int u^{3}du = \\sin^{4}x/4\\)."},
	}
	tableHard := []entry{
		{"Which technique for \\(\\int \\arctan x \\, dx\\)? (enter substitution, parts, or partial-fractions)", "parts", "A single inverse-trig factor is a product with 1: take \\(u=\\arctan x\\), \\(dv=dx\\)."},
		{"Which technique for \\(\\int \\frac{x+5}{x^{2}-x-6} \\, dx\\)? (enter substitution, parts, or partial-fractions)", "partial-fractions", "The denominator factors as \\((x-3)(x+2)\\) with positive discriminant: two logarithmic terms."},
		{"Which technique for \\(\\int \\frac{e^{\\arctan x}}{1+x^{2}} \\, dx\\)? (enter substitution, parts, or partial-fractions)", "substitution", "The derivative of \\(\\arctan x\\) is present: \\(\\int e^{u}du = e^{u}\\) with \\(u=\\arctan x\\)."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type reductionGen struct{}

func (g *reductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"The recurrence \\(J_{n} = \\frac{x^{n}e^{ax}}{a} - \\frac{n}{a}J_{n-1}\\) lowers \\(n\\) by how much? (enter a number)", "1", "Each integration by parts differentiates \\(x^{n}\\) once: the index drops by 1 until \\(J_{0}=e^{ax}/a\\)."},
		{"The recurrence \\(S_{n} = -\\frac{\\sin^{n-1}x\\cos x}{n} + \\frac{n-1}{n}S_{n-2}\\) lowers \\(n\\) by how much? (enter a number)", "2", "The trigonometric identity rejoins the family two indices down; even and odd chains end at \\(S_{0}=x\\) and \\(S_{1}=-\\cos x\\)."},
		{"What is \\(S_{0} = \\int \\sin^{0}x \\, dx\\)? (enter it)", "x", "The empty product integrates to \\(x\\): base case of the even sine chain."},
		{"What is \\(T_{1} = \\int \\tan x \\, dx\\)? (enter it)", "-ln|cos(x)|", "Base case of the odd tangent chain: \\(\\int\\tan x\\,dx = -\\ln|\\cos x|\\)."},
	}
	tableHard := []entry{
		{"Evaluate \\(\\int_{0}^{\\pi/2} \\sin^{4}x \\, dx\\). (enter in terms of pi)", "3*pi/16", "Wallis recurrence: \\(\\frac{3}{4}\\cdot\\frac{1}{2}\\cdot\\frac{\\pi}{2} = 3\\pi/16\\)."},
		{"Evaluate \\(\\int_{0}^{\\pi/2} \\sin^{5}x \\, dx\\). (enter a fraction)", "8/15", "Wallis recurrence down the odd chain: \\(\\frac{4}{5}\\cdot\\frac{2}{3} = 8/15\\)."},
		{"What is \\(I_{1} = \\int \\frac{dx}{x^{2}+a^{2}}\\)? (enter in terms of arctan)", "(1/a)*arctan(x/a)", "Base case of the quadratic-denominator family: \\(I_{1} = \\frac{1}{a}\\arctan\\frac{x}{a}\\)."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type discMethodGen struct{}

func (g *discMethodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	if scale <= 3 {
		switch rand.Intn(6) {
		case 0:
			r := rand.Intn(4) + 1
			h := rand.Intn(5) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("The region under \\(y=%d\\) from \\(x=0\\) to \\(x=%d\\) revolves about the \\(x\\)-axis (a cylinder). What is \\(V/\\pi\\)? (enter a number)", r, h),
				Answer:      fmt.Sprintf("%d", r*r*h),
				Explanation: fmt.Sprintf("Disc radius %d gives V = π∫₀^%d %d²dx = π·%d·%d, so V/π = %d.", r, h, r, r*r, h, r*r*h),
			}
		case 1:
			b := (rand.Intn(3) + 1) * 2
			return generator.Problem{
				Question:    fmt.Sprintf("The region under \\(y=\\sqrt{x}\\) from \\(x=0\\) to \\(x=%d\\) revolves about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a number)", b),
				Answer:      fmt.Sprintf("%d", b*b/2),
				Explanation: fmt.Sprintf("V = π∫₀^%d x\\,dx = π[%d²/2] = π·%d, so V/π = %d.", b, b, b*b/2, b*b/2),
			}
		case 2:
			n := rand.Intn(5) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("Disc method with \\(f(x)=x^{2}\\). What is the disc radius at \\(x=%d\\)? (enter a number)", n),
				Answer:      fmt.Sprintf("%d", n*n),
				Explanation: fmt.Sprintf("Radius is f(%d) = %d² = %d.", n, n, n*n),
			}
		case 3:
			r := rand.Intn(6) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("A disc-method cross section has radius %d. What is its area divided by \\(\\pi\\)? (enter a number)", r),
				Answer:      fmt.Sprintf("%d", r*r),
				Explanation: fmt.Sprintf("Area πr² with r=%d gives area/π = %d² = %d.", r, r, r*r),
			}
		case 4:
			c := rand.Intn(4) + 1
			h := rand.Intn(4) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("The segment \\(x=%d\\) for \\(0 \\leq y \\leq %d\\) revolves about the \\(y\\)-axis. What is \\(V/\\pi\\)? (enter a number)", c, h),
				Answer:      fmt.Sprintf("%d", c*c*h),
				Explanation: fmt.Sprintf("Discs of radius %d stacked over height %d: V = π·%d²·%d, so V/π = %d.", c, h, c, h, c*c*h),
			}
		default:
			c := rand.Intn(5) + 1
			l := rand.Intn(4) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("The region under \\(f(x)=-%d\\) from \\(x=0\\) to \\(x=%d\\) revolves about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a number)", c, l),
				Answer:      fmt.Sprintf("%d", c*c*l),
				Explanation: fmt.Sprintf("Radius is |f| = %d, so V = π·%d²·%d and V/π = %d.", c, c, l, c*c*l),
			}
		}
	}
	type entry struct {
		q, a, e string
	}
	c := rand.Intn(3) + 1
	tableHard := []entry{
		{"Revolve \\(y=x^{2}\\) from \\(x=0\\) to \\(x=2\\) about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a fraction)", "32/5", "V = π∫₀²x⁴dx = π[32/5], so V/π = 32/5."},
		{"Revolve \\(y=x^{3}\\) from \\(x=0\\) to \\(x=1\\) about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a fraction)", "1/7", "V = π∫₀¹x⁶dx = π[1/7], so V/π = 1/7."},
		{fmt.Sprintf("The region under \\(y=x+%d\\) from \\(x=0\\) to \\(x=1\\) revolves about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a fraction)", c), fmt.Sprintf("%d/3", (c+1)*(c+1)*(c+1)-c*c*c), fmt.Sprintf("Radius x+%d: V = π∫₀¹(x+%d)²dx = π((%d)³-(%d)³)/3 = %d/3.", c, c, c+1, c, (c+1)*(c+1)*(c+1)-c*c*c)},
	}
	e := tableHard[rand.Intn(len(tableHard))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type washerMethodGen struct{}

func (g *washerMethodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	if scale <= 3 {
		switch rand.Intn(6) {
		case 0:
			bigR := rand.Intn(4) + 2
			smallR := rand.Intn(bigR-1) + 1
			l := rand.Intn(4) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("Washers have outer radius %d, inner radius %d, over length %d. What is \\(V/\\pi\\)? (enter a number)", bigR, smallR, l),
				Answer:      fmt.Sprintf("%d", (bigR*bigR-smallR*smallR)*l),
				Explanation: fmt.Sprintf("V = π(%d²-%d²)·%d = π·%d, so V/π = %d.", bigR, smallR, l, (bigR*bigR-smallR*smallR)*l, (bigR*bigR-smallR*smallR)*l),
			}
		case 1:
			bigR := rand.Intn(5) + 2
			smallR := rand.Intn(bigR-1) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("A washer has outer radius %d and inner radius %d. What is its area divided by \\(\\pi\\)? (enter a number)", bigR, smallR),
				Answer:      fmt.Sprintf("%d", bigR*bigR-smallR*smallR),
				Explanation: fmt.Sprintf("Annulus π(%d²-%d²) = π·%d, so area/π = %d.", bigR, smallR, bigR*bigR-smallR*smallR, bigR*bigR-smallR*smallR),
			}
		case 2:
			return generator.Problem{
				Question:    "Revolve the region between \\(y=x\\) and \\(y=x^{2}\\) on \\([0,1]\\) about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a fraction)",
				Answer:      "2/15",
				Explanation: "V = π∫₀¹(x²-x⁴)dx = π(1/3-1/5) = 2π/15, so V/π = 2/15.",
			}
		case 3:
			n := rand.Intn(3) + 2
			return generator.Problem{
				Question:    fmt.Sprintf("Washers use outer \\(f(x)=x^{2}\\) and inner \\(g(x)=x\\). What is \\(f^{2}-g^{2}\\) at \\(x=%d\\)? (enter a number)", n),
				Answer:      fmt.Sprintf("%d", n*n*n*n-n*n),
				Explanation: fmt.Sprintf("f(%d)²-g(%d)² = %d⁴-%d² = %d-%d = %d.", n, n, n, n, n*n*n*n, n*n, n*n*n*n-n*n),
			}
		case 4:
			h := 2
			ans := "4/3"
			if rand.Intn(2) == 0 {
				h = 4
				ans = "32/3"
			}
			return generator.Problem{
				Question:    fmt.Sprintf("A spherical napkin ring has height %d. Its volume is \\(\\pi h^{3}/6\\). What is \\(V/\\pi\\)? (enter a fraction)", h),
				Answer:      ans,
				Explanation: fmt.Sprintf("V/π = %d³/6 = %s.", h, ans),
			}
		default:
			bigR := rand.Intn(5) + 1
			l := rand.Intn(4) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("Washers with inner radius 0 (a disc), outer radius %d, over length %d. What is \\(V/\\pi\\)? (enter a number)", bigR, l),
				Answer:      fmt.Sprintf("%d", bigR*bigR*l),
				Explanation: fmt.Sprintf("Inner radius 0 degenerates to a disc: V = π·%d²·%d, so V/π = %d.", bigR, l, bigR*bigR*l),
			}
		}
	}
	type entry struct {
		q, a, e string
	}
	c := rand.Intn(3) + 1
	tableHard := []entry{
		{"Rotate the region between \\(y=x\\) and \\(y=x^{2}\\) on \\([0,1]\\) about the line \\(y=-1\\). What is \\(V/\\pi\\)? (enter a fraction)", "7/15", "Outer x+1, inner x²+1: V = π∫₀¹((x+1)²-(x²+1)²)dx = 7π/15, so V/π = 7/15."},
		{fmt.Sprintf("Rotate the region between \\(y=x+%d\\) and \\(y=x\\) on \\([0,1]\\) about the \\(x\\)-axis. What is \\(V/\\pi\\)? (enter a number)", c), fmt.Sprintf("%d", c+c*c), fmt.Sprintf("Outer (x+%d)² minus inner x² = %dx+%d: V = π∫₀¹(%dx+%d)dx = π(%d), so V/π = %d.", c, 2*c, c*c, 2*c, c*c, c+c*c, c+c*c)},
		{"Describe \\(x=y\\) and \\(x=y^{2}\\) on \\([0,1]\\) rotated about the \\(y\\)-axis with washers. What is \\(V/\\pi\\)? (enter a fraction)", "2/15", "V = π∫₀¹(y²-(y²)²)dy = π(1/3-1/5) = 2π/15, so V/π = 2/15."},
	}
	e := tableHard[rand.Intn(len(tableHard))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type shellMethodGen struct{}

func (g *shellMethodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	if scale <= 3 {
		switch rand.Intn(6) {
		case 0:
			c := rand.Intn(4) + 1
			a := rand.Intn(3)
			b := a + rand.Intn(3) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("The region under \\(y=%d\\) on \\([%d,%d]\\) revolves about the \\(y\\)-axis (shells). What is \\(V/\\pi\\)? (enter a number)", c, a, b),
				Answer:      fmt.Sprintf("%d", c*(b*b-a*a)),
				Explanation: fmt.Sprintf("V = 2π∫_{%d}^{%d} x·%d\\,dx = 2π·%d·(%d-%d)/2 = π·%d, so V/π = %d.", a, b, c, c, b*b, a*a, c*(b*b-a*a), c*(b*b-a*a)),
			}
		case 1:
			n := rand.Intn(4) + 1
			num := n * n * n * n
			ans := fmt.Sprintf("%d", num/2)
			if num%2 == 1 {
				ans = fmt.Sprintf("%d/2", num)
			}
			return generator.Problem{
				Question:    fmt.Sprintf("The region under \\(y=x^{2}\\) on \\([0,%d]\\) revolves about the \\(y\\)-axis (shells). What is \\(V/\\pi\\)? (enter a fraction)", n),
				Answer:      ans,
				Explanation: fmt.Sprintf("V = 2π∫₀^%d x³\\,dx = 2π·%d⁴/4 = π·%d⁴/2, so V/π = %s.", n, n, n, ans),
			}
		case 2:
			r := rand.Intn(5) + 1
			h := rand.Intn(5) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("A cylindrical shell has radius %d and height %d. What is \\(dV/(\\pi\\,dx)\\)? (enter a number)", r, h),
				Answer:      fmt.Sprintf("%d", 2*r*h),
				Explanation: fmt.Sprintf("dV = 2π·radius·height·dx = 2π·%d·%d·dx, so dV/(π\\,dx) = %d.", r, h, 2*r*h),
			}
		case 3:
			a := rand.Intn(5) + 1
			return generator.Problem{
				Question:    fmt.Sprintf("The vertical strip at \\(x=%d\\) runs from \\(y=x\\) up to \\(y=4x\\). What is its height? (enter a number)", a),
				Answer:      fmt.Sprintf("%d", 3*a),
				Explanation: fmt.Sprintf("Height = top minus bottom = 4·%d-%d = %d.", a, a, 3*a),
			}
		case 4:
			return generator.Problem{
				Question:    "Revolve the region between \\(y=x\\) and \\(y=x^{2}\\) on \\([0,1]\\) about the \\(y\\)-axis (shells). What is \\(V/\\pi\\)? (enter a fraction)",
				Answer:      "1/6",
				Explanation: "V = 2π∫₀¹x(x-x²)dx = 2π(1/3-1/4) = π/6, so V/π = 1/6.",
			}
		default:
			a := rand.Intn(5) + 2
			return generator.Problem{
				Question:    fmt.Sprintf("A vertical strip at \\(x=%d\\) revolves about the \\(y\\)-axis. What is the shell radius? (enter a number)", a),
				Answer:      fmt.Sprintf("%d", a),
				Explanation: fmt.Sprintf("Shell radius is the distance from the strip to the axis: %d.", a),
			}
		}
	}
	type entry struct {
		q, a, e string
	}
	c := rand.Intn(3) + 2
	tableHard := []entry{
		{"Rotate the region between \\(y=x\\) and \\(y=x^{2}\\) on \\([0,1]\\) about the line \\(x=2\\) (shells). What is \\(V/\\pi\\)? (enter a fraction)", "1/2", "Radius 2-x: V = 2π∫₀¹(2-x)(x-x²)dx = π/2, so V/π = 1/2."},
		{"Rotate the region under \\(y=2x^{2}-x^{3}\\) on \\([0,2]\\) about the \\(y\\)-axis (shells). What is \\(V/\\pi\\)? (enter a fraction)", "16/5", "V = 2π∫₀²x(2x²-x³)dx = 2π(8-32/5) = 16π/5, so V/π = 16/5."},
		{fmt.Sprintf("A vertical strip at \\(x=1\\) revolves about the line \\(x=%d\\) (shells). What is the shell radius? (enter a number)", c), fmt.Sprintf("%d", c-1), fmt.Sprintf("Radius is the distance %d-1 = %d.", c, c-1)},
	}
	e := tableHard[rand.Intn(len(tableHard))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
