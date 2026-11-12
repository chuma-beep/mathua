package prealgebra

// Wave 11B prealgebra: function notation, rules read from a table, equations
// with variables on both sides, equations that need distributing, proportions,
// solution counts, compound and absolute-value inequalities, inequality checks,
// unit conversions, powers and subtraction of negatives. Every branch of a kind
// publishes the same fact keys so one authored schema resolves for all three
// knowledge points.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave11B(reg *generator.Registry) {
	reg.Register("prealg.func.notation", &w11bFuncNotationGen{})
	reg.Register("prealg.func.rule", &w11bFuncRuleGen{})
	reg.Register("prealg.eq.variables_both_sides", &w11bBothSidesGen{})
	reg.Register("prealg.eq.distribute", &w11bDistributeGen{})
	reg.Register("prealg.eq.proportion", &w11bProportionGen{})
	reg.Register("prealg.eq.no_solution", &w11bNoSolutionGen{})
	reg.Register("prealg.ineq.compound", &w11bCompoundGen{})
	reg.Register("prealg.ineq.absolute", &w11bAbsoluteGen{})
	reg.Register("prealg.ineq.check", &w11bIneqCheckGen{})
	reg.Register("prealg.measure.capacity", &w11bCapacityGen{})
	reg.Register("prealg.measure.weight", &w11bWeightGen{})
	reg.Register("prealg.measure.area_units", &w11bAreaUnitsGen{})
	reg.Register("prealg.exp.powers_of_ten", &w11bPowersOfTenGen{})
	reg.Register("prealg.exp.expanded_form", &w11bExpandedFormGen{})
	reg.Register("arith.neg.subtract", &w11bSubtractNegativeGen{})
}

func w11bP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w11bI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w11bOrd(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

// w11bClean redraws a bounded number of times when the answer appears in the
// stem as a standalone numeral. See generator.MaxDrawAttempts: after the
// ceiling the last draw is returned, so a stubborn leak stays a visible quality
// defect rather than an unbounded loop.
func w11bClean(ctx generator.GeneratorContext, build func(o int) generator.Problem) generator.Problem {
	o := w11bOrd(ctx)
	var p generator.Problem
	for i := 0; i < generator.MaxDrawAttempts; i++ {
		p = build(o)
		if !w11bLeaks(p.Question, p.Answer) {
			return p
		}
	}
	return p
}

func w11bLeaks(stem, ans string) bool {
	want, err := strconv.Atoi(ans)
	if err != nil {
		return false
	}
	for i := 0; i < len(stem); {
		if stem[i] < '0' || stem[i] > '9' {
			i++
			continue
		}
		j := i
		for j < len(stem) && stem[j] >= '0' && stem[j] <= '9' {
			j++
		}
		start := i
		if i > 0 && stem[i-1] == '-' && (i == 1 || !w11bNumByte(stem[i-2])) {
			start = i - 1
		}
		before := start == 0 || !w11bNumByte(stem[start-1])
		after := j == len(stem) || !w11bNumByte(stem[j])
		if before && after {
			n, convErr := strconv.Atoi(stem[start:j])
			if convErr == nil && n == want {
				return true
			}
		}
		i = j
	}
	return false
}

func w11bNumByte(b byte) bool {
	return b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' || b == '.'
}

func w11bSide(m, b int) string {
	s := fmt.Sprintf("%dx", m)
	if b < 0 {
		return s + fmt.Sprintf(" - %d", -b)
	}
	return s + fmt.Sprintf(" + %d", b)
}

func w11bEquation(a, b, c, d int) string {
	return fmt.Sprintf("%s = %s", w11bSide(a, b), w11bSide(c, d))
}

func w11bPow(base, exp int) int {
	v := 1
	for i := 0; i < exp; i++ {
		v *= base
	}
	return v
}

// ---------------------------------------------------------------------------
// 1. Function notation
// ---------------------------------------------------------------------------

type w11bFuncNotationGen struct{}

func (g *w11bFuncNotationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		a := 2 + rand.Intn(4)
		facts := map[string]string{}
		w11bI(facts, "a", a)
		switch o {
		case 1:
			b := 1 + rand.Intn(9)
			k := 2 + rand.Intn(7)
			ans := a*k + b
			w11bI(facts, "b", b)
			w11bI(facts, "x", k)
			return w11bP(o,
				fmt.Sprintf("For the function f(x) = %dx + %d, what is f(%d)? (enter a number)", a, b, k),
				strconv.Itoa(ans),
				fmt.Sprintf("Substitute x = %d: f(%d) = %d(%d) + %d = %d.", k, k, a, k, b, ans),
				facts)
		case 2:
			k := 2 + rand.Intn(5)
			b := 1
			for {
				b = 1 + rand.Intn(a*k-1)
				if b != (a-1)*k {
					break
				}
			}
			ans := b - a*k
			w11bI(facts, "b", b)
			w11bI(facts, "x", -k)
			return w11bP(o,
				fmt.Sprintf("For the function f(x) = %dx + %d, what is f(-%d)? (enter a number)", a, b, k),
				strconv.Itoa(ans),
				fmt.Sprintf("Substitute x = -%d: f(-%d) = %d(-%d) + %d = %d.", k, k, a, k, b, ans),
				facts)
		default:
			b := 1 + rand.Intn(9)
			x := 10 + rand.Intn(12)
			value := a*x + b
			w11bI(facts, "b", b)
			w11bI(facts, "x", x)
			return w11bP(o,
				fmt.Sprintf("For the function f(x) = %dx + %d, which input x makes f(x) = %d? (enter a number)", a, b, value),
				strconv.Itoa(x),
				fmt.Sprintf("Solve %dx + %d = %d. Subtract %d and divide by %d: x = %d.", a, b, value, b, a, x),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 2. Writing a function rule from a table
// ---------------------------------------------------------------------------

type w11bFuncRuleGen struct{}

func (g *w11bFuncRuleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		facts := map[string]string{}
		w11bI(facts, "x", 4)
		switch o {
		case 1:
			a := 2 + rand.Intn(5)
			b := 1 + rand.Intn(9)
			y1, y2, y3 := a+b, 2*a+b, 3*a+b
			ans := 4*a + b
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			return w11bP(o,
				fmt.Sprintf("A table follows a linear rule. Its rows are (1, %d), (2, %d) and (3, %d). What is the output when the input is 4? (enter a number)", y1, y2, y3),
				strconv.Itoa(ans),
				fmt.Sprintf("The outputs rise by %d for each step of 1, so the multiplier is %d. Working back, the starting value is %d. At x = 4, y = %d x 4 + %d = %d.", a, a, b, a, b, ans),
				facts)
		case 2:
			a := 5 + rand.Intn(11)
			b := 2 + rand.Intn(8)
			y1, y2, y3 := a+b, 2*a+b, 3*a+b
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			return w11bP(o,
				fmt.Sprintf("A table follows a linear rule. Its rows are (1, %d), (2, %d) and (3, %d). What is the multiplier in the rule y = ax + b? (enter a number)", y1, y2, y3),
				strconv.Itoa(a),
				fmt.Sprintf("The outputs rise by %d each time the input rises by 1, so the multiplier is %d.", y2-y1, a),
				facts)
		default:
			a := 2 + rand.Intn(5)
			b := 5 + rand.Intn(11)
			y1, y2, y3 := a+b, 2*a+b, 3*a+b
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			return w11bP(o,
				fmt.Sprintf("A table follows a linear rule. Its rows are (1, %d), (2, %d) and (3, %d). What is the starting value, the output when the input is 0? (enter a number)", y1, y2, y3),
				strconv.Itoa(b),
				fmt.Sprintf("The multiplier is %d, so from the row (1, %d) the starting value is %d - %d = %d.", a, y1, y1, a, b),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 3. Equations with variables on both sides
// ---------------------------------------------------------------------------

type w11bBothSidesGen struct{}

func (g *w11bBothSidesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		facts := map[string]string{}
		var a, b, c, d, x int
		switch o {
		case 1:
			c = 2 + rand.Intn(4)
			a = c + 1 + rand.Intn(4)
			x = 10 + rand.Intn(12)
			b = 1 + rand.Intn(6)
			d = (a-c)*x + b
		case 2:
			a = 2 + rand.Intn(4)
			c = a + 1 + rand.Intn(4)
			x = 10 + rand.Intn(12)
			d = 1 + rand.Intn(6)
			b = (c-a)*x + d
		default:
			c = 2 + rand.Intn(3)
			a = c + 1 + rand.Intn(3)
			x = 10 + rand.Intn(12)
			d = 1 + rand.Intn(6)
			b = d - (a-c)*x
		}
		w11bI(facts, "a", a)
		w11bI(facts, "b", b)
		w11bI(facts, "c", c)
		w11bI(facts, "d", d)
		q := fmt.Sprintf("Solve for x: %s. (enter a number)", w11bEquation(a, b, c, d))
		if o == 2 {
			q = fmt.Sprintf("Find the value of x in %s. (enter a number)", w11bEquation(a, b, c, d))
		}
		if o == 3 {
			q = fmt.Sprintf("Solve the equation %s, where one side has a subtraction. (enter a number)", w11bEquation(a, b, c, d))
		}
		ex := fmt.Sprintf("Collect x on one side: (%d - %d)x = %d - %d, so %dx = %d and x = %d.", a, c, d, b, a-c, d-b, x)
		return w11bP(o, q, strconv.Itoa(x), ex, facts)
	})
}

// ---------------------------------------------------------------------------
// 4. Equations that need distributing first
// ---------------------------------------------------------------------------

type w11bDistributeGen struct{}

func (g *w11bDistributeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		a := 2 + rand.Intn(5)
		x := 10 + rand.Intn(12)
		facts := map[string]string{}
		w11bI(facts, "a", a)
		switch o {
		case 1:
			b := 1 + rand.Intn(9)
			c := a * (x + b)
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "d", 0)
			return w11bP(o,
				fmt.Sprintf("Solve for x: %d(x + %d) = %d. (enter a number)", a, b, c),
				strconv.Itoa(x),
				fmt.Sprintf("Divide both sides by %d: x + %d = %d. Subtract %d: x = %d.", a, b, c/a, b, x),
				facts)
		case 2:
			b := -(2 + rand.Intn(8))
			c := a * (x + b)
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "d", 0)
			return w11bP(o,
				fmt.Sprintf("Solve for x: %d(x - %d) = %d. (enter a number)", a, -b, c),
				strconv.Itoa(x),
				fmt.Sprintf("Divide both sides by %d: x - %d = %d. Add %d: x = %d.", a, -b, c/a, -b, x),
				facts)
		default:
			b := 1 + rand.Intn(9)
			d := 1 + rand.Intn(9)
			c := a * (x + b)
			e := c + d
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "d", d)
			return w11bP(o,
				fmt.Sprintf("Solve for x: %d(x + %d) + %d = %d. (enter a number)", a, b, d, e),
				strconv.Itoa(x),
				fmt.Sprintf("Subtract %d from both sides: %d(x + %d) = %d. Divide by %d: x + %d = %d, so x = %d.", d, a, b, c, a, b, c/a, x),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 5. Equations that are proportions
// ---------------------------------------------------------------------------

type w11bProportionGen struct{}

func (g *w11bProportionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		facts := map[string]string{}
		switch o {
		case 1:
			a := 2 + rand.Intn(5)
			c := 2 + rand.Intn(5)
			for c == a {
				c = 2 + rand.Intn(5)
			}
			t := 7 + rand.Intn(12)
			x := a * t
			b := c * t
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "s", 0)
			return w11bP(o,
				fmt.Sprintf("Solve for x: x/%d = %d/%d. (enter a number)", a, b, c),
				strconv.Itoa(x),
				fmt.Sprintf("Cross multiply: x times %d equals %d times %d. Then x = %d / %d = %d.", c, a, b, a*b, c, x),
				facts)
		case 2:
			b := 2 + rand.Intn(5)
			c := 2 + rand.Intn(5)
			for c == b {
				c = 2 + rand.Intn(5)
			}
			t := 7 + rand.Intn(12)
			x := c * t
			a := b * t
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "s", 0)
			return w11bP(o,
				fmt.Sprintf("Solve for x: %d/x = %d/%d. (enter a number)", a, b, c),
				strconv.Itoa(x),
				fmt.Sprintf("Cross multiply: %d times %d equals %d times x. Then x = %d / %d = %d.", a, c, b, a*c, b, x),
				facts)
		default:
			a := 2 + rand.Intn(5)
			c := 2 + rand.Intn(5)
			t := 10 + rand.Intn(9)
			b := c * t
			s := 1 + rand.Intn(9)
			x := a*t - s
			w11bI(facts, "a", a)
			w11bI(facts, "b", b)
			w11bI(facts, "c", c)
			w11bI(facts, "s", s)
			return w11bP(o,
				fmt.Sprintf("Solve for x: (x + %d)/%d = %d/%d. (enter a number)", s, a, b, c),
				strconv.Itoa(x),
				fmt.Sprintf("Cross multiply: (x + %d) times %d equals %d times %d. So x + %d = %d and x = %d.", s, c, a, b, s, a*b/c, x),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 6. No solution or infinitely many solutions
// ---------------------------------------------------------------------------

type w11bNoSolutionGen struct{}

func (g *w11bNoSolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		facts := map[string]string{}
		var a, b, c, d, ans int
		var q, ex string
		switch o {
		case 1:
			a = 2 + rand.Intn(8)
			c = 2 + rand.Intn(8)
			for c == a {
				c = 2 + rand.Intn(8)
			}
			b = 2 + rand.Intn(8)
			d = 2 + rand.Intn(8)
			ans = 1
			q = fmt.Sprintf("How many solutions does the equation %s have? Enter 0 when there is no solution and 2 when every number is a solution, otherwise enter the exact count. (enter a number)", w11bEquation(a, b, c, d))
			ex = fmt.Sprintf("Subtract %dx from both sides: (%d - %d)x = %d - %d. The coefficient of x is not zero, so there is exactly one solution.", c, a, c, d, b)
		case 2:
			a = 2 + rand.Intn(8)
			c = a
			b = 2 + rand.Intn(8)
			d = 2 + rand.Intn(8)
			for d == b {
				d = 2 + rand.Intn(8)
			}
			ans = 0
			q = fmt.Sprintf("How many solutions does the equation %s have? Enter 1 when there is exactly one solution and 2 when every number is a solution, otherwise enter the number of solutions. (enter a number)", w11bEquation(a, b, c, d))
			ex = fmt.Sprintf("Subtract %dx from both sides: 0 = %d - %d, which is false, so no value of x works.", c, d, b)
		default:
			a = 3 + rand.Intn(7)
			b = 3 + rand.Intn(7)
			c = a
			d = b
			ans = 2
			q = fmt.Sprintf("The equation %s is true for every value of x, so it has infinitely many solutions. In the convention where 0 means no solution and 1 means exactly one solution, how many solutions does it have? (enter a number)", w11bEquation(a, b, c, d))
			ex = "Both sides are identical, so every value of x is a solution."
		}
		w11bI(facts, "a", a)
		w11bI(facts, "b", b)
		w11bI(facts, "c", c)
		w11bI(facts, "d", d)
		return w11bP(o, q, strconv.Itoa(ans), ex, facts)
	})
}

// ---------------------------------------------------------------------------
// 7. Compound inequalities
// ---------------------------------------------------------------------------

type w11bCompoundGen struct{}

func (g *w11bCompoundGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		facts := map[string]string{}
		var lo, hi, n, ans int
		var q, ex string
		switch o {
		case 1:
			lo = 2 + rand.Intn(29)
			hi = lo + 2 + rand.Intn(21)
			ans = hi - lo - 1
			q = fmt.Sprintf("How many integers x satisfy %d < x < %d? (enter a number)", lo, hi)
			ex = fmt.Sprintf("The integers strictly between %d and %d run from %d to %d, which is %d numbers.", lo, hi, lo+1, hi-1, ans)
		case 2:
			lo = 2 + rand.Intn(29)
			hi = lo + 2 + rand.Intn(21)
			ans = hi - lo + 1
			q = fmt.Sprintf("How many integers x satisfy %d ≤ x ≤ %d? (enter a number)", lo, hi)
			ex = fmt.Sprintf("Count both endpoints: %d - %d + 1 = %d.", hi, lo, ans)
		default:
			n = 20 + rand.Intn(41)
			inside := 1 + rand.Intn(15)
			maxLo := n - inside - 2
			lo = 2 + rand.Intn(maxLo-1)
			hi = lo + inside + 1
			ans = n - inside
			q = fmt.Sprintf("How many integers x with 1 ≤ x ≤ %d lie outside the open interval %d < x < %d? (enter a number)", n, lo, hi)
			ex = fmt.Sprintf("The integers inside the open interval number %d - %d - 1 = %d, so outside it in the range 1 to %d there are %d - %d = %d.", hi, lo, inside, n, n, inside, ans)
		}
		w11bI(facts, "lo", lo)
		w11bI(facts, "hi", hi)
		w11bI(facts, "n", n)
		return w11bP(o, q, strconv.Itoa(ans), ex, facts)
	})
}

// ---------------------------------------------------------------------------
// 8. Absolute value inequalities
// ---------------------------------------------------------------------------

type w11bAbsoluteGen struct{}

func (g *w11bAbsoluteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		var a int
		var q string
		switch o {
		case 1:
			a = 2 + rand.Intn(5)
			q = fmt.Sprintf("How many integer values of x satisfy |x| < %d? (enter a number)", a)
		case 2:
			a = 7 + rand.Intn(5)
			q = fmt.Sprintf("Count the integers x that satisfy |x| < %d. (enter a number)", a)
		default:
			a = 12 + rand.Intn(7)
			q = fmt.Sprintf("What is the number of integers x with |x| < %d? (enter a number)", a)
		}
		ans := 2*a - 1
		facts := map[string]string{}
		w11bI(facts, "a", a)
		return w11bP(o, q, strconv.Itoa(ans),
			fmt.Sprintf("The integers run from -%d to %d, which is %d values.", a-1, a-1, ans), facts)
	})
}

// ---------------------------------------------------------------------------
// 9. Checking a solution to an inequality
// ---------------------------------------------------------------------------

type w11bIneqCheckGen struct{}

func (g *w11bIneqCheckGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		a := 2 + rand.Intn(8)
		b := 2 + rand.Intn(8)
		x := 2 + rand.Intn(8)
		facts := map[string]string{}
		w11bI(facts, "a", a)
		w11bI(facts, "b", b)
		w11bI(facts, "x", x)
		var c, ans int
		var q string
		switch o {
		case 1:
			slack := 1 + rand.Intn(9)
			c = a*x + b + slack
			ans = 1
			q = fmt.Sprintf("Check whether x = %d satisfies %dx + %d < %d. Enter one when the statement is true and zero when it is false. (enter a number)", x, a, b, c)
		case 2:
			base := a*x + b
			slack := 1 + rand.Intn(base-1)
			c = base - slack
			ans = 0
			q = fmt.Sprintf("Test the value x = %d in the inequality %dx + %d < %d. Enter one when the statement is true and zero when it is false. (enter a number)", x, a, b, c)
		default:
			c = a*x + b
			ans = 0
			q = fmt.Sprintf("For x = %d the two sides of %dx + %d = %d are equal. Is the inequality %dx + %d < %d true? Enter one when the statement is true and zero when it is false. (enter a number)", x, a, b, c, a, b, c)
		}
		w11bI(facts, "c", c)
		return w11bP(o, q, strconv.Itoa(ans),
			fmt.Sprintf("Substitute x = %d: %d(%d) + %d = %d. Compare that with %d.", x, a, x, b, a*x+b, c), facts)
	})
}

// ---------------------------------------------------------------------------
// 10-12. Unit conversions
// ---------------------------------------------------------------------------

type w11bUnitCase struct {
	sentence     string
	fromSingular string
	targetPlural string
	factor       int
}

func w11bUnitProblem(o int, cases [3]w11bUnitCase) generator.Problem {
	c := cases[o-1]
	quantity := 2 + rand.Intn(11)
	ans := quantity * c.factor
	facts := map[string]string{}
	w11bI(facts, "quantity", quantity)
	w11bI(facts, "factor", c.factor)
	q := fmt.Sprintf("%s How many %s is that? (1 %s = %d %s) (enter a number)",
		fmt.Sprintf(c.sentence, quantity), c.targetPlural, c.fromSingular, c.factor, c.targetPlural)
	return w11bP(o, q, strconv.Itoa(ans),
		fmt.Sprintf("Multiply the quantity by the conversion factor: %d x %d = %d %s.", quantity, c.factor, ans, c.targetPlural),
		facts)
}

var w11bCapacityCases = [3]w11bUnitCase{
	{"A caterer has %d gallons of lemonade.", "gallon", "cups", 16},
	{"A recipe uses %d quarts of stock.", "quart", "cups", 4},
	{"A tank holds %d gallons of water.", "gallon", "quarts", 4},
}

type w11bCapacityGen struct{}

func (g *w11bCapacityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem { return w11bUnitProblem(o, w11bCapacityCases) })
}

var w11bWeightCases = [3]w11bUnitCase{
	{"A parcel weighs %d pounds.", "pound", "ounces", 16},
	{"A bag of flour has a mass of %d kilograms.", "kilogram", "grams", 1000},
	{"A small truck carries %d tons of gravel.", "ton", "pounds", 2000},
}

type w11bWeightGen struct{}

func (g *w11bWeightGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem { return w11bUnitProblem(o, w11bWeightCases) })
}

var w11bAreaUnitsCases = [3]w11bUnitCase{
	{"A garden measures %d square yards.", "square yard", "square feet", 9},
	{"A tile measures %d square meters.", "square meter", "square centimeters", 10000},
	{"A sheet of metal measures %d square feet.", "square foot", "square inches", 144},
}

type w11bAreaUnitsGen struct{}

func (g *w11bAreaUnitsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem { return w11bUnitProblem(o, w11bAreaUnitsCases) })
}

// ---------------------------------------------------------------------------
// 13. Powers of ten
// ---------------------------------------------------------------------------

type w11bPowersOfTenGen struct{}

func (g *w11bPowersOfTenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		var n int
		var q string
		switch o {
		case 1:
			n = rand.Intn(5)
			// "ten" rather than "10" keeps the base from printing the
			// answer when the exponent is 1.
			q = fmt.Sprintf("What is ten raised to the power %d? (enter a number)", n)
		case 2:
			n = 5 + rand.Intn(5)
			q = fmt.Sprintf("Write 10^%d as an ordinary number. (enter a number)", n)
		default:
			n = 10 + rand.Intn(5)
			q = fmt.Sprintf("A power of ten is written 10^%d. What is its value? (enter a number)", n)
		}
		ans := w11bPow(10, n)
		facts := map[string]string{}
		w11bI(facts, "exponent", n)
		ex := fmt.Sprintf("10^%d is a 1 followed by %d zeros, which is %d.", n, n, ans)
		if n == 0 {
			ex = "Any nonzero number raised to the power 0 equals 1."
		} else if n == 1 {
			ex = "The first power of ten is ten itself, so 10^1 = 10."
		}
		return w11bP(o, q, strconv.Itoa(ans), ex, facts)
	})
}

// ---------------------------------------------------------------------------
// 14. Expanded form of a power
// ---------------------------------------------------------------------------

type w11bExpandedFormGen struct{}

func (g *w11bExpandedFormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		var base, exp int
		var q string
		switch o {
		case 1:
			base, exp = 2, 3+rand.Intn(5)
			q = fmt.Sprintf("What is the value of 2^%d? (enter a number)", exp)
		case 2:
			base, exp = 3, 2+rand.Intn(8)
			q = fmt.Sprintf("A power is written 3^%d. What is its value? (enter a number)", exp)
		default:
			base, exp = 4+rand.Intn(2), 2+rand.Intn(5)
			q = fmt.Sprintf("Evaluate %d^%d. (enter a number)", base, exp)
		}
		ans := w11bPow(base, exp)
		facts := map[string]string{}
		w11bI(facts, "base", base)
		w11bI(facts, "exponent", exp)
		return w11bP(o, q, strconv.Itoa(ans),
			fmt.Sprintf("The exponent says how many copies of %d to multiply: %d^%d = %d.", base, base, exp, ans), facts)
	})
}

// ---------------------------------------------------------------------------
// 15. Subtracting negative numbers
// ---------------------------------------------------------------------------

type w11bSubtractNegativeGen struct{}

func (g *w11bSubtractNegativeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return w11bClean(ctx, func(o int) generator.Problem {
		a := 2 + rand.Intn(11)
		b := 2 + rand.Intn(11)
		facts := map[string]string{}
		w11bI(facts, "a", a)
		w11bI(facts, "b", b)
		var q, ex string
		var ans int
		switch o {
		case 1:
			ans = a + b
			q = fmt.Sprintf("Compute %d - (-%d). (enter a number)", a, b)
			ex = fmt.Sprintf("Subtracting a negative is the same as adding: %d - (-%d) = %d + %d = %d.", a, b, a, b, ans)
		case 2:
			ans = b - a
			q = fmt.Sprintf("Compute -%d - (-%d). (enter a number)", a, b)
			ex = fmt.Sprintf("Subtracting a negative is the same as adding: -%d - (-%d) = -%d + %d = %d.", a, b, a, b, ans)
		default:
			ans = -(a + b)
			q = fmt.Sprintf("Compute -%d - %d. (enter a number)", a, b)
			ex = fmt.Sprintf("Both terms are negative, so add their sizes and keep the minus sign: -%d - %d = %d.", a, b, ans)
		}
		return w11bP(o, q, strconv.Itoa(ans), ex, facts)
	})
}
