package abstract

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("abstract.group.def", &groupDefGen{})
	reg.Register("abstract.group.examples", &groupExamplesGen{})
	reg.Register("abstract.group.subgroup", &subgroupGen{})
	reg.Register("abstract.rings.def", &ringGen{})
	reg.Register("abstract.group.homomorphism", &homomorphismGen{})
	reg.Register("abstract.structures.field", &fieldGen{})
	reg.Register("abstract.structures.module", &moduleGen{})
	reg.Register("abstract.group.cyclic", &cyclicGen{})
	reg.Register("abstract.group.order", &orderGen{})
	reg.Register("abstract.rings.ideal", &idealGen{})
	reg.Register("abstract.rings.quotient", &quotientGen{})
	reg.Register("abstract.group.normal", &normalSubgroupGen{})
	reg.Register("abstract.group.lagrange", &lagrangeGen{})
	reg.Register("abstract.group.quotient_group", &quotientGroupGen{})
	reg.Register("abstract.rings.polynomial", &polynomialRingGen{})
	reg.Register("abstract.group.action", &groupActionGen{})
	reg.Register("abstract.group.sylow", &sylowGen{})
	reg.Register("abstract.rings.ufd", &ufdGen{})
	reg.Register("abstract.field.extension", &fieldExtensionGen{})
	reg.Register("abstract.structures.galois", &galoisGen{})
}

// ----- 1. group.def -----

type groupDefGen struct{}

func (g *groupDefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set       string
		operation string
		isGroup   string
		reason    string
	}
	table := []entry{
		{"integers", "addition", "yes", "Integers under addition satisfy closure, associativity, identity (0), and inverses (-a)."},
		{"natural numbers", "addition", "no", "Natural numbers have no additive inverses (e.g., there is no n such that 3+n=0 in N)."},
		{"rational numbers", "addition", "yes", "Rationals under addition satisfy closure, associativity, identity (0), and inverses (-a)."},
		{"even integers", "addition", "yes", "Even integers under addition are closed, associative, have identity 0, and each 2n has inverse -2n."},
		{"odd integers", "addition", "no", "Odd integers are not closed under addition (odd+odd=even)."},
		{"non-zero rationals", "multiplication", "yes", "Non-zero rationals under multiplication satisfy all group properties: identity 1, inverses 1/a."},
		{"integers", "multiplication", "no", "Integers under multiplication have no inverses (e.g., 2 has no integer multiplicative inverse)."},
		{"real numbers", "addition", "yes", "Real numbers under addition satisfy closure, associativity, identity (0), and inverses (-a)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is the set of %s under %s a group? (yes/no)", e.set, e.operation),
		Answer:      e.isGroup,
		Explanation: e.reason,
	}
}

// ----- 2. group.examples -----

type groupExamplesGen struct{}

func (g *groupExamplesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set       string
		operation string
		isGroup   string
		reason    string
	}
	table := []entry{
		{"Z₄ (integers mod 4)", "addition mod 4", "yes", "Z₄ under addition mod 4 is a group: closed, associative, identity 0, every a has inverse 4−a."},
		{"{0,1,2,3}", "multiplication mod 4", "no", "2 has no inverse mod 4 (2×1=2, 2×2=0, 2×3=2; never get 1)."},
		{"GL(2,R) (invertible 2x2 matrices)", "multiplication", "yes", "Invertible matrices under multiplication form the general linear group: closed, associative, identity I, inverses exist by definition."},
		{"Z₅ \\ {0}", "multiplication mod 5", "yes", "Non-zero elements of Z₅ under multiplication form a group since 5 is prime."},
		{"Z₅ (integers mod 5)", "addition mod 5", "yes", "Z₅ under addition mod 5 is a group: closed, associative, identity 0, each a has inverse 5−a."},
		{"{0,1,2,3,4,5}", "multiplication mod 6", "no", "2 has no inverse mod 6 (2×3=6≡0, never 1), and not all elements have inverses."},
		{"SL(2,R) (2x2 matrices with det=1)", "multiplication", "yes", "Matrices with determinant 1 under multiplication form a group: closed (det(AB)=det(A)det(B)=1), associative, identity I, inverses have det=1."},
		{"S₃ (permutations of 3 elements)", "composition", "yes", "S₃ under composition is the symmetric group: closed, associative, identity, and every permutation has an inverse."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s under %s a group? (yes/no)", e.set, e.operation),
		Answer:      e.isGroup,
		Explanation: e.reason,
	}
}

// ----- 3. subgroup -----

type subgroupGen struct{}

func (g *subgroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	if rand.Intn(3) == 0 {
		type orderEntry struct {
			subgroup string
			group    string
			order    string
			reason   string
		}
		orderTable := []orderEntry{
			{"{0,3}", "Z₆", "2", "The subgroup {0,3} has 2 elements, so its order is 2."},
			{"{0,2,4}", "Z₆", "3", "The subgroup {0,2,4} has 3 elements, so its order is 3."},
			{"{0}", "Z", "1", "The trivial subgroup {0} has order 1."},
			{"{0,5}", "Z₁₀", "2", "The subgroup {0,5} has 2 elements, so its order is 2."},
			{"{0,3,6,9}", "Z₁₂", "4", "The subgroup {0,3,6,9} has 4 elements, so its order is 4."},
		}
		e := orderTable[rand.Intn(len(orderTable))]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the order of the subgroup %s in %s?", e.subgroup, e.group),
			Answer:      e.order,
			Explanation: e.reason,
		}
	}

	type entry struct {
		set       string
		group     string
		operation string
		isSub     string
		reason    string
	}
	table := []entry{
		{"even integers", "Z", "addition", "yes", "Even integers are closed under addition, contain 0, and contain inverses (-2n ∈ 2Z)."},
		{"odd integers", "Z", "addition", "no", "Odd integers are not closed under addition (odd+odd=even), and do not contain 0."},
		{"{0,2,4}", "Z₆", "addition mod 6", "yes", "{0,2,4} is closed under addition mod 6 (0+2=2, 2+4=0, etc.), has identity 0, and each element is its own inverse."},
		{"{0,3}", "Z₆", "addition mod 6", "yes", "{0,3} is closed: 3+3=0, 3+0=3, 0+0=0. Contains identity 0 and inverses."},
		{"{1,2,3,4}", "Z₅", "multiplication mod 5", "yes", "{1,2,3,4} is closed under multiplication mod 5, has identity 1, and each element has an inverse."},
		{"positive integers", "Z", "addition", "no", "Positive integers do not contain the identity 0 and have no additive inverses."},
		{"{0,1,2}", "Z₆", "addition mod 6", "no", "{0,1,2} is not closed: 2+2=4 ∉ {0,1,2}."},
		{"{0,2,4,6}", "Z₈", "addition mod 8", "yes", "{0,2,4,6} is closed (even+even=even), has identity 0, and each element has an inverse (e.g., 2+6=0)."},
		{"non-zero rationals", "Q", "multiplication", "yes", "Non-zero rationals form a subgroup of Q\\{0} under multiplication: closed, contains 1, and each a/b has inverse b/a."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is the set %s a subgroup of %s under %s? (yes/no)", e.set, e.group, e.operation),
		Answer:      e.isSub,
		Explanation: e.reason,
	}
}

// ----- 4. ring -----

type ringGen struct{}

func (g *ringGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set     string
		isRing  string
		isField string
		reason  string
	}
	table := []entry{
		{"Z (integers)", "yes", "no", "Z is a ring under addition and multiplication. It is not a field because elements other than ±1 lack multiplicative inverses in Z."},
		{"2Z (even integers)", "yes", "no", "2Z is a ring under addition and multiplication. It has no multiplicative identity, and is not a field."},
		{"Q (rationals)", "yes", "yes", "Q is a ring and a field: every non-zero element has a multiplicative inverse (a/b → b/a)."},
		{"R (real numbers)", "yes", "yes", "R is a ring and a field: every non-zero real has a multiplicative inverse."},
		{"Z₆ (integers mod 6)", "yes", "no", "Z₆ is a ring but not a field: 2 and 3 have no multiplicative inverses mod 6."},
		{"Z₅ (integers mod 5)", "yes", "yes", "Z₅ is a ring and a field since 5 is prime; every non-zero element has an inverse mod 5."},
		{"Z[i] (Gaussian integers)", "yes", "no", "Gaussian integers a+bi form a ring but not a field (e.g., 2 has no inverse in Z[i])."},
		{"n×n matrices over R", "yes", "no", "n×n matrices under addition and multiplication form a ring but not a field (matrix multiplication is not commutative)."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Is %s a ring under addition and multiplication? (yes/no)", e.set),
			Answer:      e.isRing,
			Explanation: e.reason,
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s a field under addition and multiplication? (yes/no)", e.set),
		Answer:      e.isField,
		Explanation: e.reason,
	}
}

// ----- 5. homomorphism -----

type homomorphismGen struct{}

func (g *homomorphismGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		f        string
		domain   string
		codomain string
		isHomo   string
		reason   string
	}
	table := []entry{
		{"f(x)=2x", "Z", "Z", "yes", "f(a+b)=2(a+b)=2a+2b=f(a)+f(b), so it preserves the group operation."},
		{"f(x)=x+1", "Z", "Z", "no", "f(a+b)=a+b+1, but f(a)+f(b)=(a+1)+(b+1)=a+b+2. These are not equal."},
		{"f(x)=0", "Z", "Z", "yes", "f(a+b)=0 and f(a)+f(b)=0+0=0, so it is the trivial homomorphism."},
		{"f(x)=−x", "Z", "Z", "yes", "f(a+b)=−(a+b)=−a−b=(−a)+(−b)=f(a)+f(b)."},
		{"f(x)=|x|", "Z", "Z", "no", "f(1+(−1))=f(0)=0, but f(1)+f(−1)=1+1=2. These are not equal."},
		{"f(x)=x²", "Z", "Z", "no", "f(a+b)=(a+b)²=a²+2ab+b², but f(a)+f(b)=a²+b². These differ unless 2ab=0."},
		{"f(x)=eˣ", "R", "R⁺", "yes", "f(a+b)=e^{a+b}=e^a·e^b, so it maps the additive group of R to the multiplicative group of R⁺."},
		{"f(x)=ln(x)", "R⁺", "R", "yes", "f(ab)=ln(ab)=ln(a)+ln(b), so it maps the multiplicative group of R⁺ to the additive group of R."},
		{"f(x)=det(x)", "GL(2,R)", "R\\{0}", "yes", "det(AB)=det(A)det(B), so it is a group homomorphism from GL(2,R) to R\\{0}."},
		{"f(x)=x mod 2", "Z", "Z₂", "yes", "f(a+b)=(a+b) mod 2 = (a mod 2)+(b mod 2) mod 2 = f(a)+f(b)."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s from %s to %s a group homomorphism? (yes/no)", e.f, e.domain, e.codomain),
		Answer:      e.isHomo,
		Explanation: e.reason,
	}
}

// ----- 6. field -----

type fieldGen struct{}

func (g *fieldGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set     string
		isField string
		reason  string
	}
	table := []entry{
		{"Q (rational numbers)", "yes", "Q is a field: every non-zero rational a/b has inverse b/a, and all field axioms are satisfied."},
		{"R (real numbers)", "yes", "R is a field: every non-zero real has a multiplicative inverse, satisfying all field axioms."},
		{"C (complex numbers)", "yes", "C is a field: every non-zero complex number a+bi has inverse (a-bi)/(a²+b²)."},
		{"Z (integers)", "no", "Z is not a field: most integers lack multiplicative inverses (e.g., 2 has no inverse in Z)."},
		{"2Z (even integers)", "no", "2Z is not a field: it has no multiplicative identity and lacks inverses."},
		{"Z₅ (integers mod 5)", "yes", "Z₅ is a field because 5 is prime; every non-zero element has a multiplicative inverse mod 5."},
		{"Z₆ (integers mod 6)", "no", "Z₆ is not a field: 2 and 3 are zero divisors with no multiplicative inverses mod 6."},
		{"Z₁₁ (integers mod 11)", "yes", "Z₁₁ is a field because 11 is prime; every non-zero element is invertible mod 11."},
		{"n×n matrices over R (n≥2)", "no", "Matrix multiplication is not commutative and many matrices have no inverse, so it fails field axioms."},
		{"Z[i] (Gaussian integers)", "no", "Gaussian integers form a ring but not a field: 2 has no inverse in Z[i]."},
		{"F₄ (finite field with 4 elements)", "yes", "F₄ exists because 4=2² is a prime power. Every non-zero element has an inverse."},
	}
	e := table[rand.Intn(len(table))]
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    fmt.Sprintf("Is %s a field under addition and multiplication? (yes/no)", e.set),
			Answer:      e.isField,
			Explanation: e.reason,
		}
	}
	// characteristic question
	charTable := []struct {
		field string
		char  string
		reason string
	}{
		{"Q", "0", "Q has characteristic 0: no finite sum of 1s equals 0."},
		{"R", "0", "R has characteristic 0: no finite sum of 1s equals 0."},
		{"C", "0", "C has characteristic 0: no finite sum of 1s equals 0."},
		{"Z₅", "5", "Z₅ has characteristic 5: 1+1+1+1+1 = 5 ≡ 0 mod 5."},
		{"Z₁₁", "11", "Z₁₁ has characteristic 11: 1+1+...+1 (11 times) ≡ 0 mod 11."},
		{"F₄", "2", "F₄ has characteristic 2: 1+1 = 0 in any field of order 2ⁿ."},
	}
	c := charTable[rand.Intn(len(charTable))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is the characteristic of %s?", c.field),
		Answer:      c.char,
		Explanation: c.reason,
	}
}

// ----- 7. module -----

type moduleGen struct{}

func (g *moduleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set     string
		ring    string
		isModule string
		reason  string
	}
	table := []entry{
		{"R³ (3-tuples of reals)", "R", "yes", "R³ is an R-module: scalar multiplication of reals on 3-tuples satisfies all module axioms."},
		{"Z (integers)", "Z", "yes", "Z is a Z-module: integer multiplication of integers satisfies all module axioms."},
		{"any abelian group A", "Z", "yes", "Every abelian group is a Z-module with n·a = a + ... + a (n times)."},
		{"Rⁿ (n-tuples of reals)", "R", "yes", "Rⁿ is a free R-module: componentwise operations satisfy all axioms."},
		{"Z₆ (integers mod 6)", "Z", "yes", "Z₆ is a Z-module with scalar multiplication given by repeated addition mod 6."},
		{"any vector space V over F", "F", "yes", "Every vector space over a field F is an F-module; modules generalise vector spaces."},
		{"2Z (even integers)", "Z", "yes", "2Z is a Z-module: integer multiplication of even integers stays in 2Z."},
		{"R as a vector space over Q", "Q", "yes", "R is a Q-vector space, hence a Q-module, though it is infinite-dimensional over Q."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s a module over %s? (yes/no)", e.set, e.ring),
		Answer:      e.isModule,
		Explanation: e.reason,
	}
}

// ----- 8. cyclic -----

type cyclicGen struct{}

func (g *cyclicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		group   string
		cyclic  string
		reason  string
	}
	table := []entry{
		{"Z_n (integers mod n)", "yes", "Z_n under addition is cyclic, generated by 1."},
		{"Z (integers under addition)", "yes", "Z is cyclic, generated by 1 (or -1)."},
		{"Klein four-group V_4", "no", "V_4 is not cyclic: every non-identity element has order 2, none generates the whole group."},
		{"Z_2 × Z_2", "no", "Direct product of two order-2 cyclic groups is not cyclic (no element of order 4)."},
		{"S_3 (symmetric group on 3)", "no", "S_3 is not cyclic: no element has order 6; max order is 3."},
		{"U(5) = {1,2,3,4} under multiplication mod 5", "yes", "U(5) is cyclic of order 4, generated by 2 (2^1=2, 2^2=4, 2^3=3, 2^4=1)."},
		{"Q under addition", "no", "Q is not cyclic: no single rational generates all rationals."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is the group %s cyclic? (yes/no)", e.group),
		Answer:      e.cyclic,
		Explanation: e.reason,
	}
}

// ----- 9. order -----

type orderGen struct{}

func (g *orderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		elem    string
		group   string
		order   string
		reason  string
	}
	table := []entry{
		{"2", "Z_6 (addition mod 6)", "3", "2+2=4, 2+2+2=0 mod 6, so order of 2 is 3."},
		{"3", "Z_6 (addition mod 6)", "2", "3+3=0 mod 6, so order of 3 is 2."},
		{"1", "Z_6 (addition mod 6)", "6", "1 generates Z_6, so its order is 6."},
		{"2", "Z_4 (addition mod 4)", "2", "2+2=0 mod 4, so order is 2."},
		{"(1 2)", "S_3", "2", "Transposition (1 2) has order 2: applying twice returns identity."},
		{"(1 2 3)", "S_3", "3", "3-cycle (1 2 3) has order 3."},
		{"-1", "Z (addition)", "infinite", "No finite n gives n*(-1)=0 in Z (except 0), so order is infinite."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What is the order of %s in %s?", e.elem, e.group),
		Answer:      e.order,
		Explanation: e.reason,
	}
}

// ----- 10. ideal -----

type idealGen struct{}

func (g *idealGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set     string
		ring    string
		isIdeal string
		reason  string
	}
	table := []entry{
		{"2Z (even integers)", "Z", "yes", "2Z is an ideal of Z: closed under addition and absorbs multiplication by any integer."},
		{"3Z (multiples of 3)", "Z", "yes", "Multiples of 3 form an ideal: adding multiples stays multiple of 3, and integer multiples of multiples of 3 stay multiple of 3."},
		{"odd integers", "Z", "no", "Odd integers are not closed under addition (odd+odd=even) and do not contain 0, so not an ideal."},
		{"nZ (multiples of n)", "Z", "yes", "Multiples of n form an ideal of Z for any integer n."},
		{"{0,2,4}", "Z_6", "yes", "{0,2,4} is an ideal of Z_6: closed under addition mod 6 and absorbs multiplication mod 6."},
		{"{0,1,2}", "Z_6", "no", "{0,1,2} is not an ideal: 1+2=3 outside the set."},
		{"the set of polynomials with zero constant term", "R[x]", "yes", "Polynomials with zero constant term form an ideal: sums and multiples by any polynomial keep zero constant term."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s an ideal of %s? (yes/no)", e.set, e.ring),
		Answer:      e.isIdeal,
		Explanation: e.reason,
	}
}

// ----- 11. quotient -----

type quotientGen struct{}

func (g *quotientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		quotient string
		iso      string
		reason   string
	}
	table := []entry{
		{"Z / 2Z", "Z_2", "Z mod even integers has two cosets: even and odd, isomorphic to integers mod 2."},
		{"Z / 3Z", "Z_3", "Three cosets: 0+3Z, 1+3Z, 2+3Z, forming Z_3."},
		{"Z / 4Z", "Z_4", "Four cosets: multiples of 4 shifted by 0,1,2,3."},
		{"R[x] / (x)", "R", "Quotient by the ideal (x) evaluates polynomials at 0, leaving constants R."},
		{"Z / 6Z", "Z_6", "Quotient of integers by multiples of 6 gives integers mod 6."},
		{"2Z / 4Z", "Z_2", "Even integers mod multiples of 4 has two cosets."},
	}
	e := table[rand.Intn(len(table))]
	return generator.Problem{
		Question:    fmt.Sprintf("What group/ring is %s isomorphic to?", e.quotient),
		Answer:      e.iso,
		Explanation: e.reason,
	}
}

// ----- 12. normal subgroup -----

type normalSubgroupGen struct{}

func (g *normalSubgroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		subgroup string
		group    string
		normal   string
		reason   string
	}
	table := []entry{
		{"A_3 (alternating)", "S_3", "yes", "A_3 has index 2 in S_3, and all index-2 subgroups are normal."},
		{"{e, (1 2)}", "S_3", "no", "Conjugation by (1 3) sends (1 2) to (2 3), which is not in the set, so not normal."},
		{"nZ", "Z (addition)", "yes", "All subgroups of abelian groups are normal; Z is abelian."},
		{"the centre Z(G)", "any group G", "yes", "The centre is always normal: gZ(G)g^{-1}=Z(G)."},
		{"{0,3} in Z_6", "Z_6", "yes", "Z_6 is abelian, so every subgroup is normal."},
		{"{0,2,4} in Z_6", "Z_6", "yes", "Abelian groups have all subgroups normal."},
	}
	// difficulty filters: higher difficulty favors non-abelian examples
	pool := table
	if scale > 3 {
		pool = table[:2] // focus on S_3 tricky cases
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("Is %s normal in %s? (yes/no)", e.subgroup, e.group),
		Answer:      e.normal,
		Explanation: e.reason,
	}
}

// ----- 13. lagrange -----

type lagrangeGen struct{}

func (g *lagrangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type qa struct {
		orderH string
		orderG string
		valid  string
		reason string
	}
	tableSmall := []qa{
		{"2", "6", "yes", "Lagrange: |H| divides |G|, and 2 divides 6."},
		{"3", "6", "yes", "3 divides 6, so a subgroup of order 3 can occur in a group of order 6."},
		{"3", "7", "no", "3 does not divide 7, so no subgroup of order 3 in a group of order 7."},
		{"4", "6", "no", "4 does not divide 6, violating Lagrange."},
	}
	tableLarge := []qa{
		{"12", "24", "yes", "12 divides 24."},
		{"8", "24", "yes", "8 divides 24."},
		{"5", "24", "no", "5 does not divide 24."},
		{"6", "18", "yes", "6 divides 18."},
	}
	pool := tableSmall
	if scale > 3 {
		pool = append(tableSmall, tableLarge...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("Can a group of order %s have a subgroup of order %s? (yes/no)", e.orderG, e.orderH),
		Answer:      e.valid,
		Explanation: e.reason,
	}
}

// ----- 14. quotient group -----

type quotientGroupGen struct{}

func (g *quotientGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		quotient string
		iso      string
		reason   string
	}
	tableEasy := []entry{
		{"Z / 2Z", "Z_2", "Two cosets: even vs odd."},
		{"Z / 3Z", "Z_3", "Three cosets mod 3."},
		{"S_3 / A_3", "Z_2", "Quotient of S_3 by A_3 has order 2, so Z_2."},
	}
	tableHard := []entry{
		{"D_4 / {e, r^2}", "Z_2 x Z_2 (Klein four)", "Dihedral group of order 8 mod its centre of order 2."},
		{"Z_12 / {0,4,8}", "Z_4", "Quotient of Z_12 by order-3 subgroup has order 4."},
		{"Z_8 / {0,4}", "Z_4", "Order 8 mod order 2 gives order 4."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    fmt.Sprintf("What group is %s isomorphic to?", e.quotient),
		Answer:      e.iso,
		Explanation: e.reason,
	}
}

// ----- 15. polynomial ring -----

type polynomialRingGen struct{}

func (g *polynomialRingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	tableEasy := []entry{
		{"Is R[x] an integral domain? (yes/no)", "yes", "Polynomial ring over a field is an integral domain."},
		{"Is Z[x] / (x) isomorphic to Z? (yes/no)", "yes", "Evaluating at 0 kills (x), leaving constants Z."},
		{"Does R[x]/(x^2+1) give a field isomorphic to C? (yes/no)", "yes", "x^2+1 is irreducible over R, so quotient is a field extension C."},
		{"Is Q[x] a PID? (yes/no)", "yes", "Polynomial ring over a field is a PID."},
	}
	tableHard := []entry{
		{"Is R[x]/(x^2) a field? (yes/no)", "no", "x^2 is not irreducible (has repeated root), quotient has zero divisors, so not a field."},
		{"Does Z[x] contain Q? (yes/no)", "no", "Z[x] has integer coefficients only; Q requires rational coefficients."},
		{"Is the ideal (2,x) in Z[x] maximal? (yes/no)", "yes", "Quotient Z[x]/(2,x) ≅ Z_2, a field, so maximal."},
	}
	pool := tableEasy
	if scale > 3 {
		pool = append(tableEasy, tableHard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 16. group action -----

type groupActionGen struct{}

func (g *groupActionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		question string
		answer   string
		reason   string
	}
	easy := []entry{
		{"Does a group action of G on X partition X into orbits? (yes/no)", "yes", "Orbits are equivalence classes under x~y iff gx=y for some g."},
		{"Does orbit-stabilizer give |G| = |Orb(x)|·|Stab(x)|? (yes/no)", "yes", "The bijection G/Stab(x) ≅ Orb(x) gives the formula."},
		{"Is action of S_n on {1..n} transitive? (yes/no)", "yes", "Any i can be sent to any j by a permutation."},
		{"Does trivial action have all orbits of size 1? (yes/no)", "yes", "gx=x for all g gives orbits {x}."},
	}
	hard := []entry{
		{"Does conjugation action of G on itself have orbits = conjugacy classes? (yes/no)", "yes", "gxg^{-1} orbits are conjugacy classes."},
		{"Is Burnside's lemma |X/G| = (1/|G|)∑_g|Fix(g)| true? (yes/no)", "yes", "Counts orbits by averaging fixed points."},
		{"Does action of D_4 on vertices of square have orbit size 4? (yes/no)", "yes", "All 4 vertices are equivalent under dihedral symmetries."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.reason}
}

// ----- 17. sylow -----

type sylowGen struct{}

func (g *sylowGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does a group of order 12=2^2·3 have a Sylow 2-subgroup of order 4? (yes/no)", "yes", "Sylow existence: subgroup of order p^k where p^k||G|."},
		{"Does n_p ≡ 1 mod p for Sylow p-subgroups? (yes/no)", "yes", "Congruence condition: number of Sylow p-subgroups ≡1 mod p."},
		{"Does Sylow guarantee n_p divides |G|? (yes/no)", "yes", "n_p divides |G|/p^k."},
	}
	hard := []entry{
		{"If |G|=12, must n_3 divide 4 and ≡1 mod 3 so n_3∈{1,4}? (yes/no)", "yes", "n_3 divides 4 and ≡1 mod 3."},
		{"Does a group of order 15=3·5 have a normal Sylow 5-subgroup? (yes/no)", "yes", "n_5 divides 3 and ≡1 mod5 → n_5=1."},
		{"If n_p=1, is the Sylow p-subgroup normal? (yes/no)", "yes", "Unique Sylow is characteristic, hence normal."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// ----- 18. UFD -----

type ufdGen struct{}

func (g *ufdGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is Z a UFD? (yes/no)", "yes", "Integers have unique prime factorization."},
		{"Is Z[√-5] a UFD? (yes/no)", "no", "6=2·3=(1+√-5)(1-√-5) gives non-unique factorization."},
		{"Does PID imply UFD? (yes/no)", "yes", "Every PID is a UFD."},
		{"Is Q[x] a UFD? (yes/no)", "yes", "Polynomial ring over a field is a UFD."},
	}
	hard := []entry{
		{"Is Gauss's lemma used to prove Z[x] is a UFD if Z is? (yes/no)", "yes", "Content and primitive polynomials give UFD lift."},
		{"Does UFD require every irreducible is prime? (yes/no)", "yes", "Irreducible ⇒ prime characterizes UFDs among domains."},
		{"Is R[x] with R a UFD itself a UFD? (yes/no)", "yes", "Polynomial ring over a UFD is a UFD."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// ----- 19. field extension -----

type fieldExtensionGen struct{}

func (g *fieldExtensionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Is [Q(√2):Q]=2? (yes/no)", "yes", "Minimal polynomial x^2-2 degree 2."},
		{"Does Q[x]/(x^2+1) ≅ Q(i) have degree 2? (yes/no)", "yes", "Irreducible degree 2 gives degree 2 extension."},
		{"Is splitting field of x^2-2 over Q equal to Q(√2)? (yes/no)", "yes", "Adjoining √2 splits it: (x-√2)(x+√2)."},
	}
	hard := []entry{
		{"Is [Q(∛2):Q]=3? (yes/no)", "yes", "x^3-2 Eisenstein at 2, degree 3."},
		{"Does tower law give [K:Q]=[K:F][F:Q]? (yes/no)", "yes", "Degrees multiply in towers."},
		{"Is Q(√2,√3) of degree 4 over Q? (yes/no)", "yes", "Two independent quadratics give degree 4."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

// ----- 20. Galois -----

type galoisGen struct{}

func (g *galoisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{"Does Galois correspondence match intermediate fields to subgroups of Gal? (yes/no)", "yes", "Order-reversing bijection for Galois extensions."},
		{"Is Gal(Q(√2)/Q) ≅ Z_2? (yes/no)", "yes", "Automorphisms: √2→±√2, order 2."},
		{"Is Q(√2)/Q Galois? (yes/no)", "yes", "Separable and normal (splitting field of x^2-2)."},
	}
	hard := []entry{
		{"Does solvable Galois group correspond to radical solvability? (yes/no)", "yes", "Galois: polynomial solvable by radicals iff Galois group is solvable."},
		{"Is Gal(splitting field of x^3-2) ≅ S_3? (yes/no)", "yes", "Degree 6, non-abelian, is S_3."},
		{"Does fixed field of subgroup H have degree [G:H]? (yes/no)", "yes", "Fundamental theorem: [K^H : k] = |G|/|H|."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}
