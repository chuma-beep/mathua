package arithmetic

// The 2026-10-18 arithmetic batch: fourteen concepts that take arithmetic from
// 46 to the domain target of 60.
//
// Every answer is computed from the instance and graded through the router with
// the type the corpus ships: `numeric`, `tuple`, `multiple_choice` or
// `expression`. The tuple concepts (ordering, expanded form, square-root
// bounds) publish bare integers, because the tuple grader compares each part as
// a string and "4.0" is not "4". The expression concepts publish the variable
// form, because SymPy equivalence accepts any equivalent unsimplified answer;
// asking for the variable form keeps the answer from collapsing to a plain
// number.
//
// Each generator publishes every fact its schema references on every branch.
// A fact that exists only when a coin flip lands one way would drop the schema
// at runtime, so the branches share one key set.

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("arith.compare.whole", &compareWholeGen{})
	reg.Register("arith.place.expanded", &expandedFormGen{})
	reg.Register("arith.place.powers_ten", &powersTenGen{})
	reg.Register("arith.estimate.sum", &estimateSumGen{})
	reg.Register("arith.estimate.product", &estimateProductGen{})
	reg.Register("arith.word.multi_step", &multiStepWordGen{})
	reg.Register("arith.div.interpret_remainder", &interpretRemainderGen{})
	reg.Register("arith.div.zero", &divZeroGen{})
	reg.Register("arith.neg.word", &negWordGen{})
	reg.Register("arith.exp.negative_base", &negativeBaseGen{})
	reg.Register("arith.exp.power_product", &powerProductGen{})
	reg.Register("arith.exp.power_quotient", &powerQuotientGen{})
	reg.Register("arith.sqrt.cube", &cubeRootGen{})
	reg.Register("arith.sqrt.estimate", &squareRootBoundsGen{})
}

// ---------------------------------------------------------------------------
// 1. Compare and order whole numbers
// ---------------------------------------------------------------------------

type compareWholeGen struct{}

func (g *compareWholeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(9900) + 100
	b := rand.Intn(9900) + 100
	for b == a {
		b = rand.Intn(9900) + 100
	}
	c := rand.Intn(9900) + 100
	for c == a || c == b {
		c = rand.Intn(9900) + 100
	}
	nums := []int{a, b, c}
	ascending := func() bool { return nums[0] < nums[1] && nums[1] < nums[2] }
	for tries := 0; tries < 10 && ascending(); tries++ {
		rand.Shuffle(3, func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	}
	if ascending() {
		nums[0], nums[1] = nums[1], nums[0]
	}

	ordered := append([]int(nil), nums...)
	sort.Ints(ordered)
	parts := make([]string, 3)
	for i, n := range ordered {
		parts[i] = strconv.Itoa(n)
	}
	answer := strings.Join(parts, ",")
	display := fmt.Sprintf("%d; %d; %d", nums[0], nums[1], nums[2])
	return generator.Problem{
		Question:    fmt.Sprintf("Order these numbers from least to greatest: %s. Enter them separated by commas.", display),
		Answer:      answer,
		Explanation: fmt.Sprintf("Compare the values place by place; the order is %s.", strings.ReplaceAll(answer, ",", ", ")),
		Facts:       map[string]string{"numbers": display, "ordered": answer},
	}
}

// ---------------------------------------------------------------------------
// 2. Expanded form of whole numbers
// ---------------------------------------------------------------------------

type expandedFormGen struct{}

func expandedParts(n int) ([]int, []string) {
	var values []int
	var texts []string
	for place := 1000; place >= 1; place /= 10 {
		digit := (n / place) % 10
		if digit == 0 {
			continue
		}
		value := digit * place
		values = append(values, value)
		texts = append(texts, strconv.Itoa(value))
	}
	return values, texts
}

func (g *expandedFormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := rand.Intn(9000) + 1000
	var texts []string
	for tries := 0; tries < 30; tries++ {
		_, texts = expandedParts(n)
		if len(texts) >= 2 {
			break
		}
		n = rand.Intn(9000) + 1000
	}
	if len(texts) < 2 {
		n = 4305
		_, texts = expandedParts(n)
	}
	answer := strings.Join(texts, ",")
	return generator.Problem{
		Question: fmt.Sprintf("Write %d as a sum of its nonzero place values, largest first. "+
			"Skip zero digits. Enter the values separated by commas, without thousands separators.", n),
		Answer:      answer,
		Explanation: fmt.Sprintf("The nonzero place values are %s, and their sum is %d.", strings.Join(texts, " + "), n),
		Facts:       map[string]string{"n": strconv.Itoa(n), "parts": answer, "expanded": strings.Join(texts, " + ")},
	}
}

// ---------------------------------------------------------------------------
// 3. Multiply and divide by powers of ten
// ---------------------------------------------------------------------------

type powersTenGen struct{}

func (g *powersTenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxN := 99
	if ctx.Difficulty >= 0.5 {
		maxN = 999
	}
	maxPower := 2
	if ctx.Difficulty >= 0.5 {
		maxPower = 3
	}
	power := rand.Intn(maxPower) + 1
	scale := mathutil.IntPow(10, power)
	if rand.Intn(2) == 0 {
		n := rand.Intn(maxN-1) + 2
		value := n * scale
		return generator.Problem{
			Question:    fmt.Sprintf("\\(%d \\times %d = ?\\)", n, scale),
			Answer:      strconv.Itoa(value),
			Explanation: fmt.Sprintf("\\(%d \\times %d = %d\\). Multiplying by %d shifts each digit %d place%s left.", n, scale, value, scale, power, plural(power)),
			Facts:       map[string]string{"n": strconv.Itoa(n), "power": strconv.Itoa(power), "operation": "multiply", "value": strconv.Itoa(value), "shift": fmt.Sprintf("%d place%s", power, plural(power))},
		}
	}
	value := rand.Intn(maxN-1) + 2
	n := value * scale
	return generator.Problem{
		Question:    fmt.Sprintf("\\(%d \\div %d = ?\\)", n, scale),
		Answer:      strconv.Itoa(value),
		Explanation: fmt.Sprintf("\\(%d \\div %d = %d\\). Dividing by %d shifts each digit %d place%s right.", n, scale, value, scale, power, plural(power)),
		Facts:       map[string]string{"n": strconv.Itoa(n), "power": strconv.Itoa(power), "operation": "divide", "value": strconv.Itoa(value), "shift": fmt.Sprintf("%d place%s", power, plural(power))},
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// ---------------------------------------------------------------------------
// 4. Estimate sums and differences
// ---------------------------------------------------------------------------

func roundTo(n, unit int) int {
	return int(math.Round(float64(n)/float64(unit))) * unit
}

type estimateSumGen struct{}

func (g *estimateSumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	draw := func() int {
		for {
			n := rand.Intn(900) + 100
			if n%100 != 50 { // 450 has two equally near hundreds; avoid the tie
				return n
			}
		}
	}
	op := "add"
	if rand.Intn(2) == 0 {
		op = "subtract"
	}
	a, b := draw(), draw()
	if op == "subtract" && a < b {
		a, b = b, a
	}
	ra, rb := roundTo(a, 100), roundTo(b, 100)
	if op == "subtract" {
		for ra == rb {
			b = draw()
			if a < b {
				a, b = b, a
			}
			ra, rb = roundTo(a, 100), roundTo(b, 100)
		}
	}
	estimate := ra + rb
	sign := "+"
	if op == "subtract" {
		estimate = ra - rb
		sign = "-"
	}
	return generator.Problem{
		Question: fmt.Sprintf("Estimate \\(%d %s %d\\) by rounding each number to the nearest hundred.", a, sign, b),
		Answer:   strconv.Itoa(estimate),
		Explanation: fmt.Sprintf("\\(%d \\approx %d\\) and \\(%d \\approx %d\\), so \\(%d %s %d = %d\\).",
			a, ra, b, rb, ra, sign, rb, estimate),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "operation": op,
			"rounded_a": strconv.Itoa(ra), "rounded_b": strconv.Itoa(rb),
			"estimate": strconv.Itoa(estimate),
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Estimate products by rounding
// ---------------------------------------------------------------------------

type estimateProductGen struct{}

func (g *estimateProductGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	draw := func() int {
		for {
			n := rand.Intn(90) + 10
			if n%10 != 5 { // 25 is equally near 20 and 30; avoid the tie
				return n
			}
		}
	}
	a, b := draw(), draw()
	ra, rb := roundTo(a, 10), roundTo(b, 10)
	estimate := ra * rb
	return generator.Problem{
		Question: fmt.Sprintf("Estimate \\(%d \\times %d\\) by rounding each factor to the nearest ten.", a, b),
		Answer:   strconv.Itoa(estimate),
		Explanation: fmt.Sprintf("\\(%d \\approx %d\\) and \\(%d \\approx %d\\), so \\(%d \\times %d = %d\\).",
			a, ra, b, rb, ra, rb, estimate),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b),
			"rounded_a": strconv.Itoa(ra), "rounded_b": strconv.Itoa(rb),
			"estimate": strconv.Itoa(estimate),
		},
	}
}

// ---------------------------------------------------------------------------
// 6. Multi-step whole-number word problems
// ---------------------------------------------------------------------------

type multiStepWordGen struct{}

type multiStepContext struct {
	place, item, itemOne, container, verb string
}

func (g *multiStepWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	contexts := []multiStepContext{
		{"A shop", "pencils", "pencil", "boxes", "sells"},
		{"A bakery", "muffins", "muffin", "trays", "sells"},
		{"A library", "books", "book", "shelves", "removes"},
	}
	c := contexts[rand.Intn(len(contexts))]
	a := rand.Intn(8) + 2  // 2..9
	b := rand.Intn(11) + 2 // 2..12
	total := a * b
	removed := rand.Intn(total-1) + 1 // 1..total-1, so the answer stays positive
	removedItem := c.item
	if removed == 1 {
		removedItem = c.itemOne
	}
	result := total - removed
	return generator.Problem{
		Question: fmt.Sprintf("%s has %d %s with %d %s each. It %s %d %s. How many %s are left?",
			c.place, a, c.container, b, c.item, c.verb, removed, removedItem, c.item),
		Answer: strconv.Itoa(result),
		Explanation: fmt.Sprintf("\\(%d \\times %d = %d\\) %s, then \\(%d - %d = %d\\).",
			a, b, total, c.item, total, removed, result),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(removed),
			"result": strconv.Itoa(result), "item": c.item, "container": c.container,
			"place": c.place,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. Interpret a remainder in context
// ---------------------------------------------------------------------------

type interpretRemainderGen struct{}

type remainderContext struct {
	item, itemOne, container, containerOne string
}

func (g *interpretRemainderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	contexts := []remainderContext{
		{"cookies", "cookie", "bags", "bag"},
		{"marbles", "marble", "jars", "jar"},
		{"pencils", "pencil", "boxes", "box"},
	}
	c := contexts[rand.Intn(len(contexts))]
	total := rand.Intn(80) + 20 // 20..99
	per := rand.Intn(7) + 3     // 3..9
	for total%per == 0 {
		total = rand.Intn(80) + 20
	}
	quotient := total / per
	remainder := total % per
	mode := "full"
	if ctx.KPOrdinal == 1 {
		mode = "full"
	} else if ctx.KPOrdinal >= 2 {
		mode = "needed"
	} else if rand.Intn(2) == 0 {
		mode = "needed"
	}
	answer := quotient
	if mode == "needed" {
		answer = quotient + 1
	}
	if mode == "full" {
		return generator.Problem{
			Question: fmt.Sprintf("You have %d %s. You put %d %s in each %s. How many full %s can you fill?",
				total, c.item, per, c.item, c.containerOne, c.container),
			Answer: strconv.Itoa(answer),
			Explanation: fmt.Sprintf("\\(%d \\times %d = %d\\) %s fill %d %s, with %d left over, so %d full %s.",
				per, quotient, per*quotient, c.item, quotient, c.container, remainder, answer, c.container),
			Facts: map[string]string{
				"total": strconv.Itoa(total), "per": strconv.Itoa(per),
				"quotient": strconv.Itoa(quotient), "remainder": strconv.Itoa(remainder),
				"mode": mode, "answer": strconv.Itoa(answer), "item": c.item, "container": c.container, "rule": remainderRule(mode),
			},
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("You have %d %s. Each %s holds %d %s. How many %s are needed to hold all of them?",
			total, c.item, c.containerOne, per, c.item, c.container),
		Answer: strconv.Itoa(answer),
		Explanation: fmt.Sprintf("\\(%d \\times %d = %d\\) %s fill %d %s, with %d left over, so one more %s is needed: %d.",
			per, quotient, per*quotient, c.item, quotient, c.container, remainder, c.containerOne, answer),
		Facts: map[string]string{
			"total": strconv.Itoa(total), "per": strconv.Itoa(per),
			"quotient": strconv.Itoa(quotient), "remainder": strconv.Itoa(remainder),
			"mode": mode, "answer": strconv.Itoa(answer), "item": c.item, "container": c.container, "rule": remainderRule(mode),
		},
	}
}

// ---------------------------------------------------------------------------
// 8. Division by zero
// ---------------------------------------------------------------------------

type divZeroGen struct{}

// Generate asks which statement about division by zero is true.
//
// The three knowledge points teach three separate claims: division by zero is
// undefined, a zero numerator over a nonzero denominator is zero, and a
// concrete expression has no value. Claims 1 and 2 keep the general rule as the
// answer, but the stem carries a concrete expression so the prompt does not
// repeat word for word. Claim 3 draws its own dividend and asks about that
// expression. The numbers are the part that varies; the correct claim is
// stable, which is what the knowledge point is assessing.
func (g *divZeroGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	claim := ctx.KPOrdinal
	if claim < 1 || claim > 3 {
		claim = rand.Intn(3) + 1
	}
	dividend := rand.Intn(98) + 2 // 2..99, so the quotients never coincide
	divisor := rand.Intn(98) + 2
	switch claim {
	case 1:
		correct := "Division by zero is undefined."
		p := generator.MultipleChoice(
			fmt.Sprintf("Consider the expression %d ÷ 0. Which statement is true?", dividend),
			correct,
			fmt.Sprintf("%d ÷ 0 equals 0.", dividend),
			fmt.Sprintf("%d ÷ 0 equals %d.", dividend, dividend),
			fmt.Sprintf("0 ÷ %d is undefined.", dividend),
		)
		p.Explanation = fmt.Sprintf("No number multiplied by zero gives %d, so %d ÷ 0 has no value. Division by zero is undefined.", dividend, dividend)
		p.Facts = map[string]string{"claim": "undefined", "dividend": strconv.Itoa(dividend), "answer": correct, "reason": p.Explanation}
		return p
	case 2:
		correct := "Zero divided by any nonzero number equals zero."
		p := generator.MultipleChoice(
			fmt.Sprintf("Consider the expression 0 ÷ %d. Which statement is true?", divisor),
			correct,
			fmt.Sprintf("0 ÷ %d is undefined.", divisor),
			fmt.Sprintf("0 ÷ %d equals %d.", divisor, divisor),
			fmt.Sprintf("%d ÷ 0 equals 0.", divisor),
		)
		p.Explanation = fmt.Sprintf("The equation %d × q = 0 has exactly one solution, q = 0, so 0 ÷ %d = 0.", divisor, divisor)
		p.Facts = map[string]string{"claim": "zero_numerator", "divisor": strconv.Itoa(divisor), "answer": correct, "reason": p.Explanation}
		return p
	}
	correct := fmt.Sprintf("The expression %d ÷ 0 has no value.", dividend)
	p := generator.MultipleChoice("Which statement is true?", correct,
		fmt.Sprintf("The expression %d ÷ 0 equals 0.", dividend),
		fmt.Sprintf("The expression %d ÷ 0 equals %d.", dividend, dividend),
		fmt.Sprintf("The expression 0 ÷ %d has no value.", dividend))
	p.Explanation = fmt.Sprintf("No number multiplied by zero gives %d, so %d ÷ 0 has no value.", dividend, dividend)
	p.Facts = map[string]string{
		"claim":    "concrete",
		"dividend": strconv.Itoa(dividend),
		"answer":   correct,
		"reason":   p.Explanation,
	}
	return p
}

// ---------------------------------------------------------------------------
// 9. Integer word problems
// ---------------------------------------------------------------------------

type negWordGen struct{}

func (g *negWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		start := rand.Intn(26) - 15 // -15..10
		delta := rand.Intn(12) + 1  // 1..12
		rise := rand.Intn(2) == 0
		change := delta
		direction := "rose"
		if !rise {
			change = -delta
			direction = "fell"
		}
		result := start + change
		return generator.Problem{
			Question:    fmt.Sprintf("The temperature was %d°C. It %s %d°C. What is the temperature now?", start, direction, delta),
			Answer:      strconv.Itoa(result),
			Explanation: fmt.Sprintf("\\(%d %s %d = %d\\).", start, signOf(change), abs(change), result),
			Facts: map[string]string{
				"start": strconv.Itoa(start), "change": strconv.Itoa(change),
				"result": strconv.Itoa(result), "unit": "°C",
				"direction": direction, "magnitude": strconv.Itoa(delta),
			},
		}
	}
	start := -(rand.Intn(20) + 1) // -20..-1
	rise := rand.Intn(2) == 0
	delta := rand.Intn(12) + 1
	if rise {
		delta = rand.Intn(-start) + 1 // rise no farther than the surface
	}
	change := delta
	direction := "rises"
	if !rise {
		change = -delta
		direction = "descends"
	}
	result := start + change
	return generator.Problem{
		Question:    fmt.Sprintf("A diver is at %d m relative to the surface. The diver %s %d m. What is the diver's position now?", start, direction, delta),
		Answer:      strconv.Itoa(result),
		Explanation: fmt.Sprintf("\\(%d %s %d = %d\\).", start, signOf(change), abs(change), result),
		Facts: map[string]string{
			"start": strconv.Itoa(start), "change": strconv.Itoa(change),
			"result": strconv.Itoa(result), "unit": "m",
			"direction": direction, "magnitude": strconv.Itoa(delta),
		},
	}
}

func signOf(n int) string {
	if n < 0 {
		return "-"
	}
	return "+"
}

func shapeOf(parenthesized bool) string {
	if parenthesized {
		return "parenthesized"
	}
	return "not parenthesized"
}

func remainderRule(mode string) string {
	if mode == "needed" {
		return "Every item must fit, so add one to the quotient."
	}
	return "Only full groups count, so keep the quotient."
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ---------------------------------------------------------------------------
// 10. Powers with negative bases
// ---------------------------------------------------------------------------

type negativeBaseGen struct{}

func (g *negativeBaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	base := rand.Intn(8) + 2  // 2..9
	power := rand.Intn(4) + 1 // 1..4, so both parities and both signs occur
	parenthesized := rand.Intn(2) == 0
	if ctx.KPOrdinal == 1 {
		parenthesized = false
	} else if ctx.KPOrdinal >= 2 {
		parenthesized = true
	}
	var value int
	var question, explanation, rule string
	if parenthesized {
		value = mathutil.IntPow(-base, power)
		question = fmt.Sprintf("Evaluate \\((-%d)^{%d}\\).", base, power)
		explanation = fmt.Sprintf("The parentheses put the negative sign inside the base, so \\((-%d)^{%d} = %d\\).", base, power, value)
		rule = "The parentheses put the negative sign inside the base, so the sign follows the exponent."
	} else {
		value = -mathutil.IntPow(base, power)
		question = fmt.Sprintf("Evaluate \\(-%d^{%d}\\).", base, power)
		explanation = fmt.Sprintf("Without parentheses the exponent applies to %d first, so \\(-%d^{%d} = -(%d^{%d}) = %d\\).", base, base, power, base, power, value)
		rule = "Without parentheses the exponent applies to the positive base first, and the negative sign stays outside."
	}
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(value),
		Explanation: explanation,
		Facts: map[string]string{
			"base": strconv.Itoa(base), "power": strconv.Itoa(power),
			"parenthesized": strconv.FormatBool(parenthesized), "shape": shapeOf(parenthesized),
			"value": strconv.Itoa(value), "rule": rule,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. Power of a product
// ---------------------------------------------------------------------------

type powerProductGen struct{}

func (g *powerProductGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	coefficient := rand.Intn(4) + 2 // 2..5
	power := rand.Intn(2) + 2       // 2..3
	coefficientPower := mathutil.IntPow(coefficient, power)
	answer := fmt.Sprintf("%dx^%d", coefficientPower, power)
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\((%dx)^{%d}\\)", coefficient, power),
		Answer:   answer,
		Explanation: fmt.Sprintf("\\((%dx)^{%d} = %d^{%d} x^{%d} = %s\\).",
			coefficient, power, coefficient, power, power, answer),
		Facts: map[string]string{
			"coefficient": strconv.Itoa(coefficient), "power": strconv.Itoa(power),
			"coefficient_power": strconv.Itoa(coefficientPower), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 12. Power of a quotient
// ---------------------------------------------------------------------------

type powerQuotientGen struct{}

func (g *powerQuotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	denominator := rand.Intn(4) + 2 // 2..5
	power := rand.Intn(2) + 2       // 2..3
	denominatorPower := mathutil.IntPow(denominator, power)
	answer := fmt.Sprintf("x^%d/%d", power, denominatorPower)
	return generator.Problem{
		Question: fmt.Sprintf("Simplify: \\(\\left(\\frac{x}{%d}\\right)^{%d}\\)", denominator, power),
		Answer:   answer,
		Explanation: fmt.Sprintf("\\(\\left(\\frac{x}{%d}\\right)^{%d} = \\frac{x^{%d}}{%d^{%d}} = \\frac{x^{%d}}{%d}\\).",
			denominator, power, power, denominator, power, power, denominatorPower),
		Facts: map[string]string{
			"denominator": strconv.Itoa(denominator), "power": strconv.Itoa(power),
			"denominator_power": strconv.Itoa(denominatorPower), "answer": answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 13. Cube roots of perfect cubes
// ---------------------------------------------------------------------------

type cubeRootGen struct{}

func (g *cubeRootGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	magnitude := rand.Intn(10) + 1 // 1..10
	root := magnitude
	if rand.Intn(2) == 0 {
		root = -magnitude
	}
	radicand := root * root * root
	return generator.Problem{
		Question:    fmt.Sprintf("Evaluate \\(\\sqrt[3]{%d}\\).", radicand),
		Answer:      strconv.Itoa(root),
		Explanation: fmt.Sprintf("\\(\\sqrt[3]{%d} = %d\\) because \\((%d)^{3} = %d\\).", radicand, root, root, radicand),
		Facts: map[string]string{
			"radicand": strconv.Itoa(radicand), "root": strconv.Itoa(root), "cube": strconv.Itoa(radicand),
		},
	}
}

// ---------------------------------------------------------------------------
// 14. Estimate square roots between integers
// ---------------------------------------------------------------------------

type squareRootBoundsGen struct{}

func (g *squareRootBoundsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := rand.Intn(149) + 2 // 2..150
	for isPerfectSquare(n) {
		n = rand.Intn(149) + 2
	}
	lower := int(math.Sqrt(float64(n)))
	upper := lower + 1
	return generator.Problem{
		Question: fmt.Sprintf("Between which two consecutive whole numbers does \\(\\sqrt{%d}\\) lie? Enter the smaller number first.", n),
		Answer:   fmt.Sprintf("%d,%d", lower, upper),
		Explanation: fmt.Sprintf("\\(%d^2 = %d\\) and \\(%d^2 = %d\\), so \\(%d < %d < %d\\) and \\(%d < \\sqrt{%d} < %d\\).",
			lower, lower*lower, upper, upper*upper, lower*lower, n, upper*upper, lower, n, upper),
		Facts: map[string]string{
			"radicand": strconv.Itoa(n), "lower": strconv.Itoa(lower), "upper": strconv.Itoa(upper),
			"lower_square": strconv.Itoa(lower * lower), "upper_square": strconv.Itoa(upper * upper),
		},
	}
}

func isPerfectSquare(n int) bool {
	r := int(math.Sqrt(float64(n)))
	return r*r == n
}
