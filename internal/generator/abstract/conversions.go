package abstract

// The first batch of multiple-choice concepts turned into computed questions.
//
// The domain was recognition-shaped: forty-three of its forty-five concepts
// asked a learner to pick a statement. These ten have an exact answer that the
// generator can compute, an order in a cyclic group, an index, a count of
// generators, a group order, a field order, an orbit size, a kernel size, a
// count of automorphisms and a greatest common divisor, so the concept keeps
// its mathematical content and loses only the list of options.
//
// Every generator publishes the numbers its answer follows from, so the test
// recomputes the answer rather than trusting the arithmetic here, and a
// solution schema can interpolate the same facts.

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func abstractGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func abstractPow(base, exp int) int {
	out := 1
	for i := 0; i < exp; i++ {
		out *= base
	}
	return out
}

func abstractFactorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

// abstractTotient returns Euler's phi of n by counting the residues coprime to
// it. The numbers here are small, so counting is clearer than factoring.
func abstractTotient(n int) int {
	count := 0
	for k := 1; k <= n; k++ {
		if abstractGCD(k, n) == 1 {
			count++
		}
	}
	return count
}

// ---------------------------------------------------------------------------
// Order of an element in a cyclic group
// ---------------------------------------------------------------------------

type orderElementGen struct{}

func (g *orderElementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ns := []int{6, 8, 9, 10, 12, 15, 18, 20, 24, 30}
	n := ns[rand.Intn(len(ns))]
	a := 1 + rand.Intn(n-1)
	order := n / abstractGCD(a, n)
	return generator.Problem{
		Question: fmt.Sprintf("In the additive group of integers modulo %d, the order of the element %d is the smallest positive k with k times %d congruent to 0 modulo %d. What is that order?", n, a, a, n),
		Answer:   fmt.Sprintf("%d", order),
		Explanation: fmt.Sprintf("The multiples of %d modulo %d repeat with period %d divided by the greatest common divisor of %d and %d, which is %d. So the order is %d.",
			a, n, n, a, n, abstractGCD(a, n), order),
		Facts: map[string]string{
			"n": fmt.Sprintf("%d", n), "a": fmt.Sprintf("%d", a), "answer": fmt.Sprintf("%d", order),
		},
	}
}

// ---------------------------------------------------------------------------
// Lagrange's theorem: the index of a subgroup
// ---------------------------------------------------------------------------

type lagrangeIndexGen struct{}

func (g *lagrangeIndexGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ group, subgroup int }
	claims := []claim{
		{group: 8, subgroup: 4}, {group: 12, subgroup: 3}, {group: 12, subgroup: 4},
		{group: 16, subgroup: 8}, {group: 18, subgroup: 6}, {group: 20, subgroup: 5},
		{group: 24, subgroup: 8}, {group: 27, subgroup: 9}, {group: 32, subgroup: 16},
	}
	c := claims[rand.Intn(len(claims))]
	index := c.group / c.subgroup
	return generator.Problem{
		Question: fmt.Sprintf("A group G has order %d and H is a subgroup of order %d. By Lagrange's theorem, what is the index of H in G, that is, the number of distinct left cosets of H?", c.group, c.subgroup),
		Answer:   fmt.Sprintf("%d", index),
		Explanation: fmt.Sprintf("Lagrange's theorem gives the order of G as the product of the order of H and the number of cosets, so the index is %d divided by %d, which is %d.",
			c.group, c.subgroup, index),
		Facts: map[string]string{
			"group_order": fmt.Sprintf("%d", c.group), "subgroup_order": fmt.Sprintf("%d", c.subgroup),
			"answer": fmt.Sprintf("%d", index),
		},
	}
}

// ---------------------------------------------------------------------------
// Generators of a cyclic group
// ---------------------------------------------------------------------------

type cyclicGeneratorsGen struct{}

func (g *cyclicGeneratorsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 5 + rand.Intn(32)
	phi := abstractTotient(n)
	return generator.Problem{
		Question: fmt.Sprintf("The cyclic group of order %d is written additively as the integers modulo %d. How many of its elements generate the whole group?", n, n),
		Answer:   fmt.Sprintf("%d", phi),
		Explanation: fmt.Sprintf("An element generates the cyclic group exactly when it is coprime to %d, and Euler's totient counts those residues: there are %d of them.",
			n, phi),
		Facts: map[string]string{
			"n": fmt.Sprintf("%d", n), "phi": fmt.Sprintf("%d", phi), "answer": fmt.Sprintf("%d", phi),
		},
	}
}

// ---------------------------------------------------------------------------
// Order of the symmetric group
// ---------------------------------------------------------------------------

type symmetricOrderGen struct{}

func (g *symmetricOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	ns := []int{3, 4, 5, 6, 7}
	n := ns[rand.Intn(len(ns))]
	order := abstractFactorial(n)
	return generator.Problem{
		Question: fmt.Sprintf("How many permutations are there of a set with %d elements, that is, what is the order of the symmetric group on %d letters?", n, n),
		Answer:   fmt.Sprintf("%d", order),
		Explanation: fmt.Sprintf("A permutation is an arrangement of the %d letters, and there are %d of them, so the symmetric group has order %d.",
			n, order, order),
		Facts: map[string]string{
			"n": fmt.Sprintf("%d", n), "answer": fmt.Sprintf("%d", order),
		},
	}
}

// ---------------------------------------------------------------------------
// Order of a dihedral group
// ---------------------------------------------------------------------------

type dihedralOrderGen struct{}

func (g *dihedralOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 3 + rand.Intn(22)
	order := 2 * n
	return generator.Problem{
		Question: fmt.Sprintf("The dihedral group of a regular %d-gon consists of its rotations and reflections. What is the order of that group?", n),
		Answer:   fmt.Sprintf("%d", order),
		Explanation: fmt.Sprintf("There are %d rotations, one for each multiple of the central angle, and %d reflections, one through each vertex or edge midpoint, so the order is %d.",
			n, n, order),
		Facts: map[string]string{
			"n": fmt.Sprintf("%d", n), "answer": fmt.Sprintf("%d", order),
		},
	}
}

// ---------------------------------------------------------------------------
// Order of a finite field
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Orbit-stabilizer: the size of an orbit
// ---------------------------------------------------------------------------

type orbitSizeGen struct{}

func (g *orbitSizeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ group, stabilizer int }
	claims := []claim{
		{group: 6, stabilizer: 2}, {group: 8, stabilizer: 4}, {group: 12, stabilizer: 3},
		{group: 12, stabilizer: 4}, {group: 16, stabilizer: 8}, {group: 18, stabilizer: 6},
		{group: 24, stabilizer: 6}, {group: 24, stabilizer: 8},
	}
	c := claims[rand.Intn(len(claims))]
	orbit := c.group / c.stabilizer
	return generator.Problem{
		Question: fmt.Sprintf("A group of order %d acts on a set, and the stabilizer of a point has order %d. By the orbit-stabilizer theorem, how many points lie in that point's orbit?", c.group, c.stabilizer),
		Answer:   fmt.Sprintf("%d", orbit),
		Explanation: fmt.Sprintf("The orbit-stabilizer theorem says the size of the group is the size of the orbit times the size of the stabilizer, so the orbit has %d divided by %d points, which is %d.",
			c.group, c.stabilizer, orbit),
		Facts: map[string]string{
			"group_order": fmt.Sprintf("%d", c.group), "stabilizer_order": fmt.Sprintf("%d", c.stabilizer),
			"answer": fmt.Sprintf("%d", orbit),
		},
	}
}

// ---------------------------------------------------------------------------
// The first isomorphism theorem: the size of a kernel
// ---------------------------------------------------------------------------

type kernelSizeGen struct{}

func (g *kernelSizeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ group, image int }
	claims := []claim{
		{group: 8, image: 2}, {group: 12, image: 3}, {group: 12, image: 4},
		{group: 16, image: 4}, {group: 18, image: 6}, {group: 20, image: 5},
		{group: 24, image: 8}, {group: 27, image: 9},
	}
	c := claims[rand.Intn(len(claims))]
	kernel := c.group / c.image
	return generator.Problem{
		Question: fmt.Sprintf("A homomorphism is defined on a group of order %d and its image has order %d. What is the order of its kernel?", c.group, c.image),
		Answer:   fmt.Sprintf("%d", kernel),
		Explanation: fmt.Sprintf("The first isomorphism theorem identifies the image with the quotient by the kernel, so the order of the group is the order of the kernel times the order of the image. The kernel therefore has %d divided by %d elements, which is %d.",
			c.group, c.image, kernel),
		Facts: map[string]string{
			"group_order": fmt.Sprintf("%d", c.group), "image_order": fmt.Sprintf("%d", c.image),
			"answer": fmt.Sprintf("%d", kernel),
		},
	}
}

// ---------------------------------------------------------------------------
// Automorphisms of a cyclic group
// ---------------------------------------------------------------------------

type automorphismCountGen struct{}

func (g *automorphismCountGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 5 + rand.Intn(32)
	phi := abstractTotient(n)
	return generator.Problem{
		Question: fmt.Sprintf("An automorphism of the cyclic group of order %d sends a generator to another generator. How many automorphisms does that group have?", n),
		Answer:   fmt.Sprintf("%d", phi),
		Explanation: fmt.Sprintf("Choosing the image of one generator determines the automorphism, and the image must also be a generator, so the count is the number of generators, which is Euler's totient of %d, namely %d.",
			n, phi),
		Facts: map[string]string{
			"n": fmt.Sprintf("%d", n), "phi": fmt.Sprintf("%d", phi), "answer": fmt.Sprintf("%d", phi),
		},
	}
}

// ---------------------------------------------------------------------------
// The Euclidean algorithm
// ---------------------------------------------------------------------------

type euclideanGCDGen struct{}

func (g *euclideanGCDGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct{ a, b int }
	claims := []claim{
		{a: 48, b: 18}, {a: 84, b: 30}, {a: 1071, b: 462}, {a: 270, b: 192},
		{a: 252, b: 105}, {a: 144, b: 60}, {a: 391, b: 299}, {a: 1001, b: 273},
	}
	c := claims[rand.Intn(len(claims))]
	divisor := abstractGCD(c.a, c.b)
	return generator.Problem{
		Question: fmt.Sprintf("Use the Euclidean algorithm to find the greatest common divisor of %d and %d.", c.a, c.b),
		Answer:   fmt.Sprintf("%d", divisor),
		Explanation: fmt.Sprintf("Replacing the larger number by its remainder on division by the smaller, and repeating, ends at the greatest common divisor. For %d and %d that divisor is %d.",
			c.a, c.b, divisor),
		Facts: map[string]string{
			"a": fmt.Sprintf("%d", c.a), "b": fmt.Sprintf("%d", c.b), "answer": fmt.Sprintf("%d", divisor),
		},
	}
}

// ---------------------------------------------------------------------------
// Sylow p-subgroups (abstract.group.sylow)
// ---------------------------------------------------------------------------

type sylowSubgroupGen struct{}

func (g *sylowSubgroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(3)

	type sylowCase struct {
		groupOrder int
		p          int
		pk         int
		m          int
		forcedN    int
	}

	cases := []sylowCase{
		{12, 2, 4, 3, 0},
		{12, 3, 3, 4, 0},
		{15, 5, 5, 3, 1},
		{15, 3, 3, 5, 1},
		{20, 2, 4, 5, 0},
		{20, 5, 5, 4, 1},
		{21, 7, 7, 3, 1},
		{24, 2, 8, 3, 0},
		{24, 3, 3, 8, 0},
		{33, 11, 11, 3, 1},
		{35, 7, 7, 5, 1},
		{35, 5, 5, 7, 1},
		{40, 2, 8, 5, 0},
		{40, 5, 5, 8, 1},
		{45, 3, 9, 5, 1},
		{45, 5, 5, 9, 1},
		{54, 3, 27, 2, 0},
		{60, 2, 4, 15, 0},
		{60, 3, 3, 20, 0},
		{60, 5, 5, 12, 0},
		{72, 2, 8, 9, 0},
		{72, 3, 9, 8, 0},
		{77, 11, 11, 7, 1},
		{99, 11, 11, 9, 1},
		{100, 2, 4, 25, 0},
		{100, 5, 25, 4, 1},
	}

	if mode == 1 {
		var forcedCases []sylowCase
		for _, c := range cases {
			if c.forcedN == 1 {
				forcedCases = append(forcedCases, c)
			}
		}
		c := forcedCases[rand.Intn(len(forcedCases))]
		q := fmt.Sprintf("A group G has order %d (%d * %d). By the Sylow theorems, the number n_%d of Sylow %d-subgroups must divide %d and satisfy n_%d ≡ 1 (mod %d). How many Sylow %d-subgroups does G have?",
			c.groupOrder, c.pk, c.m, c.p, c.p, c.m, c.p, c.p, c.p)
		exp := fmt.Sprintf("By the third Sylow theorem, n_%d divides %d and n_%d ≡ 1 (mod %d). The only positive divisor of %d congruent to 1 modulo %d is 1, so G has exactly 1 Sylow %d-subgroup (which is therefore normal).",
			c.p, c.m, c.p, c.p, c.m, c.p, c.p)
		return generator.Problem{
			Question:    q,
			Answer:      "1",
			Explanation: exp,
			Facts: map[string]string{
				"group_order": fmt.Sprintf("%d", c.groupOrder),
				"p":           fmt.Sprintf("%d", c.p),
				"answer":      "1",
			},
		}
	}

	c := cases[rand.Intn(len(cases))]
	if mode == 0 {
		q := fmt.Sprintf("A finite group G has order %d. By the first Sylow theorem, what is the order of a Sylow %d-subgroup of G?", c.groupOrder, c.p)
		exp := fmt.Sprintf("The prime factorization of %d has %d^k = %d with gcd(%d, %d/%d) = 1. By Sylow's theorem, a Sylow %d-subgroup has maximum prime-power order %d.",
			c.groupOrder, c.p, c.pk, c.p, c.groupOrder, c.pk, c.p, c.pk)
		return generator.Problem{
			Question:    q,
			Answer:      fmt.Sprintf("%d", c.pk),
			Explanation: exp,
			Facts: map[string]string{
				"group_order": fmt.Sprintf("%d", c.groupOrder),
				"p":           fmt.Sprintf("%d", c.p),
				"sylow_order": fmt.Sprintf("%d", c.pk),
				"answer":      fmt.Sprintf("%d", c.pk),
			},
		}
	}

	// mode == 2: index [G : P]
	q := fmt.Sprintf("A group G has order %d and P is a Sylow %d-subgroup of order %d. What is the index [G : P] of P in G?", c.groupOrder, c.p, c.pk)
	exp := fmt.Sprintf("By Lagrange's theorem, the index [G : P] is |G| / |P| = %d / %d = %d.", c.groupOrder, c.pk, c.m)
	return generator.Problem{
		Question:    q,
		Answer:      fmt.Sprintf("%d", c.m),
		Explanation: exp,
		Facts: map[string]string{
			"group_order": fmt.Sprintf("%d", c.groupOrder),
			"p":           fmt.Sprintf("%d", c.p),
			"sylow_order": fmt.Sprintf("%d", c.pk),
			"index":       fmt.Sprintf("%d", c.m),
			"answer":      fmt.Sprintf("%d", c.m),
		},
	}
}

// registerConversions registers the concepts of this batch that used to be
// multiple choice. The old registrations are removed from Register, because a
// concept id may be registered only once.
func registerConversions(reg *generator.Registry) {
	reg.Register("abstract.group.order", &orderElementGen{})
	reg.Register("abstract.group.lagrange", &lagrangeIndexGen{})
	reg.Register("abstract.group.cyclic", &cyclicGeneratorsGen{})
	reg.Register("abstract.group.symmetric", &symmetricOrderGen{})
	reg.Register("abstract.group.dihedral", &dihedralOrderGen{})
	reg.Register("abstract.field.finite", &finiteFieldOrderGen{})
	reg.Register("abstract.group.action", &orbitSizeGen{})
	reg.Register("abstract.group.homomorphism", &kernelSizeGen{})
	reg.Register("abstract.group.automorphism", &automorphismCountGen{})
	reg.Register("abstract.rings.euclidean", &euclideanGCDGen{})
	reg.Register("abstract.group.sylow", &sylowSubgroupGen{})
}
