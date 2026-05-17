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
	reg.Register("dec.compare", &decCompareGen{})
	reg.Register("dec.add", &decAddSubGen{op: "+"})
	reg.Register("dec.sub", &decAddSubGen{op: "-"})
	reg.Register("dec.mult", &decMultGen{})
	reg.Register("dec.div", &decDivGen{})
	reg.Register("dec.from_frac", &decFromFracGen{})
	reg.Register("dec.to_frac", &decToFracGen{})
	reg.Register("dec.round", &decRoundGen{})

	reg.Register("pct.concept", &pctConceptGen{})
	reg.Register("pct.to_dec", &pctToDecGen{})
	reg.Register("pct.from_dec", &pctFromDecGen{})
	reg.Register("pct.of_number", &pctOfNumberGen{})
	reg.Register("pct.find_rate", &pctFindRateGen{})
	reg.Register("pct.increase", &pctIncreaseGen{})
	reg.Register("pct.discount", &pctDiscountGen{})
	reg.Register("pct.tax_tip", &pctTaxTipGen{})

	reg.Register("arith.abs_value", &absValueGen{})
	reg.Register("arith.neg.order_ops", &negOrderOpsGen{})

	reg.Register("ratio.concept", &ratioConceptGen{})
	reg.Register("ratio.simplify", &ratioSimplifyGen{})
	reg.Register("ratio.proportion", &ratioProportionGen{})
	reg.Register("ratio.rate", &ratioRateGen{})
	reg.Register("ratio.scale", &ratioScaleGen{})

	reg.Register("arith.exp.neg", &expNegGen{})
	reg.Register("arith.exp.zero", &expZeroGen{})
	reg.Register("arith.sci_notation", &sciNotationGen{})
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
}

// Decimals

type decCompareGen struct{}

func (g *decCompareGen) Generate(difficulty float64) generator.Problem {
	a := float64(rand.Intn(10000)) / 100
	b := float64(rand.Intn(10000)) / 100
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

func (g *decAddSubGen) Generate(difficulty float64) generator.Problem {
	a := float64(rand.Intn(50000)) / 100
	b := float64(rand.Intn(30000)) / 100
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

func (g *decMultGen) Generate(difficulty float64) generator.Problem {
	a := float64(rand.Intn(500)) / 10
	b := float64(rand.Intn(200)) / 10
	return generator.Problem{
		Question:    fmt.Sprintf("%.1f x %.1f = ?", a, b),
		Answer:      fmt.Sprintf("%.2f", a*b),
		Explanation: fmt.Sprintf("%.1f x %.1f = %.2f", a, b, a*b),
	}
}

type decDivGen struct{}

func (g *decDivGen) Generate(difficulty float64) generator.Problem {
	b := float64(rand.Intn(90)+10) / 10
	q := float64(rand.Intn(50)) / 10
	a := b * q
	return generator.Problem{
		Question:    fmt.Sprintf("%.2f / %.1f = ?", a, b),
		Answer:      fmt.Sprintf("%.1f", q),
		Explanation: fmt.Sprintf("%.2f / %.1f = %.1f", a, b, q),
	}
}

type decFromFracGen struct{}

func (g *decFromFracGen) Generate(difficulty float64) generator.Problem {
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

func (g *decToFracGen) Generate(difficulty float64) generator.Problem {
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

func (g *decRoundGen) Generate(difficulty float64) generator.Problem {
	n := float64(rand.Intn(100000)) / 1000
	r := math.Round(n*100) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Round %g to the nearest hundredth (2 decimal places).", n),
		Answer:      fmt.Sprintf("%.2f", r),
		Explanation: fmt.Sprintf("%g rounded to 2 decimals = %.2f", n, r),
	}
}

// Percents

type pctConceptGen struct{}

func (g *pctConceptGen) Generate(difficulty float64) generator.Problem {
	p := rand.Intn(100) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("How do you write %d percent as a fraction?", p),
		Answer:      fmt.Sprintf("%d/100", p),
		Explanation: fmt.Sprintf("%d%% = %d/100", p, p),
	}
}

type pctToDecGen struct{}

func (g *pctToDecGen) Generate(difficulty float64) generator.Problem {
	p := rand.Intn(150) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %d%% to a decimal.", p),
		Answer:      fmt.Sprintf("%g", float64(p)/100),
		Explanation: fmt.Sprintf("%d%% = %d/100 = %g", p, p, float64(p)/100),
	}
}

type pctFromDecGen struct{}

func (g *pctFromDecGen) Generate(difficulty float64) generator.Problem {
	d := float64(rand.Intn(100)) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Convert %g to a percent.", d),
		Answer:      fmt.Sprintf("%g", d*100),
		Explanation: fmt.Sprintf("%g x 100 = %g%%", d, d*100),
	}
}

type pctOfNumberGen struct{}

func (g *pctOfNumberGen) Generate(difficulty float64) generator.Problem {
	p := rand.Intn(50) + 10
	n := (rand.Intn(40) + 10) * 10
	result := float64(p) / 100 * float64(n)
	return generator.Problem{
		Question:    fmt.Sprintf("What is %d%% of %d?", p, n),
		Answer:      fmt.Sprintf("%.1f", result),
		Explanation: fmt.Sprintf("%d%% of %d = %.2f x %d = %.1f", p, n, float64(p)/100, n, result),
	}
}

type pctFindRateGen struct{}

func (g *pctFindRateGen) Generate(difficulty float64) generator.Problem {
	p := rand.Intn(40) + 20
	total := (rand.Intn(40) + 10) * 10
	part := p * total / 100
	return generator.Problem{
		Question:    fmt.Sprintf("%d is what percent of %d?", part, total),
		Answer:      fmt.Sprintf("%d", p),
		Explanation: fmt.Sprintf("%d/%d = %d/%d = %d%%", part, total, p, 100, p),
	}
}

type pctIncreaseGen struct{}

func (g *pctIncreaseGen) Generate(difficulty float64) generator.Problem {
	orig := (rand.Intn(90) + 10) * 10
	p := rand.Intn(30) + 10
	inc := float64(orig) * float64(p) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("What is a %d%% increase on %d?", p, orig),
		Answer:      fmt.Sprintf("%.1f", float64(orig)+inc),
		Explanation: fmt.Sprintf("%d + %d%% = %d + %.1f = %.1f", orig, p, orig, inc, float64(orig)+inc),
	}
}

type pctDiscountGen struct{}

func (g *pctDiscountGen) Generate(difficulty float64) generator.Problem {
	orig := (rand.Intn(90) + 10) * 10
	p := rand.Intn(40) + 10
	disc := float64(orig) * float64(p) / 100
	return generator.Problem{
		Question:    fmt.Sprintf("Price: $%d. %d%% off. What is the sale price?", orig, p),
		Answer:      fmt.Sprintf("%.2f", float64(orig)-disc),
		Explanation: fmt.Sprintf("%d%% of %d = %.2f. %d - %.2f = %.2f.", p, orig, disc, orig, disc, float64(orig)-disc),
	}
}

type pctTaxTipGen struct{}

func (g *pctTaxTipGen) Generate(difficulty float64) generator.Problem {
	amount := (rand.Intn(90) + 10)
	rate := rand.Intn(15) + 5
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

func (g *absValueGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(30) - 15
	return generator.Problem{
		Question:    fmt.Sprintf("|%d| = ?", n),
		Answer:      fmt.Sprintf("%d", mathutil.Abs(n)),
		Explanation: fmt.Sprintf("The absolute value of %d is %d.", n, mathutil.Abs(n)),
	}
}

type negOrderOpsGen struct{}

func (g *negOrderOpsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(9) - 4
	b := rand.Intn(9) - 4
	if b == 0 {
		b = 1
	}
	c := rand.Intn(5) + 2
	result := a + b*c
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate: %d + (%d) x %d", a, b, c),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("%d + (%d) x %d = %d + %d = %d", a, b, c, a, b*c, result),
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

func (g *ratioConceptGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(9) + 1
	b := rand.Intn(9) + 1
	for a == b {
		b = rand.Intn(9) + 1
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

func (g *ratioSimplifyGen) Generate(difficulty float64) generator.Problem {
	f := rand.Intn(5) + 2
	a := (rand.Intn(6) + 2) * f
	b := (rand.Intn(6) + 2) * f
	for a == b {
		b += f
	}
	d := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify the ratio %d:%d", a, b),
		Answer:      fmt.Sprintf("%d:%d", a/d, b/d),
		Explanation: fmt.Sprintf("Divide both by %d: %d/%d : %d/%d = %d:%d", d, a, d, b, d, a/d, b/d),
	}
}

type ratioProportionGen struct{}

func (g *ratioProportionGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(9) + 1
	b := rand.Intn(5) + 2
	mult := rand.Intn(3) + 2
	x := a * mult
	// a/b = ?/c where c = b*mult
	c := b * mult
	return generator.Problem{
		Question:    fmt.Sprintf("Solve the proportion: %d/%d = x/%d", a, b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%d/%d = x/%d, cross-multiply: %dx = %dx%d, x = %d", a, b, c, a, c, b, x),
	}
}

type ratioRateGen struct{}

func (g *ratioRateGen) Generate(difficulty float64) generator.Problem {
	dist := (rand.Intn(30) + 10) * 10
	time := rand.Intn(5) + 1
	speed := dist / time
	return generator.Problem{
		Question:    fmt.Sprintf("%d miles in %d hours. What is the rate in miles per hour?", dist, time),
		Answer:      fmt.Sprintf("%d", speed),
		Explanation: fmt.Sprintf("%d / %d = %d mph.", dist, time, speed),
	}
}

type ratioScaleGen struct{}

func (g *ratioScaleGen) Generate(difficulty float64) generator.Problem {
	scale := rand.Intn(5) + 2
	length := (rand.Intn(10) + 1) * 2
	actual := length * scale
	return generator.Problem{
		Question:    fmt.Sprintf("Scale is 1:%d. Model length is %d cm. What is the actual length?", scale, length),
		Answer:      fmt.Sprintf("%d", actual),
		Explanation: fmt.Sprintf("%d x %d = %d cm.", length, scale, actual),
	}
}

// Exponents

type expNegGen struct{}

func (g *expNegGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(6) + 2
	exp := rand.Intn(3) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %d^-%d", base, exp),
		Answer:      fmt.Sprintf("1/%d", mathutil.IntPow(base, exp)),
		Explanation: fmt.Sprintf("%d^-%d = 1/%d^%d = 1/%d", base, exp, base, exp, mathutil.IntPow(base, exp)),
	}
}

type expZeroGen struct{}

func (g *expZeroGen) Generate(difficulty float64) generator.Problem {
	base := rand.Intn(9) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("%d^0 = ?", base),
		Answer:      "1",
		Explanation: fmt.Sprintf("Any non-zero number raised to 0 equals 1. %d^0 = 1.", base),
	}
}

type sciNotationGen struct{}

func (g *sciNotationGen) Generate(difficulty float64) generator.Problem {
	coeff := float64(rand.Intn(90)+10) / 10
	exp := rand.Intn(5) + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Write %g x 10^%d as a standard number.", coeff, exp),
		Answer:      fmt.Sprintf("%g", coeff*float64(mathutil.IntPow(10, exp))),
		Explanation: fmt.Sprintf("%g x 10^%d = %g", coeff, exp, coeff*float64(mathutil.IntPow(10, exp))),
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

func (g *sciNotationOpsGen) Generate(difficulty float64) generator.Problem {
	a := float64(rand.Intn(90)+10) / 10
	b := float64(rand.Intn(90)+10) / 10
	ea := rand.Intn(5) + 1
	// Multiply: (a x 10^ea) x (b x 10^eb)
	eb := rand.Intn(3)
	coeff := a * b
	totalExp := ea + eb
	for coeff >= 10 {
		coeff /= 10
		totalExp++
	}
	return generator.Problem{
		Question:    fmt.Sprintf("(%g x 10^%d) x (%g x 10^%d) = ? (in scientific notation)", a, ea, b, eb),
		Answer:      fmt.Sprintf("%g x 10^%d", coeff, totalExp),
		Explanation: fmt.Sprintf("(%g x %g) x 10^(%d+%d) = %g x 10^%d", a, b, ea, eb, coeff, totalExp),
	}
}

// Pre-algebra expressions

type varConceptGen struct{}

func (g *varConceptGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(9) + 2
	return generator.Problem{
		Question:    fmt.Sprintf("If x = %d, what is 3x?", n),
		Answer:      fmt.Sprintf("%d", 3*n),
		Explanation: fmt.Sprintf("3x means 3 times x. 3 x %d = %d.", n, 3*n),
	}
}

type exprEvalGen struct{}

func (g *exprEvalGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(9) + 2
	a := rand.Intn(6) + 2
	b := rand.Intn(5) + 1
	result := a*x + b
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate %dx + %d when x = %d.", a, b, x),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("%d(%d) + %d = %d + %d = %d.", a, x, b, a*x, b, result),
	}
}

type likeTermsGen struct{}

func (g *likeTermsGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(9) + 2
	b := rand.Intn(9) + 2
	c := rand.Intn(6) + 1
	resultCoef := a + b
	return generator.Problem{
		Question:    fmt.Sprintf("Combine like terms: %dx + %dx + %d", a, b, c),
		Answer:      fmt.Sprintf("%dx + %d", resultCoef, c),
		Explanation: fmt.Sprintf("%dx + %dx = %dx, plus %d = %dx + %d", a, b, resultCoef, c, resultCoef, c),
	}
}

type distributeGen struct{}

func (g *distributeGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(8) + 2
	b := rand.Intn(5) + 1
	c := rand.Intn(5) + 1
	resultB := a * b
	resultC := a * c
	return generator.Problem{
		Question:    fmt.Sprintf("Simplify: %d(%dx + %d)", a, b, c),
		Answer:      fmt.Sprintf("%dx + %d", resultB, resultC),
		Explanation: fmt.Sprintf("%d(%dx + %d) = %d(%dx) + %d(%d) = %dx + %d", a, b, c, a, b, a, c, resultB, resultC),
	}
}

// Equations

type eqOneStepAddGen struct{}

func (g *eqOneStepAddGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(20) + 1
	b := rand.Intn(30) - 15
	a := x - b
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: x + %d = %d", b, a+b),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("x + %d = %d, subtract %d: x = %d", b, a+b, b, x),
	}
}

type eqOneStepMultGen struct{}

func (g *eqOneStepMultGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(12) + 2
	coeff := rand.Intn(9) + 2
	rhs := coeff * x
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %dx = %d", coeff, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%dx = %d, divide by %d: x = %d", coeff, rhs, coeff, x),
	}
}

type eqTwoStepGen struct{}

func (g *eqTwoStepGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(10) + 2
	a := rand.Intn(9) + 2
	b := rand.Intn(15) - 7
	rhs := a*x + b
	lhs := formatExpr(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: %s = %d", lhs, rhs),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("%s = %d -> %dx = %d -> x = %d.", lhs, rhs, a, rhs-b, x),
	}
}

type eqWordGen struct{}

func (g *eqWordGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(20) + 5
	items := []string{"apples", "coins", "tickets", "pencils"}
	item := items[rand.Intn(len(items))]
	cost := rand.Intn(3) + 2
	total := x * cost
	return generator.Problem{
		Question:    fmt.Sprintf("You bought %d %s for $%d. How much does each %s cost?", x, item, total, strings.TrimSuffix(item, "s")),
		Answer:      fmt.Sprintf("%d", cost),
		Explanation: fmt.Sprintf("Let p = price. %d x p = %d. p = %d / %d = %d.", x, total, total, x, cost),
	}
}

// Inequalities

type ineqOneStepGen struct{}

func (g *ineqOneStepGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(9) + 2
	b := rand.Intn(15) - 7
	c := x + b - 1
	return generator.Problem{
		Question:    fmt.Sprintf("Solve: x + %d > %d", b, c),
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("x + %d > %d, subtract %d: x > %d.", b, c, b, c-b),
	}
}

type ineqTwoStepGen struct{}

func (g *ineqTwoStepGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(8) + 2
	a := rand.Intn(5) + 2
	b := rand.Intn(8) - 3
	// We want ax + b < c where x is the largest integer solution.
	// x must satisfy: ax + b < c
	// x+1 must NOT satisfy: a(x+1) + b >= c
	// So c in (ax + b, a(x+1) + b]
	c := a*x + b + rand.Intn(a) + 1
	q := formatIneq(a, b, c)
	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("Solve %s -> %dx < %d -> x < %.1f -> x = %d.", q, a, c-b, float64(c-b)/float64(a), x),
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
