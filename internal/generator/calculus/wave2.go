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
	a := rand.Intn(5) + 2 // 2..6

	var question, explanation string
	answer := a * a
	switch kp {
	case 1:
		question = fmt.Sprintf("Evaluate \\(\\lim_{x→%d}(1/x - 1/%d)/(x - %d)\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a, a)
		explanation = fmt.Sprintf("Combine the fractions: \\((1/x - 1/%d)/(x - %d) = -1/(%d x)\\), whose limit is \\(-1/%d\\).", a, a, a, a*a)
	case 2:
		question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(x + %d) - 1/%d)/x\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a)
		explanation = fmt.Sprintf("Combine the numerator over \\(%d(x + %d)\\), then divide by x. The limit is \\(-1/%d\\).", a, a, a*a)
	default:
		question = fmt.Sprintf("Evaluate \\(\\lim_{x→0}(1/(x - %d) + 1/%d)/x\\). The limit is \\(-1/k\\). What is k? (enter a number)", a, a)
		explanation = fmt.Sprintf("The numerator is \\(x/(%d(x - %d))\\). Dividing by x leaves \\(1/(%d(x - %d))\\), whose limit is \\(-1/%d\\).", a, a, a, a, a*a)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Limits with absolute values
// ---------------------------------------------------------------------------

type absoluteValueLimitGen struct{}

func (g *absoluteValueLimitGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := rand.Intn(5) + 1

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = -1
		question = fmt.Sprintf("For \\(f(x) = |x - %d|/(x - %d)\\), what is the left-hand limit as \\(x → %d\\)? (enter a number)", a, a, a)
		explanation = fmt.Sprintf("For \\(x < %d\\), \\(|x - %d| = -(x - %d)\\), so the quotient is -1 throughout the left side.", a, a, a)
	case 2:
		answer = 1
		question = fmt.Sprintf("For \\(f(x) = |x - %d|/(x - %d)\\), what is the right-hand limit as \\(x → %d\\)? (enter a number)", a, a, a)
		explanation = fmt.Sprintf("For \\(x > %d\\), \\(|x - %d| = x - %d\\), so the quotient is 1 throughout the right side.", a, a, a)
	default:
		answer = 0
		question = fmt.Sprintf("For \\(f(x) = |x - %d|/(x - %d)\\), does the two-sided limit as \\(x → %d\\) exist? Enter 1 for yes and 0 for no.", a, a, a)
		explanation = "The left limit is -1 and the right limit is 1. They disagree, so the two-sided limit does not exist."
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Slopes of normal lines
// ---------------------------------------------------------------------------

type normalLineGen struct{}

func (g *normalLineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := rand.Intn(5) + 1

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = 2 * a
		question = fmt.Sprintf("For \\(f(x) = x^2\\), what is the slope of the tangent line at \\(x = %d\\)? (enter a number)", a)
		explanation = fmt.Sprintf("\\(f'(x) = 2x\\), so at \\(x = %d\\) the tangent slope is %d.", a, 2*a)
	case 2:
		answer = 2 * a
		question = fmt.Sprintf("For \\(f(x) = x^2\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
		explanation = fmt.Sprintf("The tangent slope is \\(2·%d = %d\\). The normal slope is the negative reciprocal, \\(-1/%d\\), so \\(k = %d\\).", a, 2*a, 2*a, 2*a)
	default:
		answer = 3 * a * a
		question = fmt.Sprintf("For \\(f(x) = x^3\\) at \\(x = %d\\), the normal line has slope \\(-1/k\\). What is k? (enter a number)", a)
		explanation = fmt.Sprintf("\\(f'(x) = 3x^2\\), so the tangent slope at \\(x = %d\\) is \\(3·%d^2 = %d\\). The normal slope is \\(-1/%d\\).", a, a, 3*a*a, 3*a*a)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a)},
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
	k := rand.Intn(9) + 2 // 2..10
	perimeter := 4 * k
	width := k
	length := 2 * k
	area := 2 * k * k

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = width
		question = fmt.Sprintf("A rectangular pen is built against a river, so only three sides need fence. The fence is %d meters long. What width parallel to the river maximizes the area? (enter a number)", perimeter)
		explanation = fmt.Sprintf("If the two sides perpendicular to the river are x, the side parallel is \\(%d - 2x\\). Area \\(A = x(%d - 2x)\\) is maximized at \\(x = %d/4 = %d\\).", perimeter, perimeter, perimeter, width)
	case 2:
		answer = length
		question = fmt.Sprintf("A three-sided rectangular pen uses %d meters of fence against a river. At the maximum area, what is the length of the side parallel to the river? (enter a number)", perimeter)
		explanation = fmt.Sprintf("The optimal perpendicular width is \\(%d/4 = %d\\), so the parallel side is \\(%d - 2(%d) = %d\\).", perimeter, width, perimeter, width, length)
	default:
		answer = area
		question = fmt.Sprintf("A three-sided rectangular pen with %d meters of fence is built against a river. What is its maximum area? (enter a number)", perimeter)
		explanation = fmt.Sprintf("The optimum has width %d and length %d, so the maximum area is \\(%d × %d = %d\\).", width, length, width, length, area)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"perimeter": strconv.Itoa(perimeter), "width": strconv.Itoa(width), "length": strconv.Itoa(length), "area": strconv.Itoa(area)},
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

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = t.y
		question = fmt.Sprintf("A %d-foot ladder leans against a wall with its base %d feet from the wall. The top slides down at %d feet per second. How fast is the base moving away from the wall? (enter a number)", t.l, t.x, t.x)
		explanation = fmt.Sprintf("With \\(x^2 + y^2 = %d^2\\), differentiating gives \\(x dx/dt + y dy/dt = 0\\). Here \\(y = %d\\) and \\(dy/dt = -%d\\), so \\(dx/dt = \\) %d feet per second.", t.l, t.y, t.x, t.y)
	case 2:
		answer = 1
		question = fmt.Sprintf("A %d-foot ladder leans with its base %d feet from the wall and top %d feet up. The top descends while the base stays on the floor. Enter 1 if the base moves away from the wall and -1 if it moves toward the wall.", t.l, t.x, t.y)
		explanation = "As the top descends, the base must move away from the wall to keep the ladder length constant, so the sign is positive."
	default:
		answer = t.x
		question = fmt.Sprintf("A %d-foot ladder has its base %d feet from the wall and moves outward at %d feet per second. How fast is the top sliding down? (enter a number)", t.l, t.x, t.y)
		explanation = fmt.Sprintf("\\(x^2 + y^2 = %d^2\\) gives \\(x dx/dt + y dy/dt = 0\\). With \\(y = %d\\) and \\(dx/dt = %d\\), the top moves at %d feet per second.", t.l, t.y, t.y, t.x)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"x": strconv.Itoa(t.x), "y": strconv.Itoa(t.y), "ladder": strconv.Itoa(t.l)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Related rates: conical tank
// ---------------------------------------------------------------------------

type coneRateGen struct{}

func (g *coneRateGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	h := 2
	c := rand.Intn(6) + 1
	if rand.Intn(2) == 0 {
		h = 4
		c = 4 * (rand.Intn(3) + 1)
	}
	dh := 4 * c / (h * h)

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = dh
		question = fmt.Sprintf("Water flows into a cone at %dπ cubic units per minute. The water's surface radius is half its depth. When the depth is %d, how fast is the depth rising? (enter a number)", c, h)
		explanation = fmt.Sprintf("With \\(r = h/2\\), \\(V = πh^3/12\\) and \\(dV/dt = (πh^2/4) dh/dt\\). Thus \\(%dπ = (π·%d/4) dh/dt\\), so \\(dh/dt = %d\\).", c, h*h, dh)
	case 2:
		answer = c
		question = fmt.Sprintf("A conical tank has surface radius half its depth. When the depth is %d and rising at %d units per minute, at what rate is water entering, in units of π cubic units per minute? (enter the coefficient of π)", h, dh)
		explanation = fmt.Sprintf("\\(dV/dt = (πh^2/4) dh/dt = (π·%d/4)(%d) = %dπ\\), so the coefficient is %d.", h*h, dh, c, c)
	default:
		answer = h
		question = fmt.Sprintf("A conical tank has surface radius half its depth. Water enters at %dπ cubic units per minute, and at some depth the water level rises at %d units per minute. What is that depth? (enter a number)", c, dh)
		explanation = fmt.Sprintf("\\(V = πh^3/12\\) gives \\(dV/dt = (πh^2/4) dh/dt\\). So \\(%dπ = (πh^2/4)(%d)\\), which gives \\(h^2 = %d\\) and \\(h = %d\\).", c, dh, h*h, h)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"h": strconv.Itoa(h), "c": strconv.Itoa(c), "dh": strconv.Itoa(dh)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Total distance from velocity
// ---------------------------------------------------------------------------

type totalDistanceGen struct{}

func (g *totalDistanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	a := 2 * (rand.Intn(3) + 1)
	b := a + 2*(rand.Intn(3)+1)
	left := a * a / 2
	right := (b - a) * (b - a) / 2

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a
		question = fmt.Sprintf("A particle has velocity \\(v(t) = t - a\\) on \\([0, %d]\\). At what time does it reverse direction? (enter a number)", b)
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
	a := 2 * (rand.Intn(9) + 1) // 2..18

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a * a
		question = fmt.Sprintf("What is \\(\\int_{-%d}^{%d} |x|dx\\)? (enter a number)", a, a)
		explanation = fmt.Sprintf("By symmetry the integral is twice \\(\\int_0^{%d}x dx = %d/2\\), so the total is %d.", a, a*a, a*a)
	case 2:
		answer = a * a / 2
		question = fmt.Sprintf("What is \\(\\int_0^{%d} |x|dx\\)? (enter a number)", a)
		explanation = fmt.Sprintf("On \\([0, %d]\\), \\(|x| = x\\), so the integral is \\(%d^2/2 = %d\\).", a, a, answer)
	default:
		answer = a * a / 2
		question = fmt.Sprintf("What is \\(\\int_{-%d}^0 |x|dx\\)? (enter a number)", a)
		explanation = fmt.Sprintf("On \\([-%d, 0]\\), \\(|x| = -x\\), so the integral is \\(%d^2/2 = %d\\).", a, a, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(a), "area": strconv.Itoa(a * a)},
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
