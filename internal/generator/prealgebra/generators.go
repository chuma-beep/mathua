package prealgebra

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("dec.basics.compare", &decCompareGen{})
	reg.Register("dec.ops.add", &decAddSubGen{op: "+"})
	reg.Register("dec.ops.sub", &decAddSubGen{op: "-"})
	reg.Register("dec.ops.mult", &decMultGen{})
	reg.Register("dec.ops.div", &decDivGen{})
	reg.Register("dec.convert.from_frac", &decFromFracGen{})
	reg.Register("dec.convert.to_frac", &decToFracGen{})
	reg.Register("dec.ops.round", &decRoundGen{})

	reg.Register("pct.basics.concept", &pctConceptGen{})
	reg.Register("pct.convert.to_dec", &pctToDecGen{})
	reg.Register("pct.convert.from_dec", &pctFromDecGen{})
	reg.Register("pct.ops.of_number", &pctOfNumberGen{})
	reg.Register("pct.ops.find_rate", &pctFindRateGen{})
	reg.Register("pct.ops.increase", &pctIncreaseGen{})
	reg.Register("pct.ops.discount", &pctDiscountGen{})
	reg.Register("pct.ops.tax_tip", &pctTaxTipGen{})

	reg.Register("arith.neg.abs_value", &absValueGen{})
	reg.Register("arith.neg.order_ops", &negOrderOpsGen{})

	reg.Register("ratio.basics.concept", &ratioConceptGen{})
	reg.Register("ratio.ops.simplify", &ratioSimplifyGen{})
	reg.Register("ratio.ops.proportion", &ratioProportionGen{})
	reg.Register("ratio.ops.rate", &ratioRateGen{})
	reg.Register("ratio.ops.scale", &ratioScaleGen{})

	reg.Register("arith.exp.neg", &expNegGen{})
	reg.Register("arith.exp.zero", &expZeroGen{})
	reg.Register("arith.exp.sci_notation", &sciNotationGen{})
	reg.Register("arith.sci_notation.ops", &sciNotationOpsGen{})

	reg.Register("prealg.var.concept", &varConceptGen{})
	reg.Register("prealg.expr.evaluate", &exprEvalGen{})
	reg.Register("prealg.expr.like_terms", &likeTermsGen{})
	reg.Register("prealg.expr.distribute", &distributeGen{})
	reg.Register("prealg.eq.one_step_add", &eqOneStepAddGen{})
	reg.Register("prealg.eq.one_step_mult", &eqOneStepMultGen{})
	reg.Register("prealg.eq.two_step", &eqTwoStepGen{})
	reg.Register("prealg.eq.word", &eqWordGen{})
	reg.Register("prealg.ineq.one_step", &ineqOneStepGen{})
	reg.Register("prealg.ineq.two_step", &ineqTwoStepGen{})

	reg.Register("prealg.real.concept", &realConceptGen{})
	reg.Register("prealg.real.properties", &realPropertiesGen{})
	reg.Register("prealg.types", &typesGen{})
}

// Decimals

type decCompareGen struct{}

func (g *decCompareGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := float64(rand.Intn(scale*2000)) / 100
	b := float64(rand.Intn(scale*2000)) / 100
	ans := ">"
	if b > a {
		ans = "<"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compare: %.2f __ %.2f  (enter >, <, or =)", a, b),
		Answer:      ans,
		Explanation: fmt.Sprintf("%.2f %s %.2f", a, ans, b),
	}
}

type decAddSubGen struct{ op string }

func (g *decAddSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := float64(rand.Intn(scale*10000)) / 100
	b := float64(rand.Intn(scale*6000)) / 100
	var result float64
	if g.op == "+" {
		result = a + b
	} else {
		result = a - b
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%.2f %s %.2f = ?", a, g.op, b),
		Answer:      fmt.Sprintf("%.2f", result),
		Explanation: fmt.Sprintf("%.2f %s %.2f = %.2f", a, g.op, b, result),
	}
}

type decMultGen struct{}

func (g *decMultGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := float64(rand.Intn(scale*100)) / 10
	b := float64(rand.Intn(scale*40)) / 10
	return generator.Problem{
		Question:    fmt.Sprintf("%.1f x %.1f = ?", a, b),
		Answer:      fmt.Sprintf("%.2f", a*b),
		Explanation: fmt.Sprintf("%.1f x %.1f = %.2f", a, b, a*b),
	}
}

type decDivGen struct{}

func (g *decDivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	b := float64(rand.Intn(max(1, scale*20))+10) / 10
	q := float64(rand.Intn(max(1, scale*10))) / 10
	a := b * q
	return generator.Problem{
		Question:    fmt.Sprintf("%.2f / %.1f = ?", a, b),
		Answer:      fmt.Sprintf("%.1f", q),
		Explanation: fmt.Sprintf("%.2f / %.1f = %.1f", a, b, q),
	}
}

type decFromFracGen struct{}

func (g *decFromFracGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	denList := []int{2, 4, 5, 8, 10, 20, 25, 50}
	den := denList[rand.Intn(len(denList))]
	num := rand.Intn(den-1) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %d/%d to a decimal.", num, den),
		Answer:      fmt.Sprintf("%g", float64(num)/float64(den)),
		Explanation: fmt.Sprintf("%d / %d = %g", num, den, float64(num)/float64(den)),
	}
}

type decToFracGen struct{}

func (g *decToFracGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	decimals := []float64{0.25, 0.5, 0.75, 0.2, 0.4, 0.6, 0.8, 0.125, 0.375, 0.625}
	d := decimals[rand.Intn(len(decimals))]
	den := 1000
	num := int(d * float64(den))
	gcd := mathutil.GCD(num, den)
	num /= gcd
	den /= gcd
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %g to a fraction in simplest form.", d),
		Answer:      fmt.Sprintf("%d/%d", num, den),
		Explanation: fmt.Sprintf("%g = %d/1000 = %d/%d", d, int(d*1000), num, den),
	}
}

type decRoundGen struct{}

func (g *decRoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := float64(rand.Intn(scale*20000)) / 1000
	r := math.Round(n*100) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Round %g to the nearest hundredth (2 decimal places).", n),
		Answer:      fmt.Sprintf("%.2f", r),
		Explanation: fmt.Sprintf("%g rounded to 2 decimals = %.2f", n, r),
	}
}

// Percents

type pctConceptGen struct{}

func (g *pctConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	p := rand.Intn(scale*20) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("How do you write %d percent as a fraction?", p),
		Answer:      fmt.Sprintf("%d/100", p),
		Explanation: fmt.Sprintf("%d%% = %d/100", p, p),
	}
}

type pctToDecGen struct{}

func (g *pctToDecGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	p := rand.Intn(scale*30) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %d%% to a decimal.", p),
		Answer:      fmt.Sprintf("%g", float64(p)/100),
		Explanation: fmt.Sprintf("%d%% = %d/100 = %g", p, p, float64(p)/100),
	}
}

type pctFromDecGen struct{}

func (g *pctFromDecGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	d := float64(rand.Intn(scale*20)) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %g to a percent.", d),
		Answer:      fmt.Sprintf("%g", d*100),
		Explanation: fmt.Sprintf("%g x 100 = %g%%", d, d*100),
	}
}

type pctOfNumberGen struct{}

func (g *pctOfNumberGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	p := rand.Intn(max(1, scale*10)) + 10
	n := (rand.Intn(max(1, scale*8)) + 10) * 10
	result := float64(p) / 100 * float64(n)
	return generator.Problem{
		Question:    fmt.Sprintf("What is %d%% of %d?", p, n),
		Answer:      fmt.Sprintf("%.1f", result),
		Explanation: fmt.Sprintf("%d%% of %d = %.2f x %d = %.1f", p, n, float64(p)/100, n, result),
	}
}

type pctFindRateGen struct{}

func (g *pctFindRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	p := rand.Intn(max(1, scale*8)) + 20
	total := (rand.Intn(max(1, scale*8)) + 10) * 10
	part := p * total / 100
	return generator.Problem{
		Question:    fmt.Sprintf("%d is what percent of %d?", part, total),
		Answer:      fmt.Sprintf("%d", p),
		Explanation: fmt.Sprintf("\\(\\frac{%d}{%d} = \\frac{%d}{%d} = %d\\%%\\)", part, total, p, 100, p),
	}
}

type pctIncreaseGen struct{}

func (g *pctIncreaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	orig := (rand.Intn(scale*20) + 10) * 10
	p := rand.Intn(max(1, scale*6)) + 10
	inc := float64(orig) * float64(p) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("What is a %d%% increase on %d?", p, orig),
		Answer:      fmt.Sprintf("%.1f", float64(orig)+inc),
		Explanation: fmt.Sprintf("%d + %d%% = %d + %.1f = %.1f", orig, p, orig, inc, float64(orig)+inc),
	}
}

type pctDiscountGen struct{}

func (g *pctDiscountGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	orig := (rand.Intn(scale*20) + 10) * 10
	p := rand.Intn(max(1, scale*8)) + 10
	disc := float64(orig) * float64(p) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Price: $%d. %d%% off. What is the sale price?", orig, p),
		Answer:      fmt.Sprintf("%.2f", float64(orig)-disc),
		Explanation: fmt.Sprintf("%d%% of %d = %.2f. %d - %.2f = %.2f.", p, orig, disc, orig, disc, float64(orig)-disc),
	}
}

type pctTaxTipGen struct{}

func (g *pctTaxTipGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	amount := (rand.Intn(scale*20) + 10)
	rate := rand.Intn(max(1, scale*3)) + 5
	tip := float64(amount) * float64(rate) / 100
	items := []string{"bill", "meal", "service"}
	item := items[rand.Intn(len(items))]
	return generator.Problem{
		Question:    fmt.Sprintf("Your %s is $%d. Add a %d%% tip. What is the total?", item, amount, rate),
		Answer:      fmt.Sprintf("%.2f", float64(amount)+tip),
		Explanation: fmt.Sprintf("%d%% of $%d = $%.2f. $%d + $%.2f = $%.2f.", rate, amount, tip, amount, tip, float64(amount)+tip),
	}
}

// Absolute value and negative ops

type absValueGen struct{}

func (g *absValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(scale*6) - 15
	return generator.Problem{
		Question:    fmt.Sprintf("\\(|%d| =\\) ?", n),
		Answer:      fmt.Sprintf("%d", mathutil.Abs(n)),
		Explanation: fmt.Sprintf("The absolute value of %d is %d.", n, mathutil.Abs(n)),
	}
}

type negOrderOpsGen struct{}

func (g *negOrderOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) - 4
	b := rand.Intn(scale*2) - 4
	if b == 0 {
		b = 1
	}
	c := rand.Intn(5) + 2
	result := a + b*c
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: \\(%d + (%d) \\times %d\\)", a, b, c),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("\\(%d + (%d) \\times %d = %d + %d = %d\\)", a, b, c, a, b*c, result),
	}
}

// Ratios

type ratioConceptGen struct{}

func (g *ratioConceptGen) Grade(expected, userAnswer string) grader.Result {
	eParts := strings.SplitN(expected, ":", 2)
	aParts := strings.SplitN(userAnswer, ":", 2)
	if len(eParts) != 2 || len(aParts) != 2 {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid ratio format (use a:b)"}
	}
	e1, err1 := strconv.Atoi(strings.TrimSpace(eParts[0]))
	e2, err2 := strconv.Atoi(strings.TrimSpace(eParts[1]))
	a1, err3 := strconv.Atoi(strings.TrimSpace(aParts[0]))
	a2, err4 := strconv.Atoi(strings.TrimSpace(aParts[1]))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid ratio format"}
	}
	if e1 == a1 && e2 == a2 {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect ratio"}
}

func (g *ratioConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale*2)) + 1
	for a == b {
		b = rand.Intn(max(1, scale*2)) + 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Write the ratio of %d to %d in simplest form.", a, b),
		Answer:      fmt.Sprintf("%d:%d", a, b),
		Explanation: fmt.Sprintf("%d:%d", a, b),
	}
}

type ratioSimplifyGen struct{}

func (g *ratioSimplifyGen) Grade(expected, userAnswer string) grader.Result {
	eParts := strings.SplitN(expected, ":", 2)
	aParts := strings.SplitN(userAnswer, ":", 2)
	if len(eParts) != 2 || len(aParts) != 2 {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid ratio format (use a:b)"}
	}
	e1, err1 := strconv.Atoi(strings.TrimSpace(eParts[0]))
	e2, err2 := strconv.Atoi(strings.TrimSpace(eParts[1]))
	a1, err3 := strconv.Atoi(strings.TrimSpace(aParts[0]))
	a2, err4 := strconv.Atoi(strings.TrimSpace(aParts[1]))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid ratio format"}
	}
	if e1 == a1 && e2 == a2 {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect ratio"}
}

func (g *ratioSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	f := rand.Intn(max(1, scale)) + 2
	a := (rand.Intn(max(1, scale)) + 2) * f
	b := (rand.Intn(max(1, scale)) + 2) * f
	for a == b {
		b += f
	}
	d := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify the ratio %d:%d", a, b),
		Answer:      fmt.Sprintf("%d:%d", a/d, b/d),
		Explanation: fmt.Sprintf("Divide both by %d: \\(\\frac{%d}{%d} : \\frac{%d}{%d} = %d:%d\\)", d, a, d, b, d, a/d, b/d),
	}
}

type ratioProportionGen struct{}

func (g *ratioProportionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*2)) + 1
	b := rand.Intn(max(1, scale)) + 2
	mult := rand.Intn(max(1, scale/2+1)) + 2
	x := a * mult
	c := b * mult
	return generator.Problem{
		Question:    fmt.Sprintf("Solve the proportion: \\(\\frac{%d}{%d} = \\frac{x}{%d}\\)", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(\\frac{%d}{%d} = \\frac{x}{%d}\\), cross-multiply: \\(%dx = %d \\times %d\\), \\(x = %d\\)", a, b, c, a, c, b, x),
	}
}

type ratioRateGen struct{}

func (g *ratioRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	dist := (rand.Intn(max(1, scale*6)) + 10) * 10
	time := rand.Intn(max(1, scale)) + 1
	speed := dist / time
	return generator.Problem{
		Question:    fmt.Sprintf("%d miles in %d hours. What is the rate in miles per hour?", dist, time),
		Answer:      fmt.Sprintf("%d", speed),
		Explanation: fmt.Sprintf("%d / %d = %d mph.", dist, time, speed),
	}
}

type ratioScaleGen struct{}

func (g *ratioScaleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	factor := rand.Intn(max(1, scale)) + 2
	length := (rand.Intn(scale*2) + 1) * 2
	actual := length * factor
	return generator.Problem{
		Question:    fmt.Sprintf("Scale is 1:%d. Model length is %d cm. What is the actual length?", factor, length),
		Answer:      fmt.Sprintf("%d", actual),
		Explanation: fmt.Sprintf("%d x %d = %d cm.", length, factor, actual),
	}
}

// Exponents

type expNegGen struct{}

func (g *expNegGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(max(1, scale)) + 2
	exp := rand.Intn(max(1, scale/2+1)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(%d^{-%d}\\)", base, exp),
		Answer:      fmt.Sprintf("1/%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("\\(%d^{-%d} = \\frac{1}{%d^{%d}} = \\frac{1}{%d}\\)", base, exp, base, exp, mathutil.IntPow(base, exp)),
	}
}

type expZeroGen struct{}

func (g *expZeroGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	base := rand.Intn(scale*2) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d^{0} =\\) ?", base),
		Answer:      "1",
		Explanation: fmt.Sprintf("Any non-zero number raised to 0 equals 1. %d^0 = 1.", base),
	}
}

type sciNotationGen struct{}

func (g *sciNotationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	coeff := float64(rand.Intn(scale*20)+10) / 10
	exp := rand.Intn(max(1, scale)) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Write \\(%g \\times 10^{%d}\\) as a standard number.", coeff, exp),
		Answer:      fmt.Sprintf("%g", coeff*float64(mathutil.IntPow(10, exp))),
		Explanation: fmt.Sprintf("\\(%g \\times 10^{%d} = %g\\)", coeff, exp, coeff*float64(mathutil.IntPow(10, exp))),
	}
}

type sciNotationOpsGen struct{}

func (g *sciNotationOpsGen) Grade(expected, userAnswer string) grader.Result {
	normE := normalizeSciNotation(expected)
	normA := normalizeSciNotation(userAnswer)
	if normE == "" || normA == "" {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid scientific notation format"}
	}
	eCoeff, eExp := parseSciNotationValue(normE)
	aCoeff, aExp := parseSciNotationValue(normA)
	if eCoeff == nil || aCoeff == nil {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid scientific notation format"}
	}
	if *eExp != *aExp {
		return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect exponent"}
	}
	diff := math.Abs(*eCoeff - *aCoeff)
	if diff < 1e-4 {
		return grader.Result{Correct: true, Score: 1}
	}
	return grader.Result{Correct: false, Score: 0, Feedback: "Incorrect coefficient"}
}

func normalizeSciNotation(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	parts := strings.Split(s, "x10^")
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[0]) + "x10^" + strings.TrimSpace(parts[1])
}

func parseSciNotationValue(s string) (*float64, *int) {
	parts := strings.Split(s, "x10^")
	if len(parts) != 2 {
		return nil, nil
	}
	coeff, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, nil
	}
	exp, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, nil
	}
	return &coeff, &exp
}

func (g *sciNotationOpsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := float64(rand.Intn(scale*20)+10) / 10
	b := float64(rand.Intn(scale*20)+10) / 10
	ea := rand.Intn(max(1, scale)) + 1
	eb := rand.Intn(max(1, scale/2+1))
	coeff := a * b
	totalExp := ea + eb
	for coeff >= 10 {
		coeff /= 10
		totalExp++
	}
	return generator.Problem{
		Question:    fmt.Sprintf("\\((%g \\times 10^{%d}) \\times (%g \\times 10^{%d}) =\\) ? (in scientific notation)", a, ea, b, eb),
		Answer:      fmt.Sprintf("%g x 10^%d", coeff, totalExp),
		Explanation: fmt.Sprintf("\\((%g \\times %g) \\times 10^{%d+%d} = %g \\times 10^{%d}\\)", a, b, ea, eb, coeff, totalExp),
	}
}

// Pre-algebra expressions

type varConceptGen struct{}

func (g *varConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(scale*2) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("If \\(x = %d\\), what is \\(3x\\)?", n),
		Answer:      fmt.Sprintf("%d", 3*n),
		Explanation: fmt.Sprintf("3x means 3 times x. 3 x %d = %d.", n, 3*n),
	}
}

type exprEvalGen struct{}

func (g *exprEvalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale) + 2
	b := rand.Intn(scale) + 1
	result := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate \\(%dx + %d\\) when \\(x = %d\\).", a, b, x),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("\\(%d(%d) + %d = %d + %d = %d\\)", a, x, b, a*x, b, result),
	}
}

type likeTermsGen struct{}

func (g *likeTermsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) + 2
	b := rand.Intn(scale*2) + 2
	c := rand.Intn(scale) + 1
	resultCoef := a + b
	return generator.Problem{
		Question:    fmt.Sprintf("Combine like terms: \\(%dx + %dx + %d\\)", a, b, c),
		Answer:      fmt.Sprintf("%dx + %d", resultCoef, c),
		Explanation: fmt.Sprintf("\\(%dx + %dx = %dx\\), plus \\(%d = %dx + %d\\)", a, b, resultCoef, c, resultCoef, c),
	}
}

type distributeGen struct{}

func (g *distributeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(scale*2) + 2
	b := rand.Intn(scale) + 1
	c := rand.Intn(scale) + 1
	resultB := a * b
	resultC := a * c
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: \\(%d(%dx + %d)\\)", a, b, c),
		Answer:      fmt.Sprintf("%dx + %d", resultB, resultC),
		Explanation: fmt.Sprintf("\\(%d(%dx + %d) = %d(%dx) + %d(%d) = %dx + %d\\)", a, b, c, a, b, a, c, resultB, resultC),
	}
}

// Equations

type eqOneStepAddGen struct{}

func (g *eqOneStepAddGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*4) + 1
	b := rand.Intn(scale*5) - scale*2
	a := x - b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x + %d = %d\\)", b, a+b),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(x + %d = %d\\), subtract %d: \\(x = %d\\)", b, a+b, b, x),
	}
}

type eqOneStepMultGen struct{}

func (g *eqOneStepMultGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*2) + 2
	coeff := rand.Intn(scale*2) + 2
	rhs := coeff * x
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%dx = %d\\)", coeff, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%dx = %d\\), divide by %d: \\(x = %d\\)", coeff, rhs, coeff, x),
	}
}

type eqTwoStepGen struct{}

func (g *eqTwoStepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale*2) + 2
	b := rand.Intn(scale*3) - scale
	rhs := a*x + b
	lhs := formatExpr(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %s = %d", lhs, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(%s = %d \\to %dx = %d \\to x = %d\\)", lhs, rhs, a, rhs-b, x),
	}
}

type eqWordGen struct{}

func (g *eqWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*4) + 5
	items := []string{"apples", "coins", "tickets", "pencils"}
	item := items[rand.Intn(len(items))]
	cost := rand.Intn(scale) + 2
	total := x * cost
	return generator.Problem{
		Question:    fmt.Sprintf("You bought %d %s for $%d. How much does each %s cost?", x, item, total, strings.TrimSuffix(item, "s")),
		Answer:      fmt.Sprintf("%d", cost),
		Explanation: fmt.Sprintf("Let p = price. \\(%d \\times p = %d\\). \\(p = \\frac{%d}{%d} = %d\\)", x, total, total, x, cost),
	}
}

// Inequalities

type ineqOneStepGen struct{}

func (g *ineqOneStepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*2) + 2
	b := rand.Intn(scale*3) - scale
	c := x + b - 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(x + %d > %d\\)", b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("\\(x + %d > %d\\), subtract %d: \\(x > %d\\)", b, c, b, c-b),
	}
}

type ineqTwoStepGen struct{}

func (g *ineqTwoStepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(scale*2) + 2
	a := rand.Intn(scale) + 2
	b := rand.Intn(scale*2) - scale
	c := a*x + b + rand.Intn(a) + 1
	q := formatIneq(a, b, c)
	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("Solve \\(%s \\to %dx < %d \\to x < %.1f \\to x = %d\\)", q, a, c-b, float64(c-b)/float64(a), x),
	}
}

func formatIneq(a, b, c int) string {
	lhs := formatExpr(a, b)
	return fmt.Sprintf("%s < %d", lhs, c)
}

func formatExpr(a, b int) string {
	if b == 0 {
		return fmt.Sprintf("%dx", a)
	}
	if b < 0 {
		return fmt.Sprintf("%dx - %d", a, -b)
	}
	return fmt.Sprintf("%dx + %d", a, b)
}

type realConceptGen struct{}

func (g *realConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What is the defining property of a rational number?", "can be expressed as \\(a/b\\)", "A rational number can be written as \\(a/b\\) where \\(a\\) and \\(b\\) are integers, \\(b \\neq 0\\)."},
		{"Is \\(\\pi\\) a rational number? (yes/no)", "no", "\\(\\pi\\) cannot be expressed as a ratio of two integers."},
		{"Is \\(\\sqrt{2}\\) a rational number? (yes/no)", "no", "\\(\\sqrt{2}\\) cannot be expressed as a ratio of two integers."},
		{"Is \\(0.333\\ldots\\) a rational number? (yes/no)", "yes", "\\(0.333\\ldots = 1/3\\), so it is rational."},
		{"The set of real numbers contains which two main subsets?", "rational and irrational", "Real numbers are the union of rational and irrational numbers."},
		{"Every integer is also a rational number. (true/false)", "true", "Any integer \\(n\\) can be written as \\(n/1\\)."},
		{"Is \\(\\sqrt{4}\\) a rational number? (yes/no)", "yes", "\\(\\sqrt{4} = 2 = 2/1\\), so it is rational."},
		{"What is the name for a non-repeating, non-terminating decimal?", "irrational", "Irrational numbers have decimal representations that neither terminate nor repeat."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type realPropertiesGen struct{}

func (g *realPropertiesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Which property says a + b = b + a?", "commutative", "The commutative property states that order does not matter for addition: a + b = b + a."},
		{"Which property says (a + b) + c = a + (b + c)?", "associative", "The associative property states that grouping does not matter for addition."},
		{"Which property says a(b + c) = ab + ac?", "distributive", "The distributive property distributes multiplication over addition."},
		{"What is the additive identity element?", "0", "Adding 0 to any number gives the same number."},
		{"What is the multiplicative identity element?", "1", "Multiplying any number by 1 gives the same number."},
		{"What is the additive inverse of 5?", "-5", "The additive inverse of a number is what you add to it to get 0."},
		{"What is the multiplicative inverse of 3?", "1/3", "The multiplicative inverse of a number is what you multiply it by to get 1."},
		{"Which property says (a × b) × c = a × (b × c)?", "associative", "The associative property for multiplication states that grouping does not matter."},
		{"Which property says a × b = b × a?", "commutative", "The commutative property for multiplication states that order does not matter."},
		{"The identity property of multiplication states a × 1 = ?", "a", "Multiplying by 1 gives the original number."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type typesGen struct{}

func (g *typesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What type of number is 5? (natural, integer, rational, irrational, real)", "natural", "5 is a natural number (counting number)."},
		{"What type of number is -3? (natural, integer, rational, irrational, real)", "integer", "-3 is an integer (not a natural number because it is negative)."},
		{"What type of number is 0? (natural, integer, rational, irrational, real)", "integer", "0 is an integer but not a natural number."},
		{"What type of number is 2/3? (natural, integer, rational, irrational, real)", "rational", "2/3 is a rational number because it can be expressed as a fraction."},
		{"What type of number is 0.5? (natural, integer, rational, irrational, real)", "rational", "0.5 = 1/2, so it is rational."},
		{"What type of number is \\(\\sqrt{2}\\)? (natural, integer, rational, irrational, real)", "irrational", "\\(\\sqrt{2}\\) cannot be expressed as a fraction of integers."},
		{"What type of number is \\(\\pi\\)? (natural, integer, rational, irrational, real)", "irrational", "\\(\\pi\\) is irrational (cannot be expressed as a fraction)."},
		{"All integers are also which larger set of numbers?", "rational", "Every integer n can be written as n/1."},
		{"Which set contains all others: natural, integer, rational, real?", "real", "Real numbers contain all rational and irrational numbers."},
		{"Is 0 a natural number? (yes/no)", "no", "Natural numbers are positive counting numbers: 1, 2, 3, ..."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}
