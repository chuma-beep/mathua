package discrete

import (
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6RemTestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6RemTestInts(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' })
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		v, _ := strconv.Atoi(part)
		out = append(out, v)
	}
	return out
}

func w6RemTestAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w6RemTestMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6RemTestMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func w6RemTestFact(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func w6RemTestChoose(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func w6RemTestDerange(n int) int {
	d := []int{1, 0, 1, 2, 9, 44, 265}
	if n < len(d) {
		return d[n]
	}
	return 0
}

func w6RemTestQB(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	dp[0][0] = 1
	for i := 1; i <= n; i++ {
		for j := 0; j <= k && j <= i; j++ {
			dp[i][j] = dp[i-1][j]
			if j > 0 {
				p := 1
				for t := 0; t < i-j; t++ {
					p *= 2
				}
				dp[i][j] += p * dp[i-1][j-1]
			}
		}
	}
	return dp[n][k]
}

func w6RemTestSat(formula string, mask int) (int, int) {
	clauses := strings.Split(formula, ";")
	satCount, clauseCount := 0, 0
	for _, c := range clauses {
		lits := w6RemTestInts(c)
		ok := false
		for _, lit := range lits {
			v := (mask >> (w6RemTestAbs(lit) - 1)) & 1
			if (lit < 0 && v == 0) || (lit > 0 && v == 1) {
				ok = true
			}
		}
		if ok {
			clauseCount++
		}
	}
	return satCount, clauseCount
}

func w6RemTestYoungHook(shape []int, row, col int) int {
	hook := shape[row] - col
	for r := row + 1; r < len(shape); r++ {
		if shape[r] > col {
			hook++
		}
	}
	return hook
}

func w6RemTestYoungCount(shape []int) int {
	n := 0
	for _, r := range shape {
		n += r
	}
	num := w6RemTestFact(n)
	den := 1
	for r := 0; r < len(shape); r++ {
		for c := 0; c < shape[r]; c++ {
			den *= w6RemTestYoungHook(shape, r, c)
		}
	}
	return num / den
}

func w6RemTestCountSchur(n int) int {
	c := 0
	for x := 1; x <= n; x++ {
		for y := x; y <= n; y++ {
			z := x + y
			if z > n {
				break
			}
			if x%2 == y%2 && y%2 == z%2 {
				c++
			}
		}
	}
	return c
}

func w6RemTestGreedy(t int) int {
	c := t / 4
	t %= 4
	c += t / 3
	t %= 3
	return c + t
}

func w6RemTestMinCoins(t int) int {
	dp := make([]int, t+1)
	for i := 1; i <= t; i++ {
		dp[i] = i
		for _, c := range []int{1, 3, 4} {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	return dp[t]
}

func TestWave6RemainingDiscreteAnswerRecomputation(t *testing.T) {
	kinds := []string{
		"pumping", "cfg", "pda", "turing", "decidability", "halting", "reduction", "sat",
		"np_complete", "cook_levin", "generator_matrix", "cyclic_code", "reed_solomon",
		"rook", "young", "q_binomial", "permutation_stats", "schur", "greedy",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6DiscreteRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w6RemTestDiscreteAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6RemTestDiscreteAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "pumping":
		n := w6RemTestInt(t, p, "states")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return n
		}
		return n
	case "cfg":
		n := w6RemTestInt(t, p, "n")
		if kp == 1 {
			return n + 1
		}
		if kp == 2 {
			return n
		}
		if n == w6RemTestInt(t, p, "m") {
			return 1
		}
		return 0
	case "pda":
		n, m := w6RemTestInt(t, p, "n"), w6RemTestInt(t, p, "m")
		if kp == 1 {
			return n
		}
		if kp == 2 {
			return w6RemTestMin(n, m)
		}
		return w6RemTestAbs(n - m)
	case "turing":
		n := w6RemTestInt(t, p, "n")
		if kp == 1 {
			return n + 1
		}
		if kp == 2 {
			return n
		}
		return n - w6RemTestInt(t, p, "ones")
	case "decidability":
		if kp == 1 {
			switch p.Facts["language"] {
			case "a^n b^n", "prime numbers", "the set of binary palindromes", "the emptiness problem for DFAs", "the equivalence problem for DFAs", "the set of strings of even length":
				return 1
			default:
				return 0
			}
		}
		return 1
	case "halting":
		if kp == 2 {
			return w6RemTestInt(t, p, "overhead") * w6RemTestInt(t, p, "n")
		}
		return 1
	case "reduction":
		n := w6RemTestInt(t, p, "clauses")
		if kp == 1 {
			return 3 * n
		}
		if kp == 2 {
			return n
		}
		return w6RemTestInt(t, p, "factor1") * w6RemTestInt(t, p, "factor2")
	case "sat":
		formula := p.Facts["formula"]
		if kp == 1 {
			total := 0
			for mask := 0; mask < 8; mask++ {
				_, clauses := w6RemTestSat(formula, mask)
				if clauses == len(strings.Split(formula, ";")) {
					total++
				}
			}
			return total
		}
		if kp == 2 {
			_, clauses := w6RemTestSat(formula, w6RemTestInt(t, p, "assignment"))
			return clauses
		}
		best := 0
		for mask := 0; mask < 8; mask++ {
			_, clauses := w6RemTestSat(formula, mask)
			if clauses > best {
				best = clauses
			}
		}
		return best
	case "np_complete":
		k := w6RemTestInt(t, p, "k")
		if kp == 1 {
			return k * (k - 1) / 2
		}
		if kp == 2 {
			return 1
		}
		return 2
	case "cook_levin":
		n := w6RemTestInt(t, p, "n")
		if kp == 3 {
			return 3 * n
		}
		return n * n
	case "generator_matrix":
		if kp == 2 {
			return 4
		}
		if kp == 3 {
			return 3
		}
		return w6HammingCodeword(w6RemTestInt(t, p, "message"))
	case "cyclic_code":
		if kp == 1 {
			word := w6RemTestInt(t, p, "word")
			return ((word << 1) & 127) | ((word >> 6) & 1)
		}
		r := w6RemTestInt(t, p, "r")
		if kp == 2 {
			return r
		}
		return w6RemTestInt(t, p, "n") - r
	case "reed_solomon":
		n, k := w6RemTestInt(t, p, "n"), w6RemTestInt(t, p, "k")
		if kp == 1 {
			pp := w6RemTestInt(t, p, "p")
			i := w6RemTestInt(t, p, "i")
			a0 := w6RemTestInt(t, p, "a0")
			a1 := w6RemTestInt(t, p, "a1")
			a2 := w6RemTestInt(t, p, "a2")
			return (a0 + a1*i + a2*i*i) % pp
		}
		if kp == 2 {
			return n - k + 1
		}
		return (n - k) / 2
	case "rook":
		m, n := w6RemTestInt(t, p, "m"), w6RemTestInt(t, p, "n")
		if kp == 1 {
			return m * n
		}
		if kp == 2 {
			return m * (m - 1) * n * (n - 1) / 2
		}
		return w6RemTestDerange(w6RemTestInt(t, p, "d"))
	case "young":
		shape := w6RemTestInts(p.Facts["shape"])
		if kp == 1 {
			n := 0
			for _, r := range shape {
				n += r
			}
			return n
		}
		if kp == 2 {
			return w6RemTestYoungHook(shape, w6RemTestInt(t, p, "row"), w6RemTestInt(t, p, "col"))
		}
		return w6RemTestYoungCount(shape)
	case "q_binomial":
		n, k := w6RemTestInt(t, p, "n"), w6RemTestInt(t, p, "k")
		if kp == 1 {
			return w6RemTestQB(n, k)
		}
		if kp == 2 {
			return w6RemTestChoose(n, k)
		}
		return 1
	case "permutation_stats":
		perm := w6RemTestInts(p.Facts["permutation"])
		if kp == 1 {
			inv := 0
			for i := 0; i < len(perm); i++ {
				for j := i + 1; j < len(perm); j++ {
					if perm[i] > perm[j] {
						inv++
					}
				}
			}
			return inv
		}
		if kp == 2 {
			major := 0
			for i := 0; i+1 < len(perm); i++ {
				if perm[i] > perm[i+1] {
					major += i + 1
				}
			}
			return major
		}
		return w6RemTestInt(t, p, "n") - 1
	case "schur":
		r := w6RemTestInt(t, p, "r")
		s := map[int]int{2: 5, 3: 14, 4: 45}[r]
		if kp == 1 {
			return s + w6RemTestInt(t, p, "offset")
		}
		if kp == 2 {
			return s - 1 + w6RemTestInt(t, p, "offset")
		}
		return w6RemTestCountSchur(w6RemTestInt(t, p, "n"))
	case "greedy":
		target := w6RemTestInt(t, p, "target")
		if kp == 1 {
			return w6RemTestGreedy(target)
		}
		if kp == 2 {
			return w6RemTestMinCoins(target)
		}
		if w6RemTestGreedy(target) == w6RemTestMinCoins(target) {
			return 1
		}
		return 0
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func TestWave6RemainingDiscreteDiversity(t *testing.T) {
	kinds := []string{
		"pumping", "cfg", "pda", "turing", "decidability", "halting", "reduction", "sat",
		"np_complete", "cook_levin", "generator_matrix", "cyclic_code", "reed_solomon",
		"rook", "young", "q_binomial", "permutation_stats", "schur", "greedy",
	}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6DiscreteRemainingGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Errorf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws", kind, kp, len(seen))
			}
		}
	}
}
