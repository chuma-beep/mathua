package abstract

// Counting problems for the abstract algebra concepts that used to ask a
// recognition question. Each pool holds sixteen statements, eight true and
// eight false. Drawing eight and asking for the number that are true gives
// nine possible answers without ever showing the answer in the stem.

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

type aaVerdict struct {
	label string
	yes   bool
	note  string
}

const aaDraw = 8

func countAAVerdicts(lead string, pool []aaVerdict) generator.Problem {
	var yes, no []aaVerdict
	for _, v := range pool {
		if v.yes {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	count := rand.Intn(aaDraw + 1)
	if count > len(yes) {
		count = len(yes)
	}
	if aaDraw-count > len(no) {
		count = aaDraw - len(no)
	}
	rand.Shuffle(len(yes), func(i, j int) { yes[i], yes[j] = yes[j], yes[i] })
	rand.Shuffle(len(no), func(i, j int) { no[i], no[j] = no[j], no[i] })
	chosen := append([]aaVerdict{}, yes[:count]...)
	chosen = append(chosen, no[:aaDraw-count]...)
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

func groupDefProblem() generator.Problem {
	pool := []aaVerdict{
		{"a group is a set with an associative operation, an identity and inverses", true, "these are the three group axioms"},
		{"the identity element of a group is unique", true, "multiplying two identities gives each of them"},
		{"every element of a group has a unique inverse", true, "two inverses of an element must coincide"},
		{"(ab)^-1 = b^-1 a^-1 holds in every group", true, "the reversed product undoes ab"},
		{"a subgroup is closed under the operation and under inverses", true, "closure and inverses are the subgroup test"},
		{"the order of an element divides the order of a finite group", true, "this is Lagrange's theorem"},
		{"every group of prime order is cyclic", true, "a non-identity element generates the whole group"},
		{"the intersection of two subgroups is a subgroup", true, "it contains the identity and is closed"},
		{"every group is abelian", false, "the symmetric group S_3 is not abelian"},
		{"every element of a group has finite order", false, "the infinite cyclic group has elements of infinite order"},
		{"a group may have two distinct identity elements", false, "identities are unique"},
		{"(ab)^-1 = a^-1 b^-1 holds in every group", false, "the order must be reversed"},
		{"the union of two subgroups is always a subgroup", false, "two lines through the origin have a union that is not closed"},
		{"every subset containing the identity is a subgroup", false, "it must also be closed under products and inverses"},
		{"every group of order 6 is cyclic", false, "S_3 has order 6 and is not cyclic"},
		{"cancellation fails in a group", false, "left and right cancellation both hold"},
	}
	return countAAVerdicts("How many of these statements about groups are true?", pool)
}

func groupExamplesProblem() generator.Problem {
	pool := []aaVerdict{
		{"the integers under addition", true, "the identity is 0 and the inverse of n is -n"},
		{"the rational numbers under addition", true, "the identity is 0 and inverses are negatives"},
		{"the non-zero real numbers under multiplication", true, "the identity is 1 and inverses are reciprocals"},
		{"the non-zero rational numbers under multiplication", true, "every non-zero rational has a reciprocal"},
		{"the invertible 2x2 real matrices under multiplication", true, "GL_2(R) is closed under products and inverses"},
		{"Z_5 under addition", true, "the residues modulo 5 form a cyclic group"},
		{"the non-zero residues modulo 7 under multiplication", true, "7 is prime, so every non-zero residue is invertible"},
		{"the symmetries of a square under composition", true, "the dihedral group D_4 has eight elements"},
		{"the natural numbers under addition", false, "positive natural numbers have no additive inverses"},
		{"the integers under multiplication", false, "2 has no integer multiplicative inverse"},
		{"the odd integers under addition", false, "the sum of two odd integers is even"},
		{"all 2x2 real matrices under multiplication", false, "singular matrices have no inverse"},
		{"the non-zero residues modulo 6 under multiplication", false, "2 and 3 multiply to zero modulo 6"},
		{"the positive integers under multiplication", false, "2 has no positive integer reciprocal"},
		{"the real numbers under multiplication", false, "0 has no multiplicative inverse"},
		{"the rational numbers under multiplication", false, "0 has no multiplicative inverse"},
	}
	return countAAVerdicts("How many of these structures are groups under the operation shown?", pool)
}

func ringAxiomsProblem() generator.Problem {
	pool := []aaVerdict{
		{"a ring is an abelian group under addition with associative multiplication distributive over addition", true, "this is the ring definition"},
		{"multiplication distributes over addition on both sides", true, "both distributive laws are required"},
		{"the additive identity of a ring is unique", true, "it is the identity of the underlying abelian group"},
		{"every element of a ring has an additive inverse", true, "the additive structure is a group"},
		{"a commutative ring satisfies xy = yx for all elements", true, "this is what commutative means"},
		{"a ring with unity has a multiplicative identity", true, "the unity is a two-sided identity"},
		{"the integers form a commutative ring with unity", true, "addition and multiplication satisfy the axioms"},
		{"the 2x2 real matrices form a ring with unity", true, "the identity matrix is the multiplicative identity"},
		{"every ring is commutative", false, "matrix rings are not commutative"},
		{"every ring has a multiplicative identity", false, "the even integers form a ring without unity"},
		{"multiplication in a ring must be commutative", false, "only commutative rings impose that"},
		{"a ring may lack additive inverses", false, "addition in a ring is a group operation"},
		{"distributivity is optional in a ring", false, "both distributive laws are axioms"},
		{"every ring is an integral domain", false, "Z_6 has zero divisors"},
		{"the polynomial ring R[x] is a field", false, "x has no multiplicative inverse in R[x]"},
		{"the even integers have a multiplicative identity", false, "1 is odd and does not belong to 2Z"},
	}
	return countAAVerdicts("How many of these statements about rings are true?", pool)
}

func fieldExamplesProblem() generator.Problem {
	pool := []aaVerdict{
		{"the rational numbers Q", true, "every non-zero rational has a reciprocal"},
		{"the real numbers R", true, "every non-zero real has a reciprocal"},
		{"the complex numbers C", true, "every non-zero complex number has an inverse"},
		{"Z_2", true, "2 is prime, so the non-zero residue is invertible"},
		{"Z_3", true, "3 is prime, so every non-zero residue is invertible"},
		{"Z_5", true, "5 is prime, so every non-zero residue is invertible"},
		{"Z_7", true, "7 is prime, so every non-zero residue is invertible"},
		{"Q(sqrt 2)", true, "the inverse of a+b sqrt 2 is (a-b sqrt 2)/(a^2-2b^2)"},
		{"the integers Z", false, "2 has no integer reciprocal"},
		{"the even integers 2Z", false, "there is no multiplicative identity"},
		{"Z_4", false, "2 is a zero divisor modulo 4"},
		{"Z_6", false, "2 and 3 are zero divisors modulo 6"},
		{"Z_8", false, "2 is a zero divisor modulo 8"},
		{"Z_9", false, "3 is a zero divisor modulo 9"},
		{"the 2x2 real matrices M_2(R)", false, "multiplication is not commutative and singular matrices are not invertible"},
		{"the Gaussian integers Z[i]", false, "2 has no inverse in Z[i]"},
	}
	return countAAVerdicts("How many of these structures are fields under the usual operations?", pool)
}

func moduleAxiomsProblem() generator.Problem {
	pool := []aaVerdict{
		{"a module over a ring is an abelian group with a scalar multiplication", true, "the additive group and the action are the module data"},
		{"scalar multiplication satisfies a(x+y) = ax + ay", true, "the action distributes over module addition"},
		{"scalar multiplication satisfies (a+b)x = ax + bx", true, "the action distributes over ring addition"},
		{"scalar multiplication satisfies (ab)x = a(bx)", true, "the action is compatible with ring multiplication"},
		{"a unital module satisfies 1x = x", true, "the ring identity acts as the identity"},
		{"every abelian group is a Z-module", true, "integer multiples define the action"},
		{"every vector space is a module over its field", true, "the vector space axioms are the module axioms"},
		{"a submodule is closed under addition and scalar multiplication", true, "this is the submodule test"},
		{"every module is a vector space", false, "the scalars need not form a field"},
		{"the scalars of a module may come from any set", false, "they must form a ring"},
		{"a module need not be an abelian group under addition", false, "the additive structure is an abelian group"},
		{"scalar multiplication need not distribute over addition", false, "both distributive laws are required"},
		{"every module is free", false, "Z_2 as a Z-module has torsion and no basis"},
		{"every module has a finite basis", false, "many modules are not finitely generated"},
		{"the Z-modules are exactly the finite abelian groups", false, "Z itself is an infinite Z-module"},
		{"every submodule is a direct summand", false, "2Z is a submodule of Z but not a summand"},
	}
	return countAAVerdicts("How many of these statements about modules are true?", pool)
}

func polynomialRingsProblem() generator.Problem {
	pool := []aaVerdict{
		{"Q[x] is a principal ideal domain", true, "the polynomial ring over a field is Euclidean"},
		{"F[x] is a unique factorization domain for every field F", true, "a Euclidean domain is a UFD"},
		{"R[x]/(x^2+1) is isomorphic to C", true, "x^2+1 is irreducible over R"},
		{"Z[x]/(x) is isomorphic to Z", true, "evaluating at zero kills the ideal (x)"},
		{"R[x]/(x^2) has a nonzero nilpotent element", true, "the class of x squares to zero"},
		{"the ideal (2,x) in Z[x] is maximal", true, "the quotient is Z_2, a field"},
		{"deg(fg) = deg f + deg g over an integral domain", true, "the leading coefficients multiply to a nonzero element"},
		{"Q[x]/(x^2-2) is a field", true, "x^2-2 is irreducible over Q"},
		{"Z[x] is a principal ideal domain", false, "(2,x) is not principal"},
		{"R[x]/(x^2) is a field", false, "the class of x is a nonzero zero divisor"},
		{"every quotient of a polynomial ring is an integral domain", false, "R[x]/(x^2) has nilpotents"},
		{"Q[x]/(x^2-2) is isomorphic to Q", false, "the quotient has degree two over Q"},
		{"Z[x] contains Q", false, "the coefficients must be integers"},
		{"the ideal (x) is maximal in Z[x]", false, "Z[x]/(x) is Z, which is not a field"},
		{"a polynomial ring over an infinite field has finitely many irreducible polynomials", false, "there are infinitely many irreducibles"},
		{"the units of F[x] are all the nonzero polynomials", false, "only the nonzero constants are units"},
	}
	return countAAVerdicts("How many of these statements about polynomial rings are true?", pool)
}

func ufdProblem() generator.Problem {
	pool := []aaVerdict{
		{"Z is a unique factorization domain", true, "integers have unique prime factorization"},
		{"every principal ideal domain is a unique factorization domain", true, "a PID satisfies the ascending chain condition and has prime irreducibles"},
		{"if R is a UFD then R[x] is a UFD", true, "Gauss's lemma lifts the factorization"},
		{"Q[x] is a UFD", true, "a polynomial ring over a field is Euclidean"},
		{"in a UFD every irreducible element is prime", true, "this characterizes UFDs among integral domains"},
		{"Z[i] is a UFD", true, "the Gaussian integers are Euclidean"},
		{"k[x,y] is a UFD for a field k", true, "a polynomial ring over a UFD is a UFD"},
		{"every Euclidean domain is a UFD", true, "the Euclidean algorithm gives factorization"},
		{"Z[sqrt(-5)] is a UFD", false, "6 has two distinct factorizations"},
		{"every UFD is a principal ideal domain", false, "k[x,y] is a UFD but not a PID"},
		{"every integral domain is a UFD", false, "Z[sqrt(-5)] is a domain that is not a UFD"},
		{"Z[x] is a principal ideal domain", false, "(2,x) is not principal"},
		{"in a UFD every prime element is reducible", false, "prime elements are irreducible and nonzero"},
		{"Q[x,y] is a principal ideal domain", false, "(x,y) is not principal"},
		{"a UFD must be a field", false, "Z is a UFD and not a field"},
		{"Z[sqrt(-5)] is a Euclidean domain", false, "it is not even a UFD"},
	}
	return countAAVerdicts("How many of these statements about unique factorization are true?", pool)
}

func fieldExtensionsProblem() generator.Problem {
	pool := []aaVerdict{
		{"[C:R] = 2", true, "C is a two-dimensional real vector space"},
		{"[Q(sqrt 2):Q] = 2", true, "sqrt 2 has minimal polynomial x^2-2"},
		{"[Q(sqrt 2, sqrt 3):Q] = 4", true, "the two square roots are independent over Q"},
		{"every finite extension is algebraic", true, "a finite basis gives a polynomial relation"},
		{"if [L:K] = n then L is an n-dimensional K-vector space", true, "the degree is the dimension"},
		{"[F_{p^n}:F_p] = n", true, "the extension has p^n elements"},
		{"Q(sqrt 2) is algebraic over Q", true, "sqrt 2 satisfies x^2-2"},
		{"a splitting field of x^2-2 over Q has degree 2", true, "the roots are plus and minus sqrt 2"},
		{"every algebraic extension is finite", false, "the algebraic closure of Q is infinite over Q"},
		{"[Q(pi):Q] is finite", false, "pi is transcendental"},
		{"[Q(sqrt 2, sqrt 3):Q] = 2", false, "the degree is four"},
		{"every algebraic extension is separable", false, "inseparable extensions exist in positive characteristic"},
		{"R has finite degree over Q", false, "R is uncountable over a countable field"},
		{"Q(sqrt 2) is transcendental over Q", false, "sqrt 2 is algebraic"},
		{"[C:Q] = 2", false, "C has infinite degree over Q"},
		{"every finite extension is normal", false, "Q(2^(1/3)) is not normal over Q"},
	}
	return countAAVerdicts("How many of these statements about field extensions are true?", pool)
}

func galoisTheoryProblem() generator.Problem {
	pool := []aaVerdict{
		{"Gal(Q(sqrt 2)/Q) has order 2", true, "the nontrivial automorphism sends sqrt 2 to minus sqrt 2"},
		{"a finite extension is Galois exactly when it is normal and separable", true, "this is the standard characterization"},
		{"the Galois group of a splitting field acts transitively on the roots", true, "for a separable irreducible polynomial the action is transitive"},
		{"Gal(Q(zeta_n)/Q) is isomorphic to (Z/nZ)^*", true, "the automorphisms are determined by the exponents coprime to n"},
		{"intermediate fields correspond to subgroups of the Galois group", true, "the fundamental theorem of Galois theory"},
		{"Gal(C/R) has order 2", true, "complex conjugation is the nontrivial element"},
		{"a Galois extension with group G has |G| = [L:K]", true, "the order equals the degree"},
		{"the Galois group of x^3-2 over Q has order 6", true, "the splitting field has degree six"},
		{"every finite extension is Galois", false, "Q(2^(1/3)) is not normal"},
		{"Gal(Q(sqrt 2)/Q) has order 4", false, "the degree is two"},
		{"the Galois group of x^3-2 over Q is cyclic of order 3", false, "the group is S_3"},
		{"intermediate fields correspond only to normal subgroups", false, "all subgroups correspond to intermediate fields"},
		{"Gal(Q(zeta_n)/Q) is cyclic for every n", false, "it is (Z/nZ)^*, which need not be cyclic"},
		{"a Galois extension of degree 4 always has group Z/4", false, "the Klein four group also occurs"},
		{"Gal(Q(pi)/Q) is finite", false, "pi is transcendental"},
		{"the Galois group of x^2+1 over R has order 4", false, "the polynomial splits over C with group of order 2"},
	}
	return countAAVerdicts("How many of these statements about Galois theory are true?", pool)
}

func localizationProblem() generator.Problem {
	pool := []aaVerdict{
		{"S^{-1}R consists of fractions r/s with s in S", true, "this is the construction of the localization"},
		{"R_p is a local ring", true, "it has a unique maximal ideal"},
		{"localizing a domain at all nonzero elements gives its field of fractions", true, "every nonzero element becomes invertible"},
		{"Z[1/2] inverts the powers of 2", true, "the allowed denominators are powers of two"},
		{"localization is an exact functor", true, "it preserves short exact sequences"},
		{"a localization of a Noetherian ring is Noetherian", true, "localization preserves the ascending chain condition"},
		{"the maximal ideal of R_p is pR_p", true, "the prime ideal generates the unique maximal ideal"},
		{"localization at a prime ideal has exactly one maximal ideal", true, "this is what makes R_p local"},
		{"localization always produces a field", false, "Z_(p) is local but not a field"},
		{"R_p has no maximal ideals", false, "pR_p is maximal"},
		{"localizing a domain can create zero divisors", false, "a subring of the fraction field stays a domain"},
		{"Z_(p) is the field with p elements", false, "it is a local ring containing Z"},
		{"localization at S = {1} gives the zero ring", false, "it gives R itself"},
		{"every localization map is injective", false, "zero divisors in S can be killed"},
		{"localization at a maximal ideal always gives a field", false, "it gives a local ring, not a field in general"},
		{"R_p is a principal ideal domain for every commutative ring R", false, "localization need not simplify the ideal structure that far"},
	}
	return countAAVerdicts("How many of these statements about localization are true?", pool)
}

func nilpotentGroupsProblem() generator.Problem {
	pool := []aaVerdict{
		{"every finite p-group is nilpotent", true, "its center is nontrivial, so the upper central series grows"},
		{"every nilpotent group is solvable", true, "the derived series reaches the identity"},
		{"a direct product of nilpotent groups is nilpotent", true, "the nilpotency class is the maximum of the factors"},
		{"the upper central series of a nilpotent group reaches the whole group", true, "this is one definition of nilpotency"},
		{"every abelian group is nilpotent of class 1", true, "the center is the whole group"},
		{"the Heisenberg group over a field is nilpotent of class 2", true, "its center is one-dimensional"},
		{"subgroups of nilpotent groups are nilpotent", true, "nilpotency is inherited by subgroups"},
		{"a finite group is nilpotent exactly when it is the product of its Sylow subgroups", true, "this is a standard characterization"},
		{"S_3 is nilpotent", false, "its center is trivial and it is not a p-group"},
		{"every solvable group is nilpotent", false, "S_3 is solvable but not nilpotent"},
		{"nilpotent groups have trivial center", false, "the center of a nontrivial nilpotent group is nontrivial"},
		{"the nilpotency class of an abelian group is 2", false, "it is 1"},
		{"every group of order 12 is nilpotent", false, "A_4 has order 12 and is not nilpotent"},
		{"quotients of nilpotent groups need not be nilpotent", false, "quotients of nilpotent groups are nilpotent"},
		{"D_4 has nilpotency class 3", false, "D_4 has class 2"},
		{"every nilpotent group is abelian", false, "D_4 and the Heisenberg group are nonabelian"},
	}
	return countAAVerdicts("How many of these statements about nilpotent groups are true?", pool)
}

func tensorProductsProblem() generator.Problem {
	pool := []aaVerdict{
		{"dim(V tensor W) = dim V times dim W for finite-dimensional vector spaces", true, "a basis of the tensor product is built from both bases"},
		{"Z/2 tensor Z/3 = 0", true, "the orders are coprime"},
		{"R tensor_R M is isomorphic to M", true, "the ring acts as the identity"},
		{"the tensor product is right exact", true, "it preserves the right end of a short exact sequence"},
		{"Q tensor Z/n = 0", true, "every element of Z/n is torsion and Q is torsion-free"},
		{"C tensor_R C is isomorphic to C x C", true, "the tensor product splits into the two complex structures"},
		{"free modules of ranks m and n have tensor product of rank mn", true, "the basis pairs form a basis"},
		{"Z/m tensor Z/n is isomorphic to Z/gcd(m,n)", true, "the tensor product is the cokernel of multiplication by both orders"},
		{"Z/2 tensor Z/2 = 0", false, "it is isomorphic to Z/2"},
		{"the tensor product is left exact in general", false, "tensoring with Z/2 does not preserve injectivity"},
		{"Q tensor Q = 0", false, "it is isomorphic to Q"},
		{"dim(V tensor W) = dim V + dim W", false, "the dimensions multiply"},
		{"Z tensor Z/2 = 0", false, "it is isomorphic to Z/2"},
		{"the tensor product of two integral domains is an integral domain", false, "C tensor_R C has zero divisors"},
		{"every tensor product is torsion-free", false, "Z/2 tensor Z/2 has torsion"},
		{"if M tensor N = 0 then M = 0 or N = 0", false, "Z/2 tensor Z/3 vanishes with both factors nonzero"},
	}
	return countAAVerdicts("How many of these statements about tensor products are true?", pool)
}

func simpleGroupsProblem() generator.Problem {
	pool := []aaVerdict{
		{"A_n is simple for n >= 5", true, "the alternating group has no proper nontrivial normal subgroup"},
		{"Z/p is simple for every prime p", true, "it has no proper nontrivial subgroup"},
		{"every group of prime order is simple", true, "Lagrange leaves only the trivial subgroups"},
		{"A_5 has order 60", true, "half of 5! is 60"},
		{"S_n is not simple for n >= 3", true, "A_n is a proper nontrivial normal subgroup"},
		{"a simple group has exactly two normal subgroups", true, "the trivial subgroup and the group itself"},
		{"every finite simple group is cyclic, of Lie type, or sporadic", true, "this is the classification of finite simple groups"},
		{"a nonabelian simple group has trivial center", true, "the center is a normal subgroup"},
		{"Z/6 is simple", false, "it has a subgroup of order 2 and one of order 3"},
		{"A_4 is simple", false, "the Klein four group is a normal subgroup"},
		{"every simple group is abelian", false, "A_5 is nonabelian"},
		{"a simple group has no proper nontrivial subgroup", false, "A_5 has subgroups of order 2, 3 and 5"},
		{"S_5 is simple", false, "A_5 is a normal subgroup"},
		{"every group of order 60 is simple", false, "many groups of order 60 are not simple"},
		{"the quaternion group Q_8 is simple", false, "it has the normal subgroup {1,-1}"},
		{"every simple group is finite", false, "PSL_2(Z) contains infinite simple subgroups"},
	}
	return countAAVerdicts("How many of these statements about simple groups are true?", pool)
}

func solvableGroupsProblem() generator.Problem {
	pool := []aaVerdict{
		{"S_3 is solvable", true, "the series 1 < A_3 < S_3 has abelian quotients"},
		{"every finite p-group is solvable", true, "p-groups are nilpotent"},
		{"every nilpotent group is solvable", true, "the derived series reaches the identity"},
		{"every abelian group is solvable", true, "the first derived subgroup is trivial"},
		{"S_n is solvable exactly for n <= 4", true, "S_5 contains the nonsolvable A_5"},
		{"quotients of solvable groups are solvable", true, "the derived series descends to the quotient"},
		{"subgroups of solvable groups are solvable", true, "solvability is inherited by subgroups"},
		{"a polynomial is solvable by radicals exactly when its Galois group is solvable", true, "this is the Abel-Ruffini-Galois theorem in characteristic zero"},
		{"S_5 is solvable", false, "A_5 is a nonsolvable composition factor"},
		{"A_5 is solvable", false, "A_5 is simple and nonabelian"},
		{"every finite group is solvable", false, "A_5 is a counterexample"},
		{"solvable implies nilpotent", false, "S_3 is solvable but not nilpotent"},
		{"every simple group is solvable", false, "A_5 is simple and nonsolvable"},
		{"solvable groups are always abelian", false, "S_3 is solvable and nonabelian"},
		{"A_6 is solvable", false, "A_n is nonsolvable for n >= 5"},
		{"D_4 is not solvable", false, "D_4 is nilpotent, hence solvable"},
	}
	return countAAVerdicts("How many of these statements about solvable groups are true?", pool)
}

func freeGroupsProblem() generator.Problem {
	pool := []aaVerdict{
		{"F_1 is isomorphic to Z", true, "one generator with no relations"},
		{"F_n is infinite for n >= 2", true, "words of increasing length are distinct"},
		{"every group is a quotient of a free group", true, "free groups are free objects"},
		{"subgroups of free groups are free", true, "the Nielsen-Schreier theorem"},
		{"F_2 is nonabelian", true, "the words ab and ba are distinct"},
		{"Z * Z is isomorphic to F_2", true, "the free product of two infinite cyclic groups"},
		{"the rank of F_n is n", true, "the rank is the number of free generators"},
		{"free groups are torsion-free", true, "no nonempty reduced word has finite order"},
		{"F_2 is abelian", false, "ab and ba are different reduced words"},
		{"every free group is finitely generated", false, "a free group on an infinite set is not"},
		{"a subgroup of F_2 has rank at most 2", false, "the commutator subgroup has infinite rank"},
		{"F_n has exactly n elements", false, "it is infinite for n >= 1"},
		{"free groups are finite", false, "F_1 is already infinite"},
		{"every group is free", false, "Z/2 has a relation"},
		{"F_2 has a nontrivial relation", false, "free means no relations"},
		{"F_0 is infinite cyclic", false, "the free group on no generators is trivial"},
	}
	return countAAVerdicts("How many of these statements about free groups are true?", pool)
}

func noetherianRingsProblem() generator.Problem {
	pool := []aaVerdict{
		{"Z is Noetherian", true, "a PID satisfies the ascending chain condition"},
		{"every field is Noetherian", true, "it has only two ideals"},
		{"if R is Noetherian then R[x] is Noetherian", true, "the Hilbert basis theorem"},
		{"quotients of Noetherian rings are Noetherian", true, "ideals of the quotient correspond to ideals above the kernel"},
		{"localizations of Noetherian rings are Noetherian", true, "localization preserves the ascending chain condition"},
		{"every PID is Noetherian", true, "a PID has the ascending chain condition"},
		{"k[x,y] is Noetherian for a field k", true, "it is a quotient of a polynomial ring in two variables"},
		{"a finite direct product of Noetherian rings is Noetherian", true, "ideals decompose componentwise"},
		{"k[x1,x2,...] is Noetherian", false, "the chain (x1) < (x1,x2) < ... never stabilizes"},
		{"every integral domain is Noetherian", false, "k[x1,x2,...] is a non-Noetherian domain"},
		{"every Noetherian ring is a PID", false, "k[x,y] is Noetherian but not a PID"},
		{"Z[x1,x2,...] is Noetherian", false, "the same ascending chain of ideals fails"},
		{"a Noetherian ring has finitely many ideals", false, "Z is Noetherian with infinitely many ideals"},
		{"subrings of Noetherian rings are Noetherian", false, "subrings can fail the ascending chain condition"},
		{"every Noetherian ring is Artinian", false, "Z is Noetherian but not Artinian"},
		{"the zero ring is not Noetherian", false, "it has only one ideal and satisfies the condition"},
	}
	return countAAVerdicts("How many of these statements about Noetherian rings are true?", pool)
}

func galoisGroupsProblem() generator.Problem {
	pool := []aaVerdict{
		{"Gal(Q(sqrt 2)/Q) is isomorphic to Z/2", true, "conjugation is the only nontrivial automorphism"},
		{"Gal(C/R) is isomorphic to Z/2", true, "complex conjugation generates the group"},
		{"Gal(Q(zeta_n)/Q) is isomorphic to (Z/nZ)^*", true, "the exponents coprime to n give the automorphisms"},
		{"the Galois group of x^3-2 over Q is S_3", true, "the splitting field has degree six"},
		{"the Galois group of x^4-2 over Q has order 8", true, "the splitting field has degree eight"},
		{"the Galois group of x^2+1 over R is Z/2", true, "x^2+1 splits over C and conjugation acts"},
		{"Gal(Q(sqrt 2, sqrt 3)/Q) is isomorphic to Z/2 x Z/2", true, "the three quadratic subfields give three involutions"},
		{"a Galois group of a degree n extension has at most n elements", true, "automorphisms are determined by the images of a basis"},
		{"Gal(Q(sqrt 2)/Q) is isomorphic to Z/4", false, "the group has order two"},
		{"the Galois group of x^3-2 over Q is cyclic of order 3", false, "it is S_3"},
		{"Gal(Q(zeta_n)/Q) is cyclic for every n", false, "(Z/8Z)^* is not cyclic"},
		{"Gal(C/R) is trivial", false, "complex conjugation is nontrivial"},
		{"Gal(Q(sqrt 2, sqrt 3)/Q) is isomorphic to Z/4", false, "every nonidentity element has order two"},
		{"the Galois group of x^5-2 over Q is cyclic", false, "it is a semidirect product of order 20"},
		{"every Galois group is abelian", false, "S_3 is a Galois group"},
		{"the Galois group of x^2+1 over Q has order 4", false, "the splitting field Q(i) has degree two"},
	}
	return countAAVerdicts("How many of these statements about Galois groups are true?", pool)
}

func exactSequencesProblem() generator.Problem {
	pool := []aaVerdict{
		{"exactness at B means the image of f equals the kernel of g", true, "this is the definition of exactness"},
		{"0 -> A -> B is exact exactly when A -> B is injective", true, "the kernel at A is zero"},
		{"B -> C -> 0 is exact exactly when B -> C is surjective", true, "the image is all of C"},
		{"a short exact sequence 0 -> A -> B -> C -> 0 gives B/A isomorphic to C", true, "the first isomorphism theorem"},
		{"a split short exact sequence has B isomorphic to A direct sum C", true, "a splitting exhibits the direct summand"},
		{"the alternating sum of dimensions in an exact sequence of vector spaces is zero", true, "the Euler characteristic of a finite exact complex vanishes"},
		{"for a linear map T the sequence 0 -> ker T -> V -> im T -> 0 is exact", true, "the map onto the image is surjective"},
		{"Z -> Z -> Z/2 -> 0 is exact with multiplication by 2", true, "the image is 2Z, the kernel of reduction mod 2"},
		{"exactness at B means the image of f and the kernel of g intersect trivially", false, "they must be equal, not merely disjoint"},
		{"every short exact sequence splits", false, "0 -> Z -> Z -> Z/2 -> 0 does not split"},
		{"a short exact sequence always gives B isomorphic to A direct sum C", false, "this holds only when the sequence splits"},
		{"exact sequences are defined only for vector spaces", false, "they make sense in any abelian category"},
		{"in an exact sequence every map is injective", false, "only the first map of 0 -> A -> B is forced to be injective"},
		{"a five-term exact sequence has all terms isomorphic", false, "the maps need not be isomorphisms"},
		{"exactness is preserved by every functor", false, "only left or right exact functors preserve the relevant half"},
		{"the zero map is never part of an exact sequence", false, "zero maps appear at the ends of every short exact sequence"},
	}
	return countAAVerdicts("How many of these statements about exact sequences are true?", pool)
}

func algebraicElementsProblem() generator.Problem {
	pool := []aaVerdict{
		{"sqrt 2 is algebraic over Q", true, "it satisfies x^2-2"},
		{"pi is transcendental over Q", true, "Lindemann proved pi is not algebraic"},
		{"every element of a finite extension is algebraic", true, "a finite basis gives a linear dependence"},
		{"the set of algebraic numbers is countable", true, "there are countably many polynomials with rational coefficients"},
		{"the algebraic numbers form a field", true, "sums and products of algebraic elements stay algebraic"},
		{"the minimal polynomial of sqrt 2 over Q is x^2-2", true, "x^2-2 is irreducible and vanishes at sqrt 2"},
		{"Q(sqrt 2) has degree 2 over Q", true, "the minimal polynomial has degree two"},
		{"e is transcendental over Q", true, "Hermite proved e is not algebraic"},
		{"every real number is algebraic", false, "pi and e are transcendental"},
		{"sqrt 2 is transcendental over Q", false, "it satisfies x^2-2"},
		{"the algebraic numbers are uncountable", false, "they are countable"},
		{"pi satisfies a nonzero polynomial over Q", false, "pi is transcendental"},
		{"the minimal polynomial of sqrt 2 over Q is x-2", false, "x-2 vanishes at 2, not at sqrt 2"},
		{"the sum of two transcendental numbers is always transcendental", false, "pi + (-pi) = 0 is algebraic"},
		{"every algebraic number is rational", false, "sqrt 2 is algebraic and irrational"},
		{"Q(sqrt 2) contains every algebraic number", false, "it does not contain sqrt 3"},
	}
	return countAAVerdicts("How many of these statements about algebraic elements are true?", pool)
}
