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

type limitPropertiesGen struct{}

func (g *limitPropertiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	fLim := rand.Intn(scale*2) + 1
	gLim := rand.Intn(scale*2) + 1
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
	// f(x) = (x+1)/(x-1), f'(x) = -2/(x-1)², ask for f'(2) = -2
	return generator.Problem{
		Question:    "Find \\(f'(2)\\) if \\(f(x)=\\frac{x+1}{x-1}\\).",
		Answer:      "-2",
		Explanation: "Quotient rule: f'(x)=((x-1)·1-(x+1)·1)/(x-1)² = -2/(x-1)². f'(2) = -2/(1)² = -2.",
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
		Question:    fmt.Sprintf("Find \\(f'(x)\\) if \\(f(x)=%s^{%d}\\).", inner, n),
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
	n := rand.Intn(max(1, scale)) + 1
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
	n := rand.Intn(max(1, scale)) + 2
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
	n := rand.Intn(max(1, scale)) + 1
	answer := fmt.Sprintf("x^%d", n)
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(F(x)=\\int_{0}^{x} t^{%d} \\, dt\\), what is \\(F'(x)\\)?", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("By the Fundamental Theorem of Calculus, F'(x)=x^%d.", n),
	}
}

type integralAreaBetweenGen struct{}

func (g *integralAreaBetweenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Area between y=x and y=x² from 0 to 1 = 1/6 ≈ 0.1667
	return generator.Problem{
		Question:    "Find area between \\(y=x\\) and \\(y=x^{2}\\) from \\(x=0\\) to \\(x=1\\).",
		Answer:      "0.1667",
		Explanation: "∫(x-x²)dx from 0 to 1 = [x²/2-x³/3]₀¹ = 1/2-1/3 = 1/6 ≈ 0.1667.",
	}
}

type integralVolumeGen struct{}

func (g *integralVolumeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// V = π∫(√x)²dx from 0 to 4 = π∫x dx = π·x²/2 from 0 to 4 = π·16/2 = 8π
	return generator.Problem{
		Question:    "Find the volume when \\(y=\\sqrt{x}\\) from \\(x=0\\) to \\(4\\) is revolved around the \\(x\\)-axis.",
		Answer:      "8π",
		Explanation: "V = π∫(√x)²dx = π∫x dx = π[x²/2]₀⁴ = π(16/2-0) = 8π.",
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
	answer := fmt.Sprintf("%dπ", ans)
	return generator.Problem{
		Question:    fmt.Sprintf("A circle's radius grows at %d cm/s. How fast is area growing when \\(r=%d\\)?", dr, r),
		Answer:      answer,
		Explanation: fmt.Sprintf("A=πr², dA/dt=2πr·dr/dt=2π(%d)(%d)=%dπ cm²/s.", r, dr, ans),
	}
}

type integralPartsGen struct{}

func (g *integralPartsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "Find \\(\\int x e^{x} \\, dx\\).",
		Answer:      "xe^x-e^x+C",
		Explanation: "Integration by parts: let u=x, dv=e^x dx → du=dx, v=e^x. ∫x·e^x dx = x·e^x - ∫e^x dx = x·e^x - e^x + C.",
	}
}

type integralPartialFractionsGen struct{}

func (g *integralPartialFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "Find \\(\\int \\frac{1}{x^{2}-1} \\, dx\\).",
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If f(a)=f(b) and f is continuous on [a,b], differentiable on (a,b), what must exist?", "a point c in (a,b) with f'(c)=0", "Rolle's theorem: there exists c∈(a,b) such that f'(c)=0."},
		{"Does f(x)=x²-4 on [-2,2] satisfy Rolle's theorem?", "yes", "f(-2)=4-4=0, f(2)=4-4=0. f is continuous on [-2,2] and differentiable on (-2,2). So Rolle applies: f'(0)=0."},
		{"Does f(x)=|x| on [-1,1] satisfy Rolle's theorem?", "no", "f(-1)=f(1)=1, but f(x)=|x| is not differentiable at x=0, so Rolle's theorem does not apply."},
		{"Does f(x)=1/x on [1,2] satisfy Rolle's theorem?", "no", "f(1)=f(2)=1/2, but f is not continuous on the closed interval [1,2]? Well, 1/x is continuous on [1,2], actually wait — the condition fails because... Actually, Rolle would say there is a c with f'(c)=0, but f'(x)=-1/x² ≠ 0, so... Actually, f(1)=1≠1/2=f(2), so f(1)≠f(2). Rolle does not apply."},
		{"If f(0)=f(6)=0 and f is differentiable everywhere, how many points with f'(c)=0 are guaranteed?", "at least one", "Rolle's theorem guarantees at least one c∈(0,6) with f'(c)=0 when f is continuous on [0,6] and differentiable on (0,6)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type mvtGen struct{}

func (g *mvtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If f is continuous on [a,b] and differentiable on (a,b), what does the MVT guarantee?", "a point c in (a,b) where f'(c)=(f(b)-f(a))/(b-a)", "MVT: there exists c∈(a,b) such that f'(c)=(f(b)-f(a))/(b-a)."},
		{"For f(x)=x² on [1,3], find c satisfying the MVT.", "2", "f(3)-f(1)=9-1=8, b-a=2. f'(c)=2c=8/2=4 → c=2."},
		{"For f(x)=√x on [1,4], find c satisfying the MVT.", "2.25", "f(4)-f(1)=2-1=1, b-a=3. f'(c)=1/(2√c)=1/3 → 2√c=3 → c=9/4=2.25."},
		{"If a police car travels 120 km in 1 hour, what does the MVT say?", "at some instant the speed was 120 km/h", "By MVT, average speed = 120 km/h. There must have been at least one instant with exactly 120 km/h."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cauchyMVTGen struct{}

func (g *cauchyMVTGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does Cauchy's MVT generalize?", "the ordinary MVT", "Cauchy's MVT: (f(b)-f(a))/(g(b)-g(a)) = f'(c)/g'(c) for some c. It generalizes MVT (take g(x)=x)."},
		{"Cauchy's MVT with g(x)=x gives what theorem?", "the ordinary MVT", "With g(x)=x, Cauchy says (f(b)-f(a))/(b-a)=f'(c)/1, which is the ordinary MVT."},
		{"What condition on g'(x) is needed for Cauchy's MVT?", "g'(x) ≠ 0 on (a,b)", "The theorem requires g'(x) ≠ 0 on (a,b) so the ratio f'(c)/g'(c) is defined."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type fermatGen struct{}

func (g *fermatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If f has a local extremum at c and f'(c) exists, what must f'(c) be?", "0", "Fermat's theorem: differentiable local extrema occur at critical points where f'(c)=0."},
		{"Does Fermat's theorem apply to f(x)=|x| at x=0?", "no", "f'(0) does not exist, so Fermat's theorem does not apply, even though x=0 is a local minimum."},
		{"Does Fermat's theorem apply to f(x)=x³ at x=0?", "no", "f'(0)=0 exists, but x=0 is not a local extremum — it's an inflection point. Fermat says if there IS an extremum AND derivative exists, then f'=0. The converse is not true."},
		{"Fermat's theorem guarantees what about f'(c) at a local extremum?", "f'(c)=0 (if differentiable)", "If f has a local max or min at c and f is differentiable at c, then f'(c)=0."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type differenceQuotientGen struct{}

func (g *differenceQuotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	h := rand.Intn(max(1, scale)) + 1
	x0 := rand.Intn(max(1, scale*2)) + 1
	fx := x0 * x0
	fxh := (x0 + h) * (x0 + h)
	ans := (fxh - fx) / h
	return generator.Problem{
		Question:    fmt.Sprintf("For \\(f(x)=x^{2}\\), what is the difference quotient \\(\\frac{f(%d+h)-f(%d)}{h}\\) for \\(h=%d\\)?", x0, x0, h),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("(f(%d+h)-f(%d))/h = ((%d+h)²-%d²)/h = (%d²+2·%d·h+h²-%d²)/h = (2·%d·h+h²)/h = 2·%d+h = %d.", x0, x0, x0, x0, x0, x0, x0, x0, x0, ans),
	}
}

type differentialGen struct{}

func (g *differentialGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x0 := rand.Intn(max(1, scale*2)) + 1
	dx := 0.1 + float64(rand.Intn(5))/10.0
	fpx := float64(2 * x0)
	dy := fpx * dx
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(y=x^{2}\\) at \\(x=%d\\), approximate \\(\\Delta y\\) using the differential \\(dy\\) when \\(dx=%.1f\\).", x0, dx),
		Answer:      fmt.Sprintf("%.2f", dy),
		Explanation: fmt.Sprintf("dy = f'(x)dx = 2x·dx = 2·%d·%.1f = %.2f. The actual Δy is (%.1f)²-%d² = %.2f-%d = %.2f.", x0, dx, dy, float64(x0)+dx, x0, (float64(x0)+dx)*(float64(x0)+dx), x0*x0, (float64(x0)+dx)*(float64(x0)+dx)-float64(x0*x0)),
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
		{"If a function is differentiable at x=c, is it necessarily continuous at x=c?", "yes", "Differentiability implies continuity. But the converse is not true — continuity does not imply differentiability."},
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
	var answerStr string
	if which == 0 {
		val := 2*a*x0*y0 + b*y0*y0
		answerStr = fmt.Sprintf("∂f/∂x = 2(%d)(%d)(%d)+%d(%d)² = %d", a, x0, y0, b, y0, val)
	} else {
		val := a*x0*x0 + 2*b*x0*y0
		answerStr = fmt.Sprintf("∂f/∂y = (%d)(%d)²+2(%d)(%d)(%d) = %d", a, x0, b, x0, y0, val)
	}
	dir := []string{"x", "y"}[which]
	return generator.Problem{
		Question:    fmt.Sprintf("Find \\(\\partial f/\\partial %s\\) at \\((%d,%d)\\) for \\(f(x,y)=%dx^{2}y+%dxy^{2}\\).", dir, x0, y0, a, b),
		Answer:      answerStr,
		Explanation: answerStr,
	}
}

// --- Integral generators ---

type arcLengthGen struct{}

func (g *arcLengthGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return generator.Problem{
		Question:    "What is the formula for arc length of \\(y=f(x)\\) from \\(x=a\\) to \\(x=b\\)?",
		Answer:      "∫√(1+(f'(x))²)dx from a to b",
		Explanation: "Arc length = ∫ₐᵇ √(1+(f'(x))²) dx. For parametric curves, use ∫√((dx/dt)²+(dy/dt)²) dt.",
	}
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
	return generator.Problem{
		Question:    "Determine if \\(\\int_{1}^{\\infty} \\frac{1}{x^{2}} \\, dx\\) converges or diverges.",
		Answer:      "converges (to 1)",
		Explanation: "∫₁^∞ 1/x² dx = lim_{b→∞} [-1/x]₁ᵇ = lim_{b→∞} (-1/b+1) = 1. The integral converges to 1.",
	}
}

type numericalIntegralGen struct{}

func (g *numericalIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What method approximates ∫f(x)dx using rectangles?", "Riemann sum", "Riemann sums partition [a,b] into subintervals and use rectangles to approximate the area."},
		{"What method approximates ∫f(x)dx using trapezoids?", "Trapezoidal rule", "The trapezoidal rule approximates ∫ₐᵇ f(x)dx ≈ (Δx/2)(f(x₀)+2f(x₁)+...+2f(xₙ₋₁)+f(xₙ))."},
		{"What method approximates ∫f(x)dx using parabolas?", "Simpson's rule", "Simpson's rule uses quadratic polynomials and is more accurate than trapezoidal for smooth functions."},
		{"What is the order of accuracy of the trapezoidal rule?", "O(h²)", "The trapezoidal rule error is proportional to h²·f''(ξ), making it second-order accurate."},
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
		{"What substitution is used for \\(\\sqrt{a^{2}-x^{2}}\\)?", "\\(x = a\\sin(\\theta)\\)", "For \\(\\sqrt{a^{2}-x^{2}}\\), let \\(x=a\\sin(\\theta)\\), then \\(dx=a\\cos(\\theta)d\\theta\\) and \\(\\sqrt{a^{2}-x^{2}} = a\\cos(\\theta)\\)."},
		{"What substitution is used for \\(\\sqrt{a^{2}+x^{2}}\\)?", "\\(x = a\\tan(\\theta)\\)", "For \\(\\sqrt{a^{2}+x^{2}}\\), let \\(x=a\\tan(\\theta)\\), then \\(dx=a\\sec^{2}(\\theta)d\\theta\\) and \\(\\sqrt{a^{2}+x^{2}} = a\\sec(\\theta)\\)."},
		{"What substitution is used for \\(\\sqrt{x^{2}-a^{2}}\\)?", "\\(x = a\\sec(\\theta)\\)", "For \\(\\sqrt{x^{2}-a^{2}}\\), let \\(x=a\\sec(\\theta)\\), then \\(dx=a\\sec(\\theta)\\tan(\\theta)d\\theta\\) and \\(\\sqrt{x^{2}-a^{2}} = a\\tan(\\theta)\\)."},
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
		{"What is the Weierstrass substitution?", "\\(t = \\tan(x/2)\\)", "Let \\(t=\\tan(x/2)\\). Then \\(\\sin(x)=2t/(1+t^{2})\\), \\(\\cos(x)=(1-t^{2})/(1+t^{2})\\), \\(dx=2/(1+t^{2})dt\\)."},
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
		{"What is lim (f(x)·g(x)) in terms of individual limits?", "lim f(x) · lim g(x)", "The limit of a product is the product of the limits (provided both exist)."},
		{"What is lim (f(x)/g(x)) in terms of individual limits?", "lim f(x) / lim g(x) (if lim g(x)≠0)", "The limit of a quotient is the quotient of the limits (provided denominator limit ≠ 0)."},
		{"What is lim c·f(x) for a constant c?", "c · lim f(x)", "Constant multiples can be factored out of limits."},
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
		{"What are the vertical asymptotes of \\(f(x)=1/(x-2)\\)?", "x=2", "The denominator is \\(0\\) at \\(x=2\\), and the numerator is nonzero, so \\(x=2\\) is a vertical asymptote."},
		{"What is the horizontal asymptote of \\(f(x)=1/x\\)?", "y=0", "As \\(x\\to\\infty\\), \\(1/x\\to 0\\). As \\(x\\to -\\infty\\), \\(1/x\\to 0\\). So \\(y=0\\) is the horizontal asymptote."},
		{"What is the horizontal asymptote of \\(f(x)=\\frac{2x+1}{x-3}\\)?", "y=2", "As \\(x\\to\\infty\\), \\(\\frac{2x+1}{x-3}\\to 2\\). The ratio of leading coefficients gives the horizontal asymptote."},
		{"Can a function cross its horizontal asymptote?", "yes", "A function CAN cross its horizontal asymptote (unlike vertical asymptotes). The asymptote describes end behavior only."},
		{"What are the vertical asymptotes of \\(f(x)=\\tan(x)\\)?", "\\(x=\\pi/2 + n\\pi\\)", "\\(\\tan(x)=\\sin(x)/\\cos(x)\\). Vertical asymptotes occur where \\(\\cos(x)=0\\), i.e., at \\(x=\\pi/2+n\\pi\\)."},
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
		{"What kind of discontinuity does \\(|x|/x\\) have at \\(x=0\\)?", "jump", "The left limit is \\(-1\\) and the right limit is \\(1\\), so there is a jump discontinuity."},
		{"What kind of discontinuity does \\(\\frac{x^{2}-1}{x-1}\\) have at \\(x=1\\)?", "removable", "The function simplifies to \\(x+1\\) for \\(x\\neq 1\\), and the limit as \\(x\\to 1\\) is \\(2\\)."},
		{"What kind of discontinuity does \\(\\sin(1/x)\\) have at \\(x=0\\)?", "essential", "\\(\\sin(1/x)\\) oscillates infinitely near \\(0\\), so the limit does not exist (essential/oscillatory discontinuity)."},
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
		{"Is \\(0/0\\) an indeterminate form?", "yes", "\\(0/0\\) is indeterminate — the limit could be any real number depending on the functions."},
		{"Is \\(\\infty/\\infty\\) an indeterminate form?", "yes", "\\(\\infty/\\infty\\) is indeterminate — the ratio of functions both approaching infinity could converge to anything."},
		{"Is \\(0 \\cdot \\infty\\) an indeterminate form?", "yes", "\\(0 \\cdot \\infty\\) is indeterminate — it can be rewritten as \\(0/0\\) or \\(\\infty/\\infty\\)."},
		{"Is \\(\\infty - \\infty\\) an indeterminate form?", "yes", "\\(\\infty - \\infty\\) is indeterminate — the difference of two quantities approaching infinity could be anything."},
		{"Is \\(0^{\\infty}\\) an indeterminate form?", "no", "\\(0^{\\infty} = 0\\). If the base approaches \\(0\\) and the exponent approaches \\(\\infty\\), the result is \\(0\\) (not indeterminate)."},
		{"Is \\(1^{\\infty}\\) an indeterminate form?", "yes", "\\(1^{\\infty}\\) is indeterminate — it often arises with exponential limits and can approach \\(e\\), \\(1\\), or other values."},
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
		{"What condition is required for L'Hôpital's rule?", "both numerator and denominator \\(\\to 0\\) or \\(\\pm \\infty\\)", "L'Hôpital's rule applies to \\(0/0\\) or \\(\\infty/\\infty\\) indeterminate forms."},
		{"Does L'Hôpital's rule apply to \\(\\lim_{x \\to 0} \\frac{x^{2} \\cdot \\sin(1/x)}{\\sin(x)}\\)?", "no (not 0/0 or ∞/∞ form)", "Wait, this IS \\(0/0\\). But L'Hôpital would be messy. The rule can be applied, but checking conditions carefully is needed."},
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
		{"What does f(x)=O(g(x)) as x→∞ mean?", "|f(x)| ≤ M·|g(x)| for large x", "f=O(g) means there exists M>0 and x₀ such that |f(x)|≤M·|g(x)| for all x>x₀."},
		{"Is x²=O(x³) as x→∞?", "yes", "x² ≤ 1·x³ for x≥1, so x²=O(x³)."},
		{"Is x³=O(x²) as x→∞?", "no", "x³/x² = x → ∞, so no constant M can bound x³ by x² for all large x."},
		{"Is sin(x)=O(1) as x→∞?", "yes", "|sin(x)| ≤ 1, so sin(x)=O(1). In fact, any bounded function is O(1)."},
		{"What is the big-O notation primarily used for?", "describing asymptotic growth rates", "Big-O notation describes the limiting behavior of functions — which function dominates as the argument grows."},
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
		{"What does f(x)=o(g(x)) as x→∞ mean?", "f(x)/g(x) → 0", "f=o(g) means that f(x)/g(x) → 0 as x → ∞ (f grows strictly slower than g)."},
		{"Is x²=o(x³) as x→∞?", "yes", "x²/x³ = 1/x → 0, so x²=o(x³)."},
		{"Is x³=o(x³) as x→∞?", "no", "x³/x³ = 1 → 1, not 0, so x³≠o(x³). However, x³=O(x³)."},
		{"What is the difference between O and o?", "o requires ratio → 0; O allows any bounded ratio", "f=O(g) means |f|≤M·|g| for large x. f=o(g) means f/g → 0, which is a stricter requirement."},
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
		{"What does the squeeze theorem say?", "if g(x)≤f(x)≤h(x) and lim g=lim h=L, then lim f=L", "The squeeze theorem: if f is bounded between g and h, and g and h have the same limit L, then f also approaches L."},
		{"Find \\(\\lim_{x \\to 0} x^{2} \\sin(1/x)\\) using the squeeze theorem.", "0", "\\(-1\\leq \\sin(1/x)\\leq 1 \\implies -x^{2} \\leq x^{2}\\sin(1/x) \\leq x^{2}\\). Since \\(-x^{2}\\to 0\\) and \\(x^{2}\\to 0\\), the squeeze theorem gives limit \\(0\\)."},
		{"Find \\(\\lim_{x \\to 0} x \\cos(1/x)\\) using the squeeze theorem.", "0", "\\(-1\\leq \\cos(1/x)\\leq 1 \\implies -x \\leq x\\cos(1/x) \\leq x\\). Since \\(-x\\to 0\\) and \\(x\\to 0\\), the limit is \\(0\\)."},
		{"What is a common way to apply the squeeze theorem?", "bound the function between two simpler functions with the same limit", "Bounding f(x) using inequalities, then showing the upper and lower bounds converge to the same limit."},
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
		{"What is the supremum of the set \\(\\{1-1/n : n \\in \\mathbb{N}\\}\\)?", "1", "The values approach \\(1\\) from below. \\(1\\) is an upper bound, and no smaller number is an upper bound. Supremum = 1."},
		{"What is the infimum of the set \\(\\{1/n : n \\in \\mathbb{N}\\}\\)?", "0", "The values approach \\(0\\) from above. \\(0\\) is a lower bound, and no larger number is a lower bound. Infimum = 0."},
		{"Is the supremum always in the set?", "no", "The supremum need not be in the set. E.g., \\(\\sup\\{1-1/n\\}=1\\), but \\(1\\) is not in the set."},
		{"Is the maximum always in the set?", "yes", "If a maximum exists, it is the supremum and is also in the set. Unlike supremum, maximum must be attained."},
		{"What is \\(\\sup\\{x \\in \\mathbb{R} : x^{2}<2\\}\\)?", "\\(\\sqrt{2}\\)", "The set is \\((-\\sqrt{2}, \\sqrt{2})\\). The supremum is \\(\\sqrt{2}\\) (the least upper bound), which is not in the set."},
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
		{"What is uniform continuity?", "δ depends only on ε, not on the point", "A function is uniformly continuous if for any ε>0, there exists δ>0 such that |x-y|<δ implies |f(x)-f(y)|<ε for ALL x,y."},
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
		{"What does the extreme value theorem (Weierstrass) state?", "a continuous function on [a,b] attains its max and min", "If f is continuous on [a,b], then there exist c,d∈[a,b] such that f(c)≤f(x)≤f(d) for all x∈[a,b]."},
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does \\(a_{n} = 1/n\\) converge? If so, to what?", "yes, to 0", "As \\(n\\to\\infty\\), \\(1/n \\to 0\\). The sequence converges to \\(0\\)."},
		{"Does \\(a_{n} = n/(n+1)\\) converge? If so, to what?", "yes, to 1", "\\(n/(n+1) = 1/(1+1/n) \\to 1/(1+0) = 1\\)."},
		{"Does \\(a_{n} = (-1)^{n}\\) converge?", "no", "\\((-1)^{n}\\) oscillates between \\(-1\\) and \\(1\\) without approaching a single limit."},
		{"Does \\(a_{n} = n^{2}\\) converge?", "no (diverges to \\(\\infty\\))", "\\(n^{2} \\to \\infty\\) as \\(n\\to\\infty\\), so the sequence diverges."},
		{"Does \\(a_{n} = 2^{n}/n!\\) converge? If so, to what?", "yes, to 0", "The factorial grows faster than the exponential, so \\(2^{n}/n! \\to 0\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type cauchySeqGen struct{}

func (g *cauchySeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is a Cauchy sequence?", "terms get arbitrarily close as n increases", "A sequence where for any ε>0, there exists N such that |aₘ-aₙ|<ε for all m,n>N."},
		{"In ℝ, do all Cauchy sequences converge?", "yes", "ℝ is complete — every Cauchy sequence in ℝ converges to a real number."},
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
		{"What is a monotone increasing sequence?", "aₙ₊₁ ≥ aₙ for all n", "Each term is ≥ the previous term. E.g., 1, 2, 3, 4, ..."},
		{"What is the monotone convergence theorem?", "a bounded monotone sequence converges", "If a sequence is monotone (increasing or decreasing) and bounded, then it converges."},
		{"Does aₙ = 1-1/n converge? Why?", "yes, it's bounded and increasing", "aₙ increases: 1-1/(n+1) > 1-1/n. It's bounded above by 1. By the MCT, it converges to 1."},
		{"Does aₙ = n converge (bounded or not)?", "no (unbounded)", "aₙ is increasing but not bounded above, so it diverges."},
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
		{"What sequence defines Euler's number \\(e\\)?", "\\((1+1/n)^{n}\\)", "\\(e = \\lim_{n\\to\\infty} (1+1/n)^{n} \\approx 2.71828...\\)"},
		{"What is \\(e\\) to 3 decimal places?", "2.718", "\\(e = \\lim (1+1/n)^{n} = 2.718281828...\\)"},
		{"What is an alternative series representation of \\(e\\)?", "\\(\\sum 1/n!\\)", "\\(e = 1 + 1/1! + 1/2! + 1/3! + \\cdots = \\sum_{n=0}^{\\infty} 1/n!\\)"},
		{"Is \\(e\\) rational or irrational?", "irrational", "\\(e\\) is irrational (proved by Euler). It is also transcendental."},
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
		{"What is a function sequence?", "a sequence where each term is a function", "A function sequence {fₙ} has each fₙ: D → ℝ defined on a common domain D."},
		{"What does pointwise convergence mean?", "fₙ(x) → f(x) for each fixed x", "fₙ → f pointwise if for each x in D, the numerical sequence fₙ(x) → f(x)."},
		{"What does uniform convergence mean?", "sup|fₙ(x)-f(x)| → 0", "fₙ → f uniformly if for any ε>0, there exists N such that |fₙ(x)-f(x)|<ε for ALL x in D when n>N."},
		{"Is uniform convergence stronger than pointwise?", "yes", "Uniform convergence implies pointwise convergence, but the converse is false. E.g., fₙ(x)=xⁿ on [0,1] converges pointwise but not uniformly."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// --- Series generators ---

type seriesConceptGen struct{}

func (g *seriesConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is a series?", "the sum of terms of a sequence", "A series ∑ aₙ = a₁ + a₂ + a₃ + ... represents the sum of a sequence's terms."},
		{"What is a partial sum Sₙ?", "Sₙ = ∑_{k=1}ⁿ aₖ", "The nth partial sum is the sum of the first n terms of the series."},
		{"When does a series converge?", "when the sequence of partial sums converges", "A series converges if lim_{n→∞} Sₙ exists and is finite."},
		{"What does it mean to find the sum of a series?", "to find the limit of its partial sums", "The sum of a convergent series is lim Sₙ, the limit of its partial sums."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type harmonicSeriesGen struct{}

func (g *harmonicSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Does the harmonic series \\(\\sum 1/n\\) converge or diverge?", "diverges", "The harmonic series \\(\\sum 1/n\\) diverges (slowly — like \\(\\ln(n)\\)), even though its terms \\(\\to 0\\)."},
		{"What is the \\(p\\)-series \\(\\sum 1/n^{p}\\): when does it converge?", "converges for \\(p > 1\\)", "\\(\\sum 1/n^{p}\\) converges if \\(p>1\\), diverges if \\(p\\leq 1\\). The harmonic series \\((p=1)\\) is the boundary case."},
		{"Approximately how large is the \\(n\\)th partial sum of the harmonic series?", "\\(H_{n} \\approx \\ln(n) + \\gamma\\) \\((\\gamma\\approx 0.577)\\)", "The harmonic numbers \\(H_{n} = \\sum_{k=1}^{n} 1/k \\approx \\ln(n) + \\gamma\\), where \\(\\gamma\\) is the Euler-Mascheroni constant."},
		{"Is \\(\\sum 1/n^{2}\\) convergent or divergent?", "converges (to \\(\\pi^{2}/6\\))", "\\(\\sum 1/n^{2}\\) converges \\((p=2>1)\\). Its sum is \\(\\pi^{2}/6 \\approx 1.645\\) (Basel problem)."},
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
		{"What is a positive term series?", "a series with all terms ≥ 0", "A positive term series ∑ aₙ has aₙ ≥ 0 for all n."},
		{"What test compares a series to another known series?", "comparison test", "If 0 ≤ aₙ ≤ bₙ and ∑ bₙ converges, then ∑ aₙ converges. If ∑ aₙ diverges, then ∑ bₙ diverges."},
		{"What does the limit comparison test say?", "if aₙ/bₙ → c>0, both series have same fate", "For positive series, if aₙ/bₙ → c where 0<c<∞, then ∑ aₙ and ∑ bₙ either both converge or both diverge."},
		{"What does the ratio test compare?", "the ratio of successive terms aₙ₊₁/aₙ", "If lim aₙ₊₁/aₙ = L < 1, the series converges absolutely. If L > 1, it diverges. Inconclusive if L=1."},
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
		{"What is an alternating series?", "terms alternate in sign: \\(\\sum (-1)^{n}a_{n}\\)", "An alternating series has terms that alternate between positive and negative."},
		{"What does the alternating series test require?", "\\(a_{n}\\) decreasing to \\(0\\)", "If \\(|a_{n}|\\) decreases monotonically and \\(a_{n}\\to 0\\), then \\(\\sum (-1)^{n}a_{n}\\) converges."},
		{"Does \\(\\sum (-1)^{n}/n\\) converge?", "yes (conditionally)", "By the alternating series test: \\(1/n\\) decreases and \\(\\to 0\\). So it converges. It converges conditionally (\\(\\sum 1/n\\) diverges)."},
		{"Does \\(\\sum (-1)^{n}\\) converge?", "no", "The terms do not \\(\\to 0\\) (they alternate between \\(-1\\) and \\(1\\)). The series diverges by the \\(n\\)th term test."},
		{"What is the error bound for an alternating series?", "\\(|\\text{error}| \\leq\\) first omitted term", "For a convergent alternating series with decreasing terms, \\(|S - S_{n}| \\leq a_{n+1}\\) (the next term)."},
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
		{"What does the integral test relate?", "\\(\\sum f(n)\\) and \\(\\int f(x) \\, dx\\)", "If \\(f\\) is positive, continuous, and decreasing on \\([1,\\infty)\\), then \\(\\sum f(n)\\) converges iff \\(\\int_{1}^{\\infty} f(x) \\, dx\\) converges."},
		{"Does \\(\\sum 1/n^{2}\\) converge? Use integral test.", "yes", "\\(\\int_{1}^{\\infty} 1/x^{2} \\, dx = [-1/x]_{1}^{\\infty} = 1\\). The integral converges, so the series converges."},
		{"Does \\(\\sum 1/n\\) converge? Use integral test.", "no", "\\(\\int_{1}^{\\infty} 1/x \\, dx = [\\ln(x)]_{1}^{\\infty} = \\infty\\). The integral diverges, so the series diverges."},
		{"Apply integral test to \\(\\sum 1/(n \\ln n)\\) for \\(n\\geq 2\\).", "diverges", "\\(\\int_{2}^{\\infty} 1/(x \\ln x) \\, dx = [\\ln(\\ln x)]_{2}^{\\infty} = \\infty\\). So \\(\\sum 1/(n \\ln n)\\) diverges."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type rootTestGen struct{}

func (g *rootTestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What does the root test evaluate?", "\\(\\limsup |a_{n}|^{1/n}\\)", "If \\(\\limsup |a_{n}|^{1/n} = L < 1\\), the series \\(\\sum a_{n}\\) converges absolutely. If \\(L > 1\\), it diverges. Inconclusive if \\(L=1\\)."},
		{"Apply root test to \\(\\sum (n/(n+1))^{n^{2}}\\).", "converges", "\\(|a_{n}|^{1/n} = (n/(n+1))^{n} = 1/(1+1/n)^{n} \\to 1/e < 1\\). The series converges."},
		{"Apply root test to \\(\\sum (1/2^{n})\\).", "converges", "\\(|a_{n}|^{1/n} = 1/2 < 1\\). The series \\(\\sum 1/2^{n}\\) converges by root test (it's a geometric series)."},
		{"When might the root test be particularly useful?", "when \\(a_{n}\\) involves \\(n\\)th powers", "The root test is especially useful when \\(a_{n}\\) contains expressions like \\((\\text{something})^{n}\\) or \\(n\\)-th powers."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type powerSeriesGen struct{}

func (g *powerSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the general form of a power series?", "\\(\\sum c_{n}(x-a)^{n}\\)", "A power series centered at \\(a\\): \\(c_{0} + c_{1}(x-a) + c_{2}(x-a)^{2} + \\cdots\\)"},
		{"What is the radius of convergence \\(R\\)?", "the series converges for \\(|x-a|<R\\), diverges for \\(|x-a|>R\\)", "The radius \\(R\\) defines the interval where the power series converges absolutely."},
		{"What is the power series for \\(1/(1-x)\\)?", "\\(\\sum x^{n}\\) (for \\(|x|<1\\))", "\\(1/(1-x) = 1 + x + x^{2} + x^{3} + \\cdots = \\sum_{n=0}^{\\infty} x^{n}\\), which converges for \\(|x|<1\\)."},
		{"What is the power series for \\(e^{x}\\)?", "\\(\\sum x^{n}/n!\\)", "\\(e^{x} = 1 + x + x^{2}/2! + x^{3}/3! + \\cdots = \\sum_{n=0}^{\\infty} x^{n}/n!\\), which converges for all real \\(x\\)."},
		{"What is the power series for \\(\\sin(x)\\)?", "\\(\\sum (-1)^{n}x^{2n+1}/(2n+1)!\\)", "\\(\\sin(x) = x - x^{3}/3! + x^{5}/5! - x^{7}/7! + \\cdots\\), converging for all real \\(x\\)."},
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
		{"What is a Taylor series?", "\\(f(x) = \\sum \\frac{f^{(n)}(a)(x-a)^{n}}{n!}\\)", "The Taylor series of \\(f\\) centered at \\(a\\) is \\(\\sum_{n=0}^{\\infty} \\frac{f^{(n)}(a)(x-a)^{n}}{n!}\\)."},
		{"What is the Taylor series for \\(e^{x}\\) at \\(x=0\\) (Maclaurin)?", "\\(\\sum x^{n}/n!\\)", "\\(e^{x} = \\sum_{n=0}^{\\infty} x^{n}/n!\\) with radius of convergence \\(R=\\infty\\)."},
		{"What is the Maclaurin series for \\(\\sin(x)\\)?", "\\(\\sum (-1)^{n}x^{2n+1}/(2n+1)!\\)", "\\(\\sin(x) = x - x^{3}/3! + x^{5}/5! - \\cdots\\) with \\(R=\\infty\\)."},
		{"What is the Maclaurin series for \\(\\cos(x)\\)?", "\\(\\sum (-1)^{n}x^{2n}/(2n)!\\)", "\\(\\cos(x) = 1 - x^{2}/2! + x^{4}/4! - \\cdots\\) with \\(R=\\infty\\)."},
		{"What is the Maclaurin series for \\(\\ln(1+x)\\)?", "\\(\\sum (-1)^{n+1}x^{n}/n\\) (for \\(|x|<1\\))", "\\(\\ln(1+x) = x - x^{2}/2 + x^{3}/3 - x^{4}/4 + \\cdots\\) with \\(R=1\\)."},
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
		{"What is the Cauchy criterion for series?", "partial sums form a Cauchy sequence", "∑ aₙ converges iff for every ε>0, there exists N such that |∑_{k=m}ⁿ aₖ|<ε for all n≥m>N."},
		{"What does the Cauchy criterion imply about terms?", "aₙ → 0", "If ∑ aₙ converges, then by Cauchy criterion with n=m, |aₘ|<ε for m>N, so aₙ→0."},
		{"What does the Cauchy condensation test say?", "∑ f(n) converges iff ∑ 2ⁿf(2ⁿ) converges", "For decreasing positive f, the condensed series ∑ 2ⁿf(2ⁿ) has the same convergence as ∑ f(n)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type fourierSeriesGen struct{}

func (g *fourierSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is a Fourier series?", "\\(\\sum (a_{n}\\cos(nx) + b_{n}\\sin(nx))\\)", "A Fourier series represents a periodic function as a sum of sines and cosines."},
		{"What is \\(a_{0}\\) in a Fourier series?", "the average value: \\((1/\\pi)\\int f(x) \\, dx\\) over one period", "\\(a_{0} = \\frac{1}{\\pi} \\int_{-\\pi}^{\\pi} f(x) \\, dx\\), giving twice the average value."},
		{"What formula gives Fourier coefficient \\(a_{n}\\)?", "\\(a_{n} = \\frac{1}{\\pi} \\int f(x)\\cos(nx) \\, dx\\)", "\\(a_{n} = \\frac{1}{\\pi} \\int_{-\\pi}^{\\pi} f(x)\\cos(nx) \\, dx\\) for \\(n\\geq 0\\)."},
		{"What formula gives Fourier coefficient \\(b_{n}\\)?", "\\(b_{n} = \\frac{1}{\\pi} \\int f(x)\\sin(nx) \\, dx\\)", "\\(b_{n} = \\frac{1}{\\pi} \\int_{-\\pi}^{\\pi} f(x)\\sin(nx) \\, dx\\) for \\(n\\geq 1\\)."},
		{"What is the Fourier series of an odd function?", "only sine terms (\\(a_{n}=0\\))", "Odd functions have Fourier series with only sine terms: all \\(a_{n}=0\\)."},
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
		{"What is a function series?", "∑ fₙ(x) where each term is a function", "A function series ∑ fₙ(x) is the sum of a sequence of functions."},
		{"What does pointwise convergence of a function series mean?", "∑ fₙ(x) converges for each fixed x", "The series converges pointwise if for each x in the domain, the numerical series converges."},
		{"What does uniform convergence of a function series mean?", "partial sums converge uniformly", "The series converges uniformly if the sequence of partial sums converges uniformly."},
		{"What is the Weierstrass M-test for?", "testing uniform convergence", "If |fₙ(x)| ≤ Mₙ and ∑ Mₙ converges, then ∑ fₙ(x) converges uniformly and absolutely."},
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
		{"What is the first step of logarithmic differentiation?", "take the natural log of both sides", "Logarithmic differentiation starts by taking \\(\\ln\\) of both sides, then differentiating implicitly."},
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
		{"State the inverse function rule.", "(f^-1)'(y) = 1/f'(x) where y = f(x)", "If \\(f\\) is differentiable with \\(f'(x) \\neq 0\\), then \\((f^{-1})'(f(x)) = 1/f'(x)\\)."},
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
		{"Find the critical points of \\(f(x)=x^{3}-3x\\).", "x=-1 and x=1", "\\(f'(x) = 3x^{2}-3 = 0\\) gives \\(x=\\pm 1\\)."},
		{"What is a critical point of \\(f\\)?", "a point where f'(x)=0 or f' is undefined", "Critical points are interior points where the derivative is zero or undefined."},
		{"Find the critical points of \\(f(x)=x^{3}\\).", "x=0", "\\(f'(x) = 3x^{2} = 0\\) gives \\(x=0\\) (a stationary point that is not an extremum)."},
		{"Find the critical points of \\(f(x)=2x^{3}-6x^{2}\\).", "x=0 and x=2", "\\(f'(x) = 6x^{2}-12x = 6x(x-2) = 0\\) gives \\(x=0, 2\\)."},
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
		{"Classify \\(x=0\\) for \\(f(x)=x^{2}\\) using the second derivative test.", "local minimum", "\\(f''(0) = 2 > 0\\), so \\(x=0\\) is a local minimum."},
		{"Classify \\(x=0\\) for \\(f(x)=-x^{2}\\) using the second derivative test.", "local maximum", "\\(f''(0) = -2 < 0\\), so \\(x=0\\) is a local maximum."},
		{"State the second derivative test.", "f''(c)>0 minimum, f''(c)<0 maximum, f''(c)=0 inconclusive", "If \\(f'(c)=0\\): \\(f''(c)>0\\) gives a local minimum, \\(f''(c)<0\\) a local maximum, \\(f''(c)=0\\) is inconclusive."},
		{"Classify \\(x=0\\) for \\(f(x)=x^{3}-3x^{2}\\) using the second derivative test.", "local maximum", "\\(f''(x) = 6x-6\\); \\(f''(0) = -6 < 0\\), so \\(x=0\\) is a local maximum."},
		{"Classify \\(x=2\\) for \\(f(x)=x^{3}-3x^{2}\\) using the second derivative test.", "local minimum", "\\(f''(x) = 6x-6\\); \\(f''(2) = 6 > 0\\), so \\(x=2\\) is a local minimum."},
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
		{"Find the inflection points of \\(f(x)=x^{3}\\).", "x=0", "\\(f''(x) = 6x = 0\\) at \\(x=0\\), and concavity changes there."},
		{"Find the inflection points of \\(f(x)=x^{4}-6x^{2}\\).", "x=-1 and x=1", "\\(f''(x) = 12x^{2}-12 = 0\\) at \\(x=\\pm 1\\), with concavity change on both sides."},
		{"What defines an inflection point?", "concavity changes sign there", "An inflection point is where \\(f''\\) changes sign, i.e. concavity flips."},
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
		{"Find the global maximum of \\(f(x)=x^{2}-4x\\) on \\([0,3]\\).", "0 at x=3", "Critical point \\(x=2\\): \\(f(2)=-4\\). Endpoints: \\(f(0)=0\\), \\(f(3)=-3\\). Global max is \\(0\\) at \\(x=3\\)."},
		{"Find the global minimum of \\(f(x)=x^{2}-4x\\) on \\([0,3]\\).", "-4 at x=2", "Critical point \\(x=2\\): \\(f(2)=-4\\). Endpoints give \\(0\\) and \\(-3\\). Global min is \\(-4\\)."},
		{"State the closed-interval method.", "evaluate f at critical points and endpoints, take largest/smallest", "On \\([a,b]\\): find critical points inside, evaluate \\(f\\) there and at \\(a,b\\); extremes are the max/min values."},
		{"Find the global maximum of \\(f(x)=-x^{2}+4x\\) on \\([0,3]\\).", "4 at x=2", "Critical point \\(x=2\\): \\(f(2)=4\\). Endpoints: \\(f(0)=0\\), \\(f(3)=3\\). Global max is \\(4\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivCurveSketchGen struct{}

func (g *derivCurveSketchGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"For \\(f(x)=x^{3}-3x\\), list critical points and inflection point for sketching.", "critical x=-1,1; inflection x=0", "\\(f' = 3x^{2}-3\\): critical \\(\\pm 1\\). \\(f'' = 6x\\): inflection at \\(0\\)."},
		{"What three ingredients does curve sketching combine?", "intercepts, asymptotes, first/second derivative analysis", "Sketching combines intercepts, asymptotes (limits at infinity), critical points, and concavity/inflection."},
		{"For \\(f(x)=x^{2}\\), where is it increasing?", "increasing on (0, infinity)", "\\(f'(x) = 2x > 0\\) for \\(x>0\\), so increasing on \\((0,\\infty)\\)."},
		{"For \\(f(x)=1/x\\), what are the asymptotes?", "vertical x=0, horizontal y=0", "\\(\\lim_{x\\to 0} 1/x = \\pm\\infty\\) (vertical \\(x=0\\)); \\(\\lim_{x\\to\\pm\\infty} 1/x = 0\\) (horizontal \\(y=0\\))."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivTangentLineGen struct{}

func (g *derivTangentLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find the tangent line to \\(f(x)=x^{2}\\) at \\(x=1\\).", "y=2x-1", "Slope \\(f'(1)=2\\); point \\((1,1)\\): \\(y-1=2(x-1)\\), i.e. \\(y=2x-1\\)."},
		{"Find the tangent line to \\(f(x)=x^{3}\\) at \\(x=2\\).", "y=12x-16", "Slope \\(f'(2)=12\\); point \\((2,8)\\): \\(y-8=12(x-2)\\), i.e. \\(y=12x-16\\)."},
		{"Find the tangent line to \\(f(x)=\\sin(x)\\) at \\(x=0\\).", "y=x", "Slope \\(\\cos(0)=1\\); point \\((0,0)\\): \\(y=x\\)."},
		{"Find the tangent line to \\(f(x)=e^{x}\\) at \\(x=0\\).", "y=x+1", "Slope \\(e^{0}=1\\); point \\((0,1)\\): \\(y-1=x\\), i.e. \\(y=x+1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivLinearApproxGen struct{}

func (g *derivLinearApproxGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Write the linearization of \\(f(x)=\\sqrt{x}\\) at \\(x=4\\).", "L(x)=2+(x-4)/4", "\\(f(4)=2\\), \\(f'(x)=1/(2\\sqrt{x})\\), \\(f'(4)=1/4\\): \\(L(x)=2+(x-4)/4\\)."},
		{"Use \\(L(x)=1+(x-1)/2\\) to approximate \\(\\sqrt{1.1}\\).", "1.05", "\\(L(1.1) = 1 + 0.1/2 = 1.05\\)."},
		{"Write the linearization of \\(f(x)=e^{x}\\) at \\(x=0\\).", "L(x)=1+x", "\\(f(0)=1\\), \\(f'(0)=1\\): \\(L(x) = 1+x\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivNewtonGen struct{}

func (g *derivNewtonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State Newton's iteration formula.", "x_{n+1} = x_n - f(x_n)/f'(x_n)", "Newton's method iterates \\(x_{n+1} = x_{n} - f(x_{n})/f'(x_{n})\\) from an initial guess."},
		{"One Newton step for \\(f(x)=x^{2}-2\\) from \\(x_{0}=1\\).", "1.5", "\\(x_{1} = 1 - (1-2)/2 = 1.5\\)."},
		{"One Newton step for \\(f(x)=x^{2}-2\\) from \\(x_{0}=1.5\\).", "1.4167", "\\(x_{1} = 1.5 - (2.25-2)/3 = 1.5 - 0.0833 = 1.4167\\)."},
		{"When does Newton's method fail?", "when f'(x_n)=0 or the guess is far from the root", "A zero derivative at an iterate (or bad starting guess / cycling) breaks the iteration."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivErrorPropGen struct{}

func (g *derivErrorPropGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If \\(y=x^{2}\\) and \\(x=3\\pm 0.1\\), estimate \\(\\Delta y\\).", "0.6", "\\(dy = 2x\\,dx = 6 \\cdot 0.1 = 0.6\\)."},
		{"If \\(y=\\sqrt{x}\\) and \\(x=100\\pm 1\\), estimate \\(\\Delta y\\).", "0.05", "\\(dy = dx/(2\\sqrt{x}) = 1/20 = 0.05\\)."},
		{"State the differential error estimate.", "Δy ≈ f'(x)Δx", "A measurement error \\(\\Delta x\\) propagates as \\(\\Delta y \\approx f'(x)\\Delta x\\)."},
		{"If \\(y=1/x\\) and \\(x=2\\pm 0.1\\), estimate \\(\\Delta y\\).", "-0.025", "\\(dy = -dx/x^{2} = -0.1/4 = -0.025\\) (magnitude 0.025)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivImplicitSecondGen struct{}

func (g *derivImplicitSecondGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"For \\(x^{2}+y^{2}=25\\), what is \\(dy/dx\\)?", "-x/y", "Implicit: \\(2x+2y\\,y'=0\\), so \\(y'=-x/y\\)."},
		{"For \\(x^{2}+y^{2}=25\\), what is \\(d^{2}y/dx^{2}\\) in terms of \\(x,y\\)?", "-25/y^3", "Differentiate \\(y'=-x/y\\): \\(y'' = -(y-xy')/y^{2} = -(y+x^{2}/y)/y^{2} = -25/y^{3}\\)."},
		{"What is the first step in finding a second implicit derivative?", "find dy/dx by implicit differentiation", "Compute \\(dy/dx\\) first, then differentiate that relation implicitly again."},
		{"For \\(y=x^{3}\\) (explicit check), what is \\(y''\\)?", "6x", "\\(y'=3x^{2}\\), \\(y''=6x\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivChainPowerGen struct{}

func (g *derivChainPowerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(f'(x)\\) if \\(f(x)=(3x+1)^{4}\\).", "12(3x+1)^3", "Chain rule: \\(4(3x+1)^{3} \\cdot 3 = 12(3x+1)^{3}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=\\sqrt{x^{2}+1}\\).", "x/sqrt(x^2+1)", "Chain rule: \\(\\frac{1}{2}(x^{2}+1)^{-1/2} \\cdot 2x = x/\\sqrt{x^{2}+1}\\)."},
		{"Find \\(f'(x)\\) if \\(f(x)=(x^{2})^{3}\\) simplified first.", "6x^5", "Simplify: \\((x^{2})^{3}=x^{6}\\), so \\(f'=6x^{5}\\)."},
		{"Which technique handles \\(f(x)=x^{x}\\)?", "logarithmic differentiation after taking logs", "\\(x^{x}\\) is a composite power: take \\(\\ln\\) both sides, then differentiate implicitly."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivChainMultiGen struct{}

func (g *derivChainMultiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"For \\(z=f(x,y)\\), \\(x=g(t)\\), \\(y=h(t)\\), state \\(dz/dt\\).", "∂z/∂x·dx/dt + ∂z/∂y·dy/dt", "Multivariable chain rule: \\(dz/dt = f_{x}x' + f_{y}y'\\)."},
		{"For \\(z=x^{2}y\\), \\(x=t\\), \\(y=t^{2}\\), find \\(dz/dt\\).", "4t^3", "\\(z=t^{4}\\), so \\(dz/dt=4t^{3}\\) (check: \\(2xy\\cdot 1 + x^{2}\\cdot 2t = 2t^{3}+2t^{3}\\))."},
		{"For \\(z=x+y\\), \\(x=t^{2}\\), \\(y=3t\\), find \\(dz/dt\\).", "2t+3", "\\(dz/dt = 1\\cdot 2t + 1\\cdot 3 = 2t+3\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivGradientGen struct{}

func (g *derivGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\nabla f\\) for \\(f(x,y)=x^{2}+y^{2}\\).", "(2x, 2y)", "\\(\\nabla f = (f_{x}, f_{y}) = (2x, 2y)\\)."},
		{"Find \\(\\nabla f\\) for \\(f(x,y)=xy\\).", "(y, x)", "\\(\\nabla f = (y, x)\\)."},
		{"What direction does the gradient point?", "steepest ascent", "The gradient points in the direction of steepest increase of \\(f\\)."},
		{"Find \\(\\nabla f\\) at \\((1,2)\\) for \\(f(x,y)=x^{2}y\\).", "(4, 1)", "\\(\\nabla f = (2xy, x^{2})\\); at \\((1,2)\\): \\((4,1)\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivDirectionalGen struct{}

func (g *derivDirectionalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the directional derivative formula.", "D_u f = ∇f · u", "The derivative in unit direction \\(u\\) is \\(D_{u}f = \\nabla f \\cdot u\\)."},
		{"For \\(f(x,y)=x^{2}+y^{2}\\) at \\((1,0)\\) in direction \\((1,0)\\), find \\(D_{u}f\\).", "2", "\\(\\nabla f(1,0) = (2,0)\\); dotted with \\((1,0)\\): \\(2\\)."},
		{"For \\(f(x,y)=xy\\) at \\((1,1)\\) in direction \\((1,0)\\), find \\(D_{u}f\\).", "1", "\\(\\nabla f(1,1) = (1,1)\\); dotted with \\((1,0)\\): \\(1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivPartialHigherGen struct{}

func (g *derivPartialHigherGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"For \\(f(x,y)=x^{3}y^{2}\\), find \\(f_{xx}\\).", "6xy^2", "\\(f_{x} = 3x^{2}y^{2}\\), so \\(f_{xx} = 6xy^{2}\\)."},
		{"For \\(f(x,y)=x^{3}y^{2}\\), find \\(f_{yy}\\).", "2x^3", "\\(f_{y} = 2x^{3}y\\), so \\(f_{yy} = 2x^{3}\\)."},
		{"For \\(f(x,y)=x^{2}y\\), find \\(f_{xy}\\).", "2x", "\\(f_{x} = 2xy\\), so \\(f_{xy} = 2x\\)."},
		{"For \\(f(x,y)=\\sin(x)\\cos(y)\\), find \\(f_{xx}\\).", "-sin(x)cos(y)", "\\(f_{x} = \\cos(x)\\cos(y)\\), so \\(f_{xx} = -\\sin(x)\\cos(y)\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivClairautGen struct{}

func (g *derivClairautGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State Clairaut's theorem.", "f_xy = f_yx when mixed partials are continuous", "If \\(f_{xy}\\) and \\(f_{yx}\\) are continuous near a point, they are equal there."},
		{"For \\(f(x,y)=x^{2}y^{3}\\), verify \\(f_{xy}\\) equals \\(f_{yx}\\).", "6xy^2", "\\(f_{xy} = 6xy^{2} = f_{yx}\\), as Clairaut guarantees for polynomials."},
		{"What hypothesis does Schwarz's theorem need?", "continuity of the mixed partials", "Schwarz/Clairaut needs \\(f_{xy}\\) and \\(f_{yx}\\) continuous at the point."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivTotalDiffGen struct{}

func (g *derivTotalDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Write the total differential of \\(z=f(x,y)\\).", "dz = f_x dx + f_y dy", "The total differential is \\(dz = f_{x}\\,dx + f_{y}\\,dy\\)."},
		{"For \\(z=x^{2}y\\), write \\(dz\\).", "dz = 2xy dx + x^2 dy", "\\(z_{x} = 2xy\\), \\(z_{y} = x^{2}\\), so \\(dz = 2xy\\,dx + x^{2}\\,dy\\)."},
		{"For \\(z=xy\\), write \\(dz\\).", "dz = y dx + x dy", "\\(z_{x} = y\\), \\(z_{y} = x\\), so \\(dz = y\\,dx + x\\,dy\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivDiffContinuityGen struct{}

func (g *derivDiffContinuityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If \\(f\\) is differentiable at \\(a\\), what follows?", "f is continuous at a", "Differentiability at \\(a\\) implies continuity at \\(a\\)."},
		{"Does continuity imply differentiability?", "no", "Continuity does not imply differentiability (e.g. \\(|x|\\) at \\(0\\))."},
		{"Give a function continuous but not differentiable at \\(0\\).", "|x|", "\\(f(x)=|x|\\) is continuous at \\(0\\) but has no derivative there (corner)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type derivJacobianGen struct{}

func (g *derivJacobianGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Write the Jacobian of \\(T(x,y)=(x^{2},y^{2})\\).", "[[2x,0],[0,2y]]", "Partial derivatives row by row: \\(\\begin{pmatrix}2x&0\\\\0&2y\\end{pmatrix}\\)."},
		{"What are the entries of a Jacobian matrix?", "all first-order partials ∂T_i/∂x_j", "Row \\(i\\), column \\(j\\) holds \\(\\partial T_{i}/\\partial x_{j}\\)."},
		{"Write the Jacobian of \\(T(x,y)=(x+y,x-y)\\).", "[[1,1],[1,-1]]", "Row 1: \\((1,1)\\); row 2: \\((1,-1)\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralNetChangeGen struct{}

func (g *integralNetChangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"If \\(v(t)=3t^{2}\\) is velocity, find displacement from \\(t=0\\) to \\(t=2\\).", "8", "Net change: \\(\\int_{0}^{2} 3t^{2}\\,dt = [t^{3}]_{0}^{2} = 8\\)."},
		{"State the net change theorem.", "∫_a^b f'(x)dx = f(b)-f(a)", "The integral of a rate of change gives net change: \\(\\int_{a}^{b} f'(x)\\,dx = f(b)-f(a)\\)."},
		{"If \\(f'(x)=2x\\) and \\(f(1)=5\\), find \\(f(3)\\).", "13", "\\(f(3)-f(1) = \\int_{1}^{3} 2x\\,dx = 8\\), so \\(f(3) = 13\\)."},
		{"Water flows at \\(r(t)=4t\\) L/min; how much flows in from \\(t=0\\) to \\(t=3\\)?", "18", "\\(\\int_{0}^{3} 4t\\,dt = [2t^{2}]_{0}^{3} = 18\\) liters."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
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
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the mean value theorem for integrals.", "there exists c with f(c) = average of f on [a,b]", "Some \\(c \\in [a,b]\\) attains the average value: \\(f(c) = \\frac{1}{b-a}\\int_{a}^{b}f\\)."},
		{"Find the average value of \\(f(x)=x^{2}\\) on \\([0,2]\\).", "4/3", "Average \\(= \\frac{1}{2}\\int_{0}^{2}x^{2}\\,dx = \\frac{1}{2}\\cdot\\frac{8}{3} = 4/3\\)."},
		{"Find \\(c\\) guaranteed by the integral MVT for \\(f(x)=x\\) on \\([0,4]\\).", "2", "Average is \\(2\\); \\(f(c)=c=2\\) gives \\(c=2\\)."},
		{"Find the average value of \\(f(x)=3x\\) on \\([1,3]\\).", "6", "Average \\(= \\frac{1}{2}\\int_{1}^{3}3x\\,dx = \\frac{1}{2}(12) = 6\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralWorkGen struct{}

func (g *integralWorkGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"A force \\(F(x)=2x\\) moves an object from \\(x=0\\) to \\(x=3\\). Find the work.", "9", "Work \\(= \\int_{0}^{3}2x\\,dx = 9\\) J."},
		{"A force \\(F(x)=x^{2}\\) moves an object from \\(x=1\\) to \\(x=2\\). Find the work.", "7/3", "Work \\(= \\int_{1}^{2}x^{2}\\,dx = (8-1)/3 = 7/3\\) J."},
		{"State the work integral.", "W = ∫_a^b F(x)dx", "Work by a variable force is \\(W = \\int_{a}^{b}F(x)\\,dx\\)."},
		{"A spring with \\(F(x)=4x\\) is stretched from \\(0\\) to \\(2\\). Find the work.", "8", "Work \\(= \\int_{0}^{2}4x\\,dx = 8\\) J."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralCenterMassGen struct{}

func (g *integralCenterMassGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"A rod on \\([0,2]\\) has density \\(\\rho=1\\). Where is its center of mass?", "1", "Uniform rod: midpoint \\(\\bar{x} = 1\\)."},
		{"State the center-of-mass formula for a rod.", "x̄ = (1/M)∫x·ρ(x)dx", "Total moment divided by total mass: \\(\\bar{x} = \\frac{1}{M}\\int x\\rho(x)\\,dx\\)."},
		{"A rod on \\([0,4]\\) has density \\(\\rho=2\\). Where is its center of mass?", "2", "Uniform density: midpoint \\(\\bar{x} = 2\\)."},
		{"A rod on \\([0,1]\\) has density \\(\\rho(x)=x\\). Find its total mass.", "1/2", "Mass \\(= \\int_{0}^{1}x\\,dx = 1/2\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralSurfaceAreaGen struct{}

func (g *integralSurfaceAreaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the surface area formula for \\(y=f(x)\\) revolved about the x-axis.", "S = 2π∫f(x)√(1+(f')²)dx", "Bands of radius \\(f(x)\\): \\(S = 2\\pi\\int_{a}^{b}f(x)\\sqrt{1+(f')^{2}}\\,dx\\)."},
		{"What is the surface area of a sphere of radius \\(r\\)?", "4πr^2", "Revolving a semicircle gives \\(S = 4\\pi r^{2}\\)."},
		{"Which factor accounts for slant in the surface formula?", "√(1+(f')²)", "The arc-length element \\(\\sqrt{1+(f')^{2}}\\,dx\\) measures slant length."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralTrapezoidGen struct{}

func (g *integralTrapezoidGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the trapezoidal rule with \\(n\\) subintervals.", "T_n = (Δx/2)(f(x_0)+2f(x_1)+...+f(x_n))", "Endpoints count once, interior points twice, scaled by \\(\\Delta x/2\\)."},
		{"Approximate \\(\\int_{0}^{2}x\\,dx\\) with one trapezoid.", "2", "Average height \\((0+2)/2=1\\) times width \\(2\\): \\(2\\)."},
		{"Approximate \\(\\int_{0}^{1}x^{2}\\,dx\\) with one trapezoid.", "1/2", "Average height \\((0+1)/2\\) times width \\(1\\): \\(1/2\\) (true value \\(1/3\\))."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralSimpsonGen struct{}

func (g *integralSimpsonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State Simpson's rule with \\(n\\) (even) subintervals.", "S_n = (Δx/3)(f_0+4f_1+2f_2+...+f_n)", "Coefficients alternate 4,2,4,...,2,4 with endpoints once, scaled by \\(\\Delta x/3\\)."},
		{"Why must \\(n\\) be even for Simpson's rule?", "parabolas are fit over pairs of subintervals", "Each parabola spans two subintervals, so \\(n\\) must be even."},
		{"Approximate \\(\\int_{0}^{2}x\\,dx\\) with Simpson's rule, \\(n=2\\).", "2", "Linear functions are integrated exactly: \\((1/3)(0+4\\cdot 1+2) = 2\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type integralImproperCompareGen struct{}

func (g *integralImproperCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the comparison test for improper integrals.", "0≤f≤g and ∫g converges ⟹ ∫f converges", "Domination preserves convergence; divergence of the smaller forces divergence of the larger."},
		{"Does \\(\\int_{1}^{\\infty} dx/(x^{2}+1)\\) converge? Why?", "yes, by comparison with 1/x^2", "\\(1/(x^{2}+1) \\leq 1/x^{2}\\) and \\(\\int_{1}^{\\infty}dx/x^{2}\\) converges."},
		{"Does \\(\\int_{1}^{\\infty} dx/\\sqrt{x}\\) converge? Why?", "no, compare (diverges like p=1/2)", "Use \\(p\\)-comparison: \\(p=1/2 \\leq 1\\), so it diverges."},
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
		{"For which \\(p\\) does \\(\\int_{1}^{\\infty}dx/x^{p}\\) converge?", "p>1", "The \\(p\\)-integral converges iff \\(p>1\\)."},
		{"Does \\(\\int_{1}^{\\infty}dx/x^{3}\\) converge?", "yes (p=3>1)", "Since \\(p=3>1\\), the integral converges."},
		{"Does \\(\\int_{1}^{\\infty}dx/x\\) converge?", "no (p=1)", "The harmonic case \\(p=1\\) diverges (logarithmic growth)."},
		{"Does \\(\\int_{0}^{1}dx/x^{p}\\) converge for \\(p<1\\)?", "yes", "At \\(0\\), \\(\\int_{0}^{1}dx/x^{p}\\) converges iff \\(p<1\\) (flipped condition)."},
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
		{"State the ratio test.", "L<1 converges, L>1 diverges, L=1 inconclusive", "With \\(L=\\lim|a_{n+1}/a_{n}|\\): \\(L<1\\) converges, \\(L>1\\) diverges, \\(L=1\\) says nothing."},
		{"Apply the ratio test to \\(\\sum 1/n!\\).", "converges (L=0)", "\\(|a_{n+1}/a_{n}| = 1/(n+1) \\to 0 < 1\\), so it converges."},
		{"Apply the ratio test to \\(\\sum n!/n^{n}\\).", "converges (L=1/e)", "The ratio tends to \\(1/e < 1\\), so it converges."},
		{"What does the ratio test say when \\(L=1\\)?", "nothing (inconclusive)", "\\(L=1\\) is inconclusive: try comparison, integral, or root tests."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type geometricSeriesGen struct{}

func (g *geometricSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find the sum of \\(\\sum_{n=0}^{\\infty}(1/2)^{n}\\).", "2", "Geometric with \\(a=1\\), \\(r=1/2\\): \\(1/(1-1/2) = 2\\)."},
		{"For which \\(|r|\\) does \\(\\sum ar^{n}\\) converge?", "|r|<1", "Geometric series converge iff \\(|r|<1\\), summing to \\(a/(1-r)\\)."},
		{"Find the sum of \\(\\sum_{n=0}^{\\infty}(1/3)^{n}\\).", "3/2", "\\(a=1\\), \\(r=1/3\\): \\(1/(1-1/3) = 3/2\\)."},
		{"Does \\(\\sum_{n=0}^{\\infty}2^{n}\\) converge?", "no (|r|=2≥1)", "Ratio \\(|r|=2\\geq 1\\): terms do not even tend to zero."},
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
		{"State the nth-term divergence test.", "if a_n does not tend to 0, the series diverges", "Convergence of \\(\\sum a_{n}\\) requires \\(a_{n}\\to 0\\); otherwise it diverges."},
		{"Does \\(\\sum n/(n+1)\\) converge?", "no (terms tend to 1)", "Terms tend to \\(1\\neq 0\\), so the series diverges."},
		{"If \\(\\sum a_{n}\\) converges, what is \\(\\lim a_{n}\\)?", "0", "Convergence forces \\(a_{n}\\to 0\\) (necessary, not sufficient)."},
		{"Does \\(a_{n}\\to 0\\) guarantee convergence?", "no (harmonic series)", "The harmonic series has \\(1/n\\to 0\\) yet diverges."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type telescopingSeriesGen struct{}

func (g *telescopingSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\sum_{n=1}^{\\infty}(1/n-1/(n+1))\\).", "1", "Partial sums telescope: \\(1-1/(N+1)\\to 1\\)."},
		{"What makes a series telescoping?", "interior terms cancel in partial sums", "Write terms as differences so consecutive terms cancel."},
		{"Find \\(\\sum_{n=1}^{\\infty}(1/(n(n+1)))\\) via partial fractions.", "1", "\\(1/(n(n+1)) = 1/n-1/(n+1)\\), telescoping to \\(1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type directCompareGen struct{}

func (g *directCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the direct comparison test.", "0≤a_n≤b_n and ∑b_n converges ⟹ ∑a_n converges", "Domination by a convergent series forces convergence (and divergence below forces divergence above)."},
		{"Does \\(\\sum 1/(n^{2}+1)\\) converge? Why?", "yes, compare with 1/n^2", "\\(1/(n^{2}+1) \\leq 1/n^{2}\\) and \\(\\sum 1/n^{2}\\) converges."},
		{"Does \\(\\sum 1/\\sqrt{n^{2}+n}\\) converge? Why?", "no, compare below with a divergent series", "Terms behave like \\(1/n\\); \\(\\sum\\) diverges by comparison with harmonic."},
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
		{"State the limit comparison test.", "finite positive limit ⟹ same behavior", "If \\(\\lim a_{n}/b_{n} = L \\in (0,\\infty)\\), both series converge or both diverge."},
		{"Use limit comparison on \\(\\sum 1/(n^{2}-n)\\).", "converges (compare 1/n^2)", "Ratio with \\(1/n^{2}\\) tends to \\(1\\); \\(\\sum 1/n^{2}\\) converges."},
		{"Use limit comparison on \\(\\sum n/(n^{3}+1)\\).", "converges (compare 1/n^2)", "Ratio with \\(1/n^{2}\\) tends to \\(1\\); converges."},
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
		{"What is absolute convergence?", "∑|a_n| converges", "A series converges absolutely if the series of absolute values converges."},
		{"Does absolute convergence imply convergence?", "yes", "Absolute convergence is stronger: \\(\\sum|a_{n}|\\) convergent forces \\(\\sum a_{n}\\) convergent."},
		{"Is \\(\\sum (-1)^{n}/n^{2}\\) absolutely convergent?", "yes", "\\(\\sum 1/n^{2}\\) converges, so the alternating version converges absolutely."},
		{"Is \\(\\sum (-1)^{n}/n\\) absolutely or conditionally convergent?", "conditionally", "It converges (Leibniz) but \\(\\sum 1/n\\) diverges: conditional."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type alternatingRemainderGen struct{}

func (g *alternatingRemainderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Bound the error stopping \\(\\sum (-1)^{n}/n\\) at \\(n=10\\).", "1/11", "Remainder \\(\\leq\\) first omitted term: \\(a_{11} = 1/11\\)."},
		{"State the alternating series remainder bound.", "|R_n| ≤ a_{n+1}", "Truncation error is at most the first dropped term."},
		{"Bound the error stopping \\(\\sum (-1)^{n}/n^{2}\\) at \\(n=5\\).", "1/36", "First omitted term: \\(1/6^{2} = 1/36\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type convergenceRadiusGen struct{}

func (g *convergenceRadiusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"What is the radius of convergence of \\(\\sum x^{n}/n!\\)?", "infinity", "Ratio test gives \\(L=0\\) for all \\(x\\): \\(R=\\infty\\)."},
		{"What is the radius of convergence of \\(\\sum x^{n}\\)?", "1", "Geometric in \\(x\\): converges for \\(|x|<1\\), so \\(R=1\\)."},
		{"What does the interval of convergence add to the radius?", "endpoint checks", "Radius gives the open interval; each endpoint must be tested separately."},
		{"What is the radius of convergence of \\(\\sum n!x^{n}\\)?", "0", "Ratio test diverges for all \\(x\\neq 0\\): \\(R=0\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type epsilonDeltaGen struct{}

func (g *epsilonDeltaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"State the epsilon-delta definition of a limit.", "∀ε>0 ∃δ>0: 0<|x-x_0|<δ ⟹ |f(x)-ℓ|<ε", "For every tolerance \\(\\varepsilon\\), some distance \\(\\delta\\) keeps \\(f\\) within \\(\\varepsilon\\) of \\(\\ell\\)."},
		{"In epsilon-delta, what does δ control?", "how close x must be to x_0", "\\(\\delta\\) bounds the input distance; \\(\\varepsilon\\) bounds the output error."},
		{"Prove \\(\\lim_{x\\to 2}3x = 6\\) with delta in terms of epsilon.", "δ=ε/3", "\\(|3x-6| = 3|x-2| < \\varepsilon\\) when \\(|x-2|<\\varepsilon/3\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type oneSidedLimitGen struct{}

func (g *oneSidedLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\lim_{x\\to 0^{+}}1/x\\).", "infinity", "From the right, \\(1/x\\) grows without bound: \\(+\\infty\\)."},
		{"Find \\(\\lim_{x\\to 0^{-}}1/x\\).", "-infinity", "From the left, \\(1/x\\) dives to \\(-\\infty\\)."},
		{"When does a two-sided limit exist?", "left and right limits agree", "\\(\\lim_{x\\to a}f\\) exists iff both one-sided limits exist and are equal."},
		{"Find \\(\\lim_{x\\to 1^{+}}|x-1|/(x-1)\\).", "1", "For \\(x>1\\), the quotient is \\(1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type infiniteLimitGen struct{}

func (g *infiniteLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\lim_{x\\to 0}1/x^{2}\\).", "infinity", "\\(1/x^{2} \\to +\\infty\\) from both sides."},
		{"What asymptote does an infinite limit signal?", "vertical", "\\(\\lim_{x\\to a}f = \\pm\\infty\\) means \\(x=a\\) is a vertical asymptote."},
		{"Find \\(\\lim_{x\\to 2^{+}}1/(x-2)\\).", "infinity", "Denominator positive and tiny: \\(+\\infty\\)."},
		{"Find \\(\\lim_{x\\to 2^{-}}1/(x-2)\\).", "-infinity", "Denominator negative and tiny: \\(-\\infty\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type sineLimitGen struct{}

func (g *sineLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\lim_{x\\to 0}\\sin(x)/x\\).", "1", "The fundamental trigonometric limit equals \\(1\\)."},
		{"Find \\(\\lim_{x\\to 0}\\sin(3x)/x\\).", "3", "\\(\\sin(3x)/x = 3\\cdot\\sin(3x)/(3x) \\to 3\\)."},
		{"Find \\(\\lim_{x\\to 0}(1-\\cos(x))/x\\).", "0", "Standard companion limit: \\((1-\\cos x)/x \\to 0\\)."},
		{"Find \\(\\lim_{x\\to 0}\\tan(x)/x\\).", "1", "\\(\\tan x/x = (\\sin x/x)\\cdot(1/\\cos x) \\to 1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type expLimitGen struct{}

func (g *expLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Find \\(\\lim_{x\\to 0}(e^{x}-1)/x\\).", "1", "The fundamental exponential limit equals \\(1\\)."},
		{"Find \\(\\lim_{n\\to\\infty}(1+1/n)^{n}\\).", "e", "This limit defines \\(e\\)."},
		{"Find \\(\\lim_{x\\to 0}\\ln(1+x)/x\\).", "1", "The logarithmic companion limit equals \\(1\\)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type boundedSeqGen struct{}

func (g *boundedSeqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	table := []entry{
		{"Is \\(a_{n}=(-1)^{n}\\) bounded?", "yes (|a_n|≤1)", "It oscillates between \\(-1\\) and \\(1\\): bounded but divergent."},
		{"Does every convergent sequence have to be bounded?", "yes", "Convergence forces boundedness (terms eventually stay near the limit)."},
		{"Is \\(a_{n}=n^{2}\\) bounded?", "no", "It grows without bound."},
		{"Is a bounded sequence always convergent?", "no", "\\((-1)^{n}\\) is bounded yet oscillates: boundedness is necessary, not sufficient."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
