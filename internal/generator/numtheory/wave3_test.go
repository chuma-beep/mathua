package numtheory

import (
	"strconv"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w3TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w3TestPrime(n int) bool {
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

func w3TestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func w3TestPowMod(a, e, m int) int {
	if m == 1 {
		return 0
	}
	r := 1
	a %= m
	for e > 0 {
		if e&1 == 1 {
			r = r * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return r
}

func w3TestPhi(n int) int {
	c := 0
	for k := 1; k <= n; k++ {
		if w3TestGCD(k, n) == 1 {
			c++
		}
	}
	return c
}

func w3TestOrder10(n int) int {
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

func w3TestProperSum(n int) int {
	s := 0
	for d := 1; d < n; d++ {
		if n%d == 0 {
			s += d
		}
	}
	return s
}

func w3TestTau(n int) int {
	c := 0
	for d := 1; d <= n; d++ {
		if n%d == 0 {
			c++
		}
	}
	return c
}

func w3TestDigitalRoot(n int) int {
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

func w3TestHappy(n int) bool {
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

func w3TestHCN(n int) bool {
	if n < 1 {
		return false
	}
	if n == 1 {
		return true
	}
	d := w3TestTau(n)
	for m := 1; m < n; m++ {
		if w3TestTau(m) >= d {
			return false
		}
	}
	return true
}

func w3TestCountConsecutive(n int) int {
	c := 0
	for k := 2; k*(k+1)/2 <= n; k++ {
		if (2*n)%k == 0 && (2*n/k-k+1)%2 == 0 {
			c++
		}
	}
	return c
}

func w3TestShortestConsecutive(n int) int {
	for k := 2; k*(k+1)/2 <= n; k++ {
		if (2*n)%k == 0 && (2*n/k-k+1)%2 == 0 {
			return k
		}
	}
	return 0
}

func TestWave3NumberAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"mod_power", "mod_inverse", "euler_theorem", "decimal_period", "digital_root",
		"aliquot", "perfect_numbers", "mersenne", "fermat_numbers", "repunit",
		"prime_counting", "twin_primes", "goldbach", "polygonal", "pentagonal",
		"egyptian", "highly_composite", "amicable", "happy", "consecutive_sums",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				g := &w3NumberGen{kind: kind}
				p := g.Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w3TestNumberAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w3TestNumberAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "mod_power":
		a := w3TestInt(t, p, "base")
		e := w3TestInt(t, p, "exponent")
		m := w3TestInt(t, p, "modulus")
		if kp == 3 {
			b := w3TestInt(t, p, "base2")
			x := w3TestInt(t, p, "exponent2")
			y := w3TestInt(t, p, "exponent3")
			return w3TestPowMod(w3TestPowMod(a, x, m)*w3TestPowMod(b, y, m), 1, m)
		}
		return w3TestPowMod(a, e, m)
	case "mod_inverse":
		a := w3TestInt(t, p, "a")
		m := w3TestInt(t, p, "modulus")
		b := 1
		if kp >= 2 {
			b = w3TestInt(t, p, "b")
		}
		for x := 0; x < m; x++ {
			if (a*x)%m == b%m {
				return x
			}
		}
		t.Fatalf("no modular inverse for %d mod %d", a, m)
	case "euler_theorem":
		n := w3TestInt(t, p, "n")
		a := w3TestInt(t, p, "a")
		phi := w3TestInt(t, p, "phi")
		if kp == 1 {
			return w3TestPhi(n)
		}
		e := w3TestInt(t, p, "exponent")
		if kp == 3 {
			if e%phi != 0 {
				t.Fatalf("exponent %d is not a multiple of phi %d", e, phi)
			}
			if w3TestPowMod(a, e, n) != 1 {
				t.Fatalf("Euler theorem failed for a=%d n=%d e=%d", a, n, e)
			}
			return 1
		}
		return w3TestPowMod(a, e, n)
	case "decimal_period":
		n := w3TestInt(t, p, "denominator")
		if kp == 1 {
			return w3TestOrder10(n)
		}
		if kp == 2 {
			return w3TestOrder10(n)
		}
		m := w3TestInt(t, p, "core")
		return w3TestOrder10(m)
	case "digital_root":
		if kp == 1 {
			return w3TestDigitalRoot(w3TestInt(t, p, "n"))
		}
		if kp == 2 {
			return w3TestDigitalRoot(w3TestInt(t, p, "a") + w3TestInt(t, p, "b"))
		}
		r := w3TestInt(t, p, "n") % 9
		if r == 0 {
			r = 9
		}
		return r
	case "aliquot":
		n := w3TestInt(t, p, "n")
		s := w3TestProperSum(n)
		if kp == 1 {
			return s
		}
		if kp == 2 {
			if s == n {
				return 1
			}
			if s > n {
				return 2
			}
			return 3
		}
		return s + n
	case "perfect_numbers":
		if kp == 1 {
			limit := w3TestInt(t, p, "limit")
			c := 0
			for n := 2; n <= limit; n++ {
				if w3TestProperSum(n) == n {
					c++
				}
			}
			return c
		}
		if kp == 2 {
			limit := w3TestInt(t, p, "limit")
			c := 0
			for n := 2; n <= limit; n++ {
				if w3TestProperSum(n) > n {
					c++
				}
			}
			return c
		}
		n := w3TestInt(t, p, "after") + 1
		for w3TestProperSum(n) != n {
			n++
		}
		return n
	case "mersenne":
		n := w3TestInt(t, p, "n")
		mn := (1 << n) - 1
		if kp == 1 {
			return mn
		}
		if kp == 2 {
			return mn % w3TestInt(t, p, "modulus")
		}
		if w3TestPrime(mn) {
			return 1
		}
		return 0
	case "fermat_numbers":
		n := w3TestInt(t, p, "n")
		fn := (1 << (1 << n)) + 1
		if kp == 1 {
			return fn
		}
		if kp == 2 {
			return fn % w3TestInt(t, p, "modulus")
		}
		if w3TestPrime(fn) {
			return 1
		}
		return 0
	case "repunit":
		n := w3TestInt(t, p, "n")
		r := 0
		for i := 0; i < n; i++ {
			r = r*10 + 1
		}
		if kp == 1 {
			return r
		}
		if kp == 2 {
			return r % w3TestInt(t, p, "modulus")
		}
		return n
	case "prime_counting":
		if kp == 1 {
			limit := w3TestInt(t, p, "limit")
			c := 0
			for n := 2; n <= limit; n++ {
				if w3TestPrime(n) {
					c++
				}
			}
			return c
		}
		if kp == 2 {
			k := w3TestInt(t, p, "k")
			c, n := 0, 1
			for c < k {
				n++
				if w3TestPrime(n) {
					c++
				}
			}
			return n
		}
		lo, hi := w3TestInt(t, p, "lo"), w3TestInt(t, p, "hi")
		c := 0
		for n := lo; n <= hi; n++ {
			if w3TestPrime(n) {
				c++
			}
		}
		return c
	case "twin_primes":
		if kp == 1 {
			limit := w3TestInt(t, p, "limit")
			c := 0
			for n := 2; n+2 <= limit; n++ {
				if w3TestPrime(n) && w3TestPrime(n+2) {
					c++
				}
			}
			return c
		}
		if kp == 2 {
			n := w3TestInt(t, p, "after") + 1
			for !(w3TestPrime(n) && w3TestPrime(n+2)) {
				n++
			}
			return n
		}
		n := w3TestInt(t, p, "p")
		if w3TestPrime(n) && w3TestPrime(n+2) {
			return 1
		}
		return 0
	case "goldbach":
		n := w3TestInt(t, p, "even")
		var pairs [][2]int
		for q := 2; q <= n/2; q++ {
			if w3TestPrime(q) && w3TestPrime(n-q) {
				pairs = append(pairs, [2]int{q, n - q})
			}
		}
		if kp == 1 {
			return len(pairs)
		}
		if len(pairs) == 0 {
			t.Fatalf("no Goldbach partitions for %d", n)
		}
		if kp == 2 {
			return pairs[0][0]
		}
		return pairs[len(pairs)-1][1]
	case "polygonal":
		s := w3TestInt(t, p, "sides")
		if kp == 1 || kp == 2 {
			n := w3TestInt(t, p, "n")
			val := ((s-2)*n*n - (s-4)*n) / 2
			if kp == 1 {
				return val
			}
			k := 1
			for ((s-2)*k*k-(s-4)*k)/2 < val {
				k++
			}
			return k
		}
		n := w3TestInt(t, p, "n")
		return (s-2)*n + 1
	case "pentagonal":
		if kp == 1 {
			n := w3TestInt(t, p, "n")
			return n * (3*n - 1) / 2
		}
		if kp == 2 {
			val := w3TestInt(t, p, "pentagonal")
			k := 1
			for k*(3*k-1)/2 < val {
				k++
			}
			return k
		}
		val := w3TestInt(t, p, "pentagonal")
		c := 0
		for m := 1; m*(3*m-1)/2 <= val; m++ {
			c++
		}
		return c
	case "egyptian":
		if kp == 3 {
			n := w3TestInt(t, p, "identity_n")
			return n * (n + 1)
		}
		a := w3TestInt(t, p, "numerator")
		b := w3TestInt(t, p, "denominator")
		d1 := (b + a - 1) / a
		if kp == 1 {
			return d1
		}
		rn := a*d1 - b
		rd := b * d1
		g := w3TestGCD(rn, rd)
		rn /= g
		rd /= g
		return (rd + rn - 1) / rn
	case "highly_composite":
		n := w3TestInt(t, p, "n")
		if kp == 1 {
			return w3TestTau(n)
		}
		if kp == 2 {
			if w3TestHCN(n) {
				return 1
			}
			return 0
		}
		next := n + 1
		for !w3TestHCN(next) {
			next++
		}
		return next
	case "amicable":
		n := w3TestInt(t, p, "n")
		s := w3TestProperSum(n)
		if kp == 1 {
			return s
		}
		if kp == 2 {
			if s != n && w3TestProperSum(s) == n {
				return 1
			}
			return 0
		}
		if s != n && w3TestProperSum(s) == n {
			return s
		}
		return 0
	case "happy":
		n := w3TestInt(t, p, "n")
		if kp == 1 {
			s := 0
			for m := n; m > 0; m /= 10 {
				d := m % 10
				s += d * d
			}
			return s
		}
		if kp == 2 {
			if w3TestHappy(n) {
				return 1
			}
			return 0
		}
		next := n
		for !w3TestHappy(next) {
			next++
		}
		return next
	case "consecutive_sums":
		n := w3TestInt(t, p, "n")
		if kp == 1 {
			return w3TestCountConsecutive(n)
		}
		if kp == 2 {
			return w3TestShortestConsecutive(n)
		}
		is := 0
		// A positive integer is not a power of two exactly when it has an odd factor > 1.
		m := n
		for m%2 == 0 {
			m /= 2
		}
		if m > 1 {
			is = 1
		}
		return is
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave3NumberQuestionDiversity(t *testing.T) {
	kinds := []string{
		"mod_power", "mod_inverse", "euler_theorem", "decimal_period", "digital_root",
		"aliquot", "perfect_numbers", "mersenne", "fermat_numbers", "repunit",
		"prime_counting", "twin_primes", "goldbach", "polygonal", "pentagonal",
		"egyptian", "highly_composite", "amicable", "happy", "consecutive_sums",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w3NumberGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
