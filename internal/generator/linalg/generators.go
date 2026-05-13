package linalg

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("linalg.vector.concept", &vectorConceptGen{})
	reg.Register("linalg.vector.add", &vectorAddGen{})
	reg.Register("linalg.vector.dot", &vectorDotGen{})
	reg.Register("linalg.matrix.concept", &matrixConceptGen{})
	reg.Register("linalg.matrix.add", &matrixAddGen{})
	reg.Register("linalg.matrix.mult", &matrixMultGen{})
	reg.Register("linalg.matrix.identity", &matrixIdentityGen{})
	reg.Register("linalg.det.2x2", &det2x2Gen{})
	reg.Register("linalg.det.3x3", &det3x3Gen{})
	reg.Register("linalg.systems.matrix", &systemsMatrixGen{})
	reg.Register("linalg.cramer", &cramerGen{})
	reg.Register("linalg.eigen.concept", &eigenConceptGen{})
	reg.Register("linalg.eigen.compute", &eigenComputeGen{})
	reg.Register("linalg.transformations", &transformationsGen{})
	reg.Register("linalg.span", &spanGen{})
	reg.Register("linalg.basis", &basisGen{})
}

type vectorConceptGen struct{}

func (g *vectorConceptGen) Generate(difficulty float64) generator.Problem {
	triples := [][3]int{
		{3, 4, 5}, {4, 3, 5},
		{6, 8, 10}, {8, 6, 10},
		{5, 12, 13}, {12, 5, 13},
		{8, 15, 17}, {15, 8, 17},
		{9, 12, 15}, {12, 9, 15},
		{7, 24, 25}, {24, 7, 25},
	}
	t := triples[rand.Intn(len(triples))]
	x, y, mag := t[0], t[1], t[2]
	if rand.Intn(2) == 0 {
		x = -x
	}
	if rand.Intn(2) == 0 {
		y = -y
	}
	return generator.Problem{
		Question:    fmt.Sprintf("What is the magnitude of vector (%d,%d)?", x, y),
		Answer:      fmt.Sprintf("%d", mag),
		Explanation: fmt.Sprintf("|(%d,%d)| = sqrt(%d²+%d²) = sqrt(%d+%d) = sqrt(%d) = %d", x, y, x, y, x*x, y*y, x*x+y*y, mag),
	}
}

type vectorAddGen struct{}

func (g *vectorAddGen) Generate(difficulty float64) generator.Problem {
	a, b := rand.Intn(9)-4, rand.Intn(9)-4
	c, d := rand.Intn(9)-4, rand.Intn(9)-4
	if rand.Intn(2) == 0 {
		sumX, sumY := a+c, b+d
		return generator.Problem{
			Question:    fmt.Sprintf("Add vectors: (%d,%d) + (%d,%d) = ?", a, b, c, d),
			Answer:      fmt.Sprintf("(%d,%d)", sumX, sumY),
			Explanation: fmt.Sprintf("(%d,%d) + (%d,%d) = (%d,%d)", a, b, c, d, sumX, sumY),
		}
	}
	scalar := rand.Intn(4) + 2
	if rand.Intn(2) == 0 {
		scalar = -scalar
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compute scalar multiplication: %d * (%d,%d) = ?", scalar, a, b),
		Answer:      fmt.Sprintf("(%d,%d)", scalar*a, scalar*b),
		Explanation: fmt.Sprintf("%d * (%d,%d) = (%d,%d)", scalar, a, b, scalar*a, scalar*b),
	}
}

type vectorDotGen struct{}

func (g *vectorDotGen) Generate(difficulty float64) generator.Problem {
	a, b := rand.Intn(7)-3, rand.Intn(7)-3
	c, d := rand.Intn(7)-3, rand.Intn(7)-3
	dot := a*c + b*d
	return generator.Problem{
		Question:    fmt.Sprintf("Compute the dot product: (%d,%d) · (%d,%d) = ?", a, b, c, d),
		Answer:      fmt.Sprintf("%d", dot),
		Explanation: fmt.Sprintf("(%d,%d) · (%d,%d) = %d*%d + %d*%d = %d+%d = %d", a, b, c, d, a, c, b, d, a*c, b*d, dot),
	}
}

type matrixConceptGen struct{}

func (g *matrixConceptGen) Generate(difficulty float64) generator.Problem {
	var rows, cols int
	if rand.Intn(2) == 0 {
		rows, cols = 2, 2
	} else {
		rows, cols = 3, 3
	}
	m := make([][]int, rows)
	for i := range m {
		m[i] = make([]int, cols)
		for j := range m[i] {
			m[i][j] = rand.Intn(9) + 1
		}
	}
	r := rand.Intn(rows) + 1
	c := rand.Intn(cols) + 1
	answer := m[r-1][c-1]
	return generator.Problem{
		Question:    fmt.Sprintf("In the matrix %s, what is the (%d,%d) entry?", formatMatrix(m), r, c),
		Answer:      fmt.Sprintf("%d", answer),
		Explanation: fmt.Sprintf("The entry at row %d, column %d is %d.", r, c, answer),
	}
}

type matrixAddGen struct{}

func (g *matrixAddGen) Generate(difficulty float64) generator.Problem {
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := range a {
		a[i] = make([]int, 2)
		b[i] = make([]int, 2)
		for j := range a[i] {
			a[i][j] = rand.Intn(9) - 4
			b[i][j] = rand.Intn(9) - 4
		}
	}
	if rand.Intn(2) == 0 {
		sum := make([][]int, 2)
		for i := range sum {
			sum[i] = make([]int, 2)
			for j := range sum[i] {
				sum[i][j] = a[i][j] + b[i][j]
			}
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Add matrices: %s + %s = ?", formatMatrix(a), formatMatrix(b)),
			Answer:      formatMatrix(sum),
			Explanation: fmt.Sprintf("%s + %s = %s (add corresponding entries)", formatMatrix(a), formatMatrix(b), formatMatrix(sum)),
		}
	}
	scalar := rand.Intn(4) + 2
	if rand.Intn(2) == 0 {
		scalar = -scalar
	}
	prod := make([][]int, 2)
	for i := range prod {
		prod[i] = make([]int, 2)
		for j := range prod[i] {
			prod[i][j] = scalar * a[i][j]
		}
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compute scalar multiplication: %d * %s = ?", scalar, formatMatrix(a)),
		Answer:      formatMatrix(prod),
		Explanation: fmt.Sprintf("%d * %s = %s", scalar, formatMatrix(a), formatMatrix(prod)),
	}
}

type matrixMultGen struct{}

func (g *matrixMultGen) Generate(difficulty float64) generator.Problem {
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := range a {
		a[i] = make([]int, 2)
		b[i] = make([]int, 2)
		for j := range a[i] {
			a[i][j] = rand.Intn(7) - 3
			b[i][j] = rand.Intn(7) - 3
		}
	}
	c := make([][]int, 2)
	for i := range c {
		c[i] = make([]int, 2)
	}
	c[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0]
	c[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1]
	c[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0]
	c[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1]
	return generator.Problem{
		Question:    fmt.Sprintf("Multiply matrices: %s × %s = ?", formatMatrix(a), formatMatrix(b)),
		Answer:      formatMatrix(c),
		Explanation: fmt.Sprintf("%s × %s = %s", formatMatrix(a), formatMatrix(b), formatMatrix(c)),
	}
}

type matrixIdentityGen struct{}

func (g *matrixIdentityGen) Generate(difficulty float64) generator.Problem {
	candidates := [][4]int{
		{2, 1, 1, 1}, // det=1
		{3, 2, 2, 1}, // det=-1
		{3, 4, 2, 3}, // det=1
		{4, 3, 3, 2}, // det=-1
		{5, 3, 3, 2}, // det=1
		{5, 7, 2, 3}, // det=1
		{4, 5, 3, 4}, // det=1
		{3, 1, 4, 1}, // det=-1
		{2, 5, 1, 3}, // det=1
		{1, 2, 3, 5}, // det=-1
		{2, 3, 5, 7}, // det=-1
		{7, 4, 5, 3}, // det=1
		{6, 5, 5, 4}, // det=-1
		{3, 5, 2, 3}, // det=-1
		{4, 7, 1, 2}, // det=1
		{7, 5, 4, 3}, // det=1
	}
	c := candidates[rand.Intn(len(candidates))]
	a, b, cVal, d := c[0], c[1], c[2], c[3]
	det := a*d - b*cVal
	if det < 0 {
		det = -det
	}
	inv := make([][]int, 2)
	for i := range inv {
		inv[i] = make([]int, 2)
	}
	inv[0][0] = d / det
	inv[0][1] = -b / det
	inv[1][0] = -cVal / det
	inv[1][1] = a / det
	return generator.Problem{
		Question:    fmt.Sprintf("Find the inverse of the matrix %s.", formatMatrix2([][]int{{a, b}, {cVal, d}})),
		Answer:      formatMatrix(inv),
		Explanation: fmt.Sprintf("det = %d*%d - %d*%d = %d, inverse = (1/%d)×[[%d,%d],[%d,%d]] = %s", a, d, b, cVal, a*d-b*cVal, det, d, -b, -cVal, a, formatMatrix(inv)),
	}
}

type det2x2Gen struct{}

func (g *det2x2Gen) Generate(difficulty float64) generator.Problem {
	a, b := rand.Intn(9)-4, rand.Intn(9)-4
	c, d := rand.Intn(9)-4, rand.Intn(9)-4
	det := a*d - b*c
	return generator.Problem{
		Question:    fmt.Sprintf("Find the determinant of %s.", formatMatrix2([][]int{{a, b}, {c, d}})),
		Answer:      fmt.Sprintf("%d", det),
		Explanation: fmt.Sprintf("det = %d*%d - %d*%d = %d-%d = %d", a, d, b, c, a*d, b*c, det),
	}
}

type det3x3Gen struct{}

func (g *det3x3Gen) Generate(difficulty float64) generator.Problem {
	m := make([][]int, 3)
	for i := range m {
		m[i] = make([]int, 3)
		for j := range m[i] {
			m[i][j] = rand.Intn(7) - 3
		}
	}
	a, b, c := m[0][0], m[0][1], m[0][2]
	d, e, f := m[1][0], m[1][1], m[1][2]
	gVal, h, i := m[2][0], m[2][1], m[2][2]
	det := a*(e*i-f*h) - b*(d*i-f*gVal) + c*(d*h-e*gVal)
	return generator.Problem{
		Question:    fmt.Sprintf("Find the determinant of %s.", formatMatrix(m)),
		Answer:      fmt.Sprintf("%d", det),
		Explanation: fmt.Sprintf("det = %d*(%d*%d-%d*%d) - %d*(%d*%d-%d*%d) + %d*(%d*%d-%d*%d) = %d", a, e, i, f, h, b, d, i, f, gVal, c, d, h, e, gVal, det),
	}
}

type systemsMatrixGen struct{}

func (g *systemsMatrixGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(5) - 2
	y := rand.Intn(5) - 2
	a, b := rand.Intn(5)-2, rand.Intn(5)-2
	c, d := rand.Intn(5)-2, rand.Intn(5)-2
	if a == 0 && c == 0 {
		a = 1
		c = 2
	}
	if b == 0 && d == 0 {
		b = 1
		d = 3
	}
	e1 := a*x + b*y
	e2 := c*x + d*y
	return generator.Problem{
		Question:    fmt.Sprintf("Solve the linear system using matrices: %dx + %dy = %d, %dx + %dy = %d", a, b, e1, c, d, e2),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("The solution is x=%d, y=%d.", x, y),
	}
}

type cramerGen struct{}

func (g *cramerGen) Generate(difficulty float64) generator.Problem {
	x := rand.Intn(5) - 2
	y := rand.Intn(5) - 2
	a, b := rand.Intn(5)-2, rand.Intn(5)-2
	c, d := rand.Intn(5)-2, rand.Intn(5)-2
	if a == 0 && c == 0 {
		a = 1
		c = 2
	}
	if b == 0 && d == 0 {
		b = 1
		d = 3
	}
	e1 := a*x + b*y
	e2 := c*x + d*y
	det := a*d - b*c
	if det == 0 {
		det = 1
	}
	askX := rand.Intn(2) == 0
	if askX {
		detX := e1*d - b*e2
		valX := detX / det
		return generator.Problem{
			Question:    fmt.Sprintf("Use Cramer's rule to find x: %dx + %dy = %d, %dx + %dy = %d", a, b, e1, c, d, e2),
			Answer:      fmt.Sprintf("%d", valX),
			Explanation: fmt.Sprintf("x = det([[%d,%d],[%d,%d]]) / det([[%d,%d],[%d,%d]]) = (%d*%d-%d*%d)/(%d*%d-%d*%d) = %d/%d = %d", e1, b, e2, d, a, b, c, d, e1, d, b, e2, a, d, b, c, detX, det, valX),
		}
	}
	detY := a*e2 - e1*c
	valY := detY / det
	return generator.Problem{
		Question:    fmt.Sprintf("Use Cramer's rule to find y: %dx + %dy = %d, %dx + %dy = %d", a, b, e1, c, d, e2),
		Answer:      fmt.Sprintf("%d", valY),
		Explanation: fmt.Sprintf("y = det([[%d,%d],[%d,%d]]) / det([[%d,%d],[%d,%d]]) = (%d*%d-%d*%d)/(%d*%d-%d*%d) = %d/%d = %d", a, e1, c, e2, a, b, c, d, a, e2, e1, c, a, d, b, c, detY, det, valY),
	}
}

type eigenConceptGen struct{}

func (g *eigenConceptGen) Generate(difficulty float64) generator.Problem {
	if rand.Intn(2) == 0 {
		return generator.Problem{
			Question:    "What is an eigenvalue λ of a square matrix A?",
			Answer:      "A scalar λ such that det(A-λI)=0, equivalently Av=λv for some nonzero v",
			Explanation: "An eigenvalue λ satisfies det(A-λI)=0, meaning there exists a nonzero vector v with Av=λv.",
		}
	}
	a, d := rand.Intn(5)-2, rand.Intn(5)-2
	if a == 0 {
		a = 1
	}
	if d == 0 {
		d = 2
	}
	lambda := rand.Intn(7) - 3
	det := (a - lambda) * (d - lambda)
	answer := "no"
	if det == 0 {
		answer = "yes"
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Is λ=%d an eigenvalue of matrix %s?", lambda, formatMatrix2([][]int{{a, 0}, {0, d}})),
		Answer:      answer,
		Explanation: fmt.Sprintf("det(A-λI) = det([[%d,%d],[%d,%d]]) = (%d)*(%d) = %d. Since det %s 0, λ=%d %s an eigenvalue.", a-lambda, 0, 0, d-lambda, a-lambda, d-lambda, det, map[bool]string{true: "=", false: "≠"}[det == 0], lambda, map[bool]string{true: "is", false: "is not"}[det == 0]),
	}
}

type eigenComputeGen struct{}

func (g *eigenComputeGen) Generate(difficulty float64) generator.Problem {
	a := rand.Intn(7) - 3
	c := rand.Intn(7) - 3
	d := rand.Intn(7) - 3
	if a == 0 {
		a = 2
	}
	if d == 0 {
		d = 1
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Find the eigenvalues of the triangular matrix %s.", formatMatrix2([][]int{{a, c}, {0, d}})),
		Answer:      fmt.Sprintf("%d,%d", a, d),
		Explanation: fmt.Sprintf("For a triangular matrix, eigenvalues are the diagonal entries: λ₁=%d, λ₂=%d.", a, d),
	}
}

type transformationsGen struct{}

func (g *transformationsGen) Generate(difficulty float64) generator.Problem {
	a, b := rand.Intn(5)+1, rand.Intn(5)+1
	x, y := rand.Intn(5)-2, rand.Intn(5)-2
	return generator.Problem{
		Question:    fmt.Sprintf("Under the linear transformation T(x,y) = (%dx,%dy), where does the point (%d,%d) map to?", a, b, x, y),
		Answer:      fmt.Sprintf("(%d,%d)", a*x, b*y),
		Explanation: fmt.Sprintf("T(%d,%d) = (%d*%d, %d*%d) = (%d,%d)", x, y, a, x, b, y, a*x, b*y),
	}
}

type spanGen struct{}

func (g *spanGen) Generate(difficulty float64) generator.Problem {
	if rand.Intn(2) == 0 {
		v1, v2 := rand.Intn(5)-2, rand.Intn(5)-2
		w1, w2 := rand.Intn(5)-2, rand.Intn(5)-2
		if v1 == 0 && v2 == 0 {
			v1, v2 = 1, 0
		}
		if w1 == 0 && w2 == 0 {
			w1, w2 = 0, 1
		}
		det := v1*w2 - v2*w1
		answer := "no"
		if det != 0 {
			answer = "yes"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Do the vectors v=(%d,%d) and w=(%d,%d) span R²?", v1, v2, w1, w2),
			Answer:      answer,
			Explanation: fmt.Sprintf("det([[%d,%d],[%d,%d]]) = %d. Since determinant is %s, the vectors %s span R².", v1, v2, w1, w2, det, map[bool]string{true: "nonzero", false: "zero"}[det != 0], answer),
		}
	}
	k := rand.Intn(4) + 2
	if rand.Intn(2) == 0 {
		k = -k
	}
	v1, v2 := rand.Intn(5)-2, rand.Intn(5)-2
	if v1 == 0 && v2 == 0 {
		v1, v2 = 1, 2
	}
	w1, w2 := k*v1, k*v2
	return generator.Problem{
		Question:    fmt.Sprintf("Do the vectors v=(%d,%d) and w=(%d,%d) span R²?", v1, v2, w1, w2),
		Answer:      "no",
		Explanation: fmt.Sprintf("w = %d*v, so the vectors are linearly dependent. Two dependent vectors cannot span R².", k),
	}
}

type basisGen struct{}

func (g *basisGen) Generate(difficulty float64) generator.Problem {
	qType := rand.Intn(3)
	switch qType {
	case 0:
		return generator.Problem{
			Question:    "What is the dimension of R³?",
			Answer:      "3",
			Explanation: "R³ has dimension 3 because it has a basis of 3 vectors, e.g., {(1,0,0), (0,1,0), (0,0,1)}.",
		}
	case 1:
		return generator.Problem{
			Question:    "What is the dimension of R²?",
			Answer:      "2",
			Explanation: "R² has dimension 2 because it has a basis of 2 vectors, e.g., {(1,0), (0,1)}.",
		}
	default:
		v1, v2 := rand.Intn(5)-2, rand.Intn(5)-2
		w1, w2 := rand.Intn(5)-2, rand.Intn(5)-2
		if v1 == 0 && v2 == 0 {
			v1, v2 = 1, 0
		}
		if w1 == 0 && w2 == 0 {
			w1, w2 = 0, 1
		}
		det := v1*w2 - v2*w1
		answer := "no"
		if det != 0 {
			answer = "yes"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Is the set {(%d,%d), (%d,%d)} a basis for R²?", v1, v2, w1, w2),
			Answer:      answer,
			Explanation: fmt.Sprintf("det = %d, so the vectors are %s. Since they %s span R² and are %s, they %s a basis.", det, map[bool]string{true: "linearly independent", false: "linearly dependent"}[det != 0], map[bool]string{true: "do", false: "do not"}[det != 0], map[bool]string{true: "independent", false: "dependent"}[det != 0], map[bool]string{true: "form", false: "do not form"}[det != 0]),
		}
	}
}

func formatMatrix(m [][]int) string {
	if len(m) == 2 && len(m[0]) == 2 {
		return formatMatrix2(m)
	}
	if len(m) == 3 && len(m[0]) == 3 {
		return formatMatrix3(m)
	}
	return fmt.Sprintf("%v", m)
}

func formatMatrix2(m [][]int) string {
	return fmt.Sprintf("[[%d,%d],[%d,%d]]", m[0][0], m[0][1], m[1][0], m[1][1])
}

func formatMatrix3(m [][]int) string {
	return fmt.Sprintf("[[%d,%d,%d],[%d,%d,%d],[%d,%d,%d]]",
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}
