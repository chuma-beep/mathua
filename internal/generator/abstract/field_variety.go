package abstract

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The field and structure generators below share one defect: every row of their
// pools answered "yes", and the computed branch returned a bare answer with no
// options even though the corpus grades these concepts multiple_choice. Both
// make the question unanswerable in the client. Each pool now carries genuine
// false propositions, and every computed shape is offered as multiple choice.

// ----- characteristic -----

type characteristicGen struct{}

// The test recomputes the expected characteristic from kind, prime and order,
// so those facts must stay published even though the answer is a single number.
func (g *characteristicGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19}
	p := primes[rand.Intn(len(primes))]
	switch rand.Intn(3) {
	case 0:
		return generator.Problem{
			Question:    fmt.Sprintf("What is the characteristic of the field F_%d (equivalently Z/%dZ)? (enter a number)", p, p),
			Answer:      strconv.Itoa(p),
			Explanation: fmt.Sprintf("Adding the identity 1 in F_%d returns 0 after %d additions and not before, so the characteristic is %d.", p, p, p),
			Facts: map[string]string{
				"kind": "F_p", "prime": strconv.Itoa(p), "order": strconv.Itoa(p), "answer": strconv.Itoa(p),
			},
		}
	case 1:
		n := 2 + rand.Intn(2)
		order := abstractPow(p, n)
		return generator.Problem{
			Question:    fmt.Sprintf("The field F_%d has %d = %d^%d elements. What is its characteristic? (enter a number)", order, order, p, n),
			Answer:      strconv.Itoa(p),
			Explanation: fmt.Sprintf("F_%d contains a copy of F_%d, so adding 1 exactly %d times gives 0. The field order %d is not the characteristic.", order, p, p, order),
			Facts: map[string]string{
				"kind": "F_pn", "prime": strconv.Itoa(p), "order": strconv.Itoa(order), "answer": strconv.Itoa(p),
			},
		}
	default:
		fields := []struct{ name, reason string }{
			{"Q", "No positive multiple of 1 equals 0 in the rationals."},
			{"R", "No positive multiple of 1 equals 0 in the reals."},
			{"C", "The complex numbers contain Q and no positive multiple of 1 vanishes."},
			{"Q(√2)", "Q(√2) contains Q, whose characteristic is 0."},
			{"Q(√3)", "Q(√3) contains Q, whose characteristic is 0."},
			{"Q(√5)", "Q(√5) contains Q, whose characteristic is 0."},
			{"Q(√7)", "Q(√7) contains Q, whose characteristic is 0."},
			{"Q(√2,√3)", "A composite of characteristic-0 extensions still contains Q."},
		}
		f := fields[rand.Intn(len(fields))]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the characteristic of the field %s? (enter a number)", f.name),
			Answer:      "0",
			Explanation: f.reason,
			Facts: map[string]string{
				"kind": "char0", "prime": "0", "order": "0", "answer": "0",
			},
		}
	}
}

// ----- finite field order -----

type finiteFieldOrderGen struct{}

func (g *finiteFieldOrderGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	primes := []int{2, 3, 5, 7, 11, 13}
	p := primes[rand.Intn(len(primes))]
	n := 1 + rand.Intn(4)
	order := abstractPow(p, n)
	return generator.Problem{
		Question: fmt.Sprintf("Every finite field has p to the n elements for a prime p and a positive integer n. For p = %d and n = %d, how many elements does the field have?", p, n),
		Answer:   strconv.Itoa(order),
		Explanation: fmt.Sprintf("The field is a vector space of dimension %d over its prime field with %d elements, so it has %d to the %d elements, which is %d.",
			n, p, p, n, order),
		Facts: map[string]string{
			"p": strconv.Itoa(p), "n": strconv.Itoa(n), "answer": strconv.Itoa(order),
		},
	}
}

// ----- field extension -----

type fieldExtensionGen struct{}

func (g *fieldExtensionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	easy := []entry{
		{"Is [Q(√2):Q]=2? (yes/no)", "yes", "The minimal polynomial x^2-2 is irreducible over Q, so the degree is 2."},
		{"Is [Q(∛2):Q]=2? (yes/no)", "no", "The minimal polynomial of ∛2 is x^3-2, so the extension degree is 3, not 2."},
		{"Does Q[x]/(x^2+1) ≅ Q(i)? (yes/no)", "yes", "x^2+1 is irreducible over Q, so the quotient is the field Q(i)."},
		{"Is the splitting field of x^2-2 over Q equal to Q(√2)? (yes/no)", "yes", "Adjoining √2 splits it into (x-√2)(x+√2)."},
		{"Is Q(∛2)/Q a Galois extension? (yes/no)", "no", "It is not normal: the other two cube roots of 2 are missing."},
	}
	hard := []entry{
		{"Does the tower law give [K:Q]=[K:F][F:Q]? (yes/no)", "yes", "Degrees multiply in a tower of finite extensions."},
		{"Is Q(√2,√3) of degree 4 over Q? (yes/no)", "yes", "Two independent quadratic constraints give degree 4."},
		{"Is the splitting field of x^3-2 over Q of degree 3? (yes/no)", "no", "It also needs a primitive cube root of unity, so the degree is 6."},
		{"Is [Q(√2,√2):Q]=4? (yes/no)", "no", "Adjoining √2 twice adds nothing the second time, so the degree is 2."},
	}
	pool := easy
	if int(1+ctx.Difficulty*4) > 3 {
		pool = append(easy, hard...)
	}
	switch rand.Intn(3) {
	case 0:
		primes := []int{2, 3, 5, 7}
		p := primes[rand.Intn(len(primes))]
		return groupChoice(
			fmt.Sprintf("Is x^4-%d irreducible over Q, so that [Q(∜%d):Q]=4?", p, p),
			"yes",
			fmt.Sprintf("x^4-%d is Eisenstein at %d, so it is irreducible and the degree is 4.", p, p),
			"no", "only if p is odd", "only when p is 2",
		)
	case 1:
		ps := []int{3, 5, 7, 11}
		p := ps[rand.Intn(len(ps))]
		deg := p - 1
		return groupChoice(
			fmt.Sprintf("The cyclotomic field Q(ζ_%d) adjoins a primitive %d-th root of unity. What is [Q(ζ_%d):Q]?", p, p, p),
			strconv.Itoa(deg),
			fmt.Sprintf("The Galois group is (Z/%dZ)^×, which has φ(%d) = %d elements, so the degree is %d.", p, p, deg, deg),
			strconv.Itoa(p), strconv.Itoa(p+1), strconv.Itoa(2*p),
		)
	default:
		m := 2 + rand.Intn(3)
		n := 2 + rand.Intn(3)
		return groupChoice(
			fmt.Sprintf("F is an extension of Q with [F:Q]=%d, and K is an extension of F with [K:F]=%d. What is [K:Q]?", m, n),
			strconv.Itoa(m*n),
			fmt.Sprintf("The tower law multiplies the degrees: %d × %d = %d.", m, n, m*n),
			strconv.Itoa(m+n), strconv.Itoa(m*n+1), strconv.Itoa(m+n+1),
		)
	}
	e := pool[rand.Intn(len(pool))]
	return groupYesNo(e.q, e.a, e.e)
}

// ----- Galois theory -----

type galoisGen struct{}

func (g *galoisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	easy := []entry{
		{"Does the Galois correspondence match intermediate fields to subgroups of the Galois group? (yes/no)", "yes", "It is an order-reversing bijection for a Galois extension."},
		{"Is Gal(Q(√2)/Q) ≅ Z_2? (yes/no)", "yes", "The only nontrivial automorphism sends √2 to -√2, and doing it twice is the identity."},
		{"Is Gal(Q(√2)/Q) ≅ Z_4? (yes/no)", "no", "The group has order 2, and Z_4 has order 4."},
		{"Is Q(√2)/Q Galois? (yes/no)", "yes", "It is the splitting field of the separable polynomial x^2-2."},
	}
	hard := []entry{
		{"Does a solvable Galois group correspond to solvability by radicals? (yes/no)", "yes", "A polynomial is solvable by radicals exactly when its Galois group is solvable."},
		{"Is the Galois group of the splitting field of x^3-2 isomorphic to S_3? (yes/no)", "yes", "The splitting field has degree 6 and the group is non-abelian, hence S_3."},
		{"Does the fixed field of a subgroup H have degree [G:H] over the base? (yes/no)", "yes", "The fundamental theorem gives [K^H:k] = |G|/|H|."},
		{"Is every finite extension of Q Galois? (yes/no)", "no", "Q(∛2)/Q has degree 3 but only one Q-automorphism, so it is not normal."},
	}
	pool := easy
	if int(1+ctx.Difficulty*4) > 3 {
		pool = append(easy, hard...)
	}
	switch rand.Intn(3) {
	case 0:
		return groupChoice(
			"Q(√2,√3) is the splitting field of (x^2-2)(x^2-3). Which group is Gal(Q(√2,√3)/Q)?",
			"V_4",
			"Each of the two square roots may be sent to its negative independently, so the group is Z_2 × Z_2 = V_4.",
			"Z_4", "S_3", "Z_2",
		)
	case 1:
		gorder := 2 + rand.Intn(11)
		horder := 1 + rand.Intn(gorder-1)
		for gorder%horder != 0 {
			horder = 1 + rand.Intn(gorder-1)
		}
		return groupChoice(
			fmt.Sprintf("A Galois extension has group G of order %d and a subgroup H of order %d. What is [K^H:Q]?", gorder, horder),
			strconv.Itoa(gorder/horder),
			fmt.Sprintf("The fundamental theorem gives [K^H:Q] = |G|/|H| = %d/%d = %d.", gorder, horder, gorder/horder),
			strconv.Itoa(gorder-horder), strconv.Itoa(gorder*horder), strconv.Itoa(horder),
		)
	default:
		ns := []int{3, 4, 5}
		n := ns[rand.Intn(len(ns))]
		fact := 1
		for i := 2; i <= n; i++ {
			fact *= i
		}
		return groupChoice(
			fmt.Sprintf("What is the order of the symmetric group S_%d, which is the Galois group of a generic degree-%d polynomial?", n, n),
			strconv.Itoa(fact),
			fmt.Sprintf("S_%d has %d! = %d elements.", n, n, fact),
			strconv.Itoa(n*n), strconv.Itoa(2*n), strconv.Itoa(fact+1),
		)
	}
	e := pool[rand.Intn(len(pool))]
	return groupYesNo(e.q, e.a, e.e)
}

// ----- Galois group of a specific polynomial -----

type galoisGroupGen struct{}

func (g *galoisGroupGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	easy := []entry{
		{"Is the Galois group of x^2-2 over Q isomorphic to Z_2? (yes/no)", "yes", "Q(√2)/Q has degree 2, so its group has order 2."},
		{"Is the Galois group of the cyclotomic field Q(ζ_p) isomorphic to (Z/pZ)^×? (yes/no)", "yes", "Each automorphism is ζ_p ↦ ζ_p^k for k coprime to p."},
		{"Is the Galois group of x^3-2 over Q of order 6? (yes/no)", "yes", "The splitting field has degree 6, and the group order equals the degree."},
		{"Is the Galois group of x^3-2 over Q abelian? (yes/no)", "no", "It is S_3, which is not abelian."},
	}
	hard := []entry{
		{"Is the Galois group of a generic quintic S_5? (yes/no)", "yes", "A generic quintic has full symmetric group by Hilbert irreducibility, and S_5 is not solvable."},
		{"Does the discriminant detect containment in A_n? (yes/no)", "yes", "The group lies in A_n exactly when the discriminant is a square."},
		{"Is the Galois group of every finite extension of a finite field cyclic? (yes/no)", "yes", "Gal(F_{p^n}/F_p) is cyclic, generated by the Frobenius map."},
		{"Is Gal(Q(ζ_8)/Q) cyclic? (yes/no)", "no", "(Z/8Z)^× ≅ Z_2 × Z_2, which is not cyclic."},
	}
	pool := easy
	if int(1+ctx.Difficulty*4) > 3 {
		pool = append(easy, hard...)
	}
	switch rand.Intn(3) {
	case 0:
		primes := []int{2, 3, 5, 7, 11}
		p := primes[rand.Intn(len(primes))]
		n := 2 + rand.Intn(3)
		return groupChoice(
			fmt.Sprintf("Gal(F_%d/F_%d) is generated by the Frobenius map. What is its order?", abstractPow(p, n), p),
			strconv.Itoa(n),
			fmt.Sprintf("The extension F_%d over F_%d has degree %d, and the group order equals the degree.", abstractPow(p, n), p, n),
			strconv.Itoa(n+1), strconv.Itoa(abstractPow(p, n)-p), strconv.Itoa(p),
		)
	case 1:
		ns := []int{5, 7, 8, 9, 10, 12}
		n := ns[rand.Intn(len(ns))]
		deg := abstractTotient(n)
		return groupChoice(
			fmt.Sprintf("What is [Q(ζ_%d):Q]?", n),
			strconv.Itoa(deg),
			fmt.Sprintf("The Galois group is (Z/%dZ)^× of order φ(%d) = %d, so the degree is %d.", n, n, deg, deg),
			strconv.Itoa(n), strconv.Itoa(n-1), strconv.Itoa(deg+1),
		)
	default:
		ns := []int{3, 4, 5}
		n := ns[rand.Intn(len(ns))]
		fact := 1
		for i := 2; i <= n; i++ {
			fact *= i
		}
		return groupChoice(
			fmt.Sprintf("A generic polynomial of degree %d has Galois group S_%d. What is its order?", n, n),
			strconv.Itoa(fact),
			fmt.Sprintf("|S_%d| = %d! = %d.", n, n, fact),
			strconv.Itoa(n*n), strconv.Itoa(2*n), strconv.Itoa(fact-1),
		)
	}
	e := pool[rand.Intn(len(pool))]
	return groupYesNo(e.q, e.a, e.e)
}

// ----- exact sequences -----

type exactSequenceGen struct{}

func (g *exactSequenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	easy := []entry{
		{"Is an exact sequence one where the image equals the kernel at every step? (yes/no)", "yes", "That equality is the definition of exactness."},
		{"Does a short exact sequence 0→A→B→C→0 imply |B|=|A|·|C| for finite groups? (yes/no)", "yes", "C is isomorphic to the quotient B/A, so the orders multiply."},
		{"Is 0→Z --2→ Z → Z/2Z →0 exact? (yes/no)", "yes", "The kernel of reduction is 2Z, which is the image of multiplication by 2."},
		{"Does every short exact sequence split? (yes/no)", "no", "0→Z --2→ Z → Z/2Z →0 does not split, because Z has no element of order 2."},
	}
	hard := []entry{
		{"Is the tensor product right exact but not left exact? (yes/no)", "yes", "Tensoring preserves cokernels, not kernels."},
		{"Is the five lemma a statement about exact sequences and isomorphisms? (yes/no)", "yes", "It is proved by a diagram chase."},
		{"Does a split short exact sequence give B ≅ A ⊕ C? (yes/no)", "yes", "A splitting map exhibits B as the direct sum."},
		{"Does the tensor product preserve injectivity? (yes/no)", "no", "Z --2→ Z is injective, but tensoring with Z/2 gives the zero map on Z/2."},
	}
	pool := easy
	if int(1+ctx.Difficulty*4) > 3 {
		pool = append(easy, hard...)
	}
	if rand.Intn(3) == 0 {
		a := 2 + rand.Intn(8)
		c := 2 + rand.Intn(8)
		return groupChoice(
			fmt.Sprintf("0→A→B→C→0 is a short exact sequence of finite groups with |A|=%d and |C|=%d. What is |B|?", a, c),
			strconv.Itoa(a*c),
			fmt.Sprintf("Exactness identifies C with B/A, so |B| = |A|·|C| = %d·%d = %d.", a, c, a*c),
			strconv.Itoa(a+c), strconv.Itoa(a*c+1), strconv.Itoa(a*c-1),
		)
	}
	e := pool[rand.Intn(len(pool))]
	return groupYesNo(e.q, e.a, e.e)
}

// ----- modules -----

type moduleGen struct{}

func (g *moduleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		set    string
		ring   string
		isMod  string
		reason string
	}
	table := []entry{
		{"R³ (3-tuples of reals)", "R", "yes", "R³ is an R-module: scalar multiplication of reals on 3-tuples satisfies every module axiom."},
		{"Z (integers)", "Z", "yes", "Z is a Z-module, since integer multiplication of integers satisfies every module axiom."},
		{"any abelian group A", "Z", "yes", "Every abelian group is a Z-module with n·a = a + ... + a taken n times."},
		{"Rⁿ (n-tuples of reals)", "R", "yes", "Rⁿ is a free R-module, with componentwise addition and scaling."},
		{"Z₆ (integers mod 6)", "Z", "yes", "Z₆ is a Z-module under repeated addition mod 6."},
		{"2Z (even integers)", "Z", "yes", "Integer multiples of even integers stay even, so 2Z is a Z-module."},
		{"Z (integers)", "Q", "no", "Scalar multiplication by 1/2 would send 1 to 1/2, which is not an integer, so Z is not a Q-module."},
		{"Q (rationals), with the obvious Z-scaling, as a finitely generated Z-module", "Z", "no", "Finitely many rationals span only a subgroup generated by them, never all of Q."},
		{"Z/2Z, asking whether every Z-module is free", "Z", "no", "Z/2Z has torsion, so it has no basis and is not free."},
	}
	if rand.Intn(3) == 0 {
		e := table[rand.Intn(len(table))]
		return groupYesNo(fmt.Sprintf("Is %s a module over %s? (yes/no)", e.set, e.ring), e.isMod, e.reason)
	}
	switch rand.Intn(2) {
	case 0:
		n := 2 + rand.Intn(4)
		return groupChoice(
			fmt.Sprintf("What is the dimension of R^%d as an R-module?", n),
			strconv.Itoa(n),
			fmt.Sprintf("R^%d is a free R-module of rank %d.", n, n),
			strconv.Itoa(n-1), strconv.Itoa(n+1), strconv.Itoa(n*n),
		)
	default:
		m := 5 + rand.Intn(4)
		k := 2 + rand.Intn(m-2)
		j := 2 + rand.Intn(m-2)
		return groupChoice(
			fmt.Sprintf("In Z_%d viewed as a Z-module, compute %d·%d. (enter a number)", m, k, j),
			strconv.Itoa(k*j%m),
			fmt.Sprintf("%d·%d = %d, and reducing modulo %d leaves %d.", k, j, k*j, m, k*j%m),
			strconv.Itoa(k+j), strconv.Itoa(k*j), strconv.Itoa(k*j%m+1),
		)
	}
}

// ----- tensor products -----

type tensorGen struct{}

func (g *tensorGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		q, a, e string
	}
	easy := []entry{
		{"Is the tensor product M⊗_R N characterised by a universal bilinear map? (yes/no)", "yes", "Every bilinear map factors uniquely through the tensor product."},
		{"Does Z/mZ ⊗_Z Z/nZ ≅ Z/gcd(m,n)Z? (yes/no)", "yes", "The tensor of two cyclic groups is cyclic of order the gcd."},
		{"Is M⊗_R R ≅ M? (yes/no)", "yes", "R is the identity for the tensor product."},
		{"Is Z/2Z ⊗_Z Z/3Z nonzero? (yes/no)", "no", "The gcd of 2 and 3 is 1, so the tensor is Z/1Z, the zero group."},
	}
	hard := []entry{
		{"Does the tensor product preserve exactness on the right? (yes/no)", "yes", "It is right exact, which is why Tor measures the failure on the left."},
		{"Is Q⊗_Z Z/nZ = 0? (yes/no)", "yes", "Q is divisible, so it kills torsion."},
		{"Does localization satisfy S^{-1}R ⊗_R M ≅ S^{-1}M? (yes/no)", "yes", "Localization is base change along R → S^{-1}R."},
		{"Is the tensor product of free modules of ranks m and n free of rank m+n? (yes/no)", "no", "The ranks multiply, so the rank is mn."},
	}
	pool := easy
	if int(1+ctx.Difficulty*4) > 3 {
		pool = append(easy, hard...)
	}
	if rand.Intn(3) == 0 {
		m := 3 + rand.Intn(10)
		n := 3 + rand.Intn(10)
		v := fieldGCD(m, n)
		return groupChoice(
			fmt.Sprintf("Z/%dZ ⊗_Z Z/%dZ ≅ Z/gcd(%d,%d)Z. How many elements does the tensor have?", m, n, m, n),
			strconv.Itoa(v),
			fmt.Sprintf("The gcd of %d and %d is %d, so the tensor is cyclic of order %d.", m, n, v, v),
			strconv.Itoa(m+n), strconv.Itoa(m*n), strconv.Itoa(v+1),
		)
	}
	e := pool[rand.Intn(len(pool))]
	return groupYesNo(e.q, e.a, e.e)
}

// fieldGCD is the greatest common divisor used for the order of a tensor
// product of cyclic groups.
func fieldGCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
