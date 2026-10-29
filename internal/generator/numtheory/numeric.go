package numtheory

// The numeric number-theory concepts the domain was missing: the divisor count and divisor sum, the
// exponent of a prime in a factorial, solving a linear congruence, the hypotenuse of a primitive
// Pythagorean triple, the number of quadratic residues, and reading a numeral in another base.
//
// The domain had twenty-five concepts and most of them `multiple_choice`, which is the
// recognition-shaped surface section 4.2 of docs/content-expansion.md measures. These seven are
// computed, and the two that already existed as multiple-choice questions -- divisibility and the
// Legendre symbol -- are extended rather than repeated: a divisibility question is "does it divide",
// and the divisor count is the arithmetic that follows from a factorisation.
//
// Every answer is an integer by construction: the numbers are built from known prime powers, the
// congruence is solved with a unit coefficient so the solution is unique, the triple is primitive by
// construction, and the quadratic-residue count is half of one less than the modulus.
import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerNumeric(reg *generator.Registry) {
	reg.Register("nt.divisor.count", &divisorCountGen{})
	reg.Register("nt.divisor.sum", &divisorSumGen{})
	reg.Register("nt.prime.exponent", &primeExponentGen{})
	reg.Register("nt.congruence.linear", &linearCongruenceGen{})
	reg.Register("nt.pythagorean.hypotenuse", &pythagoreanGen{})
	reg.Register("nt.residue.count", &residueCountGen{})
	reg.Register("nt.base.convert", &baseConvertGen{})
}

func smallInt(lo, hi int) int {
	if hi < lo {
		hi = lo
	}
	return lo + rand.Intn(hi-lo+1)
}

// primePower is one term of a factorisation.
type primePower struct {
	p, e int
}

// factorisation builds a number from small prime powers and returns it with its factorisation.
func factorisation(ctx generator.GeneratorContext) (int, []primePower) {
	primes := []int{2, 3, 5, 7}
	n := 1
	var fs []primePower
	for _, p := range primes {
		if rand.Intn(3) == 0 {
			continue
		}
		e := 1 + rand.Intn(1+int(ctx.Difficulty*3))
		n *= powInt(p, e)
		fs = append(fs, primePower{p: p, e: e})
	}
	if len(fs) == 0 {
		n = 4
		fs = []primePower{{p: 2, e: 2}}
	}
	return n, fs
}

func powInt(b, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= b
	}
	return r
}

// factorString renders a factorisation for the explanation. A builder rather than repeated
// concatenation: make lint-go rejects a quadratic string build (AIL071), and this loop is one.
func factorString(fs []primePower) string {
	var b strings.Builder
	for i, f := range fs {
		if i > 0 {
			b.WriteString(" * ")
		}
		fmt.Fprintf(&b, "%d^%d", f.p, f.e)
	}
	return b.String()
}

// divisorCountGen counts divisors from the exponents: the count is the product of one more than each.
type divisorCountGen struct{}

func (g *divisorCountGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n, fs := factorisation(ctx)
	count := 1
	for _, f := range fs {
		count *= f.e + 1
	}
	return generator.Problem{
		Question: fmt.Sprintf("How many positive divisors does %d have?", n),
		Answer:   fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Factorise first: %d = %s. Every divisor is built by choosing an exponent for each prime, so the choices multiply: each exponent here offers one more option than its size, and the product is %d.",
			n, factorString(fs), count),
		Facts: map[string]string{"count": fmt.Sprintf("%d", count), "n": fmt.Sprintf("%d", n)},
	}
}

// divisorSumGen sums the divisors from the same factorisation, each prime contributing a geometric
// series.
type divisorSumGen struct{}

func (g *divisorSumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n, fs := factorisation(ctx)
	sum := 1
	for _, f := range fs {
		// 1 + p + ... + p^e
		term := 0
		for i := 0; i <= f.e; i++ {
			term += powInt(f.p, i)
		}
		sum *= term
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is the sum of the positive divisors of %d?", n),
		Answer:   fmt.Sprintf("%d", sum),
		Explanation: fmt.Sprintf("With %d = %s, each prime contributes the series 1 + p + ... + p^e, and the series multiply: their product is %d. Choosing one term from each series is the same as choosing one divisor.",
			n, factorString(fs), sum),
		Facts: map[string]string{"sum": fmt.Sprintf("%d", sum), "n": fmt.Sprintf("%d", n)},
	}
}

// primeExponentGen applies Legendre's formula, which sums the floors of n divided by the powers of p.
type primeExponentGen struct{}

func (g *primeExponentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	p := []int{2, 3, 5, 7}[rand.Intn(4)]
	n := 5 + rand.Intn(40+int(ctx.Difficulty*120))
	exp := 0
	for q := p; q <= n; q *= p {
		exp += n / q
	}
	return generator.Problem{
		Question: fmt.Sprintf("What is the exponent of %d in the prime factorisation of %d! ?", p, n),
		Answer:   fmt.Sprintf("%d", exp),
		Explanation: fmt.Sprintf("Count how many multiples of each power of %d are in %d!: the floor of %d/%d, plus the floor of %d/%d^2, and so on until the power exceeds %d. The total is %d, and the floors are counted because a multiple of p^2 contributes two factors of p.",
			p, n, n, p, n, p, n, exp),
		Facts: map[string]string{"exponent": fmt.Sprintf("%d", exp), "prime": fmt.Sprintf("%d", p), "n": fmt.Sprintf("%d", n)},
	}
}

// linearCongruenceGen uses a unit coefficient, so the solution is unique modulo the modulus.
type linearCongruenceGen struct{}

func (g *linearCongruenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 5 + rand.Intn(10+int(ctx.Difficulty*20))
	a := smallInt(2, m-1)
	for gcdInt(a, m) != 1 {
		a = smallInt(2, m-1)
	}
	x := smallInt(0, m-1)
	b := (a * x) % m
	return generator.Problem{
		Question: fmt.Sprintf("Solve the congruence \\\\(%dx \\equiv %d \\pmod{%d}\\\\), giving the least non-negative solution.", a, b, m),
		Answer:   fmt.Sprintf("%d", x),
		Explanation: fmt.Sprintf("The coefficient %d is coprime to %d, so it has an inverse modulo %d and the solution is unique. Multiplying both sides by that inverse gives x ≡ %d (mod %d), and that residue is the least non-negative one.",
			a, m, m, x, m),
		Facts: map[string]string{"solution": fmt.Sprintf("%d", x), "modulus": fmt.Sprintf("%d", m)},
	}
}

// pythagoreanGen builds a primitive triple from Euclid's formula.
type pythagoreanGen struct{}

func (g *pythagoreanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 2 + rand.Intn(4+int(ctx.Difficulty*5))
	n := 1 + rand.Intn(m-1)
	for (m-n)%2 == 0 || gcdInt(m, n) != 1 {
		n = 1 + rand.Intn(m-1)
		if m-n == 0 {
			m++
		}
	}
	odd := m*m - n*n
	even := 2 * m * n
	hyp := m*m + n*n
	return generator.Problem{
		Question: fmt.Sprintf("Euclid's formula gives a primitive Pythagorean triple from %d and %d. What is its hypotenuse?", m, n),
		Answer:   fmt.Sprintf("%d", hyp),
		Explanation: fmt.Sprintf("The formula is (m^2 - n^2, 2mn, m^2 + n^2): here that is (%d, %d, %d), and %d^2 + %d^2 = %d^2. The conditions on m and n -- coprime and of opposite parity -- are what make the triple primitive rather than a multiple of a smaller one.",
			odd, even, hyp, odd, even, hyp),
		Facts: map[string]string{"hypotenuse": fmt.Sprintf("%d", hyp), "leg_odd": fmt.Sprintf("%d", odd), "leg_even": fmt.Sprintf("%d", even)},
	}
}

// residueCountGen uses an odd prime, where exactly half of the non-zero residues are squares.
type residueCountGen struct{}

func (g *residueCountGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{7, 11, 13, 17, 19, 23, 29, 31}
	p := primes[rand.Intn(len(primes))]
	count := (p - 1) / 2
	return generator.Problem{
		Question: fmt.Sprintf("How many non-zero quadratic residues are there modulo %d?", p),
		Answer:   fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Squaring pairs the non-zero residues: k and -k have the same square, so the %d non-zero residues fall into %d pairs and each pair gives one square. That is why exactly half the non-zero residues are quadratic, and why the Legendre symbol is defined as it is.",
			p, count),
		Facts: map[string]string{"count": fmt.Sprintf("%d", count), "prime": fmt.Sprintf("%d", p)},
	}
}

// baseConvertGen asks for the value of a numeral in another base.
type baseConvertGen struct{}

func (g *baseConvertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	base := 2 + rand.Intn(8)
	digits := 2 + rand.Intn(3+int(ctx.Difficulty*3))
	val := 0
	var b strings.Builder
	for i := 0; i < digits; i++ {
		d := rand.Intn(base)
		if i == 0 && d == 0 {
			d = 1
		}
		val = val*base + d
		fmt.Fprintf(&b, "%d", d)
	}
	numeral := b.String()
	return generator.Problem{
		Question: fmt.Sprintf("What is the value of the numeral %s in base %d, written in base ten?", numeral, base),
		Answer:   fmt.Sprintf("%d", val),
		Explanation: fmt.Sprintf("The numeral %s in base %d expands as %s = %d. Each digit is a coefficient on a power of the base, the same place-value rule as base ten, with powers of %d instead of powers of ten.",
			numeral, base, baseConvertExpansion(numeral, base), val, base),
		Facts: map[string]string{"value": fmt.Sprintf("%d", val), "base": fmt.Sprintf("%d", base)},
	}
}

// baseConvertExpansion spells out the place-value sum, so the explanation
// shows the arithmetic rather than only describing it.
func baseConvertExpansion(numeral string, base int) string {
	var b strings.Builder
	for i, ch := range numeral {
		if i > 0 {
			b.WriteString(" + ")
		}
		digit := int(ch - '0')
		place := 1
		for j := 0; j < len(numeral)-1-i; j++ {
			place *= base
		}
		fmt.Fprintf(&b, "%d×%d", digit, place)
	}
	return b.String()
}

// gcdInt is the Euclidean algorithm, used for the coprimality conditions.
func gcdInt(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
