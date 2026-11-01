package numtheory

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("nt.basics.divisibility", &divisibilityGen{})
	reg.Register("nt.basics.gcd_euclidean", &gcdEuclideanGen{})
	reg.Register("nt.basics.modular", &modularGen{})
	reg.Register("nt.basics.congruence", &congruenceGen{})
	reg.Register("nt.adv.fermat_little", &fermatLittleGen{})
	reg.Register("nt.adv.euler_phi", &eulerPhiGen{})
	reg.Register("nt.adv.diophantine", &diophantineGen{})
	reg.Register("nt.app.crypto", &cryptoGen{})
	reg.Register("nt.basics.bezout", &bezoutGen{})
	reg.Register("nt.basics.crt", &crtGen{})
	reg.Register("nt.adv.wilson", &wilsonGen{})
	reg.Register("nt.basics.prime_inf", &primeInfGen{})
	reg.Register("nt.adv.quadratic_reciprocity", &quadraticReciprocityGen{})
	reg.Register("nt.adv.primitive_root", &primitiveRootGen{})
	reg.Register("nt.analytic.prime_number_theorem", &pntGen{})
	reg.Register("nt.adv.dirichlet", &dirichletGen{})
	reg.Register("nt.adv.sieve", &sieveGen{})
	reg.Register("nt.analytic.zeta", &zetaGen{})
	reg.Register("nt.adv.chebyshev", &chebyshevGen{})
	reg.Register("nt.adv.elliptic_curve", &ellipticCurveGen{})
	reg.Register("nt.adv.continued_fractions", &continuedFractionsGen{})
	reg.Register("nt.adv.class_number", &classNumberGen{})
	registerNumeric(reg)
	registerBatch(reg)
	registerConversions(reg)
}

type divisibilityGen struct{}

func (g *divisibilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	d := 2 + rand.Intn(scale*10-1)
	q := 2 + rand.Intn(scale*12-1)
	a := d * q
	return generator.Problem{
		Question:    fmt.Sprintf("What is %d / %d? (enter a number)", a, d),
		Answer:      fmt.Sprintf("%d", q),
		Explanation: fmt.Sprintf("%d / %d = %d because %d = %d x %d.", a, d, q, a, d, q),
		Facts:       map[string]string{"a": fmt.Sprintf("%d", a), "d": fmt.Sprintf("%d", d), "quotient": fmt.Sprintf("%d", q)},
	}
}

type gcdEuclideanGen struct{}

func (g *gcdEuclideanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*37)) + 12
	b := rand.Intn(max(1, scale*37)) + 12
	gcd := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the GCD of %d and %d using the Euclidean algorithm.", a, b),
		Answer:      fmt.Sprintf("%d", gcd),
		Explanation: fmt.Sprintf("GCD(%d, %d) = %d (computed via Euclidean algorithm).", a, b, gcd),
	}
}

type modularGen struct{}

func (g *modularGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*30)) + 5
	m := rand.Intn(max(1, scale*16)) + 2
	ans := a % m
	return generator.Problem{
		Question:    fmt.Sprintf("Compute \\(%d \\bmod %d\\).", a, m),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("%d mod %d = %d because %d = %d x %d + %d.", a, m, ans, a, m, a/m, ans),
	}
}

type congruenceGen struct{}

func (g *congruenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	m := 2 + rand.Intn(29)
	a := rand.Intn(500)
	r := ((a % m) + m) % m
	return generator.Problem{
		Question:    fmt.Sprintf("What is the least residue of %d mod %d? (enter a number)", a, m),
		Answer:      fmt.Sprintf("%d", r),
		Explanation: fmt.Sprintf("%d = %d x %d + %d, so the least residue is %d.", a, m, (a-r)/m, r, r),
		Facts:       map[string]string{"a": fmt.Sprintf("%d", a), "m": fmt.Sprintf("%d", m), "residue": fmt.Sprintf("%d", r)},
	}
}

type fermatLittleGen struct{}

func (g *fermatLittleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	a := rand.Intn(p-2) + 2
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("By Fermat's little theorem, what is \\(%d^{%d} \\bmod %d\\)?", a, p-1, p),
			Answer:      "1",
			Explanation: fmt.Sprintf("Fermat: \\(a^{p-1} \\equiv 1 \\pmod{p}\\) for \\(a\\) coprime to \\(p=%d\\).", p),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("By Fermat's little theorem, what is \\(%d^{%d} \\bmod %d\\)?", a, p, p),
		Answer:      fmt.Sprintf("%d", a),
		Explanation: fmt.Sprintf("Multiply Fermat by \\(a\\): \\(a^{p} \\equiv a \\pmod{%d}\\), so the answer is %d.", p, a),
	}
}

type eulerPhiGen struct{}

func (g *eulerPhiGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	n := rand.Intn(max(1, scale*23)) + 8
	phi := 0
	for k := 1; k <= n; k++ {
		if mathutil.GCD(k, n) == 1 {
			phi++
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is \\(\\varphi(%d)\\) (Euler's totient)?", n),
		Answer:      fmt.Sprintf("%d", phi),
		Explanation: fmt.Sprintf("\\(\\varphi(%d) = %d\\) (numbers \\(1 \\ldots %d\\) coprime to %d).", n, phi, n, n),
	}
}

type diophantineGen struct{}

func (g *diophantineGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := 1 + rand.Intn(scale*10)
	y := 1 + rand.Intn(scale*10)
	c := 3*x + 5*y
	return generator.Problem{
		Question:    fmt.Sprintf("A solution of 3x + 5y = %d has x = %d. What is y? (enter a number)", c, x),
		Answer:      fmt.Sprintf("%d", y),
		Explanation: fmt.Sprintf("Substitute x = %d: 3 x %d + 5y = %d, so 5y = %d and y = %d.", x, x, c, c-3*x, y),
		Facts:       map[string]string{"c": fmt.Sprintf("%d", c), "x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y)},
	}
}

type cryptoGen struct{}

func (g *cryptoGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	q := primes[rand.Intn(len(primes))]
	for q == p {
		q = primes[rand.Intn(len(primes))]
	}
	n := p * q
	return generator.Problem{
		Question:    fmt.Sprintf("In RSA, if \\(p=%d\\) and \\(q=%d\\), what is \\(n\\)?", p, q),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("n = p x q = %d x %d = %d.", p, q, n),
	}
}

type bezoutGen struct{}

func (g *bezoutGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*12)) + 6
	b := rand.Intn(max(1, scale*12)) + 6
	gcd := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("By B\\u00e9zout's identity, the smallest positive integer of the form \\(%d x + %d y\\) is \\(\\gcd(%d,%d)\\). What is it?", a, b, a, b),
		Answer:      fmt.Sprintf("%d", gcd),
		Explanation: fmt.Sprintf("The set \\(\\{%d x + %d y\\}\\) contains exactly the multiples of \\(\\gcd(%d,%d)=%d\\).", a, b, a, b, gcd),
	}
}

type crtGen struct{}

func (g *crtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mods := [][2]int{{3, 5}, {3, 7}, {4, 9}, {5, 7}, {5, 9}}
	pair := mods[rand.Intn(len(mods))]
	m1, m2 := pair[0], pair[1]
	if m1 == 4 && m2 == 9 {
		// keep coprime
	}
	r1 := rand.Intn(m1)
	r2 := rand.Intn(m2)
	// brute solve
	M := m1 * m2
	sol := 0
	for x := 0; x < M; x++ {
		if x%m1 == r1 && x%m2 == r2 {
			sol = x
			break
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the smallest \\(x \\ge 0\\) with \\(x \\equiv %d \\pmod{%d}\\) and \\(x \\equiv %d \\pmod{%d}\\).", r1, m1, r2, m2),
		Answer:      fmt.Sprintf("%d", sol),
		Explanation: fmt.Sprintf("Search \\(0 \\le x < %d\\): \\(x=%d\\) satisfies both congruences.", M, sol),
	}
}

type wilsonGen struct{}

func (g *wilsonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Any prime works: the factorial is computed modulo p rather than recalled,
	// so the answer is derived for each draw. The list spans small primes where
	// Wilson is checkable by hand up to 29 where it is not.
	primes := []int{3, 5, 7, 11, 13, 17, 19, 23, 29}
	p := primes[rand.Intn(len(primes))]
	// (p-1)! mod p = p-1 by Wilson
	fact := 1
	for i := 2; i <= p-1; i++ {
		fact = (fact * i) % p
	}
	return generator.Problem{
		Question:    fmt.Sprintf("By Wilson's theorem, what is \\((%d-1)! \\bmod %d\\)? (p=%d is prime)", p, p, p),
		Answer:      fmt.Sprintf("%d", fact),
		Explanation: fmt.Sprintf("Wilson: \\((p-1)! \\equiv -1 \\equiv %d \\pmod{%d}\\).", p-1, p),
	}
}

type primeInfGen struct{}

func (g *primeInfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	sets := [][]int{{2, 3}, {2, 3, 5}, {2, 3, 5, 7}, {3, 5, 7}}
	set := sets[rand.Intn(len(sets))]
	prod := 1
	for _, p := range set {
		prod *= p
	}
	n := prod + 1
	return generator.Problem{
		Question:    fmt.Sprintf("Euclid's construction: N = %d + 1. What is N? (enter a number)", prod),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("Multiply the listed primes to get %d, then add 1, so N = %d.", prod, n),
		Facts:       map[string]string{"product": fmt.Sprintf("%d", prod), "N": fmt.Sprintf("%d", n)},
	}
}

func joinProd(primes []int) string {
	var s strings.Builder
	for i, p := range primes {
		if i > 0 {
			s.WriteString("\\times ")
		}
		fmt.Fprintf(&s, "%d", p)
	}
	return s.String()
}

type quadraticReciprocityGen struct{}

func (g *quadraticReciprocityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return quadraticReciprocityProblem()
}

type primitiveRootGen struct{}

func (g *primitiveRootGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return primitiveRootProblem()
}

type pntGen struct{}

func (g *pntGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return primeNumberTheoremProblem()
}

type sieveGen struct{}

func (g *sieveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return sieveProblem()
}

type zetaGen struct{}

func (g *zetaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return zetaProblem()
}

type classNumberGen struct{}

func (g *classNumberGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return classNumberProblem()
}
