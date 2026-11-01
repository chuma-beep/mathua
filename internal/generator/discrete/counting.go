package discrete

// Counting problems for the discrete mathematics concepts that used to ask a
// recognition question. Each pool holds sixteen statements, eight true and
// eight false. Drawing eight and asking for the number that are true gives
// nine possible answers without ever showing the answer in the stem.

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type discreteVerdict struct {
	label string
	yes   bool
	note  string
}

const discreteDraw = 8

func countDiscreteVerdicts(lead string, pool []discreteVerdict) generator.Problem {
	var yes, no []discreteVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	count := rand.Intn(discreteDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if discreteDraw-count > len(no) {
		count = discreteDraw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append([]discreteVerdict{}, yes[:count]...)
	chosen = append(chosen, no[:discreteDraw-count]...)
	rand.Shuffle(len(chosen), func(i, j int) { chosen[i], chosen[j] = chosen[j], chosen[i] })

	parts := make([]string, 0, len(chosen))
	notes := make([]string, 0, len(chosen))
	for i, v := range chosen {
		parts = append(parts, fmt.Sprintf("(%c) %s", 'a'+i, v.label))
		verdict := "no"
		if v.yes {
			verdict = "yes"
		}
		notes = append(notes, fmt.Sprintf("%s: %s, %s", v.label, verdict, v.note))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("%s %s (enter a number)", lead, strings.Join(parts, " ")),
		Answer:      fmt.Sprintf("%d", count),
		Explanation: fmt.Sprintf("Exactly %d of the eight statements are true. %s", count, strings.Join(notes, " ")),
		Facts:       map[string]string{"count": fmt.Sprintf("%d", count)},
	}
}

func setOperationsProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these set identities hold for all sets A, B and C?",
		[]discreteVerdict{
			{"A union B = B union A", true, "union is commutative"},
			{"A intersect B = B intersect A", true, "intersection is commutative"},
			{"A union (B intersect C) = (A union B) intersect (A union C)", true, "union distributes over intersection"},
			{"A intersect (B union C) = (A intersect B) union (A intersect C)", true, "intersection distributes over union"},
			{"A minus B = A intersect B complement", true, "difference is intersection with the complement"},
			{"(A complement) complement = A", true, "double complement"},
			{"A is a subset of A union B", true, "a union contains both parts"},
			{"(A union B) complement = A complement intersect B complement", true, "De Morgan"},
			{"A union B = A intersect B", false, "only when A = B"},
			{"A minus B = B minus A", false, "difference is not symmetric"},
			{"A union the empty set = the empty set", false, "the union is A"},
			{"A intersect A complement = A", false, "the intersection is empty"},
			{"(A intersect B) complement = A complement intersect B complement", false, "the dual uses union"},
			{"A is a subset of A intersect B", false, "not unless A is a subset of B"},
			{"A union A complement = the empty set", false, "the union is the universe"},
			{"A minus A = A", false, "the difference is empty"},
		})
}

func inductionProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about mathematical induction are true?",
		[]discreteVerdict{
			{"An induction proof starts with a base case", true, "the base case anchors the chain"},
			{"The inductive step assumes P(k) and proves P(k+1)", true, "ordinary induction"},
			{"Induction can prove a statement for every natural number", true, "the conclusion covers all n"},
			{"Strong induction may assume P(j) for every j at most k", true, "complete induction"},
			{"Strong induction and ordinary induction have the same strength", true, "each proves the other"},
			{"The well-ordering principle is equivalent to induction", true, "a standard equivalence"},
			{"A false base case can make the inductive step useless", true, "the chain never starts"},
			{"Induction can prove that 1 + ... + n = n(n+1)/2", true, "a standard example"},
			{"The inductive step assumes P(k+1) and proves P(k)", false, "the direction is reversed"},
			{"A base case is optional in an induction proof", false, "without it the chain has no start"},
			{"Induction can only prove statements about finite sets", false, "it proves statements about all naturals"},
			{"Strong induction needs a separate base case for every k", false, "one base case or a few suffice"},
			{"Induction cannot prove inequalities", false, "many inequalities are proved by induction"},
			{"If the inductive step holds then the statement is always true", false, "the base case is also required"},
			{"Induction requires the statement to hold for negative integers", false, "the usual form starts at a fixed integer"},
			{"The inductive step alone proves the base case", false, "the base case is separate"},
		})
}

func equivalenceProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these propositional equivalences are valid?",
		[]discreteVerdict{
			{"P implies Q is equivalent to not P or Q", true, "material implication"},
			{"not (P and Q) is equivalent to not P or not Q", true, "De Morgan"},
			{"not (P or Q) is equivalent to not P and not Q", true, "De Morgan"},
			{"P if and only if Q is equivalent to (P implies Q) and (Q implies P)", true, "biconditional definition"},
			{"P implies Q is equivalent to not Q implies not P", true, "contrapositive"},
			{"not not P is equivalent to P", true, "double negation"},
			{"P or P is equivalent to P", true, "idempotence"},
			{"P and P is equivalent to P", true, "idempotence"},
			{"P implies Q is equivalent to Q implies P", false, "that is the converse"},
			{"P and Q is equivalent to P or Q", false, "they differ when exactly one holds"},
			{"not (P implies Q) is equivalent to P and Q", false, "it is P and not Q"},
			{"P if and only if Q is equivalent to P and Q", false, "both false also satisfies the biconditional"},
			{"P or not P is equivalent to P", false, "it is a tautology"},
			{"not (P implies Q) is equivalent to not P or Q", false, "that is the implication itself"},
			{"P implies Q is equivalent to P and Q", false, "the implication is weaker"},
			{"P if and only if Q is equivalent to P or Q", false, "the biconditional is stronger"},
		})
}

func relationsProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about relations are true?",
		[]discreteVerdict{
			{"Equality on a set is an equivalence relation", true, "reflexive, symmetric and transitive"},
			{"The order at most on the real numbers is a partial order", true, "reflexive, antisymmetric and transitive"},
			{"Congruence modulo n is an equivalence relation on the integers", true, "the standard example"},
			{"Divisibility on the positive integers is a partial order", true, "reflexive, antisymmetric and transitive"},
			{"Every function from A to B is a relation from A to B", true, "a function is a special relation"},
			{"Composition of relations is associative", true, "composition is associative"},
			{"The inverse of an equivalence relation is an equivalence relation", true, "all three properties are symmetric"},
			{"A relation on A is a subset of A times A", true, "the definition"},
			{"Every relation is symmetric", false, "order relations are not"},
			{"Every reflexive relation is an equivalence relation", false, "it may fail symmetry or transitivity"},
			{"The order at most on the real numbers is an equivalence relation", false, "it is not symmetric"},
			{"Every symmetric relation is transitive", false, "a relation may be symmetric only"},
			{"The empty relation on a nonempty set is reflexive", false, "no element is related to itself"},
			{"Every relation is a function", false, "one input may have several outputs"},
			{"The inverse of a partial order need not be a partial order", false, "the inverse is always a partial order"},
			{"A relation on A is a subset of A union A", false, "the relation lives in the product"},
		})
}

func eulerianProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about Eulerian graphs are true?",
		[]discreteVerdict{
			{"K_3 has an Eulerian circuit", true, "every vertex has degree 2"},
			{"A connected graph with every vertex of even degree has an Eulerian circuit", true, "Euler's theorem"},
			{"A connected graph with exactly two vertices of odd degree has an Eulerian trail", true, "the trail runs between the odd vertices"},
			{"K_5 has an Eulerian circuit", true, "every vertex has degree 4"},
			{"Every cycle graph C_n has an Eulerian circuit", true, "the cycle itself"},
			{"An Eulerian circuit uses every edge exactly once", true, "the definition"},
			{"A connected graph with four vertices of odd degree has no Eulerian trail", true, "a trail needs zero or two odd vertices"},
			{"The complete bipartite graph K_2,2 has an Eulerian circuit", true, "it is a 4-cycle"},
			{"K_4 has an Eulerian circuit", false, "every vertex has odd degree 3"},
			{"The path A-B-C has an Eulerian circuit", false, "it has two odd endpoints"},
			{"A connected graph with exactly two odd vertices has an Eulerian circuit", false, "it has a trail but not a circuit"},
			{"A graph can have exactly one vertex of odd degree", false, "the handshake lemma forbids it"},
			{"Every Eulerian graph is a complete graph", false, "cycles are Eulerian and not complete"},
			{"An Eulerian trail must return to its starting vertex", false, "only a circuit does"},
			{"A graph with a bridge always has an Eulerian circuit", false, "a bridge blocks a circuit"},
			{"K_3,3 has an Eulerian circuit", false, "all six vertices have odd degree 3"},
		})
}

func strongInductionProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about strong induction are true?",
		[]discreteVerdict{
			{"Strong induction assumes P(j) for every j at most k", true, "the strong hypothesis"},
			{"Strong induction implies ordinary induction", true, "the strong hypothesis contains P(k)"},
			{"The well-ordering principle is equivalent to strong induction", true, "a standard equivalence"},
			{"Strong induction can prove unique prime factorization", true, "split n into smaller factors"},
			{"Strong induction may need more than one base case", true, "some recurrences need several starting values"},
			{"Strong induction is also called complete induction", true, "a common name"},
			{"A strong induction proof may split k+1 into smaller cases", true, "the usual mechanism"},
			{"Ordinary induction implies strong induction", true, "the two principles are equivalent"},
			{"Strong induction assumes only P(k)", false, "that is ordinary induction"},
			{"Strong induction needs no base case", false, "the base case is still required"},
			{"Strong induction is weaker than ordinary induction", false, "they are equivalent"},
			{"Strong induction applies only to finite sets", false, "it applies to all naturals from a starting point"},
			{"Strong induction cannot prove statements about integers", false, "it is used throughout number theory"},
			{"Strong induction requires the statement to hold for all integers", false, "only from the starting value upward"},
			{"Strong induction assumes P(k+1) to prove P(k)", false, "the direction is reversed"},
			{"The base case of strong induction must be k = 0", false, "any starting integer works"},
		})
}

func coloringProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about graph colouring are true?",
		[]discreteVerdict{
			{"A graph is bipartite exactly when it is 2-colourable", true, "the two parts are the colour classes"},
			{"The complete graph K_n needs n colours", true, "every pair is adjacent"},
			{"An odd cycle needs 3 colours", true, "it is not bipartite"},
			{"Every planar graph is 4-colourable", true, "the four colour theorem"},
			{"Every tree is 2-colourable", true, "trees are bipartite"},
			{"The complete bipartite graph K_m,n has chromatic number 2", true, "for positive m and n"},
			{"The clique number is a lower bound for the chromatic number", true, "a clique needs distinct colours"},
			{"The chromatic number is at least the clique number", true, "the same fact"},
			{"Every graph is 2-colourable", false, "odd cycles are not"},
			{"The complete graph K_n has chromatic number 2 for n at least 3", false, "it needs n colours"},
			{"Every even cycle needs 3 colours", false, "an even cycle is 2-colourable"},
			{"Three colours always suffice", false, "complete graphs need more"},
			{"Every planar graph needs 5 colours", false, "four always suffice"},
			{"The chromatic number equals the number of vertices", false, "only for complete graphs"},
			{"Every bipartite graph needs 3 colours", false, "bipartite graphs need 2"},
			{"The chromatic number is always less than the clique number", false, "it is at least the clique number"},
		})
}

func contradictionProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about proof by contradiction are true?",
		[]discreteVerdict{
			{"A proof by contradiction assumes not P", true, "the negated goal"},
			{"Deriving Q and not Q from not P proves P", true, "the contradiction rule"},
			{"The proof that the square root of 2 is irrational is usually by contradiction", true, "assume a lowest-terms fraction"},
			{"Euclid's proof that there are infinitely many primes is by contradiction", true, "assume a finite list"},
			{"Contradiction proves P by showing not P implies Q and not Q", true, "the logical form"},
			{"A contradiction proof may use any false consequence of not P", true, "any contradiction suffices"},
			{"If not P is impossible then P is true", true, "classical logic"},
			{"Proof by contradiction is a valid classical proof method", true, "it is standard"},
			{"A proof by contradiction assumes P", false, "it assumes the negation"},
			{"Proof by contrapositive is the same as proof by contradiction", false, "they are different methods"},
			{"Proof by contradiction requires a counterexample", false, "a contradiction is enough"},
			{"A contradiction proof proves not P", false, "it proves P"},
			{"Deriving P from P is a proof by contradiction", false, "there is no contradiction"},
			{"Contradiction is invalid in classical logic", false, "it is a core classical rule"},
			{"Every proof by contradiction is also a direct proof", false, "the methods differ"},
			{"A contradiction proof must assume Q and derive P", false, "it assumes not P"},
		})
}

func inferenceProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these inference rules are valid?",
		[]discreteVerdict{
			{"From P and P implies Q infer Q", true, "modus ponens"},
			{"From P implies Q and not Q infer not P", true, "modus tollens"},
			{"From P implies Q and Q implies R infer P implies R", true, "hypothetical syllogism"},
			{"From P or Q and not P infer Q", true, "disjunctive syllogism"},
			{"From P infer P or Q", true, "addition"},
			{"From P and Q infer P", true, "simplification"},
			{"From P and Q infer P and Q", true, "conjunction"},
			{"From P if and only if Q and P infer Q", true, "biconditional elimination"},
			{"From P implies Q and Q infer P", false, "affirming the consequent"},
			{"From P implies Q and not P infer not Q", false, "denying the antecedent"},
			{"From P or Q and P infer not Q", false, "both may be true"},
			{"From P implies Q infer Q implies P", false, "that is the converse"},
			{"From not (P and Q) infer not P and not Q", false, "at least one is false, not both"},
			{"From P implies Q and Q implies P infer P and Q", false, "both may be false"},
			{"From P or Q infer P", false, "Q may be the true disjunct"},
			{"From P implies Q and P infer not Q", false, "modus ponens gives Q"},
		})
}

func predicatesProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these predicate logic statements are true?",
		[]discreteVerdict{
			{"For all x P(x) implies P(a) for any a in the domain", true, "universal instantiation"},
			{"not for all x P(x) is equivalent to there exists x not P(x)", true, "quantifier negation"},
			{"not there exists x P(x) is equivalent to for all x not P(x)", true, "quantifier negation"},
			{"There exists x with P(x) and Q(x) implies there exists x P(x) and there exists x Q(x)", true, "the same witness serves both"},
			{"For all x P(x) implies there exists x P(x) when the domain is nonempty", true, "a universal claim has a witness"},
			{"There exists x for all y R(x,y) implies for all y there exists x R(x,y)", true, "the same x works for every y"},
			{"For all x (P(x) and Q(x)) implies for all x P(x)", true, "conjunction elimination under a universal"},
			{"For all x P(x) and for all x (P(x) implies Q(x)) imply for all x Q(x)", true, "universal modus ponens"},
			{"For all x there exists y R(x,y) implies there exists y for all x R(x,y)", false, "the witness may depend on x"},
			{"There exists x P(x) and there exists x Q(x) implies there exists x (P(x) and Q(x))", false, "the witnesses may differ"},
			{"For all x P(x) implies P(x) only for some x", false, "it implies it for every x"},
			{"There exists x P(x) implies for all x P(x)", false, "one witness does not cover all"},
			{"For all x (P(x) or Q(x)) implies for all x P(x) or for all x Q(x)", false, "the choice may vary with x"},
			{"not there exists x P(x) is equivalent to there exists x not P(x)", false, "the dual uses for all"},
			{"For all x P(x) is false when P(a) is true for one a", false, "one case does not refute a universal"},
			{"There exists x P(x) is true only when P(x) holds for every x", false, "one witness is enough"},
		})
}

func satisfiabilityProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about satisfiability are true?",
		[]discreteVerdict{
			{"Every tautology is satisfiable", true, "every assignment satisfies it"},
			{"A contradiction is unsatisfiable", true, "no assignment satisfies it"},
			{"P and not P is unsatisfiable", true, "the two conjuncts conflict"},
			{"P or not P is a tautology", true, "excluded middle"},
			{"A satisfiable formula need not be valid", true, "it may fail on some assignment"},
			{"2-SAT is decidable in polynomial time", true, "implication graph algorithm"},
			{"A formula with no clauses is satisfiable", true, "the empty conjunction is true"},
			{"If a formula is unsatisfiable then its negation is a tautology", true, "negation flips satisfiability"},
			{"Every satisfiable formula is a tautology", false, "satisfiable is weaker"},
			{"P and not P is satisfiable", false, "no assignment works"},
			{"A tautology is unsatisfiable", false, "every assignment works"},
			{"Every formula is either satisfiable or a tautology", false, "a contingent formula is neither"},
			{"General SAT is decidable in polynomial time", false, "it is the canonical NP-complete problem"},
			{"A formula is valid exactly when it is satisfiable", false, "validity is stronger"},
			{"An unsatisfiable formula has a satisfying assignment", false, "by definition it has none"},
			{"P or Q is unsatisfiable", false, "three assignments satisfy it"},
		})
}

func contrapositiveProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about proof by contrapositive are true?",
		[]discreteVerdict{
			{"To prove P implies Q by contrapositive, prove not Q implies not P", true, "the contrapositive"},
			{"P implies Q and not Q implies not P are equivalent", true, "contrapositive equivalence"},
			{"The contrapositive of P implies Q is not Q implies not P", true, "the definition"},
			{"The contrapositive of not P implies Q is not Q implies P", true, "double negation"},
			{"A contrapositive proof assumes not Q", true, "the hypothesis of the contrapositive"},
			{"Contrapositive and direct proof are both valid", true, "both prove the same implication"},
			{"The contrapositive of P implies not Q is Q implies not P", true, "negate and swap"},
			{"The contrapositive of not P implies not Q is Q implies P", true, "negate and swap"},
			{"The contrapositive of P implies Q is Q implies P", false, "that is the converse"},
			{"The contrapositive of P implies Q is not P implies not Q", false, "that is the inverse"},
			{"A contrapositive proof assumes P", false, "it assumes not Q"},
			{"The converse of P implies Q is not Q implies not P", false, "that is the contrapositive"},
			{"P implies Q is equivalent to its converse", false, "they differ when P is false and Q true"},
			{"A contrapositive proof derives a contradiction", false, "that is proof by contradiction"},
			{"The contrapositive of P implies Q is P and not Q", false, "that is the negation"},
			{"Contrapositive is invalid in classical logic", false, "it is valid"},
		})
}

func structuralProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these statements about structural induction are true?",
		[]discreteVerdict{
			{"Structural induction proves properties of recursively defined objects", true, "the general method"},
			{"A structural induction proof needs a base case for atomic objects", true, "the base case anchors the recursion"},
			{"Structural induction can prove properties of binary trees", true, "a standard application"},
			{"Structural induction can prove properties of strings", true, "strings are recursively built"},
			{"The inductive step handles each recursive constructor", true, "one case per constructor"},
			{"Structural induction is a form of well-founded induction", true, "it follows the construction order"},
			{"The base case covers leaves or empty structures", true, "the atomic objects"},
			{"Structural induction can prove properties of arithmetic expressions", true, "expressions are recursive"},
			{"Structural induction applies only to natural numbers", false, "it applies to any recursive type"},
			{"Structural induction needs no base case", false, "atomic objects still need a case"},
			{"Structural induction assumes the property for all objects at once", false, "it assumes it for subobjects"},
			{"Structural induction cannot prove properties of trees", false, "trees are the standard example"},
			{"The inductive step only handles one constructor", false, "every constructor needs a case"},
			{"Structural induction is invalid for mutually recursive definitions", false, "it extends to mutual recursion"},
			{"Structural induction requires a numeric measure on every object", false, "the construction order suffices"},
			{"Structural induction proves only finite objects of size one", false, "it covers all finite constructions"},
		})
}

func countableProblem() generator.Problem {
	return countDiscreteVerdicts(
		"How many of these sets are countable?",
		[]discreteVerdict{
			{"The natural numbers", true, "the defining example"},
			{"The integers", true, "interleave positive and negative"},
			{"The rational numbers", true, "a diagonal enumeration"},
			{"Every finite set", true, "list its elements"},
			{"A countable union of countable sets", true, "a standard theorem"},
			{"The product N times N", true, "enumerate by diagonals"},
			{"The algebraic numbers", true, "countably many polynomials"},
			{"The set of finite subsets of N", true, "a countable union of finite products"},
			{"The real numbers", false, "Cantor's diagonal argument"},
			{"The power set of N", false, "uncountable by Cantor's theorem"},
			{"The irrational numbers", false, "uncountable as R minus a countable set"},
			{"The open interval from 0 to 1", false, "in bijection with R"},
			{"The Cantor set", false, "it is uncountable"},
			{"Every infinite set", false, "R is a counterexample"},
			{"The product Q times R", false, "it has the cardinality of R"},
			{"The set of all subsets of N", false, "same as the power set"},
		})
}

func generatingFunctionsProblem() generator.Problem {
	return countDiscreteVerdicts("How many of these statements about generating functions are true?", []discreteVerdict{
		{"1/(1-x) = 1 + x + x^2 + ...", true, "the geometric series expansion"},
		{"the coefficient of x^3 in x/(1-x-x^2) is 2", true, "the series is x + x^2 + 2x^3 + 3x^4 + ..."},
		{"the product of two generating functions gives a convolution", true, "coefficients are sums a_k b_{n-k}"},
		{"the exponential generating function is the sum of a_n x^n/n!", true, "the factorial denominator defines the EGF"},
		{"differentiating an ordinary generating function multiplies a_n by n", true, "A'(x) has coefficients n a_n"},
		{"1/(1-2x) has coefficients 2^n", true, "the ratio is 2, so the coefficients double"},
		{"the coefficient of x^n in (1+x)^m is C(m,n)", true, "the binomial theorem"},
		{"a generating function turns a recurrence into algebra", true, "solve for A(x), then expand"},
		{"1/(1-x) has coefficients n", false, "every coefficient is 1"},
		{"the product of generating functions multiplies coefficients pointwise", false, "the product convolves them"},
		{"the exponential generating function of 1, 1, 1, ... is 1/(1-x)", false, "it is e^x"},
		{"x/(1-x) is the ordinary generating function of 1, 1, 1, ...", false, "it is 0, 1, 1, 1, ...; 1/(1-x) is the constant sequence"},
		{"the coefficient of x^2 in 1/(1-x)^2 is 1", false, "the coefficients are n+1, so it is 3"},
		{"every generating function converges for every real x", false, "formal series need not converge"},
		{"the coefficient of x^n in 1/(1-x)^2 is n", false, "it is n+1"},
		{"generating functions only apply to finite sequences", false, "their main use is infinite sequences"},
	})
}

func growthRatesProblem() generator.Problem {
	return countDiscreteVerdicts("How many of these growth-rate statements are true?", []discreteVerdict{
		{"n! grows faster than 2^n", true, "factorial eventually dominates every exponential"},
		{"log n grows slower than sqrt(n)", true, "the logarithm is dominated by every positive power"},
		{"n log n grows slower than n^2", true, "divide both sides by n"},
		{"2^n grows faster than n^3", true, "the exponential eventually overtakes the polynomial"},
		{"n^2 + 3n is Theta(n^2)", true, "the lower-order term is absorbed"},
		{"n log n is O(n^2)", true, "log n is O(n)"},
		{"2n^2 + 5 is O(n^2)", true, "constants and lower terms are absorbed"},
		{"log(n!) is Theta(n log n)", true, "Stirling's approximation"},
		{"n^2 grows faster than 2^n", false, "the exponential wins in the limit"},
		{"sqrt(n) grows faster than n log n", false, "n log n dominates sqrt(n)"},
		{"3^n is O(2^n)", false, "the ratio (3/2)^n is unbounded"},
		{"n! is O(2^n)", false, "the factorial dominates every exponential"},
		{"log n is Theta(n)", false, "the logarithm grows far slower"},
		{"n^3 is O(n^2)", false, "the ratio n is unbounded"},
		{"n! is Theta(2^n)", false, "the ratio is unbounded"},
		{"n is O(log n)", false, "the logarithm grows slower than n"},
	})
}
