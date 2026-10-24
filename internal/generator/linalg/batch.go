package linalg

// The 2026-10-18 linear-algebra batch: eleven concepts that take the domain
// from 34 to the target of 45.
//
// The survey found the domain well supplied with computations but thin on the
// structural ideas that connect them: vector-space axioms, kernel and image,
// invertibility, change of basis, orthogonal complements, projection onto a
// subspace, the spectral theorem, quadratic definiteness, Cayley-Hamilton and
// the condition number. Several obvious additions were already taught, so they
// were not repeated: determinant properties, transpose, linear combinations and
// coordinate representation, the characteristic polynomial, and algebraic and
// geometric multiplicity all live in existing lessons. Gram-Schmidt is the
// method behind the existing QR lesson and was also left out.
//
// Five concepts are numeric, one is a tuple (the cross product), and five are
// multiple choice. Every numeric generator publishes the facts its answer
// follows from, so the tests recompute the answer instead of comparing a string
// with itself; the tuple generator publishes both factors so the test can
// recompute the cross product.

import (
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("linalg.space.axioms", &spaceAxiomsGen{})
	reg.Register("linalg.transform.kernel", &kernelGen{})
	reg.Register("linalg.transform.invertible", &invertibleMapGen{})
	reg.Register("linalg.basis.change", &basisChangeGen{})
	reg.Register("linalg.orth.complement", &complementGen{})
	reg.Register("linalg.orth.projection_matrix", &projectionMatrixGen{})
	reg.Register("linalg.symmetric.spectral", &spectralGen{})
	reg.Register("linalg.quadratic.definite", &definiteGen{})
	reg.Register("linalg.vector.cross", &crossGen{})
	reg.Register("linalg.matrix.condition", &conditionGen{})
	reg.Register("linalg.eigen.cayley_hamilton", &cayleyHamiltonGen{})
}

// ---------------------------------------------------------------------------
// 1. Vector spaces and their axioms
// ---------------------------------------------------------------------------

type spaceAxiomsGen struct{}

func (g *spaceAxiomsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Which axiom guarantees that every vector v has an additive inverse?",
			correct:     "for every v there is a w with v + w = 0",
			distractors: []string{"closure under scalar multiplication", "distributivity of scalar multiplication over addition", "commutativity of scalar multiplication"},
			reason:      "The additive-inverse axiom supplies the negative of each vector, which is what makes subtraction possible.",
		},
		{
			question:    "Which operation is required of a real vector space but not of a group?",
			correct:     "multiplication by every real scalar",
			distractors: []string{"associative addition", "an additive identity", "additive inverses"},
			reason:      "A vector space is an abelian group under addition with an extra action of the field by scalar multiplication.",
		},
		{
			question:    "The zero vector in a vector space is...",
			correct:     "unique",
			distractors: []string{"one of many", "the number zero in every space", "not required to exist"},
			reason:      "If 0 and 0' are both additive identities then 0 = 0 + 0' = 0', so the zero vector is unique.",
		},
		{
			question:    "Which operation must distribute over vector addition?",
			correct:     "scalar multiplication",
			distractors: []string{"the dot product", "the determinant", "the transpose"},
			reason:      "The axiom c(u + v) = cu + cv is one of the two distributivity laws that tie the field action to addition.",
		},
		{
			question:    "A vector space over R must be closed under...",
			correct:     "addition and multiplication by every real scalar",
			distractors: []string{"addition only", "multiplication by positive scalars only", "the cross product"},
			reason:      "Closure under both operations is the first requirement; the remaining axioms describe how they behave.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 2. Kernel and image (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 3. Injective, surjective and invertible maps
// ---------------------------------------------------------------------------

type invertibleMapGen struct{}

func (g *invertibleMapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A linear map T from R^n to R^n is invertible exactly when...",
			correct:     "its kernel is the zero subspace",
			distractors: []string{"its image is the zero subspace", "it is injective but not surjective", "its matrix has a zero eigenvalue"},
			reason:      "For a square map, a trivial kernel makes the map injective, and finite dimension then makes it surjective and invertible.",
		},
		{
			question:    "For a square matrix A, which condition is equivalent to invertibility?",
			correct:     "the determinant is nonzero",
			distractors: []string{"the determinant is zero", "A has a zero eigenvalue", "the columns are linearly dependent"},
			reason:      "The determinant is the product of the eigenvalues, so it is nonzero exactly when no eigenvalue is zero.",
		},
		{
			question:    "Can a linear map from R^3 to R^2 be injective?",
			correct:     "no, because the rank is at most 2 and the nullity is at least 1",
			distractors: []string{"yes, if it is surjective", "yes, if its matrix has rank 3", "only when it is the zero map"},
			reason:      "Rank-nullity gives nullity 3 - rank, and the rank cannot exceed the dimension 2 of the codomain.",
		},
		{
			question:    "For a square matrix, injectivity of the linear map implies...",
			correct:     "surjectivity",
			distractors: []string{"nothing about surjectivity", "that the matrix is singular", "that the determinant is zero"},
			reason:      "A square matrix with trivial kernel has full rank, so its columns span the whole codomain.",
		},
		{
			question:    "The inverse of an invertible linear map is...",
			correct:     "linear",
			distractors: []string{"not necessarily linear", "always nonlinear", "the transpose"},
			reason:      "The inverse of an invertible linear map is itself linear, which is why it is represented by the inverse matrix.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 4. Coordinates and change of basis (numeric)
// ---------------------------------------------------------------------------

type basisChangeGen struct{}

func (g *basisChangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		b1, b2   [2]int
		v        [2]int
		ask      int
		reason   string
	}
	claims := []claim{
		{
			question: "In R^2, the vector (3, 5) is written as a(1, 1) + b(1, -1). What is a? (enter a number)",
			answer:   "4", b1: [2]int{1, 1}, b2: [2]int{1, -1}, v: [2]int{3, 5}, ask: 0,
			reason: "The equations are a + b = 3 and a - b = 5, whose solution is a = 4 and b = -1.",
		},
		{
			question: "In R^2, the vector (4, 9) is written as a(2, 0) + b(0, 3). What is b? (enter a number)",
			answer:   "3", b1: [2]int{2, 0}, b2: [2]int{0, 3}, v: [2]int{4, 9}, ask: 1,
			reason: "The second coordinate gives 3b = 9, so b = 3.",
		},
		{
			question: "In R^2, the vector (5, 2) is written as a(1, 0) + b(1, 1). What is a? (enter a number)",
			answer:   "3", b1: [2]int{1, 0}, b2: [2]int{1, 1}, v: [2]int{5, 2}, ask: 0,
			reason: "The second coordinate gives b = 2, and then the first gives a + 2 = 5, so a = 3.",
		},
		{
			question: "In R^2, the vector (5, 3) is written as a(2, 1) + b(1, 1). What is a? (enter a number)",
			answer:   "2", b1: [2]int{2, 1}, b2: [2]int{1, 1}, v: [2]int{5, 3}, ask: 0,
			reason: "The equations are 2a + b = 5 and a + b = 3, whose solution is a = 2 and b = 1.",
		},
		{
			question: "In R^2, the vector (3, 7) is written as a(1, 2) + b(0, 1). What is b? (enter a number)",
			answer:   "1", b1: [2]int{1, 2}, b2: [2]int{0, 1}, v: [2]int{3, 7}, ask: 1,
			reason: "The first coordinate gives a = 3; then 2a + b = 7 gives b = 1.",
		},
		{
			question: "In R^2, the vector (5, 1) is written as a(1, 1) + b(1, -1). What is b? (enter a number)",
			answer:   "2", b1: [2]int{1, 1}, b2: [2]int{1, -1}, v: [2]int{5, 1}, ask: 1,
			reason: "The equations are a + b = 5 and a - b = 1, whose solution is a = 3 and b = 2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"b1":  strconv.Itoa(c.b1[0]) + "," + strconv.Itoa(c.b1[1]),
			"b2":  strconv.Itoa(c.b2[0]) + "," + strconv.Itoa(c.b2[1]),
			"v":   strconv.Itoa(c.v[0]) + "," + strconv.Itoa(c.v[1]),
			"ask": strconv.Itoa(c.ask), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Orthogonal complements (numeric)
// ---------------------------------------------------------------------------

type complementGen struct{}

func (g *complementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		n, m     int
		reason   string
	}
	claims := []claim{
		{
			question: "In R^3, the line spanned by (1, 0, 0) has dimension 1. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "2", n: 3, m: 1,
			reason: "The orthogonal complement of a line in R^3 is the plane through the origin perpendicular to it, of dimension 3 - 1 = 2.",
		},
		{
			question: "In R^3, the plane x + y + z = 0 has dimension 2. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "1", n: 3, m: 2,
			reason: "The normal line is one-dimensional, and dim(U-perp) = 3 - 2 = 1.",
		},
		{
			question: "In R^4, a subspace has dimension 3. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "1", n: 4, m: 3,
			reason: "dim(U-perp) = 4 - 3 = 1.",
		},
		{
			question: "In R^2, a line through the origin has dimension 1. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "1", n: 2, m: 1,
			reason: "The perpendicular line is also one-dimensional: dim(U-perp) = 2 - 1 = 1.",
		},
		{
			question: "In R^5, a subspace has dimension 2. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "3", n: 5, m: 2,
			reason: "dim(U-perp) = 5 - 2 = 3.",
		},
		{
			question: "In R^3, the whole space has dimension 3. What is the dimension of its orthogonal complement? (enter a number)",
			answer:   "0", n: 3, m: 3,
			reason: "The only vector orthogonal to every vector is 0, so the orthogonal complement is the zero subspace of dimension 0.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"n": strconv.Itoa(c.n), "m": strconv.Itoa(c.m), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 6. Projection onto a subspace (numeric)
// ---------------------------------------------------------------------------

type projectionMatrixGen struct{}

func (g *projectionMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		v, u     [2]int
		axis     int
		reason   string
	}
	claims := []claim{
		{
			question: "The projection of (3, 4) onto the line spanned by (1, 0) has what x-coordinate? (enter a number)",
			answer:   "3", v: [2]int{3, 4}, u: [2]int{1, 0}, axis: 0,
			reason: "The projection is ((3*1 + 4*0)/(1*1 + 0*0))(1, 0) = (3, 0), so the x-coordinate is 3.",
		},
		{
			question: "The projection of (3, 4) onto the line spanned by (1, 1) has what x-coordinate? (enter a decimal)",
			answer:   "3.5", v: [2]int{3, 4}, u: [2]int{1, 1}, axis: 0,
			reason: "The coefficient is (3 + 4)/2 = 3.5, so the projection is 3.5(1, 1) and its x-coordinate is 3.5.",
		},
		{
			question: "The projection of (2, 6) onto the line spanned by (0, 1) has what y-coordinate? (enter a number)",
			answer:   "6", v: [2]int{2, 6}, u: [2]int{0, 1}, axis: 1,
			reason: "The projection is ((2*0 + 6*1)/(0*0 + 1*1))(0, 1) = (0, 6), so the y-coordinate is 6.",
		},
		{
			question: "The projection of (4, 2) onto the line spanned by (2, 0) has what x-coordinate? (enter a number)",
			answer:   "4", v: [2]int{4, 2}, u: [2]int{2, 0}, axis: 0,
			reason: "The coefficient is (8 + 0)/4 = 2, so the projection is 2(2, 0) = (4, 0) and its x-coordinate is 4.",
		},
		{
			question: "The projection of (1, 3) onto the line spanned by (1, 1) has what x-coordinate? (enter a number)",
			answer:   "2", v: [2]int{1, 3}, u: [2]int{1, 1}, axis: 0,
			reason: "The coefficient is (1 + 3)/2 = 2, so the projection is 2(1, 1) and its x-coordinate is 2.",
		},
		{
			question: "The projection of (5, 5) onto the line spanned by (1, 2) has what x-coordinate? (enter a number)",
			answer:   "3", v: [2]int{5, 5}, u: [2]int{1, 2}, axis: 0,
			reason: "The coefficient is (5 + 10)/5 = 3, so the projection is 3(1, 2) and its x-coordinate is 3.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	axis := "x"
	if c.axis == 1 {
		axis = "y"
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"v":    strconv.Itoa(c.v[0]) + "," + strconv.Itoa(c.v[1]),
			"u":    strconv.Itoa(c.u[0]) + "," + strconv.Itoa(c.u[1]),
			"axis": axis, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. The spectral theorem for symmetric matrices
// ---------------------------------------------------------------------------

type spectralGen struct{}

func (g *spectralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "The spectral theorem says that a real symmetric matrix...",
			correct:     "is orthogonally diagonalizable with real eigenvalues",
			distractors: []string{"is always invertible", "has only positive eigenvalues", "is always a rotation"},
			reason:      "A real symmetric matrix has an orthonormal basis of eigenvectors, so it can be written Q D Q^T with Q orthogonal.",
		},
		{
			question:    "Eigenvectors of a real symmetric matrix belonging to distinct eigenvalues are...",
			correct:     "orthogonal",
			distractors: []string{"parallel", "linearly dependent", "complex conjugates"},
			reason:      "If Ax = lambda x and Ay = mu y with lambda different from mu, symmetry forces (lambda - mu)(x dot y) = 0, so x dot y = 0.",
		},
		{
			question:    "A symmetric matrix is orthogonally diagonalizable because...",
			correct:     "it has an orthonormal basis of eigenvectors",
			distractors: []string{"it is triangular", "its determinant is positive", "it has distinct eigenvalues"},
			reason:      "The spectral theorem supplies the orthonormal eigenbasis even when eigenvalues repeat.",
		},
		{
			question:    "Which matrix is symmetric?",
			correct:     "[[2, 1], [1, 3]]",
			distractors: []string{"[[0, 1], [-1, 0]]", "[[1, 2], [0, 1]]", "[[1, 1], [2, 1]]"},
			reason:      "A matrix is symmetric when it equals its transpose; the off-diagonal entries must match.",
		},
		{
			question:    "The eigenvalues of a real symmetric matrix are...",
			correct:     "real",
			distractors: []string{"always complex", "always positive", "always zero"},
			reason:      "Symmetry rules out nonreal eigenvalues, which is why the spectral theorem works entirely over R.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 8. Quadratic forms and definiteness
// ---------------------------------------------------------------------------

type definiteGen struct{}

func (g *definiteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A symmetric matrix is positive definite exactly when...",
			correct:     "all its eigenvalues are positive",
			distractors: []string{"all its entries are positive", "its determinant is zero", "it is diagonal"},
			reason:      "The quadratic form equals a sum of eigenvalue-weighted squares in an orthonormal eigenbasis, so positivity is exactly positivity of every eigenvalue.",
		},
		{
			question:    "The quadratic form x^2 + y^2 is...",
			correct:     "positive definite",
			distractors: []string{"negative definite", "indefinite", "positive semidefinite but not definite"},
			reason:      "It is a sum of squares and vanishes only at (x, y) = (0, 0).",
		},
		{
			question:    "The quadratic form x^2 - y^2 is...",
			correct:     "indefinite",
			distractors: []string{"positive definite", "negative definite", "identically zero"},
			reason:      "It is positive at (1, 0) and negative at (0, 1), so it takes both signs.",
		},
		{
			question:    "The symmetric matrix [[a, b], [b, c]] is positive definite exactly when...",
			correct:     "a is positive and ac - b^2 is positive",
			distractors: []string{"a is positive and ac - b^2 is negative", "a is negative and c is positive", "b^2 is larger than ac"},
			reason:      "Sylvester's criterion requires the leading principal minors a and ac - b^2 to be positive.",
		},
		{
			question:    "The level sets of a positive definite quadratic form are...",
			correct:     "ellipses",
			distractors: []string{"hyperbolas", "parallel lines", "parabolas"},
			reason:      "A positive definite form is a weighted sum of squares, so its level curves are bounded and closed, hence ellipses.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 9. The cross product (tuple)
// ---------------------------------------------------------------------------

type crossGen struct{}

func (g *crossGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		u, v     [3]int
		reason   string
	}
	claims := []claim{
		{
			question: "What is the cross product (1, 0, 0) x (0, 1, 0)? (enter a vector)",
			answer:   "(0,0,1)", u: [3]int{1, 0, 0}, v: [3]int{0, 1, 0},
			reason: "The standard basis vectors satisfy e1 x e2 = e3 = (0, 0, 1).",
		},
		{
			question: "What is the cross product (0, 1, 0) x (1, 0, 0)? (enter a vector)",
			answer:   "(0,0,-1)", u: [3]int{0, 1, 0}, v: [3]int{1, 0, 0},
			reason: "Reversing the order changes the sign, so e2 x e1 = -e3 = (0, 0, -1).",
		},
		{
			question: "What is the cross product (2, 0, 0) x (0, 3, 0)? (enter a vector)",
			answer:   "(0,0,6)", u: [3]int{2, 0, 0}, v: [3]int{0, 3, 0},
			reason: "The cross product is (2*3) e3 = (0, 0, 6).",
		},
		{
			question: "What is the cross product (1, 2, 3) x (0, 1, 0)? (enter a vector)",
			answer:   "(-3,0,1)", u: [3]int{1, 2, 3}, v: [3]int{0, 1, 0},
			reason: "Compute (2*0 - 3*1, 3*0 - 1*0, 1*1 - 2*0) = (-3, 0, 1).",
		},
		{
			question: "What is the cross product (1, 0, 1) x (0, 1, 1)? (enter a vector)",
			answer:   "(-1,-1,1)", u: [3]int{1, 0, 1}, v: [3]int{0, 1, 1},
			reason: "Compute (0*1 - 1*1, 1*0 - 1*1, 1*1 - 0*0) = (-1, -1, 1).",
		},
		{
			question: "What is the cross product (2, 1, 0) x (1, 0, 1)? (enter a vector)",
			answer:   "(1,-2,-1)", u: [3]int{2, 1, 0}, v: [3]int{1, 0, 1},
			reason: "Compute (1*1 - 0*0, 0*1 - 2*1, 2*0 - 1*1) = (1, -2, -1).",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"u":      strconv.Itoa(c.u[0]) + "," + strconv.Itoa(c.u[1]) + "," + strconv.Itoa(c.u[2]),
			"v":      strconv.Itoa(c.v[0]) + "," + strconv.Itoa(c.v[1]) + "," + strconv.Itoa(c.v[2]),
			"answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 10. The condition number (numeric)
// ---------------------------------------------------------------------------

type conditionGen struct{}

func (g *conditionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		largest  float64
		smallest float64
		reason   string
	}
	claims := []claim{
		{
			question: "A matrix has singular values 10 and 2. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "5", largest: 10, smallest: 2,
			reason: "The condition number is 10/2 = 5.",
		},
		{
			question: "A matrix has singular values 6 and 3. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "2", largest: 6, smallest: 3,
			reason: "The condition number is 6/3 = 2.",
		},
		{
			question: "A matrix has singular values 8 and 1. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "8", largest: 8, smallest: 1,
			reason: "The condition number is 8/1 = 8.",
		},
		{
			question: "A matrix has singular values 5 and 0.5. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "10", largest: 5, smallest: 0.5,
			reason: "The condition number is 5/0.5 = 10.",
		},
		{
			question: "A matrix has singular values 12 and 4. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "3", largest: 12, smallest: 4,
			reason: "The condition number is 12/4 = 3.",
		},
		{
			question: "A matrix has singular values 9 and 9. What is its condition number, the largest divided by the smallest? (enter a number)",
			answer:   "1", largest: 9, smallest: 9,
			reason: "Equal singular values give condition number 9/9 = 1, the best possible conditioning.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"largest":  strconv.FormatFloat(c.largest, 'g', -1, 64),
			"smallest": strconv.FormatFloat(c.smallest, 'g', -1, 64),
			"answer":   c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. The Cayley-Hamilton theorem
// ---------------------------------------------------------------------------

type cayleyHamiltonGen struct{}

func (g *cayleyHamiltonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "The Cayley-Hamilton theorem says that every square matrix...",
			correct:     "satisfies its own characteristic polynomial",
			distractors: []string{"is diagonalizable", "has distinct eigenvalues", "is invertible"},
			reason:      "Substituting A into its characteristic polynomial gives the zero matrix, whether or not A is diagonalizable.",
		},
		{
			question:    "A 2x2 matrix has characteristic polynomial lambda^2 - 5 lambda + 6. Cayley-Hamilton gives...",
			correct:     "A^2 - 5A + 6I = 0",
			distractors: []string{"A^2 + 5A + 6I = 0", "A - 5I = 0", "A^2 = 6I only"},
			reason:      "Replace lambda by A and the constant term by that constant times I, giving A^2 - 5A + 6I = 0.",
		},
		{
			question:    "Cayley-Hamilton lets you compute which of these for a 2x2 matrix?",
			correct:     "powers of A as combinations of I and A",
			distractors: []string{"only the determinant", "only the trace", "the inverse of a non-square matrix"},
			reason:      "The relation A^2 = 5A - 6I expresses every higher power in terms of I and A.",
		},
		{
			question:    "The minimal polynomial of a matrix...",
			correct:     "divides the characteristic polynomial",
			distractors: []string{"is always the characteristic polynomial", "is always linear", "has larger degree than the matrix"},
			reason:      "The minimal polynomial is the monic polynomial of least degree that annihilates A, and Cayley-Hamilton forces it to divide the characteristic polynomial.",
		},
		{
			question:    "What is the degree of the characteristic polynomial of a 3x3 matrix?",
			correct:     "3",
			distractors: []string{"2", "9", "0"},
			reason:      "The characteristic polynomial is det(A - lambda I), a determinant of a 3x3 matrix, so its degree is 3.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}
