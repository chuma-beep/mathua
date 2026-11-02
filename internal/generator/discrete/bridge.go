package discrete

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// discBridgeGen serves the concepts added by the 2026 content expansion.
// Every branch is numeric and honours the knowledge-point ordinal so that
// each lesson section asks about the material taught in that section.
type discBridgeGen struct{ kind string }

func discOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func discProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func (g *discBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := discOrdinal(ctx)
	switch g.kind {
	case "master_theorem":
		return masterTheoremProblem(o)
	case "boolean_functions":
		return booleanFunctionsProblem(o)
	case "burnside":
		return burnsideProblem(o)
	case "egf":
		return egfProblem(o)
	case "resolution":
		return resolutionProblem(o)
	case "relation_closure":
		return relationClosureProblem(o)
	case "fsm":
		return fsmProblem(o)
	case "proof_cases":
		return proofCasesProblem(o)
	case "hamming":
		return hammingProblem(o)
	}
	return discProblem(o, "What is 1 + 1? (enter a number)", "2", "One plus one is two.")
}

func masterTheoremProblem(o int) generator.Problem {
	switch o {
	case 1:
		type pair struct {
			a, b int
			ans  string
		}
		pairs := []pair{{4, 2, "2"}, {8, 2, "3"}, {16, 2, "4"}, {2, 4, "0.5"}, {3, 9, "0.5"}, {9, 3, "2"}, {27, 3, "3"}, {8, 4, "1.5"}}
		p := pairs[rand.Intn(len(pairs))]
		return discProblem(o,
			fmt.Sprintf("A recurrence has the form T(n) = %d T(n/%d) + f(n). What is log base %d of %d? (enter a decimal)", p.a, p.b, p.b, p.a),
			p.ans,
			fmt.Sprintf("The exponent is log base %d of %d = %s.", p.b, p.a, p.ans))
	case 2:
		type triple struct{ a, b, k, class int }
		triples := []triple{{4, 2, 1, 1}, {8, 2, 2, 1}, {16, 2, 2, 1}, {2, 2, 1, 2}, {3, 3, 2, 2}, {2, 2, 3, 3}, {3, 3, 4, 3}, {4, 2, 4, 3}}
		t := triples[rand.Intn(len(triples))]
		return discProblem(o,
			fmt.Sprintf("For T(n) = %d T(n/%d) + n^%d, which master-theorem case applies? Enter 1 if the recursive term dominates, 2 for a tie, or 3 if the extra work dominates.", t.a, t.b, t.k),
			strconv.Itoa(t.class),
			fmt.Sprintf("Compare log base %d of %d with the extra-work exponent %d; the comparison gives case %d.", t.b, t.a, t.k, t.class))
	default:
		type triple struct{ a, b, k, exp int }
		triples := []triple{{4, 2, 3, 3}, {8, 2, 1, 3}, {16, 2, 2, 4}, {2, 2, 3, 3}, {9, 3, 1, 2}, {27, 3, 4, 4}, {4, 2, 0, 2}, {8, 2, 0, 3}}
		t := triples[rand.Intn(len(triples))]
		return discProblem(o,
			fmt.Sprintf("For T(n) = %d T(n/%d) + n^%d, the solution is Theta(n^p) times at most a logarithmic factor. What is p? (enter a number)", t.a, t.b, t.k),
			strconv.Itoa(t.exp),
			fmt.Sprintf("The exponent p is the larger of log base %d of %d and %d, which is %d.", t.b, t.a, t.k, t.exp))
	}
}

func booleanFunctionsProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := rand.Intn(5)
		ans := 1 << (1 << n)
		return discProblem(o,
			fmt.Sprintf("How many Boolean functions of %d variable(s) are there? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("There are 2^%d truth-table rows, and each row can output 0 or 1, so the count is 2^(2^%d) = %d.", n, n, ans))
	case 2:
		type expr struct {
			text string
			ones int
		}
		exprs := []expr{{"x AND y", 1}, {"x OR y", 3}, {"x XOR y", 2}, {"NOT x", 2}, {"x AND NOT y", 1}, {"NOT (x AND y)", 3}, {"x OR NOT y", 3}, {"x NAND y", 3}}
		e := exprs[rand.Intn(len(exprs))]
		return discProblem(o,
			fmt.Sprintf("For two inputs, how many of the four truth-table rows make %s true? (enter a number)", e.text),
			strconv.Itoa(e.ones),
			fmt.Sprintf("Exactly %d of the four assignments make %s true.", e.ones, e.text))
	default:
		type expr struct {
			text string
			ones int
		}
		exprs := []expr{{"x AND y AND z", 2}, {"x OR y", 12}, {"x XOR y", 8}, {"x AND y", 4}, {"NOT x AND y", 4}, {"x OR y OR z", 14}, {"x XOR y XOR z", 8}}
		e := exprs[rand.Intn(len(exprs))]
		return discProblem(o,
			fmt.Sprintf("For four Boolean variables, how many of the sixteen assignments make %s true? (enter a number)", e.text),
			strconv.Itoa(e.ones),
			fmt.Sprintf("Counting the satisfying rows of the sixteen gives %d for %s.", e.ones, e.text))
	}
}

func burnsideProblem(o int) generator.Problem {
	switch o {
	case 1:
		type pair struct{ n, k, ans int }
		pairs := []pair{{3, 2, 4}, {4, 2, 6}, {5, 2, 8}, {6, 2, 14}, {3, 3, 11}, {4, 3, 24}}
		p := pairs[rand.Intn(len(pairs))]
		return discProblem(o,
			fmt.Sprintf("How many distinct necklaces can be made from %d beads and %d colours, counting two necklaces as the same when one is a rotation of the other? (enter a number)", p.n, p.k),
			strconv.Itoa(p.ans),
			fmt.Sprintf("Burnside averages the fixed colourings over the %d rotations, which gives %d.", p.n, p.ans))
	case 2:
		type args struct{ n, k, d, ans int }
		all := []args{{4, 2, 1, 2}, {4, 2, 2, 4}, {4, 3, 2, 9}, {6, 2, 2, 4}, {6, 2, 3, 8}, {6, 3, 3, 27}}
		a := all[rand.Intn(len(all))]
		return discProblem(o,
			fmt.Sprintf("%d positions are coloured with %d colours. How many colourings are fixed by a rotation of the positions by %d step(s)? (enter a number)", a.n, a.k, a.d),
			strconv.Itoa(a.ans),
			fmt.Sprintf("The rotation splits the positions into gcd(%d, %d) cycles, so %d to that power gives %d fixed colourings.", a.n, a.d, a.k, a.ans))
	default:
		type args struct {
			shape  string
			k, ans int
		}
		all := []args{{"a square", 2, 6}, {"a square", 3, 21}, {"a square", 4, 100}, {"a triangle", 2, 4}, {"a triangle", 3, 11}, {"a triangle", 4, 24}}
		a := all[rand.Intn(len(all))]
		return discProblem(o,
			fmt.Sprintf("How many distinct colourings of %s with %d colours are there, up to rotation? (enter a number)", a.shape, a.k),
			strconv.Itoa(a.ans),
			fmt.Sprintf("Burnside's lemma over the rotations of %s gives %d.", a.shape, a.ans))
	}
}

func egfProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 2 + rand.Intn(2)
		n := 1 + rand.Intn(5)
		ans := 1
		for i := 0; i < n; i++ {
			ans *= k
		}
		return discProblem(o,
			fmt.Sprintf("The exponential generating function is e^(%d x). What is the coefficient a_n for n = %d in sum a_n x^n/n!? (enter a number)", k, n),
			strconv.Itoa(ans),
			fmt.Sprintf("Expanding e^(%d x) gives coefficient %d^%d = %d.", k, k, n, ans))
	case 2:
		n := 1 + rand.Intn(5)
		ans := 1
		for i := 2; i <= n; i++ {
			ans *= i
		}
		return discProblem(o,
			fmt.Sprintf("The exponential generating function is 1/(1-x). What is the coefficient a_n for n = %d in sum a_n x^n/n!? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The ordinary coefficient is 1, so a_n = n! = %d.", ans))
	default:
		n := 1 + rand.Intn(6)
		return discProblem(o,
			fmt.Sprintf("The exponential generating function is x e^x. What is the coefficient a_n for n = %d in sum a_n x^n/n!? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("x e^x has coefficient a_n = n, so for n = %d the answer is %d.", n, n))
	}
}

func resolutionProblem(o int) generator.Problem {
	switch o {
	case 1:
		p := 2 + rand.Intn(3)
		q := 2 + rand.Intn(3)
		ans := p + q - 2
		return discProblem(o,
			fmt.Sprintf("A clause has %d literals and another has %d literals. They share exactly one complementary pair. How many literals are in their resolvent? (enter a number)", p, q),
			strconv.Itoa(ans),
			fmt.Sprintf("Delete the two complementary literals: %d + %d - 2 = %d.", p, q, ans))
	case 2:
		type pair struct {
			text string
			ans  int
		}
		all := []pair{{"(x OR y) and (NOT x OR z)", 1}, {"(x OR y) and (NOT z OR w)", 0}, {"(x OR y OR z) and (NOT x OR NOT y)", 2}, {"(x) and (NOT x)", 1}, {"(x OR y) and (NOT y OR z)", 1}}
		p := all[rand.Intn(len(all))]
		return discProblem(o,
			fmt.Sprintf("How many distinct resolvents can be formed from %s? (enter a number)", p.text),
			strconv.Itoa(p.ans),
			fmt.Sprintf("Count the complementary variables between the two clauses; each pair gives one resolvent, so the answer is %d.", p.ans))
	default:
		type formula struct {
			text  string
			units int
		}
		all := []formula{{"(x) AND (y) AND (NOT x OR z)", 2}, {"(x OR y) AND (NOT x) AND (NOT y)", 0}, {"(x) AND (NOT x OR y) AND (NOT y OR z)", 1}, {"(x) AND (y) AND (z) AND (NOT x OR NOT y)", 3}, {"(x OR y) AND (NOT x OR y) AND (z)", 1}}
		f := all[rand.Intn(len(all))]
		return discProblem(o,
			fmt.Sprintf("How many unit clauses appear in %s? (enter a number)", f.text),
			strconv.Itoa(f.units),
			fmt.Sprintf("A unit clause has one literal; counting them in the formula gives %d.", f.units))
	}
}

func relationClosureProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 3 + rand.Intn(4)
		have := rand.Intn(n + 1)
		ans := n - have
		return discProblem(o,
			fmt.Sprintf("A relation on a set of %d elements already contains %d of the pairs (a, a). How many pairs must be added to form the reflexive closure? (enter a number)", n, have),
			strconv.Itoa(ans),
			fmt.Sprintf("The reflexive closure needs all %d diagonal pairs, so it adds %d - %d = %d.", n, n, have, ans))
	case 2:
		m := 1 + rand.Intn(5)
		return discProblem(o,
			fmt.Sprintf("A relation consists of %d ordered pairs, and no pair has its reverse present. How many pairs must be added to form the symmetric closure? (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("Each of the %d pairs contributes one missing reverse, so the symmetric closure adds %d pairs.", m, m))
	default:
		n := 2 + rand.Intn(5)
		ans := n * (n - 1) / 2
		return discProblem(o,
			fmt.Sprintf("On vertices 1 to %d, the relation contains only the pairs (i, i+1) for i from 1 to %d. How many ordered pairs are in its transitive closure? (enter a number)", n, n-1),
			strconv.Itoa(ans),
			fmt.Sprintf("The chain reaches every later vertex, so the closure has %d choose 2 = %d pairs.", n, ans))
	}
}

func fsmProblem(o int) generator.Problem {
	switch o {
	case 1:
		m := 2 + rand.Intn(5)
		s := rand.Intn(m)
		k := 1 + rand.Intn(12)
		ans := (s + k) % m
		return discProblem(o,
			fmt.Sprintf("A machine has states 0 to %d and advances by one state modulo %d on every input. Starting in state %d, what state is reached after %d step(s)? (enter a number)", m-1, m, s, k),
			strconv.Itoa(ans),
			fmt.Sprintf("Add %d to %d and reduce modulo %d: the result is %d.", k, s, m, ans))
	case 2:
		m := 2 + rand.Intn(5)
		n := 2 + rand.Intn(5)
		return discProblem(o,
			fmt.Sprintf("One machine has %d states and another has %d states. How many states does their product machine have? (enter a number)", m, n),
			strconv.Itoa(m*n),
			fmt.Sprintf("A product machine pairs every state of the first with every state of the second: %d x %d = %d.", m, n, m*n))
	default:
		a := 2 + rand.Intn(2)
		n := 2 + rand.Intn(4)
		ans := 1
		for i := 0; i < n; i++ {
			ans *= a
		}
		return discProblem(o,
			fmt.Sprintf("An alphabet has %d symbols. How many input strings of length %d are there? (enter a number)", a, n),
			strconv.Itoa(ans),
			fmt.Sprintf("Each of the %d positions has %d choices, so there are %d^%d = %d strings.", n, a, a, n, ans))
	}
}

func proofCasesProblem(o int) generator.Problem {
	switch o {
	case 1:
		m := 2 + rand.Intn(5)
		return discProblem(o,
			fmt.Sprintf("A proof splits integers into cases according to their remainder modulo %d. How many cases are needed? (enter a number)", m),
			strconv.Itoa(m),
			fmt.Sprintf("The possible remainders modulo %d are 0 through %d, so there are %d cases.", m, m-1, m))
	case 2:
		k := 1 + rand.Intn(5)
		ans := 1
		for i := 0; i < k; i++ {
			ans *= 2
		}
		return discProblem(o,
			fmt.Sprintf("A proof analyses the sign of a product of %d factors. Each factor is positive or negative. How many sign cases are there? (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("Each factor has two signs, so %d factors give 2^%d = %d cases.", k, k, ans))
	default:
		m := 1 + rand.Intn(5)
		return discProblem(o,
			fmt.Sprintf("A proof splits on the remainder modulo %d and also on whether the integer is positive, negative or zero. How many cases are there? (enter a number)", m),
			strconv.Itoa(3*m),
			fmt.Sprintf("There are %d residues and three signs, so the total is 3 x %d = %d.", m, m, 3*m))
	}
}

func hammingProblem(o int) generator.Problem {
	type pair struct{ m, r int }
	pairs := []pair{{1, 2}, {2, 3}, {3, 3}, {4, 3}, {5, 4}, {6, 4}, {7, 4}, {8, 4}, {9, 4}, {10, 4}, {11, 4}}
	p := pairs[rand.Intn(len(pairs))]
	switch o {
	case 1:
		return discProblem(o,
			fmt.Sprintf("A Hamming code carries %d data bits. What is the smallest number of parity bits r with 2^r at least %d + r + 1? (enter a number)", p.m, p.m),
			strconv.Itoa(p.r),
			fmt.Sprintf("Testing r from one upward, the first value with 2^r at least %d is %d.", p.m+p.r+1, p.r))
	case 2:
		return discProblem(o,
			fmt.Sprintf("A Hamming code has %d data bits and %d parity bits. What is the total codeword length? (enter a number)", p.m, p.r),
			strconv.Itoa(p.m+p.r),
			fmt.Sprintf("The code length is the data bits plus the parity bits: %d + %d = %d.", p.m, p.r, p.m+p.r))
	default:
		d := 3 + 2*rand.Intn(6)
		ans := (d - 1) / 2
		return discProblem(o,
			fmt.Sprintf("A code has minimum Hamming distance %d. How many bit errors can it correct? (enter a number)", d),
			strconv.Itoa(ans),
			fmt.Sprintf("A code with distance %d corrects floor((%d - 1)/2) = %d errors.", d, d, ans))
	}
}
