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
	"math/rand"
	"strconv"

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
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "What is the center Z(S_3) of the symmetric group S_3?",
			correct:     "Z(S_3) = {e}",
			distractors: []string{"Z(S_3) = S_3", "Z(S_3) = {e, (1 2)}", "Z(S_3) = {e, (1 2 3)}"},
			reason:      "Only the identity commutes with every permutation; a transposition and a 3-cycle do not commute, so nothing else is central.",
		},
		{
			question:    "What is the center Z(D_4) of the dihedral group D_4 of order 8?",
			correct:     "Z(D_4) = {e, r^2}",
			distractors: []string{"Z(D_4) = {e}", "Z(D_4) = D_4", "Z(D_4) = {e, r, r^2, r^3}"},
			reason:      "The half-turn r^2 commutes with every rotation and every reflection, and no other nonidentity element does.",
		},
		{
			question:    "What is the center Z(Q_8) of the quaternion group Q_8?",
			correct:     "Z(Q_8) = {1, -1}",
			distractors: []string{"Z(Q_8) = {1}", "Z(Q_8) = Q_8", "Z(Q_8) = {1, i, -1, -i}"},
			reason:      "The scalars 1 and -1 commute with every element; i and j anticommute, so no other element is central.",
		},
		{
			question:    "What is the center of the cyclic group Z_8 (addition mod 8)?",
			correct:     "Z(Z_8) = Z_8",
			distractors: []string{"Z(Z_8) = {0}", "Z(Z_8) = {0, 2, 4, 6}", "Z(Z_8) = {0, 4}"},
			reason:      "Z_8 is abelian, so every element commutes with every other and the whole group is its own center.",
		},
		{
			question:    "What is the center of the Klein four-group V_4 = {e, a, b, ab}?",
			correct:     "Z(V_4) = V_4",
			distractors: []string{"Z(V_4) = {e}", "Z(V_4) = {e, a}", "Z(V_4) = {e, a, b}"},
			reason:      "V_4 is abelian, so its center is the whole group.",
		},
		{
			question:    "For any group G, the center Z(G) is which kind of subgroup?",
			correct:     "a normal subgroup of G",
			distractors: []string{"a subgroup that need not be normal", "a quotient of G", "a subgroup of the quotient G/Z(G)"},
			reason:      "Every central element commutes with all of G, so gZ(G)g^-1 = Z(G) for every g and Z(G) is normal.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
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
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "How many factors does a composition series of the cyclic group Z_12 have? (enter a number)",
			correct:     "3",
			distractors: []string{"2", "4", "1"},
			reason:      "The composition factors of Z_12 are Z_2, Z_2 and Z_3, so the series has three factors.",
		},
		{
			question:    "How many factors does a composition series of Z_6 have? (enter a number)",
			correct:     "2",
			distractors: []string{"1", "3", "6"},
			reason:      "6 = 2 * 3, so the composition factors are Z_2 and Z_3: two of them.",
		},
		{
			question:    "What are the composition factors of Z_12, up to order?",
			correct:     "Z_2, Z_2, Z_3",
			distractors: []string{"Z_4, Z_3", "Z_6, Z_2", "Z_12"},
			reason:      "A composition series refines 12 = 2 * 2 * 3 into factors of prime order, giving Z_2, Z_2, Z_3.",
		},
		{
			question:    "The Jordan-Holder theorem says two composition series of the same finite group...",
			correct:     "have the same multiset of composition factors",
			distractors: []string{"have the same length but may have different factors", "end in the same proper normal subgroup", "are the same series"},
			reason:      "Jordan-Holder makes the multiset of composition factors an invariant of the group, even when the series itself is not unique.",
		},
		{
			question:    "The composition factors of a solvable group are...",
			correct:     "cyclic groups of prime order",
			distractors: []string{"nonabelian simple groups", "trivial groups", "fields"},
			reason:      "A solvable group has a composition series with abelian factors, and the simple abelian groups are exactly the cyclic groups of prime order.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 6. Zero divisors and cancellation
// ---------------------------------------------------------------------------

type zeroDivisorGen struct{}

func (g *zeroDivisorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Which nonzero element of Z_4 is a zero divisor?",
			correct:     "2",
			distractors: []string{"1", "3", "0"},
			reason:      "2 * 2 = 0 mod 4, so 2 is a zero divisor; 1 and 3 are units and 0 is excluded by the question.",
		},
		{
			question:    "In Z_6, which product is zero with both factors nonzero?",
			correct:     "2 * 3",
			distractors: []string{"1 * 5", "5 * 5", "2 * 2"},
			reason:      "2 * 3 = 6 = 0 mod 6, and both factors are nonzero, so 2 and 3 are zero divisors.",
		},
		{
			question:    "Which ring has no zero divisors?",
			correct:     "Z_7",
			distractors: []string{"Z_4", "Z_6", "Z_8"},
			reason:      "Z_7 is a field, so it has no zero divisors; each of Z_4, Z_6 and Z_8 has nonzero elements with zero product.",
		},
		{
			question:    "In an integral domain, if a and b are nonzero then ab is...",
			correct:     "never zero",
			distractors: []string{"sometimes zero", "always zero", "always a unit"},
			reason:      "An integral domain has no zero divisors, so a product of two nonzero elements cannot be zero.",
		},
		{
			question:    "The cancellation law (ab = ac with a nonzero implies b = c) holds in which of these rings?",
			correct:     "Z_7",
			distractors: []string{"Z_4", "Z_6", "Z_8"},
			reason:      "Cancellation is equivalent to having no zero divisors, which holds exactly in an integral domain such as Z_7.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 7. Integral domains
// ---------------------------------------------------------------------------

type integralDomainGen struct{}

func (g *integralDomainGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Which of the following is an integral domain?",
			correct:     "Z_7",
			distractors: []string{"Z_6", "Z_8", "Z_10"},
			reason:      "Z_7 is a field, and every field is an integral domain; Z_6, Z_8 and Z_10 all have zero divisors.",
		},
		{
			question:    "Is the ring M_2(R) of 2x2 real matrices an integral domain? (yes/no)",
			correct:     "no",
			distractors: []string{"yes", "only over Q", "only if the matrices commute"},
			reason:      "Take the matrices with a 1 in the top-left corner and in the bottom-right corner; their product is zero, so M_2(R) has zero divisors.",
		},
		{
			question:    "Every finite integral domain is...",
			correct:     "a field",
			distractors: []string{"a UFD but not a field", "a ring with zero divisors", "a group under multiplication"},
			reason:      "In a finite integral domain, multiplication by a nonzero a is injective and hence surjective, so a has an inverse.",
		},
		{
			question:    "Is Z[x] an integral domain? (yes/no)",
			correct:     "yes",
			distractors: []string{"no", "only if x is invertible", "only when the degree is at most one"},
			reason:      "Z is an integral domain, and the polynomial ring over an integral domain is again an integral domain.",
		},
		{
			question:    "Which condition defines an integral domain?",
			correct:     "a commutative ring with identity, no zero divisors, and 1 not equal to 0",
			distractors: []string{"a commutative ring in which every nonzero element is a unit", "a ring with no ideals", "a group under both operations"},
			reason:      "An integral domain is a nonzero commutative ring with identity that has no zero divisors.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 8. Prime and maximal ideals
// ---------------------------------------------------------------------------

type primeIdealGen struct{}

func (g *primeIdealGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "In a commutative ring R with identity, an ideal P is prime exactly when...",
			correct:     "R/P is an integral domain",
			distractors: []string{"R/P is a field", "R/P is a group under multiplication", "P is principal"},
			reason:      "R/P is an integral domain exactly when P is prime, for R commutative with identity.",
		},
		{
			question:    "In a commutative ring R with identity, an ideal M is maximal exactly when...",
			correct:     "R/M is a field",
			distractors: []string{"R/M is an integral domain but not a field", "M is principal", "M contains every prime ideal"},
			reason:      "M is maximal exactly when R/M has no nontrivial proper ideals, which is the same as every nonzero class being invertible.",
		},
		{
			question:    "Which ideal of Z is prime?",
			correct:     "(5)",
			distractors: []string{"(4)", "(6)", "(9)"},
			reason:      "Z/(n) is an integral domain exactly when n is prime, and 5 is prime.",
		},
		{
			question:    "Which ideal of Z is maximal?",
			correct:     "(7)",
			distractors: []string{"(0)", "(4)", "(6)"},
			reason:      "In Z the nonzero prime ideals are maximal; (0) is prime but not maximal because Z is not a field.",
		},
		{
			question:    "Is the ideal (x) prime in Z[x]? (yes/no)",
			correct:     "yes",
			distractors: []string{"no", "only if x is invertible", "only when Z[x] is a field"},
			reason:      "Z[x]/(x) is isomorphic to Z, an integral domain, so (x) is prime; it is not maximal because Z is not a field.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
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

type characteristicGen struct{}

func (g *characteristicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   int
		kind     string
		prime    int
		order    int
		reason   string
	}
	// kind "F_p" means order == prime; "F_pn" means order is a power of the
	// prime larger than it; "char0" means the characteristic is 0. The test
	// recomputes the expected value from kind, prime and order.
	finite := []claim{
		{
			question: "What is the characteristic of the field F_11 (equivalently Z/11Z)? (enter a number)",
			answer:   11, kind: "F_p", prime: 11, order: 11,
			reason: "Adding the identity 1 in F_11 returns 0 after 11 additions and not before, so the characteristic is 11.",
		},
		{
			question: "What is the characteristic of the field F_7 (equivalently Z/7Z)? (enter a number)",
			answer:   7, kind: "F_p", prime: 7, order: 7,
			reason: "The characteristic of F_p is p, so it is 7 for F_7.",
		},
		{
			question: "What is the characteristic of the field F_2 (equivalently Z/2Z)? (enter a number)",
			answer:   2, kind: "F_p", prime: 2, order: 2,
			reason: "In F_2, 1 + 1 = 0, so the characteristic is 2.",
		},
		{
			question: "The field F_9 has 9 = 3^2 elements. What is its characteristic? (enter a number)",
			answer:   3, kind: "F_pn", prime: 3, order: 9,
			reason: "F_9 contains a copy of F_3, and adding 1 three times gives 0; the field order 9 is not the characteristic.",
		},
		{
			question: "The field F_8 has 8 = 2^3 elements. What is its characteristic? (enter a number)",
			answer:   2, kind: "F_pn", prime: 2, order: 8,
			reason: "F_8 has prime subfield F_2, so its characteristic is 2, not the field order 8.",
		},
		{
			question: "The field F_25 has 25 = 5^2 elements. What is its characteristic? (enter a number)",
			answer:   5, kind: "F_pn", prime: 5, order: 25,
			reason: "F_25 has prime subfield F_5, so its characteristic is 5, not 25.",
		},
		{
			question: "The field F_27 has 27 = 3^3 elements. What is its characteristic? (enter a number)",
			answer:   3, kind: "F_pn", prime: 3, order: 27,
			reason: "The characteristic of F_{p^n} is p, so it is 3 here, not the field order 27.",
		},
	}
	zero := []claim{
		{
			question: "What is the characteristic of the field Q? (enter a number)",
			answer:   0, kind: "char0",
			reason: "No positive multiple of 1 equals 0 in Q, so the characteristic is 0 by definition.",
		},
		{
			question: "What is the characteristic of the field R? (enter a number)",
			answer:   0, kind: "char0",
			reason: "No positive multiple of 1 equals 0 in R, so the characteristic is 0.",
		},
		{
			question: "What is the characteristic of the field C? (enter a number)",
			answer:   0, kind: "char0",
			reason: "C contains Q and no positive multiple of 1 vanishes, so its characteristic is 0.",
		},
		{
			question: "What is the characteristic of the field Q(sqrt(2))? (enter a number)",
			answer:   0, kind: "char0",
			reason: "Q(sqrt(2)) contains Q, so it has characteristic 0.",
		},
	}
	claims := finite
	if rand.Intn(4) == 0 {
		claims = zero
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      strconv.Itoa(c.answer),
		Explanation: c.reason,
		Facts: map[string]string{
			"kind": c.kind, "prime": strconv.Itoa(c.prime),
			"order": strconv.Itoa(c.order), "answer": strconv.Itoa(c.answer),
		},
	}
}

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
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A subset N of an R-module M is a submodule exactly when...",
			correct:     "N is closed under addition and under scalar multiplication and contains 0",
			distractors: []string{"N is closed under addition only", "N is a subgroup of M under addition only", "N is closed under scalar multiplication only"},
			reason:      "A submodule is an additive subgroup of M that is also closed under the action of R.",
		},
		{
			question:    "Which is a submodule of the Z-module Z_6?",
			correct:     "{0, 2, 4}",
			distractors: []string{"{0, 1, 2}", "{2, 4, 5}", "{1, 3, 5}"},
			reason:      "{0, 2, 4} is closed under addition mod 6 and under integer scaling; the others are not even subgroups.",
		},
		{
			question:    "The quotient M/N of a module by a submodule N is...",
			correct:     "a module over the same ring",
			distractors: []string{"a ring", "a field", "a group but not a module"},
			reason:      "Addition and scalar multiplication descend to the cosets of N, giving a module over the original ring.",
		},
		{
			question:    "Which is a submodule of R^2 over R?",
			correct:     "the x-axis {(x, 0)}",
			distractors: []string{"the line {(x, 1)}", "the unit circle", "the set {(1, 2)}"},
			reason:      "The x-axis contains 0 and is closed under vector addition and real scaling.",
		},
		{
			question:    "If N is a submodule of M and the quotient map sends m to m + N, its kernel is...",
			correct:     "the submodule N",
			distractors: []string{"the zero submodule {0}", "the whole module M", "the quotient module M/N"},
			reason:      "The coset m + N is the zero coset exactly when m lies in N, so the kernel is N.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 13. Module homomorphisms
// ---------------------------------------------------------------------------

type moduleHomGen struct{}

func (g *moduleHomGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A map f: M -> N of R-modules is a module homomorphism exactly when...",
			correct:     "f(x + y) = f(x) + f(y) and f(rx) = r f(x)",
			distractors: []string{"f(x + y) = f(x) f(y)", "f(rx) = f(r) f(x)", "f(x + y) = f(x) + f(y) only"},
			reason:      "A module map preserves addition and the ring action; the scalar r is not an argument of f.",
		},
		{
			question:    "The kernel of a module homomorphism f: M -> N is...",
			correct:     "a submodule of M",
			distractors: []string{"a submodule of N", "an ideal of R", "a quotient module of N"},
			reason:      "If f(x) = f(y) = 0 then f(x + y) = 0 and f(rx) = r f(x) = 0, so the kernel is a submodule of the domain.",
		},
		{
			question:    "A module homomorphism is injective exactly when...",
			correct:     "its kernel is {0}",
			distractors: []string{"its image is {0}", "it is surjective", "its kernel is the whole module"},
			reason:      "f(x) = f(y) exactly when f(x - y) = 0, so injectivity is equivalent to a trivial kernel.",
		},
		{
			question:    "How many Z-module homomorphisms are there from Z_2 to Z_4? (enter a number)",
			correct:     "2",
			distractors: []string{"1", "4", "0"},
			reason:      "A homomorphism is fixed by the image of 1, which must satisfy 2a = 0 in Z_4; a = 0 and a = 2 both work.",
		},
		{
			question:    "The image of a module homomorphism f: M -> N is...",
			correct:     "a submodule of N",
			distractors: []string{"a submodule of M", "a quotient of N", "an ideal of R"},
			reason:      "The image is closed under addition and scalar multiplication inside the codomain.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}
