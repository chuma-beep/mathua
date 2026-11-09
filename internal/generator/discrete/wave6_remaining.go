package discrete

// Wave 6 discrete mathematics, remaining batch: automata, computability,
// complexity, coding, combinatorics and greedy algorithms. Every answer is an
// integer and every instance is published in Facts for the independent tests.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"discrete.automata.pumping_lemma":               "pumping",
		"discrete.automata.cfg":                         "cfg",
		"discrete.automata.pda":                         "pda",
		"discrete.automata.turing":                      "turing",
		"discrete.automata.decidability":                "decidability",
		"discrete.automata.halting":                     "halting",
		"discrete.complexity.reduction":                 "reduction",
		"discrete.complexity.sat":                       "sat",
		"discrete.complexity.np_complete":               "np_complete",
		"discrete.complexity.cook_levin":                "cook_levin",
		"discrete.coding.generator":                     "generator_matrix",
		"discrete.coding.cyclic":                        "cyclic_code",
		"discrete.coding.reed_solomon":                  "reed_solomon",
		"discrete.combinatorics.rook":                   "rook",
		"discrete.combinatorics.young_tableaux":         "young",
		"discrete.combinatorics.q_binomial":             "q_binomial",
		"discrete.combinatorics.permutation_statistics": "permutation_stats",
		"discrete.combinatorics.schur":                  "schur",
		"discrete.algorithms.greedy":                    "greedy",
	} {
		reg.Register(id, &w6DiscreteRemainingGen{kind: kind})
	}
}

type w6DiscreteRemainingGen struct{ kind string }

func (g *w6DiscreteRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "pumping":
		return w6Pumping(kp)
	case "cfg":
		return w6CFG(kp)
	case "pda":
		return w6PDA(kp)
	case "turing":
		return w6Turing(kp)
	case "decidability":
		return w6Decidability(kp)
	case "halting":
		return w6Halting(kp)
	case "reduction":
		return w6Reduction(kp)
	case "sat":
		return w6SAT(kp)
	case "np_complete":
		return w6NPComplete(kp)
	case "cook_levin":
		return w6CookLevin(kp)
	case "generator_matrix":
		return w6GeneratorMatrix(kp)
	case "cyclic_code":
		return w6CyclicCode(kp)
	case "reed_solomon":
		return w6ReedSolomon(kp)
	case "rook":
		return w6Rook(kp)
	case "young":
		return w6Young(kp)
	case "q_binomial":
		return w6QBinomial(kp)
	case "permutation_stats":
		return w6PermutationStats(kp)
	case "schur":
		return w6Schur(kp)
	case "greedy":
		return w6Greedy(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6DM(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w6DRemInt(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6RemPick(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6Pumping(kp int) generator.Problem {
	n := w6RemPick(2, 14)
	f := map[string]string{}
	w6DRemInt(f, "states", n)
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("A deterministic finite automaton has %d states. What pumping length does the pumping lemma guarantee?", n), strconv.Itoa(n),
			"The pumping lemma can use the number of states as its pumping length.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For the word a^%d b^%d, how many starting positions can the pumped block have inside the first %d letters?", n, n, n), strconv.Itoa(n),
			"The pumping lemma lets the pumped block start at any of the first p positions.", f)
	default:
		return w6DM(kp, fmt.Sprintf("In the proof that {a^k b^k} is not regular, pump a^%d b^%d once more on the a-block. By how much do the numbers of a and b then differ?", n, n), strconv.Itoa(n),
			"Pumping the a-block changes only the a count, so the imbalance is the pumped amount.", f)
	}
}

func w6CFG(kp int) generator.Problem {
	n := w6RemPick(2, 14)
	f := map[string]string{"n": strconv.Itoa(n), "grammar": "S -> aSb | epsilon"}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("For the grammar S -> aSb | epsilon, how many derivation steps produce a^%d b^%d?", n, n), strconv.Itoa(n+1),
			"Each matched pair uses one S -> aSb step and the final step uses S -> epsilon.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For the grammar S -> aSb | epsilon, how many internal nodes does the parse tree for a^%d b^%d have?", n, n), strconv.Itoa(n),
			"There is one internal S node for each matched pair.", f)
	default:
		m := w6RemPick(1, 6)
		f["m"] = strconv.Itoa(m)
		ok := 0
		if n == m {
			ok = 1
		}
		return w6DM(kp, fmt.Sprintf("Does the grammar S -> aSb | epsilon generate a^%d b^%d? Enter 1 for yes and 0 for no.", n, m), strconv.Itoa(ok),
			"The grammar generates exactly the words with equal numbers of a and b.", f)
	}
}

func w6PDA(kp int) generator.Problem {
	n := w6RemPick(2, 6)
	m := w6RemPick(1, 6)
	f := map[string]string{"n": strconv.Itoa(n), "m": strconv.Itoa(m)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("A pushdown automaton reads a^%d b^%d by pushing on every a and popping on every b. How many pushes occur before the first pop?", n, m), strconv.Itoa(n),
			"Every a is pushed before the first b is read.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For the same PDA reading a^%d b^%d, how many pops occur in total?", n, m), strconv.Itoa(minInt(n, m)),
			"A pop occurs only when a b is read and the stack is nonempty.", f)
	default:
		return w6DM(kp, fmt.Sprintf("For the same PDA reading a^%d b^%d, what is the final stack height?", n, m), strconv.Itoa(absInt(n-m)),
			"The stack height is the number of pushes minus pops.", f)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func w6Turing(kp int) generator.Problem {
	n := w6RemPick(3, 16)
	ones := w6RemPick(0, n)
	f := map[string]string{"n": strconv.Itoa(n), "ones": strconv.Itoa(ones)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("A Turing machine scans %d tape cells from left to right and halts on the next blank. How many transitions does it make?", n), strconv.Itoa(n+1),
			"There is one transition per cell plus the transition that reads the blank and halts.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("The machine flips every one of the %d cells it scans. How many cells does it flip?", n), strconv.Itoa(n),
			"Each of the n scanned cells is flipped once.", f)
	default:
		return w6DM(kp, fmt.Sprintf("A tape of %d cells contains %d ones. After flipping every cell, how many ones remain?", n, ones), strconv.Itoa(n-ones),
			"Flipping every cell changes each one to zero and each zero to one, so n minus the original ones remain.", f)
	}
}

func w6Decidability(kp int) generator.Problem {
	names := []string{"a^n b^n", "the halting problem", "prime numbers", "the set of all Turing machines that halt on themselves", "the set of binary palindromes", "the emptiness problem for DFAs", "the equivalence problem for DFAs", "the acceptance problem for Turing machines", "the set of strings of even length", "the Post correspondence problem"}
	codes := []int{1, 0, 1, 0, 1, 1, 1, 0, 1, 0}
	i := rand.Intn(len(names))
	f := map[string]string{"language": names[i]}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("Classify the language %q: enter 1 for decidable and 0 for undecidable.", names[i]), strconv.Itoa(codes[i]),
			"The classification follows from whether a finite algorithm decides membership for every input.", f)
	case 2:
		n := w6RemPick(2, 9)
		f["steps"] = strconv.Itoa(n)
		return w6DM(kp, fmt.Sprintf("A reduction from A to B runs in %d steps per input. If B is decidable, is A necessarily decidable? Enter 1 for yes and 0 for no.", n), "1",
			"Compose the reduction with the decider for B.", f)
	default:
		n := w6RemPick(2, 9)
		f["pair"] = strconv.Itoa(n)
		return w6DM(kp, fmt.Sprintf("For decidable languages L_%d and its complement, is the complement decidable? Enter 1 for yes and 0 for no.", n), "1",
			"Swap the accepting and rejecting states of the decider.", f)
	}
}

func w6Halting(kp int) generator.Problem {
	n := w6RemPick(3, 9)
	f := map[string]string{"n": strconv.Itoa(n)}
	switch kp {
	case 1:
		k := w6RemPick(1, 9)
		f["case"] = strconv.Itoa(k)
		return w6DM(kp, fmt.Sprintf("In case %d of the diagonal argument for the halting problem, a program loops exactly when its own code halts. What truth value does the construction force? Enter 1 for contradiction and 0 for consistent.", k), "1",
			"The construction makes the hypothetical decider wrong on its own code.", f)
	case 2:
		k := w6RemPick(2, 9)
		f["overhead"] = strconv.Itoa(k)
		return w6DM(kp, fmt.Sprintf("A reduction from the halting problem to problem B adds %d steps per simulated step. If a decider for B uses %d steps, how many steps does the composed decider use?", k, n), strconv.Itoa(k*n),
			"Multiply the reduction overhead by the decider's step count.", f)
	default:
		return w6DM(kp, fmt.Sprintf("Is halting decidable for programs that run at most %d steps? Enter 1 for yes and 0 for no.", n), "1",
			"A simulator can run the program for n steps and stop; no unbounded search is needed.", f)
	}
}

func w6Reduction(kp int) generator.Problem {
	n := w6RemPick(2, 12)
	f := map[string]string{"clauses": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("The standard reduction from 3-SAT with %d clauses to the clique problem creates how many clause-vertices?", n), strconv.Itoa(3*n),
			"Each clause contributes one vertex per literal.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For a satisfiable 3-SAT instance with %d clauses, what clique size does the reduction guarantee?", n), strconv.Itoa(n),
			"Choose one true literal from each clause; the chosen vertices form a clique of size equal to the number of clauses.", f)
	default:
		a := w6RemPick(2, 6)
		b := w6RemPick(2, 6)
		f["factor1"] = strconv.Itoa(a)
		f["factor2"] = strconv.Itoa(b)
		return w6DM(kp, fmt.Sprintf("Two reductions have overhead factors %d and %d. What is the overhead factor of their composition?", a, b), strconv.Itoa(a*b),
			"The factors multiply under composition.", f)
	}
}

func w6SAT(kp int) generator.Problem {
	clauses := make([][3]int, 3)
	for i := range clauses {
		for j := range clauses[i] {
			lit := w6RemPick(1, 3)
			if rand.Intn(2) == 0 {
				lit = -lit
			}
			clauses[i][j] = lit
		}
	}
	parts := make([]string, len(clauses))
	for i, c := range clauses {
		parts[i] = fmt.Sprintf("%d,%d,%d", c[0], c[1], c[2])
	}
	formula := strings.Join(parts, ";")
	countSat := func() int {
		count := 0
		for mask := 0; mask < 8; mask++ {
			ok := true
			for _, c := range clauses {
				sat := false
				for _, lit := range c {
					v := (mask >> (absInt(lit) - 1)) & 1
					if (lit < 0 && v == 0) || (lit > 0 && v == 1) {
						sat = true
					}
				}
				if !sat {
					ok = false
					break
				}
			}
			if ok {
				count++
			}
		}
		return count
	}
	f := map[string]string{"formula": formula}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("How many of the 8 truth assignments satisfy the 3-variable formula %s? (enter a number)", formula), strconv.Itoa(countSat()),
			"Enumerate the 8 assignments and check every clause.", f)
	case 2:
		mask := rand.Intn(8)
		f["assignment"] = strconv.Itoa(mask)
		ok := 0
		for _, c := range clauses {
			sat := false
			for _, lit := range c {
				v := (mask >> (absInt(lit) - 1)) & 1
				if (lit < 0 && v == 0) || (lit > 0 && v == 1) {
					sat = true
				}
			}
			if sat {
				ok++
			}
		}
		return w6DM(kp, fmt.Sprintf("For formula %s, how many clauses does the assignment with bit mask %d satisfy?", formula, mask), strconv.Itoa(ok),
			"Evaluate each clause under the assignment encoded by the mask.", f)
	default:
		best := 0
		for mask := 0; mask < 8; mask++ {
			ok := 0
			for _, c := range clauses {
				sat := false
				for _, lit := range c {
					v := (mask >> (absInt(lit) - 1)) & 1
					if (lit < 0 && v == 0) || (lit > 0 && v == 1) {
						sat = true
					}
				}
				if sat {
					ok++
				}
			}
			if ok > best {
				best = ok
			}
		}
		return w6DM(kp, fmt.Sprintf("For formula %s, what is the maximum number of clauses satisfied by any truth assignment?", formula), strconv.Itoa(best),
			"Check all 8 assignments and take the largest clause count.", f)
	}
}

func w6NPComplete(kp int) generator.Problem {
	k := w6RemPick(3, 12)
	f := map[string]string{"k": strconv.Itoa(k)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("A proposed certificate for a graph clique of size %d lists %d vertices. How many edges must the verifier check to confirm the clique?", k, k), strconv.Itoa(k*(k-1)/2),
			"Every pair of the k vertices must be checked.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For a reduction that asks for a clique of size %d, does a clique of that size imply the original formula is satisfiable? Enter 1 for yes and 0 for no.", k), "1",
			"The reduction preserves yes instances and no instances.", f)
	default:
		k2 := w6RemPick(3, 12)
		f["k"] = strconv.Itoa(k2)
		return w6DM(kp, fmt.Sprintf("Classify the clique problem for a graph on k=%d: enter 1 for P, 2 for NP-complete and 0 for unknown.", k2), "2",
			"Clique is NP-complete.", f)
	}
}

func w6CookLevin(kp int) generator.Problem {
	n := w6RemPick(3, 14)
	f := map[string]string{"n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("A Cook-Levin tableau has %d time steps and %d tape cells. How many cells are in the tableau?", n, n), strconv.Itoa(n*n),
			"The tableau is a time-by-space grid.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For a tableau with %d rows and %d columns, how many cell variables are needed?", n, n), strconv.Itoa(n*n),
			"One Boolean variable is used for each cell.", f)
	default:
		return w6DM(kp, fmt.Sprintf("Each of the %d rows in a Cook-Levin tableau contributes one initial, one accepting and one transition clause group. How many groups are listed?", n), strconv.Itoa(3*n),
			"Three clause groups per row.", f)
	}
}

func w6GeneratorMatrix(kp int) generator.Problem {
	// The binary [7,4] Hamming code has 16 messages and a fixed parity part.
	m := rand.Intn(16)
	f := map[string]string{"message": strconv.Itoa(m), "matrix": "1000110,0100101,0010011,0001111"}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("For the [7,4] Hamming code, encode message integer %d. Enter the codeword as an integer.", m), strconv.Itoa(w6HammingCodeword(m)),
			"Write the four message bits, then append the three parity bits.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For message integer %d in the [7,4] Hamming code, what is the code dimension k?", m), "4",
			"The code has four message positions and three parity positions.", f)
	default:
		return w6DM(kp, fmt.Sprintf("For message integer %d in the [7,4] Hamming code, how many parity-check equations are there?", m), "3",
			"The parity-check matrix has n-k = 3 rows.", f)
	}
}

func w6CyclicCode(kp int) generator.Problem {
	r := w6RemPick(2, 6)
	n := 7
	f := map[string]string{"n": strconv.Itoa(n), "r": strconv.Itoa(r)}
	switch kp {
	case 1:
		word := w6RemPick(1, 127)
		f["word"] = strconv.Itoa(word)
		rot := ((word << 1) & 127) | ((word >> 6) & 1)
		return w6DM(kp, fmt.Sprintf("In a cyclic code of length %d, rotate the 7-bit word %d one position left. Enter the resulting integer.", n, word), strconv.Itoa(rot),
			"Move the leftmost bit to the right end.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("A cyclic code has generator polynomial of degree %d. How many parity bits does it add to each message?", r), strconv.Itoa(r),
			"The degree of the generator polynomial is the number of parity bits.", f)
	default:
		return w6DM(kp, fmt.Sprintf("A cyclic code of length %d has generator polynomial degree %d. What is the message length k?", n, r), strconv.Itoa(n-r),
			"Subtract the parity length from the code length.", f)
	}
}

func w6ReedSolomon(kp int) generator.Problem {
	p := []int{7, 11, 13}[rand.Intn(3)]
	n := w6RemPick(5, 9)
	k := w6RemPick(2, n-2)
	a0, a1, a2 := w6RemPick(0, p-1), w6RemPick(0, p-1), w6RemPick(0, p-1)
	i := w6RemPick(1, n)
	f := map[string]string{"p": strconv.Itoa(p), "n": strconv.Itoa(n), "k": strconv.Itoa(k), "i": strconv.Itoa(i), "a0": strconv.Itoa(a0), "a1": strconv.Itoa(a1), "a2": strconv.Itoa(a2)}
	switch kp {
	case 1:
		val := (a0 + a1*i + a2*i*i) % p
		return w6DM(kp, fmt.Sprintf("A Reed-Solomon code over F_%d evaluates f(x)=%d+%dx+%dx^2 at x=%d. What is the value modulo %d?", p, a0, a1, a2, i, p), strconv.Itoa(val),
			"Substitute x and reduce modulo p.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("A Reed-Solomon code has length %d and dimension %d. What is its minimum distance?", n, k), strconv.Itoa(n-k+1),
			"The minimum distance is n-k+1.", f)
	default:
		return w6DM(kp, fmt.Sprintf("A Reed-Solomon code of length %d and dimension %d has minimum distance %d. How many errors can it correct?", n, k, n-k+1), strconv.Itoa((n-k)/2),
			"Correct floor((d-1)/2) errors.", f)
	}
}

func w6Rook(kp int) generator.Problem {
	m := w6RemPick(3, 8)
	n := w6RemPick(3, 8)
	f := map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("On a %d by %d board, how many ways are there to place one rook?", m, n), strconv.Itoa(m*n),
			"Choose any of the m*n cells.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("On a %d by %d board, how many ways are there to place two nonattacking rooks?", m, n), strconv.Itoa(m*(m-1)*n*(n-1)/2),
			"Choose two rows and two columns, then pair them in two ways.", f)
	default:
		d := w6RemPick(3, 9)
		f["d"] = strconv.Itoa(d)
		return w6DM(kp, fmt.Sprintf("How many permutations of %d items have no fixed point?", d), strconv.Itoa(derange(d)),
			"This is the derangement count, the rook-polynomial inclusion-exclusion on the forbidden diagonal.", f)
	}
}

func derange(n int) int {
	d := []int{1, 0, 1, 2, 9, 44, 265}
	if n < len(d) {
		return d[n]
	}
	return 0
}

func w6Young(kp int) generator.Problem {
	shapes := [][][]int{
		{{2, 1}}, {{2, 2}}, {{3, 1}}, {{3, 2}}, {{3, 2, 1}},
		{{4, 1}}, {{4, 2}}, {{4, 3}}, {{2, 2, 2}}, {{4, 2, 1}}, {{3, 3, 1}},
	}
	shape := shapes[rand.Intn(len(shapes))][0]
	n := 0
	for _, r := range shape {
		n += r
	}
	f := map[string]string{"shape": w6DJoin(shape), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("How many boxes does the Young diagram with row lengths %s have?", w6DJoin(shape)), strconv.Itoa(n),
			"Add the row lengths.", f)
	case 2:
		row := rand.Intn(len(shape))
		col := rand.Intn(shape[row])
		f["row"] = strconv.Itoa(row)
		f["col"] = strconv.Itoa(col)
		hook := (shape[row] - col) + countBelow(shape, row, col)
		return w6DM(kp, fmt.Sprintf("For the shape %s, what is the hook length of the cell in row %d, column %d (zero-indexed)?", w6DJoin(shape), row, col), strconv.Itoa(hook),
			"The hook length is the number of cells to the right plus below plus one.", f)
	default:
		count := youngCount(shape)
		f["count"] = strconv.Itoa(count)
		return w6DM(kp, fmt.Sprintf("How many standard Young tableaux have shape %s?", w6DJoin(shape)), strconv.Itoa(count),
			"The hook-length formula gives the count.", f)
	}
}

func countBelow(shape []int, row, col int) int {
	c := 0
	for r := row + 1; r < len(shape); r++ {
		if shape[r] > col {
			c++
		}
	}
	return c
}

func youngCount(shape []int) int {
	n := 0
	for _, r := range shape {
		n += r
	}
	num := fact6(n)
	den := 1
	for r := 0; r < len(shape); r++ {
		for c := 0; c < shape[r]; c++ {
			den *= (shape[r] - c) + countBelow(shape, r, c)
		}
	}
	return num / den
}

func fact6(n int) int {
	r := 1
	for i := 2; i <= n; i++ {
		r *= i
	}
	return r
}

func w6QBinomial(kp int) generator.Problem {
	n := w6RemPick(3, 6)
	k := w6RemPick(1, n-1)
	f := map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("Evaluate the Gaussian binomial coefficient [%d choose %d]_q at q=2.", n, k), strconv.Itoa(qBinom(n, k)),
			"Use the recurrence [n choose k]_q = [n-1 choose k]_q + q^(n-k)[n-1 choose k-1]_q.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("How many lattice paths go from (0,0) to (%d,%d) using right and up steps?", k, n-k), strconv.Itoa(binom6(n, k)),
			"The Gaussian binomial at q=1 counts these paths.", f)
	default:
		return w6DM(kp, fmt.Sprintf("What is the value of [%d choose %d]_q at q=0?", n, k), "1",
			"At q=0 every Gaussian binomial coefficient is 1.", f)
	}
}

func binom6(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func qBinom(n, k int) int {
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
				dp[i][j] += intPow6(2, i-j) * dp[i-1][j-1]
			}
		}
	}
	return dp[n][k]
}

func intPow6(a, e int) int {
	r := 1
	for i := 0; i < e; i++ {
		r *= a
	}
	return r
}

func w6PermutationStats(kp int) generator.Problem {
	n := w6RemPick(4, 14)
	perm := rand.Perm(n)
	inv := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if perm[i] > perm[j] {
				inv++
			}
		}
	}
	major := 0
	for i := 0; i+1 < n; i++ {
		if perm[i] > perm[i+1] {
			major += i + 1
		}
	}
	f := map[string]string{"n": strconv.Itoa(n), "permutation": w6DJoin(perm)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("What is the inversion number of the permutation %s?", w6DJoin(perm)), strconv.Itoa(inv),
			"Count pairs i<j with p_i>p_j.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("What is the major index of the permutation %s?", w6DJoin(perm)), strconv.Itoa(major),
			"Sum the positions i with p_i>p_{i+1}.", f)
	default:
		return w6DM(kp, fmt.Sprintf("How many permutations of %d items have exactly one inversion?", n), strconv.Itoa(n-1),
			"A single inversion is a descent of adjacent values in one position.", f)
	}
}

func w6Schur(kp int) generator.Problem {
	r := w6RemPick(2, 4)
	s := map[int]int{2: 5, 3: 14, 4: 45}[r]
	f := map[string]string{"r": strconv.Itoa(r)}
	switch kp {
	case 1:
		c := w6RemPick(0, 6)
		f["offset"] = strconv.Itoa(c)
		return w6DM(kp, fmt.Sprintf("What is the Schur number S(%d) plus %d?", r, c), strconv.Itoa(s+c),
			"S(2)=5, S(3)=14 and S(4)=45; add the stated offset.", f)
	case 2:
		c := w6RemPick(0, 6)
		f["offset"] = strconv.Itoa(c)
		return w6DM(kp, fmt.Sprintf("What is the largest n for which [1,n] has a valid %d-coloring with no monochromatic x+y=z, plus %d?", r, c), strconv.Itoa(s-1+c),
			"The largest such n is S(r)-1; add the stated offset.", f)
	default:
		n := w6RemPick(4, 10)
		f["n"] = strconv.Itoa(n)
		return w6DM(kp, fmt.Sprintf("In the 2-coloring of [1,%d] that colors odd numbers red and even numbers blue, how many monochromatic Schur triples are there?", n), strconv.Itoa(countSchur(n)),
			"Check x+y=z with x,y,z all the same color.", f)
	}
}

func countSchur(n int) int {
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

func w6Greedy(kp int) generator.Problem {
	target := w6RemPick(6, 15)
	g := greedyCoins(target)
	opt := minCoins(target)
	f := map[string]string{"target": strconv.Itoa(target)}
	switch kp {
	case 1:
		return w6DM(kp, fmt.Sprintf("For coins {1,3,4} and target %d, how many coins does the greedy algorithm use?", target), strconv.Itoa(g),
			"Take as many 4s as possible, then 3s, then 1s.", f)
	case 2:
		return w6DM(kp, fmt.Sprintf("For coins {1,3,4} and target %d, what is the minimum number of coins?", target), strconv.Itoa(opt),
			"Dynamic programming finds the optimum.", f)
	default:
		ok := 0
		if g == opt {
			ok = 1
		}
		return w6DM(kp, fmt.Sprintf("For coins {1,3,4} and target %d, is the greedy algorithm optimal? Enter 1 for yes and 0 for no.", target), strconv.Itoa(ok),
			"Compare the greedy count with the dynamic-programming optimum.", f)
	}
}

func greedyCoins(t int) int {
	c := 0
	c += t / 4
	t %= 4
	c += t / 3
	t %= 3
	return c + t
}

func minCoins(t int) int {
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
