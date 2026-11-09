package discrete

// Wave 6 discrete mathematics, first batch: nondeterministic automata,
// regular-language counting, linear codes and dynamic programming. Every
// answer is an integer and every generator publishes the instance in Facts so
// wave6_test.go can recompute it independently.

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"discrete.automata.nfa":                   "nfa",
		"discrete.automata.regex":                 "regex",
		"discrete.coding.linear":                  "linear_code",
		"discrete.algorithms.dynamic_programming": "dynamic_programming",
	} {
		reg.Register(id, &w6DiscreteGen{kind: kind})
	}
}

type w6DiscreteGen struct{ kind string }

func (g *w6DiscreteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "nfa":
		return w6NFA(kp)
	case "regex":
		return w6Regex(kp)
	case "linear_code":
		return w6LinearCode(kp)
	case "dynamic_programming":
		return w6DP(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6D(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w6DI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6DR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6DJoin(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = strconv.Itoa(x)
	}
	return strings.Join(parts, ",")
}

func w6DParseInts(s string) []int {
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

// ---------------------------------------------------------------------------
// Nondeterministic finite automata
// ---------------------------------------------------------------------------

type w6NFAState struct {
	trans [2][]int
	eps   []int
}

func w6NFAClosure(states []int, nfa []w6NFAState) []int {
	seen := make([]bool, len(nfa))
	var out []int
	var visit func(int)
	visit = func(q int) {
		if seen[q] {
			return
		}
		seen[q] = true
		out = append(out, q)
		for _, r := range nfa[q].eps {
			visit(r)
		}
	}
	for _, q := range states {
		visit(q)
	}
	sort.Ints(out)
	return out
}

func w6NFAMove(states []int, sym int, nfa []w6NFAState) []int {
	seen := make([]bool, len(nfa))
	var out []int
	for _, q := range states {
		for _, r := range nfa[q].trans[sym] {
			if !seen[r] {
				seen[r] = true
				out = append(out, r)
			}
		}
	}
	sort.Ints(out)
	return out
}

func w6NFAString(nfa []w6NFAState) string {
	var parts []string
	for q := range nfa {
		for sym := 0; sym < 2; sym++ {
			letter := "a"
			if sym == 1 {
				letter = "b"
			}
			parts = append(parts, fmt.Sprintf("%d%s:%s", q, letter, w6DJoin(nfa[q].trans[sym])))
		}
	}
	return strings.Join(parts, ";")
}

func w6NFAParse(s string, n int) []w6NFAState {
	nfa := make([]w6NFAState, n)
	for _, part := range strings.Split(s, ";") {
		if part == "" {
			continue
		}
		colon := strings.Index(part, ":")
		if colon < 0 {
			continue
		}
		head, tail := part[:colon], part[colon+1:]
		q, _ := strconv.Atoi(head[:len(head)-1])
		sym := 0
		if strings.HasSuffix(head, "b") {
			sym = 1
		}
		nfa[q].trans[sym] = w6DParseInts(tail)
	}
	return nfa
}

func w6NFA(kp int) generator.Problem {
	n := w6DR(2, 3)
	nfa := make([]w6NFAState, n)
	for q := 0; q < n; q++ {
		for sym := 0; sym < 2; sym++ {
			k := rand.Intn(3)
			perm := rand.Perm(n)
			for i := 0; i < k; i++ {
				nfa[q].trans[sym] = append(nfa[q].trans[sym], perm[i])
			}
			sort.Ints(nfa[q].trans[sym])
		}
		if rand.Intn(2) == 0 {
			r := rand.Intn(n)
			if r != q {
				nfa[q].eps = []int{r}
			}
		}
	}
	accept := make([]int, 0, n)
	for q := 0; q < n; q++ {
		if rand.Intn(2) == 0 {
			accept = append(accept, q)
		}
	}
	if len(accept) == 0 {
		accept = []int{n - 1}
	}
	f := map[string]string{}
	w6DI(f, "n", n)
	f["transitions"] = w6NFAString(nfa)
	f["epsilon"] = func() string {
		parts := make([]string, n)
		for q := range nfa {
			parts[q] = fmt.Sprintf("%d:%s", q, w6DJoin(nfa[q].eps))
		}
		return strings.Join(parts, ";")
	}()
	f["accept"] = w6DJoin(accept)
	w6DI(f, "start", 0)

	switch kp {
	case 1:
		// Reachable subset states from the start state.
		start := w6NFAClosure([]int{0}, nfa)
		seen := map[string]bool{}
		queue := [][]int{start}
		seen[fmt.Sprint(start)] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for sym := 0; sym < 2; sym++ {
				next := w6NFAClosure(w6NFAMove(cur, sym, nfa), nfa)
				if len(next) == 0 {
					continue
				}
				key := fmt.Sprint(next)
				if !seen[key] {
					seen[key] = true
					queue = append(queue, next)
				}
			}
		}
		return w6D(kp, fmt.Sprintf("A nondeterministic finite automaton has states 0..%d, start 0, accepting states %s, transitions %s, and epsilon edges %s. How many nonempty reachable subset states does the subset construction produce?", n-1, w6DJoin(accept), f["transitions"], f["epsilon"]), strconv.Itoa(len(seen)),
			fmt.Sprintf("Start from the epsilon closure of 0. Explore the closure of every move on a and b. The reachable subset states number %d.", len(seen)), f)
	case 2:
		length := w6DR(2, 4)
		word := make([]byte, length)
		for i := range word {
			if rand.Intn(2) == 0 {
				word[i] = 'a'
			} else {
				word[i] = 'b'
			}
		}
		cur := w6NFAClosure([]int{0}, nfa)
		for _, ch := range word {
			sym := 0
			if ch == 'b' {
				sym = 1
			}
			cur = w6NFAClosure(w6NFAMove(cur, sym, nfa), nfa)
		}
		accepted := 0
		for _, q := range cur {
			for _, a := range accept {
				if q == a {
					accepted = 1
				}
			}
		}
		f["string"] = string(word)
		return w6D(kp, fmt.Sprintf("For the NFA with states 0..%d, start 0, accepting states %s, transitions %s, and epsilon edges %s, is the string %q accepted? Enter 1 for yes and 0 for no.", n-1, w6DJoin(accept), f["transitions"], f["epsilon"], string(word)), strconv.Itoa(accepted),
			fmt.Sprintf("Simulate the NFA on %q using epsilon closures. The string is %s.", string(word), map[int]string{0: "rejected", 1: "accepted"}[accepted]), f)
	default:
		target := w6DR(0, n-1)
		closure := w6NFAClosure([]int{target}, nfa)
		w6DI(f, "target", target)
		noun := "states"
		if len(closure) == 1 {
			noun = "state"
		}
		return w6D(kp, fmt.Sprintf("For the NFA with epsilon edges %s, what is the size of the epsilon closure of state %d?", f["epsilon"], target), strconv.Itoa(len(closure)),
			fmt.Sprintf("Follow epsilon edges from %d. The closure contains %d %s.", target, len(closure), noun), f)
	}
}

// ---------------------------------------------------------------------------
// Regular expressions and finite automata
// ---------------------------------------------------------------------------

func w6DFAString(trans [][2]int) string {
	parts := make([]string, 0, len(trans)*2)
	for q := range trans {
		parts = append(parts, fmt.Sprintf("%da:%d", q, trans[q][0]), fmt.Sprintf("%db:%d", q, trans[q][1]))
	}
	return strings.Join(parts, ";")
}

func w6DFAParse(s string) [][2]int {
	var out [][2]int
	for _, part := range strings.Split(s, ";") {
		if part == "" {
			continue
		}
		colon := strings.Index(part, ":")
		if colon < 0 {
			continue
		}
		head, tail := part[:colon], part[colon+1:]
		q, _ := strconv.Atoi(head[:len(head)-1])
		sym := 0
		if strings.HasSuffix(head, "b") {
			sym = 1
		}
		r, _ := strconv.Atoi(tail)
		for len(out) <= q {
			out = append(out, [2]int{})
		}
		out[q][sym] = r
	}
	return out
}

func w6DFAAccepts(trans [][2]int, accept map[int]bool, word string) int {
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

func w6DFACount(trans [][2]int, accept map[int]bool, length int) int {
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

func w6Regex(kp int) generator.Problem {
	n := w6DR(2, 3)
	trans := make([][2]int, n)
	for q := range trans {
		trans[q][0] = rand.Intn(n)
		trans[q][1] = rand.Intn(n)
	}
	accept := map[int]bool{}
	var acceptList []int
	for q := 0; q < n; q++ {
		if rand.Intn(2) == 0 {
			accept[q] = true
			acceptList = append(acceptList, q)
		}
	}
	if len(acceptList) == 0 {
		accept[n-1] = true
		acceptList = []int{n - 1}
	}
	f := map[string]string{}
	w6DI(f, "n", n)
	f["transitions"] = w6DFAString(trans)
	f["accept"] = w6DJoin(acceptList)
	switch kp {
	case 1:
		length := w6DR(2, 4)
		w6DI(f, "length", length)
		ans := w6DFACount(trans, accept, length)
		return w6D(kp, fmt.Sprintf("A DFA has states 0..%d, start 0, accepting states %s and transitions %s. How many strings of length exactly %d does it accept?", n-1, w6DJoin(acceptList), f["transitions"], length), strconv.Itoa(ans),
			fmt.Sprintf("Propagate the count of words through the DFA for %d steps and add the accepting states. The count is %d.", length, ans), f)
	case 2:
		length := w6DR(2, 4)
		word := make([]byte, length)
		for i := range word {
			if rand.Intn(2) == 0 {
				word[i] = 'a'
			} else {
				word[i] = 'b'
			}
		}
		f["string"] = string(word)
		ans := w6DFAAccepts(trans, accept, string(word))
		return w6D(kp, fmt.Sprintf("For the DFA with states 0..%d, start 0, accepting states %s and transitions %s, is %q accepted? Enter 1 for yes and 0 for no.", n-1, w6DJoin(acceptList), f["transitions"], string(word)), strconv.Itoa(ans),
			fmt.Sprintf("Run the DFA on %q and inspect the final state. The string is %s.", string(word), map[int]string{0: "rejected", 1: "accepted"}[ans]), f)
	default:
		length := w6DR(2, 4)
		w6DI(f, "length", length)
		total := 0
		for l := 0; l <= length; l++ {
			total += w6DFACount(trans, accept, l)
		}
		return w6D(kp, fmt.Sprintf("For the DFA with states 0..%d, start 0, accepting states %s and transitions %s, how many strings of length at most %d does it accept?", n-1, w6DJoin(acceptList), f["transitions"], length), strconv.Itoa(total),
			fmt.Sprintf("Sum the accepted counts for lengths 0 through %d. The total is %d.", length, total), f)
	}
}

// ---------------------------------------------------------------------------
// Linear codes: the Hamming [7,4] code
// ---------------------------------------------------------------------------

func w6HammingCodeword(msg int) int {
	m := [4]int{(msg >> 3) & 1, (msg >> 2) & 1, (msg >> 1) & 1, msg & 1}
	c := [7]int{m[0], m[1], m[2], m[3], 0, 0, 0}
	c[4] = (m[0] + m[1] + m[3]) % 2
	c[5] = (m[0] + m[2] + m[3]) % 2
	c[6] = (m[1] + m[2] + m[3]) % 2
	out := 0
	for i := 0; i < 7; i++ {
		out = out<<1 | c[i]
	}
	return out
}

func w6HammingDistance(a, b int) int {
	x := a ^ b
	d := 0
	for x > 0 {
		d += x & 1
		x >>= 1
	}
	return d
}

func w6LinearCode(kp int) generator.Problem {
	f := map[string]string{}
	msg := rand.Intn(16)
	cw := w6HammingCodeword(msg)
	w6DI(f, "message", msg)
	w6DI(f, "codeword", cw)
	switch kp {
	case 1:
		return w6D(kp, fmt.Sprintf("The Hamming [7,4] code has generator matrix with rows 1000110, 0100101, 0010011 and 0001111. Encode the 4-bit message %04b. Enter the 7-bit codeword as an integer.", msg), strconv.Itoa(cw),
			fmt.Sprintf("The first four bits repeat the message and the last three are the parity checks. The codeword is %07b, which is %d.", cw, cw), f)
	case 2:
		msg2 := rand.Intn(16)
		for msg2 == msg {
			msg2 = rand.Intn(16)
		}
		cw2 := w6HammingCodeword(msg2)
		w6DI(f, "message2", msg2)
		w6DI(f, "codeword2", cw2)
		d := w6HammingDistance(cw, cw2)
		return w6D(kp, fmt.Sprintf("For the Hamming [7,4] code, encode messages %04b and %04b. What is the Hamming distance between the two codewords?", msg, msg2), strconv.Itoa(d),
			fmt.Sprintf("The codewords are %07b and %07b. They differ in %d positions.", cw, cw2, d), f)
	default:
		bit := rand.Intn(7)
		received := cw ^ (1 << bit)
		w6DI(f, "received", received)
		w6DI(f, "flipped_bit", bit)
		return w6D(kp, fmt.Sprintf("A Hamming [7,4] codeword was transmitted and exactly one bit was flipped, giving the 7-bit word %07b. What was the original 4-bit message, as an integer from 0 to 15?", received), strconv.Itoa(msg),
			fmt.Sprintf("Compare the received word with the 16 codewords. The unique codeword at distance 1 is %07b, whose message is %04b, or %d.", cw, msg, msg), f)
	}
}

// ---------------------------------------------------------------------------
// Dynamic programming: coin change
// ---------------------------------------------------------------------------

func w6MinCoins(coins []int, t int) (int, int) {
	const inf = 1 << 20
	dp := make([]int, t+1)
	for i := 1; i <= t; i++ {
		dp[i] = inf
	}
	for i := 1; i <= t; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	if dp[t] >= inf {
		return -1, 0
	}
	return dp[t], dp[t]
}

func w6CountWays(coins []int, t int) int {
	dp := make([]int, t+1)
	dp[0] = 1
	for _, c := range coins {
		for i := c; i <= t; i++ {
			dp[i] += dp[i-c]
		}
	}
	return dp[t]
}

func w6DP(kp int) generator.Problem {
	coins := []int{1}
	extras := []int{2, 3, 5}
	perm := rand.Perm(len(extras))
	for i := 0; i < 1+rand.Intn(3); i++ {
		coins = append(coins, extras[perm[i]])
	}
	sort.Ints(coins)
	target := w6DR(5, 12)
	f := map[string]string{"coins": w6DJoin(coins)}
	w6DI(f, "target", target)
	switch kp {
	case 1:
		t := w6DR(1, target)
		w6DI(f, "t", t)
		ans, _ := w6MinCoins(coins, t)
		noun := "coins"
		if ans == 1 {
			noun = "coin"
		}
		return w6D(kp, fmt.Sprintf("With unlimited coins of denominations %s, let dp[t] be the minimum number of coins needed to make t. What is dp[%d]?", w6DJoin(coins), t), strconv.Itoa(ans),
			fmt.Sprintf("Build dp from 1 to %d. The minimum for %d is %d %s.", t, t, ans, noun), f)
	case 2:
		ans, _ := w6MinCoins(coins, target)
		return w6D(kp, fmt.Sprintf("With unlimited coins of denominations %s, what is the minimum number of coins needed to make %d?", w6DJoin(coins), target), strconv.Itoa(ans),
			fmt.Sprintf("The dynamic program for the minimum coin count gives %d coins for %d.", ans, target), f)
	default:
		ans := w6CountWays(coins, target)
		return w6D(kp, fmt.Sprintf("With unlimited coins of denominations %s, how many ways are there to make %d? The order of the coins does not matter.", w6DJoin(coins), target), strconv.Itoa(ans),
			fmt.Sprintf("Use the unbounded-coin recurrence and treat each denomination in turn. The count is %d.", ans), f)
	}
}
