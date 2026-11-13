package fractions

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// fracBridgeGen serves the fraction concepts added by the 2026 content
// expansion. Each knowledge point asks a different question, and the answer is
// always an exact whole number or a reduced fraction.
type fracBridgeGen struct{ kind string }

func fracOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func fracProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func fracAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func fracText(n, d int) string {
	if d < 0 {
		n, d = -n, -d
	}
	g := mathutil.GCD(fracAbs(n), d)
	if g > 1 {
		n, d = n/g, d/g
	}
	if d == 1 {
		return strconv.Itoa(n)
	}
	return fmt.Sprintf("%d/%d", n, d)
}

func fracValue(n, d int) string {
	return fracText(n, d)
}

func (g *fracBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := fracOrdinal(ctx)
	switch g.kind {
	case "negative":
		switch o {
		case 1:
			b := 2 + rand.Intn(8)
			d := 2 + rand.Intn(8)
			for d == b {
				d = 2 + rand.Intn(8)
			}
			a := 1 + rand.Intn(8)
			c := 1 + rand.Intn(8)
			num := -a*d + c*b
			den := b * d
			return fracProb(o, fmt.Sprintf("What is minus %d/%d plus %d/%d? (enter a reduced fraction)", a, b, c, d), fracValue(num, den), fmt.Sprintf("Use the common denominator %d: minus %d/%d plus %d/%d = %d/%d = %s.", den, a, b, c, d, num, den, fracValue(num, den)))
		case 2:
			a := 1 + rand.Intn(8)
			b := 2 + rand.Intn(8)
			c := 1 + rand.Intn(8)
			d := 2 + rand.Intn(8)
			return fracProb(o, fmt.Sprintf("What is minus %d/%d times %d/%d? (enter a reduced fraction)", a, b, c, d), fracValue(-a*c, b*d), fmt.Sprintf("Multiply the magnitudes and count one negative factor: %d/%d.", a*c, b*d))
		default:
			terms := make([]string, 4)
			count := 0
			for i := range terms {
				a := 1 + rand.Intn(7)
				b := 2 + rand.Intn(7)
				c := 1 + rand.Intn(7)
				d := 2 + rand.Intn(7)
				if d == b {
					d++
				}
				terms[i] = fmt.Sprintf("(-%d/%d) + %d/%d", a, b, c, d)
				if -a*d+c*b < 0 {
					count++
				}
			}
			return fracProb(o, fmt.Sprintf("How many of these expressions are negative: %s, %s, %s and %s? (enter a number)", terms[0], terms[1], terms[2], terms[3]), strconv.Itoa(count), fmt.Sprintf("Compare the two terms after a common denominator. Exactly %d of the four sums are negative.", count))
		}
	case "repeating":
		type period struct{ n, p int }
		table := []period{{3, 1}, {7, 6}, {9, 1}, {11, 2}, {13, 6}, {17, 16}, {19, 18}, {21, 6}, {23, 22}}
		switch o {
		case 1:
			t := table[rand.Intn(len(table))]
			return fracProb(o, fmt.Sprintf("How many digits are in the repeating block of the decimal for 1/%d? (enter a number)", t.n), strconv.Itoa(t.p), fmt.Sprintf("The remainders in the long division of 1/%d cycle with period %d.", t.n, t.p))
		case 2:
			denoms := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 16, 20, 25}
			chosen := make([]int, 4)
			count := 0
			for i := range chosen {
				chosen[i] = denoms[rand.Intn(len(denoms))]
				n := chosen[i]
				for n%2 == 0 {
					n /= 2
				}
				for n%5 == 0 {
					n /= 5
				}
				if n == 1 {
					count++
				}
			}
			return fracProb(o, fmt.Sprintf("How many of these fractions have terminating decimals: 1/%d, 1/%d, 1/%d and 1/%d? (enter a number)", chosen[0], chosen[1], chosen[2], chosen[3]), strconv.Itoa(count), fmt.Sprintf("A fraction terminates when the denominator has only factors 2 and 5. Exactly %d of the four do.", count))
		default:
			t := table[rand.Intn(len(table))]
			m := 2 + rand.Intn(5)
			return fracProb(o, fmt.Sprintf("How many digits are in the repeating block of the decimal for %d/%d? (enter a number)", m, t.n), strconv.Itoa(t.p), fmt.Sprintf("The numerator does not change the cycle length, so the block has %d digits.", t.p))
		}
	case "mixed_borrow":
		switch o {
		case 1:
			w := 3 + rand.Intn(7)
			m := 1 + rand.Intn(w-1)
			d := 2 + rand.Intn(7)
			f := 1 + rand.Intn(d-1)
			num := w*d - (m*d + f)
			return fracProb(o, fmt.Sprintf("What is %d minus %d %d/%d? (enter a reduced fraction)", w, m, f, d), fracValue(num, d), fmt.Sprintf("Borrow one whole: %d = %d %d/%d. Subtract to get %d/%d = %s.", w, w-1, d, d, num, d, fracValue(num, d)))
		case 2:
			w := 4 + rand.Intn(6)
			d := 3 + rand.Intn(5)
			f := 1 + rand.Intn(d-1)
			m := 1 + rand.Intn(w-2)
			e := 2 + rand.Intn(5)
			g2 := 1 + rand.Intn(e-1)
			for e == d && g2 >= f {
				e = 2 + rand.Intn(5)
				g2 = 1 + rand.Intn(e-1)
			}
			first := w*d + f
			second := m*e + g2
			num := first*e - second*d
			den := d * e
			return fracProb(o, fmt.Sprintf("What is %d %d/%d minus %d %d/%d? (enter a reduced fraction)", w, f, d, m, g2, e), fracValue(num, den), fmt.Sprintf("Convert both to improper fractions, subtract and reduce: %d/%d - %d/%d = %s.", first, d, second, e, fracValue(num, den)))
		default:
			w := 5 + rand.Intn(6)
			m := 2 + rand.Intn(w-2)
			d := 3 + rand.Intn(6)
			f := 1 + rand.Intn(d-1)
			num := w*d - (m*d + f)
			return fracProb(o, fmt.Sprintf("A board is %d metres long. A piece of %d %d/%d metres is cut off. How many metres remain? (enter a reduced fraction)", w, m, f, d), fracValue(num, d), fmt.Sprintf("Subtract with borrowing, or convert to improper fractions: the remainder is %s metres.", fracValue(num, d)))
		}
	case "estimate":
		switch o {
		case 1:
			b := 3 + rand.Intn(7)
			d := 3 + rand.Intn(7)
			a := 1 + rand.Intn(b-1)
			c := 1 + rand.Intn(d-1)
			sum := a*d + c*b
			den := b * d
			nearest := (2*sum + den) / (2 * den)
			return fracProb(o, fmt.Sprintf("Estimate %d/%d plus %d/%d by rounding each fraction to 0, 1/2 or 1, then give the nearest whole number. (enter a number)", a, b, c, d), strconv.Itoa(nearest), fmt.Sprintf("The exact sum is %s, whose nearest whole number is %d.", fracValue(sum, den), nearest))
		case 2:
			b := 3 + rand.Intn(7)
			d := 3 + rand.Intn(7)
			a := 1 + rand.Intn(b-1)
			c := 1 + rand.Intn(d-1)
			diff := a*d - c*b
			den := b * d
			nearest := (2*fracAbs(diff) + den) / (2 * den)
			return fracProb(o, fmt.Sprintf("Estimate %d/%d minus %d/%d by rounding each fraction to 0, 1/2 or 1, then give the nearest whole number. (enter a number)", a, b, c, d), strconv.Itoa(nearest), fmt.Sprintf("The difference is %s, whose size is nearest to %d.", fracValue(diff, den), nearest))
		default:
			vals := make([]int, 4)
			count := 0
			for i := range vals {
				den := 5 + rand.Intn(8)
				num := 1 + rand.Intn(den-1)
				vals[i] = num * 100 / den
				// Closer to 1 than to 1/2 means above the midpoint 3/4, not
				// merely above 1/2. Reports 454 and 460 both answered 57/100
				// and 1 correctly and were marked wrong by the old test.
				if 4*vals[i] > 300 {
					count++
				}
			}
			return fracProb(o, fmt.Sprintf("How many of these fractions are closer to 1 than to 1/2: %d/100, %d/100, %d/100 and %d/100? (enter a number)", vals[0], vals[1], vals[2], vals[3]), strconv.Itoa(count), fmt.Sprintf("A fraction is closer to 1 when it exceeds 3/4. Exactly %d of the four do.", count))
		}
	case "word_multi_step":
		switch o {
		case 1:
			a := 3 + rand.Intn(5)
			b := a + 1 + rand.Intn(4)
			num := a*b - b - a
			den := a * b
			return fracProb(o, fmt.Sprintf("A tank is filled 1/%d of the way, then another 1/%d. What fraction of the tank is still empty? (enter a reduced fraction)", a, b), fracValue(num, den), fmt.Sprintf("The filled part is %s, so the empty part is 1 - %s = %s.", fracValue(a+b, den), fracValue(a+b, den), fracValue(num, den)))
		case 2:
			a := 2 + rand.Intn(4)
			b := 2 + rand.Intn(5)
			num := (a - 1) * (b - 1)
			den := a * b
			return fracProb(o, fmt.Sprintf("A student spends 1/%d of her savings, then 1/%d of what remains. What fraction of the original savings is left? (enter a reduced fraction)", a, b), fracValue(num, den), fmt.Sprintf("After the first spend %d/%d remains. Taking 1/%d of that leaves %s.", a-1, a, b, fracValue(num, den)))
		default:
			a := 3 + rand.Intn(4)
			b := 3 + rand.Intn(4)
			num := (a - 1) * (b - 1)
			den := a * b
			return fracProb(o, fmt.Sprintf("A pipe fills 1/%d of a pool, a drain empties 1/%d of what is present, then no more water moves. What fraction of the pool is full? (enter a reduced fraction)", a, b), fracValue(num, den), fmt.Sprintf("After filling, 1/%d is present. Removing 1/%d of that leaves %s of the pool.", a, b, fracValue(num, den)))
		}
	case "div_by_whole":
		switch o {
		case 1:
			a := 1 + rand.Intn(8)
			b := 2 + rand.Intn(8)
			n := 2 + rand.Intn(8)
			return fracProb(o, fmt.Sprintf("What is %d/%d divided by %d? (enter a reduced fraction)", a, b, n), fracValue(a, b*n), fmt.Sprintf("Multiplying the denominator by %d gives %d/%d, which reduces to %s.", n, a, b*n, fracValue(a, b*n)))
		case 2:
			a := 1 + rand.Intn(8)
			b := 2 + rand.Intn(8)
			n := 2 + rand.Intn(8)
			return fracProb(o, fmt.Sprintf("The fraction %d/%d is split equally into shares of %s. How many shares are there? (enter a number)", a, b, fracValue(a, b*n)), strconv.Itoa(n), fmt.Sprintf("The share is the original fraction divided by %d, so there are %d shares.", n, n))
		default:
			a := 1 + rand.Intn(5)
			n := 2 + rand.Intn(5)
			return fracProb(o, fmt.Sprintf("A recipe uses 1/%d of a bag of flour and is made in %d equal batches. What fraction of the bag is in each batch? (enter a reduced fraction)", a, n), fracValue(1, a*n), fmt.Sprintf("Split 1/%d into %d parts: 1/(%d x %d) = 1/%d.", a, n, a, n, a*n))
		}
	case "compare_forms":
		switch o {
		case 1:
			vals := make([]string, 4)
			count := 0
			for i := range vals {
				switch rand.Intn(3) {
				case 0:
					den := 4 + rand.Intn(9)
					num := 1 + rand.Intn(den-1)
					vals[i] = fmt.Sprintf("%d/%d", num, den)
					if num*10 > 6*den {
						count++
					}
				case 1:
					v := 10 + rand.Intn(80)
					vals[i] = fmt.Sprintf("0.%02d", v)
					if v > 60 {
						count++
					}
				default:
					v := 20 + rand.Intn(70)
					vals[i] = fmt.Sprintf("%d%%", v)
					if v > 60 {
						count++
					}
				}
			}
			return fracProb(o, fmt.Sprintf("How many of these values are greater than 0.6: %s, %s, %s and %s? (enter a number)", vals[0], vals[1], vals[2], vals[3]), strconv.Itoa(count), fmt.Sprintf("Convert every value to a decimal. Exactly %d of the four exceed 0.6.", count))
		case 2:
			p := 10 + rand.Intn(80)
			return fracProb(o, fmt.Sprintf("Write %d percent as a fraction over 100, reduced. What is the numerator? (enter a number)", p), strconv.Itoa(p/mathutil.GCD(p, 100)), fmt.Sprintf("%d/100 reduces to %s, whose numerator is %d.", p, fracValue(p, 100), p/mathutil.GCD(p, 100)))
		default:
			b := []int{4, 5, 10, 20, 25}[rand.Intn(5)]
			a := 1 + rand.Intn(b-1)
			d := 50 + rand.Intn(40)
			diff := a*100 - d*b
			return fracProb(o, fmt.Sprintf("By how many hundredths is %d/%d greater than 0.%02d? (enter a number)", a, b, d), strconv.Itoa(diff/b), fmt.Sprintf("Both values in hundredths: %d/%d is %d hundredths and 0.%02d is %d hundredths, a difference of %d hundredths.", a, b, a*100/b, d, d, diff/b))
		}
	}
	return generator.Problem{}
}
