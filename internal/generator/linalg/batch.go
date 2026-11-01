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
	"fmt"
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
	return spaceAxiomsProblem()
}

// ---------------------------------------------------------------------------
// 2. Kernel and image (numeric)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// 3. Injective, surjective and invertible maps
// ---------------------------------------------------------------------------

type invertibleMapGen struct{}

func (g *invertibleMapGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return invertibleMapProblem()
}

// ---------------------------------------------------------------------------
// 4. Coordinates and change of basis (numeric)
// ---------------------------------------------------------------------------

type basisChangeGen struct{}

func (g *basisChangeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Pairs of 2D unimodular basis vectors (det = 1), strictly excluding the standard basis
	unimodularPairs := [][2][2]int{
		{{1, 1}, {1, 2}},
		{{1, -1}, {1, 0}},
		{{2, 1}, {1, 1}},
		{{1, 2}, {0, 1}},
		{{1, 0}, {1, 1}},
		{{2, 3}, {1, 2}},
		{{3, 2}, {1, 1}},
		{{1, 1}, {2, 3}},
		{{1, 2}, {1, 3}},
		{{2, 1}, {3, 2}},
	}
	pair := unimodularPairs[rand.Intn(len(unimodularPairs))]
	b1 := pair[0]
	b2 := pair[1]

	var a, b, vx, vy, ans, ask int
	var varName string
	for attempt := 0; attempt < 50; attempt++ {
		a = rand.Intn(11) - 5
		b = rand.Intn(11) - 5
		if a == 0 || b == 0 {
			continue
		}
		vx = a*b1[0] + b*b2[0]
		vy = a*b1[1] + b*b2[1]
		ask = rand.Intn(2)
		ans = a
		varName = "a"
		if ask == 1 {
			ans = b
			varName = "b"
		}
		if ans != vx && ans != vy && -ans != vx && -ans != vy &&
			ans != b1[0] && ans != b1[1] && -ans != b1[0] && -ans != b1[1] &&
			ans != b2[0] && ans != b2[1] && -ans != b2[0] && -ans != b2[1] {
			break
		}
	}

	q := fmt.Sprintf("In R^2, the vector (%d, %d) is written as a(%d, %d) + b(%d, %d). What is %s? (enter a number)",
		vx, vy, b1[0], b1[1], b2[0], b2[1], varName)
	ex := fmt.Sprintf("The equations are %da + %db = %d and %da + %db = %d. "+
		"Solving the system yields a = %d and b = %d.",
		b1[0], b2[0], vx, b1[1], b2[1], vy, a, b)

	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(ans),
		Explanation: ex,
		Facts: map[string]string{
			"b1":     fmt.Sprintf("%d,%d", b1[0], b1[1]),
			"b2":     fmt.Sprintf("%d,%d", b2[0], b2[1]),
			"v":      fmt.Sprintf("%d,%d", vx, vy),
			"ask":    strconv.Itoa(ask),
			"answer": strconv.Itoa(ans),
		},
	}
}

// ---------------------------------------------------------------------------
// 5. Orthogonal complements (numeric)
// ---------------------------------------------------------------------------

type complementGen struct{}

func (g *complementGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 2 + rand.Intn(4)   // 2 to 5
	m := 1 + rand.Intn(n-1) // 1 to n-1
	want := n - m

	var desc, reason string
	if m == 1 {
		desc = "the line spanned by a non-zero vector has dimension 1."
		reason = fmt.Sprintf("The orthogonal complement of a line in R^%d has dimension %d - 1 = %d.", n, n, want)
	} else if m == n-1 {
		desc = fmt.Sprintf("a hyperplane has dimension %d.", m)
		reason = fmt.Sprintf("The orthogonal complement of a hyperplane of dimension %d in R^%d is 1-dimensional: %d - %d = %d.", m, n, n, m, want)
	} else {
		desc = fmt.Sprintf("a subspace spanned by %d linearly independent vectors has dimension %d.", m, m)
		reason = fmt.Sprintf("By the orthogonal complement theorem, dim(U-perp) = %d - %d = %d.", n, m, want)
	}

	q := fmt.Sprintf("In R^%d, %s What is the dimension of its orthogonal complement? (enter a number)", n, desc)
	ansStr := strconv.Itoa(want)
	return generator.Problem{
		Question:    q,
		Answer:      ansStr,
		Explanation: reason,
		Facts: map[string]string{
			"n":      strconv.Itoa(n),
			"m":      strconv.Itoa(m),
			"answer": ansStr,
		},
	}
}

// ---------------------------------------------------------------------------
// 6. Projection onto a subspace (numeric)
// ---------------------------------------------------------------------------

type projectionMatrixGen struct{}

func (g *projectionMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Strictly non-axis spanning lines to prevent prompt coordinate leaks
	dirs := [][2]int{
		{1, 1}, {1, -1}, {1, 2}, {2, 1}, {1, 3}, {3, 1}, {2, 3}, {3, 2},
	}
	u := dirs[rand.Intn(len(dirs))]
	uNormSq := u[0]*u[0] + u[1]*u[1]

	var c, k, vx, vy, ans int
	var coordName, axisStr string
	for attempt := 0; attempt < 50; attempt++ {
		c = rand.Intn(4) + 1
		if rand.Intn(2) == 1 {
			c = -c
		}
		k = rand.Intn(3) + 1
		if rand.Intn(2) == 1 {
			k = -k
		}
		uPerp := [2]int{-u[1], u[0]}
		vx = c*u[0] + k*uPerp[0]
		vy = c*u[1] + k*uPerp[1]

		axis := rand.Intn(2)
		coordName = "x-coordinate"
		axisStr = "x"
		ans = c * u[0]
		if axis == 1 {
			coordName = "y-coordinate"
			axisStr = "y"
			ans = c * u[1]
		}
		if ans != vx && ans != vy && -ans != vx && -ans != vy &&
			ans != u[0] && ans != u[1] && -ans != u[0] && -ans != u[1] {
			break
		}
	}

	q := fmt.Sprintf("The projection of (%d, %d) onto the line spanned by (%d, %d) has what %s? (enter a number)",
		vx, vy, u[0], u[1], coordName)
	ex := fmt.Sprintf("The coefficient is ((%d)(%d) + (%d)(%d)) / ((%d)^2 + (%d)^2) = %d / %d = %d. "+
		"So the projection is %d(%d, %d) = (%d, %d) and its %s is %d.",
		vx, u[0], vy, u[1], u[0], u[1], c*uNormSq, uNormSq, c,
		c, u[0], u[1], c*u[0], c*u[1], coordName, ans)

	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(ans),
		Explanation: ex,
		Facts: map[string]string{
			"v":      fmt.Sprintf("%d,%d", vx, vy),
			"u":      fmt.Sprintf("%d,%d", u[0], u[1]),
			"axis":   axisStr,
			"answer": strconv.Itoa(ans),
		},
	}
}

// ---------------------------------------------------------------------------
// 7. The spectral theorem for symmetric matrices
// ---------------------------------------------------------------------------

type spectralGen struct{}

func (g *spectralGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return spectralProblem()
}

// ---------------------------------------------------------------------------
// 8. Quadratic forms and definiteness
// ---------------------------------------------------------------------------

type definiteGen struct{}

func (g *definiteGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return definiteProblem()
}

// ---------------------------------------------------------------------------
// 9. The cross product (tuple)
// ---------------------------------------------------------------------------

type crossGen struct{}

func (g *crossGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for {
		u := [3]int{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		v := [3]int{rand.Intn(7) - 3, rand.Intn(7) - 3, rand.Intn(7) - 3}
		if u[0] == 0 && u[1] == 0 && u[2] == 0 {
			continue
		}
		if v[0] == 0 && v[1] == 0 && v[2] == 0 {
			continue
		}
		x := u[1]*v[2] - u[2]*v[1]
		y := u[2]*v[0] - u[0]*v[2]
		z := u[0]*v[1] - u[1]*v[0]
		if x == 0 && y == 0 && z == 0 {
			continue
		}

		q := fmt.Sprintf("What is the cross product (%d, %d, %d) x (%d, %d, %d)? (enter a vector)",
			u[0], u[1], u[2], v[0], v[1], v[2])
		ans := fmt.Sprintf("(%d,%d,%d)", x, y, z)
		ex := fmt.Sprintf("Compute ((%d)(%d) - (%d)(%d), (%d)(%d) - (%d)(%d), (%d)(%d) - (%d)(%d)) = (%d, %d, %d).",
			u[1], v[2], u[2], v[1], u[2], v[0], u[0], v[2], u[0], v[1], u[1], v[0], x, y, z)

		return generator.Problem{
			Question:    q,
			Answer:      ans,
			Explanation: ex,
			Facts: map[string]string{
				"u":      fmt.Sprintf("%d,%d,%d", u[0], u[1], u[2]),
				"v":      fmt.Sprintf("%d,%d,%d", v[0], v[1], v[2]),
				"answer": ans,
			},
		}
	}
}

// ---------------------------------------------------------------------------
// 10. The condition number (numeric)
// ---------------------------------------------------------------------------

type conditionGen struct{}

func (g *conditionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type pair struct {
		ratio int
		small float64
	}
	pairs := []pair{
		{1, 2}, {1, 3}, {1, 4}, {1, 5},
		{2, 3}, {2, 4}, {2, 5},
		{3, 2}, {3, 4}, {3, 5},
		{4, 2}, {4, 3}, {4, 5},
		{5, 2}, {5, 3}, {5, 4},
		{6, 2}, {6, 3}, {6, 4},
		{8, 0.5}, {8, 2}, {8, 3},
		{10, 0.5}, {10, 2}, {10, 3},
	}
	p := pairs[rand.Intn(len(pairs))]
	ratio := p.ratio
	smallest := p.small
	largest := float64(ratio) * smallest

	ansStr := strconv.Itoa(ratio)
	var q, ex string
	if ratio == 1 {
		q = fmt.Sprintf("A matrix has singular values %g and %g. What is its condition number, the largest divided by the smallest? (enter a number)",
			largest, smallest)
		ex = fmt.Sprintf("Equal singular values give condition number %g/%g = 1, the best possible conditioning.",
			largest, smallest)
	} else {
		q = fmt.Sprintf("A matrix has singular values %g and %g. What is its condition number, the largest divided by the smallest? (enter a number)",
			largest, smallest)
		ex = fmt.Sprintf("The condition number is %g/%g = %s.", largest, smallest, ansStr)
	}

	return generator.Problem{
		Question:    q,
		Answer:      ansStr,
		Explanation: ex,
		Facts: map[string]string{
			"largest":  strconv.FormatFloat(largest, 'g', -1, 64),
			"smallest": strconv.FormatFloat(smallest, 'g', -1, 64),
			"answer":   ansStr,
		},
	}
}

// ---------------------------------------------------------------------------
// 11. The Cayley-Hamilton theorem
// ---------------------------------------------------------------------------

type cayleyHamiltonGen struct{}

func (g *cayleyHamiltonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return cayleyHamiltonProblem()
}
