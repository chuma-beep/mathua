package discrete

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func w6TestInt(t *testing.T, p generator.Problem, key string) int {
	t.Helper()
	v, err := strconv.Atoi(p.Facts[key])
	if err != nil {
		t.Fatalf("%s: fact %q = %q is not an integer\nfacts: %v", p.Question, key, p.Facts[key], p.Facts)
	}
	return v
}

func w6TestInts(s string) []int {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	var out []int
	for _, p := range strings.Split(s, ",") {
		v, _ := strconv.Atoi(strings.TrimSpace(p))
		out = append(out, v)
	}
	return out
}

func w6TestClosure(start []int, n int, eps map[int][]int) []int {
	seen := make([]bool, n)
	var out []int
	var visit func(int)
	visit = func(q int) {
		if seen[q] {
			return
		}
		seen[q] = true
		out = append(out, q)
		for _, r := range eps[q] {
			visit(r)
		}
	}
	for _, q := range start {
		visit(q)
	}
	sort.Ints(out)
	return out
}

func w6TestNFA(t *testing.T, p generator.Problem) (int, map[string][]int, map[int][]int, map[int]bool) {
	t.Helper()
	n := w6TestInt(t, p, "n")
	trans := map[string][]int{}
	for _, part := range strings.Split(p.Facts["transitions"], ";") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		trans[kv[0]] = w6TestInts(kv[1])
	}
	eps := map[int][]int{}
	for _, part := range strings.Split(p.Facts["epsilon"], ";") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		q, _ := strconv.Atoi(kv[0])
		eps[q] = w6TestInts(kv[1])
	}
	accept := map[int]bool{}
	for _, q := range w6TestInts(p.Facts["accept"]) {
		accept[q] = true
	}
	return n, trans, eps, accept
}

func w6TestMove(cur []int, sym string, trans map[string][]int) []int {
	seen := map[int]bool{}
	var out []int
	for _, q := range cur {
		key := strconv.Itoa(q) + sym
		for _, r := range trans[key] {
			if !seen[r] {
				seen[r] = true
				out = append(out, r)
			}
		}
	}
	sort.Ints(out)
	return out
}

func w6TestDFA(t *testing.T, p generator.Problem) ([][2]int, map[int]bool, int) {
	t.Helper()
	n := w6TestInt(t, p, "n")
	trans := make([][2]int, n)
	for _, part := range strings.Split(p.Facts["transitions"], ";") {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		head := kv[0]
		q, _ := strconv.Atoi(head[:len(head)-1])
		sym := 0
		if strings.HasSuffix(head, "b") {
			sym = 1
		}
		trans[q][sym], _ = strconv.Atoi(kv[1])
	}
	accept := map[int]bool{}
	for _, q := range w6TestInts(p.Facts["accept"]) {
		accept[q] = true
	}
	return trans, accept, n
}

func w6TestDFAAccepts(trans [][2]int, accept map[int]bool, word string) int {
	q := 0
	for _, ch := range word {
		sym := 0
		if ch == 'b' {
			sym = 1
		}
		q = trans[q][sym]
	}
	if accept[q] {
		return 1
	}
	return 0
}

func w6TestDFACount(trans [][2]int, accept map[int]bool, length int) int {
	cur := make([]int, len(trans))
	cur[0] = 1
	for i := 0; i < length; i++ {
		next := make([]int, len(trans))
		for q, c := range cur {
			if c == 0 {
				continue
			}
			next[trans[q][0]] += c
			next[trans[q][1]] += c
		}
		cur = next
	}
	total := 0
	for q, c := range cur {
		if accept[q] {
			total += c
		}
	}
	return total
}

func w6TestHamming(msg int) int {
	m := [4]int{(msg >> 3) & 1, (msg >> 2) & 1, (msg >> 1) & 1, msg & 1}
	c := [7]int{m[0], m[1], m[2], m[3], 0, 0, 0}
	c[4] = (m[0] + m[1] + m[3]) % 2
	c[5] = (m[0] + m[2] + m[3]) % 2
	c[6] = (m[1] + m[2] + m[3]) % 2
	out := 0
	for _, v := range c {
		out = out<<1 | v
	}
	return out
}

func w6TestDist(a, b int) int {
	x := a ^ b
	d := 0
	for x > 0 {
		d += x & 1
		x >>= 1
	}
	return d
}

func w6TestMinCoins(coins []int, target int) int {
	const inf = 1 << 20
	dp := make([]int, target+1)
	for i := 1; i <= target; i++ {
		dp[i] = inf
	}
	for i := 1; i <= target; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	return dp[target]
}

func w6TestCountWays(coins []int, target int) int {
	dp := make([]int, target+1)
	dp[0] = 1
	for _, c := range coins {
		for i := c; i <= target; i++ {
			dp[i] += dp[i-c]
		}
	}
	return dp[target]
}

func TestWave6DiscreteAnswerRecomputation(t *testing.T) {
	kinds := []string{"nfa", "regex", "linear_code", "dynamic_programming"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			for draw := 0; draw < 30; draw++ {
				p := (&w6DiscreteGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				want := w6TestDiscreteAnswer(t, kind, kp, p)
				if p.Answer != strconv.Itoa(want) {
					t.Fatalf("%s kp%d: answer = %s, want %d\nquestion: %s\nfacts: %v", kind, kp, p.Answer, want, p.Question, p.Facts)
				}
			}
		}
	}
}

func w6TestDiscreteAnswer(t *testing.T, kind string, kp int, p generator.Problem) int {
	t.Helper()
	switch kind {
	case "nfa":
		n, trans, eps, accept := w6TestNFA(t, p)
		switch kp {
		case 1:
			start := w6TestClosure([]int{0}, n, eps)
			seen := map[string]bool{strings.Join(intsToText(start), ","): true}
			queue := [][]int{start}
			for len(queue) > 0 {
				cur := queue[0]
				queue = queue[1:]
				for _, sym := range []string{"a", "b"} {
					next := w6TestClosure(w6TestMove(cur, sym, trans), n, eps)
					if len(next) == 0 {
						continue
					}
					key := strings.Join(intsToText(next), ",")
					if !seen[key] {
						seen[key] = true
						queue = append(queue, next)
					}
				}
			}
			return len(seen)
		case 2:
			cur := w6TestClosure([]int{0}, n, eps)
			for _, ch := range p.Facts["string"] {
				cur = w6TestClosure(w6TestMove(cur, string(ch), trans), n, eps)
			}
			for _, q := range cur {
				if accept[q] {
					return 1
				}
			}
			return 0
		default:
			target := w6TestInt(t, p, "target")
			return len(w6TestClosure([]int{target}, n, eps))
		}
	case "regex":
		trans, accept, _ := w6TestDFA(t, p)
		if kp == 1 {
			return w6TestDFACount(trans, accept, w6TestInt(t, p, "length"))
		}
		if kp == 2 {
			return w6TestDFAAccepts(trans, accept, p.Facts["string"])
		}
		total := 0
		length := w6TestInt(t, p, "length")
		for i := 0; i <= length; i++ {
			total += w6TestDFACount(trans, accept, i)
		}
		return total
	case "linear_code":
		msg := w6TestInt(t, p, "message")
		cw := w6TestHamming(msg)
		if kp == 1 {
			return cw
		}
		if kp == 2 {
			return w6TestDist(cw, w6TestHamming(w6TestInt(t, p, "message2")))
		}
		received := w6TestInt(t, p, "received")
		for m := 0; m < 16; m++ {
			if w6TestDist(received, w6TestHamming(m)) == 1 {
				return m
			}
		}
		t.Fatalf("received word %07b has no distance-one codeword", received)
	case "dynamic_programming":
		coins := w6TestInts(p.Facts["coins"])
		if kp == 1 {
			return w6TestMinCoins(coins, w6TestInt(t, p, "t"))
		}
		if kp == 2 {
			return w6TestMinCoins(coins, w6TestInt(t, p, "target"))
		}
		return w6TestCountWays(coins, w6TestInt(t, p, "target"))
	}
	t.Fatalf("unknown kind %q", kind)
	return 0
}

func intsToText(xs []int) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = strconv.Itoa(x)
	}
	return out
}

func TestWave6DiscreteDiversity(t *testing.T) {
	kinds := []string{"nfa", "regex", "linear_code", "dynamic_programming"}
	for _, kind := range kinds {
		for kp := 1; kp <= 3; kp++ {
			seen := map[string]bool{}
			for i := 0; i < 60; i++ {
				p := (&w6DiscreteGen{kind: kind}).Generate(generator.GeneratorContext{Difficulty: 0.5, KPOrdinal: kp})
				seen[p.Question] = true
			}
			if len(seen) < 5 {
				t.Fatalf("%s kp%d: only %d distinct questions in 60 draws", kind, kp, len(seen))
			}
			if len(seen) < 10 {
				t.Logf("%s kp%d: %d distinct questions in 60 draws, below the 10 target", kind, kp, len(seen))
			}
		}
	}
}
