package numtheory

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// These concepts are graded multiple choice, so every shape carries options.
// The conceptual rows keep the theorem statements the concept is about, and the
// computed rows ask for a small quantity the learner must actually work out.

func advChoice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

type chebyshevGen struct{}

func (g *chebyshevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		x := 10 + rand.Intn(51)
		count := primePowerCount(x)
		return advChoice(
			fmt.Sprintf("How many prime powers \\(n \\le %d\\) are there, counting each prime and each power separately?", x),
			strconv.Itoa(count),
			fmt.Sprintf("Listing the primes up to %d together with the prime powers that are not themselves prime gives %d values, which is the count \\psi weights with log p.", x, count),
			strconv.Itoa(count-1), strconv.Itoa(count+1), strconv.Itoa(count+2),
		)
	}
	rows := []struct {
		q, a, e string
	}{
		{"Is \\(\\psi(x)=\\sum_{n\\le x}\\Lambda(n)\\) the Chebyshev function?", "yes", "The von Mangoldt weight gives the Chebyshev function psi."},
		{"Is \\(\\psi(x)\\) the same function as \\(\\pi(x)\\)?", "no", "Psi weights prime powers by log p while pi simply counts primes, so they are different functions."},
		{"Is the prime number theorem equivalent to \\(\\psi(x)\\sim x\\)?", "yes", "This is the Chebyshev form of the prime number theorem."},
		{"Is \\(\\theta(x)=\\sum_{p\\le x}\\log p\\) equivalent to the prime number theorem when \\(\\theta(x)\\sim x\\)?", "yes", "The two forms differ only by the contribution of higher prime powers, which is negligible."},
	}
	r := rows[rand.Intn(len(rows))]
	return advChoice(r.q, r.a, r.e, otherAnswer(r.a))
}

type continuedFractionsGen struct{}

func (g *continuedFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		a0 := 1 + rand.Intn(4)
		a1 := 1 + rand.Intn(4)
		a2 := 1 + rand.Intn(4)
		num := a0*(a1*a2+1) + a2
		den := a1*a2 + 1
		answer := strconv.Itoa(num)
		if den != 1 {
			answer = fmt.Sprintf("%d/%d", num, den)
		}
		wrong1 := fmt.Sprintf("%d/%d", num+1, den)
		wrong2 := fmt.Sprintf("%d/%d", num, den+1)
		wrong3 := fmt.Sprintf("%d/%d", a0*a1+a2, a1)
		return advChoice(
			fmt.Sprintf("What rational number does the continued fraction \\([%d;%d,%d]\\) equal? (type like 3/2)", a0, a1, a2),
			answer,
			fmt.Sprintf("Work from the tail: \\(1/%d\\) with %d gives \\(%d/%d\\), then add %d to get %d/%d.", a2, a1, a1*a2+1, a2, a0, num, den),
			wrong1, wrong2, wrong3,
		)
	}
	rows := []struct {
		q, a, e string
	}{
		{"Do the convergents of a continued fraction give best rational approximations?", "yes", "Each convergent is the best rational approximation with a denominator that small."},
		{"Is the continued fraction of every real number finite?", "no", "Only rationals terminate; irrationals have infinite expansions."},
		{"Is \\(\\sqrt{2}=[1;2,2,2,\\dots]\\)?", "yes", "Quadratic irrationals have periodic continued fractions, and the period here is 2."},
		{"Does Lagrange's theorem say periodic continued fractions are exactly the quadratic irrationals?", "yes", "Lagrange proved both directions of that equivalence."},
	}
	r := rows[rand.Intn(len(rows))]
	return advChoice(r.q, r.a, r.e, otherAnswer(r.a))
}

type dirichletGen struct{}

func (g *dirichletGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		moduli := []int{4, 5, 6, 7}
		m := moduli[rand.Intn(len(moduli))]
		a := 1 + rand.Intn(m-1)
		if gcdInt(a, m) != 1 {
			a = 1
		}
		p := smallestPrimeInAP(a, m)
		return advChoice(
			fmt.Sprintf("What is the smallest prime congruent to %d modulo %d?", a, m),
			strconv.Itoa(p),
			fmt.Sprintf("Test the numbers %d, %d, %d, %d and so on for divisibility; the first prime is %d. Dirichlet guarantees there are infinitely many more.", a, a+m, a+2*m, a+3*m, p),
			strconv.Itoa(p+m), strconv.Itoa(p+1), strconv.Itoa(p-1),
		)
	}
	rows := []struct {
		q, a, e string
	}{
		{"Does Dirichlet's theorem give infinitely many primes in every progression \\(a, a+m, a+2m,\\dots\\) with \\(\\gcd(a,m)=1\\)?", "yes", "That coprimality is the hypothesis, and it is exactly what rules out a common divisor."},
		{"Is the progression \\(2,4,6,\\dots\\) expected to contain infinitely many primes?", "no", "Every term beyond 2 is even and greater than 2, so it is composite."},
		{"Does Dirichlet's theorem require \\(\\gcd(a,m)=1\\)?", "yes", "Without it the progression shares a divisor and can contain at most one prime."},
		{"Is the non-vanishing of \\(L(1,\\chi)\\) the key step in Dirichlet's proof?", "yes", "The theorem follows once the Dirichlet L-series at 1 is known to be nonzero."},
	}
	r := rows[rand.Intn(len(rows))]
	return advChoice(r.q, r.a, r.e, otherAnswer(r.a))
}

type ellipticCurveGen struct{}

func (g *ellipticCurveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(2) == 0 {
		a := rand.Intn(7) - 3
		b := 1 + rand.Intn(5)
		disc := 4*a*a*a + 27*b*b
		if disc == 0 {
			b++
			disc = 4*a*a*a + 27*b*b
		}
		return advChoice(
			fmt.Sprintf("For the curve \\(y^{2}=x^{3}%+dx%+d\\), what is \\(4a^{3}+27b^{2}\\)?", a, b),
			strconv.Itoa(disc),
			fmt.Sprintf("Here a is %d and b is %d, so \\(4a^{3}=%d\\) and \\(27b^{2}=%d\\), giving %d. It is nonzero, so the cubic is nonsingular.", a, b, 4*a*a*a, 27*b*b, disc),
			strconv.Itoa(disc+1), strconv.Itoa(disc-1), strconv.Itoa(4*a*a*a-27*b*b),
		)
	}
	rows := []struct {
		q, a, e string
	}{
		{"Is the curve \\(y^{2}=x^{3}+ax+b\\) elliptic exactly when the cubic has no repeated root?", "yes", "A nonzero discriminant means the cubic is nonsingular, which is the defining condition."},
		{"Is \\(y^{2}=x^{3}\\) an elliptic curve?", "no", "Its discriminant is zero: the cubic has a cusp, so it is singular."},
		{"Do the points of an elliptic curve form an abelian group under the chord and tangent law?", "yes", "The chord and tangent construction gives a group law with the point at infinity as identity."},
		{"Does the Mordell and Weil theorem say \\(E(\\mathbb{Q})\\) is finitely generated?", "yes", "It splits as a free part of finite rank and a torsion subgroup."},
	}
	r := rows[rand.Intn(len(rows))]
	return advChoice(r.q, r.a, r.e, otherAnswer(r.a))
}

func otherAnswer(a string) string {
	if a == "yes" {
		return "no"
	}
	return "yes"
}

func primePowerCount(x int) int {
	count := 0
	for n := 2; n <= x; n++ {
		if isPrimePower(n) {
			count++
		}
	}
	return count
}

func isPrimePower(n int) bool {
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			m := n
			for m%p == 0 {
				m /= p
			}
			return m == 1
		}
	}
	return n >= 2
}

func smallestPrimeInAP(a, m int) int {
	for n := a; ; n += m {
		if n < 2 {
			continue
		}
		if isPrimePower(n) && isPrimeNumber(n) {
			return n
		}
	}
}

func isPrimeNumber(n int) bool {
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
