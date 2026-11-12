package calculus

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Wave 22 adds twenty-three calculus II concepts: limits and subsequences,
// the zeta and eta series, Bernoulli numbers, double series, Stirling's
// approximation, Abel summation, power-series uniqueness, grouping,
// generating functions, arithmetico-geometric series, volumes by cross
// sections, parametric length and area, the beta-gamma relation, unbounded
// integrands, pumping work, the three-eighths rule, Riemann-Stieltjes
// integrals, Holder and Jensen inequalities for integrals, and mass from a
// radial density.

func registerWave22(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"calc.seq.limit":                    "limit",
		"calc.seq.subsequence":              "subsequence",
		"calc.series.riemann_zeta":          "zeta",
		"calc.series.eta":                   "eta",
		"calc.series.bernoulli":             "bernoulli",
		"calc.series.double":                "double",
		"calc.series.stirling":              "stirling",
		"calc.series.abel_summation":        "abel",
		"calc.series.power_uniqueness":      "power_uniqueness",
		"calc.series.grouping":              "grouping",
		"calc.series.generating":            "generating",
		"calc.series.arithmetico_geometric": "arithmetico_geometric",
		"calc.integral.cross_sections":      "cross_sections",
		"calc.integral.parametric_length":   "parametric_length",
		"calc.integral.parametric_area":     "parametric_area",
		"calc.integral.beta_gamma":          "beta_gamma",
		"calc.integral.unbounded":           "unbounded",
		"calc.integral.pumping":             "pumping",
		"calc.integral.simpson_38":          "simpson_38",
		"calc.integral.stieltjes":           "stieltjes",
		"calc.integral.holder":              "holder",
		"calc.integral.jensen":              "jensen",
		"calc.integral.radial_mass":         "radial_mass",
	} {
		reg.Register(id, &w22Gen{kind: kind})
	}
}

type w22Gen struct{ kind string }

func w22N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w22I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w22P(kp int, q, a, e string, f map[string]string) generator.Problem {
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: e, Facts: f, KPOrdinal: kp}
}

func w22Fact(n int) int {
	p := 1
	for i := 2; i <= n; i++ {
		p *= i
	}
	return p
}

func w22GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w22Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - i + 1) / i
	}
	return r
}

func w22Fib(k int) int {
	a, b := 1, 1
	for i := 0; i < k; i++ {
		a, b = b, a+b
	}
	return a
}

func (g *w22Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "limit":
		return w22Limit(kp)
	case "subsequence":
		return w22Subsequence(kp)
	case "zeta":
		return w22Zeta(kp)
	case "eta":
		return w22Eta(kp)
	case "bernoulli":
		return w22Bernoulli(kp)
	case "double":
		return w22Double(kp)
	case "stirling":
		return w22Stirling(kp)
	case "abel":
		return w22Abel(kp)
	case "power_uniqueness":
		return w22PowerUniqueness(kp)
	case "grouping":
		return w22Grouping(kp)
	case "generating":
		return w22Generating(kp)
	case "arithmetico_geometric":
		return w22ArithmeticoGeometric(kp)
	case "cross_sections":
		return w22CrossSections(kp)
	case "parametric_length":
		return w22ParametricLength(kp)
	case "parametric_area":
		return w22ParametricArea(kp)
	case "beta_gamma":
		return w22BetaGamma(kp)
	case "unbounded":
		return w22Unbounded(kp)
	case "pumping":
		return w22Pumping(kp)
	case "simpson_38":
		return w22Simpson38(kp)
	case "stieltjes":
		return w22Stieltjes(kp)
	case "holder":
		return w22Holder(kp)
	case "jensen":
		return w22Jensen(kp)
	default:
		return w22RadialMass(kp)
	}
}

func w22Limit(kp int) generator.Problem {
	switch kp {
	case 1:
		n0 := w22N(2, 11)
		f := map[string]string{}
		w22I(f, "n0", n0)
		return w22P(kp,
			fmt.Sprintf("For \\(a_n = 1/n\\), the smallest \\(N\\) such that \\(|a_n| < 1/%d\\) for every \\(n \\ge N\\). Enter \\(N\\). (enter a number)", n0),
			strconv.Itoa(n0+1),
			fmt.Sprintf("The inequality 1/n < 1/%d is the same as n > %d, so the first integer that works is %d.", n0, n0, n0+1),
			f)
	case 2:
		a := w22N(-9, 9)
		b := w22N(-9, 9)
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("A sequence has \\(\\lim a_n = %d\\) and \\(\\lim b_n = %d\\). Enter \\(\\lim (a_n+b_n)\\). (enter a number)", a, b),
			strconv.Itoa(a+b),
			fmt.Sprintf("The limit of a sum is the sum of the limits: %d + (%d) = %d.", a, b, a+b),
			f)
	default:
		r := w22N(1, 4)
		k := w22N(1, 9)
		p := r * k
		q := w22N(-9, 9)
		s := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "k", k)
		w22I(f, "p", p)
		w22I(f, "q", q)
		w22I(f, "s", s)
		return w22P(kp,
			fmt.Sprintf("Enter \\(\\lim_{n\\to\\infty} \\frac{%d n + %d}{%d n + %d}\\). (enter a number)", p, q, r, s),
			strconv.Itoa(k),
			fmt.Sprintf("Divide numerator and denominator by n; the constants vanish and the limit is %d/%d = %d.", p, r, k),
			f)
	}
}

func w22Subsequence(kp int) generator.Problem {
	switch kp {
	case 1:
		c := w22N(1, 9)
		d := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "c", c)
		w22I(f, "d", d)
		return w22P(kp,
			fmt.Sprintf("The sequence \\(a_n = %d + (-1)^n %d\\) has two subsequential limits. Enter their sum. (enter a number)", c, d),
			strconv.Itoa(2*c),
			fmt.Sprintf("The even terms tend to %d + %d and the odd terms to %d - %d; the two limits add to %d.", c, d, c, d, 2*c),
			f)
	case 2:
		a := w22N(1, 9)
		b := a + w22N(1, 9)
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("A sequence is bounded between %d and %d. Bolzano-Weierstrass gives a convergent subsequence whose limit lies in that interval. Enter the number of integers in the closed interval. (enter a number)", a, b),
			strconv.Itoa(b-a+1),
			fmt.Sprintf("The integers from %d to %d inclusive are %d in number.", a, b, b-a+1),
			f)
	default:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The sequence \\(a_n = \\frac{%d n}{n+1}\\) converges. Every subsequence converges to the same limit. Enter that limit. (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("Dividing by n gives %d/(1+1/n), which tends to %d, and a subsequence inherits the limit.", k, k),
			f)
	}
}

func w22Zeta(kp int) generator.Problem {
	switch kp {
	case 1:
		s := w22N(2, 9)
		den := 1 << uint(s)
		f := map[string]string{}
		w22I(f, "s", s)
		return w22P(kp,
			fmt.Sprintf("The zeta series is \\(\\zeta(s)=\\sum 1/n^s\\). For \\(s=%d\\), the term at \\(n=2\\) is \\(1/2^s\\). Enter its denominator. (enter a number)", s),
			strconv.Itoa(den),
			fmt.Sprintf("The denominator is 2 to the power %d, which is %d.", s, den),
			f)
	case 2:
		p := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "p", p)
		return w22P(kp,
			fmt.Sprintf("The integral test compares \\(\\sum 1/n^p\\) with \\(\\int_1^\\infty x^{-p}\\,dx = 1/(p-1)\\). For \\(p=%d\\), enter \\(p-1\\). (enter a number)", p),
			strconv.Itoa(p-1),
			fmt.Sprintf("Subtracting one from %d gives %d, and a positive value confirms convergence.", p, p-1),
			f)
	default:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("Euler's formula writes \\(\\zeta(2k)\\) as a rational multiple of \\(\\pi^{2k}\\). For \\(k=%d\\), enter the exponent of \\(\\pi\\). (enter a number)", k),
			strconv.Itoa(2*k),
			fmt.Sprintf("The exponent is 2k, and 2 times %d is %d.", k, 2*k),
			f)
	}
}

func w22Eta(kp int) generator.Problem {
	switch kp {
	case 1:
		s := w22N(2, 9)
		num := (1 << uint(s-1)) - 1
		f := map[string]string{}
		w22I(f, "s", s)
		return w22P(kp,
			fmt.Sprintf("The eta series satisfies \\(\\eta(s)=(1-2^{1-s})\\zeta(s)\\), and \\(1-2^{1-s}=(2^{s-1}-1)/2^{s-1}\\). For \\(s=%d\\), enter the numerator. (enter a number)", s),
			strconv.Itoa(num),
			fmt.Sprintf("The numerator is 2 to the power %d minus one, which is %d.", s-1, num),
			f)
	case 2:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("\\(\\eta(2k)\\) is a rational multiple of \\(\\pi^{2k}\\). For \\(k=%d\\), enter the exponent of \\(\\pi\\). (enter a number)", k),
			strconv.Itoa(2*k),
			fmt.Sprintf("The exponent is 2k, and 2 times %d is %d.", k, 2*k),
			f)
	default:
		n := w22N(1, 10)
		dens := []int{1, 2, 6, 12, 60, 60, 420, 840, 2520, 2520}
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("For the alternating harmonic partial sum \\(S_N=\\sum_{n=1}^{N}(-1)^{n-1}/n\\), take \\(N=%d\\). Enter the denominator of \\(S_N\\) in lowest terms. (enter a number)", n),
			strconv.Itoa(dens[n-1]),
			fmt.Sprintf("Combining the first %d terms gives a fraction whose reduced denominator is %d.", n, dens[n-1]),
			f)
	}
}

func w22Bernoulli(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w22N(1, 12)
		dens := []int{6, 30, 42, 30, 66, 2730, 6, 510, 798, 330, 138, 2730}
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("Write the Bernoulli number \\(B_{2k}\\) in lowest terms. For \\(k=%d\\), enter its denominator. (enter a number)", k),
			strconv.Itoa(dens[k-1]),
			fmt.Sprintf("The reduced denominator of B at index %d is %d.", 2*k, dens[k-1]),
			f)
	case 2:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("Euler's formula expresses \\(\\zeta(2k)\\) through \\(B_{2k}\\) and a power of \\(\\pi\\). For \\(k=%d\\), enter the exponent of \\(\\pi\\). (enter a number)", k),
			strconv.Itoa(2*k),
			fmt.Sprintf("The formula contains \\((2\\pi)^{2k}\\), so the exponent is %d.", 2*k),
			f)
	default:
		n := w22N(1, 12)
		dens := []int{2, 12, 0, 720, 0, 30240, 0, 1209600, 0, 47900160, 0, 1307674368000}
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("The coefficient of \\(x^n\\) in \\(x/(e^x-1)\\) is \\(B_n/n!\\). For \\(n=%d\\), enter its denominator in lowest terms, or 0 if the coefficient vanishes. (enter a number)", n),
			strconv.Itoa(dens[n-1]),
			fmt.Sprintf("At index %d the reduced denominator is %d, with 0 marking a vanishing odd term.", n, dens[n-1]),
			f)
	}
}

func w22Double(kp int) generator.Problem {
	switch kp {
	case 1:
		a := w22N(2, 6)
		b := w22N(2, 6)
		d := (a - 1) * (b - 1)
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("For \\(a=%d\\) and \\(b=%d\\), \\(\\sum_{i,j\\ge1} a^{-i} b^{-j} = 1/D\\). Enter \\(D\\). (enter a number)", a, b),
			strconv.Itoa(d),
			fmt.Sprintf("The two geometric sums are 1/(a-1) and 1/(b-1), so D is %d times %d, which is %d.", a-1, b-1, d),
			f)
	case 2:
		p := w22N(2, 9)
		q := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "p", p)
		w22I(f, "q", q)
		return w22P(kp,
			fmt.Sprintf("The double series \\(\\sum 1/(i^p j^q)\\) converges exactly when \\(p>1\\) and \\(q>1\\). For \\(p=%d\\) and \\(q=%d\\), enter \\(p+q\\). (enter a number)", p, q),
			strconv.Itoa(p+q),
			fmt.Sprintf("Both exponents exceed one, and their sum is %d + %d = %d.", p, q, p+q),
			f)
	default:
		ps := []int{2, 4, 6, 8}
		zetaDen := map[int]int{2: 6, 4: 90, 6: 945, 8: 9450}
		p := ps[rand.Intn(len(ps))]
		q := ps[rand.Intn(len(ps))]
		d := zetaDen[p] * zetaDen[q]
		f := map[string]string{}
		w22I(f, "p", p)
		w22I(f, "q", q)
		return w22P(kp,
			fmt.Sprintf("For \\(p=%d\\) and \\(q=%d\\), \\(\\sum 1/(i^{2p} j^{2q}) = \\pi^{2p+2q}/D\\). Enter \\(D\\). (enter a number)", p, q),
			strconv.Itoa(d),
			fmt.Sprintf("The product separates into zeta factors with denominators %d and %d, giving D = %d.", zetaDen[p], zetaDen[q], d),
			f)
	}
}

func w22Stirling(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("Stirling's formula is \\(n! \\sim \\sqrt{2\\pi n}\\,(n/e)^n\\). For \\(n=%d\\), the value under the square root is \\(2\\pi n\\). Enter \\(2n\\). (enter a number)", n),
			strconv.Itoa(2*n),
			fmt.Sprintf("Twice %d is %d, which is the factor multiplying pi under the root.", n, 2*n),
			f)
	case 2:
		ns := []int{5, 10, 15, 20, 25, 30, 40, 50, 75, 100}
		digits := []int{3, 7, 13, 19, 26, 33, 48, 65, 110, 158}
		i := rand.Intn(len(ns))
		n := ns[i]
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("Use Stirling's approximation to count the decimal digits of \\(n!\\). For \\(n=%d\\), enter the number of digits. (enter a number)", n),
			strconv.Itoa(digits[i]),
			fmt.Sprintf("Stirling gives log base ten of the factorial just below %d, so the number has %d digits.", digits[i], digits[i]),
			f)
	default:
		n := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("Stirling's correction factor is \\(1+1/(12n)\\). For \\(n=%d\\), enter \\(12n\\). (enter a number)", n),
			strconv.Itoa(12*n),
			fmt.Sprintf("Twelve times %d is %d, the denominator of the correction.", n, 12*n),
			f)
	}
}

func w22Abel(kp int) generator.Problem {
	switch kp {
	case 1:
		c := w22N(1, 9)
		n := w22N(2, 12)
		f := map[string]string{}
		w22I(f, "c", c)
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("Abel summation uses the partial sums \\(A_N=\\sum_{n=1}^{N} a_n\\). If every \\(a_n=%d\\) and \\(N=%d\\), enter \\(A_N\\). (enter a number)", c, n),
			strconv.Itoa(c*n),
			fmt.Sprintf("There are %d terms of size %d, so the partial sum is %d.", n, c, c*n),
			f)
	case 2:
		n := w22N(2, 12)
		ans := n * (n + 1) / 2
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("If \\(a_n=n\\) and \\(N=%d\\), the partial sum \\(A_N\\) is the sum of the first N integers. Enter \\(A_N\\). (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The triangular sum is %d times %d divided by two, which is %d.", n, n+1, ans),
			f)
	default:
		m := w22N(1, 9)
		b := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "m", m)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("Dirichlet's test bounds \\(|\\sum a_n b_n|\\) by \\(M b_1\\) when the partial sums of \\(a_n\\) are bounded by \\(M\\) and \\(b_n\\) decreases to zero. For \\(M=%d\\) and \\(b_1=%d\\), enter the bound. (enter a number)", m, b),
			strconv.Itoa(m*b),
			fmt.Sprintf("The bound is the product %d times %d, which is %d.", m, b, m*b),
			f)
	}
}

func w22PowerUniqueness(kp int) generator.Problem {
	switch kp {
	case 1:
		a := w22N(2, 5)
		k := w22N(1, 5)
		ans := 1
		for i := 0; i < k; i++ {
			ans *= a
		}
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The series \\(1/(1-%d x)=\\sum %d^n x^n\\) is unique. Enter the coefficient of \\(x^k\\) for \\(k=%d\\). (enter a number)", a, a, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is %d to the power %d, which is %d.", a, k, ans),
			f)
	case 2:
		c := w22N(2, 5)
		n := w22N(1, 4)
		ans := 1
		for i := 0; i < n; i++ {
			ans *= c
		}
		f := map[string]string{}
		w22I(f, "c", c)
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("For \\(e^{cx}=\\sum c^n x^n/n!\\), the coefficient of \\(x^n\\) is \\(c^n/n!\\). For \\(c=%d\\) and \\(n=%d\\), enter the numerator \\(c^n\\). (enter a number)", c, n),
			strconv.Itoa(ans),
			fmt.Sprintf("The numerator is %d to the power %d, which is %d.", c, n, ans),
			f)
	default:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The identity \\(x/(1-x)^2=\\sum n x^n\\) fixes every coefficient. Enter the coefficient of \\(x^k\\) for \\(k=%d\\). (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The coefficient of the k-th power is k itself, so at %d it is %d.", k, k),
			f)
	}
}

func w22Grouping(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w22N(1, 9)
		d := (2*k - 1) * (2 * k)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The k-th group of the grouped alternating harmonic series is \\(1/(2k-1)-1/(2k)=1/D\\). For \\(k=%d\\), enter \\(D\\). (enter a number)", k),
			strconv.Itoa(d),
			fmt.Sprintf("The denominator is the product of %d and %d, which is %d.", 2*k-1, 2*k, d),
			f)
	case 2:
		n := w22N(1, 9)
		ans := 4*n - 1
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("The n-th group is \\(1/(2n-1)-1/(2n)\\). For \\(n=%d\\), enter the sum of its two denominators. (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The denominators are %d and %d, whose sum is %d.", 2*n-1, 2*n, ans),
			f)
	default:
		k := w22N(1, 9)
		ans := 1 << uint(k)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("A geometric series is grouped into blocks of k terms. The first block sums to \\((2^k-1)/2^{k-1}\\). For \\(k=%d\\), enter \\(2^k\\). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("Two to the power %d is %d.", k, ans),
			f)
	}
}

func w22Generating(kp int) generator.Problem {
	switch kp {
	case 1:
		c := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "c", c)
		return w22P(kp,
			fmt.Sprintf("The sequence \\(a_n=%d^n\\) has generating function \\(1/(1-cx)\\). Enter \\(c\\). (enter a number)", c),
			strconv.Itoa(c),
			fmt.Sprintf("The ratio between consecutive terms is %d, so the generating function is 1/(1-%d x).", c, c),
			f)
	case 2:
		m := w22N(2, 5)
		k := w22N(1, 6)
		ans := w22Binom(k+m-1, m-1)
		f := map[string]string{}
		w22I(f, "m", m)
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The expansion \\(1/(1-x)^m=\\sum \\binom{n+m-1}{m-1}x^n\\) gives every coefficient. For \\(m=%d\\) and \\(k=%d\\), enter the coefficient of \\(x^k\\). (enter a number)", m, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is the binomial number with top %d and bottom %d, which is %d.", k+m-1, m-1, ans),
			f)
	default:
		k := w22N(1, 10)
		ans := w22Fib(k)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The Fibonacci generating function \\(1/(1-x-x^2)\\) has coefficients \\(F_0=F_1=1\\) and \\(F_n=F_{n-1}+F_{n-2}\\). Enter \\(F_k\\) for \\(k=%d\\). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("Carrying the recurrence to index %d gives %d.", k, ans),
			f)
	}
}

func w22ArithmeticoGeometric(kp int) generator.Problem {
	switch kp {
	case 1:
		a := w22N(1, 9)
		d := w22N(1, 9)
		ans := 2*a + 2*d
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "d", d)
		return w22P(kp,
			fmt.Sprintf("For \\(a=%d\\) and \\(d=%d\\), evaluate \\(\\sum_{n\\ge0}(a+nd)/2^n\\). (enter a number)", a, d),
			strconv.Itoa(ans),
			fmt.Sprintf("The two standard sums are 2 and 2, so the total is 2 times %d plus 2 times %d, which is %d.", a, d, ans),
			f)
	case 2:
		A := w22N(1, 9)
		D := w22N(1, 9)
		ans := 3*A + 3*D
		f := map[string]string{}
		w22I(f, "a", 2*A)
		w22I(f, "d", 4*D)
		w22I(f, "A", A)
		w22I(f, "D", D)
		return w22P(kp,
			fmt.Sprintf("For \\(a=%d\\) and \\(d=%d\\), evaluate \\(\\sum_{n\\ge0}(a+nd)/3^n\\). (enter a number)", 2*A, 4*D),
			strconv.Itoa(ans),
			fmt.Sprintf("The sum is 3/2 times %d plus 3/4 times %d, which simplifies to %d.", 2*A, 4*D, ans),
			f)
	default:
		A := w22N(1, 9)
		D := w22N(1, 9)
		ans := 4*A + 4*D
		f := map[string]string{}
		w22I(f, "a", 3*A)
		w22I(f, "d", 9*D)
		w22I(f, "A", A)
		w22I(f, "D", D)
		return w22P(kp,
			fmt.Sprintf("For \\(a=%d\\) and \\(d=%d\\), evaluate \\(\\sum_{n\\ge0}(a+nd)/4^n\\). (enter a number)", 3*A, 9*D),
			strconv.Itoa(ans),
			fmt.Sprintf("The sum is 4/3 times %d plus 4/9 times %d, which simplifies to %d.", 3*A, 9*D, ans),
			f)
	}
}

func w22CrossSections(kp int) generator.Problem {
	switch kp {
	case 1:
		s := w22N(2, 9)
		h := w22N(2, 9)
		ans := s * s * h
		f := map[string]string{}
		w22I(f, "s", s)
		w22I(f, "h", h)
		return w22P(kp,
			fmt.Sprintf("A solid of height %d has square cross sections of side %d. Enter its volume. (enter a number)", h, s),
			strconv.Itoa(ans),
			fmt.Sprintf("Each slice has area %d times %d, and stacking %d slices gives %d.", s, s, h, ans),
			f)
	case 2:
		a := w22N(1, 5)
		b := 3 * a
		h := 3 * w22N(1, 5)
		ans := b * b * h / 3
		f := map[string]string{}
		w22I(f, "b", b)
		w22I(f, "h", h)
		return w22P(kp,
			fmt.Sprintf("A pyramid has a square base of side %d and height %d. Enter its volume. (enter a number)", b, h),
			strconv.Itoa(ans),
			fmt.Sprintf("The volume is one third of the base area times the height: %d times %d divided by three is %d.", b*b, h, ans),
			f)
	default:
		c := w22N(1, 9)
		a := 3 * w22N(1, 5)
		ans := c * a * a * a / 3
		f := map[string]string{}
		w22I(f, "c", c)
		w22I(f, "a", a)
		return w22P(kp,
			fmt.Sprintf("A solid has cross-sectional area \\(A(x)=%d x^2\\) for \\(0\\le x\\le %d\\). Enter its volume. (enter a number)", c, a),
			strconv.Itoa(ans),
			fmt.Sprintf("Integrating gives %d times %d cubed divided by three, which is %d.", c, a, ans),
			f)
	}
}

func w22ParametricLength(kp int) generator.Problem {
	switch kp {
	case 1:
		trips := [][3]int{{3, 4, 5}, {6, 8, 10}, {5, 12, 13}, {8, 15, 17}, {7, 24, 25}, {9, 40, 41}, {20, 21, 29}, {12, 35, 37}}
		t := trips[rand.Intn(len(trips))]
		f := map[string]string{}
		w22I(f, "a", t[0])
		w22I(f, "b", t[1])
		return w22P(kp,
			fmt.Sprintf("The curve \\((x,y)=(%d t,%d t)\\) runs from \\(t=0\\) to \\(t=1\\). Enter its length. (enter a number)", t[0], t[1]),
			strconv.Itoa(t[2]),
			fmt.Sprintf("The constant speed is the square root of %d squared plus %d squared, which is %d.", t[0], t[1], t[2]),
			f)
	case 2:
		a := w22N(2, 9)
		b := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("The ellipse \\((x,y)=(%d\\cos t,%d\\sin t)\\) has speed \\(\\sqrt{a^2\\sin^2 t+b^2\\cos^2 t}\\). For these values, enter the speed at \\(t=0\\). (enter a number)", a, b),
			strconv.Itoa(b),
			fmt.Sprintf("At t equal to zero the sine term vanishes and the speed is the second semiaxis, %d.", b),
			f)
	default:
		quads := [][4]int{{1, 2, 2, 3}, {2, 3, 6, 7}, {1, 4, 8, 9}, {4, 4, 7, 9}, {2, 6, 9, 11}, {6, 6, 7, 11}, {3, 4, 12, 13}, {2, 10, 11, 15}}
		q := quads[rand.Intn(len(quads))]
		f := map[string]string{}
		w22I(f, "a", q[0])
		w22I(f, "b", q[1])
		w22I(f, "c", q[2])
		return w22P(kp,
			fmt.Sprintf("The line \\((x,y,z)=(%d t,%d t,%d t)\\) runs from \\(t=0\\) to \\(t=1\\). Enter its length. (enter a number)", q[0], q[1], q[2]),
			strconv.Itoa(q[3]),
			fmt.Sprintf("The length is the square root of the sum of the three squares, which is %d.", q[3]),
			f)
	}
}

func w22ParametricArea(kp int) generator.Problem {
	switch kp {
	case 1:
		a := 3 * w22N(1, 5)
		ans := a * a * a / 3
		f := map[string]string{}
		w22I(f, "a", a)
		return w22P(kp,
			fmt.Sprintf("The curve \\((x,y)=(t,t^2)\\) runs from \\(t=0\\) to \\(t=%d\\). Enter the area under it. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("The area is the integral of y dx, which is %d cubed over three, or %d.", a, ans),
			f)
	case 2:
		a := 3 * w22N(1, 5)
		ans := 2 * a * a * a / 3
		f := map[string]string{}
		w22I(f, "a", a)
		return w22P(kp,
			fmt.Sprintf("The curve \\((x,y)=(t^2,t)\\) runs from \\(t=0\\) to \\(t=%d\\). Enter the area under it. (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("Here dx is 2t dt, so the area is twice %d cubed over three, or %d.", a, ans),
			f)
	default:
		A := w22N(1, 6)
		B := w22N(1, 6)
		a := 2 * A
		b := 2 * B
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("The ellipse with semiaxes \\(a=%d\\) and \\(b=%d\\) bounds an area \\(\\pi ab/4\\) in the first quadrant. Enter the coefficient of \\(\\pi\\). (enter a number)", a, b),
			strconv.Itoa(A*B),
			fmt.Sprintf("The coefficient is %d times %d divided by four, which is %d.", a, b, A*B),
			f)
	}
}

func w22BetaGamma(kp int) generator.Problem {
	switch kp {
	case 1:
		p := w22N(2, 6)
		q := w22N(2, 6)
		num := w22Fact(p-1) * w22Fact(q-1)
		den := w22Fact(p + q - 1)
		g := w22GCD(num, den)
		ans := den / g
		f := map[string]string{}
		w22I(f, "p", p)
		w22I(f, "q", q)
		return w22P(kp,
			fmt.Sprintf("The beta integral satisfies \\(B(p,q)=(p-1)!(q-1)!/(p+q-1)!\\). For \\(p=%d\\) and \\(q=%d\\), enter the denominator of \\(B(p,q)\\) in lowest terms. (enter a number)", p, q),
			strconv.Itoa(ans),
			fmt.Sprintf("The unreduced denominator is %d; cancelling the common factor %d leaves %d.", den, g, ans),
			f)
	case 2:
		k := w22N(2, 10)
		ans := w22Fact(k - 1)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("The gamma function satisfies \\(\\Gamma(k)=(k-1)!\\). For \\(k=%d\\), enter \\(\\Gamma(k)\\). (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("The factorial of %d is %d.", k-1, ans),
			f)
	default:
		p := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "p", p)
		return w22P(kp,
			fmt.Sprintf("For \\(p=%d\\), \\(B(p,1)=1/p\\). Enter the denominator in lowest terms. (enter a number)", p),
			strconv.Itoa(p),
			fmt.Sprintf("The fraction is 1 over %d, so the denominator is %d.", p, p),
			f)
	}
}

func w22Unbounded(kp int) generator.Problem {
	switch kp {
	case 1:
		k := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("For \\(p=1-1/k\\) with \\(k=%d\\), the improper integral \\(\\int_0^1 x^{-p}\\,dx\\) equals \\(k\\). Enter its value. (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The antiderivative is x to the 1-p over 1-p, and at this exponent the value is %d.", k),
			f)
	case 2:
		m := w22N(1, 9)
		ans := 2 * m
		f := map[string]string{}
		w22I(f, "m", m)
		return w22P(kp,
			fmt.Sprintf("Enter \\(\\int_0^{%d} x^{-1/2}\\,dx\\). (enter a number)", m*m),
			strconv.Itoa(ans),
			fmt.Sprintf("The antiderivative is 2 square root x, evaluated at %d gives twice %d, or %d.", m*m, m, ans),
			f)
	default:
		k := w22N(1, 9)
		f := map[string]string{}
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("Enter \\(\\int_{e^{-k}}^{1} \\frac{dx}{x}\\) for \\(k=%d\\). (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("The integral is minus the natural logarithm of the lower limit, which is %d.", k),
			f)
	}
}

func w22Pumping(kp int) generator.Problem {
	switch kp {
	case 1:
		r := w22N(1, 5)
		h := 2 * w22N(1, 5)
		ans := r * r * h * h / 2
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "h", h)
		return w22P(kp,
			fmt.Sprintf("A full cylindrical tank of radius %d and height %d is pumped to the top. With unit weight density, \\(W/\\pi=r^2h^2/2\\). Enter \\(W/\\pi\\). (enter a number)", r, h),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is %d squared times %d squared divided by two, which is %d.", r, h, ans),
			f)
	case 2:
		l := w22N(1, 5)
		w := w22N(1, 5)
		h := 2 * w22N(1, 5)
		ans := l * w * h * h / 2
		f := map[string]string{}
		w22I(f, "l", l)
		w22I(f, "w", w)
		w22I(f, "h", h)
		return w22P(kp,
			fmt.Sprintf("A full rectangular tank %d by %d by %d is pumped to the top. With unit weight density, \\(W=lwh^2/2\\). Enter \\(W\\). (enter a number)", l, w, h),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is %d times %d times %d squared divided by two, which is %d.", l, w, h, ans),
			f)
	default:
		r := w22N(1, 5)
		h := 2 * w22N(1, 5)
		extra := w22N(1, 5)
		ans := r * r * h * (h/2 + extra)
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "h", h)
		w22I(f, "extra", extra)
		return w22P(kp,
			fmt.Sprintf("A full cylindrical tank of radius %d and height %d is pumped %d units above its top. With unit weight density, \\(W/\\pi=r^2h(h/2+e)\\) where e is the extra height. Enter \\(W/\\pi\\). (enter a number)", r, h, extra),
			strconv.Itoa(ans),
			fmt.Sprintf("The bracket is %d/2 + %d, and multiplying by %d squared times %d gives %d.", h, extra, r, h, ans),
			f)
	}
}

func w22Simpson38(kp int) generator.Problem {
	switch kp {
	case 1:
		f0 := w22N(1, 9)
		f1 := w22N(1, 9)
		f2 := w22N(1, 9)
		f3 := w22N(1, 9)
		ans := f0 + 3*f1 + 3*f2 + f3
		f := map[string]string{}
		w22I(f, "f0", f0)
		w22I(f, "f1", f1)
		w22I(f, "f2", f2)
		w22I(f, "f3", f3)
		return w22P(kp,
			fmt.Sprintf("One panel of the three-eighths rule uses weights 1, 3, 3, 1. For values \\(f_0=%d\\), \\(f_1=%d\\), \\(f_2=%d\\) and \\(f_3=%d\\), enter the weighted sum. (enter a number)", f0, f1, f2, f3),
			strconv.Itoa(ans),
			fmt.Sprintf("The sum is %d plus three times %d plus three times %d plus %d, which is %d.", f0, f1, f2, f3, ans),
			f)
	case 2:
		v := [7]int{}
		for i := range v {
			v[i] = w22N(1, 9)
		}
		ans := v[0] + 3*v[1] + 3*v[2] + 2*v[3] + 3*v[4] + 3*v[5] + v[6]
		f := map[string]string{}
		for i := range v {
			w22I(f, fmt.Sprintf("f%d", i), v[i])
		}
		return w22P(kp,
			fmt.Sprintf("A composite three-eighths rule over six subintervals uses weights 1, 3, 3, 2, 3, 3, 1. For values \\(f_0=%d\\), \\(f_1=%d\\), \\(f_2=%d\\), \\(f_3=%d\\), \\(f_4=%d\\), \\(f_5=%d\\) and \\(f_6=%d\\), enter the weighted sum. (enter a number)", v[0], v[1], v[2], v[3], v[4], v[5], v[6]),
			strconv.Itoa(ans),
			fmt.Sprintf("Applying the weights in order gives the total %d.", ans),
			f)
	default:
		a := w22N(0, 5)
		step := w22N(1, 4)
		b := a + 3*step
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("The three-eighths rule on \\([%d,%d]\\) uses three subintervals. Enter the step size \\(h=(b-a)/3\\). (enter a number)", a, b),
			strconv.Itoa(step),
			fmt.Sprintf("The width is %d, and dividing by three gives %d.", b-a, step),
			f)
	}
}

func w22Stieltjes(kp int) generator.Problem {
	switch kp {
	case 1:
		n := w22N(2, 12)
		ans := n * (n + 1) / 2
		f := map[string]string{}
		w22I(f, "n", n)
		return w22P(kp,
			fmt.Sprintf("The integrator \\(g(x)=\\lfloor x\\rfloor\\) jumps by one at each positive integer. With \\(f(x)=x\\) and an upper limit of %d, enter \\(\\int_0^{%d} f\\,dg\\). (enter a number)", n, n),
			strconv.Itoa(ans),
			fmt.Sprintf("The jumps collect f at 1 through %d, a triangular sum equal to %d.", n, ans),
			f)
	case 2:
		c1 := w22N(1, 5)
		c2 := w22N(1, 5)
		c3 := w22N(1, 5)
		ans := c1 + 2*c2 + 3*c3
		f := map[string]string{}
		w22I(f, "c1", c1)
		w22I(f, "c2", c2)
		w22I(f, "c3", c3)
		return w22P(kp,
			fmt.Sprintf("An integrator jumps by %d, %d and %d at x equal to 1, 2 and 3. With \\(f(x)=x\\), enter \\(\\int_0^3 f\\,dg\\). (enter a number)", c1, c2, c3),
			strconv.Itoa(ans),
			fmt.Sprintf("Each jump contributes its size times its location: %d + %d + %d = %d.", c1, 2*c2, 3*c3, ans),
			f)
	default:
		a := w22N(2, 9)
		ans := a * a
		f := map[string]string{}
		w22I(f, "a", a)
		return w22P(kp,
			fmt.Sprintf("For \\(g(x)=x^2\\) and \\(f(x)=1\\), enter \\(\\int_0^{%d} f\\,dg\\). (enter a number)", a),
			strconv.Itoa(ans),
			fmt.Sprintf("The integral is g evaluated between zero and %d, which is %d.", a, ans),
			f)
	}
}

func w22Holder(kp int) generator.Problem {
	switch kp {
	case 1:
		p := w22N(2, 9)
		f := map[string]string{}
		w22I(f, "p", p)
		return w22P(kp,
			fmt.Sprintf("The conjugate exponent satisfies \\(1/p+1/q=1\\), so \\(q=p/(p-1)\\). For \\(p=%d\\), enter the denominator of \\(q\\) in lowest terms. (enter a number)", p),
			strconv.Itoa(p-1),
			fmt.Sprintf("The fraction is %d over %d in lowest terms, so the denominator is %d.", p, p-1, p-1),
			f)
	case 2:
		a := w22N(1, 9)
		p := w22N(2, 5)
		ans := a*p + 1
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "p", p)
		return w22P(kp,
			fmt.Sprintf("Holder's inequality pairs \\(f\\) with an exponent \\(p\\). For \\(f(x)=x^a\\) on the unit interval, \\(\\int_0^1 f^p\\,dx=1/(ap+1)\\). For \\(a=%d\\) and \\(p=%d\\), enter the denominator. (enter a number)", a, p),
			strconv.Itoa(ans),
			fmt.Sprintf("The exponent product is %d, and adding one gives %d.", a*p, ans),
			f)
	default:
		p := w22N(2, 9)
		ans := 2*p - 1
		f := map[string]string{}
		w22I(f, "p", p)
		return w22P(kp,
			fmt.Sprintf("For \\(p=%d\\), the conjugate is \\(q=p/(p-1)\\). Enter the sum of the numerator and denominator of \\(q\\) in lowest terms. (enter a number)", p),
			strconv.Itoa(ans),
			fmt.Sprintf("The fraction is %d over %d, and its numerator plus denominator is %d.", p, p-1, ans),
			f)
	}
}

func w22Jensen(kp int) generator.Problem {
	switch kp {
	case 1:
		a := w22N(1, 9)
		m := w22N(1, 5)
		b := a + 2*m
		ans := m * m
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("For the convex function \\(f(x)=x^2\\), the Jensen gap at %d and %d is \\((f(a)+f(b))/2-f((a+b)/2)\\). Enter the gap. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The gap equals half the difference squared, and half of %d squared is %d.", b-a, ans),
			f)
	case 2:
		A := w22N(1, 5)
		l := 6 * A
		ans := l * l / 12
		f := map[string]string{}
		w22I(f, "l", l)
		return w22P(kp,
			fmt.Sprintf("For \\(f(x)=x^2\\) on \\([0,%d]\\), the integral form of Jensen gives a gap of \\(L^2/12\\). Enter the gap. (enter a number)", l),
			strconv.Itoa(ans),
			fmt.Sprintf("Squaring %d and dividing by twelve gives %d.", l, ans),
			f)
	default:
		a := w22N(0, 9)
		m := w22N(1, 5)
		b := a + 3*m
		ans := 2 * m * m
		f := map[string]string{}
		w22I(f, "a", a)
		w22I(f, "b", b)
		return w22P(kp,
			fmt.Sprintf("For \\(f(x)=x^2\\) with weights 1/3 and 2/3 at %d and %d, the weighted Jensen gap is \\((2/9)(b-a)^2\\). Enter the gap. (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The difference is %d, and two ninths of its square is %d.", b-a, ans),
			f)
	}
}

func w22RadialMass(kp int) generator.Problem {
	switch kp {
	case 1:
		r := w22N(1, 9)
		rho := w22N(1, 9)
		ans := rho * r * r
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "rho", rho)
		return w22P(kp,
			fmt.Sprintf("A disk of radius %d has constant density %d. Its mass is \\(\\pi\\rho R^2\\). Enter \\(M/\\pi\\). (enter a number)", r, rho),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is %d times %d squared, which is %d.", rho, r, ans),
			f)
	case 2:
		A := w22N(1, 5)
		r := 3 * A
		k := w22N(1, 9)
		ans := 18 * k * A * A * A
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("A disk of radius %d has density \\(\\rho(r)=%d r\\). Its mass is \\(2\\pi k R^3/3\\). Enter \\(M/\\pi\\). (enter a number)", r, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is two thirds of %d times %d cubed, which is %d.", k, r, ans),
			f)
	default:
		A := w22N(1, 5)
		r := 3 * A
		k := w22N(1, 9)
		ans := 9 * k * A * A * A
		f := map[string]string{}
		w22I(f, "r", r)
		w22I(f, "k", k)
		return w22P(kp,
			fmt.Sprintf("A disk of radius %d has density \\(\\rho(r)=%d(R-r)\\). Its mass is \\(\\pi k R^3/3\\). Enter \\(M/\\pi\\). (enter a number)", r, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The coefficient is one third of %d times %d cubed, which is %d.", k, r, ans),
			f)
	}
}
