package prealgebra

// Wave 11A prealgebra: number-line decimals, decimal comparison, repeating
// decimals, money multi-step problems, benchmark percents, depreciation,
// concentration, tip splitting, similar figures, recipes, mixtures, speed,
// density, negative substitution, common factors and independent/dependent
// variables. Every branch of a kind publishes the same fact keys so an
// authored solution schema resolves for all three knowledge points.

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave11A(reg *generator.Registry) {
	reg.Register("dec.basics.number_line", &w11aNumberLineGen{})
	reg.Register("dec.ops.compare_order", &w11aCompareOrderGen{})
	reg.Register("dec.convert.repeating", &w11aRepeatingGen{})
	reg.Register("dec.ops.word_multi_step", &w11aWordMultiStepGen{})
	reg.Register("pct.basics.benchmarks", &w11aBenchmarkGen{})
	reg.Register("pct.ops.depreciation", &w11aDepreciationGen{})
	reg.Register("pct.ops.concentration", &w11aConcentrationGen{})
	reg.Register("pct.ops.tip_split", &w11aTipSplitGen{})
	reg.Register("ratio.ops.similar_figures", &w11aSimilarFiguresGen{})
	reg.Register("ratio.ops.recipe", &w11aRecipeGen{})
	reg.Register("ratio.ops.mixture", &w11aMixtureGen{})
	reg.Register("ratio.ops.speed", &w11aSpeedGen{})
	reg.Register("ratio.ops.density", &w11aDensityGen{})
	reg.Register("prealg.expr.substitute_negative", &w11aSubstituteNegativeGen{})
	reg.Register("prealg.expr.factor_common", &w11aFactorCommonGen{})
	reg.Register("prealg.var.independent_dependent", &w11aIndependentDependentGen{})
}

func w11aP(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w11aI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w11aOrd(ctx generator.GeneratorContext) int {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	return o
}

func w11aGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// w11aDec renders a whole number of tenths, keeping the sign in front of the
// zero for values between -1 and 0.
func w11aDec(tenths int) string {
	if tenths < 0 {
		return "-" + w11aDec(-tenths)
	}
	return fmt.Sprintf("%d.%d", tenths/10, tenths%10)
}

var (
	w11aOverline    = regexp.MustCompile(`\\overline\{[^}]*\}`)
	w11aNumberToken = regexp.MustCompile(`-?[0-9]+(\.[0-9]+)?`)
)

// w11aPrintsAnswer reports whether the stem already contains the answer as a
// standalone number. A digit inside a repeating-decimal overline is notation,
// not the answer, so that group is removed before the scan.
func w11aPrintsAnswer(question, answer string) bool {
	question = w11aOverline.ReplaceAllString(question, "")
	for _, tok := range w11aNumberToken.FindAllString(question, -1) {
		if tok == answer {
			return true
		}
	}
	return false
}

// w11aAvoidAnswer redraws a problem when its stem would print the answer.
func w11aAvoidAnswer(makeProblem func() generator.Problem) generator.Problem {
	var p generator.Problem
	for i := 0; i < generator.MaxDrawAttempts; i++ {
		p = makeProblem()
		if !w11aPrintsAnswer(p.Question, p.Answer) {
			return p
		}
	}
	return p
}

// w11aBag serves each value once per shuffle. A few ordinals have only a
// handful of possible answers, so an independent draw can repeat itself and
// leave twelve draws with fewer than five distinct answers.
type w11aBag struct {
	mu     sync.Mutex
	values []int
	next   int
}

func w11aNewBag(values ...int) *w11aBag {
	b := &w11aBag{values: append([]int(nil), values...)}
	rand.Shuffle(len(b.values), func(i, j int) { b.values[i], b.values[j] = b.values[j], b.values[i] })
	return b
}

func (b *w11aBag) Next() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.next >= len(b.values) {
		rand.Shuffle(len(b.values), func(i, j int) { b.values[i], b.values[j] = b.values[j], b.values[i] })
		b.next = 0
	}
	v := b.values[b.next]
	b.next++
	return v
}

var (
	w11aNumberTenths1 = w11aNewBag(1, 2, 3, 4, 5, 6, 7, 8, 9)
	w11aNumberTenths3 = w11aNewBag(1, 2, 3, 4, 5, 6, 7, 8, 9)
	w11aRepeatDigits  = w11aNewBag(1, 2, 3, 4, 5, 6, 7, 8)
	w11aCompareCounts = [3]*w11aBag{
		w11aNewBag(0, 1, 2, 3, 4, 5, 6, 7, 8),
		w11aNewBag(0, 1, 2, 3, 4, 5, 6, 7, 8),
		w11aNewBag(0, 1, 2, 3, 4, 5, 6, 7, 8),
	}
	w11aFactorGCFs = w11aNewBag(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
)

func w11aLinearExpr(a, b int) string {
	if b < 0 {
		return fmt.Sprintf("%dx - %d", a, -b)
	}
	return fmt.Sprintf("%dx + %d", a, b)
}

// ---------------------------------------------------------------------------
// 1. A point between whole numbers on a number line
// ---------------------------------------------------------------------------

type w11aNumberLineGen struct{}

func (g *w11aNumberLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		switch o {
		case 1:
			base := 1 + rand.Intn(30)
			tenths := w11aNumberTenths1.Next()
			value := base*10 + tenths
			facts := map[string]string{}
			facts["value"] = w11aDec(value)
			w11aI(facts, "base", base)
			w11aI(facts, "tenths", tenths)
			return w11aP(o,
				fmt.Sprintf("A point sits at %s on a number line, between two whole numbers. How many tenths above the lower whole number is the point? (enter a number)", w11aDec(value)),
				strconv.Itoa(tenths),
				fmt.Sprintf("The lower whole number is %d and the point is %s, so it is %d tenths above it.", base, w11aDec(value), tenths),
				facts)
		case 2:
			base := 1 + rand.Intn(30)
			tenths := 1 + rand.Intn(9)
			value := base*10 + tenths
			answer := base
			if tenths >= 5 {
				answer = base + 1
			}
			facts := map[string]string{}
			facts["value"] = w11aDec(value)
			w11aI(facts, "base", base)
			w11aI(facts, "tenths", tenths)
			return w11aP(o,
				fmt.Sprintf("A point sits at %s on a number line. To which whole number does the point round, with a half rounded away from zero? (enter a number)", w11aDec(value)),
				strconv.Itoa(answer),
				fmt.Sprintf("The tenths digit is %d, so the point rounds to %d.", tenths, answer),
				facts)
		default:
			magnitude := 1 + rand.Intn(30)
			tenths := w11aNumberTenths3.Next()
			value := -(magnitude*10 + tenths)
			base := -magnitude - 1
			above := 10 - tenths
			facts := map[string]string{}
			facts["value"] = w11aDec(value)
			w11aI(facts, "base", base)
			w11aI(facts, "tenths", above)
			return w11aP(o,
				fmt.Sprintf("A point sits at %s on a number line, below the whole numbers. How many tenths above the lower whole number is the point? (enter a number)", w11aDec(value)),
				strconv.Itoa(above),
				fmt.Sprintf("The lower whole number is %d, and %s is %d tenths above it.", base, w11aDec(value), above),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 2. Counting decimals against a threshold
// ---------------------------------------------------------------------------

type w11aCompareOrderGen struct{}

func w11aSplitDecimals(threshold, count, n int, countAbove bool) []int {
	vals := make([]int, 0, n)
	seen := map[int]bool{}
	high := count
	if !countAbove {
		high = n - count
	}
	for len(vals) < high {
		v := threshold + 1 + rand.Intn(80)
		if v == 0 || seen[v] {
			continue
		}
		seen[v] = true
		vals = append(vals, v)
	}
	for len(vals) < n {
		v := threshold - 1 - rand.Intn(80)
		if v == 0 || seen[v] {
			continue
		}
		seen[v] = true
		vals = append(vals, v)
	}
	rand.Shuffle(len(vals), func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	return vals
}

func w11aDecList(vals []int) string {
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = w11aDec(v)
	}
	return strings.Join(parts, ", ")
}

func (g *w11aCompareOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		n := 4 + rand.Intn(5)
		count := w11aCompareCounts[o-1].Next()
		if count > n {
			n = count
		}
		var threshold int
		countAbove := true
		if o == 1 {
			threshold = 10 + rand.Intn(40)
		} else if o == 2 {
			threshold = -(10 + rand.Intn(40))
		} else {
			threshold = 10 + rand.Intn(40)
			countAbove = false
		}
		vals := w11aSplitDecimals(threshold, count, n, countAbove)
		facts := map[string]string{}
		facts["threshold"] = w11aDec(threshold)
		facts["values"] = w11aDecList(vals)
		w11aI(facts, "count", count)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("Consider the decimals %s. How many of them are greater than the positive threshold %s? (enter a number)", w11aDecList(vals), w11aDec(threshold)),
				strconv.Itoa(count),
				fmt.Sprintf("Compare each decimal with %s and count the ones above it: %d.", w11aDec(threshold), count),
				facts)
		case 2:
			return w11aP(o,
				fmt.Sprintf("Consider the decimals %s. The threshold %s is negative. How many of the decimals are greater than it? (enter a number)", w11aDecList(vals), w11aDec(threshold)),
				strconv.Itoa(count),
				fmt.Sprintf("A decimal is greater than %s when it lies to its right on the number line. There are %d such decimals.", w11aDec(threshold), count),
				facts)
		default:
			return w11aP(o,
				fmt.Sprintf("Consider the decimals %s. How many of them are less than the positive threshold %s? (enter a number)", w11aDecList(vals), w11aDec(threshold)),
				strconv.Itoa(count),
				fmt.Sprintf("Compare each decimal with %s and count the ones below it: %d.", w11aDec(threshold), count),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 3. Repeating decimals as fractions with denominator 9, 99 or 999
// ---------------------------------------------------------------------------

type w11aRepeatingGen struct{}

func (g *w11aRepeatingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		var digits string
		var denominator int
		switch o {
		case 1:
			d := w11aRepeatDigits.Next()
			digits = strconv.Itoa(d)
			denominator = 9
		case 2:
			for {
				a := 1 + rand.Intn(9)
				b := rand.Intn(10)
				num := 10*a + b
				if num == 99 {
					continue
				}
				digits = fmt.Sprintf("%d%d", a, b)
				denominator = 99
				break
			}
		default:
			for {
				a := 1 + rand.Intn(9)
				b := rand.Intn(10)
				c := rand.Intn(10)
				num := 100*a + 10*b + c
				if num == 999 {
					continue
				}
				digits = fmt.Sprintf("%d%d%d", a, b, c)
				denominator = 999
				break
			}
		}
		num, _ := strconv.Atoi(digits)
		facts := map[string]string{}
		facts["digit"] = digits
		w11aI(facts, "denominator", denominator)
		word, noun := "two", "digits"
		if o == 1 {
			word, noun = "one", "digit"
		} else if o == 3 {
			word = "three"
		}
		return w11aP(o,
			fmt.Sprintf("The decimal \\(0.\\overline{%s}\\) repeats %s %s. With denominator %d, what numerator gives the same value? (enter a number)", digits, word, noun, denominator),
			strconv.Itoa(num),
			fmt.Sprintf("The repeating block is %s, so the fraction is %d/%d before reducing.", digits, num, denominator),
			facts)
	})
}

// ---------------------------------------------------------------------------
// 4. Multi-step money problems
// ---------------------------------------------------------------------------

type w11aWordMultiStepGen struct{}

func (g *w11aWordMultiStepGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		facts := map[string]string{}
		switch o {
		case 1:
			p1 := 2 * (2 + rand.Intn(8))
			p2 := 2 * (2 + rand.Intn(8))
			for p1 == p2 {
				p2 = 2 * (2 + rand.Intn(8))
			}
			fee := 1 + rand.Intn(5)
			items := 2
			price := (p1 + p2) / 2
			total := items*price + fee
			w11aI(facts, "items", items)
			w11aI(facts, "price", price)
			w11aI(facts, "fee", fee)
			w11aI(facts, "total", total)
			return w11aP(o,
				fmt.Sprintf("A shopper buys two items costing $%d and $%d and pays a $%d fee. What is the total cost in dollars? (enter a number)", p1, p2, fee),
				strconv.Itoa(total),
				fmt.Sprintf("The two items cost $%d + $%d = $%d, and the fee makes $%d.", p1, p2, p1+p2, total),
				facts)
		case 2:
			items := 2 + rand.Intn(4)
			price := 3 + rand.Intn(7)
			discount := 1 + rand.Intn(price-1)
			fee := -discount
			total := items*price + fee
			w11aI(facts, "items", items)
			w11aI(facts, "price", price)
			w11aI(facts, "fee", fee)
			w11aI(facts, "total", total)
			return w11aP(o,
				fmt.Sprintf("A shopper buys %d items at $%d each. A $%d discount is applied. What is the total cost in dollars? (enter a number)", items, price, discount),
				strconv.Itoa(total),
				fmt.Sprintf("The items cost %d x $%d = $%d. Subtract the $%d discount to get $%d.", items, price, items*price, discount, total),
				facts)
		default:
			items := 2 + rand.Intn(6)
			price := 3 + rand.Intn(10)
			fee := 1 + rand.Intn(3)
			cost := items*price + fee
			change := 1 + rand.Intn(30)
			total := cost + change
			w11aI(facts, "items", items)
			w11aI(facts, "price", price)
			w11aI(facts, "fee", fee)
			w11aI(facts, "total", total)
			return w11aP(o,
				fmt.Sprintf("A shopper buys %d items at $%d each and pays a $%d service fee. They hand over $%d. How much change, in dollars, do they get? (enter a number)", items, price, fee, total),
				strconv.Itoa(change),
				fmt.Sprintf("The bill is %d x $%d + $%d = $%d. The change is $%d - $%d = $%d.", items, price, fee, cost, total, cost, change),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 5. Benchmark percents
// ---------------------------------------------------------------------------

type w11aBenchmarkGen struct{}

func (g *w11aBenchmarkGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		percent, number := 0, 0
		switch o {
		case 1:
			percent = 10
			number = 10 * (1 + rand.Intn(60))
		case 2:
			if rand.Intn(2) == 0 {
				percent = 25
			} else {
				percent = 50
			}
			number = 4 * (1 + rand.Intn(60))
		default:
			percent = 1
			number = 100 * (1 + rand.Intn(60))
		}
		answer := number * percent / 100
		facts := map[string]string{}
		w11aI(facts, "percent", percent)
		w11aI(facts, "number", number)
		return w11aP(o,
			fmt.Sprintf("What is %d%% of %d? (enter a number)", percent, number),
			strconv.Itoa(answer),
			fmt.Sprintf("Write %d%% as %d/100 and multiply: %d x %d/100 = %d.", percent, percent, number, percent, answer),
			facts)
	})
}

// ---------------------------------------------------------------------------
// 6. Depreciation after one or two years
// ---------------------------------------------------------------------------

type w11aDepreciationGen struct{}

func (g *w11aDepreciationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		value, rate := 0, 0
		if o == 3 {
			value = 10000 * (1 + rand.Intn(20))
			rate = 10 * (1 + rand.Intn(8))
		} else {
			value = 100 * (1 + rand.Intn(60))
			rate = 5 * (1 + rand.Intn(18))
		}
		facts := map[string]string{}
		w11aI(facts, "value", value)
		w11aI(facts, "rate", rate)
		switch o {
		case 1:
			answer := value * (100 - rate) / 100
			return w11aP(o,
				fmt.Sprintf("A machine worth $%d loses %d%% of its value in one year. What is it worth at the end of the year? (enter a number)", value, rate),
				strconv.Itoa(answer),
				fmt.Sprintf("The loss is $%d, so the machine is worth $%d - $%d = $%d.", value*rate/100, value, value*rate/100, answer),
				facts)
		case 2:
			answer := value * rate / 100
			return w11aP(o,
				fmt.Sprintf("A machine worth $%d loses %d%% of its value in one year. How many dollars of value does it lose? (enter a number)", value, rate),
				strconv.Itoa(answer),
				fmt.Sprintf("Multiply the value by the rate: $%d x %d/100 = $%d.", value, rate, answer),
				facts)
		default:
			first := value * (100 - rate) / 100
			answer := first * (100 - rate) / 100
			return w11aP(o,
				fmt.Sprintf("A machine worth $%d loses %d%% of its value each year for two years. What is it worth after two years? (enter a number)", value, rate),
				strconv.Itoa(answer),
				fmt.Sprintf("After one year it is worth $%d, and after the second year it is worth $%d.", first, answer),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 7. Percent concentration and dilution
// ---------------------------------------------------------------------------

type w11aConcentrationGen struct{}

func (g *w11aConcentrationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		facts := map[string]string{}
		if o == 3 {
			percent := []int{5, 10, 20, 25, 50}[rand.Intn(5)]
			final := 200 * (1 + rand.Intn(50))
			solute := final * percent / 100
			added := 5 * (1 + rand.Intn(38))
			solution := final - added
			w11aI(facts, "solute", solute)
			w11aI(facts, "solution", solution)
			w11aI(facts, "percent", percent)
			return w11aP(o,
				fmt.Sprintf("A solution contains %d mL of solute dissolved in %d mL of solution. How many mL of water must be added to make a %d%% solution? (enter a number)", solute, solution, percent),
				strconv.Itoa(added),
				fmt.Sprintf("A %d%% solution with %d mL of solute has total volume %d mL. It already has %d mL, so add %d mL of water.", percent, solute, final, solution, added),
				facts)
		}
		solution := 100 * (1 + rand.Intn(50))
		percent := 2 + rand.Intn(89)
		solute := solution * percent / 100
		w11aI(facts, "solute", solute)
		w11aI(facts, "solution", solution)
		w11aI(facts, "percent", percent)
		if o == 1 {
			return w11aP(o,
				fmt.Sprintf("A solution contains %d mL of solute in %d mL of solution. What is the percent concentration? (enter a number)", solute, solution),
				strconv.Itoa(percent),
				fmt.Sprintf("Divide the solute by the solution and multiply by 100: %d / %d x 100 = %d%%.", solute, solution, percent),
				facts)
		}
		return w11aP(o,
			fmt.Sprintf("A solution has a concentration of %d%%. The total volume of the solution is %d mL. How many mL of solute are in it? (enter a number)", percent, solution),
			strconv.Itoa(solute),
			fmt.Sprintf("Multiply the volume by the rate: %d x %d/100 = %d mL.", solution, percent, solute),
			facts)
	})
}

// ---------------------------------------------------------------------------
// 8. Tip and split
// ---------------------------------------------------------------------------

type w11aTipSplitGen struct{}

func (g *w11aTipSplitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		people := 2 + rand.Intn(9)
		percent := []int{10, 15, 20, 25}[rand.Intn(4)]
		multiplier := 1 + rand.Intn(10)
		bill := 20 * people * multiplier
		tip := bill * percent / 100
		facts := map[string]string{}
		w11aI(facts, "bill", bill)
		w11aI(facts, "tip_percent", percent)
		w11aI(facts, "people", people)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("A restaurant bill is $%d. A %d%% tip is added. What is the tip in dollars? (enter a number)", bill, percent),
				strconv.Itoa(tip),
				fmt.Sprintf("Multiply the bill by the tip rate: $%d x %d/100 = $%d.", bill, percent, tip),
				facts)
		case 2:
			share := (bill + tip) / people
			return w11aP(o,
				fmt.Sprintf("A restaurant bill is $%d. A %d%% tip is added and the total is split equally among %d people. How much does each person pay in dollars? (enter a number)", bill, percent, people),
				strconv.Itoa(share),
				fmt.Sprintf("The tip is $%d and the total is $%d. Split %d ways gives $%d each.", tip, bill+tip, people, share),
				facts)
		default:
			share := tip / people
			return w11aP(o,
				fmt.Sprintf("A restaurant bill is $%d. A %d%% tip is added and the tip alone is split equally among %d people. How much tip does each person pay in dollars? (enter a number)", bill, percent, people),
				strconv.Itoa(share),
				fmt.Sprintf("The tip is $%d. Split %d ways gives $%d each.", tip, people, share),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 9. Similar figures
// ---------------------------------------------------------------------------

type w11aSimilarFiguresGen struct{}

func (g *w11aSimilarFiguresGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		facts := map[string]string{}
		var a, b, c, answer int
		switch o {
		case 1:
			a = 2 + rand.Intn(8)
			scale := 2 + rand.Intn(4)
			b = a * scale
			c = a + 1 + rand.Intn(10)
			answer = b * c / a
			w11aI(facts, "a", a)
			w11aI(facts, "b", b)
			w11aI(facts, "c", c)
			return w11aP(o,
				fmt.Sprintf("Two similar figures have corresponding sides of %d and %d. A third side of the smaller figure is %d. What is the matching side of the larger figure? (enter a number)", a, b, c),
				strconv.Itoa(answer),
				fmt.Sprintf("The scale factor is %d / %d = %d, so the matching side is %d x %d = %d.", b, a, scale, c, scale, answer),
				facts)
		case 2:
			a = 1
			scale := 2 + rand.Intn(5)
			b = scale
			c = 2 + rand.Intn(12)
			answer = b * c / a
			w11aI(facts, "a", a)
			w11aI(facts, "b", b)
			w11aI(facts, "c", c)
			return w11aP(o,
				fmt.Sprintf("The scale factor from a small figure to a larger similar figure is %d:%d. A side of the small figure is %d. What is the matching side of the larger figure? (enter a number)", a, b, c),
				strconv.Itoa(answer),
				fmt.Sprintf("Multiply the side by the scale factor: %d x %d = %d.", c, b, answer),
				facts)
		default:
			a = 1
			scale := 2 + rand.Intn(5)
			b = scale
			c = 5 + rand.Intn(20)
			answer = b * c / a
			w11aI(facts, "a", a)
			w11aI(facts, "b", b)
			w11aI(facts, "c", c)
			return w11aP(o,
				fmt.Sprintf("Two similar figures have corresponding sides in the ratio %d:%d. The smaller figure has perimeter %d. What is the perimeter of the larger figure? (enter a number)", a, b, c),
				strconv.Itoa(answer),
				fmt.Sprintf("Perimeters scale like the sides, so multiply by %d: %d x %d = %d.", b, c, b, answer),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 10. Recipe scaling
// ---------------------------------------------------------------------------

type w11aRecipeGen struct{}

func (g *w11aRecipeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		servings := 5 + rand.Intn(16)
		unit := 2 + rand.Intn(5)
		cups := servings * unit
		facts := map[string]string{}
		w11aI(facts, "servings", servings)
		w11aI(facts, "cups", cups)
		switch o {
		case 1:
			other := servings + 1 + rand.Intn(6)
			answer := cups * other / servings
			w11aI(facts, "other_servings", other)
			w11aI(facts, "given_cups", answer)
			return w11aP(o,
				fmt.Sprintf("A recipe for %d servings uses %d cups of flour. How many cups are needed to make a larger batch of %d servings? (enter a number)", servings, cups, other),
				strconv.Itoa(answer),
				fmt.Sprintf("Each serving needs %d / %d = %d cups, so %d servings need %d x %d = %d cups.", cups, servings, unit, other, other, unit, answer),
				facts)
		case 2:
			other := 1 + rand.Intn(servings-1)
			answer := cups * other / servings
			w11aI(facts, "other_servings", other)
			w11aI(facts, "given_cups", answer)
			return w11aP(o,
				fmt.Sprintf("A recipe for %d servings uses %d cups of flour. How many cups are needed to make a smaller batch of %d servings? (enter a number)", servings, cups, other),
				strconv.Itoa(answer),
				fmt.Sprintf("Each serving needs %d / %d = %d cups, so %d servings need %d x %d = %d cups.", cups, servings, unit, other, other, unit, answer),
				facts)
		default:
			other := 1 + rand.Intn(30)
			given := cups * other / servings
			w11aI(facts, "other_servings", other)
			w11aI(facts, "given_cups", given)
			return w11aP(o,
				fmt.Sprintf("A recipe for %d servings uses %d cups of flour. How many servings can you make with %d cups? (enter a number)", servings, cups, given),
				strconv.Itoa(other),
				fmt.Sprintf("Each serving needs %d cups, and %d / %d = %d servings.", unit, given, unit, other),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 11. Mixture in a given ratio
// ---------------------------------------------------------------------------

type w11aMixtureGen struct{}

func (g *w11aMixtureGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		a := 1 + rand.Intn(10)
		b := 1 + rand.Intn(10)
		for a == b {
			b = 1 + rand.Intn(10)
		}
		unit := 2 + rand.Intn(12)
		partA := a * unit
		partB := b * unit
		total := partA + partB
		facts := map[string]string{}
		w11aI(facts, "part_a", partA)
		w11aI(facts, "part_b", partB)
		w11aI(facts, "total", total)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("Paint is mixed in the ratio %d:%d to make %d litres. How many litres are in the first part? (enter a number)", a, b, total),
				strconv.Itoa(partA),
				fmt.Sprintf("There are %d + %d = %d parts, so one part is %d / %d = %d litres. The first part is %d x %d = %d litres.", a, b, a+b, total, a+b, unit, a, unit, partA),
				facts)
		case 2:
			return w11aP(o,
				fmt.Sprintf("Paint is mixed in the ratio %d:%d. The first part is %d litres. How many litres are in the second part? (enter a number)", a, b, partA),
				strconv.Itoa(partB),
				fmt.Sprintf("One part is %d / %d = %d litres, so the second part is %d x %d = %d litres.", partA, a, unit, b, unit, partB),
				facts)
		default:
			return w11aP(o,
				fmt.Sprintf("Paint is mixed in the ratio %d:%d. The first part is %d litres. What is the total number of litres? (enter a number)", a, b, partA),
				strconv.Itoa(total),
				fmt.Sprintf("One part is %d / %d = %d litres. The total is %d + %d = %d parts, or %d x %d = %d litres.", partA, a, unit, a, b, a+b, a+b, unit, total),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 12. Distance, rate and time
// ---------------------------------------------------------------------------

type w11aSpeedGen struct{}

func (g *w11aSpeedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		rate := 2 + rand.Intn(29)
		time := 2 + rand.Intn(19)
		for rate == time {
			time++
		}
		distance := rate * time
		facts := map[string]string{}
		w11aI(facts, "rate", rate)
		w11aI(facts, "time", time)
		w11aI(facts, "distance", distance)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("A car travels at %d miles per hour for %d hours. How far does it travel in miles? (enter a number)", rate, time),
				strconv.Itoa(distance),
				fmt.Sprintf("Multiply the rate by the time: %d x %d = %d miles.", rate, time, distance),
				facts)
		case 2:
			return w11aP(o,
				fmt.Sprintf("A car travels %d miles at a steady %d miles per hour. How many hours does the trip take? (enter a number)", distance, rate),
				strconv.Itoa(time),
				fmt.Sprintf("Divide the distance by the rate: %d / %d = %d hours.", distance, rate, time),
				facts)
		default:
			return w11aP(o,
				fmt.Sprintf("A car travels %d miles in %d hours. What is its speed in miles per hour? (enter a number)", distance, time),
				strconv.Itoa(rate),
				fmt.Sprintf("Divide the distance by the time: %d / %d = %d miles per hour.", distance, time, rate),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 13. Density, mass and volume
// ---------------------------------------------------------------------------

type w11aDensityGen struct{}

func (g *w11aDensityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		density := 2 + rand.Intn(19)
		volume := 2 + rand.Intn(29)
		for density == volume {
			volume++
		}
		mass := density * volume
		facts := map[string]string{}
		w11aI(facts, "density", density)
		w11aI(facts, "volume", volume)
		w11aI(facts, "mass", mass)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("A metal block has a density of %d grams per cubic centimetre and a volume of %d cubic centimetres. What is its mass in grams? (enter a number)", density, volume),
				strconv.Itoa(mass),
				fmt.Sprintf("Multiply the density by the volume: %d x %d = %d grams.", density, volume, mass),
				facts)
		case 2:
			return w11aP(o,
				fmt.Sprintf("A metal block has a mass of %d grams and a density of %d grams per cubic centimetre. What is its volume in cubic centimetres? (enter a number)", mass, density),
				strconv.Itoa(volume),
				fmt.Sprintf("Divide the mass by the density: %d / %d = %d cubic centimetres.", mass, density, volume),
				facts)
		default:
			return w11aP(o,
				fmt.Sprintf("A metal block has a mass of %d grams and a volume of %d cubic centimetres. What is its density in grams per cubic centimetre? (enter a number)", mass, volume),
				strconv.Itoa(density),
				fmt.Sprintf("Divide the mass by the volume: %d / %d = %d grams per cubic centimetre.", mass, volume, density),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 14. Substituting a negative value
// ---------------------------------------------------------------------------

type w11aSubstituteNegativeGen struct{}

func (g *w11aSubstituteNegativeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		k := 1 + rand.Intn(20)
		x := -k
		var a, b int
		switch o {
		case 1:
			a = 2 + rand.Intn(12)
			b = 1 + rand.Intn(15)
		case 2:
			a = 2 + rand.Intn(12)
			b = -(1 + rand.Intn(15))
		default:
			a = -(2 + rand.Intn(12))
			b = 1 + rand.Intn(15)
		}
		answer := a*x + b
		facts := map[string]string{}
		w11aI(facts, "a", a)
		w11aI(facts, "b", b)
		facts["x"] = strconv.Itoa(x)
		return w11aP(o,
			fmt.Sprintf("Evaluate %s when x = %d. (enter a number)", w11aLinearExpr(a, b), x),
			strconv.Itoa(answer),
			fmt.Sprintf("Substitute x = %d: %d(%d) + %d = %d.", x, a, x, b, answer),
			facts)
	})
}

// ---------------------------------------------------------------------------
// 15. Greatest common factor
// ---------------------------------------------------------------------------

type w11aFactorCommonGen struct{}

func (g *w11aFactorCommonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		gcf := w11aFactorGCFs.Next()
		x := 2 + rand.Intn(6)
		y := 2 + rand.Intn(6)
		for y == x || w11aGCD(x, y) != 1 {
			y = 2 + rand.Intn(6)
		}
		a := gcf * x
		b := gcf * y
		facts := map[string]string{}
		w11aI(facts, "a", a)
		w11aI(facts, "b", b)
		w11aI(facts, "gcf", gcf)
		switch o {
		case 1:
			return w11aP(o,
				fmt.Sprintf("What is the greatest common factor of %d and %d? (enter a number)", a, b),
				strconv.Itoa(gcf),
				fmt.Sprintf("Both numbers are multiples of %d, and no larger factor is common. The greatest common factor is %d.", gcf, gcf),
				facts)
		case 2:
			return w11aP(o,
				fmt.Sprintf("The terms %dx and %dx have a common numerical factor. What is the greatest common numerical factor? (enter a number)", a, b),
				strconv.Itoa(gcf),
				fmt.Sprintf("Factor out the shared coefficient: %dx + %dx = %d(%dx + %dx). The greatest common numerical factor is %d.", a, b, gcf, x, y, gcf),
				facts)
		default:
			third := a + b
			return w11aP(o,
				fmt.Sprintf("What is the greatest common factor of the three terms %d, %d and %d? (enter a number)", a, b, third),
				strconv.Itoa(gcf),
				fmt.Sprintf("The first two terms have greatest common factor %d, and %d is also a multiple of %d, so the answer is %d.", gcf, third, gcf, gcf),
				facts)
		}
	})
}

// ---------------------------------------------------------------------------
// 16. Independent and dependent variables
// ---------------------------------------------------------------------------

type w11aIndependentDependentGen struct{}

func (g *w11aIndependentDependentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := w11aOrd(ctx)
	return w11aAvoidAnswer(func() generator.Problem {
		rate := 2 + rand.Intn(12)
		start := 1 + rand.Intn(12)
		facts := map[string]string{}
		w11aI(facts, "rate", rate)
		w11aI(facts, "start", start)
		switch o {
		case 1:
			input := 2 + rand.Intn(30)
			output := rate*input + start
			w11aI(facts, "input", input)
			w11aI(facts, "output", output)
			return w11aP(o,
				fmt.Sprintf("A rule multiplies the input by %d and then adds %d. What is the output for an input of %d? (enter a number)", rate, start, input),
				strconv.Itoa(output),
				fmt.Sprintf("Multiply %d by %d and add %d: %d + %d = %d.", rate, input, start, rate*input, start, output),
				facts)
		case 2:
			input := 2 + rand.Intn(30)
			output := rate*input + start
			w11aI(facts, "input", input)
			w11aI(facts, "output", output)
			return w11aP(o,
				fmt.Sprintf("A rule multiplies the input by %d and then adds %d. What input gives an output of %d? (enter a number)", rate, start, output),
				strconv.Itoa(input),
				fmt.Sprintf("Subtract %d from the output and divide by %d: (%d - %d) / %d = %d.", start, rate, output, start, rate, input),
				facts)
		default:
			first := 1 + rand.Intn(10)
			input := first + 3
			output := rate*input + start
			w11aI(facts, "input", input)
			w11aI(facts, "output", output)
			return w11aP(o,
				fmt.Sprintf("A rule multiplies the input by %d and then adds %d. A table shows input %d with output %d, input %d with output %d and input %d with output %d. What is the output for input %d? (enter a number)", rate, start, first, rate*first+start, first+1, rate*(first+1)+start, first+2, rate*(first+2)+start, input),
				strconv.Itoa(output),
				fmt.Sprintf("The pattern adds %d for each increase of 1 in the input, so the output for %d is %d x %d + %d = %d.", rate, input, rate, input, start, output),
				facts)
		}
	})
}
