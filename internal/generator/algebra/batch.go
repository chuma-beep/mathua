package algebra

// The 2026-10-18 batch: fifteen concepts that take algebra from 89 to the
// domain target of 104.
//
// Every answer is computed from the instance and graded through the router with
// the type the corpus ships: `expression`, `tuple`, `numeric` or
// `multiple_choice`. The three rational-expression concepts publish the values
// their cancellation excludes, because cancelling a factor changes the domain;
// the explanation states the restriction rather than dropping it.
//
// The two factoring concepts ask for one factor rather than the whole product.
// `expression` answers are compared by SymPy equivalence, so an expanded answer
// to "factor completely" would also be accepted; asking for the factor that is
// hard to produce keeps the answer unambiguous.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("alg.linear.two_points", &twoPointsGen{})
	reg.Register("alg.linear.intercepts", &interceptsGen{})
	reg.Register("alg.factor.grouping", &groupingGen{})
	reg.Register("alg.factor.cubes", &cubesGen{})
	reg.Register("alg.quad.vertex_form", &vertexFormGen{})
	reg.Register("alg.quad.word_problem", &quadWordGen{})
	reg.Register("alg.poly.end_behavior", &endBehaviorGen{})
	reg.Register("alg.rational.simplify", &rationalSimplifyGen{})
	reg.Register("alg.rational.mult_div", &rationalMultDivGen{})
	reg.Register("alg.rational.add_sub", &rationalAddSubGen{})
	reg.Register("alg.radical.rational_exponents", &rationalExponentsGen{})
	reg.Register("alg.radical.rationalize", &rationalizeGen{})
	reg.Register("alg.func.variation", &variationGen{})
	reg.Register("alg.log.change_base", &changeBaseGen{})
	reg.Register("alg.conic.classify", &conicClassifyGen{})
}

// ---------------------------------------------------------------------------
// Formatting
// ---------------------------------------------------------------------------

type polyTerm struct{ coef, deg int }

// formatPoly renders terms in the order given, dropping zero coefficients and
// omitting a coefficient of one in front of x.
func formatPoly(terms ...polyTerm) string {
	var b strings.Builder
	first := true
	for _, t := range terms {
		if t.coef == 0 {
			continue
		}
		mag := t.coef
		if mag < 0 {
			mag = -mag
		}
		switch {
		case first && t.coef < 0:
			b.WriteString("-")
		case !first && t.coef < 0:
			b.WriteString(" - ")
		case !first:
			b.WriteString(" + ")
		}
		if t.deg == 0 {
			b.WriteString(strconv.Itoa(mag))
		} else {
			if mag != 1 {
				b.WriteString(strconv.Itoa(mag))
			}
			b.WriteString("x")
			if t.deg > 1 {
				fmt.Fprintf(&b, "^%d", t.deg)
			}
		}
		first = false
	}
	if first {
		return "0"
	}
	return b.String()
}

// formatBinom renders the inside of a binomial factor: "x + 3" / "x - 3".
func formatBinom(k int) string {
	if k < 0 {
		return fmt.Sprintf("x - %d", -k)
	}
	return fmt.Sprintf("x + %d", k)
}

// formatTwoVar renders "cx x + cy y = c" with signs.
func formatTwoVar(cx, cy, c int) string {
	var b strings.Builder
	b.WriteString(formatPoly(polyTerm{cx, 1}))
	mag := cy
	if mag < 0 {
		mag = -mag
	}
	y := "y"
	if mag != 1 {
		y = fmt.Sprintf("%dy", mag)
	}
	if cy < 0 {
		b.WriteString(" - " + y)
	} else {
		b.WriteString(" + " + y)
	}
	fmt.Fprintf(&b, " = %d", c)
	return b.String()
}

// ---------------------------------------------------------------------------
// 1. Line through two points
// ---------------------------------------------------------------------------

type twoPointsGen struct{}

func (g *twoPointsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := int(4 + ctx.Difficulty*8)
	x1 := rand.Intn(span) + 1
	dx := rand.Intn(5) + 1
	x2 := x1 + dx
	m := rand.Intn(7) - 3 // -3..3
	if m == 0 {
		m = 2
	}
	y1 := rand.Intn(2*span+1) - span
	y2 := y1 + m*dx
	b := y1 - m*x1
	rhs := formatPoly(polyTerm{m, 1}, polyTerm{b, 0})
	return generator.Problem{
		Question: fmt.Sprintf("A line passes through \\((%d, %d)\\) and \\((%d, %d)\\). "+
			"Write its equation in slope-intercept form and enter the right-hand side, the part after \\(y = \\).", x1, y1, x2, y2),
		Answer: rhs,
		Explanation: fmt.Sprintf("The slope is \\(m = \\frac{%d - %d}{%d - %d} = \\frac{%d}{%d} = %d\\). "+
			"Using \\((%d, %d)\\), \\(b = %d - (%d)(%d) = %d\\), so \\(y = %s\\).",
			y2, y1, x2, x1, y2-y1, x2-x1, m, x1, y1, y1, m, x1, b, rhs),
		Facts: map[string]string{
			"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1),
			"x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2),
			"dy": strconv.Itoa(y2 - y1), "dx": strconv.Itoa(x2 - x1),
			"m": strconv.Itoa(m), "b": strconv.Itoa(b), "rhs": rhs,
		},
	}
}

// ---------------------------------------------------------------------------
// 2. Intercepts
// ---------------------------------------------------------------------------

type interceptsGen struct{}

func (g *interceptsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	lim := int(5 + ctx.Difficulty*4)
	a := rand.Intn(lim) + 2 // x-intercept
	b := rand.Intn(lim) + 2 // y-intercept
	if a == b {
		b++
	}
	if rand.Intn(2) == 0 {
		a = -a
	}
	if rand.Intn(2) == 0 {
		b = -b
	}
	coefX := b
	coefY := a
	rhs := a * b
	eq := formatTwoVar(coefX, coefY, rhs)
	return generator.Problem{
		Question: fmt.Sprintf("Find the \\(x\\)- and \\(y\\)-intercepts of \\(%s\\). "+
			"Enter them as \\((x\\text{-intercept}, y\\text{-intercept})\\).", eq),
		Answer: fmt.Sprintf("%d,%d", a, b),
		Explanation: fmt.Sprintf("Set \\(y = 0\\): \\(%s = %d\\), so \\(x = %d\\). "+
			"Set \\(x = 0\\): \\(%s = %d\\), so \\(y = %d\\). The intercepts are \\((%d, 0)\\) and \\((0, %d)\\).",
			formatPoly(polyTerm{coefX, 1}), rhs, a, formatPoly(polyTerm{coefY, 1}), rhs, b, a, b),
		Facts: map[string]string{
			"eq": eq, "coef_x": strconv.Itoa(coefX), "coef_y": strconv.Itoa(coefY),
			"rhs": strconv.Itoa(rhs), "x_int": strconv.Itoa(a), "y_int": strconv.Itoa(b),
		},
	}
}

// ---------------------------------------------------------------------------
// 3. Factor by grouping
// ---------------------------------------------------------------------------

type groupingGen struct{}

func (g *groupingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(4) + 1 // 1..4
	k := rand.Intn(9) - 4 // -4..4
	if k == 0 {
		k = 2
	}
	c := rand.Intn(8) + 1 // positive: a x^2 + c has no real linear factor
	poly := formatPoly(polyTerm{a, 3}, polyTerm{a * k, 2}, polyTerm{c, 1}, polyTerm{c * k, 0})
	binomial := formatBinom(k)
	quadratic := formatPoly(polyTerm{a, 2}, polyTerm{c, 0})
	first := fmt.Sprintf("%s(%s)", formatPoly(polyTerm{a, 2}), binomial)
	second := fmt.Sprintf("%d(%s)", c, binomial)
	return generator.Problem{
		Question: fmt.Sprintf("Factor by grouping: \\(%s\\). What binomial factor do the two groups share?", poly),
		Answer:   binomial,
		Explanation: fmt.Sprintf("The first two terms give \\(%s\\), and the last two give \\(%s\\). "+
			"Both groups contain \\(%s\\), so the factorization is \\((%s)(%s)\\).",
			first, second, binomial, binomial, quadratic),
		Facts: map[string]string{
			"poly": poly, "a": strconv.Itoa(a), "k": strconv.Itoa(k), "c": strconv.Itoa(c),
			"binomial": binomial, "quadratic": quadratic, "first_group": first, "second_group": second,
		},
	}
}

// ---------------------------------------------------------------------------
// 4. Sum and difference of cubes
// ---------------------------------------------------------------------------

type cubesGen struct{}

func (g *cubesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := rand.Intn(11) + 2 // 2..12, so the cube stays inside a readable integer
	n3 := n * n * n
	sum := rand.Intn(2) == 0
	var kind, poly, linear, quadratic string
	if sum {
		kind = "sum"
		poly = formatPoly(polyTerm{1, 3}, polyTerm{n3, 0})
		linear = formatBinom(n)
		quadratic = formatPoly(polyTerm{1, 2}, polyTerm{-n, 1}, polyTerm{n * n, 0})
	} else {
		kind = "difference"
		poly = formatPoly(polyTerm{1, 3}, polyTerm{-n3, 0})
		linear = formatBinom(-n)
		quadratic = formatPoly(polyTerm{1, 2}, polyTerm{n, 1}, polyTerm{n * n, 0})
	}
	return generator.Problem{
		Question: fmt.Sprintf("Factor the %s of cubes \\(%s\\). Enter the quadratic factor.", kind, poly),
		Answer:   quadratic,
		Explanation: fmt.Sprintf("The %s of cubes factors as \\((%s)(%s)\\). "+
			"The quadratic factor is \\(%s\\).", kind, linear, quadratic, quadratic),
		Facts: map[string]string{
			"n": strconv.Itoa(n), "cube": strconv.Itoa(n3), "kind": kind, "poly": poly,
			"linear": linear, "quadratic": quadratic,
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Vertex form
// ---------------------------------------------------------------------------

type vertexFormGen struct{}

func (g *vertexFormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	h := rand.Intn(13) - 6 // -6..6
	if h == 0 {
		h = 2
	}
	k := rand.Intn(21) - 10 // -10..10
	b := -2 * h
	c := h*h + k
	quad := formatPoly(polyTerm{1, 2}, polyTerm{b, 1}, polyTerm{c, 0})
	vertex := formatVertex(h, k)
	square := fmt.Sprintf("%s = %s", formatVertex(h, 0),
		formatPoly(polyTerm{1, 2}, polyTerm{b, 1}, polyTerm{h * h, 0}))
	return generator.Problem{
		Question: fmt.Sprintf("Write \\(y = %s\\) in vertex form \\(y = (x - h)^2 + k\\). "+
			"Enter \\(h\\) and \\(k\\) as a comma-separated pair.", quad),
		Answer: fmt.Sprintf("%d,%d", h, k),
		Explanation: fmt.Sprintf("Half of the \\(x\\)-coefficient is %d, so the square is \\(%s\\). "+
			"The square contributes %d as its constant and the polynomial's constant is %d, "+
			"so the adjustment is %d. Thus \\(y = %s\\), and \\((h,k) = (%d,%d)\\).",
			b/2, square, h*h, c, k, vertex, h, k),
		Facts: map[string]string{
			"quad": quad, "b": strconv.Itoa(b), "c": strconv.Itoa(c),
			"h": strconv.Itoa(h), "k": strconv.Itoa(k), "half": strconv.Itoa(b / 2),
			"vertex_form": vertex, "square": formatVertex(h, 0),
		},
	}
}

// formatVertex renders "(x - 3)^2 + 5", "(x + 2)^2 - 4" and "(x - 3)^2".
func formatVertex(h, k int) string {
	var b strings.Builder
	if h >= 0 {
		fmt.Fprintf(&b, "(x - %d)^2", h)
	} else {
		fmt.Fprintf(&b, "(x + %d)^2", -h)
	}
	if k > 0 {
		fmt.Fprintf(&b, " + %d", k)
	}
	if k < 0 {
		fmt.Fprintf(&b, " - %d", -k)
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// 6. Quadratic word problem
// ---------------------------------------------------------------------------

type quadWordGen struct{}

func (g *quadWordGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	w := rand.Intn(10) + 2 // 2..11
	m := rand.Intn(3) + 1  // 1..3
	d := rand.Intn(9) + 1  // 1..9
	length := m*w + d
	area := w * length
	lengthExpr := formatPoly(polyTerm{m, 1}, polyTerm{d, 0})
	quadratic := formatPoly(polyTerm{m, 2}, polyTerm{d, 1}, polyTerm{-area, 0})
	var wording string
	if m == 1 {
		wording = fmt.Sprintf("A rectangle's length is %d m more than its width.", d)
	} else {
		wording = fmt.Sprintf("A rectangle's length is %d m more than %d times its width.", d, m)
	}
	return generator.Problem{
		Question: fmt.Sprintf("%s Its area is %d m². Find its width in meters.", wording, area),
		Answer:   strconv.Itoa(w),
		Explanation: fmt.Sprintf("With width \\(w\\), the length is \\(%s\\), so \\(w(%s) = %d\\), "+
			"that is \\(%s = 0\\). The positive root is \\(w = %d\\), and the other root is negative.",
			lengthExpr, lengthExpr, area, quadratic, w),
		Facts: map[string]string{
			"w": strconv.Itoa(w), "m": strconv.Itoa(m), "d": strconv.Itoa(d),
			"length": strconv.Itoa(length), "area": strconv.Itoa(area),
			"length_expr": lengthExpr, "quadratic": quadratic,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. End behavior
// ---------------------------------------------------------------------------

type endBehaviorGen struct{}

func (g *endBehaviorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	degree := []int{2, 3, 4, 5}[rand.Intn(4)]
	lead := rand.Intn(8) - 4 // -4..3
	if lead == 0 {
		lead = 2
	}
	terms := []polyTerm{{lead, degree}}
	for d := degree - 1; d >= 0; d-- {
		terms = append(terms, polyTerm{rand.Intn(7) - 3, d})
	}
	poly := formatPoly(terms...)
	even := degree%2 == 0
	var answer, parity string
	if even {
		parity = "even"
		if lead > 0 {
			answer = "up on both ends"
		} else {
			answer = "down on both ends"
		}
	} else {
		parity = "odd"
		if lead > 0 {
			answer = "down on the left, up on the right"
		} else {
			answer = "up on the left, down on the right"
		}
	}
	p := generator.MultipleChoice(
		fmt.Sprintf("Describe the end behavior of \\(f(x) = %s\\).", poly),
		answer,
		"up on both ends",
		"down on both ends",
		"down on the left, up on the right",
		"up on the left, down on the right",
	)
	p.Facts = map[string]string{
		"poly": poly, "degree": strconv.Itoa(degree), "leading": strconv.Itoa(lead),
		"parity": parity, "behavior": answer,
	}
	p.Explanation = fmt.Sprintf("The degree is %d (%s) and the leading coefficient is %d, so the end behavior is %s.",
		degree, parity, lead, answer)
	return p
}

// ---------------------------------------------------------------------------
// 8. Simplify a rational expression
// ---------------------------------------------------------------------------

type rationalSimplifyGen struct{}

func (g *rationalSimplifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	var numerator, denominator, factored, cancelled, remaining, validFor string
	if rand.Intn(2) == 0 {
		a := rand.Intn(9) - 4 // -4..4
		b := rand.Intn(9) - 4
		if a == 0 {
			a = 2
		}
		if b == 0 {
			b = 3
		}
		if a == b {
			b++
		}
		numerator = formatPoly(polyTerm{1, 2}, polyTerm{a + b, 1}, polyTerm{a * b, 0})
		denominator = formatBinom(b)
		factored = fmt.Sprintf("(%s)(%s)", formatBinom(a), formatBinom(b))
		cancelled = formatBinom(b)
		remaining = formatBinom(a)
		validFor = fmt.Sprintf("x \\ne %d", -b)
	} else {
		a := rand.Intn(8) + 2 // 2..9
		numerator = formatPoly(polyTerm{1, 2}, polyTerm{-a * a, 0})
		denominator = formatBinom(-a)
		factored = fmt.Sprintf("(%s)(%s)", formatBinom(-a), formatBinom(a))
		cancelled = formatBinom(-a)
		remaining = formatBinom(a)
		validFor = fmt.Sprintf("x \\ne %d", a)
	}
	return generator.Problem{
		Question: fmt.Sprintf("Simplify \\(\\frac{%s}{%s}\\).", numerator, denominator),
		Answer:   remaining,
		Explanation: fmt.Sprintf("Factor the numerator as \\(%s\\) and cancel the common factor \\(%s\\), "+
			"leaving \\(%s\\). The cancellation requires \\(%s\\).",
			factored, cancelled, remaining, validFor),
		Facts: map[string]string{
			"numerator": numerator, "denominator": denominator, "factored": factored,
			"cancelled": cancelled, "remaining": remaining, "valid_for": validFor,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. Multiply and divide rational expressions
// ---------------------------------------------------------------------------

type rationalMultDivGen struct{}

func (g *rationalMultDivGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(7) + 2 // 2..8
	b := rand.Intn(9) + 1 // 1..9
	if b == a {
		b++
	}
	numFactored := fmt.Sprintf("(%s)(%s)", formatBinom(-a), formatBinom(a))
	num1 := formatPoly(polyTerm{1, 2}, polyTerm{-a * a, 0})
	den1 := formatBinom(b)
	cancelled := fmt.Sprintf("%s and %s", formatBinom(b), formatBinom(-a))
	validFor := fmt.Sprintf("x \\ne %d,\\ x \\ne %d", -b, a)
	remaining := formatBinom(a)
	var num2, den2, question, operation, approach string
	if rand.Intn(2) == 0 {
		operation = "multiplying"
		approach = "Multiply the numerators and the denominators."
		num2 = formatBinom(b)
		den2 = formatBinom(-a)
		question = fmt.Sprintf("Multiply and simplify: \\(\\frac{%s}{%s} \\cdot \\frac{%s}{%s}\\).",
			num1, den1, num2, den2)
	} else {
		operation = "dividing"
		approach = "Invert the divisor, then multiply the numerators and the denominators."
		num2 = formatBinom(-a)
		den2 = formatBinom(b)
		question = fmt.Sprintf("Divide and simplify: \\(\\frac{%s}{%s} \\div \\frac{%s}{%s}\\).",
			num1, den1, num2, den2)
	}
	first := fmt.Sprintf("\\frac{%s}{%s}", num1, den1)
	second := fmt.Sprintf("\\frac{%s}{%s}", num2, den2)
	return generator.Problem{
		Question: question,
		Answer:   remaining,
		Explanation: fmt.Sprintf("%s Factor \\(%s\\) as \\(%s\\), then cancel \\(%s\\). "+
			"The result is \\(%s\\), valid for \\(%s\\).",
			approach, num1, numFactored, cancelled, remaining, validFor),
		Facts: map[string]string{
			"first": first, "second": second, "operation": operation, "approach": approach,
			"num1": num1, "den1": den1, "num2": num2, "den2": den2,
			"num_factored": numFactored, "cancelled": cancelled,
			"remaining": remaining, "valid_for": validFor,
		},
	}
}

// ---------------------------------------------------------------------------
// 10. Add and subtract rational expressions
// ---------------------------------------------------------------------------

type rationalAddSubGen struct{}

func (g *rationalAddSubGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := rand.Intn(6) + 2 // 2..7
	m := rand.Intn(4) + 1 // 1..4
	n := rand.Intn(4) + 1 // 1..4
	plus := rand.Intn(2) == 0
	var cx, c int
	var operation, sign, question string
	if plus {
		cx = m + n
		c = a * (m - n)
		operation = "adding"
		sign = "+"
		question = fmt.Sprintf("Add and simplify: \\(\\frac{%d}{%s} + \\frac{%d}{%s}\\).",
			m, formatBinom(-a), n, formatBinom(a))
	} else {
		cx = m - n
		c = a * (m + n)
		operation = "subtracting"
		sign = "-"
		question = fmt.Sprintf("Subtract and simplify: \\(\\frac{%d}{%s} - \\frac{%d}{%s}\\).",
			m, formatBinom(-a), n, formatBinom(a))
	}
	numerator := formatPoly(polyTerm{cx, 1}, polyTerm{c, 0})
	denominator := formatPoly(polyTerm{1, 2}, polyTerm{-a * a, 0})
	result := fmt.Sprintf("(%s)/(%s)", numerator, denominator)
	validFor := fmt.Sprintf("x \\ne %d,\\ x \\ne %d", a, -a)
	expandedLeft := fmt.Sprintf("%d(%s)", m, formatBinom(a))
	expandedRight := fmt.Sprintf("%d(%s)", n, formatBinom(-a))
	return generator.Problem{
		Question: question,
		Answer:   result,
		Explanation: fmt.Sprintf("The common denominator is \\(%s\\). Over it the numerators are "+
			"\\(%s\\) and \\(%s\\), so the combined numerator is %s. "+
			"The result is \\(%s\\), valid for \\(%s\\).",
			denominator, expandedLeft, expandedRight, numerator, result, validFor),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "m": strconv.Itoa(m), "n": strconv.Itoa(n),
			"operation": operation, "sign": sign,
			"left": formatFraction(m, formatBinom(-a)), "right": formatFraction(n, formatBinom(a)),
			"expanded_left": expandedLeft, "expanded_right": expandedRight,
			"numerator": numerator, "denominator": denominator, "result": result,
			"valid_for": validFor,
		},
	}
}

func formatFraction(num int, den string) string {
	return fmt.Sprintf("\\frac{%d}{%s}", num, den)
}

// ---------------------------------------------------------------------------
// 11. Rational exponents
// ---------------------------------------------------------------------------

type rationalExponentsGen struct{}

func (g *rationalExponentsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	pairs := [][2]int{{2, 1}, {3, 1}, {3, 2}, {4, 1}, {4, 3}, {5, 2}, {5, 3}}
	p := pairs[rand.Intn(len(pairs))]
	index, power := p[0], p[1]
	root := rand.Intn(4) + 2 // 2..5
	base := mathutil.IntPow(root, index)
	magnitude := mathutil.IntPow(root, power)
	negative := rand.Intn(2) == 0
	var exponent, value string
	if negative {
		exponent = fmt.Sprintf("-%d/%d", power, index)
		value = fracOrInt(1, magnitude)
	} else {
		exponent = fmt.Sprintf("%d/%d", power, index)
		value = strconv.Itoa(magnitude)
	}
	return generator.Problem{
		Question: fmt.Sprintf("Evaluate \\(%d^{%s}\\).", base, exponent),
		Answer:   value,
		Explanation: fmt.Sprintf("Since \\(%d = %d^{%d}\\), \\(%d^{%d/%d} = %d^{%d} = %d\\). "+
			"With the exponent %s, the value is %s.",
			base, root, index, base, power, index, root, power, magnitude, exponent, value),
		Facts: map[string]string{
			"base": strconv.Itoa(base), "index": strconv.Itoa(index), "power": strconv.Itoa(power),
			"root": strconv.Itoa(root), "exponent": exponent,
			"magnitude": strconv.Itoa(magnitude), "value": value,
		},
	}
}

// ---------------------------------------------------------------------------
// 12. Rationalize a denominator
// ---------------------------------------------------------------------------

type rationalizeGen struct{}

func (g *rationalizeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	squarefree := []int{2, 3, 5, 6, 7, 10, 11, 13, 14, 15}
	b := squarefree[rand.Intn(len(squarefree))]
	a := rand.Intn(5) + 2 // 2..6
	k := a * b
	return generator.Problem{
		Question: fmt.Sprintf("Rationalize the denominator of \\(%d/\\sqrt{%d}\\). "+
			"Write the result as \\(a\\sqrt{b}\\) with \\(b\\) squarefree. "+
			"Enter \\(a\\) and \\(b\\) as a comma-separated pair.", k, b),
		Answer: fmt.Sprintf("%d,%d", a, b),
		Explanation: fmt.Sprintf("Multiply numerator and denominator by \\(\\sqrt{%d}\\): "+
			"\\(\\frac{%d\\sqrt{%d}}{%d} = %d\\sqrt{%d}\\). So \\(a = %d\\) and \\(b = %d\\).",
			b, k, b, b, a, b, a, b),
		Facts: map[string]string{
			"k": strconv.Itoa(k), "radicand": strconv.Itoa(b), "a": strconv.Itoa(a),
		},
	}
}

// ---------------------------------------------------------------------------
// 13. Direct and inverse variation
// ---------------------------------------------------------------------------

type variationGen struct{}

func (g *variationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	x1 := rand.Intn(5) + 2 // 2..6
	x2 := rand.Intn(5) + 2 // 2..6
	if x1 == x2 {
		x2++
	}
	if rand.Intn(2) == 0 {
		k := rand.Intn(8) + 2 // 2..9
		y1 := k * x1
		y2 := k * x2
		return generator.Problem{
			Question: fmt.Sprintf("\\(y\\) varies directly as \\(x\\). \\(y = %d\\) when \\(x = %d\\). "+
				"Find \\(y\\) when \\(x = %d\\).", y1, x1, x2),
			Answer: strconv.Itoa(y2),
			Explanation: fmt.Sprintf("The constant of variation is \\(k = \\frac{%d}{%d} = %d\\). "+
				"Then \\(y = %d(%d) = %d\\).", y1, x1, k, k, x2, y2),
			Facts: map[string]string{
				"kind": "directly", "constant": strconv.Itoa(k),
				"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2),
				"value": strconv.Itoa(y2),
			},
		}
	}
	t := rand.Intn(4) + 1 // 1..4
	k := t * x1 * x2
	y1 := k / x1
	y2 := k / x2
	return generator.Problem{
		Question: fmt.Sprintf("\\(y\\) varies inversely as \\(x\\). \\(y = %d\\) when \\(x = %d\\). "+
			"Find \\(y\\) when \\(x = %d\\).", y1, x1, x2),
		Answer: strconv.Itoa(y2),
		Explanation: fmt.Sprintf("The constant is \\(k = (%d)(%d) = %d\\). "+
			"Then \\(y = \\frac{%d}{%d} = %d\\).", x1, y1, k, k, x2, y2),
		Facts: map[string]string{
			"kind": "inversely", "constant": strconv.Itoa(k),
			"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2),
			"value": strconv.Itoa(y2),
		},
	}
}

// ---------------------------------------------------------------------------
// 14. Change of base
// ---------------------------------------------------------------------------

type changeBaseGen struct{}

func (g *changeBaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := []int{2, 3, 5}[rand.Intn(3)]
	m := []int{2, 3, 4}[rand.Intn(3)]
	var candidates []int
	for i := 1; i <= 3; i++ {
		if i%m != 0 { // a whole-number result is a plain log evaluation, not a change of base
			candidates = append(candidates, i)
		}
	}
	n := candidates[rand.Intn(len(candidates))]
	base := mathutil.IntPow(a, m)
	arg := mathutil.IntPow(a, n)
	value := fracOrInt(n, m)
	return generator.Problem{
		Question: fmt.Sprintf("Use change of base to evaluate \\(\\log_{%d}(%d)\\).", base, arg),
		Answer:   value,
		Explanation: fmt.Sprintf("Both numbers are powers of %d: \\(%d = %d^{%d}\\) and \\(%d = %d^{%d}\\). "+
			"Changing the base gives \\(\\frac{\\log(%d)}{\\log(%d)} = \\frac{%d}{%d} = %s\\).",
			a, base, a, m, arg, a, n, arg, base, n, m, value),
		Facts: map[string]string{
			"base": strconv.Itoa(base), "argument": strconv.Itoa(arg), "a": strconv.Itoa(a),
			"m": strconv.Itoa(m), "n": strconv.Itoa(n), "value": value,
		},
	}
}

// ---------------------------------------------------------------------------
// 15. Classify a conic
// ---------------------------------------------------------------------------

type conicClassifyGen struct{}

func (g *conicClassifyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	var equation, class, reason string
	switch rand.Intn(4) {
	case 0:
		r := rand.Intn(7) + 2 // 2..8
		equation = fmt.Sprintf("x^2 + y^2 = %d", r*r)
		class = "circle"
		reason = "Both variables are squared with the same coefficient and there is no linear term, so the conic is a circle."
	case 1:
		a := rand.Intn(3) + 1
		if rand.Intn(2) == 0 {
			a = -a
		}
		b := rand.Intn(7) - 3
		c := rand.Intn(11) - 5
		equation = "y = " + formatPoly(polyTerm{a, 2}, polyTerm{b, 1}, polyTerm{c, 0})
		class = "parabola"
		reason = "Only one variable is squared, so the conic is a parabola."
	case 2:
		p := rand.Intn(7) + 2 // 2..8
		q := rand.Intn(7) + 2
		if p == q {
			q++
		}
		equation = fmt.Sprintf("x^2/%d + y^2/%d = 1", p*p, q*q)
		class = "ellipse"
		reason = "Both variables are squared with the same sign but unequal coefficients, so the conic is an ellipse."
	default:
		p := rand.Intn(7) + 2
		q := rand.Intn(7) + 2
		equation = fmt.Sprintf("x^2/%d - y^2/%d = 1", p*p, q*q)
		class = "hyperbola"
		reason = "The squared terms have opposite signs, so the conic is a hyperbola."
	}
	p := generator.MultipleChoice(
		fmt.Sprintf("Classify the conic \\(%s\\).", equation),
		class, "circle", "ellipse", "parabola", "hyperbola",
	)
	p.Facts = map[string]string{"equation": equation, "class": class, "reason": reason}
	p.Explanation = reason
	return p
}
