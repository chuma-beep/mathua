package abstract

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("abstract.group.def", &groupDefGen{})
	reg.Register("abstract.group.examples", &groupExamplesGen{})
	reg.Register("abstract.subgroup", &subgroupGen{})
	reg.Register("abstract.ring", &ringGen{})
	reg.Register("abstract.homomorphism", &homomorphismGen{})
}

// ----- 1. group.def -----

type groupDefGen struct{}

func (g *groupDefGen) Generate(difficulty float64) generator.Problem {
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

func (g *groupExamplesGen) Generate(difficulty float64) generator.Problem {
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

func (g *subgroupGen) Generate(difficulty float64) generator.Problem {
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

func (g *ringGen) Generate(difficulty float64) generator.Problem {
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

func (g *homomorphismGen) Generate(difficulty float64) generator.Problem {
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
