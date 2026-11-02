package prealgebra

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// prealgBridgeGen serves the concepts added by the 2026 content expansion.
// Each kind switches on the knowledge-point ordinal so the question matches
// the section the learner is reading.
type prealgBridgeGen struct{ kind string }

func prealgOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func prealgProb(o int, q, a, e string) generator.Problem {
	return generator.Problem{
		Question:    q,
		Answer:      a,
		Explanation: e,
		Facts:       map[string]string{"answer": a},
		KPOrdinal:   o,
	}
}

func prealgDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 4, 64)
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

func (g *prealgBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := prealgOrdinal(ctx)
	switch g.kind {
	case "ineq_negative":
		k := 1 + rand.Intn(8)
		a := 2 + rand.Intn(8)
		switch o {
		case 1:
			return prealgProb(o, fmt.Sprintf("Solve for x: -%d x = %d. (enter a number)", a, a*k), strconv.Itoa(-k), fmt.Sprintf("Divide both sides by -%d: x = -%d.", a, k))
		case 2:
			return prealgProb(o, fmt.Sprintf("Solve for x: -%d x < %d. What is the boundary value of x? (enter a number)", a, -a*k), strconv.Itoa(k), fmt.Sprintf("Dividing by -%d reverses the inequality: x > %d, so the boundary is %d.", a, k, k))
		default:
			return prealgProb(o, fmt.Sprintf("What is the greatest integer x that satisfies -%d x >= %d? (enter a number)", a, -a*k), strconv.Itoa(k), fmt.Sprintf("Dividing by -%d gives x <= %d, so the greatest integer is %d.", a, k, k))
		}
	case "eq_fractions":
		a := 2 + rand.Intn(6)
		k := 2 + rand.Intn(8)
		b := 2 + rand.Intn(6)
		switch o {
		case 1:
			return prealgProb(o, fmt.Sprintf("Solve for x: (1/%d) x = %d. (enter a number)", a, k), strconv.Itoa(a*k), fmt.Sprintf("Multiply both sides by %d: x = %d x %d = %d.", a, a, k, a*k))
		case 2:
			return prealgProb(o, fmt.Sprintf("Solve for x: x/%d = %d. (enter a number)", a, k), strconv.Itoa(a*k), fmt.Sprintf("Multiply both sides by %d: x = %d x %d = %d.", a, a, k, a*k))
		default:
			m := 1 + rand.Intn(6)
			num := m * b * b
			den := a
			g := ratioGCD(num, den)
			num /= g
			den /= g
			ans := strconv.Itoa(num)
			if den != 1 {
				ans = fmt.Sprintf("%d/%d", num, den)
			}
			return prealgProb(o, fmt.Sprintf("Solve for x: (%d/%d) x = %d. (enter a number)", a, b, m*b), ans, fmt.Sprintf("Multiply both sides by %d/%d: x = (%d/%d) x %d = %s.", b, a, b, a, m*b, ans))
		}
	case "prop_graph":
		k := 2 + rand.Intn(7)
		switch o {
		case 1:
			x := 1 + rand.Intn(9)
			return prealgProb(o, fmt.Sprintf("A proportional relationship has constant 3.5. What is y when x = %d? (enter a decimal)", x), prealgDecimal(3.5*float64(x)), fmt.Sprintf("Multiply the constant by the input: 3.5 x %d = %s.", x, prealgDecimal(3.5*float64(x))))
		case 2:
			x := 1 + rand.Intn(8)
			return prealgProb(o, fmt.Sprintf("A proportional graph passes through (%d, %d). What is the constant of proportionality? (enter a number)", x, k*x), strconv.Itoa(k), fmt.Sprintf("The constant is y/x = %d/%d = %d.", k*x, x, k))
		default:
			pts := [][2]int{{1, k}, {2, 2 * k}, {2, 2*k + 1}, {3, 3 * k}, {3, 3*k + 2}}
			perm := rand.Perm(len(pts))[:4]
			count := 0
			var parts []string
			for _, i := range perm {
				p := pts[i]
				parts = append(parts, fmt.Sprintf("(%d, %d)", p[0], p[1]))
				if p[1] == k*p[0] {
					count++
				}
			}
			return prealgProb(o, fmt.Sprintf("A proportional relationship has constant %d. How many of these points lie on its graph: %s? (enter a number)", k, strings.Join(parts, ", ")), strconv.Itoa(count), fmt.Sprintf("%d of the four points satisfy y = %d x.", count, k))
		}
	case "func_table":
		a := 2 + rand.Intn(6)
		b := 1 + rand.Intn(9)
		switch o {
		case 1:
			x := 1 + rand.Intn(9)
			return prealgProb(o, fmt.Sprintf("A function follows the rule y = %d x + %d. What is the output when the input is %d? (enter a number)", a, b, x), strconv.Itoa(a*x+b), fmt.Sprintf("Substitute the input: %d x %d + %d = %d.", a, x, b, a*x+b))
		case 2:
			x := 1 + rand.Intn(9)
			return prealgProb(o, fmt.Sprintf("A function follows the rule y = %d x + %d. What input gives the output %d? (enter a number)", a, b, a*x+b), strconv.Itoa(x), fmt.Sprintf("Solve %d x + %d = %d, so x = %d.", a, b, a*x+b, x))
		default:
			count := 0
			var rows []string
			for i := 0; i < 5; i++ {
				x := 1 + rand.Intn(7)
				y := a*x + b
				if i == 2 {
					y += 1 + rand.Intn(3)
				}
				rows = append(rows, fmt.Sprintf("(%d, %d)", x, y))
				if y == a*x+b {
					count++
				}
			}
			return prealgProb(o, fmt.Sprintf("A function follows the rule y = %d x + %d. How many of these rows match the rule: %s? (enter a number)", a, b, strings.Join(rows, ", ")), strconv.Itoa(count), fmt.Sprintf("%d of the five rows satisfy y = %d x + %d.", count, a, b))
		}
	case "graph_interpret":
		m := 1 + rand.Intn(6)
		b := rand.Intn(11) - 5
		switch o {
		case 1:
			x := rand.Intn(9) - 2
			return prealgProb(o, fmt.Sprintf("A line is given by y = %d x + %d. What is y when x = %d? (enter a number)", m, b, x), strconv.Itoa(m*x+b), fmt.Sprintf("Substitute the input: %d x %d + %d = %d.", m, x, b, m*x+b))
		case 2:
			k := rand.Intn(9) - 4
			return prealgProb(o, fmt.Sprintf("A line is given by y = %d x + %d. At what x does the line cross the horizontal axis? (enter a number)", m, -m*k), strconv.Itoa(k), fmt.Sprintf("Set y = 0: %d x + %d = 0, so x = %d.", m, -m*k, k))
		default:
			x1 := rand.Intn(5)
			x2 := x1 + 1 + rand.Intn(4)
			y1 := m*x1 + b
			y2 := m*x2 + b
			return prealgProb(o, fmt.Sprintf("A line passes through (%d, %d) and (%d, %d). What is its slope? (enter a number)", x1, y1, x2, y2), strconv.Itoa(m), fmt.Sprintf("The slope is the change in y over the change in x: (%d - %d)/(%d - %d) = %d.", y2, y1, x2, x1, m))
		}
	case "order_rational":
		base := 1 + rand.Intn(9)
		switch o {
		case 1:
			num := 2*base + 1
			return prealgProb(o, fmt.Sprintf("Which is larger, %d/%d or 0.%d? Enter 1 for the fraction or 2 for the decimal.", num, 2*base, base*10), "1", fmt.Sprintf("%d/%d = %s, which is greater than 0.%d.", num, 2*base, prealgDecimal(float64(num)/float64(2*base)), base))
		case 2:
			lo := float64(base)
			hi := lo + 2
			vals := []float64{lo - 0.5, lo + 0.25, lo + 1.5, hi - 0.25, hi + 0.5}
			count := 0
			var parts []string
			for _, v := range vals {
				parts = append(parts, prealgDecimal(v))
				if v > lo && v < hi {
					count++
				}
			}
			return prealgProb(o, fmt.Sprintf("How many of these numbers lie strictly between %d and %d: %s? (enter a number)", base, base+2, strings.Join(parts, ", ")), strconv.Itoa(count), fmt.Sprintf("%d of the five numbers are strictly between %d and %d.", count, base, base+2))
		default:
			threshold := base + 1
			vals := []string{fmt.Sprintf("-%d", base+2), fmt.Sprintf("%d", base), fmt.Sprintf("%d", base+3), prealgDecimal(float64(base) + 0.5), prealgDecimal(float64(base) - 0.5)}
			count := 0
			for _, v := range vals {
				f, _ := strconv.ParseFloat(v, 64)
				if int(f) < threshold {
					count++
				}
			}
			return prealgProb(o, fmt.Sprintf("How many of these numbers are less than %d: %s? (enter a number)", threshold, strings.Join(vals, ", ")), strconv.Itoa(count), fmt.Sprintf("%d of the five numbers are less than %d.", count, threshold))
		}
	case "place_value":
		w := 1 + rand.Intn(90)
		d1 := rand.Intn(10)
		d2 := rand.Intn(10)
		d3 := rand.Intn(10)
		switch o {
		case 1:
			places := []string{"tenths", "hundredths", "thousandths"}
			digits := []int{d1, d2, d3}
			i := rand.Intn(3)
			return prealgProb(o, fmt.Sprintf("In %d.%d%d%d, which digit is in the %s place? (enter a number)", w, d1, d2, d3, places[i]), strconv.Itoa(digits[i]), fmt.Sprintf("The digits after the point are tenths, hundredths and thousandths, so the %s digit is %d.", places[i], digits[i]))
		case 2:
			places := []string{"tenths", "hundredths", "thousandths"}
			digits := []int{d1, d2, d3}
			i := rand.Intn(3)
			vals := []string{"0.1", "0.01", "0.001"}
			v, _ := strconv.ParseFloat(vals[i], 64)
			return prealgProb(o, fmt.Sprintf("In %d.%d%d%d, what is the value of the digit in the %s place? (enter a decimal)", w, d1, d2, d3, places[i]), prealgDecimal(float64(digits[i])*v), fmt.Sprintf("The %s digit is %d and its place value is %s, so the value is %s.", places[i], digits[i], vals[i], prealgDecimal(float64(digits[i])*v)))
		default:
			return prealgProb(o, fmt.Sprintf("In %d.%d%d%d, how many digits are to the right of the decimal point? (enter a number)", w, d1, d2, d3), "3", "The decimal part has three digits: tenths, hundredths and thousandths.")
		}
	case "powers_ten":
		a := 1 + rand.Intn(90)
		dec := float64(a) / 10
		k := 1 + rand.Intn(3)
		switch o {
		case 1:
			val := dec * float64(intPow10(k))
			return prealgProb(o, fmt.Sprintf("Compute %s x 10^%d. (enter a decimal)", prealgDecimal(dec), k), prealgDecimal(val), fmt.Sprintf("Multiplying by 10^%d moves the decimal point %d place(s) right: %s.", k, k, prealgDecimal(val)))
		case 2:
			val := dec / float64(intPow10(k))
			return prealgProb(o, fmt.Sprintf("Compute %s / 10^%d. (enter a decimal)", prealgDecimal(dec), k), prealgDecimal(val), fmt.Sprintf("Dividing by 10^%d moves the decimal point %d place(s) left: %s.", k, k, prealgDecimal(val)))
		default:
			n := 1 + rand.Intn(8)
			return prealgProb(o, fmt.Sprintf("By how many places does the decimal point move when a number is multiplied by 10^%d? (enter a number)", n), strconv.Itoa(n), fmt.Sprintf("Each power of ten moves the point one place, so 10^%d moves it %d places.", n, n))
		}
	case "dimensional":
		switch o {
		case 1:
			feet := 2 + rand.Intn(9)
			return prealgProb(o, fmt.Sprintf("How many inches are in %d feet? (enter a number)", feet), strconv.Itoa(12*feet), fmt.Sprintf("Each foot is 12 inches, so 12 x %d = %d.", feet, 12*feet))
		case 2:
			hours := 1 + rand.Intn(5)
			return prealgProb(o, fmt.Sprintf("How many minutes are in %d hours? (enter a number)", hours), strconv.Itoa(60*hours), fmt.Sprintf("Each hour is 60 minutes, so 60 x %d = %d.", hours, 60*hours))
		default:
			hours := 1 + rand.Intn(10)
			return prealgProb(o, fmt.Sprintf("How many seconds are in %d hours? (enter a number)", hours), strconv.Itoa(3600*hours), fmt.Sprintf("Chain two factors: %d hours x 60 minutes/hour x 60 seconds/minute = %d seconds.", hours, 3600*hours))
		}
	case "ratio_table":
		a := 2 + rand.Intn(7)
		b := 2 + rand.Intn(7)
		k := 2 + rand.Intn(6)
		switch o {
		case 1:
			return prealgProb(o, fmt.Sprintf("A ratio table shows %d apples for %d dollars. At the same rate, how many dollars for %d apples? (enter a number)", a, b, a*k), strconv.Itoa(b*k), fmt.Sprintf("Multiply both columns by %d: %d apples cost %d dollars.", k, a*k, b*k))
		case 2:
			return prealgProb(o, fmt.Sprintf("A ratio table shows %d items for %d dollars. What is the cost per item? (enter a decimal)", a, b), prealgDecimal(float64(b)/float64(a)), fmt.Sprintf("Divide the cost by the count: %d/%d = %s.", b, a, prealgDecimal(float64(b)/float64(a))))
		default:
			pairs := [][2]int{{a, b}, {2 * a, 2 * b}, {2 * a, 2*b + 1}, {3 * a, 3 * b}, {3 * a, 3*b + 2}}
			perm := rand.Perm(len(pairs))[:4]
			count := 0
			var parts []string
			for _, i := range perm {
				p := pairs[i]
				parts = append(parts, fmt.Sprintf("(%d, %d)", p[0], p[1]))
				if p[0]*b == p[1]*a {
					count++
				}
			}
			return prealgProb(o, fmt.Sprintf("The ratio %d to %d is given. How many of these pairs are equivalent to it: %s? (enter a number)", a, b, strings.Join(parts, ", ")), strconv.Itoa(count), fmt.Sprintf("%d of the four pairs satisfy the same ratio.", count))
		}
	case "simplify_expr":
		a := 2 + rand.Intn(6)
		b := 2 + rand.Intn(6)
		c := 1 + rand.Intn(9)
		switch o {
		case 1:
			return prealgProb(o, fmt.Sprintf("Simplify %d x + %d x. What is the coefficient of x? (enter a number)", a, b), strconv.Itoa(a+b), fmt.Sprintf("Like terms add their coefficients: %d + %d = %d.", a, b, a+b))
		case 2:
			return prealgProb(o, fmt.Sprintf("Expand %d (x + %d). What is the constant term? (enter a number)", a, c), strconv.Itoa(a*c), fmt.Sprintf("The factor multiplies both terms, so the constant is %d x %d = %d.", a, c, a*c))
		default:
			x := 1 + rand.Intn(7)
			return prealgProb(o, fmt.Sprintf("Simplify %d x + %d x, then evaluate the result at x = %d. (enter a number)", a, b, x), strconv.Itoa((a+b)*x), fmt.Sprintf("The simplified expression is %d x, and at x = %d it equals %d.", a+b, x, (a+b)*x))
		}
	}
	return prealgProb(o, "Simplify 2 x + 3 x. What is the coefficient of x? (enter a number)", "5", "Like terms add their coefficients: 2 + 3 = 5.")
}

func intPow10(n int) int {
	p := 1
	for i := 0; i < n; i++ {
		p *= 10
	}
	return p
}
