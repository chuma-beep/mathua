package calculus

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// calc2BridgeGen serves the calculus II concepts added by the content
// expansion. The kind selects the concept and KPOrdinal selects the lesson
// section.
type calc2BridgeGen struct{ kind string }

func (g *calc2BridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := calcBridgeOrdinal(ctx)
	switch g.kind {
	case "maclaurin":
		switch o {
		case 1:
			n := 3 + rand.Intn(5)
			return calcBridgeProblem(o,
				fmt.Sprintf("The Maclaurin series for \\(e^x\\) has coefficient \\(1/n!\\). What is \\(n!\\) when \\(n=%d\\)? (enter a number)", n),
				strconv.Itoa(calcFactorial(n)),
				fmt.Sprintf("The coefficient is 1/%d!, and %d! = %d.", n, n, calcFactorial(n)),
				map[string]string{"kind": "factorial", "n": strconv.Itoa(n)})
		case 2:
			n := []int{3, 5, 7, 9, 11}[rand.Intn(5)]
			return calcBridgeProblem(o,
				fmt.Sprintf("In the Maclaurin series for \\(\\sin x\\), the term containing \\(x^{%d}\\) has denominator \\(n!\\). What is \\(n!\\)? (enter a number)", n),
				strconv.Itoa(calcFactorial(n)),
				fmt.Sprintf("The odd powers follow the factorial denominators, and %d! = %d.", n, calcFactorial(n)),
				map[string]string{"kind": "sine-denominator", "n": strconv.Itoa(n)})
		default:
			k := 1 + rand.Intn(6)
			v := 1 << k
			return calcBridgeProblem(o,
				fmt.Sprintf("Substitute \\(2x\\) into \\(1/(1-x)\\). What is the coefficient of \\(x^{%d}\\) in the resulting series? (enter a number)", k),
				strconv.Itoa(v),
				fmt.Sprintf("The series becomes the sum of (2x)^m, so the coefficient of x^%d is 2^%d = %d.", k, k, v),
				map[string]string{"kind": "substitution", "k": strconv.Itoa(k)})
		}
	case "inverse_trig":
		switch o {
		case 1:
			h := rand.Intn(5) + 1
			k := 1 + rand.Intn(4)
			b := 2 * h
			c := h*h + k
			return calcBridgeProblem(o,
				fmt.Sprintf("Write \\(x^2+%dx+%d\\) as \\((x+h)^2+k\\). What is \\(k\\)? (enter a number)", b, c),
				strconv.Itoa(k),
				fmt.Sprintf("A perfect square uses %d, so the constant becomes %d - %d = %d.", h, c, h*h, k),
				map[string]string{"kind": "complete-square", "k": strconv.Itoa(k)})
		case 2:
			pairs := [][2]string{{"1/2", "30"}, {"\\sqrt{3}/2", "60"}, {"1/\\sqrt{2}", "45"}, {"0", "0"}, {"1", "90"}}
			p := pairs[rand.Intn(len(pairs))]
			return calcBridgeProblem(o,
				fmt.Sprintf("Using \\(x=a\\sin\\theta\\), at what angle \\(\\theta\\) in degrees is \\(x/a=%s\\)? (enter a number)", p[0]),
				p[1],
				fmt.Sprintf("The sine of the angle equals %s, so theta is %s degrees.", p[0], p[1]),
				map[string]string{"kind": "angle", "angle": p[1]})
		default:
			as := []int{1, 2, 4, 5, 10}
			a := as[rand.Intn(len(as))]
			v := math.Pi / (4 * float64(a))
			return calcBridgeProblem(o,
				fmt.Sprintf("Evaluate \\(\\int_0^{%d} \\frac{dx}{x^2+%d}\\) exactly, then enter the decimal value. (enter a number)", a, a*a),
				calcBridgeDecimal(v),
				fmt.Sprintf("The integral is (1/a)arctan(x/a) evaluated at a, which equals pi/(4a) = %s.", calcBridgeDecimal(v)),
				map[string]string{"kind": "arctan-definite", "a": strconv.Itoa(a)})
		}
	case "probability_density":
		switch o {
		case 1:
			bs := []int{2, 4, 5, 8, 10}
			b := bs[rand.Intn(len(bs))]
			v := 1 / float64(b)
			return calcBridgeProblem(o,
				fmt.Sprintf("A uniform density is constant on the interval from 0 to %d and zero elsewhere. What is its height? (enter a decimal)", b),
				calcBridgeDecimal(v),
				fmt.Sprintf("The area must be one, so the height is 1/%d = %s.", b, calcBridgeDecimal(v)),
				map[string]string{"kind": "uniform-height", "b": strconv.Itoa(b)})
		case 2:
			lam := 1 + rand.Intn(3)
			t := 1 + rand.Intn(2)
			v := math.Exp(-float64(lam) * float64(t))
			return calcBridgeProblem(o,
				fmt.Sprintf("An exponential density with rate %d has survival \\(P(X>t)=e^{-\\lambda t}\\). What is \\(P(X>%d)\\)? (enter a decimal rounded to 4 places)", lam, t),
				calcBridgeDecimal(v),
				fmt.Sprintf("Substitute the rate and time: e^(-%d*%d) = %s.", lam, t, calcBridgeDecimal(v)),
				map[string]string{"kind": "exponential-tail", "lambda": strconv.Itoa(lam), "t": strconv.Itoa(t)})
		default:
			type rule struct {
				k   int
				pct string
			}
			rules := []rule{{1, "68"}, {2, "95"}, {3, "99.7"}}
			r := rules[rand.Intn(len(rules))]
			outside := rand.Intn(2) == 1
			if outside {
				tails := map[int]string{1: "32", 2: "5", 3: "0.3"}
				return calcBridgeProblem(o,
					fmt.Sprintf("A normal distribution follows the empirical rule. What percent of observations lie more than %d standard deviations from the mean? (enter a number)", r.k),
					tails[r.k],
					fmt.Sprintf("The rule leaves 100 - %s = %s percent outside %d standard deviations.", r.pct, tails[r.k], r.k),
					map[string]string{"kind": "empirical-outside", "k": strconv.Itoa(r.k)})
			}
			return calcBridgeProblem(o,
				fmt.Sprintf("A normal distribution follows the empirical rule. What percent of observations lie within %d standard deviations of the mean? (enter a number)", r.k),
				r.pct,
				fmt.Sprintf("The empirical rule gives %s percent within %d standard deviations.", r.pct, r.k),
				map[string]string{"kind": "empirical-inside", "k": strconv.Itoa(r.k)})
		}
	case "series_estimate":
		switch o {
		case 1:
			ns := []int{10, 20, 50, 100, 1000}
			n := ns[rand.Intn(len(ns))]
			v := 1 / float64(n)
			return calcBridgeProblem(o,
				fmt.Sprintf("An alternating series has terms of size \\(1/n\\). After summing through \\(n=%d\\), the remainder is at most the next term. What is that bound? (enter a decimal)", n),
				calcBridgeDecimal(v),
				fmt.Sprintf("The next term is 1/%d = %s, which bounds the remainder.", n+1, calcBridgeDecimal(1/float64(n+1))),
				map[string]string{"kind": "alternating-remainder", "n": strconv.Itoa(n)})
		case 2:
			ns := []int{10, 20, 50, 100, 1000}
			n := ns[rand.Intn(len(ns))]
			v := 1 / float64(n)
			return calcBridgeProblem(o,
				fmt.Sprintf("For the p-series with \\(p=2\\), the tail after \\(N=%d\\) is bounded above by \\(1/N\\). What is that upper bound? (enter a decimal)", n),
				calcBridgeDecimal(v),
				fmt.Sprintf("The integral test gives 1/N = 1/%d = %s as the upper bound.", n, calcBridgeDecimal(v)),
				map[string]string{"kind": "integral-remainder", "n": strconv.Itoa(n)})
		default:
			type tol struct {
				t string
				n int
			}
			tols := []tol{{"0.05", 4}, {"0.01", 5}, {"0.005", 5}, {"0.001", 6}, {"0.0001", 7}}
			x := tols[rand.Intn(len(tols))]
			return calcBridgeProblem(o,
				fmt.Sprintf("The Taylor remainder for \\(e^x\\) at \\(x=1\\) is bounded by \\(e/(n+1)!\\). What is the smallest \\(n\\) that makes the bound below %s? (enter a number)", x.t),
				strconv.Itoa(x.n),
				fmt.Sprintf("Evaluate e/(n+1)! until it first falls below %s; the first such n is %d.", x.t, x.n),
				map[string]string{"kind": "taylor-terms", "n": strconv.Itoa(x.n)})
		}
	case "integral_log":
		switch o {
		case 1:
			ks := []int{2, 3, 4, 5, 10}
			k := ks[rand.Intn(len(ks))]
			v := 1 / float64(k)
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\int \\frac{dx}{%dx}\\), the antiderivative is \\(c\\ln|x|\\). What is \\(c\\)? (enter a decimal)", k),
				calcBridgeDecimal(v),
				fmt.Sprintf("Pulling out the constant gives (1/%d) ln|x|, so c = %s.", k, calcBridgeDecimal(v)),
				map[string]string{"kind": "log-coefficient", "k": strconv.Itoa(k)})
		case 2:
			n := 1 + rand.Intn(6)
			d := (n + 1) * (n + 1)
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\int x^{%d}\\ln x\\,dx\\), the second term is \\(-x^{%d}/d\\). What is \\(d\\)? (enter a number)", n, n+1),
				strconv.Itoa(d),
				fmt.Sprintf("Integration by parts gives denominator (n+1)^2 = %d^2 = %d.", n+1, d),
				map[string]string{"kind": "log-parts", "n": strconv.Itoa(n)})
		default:
			m := 1 + rand.Intn(6)
			v := 1 / float64(m+1)
			return calcBridgeProblem(o,
				fmt.Sprintf("For \\(\\int \\frac{(\\ln x)^{%d}}{x}\\,dx\\), the antiderivative is \\(c(\\ln x)^{%d}\\). What is \\(c\\)? (enter a decimal)", m, m+1),
				calcBridgeDecimal(v),
				fmt.Sprintf("The power rule on ln x gives coefficient 1/(%d+1) = %s.", m, calcBridgeDecimal(v)),
				map[string]string{"kind": "log-power", "m": strconv.Itoa(m)})
		}
	case "rationalizing":
		switch o {
		case 1:
			k := 2 + rand.Intn(5)
			return calcBridgeProblem(o,
				fmt.Sprintf("Substitute \\(t=\\sqrt[%d]{x+a}\\), so \\(x=t^{%d}-a\\). What is the exponent of \\(t\\) in \\(dx\\)? (enter a number)", k, k),
				strconv.Itoa(k-1),
				fmt.Sprintf("Differentiating t^%d gives %d t^%d, so the exponent is %d.", k, k, k-1, k-1),
				map[string]string{"kind": "root-power", "k": strconv.Itoa(k)})
		case 2:
			bs := []int{4, 9, 16, 25, 36}
			b := bs[rand.Intn(len(bs))]
			v := 2 * (int(math.Sqrt(float64(b))) - 1)
			return calcBridgeProblem(o,
				fmt.Sprintf("Evaluate \\(\\int_1^{%d} \\frac{dx}{\\sqrt{x}}\\). (enter a number)", b),
				strconv.Itoa(v),
				fmt.Sprintf("The antiderivative is 2 sqrt(x); at %d and 1 the difference is %d.", b, v),
				map[string]string{"kind": "sqrt-integral", "b": strconv.Itoa(b)})
		default:
			bs := []int{1, 4, 9, 16, 25}
			b := bs[rand.Intn(len(bs))]
			u := math.Sqrt(float64(b))
			v := 2 * (u - math.Log(1+u))
			return calcBridgeProblem(o,
				fmt.Sprintf("With \\(u=\\sqrt{x}\\), \\(\\int_0^{%d} \\frac{dx}{1+\\sqrt{x}}\\) becomes \\(\\int_0^{%d} \\frac{2u}{1+u}\\,du\\). Enter its decimal value. (enter a number)", b, int(u)),
				calcBridgeDecimal(v),
				fmt.Sprintf("The antiderivative is 2(u - ln(1+u)); at u=%d this is %s.", int(u), calcBridgeDecimal(v)),
				map[string]string{"kind": "rationalized-definite", "b": strconv.Itoa(b)})
		}
	}
	return calcBridgeProblem(o, "Evaluate the expression. (enter a number)", "0", "The expression evaluates to zero.", nil)
}

func (g *calc2BridgeGen) String() string { return "calc2BridgeGen(" + g.kind + ")" }
