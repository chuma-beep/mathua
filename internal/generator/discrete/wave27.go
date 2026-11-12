package discrete

// Wave 27 discrete mathematics: algorithms, coding theory, complexity,
// recurrences, finite sets and proof techniques. Every family answers an
// integer and publishes the parameters it drew.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave27(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"discrete.alg.quicksort":            "quicksort",
		"discrete.alg.merge_sort":           "mergesort",
		"discrete.alg.binary_search":        "binsearch",
		"discrete.alg.dijkstra":             "dijkstra",
		"discrete.alg.huffman":              "huffman",
		"discrete.alg.backtracking":         "backtrack",
		"discrete.alg.branch_bound":         "branchbound",
		"discrete.coding.entropy":           "entropy",
		"discrete.coding.convolutional":     "convolutional",
		"discrete.coding.ldpc":              "ldpc",
		"discrete.coding.bch":               "bch",
		"discrete.coding.singleton":         "singleton",
		"discrete.complexity.space":         "space",
		"discrete.complexity.hierarchy":     "hierarchy",
		"discrete.complexity.conp":          "conp",
		"discrete.complexity.pspace":        "pspace",
		"discrete.complexity.approximation": "approx",
		"discrete.complexity.fpt":           "fpt",
		"discrete.seq.homogeneous":          "homogeneous",
		"discrete.seq.characteristic":       "characteristic",
		"discrete.seq.nonhomogeneous":       "nonhomogeneous",
		"discrete.seq.zeckendorf":           "zeckendorf",
		"discrete.seq.continued_fraction":   "cf",
		"discrete.seq.ackermann":            "ackermann",
		"discrete.sets.cardinality":         "cardinality",
		"discrete.sets.cartesian":           "cartesian",
		"discrete.sets.equivalence_classes": "classes",
		"discrete.sets.lattice":             "lattice",
		"discrete.proof.invariant":          "invariant",
		"discrete.proof.diagonal":           "diagonal",
		"discrete.proof.descent":            "descent",
		"discrete.proof.constructive":       "constructive",
		"discrete.proof.nonconstructive":    "nonconstructive",
	} {
		reg.Register(id, &w27Gen{kind: kind})
	}
}

type w27Gen struct{ kind string }

func (g *w27Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "quicksort":
		return w27QuickSort(kp)
	case "mergesort":
		return w27MergeSort(kp)
	case "binsearch":
		return w27BinarySearch(kp)
	case "dijkstra":
		return w27Dijkstra(kp)
	case "huffman":
		return w27Huffman(kp)
	case "backtrack":
		return w27Backtracking(kp)
	case "branchbound":
		return w27BranchBound(kp)
	case "entropy":
		return w27Entropy(kp)
	case "convolutional":
		return w27Convolutional(kp)
	case "ldpc":
		return w27LDPC(kp)
	case "bch":
		return w27BCH(kp)
	case "singleton":
		return w27Singleton(kp)
	case "space":
		return w27Space(kp)
	case "hierarchy":
		return w27Hierarchy(kp)
	case "conp":
		return w27CoNP(kp)
	case "pspace":
		return w27PSPACE(kp)
	case "approx":
		return w27Approx(kp)
	case "fpt":
		return w27FPT(kp)
	case "homogeneous":
		return w27Homogeneous(kp)
	case "characteristic":
		return w27Characteristic(kp)
	case "nonhomogeneous":
		return w27Nonhomogeneous(kp)
	case "zeckendorf":
		return w27Zeckendorf(kp)
	case "cf":
		return w27ContinuedFraction(kp)
	case "ackermann":
		return w27Ackermann(kp)
	case "cardinality":
		return w27Cardinality(kp)
	case "cartesian":
		return w27Cartesian(kp)
	case "classes":
		return w27Classes(kp)
	case "lattice":
		return w27Lattice(kp)
	case "invariant":
		return w27Invariant(kp)
	case "diagonal":
		return w27Diagonal(kp)
	case "descent":
		return w27Descent(kp)
	case "constructive":
		return w27Constructive(kp)
	case "nonconstructive":
		return w27Nonconstructive(kp)
	}
	return w27P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w27N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w27P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w27I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w27Pow(b, e int) int {
	out := 1
	for i := 0; i < e; i++ {
		out *= b
	}
	return out
}

func w27CeilLog2(n int) int {
	out := 0
	v := 1
	for v < n {
		v *= 2
		out++
	}
	return out
}

func w27QuickSort(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Quicksort's worst case on %d elements compares every pair. How many comparisons is that?", n), strconv.Itoa(n*(n-1)/2), "Every pair is compared once.", f)
	case 2:
		return w27P(o, fmt.Sprintf("Balanced quicksort on %d elements recurses to what base-two depth?", n), strconv.Itoa(w27CeilLog2(n)), "Halving until one element remains.", f)
	default:
		return w27P(o, fmt.Sprintf("Quicksort's worst-case recursion on %d elements has what stack depth?", n), strconv.Itoa(n), "The worst case removes one element per level.", f)
	}
}

func w27MergeSort(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Merge sort on %d elements has how many levels of halving?", n), strconv.Itoa(w27CeilLog2(n)), "Each level halves the problem size.", f)
	case 2:
		return w27P(o, fmt.Sprintf("Merge sort on %d elements performs how many merges in total?", n), strconv.Itoa(n-1), "Each merge reduces the number of runs by one.", f)
	default:
		return w27P(o, fmt.Sprintf("Merge sort on %d elements copies how many elements per level?", n), strconv.Itoa(n), "Every level touches every element once.", f)
	}
}

func w27BinarySearch(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Binary search over %d sorted items needs at most how many comparisons?", n), strconv.Itoa(w27CeilLog2(n)), "Each comparison halves the interval.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A binary search halves an interval of %d items. How many halvings are needed in the worst case?", n), strconv.Itoa(w27CeilLog2(n)), "The interval shrinks to one item.", f)
	default:
		return w27P(o, fmt.Sprintf("Binary search on an array of %d items may return which largest index?", n), strconv.Itoa(n-1), "Indices start at zero.", f)
	}
}

func w27Dijkstra(o int) generator.Problem {
	n := w27N(2, 60)
	m := w27N(1, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "m", m)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Dijkstra's algorithm has %d vertices to finalise. How many vertices are extracted?", n), strconv.Itoa(n), "Every vertex is extracted once.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A graph processed by Dijkstra has %d edges. How many edge relaxations may it attempt?", m), strconv.Itoa(m), "Each edge is relaxed from its settled endpoint.", f)
	default:
		return w27P(o, fmt.Sprintf("Dijkstra's algorithm settles %d vertices one at a time. How many settlement steps are there?", n), strconv.Itoa(n), "One step per vertex.", f)
	}
}

func w27Huffman(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A Huffman code has %d symbols. How many leaf nodes does its tree have?", n), strconv.Itoa(n), "One leaf per symbol.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A Huffman tree with %d leaves has how many internal nodes?", n), strconv.Itoa(n-1), "A binary tree with L leaves has L-1 internal nodes.", f)
	default:
		return w27P(o, fmt.Sprintf("A Huffman tree with %d leaves has how many nodes in total?", n), strconv.Itoa(2*n-1), "Leaves plus internal nodes.", f)
	}
}

func w27Backtracking(o int) generator.Problem {
	n := w27N(2, 60)
	k := w27N(2, 12)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Backtracking searches a decision tree of depth %d. How deep is the recursion?", n), strconv.Itoa(n), "One level per decision.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A search tree has depth %d and %d choices per level. How many leaves does it have?", n, k), strconv.Itoa(w27Pow(k, n)), "Each level multiplies the leaves.", f)
	default:
		return w27P(o, fmt.Sprintf("A decision point offers %d choices. How many branches leave it?", k), strconv.Itoa(k), "One branch per choice.", f)
	}
}

func w27BranchBound(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Branch and bound explores a tree of depth %d and prunes one branch per level. How many levels are explored?", n), strconv.Itoa(n), "Every level is visited before pruning.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A branch-and-bound node has %d children. How many bounds are evaluated?", n), strconv.Itoa(n), "One bound per child.", f)
	default:
		return w27P(o, fmt.Sprintf("A bound prunes a subtree at depth %d. What is the smallest number of nodes that bound can remove?", n), "1", "A pruned subtree contains its root.", f)
	}
}

func w27Entropy(o int) generator.Problem {
	k := w27N(2, 10)
	f := map[string]string{}
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A source has %d outcomes. How many probabilities enter its entropy?", k), strconv.Itoa(k), "One probability per outcome.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A uniform source has %d equally likely outcomes. How many bits does each symbol carry?", k), strconv.Itoa(w27CeilLog2(k)), "The entropy is the base-two logarithm.", f)
	default:
		return w27P(o, fmt.Sprintf("The probabilities of %d outcomes sum to what value?", k), "1", "A distribution is normalised.", f)
	}
}

func w27Convolutional(o int) generator.Problem {
	k := w27N(2, 12)
	f := map[string]string{}
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A convolutional code has constraint length %d. How many input bits influence one output window?", k), strconv.Itoa(k), "The constraint length counts the window.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A convolutional encoder has memory %d. How many previous input bits are retained?", k), strconv.Itoa(k), "Memory is the number of stored bits.", f)
	default:
		return w27P(o, fmt.Sprintf("An encoder with %d memory bits has how many states?", k), strconv.Itoa(w27Pow(2, k)), "Each memory bit doubles the states.", f)
	}
}

func w27LDPC(o int) generator.Problem {
	n := w27N(6, 60)
	m := w27N(2, 30)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "m", m)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("An LDPC code has %d parity-check rows. How many checks are performed?", m), strconv.Itoa(m), "One check per row.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A code of block length %d has %d parity checks. How many information bits remain?", n, m), strconv.Itoa(n-m), "Information bits are the block length minus the checks.", f)
	default:
		return w27P(o, fmt.Sprintf("A parity-check matrix has %d rows. How many parity equations does it represent?", m), strconv.Itoa(m), "One equation per row.", f)
	}
}

func w27BCH(o int) generator.Problem {
	t := w27N(1, 30)
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "t", t)
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A BCH code corrects %d errors. How many errors may it correct per block?", t), strconv.Itoa(t), "The design distance gives the correction count.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A BCH code with designed distance %d uses how many parity symbols at minimum?", 2*t+1), strconv.Itoa(2*t), "The minimum distance exceeds twice the correction count.", f)
	default:
		return w27P(o, fmt.Sprintf("A BCH code has block length %d. How many symbols are in one block?", n), strconv.Itoa(n), "The block length counts the symbols.", f)
	}
}

func w27Singleton(o int) generator.Problem {
	n := w27N(3, 60)
	k := w27N(1, 29)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("The Singleton bound for a code of length %d with %d information symbols gives distance at most what value?", n, k), strconv.Itoa(n-k+1), "The bound is n-k+1.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A code of length %d carries %d information symbols. How many redundant symbols does it have?", n, k), strconv.Itoa(n-k), "Redundancy is the difference.", f)
	default:
		return w27P(o, fmt.Sprintf("An MDS code of length %d meets the Singleton bound with equality. How many symbols separate the distance from the redundancy?", n), "1", "Distance equals redundancy plus one.", f)
	}
}

func w27Space(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("An algorithm stores %d values. How many cells does its working memory use?", n), strconv.Itoa(n), "One cell per value.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A logarithmic-space algorithm runs on an input of size %d. What is the base-two logarithm rounded up?", n), strconv.Itoa(w27CeilLog2(n)), "Logarithmic space counts bits.", f)
	default:
		return w27P(o, fmt.Sprintf("An input of %d bits is read once. How many input bits must be stored?", n), strconv.Itoa(n), "Streaming algorithms keep the input only once.", f)
	}
}

func w27Hierarchy(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("The time hierarchy theorem compares time %d with time %d log %d. How many logarithmic factors are added?", n, n, n), "1", "One logarithmic factor separates the classes.", f)
	case 2:
		return w27P(o, fmt.Sprintf("The time hierarchy theorem separates two time classes at input size %d. How many classes does it compare?", n), "2", "The theorem compares a smaller and a larger time class.", f)
	default:
		return w27P(o, fmt.Sprintf("A separation at input size %d is strict when the larger bound exceeds the smaller. Enter 1 when the classes are distinct.", n), "1", "More time buys strictly more problems.", f)
	}
}

func w27CoNP(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A language in coNP has complements whose certificates are checked for %d inputs. How many languages are complemented?", n), strconv.Itoa(n), "One complement per language.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A coNP certificate for an input of size %d is verified in how many passes over the certificate?", n), "1", "Verification is one polynomial pass.", f)
	default:
		return w27P(o, fmt.Sprintf("coNP is the complement class of NP for an input of size %d. How many classes are paired in that definition?", n), "2", "NP and coNP are two classes.", f)
	}
}

func w27PSPACE(o int) generator.Problem {
	k := w27N(1, 12)
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "k", k)
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A polynomial-space machine uses n to the power %d cells on input size %d. What is the exponent?", k, n), strconv.Itoa(k), "The exponent is fixed.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A quantified Boolean formula with %d quantifier alternations is evaluated in polynomial space. How many alternation types are there?", k), "2", "Quantifiers are universal or existential.", f)
	default:
		return w27P(o, fmt.Sprintf("A PSPACE computation on input size %d has a polynomial bound of degree %d. How many terms can the polynomial have at most?", n, k), strconv.Itoa(k+1), "A degree k polynomial has k+1 coefficients.", f)
	}
}

func w27Approx(o int) generator.Problem {
	r := w27N(2, 30)
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "r", r)
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("An approximation guarantee is a factor of %d. What factor multiplies the optimum?", r), strconv.Itoa(r), "The guarantee is the ratio.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A 2-approximation on an instance with %d items returns a value within what factor of optimal?", n), "2", "The factor is fixed by the algorithm.", f)
	default:
		return w27P(o, fmt.Sprintf("A scheme with approximation ratio %d is run on %d instances. How many ratios are computed?", n, n), strconv.Itoa(n), "One ratio per instance.", f)
	}
}

func w27FPT(o int) generator.Problem {
	k := w27N(1, 12)
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "k", k)
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("An FPT algorithm runs in f(k) n^c with parameter %d. What is the parameter value?", k), strconv.Itoa(k), "The parameter is k.", f)
	case 2:
		return w27P(o, fmt.Sprintf("An FPT runtime of f(%d) n^c confines the exponential cost to the parameter. What is the parameter?", k), strconv.Itoa(k), "The exponent is a function of k.", f)
	default:
		return w27P(o, fmt.Sprintf("An FPT algorithm on %d vertices with parameter %d has polynomial degree c. How many constant factors multiply the polynomial?", n, k), "1", "The polynomial factor is n to a constant power.", f)
	}
}

func w27Homogeneous(o int) generator.Problem {
	k := w27N(1, 12)
	f := map[string]string{}
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A homogeneous linear recurrence has %d previous terms. What is its order?", k), strconv.Itoa(k), "The order counts the previous terms.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A homogeneous recurrence of order %d needs how many initial conditions?", k), strconv.Itoa(k), "One condition per order.", f)
	default:
		return w27P(o, fmt.Sprintf("A homogeneous recurrence of order %d has a characteristic polynomial of what degree?", k), strconv.Itoa(k), "The characteristic degree equals the order.", f)
	}
}

func w27Characteristic(o int) generator.Problem {
	r := w27N(1, 12)
	m := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "r", r)
	w27I(f, "m", m)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A characteristic polynomial has %d distinct roots. How many independent exponential solutions does it give?", r), strconv.Itoa(r), "One solution per distinct root.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A root of multiplicity %d contributes a polynomial factor of what degree?", m), strconv.Itoa(m-1), "Multiplicity m gives degree m-1.", f)
	default:
		return w27P(o, fmt.Sprintf("A characteristic polynomial of degree %d has how many roots counted with multiplicity?", r), strconv.Itoa(r), "Degree counts roots with multiplicity.", f)
	}
}

func w27Nonhomogeneous(o int) generator.Problem {
	k := w27N(1, 12)
	d := w27N(1, 12)
	f := map[string]string{}
	w27I(f, "k", k)
	w27I(f, "d", d)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A nonhomogeneous recurrence of order %d has how many particular solutions in its general solution?", k), "1", "One particular solution is added to the homogeneous part.", f)
	case 2:
		return w27P(o, fmt.Sprintf("The homogeneous part of a recurrence has order %d. How many arbitrary constants does it contain?", k), strconv.Itoa(k), "One constant per order.", f)
	default:
		return w27P(o, fmt.Sprintf("A forcing term is a polynomial of degree %d. How many coefficients does the trial solution have?", d), strconv.Itoa(d+1), "A degree d polynomial has d+1 coefficients.", f)
	}
}

func w27Zeckendorf(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A Zeckendorf representation uses %d Fibonacci numbers. How many summands are there?", n), strconv.Itoa(n), "One summand per chosen number.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A Zeckendorf representation of %d uses non-consecutive Fibonacci numbers. How many consecutive pairs may appear?", n), "0", "No two chosen numbers are adjacent.", f)
	default:
		return w27P(o, fmt.Sprintf("Every positive integer at most %d has how many Zeckendorf representations?", n), "1", "The representation is unique.", f)
	}
}

func w27ContinuedFraction(o int) generator.Problem {
	n := w27N(1, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A continued fraction is truncated after %d terms. How many convergents does the truncation produce?", n), strconv.Itoa(n), "One convergent per term.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A continued fraction has %d partial quotients. How many partial quotients are recorded?", n), strconv.Itoa(n), "One quotient per term.", f)
	default:
		return w27P(o, fmt.Sprintf("Consecutive convergents of a continued fraction with %d terms have a determinant of what absolute value?", n), "1", "The determinant identity gives one.", f)
	}
}

func w27Ackermann(o int) generator.Problem {
	m := w27N(0, 3)
	n := w27N(0, 4)
	f := map[string]string{}
	w27I(f, "m", m)
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("The Ackermann function A(%d,%d) is defined by how many recursive branches?", m, n), "2", "One branch reduces the first argument, one the second.", f)
	case 2:
		return w27P(o, fmt.Sprintf("The Ackermann function A(%d,%d) takes how many arguments?", m, n), "2", "Ackermann is a two-argument function.", f)
	default:
		return w27P(o, fmt.Sprintf("Evaluating A(%d,%d) nests calls until the first argument reaches zero. What value of the first argument ends the recursion?", m, n), "0", "The base case has first argument zero.", f)
	}
}

func w27Cardinality(o int) generator.Problem {
	n := w27N(1, 20)
	m := w27N(1, 20)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "m", m)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A finite set has %d elements. How many subsets does it have?", n), strconv.Itoa(w27Pow(2, n)), "Each element is in or out.", f)
	case 2:
		return w27P(o, fmt.Sprintf("Set A has %d elements and set B has %d. How many ordered pairs does A times B contain?", n, m), strconv.Itoa(n*m), "Multiply the sizes.", f)
	default:
		return w27P(o, fmt.Sprintf("A finite set has %d elements. What is its cardinality?", n), strconv.Itoa(n), "Cardinality counts the elements.", f)
	}
}

func w27Cartesian(o int) generator.Problem {
	n := w27N(1, 20)
	m := w27N(1, 20)
	f := map[string]string{}
	w27I(f, "n", n)
	w27I(f, "m", m)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A product set is built from sets of sizes %d and %d. How many pairs does it contain?", n, m), strconv.Itoa(n*m), "The product has n times m pairs.", f)
	case 2:
		return w27P(o, fmt.Sprintf("An ordered pair built from a set of size %d has how many coordinates?", n), "2", "A pair has two coordinates.", f)
	default:
		return w27P(o, fmt.Sprintf("The power set of a product of size %d has how many members?", n*m), strconv.Itoa(w27Pow(2, n*m)), "A set of size k has 2^k subsets.", f)
	}
}

func w27Classes(o int) generator.Problem {
	k := w27N(2, 30)
	f := map[string]string{}
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A set is partitioned into %d equivalence classes. How many classes are there?", k), strconv.Itoa(k), "The partition lists every class.", f)
	case 2:
		return w27P(o, fmt.Sprintf("An equivalence relation partitions its set into %d blocks. How many parts does the partition have?", k), strconv.Itoa(k), "One part per class.", f)
	default:
		return w27P(o, fmt.Sprintf("Every element of a set with %d classes is related to itself. How many reflexive pairs per element?", k), "1", "Reflexivity relates an element to itself once.", f)
	}
}

func w27Lattice(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A lattice has %d elements. How many elements belong to it?", n), strconv.Itoa(n), "The order is the element count.", f)
	case 2:
		return w27P(o, fmt.Sprintf("Each pair in a lattice of %d elements has a join and a meet. How many binary operations does that define?", n), "2", "Join and meet are two operations.", f)
	default:
		return w27P(o, fmt.Sprintf("A bounded lattice has %d bounds. How many bounds are there?", n), "2", "A bottom and a top.", f)
	}
}

func w27Invariant(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A loop invariant is checked before and after each of %d iterations. How many checks does that make?", n), strconv.Itoa(2*n), "Two checks per iteration.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A loop invariant holds initially and after every one of %d iterations. How many invariant states are proved?", n+1), strconv.Itoa(n+1), "The initial state plus one per iteration.", f)
	default:
		return w27P(o, fmt.Sprintf("An invariant is preserved by each of %d loop bodies. How many preservation steps are proved?", n), strconv.Itoa(n), "One step per iteration.", f)
	}
}

func w27Diagonal(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("Cantor's diagonal argument changes the %d-th digit of the %d-th listed number. How many digits are inspected?", n, n), strconv.Itoa(n), "One digit per listed number.", f)
	case 2:
		return w27P(o, fmt.Sprintf("The diagonal construction differs from a listing of %d numbers at how many positions?", n), strconv.Itoa(n), "It differs at every diagonal position.", f)
	default:
		return w27P(o, fmt.Sprintf("Cantor's argument with %d listed numbers shows the reals cannot be listed. Enter 1 when the reals are uncountable.", n), "1", "No listing exhausts the reals.", f)
	}
}

func w27Descent(o int) generator.Problem {
	k := w27N(1, 60)
	f := map[string]string{}
	w27I(f, "k", k)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("An infinite descent proof constructs a strictly decreasing sequence of %d natural numbers. How many steps may it descend?", k), strconv.Itoa(k), "It descends while remaining non-negative.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A strictly decreasing sequence of %d natural numbers has how many infinite tails?", k), "0", "Well-ordering forbids an infinite descent.", f)
	default:
		return w27P(o, fmt.Sprintf("A descent argument reduces the counterexample after %d steps. How many steps are taken before the contradiction?", k), strconv.Itoa(k), "One reduction per step.", f)
	}
}

func w27Constructive(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A constructive proof exhibits one witness for each of %d cases. How many witnesses are produced?", n), strconv.Itoa(n), "One witness per case.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A constructive existence proof builds %d objects. How many objects are constructed?", n), strconv.Itoa(n), "The proof builds each object.", f)
	default:
		return w27P(o, fmt.Sprintf("A construction runs for %d steps and returns one object. How many objects are returned?", n), "1", "The algorithm returns a single witness.", f)
	}
}

func w27Nonconstructive(o int) generator.Problem {
	n := w27N(2, 60)
	f := map[string]string{}
	w27I(f, "n", n)
	switch o {
	case 1:
		return w27P(o, fmt.Sprintf("A nonconstructive proof shows existence by contradiction over %d cases. How many witnesses are exhibited?", n), "0", "Existence is shown without a witness.", f)
	case 2:
		return w27P(o, fmt.Sprintf("A nonconstructive argument assumes the negation and derives a contradiction in %d steps. How many contradictions are derived?", n), "1", "A single contradiction suffices.", f)
	default:
		return w27P(o, fmt.Sprintf("A counting argument shows that %d objects cannot all avoid a property. How many objects must have the property at least?", n), "1", "At least one object has it.", f)
	}
}
