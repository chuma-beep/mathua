package numtheory

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ntBridgeGen serves the number theory concepts added by the 2026 content
// expansion. Each concept has three knowledge points and every point draws its
// own parameters, so the stem changes with the ordinal and with the seed.
type ntBridgeGen struct{ kind string }

func (g *ntBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "fta":
		return ftaProblem(o)
	case "arithmetic_functions":
		return arithmeticFunctionsProblem(o)
	case "gaussian_integers":
		return gaussianIntegersProblem(o)
	case "four_squares":
		return fourSquaresProblem(o)
	case "hensel":
		return henselProblem(o)
	case "farey":
		return fareyProblem(o)
	case "key_exchange":
		return keyExchangeProblem(o)
	case "quadratic_forms":
		return quadraticFormsProblem(o)
	case "valuation":
		return valuationProblem(o)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func ntBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func ntBridgeModPow(base, exp, mod int) int {
	if mod == 1 {
		return 0
	}
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

func ntBridgeGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func ntBridgePhi(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if ntBridgeGCD(k, n) == 1 {
			count++
		}
	}
	return count
}

func ntBridgePrimeFactors(n int) []int {
	var factors []int
	for d := 2; d*d <= n; d++ {
		for n%d == 0 {
			factors = append(factors, d)
			n /= d
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

func ntBridgeOmega(n int) int { return len(ntBridgePrimeFactors(n)) }

func ntBridgeDistinct(n int) int {
	count := 0
	last := 0
	for _, p := range ntBridgePrimeFactors(n) {
		if p != last {
			count++
			last = p
		}
	}
	return count
}

func ntBridgeSigma(n int) int {
	sum := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			sum += d
		}
	}
	return sum
}

func ntBridgeSigma2(n int) int {
	sum := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			sum += d * d
		}
	}
	return sum
}

func ntBridgeTau(n int) int {
	count := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			count++
		}
	}
	return count
}

func ntBridgeMinSquares(n int) int {
	limit := 1
	for (limit+1)*(limit+1) <= n {
		limit++
	}
	for a := 0; a <= limit; a++ {
		if a*a == n {
			return 1
		}
	}
	for a := 0; a <= limit; a++ {
		for b := 0; b <= limit; b++ {
			if a*a+b*b == n {
				return 2
			}
		}
	}
	for a := 0; a <= limit; a++ {
		for b := 0; b <= limit; b++ {
			for c := 0; c <= limit; c++ {
				if a*a+b*b+c*c == n {
					return 3
				}
			}
		}
	}
	return 4
}

func ntBridgeR2(n int) int {
	count := 0
	limit := 1
	for limit*limit < n {
		limit++
	}
	for x := -limit; x <= limit; x++ {
		for y := -limit; y <= limit; y++ {
			if x*x+y*y == n {
				count++
			}
		}
	}
	return count
}

func ntBridgeR4(n int) int {
	sum := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 && d%4 != 0 {
			sum += d
		}
	}
	return 8 * sum
}

func ntBridgeInvMod(a, p int) int {
	a %= p
	if a < 0 {
		a += p
	}
	for x := 1; x < p; x++ {
		if a*x%p == 1 {
			return x
		}
	}
	return 1
}

func ntBridgeLegendre(a, p int) int {
	if a%p == 0 {
		return 0
	}
	r := ntBridgeModPow(a%p, (p-1)/2, p)
	if r == 1 {
		return 1
	}
	return -1
}

func ftaProblem(o int) generator.Problem {
	n := 12 + rand.Intn(189)
	switch o {
	case 1:
		k := ntBridgeOmega(n)
		return ntBridgeProblem(o, fmt.Sprintf("How many prime factors does %d have, counted with multiplicity? (enter a number)", n), strconv.Itoa(k), fmt.Sprintf("%d factors as %v, so there are %d prime factors with multiplicity.", n, ntBridgePrimeFactors(n), k))
	case 2:
		k := ntBridgeDistinct(n)
		return ntBridgeProblem(o, fmt.Sprintf("How many distinct prime factors does %d have? (enter a number)", n), strconv.Itoa(k), fmt.Sprintf("%d has the distinct prime factors %v.", n, ntBridgePrimeFactors(n)))
	default:
		f := ntBridgePrimeFactors(n)
		largest := f[len(f)-1]
		return ntBridgeProblem(o, fmt.Sprintf("What is the largest prime factor of %d? (enter a number)", n), strconv.Itoa(largest), fmt.Sprintf("%d factors as %v, so the largest prime factor is %d.", n, f, largest))
	}
}

func arithmeticFunctionsProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 2 + rand.Intn(29)
		v := ntBridgeSigma2(n)
		return ntBridgeProblem(o, fmt.Sprintf("What is the sum of the squares of the positive divisors of %d? (enter a number)", n), strconv.Itoa(v), fmt.Sprintf("The divisors of %d are summed after squaring, giving %d.", n, v))
	case 2:
		n := 2 + rand.Intn(39)
		v := ntBridgeSigma(n)
		return ntBridgeProblem(o, fmt.Sprintf("The function sigma sums the positive divisors of an integer. What is sigma(%d)? (enter a number)", n), strconv.Itoa(v), fmt.Sprintf("The positive divisors of %d add to %d, which is the Dirichlet convolution of the identity with the constant function one.", n, v))
	default:
		primes := []int{2, 3, 5, 7}
		pi := rand.Intn(len(primes))
		qi := rand.Intn(len(primes) - 1)
		if qi >= pi {
			qi++
		}
		p, q := primes[pi], primes[qi]
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(3)
		v := (a + 1) * (b + 1)
		return ntBridgeProblem(o, fmt.Sprintf("How many positive divisors does %d^%d * %d^%d have? (enter a number)", p, a, q, b), strconv.Itoa(v), fmt.Sprintf("A divisor chooses an exponent from 0 to %d for %d and from 0 to %d for %d, so there are %d choices.", a, p, b, q, v))
	}
}

func gaussianIntegersProblem(o int) generator.Problem {
	a := rand.Intn(11) - 5
	b := rand.Intn(11) - 5
	if a == 0 && b == 0 {
		a = 1
	}
	c := rand.Intn(11) - 5
	d := rand.Intn(11) - 5
	switch o {
	case 1:
		v := a*a + b*b
		return ntBridgeProblem(o, fmt.Sprintf("What is the norm of the Gaussian integer %d + %d i? (enter a number)", a, b), strconv.Itoa(v), fmt.Sprintf("The norm is a squared plus b squared: %d + %d = %d.", a*a, b*b, v))
	case 2:
		v := a*c - b*d
		return ntBridgeProblem(o, fmt.Sprintf("What is the real part of (%d + %d i)(%d + %d i)? (enter a number)", a, b, c, d), strconv.Itoa(v), fmt.Sprintf("Multiply out: the real part is a c minus b d, which is %d.", v))
	default:
		v := a*d + b*c
		return ntBridgeProblem(o, fmt.Sprintf("What is the imaginary part of (%d + %d i)(%d + %d i)? (enter a number)", a, b, c, d), strconv.Itoa(v), fmt.Sprintf("Multiply out: the imaginary part is a d plus b c, which is %d.", v))
	}
}

func fourSquaresProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 1 + rand.Intn(100)
		k := ntBridgeMinSquares(n)
		return ntBridgeProblem(o, fmt.Sprintf("What is the smallest number of nonzero squares needed to write %d as a sum of squares? (enter a number)", n), strconv.Itoa(k), fmt.Sprintf("By the three-square and two-square criteria, %d needs %d squares and no fewer.", n, k))
	case 2:
		n := 1 + rand.Intn(50)
		v := ntBridgeR2(n)
		return ntBridgeProblem(o, fmt.Sprintf("How many ordered pairs of integers (x, y) satisfy x^2 + y^2 = %d? (enter a number)", n), strconv.Itoa(v), fmt.Sprintf("Counting signs and order, there are %d such pairs for %d.", v, n))
	default:
		n := 1 + rand.Intn(15)
		v := ntBridgeR4(n)
		return ntBridgeProblem(o, fmt.Sprintf("How many ordered quadruples of integers (a, b, c, d) satisfy a^2 + b^2 + c^2 + d^2 = %d? (enter a number)", n), strconv.Itoa(v), fmt.Sprintf("The four-square representation count is eight times the sum of the divisors of %d that are not divisible by four, giving %d.", n, v))
	}
}

func henselProblem(o int) generator.Problem {
	primes := []int{3, 5, 7, 11}
	p := primes[rand.Intn(len(primes))]
	r := 1 + rand.Intn(p-1)
	a := r * r % p
	f := r*r - a
	correction := (f / p) * ntBridgeInvMod(2*r, p) % p
	t := (p - correction) % p
	lift := (r + p*t) % (p * p)
	switch o {
	case 1:
		return ntBridgeProblem(o, fmt.Sprintf("The congruence x^2 = %d (mod %d) has the root %d. Use Hensel's lemma to find the root modulo %d that reduces to %d. (enter a number)", a, p, r, p*p, r), strconv.Itoa(lift), fmt.Sprintf("The lift is r - (f(r)/p) times the inverse of 2r modulo p, which gives %d modulo %d.", lift, p*p))
	case 2:
		roots := 0
		switch ntBridgeLegendre(a, p) {
		case 0:
			roots = 1
		case 1:
			roots = 2
		default:
			roots = 0
		}
		return ntBridgeProblem(o, fmt.Sprintf("How many roots does x^2 = %d (mod %d) have? (enter a number)", a, p), strconv.Itoa(roots), fmt.Sprintf("The Legendre symbol decides the count: a zero residue has one root, a quadratic residue has two and a nonresidue has none. Here the count is %d.", roots))
	default:
		return ntBridgeProblem(o, fmt.Sprintf("In the Hensel lift of the root %d of x^2 = %d (mod %d), the lift has the form %d + %d t. What is the correction t modulo %d? (enter a number)", r, a, p, r, p, p), strconv.Itoa(t), fmt.Sprintf("The correction solves (f(r)/p) + 2r t = 0 modulo p, which gives t = %d.", t))
	}
}

func fareyProblem(o int) generator.Problem {
	n := 2 + rand.Intn(11)
	switch o {
	case 1:
		length := 1
		for m := 1; m <= n; m++ {
			length += ntBridgePhi(m)
		}
		return ntBridgeProblem(o, fmt.Sprintf("How many fractions are in the Farey sequence of order %d? (enter a number)", n), strconv.Itoa(length), fmt.Sprintf("The length is one plus the sum of phi(m) for m from one to %d, which is %d.", n, length))
	case 2:
		m := 2 + rand.Intn(29)
		v := ntBridgePhi(m)
		return ntBridgeProblem(o, fmt.Sprintf("How many fractions in lowest terms with denominator exactly %d lie between zero and one? (enter a number)", m), strconv.Itoa(v), fmt.Sprintf("Exactly phi(%d) = %d numerators are coprime to %d.", m, v, m))
	default:
		sum := 0
		for m := 1; m <= n; m++ {
			sum += ntBridgePhi(m)
		}
		return ntBridgeProblem(o, fmt.Sprintf("What is the sum of phi(m) for m from 1 to %d? (enter a number)", n), strconv.Itoa(sum), fmt.Sprintf("Adding the totients from one to %d gives %d, one less than the Farey length.", n, sum))
	}
}

func keyExchangeProblem(o int) generator.Problem {
	primes := []int{17, 19, 23, 29, 31, 37}
	p := primes[rand.Intn(len(primes))]
	gens := []int{2, 3, 5, 7}
	g := gens[rand.Intn(len(gens))]
	a := 2 + rand.Intn(p-3)
	b := 2 + rand.Intn(p-3)
	c := 2 + rand.Intn(p-3)
	switch o {
	case 1:
		v := ntBridgeModPow(g, a*b, p)
		return ntBridgeProblem(o, fmt.Sprintf("Diffie-Hellman uses prime %d and generator %d. Alice has private exponent %d and Bob has private exponent %d. What is their shared secret modulo %d? (enter a number)", p, g, a, b, p), strconv.Itoa(v), fmt.Sprintf("The shared secret is g raised to a b modulo p, which is %d.", v))
	case 2:
		v := ntBridgeModPow(g, a, p)
		return ntBridgeProblem(o, fmt.Sprintf("With prime %d and generator %d, Alice's private exponent is %d. What public value does she send? (enter a number)", p, g, a), strconv.Itoa(v), fmt.Sprintf("Alice sends g raised to a modulo p, which is %d.", v))
	default:
		v := ntBridgeModPow(g, a*b*c, p)
		return ntBridgeProblem(o, fmt.Sprintf("Three parties use prime %d and generator %d with private exponents %d, %d and %d. What is their common secret modulo %d? (enter a number)", p, g, a, b, c, p), strconv.Itoa(v), fmt.Sprintf("The common secret is g raised to a b c modulo p, which is %d.", v))
	}
}

func quadraticFormsProblem(o int) generator.Problem {
	a := rand.Intn(7) - 3
	if a == 0 {
		a = 1
	}
	c := rand.Intn(7) - 3
	if c == 0 {
		c = -1
	}
	b := rand.Intn(9) - 4
	disc := b*b - 4*a*c
	if disc == 0 {
		b++
		disc = b*b - 4*a*c
	}
	switch o {
	case 1:
		return ntBridgeProblem(o, fmt.Sprintf("What is the discriminant of the binary quadratic form %d x^2 + %d x y + %d y^2? (enter a number)", a, b, c), strconv.Itoa(disc), fmt.Sprintf("The discriminant is b squared minus four a c: %d - %d = %d.", b*b, 4*a*c, disc))
	case 2:
		kind := 0
		if disc < 0 {
			if a > 0 {
				kind = 1
			} else {
				kind = -1
			}
		}
		return ntBridgeProblem(o, fmt.Sprintf("Classify the form %d x^2 + %d x y + %d y^2. Enter 1 for positive definite, -1 for negative definite and 0 for indefinite. (enter a number)", a, b, c), strconv.Itoa(kind), fmt.Sprintf("The discriminant is %d, so the form is %s.", disc, map[int]string{1: "positive definite", -1: "negative definite", 0: "indefinite"}[kind]))
	default:
		x := rand.Intn(7) - 3
		y := rand.Intn(7) - 3
		v := a*x*x + b*x*y + c*y*y
		return ntBridgeProblem(o, fmt.Sprintf("What is the value of %d x^2 + %d x y + %d y^2 at x = %d and y = %d? (enter a number)", a, b, c, x, y), strconv.Itoa(v), fmt.Sprintf("Substitute and add: %d.", v))
	}
}

func valuationProblem(o int) generator.Problem {
	primes := []int{2, 3, 5, 7}
	p := primes[rand.Intn(len(primes))]
	switch o {
	case 1:
		e := 1 + rand.Intn(5)
		m := 1 + rand.Intn(20)
		for m%p == 0 {
			m++
		}
		n := 1
		for i := 0; i < e; i++ {
			n *= p
		}
		n *= m
		return ntBridgeProblem(o, fmt.Sprintf("What is the exponent of %d in the prime factorisation of %d? (enter a number)", p, n), strconv.Itoa(e), fmt.Sprintf("%d divides %d exactly %d times and %d does not divide the remaining factor %d.", p, n, e, p, m))
	case 2:
		n := 5 + rand.Intn(96)
		sum := 0
		for q := n / p; q > 0; q /= p {
			sum += q
		}
		return ntBridgeProblem(o, fmt.Sprintf("What is the exponent of %d in the prime factorisation of %d! ? (enter a number)", p, n), strconv.Itoa(sum), fmt.Sprintf("Legendre's formula adds the floors of n over p, p squared and higher powers, giving %d.", sum))
	default:
		e := 1 + rand.Intn(5)
		m := 1 + rand.Intn(20)
		for m%p == 0 {
			m++
		}
		n := 1
		for i := 0; i < e; i++ {
			n *= p
		}
		n *= m
		k := 2 + rand.Intn(4)
		v := e * k
		return ntBridgeProblem(o, fmt.Sprintf("If the exponent of %d in %d is %d, what is the exponent of %d in %d^%d? (enter a number)", p, n, e, p, n, k), strconv.Itoa(v), fmt.Sprintf("Valuations add under powers, so the exponent is %d times %d, which is %d.", k, e, v))
	}
}
