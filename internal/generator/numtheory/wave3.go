package numtheory

// Wave 3 number theory: twenty concepts. Each concept has three knowledge
// points, and every knowledge point publishes the operands and the derived
// values its answer follows from. The tests in wave3_test.go recompute the
// answer from those facts rather than trusting the generator's own arithmetic.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave3(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"nt.basics.mod_power":      "mod_power",
		"nt.basics.mod_inverse":    "mod_inverse",
		"nt.basics.euler_theorem":  "euler_theorem",
		"nt.basics.decimal_period": "decimal_period",
		"nt.basics.digital_root":   "digital_root",
		"nt.adv.aliquot":           "aliquot",
		"nt.adv.perfect_numbers":   "perfect_numbers",
		"nt.adv.mersenne":          "mersenne",
		"nt.adv.fermat_numbers":    "fermat_numbers",
		"nt.adv.repunit":           "repunit",
		"nt.adv.prime_counting":    "prime_counting",
		"nt.adv.twin_primes":       "twin_primes",
		"nt.adv.goldbach":          "goldbach",
		"nt.adv.polygonal":         "polygonal",
		"nt.adv.pentagonal":        "pentagonal",
		"nt.adv.egyptian":          "egyptian",
		"nt.adv.highly_composite":  "highly_composite",
		"nt.adv.amicable":          "amicable",
		"nt.adv.happy":             "happy",
		"nt.adv.consecutive_sums":  "consecutive_sums",
	} {
		reg.Register(id, &w3NumberGen{kind: kind})
	}
}

type w3NumberGen struct{ kind string }

func (g *w3NumberGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "mod_power":
		return w3ModPower(kp)
	case "mod_inverse":
		return w3ModInverse(kp)
	case "euler_theorem":
		return w3EulerTheorem(kp)
	case "decimal_period":
		return w3DecimalPeriod(kp)
	case "digital_root":
		return w3DigitalRootProblem(kp)
	case "aliquot":
		return w3Aliquot(kp)
	case "perfect_numbers":
		return w3PerfectNumbers(kp)
	case "mersenne":
		return w3Mersenne(kp)
	case "fermat_numbers":
		return w3FermatNumbers(kp)
	case "repunit":
		return w3Repunit(kp)
	case "prime_counting":
		return w3PrimeCounting(kp)
	case "twin_primes":
		return w3TwinPrimes(kp)
	case "goldbach":
		return w3Goldbach(kp)
	case "polygonal":
		return w3Polygonal(kp)
	case "pentagonal":
		return w3Pentagonal(kp)
	case "egyptian":
		return w3Egyptian(kp)
	case "highly_composite":
		return w3HighlyComposite(kp)
	case "amicable":
		return w3Amicable(kp)
	case "happy":
		return w3Happy(kp)
	case "consecutive_sums":
		return w3ConsecutiveSums(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w3P(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w3IntFact(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w3Small(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w3Divisors(n int) []int {
	var out []int
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			out = append(out, d)
		}
	}
	return out
}

func w3ProperDivisorSum(n int) int {
	sum := 0
	for _, d := range w3Divisors(n) {
		if d != n {
			sum += d
		}
	}
	return sum
}

func w3DivisorCount(n int) int { return len(w3Divisors(n)) }

func w3IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func w3PrimeList(limit int) []int {
	var out []int
	for n := 2; n <= limit; n++ {
		if w3IsPrime(n) {
			out = append(out, n)
		}
	}
	return out
}

func w3PowMod(base, exp, mod int) int {
	if mod == 1 {
		return 0
	}
	r := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			r = r * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return r
}

func w3Phi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if computeGCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w3Order10(n int) int {
	if n <= 1 {
		return 0
	}
	r := 10 % n
	for k := 1; k <= n; k++ {
		if r == 1 {
			return k
		}
		r = r * 10 % n
	}
	return 0
}

func w3DigitalRoot(n int) int {
	for n >= 10 {
		s := 0
		for n > 0 {
			s += n % 10
			n /= 10
		}
		n = s
	}
	return n
}

func w3IsPerfect(n int) bool   { return w3ProperDivisorSum(n) == n }
func w3IsAbundant(n int) bool  { return w3ProperDivisorSum(n) > n }
func w3IsDeficient(n int) bool { return w3ProperDivisorSum(n) < n }
func w3IsHighlyComposite(n int) bool {
	if n < 1 {
		return false
	}
	if n == 1 {
		return true
	}
	d := w3DivisorCount(n)
	for m := 1; m < n; m++ {
		if w3DivisorCount(m) >= d {
			return false
		}
	}
	return true
}

func w3IsAmicable(n int) bool {
	m := w3ProperDivisorSum(n)
	return m != n && w3ProperDivisorSum(m) == n
}

func w3IsHappy(n int) bool {
	seen := map[int]bool{}
	for n != 1 && !seen[n] {
		seen[n] = true
		s := 0
		for n > 0 {
			d := n % 10
			s += d * d
			n /= 10
		}
		n = s
	}
	return n == 1
}

func w3IsPrimePowerOfTwo(n int) bool {
	if n < 1 {
		return false
	}
	for n > 1 {
		if n%2 != 0 {
			return false
		}
		n /= 2
	}
	return true
}

func w3CountConsecutive(n int) int {
	count := 0
	// Length 1 is excluded, so the search starts at k = 2.
	for k := 2; k*(k+1)/2 <= n; k++ {
		// n = k(2a+k-1)/2 for a >= 1.
		if (2*n)%k != 0 {
			continue
		}
		if (2*n/k-k+1)%2 == 0 {
			count++
		}
	}
	return count
}

func w3ShortestConsecutive(n int) int {
	for k := 2; k*(k+1)/2 <= n; k++ {
		if (2*n)%k != 0 {
			continue
		}
		a := (2*n/k - k + 1) / 2
		if a >= 1 && (2*n/k-k+1)%2 == 0 {
			return k
		}
	}
	return 0
}

func w3ModPower(kp int) generator.Problem {
	a := w3Small(2, 19)
	e := w3Small(3, 12)
	m := w3Small(3, 23)
	f := map[string]string{}
	w3IntFact(f, "base", a)
	w3IntFact(f, "exponent", e)
	w3IntFact(f, "modulus", m)
	switch kp {
	case 1:
		ans := w3PowMod(a, e, m)
		return w3P(kp, fmt.Sprintf("Compute \\(%d^{%d} \\bmod %d\\). (enter a number)", a, e, m), strconv.Itoa(ans),
			fmt.Sprintf("Repeated squaring gives \\(%d^{%d} \\bmod %d = %d\\).", a, e, m, ans), f)
	case 2:
		for computeGCD(a, m) != 1 {
			a = w3Small(2, 19)
		}
		w3IntFact(f, "base", a)
		big := e + 20*w3Small(1, 4)
		w3IntFact(f, "exponent", big)
		ans := w3PowMod(a, big, m)
		return w3P(kp, fmt.Sprintf("Use exponent reduction to compute \\(%d^{%d} \\bmod %d\\). (enter a number)", a, big, m), strconv.Itoa(ans),
			fmt.Sprintf("Repeated squaring, or Fermat reduction when applicable, gives \\(%d^{%d} \\bmod %d = %d\\).", a, big, m, ans), f)
	default:
		b := w3Small(2, 17)
		x := w3Small(2, 9)
		y := w3Small(2, 9)
		w3IntFact(f, "base2", b)
		w3IntFact(f, "exponent2", x)
		w3IntFact(f, "exponent3", y)
		ans := w3PowMod(w3PowMod(a, x, m)*w3PowMod(b, y, m), 1, m)
		return w3P(kp, fmt.Sprintf("Compute \\(%d^{%d} \\cdot %d^{%d} \\bmod %d\\). (enter a number)", a, x, b, y, m), strconv.Itoa(ans),
			fmt.Sprintf("Reduce each power modulo %d, multiply the residues, then reduce again: the result is %d.", m, ans), f)
	}
}

func w3ModInverse(kp int) generator.Problem {
	m := w3Small(5, 29)
	a := w3Small(2, m-1)
	for computeGCD(a, m) != 1 {
		a = w3Small(2, m-1)
	}
	f := map[string]string{}
	w3IntFact(f, "a", a)
	w3IntFact(f, "modulus", m)
	inv := ntBridgeInvMod(a, m)
	w3IntFact(f, "inverse", inv)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("Find the inverse of %d modulo %d, as a least residue. (enter a number)", a, m), strconv.Itoa(inv),
			fmt.Sprintf("Because \\(\\gcd(%d,%d)=1\\), the inverse exists. \\(%d \\cdot %d \\equiv 1 \\pmod{%d}\\), so the inverse is %d.", a, m, a, inv, m, inv), f)
	case 2:
		b := w3Small(1, m-1)
		w3IntFact(f, "b", b)
		ans := inv * b % m
		return w3P(kp, fmt.Sprintf("Solve \\(%d x \\equiv %d \\pmod{%d}\\) for the least nonnegative \\(x\\). (enter a number)", a, b, m), strconv.Itoa(ans),
			fmt.Sprintf("Multiply by the inverse %d: \\(x \\equiv %d \\cdot %d \\equiv %d \\pmod{%d}\\).", inv, inv, b, ans, m), f)
	default:
		b := w3Small(1, m-1)
		w3IntFact(f, "b", b)
		ans := inv * b % m
		return w3P(kp, fmt.Sprintf("Solve the congruence \\(%d x \\equiv %d \\pmod{%d}\\) by inverting %d. (enter a number)", a, b, m, a), strconv.Itoa(ans),
			fmt.Sprintf("The inverse of %d modulo %d is %d, so \\(x \\equiv %d \\cdot %d = %d \\pmod{%d}\\).", a, m, inv, inv, b, ans, m), f)
	}
}

func w3EulerTheorem(kp int) generator.Problem {
	n := w3Small(6, 40)
	a := w3Small(2, n-1)
	for computeGCD(a, n) != 1 {
		a = w3Small(2, n-1)
	}
	phi := w3Phi(n)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	w3IntFact(f, "a", a)
	w3IntFact(f, "phi", phi)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is \\(\\varphi(%d)\\)? (enter a number)", n), strconv.Itoa(phi),
			fmt.Sprintf("Count the integers from 1 to %d that are coprime to %d: there are %d.", n, n, phi), f)
	case 2:
		e := phi + w3Small(1, 15)
		w3IntFact(f, "exponent", e)
		ans := w3PowMod(a, e, n)
		return w3P(kp, fmt.Sprintf("Euler's theorem says \\(%d^{\\varphi(%d)} \\equiv 1 \\pmod{%d}\\). Use it to compute \\(%d^{%d} \\bmod %d\\). (enter a number)", a, n, n, a, e, n), strconv.Itoa(ans),
			fmt.Sprintf("Since \\(\\varphi(%d)=%d\\), reduce the exponent modulo %d. The resulting residue is %d.", n, phi, phi, ans), f)
	default:
		e := w3Small(2, 9) * phi
		w3IntFact(f, "exponent", e)
		ans := w3PowMod(a, e, n)
		return w3P(kp, fmt.Sprintf("Compute \\(%d^{%d} \\bmod %d\\) using Euler's theorem. (enter a number)", a, e, n), strconv.Itoa(ans),
			fmt.Sprintf("The exponent is a multiple of \\(\\varphi(%d)=%d\\), so the power is \\(1 \\bmod %d\\).", n, phi, n), f)
	}
}

func w3DecimalPeriod(kp int) generator.Problem {
	n := w3Small(3, 37)
	if kp == 3 {
		// Keep a nonterminating core so the requested order is defined.
		for {
			m := n
			for m%2 == 0 {
				m /= 2
			}
			for m%5 == 0 {
				m /= 5
			}
			if m > 1 {
				break
			}
			n = w3Small(3, 60)
		}
	} else {
		for n%2 == 0 || n%5 == 0 {
			n = w3Small(3, 37)
		}
	}
	f := map[string]string{}
	w3IntFact(f, "denominator", n)
	period := w3Order10(n)
	w3IntFact(f, "period", period)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the repeating period of the decimal expansion of \\(1/%d\\)? (enter a number)", n), strconv.Itoa(period),
			fmt.Sprintf("The period is the multiplicative order of 10 modulo %d. The first positive \\(k\\) with \\(10^k \\equiv 1 \\pmod{%d}\\) is %d.", n, n, period), f)
	case 2:
		k := w3Small(1, n-1)
		for computeGCD(k, n) != 1 {
			k = w3Small(1, n-1)
		}
		w3IntFact(f, "numerator", k)
		p := w3Order10(n)
		return w3P(kp, fmt.Sprintf("What is the repeating period of the decimal expansion of \\(%d/%d\\)? (enter a number)", k, n), strconv.Itoa(p),
			fmt.Sprintf("The fraction is already reduced, so the period is the order of 10 modulo %d, namely %d.", n, p), f)
	default:
		twos, fives := 0, 0
		m := n
		for m%2 == 0 {
			twos++
			m /= 2
		}
		for m%5 == 0 {
			fives++
			m /= 5
		}
		w3IntFact(f, "twos", twos)
		w3IntFact(f, "fives", fives)
		w3IntFact(f, "core", m)
		p := w3Order10(m)
		if twos == 0 && fives == 0 {
			return w3P(kp, fmt.Sprintf("The denominator %d already has no factors of 2 or 5. Its decimal period is the order of 10 modulo %d. What is that order? (enter a number)", n, n), strconv.Itoa(p),
				fmt.Sprintf("The period is the smallest positive k with 10 to the k congruent to 1 modulo %d. That k is %d.", n, p), f)
		}
		twosText := fmt.Sprintf("%d factors", twos)
		if twos == 1 {
			twosText = "1 factor"
		}
		fivesText := fmt.Sprintf("%d factors", fives)
		if fives == 1 {
			fivesText = "1 factor"
		}
		return w3P(kp, fmt.Sprintf("After removing the factors 2 and 5 from the denominator, the decimal period of \\(1/%d\\) is the order of 10 modulo %d. What is that order? (enter a number)", n, m), strconv.Itoa(p),
			fmt.Sprintf("There are %s of 2 and %s of 5. They contribute only a terminating prefix. The period is the order of 10 modulo %d, which is %d.", twosText, fivesText, m, p), f)
	}
}

func w3DigitalRootProblem(kp int) generator.Problem {
	n := w3Small(12, 999)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	root := w3DigitalRoot(n)
	w3IntFact(f, "digital_root", root)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the digital root of %d? (enter a number)", n), strconv.Itoa(root),
			fmt.Sprintf("Sum the digits repeatedly until one digit remains. The digital root of %d is %d.", n, root), f)
	case 2:
		a := w3Small(10, 199)
		b := w3Small(10, 199)
		w3IntFact(f, "a", a)
		w3IntFact(f, "b", b)
		r := w3DigitalRoot(a + b)
		return w3P(kp, fmt.Sprintf("What is the digital root of \\(%d + %d\\)? (enter a number)", a, b), strconv.Itoa(r),
			fmt.Sprintf("A digital root is congruent to the number modulo 9. Here \\(%d + %d = %d\\), whose digital root is %d.", a, b, a+b, r), f)
	default:
		r := n % 9
		if r == 0 {
			r = 9
		}
		return w3P(kp, fmt.Sprintf("Casting out nines, what is the least positive residue of %d modulo 9? (enter a number)", n), strconv.Itoa(r),
			fmt.Sprintf("Reduce the digit sum modulo 9. The least positive residue of %d modulo 9 is %d.", n, r), f)
	}
}

func w3Aliquot(kp int) generator.Problem {
	n := w3Small(10, 120)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	s := w3ProperDivisorSum(n)
	sigma := s + n
	w3IntFact(f, "aliquot_sum", s)
	w3IntFact(f, "sigma", sigma)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the aliquot sum \\(s(%d)\\), the sum of the proper divisors? (enter a number)", n), strconv.Itoa(s),
			fmt.Sprintf("The proper divisors of %d add to %d.", n, s), f)
	case 2:
		kind := 3
		if w3IsPerfect(n) {
			kind = 1
		} else if w3IsAbundant(n) {
			kind = 2
		}
		w3IntFact(f, "class", kind)
		return w3P(kp, fmt.Sprintf("Classify %d: enter 1 for perfect, 2 for abundant, 3 for deficient.", n), strconv.Itoa(kind),
			fmt.Sprintf("The proper divisors sum to %d. Compared with %d, the number is %s.", s, n, []string{"", "perfect", "abundant", "deficient"}[kind]), f)
	default:
		return w3P(kp, fmt.Sprintf("What is \\(\\sigma(%d)\\), the sum of all positive divisors? (enter a number)", n), strconv.Itoa(sigma),
			fmt.Sprintf("The sum of all divisors is the aliquot sum plus the number itself: \\(%d + %d = %d\\).", s, n, sigma), f)
	}
}

func w3PerfectNumbers(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w3Small(10, 50)
		count := 0
		for m := 2; m <= n; m++ {
			if w3IsPerfect(m) {
				count++
			}
		}
		w3IntFact(f, "limit", n)
		return w3P(kp, fmt.Sprintf("How many perfect numbers are at most %d? (enter a number)", n), strconv.Itoa(count),
			fmt.Sprintf("Test each candidate up to %d by comparing its proper-divisor sum with the number. The count of perfect numbers is %d.", n, count), f)
	case 2:
		n := w3Small(20, 80)
		count := 0
		for m := 2; m <= n; m++ {
			if w3IsAbundant(m) {
				count++
			}
		}
		w3IntFact(f, "limit", n)
		return w3P(kp, fmt.Sprintf("How many abundant numbers are at most %d? (enter a number)", n), strconv.Itoa(count),
			fmt.Sprintf("Test each number up to %d by comparing its proper-divisor sum with the number. The count of abundant numbers is %d.", n, count), f)
	default:
		n := w3Small(10, 100)
		next := n + 1
		for !w3IsPerfect(next) {
			next++
		}
		w3IntFact(f, "after", n)
		return w3P(kp, fmt.Sprintf("What is the smallest perfect number greater than %d? (enter a number)", n), strconv.Itoa(next),
			fmt.Sprintf("The perfect numbers are 6, 28, 496, ... The smallest one greater than %d is %d.", n, next), f)
	}
}

func w3Mersenne(kp int) generator.Problem {
	n := w3Small(2, 13)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	mn := (1 << n) - 1
	w3IntFact(f, "mersenne", mn)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the Mersenne number \\(M_{%d} = 2^{%d}-1\\)? (enter a number)", n, n), strconv.Itoa(mn),
			fmt.Sprintf("\\(2^{%d}=%d\\), so \\(M_{%d}=%d\\).", n, 1<<n, n, mn), f)
	case 2:
		m := w3Small(3, 17)
		w3IntFact(f, "modulus", m)
		r := mn % m
		return w3P(kp, fmt.Sprintf("What is \\(M_{%d} \\bmod %d\\)? (enter a number)", n, m), strconv.Itoa(r),
			fmt.Sprintf("\\(M_{%d}=%d\\), and \\(%d \\bmod %d = %d\\).", n, mn, mn, m, r), f)
	default:
		is := 0
		if w3IsPrime(mn) {
			is = 1
		}
		w3IntFact(f, "prime", is)
		return w3P(kp, fmt.Sprintf("Is \\(M_{%d}=%d\\) prime? Enter 1 for yes and 0 for no.", n, mn), strconv.Itoa(is),
			fmt.Sprintf("A Mersenne number \\(2^n-1\\) can be prime only when \\(n\\) is prime, and not every prime exponent works. For \\(n=%d\\), the number %d is %s.", n, mn, map[int]string{0: "composite", 1: "prime"}[is]), f)
	}
}

func w3FermatNumbers(kp int) generator.Problem {
	// F_5 is the first composite Fermat number, so the range must include it.
	n := w3Small(0, 5)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	fn := (1 << (1 << n)) + 1
	w3IntFact(f, "fermat", fn)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the Fermat number \\(F_{%d}=2^{2^{%d}}+1\\)? (enter a number)", n, n), strconv.Itoa(fn),
			fmt.Sprintf("\\(2^{2^{%d}}=2^{%d}=%d\\), so \\(F_{%d}=%d\\).", n, 1<<n, 1<<(1<<n), n, fn), f)
	case 2:
		m := w3Small(3, 31)
		w3IntFact(f, "modulus", m)
		r := fn % m
		return w3P(kp, fmt.Sprintf("What is \\(F_{%d} \\bmod %d\\)? (enter a number)", n, m), strconv.Itoa(r),
			fmt.Sprintf("\\(F_{%d}=%d\\), and \\(%d \\bmod %d = %d\\).", n, fn, fn, m, r), f)
	default:
		is := 0
		if w3IsPrime(fn) {
			is = 1
		}
		w3IntFact(f, "prime", is)
		return w3P(kp, fmt.Sprintf("Is \\(F_{%d}=%d\\) prime? Enter 1 for yes and 0 for no.", n, fn), strconv.Itoa(is),
			fmt.Sprintf("The first five Fermat numbers \\(F_0\\) through \\(F_4\\) are prime, while \\(F_5\\) is composite. Here \\(F_{%d}=%d\\) is %s.", n, fn, map[int]string{0: "composite", 1: "prime"}[is]), f)
	}
}

func w3Repunit(kp int) generator.Problem {
	n := w3Small(2, 9)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	r := 0
	for i := 0; i < n; i++ {
		r = r*10 + 1
	}
	w3IntFact(f, "repunit", r)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the repunit \\(R_{%d}\\) with %d digits? (enter a number)", n, n), strconv.Itoa(r),
			fmt.Sprintf("\\(R_{%d}\\) is the number made of %d ones, namely %d.", n, n, r), f)
	case 2:
		m := w3Small(3, 19)
		w3IntFact(f, "modulus", m)
		return w3P(kp, fmt.Sprintf("What is \\(R_{%d} \\bmod %d\\)? (enter a number)", n, m), strconv.Itoa(r%m),
			fmt.Sprintf("\\(R_{%d}=%d\\), so the residue modulo %d is %d.", n, r, m, r%m), f)
	default:
		return w3P(kp, fmt.Sprintf("What is the sum of the decimal digits of \\(R_{%d}\\)? (enter a number)", n), strconv.Itoa(n),
			fmt.Sprintf("A repunit with %d digits has %d ones, so its digit sum is %d.", n, n, n), f)
	}
}

func w3PrimeCounting(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w3Small(10, 100)
		ps := w3PrimeList(n)
		w3IntFact(f, "limit", n)
		return w3P(kp, fmt.Sprintf("What is \\(\\pi(%d)\\), the number of primes at most %d? (enter a number)", n, n), strconv.Itoa(len(ps)),
			fmt.Sprintf("List the primes up to %d and count them: there are %d.", n, len(ps)), f)
	case 2:
		k := w3Small(3, 25)
		ps := w3PrimeList(200)
		w3IntFact(f, "k", k)
		return w3P(kp, fmt.Sprintf("What is the %dth prime? (enter a number)", k), strconv.Itoa(ps[k-1]),
			fmt.Sprintf("Counting the primes in order, the %dth prime is %d.", k, ps[k-1]), f)
	default:
		a := w3Small(10, 40)
		b := a + w3Small(10, 40)
		count := 0
		for n := a; n <= b; n++ {
			if w3IsPrime(n) {
				count++
			}
		}
		w3IntFact(f, "lo", a)
		w3IntFact(f, "hi", b)
		return w3P(kp, fmt.Sprintf("How many primes lie between %d and %d inclusive? (enter a number)", a, b), strconv.Itoa(count),
			fmt.Sprintf("Count the primes in the interval \\([%d,%d]\\): there are %d.", a, b, count), f)
	}
}

func w3TwinPrimes(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w3Small(20, 120)
		count := 0
		for p := 2; p+2 <= n; p++ {
			if w3IsPrime(p) && w3IsPrime(p+2) {
				count++
			}
		}
		w3IntFact(f, "limit", n)
		return w3P(kp, fmt.Sprintf("How many twin-prime pairs \\((p,p+2)\\) have both members at most %d? (enter a number)", n), strconv.Itoa(count),
			fmt.Sprintf("Scan the primes up to %d and count the pairs separated by 2: there are %d.", n, count), f)
	case 2:
		n := w3Small(5, 80)
		p := n + 1
		for !(w3IsPrime(p) && w3IsPrime(p+2)) {
			p++
		}
		w3IntFact(f, "after", n)
		return w3P(kp, fmt.Sprintf("What is the smaller member of the first twin-prime pair whose smaller member is greater than %d? (enter a number)", n), strconv.Itoa(p),
			fmt.Sprintf("Search upward from %d for a prime \\(p\\) with \\(p+2\\) also prime. The first such \\(p\\) is %d.", n+1, p), f)
	default:
		p := w3Small(3, 80)
		twin := 0
		if w3IsPrime(p) && w3IsPrime(p+2) {
			twin = 1
		}
		w3IntFact(f, "p", p)
		return w3P(kp, fmt.Sprintf("Is %d the smaller member of a twin-prime pair? Enter 1 for yes and 0 for no.", p), strconv.Itoa(twin),
			fmt.Sprintf("A twin-prime pair needs both %d and %d prime. The answer is %d.", p, p+2, twin), f)
	}
}

func w3Goldbach(kp int) generator.Problem {
	n := 2 * w3Small(5, 40)
	f := map[string]string{}
	w3IntFact(f, "even", n)
	var pairs [][2]int
	for p := 2; p <= n/2; p++ {
		if w3IsPrime(p) && w3IsPrime(n-p) {
			pairs = append(pairs, [2]int{p, n - p})
		}
	}
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("How many unordered Goldbach partitions does %d have? (enter a number)", n), strconv.Itoa(len(pairs)),
			fmt.Sprintf("Count the ways to write %d as a sum of two primes, with order ignored: there are %d.", n, len(pairs)), f)
	case 2:
		if len(pairs) == 0 {
			return w3P(kp, fmt.Sprintf("What is the smallest prime in a Goldbach partition of %d? (enter a number)", n), "0",
				fmt.Sprintf("No partition was found, which should not happen for this even number."), f)
		}
		return w3P(kp, fmt.Sprintf("What is the smallest prime in a Goldbach partition of %d? (enter a number)", n), strconv.Itoa(pairs[0][0]),
			fmt.Sprintf("The smallest prime that works is %d, because %d is also prime.", pairs[0][0], pairs[0][1]), f)
	default:
		if len(pairs) == 0 {
			return w3P(kp, fmt.Sprintf("What is the largest prime in a Goldbach partition of %d? (enter a number)", n), "0",
				fmt.Sprintf("No partition was found, which should not happen for this even number."), f)
		}
		return w3P(kp, fmt.Sprintf("What is the largest prime in a Goldbach partition of %d? (enter a number)", n), strconv.Itoa(pairs[len(pairs)-1][1]),
			fmt.Sprintf("The largest prime in any Goldbach partition is %d, paired with %d.", pairs[len(pairs)-1][1], pairs[len(pairs)-1][0]), f)
	}
}

func w3Polygonal(kp int) generator.Problem {
	s := w3Small(3, 9)
	n := w3Small(3, 15)
	f := map[string]string{}
	w3IntFact(f, "sides", s)
	w3IntFact(f, "n", n)
	p := ((s-2)*n*n - (s-4)*n) / 2
	w3IntFact(f, "polygonal", p)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the %dth %d-gonal number? (enter a number)", n, s), strconv.Itoa(p),
			fmt.Sprintf("Use \\(P(s,n)=((s-2)n^2-(s-4)n)/2\\) with \\(s=%d\\) and \\(n=%d\\): the value is %d.", s, n, p), f)
	case 2:
		k := 1
		for ((s-2)*k*k-(s-4)*k)/2 < p {
			k++
		}
		w3IntFact(f, "index", k)
		return w3P(kp, fmt.Sprintf("The %d-gonal number %d is which index? (enter a number)", s, p), strconv.Itoa(k),
			fmt.Sprintf("Solve \\(((%d-2)n^2-(%d-4)n)/2=%d\\). The positive integer solution is \\(n=%d\\).", s, s, p, k), f)
	default:
		diff := (s-2)*n + 1
		w3IntFact(f, "difference", diff)
		return w3P(kp, fmt.Sprintf("For %d-gonal numbers, what is \\(P(%d,%d)-P(%d,%d)\\)? (enter a number)", s, s, n+1, s, n), strconv.Itoa(diff),
			fmt.Sprintf("Consecutive polygonal numbers differ by \\((s-2)n+1\\). With \\(s=%d\\) and \\(n=%d\\), the difference is %d.", s, n, diff), f)
	}
}

func w3Pentagonal(kp int) generator.Problem {
	n := w3Small(3, 20)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	p := n * (3*n - 1) / 2
	w3IntFact(f, "pentagonal", p)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the %dth pentagonal number \\(P_{%d}\\)? (enter a number)", n, n), strconv.Itoa(p),
			fmt.Sprintf("\\(P_n=n(3n-1)/2\\). With \\(n=%d\\), the value is %d.", n, p), f)
	case 2:
		k := 1
		for k*(3*k-1)/2 < p {
			k++
		}
		w3IntFact(f, "index", k)
		return w3P(kp, fmt.Sprintf("The pentagonal number %d has which index? (enter a number)", p), strconv.Itoa(k),
			fmt.Sprintf("Solve \\(n(3n-1)/2=%d\\). The positive integer solution is \\(n=%d\\).", p, k), f)
	default:
		count := 0
		for m := 1; m*(3*m-1)/2 <= p; m++ {
			count++
		}
		w3IntFact(f, "count", count)
		return w3P(kp, fmt.Sprintf("How many pentagonal numbers are at most %d? (enter a number)", p), strconv.Itoa(count),
			fmt.Sprintf("The pentagonal numbers grow quadratically. The ones at most %d are the first %d.", p, count), f)
	}
}

func w3Egyptian(kp int) generator.Problem {
	b := w3Small(5, 24)
	a := w3Small(2, b-1)
	for computeGCD(a, b) != 1 {
		a = w3Small(2, b-1)
	}
	f := map[string]string{}
	w3IntFact(f, "numerator", a)
	w3IntFact(f, "denominator", b)
	// Greedy Egyptian expansion: ceil(b/a).
	d1 := (b + a - 1) / a
	w3IntFact(f, "first_denominator", d1)
	// Remainder a/b - 1/d1 = (a*d1-b)/(b*d1), then reduce.
	rn := a*d1 - b
	rd := b * d1
	g := computeGCD(rn, rd)
	rn /= g
	rd /= g
	d2 := 0
	if rn > 0 {
		d2 = (rd + rn - 1) / rn
	}
	w3IntFact(f, "second_denominator", d2)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("For \\(%d/%d\\), what is the first denominator produced by the greedy Egyptian-fraction algorithm? (enter a number)", a, b), strconv.Itoa(d1),
			fmt.Sprintf("The first denominator is \\(\\lceil b/a \\rceil = \\lceil %d/%d \\rceil = %d\\).", b, a, d1), f)
	case 2:
		return w3P(kp, fmt.Sprintf("After the first greedy step for \\(%d/%d\\), what is the next denominator? (enter a number)", a, b), strconv.Itoa(d2),
			fmt.Sprintf("Subtract \\(1/%d\\) from \\(%d/%d\\), reduce, and take the ceiling of the reciprocal. The next denominator is %d.", d1, a, b, d2), f)
	default:
		n := w3Small(2, 15)
		w3IntFact(f, "identity_n", n)
		second := n * (n + 1)
		w3IntFact(f, "identity_second", second)
		return w3P(kp, fmt.Sprintf("Using \\(1/n = 1/(n+1) + 1/(n(n+1))\\), what is the larger denominator when \\(n=%d\\)? (enter a number)", n), strconv.Itoa(second),
			fmt.Sprintf("The identity gives \\(1/%d = 1/%d + 1/%d\\), so the larger denominator is %d.", n, n+1, second, second), f)
	}
}

func w3HighlyComposite(kp int) generator.Problem {
	n := w3Small(2, 240)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("How many positive divisors does %d have? (enter a number)", n), strconv.Itoa(w3DivisorCount(n)),
			fmt.Sprintf("Factor %d and multiply one more than each exponent. The divisor count is %d.", n, w3DivisorCount(n)), f)
	case 2:
		is := 0
		if w3IsHighlyComposite(n) {
			is = 1
		}
		w3IntFact(f, "highly_composite", is)
		return w3P(kp, fmt.Sprintf("Is %d highly composite? Enter 1 for yes and 0 for no.", n), strconv.Itoa(is),
			fmt.Sprintf("A number is highly composite when it has more divisors than every smaller positive integer. The answer for %d is %d.", n, is), f)
	default:
		next := n + 1
		for !w3IsHighlyComposite(next) {
			next++
		}
		w3IntFact(f, "next", next)
		return w3P(kp, fmt.Sprintf("What is the smallest highly composite number greater than %d? (enter a number)", n), strconv.Itoa(next),
			fmt.Sprintf("Search upward from %d for the first number whose divisor count exceeds every smaller number's. It is %d.", n+1, next), f)
	}
}

func w3Amicable(kp int) generator.Problem {
	pairs := [][2]int{{220, 284}, {1184, 1210}, {2620, 2924}, {5020, 5564}, {6232, 6368}}
	pair := pairs[rand.Intn(len(pairs))]
	n, partner := pair[0], pair[1]
	if rand.Intn(2) == 0 {
		n, partner = partner, n
	}
	f := map[string]string{}
	w3IntFact(f, "n", n)
	s := w3ProperDivisorSum(n)
	w3IntFact(f, "aliquot_sum", s)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the aliquot sum \\(s(%d)\\)? (enter a number)", n), strconv.Itoa(s),
			fmt.Sprintf("The proper divisors of %d sum to %d.", n, s), f)
	case 2:
		is := 0
		if w3IsAmicable(n) {
			is = 1
		}
		w3IntFact(f, "amicable", is)
		return w3P(kp, fmt.Sprintf("Is %d part of an amicable pair? Enter 1 for yes and 0 for no.", n), strconv.Itoa(is),
			fmt.Sprintf("Apply the aliquot sum twice. For %d, the second application returns %d, so the answer is %d.", n, n, is), f)
	default:
		return w3P(kp, fmt.Sprintf("What is the amicable partner of %d? (enter a number, or 0 if it has none)", n), strconv.Itoa(partner),
			fmt.Sprintf("The aliquot sum of %d is %d, and the aliquot sum of %d is %d, so the partner is %d.", n, s, s, n, partner), f)
	}
}

func w3Happy(kp int) generator.Problem {
	n := w3Small(10, 200)
	f := map[string]string{}
	w3IntFact(f, "n", n)
	sumSq := 0
	for m := n; m > 0; m /= 10 {
		d := m % 10
		sumSq += d * d
	}
	w3IntFact(f, "sum_squares", sumSq)
	switch kp {
	case 1:
		return w3P(kp, fmt.Sprintf("What is the sum of the squares of the digits of %d? (enter a number)", n), strconv.Itoa(sumSq),
			fmt.Sprintf("Square each digit of %d and add: the sum is %d.", n, sumSq), f)
	case 2:
		is := 0
		if w3IsHappy(n) {
			is = 1
		}
		w3IntFact(f, "happy", is)
		return w3P(kp, fmt.Sprintf("Is %d happy? Enter 1 for yes and 0 for no.", n), strconv.Itoa(is),
			fmt.Sprintf("Iterate the sum-of-squared-digits map from %d. It reaches 1 if and only if the number is happy. The answer is %d.", n, is), f)
	default:
		next := n
		for !w3IsHappy(next) {
			next++
		}
		w3IntFact(f, "next_happy", next)
		return w3P(kp, fmt.Sprintf("What is the smallest happy number at least %d? (enter a number)", n), strconv.Itoa(next),
			fmt.Sprintf("Search upward from %d. The first happy number is %d.", n, next), f)
	}
}

func w3ConsecutiveSums(kp int) generator.Problem {
	n := w3Small(9, 200)
	if kp != 3 {
		for w3IsPrimePowerOfTwo(n) {
			n++
		}
	}
	f := map[string]string{}
	w3IntFact(f, "n", n)
	switch kp {
	case 1:
		c := w3CountConsecutive(n)
		return w3P(kp, fmt.Sprintf("In how many ways can %d be written as a sum of two or more consecutive positive integers? (enter a number)", n), strconv.Itoa(c),
			fmt.Sprintf("The number of such representations equals the number of odd divisors of %d greater than 1. Here that count is %d.", n, c), f)
	case 2:
		k := w3ShortestConsecutive(n)
		return w3P(kp, fmt.Sprintf("What is the shortest length of a sum of two or more consecutive positive integers equal to %d? (enter a number)", n), strconv.Itoa(k),
			fmt.Sprintf("Try lengths \\(k=2,3,\\ldots\\). The first length for which \\(n=k(2a+k-1)/2\\) has an integer \\(a\\ge1\\) is %d.", k), f)
	default:
		is := 0
		if !w3IsPrimePowerOfTwo(n) {
			is = 1
		}
		w3IntFact(f, "representable", is)
		return w3P(kp, fmt.Sprintf("Can %d be written as a sum of two or more consecutive positive integers? Enter 1 for yes and 0 for no.", n), strconv.Itoa(is),
			fmt.Sprintf("A positive integer has such a representation exactly when it is not a power of two. For %d the answer is %d.", n, is), f)
	}
}

// w3FormatInts is used by the tests and by a few explanations to keep list
// rendering consistent.
func w3FormatInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ", ")
}
