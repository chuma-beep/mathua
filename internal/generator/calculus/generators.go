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

type limitInfinityGen struct{}

func (g *limitInfinityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"\\(\\lim_{x \\to \\infty} 1/x\\)", "0", "As \\(x\\to\\infty\\), \\(1/x \\to 0\\)."},
		{"\\(\\lim_{x \\to \\infty} 1/x^{2}\\)", "0", "As \\(x\\to\\infty\\), \\(1/x^{2} \\to 0\\)."},
		{"\\(\\lim_{x \\to \\infty} x/(x+1)\\)", "1", "Divide numerator and denominator by \\(x\\): \\(1/(1+1/x) \\to 1\\)."},
		{"\\(\\lim_{x \\to \\infty} 2x/(x+1)\\)", "2", "Divide numerator and denominator by \\(x\\): \\(2/(1+1/x) \\to 2\\)."},
		{"\\(\\lim_{x \\to \\infty} (x^{2}+1)/x^{2}\\)", "1", "\\((x^{2}+1)/x^{2} = 1+1/x^{2} \\to 1\\)."},
		{"\\(\\lim_{x \\to \\infty} 3\\)", "3", "The limit of a constant is the constant itself."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is %s?", e.q),
		Answer:      e.a,
		Explanation: e.e,
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
	table := []entry{
		{"What is the slope of \\(f(x)=x^{2}\\) at \\(x=3\\)?", 6, "\\(f'(x)=2x\\), so \\(f'(3)=2 \\cdot 3 = 6\\)."},
		{"What is the slope of \\(f(x)=x^{2}\\) at \\(x=2\\)?", 4, "\\(f'(x)=2x\\), so \\(f'(2)=2 \\cdot 2 = 4\\)."},
		{"What is the slope of \\(f(x)=x^{2}\\) at \\(x=1\\)?", 2, "\\(f'(x)=2x\\), so \\(f'(1)=2 \\cdot 1 = 2\\)."},
		{"What is the slope of \\(f(x)=2x^{2}\\) at \\(x=2\\)?", 8, "\\(f'(x)=4x\\), so \\(f'(2)=4 \\cdot 2 = 8\\)."},
		{"What is the slope of \\(f(x)=x^{3}\\) at \\(x=2\\)?", 12, "\\(f'(x)=3x^{2}\\), so \\(f'(2)=3 \\cdot 4 = 12\\)."},
		{"What is the slope of \\(f(x)=x^{3}\\) at \\(x=1\\)?", 3, "\\(f'(x)=3x^{2}\\), so \\(f'(1)=3 \\cdot 1 = 3\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      fmt.Sprintf("%d", e.a),
		Explanation: e.e,
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
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Find \\(f'(2)\\) if \\(f(x)=\\frac{x+1}{x-1}\\).", "-2", "Quotient rule: f'(x)=((x-1)·1-(x+1)·1)/(x-1)² = -2/(x-1)². f'(2) = -2/(1)² = -2."},
		{"Find \\(f'(3)\\) if \\(f(x)=\\frac{x+2}{x-2}\\).", "-4", "f'(x)=((x-2)-(x+2))/(x-2)² = -4/(x-2)². f'(3) = -4/1 = -4."},
		{"Find \\(f'(0)\\) if \\(f(x)=\\frac{x}{x+1}\\).", "1", "f'(x)=((x+1)-x)/(x+1)² = 1/(x+1)². f'(0) = 1/1 = 1."},
		{"Find \\(f'(0)\\) if \\(f(x)=\\frac{x-1}{x+1}\\).", "2", "f'(x)=((x+1)-(x-1))/(x+1)² = 2/(x+1)². f'(0) = 2/1 = 2."},
		{"For \\(f(x)=\\frac{x+1}{x-1}\\), what is \\(f'(x)\\) in terms of \\(x\\)?", "-2/(x-1)^2", "Quotient rule gives f'(x) = -2/(x-1)²."},
	}
	tableHard := []entry{
		{"Find \\(f'(2)\\) if \\(f(x)=\\frac{2x+1}{x-1}\\).", "-3", "f'(x)=((x-1)·2-(2x+1))/(x-1)² = -3/(x-1)². f'(2) = -3."},
		{"Find \\(f'(1)\\) if \\(f(x)=\\frac{x+3}{x+1}\\).", "-1/2", "f'(x)=((x+1)-(x+3))/(x+1)² = -2/(x+1)². f'(1) = -2/4 = -1/2."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=\\sin(x)\\).", "cos(x)", "\\(\\frac{d}{dx} \\sin(x) = \\cos(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\cos(x)\\).", "-sin(x)", "\\(\\frac{d}{dx} \\cos(x) = -\\sin(x)\\)."},
		{"What is \\(\\frac{d}{dx} \\sin(x)\\)?", "cos(x)", "The derivative of sine is cosine."},
		{"What is \\(\\frac{d}{dx} \\cos(x)\\)?", "-sin(x)", "The derivative of cosine is \\(-\\sin(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\sin(2x)\\).", "2cos(2x)", "Chain rule: \\(\\frac{d}{dx} \\sin(2x) = \\cos(2x) \\cdot 2 = 2\\cos(2x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\cos(3x)\\).", "-3sin(3x)", "Chain rule: \\(\\frac{d}{dx} \\cos(3x) = -\\sin(3x) \\cdot 3 = -3\\sin(3x)\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
	}
}

type derivExpLogGen struct{}

func (g *derivExpLogGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=e^{x}\\).", "e^x", "\\(\\frac{d}{dx} e^{x} = e^{x}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\ln(x)\\).", "1/x", "\\(\\frac{d}{dx} \\ln(x) = 1/x\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=e^{3x}\\).", "3e^(3x)", "Chain rule: \\(\\frac{d}{dx} e^{3x} = e^{3x} \\cdot 3 = 3e^{3x}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=e^{2x}\\).", "2e^(2x)", "Chain rule: \\(\\frac{d}{dx} e^{2x} = e^{2x} \\cdot 2 = 2e^{2x}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\ln(2x)\\).", "1/x", "Chain rule: \\(\\frac{d}{dx} \\ln(2x) = \\frac{1}{2x} \\cdot 2 = 1/x\\)."},
		{"What is \\(\\frac{d}{dx} e^{x}\\)?", "e^x", "The derivative of \\(e^{x}\\) is \\(e^{x}\\)."},
		{"What is \\(\\frac{d}{dx} \\ln(x)\\)?", "1/x", "The derivative of \\(\\ln(x)\\) is \\(1/x\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    e.q,
		Answer:      e.a,
		Explanation: e.e,
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
	n := rand.Intn(max(2, scale)) + 1
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
		Question:    fmt.Sprintf("Find \\(\\int %d x^{%d} \\, dx\\) (use power rule for integration, omit \\(+C\\)).", a, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("∫%dx^%d dx = %d/(%d+1)·x^(%d+1)+C = %s+C.", a, n, a, n, n, answer),
	}
}

type integralSubstitutionGen struct{}

func (g *integralSubstitutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(2, scale)) + 2
	coef := n
	exp := n - 1
	answer := fmt.Sprintf("e^(x^%d)+C", n)
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\int %d x^{%d} e^{x^{%d}} \\, dx\\).", coef, exp, n),
		Answer:      answer,
		Explanation: fmt.Sprintf("Let u=x^%d, du=%dx^%d dx. Integral becomes ∫e^u du = e^u+C = e^(x^%d)+C.", n, n, exp, n),
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
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(2, scale)) + 1
	answer := fmt.Sprintf("x^%d", n)
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(F(x)=\\int_{0}^{x} t^{%d} \\, dt\\), what is \\(F'(x)\\)?", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("By the Fundamental Theorem of Calculus, F'(x)=x^%d.", n),
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for the area between curves: \\(A = \\int_{a}^{b} (f(x) - g(x))\\,dx\\).\n"+
				"Subgoal 2: Substitute parameters \\(f(x) = %d x\\), \\(g(x) = 0\\), \\(a = 0\\) and \\(b = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: find roots \\(\\pm %d\\) and integrate \\(A = \\int_{-%d}^{%d} (%d - x^2)\\,dx\\).\n"+
				"Subgoal 2: Use symmetry: \\(A = 2 \\int_{0}^{%d} (%d - x^2)\\,dx\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: intersections occur at \\(x = 0\\) and \\(x = %d\\), so \\(A = \\int_{0}^{%d} (%d x - x^2)\\,dx\\).\n"+
				"Subgoal 2: Find antiderivative: \\([%d x^2 / 2 - x^3 / 3]_{0}^{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(A = \\int_{0}^{%d} %d x^2\\,dx\\).\n"+
				"Subgoal 2: Find the antiderivative: \\([%d x^3]_{0}^{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(A = \\int_{a}^{b} (f(x) - g(x))\\,dx\\).\n"+
				"Subgoal 2: Substitute \\(f(x) - g(x) = (%d x + %d) - %d x = %d\\).\n"+
				"Subgoal 3: Evaluate: \\(\\int_{0}^{%d} %d\\,dx = %d \\cdot %d = %d\\).", m, b, m, b, L, b, b, L, ans),
		}
	}
}

type integralVolumeGen struct{}

func (g *integralVolumeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q, a, e string
	}
	tableEasy := []entry{
		{"Find the volume when \\(y=\\sqrt{x}\\) from \\(x=0\\) to 4 is revolved around the \\(x\\)-axis.", "8pi", "V = π∫(√x)²dx = π∫x dx = π[x²/2]₀⁴ = π(16/2-0) = 8π."},
		{"Find the volume when \\(y=x\\) from \\(x=0\\) to 1 is revolved around the \\(x\\)-axis.", "pi/3", "V = π∫x²dx = π[x³/3]₀¹ = pi/3."},
		{"Find the volume when \\(y=x\\) from \\(x=0\\) to 2 is revolved around the \\(x\\)-axis.", "8pi/3", "V = π∫x²dx = π[x³/3]₀² = 8pi/3."},
		{"Find the volume when \\(y=2\\) from \\(x=0\\) to 3 is revolved around the \\(x\\)-axis.", "12pi", "V = π∫4dx = π[4x]₀³ = 12pi (a cylinder of radius 2)."},
	}
	tableHard := []entry{
		{"Find the volume when \\(y=x^{2}\\) from \\(x=0\\) to 2 is revolved around the \\(x\\)-axis.", "32pi/5", "V = π∫x⁴dx = π[x⁵/5]₀² = 32pi/5."},
		{"Find the volume when \\(y=\\sqrt{x}\\) from \\(x=0\\) to 1 is revolved around the \\(x\\)-axis.", "pi/2", "V = π∫x dx = π[x²/2]₀¹ = pi/2."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: if \\(f\\) is continuous on \\([a, b]\\), differentiable on \\((a, b)\\) and \\(f(a) = f(b)\\), there exists \\(c \\in (a, b)\\) with \\(f'(c) = 0\\).\n"+
				"Subgoal 2: Compute the derivative \\(f'(x) = 2x - %d\\) and set it to zero.\n"+
				"Subgoal 3: Solve for \\(c\\): \\(2c - %d = 0 \\implies c = %d\\).", b, b, cMid),
		}
	case 1:
		// Inverted quadratic f(x) = -x^2 + 2kx on [0, 2k]
		k := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		b := 2 * k
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = -x^2 + %d x\\) on \\([0, %d]\\), find the critical point \\(c\\) guaranteed by Rolle's Theorem.", b, b),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Subgoal 1: Confirm the boundary condition \\(f(0) = 0\\) and \\(f(%d) = -(%d)^2 + %d(%d) = 0\\).\n"+
				"Subgoal 2: Compute \\(f'(x) = -2x + %d\\) and equate to zero.\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: since \\(f(%d) = f(%d) = 0\\), a point \\(c \\in (%d, %d)\\) satisfies \\(f'(c) = 0\\).\n"+
				"Subgoal 2: Differentiate: \\(f'(x) = -2x + %d\\).\n"+
				"Subgoal 3: Set \\(-2c + %d = 0\\) and solve for \\(c = %d\\).", aEnd, bEnd, aEnd, bEnd, lin, lin, c0),
		}
	case 3:
		// Value of derivative guaranteed by Rolle
		x1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		x2 := x1 + mathutil.ScaleInt(ctx.Difficulty, 2, 6)
		return generator.Problem{
			Question: fmt.Sprintf("If a differentiable function \\(f(x)\\) satisfies \\(f(%d) = f(%d)\\) on \\([%d, %d]\\), what is the value of \\(f'(c)\\) at the point \\(c\\) guaranteed by Rolle's Theorem?", x1, x2, x1, x2),
			Answer:   "0",
			Explanation: "Subgoal 1: State the conclusion of Rolle's Theorem: the derivative must vanish at some interior point \\(c \\in (a, b)\\).\n" +
				"Subgoal 2: Identify the critical condition: \\(f'(c) = 0\\).\n" +
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
			Explanation: fmt.Sprintf("Subgoal 1: State Rolle's Theorem: since \\(f(%d) = f(%d)\\), there exists \\(c \\in (%d, %d)\\) with \\(f'(c) = 0\\).\n"+
				"Subgoal 2: Differentiate \\(f(x)\\): \\(f'(x) = 4x - %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem formula: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\n"+
				"Subgoal 2: Compute the average rate of change: \\(\\frac{(%d)^2 - (%d)^2}{%d - %d} = \\frac{%d - %d}{%d} = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\n"+
				"Subgoal 2: For any quadratic function, the instantaneous rate of change equals the average rate of change at the midpoint \\(c = \\frac{a + b}{2}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the average rate of change formula: \\(\\frac{f(b) - f(a)}{b - a}\\).\n"+
				"Subgoal 2: Evaluate \\(f(%d) = (%d)^2 + %d(%d) = %d\\) and \\(f(0) = 0\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem in kinematics: \\(v_{\\text{avg}} = \\frac{s(t_2) - s(t_1)}{t_2 - t_1} = s'(c)\\).\n"+
				"Subgoal 2: Substitute distance \\(\\Delta s = %d\\) and time \\(\\Delta t = %d\\): \\(\\frac{%d}{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Mean Value Theorem: \\(f'(c) = \\frac{f(b) - f(a)}{b - a}\\).\n"+
				"Subgoal 2: Evaluate \\(f(%d) = (%d)^3 = %d\\) and \\(f(0) = 0\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State Cauchy's Mean Value Theorem: \\(\\frac{f(b) - f(a)}{g(b) - g(a)} = \\frac{f'(c)}{g'(c)}\\).\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 2x\\) and \\(g'(x) = 1\\), giving quotient \\(\\frac{2c}{1} = \\frac{(%d)^2 - (%d)^2}{%d - %d} = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State Cauchy's Mean Value Theorem: \\(\\frac{f(b) - f(a)}{g(b) - g(a)} = \\frac{f'(c)}{g'(c)}\\).\n"+
				"Subgoal 2: Substitute derivatives \\(f'(c) = 3c^2\\) and \\(g'(c) = 2c\\), so \\(\\frac{3c^2}{2c} = \\frac{3}{2}c\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Cauchy difference quotient: \\(\\frac{f(b) - f(a)}{g(b) - g(a)}\\).\n"+
				"Subgoal 2: Evaluate \\(f(%d) = %d(%d)^2 = %d\\) and \\(g(%d) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Differentiate the linear functions: \\(f'(x) = %d\\) and \\(g'(x) = %d\\).\n"+
				"Subgoal 2: Form the derivative ratio: \\(\\frac{f'(c)}{g'(c)} = \\frac{%d}{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Substitute \\(g(x) = x\\) into the denominator \\(g(b) - g(a)\\).\n"+
				"Subgoal 2: Compute \\(b - a = %d - %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State Fermat's Theorem: differentiable local extrema satisfy \\(f'(c) = 0\\).\n"+
				"Subgoal 2: Differentiate: \\(f'(x) = 2(%d)x - %d = %d x - %d\\).\n"+
				"Subgoal 3: Set \\(%d c - %d = 0\\) and solve for \\(c = %d\\).", a, b, 2*a, b, 2*a, b, h),
		}
	case 1:
		// Value of derivative at local extremum
		x0 := mathutil.ScaleInt(ctx.Difficulty, 1, 12)
		return generator.Problem{
			Question: fmt.Sprintf("If a differentiable function \\(f(x)\\) attains a local extremum at \\(x = %d\\), what is the value of \\(f'(%d)\\) by Fermat's Theorem?", x0, x0),
			Answer:   "0",
			Explanation: "Subgoal 1: State Fermat's Theorem on stationary points: if \\(f\\) has a local extremum at \\(c\\) and \\(f'(c)\\) exists, then \\(f'(c) = 0\\).\n" +
				"Subgoal 2: Identify that \\(f\\) is differentiable at the local extremum.\n" +
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
			Explanation: fmt.Sprintf("Subgoal 1: Differentiate \\(f(x)\\): \\(f'(x) = 6x^2 - %d x + %d = 6(x - %d)(x - %d)\\).\n"+
				"Subgoal 2: Set \\(f'(x) = 0\\) to find stationary points \\(x = %d\\) and \\(x = %d\\).\n"+
				"Subgoal 3: Select the smaller critical point: %d.", 2*b, cVal, h1, h2, h1, h2, h1),
		}
	case 3:
		// Maximum value at stationary point for f(x) = -(x - h)^2 + K
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 7)
		k := mathutil.ScaleInt(ctx.Difficulty, 5, 25)
		return generator.Problem{
			Question: fmt.Sprintf("The function \\(f(x) = -(x - %d)^2 + %d\\) attains its local maximum at its stationary point. What is this maximum value?", h, k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf("Subgoal 1: Apply Fermat's Theorem: the stationary point occurs where \\(f'(x) = -2(x - %d) = 0 \\implies x = %d\\).\n"+
				"Subgoal 2: Substitute \\(x = %d\\) into \\(f(x)\\): \\(f(%d) = -(%d - %d)^2 + %d\\).\n"+
				"Subgoal 3: Evaluate: \\(0 + %d = %d\\).", h, h, h, h, h, h, k, k, k),
		}
	default:
		// Count of stationary points for non-degenerate quadratic
		return generator.Problem{
			Question: "How many stationary points (where \\(f'(c) = 0\\)) does a non-degenerate quadratic polynomial \\(f(x) = ax^2 + bx + c\\) possess?",
			Answer:   "1",
			Explanation: "Subgoal 1: Compute the derivative of a quadratic: \\(f'(x) = 2ax + b\\).\n" +
				"Subgoal 2: Set \\(2ax + b = 0\\).\n" +
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
			Explanation: fmt.Sprintf("Subgoal 1: State the difference quotient formula: \\(\\frac{f(x_0 + h) - f(x_0)}{h}\\).\n"+
				"Subgoal 2: Evaluate \\(f(%d) = %d\\) and \\(f(%d) = %d\\), so \\(f(%d) - f(%d) = %d - %d = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the difference quotient for a linear function: \\(\\frac{f(x + h) - f(x)}{h}\\).\n"+
				"Subgoal 2: Substitute \\(f(x + h) - f(x) = [%d(x + h) + %d] - [%d x + %d] = %d h\\).\n"+
				"Subgoal 3: Divide by \\(h\\): \\(\\frac{%d h}{h} = %d\\).", m, b, m, b, m, m, m),
		}
	case 2:
		// f(x) = x^3 at x0 = 1
		h := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		ans := 3 + 3*h + h*h
		return generator.Problem{
			Question: fmt.Sprintf("For \\(f(x) = x^3\\), compute the difference quotient \\(\\frac{f(1 + %d) - f(1)}{%d}\\).", h, h),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: Expand \\(f(1 + h) = (1 + h)^3 = 1 + 3h + 3h^2 + h^3\\).\n"+
				"Subgoal 2: Subtract \\(f(1) = 1\\) to get \\(3h + 3h^2 + h^3\\), and divide by \\(h\\) to obtain \\(3 + 3h + h^2\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Recognize that the limit of the difference quotient defines the derivative: \\(f'(x_0) = \\lim_{h \\to 0} \\frac{f(x_0 + h) - f(x_0)}{h}\\).\n"+
				"Subgoal 2: Differentiate \\(f(x) = %d x^2\\): \\(f'(x) = 2(%d)x = %d x\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Substitute \\(x = %d\\) into \\(f(x)\\): \\(f(%d) = %d(%d)^2 + %d(%d) = %d\\).\n"+
				"Subgoal 2: Subtract \\(f(0) = 0\\), giving numerator %d.\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the differential definition: \\(dy = f'(x) dx\\).\n"+
				"Subgoal 2: Compute derivative \\(f'(x) = 2(%d)x = %d x\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the differential formula: \\(dy = f'(x) dx\\).\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 3(%d)x^2 = %d x^2\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d(%d)^2 = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the linear approximation differential: \\(dy = f'(x) dx\\).\n"+
				"Subgoal 2: Compute \\(f'(x) = 2(%d)x + %d = %d x + %d\\) and evaluate at \\(x = %d\\): \\(f'(%d) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Recall that the ratio of differentials equals the derivative: \\(\\frac{dy}{dx} = f'(x)\\).\n"+
				"Subgoal 2: Differentiate \\(f'(x) = 2(%d)x + %d = %d x + %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the differential rule: \\(dy = f'(x) dx\\).\n"+
				"Subgoal 2: Since \\(f(x) = %d x + 7\\), its derivative is constant: \\(f'(x) = %d\\).\n"+
				"Subgoal 3: Multiply by \\(dx = %d\\): \\(dy = %d \\times %d = %d\\).", m, m, dx, m, dx, dy),
		}
	}
}

type convexityGen struct{}

func (g *convexityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If f''(x) > 0 for all x, the function is...", "convex (concave up)", "A positive second derivative means the function is convex (curving upward)."},
		{"If f''(x) < 0 for all x, the function is...", "concave (concave down)", "A negative second derivative means the function is concave (curving downward)."},
		{"Is f(x)=x² convex or concave?", "convex", "f''(x)=2 > 0, so f(x)=x² is convex everywhere."},
		{"Is f(x)=-x² convex or concave?", "concave", "f''(x)=-2 < 0, so f(x)=-x² is concave everywhere."},
		{"Is f(x)=x³ convex or concave at x=0?", "neither (inflection point)", "f''(x)=6x, so f''(0)=0 and the concavity changes sign at 0. x=0 is an inflection point."},
		{"What is an inflection point?", "a point where concavity changes", "An inflection point is where the function changes from convex to concave or vice versa."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type nonDiffGen struct{}

func (g *nonDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is f(x)=|x| differentiable at x=0?", "no", "|x| has a corner at x=0: left derivative = -1, right derivative = 1."},
		{"Is f(x)=x^(1/3) differentiable at x=0?", "no", "f'(x)=1/(3x^(2/3)) → ∞ as x→0, so the derivative does not exist (vertical tangent)."},
		{"Is f(x)=√|x| differentiable at x=0?", "no", "√|x| has a cusp at x=0: the derivative approaches ±∞ from each side."},
		{"Name one way a function can fail to be differentiable.", "corner/cusp/vertical tangent/discontinuity", "Functions can fail to be differentiable at corners (|x|), cusps (√|x|), vertical tangents (x^(1/3)), or points of discontinuity."},
		{"If a function is differentiable at x=c, is it necessarily continuous at x=c?", "yes", "Differentiability implies continuity. But the converse is not true, continuity does not imply differentiability."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What method approximates the integral using rectangles? (enter trapezoids, parabolas, or rectangles)", "rectangles", "Riemann sums partition [a,b] into subintervals and use rectangles to approximate the area."},
		{"What method approximates the integral using trapezoids? (enter trapezoids, parabolas, or rectangles)", "trapezoids", "The trapezoidal rule approximates the integral with (Δx/2)(f(x₀)+2f(x₁)+...+2f(xₙ₋₁)+f(xₙ))."},
		{"What method approximates the integral using parabolas? (enter trapezoids, parabolas, or rectangles)", "parabolas", "Simpson's rule uses quadratic polynomials and is more accurate than trapezoidal for smooth functions."},
		{"What power of h gives the accuracy order of the trapezoidal rule? (enter a number)", "2", "The trapezoidal rule error is proportional to h^2 times f'' at some point, making it second-order accurate: O(h^2)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type riemannCriteriaGen struct{}

func (g *riemannCriteriaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is f(x)=x² Riemann integrable on [0,1]?", "yes", "x² is continuous on [0,1], and continuous functions on closed intervals are Riemann integrable."},
		{"Is f(x)={0 for x rational, 1 for x irrational} Riemann integrable on [0,1]?", "no", "The upper sum = 1 and lower sum = 0 for any partition, so the integral does not exist (Dirichlet function)."},
		{"Is f(x)=1/x Riemann integrable on [0,1]?", "no", "1/x is unbounded near 0, so it is not Riemann integrable on [0,1] (though the improper integral exists)."},
		{"Is a function with finitely many discontinuities Riemann integrable?", "yes", "A bounded function with finitely many discontinuities on [a,b] is Riemann integrable."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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

type trigSubstitutionGen struct{}

func (g *trigSubstitutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What substitution is used for \\(\\sqrt{a^{2}-x^{2}}\\)? (enter as x = ...)", "x = a*sin(theta)", "For \\(\\sqrt{a^{2}-x^{2}}\\), let \\(x=a\\sin(\\theta)\\), then \\(dx=a\\cos(\\theta)d\\theta\\) and \\(\\sqrt{a^{2}-x^{2}} = a\\cos(\\theta)\\)."},
		{"What substitution is used for \\(\\sqrt{a^{2}+x^{2}}\\)? (enter as x = ...)", "x = a*tan(theta)", "For \\(\\sqrt{a^{2}+x^{2}}\\), let \\(x=a\\tan(\\theta)\\), then \\(dx=a\\sec^{2}(\\theta)d\\theta\\) and \\(\\sqrt{a^{2}+x^{2}} = a\\sec(\\theta)\\)."},
		{"What substitution is used for \\(\\sqrt{x^{2}-a^{2}}\\)? (enter as x = ...)", "x = a*sec(theta)", "For \\(\\sqrt{x^{2}-a^{2}}\\), let \\(x=a\\sec(\\theta)\\), then \\(dx=a\\sec(\\theta)\\tan(\\theta)d\\theta\\) and \\(\\sqrt{x^{2}-a^{2}} = a\\tan(\\theta)\\)."},
		{"Find \\(\\int \\frac{1}{\\sqrt{1-x^{2}}} \\, dx\\).", "arcsin(x)+C", "Let \\(x=\\sin(\\theta), dx=\\cos(\\theta)d\\theta\\). \\(\\int \\frac{\\cos(\\theta)}{\\cos(\\theta)} d\\theta = \\int d\\theta = \\theta+C = \\arcsin(x)+C\\)."},
		{"Find \\(\\int \\frac{1}{1+x^{2}} \\, dx\\).", "arctan(x)+C", "Let \\(x=\\tan(\\theta), dx=\\sec^{2}(\\theta)d\\theta\\). \\(\\int \\frac{\\sec^{2}(\\theta)}{1+\\tan^{2}(\\theta)} d\\theta = \\int d\\theta = \\theta+C = \\arctan(x)+C\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type weierstrassSubGen struct{}

func (g *weierstrassSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the Weierstrass substitution? (enter as t = ...)", "t = tan(x/2)", "Let \\(t=\\tan(x/2)\\). Then \\(\\sin(x)=2t/(1+t^{2})\\), \\(\\cos(x)=(1-t^{2})/(1+t^{2})\\), \\(dx=2/(1+t^{2})dt\\)."},
		{"What is \\(\\sin(x)\\) in terms of \\(t=\\tan(x/2)\\)?", "2t/(1+t^2)", "\\(\\sin(x) = 2t/(1+t^{2})\\) where \\(t=\\tan(x/2)\\)."},
		{"What is \\(\\cos(x)\\) in terms of \\(t=\\tan(x/2)\\)?", "(1-t^2)/(1+t^2)", "\\(\\cos(x) = (1-t^{2})/(1+t^{2})\\) where \\(t=\\tan(x/2)\\)."},
		{"What is \\(dx\\) in terms of \\(t=\\tan(x/2)\\)?", "2/(1+t^2) dt", "\\(dx = 2/(1+t^{2}) dt\\) where \\(t=\\tan(x/2)\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// --- Limit generators ---

type limitAlgebraGen struct{}

func (g *limitAlgebraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is lim (f(x)+g(x)) in terms of individual limits?", "lim f(x) + lim g(x)", "The limit of a sum is the sum of the limits (provided both exist)."},
		{"Is the limit of a product the product of the limits? (yes/no)", "yes", "The limit of a product is lim f(x) times lim g(x) (provided both exist)."},
		{"Is the limit of a quotient the quotient of the limits when the denominator limit is nonzero? (yes/no)", "yes", "The limit of a quotient is lim f(x) / lim g(x), provided the denominator limit is nonzero."},
		{"Can a constant multiple be factored out of a limit? (yes/no)", "yes", "Constant multiples can be factored out: lim c f(x) = c times lim f(x)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type asymptotesGen struct{}

func (g *asymptotesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What are the vertical asymptotes of \\(f(x)=1/(x-2)\\)?", "x=2", "The denominator is 0 at \\(x=2\\), and the numerator is nonzero, so \\(x=2\\) is a vertical asymptote."},
		{"What is the horizontal asymptote of \\(f(x)=1/x\\)?", "y=0", "As \\(x\\to\\infty\\), \\(1/x\\to 0\\). As \\(x\\to -\\infty\\), \\(1/x\\to 0\\). So \\(y=0\\) is the horizontal asymptote."},
		{"What is the horizontal asymptote of \\(f(x)=\\frac{2x+1}{x-3}\\)?", "y=2", "As \\(x\\to\\infty\\), \\(\\frac{2x+1}{x-3}\\to 2\\). The ratio of leading coefficients gives the horizontal asymptote."},
		{"Can a function cross its horizontal asymptote?", "yes", "A function CAN cross its horizontal asymptote (unlike vertical asymptotes). The asymptote describes end behavior only."},
		{"What are the vertical asymptotes of \\(f(x)=\\tan(x)\\)? (enter as x = ...)", "x=pi/2 + n*pi", "\\(\\tan(x)=\\sin(x)/\\cos(x)\\). Vertical asymptotes occur where \\(\\cos(x)=0\\), i.e., at \\(x=\\pi/2+n\\pi\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type discontinuityGen struct{}

func (g *discontinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What kind of discontinuity does \\(1/x\\) have at \\(x=0\\)?", "infinite", "\\(1/x\\) has an infinite discontinuity (vertical asymptote) at \\(x=0\\)."},
		{"What kind of discontinuity does \\(|x|/x\\) have at \\(x=0\\)?", "jump", "The left limit is -1 and the right limit is 1, so there is a jump discontinuity."},
		{"What kind of discontinuity does \\(\\frac{x^{2}-1}{x-1}\\) have at \\(x=1\\)?", "removable", "The function simplifies to \\(x+1\\) for \\(x\\neq 1\\), and the limit as \\(x\\to 1\\) is 2."},
		{"What kind of discontinuity does \\(\\sin(1/x)\\) have at \\(x=0\\)?", "essential", "\\(\\sin(1/x)\\) oscillates infinitely near 0, so the limit does not exist (essential/oscillatory discontinuity)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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

type lhopitalGen struct{}

func (g *lhopitalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\lim_{x \\to 0} \\frac{\\sin(x)}{x}\\) using L'Hôpital's rule.", "1", "Both numerator and denominator \\(\\to 0\\). L'Hôpital: \\(\\lim_{x \\to 0} \\frac{\\cos(x)}{1} = 1\\)."},
		{"Find \\(\\lim_{x \\to 0} \\frac{e^{x}-1}{x}\\) using L'Hôpital's rule.", "1", "Both \\(\\to 0\\). L'Hôpital: \\(\\lim_{x \\to 0} \\frac{e^{x}}{1} = 1\\)."},
		{"Find \\(\\lim_{x \\to \\infty} \\frac{x}{e^{x}}\\) using L'Hôpital's rule.", "0", "Both \\(\\to \\infty\\). L'Hôpital: \\(\\lim_{x \\to \\infty} \\frac{1}{e^{x}} = 0\\)."},
		{"What condition is required for L'Hôpital's rule? (yes/no: a 0/0 or infinity/infinity form)", "yes", "L'Hôpital's rule applies to 0/0 or infinity/infinity indeterminate forms: both numerator and denominator tend to 0 or to plus/minus infinity."},
		{"Is x^2*sin(1/x)/sin(x) as x→0 a 0/0 form? (yes/no)", "yes", "Both numerator and denominator tend to 0, so it is 0/0; L'Hôpital applies in principle but the derivatives oscillate, so the squeeze theorem is the practical tool."},
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does f(x)=o(g(x)) as x→∞ mean? (yes/no: f(x)/g(x) tends to 0)", "yes", "f=o(g) means that f(x)/g(x) → 0 as x → ∞ (f grows strictly slower than g)."},
		{"Is x²=o(x³) as x→∞?", "yes", "x²/x³ = 1/x → 0, so x²=o(x³)."},
		{"Is x³=o(x³) as x→∞?", "no", "x³/x³ = 1 → 1, not 0, so x³≠o(x³). However, x³=O(x³)."},
		{"What is the difference between O and o? (yes/no: o needs ratio → 0, O allows any bounded ratio)", "yes", "f=O(g) means |f|≤M·|g| for large x. f=o(g) means f/g → 0, which is a stricter requirement: o requires ratio → 0; O allows any bounded ratio."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type squeezeGen struct{}

func (g *squeezeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does the squeeze theorem say? (yes/no: squeezed f takes the common limit L)", "yes", "The squeeze theorem: if g(x)≤f(x)≤h(x) and lim g=lim h=L, then lim f=L."},
		{"Find \\(\\lim_{x \\to 0} x^{2} \\sin(1/x)\\) using the squeeze theorem.", "0", "\\(-1\\leq \\sin(1/x)\\leq 1 \\implies -x^{2} \\leq x^{2}\\sin(1/x) \\leq x^{2}\\). Since \\(-x^{2}\\to 0\\) and \\(x^{2}\\to 0\\), the squeeze theorem gives limit 0."},
		{"Find \\(\\lim_{x \\to 0} x \\cos(1/x)\\) using the squeeze theorem.", "0", "\\(-1\\leq \\cos(1/x)\\leq 1 \\implies -x \\leq x\\cos(1/x) \\leq x\\). Since \\(-x\\to 0\\) and \\(x\\to 0\\), the limit is 0."},
		{"What is a common way to apply the squeeze theorem? (yes/no: bound f between simpler functions with the same limit)", "yes", "Bound the function between two simpler functions with the same limit, then show the upper and lower bounds converge to that limit."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type supremumGen struct{}

func (g *supremumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the supremum of the set \\(\\{1-1/n : n \\in \\mathbb{N}\\}\\)?", "1", "The values approach 1 from below. 1 is an upper bound, and no smaller number is an upper bound. Supremum = 1."},
		{"What is the infimum of the set \\(\\{1/n : n \\in \\mathbb{N}\\}\\)?", "0", "The values approach 0 from above. 0 is a lower bound, and no larger number is a lower bound. Infimum = 0."},
		{"Is the supremum always in the set?", "no", "The supremum need not be in the set. E.g., \\(\\sup\\{1-1/n\\}=1\\), but 1 is not in the set."},
		{"Is the maximum always in the set?", "yes", "If a maximum exists, it is the supremum and is also in the set. Unlike supremum, maximum must be attained."},
		{"What is \\(\\sup\\{x \\in \\mathbb{R} : x^{2}<2\\}\\)? (enter as sqrt(...))", "sqrt(2)", "The set is \\((-\\sqrt{2}, \\sqrt{2})\\). The supremum is \\(\\sqrt{2}\\) (the least upper bound), which is not in the set."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does the extreme value theorem (Weierstrass) state? (yes/no: continuous f on [a,b] attains max and min)", "yes", "If f is continuous on [a,b], then there exist c,d∈[a,b] such that f(c)≤f(x)≤f(d) for all x∈[a,b]: a continuous function on [a,b] attains its max and min."},
		{"Does the extreme value theorem apply to f(x)=1/x on (0,1]?", "no", "The interval must be closed. (0,1] is not closed, and in fact 1/x → ∞ as x→0+."},
		{"Does the extreme value theorem apply to f(x)=x on [0,1]?", "yes", "f(x)=x is continuous on the closed bounded interval [0,1], so it attains max=1 at x=1 and min=0 at x=0."},
		{"Does the extreme value theorem apply to f(x)={x for x<1, 2 for x=1} on [0,1]?", "no", "f is not continuous on [0,1] (jump at x=1), so the EVT does not apply."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// --- Sequence generators ---

type seqConceptGen struct{}

func (g *seqConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is a sequence?", "an ordered list of numbers (a function ℕ→ℝ)", "A sequence is a function from natural numbers (or a subset) to real numbers, denoted {aₙ}."},
		{"What is aₙ = n/(n+1): write first 3 terms.", "1/2, 2/3, 3/4", "a₁=1/2, a₂=2/3, a₃=3/4."},
		{"What is aₙ = (-1)ⁿ: write first 4 terms.", "-1, 1, -1, 1", "a₁=-1, a₂=1, a₃=-1, a₄=1. This sequence alternates."},
		{"What is a recursive sequence?", "a sequence defined by relating aₙ to previous terms", "E.g., a₁=1, aₙ=aₙ₋₁+2 defines the odd numbers: 1, 3, 5, 7, ..."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type seqConvergenceGen struct{}

func (g *seqConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The answers used to be prose ("no (diverges to \(\infty\))") on a concept
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does Cauchy mean terms get arbitrarily close as n increases? (yes/no)", "yes", "A sequence is Cauchy where for any ε>0, there exists N such that |aₘ-aₙ|<ε for all m,n>N; terms get arbitrarily close as n increases."},
		{"In ℝ, do all Cauchy sequences converge?", "yes", "ℝ is complete, every Cauchy sequence in ℝ converges to a real number."},
		{"In ℚ, do all Cauchy sequences converge?", "no", "ℚ is not complete. E.g., a sequence of rationals converging to √2 is Cauchy in ℚ but does not converge in ℚ."},
		{"Is every convergent sequence Cauchy?", "yes", "Every convergent sequence is Cauchy. In ℝ, the converse also holds: Cauchy ⟺ convergent."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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

type eulerSeqGen struct{}

func (g *eulerSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Which expression defines Euler's number e? (enter (1+1/n)^n)", "(1+1/n)^n", "e = lim (1+1/n)^n ≈ 2.71828..."},
		{"What is \\(e\\) to 3 decimal places?", "2.718", "\\(e = \\lim (1+1/n)^{n} = 2.718281828...\\)"},
		{"What is the limit of (1+1/n)^n? (enter e)", "e", "e = 1 + 1/1! + 1/2! + 1/3! + ... = sum 1/n!; the sequence (1+1/n)^n tends to e."},
		{"Is e irrational? (yes/no)", "yes", "e is irrational (proved by Euler). It is also transcendental."},
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

func (g *seriesConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is a series the sum of terms of a sequence? (yes/no)", "yes", "A series ∑ aₙ = a₁ + a₂ + a₃ + ... represents the sum of a sequence's terms; the sum of terms of a sequence."},
		{"Is the nth partial sum the sum of the first n terms? (yes/no)", "yes", "The nth partial sum Sₙ = ∑_{k=1}ⁿ aₖ is the sum of the first n terms of the series."},
		{"Does a series converge when its partial sums converge? (yes/no)", "yes", "A series converges if lim Sₙ exists and is finite; when the sequence of partial sums converges."},
		{"Does every series converge? (yes/no)", "no", "No. E.g., the harmonic series diverges. To find the sum of a series means to find the limit of its partial sums; the sum is lim Sₙ."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    "The series 1+2+3 has partial sum S_3. What is S_3? (enter a number)",
			Answer:      "6",
			Explanation: "S_3 = 1+2+3 = 6.",
		}
	}
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type harmonicSeriesGen struct{}

func (g *harmonicSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does the harmonic series \\(\\sum 1/n\\) converge or diverge?", "diverges", "The harmonic series \\(\\sum 1/n\\) diverges (slowly, like \\(\\ln(n)\\)), even though its terms \\(\\to 0\\)."},
		{"For the p-series \\(\\sum 1/n^{p}\\): does p=2 converge? (yes/no)", "yes", "\\(\\sum 1/n^{p}\\) converges if \\(p>1\\), diverges if \\(p\\leq 1\\). The harmonic series \\((p=1)\\) is the boundary case."},
		{"The nth partial sum H_n grows like ln(n). Does H_n stay bounded? (yes/no)", "no", "The harmonic numbers \\(H_{n} = \\sum_{k=1}^{n} 1/k \\approx \\ln(n) + \\gamma\\) grow without bound, where \\(\\gamma\\) is the Euler-Mascheroni constant."},
		{"Does \\(\\sum 1/n^{2}\\) converge? (yes/no)", "yes", "\\(\\sum 1/n^{2}\\) converges \\((p=2>1)\\). Its sum is \\(\\pi^{2}/6 \\approx 1.645\\) (Basel problem)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type positiveTermsGen struct{}

func (g *positiveTermsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does a positive term series have all terms >= 0? (yes/no)", "yes", "A positive term series ∑ aₙ has aₙ ≥ 0 for all n."},
		{"Which test compares a series to another known series? (enter comparison)", "comparison", "Comparison test: if 0 ≤ aₙ ≤ bₙ and ∑ bₙ converges, then ∑ aₙ converges. If ∑ aₙ diverges, then ∑ bₙ diverges."},
		{"If a_n/b_n tends to a positive constant, do both series share the same fate? (yes/no)", "yes", "Limit comparison: for positive series, if aₙ/bₙ → c where 0<c<∞, then ∑ aₙ and ∑ bₙ either both converge or both diverge; if aₙ/bₙ → c>0, both series have same fate."},
		{"Which test uses the ratio of successive terms? (enter ratio)", "ratio", "Ratio test compares the ratio of successive terms a(n+1)/a(n); if lim = L < 1, converges absolutely. If L > 1, diverges. Inconclusive if L=1."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type alternatingSeriesGen struct{}

func (g *alternatingSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does an alternating series have terms that alternate in sign? (yes/no)", "yes", "An alternating series has terms that alternate between positive and negative: sum (-1)^n a(n)."},
		{"Does the alternating series test need decreasing terms to 0? (yes/no)", "yes", "If |a(n)| decreases monotonically and a(n) → 0, then sum (-1)^n a(n) converges; a(n) decreasing to 0."},
		{"Does \\(\\sum (-1)^{n}/n\\) converge?", "yes", "By the alternating series test: 1/n decreases and → 0. So it converges. It converges conditionally (sum 1/n diverges)."},
		{"Does \\(\\sum (-1)^{n}\\) converge?", "no", "The terms do not → 0 (they alternate between -1 and 1). The series diverges by the nth term test."},
		{"Bound the error stopping sum (-1)^n/n at n=10. (enter 1/11)", "1/11", "For a convergent alternating series with decreasing terms, |S - S(n)| ≤ a(n+1); remainder ≤ first omitted term 1/11."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralTestGen struct{}

func (g *integralTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does the integral test compare a series against an improper integral? (yes/no)", "yes", "If \\(f\\) is positive, continuous, and decreasing on \\([1,\\infty)\\), then \\(\\sum f(n)\\) converges iff \\(\\int_{1}^{\\infty} f(x) \\, dx\\) converges."},
		{"Does \\(\\sum 1/n^{2}\\) converge? Use integral test.", "yes", "\\(\\int_{1}^{\\infty} 1/x^{2} \\, dx = [-1/x]_{1}^{\\infty} = 1\\). The integral converges, so the series converges."},
		{"Does \\(\\sum 1/n\\) converge? Use integral test.", "no", "\\(\\int_{1}^{\\infty} 1/x \\, dx = [\\ln(x)]_{1}^{\\infty} = \\infty\\). The integral diverges, so the series diverges."},
		{"Apply integral test to \\(\\sum 1/(n \\ln n)\\) for \\(n\\geq 2\\).", "diverges", "\\(\\int_{2}^{\\infty} 1/(x \\ln x) \\, dx = [\\ln(\\ln x)]_{2}^{\\infty} = \\infty\\). So \\(\\sum 1/(n \\ln n)\\) diverges."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is a function series one where each term is a function? (yes/no)", "yes", "A function series sum f(n)(x) is the sum of a sequence of functions; sum f_n(x) where each term is a function."},
		{"Does pointwise convergence mean convergence for each fixed x? (yes/no)", "yes", "The series converges pointwise if for each x in the domain, the numerical series converges; sum f_n(x) converges for each fixed x."},
		{"Does uniform convergence mean partial sums converge uniformly? (yes/no)", "yes", "The series converges uniformly if the sequence of partial sums converges uniformly; partial sums converge uniformly."},
		{"Which convergence does the Weierstrass M-test give? (enter uniform)", "uniform", "Weierstrass M-test for testing uniform convergence: if |f(n)(x)| ≤ M(n) and sum M(n) converges, then sum f(n)(x) converges uniformly and absolutely."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivTrigOtherGen struct{}

func (g *derivTrigOtherGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=\\tan(x)\\).", "sec(x)^2", "\\(\\frac{d}{dx} \\tan(x) = \\sec^{2}(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\sec(x)\\).", "sec(x)tan(x)", "\\(\\frac{d}{dx} \\sec(x) = \\sec(x)\\tan(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\csc(x)\\).", "-csc(x)cot(x)", "\\(\\frac{d}{dx} \\csc(x) = -\\csc(x)\\cot(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\cot(x)\\).", "-csc(x)^2", "\\(\\frac{d}{dx} \\cot(x) = -\\csc^{2}(x)\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\tan(2x)\\).", "2sec(2x)^2", "Chain rule: \\(\\frac{d}{dx} \\tan(2x) = \\sec^{2}(2x) \\cdot 2\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\sec(3x)\\).", "3sec(3x)tan(3x)", "Chain rule: \\(\\frac{d}{dx} \\sec(3x) = \\sec(3x)\\tan(3x) \\cdot 3\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivInverseTrigGen struct{}

func (g *derivInverseTrigGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=\\arcsin(x)\\).", "1/sqrt(1-x^2)", "\\(\\frac{d}{dx} \\arcsin(x) = \\frac{1}{\\sqrt{1-x^{2}}}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\arccos(x)\\).", "-1/sqrt(1-x^2)", "\\(\\frac{d}{dx} \\arccos(x) = -\\frac{1}{\\sqrt{1-x^{2}}}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\arctan(x)\\).", "1/(1+x^2)", "\\(\\frac{d}{dx} \\arctan(x) = \\frac{1}{1+x^{2}}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\arctan(2x)\\).", "2/(1+4x^2)", "Chain rule: \\(\\frac{d}{dx} \\arctan(2x) = \\frac{1}{1+(2x)^{2}} \\cdot 2\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\arcsin(3x)\\).", "3/sqrt(1-9x^2)", "Chain rule: \\(\\frac{d}{dx} \\arcsin(3x) = \\frac{1}{\\sqrt{1-(3x)^{2}}} \\cdot 3\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivGeneralExpGen struct{}

func (g *derivGeneralExpGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=2^{x}\\).", "2^x*ln(2)", "\\(\\frac{d}{dx} 2^{x} = 2^{x}\\ln 2\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=10^{x}\\).", "10^x*ln(10)", "\\(\\frac{d}{dx} 10^{x} = 10^{x}\\ln 10\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=3^{x}\\).", "3^x*ln(3)", "\\(\\frac{d}{dx} 3^{x} = 3^{x}\\ln 3\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=2^{3x}\\).", "3*2^(3x)*ln(2)", "Chain rule: \\(\\frac{d}{dx} 2^{3x} = 2^{3x}\\ln 2 \\cdot 3\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=5^{2x}\\).", "2*5^(2x)*ln(5)", "Chain rule: \\(\\frac{d}{dx} 5^{2x} = 5^{2x}\\ln 5 \\cdot 2\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivLogDiffGen struct{}

func (g *derivLogDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=x^{x}\\).", "x^x*(ln(x)+1)", "Take logs: \\(\\ln f = x\\ln x\\), so \\(f'/f = \\ln x + 1\\)."},
		{"Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=x^{2x}\\).", "x^(2x)*(2ln(x)+2)", "Take logs: \\(\\ln f = 2x\\ln x\\), so \\(f'/f = 2\\ln x + 2\\)."},
		{"What is the first step of logarithmic differentiation? (yes/no: take the natural log of both sides)", "yes", "Logarithmic differentiation starts by taking \\(\\ln\\) of both sides (take the natural log of both sides), then differentiating implicitly."},
		{"Use logarithmic differentiation: find \\(f'(x)\\) if \\(f(x)=(2x)^{x}\\).", "(2x)^x*(ln(2x)+1)", "Take logs: \\(\\ln f = x\\ln(2x)\\), so \\(f'/f = \\ln(2x) + 1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivHigherOrderGen struct{}

func (g *derivHigherOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f''(x)\\) if \\(f(x)=x^{4}\\).", "12x^2", "\\(f'(x) = 4x^{3}\\), so \\(f''(x) = 12x^{2}\\)."},
		{"Find \\(f''(x)\\) if \\(f(x)=\\sin(x)\\).", "-sin(x)", "\\(f'(x) = \\cos(x)\\), so \\(f''(x) = -\\sin(x)\\)."},
		{"Find \\(f''(x)\\) if \\(f(x)=e^{2x}\\).", "4e^(2x)", "\\(f'(x) = 2e^{2x}\\), so \\(f''(x) = 4e^{2x}\\)."},
		{"Find the third derivative of \\(f(x)=x^{5}\\).", "60x^2", "\\(f' = 5x^{4}\\), \\(f'' = 20x^{3}\\), \\(f''' = 60x^{2}\\)."},
		{"Find \\(f''(x)\\) if \\(f(x)=\\ln(x)\\).", "-1/x^2", "\\(f'(x) = 1/x\\), so \\(f''(x) = -1/x^{2}\\)."},
		{"Find \\(f''(x)\\) if \\(f(x)=x^{3}+2x\\).", "6x", "\\(f'(x) = 3x^{2}+2\\), so \\(f''(x) = 6x\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivInverseFuncGen struct{}

func (g *derivInverseFuncGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If \\(f\\) is invertible and \\(f'(2)=5\\), what is \\((f^{-1})'(f(2))\\)?", "1/5", "Inverse function rule: \\((f^{-1})'(y) = 1/f'(x)\\) where \\(y=f(x)\\)."},
		{"If \\(f\\) is invertible and \\(f'(1)=4\\), what is \\((f^{-1})'(f(1))\\)?", "1/4", "By the inverse function rule, \\((f^{-1})'(f(1)) = 1/f'(1) = 1/4\\)."},
		{"If f is differentiable with nonzero derivative, does (f⁻¹)'(f(x)) equal 1/f'(x)? (yes/no)", "yes", "If \\(f\\) is differentiable with \\(f'(x) \\neq 0\\), then \\((f^{-1})'(f(x)) = 1/f'(x)\\)."},
		{"If \\(f(x)=x^{3}\\) and \\(g\\) is its inverse, what is \\(g'(8)\\)?", "1/12", "\\(f'(x) = 3x^{2}\\); at \\(x=2\\), \\(f(2)=8\\) and \\(f'(2)=12\\), so \\(g'(8)=1/12\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivCriticalPointsGen struct{}

func (g *derivCriticalPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find the critical points of \\(f(x)=x^{2}-4x\\).", "x=2", "\\(f'(x) = 2x-4 = 0\\) gives \\(x=2\\)."},
		{"How many critical points does \\(f(x)=x^{3}-3x\\) have? (enter a number)", "2", "\\(f'(x) = 3x^{2}-3 = 0\\) gives \\(x=-1\\) and \\(x=1\\): 2 critical points."},
		{"Is \\(x=2\\) a critical point of \\(f(x)=x^{2}-4x\\)? (yes/no)", "yes", "A critical point is an interior point where f'(x)=0 or f' is undefined. Here \\(f'(x)=2x-4\\), so \\(f'(2)=0\\)."},
		{"Find the critical points of \\(f(x)=x^{3}\\).", "x=0", "\\(f'(x) = 3x^{2} = 0\\) gives \\(x=0\\) (a stationary point that is not an extremum)."},
		{"How many critical points does \\(f(x)=2x^{3}-6x^{2}\\) have? (enter a number)", "2", "\\(f'(x) = 6x^{2}-12x = 6x(x-2) = 0\\) gives \\(x=0\\) and \\(x=2\\): 2 critical points."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	table := []entry{
		{"How many inflection points does \\(f(x)=x^{3}\\) have? (enter a number)", "1", "\\(f''(x) = 6x = 0\\) at \\(x=0\\), with concavity change there: one inflection point."},
		{"How many inflection points does \\(f(x)=x^{4}-6x^{2}\\) have? (enter a number)", "2", "\\(f''(x) = 12x^{2}-12 = 0\\) at \\(x=\\pm 1\\), with concavity change on both sides: two points."},
		{"Does an inflection point require concavity to change sign? (yes/no)", "yes", "An inflection point is where \\(f''\\) changes sign, i.e. concavity flips."},
		{"Does \\(f(x)=x^{4}\\) have an inflection point at \\(x=0\\)?", "no", "\\(f''(0)=0\\) but \\(f''(x) = 12x^{2} \\geq 0\\) everywhere, so no sign change."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivGlobalExtremaGen struct{}

func (g *derivGlobalExtremaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the global maximum VALUE of \\(f(x)=x^{2}-4x\\) on \\([0,3]\\)? (enter a number)", "0", "Critical point \\(x=2\\): \\(f(2)=-4\\). Endpoints: \\(f(0)=0\\), \\(f(3)=-3\\). Global max value is 0 at \\(x=3\\)."},
		{"What is the global minimum VALUE of \\(f(x)=x^{2}-4x\\) on \\([0,3]\\)? (enter a number)", "-4", "Critical point \\(x=2\\): \\(f(2)=-4\\). Endpoints give 0 and -3. Global min value is -4 at \\(x=2\\)."},
		{"Does the closed-interval method evaluate f at critical points and endpoints? (yes/no)", "yes", "On \\([a,b]\\): find critical points inside, evaluate \\(f\\) there and at \\(a,b\\); extremes are the largest/smallest values."},
		{"What is the global maximum VALUE of \\(f(x)=-x^{2}+4x\\) on \\([0,3]\\)? (enter a number)", "4", "Critical point \\(x=2\\): \\(f(2)=4\\). Endpoints: \\(f(0)=0\\), \\(f(3)=3\\). Global max value is 4 at \\(x=2\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivCurveSketchGen struct{}

// Every entry is a choice now. Three of the four used to answer "yes"/"no" on a
// concept graded symbolic, so the client showed a maths field: a learner typed
// y for yes and was marked wrong (report 118), and no form of the answer was
// typeable in a maths editor.
func (g *derivCurveSketchGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q       string
		a       string
		e       string
		options []string
	}
	table := []entry{
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
	e := table[rand.Intn(len(table))]
	distractors := make([]string, 0, len(e.options))
	for _, o := range e.options {
		if o != e.a {
			distractors = append(distractors, o)
		}
	}
	p := generator.MultipleChoice(e.q, e.a, distractors...)
	p.Explanation = e.e
	return p
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
			"Subgoal 1: Differentiate: \\(f'(x)=%d%s\\).\n"+
				"Subgoal 2: The slope at \\(x=%d\\) is \\(f'(%d)=%d\\) and the point of tangency is \\((%d,%d)\\).\n"+
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
			"Subgoal 1: \\(f(%d)=%d\\) and \\(f'(x)=%d%s\\), so \\(f'(%d)=%d\\).\n"+
				"Subgoal 2: The linearization is \\(L(x)=%d+%d(x-%d)\\).\n"+
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
			"Subgoal 1: Newton's method uses \\(x_{1}=x_{0}-f(x_{0})/f'(x_{0})\\) with \\(f'(x)=2x\\).\n"+
				"Subgoal 2: \\(f(%d)=%d-%d=%d\\) and \\(f'(%d)=%d\\).\n"+
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
			"Subgoal 1: For small changes \\(\\Delta y \\approx dy = f'(x)\\,dx\\).\n"+
				"Subgoal 2: \\(f'(x)=%d%s\\), so \\(f'(%d)=%d\\).\n"+
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
				"Subgoal 1: Differentiate both sides implicitly: \\(2x+2y\\,y'=0\\).\n"+
					"Subgoal 2: Solve for \\(y'\\): \\(y'=-x/y\\).\n"+
					"Subgoal 3: This holds on the circle \\(x^{2}+y^{2}=%d\\) wherever \\(y\\neq 0\\).", r*r),
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("For \\(x^{2}+y^{2}=%d\\), what is \\(d^{2}y/dx^{2}\\) in terms of \\(x,y\\)?", r*r),
		Answer:   fmt.Sprintf("-%d/y^3", r*r),
		Explanation: fmt.Sprintf(
			"Subgoal 1: From \\(y'=-x/y\\), differentiate again: \\(y''=-(y-xy')/y^{2}\\).\n"+
				"Subgoal 2: Substitute \\(y'=-x/y\\): \\(y''=-(y+x^{2}/y)/y^{2}=-(y^{2}+x^{2})/y^{3}\\).\n"+
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
			"Subgoal 1: The outer power rule gives \\(%d(%s)^{%d}\\).\n"+
				"Subgoal 2: The inner derivative is %d.\n"+
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
			"Subgoal 1: The multivariable chain rule is \\(dz/dt = z_{x}\\,x'(t) + z_{y}\\,y'(t)\\).\n"+
				"Subgoal 2: \\(z_{x}=%d\\), \\(z_{y}=%d\\), \\(x'(t)=1\\), \\(y'(t)=2t\\).\n"+
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
			"Subgoal 1: \\(\\nabla f=(f_{x},f_{y})\\).\n"+
				"Subgoal 2: Differentiate in \\(x\\) holding \\(y\\) fixed: \\(f_{x}=%s\\).\n"+
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
			"Subgoal 1: \\(D_{u}f=\\nabla f\\cdot u\\) for a unit direction \\(u\\).\n"+
				"Subgoal 2: \\(\\nabla f=(%s, %s)\\), so at \\((%d,%d)\\) it is \\((%d, %d)\\).\n"+
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
			"Subgoal 1: Differentiate once in each variable the label names, holding the other fixed.\n"+
				"Subgoal 2: Powers drop by one per derivative, so the surviving exponents are \\(x^{%d}y^{%d}\\).\n"+
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
			"Subgoal 1: Differentiate in \\(x\\) first: \\(f_{x}=%s\\).\n"+
				"Subgoal 2: Then differentiate that in \\(y\\): \\(f_{xy}=%s\\).\n"+
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
			"Subgoal 1: The total differential is \\(dz=z_{x}\\,dx+z_{y}\\,dy\\).\n"+
				"Subgoal 2: At \\((%d,%d)\\), \\(z_{x}=%d\\) and \\(z_{y}=%d\\).\n"+
				"Subgoal 3: \\(dz\\approx %d(%.1f)+%d(%.1f)=%s\\).",
			x0, y0, fx, fy, fx, dx, fy, dy, decimal(dz, 4)),
	}
}

type derivDiffContinuityGen struct{}

func (g *derivDiffContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If \\(f\\) is differentiable at \\(a\\), must \\(f\\) be continuous at \\(a\\)? (yes/no)", "yes", "Differentiability at \\(a\\) implies continuity at \\(a\\): f is continuous at a."},
		{"Does continuity imply differentiability?", "no", "Continuity does not imply differentiability (e.g. \\(|x|\\) at 0)."},
		{"Give a function continuous but not differentiable at 0.", "|x|", "\\(f(x)=|x|\\) is continuous at 0 but has no derivative there (corner)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
				"Subgoal 1: Row 1 holds the partials of the first component \\(x^{%d}\\): \\(%s\\) in \\(x\\) and 0 in \\(y\\).\n"+
					"Subgoal 2: Row 2 holds the partials of \\(y^{%d}\\): 0 in \\(x\\) and \\(%s\\) in \\(y\\).\n"+
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
			"Subgoal 1: The first row holds \\(\\partial T_{1}/\\partial x=%d\\) and \\(\\partial T_{1}/\\partial y=%d\\).\n"+
				"Subgoal 2: The second row holds \\(\\partial T_{2}/\\partial x=%d\\) and \\(\\partial T_{2}/\\partial y=%d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Net Change Theorem: displacement equals \\(\\int_{t_1}^{t_2} v(t)\\,dt\\).\n"+
				"Subgoal 2: Substitute \\(v(t) = %d t^2\\), \\(t_1 = 0\\) and \\(t_2 = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: net volume change is \\(\\int_{a}^{b} r(t)\\,dt\\).\n"+
				"Subgoal 2: Substitute integrand \\(r(t) = %d t + %d\\) with limits from 0 to %d.\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: Apply the Net Change Theorem: \\(f(b) = f(a) + \\int_{a}^{b} f'(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute parameters \\(f(0) = %d\\) and \\(\\int_{0}^{%d} %d x\\,dx\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: net accumulation is \\(\\Delta V = \\int_{0}^{T} r(t)\\,dt\\).\n"+
				"Subgoal 2: Integrate \\(\\int_{0}^{%d} %d t^2\\,dt = [%d t^3]_{0}^{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Net Change Theorem for acceleration: \\(v(T) = v(0) + \\int_{0}^{T} a(t)\\,dt\\).\n"+
				"Subgoal 2: Substitute \\(v(0) = %d\\) and \\(\\int_{0}^{%d} %d\\,dt\\).\n"+
				"Subgoal 3: Evaluate and compute: \\(%d + (%d \\cdot %d) = %d\\).", v0, T, coeff, v0, coeff, T, ans),
		}
	}
}

type integralVariableLimitsGen struct{}

func (g *integralVariableLimitsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\frac{d}{dx}\\int_{0}^{x} t^{2}\\,dt\\).", "x^2", "FTC part 1: derivative of the accumulation is the integrand at \\(x\\)."},
		{"Find \\(\\frac{d}{dx}\\int_{1}^{x} \\sin(t)\\,dt\\).", "sin(x)", "By FTC1, \\(\\frac{d}{dx}\\int_{1}^{x}\\sin t\\,dt = \\sin x\\)."},
		{"Find \\(\\frac{d}{dx}\\int_{0}^{x^{2}} t\\,dt\\).", "2x^3", "Chain rule + FTC1: \\(x^{2} \\cdot 2x = 2x^{3}\\)."},
		{"Find \\(\\frac{d}{dx}\\int_{x}^{2} e^{t}\\,dt\\).", "-e^x", "Flipped limit: \\(-\\frac{d}{dx}\\int_{2}^{x}e^{t}\\,dt = -e^{x}\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute parameters \\(a = 0\\), \\(b = %d\\) and \\(f(x) = %d x\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the average value formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute limits from 0 to %d with integrand \\(%d x^2\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the theorem: there exists \\(c \\in [a, b]\\) such that \\(f(c) = f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\n"+
				"Subgoal 2: Compute \\(f_{\\text{avg}} = \\frac{1}{%d} \\int_{0}^{%d} 2x\\,dx = \\frac{1}{%d} [x^2]_{0}^{%d} = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\n"+
				"Subgoal 2: Compute \\(\\int_{0}^{%d} (%d x + %d)\\,dx = [%d x^2 + %d x]_{0}^{%d} = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(f_{\\text{avg}} = \\frac{1}{b-a}\\int_{a}^{b} f(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute integrand \\(%d x^3\\) on interval [0, %d].\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for work: \\(W = \\int_{a}^{b} F(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute \\(F(x) = %d x\\) with limits \\(a = 0\\) and \\(b = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{0}^{D} F(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute force \\(F(x) = %d x^2\\) and distance \\(D = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{x_1}^{x_2} F(x)\\,dx\\).\n"+
				"Subgoal 2: Substitute \\(F(x) = %d x\\), \\(x_1 = %d\\) and \\(x_2 = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = \\int_{0}^{D} (mx + C)\\,dx\\).\n"+
				"Subgoal 2: Compute antiderivative \\([%d x^2 + %d x]_{0}^{%d}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(W = F \\cdot d\\).\n"+
				"Subgoal 2: Substitute constant force \\(F = %d\\) N and displacement \\(d = %d\\) m.\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula for a uniform rod: \\(\\bar{x} = \\frac{a + b}{2}\\).\n"+
				"Subgoal 2: Substitute endpoints \\(a = 0\\) and \\(b = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{M_y}{M} = \\frac{\\int_{0}^{L} x \\rho(x)\\,dx}{\\int_{0}^{L} \\rho(x)\\,dx}\\).\n"+
				"Subgoal 2: Compute total mass \\(M = \\int_{0}^{%d} x\\,dx = \\frac{%d^2}{2}\\) and moment \\(M_y = \\int_{0}^{%d} x^2\\,dx = \\frac{%d^3}{3}\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{\\int_{0}^{L} x^3\\,dx}{\\int_{0}^{L} x^2\\,dx}\\).\n"+
				"Subgoal 2: Evaluate integrals: numerator is \\(L^4 / 4\\) and denominator is \\(L^3 / 3\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: by symmetry, \\(\\bar{x} = \\frac{a + b}{2}\\).\n"+
				"Subgoal 2: Substitute endpoints \\(a = %d\\) and \\(b = %d\\).\n"+
				"Subgoal 3: Calculate: \\(\\bar{x} = \\frac{%d + %d}{2} = %d\\).", a, b, a, b, ans),
		}
	default:
		// Two equal point masses m at x = 0 and x = 2d -> xbar = d
		d := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		ans := d
		return generator.Problem{
			Question: fmt.Sprintf("Two equal point masses \\(m\\) are placed at positions \\(x = 0\\) and \\(x = %d\\). Find the center of mass of the system.", 2*d),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\bar{x} = \\frac{\\sum m_i x_i}{\\sum m_i}\\).\n"+
				"Subgoal 2: Substitute \\(m_1 = m_2 = m\\), \\(x_1 = 0\\) and \\(x_2 = %d\\).\n"+
				"Subgoal 3: Evaluate: \\(\\bar{x} = \\frac{m(0) + m(%d)}{2m} = \\frac{%d}{2} = %d\\).", 2*d, 2*d, 2*d, ans),
		}
	}
}

type integralSurfaceAreaGen struct{}

func (g *integralSurfaceAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the surface area of a sphere of radius 2? (enter in terms of pi)", "16pi", "Revolving a semicircle gives S = 4 pi r^2; with r=2: 16pi. Bands of radius f(x): S = 2 pi times integral of f(x)sqrt(1+(f')^2)dx."},
		{"What is the surface area of a sphere of radius \\(r\\)? (enter in terms of pi)", "4*pi*r^2", "Revolving a semicircle gives \\(S = 4\\pi r^{2}\\)."},
		{"For \\(y=3x\\), what is \\(1+(f'(x))^{2}\\)? (enter a number)", "10", "f'=3, so 1+(f')^2=10. The arc-length element sqrt(1+(f')^2)dx measures slant length in the surface formula."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule for \\(n = 1\\): \\(T_1 = \\frac{b-a}{2} (f(a) + f(b))\\).\n"+
				"Subgoal 2: Substitute \\(a = 0\\), \\(b = %d\\), \\(f(0) = 0\\) and \\(f(%d) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule: \\(T_n = \\frac{\\Delta x}{2} [f(x_0) + 2f(x_1) + f(x_2)]\\).\n"+
				"Subgoal 2: Note that the Trapezoidal Rule is exact for linear polynomials, so \\(T_2 = \\int_{0}^{%d} (%d x + %d)\\,dx\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\Delta x = \\frac{b-a}{n} = %d\\), \\(T_2 = \\frac{%d}{2} [f(0) + 2f(%d) + f(%d)]\\).\n"+
				"Subgoal 2: Substitute function values \\(f(0) = 0\\), \\(f(%d) = %d\\) and \\(f(%d) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: \\(\\Delta x = \\frac{b - a}{n} \\implies n = \\frac{b - a}{\\Delta x}\\).\n"+
				"Subgoal 2: Substitute parameters \\(a = %d\\), \\(b = %d\\) and \\(\\Delta x = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State the Trapezoidal Rule: \\(T_1 = \\frac{b-a}{2} (f(a) + f(b))\\).\n"+
				"Subgoal 2: For \\([0, 2]\\), \\(\\frac{b-a}{2} = 1\\), \\(f(0) = 0\\) and \\(f(2) = %d(4) = %d\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for polynomials of degree up to 3.\n"+
				"Subgoal 2: Recognize that \\(S_2 = \\int_{0}^{%d} x^3\\,dx\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for quadratic polynomials.\n"+
				"Subgoal 2: Evaluate \\(S_2 = \\int_{0}^{%d} 3x^2\\,dx = [x^3]_{0}^{%d}\\).\n"+
				"Subgoal 3: Compute: \\(%d^3 = %d\\).", b, b, b, ans),
		}
	case 2:
		// What is the smallest even number of subintervals n usable in Simpson's rule?
		return generator.Problem{
			Question: "What is the smallest number of subintervals \\(n\\) required to apply Simpson's Rule?",
			Answer:   "2",
			Explanation: "Subgoal 1: State the requirement for Simpson's Rule: parabolas are fit across pairs of adjacent subintervals.\n" +
				"Subgoal 2: Deduce that \\(n\\) must be an even positive integer.\n" +
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule \\(S_2 = \\frac{\\Delta x}{3} [f(x_0) + 4f(x_1) + f(x_2)]\\).\n"+
				"Subgoal 2: Because the integrand is linear, Simpson's Rule gives the exact integral value \\(\\int_{0}^{%d} %d x\\,dx\\).\n"+
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
			Explanation: fmt.Sprintf("Subgoal 1: State governing formula: Simpson's Rule is exact for quadratic polynomials.\n"+
				"Subgoal 2: Compute exact integral \\(\\int_{0}^{%d} x^2\\,dx = [\\frac{x^3}{3}]_{0}^{%d}\\).\n"+
				"Subgoal 3: Calculate: \\(\\frac{%d^3}{3} = \\frac{%d}{3} = %d\\).", b, b, b, b*b*b, ans),
		}
	}
}

type integralImproperCompareGen struct{}

func (g *integralImproperCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"For \\(\\int_{1}^{\\infty} dx/x^{p}\\), what is the smallest integer \\(p\\) giving convergence? (enter a number)", "2", "Comparison test: the p-integral converges iff p>1, so the smallest integer is 2. Domination (0≤f≤g with ∫g convergent) preserves convergence; divergence of the smaller forces divergence of the larger."},
		{"Does \\(\\int_{1}^{\\infty} dx/(x^{2}+1)\\) converge? (yes/no)", "yes", "\\(1/(x^{2}+1) \\leq 1/x^{2}\\) and \\(\\int_{1}^{\\infty}dx/x^{2}\\) converges, so yes by comparison with 1/x^2."},
		{"Does \\(\\int_{1}^{\\infty} dx/\\sqrt{x}\\) converge? (yes/no)", "no", "By p-comparison with p=1/2 ≤ 1 it diverges: no, it diverges like p=1/2."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralPTestGen struct{}

func (g *integralPTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the smallest integer p for which \\(\\int_{1}^{\\infty}dx/x^{p}\\) converges? (enter a number)", "2", "The \\(p\\)-integral converges iff \\(p>1\\): smallest integer 2."},
		{"Does \\(\\int_{1}^{\\infty}dx/x^{3}\\) converge? (yes/no)", "yes", "Since \\(p=3>1\\), the integral converges."},
		{"Does \\(\\int_{1}^{\\infty}dx/x\\) converge? (yes/no)", "no", "The harmonic case \\(p=1\\) diverges (logarithmic growth)."},
		{"Does \\(\\int_{0}^{1}dx/x^{p}\\) converge for \\(p<1\\)? (yes/no)", "yes", "At 0, \\(\\int_{0}^{1}dx/x^{p}\\) converges iff \\(p<1\\) (flipped condition)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find the sum of \\(\\sum_{n=0}^{\\infty}(1/2)^{n}\\).", "2", "Geometric with \\(a=1\\), \\(r=1/2\\): \\(1/(1-1/2) = 2\\)."},
		{"Does sum r^n converge for r=1/2? (yes/no)", "yes", "Geometric sum a r^n converges iff |r|<1, summing to a/(1-r); for r=1/2 it converges."},
		{"Find the sum of \\(\\sum_{n=0}^{\\infty}(1/3)^{n}\\).", "3/2", "\\(a=1\\), \\(r=1/3\\): \\(1/(1-1/3) = 3/2\\)."},
		{"Does sum 2^n converge? (yes/no)", "no", "Ratio |r|=2>=1: terms do not even tend to zero, so it diverges."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type nthTermGen struct{}

func (g *nthTermGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If a_n does not tend to 0, does the series diverge? (yes/no)", "yes", "Nth-term divergence test: convergence of sum a(n) requires a(n)→0; if a_n does not tend to 0, the series diverges."},
		{"Does \\(\\sum n/(n+1)\\) converge?", "no", "Terms tend to 1, not 0, so the series diverges by the nth-term test."},
		{"If \\(\\sum a_{n}\\) converges, what is \\(\\lim a_{n}\\)?", "0", "Convergence forces \\(a_{n}\\to 0\\) (necessary, not sufficient)."},
		{"Does a_n to 0 guarantee convergence? (yes/no)", "no", "No. The harmonic series has 1/n→0 yet diverges."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Which test uses domination by a convergent series? (enter comparison)", "comparison", "Direct comparison: 0≤a(n)≤b(n) and sum b(n) converges implies sum a(n) converges; divergence below forces divergence above."},
		{"Does sum 1/(n^2+1) converge? (yes/no)", "yes", "1/(n^2+1) ≤ 1/n^2 and sum 1/n^2 converges; yes, compare with 1/n^2."},
		{"Does sum 1/sqrt(n^2+n) converge? (yes/no)", "no", "Terms behave like 1/n; sum diverges by comparison below with a divergent harmonic series; no, compare below with a divergent series."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type limitCompareGen struct{}

func (g *limitCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Which test takes the limit of a_n/b_n? (enter comparison)", "comparison", "Limit comparison: finite positive limit implies same behavior; if lim a(n)/b(n) = L in (0,infinity), both series converge or both diverge."},
		{"Does sum 1/(n^2-n) converge? (yes/no)", "yes", "Limit comparison on sum 1/(n^2-n): ratio with 1/n^2 tends to 1; sum 1/n^2 converges, so it converges (compare 1/n^2)."},
		{"Does sum n/(n^3+1) converge? (yes/no)", "yes", "Limit comparison on sum n/(n^3+1): ratio with 1/n^2 tends to 1; converges (compare 1/n^2)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type absoluteConvergenceGen struct{}

func (g *absoluteConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does absolute convergence mean sum |a_n| converges? (yes/no)", "yes", "A series converges absolutely if the series of absolute values converges; sum |a(n)| converges."},
		{"Does absolute convergence imply convergence?", "yes", "Absolute convergence is stronger: \\(\\sum|a_{n}|\\) convergent forces \\(\\sum a_{n}\\) convergent."},
		{"Is \\(\\sum (-1)^{n}/n^{2}\\) absolutely convergent?", "yes", "\\(\\sum 1/n^{2}\\) converges, so the alternating version converges absolutely."},
		{"Does sum (-1)^n/n converge absolutely? (yes/no)", "no", "sum 1/n diverges, so convergence is conditional only."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(3) == 0 {
		return generator.Problem{
			Question:    "sum (-1)^n/n converges but sum 1/n diverges. What kind of convergence? (type absolute or conditional)",
			Answer:      "conditional",
			Explanation: "Converges (Leibniz) but not absolutely: conditional convergence.",
		}
	}
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is \\(a_{n}=(-1)^{n}\\) bounded? (yes/no)", "yes", "It oscillates between -1 and 1: bounded (|a_n|≤1) but divergent."},
		{"Does every convergent sequence have to be bounded?", "yes", "Convergence forces boundedness (terms eventually stay near the limit)."},
		{"Is \\(a_{n}=n^{2}\\) bounded?", "no", "It grows without bound."},
		{"Is a bounded sequence always convergent?", "no", "\\((-1)^{n}\\) is bounded yet oscillates: boundedness is necessary, not sufficient."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
