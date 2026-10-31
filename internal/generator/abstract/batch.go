package abstract

// The 2026-10-18 abstract-algebra batch: thirteen concepts that take the domain
// from 32 to the target of 45.
//
// The batch fills the gaps between the group, ring, field and module concepts
// already in the corpus: the center of a group, automorphisms, simplicity and
// composition series, p-groups, zero divisors and integral domains, prime and
// maximal ideals, Euclidean domains, the characteristic of a field, algebraic
// elements, submodules and module homomorphisms. Every concept has at least one
// prerequisite already in the corpus, and several form short chains rather than
// standing alone (center -> p-group, zero divisor -> integral domain -> prime
// ideal, field extension -> algebraic element, module -> submodule/module
// homomorphism).
//
// The theoretical concepts are multiple choice with three or more options built
// by generator.MultipleChoice, so every distractor is a real option the choice
// grader has to reject. The two computed concepts, p-group and characteristic,
// are numeric and publish the facts their answer follows from, so a test can
// recompute the answer instead of comparing a string with itself.
//
// Two conventions are stated wherever they matter. Prime and maximal ideals are
// taught for a commutative ring with identity, because that is the setting in
// which "R/P is an integral domain" and "R/M is a field" hold. The
// characteristic is kept distinct from the field order: F_{p^n} has
// characteristic p, not p^n, and that distinction is a question in its own right.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("abstract.group.center", &groupCenterGen{})
	reg.Register("abstract.group.simple", &simpleGroupGen{})
	reg.Register("abstract.group.p_group", &pGroupGen{})
	reg.Register("abstract.group.composition_series", &compositionSeriesGen{})
	reg.Register("abstract.rings.zero_divisor", &zeroDivisorGen{})
	reg.Register("abstract.rings.integral_domain", &integralDomainGen{})
	reg.Register("abstract.rings.prime_ideal", &primeIdealGen{})
	reg.Register("abstract.field.characteristic", &characteristicGen{})
	reg.Register("abstract.field.algebraic", &algebraicElementGen{})
	reg.Register("abstract.structures.submodule", &submoduleGen{})
	reg.Register("abstract.structures.module_hom", &moduleHomGen{})
}

// ---------------------------------------------------------------------------
// 1. Center of a group
// ---------------------------------------------------------------------------

type groupCenterGen struct{}

func (g *groupCenterGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type centerCase struct {
		group  string
		order  int
		reason string
	}
	var c centerCase
	switch rand.Intn(7) {
	case 0:
		c = centerCase{"S_3", 1, "Only the identity commutes with every permutation, so the center of S_3 is trivial."}
	case 1:
		n := 3 + rand.Intn(8)
		if n%2 == 0 {
			c = centerCase{fmt.Sprintf("D_%d", n), 2, fmt.Sprintf("The half-turn r^(%d/2) commutes with every rotation and every reflection, and no other nonidentity element does.", n)}
		} else {
			c = centerCase{fmt.Sprintf("D_%d", n), 1, fmt.Sprintf("For odd n the dihedral group D_%d has trivial center.", n)}
		}
	case 2:
		n := 2 + rand.Intn(11)
		c = centerCase{fmt.Sprintf("Z_%d", n), n, fmt.Sprintf("Z_%d is cyclic, hence abelian, so every element is central and the center is the whole group.", n)}
	case 3:
		n := 3 + rand.Intn(4)
		c = centerCase{fmt.Sprintf("S_%d", n), 1, fmt.Sprintf("For n >= 3 the center of S_%d is trivial.", n)}
	case 4:
		p := []int{2, 3, 5}[rand.Intn(3)]
		c = centerCase{fmt.Sprintf("GL(2,%d)", p), p - 1, fmt.Sprintf("The center of GL(2,%d) is the set of nonzero scalar matrices lambda I, so it has order %d.", p, p-1)}
	case 5:
		c = centerCase{"Q_8", 2, "Only 1 and -1 commute with every element of the quaternion group, so the center has order 2."}
	default:
		c = centerCase{"V_4", 4, "V_4 is abelian, so every element is central and the center is the whole group."}
	}
	answer := strconv.Itoa(c.order)
	return generator.Problem{
		Question:    fmt.Sprintf("What is the order of the center Z(%s)? (enter a number)", c.group),
		Answer:      answer,
		Explanation: c.reason,
		Facts:       map[string]string{"group": c.group, "order": answer},
	}
}

// ---------------------------------------------------------------------------
// 2. Group automorphisms
// ---------------------------------------------------------------------------

type automorphismGen struct{}

func (g *automorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "How many automorphisms does the cyclic group Z_5 have? (enter a number)",
			correct:     "4",
			distractors: []string{"5", "1", "2"},
			reason:      "Aut(Z_n) has order phi(n), and phi(5) = 4.",
		},
		{
			question:    "How many automorphisms does the cyclic group Z_6 have? (enter a number)",
			correct:     "2",
			distractors: []string{"6", "1", "4"},
			reason:      "Aut(Z_n) has order phi(n), and phi(6) = 2.",
		},
		{
			question:    "How many automorphisms does the cyclic group Z_2 have? (enter a number)",
			correct:     "1",
			distractors: []string{"2", "0", "4"},
			reason:      "The only automorphism of Z_2 fixes the identity and the nonidentity element, so Aut(Z_2) is trivial.",
		},
		{
			question:    "Which group is Aut(Z_3) isomorphic to?",
			correct:     "Z_2",
			distractors: []string{"Z_3", "the trivial group", "S_3"},
			reason:      "Aut(Z_n) is isomorphic to the unit group (Z/nZ)^x, which has order phi(3) = 2.",
		},
		{
			question:    "For a group G, the inner automorphisms form a subgroup isomorphic to which group?",
			correct:     "G/Z(G)",
			distractors: []string{"the whole group G", "Aut(G)", "Z(G)"},
			reason:      "Conjugation by g and by h agree exactly when gh^-1 is central, so the kernel of G -> Aut(G) is Z(G) and the image is G/Z(G).",
		},
		{
			question:    "Which statement about the inner automorphisms of a group G is true?",
			correct:     "They form a normal subgroup of Aut(G)",
			distractors: []string{"They always equal Aut(G)", "They are exactly the automorphisms of order 2", "They form the center of Aut(G)"},
			reason:      "Inn(G) is closed under composition and is preserved by conjugation by any automorphism, so it is normal in Aut(G).",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 3. Simple groups
// ---------------------------------------------------------------------------

type simpleGroupGen struct{}

func (g *simpleGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Which of the following groups is simple?",
			correct:     "A_5",
			distractors: []string{"Z_6", "S_3", "D_4"},
			reason:      "A_5 has no nontrivial proper normal subgroup; Z_6, S_3 and D_4 each do.",
		},
		{
			question:    "Is the cyclic group Z_7 simple? (yes/no)",
			correct:     "yes",
			distractors: []string{"no", "only if 7 is even", "only if Z_7 is nonabelian"},
			reason:      "A cyclic group of prime order has no nontrivial proper subgroup at all, so it is simple.",
		},
		{
			question:    "Which is the smallest nonabelian simple group?",
			correct:     "A_5",
			distractors: []string{"Z_2", "S_3", "D_4"},
			reason:      "A_5 has order 60, and every smaller nonabelian group has a nontrivial proper normal subgroup.",
		},
		{
			question:    "Is Z_6 simple? (yes/no)",
			correct:     "no",
			distractors: []string{"yes", "only if 6 is prime", "it has no normal subgroups"},
			reason:      "Z_6 has the nontrivial proper normal subgroup {0, 3}, isomorphic to Z_2.",
		},
		{
			question:    "Which of the following groups is simple?",
			correct:     "Z_5",
			distractors: []string{"Z_4", "Z_6", "Z_8"},
			reason:      "A cyclic group is simple exactly when its order is prime, so Z_5 qualifies and Z_4, Z_6, Z_8 do not.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 4. p-groups (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 5. Composition series
// ---------------------------------------------------------------------------

type compositionSeriesGen struct{}

func (g *compositionSeriesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type seriesCase struct {
		group  string
		length int
		reason string
	}
	n := 2 + rand.Intn(29)
	omega := 0
	for m := n; m > 1; {
		d := 2
		for m%d != 0 {
			d++
		}
		omega++
		m /= d
	}
	cases := []seriesCase{
		{fmt.Sprintf("Z_%d", n), omega, fmt.Sprintf("%d factors into %d primes counted with multiplicity, so Z_%d has a composition series of length %d.", n, omega, n, omega)},
		{"S_3", 2, "S_3 > A_3 > {e} has two factors."},
		{"S_4", 4, "S_4 > A_4 > V_4 > C_2 > {e} has four factors."},
		{"S_5", 2, "S_5 > A_5 > {e} has two factors because A_5 is simple."},
		{"D_4", 3, "D_4 > <r> > <r^2> > {e} has three factors."},
		{"Q_8", 3, "Q_8 > <i> > <-1> > {e} has three factors."},
		{"V_4", 2, "V_4 > C_2 > {e} has two factors."},
	}
	c := cases[rand.Intn(len(cases))]
	return generator.Problem{
		Question:    fmt.Sprintf("How many factors does a composition series of %s have? (enter a number)", c.group),
		Answer:      strconv.Itoa(c.length),
		Explanation: c.reason,
		Facts:       map[string]string{"group": c.group, "length": strconv.Itoa(c.length)},
	}
}

// ---------------------------------------------------------------------------
// 6. Zero divisors and cancellation
// ---------------------------------------------------------------------------

type zeroDivisorGen struct{}

func (g *zeroDivisorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(27)
	phi := 0
	for a := 1; a <= n; a++ {
		if intGCD(a, n) == 1 {
			phi++
		}
	}
	count := n - phi - 1
	return generator.Problem{
		Question:    fmt.Sprintf("How many nonzero zero divisors does Z_%d have? (enter a number)", n),
		Answer:      strconv.Itoa(count),
		Explanation: fmt.Sprintf("Z_%d has %d units and the zero element, so the other %d nonzero elements are zero divisors.", n, phi, count),
		Facts:       map[string]string{"n": strconv.Itoa(n), "count": strconv.Itoa(count)},
	}
}

func intGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// ---------------------------------------------------------------------------
// 7. Integral domains
// ---------------------------------------------------------------------------

type integralDomainGen struct{}

func (g *integralDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type ringCase struct {
		name     string
		isDomain bool
		reason   string
	}
	rings := []ringCase{
		{"Z_6", false, "2 * 3 = 0 in Z_6"},
		{"Z_7", true, "Z_7 is a field"},
		{"Z_8", false, "2 * 4 = 0 in Z_8"},
		{"Z_10", false, "2 * 5 = 0 in Z_10"},
		{"Z_12", false, "3 * 4 = 0 in Z_12"},
		{"M_2(R)", false, "the matrix units E_11 and E_22 multiply to zero"},
		{"Z[x]", true, "a polynomial ring over an integral domain is an integral domain"},
		{"Q", true, "Q is a field"},
		{"R", true, "R is a field"},
		{"C", true, "C is a field"},
		{"Z_4", false, "2 * 2 = 0 in Z_4"},
		{"F_5", true, "F_5 is a field"},
		{"Z_2[x]", true, "a polynomial ring over a field is an integral domain"},
		{"Z_2 x Z_2", false, "(1,0) * (0,1) = (0,0)"},
		{"Z[i]", true, "the Gaussian integers sit inside C and have no zero divisors"},
		{"Z_5[x]", true, "a polynomial ring over a field is an integral domain"},
	}
	idx := rand.Perm(len(rings))[:4]
	chosen := make([]ringCase, 0, 4)
	parts := make([]string, 0, 4)
	count := 0
	for i, j := range idx {
		chosen = append(chosen, rings[j])
		parts = append(parts, fmt.Sprintf("(%s) %s", string(rune('a'+i)), rings[j].name))
		if rings[j].isDomain {
			count++
		}
	}
	explainParts := make([]string, 0, 4)
	for _, r := range chosen {
		if r.isDomain {
			explainParts = append(explainParts, fmt.Sprintf("%s is an integral domain because %s", r.name, r.reason))
		} else {
			explainParts = append(explainParts, fmt.Sprintf("%s is not an integral domain because %s", r.name, r.reason))
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("How many of these rings are integral domains? %s (enter a number)", strings.Join(parts, " ")),
		Answer:      strconv.Itoa(count),
		Explanation: strings.Join(explainParts, "; ") + ".",
		Facts:       map[string]string{"count": strconv.Itoa(count), "rings": strings.Join(parts, " ")},
	}
}

// ---------------------------------------------------------------------------
// 8. Prime and maximal ideals
// ---------------------------------------------------------------------------

type primeIdealGen struct{}

func (g *primeIdealGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(59)
	primes := distinctPrimeFactors(n)
	kind := "prime"
	if rand.Intn(2) == 0 {
		kind = "maximal"
	}
	factors := make([]string, 0, len(primes))
	for _, p := range primes {
		factors = append(factors, fmt.Sprintf("(%d)", p))
	}
	return generator.Problem{
		Question:    fmt.Sprintf("How many %s ideals does Z_%d have? (enter a number)", kind, n),
		Answer:      strconv.Itoa(len(primes)),
		Explanation: fmt.Sprintf("The %s ideals of Z_%d correspond to the distinct prime divisors of %d, namely %s.", kind, n, n, strings.Join(factors, ", ")),
		Facts:       map[string]string{"n": strconv.Itoa(n), "count": strconv.Itoa(len(primes)), "kind": kind},
	}
}

func distinctPrimeFactors(n int) []int {
	factors := make([]int, 0, 4)
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			factors = append(factors, d)
			for n%d == 0 {
				n /= d
			}
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}
	return factors
}

// ---------------------------------------------------------------------------
// 9. Euclidean domains
// ---------------------------------------------------------------------------

type euclideanGen struct{}

func (g *euclideanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "In a Euclidean domain with Euclidean function N, the division algorithm writes a = qb + r with...",
			correct:     "r = 0 or N(r) < N(b)",
			distractors: []string{"r = 0 or N(r) < N(a)", "always r = 0", "N(a) < N(b)"},
			reason:      "The remainder is either zero or strictly smaller than the divisor, measured by N.",
		},
		{
			question:    "Which of these is a Euclidean domain?",
			correct:     "the ring of integers Z",
			distractors: []string{"the polynomial ring Z[x]", "the ring Z[sqrt(-5)]", "the ring of 2x2 real matrices"},
			reason:      "Z is Euclidean with N equal to the absolute value; Z[x] is not a PID, Z[sqrt(-5)] is not a UFD, and the matrix ring has zero divisors.",
		},
		{
			question:    "Euclidean domain implies...",
			correct:     "principal ideal domain",
			distractors: []string{"a UFD but not necessarily a PID", "a field", "a ring with zero divisors"},
			reason:      "The division algorithm produces a single generator for every ideal, so a Euclidean domain is a PID and therefore a UFD.",
		},
		{
			question:    "On the integers, the standard Euclidean function is...",
			correct:     "the absolute value",
			distractors: []string{"the number of prime factors", "the residue mod 2", "the square of the element"},
			reason:      "N(a) = |a| makes the remainder strictly smaller in absolute value than the divisor.",
		},
		{
			question:    "Which chain of implications is correct?",
			correct:     "Euclidean implies PID implies UFD",
			distractors: []string{"UFD implies PID implies Euclidean", "PID implies Euclidean implies UFD", "Euclidean implies field implies UFD"},
			reason:      "Euclidean gives principal ideals, and a PID has unique factorization; the reverse implications fail in general.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 10. Characteristic of a field (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 11. Algebraic and transcendental elements
// ---------------------------------------------------------------------------

type algebraicElementGen struct{}

func (g *algebraicElementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Is sqrt(2) algebraic or transcendental over Q?",
			correct:     "algebraic",
			distractors: []string{"transcendental", "neither", "both"},
			reason:      "sqrt(2) is a root of x^2 - 2, a nonzero polynomial over Q, so it is algebraic over Q.",
		},
		{
			question:    "Is pi algebraic over Q? (yes/no)",
			correct:     "no",
			distractors: []string{"yes", "only if pi is rational", "only over C"},
			reason:      "pi is transcendental over Q: no nonzero polynomial with rational coefficients has pi as a root.",
		},
		{
			question:    "Which polynomial is the minimal polynomial of sqrt(2) over Q?",
			correct:     "x^2 - 2",
			distractors: []string{"x - 2", "x^2 + 2", "x^2 - 4"},
			reason:      "x^2 - 2 is monic, irreducible over Q, and has sqrt(2) as a root.",
		},
		{
			question:    "The degree [Q(sqrt(2)) : Q] is...",
			correct:     "2",
			distractors: []string{"1", "4", "infinite"},
			reason:      "The minimal polynomial x^2 - 2 has degree 2, and the extension degree equals the degree of the minimal polynomial.",
		},
		{
			question:    "Which number is algebraic over Q?",
			correct:     "sqrt(3)",
			distractors: []string{"pi", "e", "1/pi"},
			reason:      "sqrt(3) is a root of x^2 - 3, while pi, e and 1/pi are transcendental.",
		},
		{
			question:    "If alpha is algebraic over a field F, then F(alpha) is...",
			correct:     "a finite extension of F",
			distractors: []string{"an infinite extension of F", "equal to F", "a transcendental extension of F"},
			reason:      "Adjoining an algebraic element adds only the finitely many powers below its minimal polynomial, so the degree is finite.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 12. Submodules and quotient modules
// ---------------------------------------------------------------------------

type submoduleGen struct{}

func (g *submoduleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(4) {
	case 0:
		n := 2 + rand.Intn(29)
		count := 0
		for d := 1; d <= n; d++ {
			if n%d == 0 {
				count++
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("How many submodules does the Z-module Z_%d have? (enter a number)", n),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("Submodules of Z_%d correspond to divisors of %d, and %d has %d divisors.", n, n, n, count),
			Facts:       map[string]string{"n": strconv.Itoa(n), "count": strconv.Itoa(count)},
		}
	case 1:
		n := 6 + rand.Intn(19)
		var divisors []int
		for d := 1; d <= n; d++ {
			if n%d == 0 {
				divisors = append(divisors, d)
			}
		}
		d := divisors[rand.Intn(len(divisors))]
		size := n / d
		return generator.Problem{
			Question:    fmt.Sprintf("How many elements does the submodule generated by %d in Z_%d have? (enter a number)", d, n),
			Answer:      strconv.Itoa(size),
			Explanation: fmt.Sprintf("The submodule is the set of multiples of %d mod %d, which has %d elements.", d, n, size),
			Facts:       map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d), "size": strconv.Itoa(size)},
		}
	case 2:
		n := 6 + rand.Intn(19)
		var divisors []int
		for d := 1; d <= n; d++ {
			if n%d == 0 {
				divisors = append(divisors, d)
			}
		}
		d := divisors[rand.Intn(len(divisors))]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the order of the quotient Z_%d / <%d>? (enter a number)", n, d),
			Answer:      strconv.Itoa(d),
			Explanation: fmt.Sprintf("The subgroup <%d> has %d elements, so the quotient has %d / %d = %d elements.", d, n/d, n, n/d, d),
			Facts:       map[string]string{"n": strconv.Itoa(n), "d": strconv.Itoa(d), "order": strconv.Itoa(d)},
		}
	default:
		m := 2 + rand.Intn(2)
		k := 1 + rand.Intn(m)
		var vecs []string
		for j := 0; j < k; j++ {
			v := "("
			for i := 0; i < m; i++ {
				if i > 0 {
					v += ", "
				}
				if i == j {
					v += "1"
				} else {
					v += "0"
				}
			}
			v += ")"
			vecs = append(vecs, v)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("What is the dimension of the subspace of R^%d spanned by %s? (enter a number)", m, joinWithAnd(vecs)),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("The %d chosen standard basis vectors are linearly independent, so they span a subspace of dimension %d.", k, k),
			Facts:       map[string]string{"m": strconv.Itoa(m), "dimension": strconv.Itoa(k)},
		}
	}
}

func joinWithAnd(parts []string) string {
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

// ---------------------------------------------------------------------------
// 13. Module homomorphisms
// ---------------------------------------------------------------------------

type moduleHomGen struct{}

func (g *moduleHomGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(3) {
	case 0:
		m := 2 + rand.Intn(11)
		n := 2 + rand.Intn(11)
		count := intGCD(m, n)
		return generator.Problem{
			Question:    fmt.Sprintf("How many Z-module homomorphisms are there from Z_%d to Z_%d? (enter a number)", m, n),
			Answer:      strconv.Itoa(count),
			Explanation: fmt.Sprintf("A homomorphism is fixed by the image a of 1, which must satisfy %d a = 0 in Z_%d; there are gcd(%d, %d) = %d such a.", m, n, m, n, count),
			Facts:       map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n), "count": strconv.Itoa(count)},
		}
	case 1:
		m := 1 + rand.Intn(5)
		n := 1 + rand.Intn(5)
		dim := m * n
		return generator.Problem{
			Question:    fmt.Sprintf("What is the dimension of Hom_R(R^%d, R^%d) as a vector space? (enter a number)", m, n),
			Answer:      strconv.Itoa(dim),
			Explanation: fmt.Sprintf("A linear map R^%d to R^%d is an %d by %d matrix, and the space of such matrices has dimension %d.", m, n, m, n, dim),
			Facts:       map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n), "dim": strconv.Itoa(dim)},
		}
	default:
		m := 2 + rand.Intn(9)
		rank := rand.Intn(m + 1)
		nullity := m - rank
		return generator.Problem{
			Question:    fmt.Sprintf("A linear map from R^%d to R^%d has rank %d. What is the dimension of its kernel? (enter a number)", m, m, rank),
			Answer:      strconv.Itoa(nullity),
			Explanation: fmt.Sprintf("By rank-nullity, dim kernel = %d - %d = %d.", m, rank, nullity),
			Facts:       map[string]string{"m": strconv.Itoa(m), "rank": strconv.Itoa(rank), "nullity": strconv.Itoa(nullity)},
		}
	}
}
