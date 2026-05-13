package numtheory

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func Register(reg *generator.Registry) {
	reg.Register("nt.divisibility", &divisibilityGen{})
	reg.Register("nt.gcd_euclidean", &gcdEuclideanGen{})
	reg.Register("nt.modular", &modularGen{})
	reg.Register("nt.congruence", &congruenceGen{})
	reg.Register("nt.fermat_little", &fermatLittleGen{})
	reg.Register("nt.euler_phi", &eulerPhiGen{})
	reg.Register("nt.diophantine", &diophantineGen{})
	reg.Register("nt.crypto", &cryptoGen{})
}

type divisibilityGen struct{}

func (g *divisibilityGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(45) + 6
	d := rand.Intn(10) + 2
	divisible := a%d == 0
	if rand.Intn(2) == 0 {
		for a%d == 0 {
			a = rand.Intn(45) + 6
		}
		divisible = false
	} else {
		for a%d != 0 {
			a = rand.Intn(45) + 6
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

func (g *gcdEuclideanGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(37) + 12
	b := rand.Intn(37) + 12
	gcd := mathutil.GCD(a, b)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the GCD of %d and %d using the Euclidean algorithm.", a, b),
		Answer:      fmt.Sprintf("%d", gcd),
		Explanation: fmt.Sprintf("GCD(%d, %d) = %d (computed via Euclidean algorithm).", a, b, gcd),
	}
}

type modularGen struct{}

func (g *modularGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(30) + 5
	m := rand.Intn(16) + 2
	ans := a % m
	return generator.Problem{
		Question:    fmt.Sprintf("Compute %d mod %d.", a, m),
		Answer:      fmt.Sprintf("%d", ans),
		Explanation: fmt.Sprintf("%d mod %d = %d because %d = %d x %d + %d.", a, m, ans, a, m, a/m, ans),
	}
}

type congruenceGen struct{}

func (g *congruenceGen) Generate(difficulty float64) generator.Problem {
	m := rand.Intn(10) + 2
	a := rand.Intn(30) + 1
	b := a + m*rand.Intn(5)
	congruent := true
	if rand.Intn(2) == 0 {
		b = a + m*rand.Intn(5) + rand.Intn(m-1) + 1
		congruent = false
	}
	ans := "no"
	exp := fmt.Sprintf("%d - %d = %d, which is not divisible by %d, so %d ≢ %d (mod %d).", a, b, a-b, m, a, b, m)
	if congruent {
		ans = "yes"
		exp = fmt.Sprintf("%d - %d = %d, which is divisible by %d, so %d ≡ %d (mod %d).", a, b, a-b, m, a, b, m)
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is %d ≡ %d (mod %d)? (yes/no)", a, b, m),
		Answer:      ans,
		Explanation: exp,
	}
}

type fermatLittleGen struct{}

func (g *fermatLittleGen) Generate(difficulty float64) generator.Problem {
	primes := []int{3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	a := rand.Intn(p-2) + 2
	result := mathutil.IntPow(a, p-1) % p
	return generator.Problem{
		Question:    fmt.Sprintf("By Fermat's little theorem, what is %d^%d mod %d?", a, p-1, p),
		Answer:      fmt.Sprintf("%d", result),
		Explanation: fmt.Sprintf("Fermat's little theorem: %d^(%d) ≡ 1 (mod %d), so %d^%d mod %d = %d.", a, p-1, p, a, p-1, p, result),
	}
}

type eulerPhiGen struct{}

func (g *eulerPhiGen) Generate(difficulty float64) generator.Problem {
	n := rand.Intn(23) + 8
	phi := 0
	for k := 1; k <= n; k++ {
		if mathutil.GCD(k, n) == 1 {
			phi++
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is φ(%d) (Euler's totient)?", n),
		Answer:      fmt.Sprintf("%d", phi),
		Explanation: fmt.Sprintf("φ(%d) = %d (numbers 1..%d coprime to %d).", n, phi, n, n),
	}
}

type diophantineGen struct{}

func (g *diophantineGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(10) + 1
	y := rand.Intn(10) + 1
	c := 3*x + 5*y
	return generator.Problem{
		Question:    fmt.Sprintf("Find an integer solution to 3x + 5y = %d.", c),
		Answer:      fmt.Sprintf("x=%d,y=%d", x, y),
		Explanation: fmt.Sprintf("3(%d) + 5(%d) = %d + %d = %d.", x, y, 3*x, 5*y, c),
	}
}

type cryptoGen struct{}

func (g *cryptoGen) Generate(difficulty float64) generator.Problem {
	primes := []int{3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	q := primes[rand.Intn(len(primes))]
	for q == p {
		q = primes[rand.Intn(len(primes))]
	}
	n := p * q
	return generator.Problem{
		Question:    fmt.Sprintf("In RSA, if p=%d and q=%d, what is n?", p, q),
		Answer:      fmt.Sprintf("%d", n),
		Explanation: fmt.Sprintf("n = p x q = %d x %d = %d.", p, q, n),
	}
}
