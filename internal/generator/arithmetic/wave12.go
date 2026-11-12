package arithmetic

// Wave 12 arithmetic concepts. Every answer is an integer and every operand
// needed to recompute it is published in Facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave12(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"arith.add.three_addends":       "add_three",
		"arith.add.four_digit":          "add_four",
		"arith.add.regroup_twice":       "add_regroup",
		"arith.sub.four_digit":          "sub_four",
		"arith.sub.regroup_twice":       "sub_regroup",
		"arith.sub.missing_subtrahend":  "sub_missing",
		"arith.mult.three_by_two":       "mult_three_two",
		"arith.mult.powers_ten":         "mult_powers",
		"arith.mult.three_by_three":     "mult_three_three",
		"arith.div.three_digit_divisor": "div_three",
		"arith.div.average":             "div_average",
		"arith.div.missing_dividend":    "div_missing",
		"arith.estimate.compatible":     "estimate_compatible",
		"arith.estimate.reasonableness": "estimate_reason",
		"arith.order_ops.exponents":     "order_exponents",
		"arith.neg.temperature":         "neg_temperature",
		"arith.neg.elevation":           "neg_elevation",
		"arith.order_ops.brackets":      "order_brackets",
		"arith.order_ops.negatives":     "order_negatives",
		"arith.word.unit_rate":          "word_unit_rate",
		"arith.word.comparison":         "word_comparison",
		"arith.dec.compare":             "dec_compare",
		"arith.dec.add":                 "dec_add",
		"arith.money.total":             "money_total",
		"arith.money.discount":          "money_discount",
		"arith.money.unit_price":        "money_unit_price",
		"arith.time.convert":            "time_convert",
		"arith.time.duration":           "time_duration",
		"arith.time.pace":               "time_pace",
		"arith.factor.divisor_count":    "divisor_count",
	} {
		reg.Register(id, &w12Gen{kind: kind})
	}
}

type w12Gen struct{ kind string }

func (g *w12Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "add_three":
		return w12AddThree(kp)
	case "add_four":
		return w12AddFour(kp)
	case "add_regroup":
		return w12AddRegroup(kp)
	case "sub_four":
		return w12SubFour(kp)
	case "sub_regroup":
		return w12SubRegroup(kp)
	case "sub_missing":
		return w12SubMissing(kp)
	case "mult_three_two":
		return w12MultThreeTwo(kp)
	case "mult_powers":
		return w12MultPowers(kp)
	case "mult_three_three":
		return w12MultThreeThree(kp)
	case "div_three":
		return w12DivThree(kp)
	case "div_average":
		return w12DivAverage(kp)
	case "div_missing":
		return w12DivMissing(kp)
	case "estimate_compatible":
		return w12EstimateCompatible(kp)
	case "estimate_reason":
		return w12EstimateReason(kp)
	case "order_exponents":
		return w12OrderExponents(kp)
	case "neg_temperature":
		return w12NegTemperature(kp)
	case "neg_elevation":
		return w12NegElevation(kp)
	case "order_brackets":
		return w12OrderBrackets(kp)
	case "order_negatives":
		return w12OrderNegatives(kp)
	case "word_unit_rate":
		return w12WordUnitRate(kp)
	case "word_comparison":
		return w12WordComparison(kp)
	case "dec_compare":
		return w12DecCompare(kp)
	case "dec_add":
		return w12DecAdd(kp)
	case "money_total":
		return w12MoneyTotal(kp)
	case "money_discount":
		return w12MoneyDiscount(kp)
	case "money_unit_price":
		return w12MoneyUnitPrice(kp)
	case "time_convert":
		return w12TimeConvert(kp)
	case "time_duration":
		return w12TimeDuration(kp)
	case "time_pace":
		return w12TimePace(kp)
	case "divisor_count":
		return w12DivisorCount(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4.", KPOrdinal: kp}
}

func w12R(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w12P(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w12I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w12AddThree(kp int) generator.Problem {
	var a, b, c int
	switch kp {
	case 1:
		a, b, c = w12R(10, 99), w12R(10, 99), w12R(10, 99)
	case 2:
		a, b, c = w12R(100, 999), w12R(100, 999), w12R(100, 999)
	default:
		a, b, c = w12R(1000, 9999), w12R(1000, 9999), w12R(1000, 9999)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	w12I(f, "c", c)
	sum := a + b + c
	return w12P(kp, fmt.Sprintf("What is %d + %d + %d? (enter a number)", a, b, c),
		strconv.Itoa(sum), fmt.Sprintf("%d + %d = %d, then %d + %d = %d.", a, b, a+b, a+b, c, sum), f)
}

func w12AddFour(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(1000, 9999), w12R(100, 999)
	case 2:
		a, b = w12R(1000, 9999), w12R(1000, 9999)
	default:
		a, b = w12R(10000, 99999), w12R(1000, 9999)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d + %d? (enter a number)", a, b),
		strconv.Itoa(a+b), fmt.Sprintf("Line up the places and add: %d + %d = %d.", a, b, a+b), f)
}

func w12AddRegroup(kp int) generator.Problem {
	maxHead := 4
	if kp >= 2 {
		maxHead = 8
	}
	if kp >= 3 {
		maxHead = 9
	}
	h1, h2 := w12R(1, maxHead), w12R(1, maxHead)
	t1, t2 := w12R(5, 9), w12R(5, 9)
	u1, u2 := w12R(5, 9), w12R(5, 9)
	a := h1*100 + t1*10 + u1
	b := h2*100 + t2*10 + u2
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d + %d? (enter a number)", a, b),
		strconv.Itoa(a+b), fmt.Sprintf("The units add to %d, so carry one ten. The tens add to %d plus the carry, so carry one hundred. The total is %d.", u1+u2, t1+t2, a+b), f)
}

func w12SubFour(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(1000, 9999), w12R(100, 999)
	case 2:
		a, b = w12R(2000, 9999), w12R(1000, 1999)
	default:
		a, b = w12R(20000, 99999), w12R(1000, 9999)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d - %d? (enter a number)", a, b),
		strconv.Itoa(a-b), fmt.Sprintf("Subtract place by place: %d - %d = %d.", a, b, a-b), f)
}

func w12SubRegroup(kp int) generator.Problem {
	head := 2
	if kp >= 2 {
		head = 5
	}
	if kp >= 3 {
		head = 9
	}
	h := w12R(head, 9)
	t1, t2 := w12R(1, 4), w12R(5, 9)
	u1, u2 := w12R(1, 4), w12R(5, 9)
	a := h*100 + t1*10 + u1
	b := t2*10 + u2
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d - %d? (enter a number)", a, b),
		strconv.Itoa(a-b), fmt.Sprintf("The units need a borrow and the tens need a borrow too. After regrouping, %d - %d = %d.", a, b, a-b), f)
}

func w12SubMissing(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(30, 99), w12R(10, 29)
	case 2:
		a, b = w12R(300, 999), w12R(100, 299)
	default:
		a, b = w12R(3000, 9999), w12R(1000, 2999)
	}
	c := a - b
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "c", c)
	return w12P(kp, fmt.Sprintf("%d - ? = %d. What is the missing number? (enter a number)", a, c),
		strconv.Itoa(b), fmt.Sprintf("Subtract the result from the starting number: %d - %d = %d.", a, c, b), f)
}

func w12MultThreeTwo(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(100, 999), w12R(10, 49)
	case 2:
		a, b = w12R(100, 999), w12R(50, 99)
	default:
		a, b = w12R(1000, 9999), w12R(10, 99)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d x %d? (enter a number)", a, b),
		strconv.Itoa(a*b), fmt.Sprintf("Split %d into tens and units, multiply each part by %d, then add: %d x %d = %d.", b, a, a, b, a*b), f)
}

func w12MultPowers(kp int) generator.Problem {
	var n, k int
	switch kp {
	case 1:
		n, k = w12R(2, 99), 1
	case 2:
		n, k = w12R(2, 99), 2
	default:
		n, k = w12R(2, 999), 3
	}
	p := 1
	for i := 0; i < k; i++ {
		p *= 10
	}
	f := map[string]string{}
	w12I(f, "n", n)
	w12I(f, "k", k)
	return w12P(kp, fmt.Sprintf("What is %d x %d? (enter a number)", n, p),
		strconv.Itoa(n*p), fmt.Sprintf("Multiplying by %d appends %d zero(s): %d x %d = %d.", p, k, n, p, n*p), f)
}

func w12MultThreeThree(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(100, 399), w12R(100, 399)
	case 2:
		a, b = w12R(400, 999), w12R(400, 999)
	default:
		a, b = w12R(1000, 9999), w12R(100, 999)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d x %d? (enter a number)", a, b),
		strconv.Itoa(a*b), fmt.Sprintf("Use partial products: %d x %d = %d.", a, b, a*b), f)
}

func w12DivThree(kp int) generator.Problem {
	var d, q int
	switch kp {
	case 1:
		d, q = w12R(100, 499), w12R(2, 9)
	case 2:
		d, q = w12R(100, 999), w12R(10, 99)
	default:
		d, q = w12R(100, 999), w12R(100, 999)
	}
	dividend := d * q
	f := map[string]string{}
	w12I(f, "dividend", dividend)
	w12I(f, "divisor", d)
	return w12P(kp, fmt.Sprintf("What is %d / %d? (enter a number)", dividend, d),
		strconv.Itoa(q), fmt.Sprintf("Because %d x %d = %d, the quotient is %d.", d, q, dividend, q), f)
}

func w12DivAverage(kp int) generator.Problem {
	var m, s int
	switch kp {
	case 1:
		m, s = w12R(10, 40), 2
	case 2:
		m, s = w12R(40, 120), 6
	default:
		m, s = w12R(120, 400), 15
	}
	n1, n2, n3, n4 := m-s, m-s/2, m+s/2, m+s
	if s < 4 {
		n1, n2, n3, n4 = m-2, m-1, m+1, m+2
	}
	f := map[string]string{}
	w12I(f, "n1", n1)
	w12I(f, "n2", n2)
	w12I(f, "n3", n3)
	w12I(f, "n4", n4)
	sum := n1 + n2 + n3 + n4
	return w12P(kp, fmt.Sprintf("What is the mean of %d, %d, %d and %d? (enter a number)", n1, n2, n3, n4),
		strconv.Itoa(m), fmt.Sprintf("The four numbers add to %d. Dividing by 4 gives %d.", sum, m), f)
}

func w12DivMissing(kp int) generator.Problem {
	var d, q int
	switch kp {
	case 1:
		d, q = w12R(2, 9), w12R(2, 9)
	case 2:
		d, q = w12R(10, 99), w12R(2, 9)
	default:
		d, q = w12R(100, 999), w12R(10, 99)
	}
	r := w12R(0, d-1)
	dividend := d*q + r
	f := map[string]string{}
	w12I(f, "divisor", d)
	w12I(f, "quotient", q)
	w12I(f, "remainder", r)
	return w12P(kp, fmt.Sprintf("A division has divisor %d, quotient %d and remainder %d. What is the dividend? (enter a number)", d, q, r),
		strconv.Itoa(dividend), fmt.Sprintf("Multiply the divisor by the quotient, then add the remainder: %d x %d + %d = %d.", d, q, r, dividend), f)
}

func w12EstimateCompatible(kp int) generator.Problem {
	var friendly, q int
	switch kp {
	case 1:
		friendly, q = w12R(2, 9)*10, w12R(10, 99)
	case 2:
		friendly, q = w12R(2, 9)*100, w12R(10, 99)
	default:
		friendly, q = w12R(2, 9)*100, w12R(100, 999)
	}
	delta := w12R(1, 4)
	if rand.Intn(2) == 0 {
		delta = -delta
	}
	divisor := friendly + delta
	if divisor <= 1 {
		divisor = friendly + 1
	}
	dividend := friendly * q
	f := map[string]string{}
	w12I(f, "dividend", dividend)
	w12I(f, "divisor", divisor)
	w12I(f, "friendly", friendly)
	return w12P(kp, fmt.Sprintf("Estimate %d / %d by using the compatible number %d. What is the estimated quotient? (enter a number)", dividend, divisor, friendly),
		strconv.Itoa(q), fmt.Sprintf("Treat the divisor as %d. Since %d / %d = %d, the estimate is %d.", friendly, dividend, friendly, q, q), f)
}

func w12EstimateReason(kp int) generator.Problem {
	var a, b, unit int
	switch kp {
	case 1:
		a, b, unit = w12R(100, 999), w12R(100, 999), 100
	case 2:
		a, b, unit = w12R(1000, 9999), w12R(1000, 9999), 1000
	default:
		a, b, unit = w12R(10000, 99999), w12R(10000, 99999), 10000
	}
	round := func(n int) int {
		return ((n + unit/2) / unit) * unit
	}
	est := round(a) + round(b)
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	w12I(f, "unit", unit)
	return w12P(kp, fmt.Sprintf("Estimate %d + %d by rounding each number to the nearest %d. What is the estimated sum? (enter a number)", a, b, unit),
		strconv.Itoa(est), fmt.Sprintf("%d rounds to %d and %d rounds to %d, so the estimate is %d.", a, round(a), b, round(b), est), f)
}

func w12OrderExponents(kp int) generator.Problem {
	var a, b, c, d int
	var answer int
	var q, ex string
	switch kp {
	case 1:
		a, b, c, d = w12R(1, 20), w12R(2, 5), w12R(2, 3), 0
		p := 1
		for i := 0; i < c; i++ {
			p *= b
		}
		answer = a + p
		q = fmt.Sprintf("What is %d + %d^%d? (enter a number)", a, b, c)
		ex = fmt.Sprintf("The power first: %d^%d = %d. Then %d + %d = %d.", b, c, p, a, p, answer)
	case 2:
		a, b, c, d = w12R(10, 40), w12R(2, 5), w12R(2, 3), w12R(2, 6)
		p := 1
		for i := 0; i < c; i++ {
			p *= b
		}
		answer = a + p*d
		q = fmt.Sprintf("What is %d + %d^%d x %d? (enter a number)", a, b, c, d)
		ex = fmt.Sprintf("The power first: %d^%d = %d. Then %d x %d = %d, and %d + %d = %d.", b, c, p, p, d, p*d, a, p*d, answer)
	default:
		a, b, c, d = w12R(2, 9), w12R(2, 9), w12R(2, 3), w12R(1, 50)
		base := a + b
		p := 1
		for i := 0; i < c; i++ {
			p *= base
		}
		answer = p - d
		q = fmt.Sprintf("What is (%d + %d)^%d - %d? (enter a number)", a, b, c, d)
		ex = fmt.Sprintf("Brackets first: %d + %d = %d. Then %d^%d = %d, and %d - %d = %d.", a, b, base, base, c, p, p, d, answer)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	w12I(f, "c", c)
	w12I(f, "d", d)
	return w12P(kp, q, strconv.Itoa(answer), ex, f)
}

func w12NegTemperature(kp int) generator.Problem {
	var start, change int
	switch kp {
	case 1:
		start, change = w12R(-10, 10), w12R(-10, 10)
	case 2:
		start, change = w12R(-30, 30), w12R(-20, 20)
	default:
		start, change = w12R(-40, 5), w12R(5, 30)
	}
	f := map[string]string{}
	w12I(f, "start", start)
	w12I(f, "change", change)
	return w12P(kp, fmt.Sprintf("The temperature starts at %d degrees and changes by %d degrees. What is the final temperature? (enter a number)", start, change),
		strconv.Itoa(start+change), fmt.Sprintf("%d + (%d) = %d.", start, change, start+change), f)
}

func w12NegElevation(kp int) generator.Problem {
	var high, low int
	switch kp {
	case 1:
		high, low = w12R(0, 100), w12R(-100, -1)
	case 2:
		high, low = w12R(100, 500), w12R(-400, -100)
	default:
		high, low = w12R(1000, 8000), w12R(-500, -50)
	}
	f := map[string]string{}
	w12I(f, "high", high)
	w12I(f, "low", low)
	return w12P(kp, fmt.Sprintf("A peak is %d metres above sea level and a valley is %d metres below sea level. How many metres apart are they? (enter a number)", high, low),
		strconv.Itoa(high-low), fmt.Sprintf("Subtract the lower elevation from the higher one: %d - (%d) = %d.", high, low, high-low), f)
}

func w12OrderBrackets(kp int) generator.Problem {
	var a, b, c, d int
	switch kp {
	case 1:
		a, b, c, d = w12R(2, 9), w12R(2, 9), w12R(2, 5), w12R(1, 9)
	case 2:
		a, b, c, d = w12R(10, 40), w12R(10, 40), w12R(2, 9), w12R(5, 40)
	default:
		a, b, c, d = w12R(20, 90), w12R(20, 90), w12R(3, 12), w12R(10, 90)
	}
	answer := (a+b)*c - d
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	w12I(f, "c", c)
	w12I(f, "d", d)
	return w12P(kp, fmt.Sprintf("What is (%d + %d) x %d - %d? (enter a number)", a, b, c, d),
		strconv.Itoa(answer), fmt.Sprintf("Brackets first: %d + %d = %d. Then %d x %d = %d, and %d - %d = %d.", a, b, a+b, a+b, c, (a+b)*c, (a+b)*c, d, answer), f)
}

func w12OrderNegatives(kp int) generator.Problem {
	a, b, c, d := w12R(1, 20), w12R(1, 20), w12R(2, 9), w12R(1, 20)
	var answer int
	var q, ex string
	switch kp {
	case 1:
		answer = -a + b*c - d
		q = fmt.Sprintf("What is -%d + %d x %d - %d? (enter a number)", a, b, c, d)
		ex = fmt.Sprintf("Multiply first: %d x %d = %d. Then -%d + %d - %d = %d.", b, c, b*c, a, b*c, d, answer)
	case 2:
		answer = (a-b)*c + d
		q = fmt.Sprintf("What is (%d - %d) x %d + %d? (enter a number)", a, b, c, d)
		ex = fmt.Sprintf("Brackets first: %d - %d = %d. Then %d x %d + %d = %d.", a, b, a-b, a-b, c, d, answer)
	default:
		answer = a - (b+c)*d
		q = fmt.Sprintf("What is %d - (%d + %d) x %d? (enter a number)", a, b, c, d)
		ex = fmt.Sprintf("Brackets first: %d + %d = %d. Then %d x %d = %d, and %d - %d = %d.", b, c, b+c, b+c, d, (b+c)*d, a, (b+c)*d, answer)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	w12I(f, "c", c)
	w12I(f, "d", d)
	return w12P(kp, q, strconv.Itoa(answer), ex, f)
}

func w12WordUnitRate(kp int) generator.Problem {
	var n, unit int
	switch kp {
	case 1:
		n, unit = w12R(2, 6), w12R(5, 40)
	case 2:
		n, unit = w12R(3, 12), w12R(10, 90)
	default:
		n, unit = w12R(5, 20), w12R(50, 400)
	}
	cost := n * unit
	f := map[string]string{}
	w12I(f, "items", n)
	w12I(f, "cost", cost)
	return w12P(kp, fmt.Sprintf("%d identical items cost %d cents. What is the cost of one item in cents? (enter a number)", n, cost),
		strconv.Itoa(unit), fmt.Sprintf("Divide the total cost by the number of items: %d / %d = %d cents.", cost, n, unit), f)
}

func w12WordComparison(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(20, 99), w12R(1, 19)
	case 2:
		a, b = w12R(100, 999), w12R(10, 99)
	default:
		a, b = w12R(1000, 9999), w12R(100, 999)
	}
	f := map[string]string{}
	w12I(f, "first", a)
	w12I(f, "second", b)
	return w12P(kp, fmt.Sprintf("One shelf holds %d books and another holds %d books. How many more books does the first shelf hold? (enter a number)", a, b),
		strconv.Itoa(a-b), fmt.Sprintf("Subtract the smaller count from the larger: %d - %d = %d.", a, b, a-b), f)
}

func w12DecCompare(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(2, 9)*10, w12R(1, 8)*10
	case 2:
		a, b = w12R(11, 99), w12R(1, 10)
	default:
		a, b = w12R(50, 99), w12R(5, 49)
	}
	if a < b {
		a, b = b, a
	}
	if a == b {
		a++
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("How much larger is %d hundredths than %d hundredths? Enter the difference in hundredths. (enter a number)", a, b),
		strconv.Itoa(a-b), fmt.Sprintf("Subtract: %d - %d = %d hundredths.", a, b, a-b), f)
}

func w12DecAdd(kp int) generator.Problem {
	var a, b int
	switch kp {
	case 1:
		a, b = w12R(1, 9)*10, w12R(1, 9)*10
	case 2:
		a, b = w12R(11, 99), w12R(11, 99)
	default:
		a, b = w12R(100, 999), w12R(100, 999)
	}
	f := map[string]string{}
	w12I(f, "a", a)
	w12I(f, "b", b)
	return w12P(kp, fmt.Sprintf("What is %d hundredths plus %d hundredths? Enter the sum in hundredths. (enter a number)", a, b),
		strconv.Itoa(a+b), fmt.Sprintf("%d + %d = %d hundredths.", a, b, a+b), f)
}

func w12MoneyTotal(kp int) generator.Problem {
	var p1, p2, p3 int
	switch kp {
	case 1:
		p1, p2, p3 = w12R(5, 99), w12R(5, 99), w12R(5, 99)
	case 2:
		p1, p2, p3 = w12R(100, 999), w12R(100, 999), w12R(100, 999)
	default:
		p1, p2, p3 = w12R(1000, 9999), w12R(1000, 9999), w12R(1000, 9999)
	}
	f := map[string]string{}
	w12I(f, "p1", p1)
	w12I(f, "p2", p2)
	w12I(f, "p3", p3)
	total := p1 + p2 + p3
	return w12P(kp, fmt.Sprintf("Three items cost %d, %d and %d cents. What is the total in cents? (enter a number)", p1, p2, p3),
		strconv.Itoa(total), fmt.Sprintf("%d + %d + %d = %d cents.", p1, p2, p3, total), f)
}

func w12MoneyDiscount(kp int) generator.Problem {
	var price int
	switch kp {
	case 1:
		price = w12R(1, 20) * 100
	case 2:
		price = w12R(10, 100) * 100
	default:
		price = w12R(50, 500) * 100
	}
	percents := []int{10, 20, 25, 50}
	percent := percents[rand.Intn(len(percents))]
	discount := price * percent / 100
	f := map[string]string{}
	w12I(f, "price", price)
	w12I(f, "percent", percent)
	return w12P(kp, fmt.Sprintf("An item costs %d cents and is discounted by %d percent. What is the discount in cents? (enter a number)", price, percent),
		strconv.Itoa(discount), fmt.Sprintf("%d percent of %d cents is %d x %d / 100 = %d cents.", percent, price, price, percent, discount), f)
}

func w12MoneyUnitPrice(kp int) generator.Problem {
	var n, unit int
	switch kp {
	case 1:
		n, unit = w12R(2, 6), w12R(50, 300)
	case 2:
		n, unit = w12R(3, 12), w12R(100, 800)
	default:
		n, unit = w12R(5, 20), w12R(200, 2500)
	}
	total := n * unit
	f := map[string]string{}
	w12I(f, "total", total)
	w12I(f, "items", n)
	return w12P(kp, fmt.Sprintf("%d items cost %d cents in total. What is the unit price in cents? (enter a number)", n, total),
		strconv.Itoa(unit), fmt.Sprintf("Divide the total by the count: %d / %d = %d cents.", total, n, unit), f)
}

func w12TimeConvert(kp int) generator.Problem {
	var h, m int
	switch kp {
	case 1:
		h, m = w12R(1, 5), w12R(0, 59)
	case 2:
		h, m = w12R(6, 12), w12R(0, 59)
	default:
		h, m = w12R(13, 48), w12R(0, 59)
	}
	f := map[string]string{}
	w12I(f, "hours", h)
	w12I(f, "minutes", m)
	total := h*60 + m
	return w12P(kp, fmt.Sprintf("How many minutes are in %d hours and %d minutes? (enter a number)", h, m),
		strconv.Itoa(total), fmt.Sprintf("%d hours is %d minutes, and adding %d gives %d minutes.", h, h*60, m, total), f)
}

func w12TimeDuration(kp int) generator.Problem {
	var start, end int
	switch kp {
	case 1:
		start = w12R(0, 300)
		end = start + w12R(20, 180)
	case 2:
		start = w12R(300, 900)
		end = start + w12R(30, 240)
	default:
		start = w12R(600, 1200)
		end = start + w12R(45, 300)
	}
	f := map[string]string{}
	w12I(f, "start", start)
	w12I(f, "end", end)
	return w12P(kp, fmt.Sprintf("An activity starts at minute %d of the day and ends at minute %d. How many minutes does it last? (enter a number)", start, end),
		strconv.Itoa(end-start), fmt.Sprintf("Subtract the start from the end: %d - %d = %d minutes.", end, start, end-start), f)
}

func w12TimePace(kp int) generator.Problem {
	var distance, pace int
	switch kp {
	case 1:
		distance, pace = w12R(2, 6), w12R(5, 15)
	case 2:
		distance, pace = w12R(5, 15), w12R(4, 12)
	default:
		distance, pace = w12R(10, 42), w12R(3, 10)
	}
	minutes := distance * pace
	f := map[string]string{}
	w12I(f, "minutes", minutes)
	w12I(f, "distance", distance)
	return w12P(kp, fmt.Sprintf("A runner covers %d kilometres in %d minutes. What is the pace in minutes per kilometre? (enter a number)", distance, minutes),
		strconv.Itoa(pace), fmt.Sprintf("Divide the time by the distance: %d / %d = %d minutes per kilometre.", minutes, distance, pace), f)
}

func w12DivisorCount(kp int) generator.Problem {
	var n int
	switch kp {
	case 1:
		n = w12R(12, 60)
	case 2:
		n = w12R(61, 200)
	default:
		n = w12R(201, 1000)
	}
	count := 0
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			count++
			if i*i != n {
				count++
			}
		}
	}
	f := map[string]string{}
	w12I(f, "n", n)
	return w12P(kp, fmt.Sprintf("How many positive divisors does %d have? (enter a number)", n),
		strconv.Itoa(count), fmt.Sprintf("Pair each divisor below the square root with one above it. Counting the pairs and the square root if it is exact gives %d divisors.", count), f)
}
