package prealgebra

// Wave 10 prealgebra, part two. Each generator publishes a, b, c and answer so
// the shared authored schema resolves on every branch.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave10Rest(reg *generator.Registry) {
	reg.Register("prealg.ratio.tape", &w10bTapeGen{})
	reg.Register("prealg.ratio.constant", &w10bConstantGen{})
	reg.Register("pct.ops.percent_points", &w10bPercentPointsGen{})
	reg.Register("pct.ops.successive", &w10bSuccessiveGen{})
	reg.Register("pct.ops.reverse", &w10bReverseGen{})
	reg.Register("dec.ops.div_decimal", &w10bDivDecimalGen{})
	reg.Register("dec.ops.estimate", &w10bEstimateGen{})
	reg.Register("dec.basics.expanded", &w10bExpandedGen{})
	reg.Register("prealg.expr.two_variables", &w10bTwoVariablesGen{})
	reg.Register("prealg.expr.coefficient", &w10bCoefficientGen{})
	reg.Register("prealg.expr.distribute_simplify", &w10bDistributeSimplifyGen{})
	reg.Register("prealg.eq.check", &w10bCheckGen{})
	reg.Register("prealg.measure.time", &w10bTimeGen{})
	reg.Register("prealg.measure.precision", &w10bPrecisionGen{})
}

func w10bP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w10bI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w10bDec(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

type w10bTapeGen struct{}

func (g *w10bTapeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a, b, c := 1+rand.Intn(5), 1+rand.Intn(5), 1+rand.Intn(5)
	unit := 2 + rand.Intn(8)
	total := (a + b + c) * unit
	facts := map[string]string{}
	w10bI(facts, "a", a)
	w10bI(facts, "b", b)
	w10bI(facts, "c", c)
	w10bI(facts, "unit", unit)
	w10bI(facts, "total", total)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("A tape diagram has %d equal units, and each unit is %d. What is the total?", a+b+c, unit), strconv.Itoa(total), fmt.Sprintf("Multiply the number of units by the unit value: %d x %d = %d.", a+b+c, unit, total), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("A tape diagram has %d equal units and a total of %d. What is the value of one unit?", a+b+c, total), strconv.Itoa(unit), fmt.Sprintf("Divide the total by the number of units: %d / %d = %d.", total, a+b+c, unit), facts)
	default:
		x := 1 + rand.Intn(6)
		y := 1 + rand.Intn(6)
		z := 1 + rand.Intn(6)
		t := (x + y + z) * unit
		w10bI(facts, "a", x)
		w10bI(facts, "b", y)
		w10bI(facts, "c", z)
		w10bI(facts, "total", t)
		return w10bP(o, fmt.Sprintf("Three tape sections contain %d, %d and an unknown number of equal units. Each unit is %d and the total is %d. How many units are in the unknown section?", x, y, unit, t), strconv.Itoa(z), fmt.Sprintf("The total is %d units. Subtract the known sections: %d - %d - %d = %d.", t/unit, t/unit, x, y, z), facts)
	}
}

type w10bConstantGen struct{}

func (g *w10bConstantGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	k := 2 + rand.Intn(8)
	x := 2 + rand.Intn(8)
	y := k * x
	x2 := 3 + rand.Intn(7)
	facts := map[string]string{}
	w10bI(facts, "a", k)
	w10bI(facts, "b", x)
	w10bI(facts, "c", y)
	w10bI(facts, "x2", x2)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("A proportional table contains the pair (%d, %d). What is the constant of proportionality?", x, y), strconv.Itoa(k), fmt.Sprintf("Divide y by x: %d / %d = %d.", y, x, k), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("A proportional relationship has constant %d. What is y when x = %d?", k, x2), strconv.Itoa(k*x2), fmt.Sprintf("Multiply the constant by x: %d x %d = %d.", k, x2, k*x2), facts)
	default:
		target := k * x2
		return w10bP(o, fmt.Sprintf("A proportional relationship has constant %d. For which x is y = %d?", k, target), strconv.Itoa(x2), fmt.Sprintf("Divide y by the constant: %d / %d = %d.", target, k, x2), facts)
	}
}

type w10bPercentPointsGen struct{}

func (g *w10bPercentPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	old := []int{20, 25, 50}[rand.Intn(3)]
	step := 5 * (1 + rand.Intn(4))
	newer := old + step
	if o != 3 && rand.Intn(2) == 0 && old-step >= 5 {
		newer = old - step
	}
	points := newer - old
	facts := map[string]string{}
	w10bI(facts, "a", old)
	w10bI(facts, "b", newer)
	w10bI(facts, "c", points)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("A rate changes from %d%% to %d%%. What is the change in percentage points?", old, newer), strconv.Itoa(points), fmt.Sprintf("Subtract the old rate from the new rate: %d - %d = %d percentage points.", newer, old, points), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("A rate changes from %d%% to %d%%. What is the percent change in the rate itself?", old, newer), strconv.Itoa(points*100/old), fmt.Sprintf("The change is %d percentage points. Divide by the old rate and multiply by 100: %d / %d x 100 = %d%%.", points, points, old, points*100/old), facts)
	default:
		return w10bP(o, fmt.Sprintf("A rate rises by %d percentage points to %d%%. What was the old rate?", step, newer), strconv.Itoa(old), fmt.Sprintf("Subtract the point change from the new rate: %d - %d = %d%%.", newer, step, old), facts)
	}
}

type w10bSuccessiveGen struct{}

func (g *w10bSuccessiveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	base := 100 * (1 + rand.Intn(10))
	p1 := []int{10, 20, 25, 50}[rand.Intn(4)]
	p2 := []int{10, 20, 25, 50}[rand.Intn(4)]
	after1 := base * (100 + p1) / 100
	after2 := after1 * (100 + p2) / 100
	facts := map[string]string{}
	w10bI(facts, "a", base)
	w10bI(facts, "b", p1)
	w10bI(facts, "c", p2)
	w10bI(facts, "after1", after1)
	w10bI(facts, "after2", after2)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("A price of %d is increased by %d%% and then by another %d%%. What is the final price?", base, p1, p2), strconv.Itoa(after2), fmt.Sprintf("After the first increase the price is %d. After the second increase it is %d.", after1, after2), facts)
	case 2:
		down := after1 * (100 - p2) / 100
		facts["after2"] = strconv.Itoa(down)
		return w10bP(o, fmt.Sprintf("A price of %d is increased by %d%% and then decreased by %d%%. What is the final price?", base, p1, p2), strconv.Itoa(down), fmt.Sprintf("After the increase the price is %d. A %d%% decrease leaves %d.", after1, p2, down), facts)
	default:
		net := ((100+p1)*(100+p2) - 10000) / 100
		return w10bP(o, fmt.Sprintf("A quantity is increased by %d%% and then by %d%%. What is the net percent change?", p1, p2), strconv.Itoa(net), fmt.Sprintf("The combined multiplier is %d/100 x %d/100, so the net change is %d%%.", 100+p1, 100+p2, net), facts)
	}
}

type w10bReverseGen struct{}

func (g *w10bReverseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	p := []int{20, 25, 40, 50, 80}[rand.Intn(5)]
	whole := 100 * (1 + rand.Intn(10))
	part := whole * p / 100
	facts := map[string]string{}
	w10bI(facts, "a", p)
	w10bI(facts, "b", part)
	w10bI(facts, "c", whole)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("%d is %d%% of what number?", part, p), strconv.Itoa(whole), fmt.Sprintf("Divide the part by the rate: %d / 0.%02d = %d.", part, p, whole), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("What is %d%% of %d?", p, whole), strconv.Itoa(part), fmt.Sprintf("Multiply the whole by the rate: %d x %d/100 = %d.", whole, p, part), facts)
	default:
		return w10bP(o, fmt.Sprintf("%d is what percent of %d?", part, whole), strconv.Itoa(p), fmt.Sprintf("Divide the part by the whole and multiply by 100: %d / %d x 100 = %d%%.", part, whole, p), facts)
	}
}

type w10bDivDecimalGen struct{}

func (g *w10bDivDecimalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	tenths := []int{2, 4, 5, 8, 15, 25}[rand.Intn(6)]
	q := 2 + rand.Intn(18)
	divisor := float64(tenths) / 10
	dividendTenths := tenths * q
	dividend := float64(dividendTenths) / 10
	facts := map[string]string{}
	w10bI(facts, "a", tenths)
	w10bI(facts, "b", q)
	w10bI(facts, "c", q)
	facts["divisor"] = w10bDec(divisor)
	facts["dividend"] = w10bDec(dividend)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("Compute %s ÷ %s.", w10bDec(dividend), w10bDec(divisor)), strconv.Itoa(q), fmt.Sprintf("Multiply both numbers by 10 to remove the decimal: %d ÷ %d = %d.", dividendTenths, tenths, q), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("Rewrite %s ÷ %s by multiplying both numbers by 10. What is the new dividend?", w10bDec(dividend), w10bDec(divisor)), strconv.Itoa(dividendTenths), fmt.Sprintf("Multiplying the dividend by 10 gives %d.", dividendTenths), facts)
	default:
		return w10bP(o, fmt.Sprintf("A quotient is %d and the divisor is %s. What is the dividend?", q, w10bDec(divisor)), w10bDec(dividend), fmt.Sprintf("Multiply the quotient by the divisor: %d x %s = %s.", q, w10bDec(divisor), w10bDec(dividend)), facts)
	}
}

type w10bEstimateGen struct{}

func (g *w10bEstimateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a := float64(1+rand.Intn(90)) / 10
	b := float64(1+rand.Intn(90)) / 10
	ea := float64(int(a + 0.5))
	eb := float64(int(b + 0.5))
	facts := map[string]string{}
	facts["a"] = w10bDec(a)
	facts["b"] = w10bDec(b)
	facts["c"] = w10bDec(a + b)
	switch o {
	case 1:
		est := ea + eb
		facts["estimate"] = w10bDec(est)
		facts["exact"] = w10bDec(a + b)
		return w10bP(o, fmt.Sprintf("Estimate %s + %s by rounding each number to the nearest whole number.", w10bDec(a), w10bDec(b)), w10bDec(est), fmt.Sprintf("Rounding gives %s + %s = %s.", w10bDec(ea), w10bDec(eb), w10bDec(est)), facts)
	case 2:
		est := ea * eb
		facts["estimate"] = w10bDec(est)
		facts["exact"] = w10bDec(a * b)
		return w10bP(o, fmt.Sprintf("Estimate %s x %s by rounding each number to the nearest whole number.", w10bDec(a), w10bDec(b)), w10bDec(est), fmt.Sprintf("Rounding gives %s x %s = %s.", w10bDec(ea), w10bDec(eb), w10bDec(est)), facts)
	default:
		exact := a + b
		facts["estimate"] = w10bDec(ea + eb)
		facts["exact"] = w10bDec(exact)
		return w10bP(o, fmt.Sprintf("What is the exact value of %s + %s?", w10bDec(a), w10bDec(b)), w10bDec(exact), fmt.Sprintf("Add the decimals: %s + %s = %s.", w10bDec(a), w10bDec(b), w10bDec(exact)), facts)
	}
}

type w10bExpandedGen struct{}

func (g *w10bExpandedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	h, t, u := rand.Intn(10), rand.Intn(10), rand.Intn(10)
	tenths, hundredths := rand.Intn(10), rand.Intn(10)
	facts := map[string]string{}
	w10bI(facts, "a", h)
	w10bI(facts, "b", t)
	w10bI(facts, "c", u)
	w10bI(facts, "tenths", tenths)
	w10bI(facts, "hundredths", hundredths)
	switch o {
	case 1:
		asked := "tenths"
		val := float64(tenths) / 10
		if rand.Intn(2) == 0 {
			asked = "hundredths"
			val = float64(hundredths) / 100
		}
		facts["asked"] = asked
		return w10bP(o, fmt.Sprintf("What is the place value of the digit in the %s place of %d.%02d?", asked, h*100+t*10+u, tenths*10+hundredths), w10bDec(val), fmt.Sprintf("The %s digit contributes %s to the number.", asked, w10bDec(val)), facts)
	case 2:
		val := float64(h*10000+t*1000+u*100+tenths*10+hundredths) / 100
		return w10bP(o, fmt.Sprintf("Write %d.%02d as the sum of its place values. What is the sum?", h*100+t*10+u, tenths*10+hundredths), w10bDec(val), fmt.Sprintf("The place values add to %s.", w10bDec(val)), facts)
	default:
		missing := float64(hundredths) / 100
		return w10bP(o, fmt.Sprintf("The number %d.%02d is built from place values. What value does the hundredths digit contribute?", h*100+t*10+u, tenths*10+hundredths), w10bDec(missing), fmt.Sprintf("The hundredths digit contributes %s.", w10bDec(missing)), facts)
	}
}

type w10bTwoVariablesGen struct{}

func (g *w10bTwoVariablesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a, b, c := 2+rand.Intn(6), 2+rand.Intn(6), rand.Intn(10)
	x, y := 1+rand.Intn(9), 1+rand.Intn(9)
	val := a*x + b*y + c
	facts := map[string]string{}
	w10bI(facts, "a", a)
	w10bI(facts, "b", b)
	w10bI(facts, "c", c)
	w10bI(facts, "x", x)
	w10bI(facts, "y", y)
	w10bI(facts, "value", val)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("Evaluate %dx + %dy when x = %d and y = %d.", a, b, x, y), strconv.Itoa(a*x+b*y), fmt.Sprintf("Substitute: %d(%d) + %d(%d) = %d.", a, x, b, y, a*x+b*y), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("Evaluate %dx + %dy + %d when x = %d and y = %d.", a, b, c, x, y), strconv.Itoa(val), fmt.Sprintf("Substitute and add: %d + %d = %d.", a*x+b*y, c, val), facts)
	default:
		return w10bP(o, fmt.Sprintf("For %dx + %dy + %d = %d and y = %d, what is x?", a, b, c, val, y), strconv.Itoa(x), fmt.Sprintf("Subtract the known terms: %d - %d - %d = %d. Divide by %d to get x = %d.", val, b*y, c, a*x, a, x), facts)
	}
}

type w10bCoefficientGen struct{}

func (g *w10bCoefficientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a, b, c := 2+rand.Intn(7), 2+rand.Intn(7), rand.Intn(10)
	facts := map[string]string{}
	w10bI(facts, "a", a)
	w10bI(facts, "b", b)
	w10bI(facts, "c", c)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("Simplify the coefficient of x in %dx + %dx + %d. What is the coefficient?", a, b, c), strconv.Itoa(a+b), fmt.Sprintf("Add the like terms: %d + %d = %d, so the coefficient is %d.", a, b, a+b, a+b), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("In the expression %dx + %dx + %d, what is the constant term?", a, b, c), strconv.Itoa(c), fmt.Sprintf("The constant term is the term without x: %d.", c), facts)
	default:
		terms := 2
		if c == 0 {
			terms = 1
		}
		return w10bP(o, fmt.Sprintf("After combining like terms, how many terms does %dx + %dx + %d have?", a, b, c), strconv.Itoa(terms), fmt.Sprintf("The x terms combine into one term and the constant is %d, giving %d term(s).", c, terms), facts)
	}
}

type w10bDistributeSimplifyGen struct{}

func (g *w10bDistributeSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a, b, c := 2+rand.Intn(7), 1+rand.Intn(8), 1+rand.Intn(9)
	x := 1 + rand.Intn(9)
	facts := map[string]string{}
	w10bI(facts, "a", a)
	w10bI(facts, "b", b)
	w10bI(facts, "c", c)
	w10bI(facts, "x", x)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("Distribute %d(x + %d) + %d. What is the coefficient of x?", a, b, c), strconv.Itoa(a), fmt.Sprintf("The x term is %dx, so the coefficient is %d.", a, a), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("Distribute %d(x + %d) + %d. What is the constant term?", a, b, c), strconv.Itoa(a*b+c), fmt.Sprintf("Multiply %d by %d and add %d: %d + %d = %d.", a, b, c, a*b, c, a*b+c), facts)
	default:
		return w10bP(o, fmt.Sprintf("Evaluate %d(x + %d) + %d at x = %d.", a, b, c, x), strconv.Itoa(a*(x+b)+c), fmt.Sprintf("First x + %d = %d. Then %d x %d + %d = %d.", b, x+b, a, x+b, c, a*(x+b)+c), facts)
	}
}

type w10bCheckGen struct{}

func (g *w10bCheckGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a, b, x := 1+rand.Intn(8), rand.Intn(15), 1+rand.Intn(9)
	c := a*x + b
	proposed := x
	if rand.Intn(3) > 0 {
		proposed = x + []int{-3, -2, -1, 1, 2, 3}[rand.Intn(6)]
	}
	lhs := a*proposed + b
	facts := map[string]string{}
	w10bI(facts, "a", a)
	w10bI(facts, "b", b)
	w10bI(facts, "c", c)
	w10bI(facts, "x", x)
	w10bI(facts, "proposed", proposed)
	w10bI(facts, "lhs", lhs)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("To check whether x = %d solves %dx + %d = %d, what is the value of the left side?", proposed, a, b, c), strconv.Itoa(lhs), fmt.Sprintf("Substitute x = %d: %d(%d) + %d = %d.", proposed, a, proposed, b, lhs), facts)
	case 2:
		return w10bP(o, fmt.Sprintf("For %dx + %d = %d and a proposed solution x = %d, what is left side minus right side?", a, b, c, proposed), strconv.Itoa(lhs-c), fmt.Sprintf("The left side is %d and the right side is %d, so the difference is %d.", lhs, c, lhs-c), facts)
	default:
		return w10bP(o, fmt.Sprintf("Solve %dx + %d = %d for x.", a, b, c), strconv.Itoa(x), fmt.Sprintf("Subtract %d from both sides to get %d, then divide by %d to get %d.", b, a*x, a, x), facts)
	}
}

type w10bTimeGen struct{}

func (g *w10bTimeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	h, m := 1+rand.Intn(6), rand.Intn(60)
	facts := map[string]string{}
	w10bI(facts, "a", h)
	w10bI(facts, "b", m)
	w10bI(facts, "c", m)
	switch o {
	case 1:
		return w10bP(o, fmt.Sprintf("How many minutes are in %d hours and %d minutes?", h, m), strconv.Itoa(60*h+m), fmt.Sprintf("Each hour has 60 minutes: 60 x %d + %d = %d.", h, m, 60*h+m), facts)
	case 2:
		total := 60*h + m
		return w10bP(o, fmt.Sprintf("How many seconds are in %d minutes?", total), strconv.Itoa(60*total), fmt.Sprintf("Multiply by 60: %d x 60 = %d.", total, 60*total), facts)
	default:
		return w10bP(o, fmt.Sprintf("An event starts at %d:%02d and lasts %d minutes. How many minutes after midnight does it end?", h, m, 45), strconv.Itoa(60*h+m+45), fmt.Sprintf("Convert the start to minutes, %d, and add 45 to get %d.", 60*h+m, 60*h+m+45), facts)
	}
}

type w10bPrecisionGen struct{}

func (g *w10bPrecisionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w10Ord(ctx)
	a := float64(1+rand.Intn(900)) / 100
	b := float64(1+rand.Intn(900)) / 100
	facts := map[string]string{}
	facts["a"] = w10bDec(a)
	facts["b"] = w10bDec(b)
	facts["c"] = w10bDec(a + b)
	switch o {
	case 1:
		digits := 0
		s := w10bDec(a)
		seen := false
		for _, r := range s {
			if r >= '1' && r <= '9' {
				seen = true
			}
			if seen && r >= '0' && r <= '9' {
				digits++
			}
		}
		if digits < 1 {
			digits = 1
		}
		w10bI(facts, "digits", digits)
		return w10bP(o, fmt.Sprintf("How many significant digits are in %s?", s), strconv.Itoa(digits), fmt.Sprintf("Count from the first nonzero digit through the final digit: %d significant digits.", digits), facts)
	case 2:
		rounded := float64(int(a*10+0.5)) / 10
		return w10bP(o, fmt.Sprintf("Round %s to the nearest tenth.", w10bDec(a)), w10bDec(rounded), fmt.Sprintf("The hundredths digit decides the rounding, giving %s.", w10bDec(rounded)), facts)
	default:
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		return w10bP(o, fmt.Sprintf("What is the difference between the two measurements %s and %s?", w10bDec(a), w10bDec(b)), w10bDec(diff), fmt.Sprintf("Subtract the smaller measurement from the larger: %s.", w10bDec(diff)), facts)
	}
}
