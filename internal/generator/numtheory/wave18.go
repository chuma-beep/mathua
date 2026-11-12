package numtheory

// Wave 18 number theory: analytic number theory, applications, and Diophantine
// equations. Every family publishes the parameters it drew in Facts so
// wave18_test.go can recompute the answer.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave18(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"nt.analytic.dirichlet_series":    "dirichlet_series",
		"nt.analytic.euler_product":       "euler_product",
		"nt.analytic.continuation":        "continuation",
		"nt.analytic.functional_equation": "functional_equation",
		"nt.analytic.zero_count":          "zero_count",
		"nt.analytic.explicit_formula":    "explicit_formula",
		"nt.analytic.mertens":             "mertens",
		"nt.analytic.abel":                "abel",
		"nt.analytic.dirichlet_l":         "dirichlet_l",
		"nt.analytic.bias":                "bias",
		"nt.app.rsa_keygen":               "rsa_keygen",
		"nt.app.elgamal":                  "elgamal",
		"nt.app.secret_sharing":           "secret_sharing",
		"nt.app.checksum":                 "checksum",
		"nt.app.zero_knowledge":           "zero_knowledge",
		"nt.app.ecc":                      "ecc",
		"nt.app.primality_cert":           "primality_cert",
		"nt.app.reed_solomon":             "reed_solomon",
		"nt.app.factorization":            "factorization",
		"nt.app.universal_hashing":        "universal_hashing",
		"nt.dio.three_squares":            "three_squares",
		"nt.dio.waring":                   "waring",
		"nt.dio.markov":                   "markov",
		"nt.dio.mordell":                  "mordell",
		"nt.dio.catalan":                  "catalan",
		"nt.dio.abc":                      "abc",
		"nt.dio.taxicab":                  "taxicab",
		"nt.dio.frobenius":                "frobenius",
		"nt.dio.fermat_last":              "fermat_last",
		"nt.dio.four_cubes":               "four_cubes",
	} {
		reg.Register(id, &w18Gen{kind: kind})
	}
}

type w18Gen struct{ kind string }

func (g *w18Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "dirichlet_series":
		return w18DirichletSeries(o)
	case "euler_product":
		return w18EulerProduct(o)
	case "continuation":
		return w18Continuation(o)
	case "functional_equation":
		return w18FunctionalEquation(o)
	case "zero_count":
		return w18ZeroCount(o)
	case "explicit_formula":
		return w18ExplicitFormula(o)
	case "mertens":
		return w18Mertens(o)
	case "abel":
		return w18Abel(o)
	case "dirichlet_l":
		return w18DirichletL(o)
	case "bias":
		return w18Bias(o)
	case "rsa_keygen":
		return w18RSA(o)
	case "elgamal":
		return w18ElGamal(o)
	case "secret_sharing":
		return w18SecretSharing(o)
	case "checksum":
		return w18Checksum(o)
	case "zero_knowledge":
		return w18ZeroKnowledge(o)
	case "ecc":
		return w18ECC(o)
	case "primality_cert":
		return w18PrimalityCert(o)
	case "reed_solomon":
		return w18ReedSolomon(o)
	case "factorization":
		return w18Factorization(o)
	case "universal_hashing":
		return w18UniversalHashing(o)
	case "three_squares":
		return w18ThreeSquares(o)
	case "waring":
		return w18Waring(o)
	case "markov":
		return w18Markov(o)
	case "mordell":
		return w18Mordell(o)
	case "catalan":
		return w18Catalan(o)
	case "abc":
		return w18ABC(o)
	case "taxicab":
		return w18Taxicab(o)
	case "frobenius":
		return w18Frobenius(o)
	case "fermat_last":
		return w18FermatLast(o)
	case "four_cubes":
		return w18FourCubes(o)
	}
	return w18RestGenerate(g.kind, o)
}

func w18N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w18P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w18I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w18GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w18ExtGCD(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	g, x, y := w18ExtGCD(b, a%b)
	return g, y, x - (a/b)*y
}

func w18ModInverse(a, m int) int {
	g, x, _ := w18ExtGCD(((a%m)+m)%m, m)
	if g != 1 {
		return 0
	}
	return ((x % m) + m) % m
}

func w18PowMod(base, exp, mod int) int {
	result := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return result
}

func w18IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func w18PrimesUpTo(n int) []int {
	var out []int
	for i := 2; i <= n; i++ {
		if w18IsPrime(i) {
			out = append(out, i)
		}
	}
	return out
}

func w18FactorCounts(n int) map[int]int {
	out := map[int]int{}
	for d := 2; d*d <= n; d++ {
		for n%d == 0 {
			out[d]++
			n /= d
		}
	}
	if n > 1 {
		out[n]++
	}
	return out
}

func w18Omega(n int) int {
	total := 0
	for _, e := range w18FactorCounts(n) {
		total += e
	}
	return total
}

func w18DistinctFactors(n int) int { return len(w18FactorCounts(n)) }

func w18Totient(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if w18GCD(k, n) == 1 {
			count++
		}
	}
	return count
}

func w18FloorSqrt(n int) int {
	r := int(math.Sqrt(float64(n)))
	for (r+1)*(r+1) <= n {
		r++
	}
	for r*r > n {
		r--
	}
	return r
}

func w18CeilSqrt(n int) int {
	r := w18FloorSqrt(n)
	if r*r < n {
		r++
	}
	return r
}

func w18Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func w18IPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func w18Mobius(n int) int {
	counts := w18FactorCounts(n)
	for _, e := range counts {
		if e > 1 {
			return 0
		}
	}
	if len(counts)%2 == 0 {
		return 1
	}
	return -1
}

func w18SquarefreeCount(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if w18Mobius(k) != 0 {
			count++
		}
	}
	return count
}

// ---------------------------------------------------------------------------
// Analytic number theory
// ---------------------------------------------------------------------------

func w18DirichletSeries(o int) generator.Problem {
	switch o {
	case 1:
		p := []int{2, 3, 5, 7, 11}[rand.Intn(5)]
		a := w18N(1, 5)
		k := w18N(1, 4)
		ans := w18Binom(a+k-1, k-1)
		f := map[string]string{}
		w18I(f, "p", p)
		w18I(f, "a", a)
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("In the Dirichlet series for the k-th power of zeta with k=%d, what is the coefficient of n^{-s} when n is %d to the power %d?", k, p, a), strconv.Itoa(ans), fmt.Sprintf("The local factor is (1 + p^{-s} + p^{-2s} + ...)^k, so the coefficient at p^{-a s} is C(a+k-1, k-1) = %d.", ans), f)
	case 2:
		primes := []int{2, 3, 5, 7, 11}
		p := primes[rand.Intn(len(primes))]
		n := 1
		if rand.Intn(2) == 0 {
			n = w18IPow(p, w18N(1, 4))
		} else {
			q := primes[rand.Intn(len(primes))]
			for q == p {
				q = primes[rand.Intn(len(primes))]
			}
			n = p * q
			if rand.Intn(2) == 0 {
				r := primes[rand.Intn(len(primes))]
				for r == p || r == q {
					r = primes[rand.Intn(len(primes))]
				}
				n *= r
			}
		}
		mu := w18Mobius(n)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("The coefficient of n^{-s} in 1/zeta(s) is the Mobius value mu(n). What is mu(%d)?", n), strconv.Itoa(mu), fmt.Sprintf("Factor %d and count its distinct primes; mu is zero on a repeated prime, otherwise minus one to the number of distinct primes.", n), f)
	default:
		p := []int{2, 3, 5}[rand.Intn(3)]
		a := w18N(1, 5)
		if rand.Intn(2) == 0 {
			n := w18IPow(p, a)
			ans := a + 1
			f := map[string]string{}
			w18I(f, "n", n)
			return w18P(o, fmt.Sprintf("The coefficient of n^{-s} in the square of zeta is the divisor count d(n). What is d(%d)?", n), strconv.Itoa(ans), fmt.Sprintf("The divisor count of p^a is a+1, and here the exponent is %d.", a), f)
		}
		q := []int{3, 5, 7}[rand.Intn(3)]
		if q == p {
			q = 7
		}
		b := w18N(1, 3)
		n := w18IPow(p, a) * w18IPow(q, b)
		ans := (a + 1) * (b + 1)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("The coefficient of n^{-s} in the square of zeta is the divisor count d(n). What is d(%d)?", n), strconv.Itoa(ans), fmt.Sprintf("Factor n and multiply one plus each prime exponent, giving %d.", ans), f)
	}
}

func w18EulerProduct(o int) generator.Problem {
	switch o {
	case 1:
		n := 10 + 5*w18N(0, 38)
		ans := len(w18PrimesUpTo(n))
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Truncating the Euler product at the primes up to %d leaves one local factor per such prime. How many factors are there?", n), strconv.Itoa(ans), fmt.Sprintf("Count the primes at most %d, which is %d.", n, ans), f)
	case 2:
		bound := []int{3, 5, 7, 11}[rand.Intn(4)]
		k := w18N(1, 4)
		primes := w18PrimesUpTo(bound)
		ans := 1
		for range primes {
			ans *= k + 1
		}
		f := map[string]string{}
		w18I(f, "bound", bound)
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("In the truncated Euler product over the primes up to %d, each exponent ranges from 0 to %d. How many exponent choices are there?", bound, k), strconv.Itoa(ans), fmt.Sprintf("Each of the %d primes contributes %d choices, so the product is %d.", len(primes), k+1, ans), f)
	default:
		k := w18N(1, 20)
		ans := w18PrimesUpTo(100)[k]
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("If the Euler product is truncated after the first %d primes, what is the next prime that is omitted?", k), strconv.Itoa(ans), fmt.Sprintf("List the primes in order and take position %d counting from one.", k+1), f)
	}
}

func w18Continuation(o int) generator.Problem {
	switch o {
	case 1:
		m := w18N(1, 30)
		ans := 2 * m
		f := map[string]string{}
		w18I(f, "m", m)
		return w18P(o, fmt.Sprintf("The value zeta(2m) is a rational multiple of pi to some power. For m=%d, what is that exponent?", m), strconv.Itoa(ans), fmt.Sprintf("Even zeta values carry pi to the power 2m, which is %d here.", ans), f)
	case 2:
		m := w18N(1, 20)
		ans := 2 * m
		f := map[string]string{}
		w18I(f, "m", m)
		return w18P(o, fmt.Sprintf("The zeta value at minus %d is minus the Bernoulli number B at 2m divided by 2m, with m=%d. What is the divisor 2m, before any cancellation?", 2*m-1, m), strconv.Itoa(ans), fmt.Sprintf("Twice %d is %d; the reduced denominator divides this value.", m, ans), f)
	default:
		k := w18N(1, 40)
		ans := 2 * k
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("The k-th trivial zero of zeta is the negative even integer minus twice k. For k=%d, what is the absolute value of that zero?", k), strconv.Itoa(ans), fmt.Sprintf("The trivial zeros are minus 2, minus 4, minus 6 and so on, so the absolute value is %d.", ans), f)
	}
}

func w18FunctionalEquation(o int) generator.Problem {
	switch o {
	case 1:
		s := w18N(2, 30)
		f := map[string]string{}
		w18I(f, "s", s)
		return w18P(o, fmt.Sprintf("The completed zeta function contains the factor pi raised to minus s over two. For s=%d, write that exponent as minus n over two. What is n?", s), strconv.Itoa(s), fmt.Sprintf("The exponent is minus s/2, so n equals s, which is %d.", s), f)
	case 2:
		s := w18N(-20, 20)
		ans := 1 - s
		f := map[string]string{}
		w18I(f, "s", s)
		return w18P(o, fmt.Sprintf("The functional equation relates the zeta function at s to its value at 1 minus s. For s=%d, what is the reflected argument?", s), strconv.Itoa(ans), fmt.Sprintf("Subtract %d from one, giving %d.", s, ans), f)
	default:
		s := 2 * w18N(1, 30)
		ans := s / 2
		f := map[string]string{}
		w18I(f, "s", s)
		return w18P(o, fmt.Sprintf("The gamma factor in the completed zeta function is evaluated at s over two. For s=%d, what is that argument?", s), strconv.Itoa(ans), fmt.Sprintf("Halve %d to get %d.", s, ans), f)
	}
}

func w18ZeroCount(o int) generator.Problem {
	switch o {
	case 1:
		k := w18N(2, 20)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("At height T equal to 2 pi times %d, the leading factor T over 2 pi equals what integer?", k), strconv.Itoa(k), fmt.Sprintf("Divide T by 2 pi, which returns %d.", k), f)
	case 2:
		k := w18N(2, 20)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("How many trivial zeros of the zeta function lie on the real axis between minus twice %d and zero, counting the endpoint?", k), strconv.Itoa(k), fmt.Sprintf("The zeros are minus 2, minus 4, up to minus twice %d, so there are %d.", k, k), f)
	default:
		k := w18N(2, 20)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("The first %d nontrivial zeros above the real axis pair with their conjugates below it. How many conjugate pairs is that?", k), strconv.Itoa(k), fmt.Sprintf("Each zero above pairs with exactly one below, so there are %d pairs.", k), f)
	}
}

func w18ExplicitFormula(o int) generator.Problem {
	switch o {
	case 1:
		k := w18N(1, 20)
		ans := w18IPow(10, k)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("The explicit formula has main term x. For x equal to ten to the power %d, what is the main term?", k), strconv.Itoa(ans), fmt.Sprintf("Ten to the power %d is %d.", k, ans), f)
	case 2:
		k := w18N(1, 25)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("Pairing the first %d nontrivial zeros with their conjugates leaves how many terms in the zero sum?", k), strconv.Itoa(k), fmt.Sprintf("Conjugate pairs combine into one real term each, so the count stays %d.", k), f)
	default:
		k := 2 * w18N(1, 20)
		ans := w18IPow(10, k/2)
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("A zero on the critical line contributes a term of size x to the power one half. For x equal to ten to the power %d, what is that size?", k), strconv.Itoa(ans), fmt.Sprintf("The exponent is half of %d, so the value is ten to the power %d, which is %d.", k, k/2, ans), f)
	}
}

func w18Mertens(o int) generator.Problem {
	switch o {
	case 1:
		n := w18N(5, 40)
		sum := 0
		for k := 1; k <= n; k++ {
			sum += w18Mobius(k)
		}
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("The Mertens function M(n) sums mu(k) for k from one to n. What is M(%d)?", n), strconv.Itoa(sum), fmt.Sprintf("Add the Mobius values up to %d, cancelling squarefree integers of even parity against odd parity.", n), f)
	case 2:
		n := w18N(10, 60)
		ans := w18SquarefreeCount(n)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("How many integers from one to %d are squarefree?", n), strconv.Itoa(ans), fmt.Sprintf("An integer is squarefree exactly when its Mobius value is nonzero. Sieve up to %d and count them.", n), f)
	default:
		n := w18N(10, 120)
		ans := w18FloorSqrt(n)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("The elementary Mertens bound keeps the absolute value of M(n) below the square root of n. What is the floor of that bound for n=%d?", n), strconv.Itoa(ans), fmt.Sprintf("The integer square root of %d is %d.", n, ans), f)
	}
}

func w18Abel(o int) generator.Problem {
	switch o {
	case 1:
		n := w18N(5, 40)
		ans := n * (n + 1) / 2
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Abel summation with unit coefficients and weight f(n)=n reduces to the sum one plus two up to %d. What is that sum?", n), strconv.Itoa(ans), fmt.Sprintf("The sum is n(n+1)/2, which is %d.", ans), f)
	case 2:
		n := w18N(3, 20)
		ans := n * (n + 1) * (2*n + 1) / 6
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("With unit coefficients and weight f(n)=n squared, the summation identity gives the sum of squares up to %d. What is that sum?", n), strconv.Itoa(ans), fmt.Sprintf("The sum of squares is n(n+1)(2n+1)/6, which is %d.", ans), f)
	default:
		n := 2 * w18N(1, 25)
		ans := n * n / 2
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Applying Abel summation to the sum of n up to N=%d replaces the sum by the boundary term N squared minus the integral of t from zero to N. What is the result?", n), strconv.Itoa(ans), fmt.Sprintf("The boundary is N squared and the integral is N squared over two, leaving %d.", ans), f)
	}
}

func w18DirichletL(o int) generator.Problem {
	switch o {
	case 1:
		q := w18N(5, 40)
		ans := w18Totient(q)
		f := map[string]string{}
		w18I(f, "q", q)
		return w18P(o, fmt.Sprintf("How many Dirichlet characters are there modulo %d?", q), strconv.Itoa(ans), fmt.Sprintf("The number of characters modulo q is the totient of q, which is %d.", ans), f)
	case 2:
		q := 9 + 2*w18N(0, 90)
		omega := w18DistinctFactors(q)
		ans := 1 << omega
		f := map[string]string{}
		w18I(f, "q", q)
		return w18P(o, fmt.Sprintf("For odd q=%d, the number of real Dirichlet characters modulo q is two to the number of distinct prime factors. What is that number?", q), strconv.Itoa(ans), fmt.Sprintf("The modulus has %d distinct prime factors, so there are %d real characters.", omega, ans), f)
	default:
		q := w18N(10, 200)
		ans := w18DistinctFactors(q)
		f := map[string]string{}
		w18I(f, "q", q)
		return w18P(o, fmt.Sprintf("The principal character modulo %d differs from the constant character only at primes dividing the modulus. How many such primes are there?", q), strconv.Itoa(ans), fmt.Sprintf("Factor %d and count distinct primes, giving %d.", q, ans), f)
	}
}

func w18Bias(o int) generator.Problem {
	bound := 20 + 10*w18N(0, 58)
	primes := w18PrimesUpTo(bound)
	one, three := 0, 0
	for _, p := range primes {
		if p%4 == 1 {
			one++
		} else if p%4 == 3 {
			three++
		}
	}
	f := map[string]string{}
	w18I(f, "bound", bound)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("How many primes at most %d are congruent to three modulo four?", bound), strconv.Itoa(three), fmt.Sprintf("Sieve the primes up to %d and count those congruent to three modulo four, which gives %d.", bound, three), f)
	case 2:
		return w18P(o, fmt.Sprintf("How many primes at most %d are congruent to one modulo four?", bound), strconv.Itoa(one), fmt.Sprintf("Sieve the primes up to %d and count those congruent to one modulo four, which gives %d.", bound, one), f)
	default:
		ans := three - one
		return w18P(o, fmt.Sprintf("Among the primes at most %d, what is the count congruent to three modulo four minus the count congruent to one modulo four?", bound), strconv.Itoa(ans), fmt.Sprintf("The two counts are %d and %d, so the difference is %d.", three, one, ans), f)
	}
}

// ---------------------------------------------------------------------------
// Applications
// ---------------------------------------------------------------------------

func w18RSA(o int) generator.Problem {
	primes := []int{11, 13, 17, 19, 23, 29, 31, 37, 41, 43}
	var p, q int
	for {
		p = primes[rand.Intn(len(primes))]
		q = primes[rand.Intn(len(primes))]
		if p != q {
			break
		}
	}
	n := p * q
	phi := (p - 1) * (q - 1)
	e := 17
	for w18GCD(e, phi) != 1 {
		q = primes[rand.Intn(len(primes))]
		if q == p {
			continue
		}
		n = p * q
		phi = (p - 1) * (q - 1)
	}
	f := map[string]string{}
	w18I(f, "p", p)
	w18I(f, "q", q)
	w18I(f, "e", e)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("RSA uses primes p=%d and q=%d. What is the public modulus n?", p, q), strconv.Itoa(n), fmt.Sprintf("Multiply the primes: %d times %d is %d.", p, q, n), f)
	case 2:
		return w18P(o, fmt.Sprintf("RSA uses primes p=%d and q=%d. What is the totient of the modulus?", p, q), strconv.Itoa(phi), fmt.Sprintf("The totient is (p-1)(q-1), which is %d.", phi), f)
	default:
		d := w18ModInverse(e, phi)
		return w18P(o, fmt.Sprintf("RSA uses p=%d and q=%d with public exponent %d. What is the private exponent, the inverse of %d modulo the totient?", p, q, e, e), strconv.Itoa(d), fmt.Sprintf("The totient is %d, and the inverse of %d modulo it is %d.", phi, e, d), f)
	}
}

func w18ElGamal(o int) generator.Problem {
	type params struct{ p, g int }
	all := []params{{23, 5}, {29, 2}, {31, 3}, {37, 2}, {41, 6}, {43, 3}, {47, 5}}
	pr := all[rand.Intn(len(all))]
	x := w18N(2, pr.p-2)
	k := w18N(2, pr.p-2)
	h := w18PowMod(pr.g, x, pr.p)
	c1 := w18PowMod(pr.g, k, pr.p)
	shared := w18PowMod(h, k, pr.p)
	f := map[string]string{}
	w18I(f, "p", pr.p)
	w18I(f, "g", pr.g)
	w18I(f, "x", x)
	w18I(f, "k", k)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("ElGamal uses prime %d, generator %d and private exponent %d. What is the public value g to the x modulo p?", pr.p, pr.g, x), strconv.Itoa(h), fmt.Sprintf("Raise %d to the power %d modulo %d, giving %d.", pr.g, x, pr.p, h), f)
	case 2:
		return w18P(o, fmt.Sprintf("ElGamal uses prime %d and generator %d with ephemeral exponent %d. What is the first ciphertext component g to the k modulo p?", pr.p, pr.g, k), strconv.Itoa(c1), fmt.Sprintf("Raise %d to the power %d modulo %d, giving %d.", pr.g, k, pr.p, c1), f)
	default:
		return w18P(o, fmt.Sprintf("ElGamal uses prime %d, generator %d, private exponent %d and ephemeral exponent %d. What is the shared secret h to the k modulo p?", pr.p, pr.g, x, k), strconv.Itoa(shared), fmt.Sprintf("The public value is %d, and raising it to %d modulo %d gives %d.", h, k, pr.p, shared), f)
	}
}

func w18SecretSharing(o int) generator.Problem {
	p := []int{17, 19, 23, 29, 31}[rand.Intn(5)]
	switch o {
	case 1:
		n := w18N(3, 12)
		k := w18N(2, n)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("A %d-of-%d Shamir sharing scheme is set up. How many shares are needed to reconstruct the secret?", k, n), strconv.Itoa(k), fmt.Sprintf("The threshold is the number of points that determine the polynomial, which is %d.", k), f)
	case 2:
		s := w18N(1, p-1)
		a := w18N(1, p-1)
		f1 := (s + a) % p
		f2 := (s + 2*a) % p
		secret := ((2*f1-f2)%p + p) % p
		f := map[string]string{}
		w18I(f, "p", p)
		w18I(f, "f1", f1)
		w18I(f, "f2", f2)
		return w18P(o, fmt.Sprintf("A 2-of-3 sharing scheme works modulo %d and gives share one equal to %d and share two equal to %d. What is the secret?", p, f1, f2), strconv.Itoa(secret), fmt.Sprintf("Interpolate the line at zero: twice the first share minus the second, modulo %d, gives %d.", p, secret), f)
	default:
		s := w18N(1, p-1)
		a := w18N(1, p-1)
		b := w18N(1, p-1)
		y1 := (s + a + b) % p
		y2 := (s + 2*a + 4*b) % p
		y3 := (s + 3*a + 9*b) % p
		secret := w18LagrangeAtZero([]int{1, 2, 3}, []int{y1, y2, y3}, p)
		f := map[string]string{}
		w18I(f, "p", p)
		w18I(f, "y1", y1)
		w18I(f, "y2", y2)
		w18I(f, "y3", y3)
		return w18P(o, fmt.Sprintf("A 3-of-5 sharing scheme works modulo %d. Shares at x=1, x=2 and x=3 are %d, %d and %d. What is the secret?", p, y1, y2, y3), strconv.Itoa(secret), fmt.Sprintf("Lagrange interpolation at zero modulo %d gives %d.", p, secret), f)
	}
}

func w18LagrangeAtZero(xs, ys []int, p int) int {
	total := 0
	for i := range xs {
		num, den := 1, 1
		for j := range xs {
			if i == j {
				continue
			}
			num = num * (-xs[j] % p) % p
			den = den * ((xs[i] - xs[j]) % p) % p
		}
		term := ys[i] * num % p * w18ModInverse(((den%p)+p)%p, p) % p
		total = (total + term) % p
	}
	return ((total % p) + p) % p
}

func w18Checksum(o int) generator.Problem {
	digits := make([]int, 12)
	for i := range digits {
		digits[i] = rand.Intn(10)
	}
	ds := make([]string, len(digits))
	for i, d := range digits {
		ds[i] = strconv.Itoa(d)
	}
	f := map[string]string{}
	f["digits"] = strings.Join(ds, "")
	switch o {
	case 1:
		sum := 0
		for i := 0; i < 9; i++ {
			sum += (10 - i) * digits[i]
		}
		check := (11 - sum%11) % 11
		return w18P(o, fmt.Sprintf("For the ISBN-10 prefix %s, what is the check digit? Enter 10 for X.", strings.Join(ds[:9], "")), strconv.Itoa(check), fmt.Sprintf("Weight the nine digits from ten down to two, reduce modulo eleven and complete the sum to zero, giving %d.", check), f)
	case 2:
		sum := 0
		for i := 0; i < 12; i++ {
			w := 1
			if i%2 == 1 {
				w = 3
			}
			sum += w * digits[i]
		}
		check := (10 - sum%10) % 10
		return w18P(o, fmt.Sprintf("For the ISBN-13 prefix %s, what is the check digit?", strings.Join(ds[:12], "")), strconv.Itoa(check), fmt.Sprintf("Weight alternate digits by one and three, reduce modulo ten and complete the sum, giving %d.", check), f)
	default:
		check := w18LuhnCheck(digits[:9])
		return w18P(o, fmt.Sprintf("For the nine-digit prefix %s, what Luhn check digit completes a valid number?", strings.Join(ds[:9], "")), strconv.Itoa(check), fmt.Sprintf("Double every second digit from the right, subtract nine above nine, and complete the total modulo ten, giving %d.", check), f)
	}
}

func w18LuhnCheck(prefix []int) int {
	full := append(append([]int{}, prefix...), 0)
	sum := 0
	for i, d := range full {
		pos := len(full) - 1 - i
		if pos%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return (10 - sum%10) % 10
}

func w18ZeroKnowledge(o int) generator.Problem {
	switch o {
	case 1:
		k := w18N(5, 40)
		ans := 1 << k
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("A zero-knowledge protocol draws its challenge from k bits with k=%d. How many challenge values are possible?", k), strconv.Itoa(ans), fmt.Sprintf("There are two to the power %d challenges, which is %d.", k, ans), f)
	case 2:
		k := w18N(5, 40)
		ans := 1 << k
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("One round of a zero-knowledge proof has soundness error one over two to the k with k=%d. What is the denominator of that error?", k), strconv.Itoa(ans), fmt.Sprintf("The denominator is two to the power %d, which is %d.", k, ans), f)
	default:
		t := w18N(10, 40)
		f := map[string]string{}
		w18I(f, "t", t)
		return w18P(o, fmt.Sprintf("To reach a soundness error below one over two to the t with t=%d, how many independent repetitions are needed?", t), strconv.Itoa(t), fmt.Sprintf("Each repetition halves the error, so %d repetitions give the bound.", t), f)
	}
}

func w18ECC(o int) generator.Problem {
	small := w18PrimesUpTo(97)
	var curves []int
	for _, q := range small {
		if q >= 17 {
			curves = append(curves, q)
		}
	}
	p := curves[rand.Intn(len(curves))]
	a := w18N(1, p-1)
	b := w18N(1, p-1)
	x := w18N(0, p-1)
	f := map[string]string{}
	w18I(f, "p", p)
	w18I(f, "a", a)
	w18I(f, "b", b)
	w18I(f, "x", x)
	switch o {
	case 1:
		rhs := ((x*x%p)*x + a*x + b) % p
		rhs = ((rhs % p) + p) % p
		return w18P(o, fmt.Sprintf("An elliptic curve modulo %d has equation y squared equals x cubed plus %d x plus %d. What is the right-hand side at x=%d?", p, a, b, x), strconv.Itoa(rhs), fmt.Sprintf("Evaluate the cubic modulo %d, giving %d.", p, rhs), f)
	case 2:
		num := (3*x*x + a) % p
		num = ((num % p) + p) % p
		return w18P(o, fmt.Sprintf("For the curve y squared equals x cubed plus %d x plus %d modulo %d, what is the slope numerator 3x squared plus a at x=%d?", a, b, p, x), strconv.Itoa(num), fmt.Sprintf("Substitute x=%d and reduce modulo %d, giving %d.", x, p, num), f)
	default:
		ans := int(math.Floor(2 * math.Sqrt(float64(p))))
		return w18P(o, fmt.Sprintf("Hasse's bound says the number of points on an elliptic curve modulo %d differs from %d by at most 2 times the square root of %d. What is the floor of that half-width?", p, p+1, p), strconv.Itoa(ans), fmt.Sprintf("Twice the square root of %d is about %.2f, whose floor is %d.", p, 2*math.Sqrt(float64(p)), ans), f)
	}
}

func w18PrimalityCert(o int) generator.Problem {
	primes := w18PrimesUpTo(997)
	p := primes[w18N(25, len(primes)-1)]
	counts := w18FactorCounts(p - 1)
	omega := len(counts)
	Omega := w18Omega(p - 1)
	largest := 0
	for q := range counts {
		if q > largest {
			largest = q
		}
	}
	f := map[string]string{}
	w18I(f, "p", p)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("A Pratt certificate for the prime %d factors p minus one. How many distinct prime factors does %d have?", p, p-1), strconv.Itoa(omega), fmt.Sprintf("Factoring %d gives %d distinct primes.", p-1, omega), f)
	case 2:
		return w18P(o, fmt.Sprintf("A Pratt certificate for the prime %d needs the largest prime factor of p minus one. What is it for %d?", p, p), strconv.Itoa(largest), fmt.Sprintf("Factor %d and take the largest prime factor, which is %d.", p-1, largest), f)
	default:
		return w18P(o, fmt.Sprintf("Counting multiplicity, how many prime factors does %d have?", p-1), strconv.Itoa(Omega), fmt.Sprintf("Factor %d with multiplicity; there are %d factors in total.", p-1, Omega), f)
	}
}

func w18ReedSolomon(o int) generator.Problem {
	fieldPrimes := w18PrimesUpTo(97)
	var candidates []int
	for _, q := range fieldPrimes {
		if q >= 17 {
			candidates = append(candidates, q)
		}
	}
	q := candidates[rand.Intn(len(candidates))]
	n := q - 1
	k := w18N(2, n-2)
	f := map[string]string{}
	w18I(f, "q", q)
	w18I(f, "k", k)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("A Reed-Solomon code over the field with %d elements uses every nonzero field element as an evaluation point. What is its block length?", q), strconv.Itoa(n), fmt.Sprintf("The block length is q minus one, which is %d.", n), f)
	case 2:
		return w18P(o, fmt.Sprintf("A Reed-Solomon code over %d elements has block length %d and dimension %d. How many parity symbols does it carry?", q, n, k), strconv.Itoa(n-k), fmt.Sprintf("Subtract the dimension from the block length: %d.", n-k), f)
	default:
		return w18P(o, fmt.Sprintf("A Reed-Solomon code has block length %d and dimension %d. How many errors can it correct?", n, k), strconv.Itoa((n-k)/2), fmt.Sprintf("Half the redundancy, rounded down, is %d.", (n-k)/2), f)
	}
}

func w18Factorization(o int) generator.Problem {
	primes := []int{11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	p := primes[rand.Intn(len(primes))]
	q := primes[rand.Intn(len(primes))]
	for q == p {
		q = primes[rand.Intn(len(primes))]
	}
	n := p * q
	a := w18CeilSqrt(n)
	b := w18FloorSqrt(a*a - n)
	small := p
	if q < small {
		small = q
	}
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("Fermat factorisation of %d starts with the ceiling of its square root. What is that ceiling?", n), strconv.Itoa(a), fmt.Sprintf("The integer square root of %d rounded up is %d.", n, a), f)
	case 2:
		return w18P(o, fmt.Sprintf("For n=%d, the Fermat step writes n as a squared minus b squared. With a=%d, what is b?", n, a), strconv.Itoa(b), fmt.Sprintf("Compute a squared minus n and take its square root, giving %d.", b), f)
	default:
		return w18P(o, fmt.Sprintf("What is the smallest prime factor of %d?", n), strconv.Itoa(small), fmt.Sprintf("The factorisation is %d times %d, so the smallest prime factor is %d.", p, q, small), f)
	}
}

func w18UniversalHashing(o int) generator.Problem {
	large := w18PrimesUpTo(200)
	var primes []int
	for _, q := range large {
		if q >= 100 {
			primes = append(primes, q)
		}
	}
	p := primes[rand.Intn(len(primes))]
	m := w18N(7, 40)
	a := w18N(1, p-1)
	b := w18N(0, p-1)
	x := w18N(0, p-1)
	h := ((a*x + b) % p) % m
	f := map[string]string{}
	w18I(f, "p", p)
	w18I(f, "m", m)
	w18I(f, "a", a)
	w18I(f, "b", b)
	w18I(f, "x", x)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("A universal hash family uses h(x) = ((%d x + %d) mod %d) mod %d. What is h(%d)?", a, b, p, m, x), strconv.Itoa(h), fmt.Sprintf("Reduce %d x plus %d modulo %d, then reduce the result modulo %d, giving %d.", a, b, p, m, h), f)
	case 2:
		ans := p * (p - 1)
		return w18P(o, fmt.Sprintf("The family h(x) = ((a x + b) mod p) mod m uses a prime p=%d, a nonzero a and any b. How many functions are in the family?", p), strconv.Itoa(ans), fmt.Sprintf("There are p minus one choices for a and p choices for b, giving %d.", ans), f)
	default:
		return w18P(o, fmt.Sprintf("In a universal family with table size m=%d, what is the upper bound on the collision probability of two distinct keys, written with denominator m?", m), strconv.Itoa(m), fmt.Sprintf("The bound is one over m, so the denominator is %d.", m), f)
	}
}

// ---------------------------------------------------------------------------
// Diophantine equations
// ---------------------------------------------------------------------------

func w18ThreeSquares(o int) generator.Problem {
	n := w18N(10, 200)
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("What is %d modulo eight?", n), strconv.Itoa(n%8), fmt.Sprintf("Divide %d by eight; the remainder is %d.", n, n%8), f)
	case 2:
		m := n
		for m%4 == 0 {
			m /= 4
		}
		ans := 1
		if m%8 == 7 {
			ans = 0
		}
		return w18P(o, fmt.Sprintf("Legendre's theorem says %d is a sum of three squares unless its odd part is seven modulo eight. Is it representable? Answer 1 for yes and 0 for no.", n), strconv.Itoa(ans), fmt.Sprintf("Divide out powers of four, leaving %d, whose residue modulo eight is %d.", m, m%8), f)
	default:
		m := n
		for m%4 == 0 {
			m /= 4
		}
		return w18P(o, fmt.Sprintf("After removing every factor of four from %d, what is the remaining residue modulo eight?", n), strconv.Itoa(m%8), fmt.Sprintf("The odd part is %d, and its residue modulo eight is %d.", m, m%8), f)
	}
}

func w18Waring(o int) generator.Problem {
	k := w18N(2, 20)
	g := w18IPow(2, k) + w18IPow(3, k)/w18IPow(2, k) - 2
	f := map[string]string{}
	w18I(f, "k", k)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("Waring's function g(k) is two to the k plus the floor of three halves to the k minus two. What is g(%d)?", k), strconv.Itoa(g), fmt.Sprintf("Compute two to the %d plus the floor of three halves to the %d minus two, giving %d.", k, k, g), f)
	case 2:
		ans := w18IPow(2, k)
		return w18P(o, fmt.Sprintf("In the formula for g(%d), what is two to the k?", k), strconv.Itoa(ans), fmt.Sprintf("Two to the power %d is %d.", k, ans), f)
	default:
		ans := w18IPow(3, k) / w18IPow(2, k)
		return w18P(o, fmt.Sprintf("In the formula for g(%d), what is the floor of three halves to the k?", k), strconv.Itoa(ans), fmt.Sprintf("Three to the %d divided by two to the %d, rounded down, is %d.", k, k, ans), f)
	}
}

func w18Markov(o int) generator.Problem {
	triples := [][3]int{{1, 1, 1}, {1, 1, 2}, {1, 2, 5}, {1, 5, 13}, {2, 5, 29}, {1, 13, 34}, {1, 34, 89}, {5, 13, 194}, {2, 29, 169}, {5, 29, 433},
		{1, 89, 233}, {2, 169, 985}, {5, 433, 6466}, {13, 34, 1325}, {13, 194, 7561}, {1, 233, 610}}
	t := triples[rand.Intn(len(triples))]
	x, y, z := t[0], t[1], t[2]
	f := map[string]string{}
	w18I(f, "x", x)
	w18I(f, "y", y)
	w18I(f, "z", z)
	switch o {
	case 1:
		ans := x*x + y*y + z*z
		return w18P(o, fmt.Sprintf("Check the Markov triple (%d, %d, %d): what is the sum of the three squares?", x, y, z), strconv.Itoa(ans), fmt.Sprintf("Square each entry and add, giving %d.", ans), f)
	case 2:
		ans := 3 * x * y
		return w18P(o, fmt.Sprintf("For the Markov triple (%d, %d, %d), what is three times the product of the first two entries, the quantity the third entry must divide?", x, y, z), strconv.Itoa(ans), fmt.Sprintf("Three times %d times %d is %d.", x, y, ans), f)
	default:
		ans := 3*x*y - z
		return w18P(o, fmt.Sprintf("Starting from the Markov triple (%d, %d, %d), Vieta jumping fixes the first two entries and replaces the third by the other root. What is that new entry?", x, y, z), strconv.Itoa(ans), fmt.Sprintf("The two roots sum to 3xy = %d, so the other root is %d.", 3*x*y, ans), f)
	}
}

func w18Mordell(o int) generator.Problem {
	k := w18N(1, 20)
	x := w18N(0, 30)
	f := map[string]string{}
	w18I(f, "k", k)
	w18I(f, "x", x)
	switch o {
	case 1:
		ans := x*x*x + k
		return w18P(o, fmt.Sprintf("For the Mordell equation y squared equals x cubed plus %d, what is the right-hand side at x=%d?", k, x), strconv.Itoa(ans), fmt.Sprintf("%d cubed plus %d is %d.", x, k, ans), f)
	case 2:
		ans := 432 * k * k
		return w18P(o, fmt.Sprintf("The discriminant of the elliptic curve y squared equals x cubed plus %d has absolute value 432 k squared. What is it?", k), strconv.Itoa(ans), fmt.Sprintf("Multiply 432 by %d squared, giving %d.", k, ans), f)
	default:
		count := 0
		for xx := 0; xx <= 10; xx++ {
			v := xx*xx*xx + k
			r := w18FloorSqrt(v)
			if r*r == v {
				count++
			}
		}
		return w18P(o, fmt.Sprintf("For k=%d and x from zero to ten, how many values of x make x cubed plus k a perfect square?", k), strconv.Itoa(count), fmt.Sprintf("Testing the eleven values and checking each square gives %d solutions.", count), f)
	}
}

func w18Catalan(o int) generator.Problem {
	a := w18N(2, 6)
	b := w18N(2, 5)
	f := map[string]string{}
	w18I(f, "a", a)
	w18I(f, "b", b)
	switch o {
	case 1:
		ans := w18IPow(a, 2)
		return w18P(o, fmt.Sprintf("For the Catalan-type expression a squared minus b cubed with a=%d and b=%d, what is a squared?", a, b), strconv.Itoa(ans), fmt.Sprintf("%d squared is %d.", a, ans), f)
	case 2:
		ans := w18IPow(b, 3)
		return w18P(o, fmt.Sprintf("For a=%d and b=%d, what is b cubed?", a, b), strconv.Itoa(ans), fmt.Sprintf("%d cubed is %d.", b, ans), f)
	default:
		ans := w18IPow(a, 2) - w18IPow(b, 3)
		return w18P(o, fmt.Sprintf("For a=%d and b=%d, what is a squared minus b cubed, the quantity that Catalan's conjecture controls?", a, b), strconv.Itoa(ans), fmt.Sprintf("%d minus %d is %d.", w18IPow(a, 2), w18IPow(b, 3), ans), f)
	}
}

func w18ABC(o int) generator.Problem {
	triples := [][3]int{{1, 8, 9}, {1, 15, 16}, {2, 25, 27}, {3, 125, 128}, {1, 48, 49}, {5, 27, 32}, {1, 80, 81}, {4, 121, 125}, {7, 9, 16}, {9, 16, 25},
		{1, 3, 4}, {1, 7, 8}, {2, 3, 5}, {1, 24, 25}, {1, 63, 64}, {3, 5, 8}, {8, 9, 17}}
	t := triples[rand.Intn(len(triples))]
	a, b, c := t[0], t[1], t[2]
	counts := map[int]int{}
	for _, v := range []int{a, b, c} {
		for p, e := range w18FactorCounts(v) {
			counts[p] += e
		}
	}
	rad := 1
	for p := range counts {
		rad *= p
	}
	f := map[string]string{}
	w18I(f, "a", a)
	w18I(f, "b", b)
	w18I(f, "c", c)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("For the triple %d plus %d equals %d, what is the radical of the product abc?", a, b, c), strconv.Itoa(rad), fmt.Sprintf("Multiply the distinct primes dividing %d, %d and %d, giving %d.", a, b, c, rad), f)
	case 2:
		ans := len(counts)
		return w18P(o, fmt.Sprintf("For the triple %d plus %d equals %d, how many distinct primes divide the product abc?", a, b, c), strconv.Itoa(ans), fmt.Sprintf("Factoring the three parts together gives %d distinct primes.", ans), f)
	default:
		ans := c - rad
		return w18P(o, fmt.Sprintf("For the triple %d plus %d equals %d, what is c minus the radical of abc?", a, b, c), strconv.Itoa(ans), fmt.Sprintf("The radical is %d, so c minus it is %d.", rad, ans), f)
	}
}

func w18Taxicab(o int) generator.Problem {
	pairs := [][2]int{{1, 12}, {9, 10}, {1, 6}, {2, 17}, {3, 10}, {4, 11},
		{2, 3}, {5, 7}, {8, 9}, {6, 10}, {7, 12}, {4, 13}, {5, 14}, {6, 11},
		{2, 5}, {3, 4}, {7, 8}, {1, 13}, {5, 12}, {9, 11}}
	pr := pairs[rand.Intn(len(pairs))]
	switch o {
	case 1:
		ans := pr[0]*pr[0]*pr[0] + pr[1]*pr[1]*pr[1]
		f := map[string]string{}
		w18I(f, "a", pr[0])
		w18I(f, "b", pr[1])
		return w18P(o, fmt.Sprintf("What is %d cubed plus %d cubed?", pr[0], pr[1]), strconv.Itoa(ans), fmt.Sprintf("The two cubes add to %d.", ans), f)
	case 2:
		ans := pr[1]*pr[1]*pr[1] - pr[0]*pr[0]*pr[0]
		f := map[string]string{}
		w18I(f, "a", pr[0])
		w18I(f, "b", pr[1])
		return w18P(o, fmt.Sprintf("In the representation %d cubed plus %d cubed, what is the difference between the larger cube and the smaller one?", pr[0], pr[1]), strconv.Itoa(ans), fmt.Sprintf("The cube of %d minus the cube of %d is %d.", pr[1], pr[0], ans), f)
	default:
		ans := pr[1] - pr[0] - 1
		f := map[string]string{}
		w18I(f, "a", pr[0])
		w18I(f, "b", pr[1])
		return w18P(o, fmt.Sprintf("How many perfect cubes lie strictly between %d cubed and %d cubed?", pr[0], pr[1]), strconv.Itoa(ans), fmt.Sprintf("The cubes are those of the integers from %d to %d exclusive, which is %d values.", pr[0]+1, pr[1], ans), f)
	}
}

func w18Frobenius(o int) generator.Problem {
	a := w18N(2, 15)
	b := w18N(2, 15)
	for w18GCD(a, b) != 1 || a == b {
		b = w18N(2, 15)
	}
	f := map[string]string{}
	w18I(f, "a", a)
	w18I(f, "b", b)
	switch o {
	case 1:
		ans := a*b - a - b
		return w18P(o, fmt.Sprintf("With coprime coin values %d and %d, what is the largest amount that cannot be formed?", a, b), strconv.Itoa(ans), fmt.Sprintf("The Frobenius number is ab minus a minus b, which is %d.", ans), f)
	case 2:
		ans := (a - 1) * (b - 1) / 2
		return w18P(o, fmt.Sprintf("With coprime coin values %d and %d, how many positive amounts cannot be formed?", a, b), strconv.Itoa(ans), fmt.Sprintf("The gap count is (a-1)(b-1)/2, which is %d.", ans), f)
	default:
		ans := a*b - a - b + 1
		return w18P(o, fmt.Sprintf("With coprime coin values %d and %d, what is the smallest amount above the largest gap that can be formed?", a, b), strconv.Itoa(ans), fmt.Sprintf("The conductor is one more than the Frobenius number, giving %d.", ans), f)
	}
}

func w18FermatLast(o int) generator.Problem {
	switch o {
	case 1:
		n := w18N(2, 30)
		ans := 0
		if n == 2 {
			ans = 1
		}
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Does x to the %d plus y to the %d equal z to the %d have a solution in positive integers for exponent %d? Answer 1 for yes and 0 for no.", n, n, n, n), strconv.Itoa(ans), fmt.Sprintf("Only exponent two has solutions, so the answer is %d.", ans), f)
	case 2:
		n := w18N(3, 40)
		ans := n - 2
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Fermat's last theorem excludes every exponent from three to %d. How many exponents does that rule out?", n), strconv.Itoa(ans), fmt.Sprintf("The exponents three through %d number %d.", n, ans), f)
	default:
		z := w18N(10, 40)
		count := 0
		for x := 1; x <= z; x++ {
			for y := x + 1; y <= z; y++ {
				for zz := y + 1; zz <= z; zz++ {
					if x*x+y*y == zz*zz {
						count++
					}
				}
			}
		}
		f := map[string]string{}
		w18I(f, "z", z)
		return w18P(o, fmt.Sprintf("How many Pythagorean triples with 1 <= x < y < z and z at most %d are there?", z), strconv.Itoa(count), fmt.Sprintf("Enumerating the triples up to hypotenuse %d gives %d.", z, count), f)
	}
}

func w18FourCubes(o int) generator.Problem {
	switch o {
	case 1:
		n := w18N(1, 100)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("A cube is congruent to zero, one or minus one modulo nine. What is %d modulo nine?", n), strconv.Itoa(n%9), fmt.Sprintf("The remainder of %d modulo nine is %d.", n, n%9), f)
	case 2:
		n := w18N(1, 100)
		ans := 0
		if n%9 == 4 || n%9 == 5 {
			ans = 1
		}
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Is %d excluded from being a sum of four cubes by the modulo-nine obstruction? Answer 1 for yes and 0 for no.", n), strconv.Itoa(ans), fmt.Sprintf("The residue modulo nine is %d, and only four and five are excluded.", n%9), f)
	default:
		x := w18N(-3, 3)
		y := w18N(-3, 3)
		z := w18N(-3, 3)
		w := w18N(-3, 3)
		ans := x*x*x + y*y*y + z*z*z + w*w*w
		f := map[string]string{}
		w18I(f, "x", x)
		w18I(f, "y", y)
		w18I(f, "z", z)
		w18I(f, "w", w)
		return w18P(o, fmt.Sprintf("What is %d cubed plus %d cubed plus %d cubed plus %d cubed?", x, y, z, w), strconv.Itoa(ans), fmt.Sprintf("The four signed cubes sum to %d.", ans), f)
	}
}
