package prealgebra

// The 2026-10-18 prealgebra batch: fourteen concepts that take the domain from
// 46 to the target of 60.
//
// The batch fills the gaps a prealgebra course leaves when it stops at
// two-step equations: percent change and its applications, the proportion
// method for a missing whole, unit-price comparison, translating words into
// expressions, evaluating formulas, inequality word problems and number-line
// graphs, and the three measurement conversions. Every concept has at least
// one prerequisite already in the corpus, and several form short chains
// (percent change -> percent error, translate -> inequality words) rather than
// standing alone.
//
// Every answer is computed from integer quantities and published as an exact
// string. Money stays in cents until it is rendered, because a dollar amount
// built with floating point can print as 7.030000000000001 and then fail the
// numeric grader. Every branch of a generator publishes the same fact keys, so
// a schema cannot resolve on one shape and silently fall back on another.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("pct.ops.change", &pctChangeGen{})
	reg.Register("pct.ops.error", &pctErrorGen{})
	reg.Register("pct.ops.commission", &pctCommissionGen{})
	reg.Register("pct.ops.markup", &pctMarkupGen{})
	reg.Register("pct.ops.simple_interest", &pctSimpleInterestGen{})
	reg.Register("ratio.ops.percent_proportion", &pctProportionGen{})
	reg.Register("ratio.ops.better_buy", &betterBuyGen{})
	reg.Register("prealg.expr.translate", &translateExprGen{})
	reg.Register("prealg.expr.formula", &formulaEvalGen{})
	reg.Register("prealg.ineq.word", &ineqWordGen{})
	reg.Register("prealg.ineq.graph", &ineqGraphGen{})
	reg.Register("prealg.measure.convert", &customaryConvertGen{})
	reg.Register("prealg.measure.metric", &metricConvertGen{})
	reg.Register("prealg.measure.temperature", &temperatureConvertGen{})
}

// ---------------------------------------------------------------------------
// 1. Percent change from an old value and a new value
// ---------------------------------------------------------------------------

type pctChangeGen struct{}

func (g *pctChangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	old := (rand.Intn(10) + 1) * 20   // 20..200
	percent := (rand.Intn(9) + 1) * 5 // 5..45
	change := old * percent / 100
	direction := "increase"
	newValue := old + change
	if rand.Intn(2) == 0 {
		direction = "decrease"
		newValue = old - change
	}
	return generator.Problem{
		Question: fmt.Sprintf("A quantity changes from %d to %d. What is the percent %s? Enter a positive percent.",
			old, newValue, direction),
		Answer: strconv.Itoa(percent),
		Explanation: fmt.Sprintf("The amount of change is %d. Divide by the original value and multiply by 100: %d / %d x 100 = %d%%.",
			change, change, old, percent),
		Facts: map[string]string{
			"old": strconv.Itoa(old), "new": strconv.Itoa(newValue),
			"change": strconv.Itoa(change), "percent": strconv.Itoa(percent),
			"direction": direction,
			"formula":   "percent change = |new - old| / old x 100",
		},
	}
}

// ---------------------------------------------------------------------------
// 2. Percent error
// ---------------------------------------------------------------------------

type pctErrorGen struct{}

func (g *pctErrorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	actual := (rand.Intn(10) + 1) * 20 // 20..200
	percent := (rand.Intn(5) + 1) * 5  // 5..25
	difference := actual * percent / 100
	measured := actual + difference
	if rand.Intn(2) == 0 {
		measured = actual - difference
	}
	return generator.Problem{
		Question: fmt.Sprintf("A measurement is %d and the accepted value is %d. What is the percent error? Enter a positive percent.",
			measured, actual),
		Answer: strconv.Itoa(percent),
		Explanation: fmt.Sprintf("The absolute difference is |%d - %d| = %d. Divide by the accepted value and multiply by 100: %d / %d x 100 = %d%%.",
			measured, actual, difference, difference, actual, percent),
		Facts: map[string]string{
			"measured": strconv.Itoa(measured), "actual": strconv.Itoa(actual),
			"difference": strconv.Itoa(difference), "percent": strconv.Itoa(percent),
			"formula": "percent error = |measured - actual| / actual x 100",
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Commission
// ---------------------------------------------------------------------------

type pctCommissionGen struct{}

func (g *pctCommissionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	salesDollars := (rand.Intn(20) + 5) * 100 // 500..2400
	rate := rand.Intn(8) + 3                  // 3..10
	sales := moneyCents(salesDollars * 100)
	commission := moneyCents(salesDollars * rate)
	return generator.Problem{
		Question: fmt.Sprintf("A salesperson earns a %d%% commission on $%s in sales. What is the commission?", rate, sales),
		Answer:   commission,
		Explanation: fmt.Sprintf("%d%% of $%s = $%s. Multiply the sales by the rate written as a decimal.",
			rate, sales, commission),
		Facts: map[string]string{
			"sales": sales, "rate": strconv.Itoa(rate), "commission": commission,
			"formula": "commission = rate x sales",
		},
	}
}

// ---------------------------------------------------------------------------
// 4. Markup and selling price
// ---------------------------------------------------------------------------

type pctMarkupGen struct{}

func (g *pctMarkupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	costDollars := (rand.Intn(20) + 5) * 10 // 50..240
	rate := rand.Intn(21) + 10              // 10..30
	markup := moneyCents(costDollars * rate)
	price := moneyCents(costDollars*100 + costDollars*rate)
	cost := moneyCents(costDollars * 100)
	return generator.Problem{
		Question: fmt.Sprintf("A store buys an item for $%s and marks it up %d%%. What is the selling price?", cost, rate),
		Answer:   price,
		Explanation: fmt.Sprintf("The markup is %d%% of $%s = $%s. Add the markup to the cost: $%s + $%s = $%s.",
			rate, cost, markup, cost, markup, price),
		Facts: map[string]string{
			"cost": cost, "rate": strconv.Itoa(rate), "markup": markup, "price": price,
			"formula": "selling price = cost + markup",
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Simple interest
// ---------------------------------------------------------------------------

type pctSimpleInterestGen struct{}

func (g *pctSimpleInterestGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	principalDollars := (rand.Intn(20) + 5) * 100 // 500..2400
	rate := rand.Intn(8) + 2                      // 2..9
	time := rand.Intn(5) + 1                      // 1..5
	interest := moneyCents(principalDollars * rate * time)
	principal := moneyCents(principalDollars * 100)
	timeUnit := "years"
	if time == 1 {
		timeUnit = "year"
	}
	return generator.Problem{
		Question: fmt.Sprintf("Find the simple interest on $%s at %d%% for %d %s.", principal, rate, time, timeUnit),
		Answer:   interest,
		Explanation: fmt.Sprintf("Use I = P x r x t. I = $%s x %d/100 x %d = $%s.",
			principal, rate, time, interest),
		Facts: map[string]string{
			"principal": principal, "rate": strconv.Itoa(rate), "time": strconv.Itoa(time),
			"time_unit": timeUnit, "interest": interest, "formula": "I = P x r x t",
		},
	}
}

// ---------------------------------------------------------------------------
// 6. Percent problems solved with a proportion
// ---------------------------------------------------------------------------

type pctProportionGen struct{}

func (g *pctProportionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	whole := (rand.Intn(10) + 1) * 20 // 20..200
	percent := (rand.Intn(9) + 1) * 5 // 5..45
	part := whole * percent / 100
	return generator.Problem{
		Question: fmt.Sprintf("%d is %d%% of what number?", part, percent),
		Answer:   strconv.Itoa(whole),
		Explanation: fmt.Sprintf("Write the proportion %d / w = %d / 100. Cross-multiply: 100 x %d = %d x w, so w = %d.",
			part, percent, part, percent, whole),
		Facts: map[string]string{
			"part": strconv.Itoa(part), "percent": strconv.Itoa(percent), "whole": strconv.Itoa(whole),
			"proportion": "part / whole = percent / 100",
		},
	}
}

// ---------------------------------------------------------------------------
// 7. Better buy by unit price
// ---------------------------------------------------------------------------

type betterBuyGen struct{}

func (g *betterBuyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return betterBuyProblem()
}

// ---------------------------------------------------------------------------
// 8. Translate words into an algebraic expression
// ---------------------------------------------------------------------------

type translateExprGen struct{}

func (g *translateExprGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		phrase string
		answer string
		reason string
	}
	entries := []entry{
		{"five more than three times a number x", "3x + 5", "Three times x is 3x, and five more means add 5."},
		{"the sum of twice a number x and seven", "2x + 7", "Twice x is 2x; the sum with seven is 2x + 7."},
		{"four less than the product of six and a number x", "6x - 4", "The product is 6x; four less means subtract 4."},
		{"the quotient of a number x and five", "x/5", "The quotient of x and five is x divided by 5."},
		{"eight times the difference of a number x and two", "8(x - 2)", "The difference is x - 2; eight times it is 8(x - 2)."},
		{"the difference of nine and twice a number x", "9 - 2x", "Twice x is 2x; the difference of nine and 2x is 9 - 2x."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{
		Question:    fmt.Sprintf("Write an expression for: %s. Use x as the variable.", e.phrase),
		Answer:      e.answer,
		Explanation: e.reason,
		Facts: map[string]string{
			"phrase": e.phrase, "answer": e.answer, "reason": e.reason,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. Evaluate a formula
// ---------------------------------------------------------------------------

type formulaEvalGen struct{}

func (g *formulaEvalGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(3) {
	case 0:
		length := rand.Intn(10) + 3 // 3..12
		width := rand.Intn(10) + 3  // 3..12
		result := 2*length + 2*width
		return generator.Problem{
			Question: fmt.Sprintf("Use the formula \\(P = 2l + 2w\\) to find \\(P\\) when \\(l = %d\\) and \\(w = %d\\).", length, width),
			Answer:   strconv.Itoa(result),
			Explanation: fmt.Sprintf("\\(P = 2(%d) + 2(%d) = %d + %d = %d\\).",
				length, width, 2*length, 2*width, result),
			Facts: map[string]string{
				"formula": "P = 2l + 2w", "quantity": "P",
				"v1_name": "l", "v1": strconv.Itoa(length),
				"v2_name": "w", "v2": strconv.Itoa(width),
				"substituted": fmt.Sprintf("2(%d) + 2(%d)", length, width),
				"result":      strconv.Itoa(result),
			},
		}
	case 1:
		base := (rand.Intn(9) + 2) * 2 // 4..20, even
		height := rand.Intn(10) + 3    // 3..12
		result := base * height / 2
		return generator.Problem{
			Question: fmt.Sprintf("Use the formula \\(A = \\frac{bh}{2}\\) to find \\(A\\) when \\(b = %d\\) and \\(h = %d\\).", base, height),
			Answer:   strconv.Itoa(result),
			Explanation: fmt.Sprintf("\\(A = \\frac{(%d)(%d)}{2} = \\frac{%d}{2} = %d\\).",
				base, height, base*height, result),
			Facts: map[string]string{
				"formula": "A = bh/2", "quantity": "A",
				"v1_name": "b", "v1": strconv.Itoa(base),
				"v2_name": "h", "v2": strconv.Itoa(height),
				"substituted": fmt.Sprintf("(%d)(%d)/2", base, height),
				"result":      strconv.Itoa(result),
			},
		}
	default:
		rate := (rand.Intn(5) + 4) * 10 // 40..80
		time := rand.Intn(5) + 2        // 2..6
		result := rate * time
		return generator.Problem{
			Question: fmt.Sprintf("Use the formula \\(d = rt\\) to find \\(d\\) when \\(r = %d\\) and \\(t = %d\\).", rate, time),
			Answer:   strconv.Itoa(result),
			Explanation: fmt.Sprintf("\\(d = (%d)(%d) = %d\\).",
				rate, time, result),
			Facts: map[string]string{
				"formula": "d = rt", "quantity": "d",
				"v1_name": "r", "v1": strconv.Itoa(rate),
				"v2_name": "t", "v2": strconv.Itoa(time),
				"substituted": fmt.Sprintf("(%d)(%d)", rate, time),
				"result":      strconv.Itoa(result),
			},
		}
	}
}

// ---------------------------------------------------------------------------
// 10. Inequality word problems
// ---------------------------------------------------------------------------

type ineqWordGen struct{}

type ineqWordContext struct {
	context  string
	relation string
	ask      string
	build    func(fixed, rate, limit int) string
}

func (g *ineqWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	contexts := []ineqWordContext{
		{
			context: "taxi fare", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A taxi charges $%d plus $%d per mile. You have at most $%d. What is the greatest whole number of miles you can travel?", fixed, rate, limit)
			},
		},
		{
			context: "elevator capacity", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("An elevator has a limit of %d kg. It already carries %d kg of equipment, and each person weighs %d kg. What is the greatest whole number of people who can ride?", limit, fixed, rate)
			},
		},
		{
			context: "fundraiser", relation: "at least", ask: "least",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A fundraiser has %d points and earns %d points for each item sold. The goal is at least %d points. What is the least whole number of items needed?", fixed, rate, limit)
			},
		},
		{
			context: "equipment rental", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A rental company charges a $%d fee plus $%d per hour. You can spend at most $%d. What is the greatest whole number of hours you can rent?", fixed, rate, limit)
			},
		},
		{
			context: "printing order", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A print shop charges a $%d setup fee plus $%d per page. Your budget is at most $%d. What is the greatest whole number of pages you can print?", fixed, rate, limit)
			},
		},
		{
			context: "streaming plan", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A streaming service charges $%d to join plus $%d per month. You can spend at most $%d in total. What is the greatest whole number of months you can subscribe?", fixed, rate, limit)
			},
		},
		{
			context: "savings goal", relation: "at least", ask: "least",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("You have saved $%d and save $%d each week. Your goal is at least $%d. What is the least whole number of weeks needed?", fixed, rate, limit)
			},
		},
		{
			context: "catering order", relation: "at most", ask: "greatest",
			build: func(fixed, rate, limit int) string {
				return fmt.Sprintf("A caterer charges a $%d delivery fee plus $%d per guest. You can spend at most $%d. What is the greatest whole number of guests you can invite?", fixed, rate, limit)
			},
		},
	}
	c := contexts[rand.Intn(len(contexts))]
	bound := rand.Intn(8) + 3 // 3..10
	rate := rand.Intn(6) + 2  // 2..7
	fixed := rand.Intn(20) + 5
	limit := rate*bound + fixed
	symbol := "\\le"
	if c.relation == "at least" {
		symbol = "\\ge"
	}
	inequality := fmt.Sprintf("%dx + %d %s %d", rate, fixed, symbol, limit)
	solveStep := fmt.Sprintf("%dx %s %d", rate, symbol, limit-fixed)
	division := fmt.Sprintf("x %s %d", symbol, bound)
	return generator.Problem{
		Question: c.build(fixed, rate, limit),
		Answer:   strconv.Itoa(bound),
		Explanation: fmt.Sprintf("Let x be the unknown. The inequality is \\(%s\\). Subtract %d: \\(%s\\). Divide by %d: \\(%s\\).",
			inequality, fixed, solveStep, rate, division),
		Facts: map[string]string{
			"context": c.context, "fixed": strconv.Itoa(fixed), "rate": strconv.Itoa(rate),
			"limit": strconv.Itoa(limit), "relation": c.relation, "ask": c.ask,
			"bound": strconv.Itoa(bound), "answer": strconv.Itoa(bound),
			"inequality": inequality, "solve_step": solveStep, "division": division,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. Graph an inequality on a number line
// ---------------------------------------------------------------------------

type ineqGraphGen struct{}

func (g *ineqGraphGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return inequalityIntegerProblem()
}

// ---------------------------------------------------------------------------
// 12. Customary unit conversion
// ---------------------------------------------------------------------------

type customaryConvertGen struct{}

func (g *customaryConvertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type conversion struct {
		from, to string
		factor   int
	}
	conversions := []conversion{
		{"feet", "inches", 12},
		{"yards", "feet", 3},
		{"miles", "feet", 5280},
		{"pounds", "ounces", 16},
		{"gallons", "quarts", 4},
		{"quarts", "pints", 2},
		{"hours", "minutes", 60},
		{"days", "hours", 24},
	}
	c := conversions[rand.Intn(len(conversions))]
	tenths := rand.Intn(36) + 5 // 0.5..4.0
	if tenths == 10 {           // one unit takes a singular noun; skip it
		tenths = 15
	}
	return generator.Problem{
		Question: fmt.Sprintf("Convert %s %s to %s.", tenthsText(tenths), c.from, c.to),
		Answer:   tenthsText(tenths * c.factor),
		Explanation: fmt.Sprintf("Each unit of %s equals %d %s, so multiply %s x %d = %s %s.",
			c.from, c.factor, c.to, tenthsText(tenths), c.factor, tenthsText(tenths*c.factor), c.to),
		Facts: map[string]string{
			"value": tenthsText(tenths), "from_unit": c.from, "to_unit": c.to,
			"factor": strconv.Itoa(c.factor), "answer": tenthsText(tenths * c.factor),
			"operation": "multiply by the conversion factor",
		},
	}
}

// ---------------------------------------------------------------------------
// 13. Metric unit conversion
// ---------------------------------------------------------------------------

type metricConvertGen struct{}

func (g *metricConvertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type conversion struct {
		from, to string
		factor   int
		power    int
	}
	conversions := []conversion{
		{"kilometres", "metres", 1000, 3},
		{"metres", "centimetres", 100, 2},
		{"centimetres", "millimetres", 10, 1},
		{"kilograms", "grams", 1000, 3},
		{"grams", "milligrams", 1000, 3},
		{"litres", "millilitres", 1000, 3},
	}
	c := conversions[rand.Intn(len(conversions))]
	tenths := rand.Intn(36) + 5 // 0.5..4.0
	if tenths == 10 {
		tenths = 15
	}
	return generator.Problem{
		Question: fmt.Sprintf("Convert %s %s to %s.", tenthsText(tenths), c.from, c.to),
		Answer:   tenthsText(tenths * c.factor),
		Explanation: fmt.Sprintf("The metric prefix changes the unit by 10 to the power %d, so move the decimal point %d places: %s %s = %s %s.",
			c.power, c.power, tenthsText(tenths), c.from, tenthsText(tenths*c.factor), c.to),
		Facts: map[string]string{
			"value": tenthsText(tenths), "from_unit": c.from, "to_unit": c.to,
			"factor": strconv.Itoa(c.factor), "power": strconv.Itoa(c.power),
			"answer": tenthsText(tenths * c.factor),
		},
	}
}

// ---------------------------------------------------------------------------
// 14. Celsius and Fahrenheit conversion
// ---------------------------------------------------------------------------

type temperatureConvertGen struct{}

func (g *temperatureConvertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	celsius := (rand.Intn(13) - 6) * 5 // -30..30
	fahrenheit := 9*celsius/5 + 32
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question: fmt.Sprintf("Convert %d°C to degrees Fahrenheit.", celsius),
			Answer:   strconv.Itoa(fahrenheit),
			Explanation: fmt.Sprintf("Use F = 9C/5 + 32. F = 9(%d)/5 + 32 = %d.",
				celsius, fahrenheit),
			Facts: map[string]string{
				"value": strconv.Itoa(celsius), "from_unit": "Celsius", "to_unit": "Fahrenheit",
				"formula": "F = 9C/5 + 32", "substituted": fmt.Sprintf("9(%d)/5 + 32", celsius),
				"answer": strconv.Itoa(fahrenheit), "direction": "CtoF",
			},
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("Convert %d°F to degrees Celsius.", fahrenheit),
		Answer:   strconv.Itoa(celsius),
		Explanation: fmt.Sprintf("Use C = 5(F - 32)/9. C = 5(%d - 32)/9 = %d.",
			fahrenheit, celsius),
		Facts: map[string]string{
			"value": strconv.Itoa(fahrenheit), "from_unit": "Fahrenheit", "to_unit": "Celsius",
			"formula": "C = 5(F - 32)/9", "substituted": fmt.Sprintf("5(%d - 32)/9", fahrenheit),
			"answer": strconv.Itoa(celsius), "direction": "FtoC",
		},
	}
}

// tenthsText renders a whole number of tenths without binary-float noise.
// 5 -> "0.5", 15 -> "1.5", 20 -> "2".
func tenthsText(tenths int) string {
	if tenths%10 == 0 {
		return strconv.Itoa(tenths / 10)
	}
	return fmt.Sprintf("%d.%d", tenths/10, tenths%10)
}
