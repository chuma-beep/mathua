package arithmetic

// Wave 8 arithmetic remaining concepts. Every answer is an integer and every
// operand needed to recompute it is published in Facts.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave8Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"arith.place.base_five":           "base_five",
		"arith.place.roman":               "roman",
		"arith.place.tally":               "tally",
		"arith.ops.doubling_halving":      "doubling_halving",
		"arith.ops.compensation":          "compensation",
		"arith.ops.near_doubles":          "near_doubles",
		"arith.div.partial_quotients":     "partial_quotients",
		"arith.mult.area_model":           "area_model",
		"arith.mult.lattice":              "lattice",
		"arith.factor.square_numbers":     "square_numbers",
		"arith.factor.triangular_numbers": "triangular_numbers",
		"arith.ops.factorial_notation":    "factorial",
		"arith.factor.gcf_lcm_word":       "gcf_lcm",
		"arith.time.elapsed":              "elapsed",
		"arith.money.change":              "change",
	} {
		reg.Register(id, &w8ArithRemainingGen{kind: kind})
	}
}

type w8ArithRemainingGen struct{ kind string }

func (g *w8ArithRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "base_five":
		return w8BaseFive(kp)
	case "roman":
		return w8Roman(kp)
	case "tally":
		return w8Tally(kp)
	case "doubling_halving":
		return w8DoublingHalving(kp)
	case "compensation":
		return w8Compensation(kp)
	case "near_doubles":
		return w8NearDoubles(kp)
	case "partial_quotients":
		return w8PartialQuotients(kp)
	case "area_model":
		return w8AreaModel(kp)
	case "lattice":
		return w8Lattice(kp)
	case "square_numbers":
		return w8SquareNumbers(kp)
	case "triangular_numbers":
		return w8TriangularNumbers(kp)
	case "factorial":
		return w8Factorial(kp)
	case "gcf_lcm":
		return w8GCFLCM(kp)
	case "elapsed":
		return w8Elapsed(kp)
	case "change":
		return w8Change(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w8RGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w8R(o int, q, a, ex string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w8RI(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w8RR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w8BaseString(n, base int) string {
	if n == 0 {
		return "0"
	}
	digits := "0123456789"
	var b strings.Builder
	for n > 0 {
		b.WriteByte(digits[n%base])
		n /= base
	}
	s := b.String()
	r := []byte(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func w8BaseFive(kp int) generator.Problem {
	n := w8RR(5, 124)
	numeral := w8BaseString(n, 5)
	f := map[string]string{"numeral": numeral}
	w8RI(f, "n", n)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the base-five numeral \\(%s_5\\) in decimal? (enter a number)", numeral), strconv.Itoa(n),
			fmt.Sprintf("The place values are powers of five. The digits of %s evaluate to %d.", numeral, n), f)
	case 2:
		p := 1
		for p*5 <= n {
			p *= 5
		}
		w8RI(f, "power", p)
		return w8R(kp, fmt.Sprintf("What is the largest power of five that is at most %d? (enter a number)", n), strconv.Itoa(p),
			fmt.Sprintf("Multiply by five until the next power would exceed %d. The largest power is %d.", n, p), f)
	default:
		sum := 0
		for _, r := range numeral {
			sum += int(r - '0')
		}
		w8RI(f, "digit_sum", sum)
		return w8R(kp, fmt.Sprintf("What is the sum of the digits of the base-five numeral \\(%s_5\\)? (enter a number)", numeral), strconv.Itoa(sum),
			fmt.Sprintf("The digits of %s add to %d.", numeral, sum), f)
	}
}

func w8RomanString(n int) string {
	vals := []int{100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	var b strings.Builder
	for i, v := range vals {
		for n >= v {
			b.WriteString(syms[i])
			n -= v
		}
	}
	return b.String()
}

func w8Roman(kp int) generator.Problem {
	n := w8RR(4, 399)
	roman := w8RomanString(n)
	f := map[string]string{"roman": roman}
	w8RI(f, "n", n)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the value of the Roman numeral %s? (enter a number)", roman), strconv.Itoa(n),
			fmt.Sprintf("Read the numeral from left to right with subtractive pairs. %s is %d.", roman, n), f)
	case 2:
		count := strings.Count(roman, "X")
		w8RI(f, "x_count", count)
		return w8R(kp, fmt.Sprintf("How many X symbols are in the Roman numeral %s? (enter a number)", roman), strconv.Itoa(count),
			fmt.Sprintf("Count the letter X in %s: there %s %d.", roman, map[bool]string{true: "is", false: "are"}[count == 1], count), f)
	default:
		code := 0
		for _, pair := range []string{"IV", "IX", "XL", "XC", "CD", "CM"} {
			if strings.Contains(roman, pair) {
				code = 1
				break
			}
		}
		w8RI(f, "subtractive", code)
		return w8R(kp, fmt.Sprintf("Does the Roman numeral %s use a subtractive pair? Enter 1 for yes and 0 for no.", roman), strconv.Itoa(code),
			fmt.Sprintf("A subtractive pair is one of IV, IX, XL, XC, CD or CM. The answer for %s is %d.", roman, code), f)
	}
}

func w8Tally(kp int) generator.Problem {
	n := w8RR(6, 99)
	groups := n / 5
	rem := n % 5
	f := map[string]string{}
	w8RI(f, "n", n)
	w8RI(f, "groups", groups)
	w8RI(f, "remainder", rem)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("A tally has %d complete groups of five and %d extra marks. How many marks is that? (enter a number)", groups, rem), strconv.Itoa(n),
			fmt.Sprintf("\\(%d\\times5+%d=%d\\).", groups, rem, n), f)
	case 2:
		return w8R(kp, fmt.Sprintf("How many complete groups of five are in a tally of %d marks? (enter a number)", n), strconv.Itoa(groups),
			fmt.Sprintf("\\(%d/5=%d\\) complete groups, with %d %s left over.", n, groups, rem, map[bool]string{true: "mark", false: "marks"}[rem == 1]), f)
	default:
		return w8R(kp, fmt.Sprintf("After grouping %d tally marks into fives, how many marks are left over? (enter a number)", n), strconv.Itoa(rem),
			fmt.Sprintf("\\(%d=%d\\times5+%d\\), so the remainder is %d.", n, groups, rem, rem), f)
	}
}

func w8DoublingHalving(kp int) generator.Problem {
	a := w8RR(2, 30)
	b := w8RR(2, 12)
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the double of %d? (enter a number)", a), strconv.Itoa(2*a),
			fmt.Sprintf("\\(2\\times%d=%d\\).", a, 2*a), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is half of %d? (enter a number)", 2*a), strconv.Itoa(a),
			fmt.Sprintf("Half of \\(2\\times%d\\) is %d.", a, a), f)
	default:
		return w8R(kp, fmt.Sprintf("Use doubling and halving to compute \\(%d\\times%d\\). What is the product? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("Double one factor while halving the other until the halving is easy. The product is %d.", a*b), f)
	}
}

func w8Compensation(kp int) generator.Problem {
	a := w8RR(20, 90)
	delta := w8RR(7, 10)
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "delta", delta)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("To add %d to %d by compensation, first add 10. How much must be subtracted? (enter a number)", delta, a), strconv.Itoa(10-delta),
			fmt.Sprintf("Adding 10 overshoots by \\(10-%d=%d\\).", delta, 10-delta), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is \\(%d+%d\\)? (enter a number)", a, delta), strconv.Itoa(a+delta),
			fmt.Sprintf("Add ten to get %d, then subtract %d to get %d.", a+10, 10-delta, a+delta), f)
	default:
		return w8R(kp, fmt.Sprintf("Adding 10 to %d gives what? (enter a number)", a), strconv.Itoa(a+10),
			fmt.Sprintf("\\(%d+10=%d\\).", a, a+10), f)
	}
}

func w8NearDoubles(kp int) generator.Problem {
	a := w8RR(2, 20)
	f := map[string]string{}
	w8RI(f, "a", a)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the double of %d? (enter a number)", a), strconv.Itoa(2*a),
			fmt.Sprintf("\\(2\\times%d=%d\\).", a, 2*a), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is \\(%d+%d\\)? (enter a number)", a, a+1), strconv.Itoa(2*a+1),
			fmt.Sprintf("This is a near double: \\(%d+%d=2\\times%d+1=%d\\).", a, a+1, a, 2*a+1), f)
	default:
		return w8R(kp, fmt.Sprintf("How much more than the double of %d is \\(%d+%d\\)? (enter a number)", a, a, a+1), "1",
			fmt.Sprintf("\\(%d+%d\\) is one more than \\(2\\times%d\\).", a, a+1, a), f)
	}
}

func w8PartialQuotients(kp int) generator.Problem {
	divisor := w8RR(3, 9)
	q := w8RR(3, 25)
	rem := w8RR(0, divisor-1)
	dividend := divisor*q + rem
	f := map[string]string{}
	w8RI(f, "dividend", dividend)
	w8RI(f, "divisor", divisor)
	w8RI(f, "quotient", q)
	w8RI(f, "remainder", rem)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("In the partial-quotients division of %d by %d, what is the final quotient? (enter a number)", dividend, divisor), strconv.Itoa(q),
			fmt.Sprintf("\\(%d=%d\\times%d+%d\\), so the quotient is %d.", dividend, divisor, q, rem, q), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is the remainder when %d is divided by %d? (enter a number)", dividend, divisor), strconv.Itoa(rem),
			fmt.Sprintf("\\(%d=%d\\times%d+%d\\), so the remainder is %d.", dividend, divisor, q, rem, rem), f)
	default:
		return w8R(kp, fmt.Sprintf("What is \\(%d\\times%d\\) in the partial-quotients work for %d divided by %d? (enter a number)", divisor, q, dividend, divisor), strconv.Itoa(divisor*q),
			fmt.Sprintf("\\(%d\\times%d=%d\\), leaving remainder %d.", divisor, q, divisor*q, rem), f)
	}
}

func w8AreaModel(kp int) generator.Problem {
	a := w8RR(2, 12)
	b := w8RR(2, 12)
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("A rectangle has side lengths %d and %d. What is its area? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("\\(%d\\times%d=%d\\).", a, b, a*b), f)
	case 2:
		return w8R(kp, fmt.Sprintf("A rectangle has area %d and one side length %d. What is the other side length? (enter a number)", a*b, a), strconv.Itoa(b),
			fmt.Sprintf("\\(%d/%d=%d\\).", a*b, a, b), f)
	default:
		return w8R(kp, fmt.Sprintf("What is the missing factor in \\(%d\\times\\square=%d\\)? (enter a number)", a, a*b), strconv.Itoa(b),
			fmt.Sprintf("The missing factor is \\(%d/%d=%d\\).", a*b, a, b), f)
	}
}

func w8Lattice(kp int) generator.Problem {
	t := w8RR(1, 8)
	u := w8RR(1, 9)
	v := w8RR(1, 8)
	w := w8RR(1, 9)
	a := 10*t + u
	b := 10*v + w
	cross := t*w + u*v
	units := u * w
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "cross", cross)
	w8RI(f, "units", units)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("For lattice multiplication of %d by %d, what is the cross-diagonal sum \\(t w+u v\\)? (enter a number)", a, b), strconv.Itoa(cross),
			fmt.Sprintf("The cross terms are \\(%d\\times%d+%d\\times%d=%d\\).", t, w, u, v, cross), f)
	case 2:
		return w8R(kp, fmt.Sprintf("For lattice multiplication of %d by %d, what is the units-diagonal sum \\(u w\\)? (enter a number)", a, b), strconv.Itoa(units),
			fmt.Sprintf("The units term is \\(%d\\times%d=%d\\).", u, w, units), f)
	default:
		return w8R(kp, fmt.Sprintf("Use lattice multiplication to compute \\(%d\\times%d\\). What is the product? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The lattice diagonals combine to \\(%d\\times%d=%d\\).", a, b, a*b), f)
	}
}

func w8SquareNumbers(kp int) generator.Problem {
	n := w8RR(2, 20)
	f := map[string]string{}
	w8RI(f, "n", n)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is \\(%d^2\\)? (enter a number)", n), strconv.Itoa(n*n),
			fmt.Sprintf("\\(%d\\times%d=%d\\).", n, n, n*n), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is \\(\\sqrt{%d}\\)? (enter a number)", n*n), strconv.Itoa(n),
			fmt.Sprintf("\\(%d\\times%d=%d\\), so the square root is %d.", n, n, n*n, n), f)
	default:
		return w8R(kp, fmt.Sprintf("What is the difference \\(( %d+1)^2-%d^2\\)? (enter a number)", n, n), strconv.Itoa(2*n+1),
			fmt.Sprintf("The difference of consecutive squares is \\(2n+1=2\\times%d+1=%d\\).", n, 2*n+1), f)
	}
}

func w8TriangularNumbers(kp int) generator.Problem {
	n := w8RR(2, 20)
	f := map[string]string{}
	w8RI(f, "n", n)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the %dth triangular number \\(T_{%d}\\)? (enter a number)", n, n), strconv.Itoa(n*(n+1)/2),
			fmt.Sprintf("\\(T_n=n(n+1)/2=%d\\times%d/2=%d\\).", n, n+1, n*(n+1)/2), f)
	case 2:
		return w8R(kp, fmt.Sprintf("The triangular number %d has which index? (enter a number)", n*(n+1)/2), strconv.Itoa(n),
			fmt.Sprintf("Solve \\(k(k+1)/2=%d\\). The positive integer solution is %d.", n*(n+1)/2, n), f)
	default:
		return w8R(kp, fmt.Sprintf("What is the sum of the first %d positive integers? (enter a number)", n), strconv.Itoa(n*(n+1)/2),
			fmt.Sprintf("The sum is the triangular number \\(T_{%d}=%d\\).", n, n*(n+1)/2), f)
	}
}

func w8Factorial(kp int) generator.Problem {
	n := w8RR(3, 8)
	fact := 1
	for i := 2; i <= n; i++ {
		fact *= i
	}
	f := map[string]string{}
	w8RI(f, "n", n)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is \\(%d!\\)? (enter a number)", n), strconv.Itoa(fact),
			fmt.Sprintf("\\(%d!=1\\times2\\times\\cdots\\times%d=%d\\).", n, n, fact), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is \\(%d!/(%d-1)!\\)? (enter a number)", n, n), strconv.Itoa(n),
			fmt.Sprintf("All factors except %d cancel, leaving %d.", n, n), f)
	default:
		return w8R(kp, fmt.Sprintf("What is \\(%d!/(%d-2)!\\)? (enter a number)", n, n), strconv.Itoa(n*(n-1)),
			fmt.Sprintf("The last two factors remain: \\(%d\\times%d=%d\\).", n, n-1, n*(n-1)), f)
	}
}

func w8GCFLCM(kp int) generator.Problem {
	a := w8RR(4, 24)
	b := w8RR(4, 24)
	g := w8RGCD(a, b)
	l := a / g * b
	f := map[string]string{}
	w8RI(f, "a", a)
	w8RI(f, "b", b)
	w8RI(f, "gcd", g)
	w8RI(f, "lcm", l)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("What is the greatest common factor of %d and %d? (enter a number)", a, b), strconv.Itoa(g),
			fmt.Sprintf("The largest integer dividing both %d and %d is %d.", a, b, g), f)
	case 2:
		return w8R(kp, fmt.Sprintf("What is the least common multiple of %d and %d? (enter a number)", a, b), strconv.Itoa(l),
			fmt.Sprintf("\\(\\operatorname{lcm}(%d,%d)=%d\\times%d/%d=%d\\).", a, b, a, b, g, l), f)
	default:
		return w8R(kp, fmt.Sprintf("For %d and %d, what is \\(\\gcd\\times\\operatorname{lcm}\\)? (enter a number)", a, b), strconv.Itoa(a*b),
			fmt.Sprintf("The product of the gcd and lcm is the product of the numbers: \\(%d\\times%d=%d\\).", a, b, a*b), f)
	}
}

func w8Elapsed(kp int) generator.Problem {
	start := w8RR(0, 20) * 60
	duration := w8RR(20, 180)
	end := start + duration
	f := map[string]string{}
	w8RI(f, "start", start)
	w8RI(f, "duration", duration)
	w8RI(f, "end", end)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("An event starts at %02d:%02d and ends at %02d:%02d. How many minutes does it last? (enter a number)", start/60, start%60, end/60, end%60), strconv.Itoa(duration),
			fmt.Sprintf("From %02d:%02d to %02d:%02d is %d minutes.", start/60, start%60, end/60, end%60, duration), f)
	case 2:
		return w8R(kp, fmt.Sprintf("An event starts at %02d:%02d and lasts %d minutes. What is the end time in minutes after midnight? (enter a number)", start/60, start%60, duration), strconv.Itoa(end),
			fmt.Sprintf("\\(%d+%d=%d\\) minutes after midnight.", start, duration, end), f)
	default:
		return w8R(kp, fmt.Sprintf("An event ends at %02d:%02d after lasting %d minutes. What was the start time in minutes after midnight? (enter a number)", end/60, end%60, duration), strconv.Itoa(start),
			fmt.Sprintf("\\(%d-%d=%d\\) minutes after midnight.", end, duration, start), f)
	}
}

func w8Change(kp int) generator.Problem {
	price := w8RR(1, 9)*100 + w8RR(0, 99)
	paid := price + w8RR(1, 20)*25
	change := paid - price
	f := map[string]string{}
	w8RI(f, "price", price)
	w8RI(f, "paid", paid)
	w8RI(f, "change", change)
	switch kp {
	case 1:
		return w8R(kp, fmt.Sprintf("An item costs $%d.%02d and the customer pays $%d.%02d. What is the change in cents? (enter a number)", price/100, price%100, paid/100, paid%100), strconv.Itoa(change),
			fmt.Sprintf("\\(%d-%d=%d\\) cents.", paid, price, change), f)
	case 2:
		return w8R(kp, fmt.Sprintf("The change is %d cents. How many quarters is that? (enter a number)", change), strconv.Itoa(change/25),
			fmt.Sprintf("\\(%d/25=%d\\) quarters, with %d cents left over.", change, change/25, change%25), f)
	default:
		return w8R(kp, fmt.Sprintf("An item costs $%d.%02d and the change is %d cents. How much did the customer pay, in cents? (enter a number)", price/100, price%100, change), strconv.Itoa(paid),
			fmt.Sprintf("\\(%d+%d=%d\\) cents.", price, change, paid), f)
	}
}
