package abstract

// Wave 23 abstract algebra: group structure, commutative algebra, homological
// algebra, field theory and representation theory. Every family returns an
// integer and publishes its drawn parameters in Facts.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave23(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"abstract.group.conjugacy":              "conjugacy",
		"abstract.group.normalizer":             "normalizer",
		"abstract.group.derived_series":         "derived",
		"abstract.group.frattini":               "frattini",
		"abstract.group.hall":                   "hall",
		"abstract.group.exponent":               "exponent",
		"abstract.group.transfer":               "transfer",
		"abstract.group.cohomology":             "cohomology",
		"abstract.group.burnside_lemma":         "burnside",
		"abstract.group.permutation_character":  "permchar",
		"abstract.rings.groebner":               "groebner",
		"abstract.rings.regular":                "regular",
		"abstract.rings.krull":                  "krull",
		"abstract.rings.integral_closure":       "integral",
		"abstract.rings.dvr":                    "dvr",
		"abstract.rings.artinian":               "artinian",
		"abstract.rings.local":                  "localt",
		"abstract.rings.symmetric_algebra":      "symalg",
		"abstract.rings.flat":                   "flat",
		"abstract.rings.completion":             "completion",
		"abstract.structures.projective":        "projective",
		"abstract.structures.injective":         "injective",
		"abstract.structures.spectral_sequence": "spectral",
		"abstract.structures.chain_homotopy":    "homotopy",
		"abstract.structures.derived_functor":   "derivedfun",
		"abstract.structures.exact_functor":     "exactfun",
		"abstract.structures.direct_limit":      "directlimit",
		"abstract.structures.free_resolution":   "resolution",
		"abstract.field.perfect":                "perfect",
		"abstract.field.normal":                 "normal",
		"abstract.field.primitive_element":      "primitive",
		"abstract.field.artin":                  "artin",
		"abstract.field.kummer":                 "kummer",
		"abstract.field.hilbert90":              "hilbert90",
		"abstract.field.local":                  "localfield",
		"abstract.field.padic":                  "padic",
		"abstract.rep.orthogonality":            "orthogonality",
		"abstract.rep.regular":                  "regularrep",
		"abstract.rep.tensor":                   "tensor",
		"abstract.rep.restriction":              "restriction",
		"abstract.rep.frobenius":                "frobenius",
		"abstract.rep.fourier":                  "fourier",
		"abstract.rep.young":                    "young",
		"abstract.rep.burnside_pq":              "burnsidepq",
	} {
		reg.Register(id, &w23Gen{kind: kind})
	}
}

type w23Gen struct{ kind string }

func (g *w23Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "conjugacy":
		return w23Conjugacy(kp)
	case "normalizer":
		return w23Normalizer(kp)
	case "derived":
		return w23Derived(kp)
	case "frattini":
		return w23Frattini(kp)
	case "hall":
		return w23Hall(kp)
	case "exponent":
		return w23Exponent(kp)
	case "transfer":
		return w23Transfer(kp)
	case "cohomology":
		return w23Cohomology(kp)
	case "burnside":
		return w23Burnside(kp)
	case "permchar":
		return w23PermChar(kp)
	case "groebner":
		return w23Groebner(kp)
	case "regular":
		return w23Regular(kp)
	case "krull":
		return w23Krull(kp)
	case "integral":
		return w23Integral(kp)
	case "dvr":
		return w23DVR(kp)
	case "artinian":
		return w23Artinian(kp)
	case "localt":
		return w23Local(kp)
	case "symalg":
		return w23SymAlg(kp)
	case "flat":
		return w23Flat(kp)
	case "completion":
		return w23Completion(kp)
	case "projective":
		return w23Projective(kp)
	case "injective":
		return w23Injective(kp)
	case "spectral":
		return w23Spectral(kp)
	case "homotopy":
		return w23Homotopy(kp)
	case "derivedfun":
		return w23DerivedFunctor(kp)
	case "exactfun":
		return w23ExactFunctor(kp)
	case "directlimit":
		return w23DirectLimit(kp)
	case "resolution":
		return w23Resolution(kp)
	case "perfect":
		return w23Perfect(kp)
	case "normal":
		return w23Normal(kp)
	case "primitive":
		return w23Primitive(kp)
	case "artin":
		return w23Artin(kp)
	case "kummer":
		return w23Kummer(kp)
	case "hilbert90":
		return w23Hilbert90(kp)
	case "localfield":
		return w23LocalField(kp)
	case "padic":
		return w23PAdic(kp)
	case "orthogonality":
		return w23Orthogonality(kp)
	case "regularrep":
		return w23RegularRep(kp)
	case "tensor":
		return w23Tensor(kp)
	case "restriction":
		return w23Restriction(kp)
	case "frobenius":
		return w23Frobenius(kp)
	case "fourier":
		return w23Fourier(kp)
	case "young":
		return w23Young(kp)
	case "burnsidepq":
		return w23BurnsidePQ(kp)
	}
	return w23P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w23N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w23P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w23I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w23Pow(b, e int) int {
	out := 1
	for i := 0; i < e; i++ {
		out *= b
	}
	return out
}

func w23Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func w23Factorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func w23LCM(a, b int) int {
	x, y := a, b
	for y != 0 {
		x, y = y, x%y
	}
	return a / x * b
}

func w23Totient(n int) int {
	out := n
	m := n
	for p := 2; p*p <= m; p++ {
		if m%p == 0 {
			for m%p == 0 {
				m /= p
			}
			out -= out / p
		}
	}
	if m > 1 {
		out -= out / m
	}
	return out
}

func w23Prime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func w23Primes(lo, hi int) []int {
	out := []int{}
	for i := lo; i <= hi; i++ {
		if w23Prime(i) {
			out = append(out, i)
		}
	}
	return out
}

func w23Partitions(n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	for part := 1; part <= n; part++ {
		for s := part; s <= n; s++ {
			dp[s] += dp[s-part]
		}
	}
	return dp[n]
}

func w23Conjugacy(o int) generator.Problem {
	n := w23N(2, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("How many conjugacy classes does the symmetric group S_%d have?", n), strconv.Itoa(w23Partitions(n)), "Conjugacy classes of S_n are indexed by partitions of n.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many elements lie in the conjugacy class of transpositions in S_%d?", n), strconv.Itoa(w23Binom(n, 2)), "A transposition is chosen by its unordered pair of points.", f)
	default:
		return w23P(o, fmt.Sprintf("How many elements does the symmetric group S_%d have?", n), strconv.Itoa(w23Factorial(n)), "The order is n factorial.", f)
	}
}

func w23Normalizer(o int) generator.Problem {
	n := w23N(3, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("How many elements does the centraliser of an n-cycle in S_%d have?", n), strconv.Itoa(n), "The centraliser of an n-cycle is the cyclic group it generates.", f)
	case 2:
		return w23P(o, fmt.Sprintf("What is the index of the centraliser of an n-cycle in S_%d?", n), strconv.Itoa(w23Factorial(n)/n), "The index is the class size.", f)
	default:
		return w23P(o, fmt.Sprintf("What is the order of the centre of S_%d?", n), "1", "The symmetric group has trivial centre.", f)
	}
}

func w23Derived(o int) generator.Problem {
	n := w23N(3, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	length := 2
	if n == 4 {
		length = 3
	}
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the derived length of the symmetric group S_%d?", n), strconv.Itoa(length), "S_4 needs three steps; every other S_n needs two.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many elements does the derived subgroup of S_%d have?", n), strconv.Itoa(w23Factorial(n)/2), "The derived subgroup is the alternating group.", f)
	default:
		return w23P(o, fmt.Sprintf("What is the derived length of an abelian group acting on S_%d?", n), "1", "An abelian group has trivial derived subgroup.", f)
	}
}

func w23Frattini(o int) generator.Problem {
	p := w23Primes(2, 13)[rand.Intn(len(w23Primes(2, 13)))]
	k := w23N(1, 7)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A group of order p^%d has Frattini subgroup of index p. If p=%d, what is the order of the Frattini subgroup?", k, p), strconv.Itoa(w23Pow(p, k-1)), "The quotient has order p.", f)
	case 2:
		return w23P(o, fmt.Sprintf("For a group of order p^%d with p=%d, what is the order of the Frattini quotient?", k, p), strconv.Itoa(p), "The quotient is elementary abelian of rank k.", f)
	default:
		return w23P(o, fmt.Sprintf("How many generators does a group of order p^%d with p=%d need at minimum?", k, p), strconv.Itoa(k), "The minimal generator count equals the rank.", f)
	}
}

func w23Hall(o int) generator.Problem {
	ps := w23Primes(2, 13)
	p := ps[rand.Intn(len(ps))]
	q := ps[rand.Intn(len(ps))]
	a := w23N(1, 4)
	b := w23N(1, 4)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "q", q)
	w23I(f, "a", a)
	w23I(f, "b", b)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A group has order p^%d q^%d with p=%d. What is the order of its Hall p-subgroup?", a, b, p), strconv.Itoa(w23Pow(p, a)), "The Hall p-subgroup collects the full p-part.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A group of order p^%d q^%d with q=%d has a Hall p-subgroup. What is its index?", a, b, q), strconv.Itoa(w23Pow(q, b)), "The index is the q-part.", f)
	default:
		return w23P(o, fmt.Sprintf("How many distinct primes divide the order of a group of order p^%d q^%d?", a, b), "2", "Two prime powers appear.", f)
	}
}

func w23Exponent(o int) generator.Problem {
	n := w23N(2, 15)
	f := map[string]string{}
	w23I(f, "n", n)
	l := 1
	for i := 2; i <= n; i++ {
		l = w23LCM(l, i)
	}
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the exponent of the cyclic group C_%d?", n), strconv.Itoa(n), "A cyclic group of order n has exponent n.", f)
	case 2:
		return w23P(o, fmt.Sprintf("What is the exponent of the symmetric group S_%d?", n), strconv.Itoa(l), "The exponent is the least common multiple of the cycle lengths.", f)
	default:
		return w23P(o, fmt.Sprintf("What is the exponent of the dihedral group D_%d?", n), strconv.Itoa(w23LCM(2, n)), "Rotations and reflections have orders dividing n and 2.", f)
	}
}

func w23Transfer(o int) generator.Problem {
	p := w23Primes(3, 97)[rand.Intn(len(w23Primes(3, 97)))]
	n := w23N(3, 30)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A Sylow p-subgroup with p=%d has index congruent to one modulo p. What is the smallest index greater than one?", p), strconv.Itoa(p+1), "The next integer congruent to one modulo p is p+1.", f)
	case 2:
		return w23P(o, fmt.Sprintf("What is the order of the abelianisation of S_%d?", n), "2", "The abelianisation of S_n is cyclic of order two.", f)
	default:
		return w23P(o, fmt.Sprintf("The transfer map for S_%d lands in an abelian group. What is the exponent of that target?", n), "2", "The sign character generates the abelianisation.", f)
	}
}

func w23Cohomology(o int) generator.Problem {
	r := w23N(1, 15)
	f := map[string]string{}
	w23I(f, "r", r)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A free group has rank %d. What is the dimension of its first cohomology with trivial integer coefficients?", r), strconv.Itoa(r), "H^1 of a free group has rank r.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A free group of rank %d acts trivially. What is the dimension of its zeroth cohomology?", r), "1", "H^0 is the invariant module, here Z.", f)
	default:
		return w23P(o, fmt.Sprintf("What is the Euler characteristic of a free group of rank %d?", r), strconv.Itoa(1-r), "A graph with r loops has characteristic 1-r.", f)
	}
}

func w23Burnside(o int) generator.Problem {
	n := w23N(2, 30)
	m := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "m", m)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A group of order %d acts transitively on %d points. How many orbits does the action have?", n, m), "1", "Transitivity means one orbit.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The identity of a group of order %d acts on a set of %d points. How many points does it fix?", n, m), strconv.Itoa(m), "The identity fixes every point.", f)
	default:
		return w23P(o, fmt.Sprintf("Burnside's lemma divides the fixed-point sum by the group order. For a group of order %d acting on %d points, what divisor is used?", n, m), strconv.Itoa(n), "The divisor is the group order.", f)
	}
}

func w23PermChar(o int) generator.Problem {
	n := w23N(3, 15)
	k := w23N(1, n-1)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("The natural permutation character of S_%d has what degree?", n), strconv.Itoa(n), "The degree is the number of points.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The permutation character on k-subsets of an n-set with n=%d and k=%d has what degree?", n, k), strconv.Itoa(w23Binom(n, k)), "The degree counts the k-subsets.", f)
	default:
		return w23P(o, fmt.Sprintf("The natural permutation character of S_%d contains the trivial character with what multiplicity?", n), "1", "The action is transitive, so the trivial multiplicity is one.", f)
	}
}

func w23Groebner(o int) generator.Problem {
	a := w23N(2, 20)
	b := w23N(2, 20)
	f := map[string]string{}
	w23I(f, "a", a)
	w23I(f, "b", b)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A monomial ideal is generated by x^%d and y^%d. How many standard monomials are there?", a, b), strconv.Itoa(a*b), "Standard monomials are x^i y^j with i<a and j<b.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A monomial ideal in two variables has %d standard monomials in x. How many standard monomials of the form x^i y^j does the staircase hold in one row?", a), strconv.Itoa(a), "One per permitted x exponent.", f)
	default:
		return w23P(o, fmt.Sprintf("A zero-dimensional monomial ideal has Hilbert function one at degree zero regardless of its %d rows. What is that value?", b), "1", "The degree-zero component is the field.", f)
	}
}

func w23Regular(o int) generator.Problem {
	d := w23N(1, 15)
	f := map[string]string{}
	w23I(f, "d", d)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A regular local ring has dimension %d. What is its embedding dimension?", d), strconv.Itoa(d), "Regularity makes embedding dimension equal dimension.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many minimal generators does the maximal ideal of a regular local ring of dimension %d need?", d), strconv.Itoa(d), "The maximal ideal is generated by a regular system of parameters.", f)
	default:
		return w23P(o, fmt.Sprintf("A regular local ring of dimension %d has residue field of what dimension over itself?", d), "1", "Every field is one-dimensional over itself.", f)
	}
}

func w23Krull(o int) generator.Problem {
	n := w23N(1, 15)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the Krull dimension of the polynomial ring k[x_1,...,x_%d]?", n), strconv.Itoa(n), "Each variable adds one to the dimension.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The ring of integers has Krull dimension one. How many prime ideals does a chain of maximal length have beyond the zero ideal for a number field of degree %d?", n), "1", "The ring of integers has dimension one regardless of degree.", f)
	default:
		return w23P(o, fmt.Sprintf("A field of transcendence degree %d over its prime field has how many prime ideals in a maximal chain?", n), "1", "A field has dimension zero, so one prime ideal.", f)
	}
}

func w23Integral(o int) generator.Problem {
	n := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("How many Z-basis elements does the integral closure of Z in a quadratic field of degree %d have?", n), "2", "Every quadratic order is a rank-two Z-module.", f)
	case 2:
		return w23P(o, fmt.Sprintf("An algebraic integer of degree %d generates an extension of what degree over Q?", n), strconv.Itoa(n), "The degree is the degree of the minimal polynomial.", f)
	default:
		return w23P(o, fmt.Sprintf("The integral closure of Z in Q is Z. For an algebraic integer of degree %d, how many Z-basis elements does that closure have?", n), "1", "Z is free of rank one.", f)
	}
}

func w23DVR(o int) generator.Problem {
	ps := w23Primes(2, 97)
	p := ps[rand.Intn(len(ps))]
	f := map[string]string{}
	w23I(f, "p", p)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A discrete valuation ring has residue field of size %d. What is the valuation of a uniformiser?", p), "1", "A uniformiser has valuation one.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A discrete valuation ring has residue field F_p with p=%d. How many elements does the residue field have?", p), strconv.Itoa(p), "The residue field is the prime field.", f)
	default:
		return w23P(o, fmt.Sprintf("The value group of a discrete valuation ring with residue field of size %d is isomorphic to Z. What is its rank?", p), "1", "The value group has rank one.", f)
	}
}

func w23Artinian(o int) generator.Problem {
	n := w23N(1, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the length of the Artinian ring k[x]/(x^%d)?", n), strconv.Itoa(n), "The composition series has one step per power.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many maximal ideals does the local Artinian ring k[x]/(x^%d) have?", n), "1", "A local ring has one maximal ideal.", f)
	default:
		return w23P(o, fmt.Sprintf("The ring k^%d decomposes into how many simple components?", n), strconv.Itoa(n), "One field factor per component.", f)
	}
}

func w23Local(o int) generator.Problem {
	n := w23N(1, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A local ring of dimension %d has how many maximal ideals?", n), "1", "Locality means a unique maximal ideal.", f)
	case 2:
		primes := w23Primes(2, 97)
		p := primes[w23N(0, len(primes)-1)]
		w23I(f, "p", p)
		return w23P(o, fmt.Sprintf("A local ring has residue field with %d elements. How many elements does that field have?", p), strconv.Itoa(p), "The residue field is the stated finite field.", f)
	default:
		return w23P(o, fmt.Sprintf("A local ring of dimension %d has a maximal ideal whose height equals the dimension. What is that height?", n), strconv.Itoa(n), "In a local ring the maximal ideal has height equal to the dimension.", f)
	}
}

func w23SymAlg(o int) generator.Problem {
	n := w23N(2, 30)
	d := w23N(1, 6)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "d", d)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the dimension of the degree-%d symmetric power of a %d-dimensional space?", d, n), strconv.Itoa(w23Binom(n+d-1, d)), "Symmetric powers are counted by multiset combinations.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The first symmetric power of a %d-dimensional space is the space itself. What is its dimension?", n), strconv.Itoa(n), "The first symmetric power is the original space.", f)
	default:
		return w23P(o, fmt.Sprintf("The degree-zero symmetric power of a %d-dimensional space has what dimension?", n), "1", "The zeroth symmetric power is the base field.", f)
	}
}

func w23Flat(o int) generator.Problem {
	n := w23N(2, 40)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("Is the ring Z flat as a module over itself? Enter 1 for yes and 0 for no (modulus %d).", n), "1", "A ring is flat over itself.", f)
	case 2:
		return w23P(o, fmt.Sprintf("Is the field Q flat as a Z-module? Enter 1 for yes and 0 for no (compare with Z/%d).", n), "1", "Q is a localisation of Z and therefore flat.", f)
	default:
		return w23P(o, fmt.Sprintf("Is Z/%d flat as a Z-module? Enter 1 for yes and 0 for no.", n), "0", "A nontrivial finite quotient has torsion and is not flat.", f)
	}
}

func w23Completion(o int) generator.Problem {
	ps := w23Primes(2, 97)
	p := ps[rand.Intn(len(ps))]
	k := w23N(1, 8)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("The ring Z/p^k has order p^k. With p=%d and k=%d, what is its order?", p, k), strconv.Itoa(w23Pow(p, k)), "Raise p to the k-th power.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The residue field of the p-adic integers with p=%d has how many elements?", p), strconv.Itoa(p), "The residue field is F_p.", f)
	default:
		return w23P(o, fmt.Sprintf("The k-th stage of the p-adic completion has p=%d and k=%d. How many stages precede the limit?", k, p), strconv.Itoa(k), "One stage per exponent.", f)
	}
}

func w23Projective(o int) generator.Problem {
	n := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A free module has rank %d. What is its projective dimension?", n), "0", "Free modules are projective.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A free module of rank %d has how many basis elements?", n), strconv.Itoa(n), "Rank counts the basis.", f)
	default:
		return w23P(o, fmt.Sprintf("A projective module of rank %d over a local ring is free. How many generators does it need at minimum?", n), strconv.Itoa(n), "Over a local ring projective implies free.", f)
	}
}

func w23Injective(o int) generator.Problem {
	n := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("How many of %d copies of Q are injective as Z-modules?", n), strconv.Itoa(n), "Q is divisible and therefore injective.", f)
	case 2:
		return w23P(o, fmt.Sprintf("Is Z injective as a Z-module? Enter 1 for yes and 0 for no (compare with %d copies of Q).", n), "0", "Z is not divisible.", f)
	default:
		return w23P(o, fmt.Sprintf("The injective dimension of a field of characteristic zero is zero. What value is that, checked against %d copies of Q?", n), "0", "A field is injective over itself.", f)
	}
}

func w23Spectral(o int) generator.Problem {
	r := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "r", r)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("On page r=%d of a spectral sequence, what is the square of the page index?", r), strconv.Itoa(r*r), "The differential on page r has bidegree of total degree one.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The differential d_r on page r=%d has what total degree?", r), "1", "Every differential raises total degree by one.", f)
	default:
		return w23P(o, fmt.Sprintf("A spectral sequence stabilises after page %d. How many pages does it print before stabilising?", r), strconv.Itoa(r), "One page per index.", f)
	}
}

func w23Homotopy(o int) generator.Problem {
	n := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A chain complex has %d terms. How many components does a chain homotopy have?", n), strconv.Itoa(n), "One component per term.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A chain complex with %d terms has how many boundary maps?", n), strconv.Itoa(n-1), "Maps join consecutive terms.", f)
	default:
		return w23P(o, fmt.Sprintf("A chain complex with %d terms has how many homology groups?", n), strconv.Itoa(n), "One homology group per term.", f)
	}
}

func w23DerivedFunctor(o int) generator.Problem {
	n := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A module has a projective resolution of length %d. How many derived functors are nonzero at most?", n), strconv.Itoa(n+1), "Derived functors run from degree zero to n.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A resolution of length %d has how many projective terms?", n), strconv.Itoa(n+1), "The resolution starts at degree zero.", f)
	default:
		return w23P(o, fmt.Sprintf("A module of projective dimension %d has Ext vanishing above which degree? Enter that degree.", n), strconv.Itoa(n), "Ext vanishes above the projective dimension.", f)
	}
}

func w23ExactFunctor(o int) generator.Problem {
	n := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("An exact sequence has %d terms. How many maps does it contain?", n), strconv.Itoa(n-1), "Maps join consecutive terms.", f)
	case 2:
		return w23P(o, fmt.Sprintf("An exact sequence has %d terms. At how many inner terms can exactness be checked?", n), strconv.Itoa(n-2), "The two end terms have no incoming or outgoing pair.", f)
	default:
		return w23P(o, fmt.Sprintf("A functor preserves a sequence of %d terms. How many terms are compared?", n), strconv.Itoa(n), "Exactness is term by term.", f)
	}
}

func w23DirectLimit(o int) generator.Problem {
	n := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A direct limit has %d stages. How many transition maps does it use?", n), strconv.Itoa(n-1), "Maps join consecutive stages.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A direct system has %d objects. How many objects does it contain?", n), strconv.Itoa(n), "One per stage.", f)
	default:
		return w23P(o, fmt.Sprintf("A direct limit of %d copies of Z under multiplication by two has what rank?", n), "1", "The colimit is the dyadic rationals, of rank one.", f)
	}
}

func w23Resolution(o int) generator.Problem {
	n := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A free resolution of length %d has how many free modules?", n), strconv.Itoa(n+1), "The resolution starts at degree zero.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A free resolution of length %d has how many maps between its modules?", n), strconv.Itoa(n), "One map per positive degree.", f)
	default:
		return w23P(o, fmt.Sprintf("The alternating sum of Betti numbers of a length-%d resolution of a module equals its Euler characteristic. What is that value for a module of rank one?", n), "1", "The Euler characteristic of a rank-one module is one.", f)
	}
}

func w23Perfect(o int) generator.Problem {
	ps := w23Primes(2, 97)
	p := ps[rand.Intn(len(ps))]
	e := w23N(1, 5)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "e", e)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A perfect field of characteristic p=%d has inseparable degree one for a finite extension. What is that degree?", p), "1", "Perfect fields have no inseparable part.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A finite field of characteristic p=%d is perfect. Enter 1 for yes and 0 for no.", p), "1", "Finite fields are perfect.", f)
	default:
		return w23P(o, fmt.Sprintf("An extension has inseparable degree p^e with p=%d and e=%d. What is that degree?", p, e), strconv.Itoa(w23Pow(p, e)), "Raise p to the e-th power.", f)
	}
}

func w23Normal(o int) generator.Problem {
	n := w23N(3, 40)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the degree of the cyclotomic field Q(zeta_%d) over Q?", n), strconv.Itoa(w23Totient(n)), "The degree is Euler's totient.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many automorphisms does the cyclotomic field Q(zeta_%d) have?", n), strconv.Itoa(w23Totient(n)), "Each unit modulo n gives one automorphism.", f)
	default:
		return w23P(o, fmt.Sprintf("The normal closure of Q(zeta_%d) is itself. What is its degree over Q?", n), strconv.Itoa(w23Totient(n)), "Cyclotomic extensions are normal.", f)
	}
}

func w23Primitive(o int) generator.Problem {
	n := w23N(2, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A finite separable extension of degree %d has how many primitive elements at minimum?", n), "1", "The primitive element theorem guarantees one.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A simple extension generated by an element of degree %d has what degree?", n), strconv.Itoa(n), "The degree is the degree of the minimal polynomial.", f)
	default:
		return w23P(o, fmt.Sprintf("A simple extension of degree %d is generated by how many elements?", n), "1", "Simple means one generator.", f)
	}
}

func w23Artin(o int) generator.Problem {
	n := w23N(2, 20)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A finite group of order %d acts faithfully on a field. What is the degree of the fixed field extension?", n), strconv.Itoa(n), "Artin's theorem makes the degree the group order.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A group of order %d acts on a field. How many elements does the group have?", n), strconv.Itoa(n), "The order counts the elements.", f)
	default:
		return w23P(o, fmt.Sprintf("The fixed field of a group of order %d has index how large in the ambient field extension?", n), strconv.Itoa(n), "The index equals the group order.", f)
	}
}

func w23Kummer(o int) generator.Problem {
	k := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A Kummer extension adjoins a k-th root with k=%d. What is its degree over the base field?", k), strconv.Itoa(k), "The degree divides k and equals k in the generic case.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The Galois group of a Kummer extension of degree %d has how many elements?", k), strconv.Itoa(k), "The group is cyclic of order k.", f)
	case 3:
		return w23P(o, fmt.Sprintf("A k-th root equation with k=%d has how many roots in a splitting field?", k), strconv.Itoa(k), "One root for each k-th root of unity phase.", f)
	}
	return w23Kummer(1)
}

func w23Hilbert90(o int) generator.Problem {
	n := w23N(2, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A cyclic extension of degree %d has Hilbert 90 kernel of what size?", n), strconv.Itoa(n), "The kernel is the image of the norm-one subgroup, of size n.", f)
	case 2:
		return w23P(o, fmt.Sprintf("Hilbert's theorem 90 makes a first cohomology group vanish for a cyclic extension of degree %d. What is its size?", n), "1", "A vanishing group has one element.", f)
	default:
		return w23P(o, fmt.Sprintf("The norm map of a cyclic extension of degree %d has image of what index in the base field units?", n), strconv.Itoa(n), "The cokernel has order n.", f)
	}
}

func w23LocalField(o int) generator.Problem {
	ps := w23Primes(2, 13)
	p := ps[rand.Intn(len(ps))]
	e := w23N(1, 20)
	ff := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "e", e)
	w23I(f, "f", ff)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A local field extension has ramification index %d and residue degree %d. What is its total degree?", e, ff), strconv.Itoa(e*ff), "The total degree is the product.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A local field has residue field with p=%d and residue degree %d. How many elements does the residue field have?", p, ff), strconv.Itoa(w23Pow(p, ff)), "The residue field has p^f elements.", f)
	default:
		return w23P(o, fmt.Sprintf("A totally ramified local field extension has residue degree one and ramification index %d. What is its total degree?", e), strconv.Itoa(e), "The residue degree is one.", f)
	}
}

func w23PAdic(o int) generator.Problem {
	ps := w23Primes(2, 97)
	p := ps[rand.Intn(len(ps))]
	k := w23N(1, 12)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("What is the p-adic valuation of p^k with p=%d and k=%d?", p, k), strconv.Itoa(k), "The valuation counts the power.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The p-adic integers with p=%d have residue field of what size?", p), strconv.Itoa(p), "The residue field is F_p.", f)
	default:
		return w23P(o, fmt.Sprintf("A p-adic expansion is truncated after k=%d digits with p=%d. How many digits are kept?", k, p), strconv.Itoa(k), "The truncation keeps k digits.", f)
	}
}

func w23Orthogonality(o int) generator.Problem {
	n := w23N(2, 60)
	c := w23N(1, 12)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "c", c)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("For a finite group of order %d, what is the sum of the squares of its irreducible character degrees?", n), strconv.Itoa(n), "The sum of squares equals the group order.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A finite group has %d conjugacy classes. How many irreducible representations does it have over the complex numbers?", c), strconv.Itoa(c), "Irreducible representations are indexed by conjugacy classes.", f)
	default:
		return w23P(o, fmt.Sprintf("A group of order %d has a trivial representation of what dimension?", n), "1", "The trivial representation is one-dimensional.", f)
	}
}

func w23RegularRep(o int) generator.Problem {
	n := w23N(2, 30)
	d := w23N(1, 30)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "d", d)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("The regular representation of a group of order %d has what dimension?", n), strconv.Itoa(n), "The regular representation acts on the group algebra.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The trivial representation appears in the regular representation of a group of order %d with what multiplicity?", n), "1", "The trivial character appears once.", f)
	default:
		return w23P(o, fmt.Sprintf("An irreducible representation of dimension %d appears in the regular representation with what multiplicity?", d), strconv.Itoa(d), "Each irreducible appears with multiplicity equal to its degree.", f)
	}
}

func w23Tensor(o int) generator.Problem {
	m := w23N(2, 20)
	n := w23N(2, 20)
	f := map[string]string{}
	w23I(f, "m", m)
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("Representations of dimensions %d and %d are tensored. What is the dimension of the tensor product?", m, n), strconv.Itoa(m*n), "Dimensions multiply.", f)
	case 2:
		return w23P(o, fmt.Sprintf("What is the dimension of the symmetric square of a representation of dimension %d?", m), strconv.Itoa(m*(m+1)/2), "Symmetric pairs are counted without order.", f)
	default:
		return w23P(o, fmt.Sprintf("What is the dimension of the exterior square of a representation of dimension %d?", m), strconv.Itoa(m*(m-1)/2), "Exterior pairs are counted without repetition.", f)
	}
}

func w23Restriction(o int) generator.Problem {
	k := w23N(1, 30)
	n := k * w23N(1, 10)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A representation of a group of order %d is restricted to a subgroup of index %d. What is the dimension after restriction?", n, k), strconv.Itoa(n/k), "The dimension is unchanged; the index counts cosets.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A subgroup has index %d in a group of order %d. How many cosets does it have?", k, n), strconv.Itoa(k), "The index counts the cosets.", f)
	default:
		return w23P(o, fmt.Sprintf("A representation is restricted to a subgroup of index %d. How many copies of the subgroup act on the restricted space?", k), "1", "Restriction keeps one copy of the subgroup.", f)
	}
}

func w23Frobenius(o int) generator.Problem {
	d := w23N(1, 30)
	k := w23N(1, 30)
	f := map[string]string{}
	w23I(f, "d", d)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A representation of dimension %d is induced from a subgroup of index %d. What is the induced dimension?", d, k), strconv.Itoa(d*k), "Induction multiplies by the index.", f)
	case 2:
		return w23P(o, fmt.Sprintf("A subgroup has index %d. How many cosets appear in the induced representation?", k), strconv.Itoa(k), "One copy per coset.", f)
	default:
		return w23P(o, fmt.Sprintf("Frobenius reciprocity pairs an induced representation of dimension %d with a restriction. What is the multiplicity of the trivial constituent?", d), "1", "The trivial character pairs once.", f)
	}
}

func w23Fourier(o int) generator.Problem {
	n := w23N(2, 40)
	f := map[string]string{}
	w23I(f, "n", n)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A finite abelian group of order %d has how many characters?", n), strconv.Itoa(n), "The dual group has the same order.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The inverse Fourier transform on a group of order %d sums how many terms per coefficient?", n), strconv.Itoa(n), "One term per character.", f)
	case 3:
		return w23P(o, fmt.Sprintf("Plancherel's identity for a group of order %d sums how many squared coefficients?", n), strconv.Itoa(n), "One per group element.", f)
	}
	return w23Fourier(1)
}

func w23Young(o int) generator.Problem {
	n := w23N(2, 15)
	k := w23N(1, n-1)
	f := map[string]string{}
	w23I(f, "n", n)
	w23I(f, "k", k)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A two-row Young diagram has shape (%d,%d). How many boxes does it contain?", n-k, k), strconv.Itoa(n), "The parts sum to n.", f)
	case 2:
		return w23P(o, fmt.Sprintf("The number of standard tableaux of two-row shape (%d,%d) is the hook-length formula. What is the binomial coefficient C(%d,%d) that bounds it?", n-k, k, n, k), strconv.Itoa(w23Binom(n, k)), "The hook-length count is a difference of binomial coefficients.", f)
	default:
		return w23P(o, fmt.Sprintf("A two-row tableau of shape (%d,%d) has how many rows?", n-k, k), "2", "Two parts means two rows.", f)
	}
}

func w23BurnsidePQ(o int) generator.Problem {
	ps := w23Primes(2, 13)
	p := ps[rand.Intn(len(ps))]
	q := ps[rand.Intn(len(ps))]
	a := w23N(1, 20)
	b := w23N(1, 20)
	f := map[string]string{}
	w23I(f, "p", p)
	w23I(f, "q", q)
	w23I(f, "a", a)
	w23I(f, "b", b)
	switch o {
	case 1:
		return w23P(o, fmt.Sprintf("A group of order p^%d q^%d is solvable. Enter 1 for yes and 0 for no.", a, b), "1", "Burnside's p^q theorem makes every such group solvable.", f)
	case 2:
		return w23P(o, fmt.Sprintf("How many distinct primes divide the order of a group of order p^%d q^%d?", a, b), "2", "Two prime powers appear.", f)
	default:
		return w23P(o, fmt.Sprintf("A group has order p^%d q^%d with p=%d and q=%d. What is its order?", a, b, p, q), strconv.Itoa(w23Pow(p, a)*w23Pow(q, b)), "Multiply the two prime powers.", f)
	}
}
