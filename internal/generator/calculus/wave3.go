package calculus

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

// ---------------------------------------------------------------------------
// Shared exact-arithmetic helpers for the Wave 2 Calculus II batch.
//
// The answers here are integers or reduced fractions. A small rational type
// keeps the generator and its explanation from rounding, and the test rebuilds
// the same value through math/big.Rat rather than through this type.
// ---------------------------------------------------------------------------

type w3Frac struct {
	num int
	den int
}

func w3AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func w3MakeFrac(num, den int) w3Frac {
	if den == 0 {
		return w3Frac{num: 0, den: 1}
	}
	if den < 0 {
		num, den = -num, -den
	}
	g := mathutil.GCD(w3AbsInt(num), den)
	if g == 0 {
		g = 1
	}
	return w3Frac{num: num / g, den: den / g}
}

func (f w3Frac) add(o w3Frac) w3Frac {
	return w3MakeFrac(f.num*o.den+o.num*f.den, f.den*o.den)
}

func (f w3Frac) sub(o w3Frac) w3Frac {
	return w3MakeFrac(f.num*o.den-o.num*f.den, f.den*o.den)
}

func (f w3Frac) mul(o w3Frac) w3Frac {
	return w3MakeFrac(f.num*o.num, f.den*o.den)
}

func (f w3Frac) div(o w3Frac) w3Frac {
	return w3MakeFrac(f.num*o.den, f.den*o.num)
}

func (f w3Frac) String() string {
	if f.den == 1 {
		return strconv.Itoa(f.num)
	}
	return fmt.Sprintf("%d/%d", f.num, f.den)
}

func w3FracTex(num, den int) string {
	f := w3MakeFrac(num, den)
	if f.den == 1 {
		return strconv.Itoa(f.num)
	}
	return fmt.Sprintf("\\frac{%d}{%d}", f.num, f.den)
}

func w3Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	res := 1
	for i := 1; i <= k; i++ {
		res = res * (n - k + i) / i
	}
	return res
}

func w3Factorial(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		res *= i
	}
	return res
}

func w3Fixed(x float64, places int) string {
	return strconv.FormatFloat(x, 'f', places, 64)
}

func w3Round(x float64, places int) float64 {
	scale := math.Pow(10, float64(places))
	return math.Round(x*scale) / scale
}

func w3FloatFactorial(n int) float64 {
	res := 1.0
	for i := 2; i <= n; i++ {
		res *= float64(i)
	}
	return res
}

// w3FresnelSeries sums the Maclaurin series of the Fresnel cosine or sine
// integral. The caller states that the answer is an approximation and how many
// terms were used.
func w3FresnelSeries(x float64, sine bool) float64 {
	sum := 0.0
	sign := 1.0
	for n := 0; n < 30; n++ {
		var term float64
		if sine {
			term = sign * math.Pow(math.Pi/2, float64(2*n+1)) * math.Pow(x, float64(4*n+3)) /
				(w3FloatFactorial(2*n+1) * float64(4*n+3))
		} else {
			term = sign * math.Pow(math.Pi/2, float64(2*n)) * math.Pow(x, float64(4*n+1)) /
				(w3FloatFactorial(2*n) * float64(4*n+1))
		}
		sum += term
		sign = -sign
	}
	return sum
}

func w3EllipticCoefficient(n int) float64 {
	num := w3FloatFactorial(2 * n)
	den := math.Pow(2, float64(2*n)) * math.Pow(w3FloatFactorial(n), 2)
	a := num / den
	return a * a
}

func w3EllipticK(k float64) float64 {
	sum := 0.0
	for n := 0; n < 30; n++ {
		sum += w3EllipticCoefficient(n) * math.Pow(k, float64(2*n))
	}
	return math.Pi / 2 * sum
}

func w3EllipticE(k float64) float64 {
	sum := 1.0
	for n := 1; n < 30; n++ {
		sum -= w3EllipticCoefficient(n) * math.Pow(k, float64(2*n)) / float64(2*n-1)
	}
	return math.Pi / 2 * sum
}

func w3AGM(a, b float64) float64 {
	for i := 0; i < 100; i++ {
		nextA := (a + b) / 2
		nextB := math.Sqrt(a * b)
		if math.Abs(nextA-nextB) < 1e-15 {
			return (nextA + nextB) / 2
		}
		a, b = nextA, nextB
	}
	return (a + b) / 2
}

func w3Choice(kp int, question, correct, wrong, explanation string, facts map[string]string) generator.Problem {
	p := generator.MultipleChoice(question, correct, wrong)
	p.Explanation = explanation
	p.Facts = facts
	p.KPOrdinal = kp
	return p
}

// ---------------------------------------------------------------------------
// Cauchy condensation test
// ---------------------------------------------------------------------------

type w3CauchyCondensationGen struct{}

func (g *w3CauchyCondensationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	p := rand.Intn(30) + 2 // 2..31
	switch kp {
	case 1:
		exponent := 1 - p
		question := fmt.Sprintf("For \\(a_n=1/n^{%d}\\), Cauchy condensation gives the term \\(b_n=2^n a_{2^n}=2^{k n}\\). What is \\(k\\)?", p)
		explanation := fmt.Sprintf("Substitute \\(a_{2^n}=1/(2^n)^{%d}=2^{-%d n}\\). Then \\(b_n=2^n\\cdot2^{-%d n}=2^{%d n}\\), so \\(k=%d\\).", p, p, p, 1-p, 1-p)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(exponent),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "exponent": strconv.Itoa(exponent)},
			KPOrdinal:   kp,
		}
	case 2:
		m := mathutil.IntPow(2, p-1)
		question := fmt.Sprintf("For \\(a_n=1/n^{%d}\\), the condensed series is geometric with ratio \\(1/m\\). What is \\(m\\)?", p)
		explanation := fmt.Sprintf("The condensed term is \\(2^n a_{2^n}=2^n/2^{%d n}=2^{%d n}\\), whose ratio is \\(2^{%d}=1/%d\\). Hence \\(m=%d\\).", p, 1-p, 1-p, m, m)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(m),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "ratio_denominator": strconv.Itoa(m)},
			KPOrdinal:   kp,
		}
	default:
		num := mathutil.IntPow(2, p-1)
		den := num - 1
		sum := w3MakeFrac(num, den)
		question := fmt.Sprintf("For \\(a_n=1/n^{%d}\\), the condensed series is geometric. What is its sum \\(S\\)? (enter a fraction or an integer)", p)
		explanation := fmt.Sprintf("The first condensed term is \\(b_0=a_1=1\\) and the ratio is \\(r=2^{%d}=1/%d\\). The geometric sum is \\(1/(1-r)=%s\\).", 1-p, num, sum.String())
		return generator.Problem{
			Question:    question,
			Answer:      sum.String(),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "numerator": strconv.Itoa(num), "denominator": strconv.Itoa(den)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Raabe's test
// ---------------------------------------------------------------------------

type w3RaabeGen struct{}

func (g *w3RaabeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		p := rand.Intn(20) + 1
		question := fmt.Sprintf("For the positive-term series with \\(a_n=1/n^{%d}\\), compute Raabe's limit \\(L=\\lim_{n\\to\\infty} n(1-a_{n+1}/a_n)\\).", p)
		explanation := fmt.Sprintf("\\(a_{n+1}/a_n=(n/(n+1))^{%d}=(1+1/n)^{-%d}=1-%d/n+O(1/n^2)\\). Thus \\(n(1-a_{n+1}/a_n)\\to %d\\), so \\(L=%d\\).", p, p, p, p, p)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(p),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "q": "0"},
			KPOrdinal:   kp,
		}
	case 2:
		p := rand.Intn(5) + 2
		q := rand.Intn(2)
		limit := p - q
		term := fmt.Sprintf("1/n^{%d}", p)
		if q == 1 {
			term = fmt.Sprintf("(n+1)/n^{%d}", p)
		}
		question := fmt.Sprintf("For the positive-term series with \\(a_n=%s\\), compute Raabe's limit \\(L=\\lim_{n\\to\\infty} n(1-a_{n+1}/a_n)\\).", term)
		explanation := fmt.Sprintf("\\(a_{n+1}/a_n=(1+1/n)^{-%d}\\), so \\(n(1-a_{n+1}/a_n)\\to %d\\). Therefore \\(L=%d\\).", p, p, limit)
		if q == 1 {
			explanation = fmt.Sprintf("Write \\(a_n=(n+1)/n^{%d}\\). Then \\(a_{n+1}/a_n=\\frac{n+2}{n+1}(1+1/n)^{-%d}=1-%d/n+O(1/n^2)\\), so \\(n(1-a_{n+1}/a_n)\\to %d\\). Therefore \\(L=%d\\).", p, p, p-1, limit, limit)
		}
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "q": strconv.Itoa(q), "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	default:
		p := rand.Intn(20)
		limit := p + 1
		question := fmt.Sprintf("For the positive-term series with \\(a_n=1/(n^{%d}(n+1))\\), compute Raabe's limit \\(L=\\lim_{n\\to\\infty} n(1-a_{n+1}/a_n)\\).", p)
		explanation := fmt.Sprintf("\\(a_{n+1}/a_n=(1+1/n)^{%d}(n+1)/(n+2)\\). The first factor contributes \\(%d/n\\) and the second contributes \\(1/n\\), so \\(L=%d+1=%d\\).", p, p, p, limit)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Bertrand's test
// ---------------------------------------------------------------------------

type w3BertrandGen struct{}

func (g *w3BertrandGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		p := rand.Intn(20) + 1
		question := fmt.Sprintf("For \\(a_n=1/(n(\\ln n)^{%d})\\) and \\(n\\ge2\\), compute Bertrand's limit \\(L=\\lim_{n\\to\\infty}\\ln n\\,[n(1-a_{n+1}/a_n)-1]\\).", p)
		explanation := fmt.Sprintf("Expanding \\(a_{n+1}/a_n\\) gives \\(n(1-a_{n+1}/a_n)=1+%d/\\ln n+O(1/n)\\). Subtracting 1 and multiplying by \\(\\ln n\\) leaves %d, so \\(L=%d\\).", p, p, p)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(p),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "form": "log_denominator"},
			KPOrdinal:   kp,
		}
	case 2:
		p := rand.Intn(20) + 1
		limit := -p
		question := fmt.Sprintf("For \\(a_n=(\\ln n)^{%d}/n\\) and \\(n\\ge2\\), compute Bertrand's limit \\(L=\\lim_{n\\to\\infty}\\ln n\\,[n(1-a_{n+1}/a_n)-1]\\).", p)
		explanation := fmt.Sprintf("Here \\(n(1-a_{n+1}/a_n)=1-%d/\\ln n+O(1/n)\\). Subtracting 1 and multiplying by \\(\\ln n\\) gives \\(-%d\\), so \\(L=%d\\).", p, p, limit)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "form": "log_numerator", "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	default:
		p := rand.Intn(5) + 1
		c := rand.Intn(3) + 2
		question := fmt.Sprintf("For \\(a_n=%d/(n(\\ln n)^{%d})\\) and \\(n\\ge2\\), compute Bertrand's limit \\(L=\\lim_{n\\to\\infty}\\ln n\\,[n(1-a_{n+1}/a_n)-1]\\).", c, p)
		explanation := fmt.Sprintf("The constant factor %d cancels in \\(a_{n+1}/a_n\\). The logarithmic denominator contributes \\(%d/\\ln n\\) to \\(n(1-a_{n+1}/a_n)\\), so after subtracting 1 and multiplying by \\(\\ln n\\) the limit is %d.", c, p, p)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(p),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "c": strconv.Itoa(c), "form": "scaled_log_denominator"},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Kummer's test
// ---------------------------------------------------------------------------

type w3KummerGen struct{}

func (g *w3KummerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		p := rand.Intn(5) + 2
		c := rand.Intn(4) + 1
		limit := c * (p - 1)
		bn := fmt.Sprintf("%dn", c)
		if c == 1 {
			bn = "n"
		}
		question := fmt.Sprintf("Kummer's test uses \\(b_n=%s\\) and \\(a_n=1/n^{%d}\\). Compute \\(L=\\lim_{n\\to\\infty}(b_n a_n/a_{n+1}-b_{n+1})\\).", bn, p)
		explanation := fmt.Sprintf("\\(b_n a_n/a_{n+1}=%s(1+1/n)^{%d}=%s+%d+O(1/n)\\), while \\(b_{n+1}=%s+%d\\). The difference tends to \\(%d(%d-1)=%d\\).", bn, p, bn, c*p, bn, c, c, p, limit)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "c": strconv.Itoa(c), "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	case 2:
		p := rand.Intn(5) + 2
		c := rand.Intn(3) + 2
		limit := c * (p - 1)
		question := fmt.Sprintf("Kummer's test uses \\(b_n=%dn\\ln n\\) and \\(a_n=1/(n(\\ln n)^{%d})\\). Compute \\(L=\\lim_{n\\to\\infty}(b_n a_n/a_{n+1}-b_{n+1})\\).", c, p)
		explanation := fmt.Sprintf("\\(a_n/a_{n+1}=(1+1/n)(\\ln(n+1)/\\ln n)^{%d}=1+1/n+%d/(n\\ln n)+O(1/n^2)\\). Multiplying by \\(%dn\\ln n\\) and subtracting \\(%d(n+1)\\ln(n+1)\\) leaves \\(%d(%d-1)=%d\\).", p, p, c, c, c, p, limit)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "c": strconv.Itoa(c), "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	default:
		p := rand.Intn(5) + 2
		c := rand.Intn(3) + 2
		limit := c * p
		question := fmt.Sprintf("Kummer's test uses \\(b_n=%dn\\) and \\(a_n=1/(n^{%d}(n+1))\\). Compute \\(L=\\lim_{n\\to\\infty}(b_n a_n/a_{n+1}-b_{n+1})\\).", c, p)
		explanation := fmt.Sprintf("\\(a_n/a_{n+1}=\\frac{n+2}{n}(1+1/n)^{%d-1}=1+%d/n+O(1/n^2)\\). Multiplying by \\(%dn\\) and subtracting \\(%d(n+1)\\) leaves \\(%d\\cdot%d=%d\\).", p, p+1, c, c, c, p, limit)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(limit),
			Explanation: explanation,
			Facts:       map[string]string{"p": strconv.Itoa(p), "c": strconv.Itoa(c), "limit": strconv.Itoa(limit)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Gauss's test
// ---------------------------------------------------------------------------

type w3GaussGen struct{}

func (g *w3GaussGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		p := rand.Intn(30) + 1
		question := fmt.Sprintf("Gauss's test writes \\(a_n/a_{n+1}=1+h/n+O(1/n^2)\\). For \\(a_n=1/n^{%d}\\), what is \\(h\\)?", p)
		explanation := fmt.Sprintf("\\(a_n/a_{n+1}=(1+1/n)^{%d}=1+%d/n+O(1/n^2)\\), so \\(h=%d\\).", p, p, p)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(p),
			Explanation: explanation,
			Facts:       map[string]string{"h": strconv.Itoa(p), "p": strconv.Itoa(p), "form": "power"},
			KPOrdinal:   kp,
		}
	case 2:
		alpha := rand.Intn(5) + 1
		h := rand.Intn(5) + 1
		beta := alpha + h
		question := fmt.Sprintf("For \\(a_n=\\Gamma(n+%d)/\\Gamma(n+%d)\\), Gauss's test writes \\(a_n/a_{n+1}=1+h/n+O(1/n^2)\\). What is \\(h\\)?", alpha, beta)
		explanation := fmt.Sprintf("\\(a_n/a_{n+1}=(n+%d)/(n+%d)=1+%d/(n+%d)=1+%d/n+O(1/n^2)\\), so \\(h=%d\\).", beta, alpha, beta-alpha, alpha, h, h)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(h),
			Explanation: explanation,
			Facts:       map[string]string{"alpha": strconv.Itoa(alpha), "beta": strconv.Itoa(beta), "h": strconv.Itoa(h), "form": "gamma_ratio"},
			KPOrdinal:   kp,
		}
	default:
		d := rand.Intn(5) + 2
		alpha := rand.Intn(5) + 1
		h := d - 1
		question := fmt.Sprintf("For \\(a_n=(n+%d)/(n+1)^{%d}\\), Gauss's test writes \\(a_n/a_{n+1}=1+h/n+O(1/n^2)\\). What is \\(h\\)?", alpha, d)
		explanation := fmt.Sprintf("The numerator contributes \\(-1/n\\) to the logarithm of the ratio and the denominator contributes \\(%d/n\\), so the coefficient of \\(1/n\\) is \\(%d-1=%d\\).", d, d, h)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(h),
			Explanation: explanation,
			Facts:       map[string]string{"d": strconv.Itoa(d), "alpha": strconv.Itoa(alpha), "h": strconv.Itoa(h), "form": "rational"},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Abel's test
// ---------------------------------------------------------------------------

type w3AbelGen struct{}

func (g *w3AbelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		factors := []string{
			"1/n", "1/\\sqrt{n}", "1/(n+1)", "e^{-n}", "1/(n^2+1)", "2^{-n}", "1/(n+3)",
		}
		b := factors[rand.Intn(len(factors))]
		question := fmt.Sprintf("Suppose \\(\\sum a_n\\) converges and \\(b_n=%s\\), which is bounded and monotone. What can you conclude about \\(\\sum a_n b_n\\)?", b)
		explanation := "Abel's test applies: a convergent series multiplied term by term by a bounded monotone factor still converges."
		return w3Choice(kp, question, "converges", "diverges", explanation, map[string]string{
			"a_converges": "1", "b_bounded": "1", "b_monotone": "1", "b_form": b,
		})
	case 2:
		factors := []string{
			"1+1/n", "3-1/n", "2+1/n^2", "1+1/\\sqrt{n}", "4-1/n^2", "5-2/n",
			"6-1/n^3", "2-1/(n+1)", "7-3/n",
		}
		b := factors[rand.Intn(len(factors))]
		question := fmt.Sprintf("Suppose \\(\\sum a_n\\) converges and \\(b_n=%s\\), which is bounded and monotone. What can you conclude about \\(\\sum a_n b_n\\)?", b)
		explanation := "A bounded monotone factor need not tend to zero. Abel's test still applies, because only boundedness and monotonicity are required, so the product series converges."
		return w3Choice(kp, question, "converges", "diverges", explanation, map[string]string{
			"a_converges": "1", "b_bounded": "1", "b_monotone": "1", "b_form": b,
		})
	default:
		type factor struct {
			text     string
			bounded  bool
			monotone bool
		}
		factors := []factor{
			{"1/n", true, true},
			{"n", false, true},
			{"(-1)^n", true, false},
			{"\\sin n", true, false},
			{"1+1/n", true, true},
			{"\\arctan n", true, true},
			{"(-1)^n/n", true, false},
			{"n/(n+1)", true, true},
			{"1/n^2", true, true},
			{"(-1)^n n", false, false},
		}
		f := factors[rand.Intn(len(factors))]
		correct := "no"
		if f.bounded && f.monotone {
			correct = "yes"
		}
		question := fmt.Sprintf("Is \\(b_n=%s\\) bounded and monotone for \\(n\\ge1\\)?", f.text)
		explanation := "Abel's test needs both properties. A factor that is unbounded, or one that changes direction, does not satisfy the hypothesis."
		if correct == "yes" {
			explanation = "The sequence stays inside a finite range and moves in one direction, so it is bounded and monotone and Abel's test applies."
		}
		bounded := "0"
		if f.bounded {
			bounded = "1"
		}
		monotone := "0"
		if f.monotone {
			monotone = "1"
		}
		return w3Choice(kp, question, correct, map[string]string{"yes": "no", "no": "yes"}[correct], explanation, map[string]string{
			"b_form": f.text, "bounded": bounded, "monotone": monotone,
		})
	}
}

// ---------------------------------------------------------------------------
// Uniform convergence
// ---------------------------------------------------------------------------

type w3UniformConvergenceGen struct{}

func (g *w3UniformConvergenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		type seqCase struct {
			fn       string
			interval string
			sup      string
		}
		cases := []seqCase{
			{"x^n", "[0,1]", "1"},
			{"x^n", "[0,1/2]", "0"},
			{"x^n", "[0,1/3]", "0"},
			{"x^n", "[0,3/4]", "0"},
			{"x/n", "[0,1]", "0"},
			{"x/n", "[0,\\infty)", "inf"},
			{"\\sin(x)/n", "\\mathbb{R}", "0"},
			{"\\cos(x)/n", "[0,1]", "0"},
			{"\\sin(nx)/n", "\\mathbb{R}", "0"},
			{"x/(1+nx)", "[0,1]", "0"},
			{"1/(1+nx)", "(0,1]", "1"},
		}
		c := cases[rand.Intn(len(cases))]
		correct := "no"
		if c.sup == "0" {
			correct = "yes"
		}
		question := fmt.Sprintf("Does \\(f_n(x)=%s\\) converge uniformly on \\(%s\\)? Enter yes or no.", c.fn, c.interval)
		explanation := fmt.Sprintf("The supremum of \\(|f_n(x)-f(x)|\\) on the interval behaves as %s. Uniform convergence requires that supremum to tend to 0, so the answer is %s.", c.sup, correct)
		return w3Choice(kp, question, correct, map[string]string{"yes": "no", "no": "yes"}[correct], explanation, map[string]string{
			"sup_limit": c.sup, "fn": c.fn, "interval": c.interval,
		})
	case 2:
		type seriesCase struct {
			id       string
			series   string
			interval string
			uniform  string
		}
		cases := []seriesCase{
			{"power_half", "\\sum_{n=0}^{\\infty} x^n", "[0,1/2]", "1"},
			{"power_third", "\\sum_{n=0}^{\\infty} x^n", "[0,1/3]", "1"},
			{"power_unit", "\\sum_{n=0}^{\\infty} x^n", "[0,1)", "0"},
			{"p2", "\\sum_{n=1}^{\\infty} x^n/n^2", "[0,1]", "1"},
			{"harmonic_open", "\\sum_{n=1}^{\\infty} x^n/n", "[0,1)", "0"},
			{"alt_harmonic", "\\sum_{n=1}^{\\infty} (-1)^n x^n/n", "[0,1]", "1"},
			{"p3", "\\sum_{n=1}^{\\infty} x^n/n^3", "[0,1]", "1"},
			{"p4", "\\sum_{n=1}^{\\infty} x^n/n^4", "[0,1]", "1"},
		}
		c := cases[rand.Intn(len(cases))]
		correct := "no"
		if c.uniform == "1" {
			correct = "yes"
		}
		question := fmt.Sprintf("Does \\(%s\\) converge uniformly on \\(%s\\)? Enter yes or no.", c.series, c.interval)
		explanation := ""
		switch c.id {
		case "power_half", "power_third":
			explanation = "On this interval \\(|x|\\) is bounded below 1, so the geometric remainder is bounded by a constant times the corresponding power of the endpoint. The bound is independent of \\(x\\) and tends to 0."
		case "power_unit":
			explanation = "On \\([0,1)\\) the geometric sum is \\(1/(1-x)\\), which is unbounded as \\(x\\) approaches 1. Uniform convergence would force a bounded sum, so the convergence is not uniform."
		case "harmonic_open":
			explanation = "The sum behaves like \\(-\\ln(1-x)\\), which is unbounded as \\(x\\) approaches 1. Uniform convergence would keep the sum bounded, so the convergence is not uniform."
		case "alt_harmonic":
			explanation = "For each \\(x\\in[0,1]\\) the alternating series has decreasing terms. Its remainder is at most the first omitted term, \\(x^{n+1}/(n+1)\\le 1/(n+1)\\), a bound independent of \\(x\\) that tends to 0."
		default:
			explanation = "For every \\(x\\in[0,1]\\) the terms satisfy \\(|x^n/n^p|\\le 1/n^p\\), and \\(\\sum 1/n^p\\) converges. The Weierstrass M-test therefore applies."
		}
		explanation += fmt.Sprintf(" The answer is %s.", correct)
		return w3Choice(kp, question, correct, map[string]string{"yes": "no", "no": "yes"}[correct], explanation, map[string]string{
			"kind": c.id, "interval": c.interval,
		})
	default:
		bounds := []string{"1/n^2", "1/2^n", "1/n^3", "1/(n^2+1)", "1/n!", "1/3^n", "1/n^4", "1/4^n", "1/(n^2+n)"}
		m := bounds[rand.Intn(len(bounds))]
		question := fmt.Sprintf("A series \\(\\sum f_n\\) on a set satisfies \\(|f_n(x)|\\le M_n=%s\\) for every \\(n\\) and \\(x\\). Does the series converge uniformly on the set? Enter yes or no.", m)
		explanation := "The bounding series is convergent, so the Weierstrass M-test applies and the series converges uniformly."
		return w3Choice(kp, question, "yes", "no", explanation, map[string]string{
			"m_series": m, "m_converges": "1",
		})
	}
}

// ---------------------------------------------------------------------------
// Euler's transformation of series
// ---------------------------------------------------------------------------

type w3EulerTransformGen struct{}

func (g *w3EulerTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		type ratio struct{ num, den int }
		ratios := []ratio{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {3, 4}}
		r := ratios[rand.Intn(len(ratios))]
		k := rand.Intn(3)
		coef := w3MakeFrac(mathutil.IntPow(r.den-r.num, k), mathutil.IntPow(r.den, k)*mathutil.IntPow(2, k+1))
		question := fmt.Sprintf("Euler's transformation writes \\(\\sum_{n\\ge0}(-1)^n a_n=\\sum_{k\\ge0}c_k\\). For \\(a_n=(%s)^n\\), what is \\(c_{%d}\\)?", w3FracTex(r.num, r.den), k)
		explanation := fmt.Sprintf("For a geometric term, \\(\\Delta^k a_0=(1-r)^k\\), so \\(c_k=(1-r)^k/2^{k+1}\\). With \\(r=%s\\) and \\(k=%d\\), \\(c_k=%s\\).", w3FracTex(r.num, r.den), k, coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"r_num": strconv.Itoa(r.num), "r_den": strconv.Itoa(r.den), "k": strconv.Itoa(k)},
			KPOrdinal:   kp,
		}
	case 2:
		type ratio struct{ num, den int }
		ratios := []ratio{{1, 2}, {1, 3}, {2, 3}, {1, 4}, {3, 4}, {1, 5}}
		r := ratios[rand.Intn(len(ratios))]
		sum := w3MakeFrac(r.den, r.den+r.num)
		question := fmt.Sprintf("For \\(a_n=(%s)^n\\), Euler's transformation sums to \\(1/(1+r)\\). What is the sum?", w3FracTex(r.num, r.den))
		explanation := fmt.Sprintf("The transformed geometric series has ratio \\((1-r)/2=%s\\), so its sum is \\(1/(1-(1-r)/2)=1/(1+r)=%s\\).", w3FracTex(r.den-r.num, 2*r.den), sum.String())
		return generator.Problem{
			Question:    question,
			Answer:      sum.String(),
			Explanation: explanation,
			Facts:       map[string]string{"r_num": strconv.Itoa(r.num), "r_den": strconv.Itoa(r.den)},
			KPOrdinal:   kp,
		}
	default:
		k := rand.Intn(30)
		coef := w3MakeFrac(1, (k+1)*mathutil.IntPow(2, k+1))
		question := fmt.Sprintf("For the alternating harmonic series \\(\\sum_{n\\ge0}(-1)^n/(n+1)\\), Euler's transformation gives \\(c_k=1/((k+1)2^{k+1})\\). What is \\(c_{%d}\\)?", k)
		explanation := fmt.Sprintf("The forward differences of \\(a_n=1/(n+1)\\) are \\(\\Delta^k a_0=1/(k+1)\\). Divide by \\(2^{k+1}\\): \\(c_k=1/((k+1)2^{k+1})=%s\\).", coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"k": strconv.Itoa(k)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Integrating power series
// ---------------------------------------------------------------------------

type w3PowerIntegrationGen struct{}

func (g *w3PowerIntegrationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		n := rand.Intn(5) + 1
		k := rand.Intn(4) + 1
		c := k * (n + 1)
		question := fmt.Sprintf("A power series \\(\\sum a_m x^m\\) has \\(a_{%d}=%d\\). Term-by-term integration gives the coefficient of \\(x^{%d}\\) equal to \\(k\\). What is \\(k\\)?", n, c, n+1)
		explanation := fmt.Sprintf("\\(\\int a_n x^n\\,dx=a_n x^{n+1}/(n+1)\\). With \\(a_{%d}=%d\\) the coefficient is \\(%d/%d=%d\\), so \\(k=%d\\).", n, c, c, n+1, k, k)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(k),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "c": strconv.Itoa(c)},
			KPOrdinal:   kp,
		}
	case 2:
		m := rand.Intn(30) + 2
		coef := w3MakeFrac(1, m)
		question := fmt.Sprintf("Integrating \\(\\sum_{n\\ge0}x^n=1/(1-x)\\) term by term gives \\(-\\ln(1-x)=\\sum_{m\\ge1}x^m/m\\). What is the coefficient of \\(x^{%d}\\)?", m)
		explanation := fmt.Sprintf("Integrating \\(x^{m-1}\\) gives \\(x^m/m\\), so the coefficient of \\(x^{%d}\\) is \\(1/%d=%s\\).", m, m, coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"m": strconv.Itoa(m)},
			KPOrdinal:   kp,
		}
	default:
		k := rand.Intn(30)
		coef := w3MakeFrac(1, 2*k+1)
		question := fmt.Sprintf("Integrating \\(\\sum_{n\\ge0}x^{2n}=1/(1-x^2)\\) term by term gives \\(\\sum_{n\\ge0}x^{2n+1}/(2n+1)\\). What is the coefficient of \\(x^{%d}\\)?", 2*k+1)
		explanation := fmt.Sprintf("The term with \\(2n=%d\\) is \\(x^{%d}/%d\\), so the coefficient is \\(1/%d=%s\\).", 2*k, 2*k+1, 2*k+1, 2*k+1, coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"k": strconv.Itoa(k)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Differentiating power series
// ---------------------------------------------------------------------------

type w3PowerDifferentiationGen struct{}

func (g *w3PowerDifferentiationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		n := rand.Intn(6) + 1
		c := rand.Intn(5) + 1
		coef := n * c
		question := fmt.Sprintf("A power series \\(\\sum a_m x^m\\) has \\(a_{%d}=%d\\). Term-by-term differentiation gives the coefficient of \\(x^{%d}\\) equal to \\(k\\). What is \\(k\\)?", n, c, n-1)
		explanation := fmt.Sprintf("Differentiating \\(a_n x^n\\) gives \\(n a_n x^{n-1}\\), so the coefficient is \\(%d\\cdot%d=%d\\).", n, c, coef)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(coef),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "c": strconv.Itoa(c)},
			KPOrdinal:   kp,
		}
	case 2:
		k := rand.Intn(30)
		coef := k + 1
		question := fmt.Sprintf("Differentiating \\(\\sum_{n\\ge0}x^n=1/(1-x)\\) gives \\(\\sum_{n\\ge1}n x^{n-1}=1/(1-x)^2\\). What is the coefficient of \\(x^{%d}\\)?", k)
		explanation := fmt.Sprintf("The coefficient of \\(x^{%d}\\) comes from \\(n=%d\\), so it is %d.", k, k+1, coef)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(coef),
			Explanation: explanation,
			Facts:       map[string]string{"k": strconv.Itoa(k)},
			KPOrdinal:   kp,
		}
	default:
		k := rand.Intn(30) + 1
		coef := 2 * k
		question := fmt.Sprintf("Differentiating \\(\\sum_{n\\ge0}x^{2n}=1/(1-x^2)\\) gives \\(\\sum_{n\\ge1}2n x^{2n-1}=2x/(1-x^2)^2\\). What is the coefficient of \\(x^{%d}\\)?", 2*k-1)
		explanation := fmt.Sprintf("The coefficient of \\(x^{2n-1}\\) is \\(2n\\). With \\(2n-1=%d\\) we have \\(n=%d\\), so the coefficient is %d.", 2*k-1, k, coef)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(coef),
			Explanation: explanation,
			Facts:       map[string]string{"k": strconv.Itoa(k)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Lagrange inversion
// ---------------------------------------------------------------------------

type w3LagrangeInversionGen struct{}

func (g *w3LagrangeInversionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		n := rand.Intn(8) + 1
		p := rand.Intn(4) + 1
		coef := w3Binom(p*n, n-1) / n
		question := fmt.Sprintf("Lagrange inversion gives \\([z^n]w(z)=\\frac1n[t^{n-1}]\\phi(t)^n\\). If \\(w=z(1+w)^{%d}\\), what is the coefficient of \\(z^{%d}\\)?", p, n)
		explanation := fmt.Sprintf("\\(\\phi(t)^n=(1+t)^{%d}\\), whose \\(t^{%d}\\) coefficient is \\(\\binom{%d}{%d}=%d\\). Divide by \\(n=%d\\): the coefficient is %d.", p*n, n-1, p*n, n-1, w3Binom(p*n, n-1), n, coef)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(coef),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "p": strconv.Itoa(p)},
			KPOrdinal:   kp,
		}
	case 2:
		n := rand.Intn(15) + 1
		coef := w3MakeFrac(mathutil.IntPow(n, n-1), w3Factorial(n))
		question := fmt.Sprintf("For \\(w=z e^w\\), Lagrange inversion gives \\([z^n]w=n^{n-1}/n!\\). What is the coefficient of \\(z^{%d}\\)?", n)
		explanation := fmt.Sprintf("\\(n^{n-1}/n!=%d/%d=%s\\).", mathutil.IntPow(n, n-1), w3Factorial(n), coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n)},
			KPOrdinal:   kp,
		}
	default:
		n := rand.Intn(16)
		coef := w3Binom(2*n, n) / (n + 1)
		question := fmt.Sprintf("For \\(w=z/(1-w)\\), the coefficient \\([z^n]w\\) is the \\(n\\)-th Catalan number. Compute it for \\(n=%d\\).", n)
		explanation := fmt.Sprintf("The Catalan number is \\(\\binom{2n}{n}/(n+1)=\\binom{%d}{%d}/%d=%d/%d=%d\\).", 2*n, n, n+1, w3Binom(2*n, n), n+1, coef)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(coef),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Romberg integration
// ---------------------------------------------------------------------------

type w3RombergGen struct{}

func (g *w3RombergGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		t1 := rand.Intn(20) + 1
		t2 := rand.Intn(20) + 1
		if t1 == t2 {
			t2 = t1%20 + 1
		}
		r := w3MakeFrac(4*t2-t1, 3)
		question := fmt.Sprintf("A Romberg table has \\(T(h)=%d\\) and \\(T(h/2)=%d\\). The first extrapolated value is \\(R_{1,1}=T(h/2)+[T(h/2)-T(h)]/3\\). What is \\(R_{1,1}\\)?", t1, t2)
		explanation := fmt.Sprintf("\\(R_{1,1}=%d+(%d-%d)/3=(4\\cdot%d-%d)/3=%s\\).", t2, t2, t1, t2, t1, r.String())
		return generator.Problem{
			Question:    question,
			Answer:      r.String(),
			Explanation: explanation,
			Facts:       map[string]string{"t_h": strconv.Itoa(t1), "t_half": strconv.Itoa(t2)},
			KPOrdinal:   kp,
		}
	case 2:
		r1 := rand.Intn(20) + 1
		r2 := rand.Intn(20) + 1
		if r1 == r2 {
			r2 = r1%20 + 1
		}
		r := w3MakeFrac(16*r2-r1, 15)
		question := fmt.Sprintf("A Romberg table has \\(R_{k-1,1}=%d\\) and \\(R_{k,1}=%d\\). The next column is \\(R_{k,2}=R_{k,1}+(R_{k,1}-R_{k-1,1})/15\\). What is \\(R_{k,2}\\)?", r1, r2)
		explanation := fmt.Sprintf("\\(R_{k,2}=%d+(%d-%d)/15=(16\\cdot%d-%d)/15=%s\\).", r2, r2, r1, r2, r1, r.String())
		return generator.Problem{
			Question:    question,
			Answer:      r.String(),
			Explanation: explanation,
			Facts:       map[string]string{"r_prev": strconv.Itoa(r1), "r_cur": strconv.Itoa(r2)},
			KPOrdinal:   kp,
		}
	default:
		t0 := rand.Intn(12) + 1
		t1 := rand.Intn(12) + 1
		t2 := rand.Intn(12) + 1
		r11 := w3MakeFrac(4*t1-t0, 3)
		r21 := w3MakeFrac(4*t2-t1, 3)
		r22 := r21.add(r21.sub(r11).div(w3MakeFrac(15, 1)))
		question := fmt.Sprintf("A Romberg table starts with \\(T(h)=%d\\), \\(T(h/2)=%d\\) and \\(T(h/4)=%d\\). Compute \\(R_{2,2}\\) in the order \\(R_{1,1}\\), \\(R_{2,1}\\), then \\(R_{2,2}=R_{2,1}+(R_{2,1}-R_{1,1})/15\\).", t0, t1, t2)
		explanation := fmt.Sprintf("\\(R_{1,1}=(4\\cdot%d-%d)/3=%s\\) and \\(R_{2,1}=(4\\cdot%d-%d)/3=%s\\). Therefore \\(R_{2,2}=%s+(%s-(%s))/15=%s\\).", t1, t0, r11.String(), t2, t1, r21.String(), r21.String(), r21.String(), r11.String(), r22.String())
		return generator.Problem{
			Question:    question,
			Answer:      r22.String(),
			Explanation: explanation,
			Facts:       map[string]string{"t_h": strconv.Itoa(t0), "t_half": strconv.Itoa(t1), "t_quarter": strconv.Itoa(t2)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Monte Carlo integration
// ---------------------------------------------------------------------------

type w3MonteCarloGen struct{}

func (g *w3MonteCarloGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		n := rand.Intn(3) + 3
		a := rand.Intn(3)
		b := a + rand.Intn(4) + 1
		values := make([]int, n)
		sum := 0
		parts := make([]string, n)
		for i := range values {
			values[i] = rand.Intn(9) + 1
			sum += values[i]
			parts[i] = strconv.Itoa(values[i])
		}
		estimate := w3MakeFrac((b-a)*sum, n)
		question := fmt.Sprintf("A Monte Carlo estimate of \\(\\int_{%d}^{%d} f(x)\\,dx\\) uses %d uniform samples with values %s. The estimate is \\(\\frac{b-a}{N}\\sum f(x_i)\\). Compute it.", a, b, n, strings.Join(parts, ", "))
		explanation := fmt.Sprintf("The sample sum is %d and the interval has length %d, so the estimate is \\((b-a)\\cdot\\text{sum}/N=%d\\cdot%d/%d=%s\\).", sum, b-a, b-a, sum, n, estimate.String())
		return generator.Problem{
			Question:    question,
			Answer:      estimate.String(),
			Explanation: explanation,
			Facts: map[string]string{
				"a": strconv.Itoa(a), "b": strconv.Itoa(b), "n": strconv.Itoa(n), "samples": strings.Join(parts, ","),
			},
			KPOrdinal: kp,
		}
	case 2:
		type pair struct{ n, m int }
		pairs := []pair{
			{1000, 785}, {1000, 786}, {500, 392}, {500, 393}, {400, 314},
			{400, 315}, {250, 196}, {250, 197}, {200, 157}, {200, 158},
		}
		p := pairs[rand.Intn(len(pairs))]
		estimate := w3Fixed(w3Round(4*float64(p.m)/float64(p.n), 3), 3)
		question := fmt.Sprintf("In a Monte Carlo estimate of \\(\\pi\\), %d of %d random points fall inside a quarter circle. The estimate is \\(4m/N\\). Compute it to three decimal places.", p.m, p.n)
		explanation := fmt.Sprintf("\\(4m/N=4\\cdot%d/%d=%s\\), rounded to three decimal places. This is a sampled approximation, not an exact value of \\(\\pi\\).", p.m, p.n, estimate)
		return generator.Problem{
			Question:    question,
			Answer:      estimate,
			Explanation: explanation,
			Facts:       map[string]string{"m": strconv.Itoa(p.m), "n": strconv.Itoa(p.n)},
			KPOrdinal:   kp,
		}
	default:
		s := rand.Intn(5) + 1
		sizes := []int{10, 20, 50, 100, 200}
		n := sizes[rand.Intn(len(sizes))]
		se := w3Fixed(w3Round(float64(s)/math.Sqrt(float64(n)), 3), 3)
		question := fmt.Sprintf("A Monte Carlo average has sample standard deviation \\(s=%d\\) and \\(N=%d\\) samples. Its standard error is \\(s/\\sqrt{N}\\). Compute it to three decimal places.", s, n)
		explanation := fmt.Sprintf("\\(s/\\sqrt{N}=%d/\\sqrt{%d}=%s\\), rounded to three decimal places. The square root is irrational, so the answer is an approximation.", s, n, se)
		return generator.Problem{
			Question:    question,
			Answer:      se,
			Explanation: explanation,
			Facts:       map[string]string{"s": strconv.Itoa(s), "n": strconv.Itoa(n)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Gauss-Legendre quadrature
// ---------------------------------------------------------------------------

type w3GaussLegendreGen struct{}

func (g *w3GaussLegendreGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		a := rand.Intn(30) + 1
		value := w3MakeFrac(2*a*a*a, 3)
		question := fmt.Sprintf("Apply the two-point Gauss-Legendre rule, with nodes \\(\\pm1/\\sqrt3\\), to \\(f(x)=x^2\\) on \\([-%d,%d]\\). What value does the rule return?", a, a)
		explanation := fmt.Sprintf("On \\([-%d,%d]\\) the nodes are \\(\\pm %d/\\sqrt3\\) and each weight is %d. Each node contributes \\(%d\\cdot %d^2/3=%s\\), so the rule returns \\(2\\cdot%s=%s\\).", a, a, a, a, a, a, w3MakeFrac(a*a*a, 3).String(), w3MakeFrac(a*a*a, 3).String(), value.String())
		return generator.Problem{
			Question:    question,
			Answer:      value.String(),
			Explanation: explanation,
			Facts:       map[string]string{"a": strconv.Itoa(a), "points": "2"},
			KPOrdinal:   kp,
		}
	case 2:
		a := rand.Intn(30) + 1
		value := w3MakeFrac(2*mathutil.IntPow(a, 5), 5)
		question := fmt.Sprintf("Apply the three-point Gauss-Legendre rule, with nodes \\(-\\sqrt{3/5},0,\\sqrt{3/5}\\) and weights \\(5/9,8/9,5/9\\), to \\(f(x)=x^4\\) on \\([-%d,%d]\\). What value does the rule return?", a, a)
		explanation := fmt.Sprintf("On \\([-%d,%d]\\) each outer node contributes \\(%d\\cdot(5/9)\\cdot(%d^4\\cdot9/25)=%s\\), and the center term is zero. The rule returns \\(2\\cdot%s=%s\\).", a, a, a, a, w3MakeFrac(mathutil.IntPow(a, 5), 5).String(), w3MakeFrac(mathutil.IntPow(a, 5), 5).String(), value.String())
		return generator.Problem{
			Question:    question,
			Answer:      value.String(),
			Explanation: explanation,
			Facts:       map[string]string{"a": strconv.Itoa(a), "points": "3"},
			KPOrdinal:   kp,
		}
	default:
		a := rand.Intn(4) + 1
		b := a + rand.Intn(5) + 1
		value := w3MakeFrac((b-a)*(a+b), 2)
		question := fmt.Sprintf("Apply the two-point Gauss-Legendre rule to \\(f(x)=x\\) on \\([%d,%d]\\). What value does the rule return?", a, b)
		explanation := fmt.Sprintf("The midpoint is %d and each weight is \\(%d/2\\). The two nodes are symmetric, so the rule returns \\((b-a)(a+b)/2=%d\\cdot%d/2=%s\\).", (a+b)/2, b-a, b-a, a+b, value.String())
		return generator.Problem{
			Question:    question,
			Answer:      value.String(),
			Explanation: explanation,
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "points": "2"},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Fresnel integrals
// ---------------------------------------------------------------------------

type w3FresnelGen struct{}

func (g *w3FresnelGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	xs := []float64{0.2, 0.4, 0.5, 0.8, 1.0, 1.2, 1.5, 1.8, 2.0, 2.5}
	switch kp {
	case 1:
		x := xs[rand.Intn(len(xs))]
		value := w3Fixed(w3Round(w3FresnelSeries(x, false), 4), 4)
		question := fmt.Sprintf("The Fresnel cosine integral is \\(C(x)=\\int_0^x\\cos(\\pi t^2/2)\\,dt\\). Compute \\(C(%s)\\) to four decimal places.", w3Fixed(x, 1))
		explanation := fmt.Sprintf("Summing the Maclaurin series \\(C(x)=\\sum_{n\\ge0}(-1)^n(\\pi/2)^{2n}x^{4n+1}/((2n)!(4n+1))\\) through thirty terms at \\(x=%s\\) gives %s. The answer is a numerical approximation.", w3Fixed(x, 1), value)
		return generator.Problem{
			Question:    question,
			Answer:      value,
			Explanation: explanation,
			Facts:       map[string]string{"x": w3Fixed(x, 1), "function": "C", "terms": "30"},
			KPOrdinal:   kp,
		}
	case 2:
		x := xs[rand.Intn(len(xs))]
		value := w3Fixed(w3Round(w3FresnelSeries(x, true), 4), 4)
		question := fmt.Sprintf("The Fresnel sine integral is \\(S(x)=\\int_0^x\\sin(\\pi t^2/2)\\,dt\\). Compute \\(S(%s)\\) to four decimal places.", w3Fixed(x, 1))
		explanation := fmt.Sprintf("Summing the Maclaurin series \\(S(x)=\\sum_{n\\ge0}(-1)^n(\\pi/2)^{2n+1}x^{4n+3}/((2n+1)!(4n+3))\\) through thirty terms at \\(x=%s\\) gives %s. The answer is a numerical approximation.", w3Fixed(x, 1), value)
		return generator.Problem{
			Question:    question,
			Answer:      value,
			Explanation: explanation,
			Facts:       map[string]string{"x": w3Fixed(x, 1), "function": "S", "terms": "30"},
			KPOrdinal:   kp,
		}
	default:
		n := rand.Intn(6)
		value := math.Pow(-1, float64(n)) * math.Pow(math.Pi/2, float64(2*n)) /
			(w3FloatFactorial(2*n) * float64(4*n+1))
		out := w3Fixed(w3Round(value, 6), 6)
		question := fmt.Sprintf("The Maclaurin series of \\(C(x)\\) is \\(\\sum_{n\\ge0}(-1)^n(\\pi/2)^{2n}x^{4n+1}/((2n)!(4n+1))\\). What is the coefficient of \\(x^{%d}\\) to six decimal places?", 4*n+1)
		explanation := fmt.Sprintf("Set \\(n=%d\\) in \\((-1)^n(\\pi/2)^{2n}/((2n)!(4n+1))\\): the value is %s. It is written to six decimal places because it involves powers of \\(\\pi\\).", n, out)
		return generator.Problem{
			Question:    question,
			Answer:      out,
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "terms": "30"},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Elliptic integrals
// ---------------------------------------------------------------------------

type w3EllipticGen struct{}

func (g *w3EllipticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	ks := []float64{
		0.10, 0.15, 0.20, 0.25, 0.30, 0.35, 0.40, 0.45, 0.50,
		0.55, 0.60, 0.65, 0.70, 0.75, 0.80, 0.85, 0.90,
	}
	switch kp {
	case 1:
		k := ks[rand.Intn(len(ks))]
		value := w3Fixed(w3Round(w3EllipticK(k), 4), 4)
		question := fmt.Sprintf("The complete elliptic integral of the first kind is \\(K(k)=\\int_0^{\\pi/2}d\\theta/\\sqrt{1-k^2\\sin^2\\theta}\\). Compute \\(K(%s)\\) to four decimal places.", w3Fixed(k, 2))
		explanation := fmt.Sprintf("Summing \\(K(k)=\\frac{\\pi}{2}\\sum_{n\\ge0}\\left(\\frac{(2n)!}{2^{2n}(n!)^2}\\right)^2k^{2n}\\) through thirty terms at \\(k=%s\\) gives %s. The answer is a numerical approximation.", w3Fixed(k, 2), value)
		return generator.Problem{
			Question:    question,
			Answer:      value,
			Explanation: explanation,
			Facts:       map[string]string{"k": w3Fixed(k, 2), "kind": "K", "terms": "30"},
			KPOrdinal:   kp,
		}
	case 2:
		k := ks[rand.Intn(len(ks))]
		value := w3Fixed(w3Round(w3EllipticE(k), 4), 4)
		question := fmt.Sprintf("The complete elliptic integral of the second kind is \\(E(k)=\\int_0^{\\pi/2}\\sqrt{1-k^2\\sin^2\\theta}\\,d\\theta\\). Compute \\(E(%s)\\) to four decimal places.", w3Fixed(k, 2))
		explanation := fmt.Sprintf("Summing \\(E(k)=\\frac{\\pi}{2}\\left[1-\\sum_{n\\ge1}\\left(\\frac{(2n)!}{2^{2n}(n!)^2}\\right)^2\\frac{k^{2n}}{2n-1}\\right]\\) through thirty terms at \\(k=%s\\) gives %s. The answer is a numerical approximation.", w3Fixed(k, 2), value)
		return generator.Problem{
			Question:    question,
			Answer:      value,
			Explanation: explanation,
			Facts:       map[string]string{"k": w3Fixed(k, 2), "kind": "E", "terms": "30"},
			KPOrdinal:   kp,
		}
	default:
		k := ks[rand.Intn(len(ks))]
		agm := w3AGM(1, math.Sqrt(1-k*k))
		value := w3Fixed(w3Round(math.Pi/(2*agm), 4), 4)
		question := fmt.Sprintf("Using the arithmetic-geometric mean, \\(K(k)=\\pi/(2\\,\\mathrm{AGM}(1,\\sqrt{1-k^2}))\\). Compute \\(K(%s)\\) to four decimal places.", w3Fixed(k, 2))
		explanation := fmt.Sprintf("Iterating \\(a_{n+1}=(a_n+b_n)/2\\) and \\(b_{n+1}=\\sqrt{a_nb_n}\\) from \\((1,\\sqrt{1-k^2})\\) at \\(k=%s\\) gives \\(\\mathrm{AGM}=%s\\). Then \\(K=\\pi/(2\\cdot%s)=%s\\), a numerical approximation.", w3Fixed(k, 2), w3Fixed(agm, 10), w3Fixed(agm, 10), value)
		return generator.Problem{
			Question:    question,
			Answer:      value,
			Explanation: explanation,
			Facts:       map[string]string{"k": w3Fixed(k, 2), "kind": "AGM", "agm": w3Fixed(agm, 10)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Wallis integrals
// ---------------------------------------------------------------------------

type w3WallisGen struct{}

func (g *w3WallisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		m := rand.Intn(13)
		n := 2*m + 1
		value := w3MakeFrac(1, 1)
		for j := 1; j <= m; j++ {
			value = value.mul(w3MakeFrac(2*j, 2*j+1))
		}
		question := fmt.Sprintf("The Wallis integral is \\(W_n=\\int_0^{\\pi/2}\\sin^n x\\,dx\\). Compute \\(W_{%d}\\).", n)
		explanation := fmt.Sprintf("For odd \\(n=2m+1\\), \\(W_n=\\prod_{j=1}^{m}\\frac{2j}{2j+1}\\). With \\(m=%d\\) the product is %s.", m, value.String())
		return generator.Problem{
			Question:    question,
			Answer:      value.String(),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m), "parity": "odd"},
			KPOrdinal:   kp,
		}
	case 2:
		m := rand.Intn(13) + 1
		n := 2 * m
		coef := w3MakeFrac(1, 2)
		for j := 1; j <= m; j++ {
			coef = coef.mul(w3MakeFrac(2*j-1, 2*j))
		}
		question := fmt.Sprintf("For the Wallis integral \\(W_n=\\int_0^{\\pi/2}\\sin^n x\\,dx\\), \\(W_{%d}=c\\pi\\). What is \\(c\\)?", n)
		explanation := fmt.Sprintf("For even \\(n=2m\\), \\(W_n=\\frac{\\pi}{2}\\prod_{j=1}^{m}\\frac{2j-1}{2j}\\). With \\(m=%d\\) the coefficient is %s.", m, coef.String())
		return generator.Problem{
			Question:    question,
			Answer:      coef.String(),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m), "parity": "even"},
			KPOrdinal:   kp,
		}
	default:
		n := rand.Intn(29) + 2
		factor := w3MakeFrac(n-1, n)
		question := fmt.Sprintf("Wallis' recurrence is \\(W_n=\\frac{n-1}{n}W_{n-2}\\). For \\(n=%d\\), what is the factor multiplying \\(W_{n-2}\\)?", n)
		explanation := fmt.Sprintf("Substitute \\(n=%d\\) into \\((n-1)/n\\): the factor is \\(%d/%d=%s\\).", n, n-1, n, factor.String())
		return generator.Problem{
			Question:    question,
			Answer:      factor.String(),
			Explanation: explanation,
			Facts:       map[string]string{"n": strconv.Itoa(n)},
			KPOrdinal:   kp,
		}
	}
}

// ---------------------------------------------------------------------------
// Centroid of a curve
// ---------------------------------------------------------------------------

type w3CentroidCurveGen struct{}

func (g *w3CentroidCurveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := calcKP(ctx)
	switch kp {
	case 1:
		x1 := rand.Intn(9) + 1
		x2 := x1 + rand.Intn(9) + 1
		y1 := rand.Intn(12) + 1
		y2 := rand.Intn(12) + 1
		xbar := w3MakeFrac(x1+x2, 2)
		question := fmt.Sprintf("A straight curve runs from \\((%d,%d)\\) to \\((%d,%d)\\). Its centroid \\(x\\)-coordinate is the average of the endpoint \\(x\\)-coordinates. Compute \\(\\bar x\\).", x1, y1, x2, y2)
		explanation := fmt.Sprintf("The segment is a straight line, so its centroid is its midpoint. Thus \\(\\bar x=(%d+%d)/2=%s\\).", x1, x2, xbar.String())
		return generator.Problem{
			Question:    question,
			Answer:      xbar.String(),
			Explanation: explanation,
			Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2)},
			KPOrdinal:   kp,
		}
	case 2:
		x1 := rand.Intn(12) + 1
		x2 := rand.Intn(12) + 1
		y1 := rand.Intn(9) + 1
		y2 := y1 + rand.Intn(9) + 1
		ybar := w3MakeFrac(y1+y2, 2)
		question := fmt.Sprintf("A straight curve runs from \\((%d,%d)\\) to \\((%d,%d)\\). Its centroid \\(y\\)-coordinate is the average of the endpoint \\(y\\)-coordinates. Compute \\(\\bar y\\).", x1, y1, x2, y2)
		explanation := fmt.Sprintf("The centroid of a straight segment is its midpoint. Thus \\(\\bar y=(%d+%d)/2=%s\\).", y1, y2, ybar.String())
		return generator.Problem{
			Question:    question,
			Answer:      ybar.String(),
			Explanation: explanation,
			Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "x2": strconv.Itoa(x2), "y2": strconv.Itoa(y2)},
			KPOrdinal:   kp,
		}
	default:
		type triple struct{ a, b, c int }
		triples := []triple{{3, 4, 5}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}}
		t := triples[rand.Intn(len(triples))]
		scale := rand.Intn(3) + 1
		dx := scale * t.a
		dy := scale * t.b
		length := scale * t.c
		x1 := rand.Intn(9) + 1
		y1 := rand.Intn(9) + 1
		question := fmt.Sprintf("A straight curve runs from \\((%d,%d)\\) to \\((%d,%d)\\). Compute its arc length.", x1, y1, x1+dx, y1+dy)
		explanation := fmt.Sprintf("The displacement is \\((%d,%d)\\), whose length is \\(\\sqrt{%d^2+%d^2}=\\sqrt{%d}=%d\\). The centroid formula divides by this length.", dx, dy, dx, dy, dx*dx+dy*dy, length)
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(length),
			Explanation: explanation,
			Facts:       map[string]string{"x1": strconv.Itoa(x1), "y1": strconv.Itoa(y1), "dx": strconv.Itoa(dx), "dy": strconv.Itoa(dy)},
			KPOrdinal:   kp,
		}
	}
}
