package discrete

// Wave 2g of the generator variety pass.
//
// Fourteen discrete-mathematics concepts answered at least one knowledge point
// from a fixed claim table, a nine-value parameter pair or a generator that
// ignored the knowledge-point ordinal. Each one now draws its instance from a
// range wide enough that a single ordinal supplies more than ten distinct
// questions, and each branch serves the material its lesson section teaches.
//
// The maths is unchanged where it was right. One table was not: burnsideProblem
// answered 21 for the square with three colours and 100 for four colours, while
// the rotation group gives 24 and 70. Those entries were the dihedral count for
// three colours and an unidentifiable number for four, so the branch now
// computes the rotation count from the group action rather than reading a
// table.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// w2gN draws an integer in [lo, hi].
func w2gN(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w2gGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func w2gPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

// w2gNecklace counts the colourings of n positions with k colours up to
// rotation, by Burnside's lemma over the cyclic group: the average of the
// colourings fixed by each rotation, which is k^gcd(n,d) for a rotation by d.
func w2gNecklace(n, k int) int {
	total := 0
	for d := 0; d < n; d++ {
		total += w2gPow(k, w2gGCD(n, d))
	}
	return total / n
}

// w2gBooleanFunctions covers three knowledge points: the number of Boolean
// functions, the number of satisfying assignments of a two-variable formula and
// the same count for four variables. The formulas are drawn as expression
// trees, so the tables of eight and seven expressions are gone.
func w2gBooleanFunctions(o int) generator.Problem {
	switch o {
	case 1:
		n := w2gN(1, 5)
		answer := 1 << uint(1<<uint(n))
		var question string
		switch rand.Intn(3) {
		case 0:
			question = fmt.Sprintf("How many distinct Boolean functions of %d variable(s) are there? (enter a number)", n)
		case 1:
			question = fmt.Sprintf("A Boolean function of %d variable(s) is specified by its truth table. How many different such functions are there? (enter a number)", n)
		default:
			question = fmt.Sprintf("How many Boolean functions have %d inputs? (enter a number)", n)
		}
		return discProblem(o, question, strconv.Itoa(answer),
			fmt.Sprintf("There are 2^%d = %d truth-table rows, and each row can output 0 or 1, so the count is 2^(2^%d) = %d.", n, 1<<uint(n), n, answer))
	case 2:
		vars := []string{"x", "y"}
		expr := w2gRandomBoolExpr(vars, 2)
		ones := w2gSatisfying(expr, vars)
		return discProblem(o,
			fmt.Sprintf("For two Boolean variables x and y, how many of the four assignments make %s true? (enter a number)", expr.String()),
			strconv.Itoa(ones),
			fmt.Sprintf("Evaluating %s over the four assignments gives %d true rows.", expr.String(), ones))
	default:
		vars := []string{"x", "y", "z", "w"}
		expr := w2gRandomBoolExpr(vars, 3)
		ones := w2gSatisfying(expr, vars)
		return discProblem(o,
			fmt.Sprintf("For four Boolean variables x, y, z and w, how many of the sixteen assignments make %s true? (enter a number)", expr.String()),
			strconv.Itoa(ones),
			fmt.Sprintf("Evaluating %s over the sixteen assignments gives %d true rows.", expr.String(), ones))
	}
}

// w2gBoolExpr is a small Boolean expression tree. Leaves are a variable or its
// negation, internal nodes are a binary connective or a negation.
type w2gBoolExpr struct {
	leaf        string
	neg         bool
	op          string
	left, right *w2gBoolExpr
}

func (e *w2gBoolExpr) String() string {
	if e.left == nil {
		if e.neg {
			return "NOT " + e.leaf
		}
		return e.leaf
	}
	if e.op == "NOT" {
		return "NOT (" + e.left.String() + ")"
	}
	return "(" + e.left.String() + " " + e.op + " " + e.right.String() + ")"
}

func (e *w2gBoolExpr) eval(assign map[string]bool) bool {
	if e.left == nil {
		v := assign[e.leaf]
		if e.neg {
			return !v
		}
		return v
	}
	if e.op == "NOT" {
		return !e.left.eval(assign)
	}
	l, r := e.left.eval(assign), e.right.eval(assign)
	switch e.op {
	case "AND":
		return l && r
	case "OR":
		return l || r
	case "XOR":
		return l != r
	case "->":
		return !l || r
	case "<->":
		return l == r
	case "NAND":
		return !(l && r)
	case "NOR":
		return !(l || r)
	}
	return false
}

func w2gRandomBoolExpr(vars []string, depth int) *w2gBoolExpr {
	if depth <= 0 || rand.Intn(100) < 30 {
		return &w2gBoolExpr{leaf: vars[rand.Intn(len(vars))], neg: rand.Intn(3) == 0}
	}
	if rand.Intn(5) == 0 {
		return &w2gBoolExpr{op: "NOT", left: w2gRandomBoolExpr(vars, depth-1)}
	}
	ops := []string{"AND", "OR", "XOR", "->", "<->", "NAND", "NOR"}
	return &w2gBoolExpr{
		op:    ops[rand.Intn(len(ops))],
		left:  w2gRandomBoolExpr(vars, depth-1),
		right: w2gRandomBoolExpr(vars, depth-1),
	}
}

func w2gSatisfying(e *w2gBoolExpr, vars []string) int {
	ones := 0
	for mask := 0; mask < 1<<uint(len(vars)); mask++ {
		assign := make(map[string]bool, len(vars))
		for i, name := range vars {
			assign[name] = mask&(1<<uint(i)) != 0
		}
		if e.eval(assign) {
			ones++
		}
	}
	return ones
}

func w2gResolution(o int) generator.Problem {
	switch o {
	case 1:
		p, q := w2gN(2, 8), w2gN(2, 8)
		answer := p + q - 2
		return discProblem(o,
			fmt.Sprintf("A clause has %d literals and another has %d literals. They share exactly one complementary pair. How many literals are in their resolvent? (enter a number)", p, q),
			strconv.Itoa(answer),
			fmt.Sprintf("Delete the two complementary literals: %d + %d - 2 = %d.", p, q, answer))
	case 2:
		left, right := w2gResolutionClauses()
		pairs := 0
		for _, name := range []string{"x", "y", "z", "w"} {
			lp, lok := left.polarity[name]
			rp, rok := right.polarity[name]
			if lok && rok && lp != rp {
				pairs++
			}
		}
		return discProblem(o,
			fmt.Sprintf("How many resolvents can be formed from the clauses (%s) and (%s), counting one per complementary pair of literals? (enter a number)", left.text, right.text),
			strconv.Itoa(pairs),
			fmt.Sprintf("A resolvent needs a variable that is positive in one clause and negative in the other. The two clauses contain %d such pair(s), so they give %d resolvent(s).", pairs, pairs))
	default:
		clauses := w2gN(3, 5)
		units := 0
		parts := make([]string, 0, clauses)
		for i := 0; i < clauses; i++ {
			vars := []string{"x", "y", "z", "w"}
			rand.Shuffle(len(vars), func(a, b int) { vars[a], vars[b] = vars[b], vars[a] })
			width := w2gN(1, 3)
			lits := make([]string, 0, width)
			for _, name := range vars[:width] {
				if rand.Intn(2) == 0 {
					lits = append(lits, "NOT "+name)
				} else {
					lits = append(lits, name)
				}
			}
			if width == 1 {
				units++
			}
			parts = append(parts, "("+strings.Join(lits, " OR ")+")")
		}
		formula := strings.Join(parts, " AND ")
		return discProblem(o,
			fmt.Sprintf("How many unit clauses appear in the formula %s? (enter a number)", formula),
			strconv.Itoa(units),
			fmt.Sprintf("A unit clause has exactly one literal. Counting the single-literal clauses in the formula gives %d.", units))
	}
}

type w2gClause struct {
	text     string
	polarity map[string]bool
}

// w2gResolutionClauses builds two non-empty clauses over a random subset of
// x, y, z and w. Each variable appears at most once per clause, so the count of
// complementary pairs is the count of resolvents.
func w2gResolutionClauses() (w2gClause, w2gClause) {
	for {
		names := []string{"x", "y", "z", "w"}
		rand.Shuffle(len(names), func(a, b int) { names[a], names[b] = names[b], names[a] })
		n := w2gN(2, 4)
		names = names[:n]
		left := w2gClause{polarity: map[string]bool{}}
		right := w2gClause{polarity: map[string]bool{}}
		var leftLits, rightLits []string
		for _, name := range names {
			if rand.Intn(2) == 0 {
				pos := rand.Intn(2) == 0
				left.polarity[name] = pos
				leftLits = append(leftLits, w2gLiteral(name, pos))
			}
			if rand.Intn(2) == 0 {
				pos := rand.Intn(2) == 0
				right.polarity[name] = pos
				rightLits = append(rightLits, w2gLiteral(name, pos))
			}
		}
		if len(leftLits) == 0 || len(rightLits) == 0 {
			continue
		}
		left.text = strings.Join(leftLits, " OR ")
		right.text = strings.Join(rightLits, " OR ")
		return left, right
	}
}

func w2gLiteral(name string, positive bool) string {
	if positive {
		return name
	}
	return "NOT " + name
}

func w2gProofCases(o int) generator.Problem {
	switch o {
	case 1:
		m := w2gN(2, 15)
		return discProblem(o,
			fmt.Sprintf("A proof splits integers into cases according to their remainder modulo %d. How many cases are needed? (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("The possible remainders modulo %d are 0 through %d, so there are %d cases.", m, m-1, m))
	case 2:
		k := w2gN(1, 16)
		answer := w2gPow(2, k)
		return discProblem(o,
			fmt.Sprintf("A proof analyses the sign of a product of %d factors. Each factor is positive or negative. How many sign cases are there? (enter a number)", k),
			strconv.Itoa(answer),
			fmt.Sprintf("Each factor has two signs, so %d factors give 2^%d = %d cases.", k, k, answer))
	default:
		m := w2gN(1, 20)
		return discProblem(o,
			fmt.Sprintf("A proof splits on the remainder modulo %d and also on whether the integer is positive, negative or zero. How many cases are there? (enter a number)", m),
			strconv.Itoa(3*m),
			fmt.Sprintf("There are %d residues and three signs, so the total is 3 x %d = %d.", m, m, 3*m))
	}
}

func w2gBurnside(o int) generator.Problem {
	switch o {
	case 1:
		n, k := w2gN(3, 10), w2gN(2, 6)
		answer := w2gNecklace(n, k)
		return discProblem(o,
			fmt.Sprintf("How many distinct necklaces can be made from %d beads and %d colours, counting two necklaces as the same when one is a rotation of the other? (enter a number)", n, k),
			strconv.Itoa(answer),
			fmt.Sprintf("Burnside averages the fixed colourings over the %d rotations, which gives %d.", n, answer))
	case 2:
		n, k := w2gN(3, 12), w2gN(2, 4)
		d := rand.Intn(n)
		answer := w2gPow(k, w2gGCD(n, d))
		return discProblem(o,
			fmt.Sprintf("%d positions are coloured with %d colours. How many colourings are fixed by a rotation of the positions by %d step(s)? (enter a number)", n, k, d),
			strconv.Itoa(answer),
			fmt.Sprintf("The rotation splits the positions into gcd(%d, %d) = %d cycles, so %d to that power gives %d fixed colourings.", n, d, w2gGCD(n, d), k, answer))
	default:
		names := map[int]string{3: "triangle", 4: "square", 5: "pentagon", 6: "hexagon", 7: "heptagon", 8: "octagon"}
		n, k := w2gN(3, 8), w2gN(2, 6)
		answer := w2gNecklace(n, k)
		return discProblem(o,
			fmt.Sprintf("How many distinct colourings of the vertices of a regular %s with %d colours are there, up to rotation? (enter a number)", names[n], k),
			strconv.Itoa(answer),
			fmt.Sprintf("Burnside's lemma averages the colourings fixed by each of the %d rotations of the %s, which gives %d.", n, names[n], answer))
	}
}

func w2gPascal(o int) generator.Problem {
	switch o {
	case 1:
		n := w2gN(0, 14)
		row := pascalRow(n)
		parts := make([]string, len(row))
		for i, v := range row {
			parts[i] = strconv.Itoa(v)
		}
		answer := strings.Join(parts, ",")
		return discProblem(o,
			fmt.Sprintf("What is row %d of Pascal's triangle, counting the top 1 as row 0? (answer as comma-separated numbers)", n),
			answer,
			fmt.Sprintf("Row %d, counting the top 1 as row 0, is %s.", n, answer))
	case 2:
		n := w2gN(2, 14)
		k := w2gN(1, n-1)
		answer := nCr(n, k)
		return discProblem(o,
			fmt.Sprintf("What is C(%d,%d), an entry of row %d of Pascal's triangle counting the top 1 as row 0? (enter a number)", n, k, n),
			strconv.Itoa(answer),
			fmt.Sprintf("In row %d, counting the top 1 as row 0, the entry at position %d is C(%d,%d) = %d.", n, k, n, k, answer))
	default:
		n := w2gN(1, 16)
		answer := w2gPow(2, n)
		return discProblem(o,
			fmt.Sprintf("What is the sum of the entries in row %d of Pascal's triangle, counting the top 1 as row 0? (enter a number)", n),
			strconv.Itoa(answer),
			fmt.Sprintf("Row %d, counting the top 1 as row 0, has %d entries and their sum is 2^%d = %d.", n, n+1, n, answer))
	}
}

func w2gPermutations(o int) generator.Problem {
	switch o {
	case 1:
		n := w2gN(2, 14)
		answer := factorial(n)
		return discProblem(o,
			fmt.Sprintf("How many ways can %d distinct books be arranged in a row on a shelf? (enter a number)", n),
			strconv.Itoa(answer),
			fmt.Sprintf("Every ordering of the %d distinct books is a permutation, so there are %d! = %d.", n, n, answer))
	case 2:
		n, k := w2gN(4, 15), w2gN(2, 4)
		answer := 1
		for i := 0; i < k; i++ {
			answer *= n - i
		}
		return discProblem(o,
			fmt.Sprintf("How many ways can %d distinct prizes be awarded to %d contestants, with at most one prize per contestant? (enter a number)", k, n),
			strconv.Itoa(answer),
			fmt.Sprintf("The prizes go to an ordered selection of %d of the %d contestants, so the count is P(%d,%d) = %d x ... x %d = %d.", k, n, n, k, n, n-k+1, answer))
	default:
		n := w2gN(3, 15)
		answer := factorial(n - 1)
		return discProblem(o,
			fmt.Sprintf("How many ways can %d people sit around a circular table, counting rotations of the table as the same seating? (enter a number)", n),
			strconv.Itoa(answer),
			fmt.Sprintf("Fix one person to remove the %d rotations, then arrange the other %d people: %d! = %d.", n, n-1, n-1, answer))
	}
}

func w2gMasterTheorem(o int) generator.Problem {
	switch o {
	case 1:
		b := w2gN(2, 9)
		p := w2gN(1, 4)
		a := w2gPow(b, p)
		return discProblem(o,
			fmt.Sprintf("A recurrence has the form T(n) = %d T(n/%d) + f(n). What is log base %d of %d? (enter a number)", a, b, b, a),
			strconv.Itoa(p),
			fmt.Sprintf("%d^%d = %d, so log base %d of %d = %d.", b, p, a, b, a, p))
	case 2:
		b := w2gN(2, 5)
		k := w2gN(0, 5)
		bk := w2gPow(b, k)
		a := w2gN(2, bk+4)
		class := 2
		verdict := "The two terms are equal"
		switch {
		case a > bk:
			class = 1
			verdict = "The recursive term dominates the extra work"
		case a < bk:
			class = 3
			verdict = "The extra work dominates the recursive term"
		}
		return discProblem(o,
			fmt.Sprintf("For T(n) = %d T(n/%d) + n^%d, which master-theorem case applies? Enter 1 if the recursive term dominates, 2 for a tie, or 3 if the extra work dominates.", a, b, k),
			strconv.Itoa(class),
			fmt.Sprintf("Compare a = %d with b^k = %d^%d = %d. %s, so the comparison gives case %d.", a, b, k, bk, verdict, class))
	default:
		b := w2gN(2, 6)
		e := w2gN(1, 4)
		a := w2gPow(b, e)
		k := w2gN(0, 5)
		p := e
		if k > p {
			p = k
		}
		return discProblem(o,
			fmt.Sprintf("For T(n) = %d T(n/%d) + n^%d, the solution is Theta(n^p) times at most a logarithmic factor. What is p? (enter a number)", a, b, k),
			strconv.Itoa(p),
			fmt.Sprintf("The exponent p is the larger of log base %d of %d = %d and %d, which is %d.", b, a, e, k, p))
	}
}

func w2gRecurrence(o int) generator.Problem {
	switch o {
	case 1:
		d := w2gN(2, 9)
		first := w2gN(1, 12)
		n := w2gN(3, 9)
		answer := first + (n-1)*d
		return discProblem(o,
			fmt.Sprintf("If a(n) = a(n-1) + %d and a(1) = %d, what is a(%d)? (enter a number)", d, first, n),
			strconv.Itoa(answer),
			fmt.Sprintf("Each step adds %d, so a(%d) = a(1) + %d x %d = %d + %d = %d.", d, n, n-1, d, first, (n-1)*d, answer))
	case 2:
		r := w2gN(2, 4)
		first := w2gN(1, 8)
		n := w2gN(3, 6)
		answer := first * w2gPow(r, n-1)
		return discProblem(o,
			fmt.Sprintf("If a(n) = %d a(n-1) and a(1) = %d, what is a(%d)? (enter a number)", r, first, n),
			strconv.Itoa(answer),
			fmt.Sprintf("Each step multiplies by %d, so a(%d) = %d x %d^%d = %d.", r, n, first, r, n-1, answer))
	default:
		first := w2gN(1, 5)
		second := w2gN(1, 5)
		n := w2gN(4, 10)
		a, b := first, second
		for i := 3; i <= n; i++ {
			a, b = b, a+b
		}
		answer := b
		return discProblem(o,
			fmt.Sprintf("If a(n) = a(n-1) + a(n-2), a(1) = %d and a(2) = %d, what is a(%d)? (enter a number)", first, second, n),
			strconv.Itoa(answer),
			fmt.Sprintf("Add the two previous terms at each step, starting from a(1) = %d and a(2) = %d, to reach a(%d) = %d.", first, second, n, answer))
	}
}

func w2gEntropy(o int) generator.Problem {
	switch o {
	case 1:
		k := w2gN(2, 20)
		return w2gEntropyProblem(o, k,
			fmt.Sprintf("A source has %d outcomes. How many probabilities enter its entropy? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("One probability per outcome, so %d probabilities enter the entropy of the source.", k))
	case 2:
		k := 1 << uint(w2gN(1, 10))
		bits := 0
		for v := 1; v < k; v *= 2 {
			bits++
		}
		if rand.Intn(2) == 0 {
			return w2gEntropyProblem(o, k,
				fmt.Sprintf("A uniform source has %d equally likely outcomes. How many bits does each symbol carry? (enter a number)", k),
				strconv.Itoa(bits),
				fmt.Sprintf("The entropy is log base two of %d = %d bits, because 2^%d = %d.", k, bits, bits, k))
		}
		return w2gEntropyProblem(o, k,
			fmt.Sprintf("A uniform source has %d equally likely outcomes, so its entropy is log base two of %d bits. How many bits is that? (enter a number)", k, k),
			strconv.Itoa(bits),
			fmt.Sprintf("log base two of %d = %d, because 2^%d = %d.", k, bits, bits, k))
	default:
		k := w2gN(2, 30)
		var question string
		switch rand.Intn(3) {
		case 0:
			question = fmt.Sprintf("The probabilities of %d outcomes sum to what value? (enter a number)", k)
		case 1:
			question = fmt.Sprintf("A source with %d outcomes has one probability for each. What is the sum of those %d probabilities? (enter a number)", k, k)
		default:
			question = fmt.Sprintf("The probability masses of the %d outcomes of a source form a distribution. By how much must they total? (enter a number)", k)
		}
		return w2gEntropyProblem(o, k, question, "1",
			fmt.Sprintf("A probability distribution is normalised, so its %d masses sum to 1.", k))
	}
}

func w2gEntropyProblem(o, k int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer, "k": strconv.Itoa(k)},
		KPOrdinal:   o,
	}
}

func w2gComplement(o int) generator.Problem {
	m := w2gN(2, 9)
	length := w2gN(2, 9)
	total := w2gPow(m, length)
	none := w2gPow(m-1, length)
	answer := total - none
	var question string
	switch o {
	case 1:
		question = fmt.Sprintf("How many strings of length %d over an alphabet of %d letters contain at least one \u2018a\u2019? (enter a number)", length, m)
	case 2:
		question = fmt.Sprintf("A string of length %d is built from an alphabet of %d letters, one of which is \u2018a\u2019. How many such strings contain \u2018a\u2019 at least once? (enter a number)", length, m)
	default:
		question = fmt.Sprintf("Of the strings of length %d over an alphabet of %d letters, how many contain the letter \u2018a\u2019 at least once? (enter a number)", length, m)
	}
	return generator.Problem{
		Question: question,
		Answer:   strconv.Itoa(answer),
		Explanation: fmt.Sprintf("Count the complement instead: %d^%d = %d strings in all, and %d^%d = %d contain no \u2018a\u2019. The difference is %d.",
			m, length, total, m-1, length, none, answer),
		Facts: map[string]string{
			"letters": strconv.Itoa(m), "length": strconv.Itoa(length),
			"total": strconv.Itoa(total), "none": strconv.Itoa(none), "value": strconv.Itoa(answer),
		},
		KPOrdinal: o,
	}
}

func w2gDegreeSequence(o int) generator.Problem {
	a, b := w2gN(2, 10), w2gN(2, 10)
	degrees := make([]int, 0, a+b)
	for i := 0; i < a; i++ {
		degrees = append(degrees, b)
	}
	for i := 0; i < b; i++ {
		degrees = append(degrees, a)
	}
	sortInts(degrees)
	for i, j := 0, len(degrees)-1; i < j; i, j = i+1, j-1 {
		degrees[i], degrees[j] = degrees[j], degrees[i]
	}
	answer := intsString(degrees, ",")
	sum := 2 * a * b
	var question string
	switch o {
	case 1:
		question = fmt.Sprintf("Give the degree sequence of the complete bipartite graph K(%d, %d) in non-increasing order, separated by commas (for example 3,3,2).", a, b)
	case 2:
		question = fmt.Sprintf("List the vertex degrees of the complete bipartite graph K(%d, %d) in non-increasing order, separated by commas.", a, b)
	default:
		question = fmt.Sprintf("Write the degree sequence of K(%d, %d) from largest to smallest, separated by commas.", a, b)
	}
	return generator.Problem{
		Question: question,
		Answer:   answer,
		Explanation: fmt.Sprintf("Each of the %d vertices in the first part has degree %d and each of the %d in the second has degree %d, so the sequence is %s. Its sum is %d = 2 \u00b7 %d, as the handshake lemma requires.",
			a, b, b, a, answer, sum, a*b),
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "sequence": answer,
			"sum": strconv.Itoa(sum), "edges": strconv.Itoa(a * b),
		},
		KPOrdinal: o,
	}
}
