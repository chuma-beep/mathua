package linalg

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func linalgBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func linalgBridgeProblem(ordinal int, question, answer, explanation string, facts map[string]string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       facts,
		KPOrdinal:   ordinal,
	}
}

type rowReductionGen struct{}

func (g *rowReductionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 1 + rand.Intn(5)
		b := 1 + rand.Intn(5)
		c := 1 + rand.Intn(5)
		d := 1 + rand.Intn(5)
		k := 2 + rand.Intn(3)
		ans := c + k*a
		return linalgBridgeProblem(o,
			fmt.Sprintf("The matrix [[%d, %d], [%d, %d]] is changed by replacing R2 with R2 + %d R1. What is the new entry in row 2, column 1? (enter a number)", a, b, c, d, k),
			strconv.Itoa(ans),
			fmt.Sprintf("The operation adds %d times the first-row entry %d to the old entry %d, giving %d + %d x %d = %d.", k, a, c, c, k, a, ans),
			map[string]string{"kind": "row_operation", "answer": strconv.Itoa(ans)},
		)
	case 2:
		rank := 1 + rand.Intn(3)
		var m [3][3]int
		for i := 0; i < rank; i++ {
			m[i][i] = 1
			for j := i + 1; j < 3; j++ {
				m[i][j] = rand.Intn(5)
			}
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("How many pivots does the echelon matrix [[%d, %d, %d], [%d, %d, %d], [%d, %d, %d]] have? (enter a number)",
				m[0][0], m[0][1], m[0][2], m[1][0], m[1][1], m[1][2], m[2][0], m[2][1], m[2][2]),
			strconv.Itoa(rank),
			fmt.Sprintf("The matrix has %d nonzero leading entries, so it has %d pivots and rank %d.", rank, rank, rank),
			map[string]string{"kind": "pivot_count", "answer": strconv.Itoa(rank)},
		)
	default:
		x := 1 + rand.Intn(6)
		y := 1 + rand.Intn(6)
		return linalgBridgeProblem(o,
			fmt.Sprintf("A system has reduced row echelon form [[1, 0, %d], [0, 1, %d]]. What is x? (enter a number)", x, y),
			strconv.Itoa(x),
			fmt.Sprintf("The first row reads x = %d and the second row reads y = %d.", x, y),
			map[string]string{"kind": "read_solution", "answer": strconv.Itoa(x)},
		)
	}
}

type fundamentalSubspacesGen struct{}

func (g *fundamentalSubspacesGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	m := 2 + rand.Intn(4)
	n := 2 + rand.Intn(4)
	maxRank := m
	if n < maxRank {
		maxRank = n
	}
	rank := 1 + rand.Intn(maxRank)
	switch o {
	case 1:
		return linalgBridgeProblem(o,
			fmt.Sprintf("A %d by %d matrix has rank %d. What is the dimension of its column space? (enter a number)", m, n, rank),
			strconv.Itoa(rank),
			fmt.Sprintf("The column space has dimension equal to the rank, which is %d.", rank),
			map[string]string{"kind": "column_space", "answer": strconv.Itoa(rank)},
		)
	case 2:
		ans := m - rank
		return linalgBridgeProblem(o,
			fmt.Sprintf("A %d by %d matrix has rank %d. What is the dimension of its left null space? (enter a number)", m, n, rank),
			strconv.Itoa(ans),
			fmt.Sprintf("The left null space lives in R^%d and has dimension %d - %d = %d.", m, m, rank, ans),
			map[string]string{"kind": "left_null_space", "answer": strconv.Itoa(ans)},
		)
	default:
		ans := n - rank
		return linalgBridgeProblem(o,
			fmt.Sprintf("A %d by %d matrix has rank %d. What is its nullity? (enter a number)", m, n, rank),
			strconv.Itoa(ans),
			fmt.Sprintf("Nullity is the number of columns minus the rank: %d - %d = %d.", n, rank, ans),
			map[string]string{"kind": "nullity", "answer": strconv.Itoa(ans)},
		)
	}
}

type gramSchmidtGen struct{}

func (g *gramSchmidtGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	a := 1 + rand.Intn(5)
	b := 1 + rand.Intn(5)
	c := 1 + rand.Intn(6)
	d := 1 + rand.Intn(6)
	switch o {
	case 1:
		ans := c*a + d*b
		return linalgBridgeProblem(o,
			fmt.Sprintf("Let v1 = (%d, %d) and v2 = (%d, %d). What is v2 dot v1? (enter a number)", a, b, c, d),
			strconv.Itoa(ans),
			fmt.Sprintf("The dot product is %d x %d + %d x %d = %d, the numerator of the projection coefficient.", c, a, d, b, ans),
			map[string]string{"kind": "dot_product", "answer": strconv.Itoa(ans)},
		)
	case 2:
		ans := a*a + b*b
		return linalgBridgeProblem(o,
			fmt.Sprintf("Let v1 = (%d, %d). What is v1 dot v1, the denominator of the projection coefficient? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The squared length is %d^2 + %d^2 = %d.", a, b, ans),
			map[string]string{"kind": "squared_length", "answer": strconv.Itoa(ans)},
		)
	default:
		k := 1 + rand.Intn(4)
		ans := a*a + b*b
		return linalgBridgeProblem(o,
			fmt.Sprintf("Let u = (%d, %d) and v = %d u + (%d, %d), where the added vector is perpendicular to u. What is the squared length of the component of v orthogonal to u? (enter a number)", a, b, k, -b, a),
			strconv.Itoa(ans),
			fmt.Sprintf("The perpendicular component is (%d, %d), whose squared length is %d^2 + %d^2 = %d.", -b, a, b, a, ans),
			map[string]string{"kind": "orthogonal_norm", "answer": strconv.Itoa(ans)},
		)
	}
}

type transformMatrixGen struct{}

func (g *transformMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	a := 1 + rand.Intn(5)
	b := 1 + rand.Intn(5)
	c := 1 + rand.Intn(5)
	d := 1 + rand.Intn(5)
	switch o {
	case 1:
		if rand.Intn(2) == 0 {
			return linalgBridgeProblem(o,
				fmt.Sprintf("A linear map sends e1 to (%d, %d) and e2 to (%d, %d). What is the entry in row 1, column 1 of its standard matrix? (enter a number)", a, b, c, d),
				strconv.Itoa(a),
				fmt.Sprintf("The first column is the image of e1, so the top entry is %d.", a),
				map[string]string{"kind": "matrix_entry", "answer": strconv.Itoa(a)},
			)
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("A linear map sends e1 to (%d, %d) and e2 to (%d, %d). What is the entry in row 2, column 1 of its standard matrix? (enter a number)", a, b, c, d),
			strconv.Itoa(b),
			fmt.Sprintf("The first column is the image of e1, so its second entry is %d.", b),
			map[string]string{"kind": "matrix_entry", "answer": strconv.Itoa(b)},
		)
	case 2:
		x := 1 + rand.Intn(4)
		y := 1 + rand.Intn(4)
		first := a*x + b*y
		second := c*x + d*y
		if rand.Intn(2) == 0 {
			return linalgBridgeProblem(o,
				fmt.Sprintf("The matrix [[%d, %d], [%d, %d]] acts on (%d, %d). What is the first component of the image? (enter a number)", a, b, c, d, x, y),
				strconv.Itoa(first),
				fmt.Sprintf("The first component is %d x %d + %d x %d = %d.", a, x, b, y, first),
				map[string]string{"kind": "image_component", "answer": strconv.Itoa(first)},
			)
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("The matrix [[%d, %d], [%d, %d]] acts on (%d, %d). What is the second component of the image? (enter a number)", a, b, c, d, x, y),
			strconv.Itoa(second),
			fmt.Sprintf("The second component is %d x %d + %d x %d = %d.", c, x, d, y, second),
			map[string]string{"kind": "image_component", "answer": strconv.Itoa(second)},
		)
	default:
		if rand.Intn(2) == 0 {
			ans := a*d - b*c
			return linalgBridgeProblem(o,
				fmt.Sprintf("A linear map has standard matrix [[%d, %d], [%d, %d]]. What is its determinant? (enter a number)", a, b, c, d),
				strconv.Itoa(ans),
				fmt.Sprintf("The determinant is %d x %d - %d x %d = %d.", a, d, b, c, ans),
				map[string]string{"kind": "determinant", "answer": strconv.Itoa(ans)},
			)
		}
		ans := a + d
		return linalgBridgeProblem(o,
			fmt.Sprintf("A linear map has standard matrix [[%d, %d], [%d, %d]]. What is its trace? (enter a number)", a, b, c, d),
			strconv.Itoa(ans),
			fmt.Sprintf("The trace is the sum of the diagonal entries: %d + %d = %d.", a, d, ans),
			map[string]string{"kind": "trace", "answer": strconv.Itoa(ans)},
		)
	}
}

type eigenspaceGen struct{}

func (g *eigenspaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 1 + rand.Intn(4)
		b := 1 + rand.Intn(4)
		lambda := a
		if rand.Intn(2) == 0 {
			lambda = b
		}
		count := 0
		if a == lambda {
			count++
		}
		if b == lambda {
			count++
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("The diagonal matrix [[%d, 0], [0, %d]] has eigenvalue %d. What is the dimension of its eigenspace? (enter a number)", a, b, lambda),
			strconv.Itoa(count),
			fmt.Sprintf("The eigenvalue appears %d time%s on the diagonal, so the eigenspace has dimension %d.", count, plural(count), count),
			map[string]string{"kind": "eigenspace_dimension", "answer": strconv.Itoa(count)},
		)
	case 2:
		n := 3 + rand.Intn(3)
		rank := 1 + rand.Intn(n-1)
		ans := n - rank
		return linalgBridgeProblem(o,
			fmt.Sprintf("An %d by %d matrix has an eigenvalue lambda, and the matrix A - lambda I has rank %d. What is the geometric multiplicity of lambda? (enter a number)", n, n, rank),
			strconv.Itoa(ans),
			fmt.Sprintf("The eigenspace has dimension %d - %d = %d.", n, rank, ans),
			map[string]string{"kind": "geometric_multiplicity", "answer": strconv.Itoa(ans)},
		)
	default:
		a := 1 + rand.Intn(3)
		b := 1 + rand.Intn(3)
		p := 1 + rand.Intn(3)
		q := 1 + rand.Intn(2)
		if a == b {
			b = a%3 + 1
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("A 3 by 3 matrix has characteristic polynomial (t - %d)^%d (t - %d)^%d. What is the algebraic multiplicity of eigenvalue %d? (enter a number)", a, p, b, q, a),
			strconv.Itoa(p),
			fmt.Sprintf("The exponent on the factor (t - %d) is %d, so the algebraic multiplicity is %d.", a, p, p),
			map[string]string{"kind": "algebraic_multiplicity", "answer": strconv.Itoa(p)},
		)
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

type homogeneousSystemsGen struct{}

func (g *homogeneousSystemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 1 + rand.Intn(5)
		b := 1 + rand.Intn(5)
		c := 1 + rand.Intn(5)
		d := 1 + rand.Intn(5)
		det := a*d - b*c
		return linalgBridgeProblem(o,
			fmt.Sprintf("For A = [[%d, %d], [%d, %d]], the system A x = 0 has a nontrivial solution exactly when det A = 0. What is det A? (enter a number)", a, b, c, d),
			strconv.Itoa(det),
			fmt.Sprintf("det A = (%d)(%d) - (%d)(%d) = %d.", a, d, b, c, det),
			map[string]string{"kind": "determinant", "answer": strconv.Itoa(det)},
		)
	case 2:
		n := 2 + rand.Intn(4)
		rank := 1 + rand.Intn(n)
		nullity := n - rank
		return linalgBridgeProblem(o,
			fmt.Sprintf("A homogeneous system has %d variables and coefficient matrix of rank %d. What is the dimension of its solution space? (enter a number)", n, rank),
			strconv.Itoa(nullity),
			fmt.Sprintf("The nullity is n - rank = %d - %d = %d.", n, rank, nullity),
			map[string]string{"kind": "nullity", "answer": strconv.Itoa(nullity)},
		)
	default:
		k := 1 + rand.Intn(5)
		t := 2 + rand.Intn(4)
		ans := -k * t
		return linalgBridgeProblem(o,
			fmt.Sprintf("The homogeneous system x + %d y = 0 has free variable y. If y = %d, what is x? (enter a number)", k, t),
			strconv.Itoa(ans),
			fmt.Sprintf("Substitute y = %d: x + %d(%d) = 0, so x = %d.", t, k, t, ans),
			map[string]string{"kind": "free_variable", "answer": strconv.Itoa(ans)},
		)
	}
}

type leontiefGen struct{}

func linalgDecimal4(numer int) string {
	s := fmt.Sprintf("%d.%04d", numer/10000, numer%10000)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}

func (g *leontiefGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := linalgBridgeOrdinal(ctx)
	switch o {
	case 1:
		a := 10 + rand.Intn(31)
		b := 10 + rand.Intn(31)
		ans := a + b
		return linalgBridgeProblem(o,
			fmt.Sprintf("A two-sector technology matrix has first row [%d%%, %d%%]. What percentage of sector 1 output is consumed internally per unit of total production? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("Add the first-row coefficients: %d + %d = %d percent.", a, b, ans),
			map[string]string{"kind": "row_sum", "answer": strconv.Itoa(ans)},
		)
	case 2:
		a := 10 + 10*rand.Intn(5)
		d := 20 + 20*rand.Intn(5)
		val := float64(100*d) / float64(100-a)
		ans := strconv.FormatFloat(val, 'f', 2, 64)
		if len(ans) > 3 && ans[len(ans)-3:] == ".00" {
			ans = ans[:len(ans)-3]
		}
		return linalgBridgeProblem(o,
			fmt.Sprintf("A one-sector economy consumes %d%% of its own output and must meet final demand %d. What total output is required? (enter a number)", a, d),
			ans,
			fmt.Sprintf("Output is d/(1-a) = %d/(1-%d/100) = %s.", d, a, ans),
			map[string]string{"kind": "one_sector_output", "answer": ans},
		)
	default:
		a := 10 + rand.Intn(31)
		b := 10 + rand.Intn(31)
		c := 10 + rand.Intn(31)
		d := 10 + rand.Intn(31)
		numer := (100-a)*(100-d) - b*c
		ans := linalgDecimal4(numer)
		return linalgBridgeProblem(o,
			fmt.Sprintf("A two-sector technology matrix is A = [[%d%%, %d%%], [%d%%, %d%%]]. What is det(I - A)? (enter a decimal)", a, b, c, d),
			ans,
			fmt.Sprintf("det(I - A) = (1 - %d/100)(1 - %d/100) - (%d/100)(%d/100) = %s.", a, d, b, c, ans),
			map[string]string{"kind": "determinant_i_minus_a", "answer": ans},
		)
	}
}
