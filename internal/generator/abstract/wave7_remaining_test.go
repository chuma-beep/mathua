package abstract

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w7RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w7RemTestInts(s string) []int {
	s = strings.TrimSpace(strings.Trim(s, "[]"))
	fields := strings.Fields(s)
	out := make([]int, len(fields))
	for i, f := range fields {
		out[i], _ = strconv.Atoi(f)
	}
	return out
}

func w7RemTestGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w7RemTestPow(a, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= a
	}
	return r
}

func w7RemTestBinom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func TestWave7RemainingAbstractAnswerRecomputation(t *testing.T) {
	kinds := []string{"wreath", "free_product", "extension", "character", "maschke", "schur", "induced", "resultant", "discriminant", "eisenstein", "fractional_ideal", "dedekind", "valuation", "primary", "nakayama", "exterior"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 40; draw++ {
				p := (&w7AbstractRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w7RemTestAbstractAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w7RemTestAbstractAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "wreath":
		m, n := w7RemTestInt(t, p, "m"), w7RemTestInt(t, p, "n")
		base := w7RemTestPow(m, n)
		if kp == 1 {
			return base * n
		}
		if kp == 2 {
			return base
		}
		return n
	case "free_product":
		r, s := w7RemTestInt(t, p, "r"), w7RemTestInt(t, p, "s")
		if kp == 1 || kp == 2 {
			return r + s
		}
		return 1 - (r + s)
	case "extension":
		n, m := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "m")
		if kp == 1 {
			return n * m
		}
		if kp == 2 {
			return m
		}
		if w7RemTestGCD(n, m) == 1 {
			return 1
		}
		return 0
	case "character":
		if kp == 1 {
			return w7RemTestInt(t, p, "a")
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "a") + w7RemTestInt(t, p, "b")
		}
		return w7RemTestInt(t, p, "n")
	case "maschke":
		n, prime := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "p")
		if kp == 1 {
			if n%prime == 0 {
				return 0
			}
			return 1
		}
		if kp == 2 {
			return n
		}
		return n % prime
	case "schur":
		if kp == 1 {
			return 1
		}
		if kp == 2 {
			return []int{1, 2, 4}[w7RemTestInt(t, p, "type")]
		}
		return 0
	case "induced":
		idx, deg := w7RemTestInt(t, p, "index"), w7RemTestInt(t, p, "degree")
		if kp == 1 || kp == 2 {
			return idx * deg
		}
		return 1
	case "resultant":
		a, b, c, d := w7RemTestInt(t, p, "a"), w7RemTestInt(t, p, "b"), w7RemTestInt(t, p, "c"), w7RemTestInt(t, p, "d")
		det := a*d - b*c
		if kp == 1 {
			return det
		}
		if kp == 2 {
			if det == 0 {
				return 1
			}
			return 0
		}
		return a * c
	case "discriminant":
		a, b, c := w7RemTestInt(t, p, "a"), w7RemTestInt(t, p, "b"), w7RemTestInt(t, p, "c")
		disc := b*b - 4*a*c
		if kp == 1 {
			return disc
		}
		if kp == 2 {
			if disc == 0 {
				return 1
			}
			return 0
		}
		if disc > 0 {
			return 1
		}
		if disc < 0 {
			return -1
		}
		return 0
	case "eisenstein":
		pv := w7RemTestInt(t, p, "p")
		coefs := w7RemTestInts(p.Facts["coefs"])
		if kp == 1 {
			c := 0
			for _, v := range coefs[:len(coefs)-1] {
				if v%pv == 0 {
					c++
				}
			}
			return c
		}
		if kp == 2 {
			return 1
		}
		if coefs[0]%(pv*pv) == 0 {
			return 1
		}
		return 0
	case "fractional_ideal":
		a, b := w7RemTestInt(t, p, "a"), w7RemTestInt(t, p, "b")
		c := w7RemTestInt(t, p, "c")
		if kp == 1 {
			return b
		}
		if kp == 2 {
			return a * c
		}
		if b == 1 {
			return 1
		}
		return 0
	case "dedekind":
		n := w7RemTestInt(t, p, "n")
		count := 0
		for m := n; m > 1; {
			for q := 2; q <= m; q++ {
				if m%q == 0 {
					count++
					m /= q
					break
				}
			}
		}
		if kp == 1 {
			return count
		}
		if kp == 2 {
			return w7RemTestInt(t, p, "distinct_primes")
		}
		return 2 * count
	case "valuation":
		a, b := w7RemTestInt(t, p, "a"), w7RemTestInt(t, p, "b")
		if kp == 1 {
			return a
		}
		if kp == 2 {
			return a + b
		}
		return 1
	case "primary":
		if kp == 1 {
			return w7RemTestInt(t, p, "n")
		}
		if kp == 2 {
			return 1
		}
		return 1
	case "nakayama":
		d, j := w7RemTestInt(t, p, "d"), w7RemTestInt(t, p, "j")
		if kp == 1 || kp == 3 {
			return d - j
		}
		if j < d {
			return 1
		}
		return 0
	case "exterior":
		n, k := w7RemTestInt(t, p, "n"), w7RemTestInt(t, p, "k")
		if kp == 1 {
			return w7RemTestBinom(n, k)
		}
		if kp == 2 {
			return w7RemTestBinom(n-1, k-1)
		}
		return w7RemTestPow(2, n)
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave7RemainingAbstractDiversity(t *testing.T) {
	kinds := []string{"wreath", "free_product", "extension", "character", "maschke", "schur", "induced", "resultant", "discriminant", "eisenstein", "fractional_ideal", "dedekind", "valuation", "primary", "nakayama", "exterior"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w7AbstractRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
