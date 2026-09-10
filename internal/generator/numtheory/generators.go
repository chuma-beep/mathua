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
	reg.Register("nt.adv.legendre", &legendreGen{})
	reg.Register("nt.adv.quadratic_reciprocity", &quadraticReciprocityGen{})
	reg.Register("nt.basics.order_mod", &orderModGen{})
	reg.Register("nt.adv.primitive_root", &primitiveRootGen{})
	reg.Register("nt.adv.mobius", &mobiusGen{})
	reg.Register("nt.analytic.prime_number_theorem", &pntGen{})
	reg.Register("nt.adv.dirichlet", &dirichletGen{})
	reg.Register("nt.adv.sieve", &sieveGen{})
	reg.Register("nt.analytic.zeta", &zetaGen{})
	reg.Register("nt.adv.chebyshev", &chebyshevGen{})
	reg.Register("nt.adv.elliptic_curve", &ellipticCurveGen{})
	reg.Register("nt.adv.continued_fractions", &continuedFractionsGen{})
	reg.Register("nt.adv.class_number", &classNumberGen{})
}

type divisibilityGen struct{}

func (g *divisibilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := rand.Intn(max(1, scale*45)) + 6
	d := rand.Intn(max(1, scale*10)) + 2
	divisible := a%d == 0
	if rand.Intn(max(1, scale*2)) == 0 {
		for a%d == 0 {
			a = rand.Intn(max(1, scale*45)) + 6
		}
		divisible = false
	} else {
		for a%d != 0 {
			a = rand.Intn(max(1, scale*45)) + 6
		}
		divisible = true
	}
	ans := "no"
	exp := fmt.Sprintf("%d / %d = %d, which is not an integer.", a, d, a/d)
	if divisible {
		ans = "yes"
		exp = fmt.Sprintf("%d / %d = %d, which is an integer.", a, d, a/d)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is %d divisible by %d? (yes/no)", a, d),
		Answer:      ans,
		Explanation: exp,
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
	scale := int(1 + ctx.Difficulty*5)
	m := rand.Intn(max(1, scale*10)) + 2
	a := rand.Intn(max(1, scale*30)) + 1
	b := a + m*rand.Intn(max(1, scale*5))
	congruent := true
	if rand.Intn(max(1, scale*2)) == 0 {
		b = a + m*rand.Intn(max(1, scale*5)) + rand.Intn(m-1) + 1
		congruent = false
	}
	ans := "no"
	exp := fmt.Sprintf("\\(%d - %d = %d\\), which is not divisible by \\(%d\\), so \\(%d \\not\\equiv %d \\pmod{%d}\\).", a, b, a-b, m, a, b, m)
	if congruent {
		ans = "yes"
		exp = fmt.Sprintf("\\(%d - %d = %d\\), which is divisible by \\(%d\\), so \\(%d \\equiv %d \\pmod{%d}\\).", a, b, a-b, m, a, b, m)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is \\(%d \\equiv %d \\pmod{%d}\\)? (yes/no)", a, b, m),
		Answer:      ans,
		Explanation: exp,
	}
}

type fermatLittleGen struct{}

func (g *fermatLittleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	a := rand.Intn(p-2) + 2
	result := mathutil.IntPow(a, p-1) % p
	return generator.Problem{
		Question:    fmt.Sprintf("By Fermat's little theorem, what is \\(%d^{%d} \\bmod %d\\)?", a, p-1, p),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("Fermat's little theorem: \\(%d^{%d} \\equiv 1 \\pmod{%d}\\), so \\(%d^{%d} \\bmod %d = %d\\).", a, p-1, p, a, p-1, p, result),
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
	x := rand.Intn(max(1, scale*10)) + 1
	y := rand.Intn(max(1, scale*10)) + 1
	c := 3*x + 5*y
	return generator.Problem{
		Question:    fmt.Sprintf("Find an integer solution to \\(3x + 5y = %d\\).", c),
		Answer:      fmt.Sprintf("x=%d,y=%d", x, y),
		Explanation: fmt.Sprintf("3(%d) + 5(%d) = %d + %d = %d.", x, y, 3*x, 5*y, c),
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
	primes := []int{5, 7, 11, 13}
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
	primes := sets[rand.Intn(len(sets))]
	prod := 1
	for _, p := range primes {
		prod *= p
	}
	n := prod + 1
	// question: is n divisible by any listed prime? always no
	return generator.Problem{
		Question:    fmt.Sprintf("Euclid's construction: \\(N = %s + 1 = %d\\) where the product is over \\(%v\\). Is \\(N\\) divisible by any of \\(%v\\)? (yes/no)", joinProd(primes), n, primes, primes),
		Answer:      "no",
		Explanation: fmt.Sprintf("If \\(N\\) were divisible by any \\(p\\) in the set, \\(p\\) would divide \\(N - \\text{product}=1\\), impossible. So \\(N\\) has a new prime factor."),
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

type legendreGen struct{}

func (g *legendreGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"What is the Legendre symbol (1/5)? (enter -1,0,1)", "1", "1 is a quadratic residue mod 5 (1^2=1)."},
		{"What is the Legendre symbol (4/5)? (enter -1,0,1)", "1", "4=2^2 mod5, so residue."},
		{"What is the Legendre symbol (2/5)? (enter -1,0,1)", "-1", "2 is non-residue mod5 (no x with x^2=2)."},
		{"Is (a/p)=0 iff p divides a? (yes/no)", "yes", "Legendre 0 exactly when p|a."},
	}
	hard := []entry{
		{"What is the Legendre symbol (3/7)? (enter -1,0,1)", "-1", "Quadratic residues mod7 are 1,2,4; 3 not in set."},
		{"Does Euler's criterion give (a/p)≡a^{(p-1)/2} mod p? (yes/no)", "yes", "Euler's criterion for odd prime p."},
		{"Is (2/p)=1 iff p≡±1 mod8? (yes/no)", "yes", "Supplement to quadratic reciprocity."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type quadraticReciprocityGen struct{}

func (g *quadraticReciprocityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does quadratic reciprocity relate (p/q) and (q/p) for odd primes p,q? (yes/no)", "yes", "(p/q)(q/p)=(-1)^{(p-1)(q-1)/4}."},
		{"If p=5,q=7, does (5/7)=(7/5) hold? (yes/no)", "yes", "Both ≡1 mod4, so sign is +1."},
		{"Is (3/11) = -(11/3) ? (yes/no)", "yes", "3≡3 mod4,11≡3 mod4 → sign -1."},
	}
	hard := []entry{
		{"What is (3/5) using reciprocity? (enter -1,0,1)", "-1", "(3/5)=-(5/3)=-(2/3)=-( -1)= -1? Direct: residues mod5 1,4, so 3 non-residue."},
		{"Does (p/q)=1 imply q is quadratic residue mod p? (yes/no)", "yes", "Definition of Legendre symbol."},
		{"Is law of quadratic reciprocity used to compute (a/p) quickly? (yes/no)", "yes", "It swaps numerator/denominator to reduce."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type orderModGen struct{}

func (g *orderModGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	// order of 2 mod small primes
	entries := []struct {
		a, m, ord int
	}{
		{2, 5, 4}, {2, 7, 3}, {3, 7, 6}, {2, 9, 6}, {3, 8, 2}, {2, 11, 10},
	}
	e := entries[rand.Intn(len(entries))]
	if scale > 3 {
		// for harder, ask for order value; easy asks yes/no about divisibility of phi
		return generator.Problem{
			Question:    fmt.Sprintf("What is the order of %d modulo %d? (smallest k with %d^k≡1 mod %d)", e.a, e.m, e.a, e.m),
			Answer:      fmt.Sprintf("%d", e.ord),
			Explanation: fmt.Sprintf("Powers of %d mod %d cycle with period %d.", e.a, e.m, e.ord),
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Does the order of %d mod %d divide φ(%d)=%d? (yes/no)", e.a, e.m, e.m, e.ord+2),
		Answer:      "yes",
		Explanation: "Order always divides φ(m) (Euler) or p-1 for prime.",
	}
}

type primitiveRootGen struct{}

func (g *primitiveRootGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is 2 a primitive root mod 5? (yes/no)", "yes", "Powers of 2 mod5: 2,4,3,1 cover all non-zero residues."},
		{"Is 2 a primitive root mod 7? (yes/no)", "no", "Order of 2 mod7 is 3, not 6=φ(7), so not primitive."},
		{"Does primitive root mod p have order p-1? (yes/no)", "yes", "By definition order is φ(p)=p-1."},
	}
	hard := []entry{
		{"How many primitive roots mod 7 are there? (enter a number)", "2", "Count = φ(φ(7))=φ(6)=2 (3 and 5)."},
		{"Does every prime have a primitive root? (yes/no)", "yes", "Primitive root theorem for primes."},
		{"Is 3 a primitive root mod 7? (yes/no)", "yes", "Powers of 3: 3,2,6,4,5,1 cover all residues."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type mobiusGen struct{}

func (g *mobiusGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		n int
		mu int
	}
	easy := []entry{
		{1, 1}, {2, -1}, {3, -1}, {5, -1}, {6, 1}, {30, -1},
	}
	hard := []entry{
		{4, 0}, {8, 0}, {12, 0}, {18, 0}, {7, -1}, {10, 1},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is μ(%d) (Möbius function)? (enter -1,0,1)", e.n),
		Answer:      fmt.Sprintf("%d", e.mu),
		Explanation: fmt.Sprintf("μ(n)=0 if squared prime divides n, otherwise (-1)^k where k is number of prime factors; μ(%d)=%d.", e.n, e.mu),
	}
}

type pntGen struct{}

func (g *pntGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does prime number theorem state π(x) ~ x/log x? (yes/no)", "yes", "π(x) ~ x/log x as x→∞."},
		{"Does π(x) denote number of primes ≤ x? (yes/no)", "yes", "Prime counting function."},
		{"Is density of primes near x about 1/log x? (yes/no)", "yes", "PNT gives density."},
	}
	hard := []entry{
		{"Does PNT imply n-th prime p_n ~ n log n? (yes/no)", "yes", "Equivalent form of PNT."},
		{"Is error in PNT related to Riemann hypothesis (O(√x log x))? (yes/no)", "yes", "RH gives sharp error."},
		{"Does π(10)=4? (yes/no)", "yes", "Primes ≤10: 2,3,5,7."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type dirichletGen struct{}

func (g *dirichletGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does Dirichlet's theorem guarantee infinitely many primes ≡1 mod4? (yes/no)", "yes", "Primes ≡a mod m with gcd(a,m)=1 are infinite."},
		{"Does arithmetic progression 3,7,11,... (≡3 mod4) have infinitely many primes? (yes/no)", "yes", "Dirichlet with a=3,m=4."},
		{"Is progression 2,4,6,... expected to have infinitely many primes? (yes/no)", "no", "All even >2 are composite; gcd(2,2)≠1."},
	}
	hard := []entry{
		{"Does Dirichlet require gcd(a,m)=1? (yes/no)", "yes", "Otherwise progression has common divisor."},
		{"Is L(1,χ)≠0 the key non-vanishing for Dirichlet? (yes/no)", "yes", "Non-vanishing of L at 1 proves theorem."},
		{"Does Dirichlet generalize Euclid's infinitude? (yes/no)", "yes", "Euclid is case m=1 or a=1 mod m."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type sieveGen struct{}

func (g *sieveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does sieve of Eratosthenes mark multiples of p starting at p^2? (yes/no)", "yes", "Smaller multiples already marked by smaller primes."},
		{"Is sieve time O(n log log n)? (yes/no)", "yes", "Eratosthenes complexity."},
		{"Does Brun's sieve give upper bound for twin primes? (yes/no)", "yes", "Brun: sum of reciprocals of twin primes converges."},
	}
	hard := []entry{
		{"Does sieve principle use inclusion-exclusion with Möbius? (yes/no)", "yes", "Count via μ."},
		{"Is Brun's theorem that twin primes are finite? (no)", "no", "Brun shows sum converges, not finiteness."},
		{"Does segmented sieve handle large n with limited memory? (yes/no)", "yes", "Process intervals."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type zetaGen struct{}

func (g *zetaGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does ζ(s)=∑_{n≥1}1/n^s converge for Re(s)>1? (yes/no)", "yes", "Zeta converges for Re(s)>1."},
		{"Is ζ(2)=π^2/6? (yes/no)", "yes", "Basel problem."},
		{"Does Euler product ζ(s)=∏_p(1-p^{-s})^{-1} hold? (yes/no)", "yes", "Euler product over primes."},
	}
	hard := []entry{
		{"Are nontrivial zeros of ζ conjectured to lie on Re(s)=1/2? (yes/no)", "yes", "Riemann hypothesis."},
		{"Does ζ(s) have pole at s=1? (yes/no)", "yes", "Simple pole with residue 1."},
		{"Is ζ(-1)=-1/12 via analytic continuation? (yes/no)", "yes", "Regularized sum."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type chebyshevGen struct{}

func (g *chebyshevGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is ψ(x)=∑_{n≤x}Λ(n) Chebyshev function? (yes/no)", "yes", "Von Mangoldt weight."},
		{"Does θ(x)=∑_{p≤x}log p count primes with log weight? (yes/no)", "yes", "Chebyshev θ."},
		{"Is PNT equivalent to ψ(x)∼x? (yes/no)", "yes", "Chebyshev form."},
	}
	hard := []entry{
		{"Does explicit formula ψ(x)=x-∑_ρ x^ρ/ρ - log(2π) hold? (yes/no)", "yes", "Over zeros of ζ."},
		{"Is θ(x)∼x equivalent to PNT? (yes/no)", "yes", "Chebyshev equivalence."},
		{"Does Chebyshev show π(x)≍x/log x with constants? (yes/no)", "yes", "Chebyshev bounds."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type ellipticCurveGen struct{}

func (g *ellipticCurveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is elliptic curve E: y^2=x^3+ax+b with discriminant Δ≠0? (yes/no)", "yes", "Nonsingular cubic."},
		{"Does E have group law via chord-tangent? (yes/no)", "yes", "Points form abelian group."},
		{"Is BSD conjecture about rank of E(Q) and L(E,1)? (yes/no)", "yes", "Birch and Swinnerton-Dyer."},
	}
	hard := []entry{
		{"Does Hasse bound |#E(F_p)-(p+1)|≤2√p hold? (yes/no)", "yes", "Hasse-Weil."},
		{"Is elliptic curve used in factorization (ECM) and cryptography? (yes/no)", "yes", "ECM and ECC."},
		{"Does Mordell-Weil state E(Q) finitely generated? (yes/no)", "yes", "E(Q)≅Z^r×torsion."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type continuedFractionsGen struct{}

func (g *continuedFractionsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is continued fraction [a0;a1,a2,...] = a0+1/(a1+1/(a2+...))? (yes/no)", "yes", "Definition."},
		{"Do convergents p_k/q_k give best rational approximations? (yes/no)", "yes", "Best approximations."},
		{"Is √2 = [1;2,2,2,...]? (yes/no)", "yes", "Periodic for quadratic irrationals."},
	}
	hard := []entry{
		{"Does Lagrange show periodic CF ⇔ quadratic irrational? (yes/no)", "yes", "Lagrange's theorem."},
		{"Is e = [2;1,2,1,1,4,1,1,6,...] pattern? (yes/no)", "yes", "Continued fraction for e."},
		{"Does CF of π start [3;7,15,1,292,...]? (yes/no)", "yes", "Approximation 22/7 from first term."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type classNumberGen struct{}

func (g *classNumberGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is class number h(D) count of ideal classes in quadratic order? (yes/no)", "yes", "Class group size."},
		{"Does h(-3)=1 mean Z[ω] is UFD? (yes/no)", "yes", "Class number 1 ⇔ PID ⇔ UFD for rings of integers."},
		{"Is Gauss class number problem about h(D)=1? (yes/no)", "yes", "Heegner completed."},
	}
	hard := []entry{
		{"Does class number formula relate h(D) to L(1,χ_D)? (yes/no)", "yes", "Dirichlet class number formula."},
		{"Is h(-163)=1 the largest Heegner number? (yes/no)", "yes", "Heegner numbers -1,-2,-3,-7,-11,-19,-43,-67,-163."},
		{"Does Stark-Heegner show h(-163)=1? (yes/no)", "yes", "Stark-Heegner theorem."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
