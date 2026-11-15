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
	return w2gMasterTheorem(o)
}

func booleanFunctionsProblem(o int) generator.Problem {
	return w2gBooleanFunctions(o)
}

func burnsideProblem(o int) generator.Problem {
	return w2gBurnside(o)
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
	return w2gResolution(o)
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
	return w2gProofCases(o)
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
