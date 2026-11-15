package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func calcKP(ctx generator.GeneratorContext) int {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}
	return kp
}

// ---------------------------------------------------------------------------
// Limits by conjugates
// ---------------------------------------------------------------------------

type conjugateLimitGen struct{}

func (g *conjugateLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	m := rand.Intn(5) + 2 // 2..6
	a := rand.Intn(m*m-1) + 1
	k := m*m - a

	var question, explanation string
	var answer int
	base := 0
	switch kp {
	case 1:
		answer = m
		question = fmt.Sprintf("The limit \\(\\lim_{x→%d}(\\sqrt{x + %d} - %d)/(x - %d)\\) equals \\(1/(2m)\\). What is m? (enter a number)", a, k, m, a)
		explanation = fmt.Sprintf("Multiply by the conjugate. The limit is \\(1/(2\\sqrt{%d}) = 1/(2·%d)\\), so \\(m = %d\\).", m*m, m, m)
	case 2:
		answer = 2 * m
		question = fmt.Sprintf("Evaluate \\(\\lim_{x→%d}(\\sqrt{x + %d} - %d)/(x - %d)\\) in the form \\(1/k\\). What is k? (enter a number)", a, k, m, a)
		explanation = fmt.Sprintf("After multiplying by the conjugate, the limit is \\(1/(2·%d) = 1/%d\\), so \\(k = %d\\).", m, 2*m, 2*m)
	default:
		base = rand.Intn(5) + 2
		answer = 2 * base
		question = fmt.Sprintf("Evaluate \\(\\lim_{x→%d}(x - %d)/(\\sqrt{x} - %d)\\). (enter a number)", base*base, base*base, base)
		explanation = fmt.Sprintf("Factor \\(x - %d = (\\sqrt{x} - %d)(\\sqrt{x} + %d)\\). The quotient tends to \\(2·%d = %d\\).", base*base, base, base, base, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"m": strconv.Itoa(m), "a": strconv.Itoa(a), "base": strconv.Itoa(base)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Limits of complex fractions
// ---------------------------------------------------------------------------

type complexFractionLimitGen struct{}

func (g *complexFractionLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := randRange(ctx.Difficulty, 2, 20)

	var question, explanation string
	answer := a * a
	b := 0
	switch kp {
	case 1:
		if rand.Intn(2) == 0 {
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→%d}(1/x - 1/%d)/(x - %d)\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a, a)
			explanation = fmt.Sprintf("Combine the fractions: \\((1/x - 1/%d)/(x - %d) = -1/(%d x)\\), whose limit is \\(-1/%d\\).", a, a, a, a*a)
		} else {
			// A squared reciprocal has a steeper singular part, so the limit is
			// -2/a^3. The denominator a is even, which keeps k = a^3/2 whole.
			b = 2 * randRange(ctx.Difficulty, 1, 10)
			answer = b * b * b / 2
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→%d}(1/x^2 - 1/%d^2)/(x - %d)\\). The limit is \\(-1/k\\). What is k? (enter a number)", b, b, b)
			explanation = fmt.Sprintf("Factor \\(1/x^2 - 1/%d^2 = -(x - %d)(x + %d)/(%d^2x^2)\\). After cancelling \\(x - %d\\) the limit is \\(-2/%d^3 = -1/%d\\).", b, b, b, b, b, b, answer)
		}
	case 2:
		if rand.Intn(2) == 0 {
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(x + %d) - 1/%d)/x\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a)
			explanation = fmt.Sprintf("Combine the numerator over \\(%d(x + %d)\\), then divide by x. The limit is \\(-1/%d\\).", a, a, a*a)
		} else {
			b = 2 * randRange(ctx.Difficulty, 1, 10)
			answer = b * b * b / 2
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(x + %d)^2 - 1/%d^2)/x\\). The limit is \\(-1/k\\). What is k? (enter a number)", b, b)
			explanation = fmt.Sprintf("The numerator is \\(-%d(2x + %d)/(%d^2(x + %d)^2)\\). Dividing by x leaves \\(-2/%d^3 = -1/%d\\) in the limit.", b, b, b, b, b, answer)
		}
	default:
		if rand.Intn(2) == 0 {
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(x - %d) + 1/%d)/x\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a)
			explanation = fmt.Sprintf("The numerator is \\(x/(%d(x - %d))\\). Dividing by x leaves \\(1/(%d(x - %d))\\), whose limit is \\(-1/%d\\).", a, a, a, a, a*a)
		} else {
			// The mirror-image difference quotient is positive, so the sign in
			// the question changes but the denominator k is still a squared.
			question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(%d - x) - 1/%d)/x\\). The limit is \\(1/k\\). What is k? (enter a number)", a, a)
			explanation = fmt.Sprintf("The numerator is \\(x/(%d(%d - x))\\). Dividing by x leaves \\(1/(%d(%d - x))\\), whose limit is \\(1/%d\\).", a, a, a, a, a*a)
		}
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "answer": strconv.Itoa(answer)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Limits with absolute values
// ---------------------------------------------------------------------------

type absoluteValueLimitGen struct{}

func (g *absoluteValueLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := randRange(ctx.Difficulty, 1, 20)
	// The mirrored quotient |x-a|/(a-x) reverses which side is positive, so a
	// branch that fixes the side still asks about a different function.
	mirrored := rand.Intn(2) == 0

	var question, explanation string
	var answer int
	form := "quotient"
	switch kp {
	case 1:
		answer = -1
		if mirrored {
			answer = 1
		}
		form := "|x - %d|/(x - %d)"
		if mirrored {
			form = "|x - %d|/(%d - x)"
		}
		question = fmt.Sprintf("For \\(f(x) = "+form+"\\), what is the left-hand limit as \\(x → %d\\)? (enter a number)", a, a, a)
		if mirrored {
			explanation = fmt.Sprintf("For \\(x < %d\\), \\(|x - %d| = %d - x\\), so the quotient is 1 throughout the left side.", a, a, a)
		} else {
			explanation = fmt.Sprintf("For \\(x < %d\\), \\(|x - %d| = -(x - %d)\\), so the quotient is -1 throughout the left side.", a, a, a)
		}
	case 2:
		answer = 1
		if mirrored {
			answer = -1
		}
		form := "|x - %d|/(x - %d)"
		if mirrored {
			form = "|x - %d|/(%d - x)"
		}
		question = fmt.Sprintf("For \\(f(x) = "+form+"\\), what is the right-hand limit as \\(x → %d\\)? (enter a number)", a, a, a)
		if mirrored {
			explanation = fmt.Sprintf("For \\(x > %d\\), \\(|x - %d| = x - %d\\), so the quotient is -1 throughout the right side.", a, a, a)
		} else {
			explanation = fmt.Sprintf("For \\(x > %d\\), \\(|x - %d| = x - %d\\), so the quotient is 1 throughout the right side.", a, a, a)
		}
	default:
		switch rand.Intn(3) {
		case 0:
			form = "jump"
			question = fmt.Sprintf("For \\(f(x) = |x - %d|/(x - %d)\\), does the two-sided limit as \\(x → %d\\) exist? Enter 1 for yes and 0 for no.", a, a, a)
			explanation = "The left limit is -1 and the right limit is 1. They disagree, so the two-sided limit does not exist."
		case 1:
			// The absolute value itself is continuous, so this limit exists even
			// though the quotient above jumps.
			answer = 0
			form = "continuous"
			question = fmt.Sprintf("For \\(f(x) = |x - %d|\\), what is the two-sided limit as \\(x → %d\\)? (enter a number)", a, a)
			explanation = fmt.Sprintf("The absolute value is continuous at \\(x = %d\\), so the limit is \\(f(%d) = 0\\).", a, a)
		default:
			// Away from the corner the ratio of the two absolute values is
			// identically one, so the hole does not affect the limit.
			answer = 1
			form = "ratio"
			question = fmt.Sprintf("For \\(f(x) = |x - %d|/|x - %d|\\) with \\(x \\ne %d\\), what is the two-sided limit as \\(x → %d\\)? (enter a number)", a, a, a, a)
			explanation = "For every point other than the corner both absolute values agree, so the quotient is 1 and the limit is 1."
		}
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "mirrored": strconv.FormatBool(mirrored), "form": form},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Slopes of normal lines
// ---------------------------------------------------------------------------

type normalLineGen struct{}

func (g *normalLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := randRange(ctx.Difficulty, 1, 15)

	var question, explanation string
	answer := 2 * a
	fn := "x^2"
	switch kp {
	case 1:
		if rand.Intn(2) == 0 {
			question = fmt.Sprintf("For \\(f(x) = x^2\\), what is the slope of the tangent line at \\(x = %d\\)? (enter a number)", a)
			explanation = fmt.Sprintf("\\(f'(x) = 2x\\), so at \\(x = %d\\) the tangent slope is %d.", a, 2*a)
		} else {
			fn = "x^3"
			answer = 3 * a * a
			question = fmt.Sprintf("For \\(f(x) = x^3\\), what is the slope of the tangent line at \\(x = %d\\)? (enter a number)", a)
			explanation = fmt.Sprintf("\\(f'(x) = 3x^2\\), so at \\(x = %d\\) the tangent slope is \\(3·%d^2 = %d\\).", a, a, answer)
		}
	case 2:
		if rand.Intn(2) == 0 {
			question = fmt.Sprintf("For \\(f(x) = x^2\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
			explanation = fmt.Sprintf("The tangent slope is \\(2·%d = %d\\). The normal slope is the negative reciprocal, \\(-1/%d\\), so \\(k = %d\\).", a, 2*a, 2*a, 2*a)
		} else {
			fn = "x^2+x"
			answer = 2*a + 1
			question = fmt.Sprintf("For \\(f(x) = x^2 + x\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
			explanation = fmt.Sprintf("\\(f'(x) = 2x + 1\\), so the tangent slope at \\(x = %d\\) is \\(2·%d + 1 = %d\\). The normal slope is \\(-1/%d\\).", a, a, answer, answer)
		}
	default:
		if rand.Intn(2) == 0 {
			fn = "x^3"
			answer = 3 * a * a
			question = fmt.Sprintf("For \\(f(x) = x^3\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
			explanation = fmt.Sprintf("\\(f'(x) = 3x^2\\), so the tangent slope at \\(x = %d\\) is \\(3·%d^2 = %d\\). The normal slope is \\(-1/%d\\).", a, a, answer, answer)
		} else {
			fn = "x^4"
			answer = 4 * a * a * a
			question = fmt.Sprintf("For \\(f(x) = x^4\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
			explanation = fmt.Sprintf("\\(f'(x) = 4x^3\\), so the tangent slope at \\(x = %d\\) is \\(4·%d^3 = %d\\). The normal slope is \\(-1/%d\\).", a, a, answer, answer)
		}
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "function": fn},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Derivatives of piecewise functions
// ---------------------------------------------------------------------------

type piecewiseDerivativeGen struct{}

func (g *piecewiseDerivativeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	boundary := rand.Intn(4) + 1
	a := rand.Intn(5) + 1
	c := rand.Intn(5) + 1
	b := rand.Intn(9) - 4
	d := rand.Intn(9) - 4
	if kp == 3 {
		c = a
		d = (a-c)*boundary + b
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a
		question = fmt.Sprintf("For \\(f(x) = %s\\) when \\(x < %d\\), and \\(f(x) = %s\\) when \\(x ≥ %d\\), what is the left-hand derivative at \\(x = %d\\)? (enter a number)",
			formatCalcLine(a, b), boundary, formatCalcLine(c, d), boundary, boundary)
		explanation = fmt.Sprintf("The left branch has slope %d, so the left-hand derivative is %d.", a, a)
	case 2:
		answer = c
		question = fmt.Sprintf("For \\(f(x) = %s\\) when \\(x < %d\\), and \\(f(x) = %s\\) when \\(x ≥ %d\\), what is the right-hand derivative at \\(x = %d\\)? (enter a number)",
			formatCalcLine(a, b), boundary, formatCalcLine(c, d), boundary, boundary)
		explanation = fmt.Sprintf("The right branch has slope %d, so the right-hand derivative is %d.", c, c)
	case 3:
		answer = b
		question = fmt.Sprintf("For \\(f(x) = %s\\) when \\(x < %d\\), and \\(f(x) = %s + c\\) when \\(x ≥ %d\\), what value of c makes f differentiable at \\(x = %d\\)? (enter a number)",
			formatCalcLine(a, b), boundary, formatCalcLine(a, 0), boundary, boundary)
		explanation = fmt.Sprintf("At \\(x = %d\\) the branches must meet: \\(%d·%d + (%d) = %d·%d + c\\), so c = %d. Their slopes are already both %d.", boundary, a, boundary, b, a, boundary, b, a)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "boundary": strconv.Itoa(boundary),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Fence optimization
// ---------------------------------------------------------------------------

type fenceOptimizationGen struct{}

func (g *fenceOptimizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	k := randRange(ctx.Difficulty, 2, 25)
	perimeter := 4 * k
	// Two pen layouts share the same fence budget: three sides against a river,
	// or a divided pen whose river side and internal divider both need fence.
	divided := rand.Intn(2) == 0
	width := k
	length := 2 * k
	area := 2 * k * k
	layout := "A rectangular pen is built against a river, so only three sides need fence."
	if divided {
		length = k
		area = k * k
		layout = "A rectangular pen is built against a river with one internal fence parallel to the river, splitting it into two equal strips."
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = width
		question = fmt.Sprintf("%s The fence is %d meters long. What width perpendicular to the river maximizes the area? (enter a number)", layout, perimeter)
		explanation = fmt.Sprintf("If the two sides perpendicular to the river are x, the fenced parallel run is \\(%d - 2x\\). Area \\(A = x(%d - 2x)\\) is maximized at \\(x = %d/4 = %d\\).", perimeter, perimeter, perimeter, width)
	case 2:
		answer = length
		question = fmt.Sprintf("%s The fence is %d meters long. At the maximum area, what is the length of the side parallel to the river? (enter a number)", layout, perimeter)
		if divided {
			explanation = fmt.Sprintf("The optimal perpendicular width is \\(%d/4 = %d\\), and the parallel run is also %d because the two fenced parallel sides share the budget equally.", perimeter, width, length)
		} else {
			explanation = fmt.Sprintf("The optimal perpendicular width is \\(%d/4 = %d\\), so the parallel side is \\(%d - 2(%d) = %d\\).", perimeter, width, perimeter, width, length)
		}
	default:
		answer = area
		question = fmt.Sprintf("%s The fence is %d meters long. What is the maximum area? (enter a number)", layout, perimeter)
		explanation = fmt.Sprintf("The optimum has width %d and length %d, so the maximum area is \\(%d × %d = %d\\).", width, length, width, length, area)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"perimeter": strconv.Itoa(perimeter), "width": strconv.Itoa(width), "length": strconv.Itoa(length), "area": strconv.Itoa(area), "divided": strconv.FormatBool(divided)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Revenue optimization
// ---------------------------------------------------------------------------

type revenueOptimizationGen struct{}

func (g *revenueOptimizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	b := rand.Intn(3) + 1
	k := rand.Intn(5) + 2
	price := b * k
	intercept := 2 * b * b * k
	quantity := b * b * k
	revenue := b * b * b * k * k

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = price
		question = fmt.Sprintf("Demand is \\(q = %s\\), where p is the price. What price maximizes revenue? (enter a number)", demandText(intercept, b))
		explanation = fmt.Sprintf("Revenue is \\(R = p(%d - %dp)\\). Its derivative is zero at \\(p = %d/(2·%d) = %d\\).", intercept, b, intercept, b, price)
	case 2:
		answer = quantity
		question = fmt.Sprintf("Demand is \\(q = %s\\). At the revenue-maximizing price, what quantity is sold? (enter a number)", demandText(intercept, b))
		explanation = fmt.Sprintf("The optimal price is %d, so \\(q = %d - %d(%d) = %d\\).", price, intercept, b, price, quantity)
	default:
		answer = revenue
		question = fmt.Sprintf("Demand is \\(q = %s\\). What is the maximum revenue? (enter a number)", demandText(intercept, b))
		explanation = fmt.Sprintf("At the optimal price %d the quantity is %d, so revenue is \\(%d × %d = %d\\).", price, quantity, price, quantity, revenue)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"b": strconv.Itoa(b), "intercept": strconv.Itoa(intercept), "price": strconv.Itoa(price), "quantity": strconv.Itoa(quantity), "revenue": strconv.Itoa(revenue)},
		KPOrdinal:   kp,
	}
}

func demandText(intercept, b int) string {
	if b == 1 {
		return fmt.Sprintf("%d - p", intercept)
	}
	return fmt.Sprintf("%d - %dp", intercept, b)
}

// ---------------------------------------------------------------------------
// Related rates: sliding ladder
// ---------------------------------------------------------------------------

type ladderRateGen struct{}

func (g *ladderRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	type triple struct{ x, y, l int }
	triples := []triple{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {9, 12, 15}}
	t := triples[rand.Intn(len(triples))]
	// The primitive triple is scaled and the sliding rate is a multiple of the
	// base distance, which keeps every speed an integer.
	scale := randRange(ctx.Difficulty, 1, 3)
	factor := randRange(ctx.Difficulty, 1, 4)
	t.x *= scale
	t.y *= scale
	t.l *= scale

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		rate := t.x * factor
		answer = t.y * factor
		question = fmt.Sprintf("A %d-foot ladder leans against a wall with its base %d feet from the wall. The top slides down at %d feet per second. How fast is the base moving away from the wall? (enter a number)", t.l, t.x, rate)
		explanation = fmt.Sprintf("With \\(x^2 + y^2 = %d^2\\), differentiating gives \\(x dx/dt + y dy/dt = 0\\). Here \\(y = %d\\) and \\(dy/dt = -%d\\), so \\(dx/dt = \\) %d feet per second.", t.l, t.y, rate, answer)
	case 2:
		if rand.Intn(2) == 0 {
			answer = 1
			question = fmt.Sprintf("A %d-foot ladder leans with its base %d feet from the wall and top %d feet up. The top descends while the base stays on the floor. Enter 1 if the base moves away from the wall and -1 if it moves toward the wall.", t.l, t.x, t.y)
			explanation = "As the top descends, the base must move away from the wall to keep the ladder length constant, so the sign is positive."
		} else {
			answer = -1
			question = fmt.Sprintf("A %d-foot ladder leans with its base %d feet from the wall and top %d feet up. The top is pushed upward while the base stays on the floor. Enter 1 if the base moves away from the wall and -1 if it moves toward the wall.", t.l, t.x, t.y)
			explanation = "As the top rises, the base must move toward the wall to keep the ladder length constant, so the sign is negative."
		}
	default:
		rate := t.y * factor
		answer = t.x * factor
		question = fmt.Sprintf("A %d-foot ladder has its base %d feet from the wall and moves outward at %d feet per second. How fast is the top sliding down? (enter a number)", t.l, t.x, rate)
		explanation = fmt.Sprintf("\\(x^2 + y^2 = %d^2\\) gives \\(x dx/dt + y dy/dt = 0\\). With \\(y = %d\\) and \\(dx/dt = %d\\), the top moves at %d feet per second.", t.l, t.y, rate, answer)
	}

	facts := map[string]string{"x": strconv.Itoa(t.x), "y": strconv.Itoa(t.y), "ladder": strconv.Itoa(t.l), "factor": strconv.Itoa(factor)}
	if kp == 2 {
		facts["sign"] = strconv.Itoa(answer)
	}
	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Related rates: conical tank
// ---------------------------------------------------------------------------

type coneRateGen struct{}

func (g *coneRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	// The tank's surface radius is either half or one third of its depth. The
	// rate coefficient c = h^2 dh / ratio^2 is an integer because h is a
	// multiple of the radius ratio.
	ratio := 2
	h := 2 * randRange(ctx.Difficulty, 1, 6)
	if rand.Intn(2) == 0 {
		ratio = 3
		h = 3 * randRange(ctx.Difficulty, 1, 5)
	}
	dh := randRange(ctx.Difficulty, 1, 8)
	c := h * h * dh / (ratio * ratio)
	relation := "half"
	if ratio == 3 {
		relation = "one third"
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = dh
		question = fmt.Sprintf("Water flows into a cone at %dπ cubic units per minute. The water's surface radius is %s its depth. When the depth is %d, how fast is the depth rising? (enter a number)", c, relation, h)
		explanation = fmt.Sprintf("With \\(r = h/%d\\), \\(V = πh^3/%d\\) and \\(dV/dt = (πh^2/%d) dh/dt\\). Thus \\(%dπ = (π·%d/%d) dh/dt\\), so \\(dh/dt = %d\\).", ratio, 3*ratio*ratio, ratio*ratio, c, h*h, ratio*ratio, dh)
	case 2:
		answer = c
		question = fmt.Sprintf("A conical tank has surface radius %s its depth. When the depth is %d and rising at %d units per minute, at what rate is water entering, in units of π cubic units per minute? (enter the coefficient of π)", relation, h, dh)
		explanation = fmt.Sprintf("\\(dV/dt = (πh^2/%d) dh/dt = (π·%d/%d)(%d) = %dπ\\), so the coefficient is %d.", ratio*ratio, h*h, ratio*ratio, dh, c, c)
	default:
		answer = h
		question = fmt.Sprintf("A conical tank has surface radius %s its depth. Water enters at %dπ cubic units per minute, and at some depth the water level rises at %d units per minute. What is that depth? (enter a number)", relation, c, dh)
		explanation = fmt.Sprintf("\\(V = πh^3/%d\\) gives \\(dV/dt = (πh^2/%d) dh/dt\\). So \\(%dπ = (πh^2/%d)(%d)\\), which gives \\(h^2 = %d\\) and \\(h = %d\\).", 3*ratio*ratio, ratio*ratio, c, ratio*ratio, dh, h*h, h)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"h": strconv.Itoa(h), "c": strconv.Itoa(c), "dh": strconv.Itoa(dh), "ratio": strconv.Itoa(ratio)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Total distance from velocity
// ---------------------------------------------------------------------------

type totalDistanceGen struct{}

func (g *totalDistanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := 2 * randRange(ctx.Difficulty, 1, 12)
	b := a + 2*randRange(ctx.Difficulty, 1, 8)
	left := a * a / 2
	right := (b - a) * (b - a) / 2

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a
		question = fmt.Sprintf("A particle has velocity \\(v(t) = t - %d\\) on \\([0, %d]\\). At what time does it reverse direction? (enter a number)", a, b)
		explanation = fmt.Sprintf("\\(v(t) = 0\\) at \\(t = %d\\), and the velocity changes sign there.", a)
	case 2:
		answer = left
		question = fmt.Sprintf("A particle has velocity \\(v(t) = t - %d\\) on \\([0, %d]\\). What distance does it travel while \\(v ≤ 0\\)? (enter a number)", a, b)
		explanation = fmt.Sprintf("On \\([0, %d]\\) the speed is \\(%d - t\\), so the distance is \\(\\int_0^{%d}(%d - t)dt = %d\\).", a, a, a, a, left)
	default:
		answer = left + right
		question = fmt.Sprintf("A particle has velocity \\(v(t) = t - %d\\) on \\([0, %d]\\). What total distance does it travel? (enter a number)", a, b)
		explanation = fmt.Sprintf("The motion reverses at \\(t = %d\\). The two distances are %d and %d, so the total is %d.", a, left, right, left+right)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "distance": strconv.Itoa(left + right)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Integrals of piecewise functions
// ---------------------------------------------------------------------------

type piecewiseIntegralGen struct{}

func (g *piecewiseIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := rand.Intn(6) + 1
	b := rand.Intn(6) + 1
	c := rand.Intn(5) + 1
	d := c + rand.Intn(4) + 1
	left := a * c
	right := b * (d - c)

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = left
		question = fmt.Sprintf("For \\(f(x) = %d\\) when \\(x < %d\\) and \\(f(x) = %d\\) when \\(x ≥ %d\\), what is \\(\\int_0^{%d} f(x)dx\\)? (enter a number)", a, c, b, c, c)
		explanation = fmt.Sprintf("On \\([0, %d]\\) the function is constantly %d, so the integral is \\(%d × %d = %d\\).", c, a, a, c, left)
	case 2:
		answer = right
		question = fmt.Sprintf("For \\(f(x) = %d\\) when \\(x < %d\\) and \\(f(x) = %d\\) when \\(x ≥ %d\\), what is \\(\\int_%d^{%d} f(x)dx\\)? (enter a number)", a, c, b, c, c, d)
		explanation = fmt.Sprintf("On \\([%d, %d]\\) the function is constantly %d, so the integral is \\(%d × %d = %d\\).", c, d, b, b, d-c, right)
	default:
		answer = left + right
		question = fmt.Sprintf("For \\(f(x) = %d\\) when \\(x < %d\\) and \\(f(x) = %d\\) when \\(x ≥ %d\\), what is \\(\\int_0^{%d} f(x)dx\\)? (enter a number)", a, c, b, c, d)
		explanation = fmt.Sprintf("Split at %d: the left part is %d and the right part is %d, so the total is %d.", c, left, right, left+right)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Integrals of absolute value
// ---------------------------------------------------------------------------

type absoluteIntegralGen struct{}

func (g *absoluteIntegralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := 2 * randRange(ctx.Difficulty, 1, 12) // even, so half-integers stay whole
	shift := randRange(ctx.Difficulty, 0, 10)

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a * a
		if shift == 0 {
			question = fmt.Sprintf("What is \\(\\int_{-%d}^{%d} |x|dx\\)? (enter a number)", a, a)
			explanation = fmt.Sprintf("By symmetry the integral is twice the half integral \\(\\int_0^{%d}x dx = %d/2\\), so the total is %d.", a, a*a, a*a)
		} else {
			question = fmt.Sprintf("What is \\(\\int_{%d}^{%d} |x - %d|dx\\)? (enter a number)", shift-a, shift+a, shift)
			explanation = fmt.Sprintf("Shift the centre to the origin. The integral becomes \\(\\int_{-%d}^{%d}|u|du\\), which is twice the half integral \\(\\int_0^{%d}u du = %d/2\\), so the total is %d.", a, a, a, a*a, a*a)
		}
	case 2:
		if shift == 0 {
			answer = a * a / 2
			question = fmt.Sprintf("What is \\(\\int_0^{%d} |x|dx\\)? (enter a number)", a)
			explanation = fmt.Sprintf("On \\([0, %d]\\), \\(|x| = x\\), so the integral is \\(%d^2/2 = %d\\).", a, a, answer)
		} else {
			answer = a*a/2 + shift*a
			question = fmt.Sprintf("What is \\(\\int_0^{%d} |x + %d|dx\\)? (enter a number)", a, shift)
			explanation = fmt.Sprintf("On \\([0, %d]\\) the argument is non-negative, so \\(|x + %d| = x + %d\\). The integral is \\(%d^2/2 + %d·%d = %d\\).", a, shift, shift, a, shift, a, answer)
		}
	default:
		answer = a * a / 2
		if shift == 0 {
			question = fmt.Sprintf("What is \\(\\int_{-%d}^0 |x|dx\\)? (enter a number)", a)
			explanation = fmt.Sprintf("On \\([-%d, 0]\\), \\(|x| = -x\\), so the integral is \\(%d^2/2 = %d\\).", a, a, answer)
		} else {
			answer = a*a/2 + a*shift
			question = fmt.Sprintf("What is \\(\\int_{%d}^{%d} |x|dx\\)? (enter a number)", -a-shift, -shift)
			explanation = fmt.Sprintf("The interval lies left of the origin, where \\(|x| = -x\\). Then \\(\\int_{%d}^{%d}(-x)dx = \\frac{(%d)^2 - (%d)^2}{2} = %d\\).", -a-shift, -shift, a+shift, shift, answer)
		}
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "area": strconv.Itoa(a * a), "shift": strconv.Itoa(shift)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Accumulation functions
// ---------------------------------------------------------------------------

type accumulationGen struct{}

func (g *accumulationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := rand.Intn(4) + 1
	b := a + rand.Intn(6) + 1

	var question, explanation string
	var answer int
	solveX := 0
	target := 0
	switch kp {
	case 1:
		answer = b*b - a*a
		question = fmt.Sprintf("Let \\(F(x) = \\int_%d^x 2t\\,dt\\). What is \\(F(%d)\\)? (enter a number)", a, b)
		explanation = fmt.Sprintf("\\(F(x) = x^2 - %d^2\\), so \\(F(%d) = %d - %d = %d\\).", a, b, b*b, a*a, answer)
	case 2:
		answer = 2 * b
		question = fmt.Sprintf("Let \\(F(x) = \\int_%d^x 2t\\,dt\\). What is \\(F'(%d)\\)? (enter a number)", a, b)
		explanation = fmt.Sprintf("By the Fundamental Theorem of Calculus, \\(F'(x) = 2x\\), so \\(F'(%d) = %d\\).", b, 2*b)
	default:
		c := a + rand.Intn(4) + 1
		solveX = c
		answer = c
		v := c*c - a*a
		target = v
		question = fmt.Sprintf("Let \\(F(x) = \\int_%d^x 2t\\,dt\\). For what \\(x > %d\\) is \\(F(x) = %d\\)? (enter a number)", a, a, v)
		explanation = fmt.Sprintf("\\(F(x) = x^2 - %d\\), so \\(x^2 = %d + %d = %d\\) and \\(x = %d\\).", a*a, v, a*a, c*c, c)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "solve_x": strconv.Itoa(solveX), "target": strconv.Itoa(target)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Motion problems
// ---------------------------------------------------------------------------

type motionGen struct{}

func (g *motionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	b := rand.Intn(5)
	a := rand.Intn(4)
	c := a + rand.Intn(5) + 1
	velocity := 2*c + b
	displacement := (c*c - a*a) + b*(c-a)
	average := (c + a) + b

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = velocity
		question = fmt.Sprintf("A particle has position \\(s(t) = t^2 + %dt\\). What is its velocity at \\(t = %d\\)? (enter a number)", b, c)
		explanation = fmt.Sprintf("\\(v(t) = s'(t) = 2t + %d\\), so \\(v(%d) = %d\\).", b, c, velocity)
	case 2:
		answer = displacement
		question = fmt.Sprintf("A particle has position \\(s(t) = t^2 + %dt\\). What is its displacement from \\(t = %d\\) to \\(t = %d\\)? (enter a number)", b, a, c)
		explanation = fmt.Sprintf("\\(s(%d) - s(%d) = (%d + %d) - (%d + %d) = %d\\).", c, a, c*c, b*c, a*a, b*a, displacement)
	default:
		answer = average
		question = fmt.Sprintf("A particle has position \\(s(t) = t^2 + %dt\\). What is its average velocity on \\([%d, %d]\\)? (enter a number)", b, a, c)
		explanation = fmt.Sprintf("Average velocity is displacement over elapsed time: \\(%d/%d = %d\\).", displacement, c-a, average)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c)},
		KPOrdinal:   kp,
	}
}
