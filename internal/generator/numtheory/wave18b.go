package numtheory

// Wave 18 number theory, second batch: computational and additive topics.

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave18Rest(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"nt.comp.trial_division":  "trial_division",
		"nt.comp.pollard_rho":     "pollard_rho",
		"nt.comp.quadratic_sieve": "quadratic_sieve",
		"nt.comp.aks":             "aks",
		"nt.comp.bsgs":            "bsgs",
		"nt.comp.pohlig_hellman":  "pohlig_hellman",
		"nt.comp.lll":             "lll",
		"nt.comp.pratt":           "pratt",
		"nt.comp.cf_factor":       "cf_factor",
		"nt.comp.nfs":             "nfs",
		"nt.add.partitions":       "partitions",
		"nt.add.vinogradov":       "vinogradov",
		"nt.add.circle_method":    "circle_method",
		"nt.add.schnirelmann":     "schnirelmann",
		"nt.add.sumsets":          "sumsets",
		"nt.add.cauchy_davenport": "cauchy_davenport",
		"nt.add.freiman":          "freiman",
		"nt.add.additive_bases":   "additive_bases",
		"nt.add.goldbach_count":   "goldbach_count",
		"nt.add.additive_energy":  "additive_energy",
	} {
		reg.Register(id, &w18Gen{kind: kind})
	}
}

func w18RestGenerate(kind string, o int) generator.Problem {
	switch kind {
	case "trial_division":
		return w18TrialDivision(o)
	case "pollard_rho":
		return w18PollardRho(o)
	case "quadratic_sieve":
		return w18QuadraticSieve(o)
	case "aks":
		return w18AKS(o)
	case "bsgs":
		return w18BSGS(o)
	case "pohlig_hellman":
		return w18PohligHellman(o)
	case "lll":
		return w18LLL(o)
	case "pratt":
		return w18Pratt(o)
	case "cf_factor":
		return w18CFFactor(o)
	case "nfs":
		return w18NFS(o)
	case "partitions":
		return w18Partitions(o)
	case "vinogradov":
		return w18Vinogradov(o)
	case "circle_method":
		return w18CircleMethod(o)
	case "schnirelmann":
		return w18Schnirelmann(o)
	case "sumsets":
		return w18Sumsets(o)
	case "cauchy_davenport":
		return w18CauchyDavenport(o)
	case "freiman":
		return w18Freiman(o)
	case "additive_bases":
		return w18AdditiveBases(o)
	case "goldbach_count":
		return w18GoldbachCount(o)
	case "additive_energy":
		return w18AdditiveEnergy(o)
	}
	return generator.Problem{
		Question:    "What is 2 + 2? (enter a number)",
		Answer:      "4",
		Explanation: "2 + 2 = 4.",
	}
}

func w18SmallestFactor(n int) int {
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return d
		}
	}
	return n
}

func w18LargestFactor(n int) int {
	counts := w18FactorCounts(n)
	largest := 1
	for p := range counts {
		if p > largest {
			largest = p
		}
	}
	return largest
}

func w18Composite() int {
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43}
	for {
		n := w18N(50, 900)
		if !w18IsPrime(n) && n > 1 {
			_ = primes
			return n
		}
	}
}

func w18TrialDivision(o int) generator.Problem {
	n := w18Composite()
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		ans := w18SmallestFactor(n)
		return w18P(o, fmt.Sprintf("Trial division starts at two. What is the smallest prime factor of %d?", n), strconv.Itoa(ans), fmt.Sprintf("The smallest divisor above one is %d.", ans), f)
	case 2:
		ans := len(w18PrimesUpTo(w18FloorSqrt(n)))
		return w18P(o, fmt.Sprintf("Trial division of %d tests every prime up to the square root. How many primes are tested?", n), strconv.Itoa(ans), fmt.Sprintf("The integer square root of %d is %d, and there are %d primes at most it.", n, w18FloorSqrt(n), ans), f)
	default:
		ans := w18LargestFactor(n)
		return w18P(o, fmt.Sprintf("What is the largest prime factor of %d?", n), strconv.Itoa(ans), fmt.Sprintf("Factor %d completely; the largest prime factor is %d.", n, ans), f)
	}
}

func w18PollardRho(o int) generator.Problem {
	var primes []int
	for _, q := range w18PrimesUpTo(97) {
		if q >= 11 {
			primes = append(primes, q)
		}
	}
	p := primes[rand.Intn(len(primes))]
	q := primes[rand.Intn(len(primes))]
	for q == p {
		q = primes[rand.Intn(len(primes))]
	}
	n := p * q
	c := w18N(1, 4)
	x0 := 2
	x1 := (x0*x0 + c) % n
	x2 := (x1*x1 + c) % n
	f := map[string]string{}
	w18I(f, "n", n)
	w18I(f, "c", c)
	w18I(f, "x0", x0)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("Pollard rho iterates x squared plus %d modulo %d starting from x=%d. What is the next value?", c, n, x0), strconv.Itoa(x1), fmt.Sprintf("Square %d, add %d and reduce modulo %d, giving %d.", x0, c, n, x1), f)
	case 2:
		return w18P(o, fmt.Sprintf("Pollard rho iterates x squared plus %d modulo %d. Starting from %d, what is the value after two steps?", c, n, x0), strconv.Itoa(x2), fmt.Sprintf("One step gives %d, and squaring it again gives %d.", x1, x2), f)
	default:
		ans := w18FloorSqrt(p)
		return w18P(o, fmt.Sprintf("Pollard rho is expected to find a factor p after about the square root of p steps. For p=%d, what is the floor of that square root?", p), strconv.Itoa(ans), fmt.Sprintf("The integer square root of %d is %d.", p, ans), f)
	}
}

func w18QuadraticSieve(o int) generator.Problem {
	b := 20 + 10*w18N(0, 18)
	pi := len(w18PrimesUpTo(b))
	f := map[string]string{}
	w18I(f, "b", b)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("The quadratic sieve uses a smoothness bound B=%d. How many primes are at most B?", b), strconv.Itoa(pi), fmt.Sprintf("Counting the primes up to %d gives %d.", b, pi), f)
	case 2:
		ans := pi + 1
		return w18P(o, fmt.Sprintf("With %d primes in the factor base, how many smooth relations are needed before a dependency is guaranteed?", pi), strconv.Itoa(ans), fmt.Sprintf("One more relation than the number of primes, which is %d.", ans), f)
	default:
		ans := pi - 1
		return w18P(o, fmt.Sprintf("The quadratic sieve factor base for B=%d excludes the prime two when n is odd. How many primes remain?", b), strconv.Itoa(ans), fmt.Sprintf("There are %d primes at most %d, so excluding two leaves %d.", pi, b, ans), f)
	}
}

func w18AKS(o int) generator.Problem {
	n := w18N(50, 500)
	rPrimes := w18PrimesUpTo(60)
	r := rPrimes[1+rand.Intn(len(rPrimes)-1)]
	phi := w18Totient(r)
	bound := int(math.Floor(math.Sqrt(float64(phi)) * math.Log2(float64(n))))
	f := map[string]string{}
	w18I(f, "n", n)
	w18I(f, "r", r)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("The AKS test checks polynomial identities for shifts a up to the floor of the square root of phi(r) times log base two of n. For r=%d and n=%d, what is that bound?", r, n), strconv.Itoa(bound), fmt.Sprintf("The totient of %d is %d, and its square root times log base two of %d gives a bound of %d.", r, phi, n, bound), f)
	case 2:
		return w18P(o, fmt.Sprintf("For the AKS modulus r=%d, what is phi(r)?", r), strconv.Itoa(phi), fmt.Sprintf("Counting integers up to %d coprime to it gives %d.", r, phi), f)
	default:
		return w18P(o, fmt.Sprintf("The AKS test works modulo X to the r minus one. For r=%d, what is the degree of that polynomial?", r), strconv.Itoa(r), fmt.Sprintf("The polynomial X to the %d minus one has degree %d.", r, r), f)
	}
}

func w18BSGS(o int) generator.Problem {
	n := 100 + 10*w18N(0, 90)
	m := w18CeilSqrt(n)
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("Baby-step giant-step sets m to the ceiling of the square root of the group order. For order %d, what is m?", n), strconv.Itoa(m), fmt.Sprintf("The integer square root of %d rounded up is %d.", n, m), f)
	case 2:
		ans := 2 * m
		return w18P(o, fmt.Sprintf("With m=%d baby steps and m giant steps, how many group operations does baby-step giant-step perform in total?", m), strconv.Itoa(ans), fmt.Sprintf("Two tables of size %d give %d operations.", m, ans), f)
	default:
		ans := m * m
		return w18P(o, fmt.Sprintf("Baby-step giant-step compares m baby steps against m giant steps. With m=%d, how many comparisons are made in the worst case?", m), strconv.Itoa(ans), fmt.Sprintf("Multiplying %d by itself gives %d.", m, ans), f)
	}
}

func w18PohligHellman(o int) generator.Problem {
	n := 1
	for n < 100 || n > 5000 {
		n = 1
		for _, p := range []int{2, 3, 5, 7} {
			n *= w18IPow(p, w18N(0, 3))
		}
	}
	largest := w18LargestFactor(n)
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("Pohlig-Hellman works hardest on the largest prime factor of the group order. For order %d, what is that factor?", n), strconv.Itoa(largest), fmt.Sprintf("Factor %d; the largest prime factor is %d.", n, largest), f)
	case 2:
		ans := w18Omega(n)
		return w18P(o, fmt.Sprintf("For Pohlig-Hellman on a group of order %d, how many prime factors does the order have counting multiplicity?", n), strconv.Itoa(ans), fmt.Sprintf("The factorisation of %d has %d prime factors in total.", n, ans), f)
	default:
		ans := w18DistinctFactors(n)
		return w18P(o, fmt.Sprintf("How many distinct prime factors does the group order %d have?", n), strconv.Itoa(ans), fmt.Sprintf("Factoring %d gives %d distinct primes.", n, ans), f)
	}
}

func w18LLL(o int) generator.Problem {
	d := w18N(3, 30)
	f := map[string]string{}
	w18I(f, "d", d)
	switch o {
	case 1:
		ans := d - 1
		return w18P(o, fmt.Sprintf("LLL guarantees a vector within a factor of two to the (d-1)/2 in dimension d. For d=%d, what is the numerator of that exponent?", d), strconv.Itoa(ans), fmt.Sprintf("The exponent is (d-1)/2, so the numerator is %d.", ans), f)
	case 2:
		return w18P(o, fmt.Sprintf("A lattice basis in dimension %d is reduced by LLL. How many basis vectors are there?", d), strconv.Itoa(d), fmt.Sprintf("The dimension equals the number of basis vectors, which is %d.", d), f)
	default:
		ans := d * (d - 1) / 2
		return w18P(o, fmt.Sprintf("A lattice in dimension %d has how many unordered pairs of basis vectors?", d), strconv.Itoa(ans), fmt.Sprintf("Choose two of %d, giving %d.", d, ans), f)
	}
}

func w18Pratt(o int) generator.Problem {
	primes := w18PrimesUpTo(997)
	p := primes[w18N(25, len(primes)-1)]
	counts := w18FactorCounts(p - 1)
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
		ans := len(counts)
		return w18P(o, fmt.Sprintf("A Pratt certificate for %d recursively certifies each prime factor of p minus one. How many distinct factors does %d have?", p, p-1), strconv.Itoa(ans), fmt.Sprintf("Factoring %d gives %d distinct primes.", p-1, ans), f)
	case 2:
		ans := w18Omega(p - 1)
		return w18P(o, fmt.Sprintf("Counting multiplicity, how many prime factors does %d have?", p-1), strconv.Itoa(ans), fmt.Sprintf("The factorisation of %d contains %d prime factors in total.", p-1, ans), f)
	default:
		return w18P(o, fmt.Sprintf("What is the largest prime factor of %d, the first recursive step of a Pratt certificate for %d?", p-1, p), strconv.Itoa(largest), fmt.Sprintf("Factoring %d gives a largest prime factor of %d.", p-1, largest), f)
	}
}

func w18CFFactor(o int) generator.Problem {
	n := w18N(20, 200)
	for w18FloorSqrt(n)*w18FloorSqrt(n) == n {
		n++
	}
	m, d, a := 0, 1, w18FloorSqrt(n)
	a0 := a
	period := 0
	firstA1 := 0
	for {
		m = d*a - m
		d = (n - m*m) / d
		a = (a0 + m) / d
		period++
		if period == 1 {
			firstA1 = a
		}
		if a == 2*a0 {
			break
		}
		if period > 200 {
			break
		}
	}
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("What is the period length of the continued fraction of the square root of %d?", n), strconv.Itoa(period), fmt.Sprintf("Expanding the continued fraction returns to twice the floor after %d steps.", period), f)
	case 2:
		return w18P(o, fmt.Sprintf("What is the integer part of the square root of %d, the first term of its continued fraction?", n), strconv.Itoa(a0), fmt.Sprintf("The integer square root of %d is %d.", n, a0), f)
	default:
		ans := a0*firstA1 + 1
		return w18P(o, fmt.Sprintf("The continued fraction of the square root of %d begins with terms %d and %d. What is the numerator of the second convergent?", n, a0, firstA1), strconv.Itoa(ans), fmt.Sprintf("The numerator is a0 times a1 plus one, which is %d.", ans), f)
	}
}

func w18NFS(o int) generator.Problem {
	d := w18N(20, 120)
	n := math.Pow(10, float64(d))
	ln := math.Log(n)
	lnln := math.Log(ln)
	degree := int(math.Floor(math.Pow(3*ln/lnln, 1.0/3.0)))
	root := int(math.Floor(math.Pow(ln, 1.0/3.0)))
	f := map[string]string{}
	w18I(f, "d", d)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("The number field sieve chooses a polynomial degree near the cube root of three log n over log log n. For n with %d digits, what is that degree rounded down?", d), strconv.Itoa(degree), fmt.Sprintf("Evaluating the formula at ten to the %d gives %d.", d, degree), f)
	case 2:
		return w18P(o, fmt.Sprintf("The subexponential exponent of the number field sieve contains the cube root of log n. For n with %d digits, what is that cube root rounded down?", d), strconv.Itoa(root), fmt.Sprintf("The cube root of the natural logarithm of ten to the %d is about %d.", d, root), f)
	default:
		return w18P(o, fmt.Sprintf("For n equal to ten to the %d, how many decimal digits does n have?", d), strconv.Itoa(d+1), fmt.Sprintf("Ten to the %d has %d digits.", d, d+1), f)
	}
}

func w18Partitions(o int) generator.Problem {
	n := w18N(5, 40)
	p := make([]int, n+1)
	p[0] = 1
	for part := 1; part <= n; part++ {
		for total := part; total <= n; total++ {
			p[total] += p[total-part]
		}
	}
	q := make([]int, n+1)
	q[0] = 1
	for part := 1; part <= n; part++ {
		for total := n; total >= part; total-- {
			q[total] += q[total-part]
		}
	}
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("How many partitions does the integer %d have?", n), strconv.Itoa(p[n]), fmt.Sprintf("The partition count p(%d) is %d.", n, p[n]), f)
	case 2:
		return w18P(o, fmt.Sprintf("How many partitions of %d into distinct parts are there?", n), strconv.Itoa(q[n]), fmt.Sprintf("Counting partitions with no repeated part gives %d.", q[n]), f)
	default:
		ans := n/2 + 1
		return w18P(o, fmt.Sprintf("How many partitions of %d have at most two parts?", n), strconv.Itoa(ans), fmt.Sprintf("The parts are the integers from zero up to half of %d, giving %d partitions.", n, ans), f)
	}
}

func w18Vinogradov(o int) generator.Problem {
	primes := w18PrimesUpTo(60)
	switch o {
	case 1:
		n := 7 + 2*rand.Intn(23)
		count := 0
		for i := 0; i < len(primes); i++ {
			for j := i; j < len(primes); j++ {
				for k := j; k < len(primes); k++ {
					if primes[i]+primes[j]+primes[k] == n {
						count++
					}
				}
			}
		}
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("How many unordered ways are there to write %d as a sum of three primes?", n), strconv.Itoa(count), fmt.Sprintf("Counting unordered triples of primes with sum %d gives %d.", n, count), f)
	case 2:
		n := 7 + 2*rand.Intn(23)
		ans := len(w18PrimesUpTo(n))
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("How many primes are at most %d?", n), strconv.Itoa(ans), fmt.Sprintf("Counting the primes up to %d gives %d.", n, ans), f)
	default:
		n := 9 + 2*rand.Intn(22)
		ans := 0
		if w18IsPrime(n - 2) {
			ans = 1
		}
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("How many unordered ways are there to write the odd number %d as a sum of two primes?", n), strconv.Itoa(ans), fmt.Sprintf("One summand must be two, so the count is one exactly when %d is prime; here it is %d.", n-2, ans), f)
	}
}

func w18CircleMethod(o int) generator.Problem {
	switch o {
	case 1:
		q := w18N(5, 40)
		ans := w18Totient(q)
		f := map[string]string{}
		w18I(f, "q", q)
		return w18P(o, fmt.Sprintf("The major arcs of the circle method are indexed by reduced fractions with denominator %d. How many such residues are there?", q), strconv.Itoa(ans), fmt.Sprintf("The count of residues coprime to %d is %d.", q, ans), f)
	case 2:
		q := w18N(3, 40)
		ans := 1
		for i := 1; i <= q; i++ {
			ans += w18Totient(i)
		}
		f := map[string]string{}
		w18I(f, "q", q)
		return w18P(o, fmt.Sprintf("How many fractions lie in the Farey sequence of order %d, including zero and one?", q), strconv.Itoa(ans), fmt.Sprintf("The length is one plus the sum of the totients up to %d, which is %d.", q, ans), f)
	default:
		q := w18N(5, 40)
		a := w18N(1, q-1)
		for w18GCD(a, q) != 1 {
			a = w18N(1, q-1)
		}
		f := map[string]string{}
		w18I(f, "q", q)
		w18I(f, "a", a)
		return w18P(o, fmt.Sprintf("The fraction %d over %d is written in lowest terms as a major-arc point. What is its denominator?", a, q), strconv.Itoa(q), fmt.Sprintf("The numerator %d is coprime to %d, so the denominator stays %d.", a, q, q), f)
	}
}

func w18Schnirelmann(o int) generator.Problem {
	m := w18N(2, 40)
	f := map[string]string{}
	w18I(f, "m", m)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("The set of multiples of %d, together with zero, has Schnirelmann density one over m. What is m?", m), strconv.Itoa(m), fmt.Sprintf("The density is one over %d, so the denominator is %d.", m, m), f)
	case 2:
		return w18P(o, fmt.Sprintf("The set of integers congruent to a fixed residue modulo %d has Schnirelmann density one over m. What is m?", m), strconv.Itoa(m), fmt.Sprintf("The density is one over %d.", m), f)
	default:
		ans := m
		return w18P(o, fmt.Sprintf("A set with Schnirelmann density one over %d needs at least how many copies of itself in a sumset to cover all positive integers, by the elementary bound?", m), strconv.Itoa(ans), fmt.Sprintf("The bound is the reciprocal of the density, which is %d.", ans), f)
	}
}

func w18Sumsets(o int) generator.Problem {
	switch o {
	case 1:
		a := w18N(3, 15)
		b := w18N(3, 15)
		ans := a + b - 1
		f := map[string]string{}
		w18I(f, "a", a)
		w18I(f, "b", b)
		return w18P(o, fmt.Sprintf("An interval of %d integers is added to an interval of %d integers. How many integers are in the sumset?", a, b), strconv.Itoa(ans), fmt.Sprintf("The sumset is an interval of size a+b-1, which is %d.", ans), f)
	case 2:
		a := w18N(3, 30)
		ans := 2*a - 1
		f := map[string]string{}
		w18I(f, "a", a)
		return w18P(o, fmt.Sprintf("An interval of %d integers is added to itself. How many integers are in the sumset?", a), strconv.Itoa(ans), fmt.Sprintf("Doubling an interval of size %d gives size 2a-1, which is %d.", a, ans), f)
	default:
		m := w18N(5, 12)
		setA := w18RandomSet(m, w18N(2, 4))
		setB := w18RandomSet(m, w18N(2, 4))
		seen := map[int]bool{}
		for _, x := range setA {
			for _, y := range setB {
				seen[(x+y)%m] = true
			}
		}
		ans := len(seen)
		f := map[string]string{}
		w18I(f, "m", m)
		f["a"] = w18JoinInts(setA)
		f["b"] = w18JoinInts(setB)
		return w18P(o, fmt.Sprintf("In the integers modulo %d, set A is %v and set B is %v. How many residues are in the sumset A+B?", m, setA, setB), strconv.Itoa(ans), fmt.Sprintf("Adding every pair modulo %d and removing duplicates gives %d residues.", m, ans), f)
	}
}

func w18RandomSet(m, size int) []int {
	if size > m {
		size = m
	}
	chosen := map[int]bool{}
	var out []int
	for len(out) < size {
		v := rand.Intn(m)
		if chosen[v] {
			continue
		}
		chosen[v] = true
		out = append(out, v)
	}
	return out
}

func w18JoinInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w18CauchyDavenport(o int) generator.Problem {
	var fieldPrimes []int
	for _, q := range w18PrimesUpTo(97) {
		if q >= 5 {
			fieldPrimes = append(fieldPrimes, q)
		}
	}
	p := fieldPrimes[rand.Intn(len(fieldPrimes))]
	a := w18N(1, p)
	b := w18N(1, p)
	f := map[string]string{}
	w18I(f, "p", p)
	w18I(f, "a", a)
	w18I(f, "b", b)
	raw := a + b - 1
	ans := raw
	if ans > p {
		ans = p
	}
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("In the field modulo %d, sets of sizes %d and %d are added. Cauchy-Davenport gives a lower bound of the minimum of p and a+b-1. What is that bound?", p, a, b), strconv.Itoa(ans), fmt.Sprintf("The raw bound is %d, capped at %d, giving %d.", raw, p, ans), f)
	case 2:
		return w18P(o, fmt.Sprintf("For set sizes %d and %d, what is a+b-1 before the field-size cap?", a, b), strconv.Itoa(raw), fmt.Sprintf("Adding the sizes and subtracting one gives %d.", raw), f)
	default:
		return w18P(o, fmt.Sprintf("In a prime field of size p, a sumset cannot exceed the whole field. What is the largest possible sumset size when p=%d?", p), strconv.Itoa(p), fmt.Sprintf("The prime is %d, so the sumset is capped there.", p), f)
	}
}

func w18Freiman(o int) generator.Problem {
	k := w18N(3, 20)
	f := map[string]string{}
	w18I(f, "k", k)
	switch o {
	case 1:
		ans := 3*k - 4
		return w18P(o, fmt.Sprintf("Freiman's 3k-4 theorem bounds a set of %d integers with a small sumset. What is 3k-4 for k=%d?", k, k), strconv.Itoa(ans), fmt.Sprintf("Three times %d minus four is %d.", k, ans), f)
	case 2:
		ans := 2*k - 3
		return w18P(o, fmt.Sprintf("The conclusion of the 3k-4 theorem places the set in an interval of length at most 2k-3. For k=%d, what is that length?", k), strconv.Itoa(ans), fmt.Sprintf("Two times %d minus three is %d.", k, ans), f)
	default:
		ans := 2*k - 1
		return w18P(o, fmt.Sprintf("A set of %d integers has a sumset of size at least 2k-1. What is that lower bound?", k), strconv.Itoa(ans), fmt.Sprintf("Two times %d minus one is %d.", k, ans), f)
	}
}

func w18AdditiveBases(o int) generator.Problem {
	switch o {
	case 1:
		k := w18N(2, 20)
		ans := w18IPow(2, k) + w18IPow(3, k)/w18IPow(2, k) - 2
		f := map[string]string{}
		w18I(f, "k", k)
		return w18P(o, fmt.Sprintf("The set of k-th powers is an additive basis of order g(k). For k=%d, what is g(k)?", k), strconv.Itoa(ans), fmt.Sprintf("Waring's formula for %d gives %d.", k, ans), f)
	case 2:
		n := w18N(10, 200)
		f := map[string]string{}
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("Lagrange's theorem represents the integer %d as a sum of at most how many squares?", n), "4", "Four squares always suffice, and some integers need all four, so the order of the squares as an additive basis is four.", f)
	default:
		k := w18N(2, 4)
		n := w18N(10, 200)
		ans := 0
		for v := 1; w18IPow(v, k) <= n; v++ {
			ans++
		}
		f := map[string]string{}
		w18I(f, "k", k)
		w18I(f, "n", n)
		return w18P(o, fmt.Sprintf("How many positive %d-th powers are at most %d?", k, n), strconv.Itoa(ans), fmt.Sprintf("Count the bases whose %d-th power does not exceed %d, giving %d.", k, n, ans), f)
	}
}

func w18GoldbachCount(o int) generator.Problem {
	n := 8 + 2*rand.Intn(27)
	primes := w18PrimesUpTo(n)
	count := 0
	firstQ := 0
	for _, p := range primes {
		if p > n-p {
			break
		}
		if w18IsPrime(n - p) {
			count++
			if firstQ == 0 {
				firstQ = n - p
			}
		}
	}
	f := map[string]string{}
	w18I(f, "n", n)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("How many unordered ways are there to write %d as a sum of two primes?", n), strconv.Itoa(count), fmt.Sprintf("Counting prime pairs summing to %d gives %d.", n, count), f)
	case 2:
		ans := len(primes)
		return w18P(o, fmt.Sprintf("How many primes are at most %d?", n), strconv.Itoa(ans), fmt.Sprintf("Counting the primes up to %d gives %d.", n, ans), f)
	default:
		return w18P(o, fmt.Sprintf("In the first Goldbach partition of %d, where the smaller prime is as small as possible, what is the larger prime?", n), strconv.Itoa(firstQ), fmt.Sprintf("The smallest prime p with n-p prime gives the larger partner %d.", firstQ), f)
	}
}

func w18AdditiveEnergy(o int) generator.Problem {
	m := w18N(5, 30)
	set := w18RandomSet(m, w18N(2, 25))
	counts := map[int]int{}
	for _, a := range set {
		for _, b := range set {
			counts[(a+b)%m]++
		}
	}
	energy := 0
	for _, c := range counts {
		energy += c * c
	}
	size := len(set)
	f := map[string]string{}
	w18I(f, "m", m)
	f["a"] = w18JoinInts(set)
	switch o {
	case 1:
		return w18P(o, fmt.Sprintf("In the integers modulo %d, the set A is %v. What is its additive energy, the number of ordered quadruples with a+b equal to c+d?", m, set), strconv.Itoa(energy), fmt.Sprintf("Grouping ordered pairs by their sum and adding the squares of the group sizes gives %d.", energy), f)
	case 2:
		ans := size * size * size
		return w18P(o, fmt.Sprintf("A set A has %d elements. What is the cube of its size, the trivial lower bound for additive energy?", size), strconv.Itoa(ans), fmt.Sprintf("The cube of %d is %d.", size, ans), f)
	default:
		ans := size * size
		return w18P(o, fmt.Sprintf("A set A has %d elements. How many ordered pairs of elements does A have?", size), strconv.Itoa(ans), fmt.Sprintf("Squaring the size gives %d.", ans), f)
	}
}
