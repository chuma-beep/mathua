package linalg

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

var matrixRe = regexp.MustCompile(`\[\[([^\]]+)\],\[([^\]]+)\]\]`)

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
	reg.Register("linalg.sys.cramer", &cramerGen{})
	reg.Register("linalg.eigen.concept", &eigenConceptGen{})
	reg.Register("linalg.eigen.compute", &eigenComputeGen{})
	reg.Register("linalg.lintrans.concept", &transformationsGen{})
	reg.Register("linalg.vec.span", &spanGen{})
	reg.Register("linalg.vec.basis", &basisGen{})

	reg.Register("linalg.eigen.diagonalization", &diagonalizationGen{})
	reg.Register("linalg.matrix.rank", &rankGen{})
	reg.Register("linalg.vector.cosine_similarity", &cosineSimilarityGen{})
	reg.Register("linalg.vector.parametric", &parametricGen{})
	reg.Register("linalg.inner_product", &innerProductGen{})
	reg.Register("linalg.decomp.lu", &luGen{})
	reg.Register("linalg.decomp.qr", &qrGen{})
	reg.Register("linalg.decomp.svd", &svdGen{})
	reg.Register("linalg.decomp.jordan", &jordanGen{})
	reg.Register("linalg.subspace", &subspaceGen{})
	registerFoundations(reg)
	registerBatch(reg)
}

func containsWord(q, word string) bool {
	fields := strings.FieldsFunc(q, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-'
	})
	for _, f := range fields {
		if f == word {
			return true
		}
	}
	return false
}

// retryable reports whether another redraw is allowed. A generator guard that
// can never be satisfied must cost a weaker problem, never a request that never
// returns.
func retryable(attempt int) bool { return attempt < generator.MaxDrawAttempts }

type vectorConceptGen struct{}

// The concept has three knowledge points: the geometric picture, components
// and unit direction, and vector operations. Practice passes the current
// knowledge-point ordinal, so each rung must ask its own question rather than
// repeating the magnitude prompt from the first rung.
func (g *vectorConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 2:
		return vectorComponentProblem(ctx)
	case 3:
		return vectorOperationProblem(ctx)
	default:
		return vectorMagnitudeProblem(ctx)
	}
}

func vectorMagnitudeProblem(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var x, y, mag int
	var q, ex string
	for attempt := 0; ; attempt++ {
		m := smallInt(2, 4+tier*2)
		n := smallInt(1, m-1)
		k := smallInt(1, 2+tier)
		x = k * (m*m - n*n)
		y = k * (2 * m * n)
		mag = k * (m*m + n*n)
		if rand.Intn(2) == 0 {
			x, y = y, x
		}
		if tier > 0 && rand.Intn(2) == 0 {
			x = -x
		}
		if tier > 0 && rand.Intn(2) == 0 {
			y = -y
		}
		ansStr := strconv.Itoa(mag)
		q = fmt.Sprintf("What is the magnitude of vector \\((%d,%d)\\)?", x, y)
		if mag == x || mag == -x || mag == y || mag == -y || containsWord(q, ansStr) {
			if retryable(attempt) {
				continue
			}
		}
		ex = fmt.Sprintf("Subgoal 1: State the Euclidean vector magnitude definition \\(\\|(x,y)\\| = \\sqrt{x^2 + y^2}\\).\n"+
			"Subgoal 2: Substitute components \\(x = %d\\) and \\(y = %d\\): \\(\\sqrt{(%d)^2 + (%d)^2}\\).\n"+
			"Subgoal 3: Evaluate squares and root: \\(\\sqrt{%d + %d} = \\sqrt{%d} = %d\\).",
			x, y, x, y, x*x, y*y, x*x+y*y, mag)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(mag),
		Explanation: ex,
		Facts: map[string]string{
			"x":      strconv.Itoa(x),
			"y":      strconv.Itoa(y),
			"answer": strconv.Itoa(mag),
		},
	}
}

// vectorComponentProblem asks for one component after the direction vector is
// scaled to a given magnitude. It keeps the answer an integer while exercising
// the unit-direction subgoal from the components knowledge point.
func vectorComponentProblem(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		m := smallInt(2, 4+tier*2)
		n := smallInt(1, m-1)
		k := smallInt(1, 2+tier)
		x := k * (m*m - n*n)
		y := k * (2 * m * n)
		mag := k * (m*m + n*n)
		if rand.Intn(2) == 0 {
			x, y = y, x
		}
		if tier > 0 && rand.Intn(2) == 0 {
			x = -x
		}
		if tier > 0 && rand.Intn(2) == 0 {
			y = -y
		}
		scale := smallInt(2, 4+tier)
		wanted := x * scale
		question := fmt.Sprintf("A vector of magnitude \\(%d\\) points in the direction of \\((%d,%d)\\). What is its \\(x\\)-component?", mag*scale, x, y)
		if containsWord(question, strconv.Itoa(wanted)) {
			if retryable(attempt) {
				continue
			}
		}
		return generator.Problem{
			Question: question,
			Answer:   strconv.Itoa(wanted),
			Explanation: fmt.Sprintf("The direction vector \\((%d,%d)\\) has magnitude \\(%d\\). Scaling it to magnitude \\(%d\\) multiplies each component by \\(%d/%d = %d\\), so the \\(x\\)-component is \\(%d \\cdot %d = %d\\).",
				x, y, mag, mag*scale, mag*scale, mag, scale, scale, x, wanted),
			Facts: map[string]string{
				"x":         strconv.Itoa(x),
				"y":         strconv.Itoa(y),
				"magnitude": strconv.Itoa(mag),
				"scale":     strconv.Itoa(scale),
				"answer":    strconv.Itoa(wanted),
			},
		}
	}
}

// vectorOperationProblem asks for one coordinate of a sum. The concept is
// graded as a number, so the question targets a single component rather than
// the tuple that the separate vector-addition concept uses.
func vectorOperationProblem(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	bound := 3 + tier*3
	for attempt := 0; ; attempt++ {
		a := smallInt(1, bound)
		b := smallInt(1, bound)
		c := smallInt(1, bound)
		d := smallInt(1, bound)
		if rand.Intn(2) == 0 {
			a = -a
		}
		if rand.Intn(2) == 0 {
			b = -b
		}
		if rand.Intn(2) == 0 {
			c = -c
		}
		if rand.Intn(2) == 0 {
			d = -d
		}
		answer := a + c
		question := fmt.Sprintf("Given \\(\\mathbf{u}=(%d,%d)\\) and \\(\\mathbf{v}=(%d,%d)\\), what is the \\(x\\)-component of \\(\\mathbf{u}+\\mathbf{v}\\)?", a, b, c, d)
		if containsWord(question, strconv.Itoa(answer)) {
			if retryable(attempt) {
				continue
			}
		}
		return generator.Problem{
			Question:    question,
			Answer:      strconv.Itoa(answer),
			Explanation: fmt.Sprintf("Add the \\(x\\)-components: \\(%d + (%d) = %d\\). The \\(y\\)-components add independently: \\(%d + (%d) = %d\\).", a, c, answer, b, d, b+d),
			Facts: map[string]string{
				"u1":     strconv.Itoa(a),
				"u2":     strconv.Itoa(b),
				"v1":     strconv.Itoa(c),
				"v2":     strconv.Itoa(d),
				"answer": strconv.Itoa(answer),
			},
		}
	}
}

type vectorAddGen struct{}

func (g *vectorAddGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a, b := rand.Intn(max(1, scale*2))-scale, rand.Intn(max(1, scale*2))-scale
	c, d := rand.Intn(max(1, scale*2))-scale, rand.Intn(max(1, scale*2))-scale
	if rand.Intn(2) == 0 {
		sumX, sumY := a+c, b+d
		return generator.Problem{
			Question:    fmt.Sprintf("Add vectors: \\((%d,%d) + (%d,%d) =\\) ?", a, b, c, d),
			Answer:      fmt.Sprintf("(%d,%d)", sumX, sumY),
			Explanation: fmt.Sprintf("(%d,%d) + (%d,%d) = (%d,%d)", a, b, c, d, sumX, sumY),
		}
	}
	scalar := rand.Intn(max(1, scale)) + 2
	if rand.Intn(2) == 0 {
		scalar = -scalar
	}
	return generator.Problem{
		Question:    fmt.Sprintf("Compute scalar multiplication: \\(%d \\cdot (%d,%d) =\\) ?", scalar, a, b),
		Answer:      fmt.Sprintf("(%d,%d)", scalar*a, scalar*b),
		Explanation: fmt.Sprintf("%d * (%d,%d) = (%d,%d)", scalar, a, b, scalar*a, scalar*b),
	}
}

type vectorDotGen struct{}

func (g *vectorDotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var a, b, c, d, dot int
	var q, ex string
	for attempt := 0; ; attempt++ {
		bound := 4 + tier*3
		a = smallInt(1, bound)
		b = smallInt(1, bound)
		c = smallInt(1, bound)
		d = smallInt(1, bound)
		if tier > 0 && rand.Intn(2) == 0 {
			a = -a
		}
		if tier > 0 && rand.Intn(2) == 0 {
			b = -b
		}
		if tier > 0 && rand.Intn(2) == 0 {
			c = -c
		}
		if tier > 0 && rand.Intn(2) == 0 {
			d = -d
		}
		dot = a*c + b*d
		ansStr := strconv.Itoa(dot)
		q = fmt.Sprintf("Compute the dot product \\((%d,%d) \\cdot (%d,%d)\\).", a, b, c, d)
		if dot == a || dot == -a || dot == b || dot == -b || dot == c || dot == -c || dot == d || dot == -d || containsWord(q, ansStr) {
			if retryable(attempt) {
				continue
			}
		}
		ex = fmt.Sprintf("Subgoal 1: State the definition of the dot product in \\(\\mathbb{R}^2\\): \\((u_1, u_2) \\cdot (v_1, v_2) = u_1 v_1 + u_2 v_2\\).\n"+
			"Subgoal 2: Substitute vectors \\((%d,%d)\\) and \\((%d,%d)\\): \\((%d)(%d) + (%d)(%d)\\).\n"+
			"Subgoal 3: Compute products and sum: \\(%d + (%d) = %d\\).",
			a, b, c, d, a, c, b, d, a*c, b*d, dot)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      strconv.Itoa(dot),
		Explanation: ex,
		Facts: map[string]string{
			"u1":     strconv.Itoa(a),
			"u2":     strconv.Itoa(b),
			"v1":     strconv.Itoa(c),
			"v2":     strconv.Itoa(d),
			"answer": strconv.Itoa(dot),
		},
	}
}

type matrixConceptGen struct{}

func (g *matrixConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
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
			m[i][j] = rand.Intn(max(1, scale*3)) + 1
		}
	}
	r := rand.Intn(rows) + 1
	c := rand.Intn(cols) + 1
	answer := m[r-1][c-1]
	return generator.Problem{
		Question:    fmt.Sprintf("In the matrix \\(%s\\), what is the \\((%d,%d)\\) entry?", formatMatrix(m), r, c),
		Answer:      fmt.Sprintf("%d", answer),
		Explanation: fmt.Sprintf("The entry at row %d, column %d is %d.", r, c, answer),
	}
}

type matrixAddGen struct{}

// A matrix answer is written with brackets and commas. The concept is numeric,
// so the field offered digits only: the syntax existed in the accepted answer
// and nowhere in the interface.
func (g *matrixAddGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Matrix with rows in brackets, e.g. [[1,2],[3,4]]", InputMode: "text"}
}

func (g *matrixAddGen) Grade(expected, userAnswer string) grader.Result {
	return gradeMatrix(expected, userAnswer)
}

func (g *matrixAddGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := range a {
		a[i] = make([]int, 2)
		b[i] = make([]int, 2)
		for j := range a[i] {
			a[i][j] = rand.Intn(max(1, scale*2)) - scale
			b[i][j] = rand.Intn(max(1, scale*2)) - scale
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
			Question:    fmt.Sprintf("Add matrices: \\(%s + %s =\\) ?", formatMatrix(a), formatMatrix(b)),
			Answer:      formatMatrix(sum),
			Explanation: fmt.Sprintf("%s + %s = %s (add corresponding entries)", formatMatrix(a), formatMatrix(b), formatMatrix(sum)),
		}
	}
	scalar := rand.Intn(max(1, scale)) + 2
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
		Question:    fmt.Sprintf("Compute scalar multiplication: \\(%d \\cdot %s =\\) ?", scalar, formatMatrix(a)),
		Answer:      formatMatrix(prod),
		Explanation: fmt.Sprintf("%d * %s = %s", scalar, formatMatrix(a), formatMatrix(prod)),
	}
}

type matrixMultGen struct{}

func (g *matrixMultGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Matrix with rows in brackets, e.g. [[1,2],[3,4]]", InputMode: "text"}
}

func (g *matrixMultGen) Grade(expected, userAnswer string) grader.Result {
	return gradeMatrix(expected, userAnswer)
}

func (g *matrixMultGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a := make([][]int, 2)
	b := make([][]int, 2)
	for i := range a {
		a[i] = make([]int, 2)
		b[i] = make([]int, 2)
		for j := range a[i] {
			a[i][j] = rand.Intn(max(1, scale*2)) - scale
			b[i][j] = rand.Intn(max(1, scale*2)) - scale
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
		Question:    fmt.Sprintf("Multiply matrices: \\(%s \\times %s =\\) ?", formatMatrix(a), formatMatrix(b)),
		Answer:      formatMatrix(c),
		Explanation: fmt.Sprintf("%s × %s = %s", formatMatrix(a), formatMatrix(b), formatMatrix(c)),
	}
}

type matrixIdentityGen struct{}

func (g *matrixIdentityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Identity vector multiplication I_2 v = v
		v1 := mathutil.ScaleInt(ctx.Difficulty, 1, 10)
		v2 := mathutil.ScaleInt(ctx.Difficulty, 1, 10)
		if rand.Intn(2) == 0 {
			v1 = -v1
		}
		if rand.Intn(2) == 0 {
			v2 = -v2
		}
		comp := 1
		ans := v1
		if rand.Intn(2) == 0 {
			comp = 2
			ans = v2
		}
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(I_2\\) be the \\(2 \\times 2\\) identity matrix and \\(v = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\). Find component %d of the product vector \\(I_2 v\\).", v1, v2, comp),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the governing property of the identity matrix: \\(I_n v = v\\) for any vector \\(v\\).\n"+
				"Subgoal 2: Substitute vector components: \\(I_2 \\begin{pmatrix} %d \\\\ %d \\end{pmatrix} = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\).\n"+
				"Subgoal 3: Extract component %d, giving %d.", v1, v2, v1, v2, comp, ans),
		}
	case 1:
		// Determinant of scaled identity det(c * I_n) = c^n
		n := 2
		if ctx.Difficulty > 0.5 {
			n = 3
		}
		c := mathutil.ScaleInt(ctx.Difficulty, 2, 5)
		ans := 1
		for i := 0; i < n; i++ {
			ans *= c
		}
		return generator.Problem{
			Question: fmt.Sprintf("Find the determinant of the scaled \\(%d \\times %d\\) identity matrix \\(%d I_%d\\).", n, n, c, n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the determinant scaling identity for identity matrices: \\(\\det(c I_n) = c^n\\).\n"+
				"Subgoal 2: Substitute dimension \\(n = %d\\) and scalar \\(c = %d\\): \\((%d)^{%d}\\).\n"+
				"Subgoal 3: Evaluate the exponent to obtain %d.", n, c, c, n, ans),
		}
	case 2:
		// Trace of identity matrix tr(I_n) = n
		n := mathutil.ScaleInt(ctx.Difficulty, 2, 8)
		return generator.Problem{
			Question: fmt.Sprintf("Compute the trace of the \\(%d \\times %d\\) identity matrix \\(I_%d\\).", n, n, n),
			Answer:   fmt.Sprintf("%d", n),
			Explanation: fmt.Sprintf("Subgoal 1: State the trace definition as the sum of all diagonal elements: \\(\\operatorname{tr}(A) = \\sum_{i=1}^n a_{ii}\\).\n"+
				"Subgoal 2: Substitute the diagonal elements of \\(I_%d\\), each equal to 1: \\(\\sum_{i=1}^{%d} 1\\).\n"+
				"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(I_%d) = %d\\).", n, n, n, n),
		}
	case 3:
		// Unimodular matrix inverse entry
		rng := rand.New(rand.NewSource(rand.Int63()))
		mat := mathutil.RandomUnimodular2x2(rng, 1+int(ctx.Difficulty*2))
		if mat.Det() == -1 {
			mat.A = -mat.A
			mat.B = -mat.B
		}
		inv, _ := mat.Inverse()
		pos := rand.Intn(4)
		row, col, ans := 1, 1, inv.A
		switch pos {
		case 1:
			row, col, ans = 1, 2, inv.B
		case 2:
			row, col, ans = 2, 1, inv.C
		case 3:
			row, col, ans = 2, 2, inv.D
		}
		return generator.Problem{
			Question: fmt.Sprintf("For the invertible matrix \\(A = %s\\) with \\(\\det(A) = 1\\), what is the \\((%d,%d)\\) entry of its inverse \\(A^{-1}\\)?", mat.LaTeX(), row, col),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("Subgoal 1: State the \\(2 \\times 2\\) matrix inverse formula: \\(A^{-1} = \\frac{1}{\\det(A)} \\begin{pmatrix} d & -b \\\\ -c & a \\end{pmatrix}\\).\n"+
				"Subgoal 2: Substitute matrix parameters with \\(\\det(A) = 1\\) to construct \\(A^{-1} = %s\\).\n"+
				"Subgoal 3: Read off entry \\((%d,%d)\\), which equals %d.", inv.LaTeX(), row, col, ans),
		}
	default:
		// Determinant of inverse matrix reciprocal
		d := mathutil.ScaleInt(ctx.Difficulty, 2, 7)
		if rand.Intn(2) == 0 {
			d = -d
		}
		return generator.Problem{
			Question: fmt.Sprintf("If an invertible matrix \\(A\\) has \\(\\det(A) = %d\\), compute the value of \\(%d \\cdot \\det(A^{-1})\\).", d, d),
			Answer:   "1",
			Explanation: fmt.Sprintf("Subgoal 1: State the determinant reciprocal rule for invertible matrices: \\(\\det(A^{-1}) = \\frac{1}{\\det(A)}\\).\n"+
				"Subgoal 2: Substitute \\(\\det(A) = %d\\) into the identity: \\(\\det(A^{-1}) = \\frac{1}{%d}\\).\n"+
				"Subgoal 3: Multiply by %d: \\(%d \\cdot \\frac{1}{%d} = 1\\).", d, d, d, d, d),
		}
	}
}

type det2x2Gen struct{}

func (g *det2x2Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 3, 12)
	a := rand.Intn(2*span+1) - span
	b := rand.Intn(2*span+1) - span
	c := rand.Intn(2*span+1) - span
	d := rand.Intn(2*span+1) - span
	det := a*d - b*c
	return generator.Problem{
		Question: fmt.Sprintf("Find the determinant of the matrix \\(%s\\).", formatMatrix2([][]int{{a, b}, {c, d}})),
		Answer:   fmt.Sprintf("%d", det),
		Explanation: fmt.Sprintf("Subgoal 1: State the governing formula for a \\(2 \\times 2\\) determinant: \\(\\det(A) = ad - bc\\).\n"+
			"Subgoal 2: Substitute matrix parameters \\(a = %d, b = %d, c = %d, d = %d\\): \\((%d)(%d) - (%d)(%d)\\).\n"+
			"Subgoal 3: Evaluate products and subtract: %d - (%d) = %d.", a, b, c, d, a, d, b, c, a*d, b*c, det),
	}
}

type det3x3Gen struct{}

func (g *det3x3Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	span := mathutil.ScaleInt(ctx.Difficulty, 2, 6)
	for attempt := 0; ; attempt++ {
		m := make([][]int, 3)
		for i := range m {
			m[i] = make([]int, 3)
			for j := range m[i] {
				m[i][j] = rand.Intn(2*span+1) - span
			}
		}
		a, b, c := m[0][0], m[0][1], m[0][2]
		d, e, f := m[1][0], m[1][1], m[1][2]
		gVal, h, i := m[2][0], m[2][1], m[2][2]
		m11 := e*i - f*h
		m12 := d*i - f*gVal
		m13 := d*h - e*gVal
		det := a*m11 - b*m12 + c*m13
		ansStr := fmt.Sprintf("%d", det)
		matStr := formatMatrix(m)
		if det == 0 || containsWord(matStr, ansStr) {
			if retryable(attempt) {
				continue
			}
		}
		return generator.Problem{
			Question: fmt.Sprintf("Find the determinant of the matrix \\(%s\\).", matStr),
			Answer:   ansStr,
			Explanation: fmt.Sprintf("Subgoal 1: State Laplace cofactor expansion along the first row: \\(\\det(A) = a(ei - fh) - b(di - fg) + c(dh - eg)\\).\n"+
				"Subgoal 2: Substitute row entries and evaluate \\(2 \\times 2\\) minors: \\(M_{11} = %d, M_{12} = %d, M_{13} = %d\\).\n"+
				"Subgoal 3: Combine with alternating signs: \\(%d(%d) - (%d)(%d) + (%d)(%d) = %d\\).", m11, m12, m13, a, m11, b, m12, c, m13, det),
		}
	}
}

type systemsMatrixGen struct{}

// randomMatrixCoefficients returns a 2x2 coefficient matrix with a non-zero
// determinant. A singular matrix has either no solution or infinitely many
// solutions, so it cannot be represented by the single tuple this generator
// asks for. The two zero-row fallbacks keep the question non-degenerate, and
// the final fixed matrix guarantees termination even if the sampled range were
// to contain no invertible pair.
func randomMatrixCoefficients(scale int) (a, b, c, d int) {
	span := max(1, scale*2)
	for i := 0; i < 100; i++ {
		a, b = rand.Intn(span)-scale, rand.Intn(span)-scale
		c, d = rand.Intn(span)-scale, rand.Intn(span)-scale
		if a == 0 && c == 0 {
			a = 1
			c = 2
		}
		if b == 0 && d == 0 {
			b = 1
			d = 3
		}
		if a*d-b*c != 0 {
			return a, b, c, d
		}
	}
	return 1, 1, 1, 2
}

func (g *systemsMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	x := rand.Intn(max(1, scale*2)) - scale
	y := rand.Intn(max(1, scale*2)) - scale
	a, b, c, d := randomMatrixCoefficients(scale)
	e1 := a*x + b*y
	e2 := c*x + d*y
	return generator.Problem{
		Question:    fmt.Sprintf("Solve the linear system using matrices: \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\)", a, b, e1, c, d, e2),
		Answer:      fmt.Sprintf("(%d,%d)", x, y),
		Explanation: fmt.Sprintf("The solution is x=%d, y=%d.", x, y),
	}
}

type cramerGen struct{}

func (g *cramerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	rng := rand.New(rand.NewSource(rand.Int63()))
	for attempt := 0; attempt < 100; attempt++ {
		mat := mathutil.RandomUnimodular2x2(rng, 1+int(ctx.Difficulty*3))
		a, b := mat.A, mat.B
		c, d := mat.C, mat.D
		if b == 0 || c == 0 {
			continue
		}
		detA := a*d - b*c

		x := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		y := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		if rand.Intn(2) == 0 {
			x = -x
		}
		if rand.Intn(2) == 0 {
			y = -y
		}
		e1 := a*x + b*y
		e2 := c*x + d*y

		askX := rand.Intn(2) == 0
		if askX {
			detX := e1*d - b*e2
			valX := detX / detA
			ansStr := fmt.Sprintf("%d", valX)
			q := fmt.Sprintf("Use Cramer's rule to find \\(x\\) for the linear system: \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\).", a, b, e1, c, d, e2)
			if strings.Contains(q, ansStr) || valX == e1 || valX == e2 || valX == a || valX == b || valX == c || valX == d {
				continue
			}
			ex := fmt.Sprintf("Subgoal 1: State Cramer's rule for coordinate \\(x\\): \\(x = \\frac{\\det(A_x)}{\\det(A)}\\).\n"+
				"Subgoal 2: Compute coefficient determinant \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d\\) and modified determinant \\(\\det(A_x) = (%d)(%d) - (%d)(%d) = %d\\).\n"+
				"Subgoal 3: Divide determinants: \\(x = \\frac{%d}{%d} = %d\\).", a, d, b, c, detA, e1, d, b, e2, detX, detX, detA, valX)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
			}
		}
		detY := a*e2 - e1*c
		valY := detY / detA
		ansStr := fmt.Sprintf("%d", valY)
		q := fmt.Sprintf("Use Cramer's rule to find \\(y\\) for the linear system: \\(%dx + %dy = %d\\), \\(%dx + %dy = %d\\).", a, b, e1, c, d, e2)
		if strings.Contains(q, ansStr) || valY == e1 || valY == e2 || valY == a || valY == b || valY == c || valY == d {
			continue
		}
		ex := fmt.Sprintf("Subgoal 1: State Cramer's rule for coordinate \\(y\\): \\(y = \\frac{\\det(A_y)}{\\det(A)}\\).\n"+
			"Subgoal 2: Compute coefficient determinant \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d\\) and modified determinant \\(\\det(A_y) = (%d)(%d) - (%d)(%d) = %d\\).\n"+
			"Subgoal 3: Divide determinants: \\(y = \\frac{%d}{%d} = %d\\).", a, d, b, c, detA, a, e2, e1, c, detY, detY, detA, valY)
		return generator.Problem{
			Question:    q,
			Answer:      ansStr,
			Explanation: ex,
		}
	}
	return generator.Problem{
		Question:    "Use Cramer's rule to find \\(x\\) for the linear system: \\(2x + 1y = 7\\), \\(1x + 1y = 4\\).",
		Answer:      "3",
		Explanation: "det(A) = 1, det(Ax) = 3, so x = 3.",
	}
}

type eigenConceptGen struct{}

func (g *eigenConceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Av = lambda v
		lambda := smallInt(1, 5+int(ctx.Difficulty*3))
		k := smallInt(1, 3)
		b := smallInt(-2, 3)
		a := lambda - b*k
		d := smallInt(-2, 3)
		c := k * (lambda - d)
		return generator.Problem{
			Question:    fmt.Sprintf("The vector \\(v = \\begin{pmatrix} 1 \\\\ %d \\end{pmatrix}\\) is an eigenvector of matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\). Find the corresponding eigenvalue \\(\\lambda\\).", k, a, b, c, d),
			Answer:      fmt.Sprintf("%d", lambda),
			Explanation: fmt.Sprintf("An eigenvector satisfies \\(Av = \\lambda v\\). Computing the first component of \\(Av\\) gives \\(%d(1) + (%d)(%d) = %d\\). Therefore \\(\\lambda = %d\\).", a, b, k, lambda, lambda),
		}
	case 1:
		// tr(A) = lambda1 + lambda2, with A = P diag(lambda1, lambda2) P^-1 for
		// P = [[1,1],[1,2]]. Building A this way makes the stated lambda1 an
		// actual eigenvalue; arbitrary entries would make the premise false.
		l1 := smallInt(-4, 4)
		l2 := smallInt(-4, 4)
		a := 2*l1 - l2
		b := l2 - l1
		c := 2*l1 - 2*l2
		d := 2*l2 - l1
		tr := a + d
		return generator.Problem{
			Question:    fmt.Sprintf("The matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\) has one eigenvalue \\(\\lambda_1 = %d\\). Using the trace formula \\(\\operatorname{tr}(A) = \\lambda_1 + \\lambda_2\\), compute the second eigenvalue \\(\\lambda_2\\).", a, b, c, d, l1),
			Answer:      fmt.Sprintf("%d", l2),
			Explanation: fmt.Sprintf("The trace of a square matrix is the sum of its eigenvalues: \\(\\operatorname{tr}(A) = %d + (%d) = %d\\). Thus \\(\\lambda_2 = \\operatorname{tr}(A) - \\lambda_1 = %d - (%d) = %d\\).", a, d, tr, tr, l1, l2),
		}
	case 2:
		// det(A) = lambda1 * lambda2
		a := smallInt(-4, 4)
		d := smallInt(-4, 4)
		b := smallInt(-3, 3)
		c := smallInt(-3, 3)
		det := a*d - b*c
		return generator.Problem{
			Question:    fmt.Sprintf("For the matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\), what is the product of its eigenvalues \\(\\lambda_1 \\lambda_2\\)?", a, b, c, d),
			Answer:      fmt.Sprintf("%d", det),
			Explanation: fmt.Sprintf("The product of all eigenvalues of a matrix equals its determinant: \\(\\lambda_1 \\lambda_2 = \\det(A) = (%d)(%d) - (%d)(%d) = %d\\).", a, d, b, c, det),
		}
	case 3:
		// c * lambda
		lambda := smallInt(-5, 5)
		for lambda == 0 {
			lambda = smallInt(-5, 5)
		}
		c := smallInt(2, 5)
		ans := c * lambda
		return generator.Problem{
			Question:    fmt.Sprintf("If matrix \\(A\\) has an eigenvalue \\(\\lambda = %d\\), what is the corresponding eigenvalue of the scaled matrix \\(%d A\\)?", lambda, c),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("For any scalar \\(c\\), the eigenvalues of \\(cA\\) are \\(c\\lambda\\). Here \\(%d \\cdot (%d) = %d\\).", c, lambda, ans),
		}
	default:
		// lambda + c
		lambda := smallInt(-5, 5)
		c := smallInt(-5, 5)
		for c == 0 {
			c = smallInt(-5, 5)
		}
		ans := lambda + c
		return generator.Problem{
			Question:    fmt.Sprintf("If matrix \\(A\\) has an eigenvalue \\(\\lambda = %d\\), what is the corresponding eigenvalue of the shifted matrix \\(A + %d I\\)?", lambda, c),
			Answer:      fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf("If \\(Av = \\lambda v\\), then \\((A + %d I)v = Av + %d v = (%d + %d)v = %d v\\). The eigenvalue is %d.", c, c, lambda, c, ans, ans),
		}
	}
}

type eigenComputeGen struct{}

func (g *eigenComputeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	shears := [][4]int{
		{1, 1, 0, 1},
		{1, -1, 0, 1},
		{1, 0, 1, 1},
		{1, 0, -1, 1},
		{1, 2, 0, 1},
		{1, 0, 2, 1},
		{1, 3, 0, 1},
		{1, 0, 3, 1},
		{1, 0, -2, 1},
		{1, -2, 0, 1},
	}

	for attempt := 0; ; attempt++ {
		if ctx.Difficulty >= 0.4 {
			rng := rand.New(rand.NewSource(rand.Int63()))
			l1 := rand.Intn(5) - 3 // -3 to 1
			if l1 == 0 {
				l1 = -2
			}
			l2 := l1 + 1 + rand.Intn(3)
			if l2 == 0 {
				l2 = 1
			}
			l3 := l2 + 1 + rand.Intn(3)
			if l3 == 0 {
				l3 = 2
			}
			mat3, _, ok := mathutil.IntegerEigenvalues3x3(rng, l1, l2, l3, 2)
			if ok {
				ansStr := fmt.Sprintf("%d,%d,%d", l1, l2, l3)
				tr := mat3.Trace()
				det := mat3.Det()
				cof := mat3.CofactorMatrix()
				c1 := cof.Entry[0][0] + cof.Entry[1][1] + cof.Entry[2][2]
				return generator.Problem{
					Question: fmt.Sprintf("Find the eigenvalues of the matrix \\(%s\\). (Enter separated by commas)", mat3.LaTeX()),
					Answer:   ansStr,
					Explanation: fmt.Sprintf(
						"The characteristic polynomial is \\(p(\\lambda) = \\det(A - \\lambda I) = -\\lambda^3 + \\operatorname{tr}(A)\\lambda^2 - c_1\\lambda + \\det(A) = 0\\). "+
							"Here \\(\\operatorname{tr}(A) = %d\\), the principal minor sum is \\(c_1 = %d\\) and \\(\\det(A) = %d\\). "+
							"Factoring yields \\(-(\\lambda - %d)(\\lambda - %d)(\\lambda - %d) = 0\\), giving eigenvalues %d, %d and %d.",
						tr, c1, det, l1, l2, l3, l1, l2, l3,
					),
				}
			}
		}

		// Pick distinct integer eigenvalues
		l1 := rand.Intn(7) - 3 // -3 to 3
		if l1 == 0 {
			l1 = -2
		}
		l2 := l1 + rand.Intn(4) + 1
		if l2 == 0 {
			l2 = 2
		}
		if l1 == l2 {
			l2 = l1 + 1
		}
		if l1 > l2 {
			l1, l2 = l2, l1
		}

		s1 := shears[rand.Intn(len(shears))]
		s2 := shears[rand.Intn(len(shears))]
		p := s1[0]*s2[0] + s1[1]*s2[2]
		q := s1[0]*s2[1] + s1[1]*s2[3]
		r := s1[2]*s2[0] + s1[3]*s2[2]
		s := s1[2]*s2[1] + s1[3]*s2[3]

		// A = P * D * P^{-1} with det(P) = 1
		a11 := p*s*l1 - q*r*l2
		a12 := p * q * (l2 - l1)
		a21 := r * s * (l1 - l2)
		a22 := -q*r*l1 + p*s*l2

		// Reject triangular matrices or matrices where eigenvalues match any entry
		if a12 == 0 || a21 == 0 {
			if retryable(attempt) {
				continue
			}
		}
		if a11 == l1 || a11 == l2 || a22 == l1 || a22 == l2 || a12 == l1 || a12 == l2 || a21 == l1 || a21 == l2 {
			if retryable(attempt) {
				continue
			}
		}
		if mathutil.Abs(a11) > 25 || mathutil.Abs(a22) > 25 || mathutil.Abs(a12) > 25 || mathutil.Abs(a21) > 25 {
			if retryable(attempt) {
				continue
			}
		}

		ansStr := fmt.Sprintf("%d,%d", l1, l2)
		matStr := formatMatrix2([][]int{{a11, a12}, {a21, a22}})
		if strings.Contains(matStr, ansStr) {
			if retryable(attempt) {
				continue
			}
		}

		tr := a11 + a22
		det := a11*a22 - a12*a21

		return generator.Problem{
			Question: fmt.Sprintf("Find the eigenvalues of the matrix \\(%s\\). (Enter from least to greatest separated by a comma)", matStr),
			Answer:   ansStr,
			Explanation: fmt.Sprintf(
				"The characteristic equation is \\(\\det(A - \\lambda I) = \\lambda^2 - \\operatorname{tr}(A)\\lambda + \\det(A) = 0\\). "+
					"Here \\(\\operatorname{tr}(A) = %d + %d = %d\\) and \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d\\). "+
					"Solving \\(\\lambda^2 - %d\\lambda + %d = (\\lambda - %d)(\\lambda - %d) = 0\\) gives eigenvalues %d and %d.",
				a11, a22, tr, a11, a22, a12, a21, det, tr, det, l1, l2, l1, l2,
			),
		}
	}
}

type transformationsGen struct{}

func (g *transformationsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*5)
	a, b := rand.Intn(max(1, scale*2))+1, rand.Intn(max(1, scale*2))+1
	x, y := rand.Intn(max(1, scale*2))-scale, rand.Intn(max(1, scale*2))-scale
	return generator.Problem{
		Question:    fmt.Sprintf("Under the linear transformation \\(T(x,y) = (%dx,%dy)\\), where does the point \\((%d,%d)\\) map to?", a, b, x, y),
		Answer:      fmt.Sprintf("(%d,%d)", a*x, b*y),
		Explanation: fmt.Sprintf("T(%d,%d) = (%d*%d, %d*%d) = (%d,%d)", x, y, a, x, b, y, a*x, b*y),
	}
}

type spanGen struct{}

func (g *spanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Dimension of span of two vectors in R^2
			dependent := rand.Intn(2) == 0
			if dependent {
				v1 := smallInt(2, 5+tier*2)
				v2 := smallInt(2, 5+tier*2)
				m := smallInt(2, 3+tier)
				w1 := m * v1
				w2 := m * v2
				q := fmt.Sprintf("What is the dimension of the subspace \\(\\operatorname{span}\\{(%d,%d), (%d,%d)\\}\\) in \\(\\mathbb{R}^2\\)?", v1, v2, w1, w2)
				ex := fmt.Sprintf("Subgoal 1: State the linear dependence criterion for span in \\(\\mathbb{R}^2\\): two non-zero vectors span a 1-dimensional line if they are proportional.\n"+
					"Subgoal 2: Observe that vector \\((%d,%d)\\) equals %d times vector \\((%d,%d)\\).\n"+
					"Subgoal 3: Conclude that the vectors are linearly dependent, so the subspace dimension is 1.",
					w1, w2, m, v1, v2)
				return generator.Problem{
					Question:    q,
					Answer:      "1",
					Explanation: ex,
					Facts: map[string]string{
						"v1":     strconv.Itoa(v1),
						"v2":     strconv.Itoa(v2),
						"w1":     strconv.Itoa(w1),
						"w2":     strconv.Itoa(w2),
						"answer": "1",
					},
				}
			}
			v1 := smallInt(2, 6+tier*2)
			w2 := smallInt(2, 6+tier*2)
			q := fmt.Sprintf("What is the dimension of the subspace \\(\\operatorname{span}\\{(%d,0), (0,%d)\\}\\) in \\(\\mathbb{R}^2\\)?", v1, w2)
			ex := fmt.Sprintf("Subgoal 1: State the linear independence criterion for span in \\(\\mathbb{R}^2\\): two orthogonal non-zero vectors span the entire 2-dimensional space.\n"+
				"Subgoal 2: Compute the determinant of matrix \\(\\begin{pmatrix} %d & 0 \\\\ 0 & %d \\end{pmatrix}\\), giving \\((%d)(%d) - 0 = %d \\ne 0\\).\n"+
				"Subgoal 3: Conclude that the vectors are linearly independent, so the subspace dimension is 2.",
				v1, w2, v1, w2, v1*w2)
			return generator.Problem{
				Question:    q,
				Answer:      "2",
				Explanation: ex,
				Facts: map[string]string{
					"v1":     strconv.Itoa(v1),
					"w2":     strconv.Itoa(w2),
					"answer": "2",
				},
			}
		case 1:
			// Linear combination scalar c
			c := smallInt(2, 5+tier*3)
			x := smallInt(2, 6+tier*2)
			y := smallInt(2, 6+tier*2)
			cx := c * x
			cy := c * y
			ansStr := strconv.Itoa(c)
			q := fmt.Sprintf("Find the scalar \\(c\\) such that \\((%d,%d) = c(%d,%d)\\), showing that \\((%d,%d)\\) lies in \\(\\operatorname{span}\\{(%d,%d)\\}\\).", cx, cy, x, y, cx, cy, x, y)
			if c == x || c == y || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the coordinate-wise scalar multiplication requirement \\((%d,%d) = c(%d,%d)\\), requiring \\(%d = c(%d)\\) and \\(%d = c(%d)\\).\n"+
				"Subgoal 2: Solve the first coordinate equation for the scalar: \\(c = \\frac{%d}{%d}\\).\n"+
				"Subgoal 3: Evaluate the quotient to find \\(c = %d\\).",
				cx, cy, x, y, cx, x, cy, y, cx, x, c)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"cx":     strconv.Itoa(cx),
					"cy":     strconv.Itoa(cy),
					"x":      strconv.Itoa(x),
					"y":      strconv.Itoa(y),
					"answer": ansStr,
				},
			}
		case 2:
			// Parameter k making vectors linearly dependent
			a := smallInt(2, 5+tier*2)
			b := smallInt(2, 5+tier*2)
			m := smallInt(2, 4+tier)
			cVal := m * a
			k := m * b
			ansStr := strconv.Itoa(k)
			q := fmt.Sprintf("For what integer value of \\(k\\) do the vectors \\(v = (%d,%d)\\) and \\(w = (%d,k)\\) fail to span \\(\\mathbb{R}^2\\)?", a, b, cVal)
			if k == a || k == b || k == cVal || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the spanning condition: two vectors in \\(\\mathbb{R}^2\\) fail to span if and only if their determinant vanishes: \\((%d)(k) - (%d)(%d) = 0\\).\n"+
				"Subgoal 2: Substitute parameters into the linear equation: \\(%d k - %d = 0\\), so \\(%d k = %d\\).\n"+
				"Subgoal 3: Divide by coefficient to solve for \\(k\\): \\(k = \\frac{%d}{%d} = %d\\).",
				a, b, cVal, a, b*cVal, a, b*cVal, b*cVal, a, k)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"c":      strconv.Itoa(cVal),
					"answer": ansStr,
				},
			}
		case 3:
			// Dimension of 3 vectors with 0 in z coordinate
			x := smallInt(2, 7+tier*3)
			y := smallInt(2, 7+tier*3)
			q := fmt.Sprintf("Consider vectors \\(v_1 = (1,0,0)\\), \\(v_2 = (0,1,0)\\) and \\(v_3 = (%d,%d,0)\\). What is the dimension of \\(\\operatorname{span}\\{v_1, v_2, v_3\\}\\)?", x, y)
			ex := fmt.Sprintf("Subgoal 1: State subspace confinement: all three vectors have zero \\(z\\)-component and thus lie in the \\(xy\\)-plane.\n"+
				"Subgoal 2: Observe that \\(v_1 = (1,0,0)\\) and \\(v_2 = (0,1,0)\\) form an orthonormal basis for the \\(xy\\)-plane, while \\(v_3 = %d v_1 + %d v_2\\) is a linear combination of them.\n"+
				"Subgoal 3: Conclude that the span equals the 2-dimensional \\(xy\\)-plane, so the dimension is 2.",
				x, y)
			return generator.Problem{
				Question:    q,
				Answer:      "2",
				Explanation: ex,
				Facts: map[string]string{
					"x":      strconv.Itoa(x),
					"y":      strconv.Itoa(y),
					"answer": "2",
				},
			}
		default:
			// Span of k standard basis vectors in R^n: dimension of orthogonal complement
			k := smallInt(2, 4+tier*2)
			n := k + smallInt(2, 4+tier)
			ans := n - k
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("Let \\(S\\) be a subset of %d distinct standard unit basis vectors in \\(\\mathbb{R}^{%d}\\). What is the dimension of the orthogonal complement \\((\\operatorname{span}(S))^\\perp\\)?", k, n)
			if ans == k || ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the orthogonal complement dimension theorem: \\(\\dim(W^\\perp) = n - \\dim(W)\\) in \\(\\mathbb{R}^n\\).\n"+
				"Subgoal 2: Substitute ambient dimension \\(n = %d\\) and span dimension \\(\\dim(\\operatorname{span}(S)) = %d\\): %d - %d.\n"+
				"Subgoal 3: Evaluate dimension: %d.",
				n, k, n, k, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"k":      strconv.Itoa(k),
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
		}
	}
}

type basisGen struct{}

func (g *basisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Polynomial space P_n
			n := smallInt(2, 6+tier*4)
			ans := n + 1
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the vector space \\(\\mathcal{P}_{%d}\\) of real polynomials of degree at most %d?", n, n)
			if ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the standard monomial basis for \\(\\mathcal{P}_{%d}\\): \\(\\{1, x, x^2, \\dots, x^{%d}\\}\\).\n"+
				"Subgoal 2: Count the basis elements: %d degree terms plus 1 constant term: %d + 1.\n"+
				"Subgoal 3: Conclude that \\(\\dim(\\mathcal{P}_{%d}) = %d\\).",
				n, n, n, n, n, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
		case 1:
			// Matrix space M_{m x n}
			m := smallInt(2, 4+tier*2)
			n := smallInt(2, 4+tier*2)
			ans := m * n
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the vector space \\(\\mathcal{M}_{%d \\times %d}(\\mathbb{R})\\) of all %d-by-%d real matrices?", m, n, m, n)
			if ans == m || ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the standard matrix unit basis \\(E_{ij}\\) for \\(\\mathcal{M}_{%d \\times %d}(\\mathbb{R})\\).\n"+
				"Subgoal 2: Multiply the row and column dimensions: \\((%d) \\times (%d)\\).\n"+
				"Subgoal 3: Evaluate the product to find \\(\\dim = %d\\).",
				m, n, m, n, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"m":      strconv.Itoa(m),
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
		case 2:
			// Trace-zero diagonal matrices
			n := smallInt(3, 7+tier*3)
			ans := n - 1
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the subspace of all %d-by-%d diagonal matrices with trace zero in \\(\\mathcal{M}_{%d \\times %d}(\\mathbb{R})\\)?", n, n, n, n)
			if ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the trace constraint on diagonal matrices: an \\(n \\times n\\) diagonal matrix has \\(n\\) free diagonal entries and the trace condition \\(d_{11} + \\dots + d_{nn} = 0\\) imposes 1 linear restriction.\n"+
				"Subgoal 2: Subtract the linear restriction from total diagonal dimension: %d - 1.\n"+
				"Subgoal 3: Evaluate dimension: %d.",
				n, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
		case 3:
			// Coordinates with respect to an upper triangular basis
			a := smallInt(2, 4+tier)
			c1 := smallInt(2, 6+tier*3)
			c2 := smallInt(2, 5+tier*2)
			w1 := c1 + a*c2
			w2 := c2
			ansStr := strconv.Itoa(c1)
			q := fmt.Sprintf("Given the ordered basis \\(\\mathcal{B} = \\{(1,0), (%d,1)\\}\\), find the first coordinate \\(c_1\\) of the vector \\((%d,%d)\\) with respect to \\(\\mathcal{B}\\).", a, w1, w2)
			if c1 == a || c1 == w1 || c1 == w2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: Set up the vector equation \\(c_1(1,0) + c_2(%d,1) = (%d,%d)\\).\n"+
				"Subgoal 2: Solve the second component for \\(c_2 = %d\\) and substitute into the first equation: \\(c_1 + (%d)(%d) = %d\\).\n"+
				"Subgoal 3: Solve for the first coordinate: \\(c_1 = %d - %d = %d\\).",
				a, w1, w2, w2, a, w2, w1, w1, a*w2, c1)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"w1":     strconv.Itoa(w1),
					"w2":     strconv.Itoa(w2),
					"answer": ansStr,
				},
			}
		default:
			// Symmetric matrices Sym_n
			n := smallInt(2, 5+tier*2)
			ans := n * (n + 1) / 2
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the subspace of all %d-by-%d real symmetric matrices in \\(\\mathcal{M}_{%d \\times %d}(\\mathbb{R})\\)?", n, n, n, n)
			if ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the symmetric matrix dimension theorem: \\(\\dim(\\operatorname{Sym}_n) = \\frac{n(n+1)}{2}\\).\n"+
				"Subgoal 2: Substitute the matrix dimension \\(n = %d\\): \\(\\frac{(%d)(%d + 1)}{2} = \\frac{(%d)(%d)}{2}\\).\n"+
				"Subgoal 3: Evaluate the product and divide by 2: \\(\\frac{%d}{2} = %d\\).",
				n, n, n, n, n+1, n*(n+1), ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
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

func parseMatrix(s string) ([][]int, bool) {
	s = strings.TrimSpace(s)
	// Strip outer math delimiters
	s = strings.TrimPrefix(s, `\(`)
	s = strings.TrimSuffix(s, `\)`)
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimSuffix(s, "$")
	s = strings.TrimSpace(s)

	// Replace unicode minus and LaTeX minus commands
	s = strings.ReplaceAll(s, "−", "-")
	s = strings.ReplaceAll(s, `\minus`, "-")

	// 1. Check for LaTeX environments: \begin{pmatrix} ... \end{pmatrix}, etc.
	var content string
	for _, env := range []string{"pmatrix", "bmatrix", "matrix", "vmatrix", "Vmatrix"} {
		begin := `\begin{` + env + `}`
		end := `\end{` + env + `}`
		if idx := strings.Index(s, begin); idx != -1 {
			if endIdx := strings.Index(s, end); endIdx > idx {
				content = s[idx+len(begin) : endIdx]
				break
			}
		}
	}

	// Also check \left[ ... \right] or \left( ... \right)
	if content == "" {
		if strings.HasPrefix(s, `\left[`) && strings.HasSuffix(s, `\right]`) {
			content = s[len(`\left[`) : len(s)-len(`\right]`)]
		} else if strings.HasPrefix(s, `\left(`) && strings.HasSuffix(s, `\right)`) {
			content = s[len(`\left(`) : len(s)-len(`\right)`)]
		}
	}

	if content != "" {
		content = strings.ReplaceAll(content, `\cr`, `\\`)
		rawRows := strings.Split(content, `\\`)
		var rows [][]int
		for _, rStr := range rawRows {
			rStr = strings.TrimSpace(rStr)
			if rStr == "" {
				continue
			}
			var rawEntries []string
			if strings.Contains(rStr, "&") {
				rawEntries = strings.Split(rStr, "&")
			} else {
				rawEntries = strings.Split(rStr, ",")
			}
			var row []int
			for _, entry := range rawEntries {
				entry = strings.TrimSpace(entry)
				if entry == "" {
					continue
				}
				val, err := strconv.Atoi(entry)
				if err != nil {
					return nil, false
				}
				row = append(row, val)
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		if len(rows) > 0 {
			cols := len(rows[0])
			for _, r := range rows[1:] {
				if len(r) != cols {
					return nil, false
				}
			}
			return rows, true
		}
	}

	// 2. Check for bracketed format: [[1, 2], [3, 4]] or [[1, 2, 3], [4, 5, 6], [7, 8, 9]]
	if strings.Contains(s, "[") && strings.Contains(s, "]") {
		rowRegex := regexp.MustCompile(`\[([^\[\]]+)\]`)
		matches := rowRegex.FindAllStringSubmatch(s, -1)
		if len(matches) > 0 {
			var rows [][]int
			for _, m := range matches {
				part := m[1]
				parts := strings.Split(part, ",")
				if len(parts) == 1 && strings.Contains(part, " ") {
					parts = strings.Fields(part)
				}
				var row []int
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p == "" {
						continue
					}
					val, err := strconv.Atoi(p)
					if err != nil {
						return nil, false
					}
					row = append(row, val)
				}
				if len(row) > 0 {
					rows = append(rows, row)
				}
			}
			if len(rows) > 0 {
				cols := len(rows[0])
				for _, r := range rows[1:] {
					if len(r) != cols {
						return nil, false
					}
				}
				return rows, true
			}
		}
	}

	// 3. Raw LaTeX rows with & and \\ without begin/end: e.g. "1 & 2 \\ 3 & 4"
	if strings.Contains(s, `\\`) || strings.Contains(s, "&") {
		rawRows := strings.Split(s, `\\`)
		var rows [][]int
		for _, rStr := range rawRows {
			rStr = strings.TrimSpace(rStr)
			if rStr == "" {
				continue
			}
			rawEntries := strings.Split(rStr, "&")
			var row []int
			for _, entry := range rawEntries {
				entry = strings.TrimSpace(entry)
				if entry == "" {
					continue
				}
				val, err := strconv.Atoi(entry)
				if err != nil {
					return nil, false
				}
				row = append(row, val)
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
		if len(rows) > 0 {
			cols := len(rows[0])
			for _, r := range rows[1:] {
				if len(r) != cols {
					return nil, false
				}
			}
			return rows, true
		}
	}

	return nil, false
}

func gradeMatrix(expected, userAnswer string) grader.Result {
	eMat, ok := parseMatrix(expected)
	if !ok {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid matrix format"}
	}
	aMat, ok := parseMatrix(userAnswer)
	if !ok {
		return grader.Result{Correct: false, Score: 0, Feedback: "Invalid matrix format"}
	}
	if len(eMat) != len(aMat) || len(eMat[0]) != len(aMat[0]) {
		return grader.Result{Correct: false, Score: 0, Feedback: "Matrix dimension mismatch"}
	}
	for i := range eMat {
		for j := range eMat[i] {
			if eMat[i][j] != aMat[i][j] {
				return grader.Result{Correct: false, Score: 0, Feedback: fmt.Sprintf("Mismatch at (%d,%d)", i+1, j+1)}
			}
		}
	}
	return grader.Result{Correct: true, Score: 1}
}

// Diagonalization: algebraic properties of diagonalizable matrices PDP⁻¹
type diagonalizationGen struct{}

func (g *diagonalizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	rng := rand.New(rand.NewSource(rand.Int63()))

	for attempt := 0; attempt < 50; attempt++ {
		if ctx.Difficulty >= 0.4 {
			l1 := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
			l2 := l1 + 1 + rand.Intn(3)
			l3 := l2 + 1 + rand.Intn(3)
			if rand.Intn(2) == 0 {
				l1 = -l1
			}
			mat3, _, ok := mathutil.IntegerEigenvalues3x3(rng, l1, l2, l3, 2)
			if ok {
				mode := rand.Intn(4)
				switch mode {
				case 0:
					// Trace of diagonal matrix D
					tr := l1 + l2 + l3
					ansStr := fmt.Sprintf("%d", tr)
					q := fmt.Sprintf("The 3x3 matrix \\(A = %s\\) is diagonalizable as \\(A = PDP^{-1}\\) with \\(D\\) diagonal. What is the trace of \\(D\\)?", mat3.LaTeX())
					if strings.Contains(q, ansStr) {
						continue
					}
					return generator.Problem{
						Question: q,
						Answer:   ansStr,
						Explanation: fmt.Sprintf("Subgoal 1: State the similarity invariance of trace: \\(\\operatorname{tr}(D) = \\operatorname{tr}(A) = \\lambda_1 + \\lambda_2 + \\lambda_3\\).\n"+
							"Subgoal 2: Compute the trace of \\(A\\) by summing its diagonal entries: %d + %d + %d.\n"+
							"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(D) = %d\\).", mat3.Entry[0][0], mat3.Entry[1][1], mat3.Entry[2][2], tr),
					}
				case 1:
					// Determinant of diagonal matrix D
					det := l1 * l2 * l3
					ansStr := fmt.Sprintf("%d", det)
					q := fmt.Sprintf("The 3x3 matrix \\(A = %s\\) is diagonalizable as \\(A = PDP^{-1}\\) with \\(D\\) diagonal. What is the determinant of \\(D\\)?", mat3.LaTeX())
					if strings.Contains(q, ansStr) {
						continue
					}
					return generator.Problem{
						Question: q,
						Answer:   ansStr,
						Explanation: fmt.Sprintf("Subgoal 1: State the similarity invariance of determinant: \\(\\det(D) = \\det(A) = \\lambda_1 \\lambda_2 \\lambda_3\\).\n"+
							"Subgoal 2: Compute \\(\\det(A) = %d\\).\n"+
							"Subgoal 3: Conclude that \\(\\det(D) = %d\\).", det, det),
					}
				case 2:
					// Power trace tr(A^2)
					ans := l1*l1 + l2*l2 + l3*l3
					ansStr := fmt.Sprintf("%d", ans)
					q := fmt.Sprintf("Let \\(A = PDP^{-1}\\) be a 3x3 diagonalizable matrix with diagonal matrix \\(D = \\operatorname{diag}(%d, %d, %d)\\). Compute \\(\\operatorname{tr}(A^2)\\).", l1, l2, l3)
					if strings.Contains(q, ansStr) {
						continue
					}
					return generator.Problem{
						Question: q,
						Answer:   ansStr,
						Explanation: fmt.Sprintf("Subgoal 1: State the power trace formula: \\(\\operatorname{tr}(A^2) = \\operatorname{tr}(D^2) = \\lambda_1^2 + \\lambda_2^2 + \\lambda_3^2\\).\n"+
							"Subgoal 2: Substitute the eigenvalues: \\((%d)^2 + (%d)^2 + (%d)^2\\).\n"+
							"Subgoal 3: Evaluate squares and sum: %d + %d + %d = %d.", l1, l2, l3, l1*l1, l2*l2, l3*l3, ans),
					}
				default:
					// Largest eigenvalue
					larger := l1
					if l2 > larger {
						larger = l2
					}
					if l3 > larger {
						larger = l3
					}
					ansStr := fmt.Sprintf("%d", larger)
					q := fmt.Sprintf("The 3x3 matrix \\(A = %s\\) has distinct integer eigenvalues. Find its largest eigenvalue.", mat3.LaTeX())
					if strings.Contains(q, ansStr) {
						continue
					}
					return generator.Problem{
						Question: q,
						Answer:   ansStr,
						Explanation: fmt.Sprintf("Subgoal 1: Compute the characteristic polynomial \\(\\det(A - \\lambda I) = 0\\) to find eigenvalues \\(\\lambda_1 = %d, \\lambda_2 = %d\\) and \\(\\lambda_3 = %d\\).\n"+
							"Subgoal 2: Compare the eigenvalues to identify the maximum.\n"+
							"Subgoal 3: Select the largest eigenvalue: %d.", l1, l2, l3, larger),
					}
				}
			}
		}
	}

	for attempt2 := 0; attempt2 < 50; attempt2++ {
		l1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		l2 := l1 + mathutil.ScaleInt(ctx.Difficulty, 1, 4)
		if rand.Intn(2) == 0 {
			l1 = -l1
		}
		mat := mathutil.IntegerEigenvalues2x2(rng, l1, l2, 1+int(ctx.Difficulty*2))

		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Trace of diagonal matrix D
			tr := l1 + l2
			ansStr := fmt.Sprintf("%d", tr)
			q := fmt.Sprintf("The matrix \\(A = %s\\) is diagonalizable as \\(A = PDP^{-1}\\) with \\(D\\) diagonal. What is the trace of \\(D\\)?", mat.LaTeX())
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf("Subgoal 1: State the similarity invariance of trace: \\(\\operatorname{tr}(D) = \\operatorname{tr}(A) = \\lambda_1 + \\lambda_2\\).\n"+
					"Subgoal 2: Compute the trace of \\(A\\) by summing its diagonal entries: \\((%d) + (%d)\\).\n"+
					"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(D) = %d\\).", mat.A, mat.D, tr),
			}
		case 1:
			// Determinant of diagonal matrix D
			det := l1 * l2
			ansStr := fmt.Sprintf("%d", det)
			q := fmt.Sprintf("The matrix \\(A = %s\\) is diagonalizable as \\(A = PDP^{-1}\\) with \\(D\\) diagonal. What is the determinant of \\(D\\)?", mat.LaTeX())
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf("Subgoal 1: State the similarity invariance of determinant: \\(\\det(D) = \\det(A) = \\lambda_1 \\lambda_2\\).\n"+
					"Subgoal 2: Compute the determinant of \\(A\\): \\((%d)(%d) - (%d)(%d)\\).\n"+
					"Subgoal 3: Evaluate products to find \\(\\det(D) = %d\\).", mat.A, mat.D, mat.B, mat.C, det),
			}
		case 2:
			// Power trace tr(A^2)
			ans := l1*l1 + l2*l2
			ansStr := fmt.Sprintf("%d", ans)
			q := fmt.Sprintf("Let \\(A = PDP^{-1}\\) be diagonalizable with diagonal matrix \\(D = \\begin{pmatrix} %d & 0 \\\\ 0 & %d \\end{pmatrix}\\). Compute \\(\\operatorname{tr}(A^2)\\).", l1, l2)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf("Subgoal 1: State the power trace formula: \\(\\operatorname{tr}(A^2) = \\operatorname{tr}(D^2) = \\lambda_1^2 + \\lambda_2^2\\).\n"+
					"Subgoal 2: Substitute the eigenvalues \\(\\lambda_1 = %d\\) and \\(\\lambda_2 = %d\\): \\((%d)^2 + (%d)^2\\).\n"+
					"Subgoal 3: Evaluate squares and sum: %d + %d = %d.", l1, l2, l1, l2, l1*l1, l2*l2, ans),
			}
		case 3:
			// Larger eigenvalue of diagonalizable matrix
			larger := l1
			if l2 > larger {
				larger = l2
			}
			ansStr := fmt.Sprintf("%d", larger)
			q := fmt.Sprintf("The matrix \\(A = %s\\) has distinct integer eigenvalues. Find its largest eigenvalue.", mat.LaTeX())
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf("Subgoal 1: State the characteristic equation \\(\\det(A - \\lambda I) = \\lambda^2 - \\operatorname{tr}(A)\\lambda + \\det(A) = 0\\).\n"+
					"Subgoal 2: Substitute \\(\\operatorname{tr}(A) = %d\\) and \\(\\det(A) = %d\\) to solve for roots \\(\\lambda_1 = %d\\) and \\(\\lambda_2 = %d\\).\n"+
					"Subgoal 3: Select the maximum root: %d.", mat.Trace(), mat.Det(), l1, l2, larger),
			}
		default:
			// Shifted determinant det(A - c I)
			c := mathutil.ScaleInt(ctx.Difficulty, 1, 4)
			ans := (l1 - c) * (l2 - c)
			ansStr := fmt.Sprintf("%d", ans)
			q := fmt.Sprintf("A diagonalizable matrix \\(A\\) has eigenvalues \\(\\lambda_1 = %d\\) and \\(\\lambda_2 = %d\\). Compute \\(\\det(A - %d I)\\).", l1, l2, c)
			if strings.Contains(q, ansStr) {
				continue
			}
			return generator.Problem{
				Question: q,
				Answer:   ansStr,
				Explanation: fmt.Sprintf("Subgoal 1: State the spectral shift property: eigenvalues of \\(A - cI\\) are \\(\\lambda_1 - c\\) and \\(\\lambda_2 - c\\), so \\(\\det(A - cI) = (\\lambda_1 - c)(\\lambda_2 - c)\\).\n"+
					"Subgoal 2: Substitute eigenvalues and shift scalar: \\((%d - %d)(%d - %d)\\).\n"+
					"Subgoal 3: Multiply the shifted factors to obtain %d.", l1, c, l2, c, ans),
			}
		}
	}

	return generator.Problem{
		Question:    "Let \\(A = PDP^{-1}\\) be diagonalizable with diagonal matrix \\(D = \\begin{pmatrix} 2 & 0 \\\\ 0 & 5 \\end{pmatrix}\\). Compute \\(\\operatorname{tr}(A)\\).",
		Answer:      "7",
		Explanation: "The trace of D is 2 + 5 = 7.",
	}
}

// Rank of a matrix
type rankGen struct{}

func (g *rankGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	rng := rand.New(rand.NewSource(rand.Int63()))
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Full rank 2x2 matrix via RandomUnimodular2x2
		mat := mathutil.RandomUnimodular2x2(rng, 1+int(ctx.Difficulty*2))
		return generator.Problem{
			Question: fmt.Sprintf("What is the rank of matrix \\(%s\\)?", mat.LaTeX()),
			Answer:   "2",
			Explanation: fmt.Sprintf("Subgoal 1: State the full-rank criterion: an \\(n \\times n\\) matrix has rank \\(n\\) if and only if \\(\\det(A) \\ne 0\\).\n"+
				"Subgoal 2: Compute \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d \\ne 0\\).\n"+
				"Subgoal 3: Conclude that both rows are linearly independent, so \\(\\operatorname{rank}(A) = 2\\).", mat.A, mat.D, mat.B, mat.C, mat.Det()),
		}
	case 1:
		// Rank 1 outer product 2x2 matrix
		u1 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		u2 := mathutil.ScaleInt(ctx.Difficulty, 1, 5)
		k := mathutil.ScaleInt(ctx.Difficulty, 2, 4)
		if rand.Intn(2) == 0 {
			k = -k
		}
		m := [][]int{{u1, u2}, {k * u1, k * u2}}
		return generator.Problem{
			Question: fmt.Sprintf("What is the rank of matrix \\(%s\\)?", formatMatrix2(m)),
			Answer:   "1",
			Explanation: fmt.Sprintf("Subgoal 1: State the rank definition as the count of linearly independent rows.\n"+
				"Subgoal 2: Observe that row 2 is a scalar multiple of row 1 (ratio %d).\n"+
				"Subgoal 3: Conclude there is only 1 linearly independent row, so \\(\\operatorname{rank}(A) = 1\\).", k),
		}
	case 2:
		// Upper triangular 3x3 matrix with controlled pivots
		pivots := []int{mathutil.ScaleInt(ctx.Difficulty, 1, 4), mathutil.ScaleInt(ctx.Difficulty, 1, 4), mathutil.ScaleInt(ctx.Difficulty, 1, 4)}
		targetRank := 1 + rand.Intn(3)
		for i := targetRank; i < 3; i++ {
			pivots[i] = 0
		}
		m := [][]int{
			{pivots[0], smallInt(0, 3), smallInt(0, 3)},
			{0, pivots[1], smallInt(0, 3)},
			{0, 0, pivots[2]},
		}
		return generator.Problem{
			Question: fmt.Sprintf("What is the rank of the upper triangular matrix \\(%s\\)?", formatMatrix(m)),
			Answer:   fmt.Sprintf("%d", targetRank),
			Explanation: fmt.Sprintf("Subgoal 1: State the echelon rank property: for an upper triangular matrix, rank equals the count of nonzero diagonal pivots.\n"+
				"Subgoal 2: Inspect diagonal entries: %d, %d and %d.\n"+
				"Subgoal 3: Count the %d nonzero pivot(s) to determine \\(\\operatorname{rank}(A) = %d\\).", pivots[0], pivots[1], pivots[2], targetRank, targetRank),
		}
	case 3:
		// Rank-Nullity theorem application
		n := mathutil.ScaleInt(ctx.Difficulty, 3, 6)
		nullity := 1 + rand.Intn(n-1)
		rank := n - nullity
		return generator.Problem{
			Question: fmt.Sprintf("Let \\(T: \\mathbb{R}^{%d} \\to \\mathbb{R}^{%d}\\) be a linear transformation with nullity \\(\\dim(\\ker(T)) = %d\\). What is the rank \\(\\dim(\\operatorname{im}(T))\\)?", n, n, nullity),
			Answer:   fmt.Sprintf("%d", rank),
			Explanation: fmt.Sprintf("Subgoal 1: State the Rank-Nullity Theorem: \\(\\operatorname{rank}(T) + \\operatorname{nullity}(T) = n\\), where \\(n\\) is the domain dimension.\n"+
				"Subgoal 2: Substitute domain dimension \\(n = %d\\) and nullity = %d: \\(\\operatorname{rank}(T) + %d = %d\\).\n"+
				"Subgoal 3: Solve for rank: \\(%d - %d = %d\\).", n, nullity, nullity, n, n, nullity, rank),
		}
	default:
		// Collinear 3x3 rows
		a, b, cVal := smallInt(1, 3), smallInt(1, 3), smallInt(1, 3)
		m := [][]int{
			{a, b, cVal},
			{2 * a, 2 * b, 2 * cVal},
			{3 * a, 3 * b, 3 * cVal},
		}
		return generator.Problem{
			Question: fmt.Sprintf("What is the rank of matrix \\(%s\\)?", formatMatrix(m)),
			Answer:   "1",
			Explanation: fmt.Sprintf("Subgoal 1: Recall that the rank is the dimension of the subspace spanned by the row vectors.\n" +
				"Subgoal 2: Observe that row 2 = 2 \\times row 1 and row 3 = 3 \\times row 1, meaning all rows are scalar multiples of row 1.\n" +
				"Subgoal 3: Conclude the row space is 1-dimensional, so \\(\\operatorname{rank}(A) = 1\\)."),
		}
	}
}

// Cosine similarity between two vectors
type cosineSimilarityGen struct{}

func (g *cosineSimilarityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Parallel vectors with positive scalar k > 0
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		k := mathutil.ScaleInt(ctx.Difficulty, 2, 4)
		uStr := fmt.Sprintf("(%d,%d)", a, b)
		vStr := fmt.Sprintf("(%d,%d)", k*a, k*b)
		return generator.Problem{
			Question: fmt.Sprintf("Compute the cosine similarity between \\(u=%s\\) and \\(v=%s\\).", uStr, vStr),
			Answer:   "1",
			Explanation: fmt.Sprintf("Subgoal 1: State the cosine similarity formula: \\(\\cos(\\theta) = \\frac{u \\cdot v}{\\|u\\| \\|v\\|}\\).\n"+
				"Subgoal 2: Note \\(v = %d u\\) with positive scalar \\(k = %d\\), so \\(u\\) and \\(v\\) point in the identical direction (\\(\\theta = 0\\)).\n"+
				"Subgoal 3: Evaluate \\(\\cos(0) = 1\\).", k, k),
		}
	case 1:
		// Anti-parallel vectors with negative scalar k < 0
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		k := -mathutil.ScaleInt(ctx.Difficulty, 2, 4)
		uStr := fmt.Sprintf("(%d,%d)", a, b)
		vStr := fmt.Sprintf("(%d,%d)", k*a, k*b)
		return generator.Problem{
			Question: fmt.Sprintf("Compute the cosine similarity between \\(u=%s\\) and \\(v=%s\\).", uStr, vStr),
			Answer:   "-1",
			Explanation: fmt.Sprintf("Subgoal 1: State the cosine similarity formula: \\(\\cos(\\theta) = \\frac{u \\cdot v}{\\|u\\| \\|v\\|}\\).\n"+
				"Subgoal 2: Note \\(v = %d u\\) with negative scalar \\(k = %d\\), so \\(u\\) and \\(v\\) point in opposite directions (\\(\\theta = \\pi\\)).\n"+
				"Subgoal 3: Evaluate \\(\\cos(\\pi) = -1\\).", k, k),
		}
	case 2:
		// Orthogonal 2D vectors (a, b) and (-b, a)
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 8)
		uStr := fmt.Sprintf("(%d,%d)", a, b)
		vStr := fmt.Sprintf("(%d,%d)", -b, a)
		return generator.Problem{
			Question: fmt.Sprintf("Compute the cosine similarity between orthogonal vectors \\(u=%s\\) and \\(v=%s\\).", uStr, vStr),
			Answer:   "0",
			Explanation: fmt.Sprintf("Subgoal 1: State the cosine similarity formula: \\(\\cos(\\theta) = \\frac{u \\cdot v}{\\|u\\| \\|v\\|}\\).\n"+
				"Subgoal 2: Compute the dot product: \\(u \\cdot v = (%d)(%d) + (%d)(%d) = %d + (%d) = 0\\).\n"+
				"Subgoal 3: Divide numerator zero by the nonzero norms to obtain \\(\\cos(\\theta) = 0\\).", a, -b, b, a, -a*b, a*b),
		}
	case 3:
		// Orthogonal 3D vectors
		a := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		b := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		cVal := mathutil.ScaleInt(ctx.Difficulty, 1, 6)
		uStr := fmt.Sprintf("(%d,0,0)", a)
		vStr := fmt.Sprintf("(0,%d,%d)", b, cVal)
		return generator.Problem{
			Question: fmt.Sprintf("Compute the cosine similarity between vectors \\(u=%s\\) and \\(v=%s\\).", uStr, vStr),
			Answer:   "0",
			Explanation: fmt.Sprintf("Subgoal 1: State the cosine similarity formula: \\(\\cos(\\theta) = \\frac{u \\cdot v}{\\|u\\| \\|v\\|}\\).\n"+
				"Subgoal 2: Compute dot product: \\((%d)(0) + (0)(%d) + (0)(%d) = 0\\).\n"+
				"Subgoal 3: Evaluate \\(\\frac{0}{\\|u\\| \\|v\\|} = 0\\).", a, b, cVal),
		}
	default:
		// Identical standard unit vectors
		return generator.Problem{
			Question: "Compute the cosine similarity between unit vectors \\(u=(1,0)\\) and \\(v=(1,0)\\).",
			Answer:   "1",
			Explanation: "Subgoal 1: State the cosine similarity formula for identical unit vectors: \\(\\cos(\\theta) = \\frac{u \\cdot v}{\\|u\\| \\|v\\|}\\).\n" +
				"Subgoal 2: Compute \\(u \\cdot v = (1)(1) + (0)(0) = 1\\) and \\(\\|u\\| = 1, \\|v\\| = 1\\).\n" +
				"Subgoal 3: Evaluate \\(\\frac{1}{1 \\times 1} = 1\\).",
		}
	}
}

// Parametric form of a line
type parametricGen struct{}

// parametricEquation writes the coordinate equation with the sign of the
// direction component folded into the term. A bare "%d+%dt" produced
// "x=3+-2t" while the learner typed the equivalent "x=3-2t".
func parametricEquation(start, direction int) string {
	if direction < 0 {
		return fmt.Sprintf("%d-%dt", start, -direction)
	}
	return fmt.Sprintf("%d+%dt", start, direction)
}

func (g *parametricGen) AnswerFormat() generator.AnswerFormat {
	return generator.AnswerFormat{Hint: "Both equations, e.g. x=3-2t, y=2+t", InputMode: "text"}
}

func (g *parametricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	var x0, y0, dx, dy int
	var q, ans, ex string
	for attempt := 0; ; attempt++ {
		x0 = smallInt(1, 4+tier*3)
		y0 = smallInt(1, 4+tier*3)
		dx = smallInt(1, 3+tier*2)
		dy = smallInt(1, 3+tier*2)
		if tier > 0 && rand.Intn(2) == 0 {
			x0 = -x0
		}
		if tier > 0 && rand.Intn(2) == 0 {
			y0 = -y0
		}
		if rand.Intn(2) == 0 {
			dx = -dx
		}
		if rand.Intn(2) == 0 {
			dy = -dy
		}
		q = fmt.Sprintf("Find the parametric equations for the line through \\((%d,%d)\\) with direction vector \\((%d,%d)\\).", x0, y0, dx, dy)
		xEq := parametricEquation(x0, dx)
		yEq := parametricEquation(y0, dy)
		ans = fmt.Sprintf("x=%s,y=%s", xEq, yEq)
		ex = fmt.Sprintf("Subgoal 1: State the vector form of a line: \\(\\mathbf{r}(t) = \\mathbf{r}_0 + t \\mathbf{v}\\).\n"+
			"Subgoal 2: Substitute initial point \\((%d,%d)\\) and direction vector \\((%d,%d)\\): \\((x,y) = (%d,%d) + t(%d,%d)\\).\n"+
			"Subgoal 3: Separate into coordinate equations: \\(x = %s\\) and \\(y = %s\\).",
			x0, y0, dx, dy, x0, y0, dx, dy, xEq, yEq)
		break
	}
	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: ex,
		Facts: map[string]string{
			"x0":     strconv.Itoa(x0),
			"y0":     strconv.Itoa(y0),
			"dx":     strconv.Itoa(dx),
			"dy":     strconv.Itoa(dy),
			"answer": ans,
		},
	}
}

type innerProductGen struct{}

func (g *innerProductGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Weighted inner product
			w1 := smallInt(2, 4+tier)
			w2 := smallInt(2, 4+tier)
			u1 := smallInt(1, 4+tier*2)
			u2 := smallInt(1, 4+tier*2)
			v1 := smallInt(1, 4+tier*2)
			v2 := smallInt(1, 4+tier*2)
			ans := w1*u1*v1 + w2*u2*v2
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("Compute the weighted inner product \\(\\langle u, v \\rangle = %d u_1 v_1 + %d u_2 v_2\\) for \\(u = (%d,%d)\\) and \\(v = (%d,%d)\\).", w1, w2, u1, u2, v1, v2)
			if ans == w1 || ans == w2 || ans == u1 || ans == u2 || ans == v1 || ans == v2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the weighted inner product definition: \\(\\langle u, v \\rangle = %d u_1 v_1 + %d u_2 v_2\\).\n"+
				"Subgoal 2: Substitute components: \\((%d)(%d)(%d) + (%d)(%d)(%d)\\).\n"+
				"Subgoal 3: Evaluate products and sum: %d + %d = %d.",
				w1, w2, w1, u1, v1, w2, u2, v2, w1*u1*v1, w2*u2*v2, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"w1":     strconv.Itoa(w1),
					"w2":     strconv.Itoa(w2),
					"u1":     strconv.Itoa(u1),
					"u2":     strconv.Itoa(u2),
					"v1":     strconv.Itoa(v1),
					"v2":     strconv.Itoa(v2),
					"answer": ansStr,
				},
			}
		case 1:
			// Orthogonality parameter k
			a := smallInt(2, 5+tier*2)
			b := smallInt(2, 5+tier*2)
			m := smallInt(2, 4+tier)
			c := m * b
			k := -m * a
			ansStr := strconv.Itoa(k)
			q := fmt.Sprintf("For what integer value of \\(k\\) are the vectors \\(u = (%d,%d)\\) and \\(v = (%d,k)\\) orthogonal under the standard dot product?", a, b, c)
			if k == a || k == -a || k == b || k == -b || k == c || k == -c || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the orthogonality condition under dot product: \\(u \\cdot v = (%d)(%d) + (%d)(k) = 0\\).\n"+
				"Subgoal 2: Set up linear equation: %d + %d k = 0, so %d k = -%d.\n"+
				"Subgoal 3: Divide by coefficient to solve: \\(k = \\frac{-%d}{%d} = %d\\).",
				a, c, b, a*c, b, b, a*c, a*c, b, k)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"c":      strconv.Itoa(c),
					"answer": ansStr,
				},
			}
		case 2:
			// Matrix-induced inner product
			d1 := smallInt(2, 4+tier)
			d2 := smallInt(2, 4+tier)
			u1 := smallInt(1, 4+tier)
			u2 := smallInt(1, 4+tier)
			v1 := smallInt(1, 4+tier)
			v2 := smallInt(1, 4+tier)
			ans := d1*u1*v1 + d2*u2*v2
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("For the diagonal matrix \\(A = \\begin{pmatrix} %d & 0 \\\\ 0 & %d \\end{pmatrix}\\), compute the induced inner product \\(\\langle u, v \\rangle_A = u^T A v\\) for \\(u = (%d,%d)\\) and \\(v = (%d,%d)\\).", d1, d2, u1, u2, v1, v2)
			if ans == d1 || ans == d2 || ans == u1 || ans == u2 || ans == v1 || ans == v2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the matrix-induced inner product definition: \\(\\langle u, v \\rangle_A = u^T A v\\).\n"+
				"Subgoal 2: Expand quadratic form for diagonal \\(A\\): \\((%d)(%d)(%d) + (%d)(%d)(%d)\\).\n"+
				"Subgoal 3: Compute products and sum: %d + %d = %d.",
				d1, u1, v1, d2, u2, v2, d1*u1*v1, d2*u2*v2, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"d1":     strconv.Itoa(d1),
					"d2":     strconv.Itoa(d2),
					"u1":     strconv.Itoa(u1),
					"u2":     strconv.Itoa(u2),
					"v1":     strconv.Itoa(v1),
					"v2":     strconv.Itoa(v2),
					"answer": ansStr,
				},
			}
		case 3:
			// Polynomial inner product
			m := smallInt(2, 5+tier*2)
			a := 2 * m
			b := smallInt(2, 5+tier*2)
			ans := m * b
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("On the vector space of polynomials with inner product \\(\\langle f, g \\rangle = \\int_0^1 f(x) g(x) \\, dx\\), compute \\(\\langle %d x, %d \\rangle\\).", a, b)
			if ans == a || ans == b || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the polynomial inner product integral: \\(\\langle f, g \\rangle = \\int_0^1 (%d x)(%d) \\, dx = %d \\int_0^1 x \\, dx\\).\n"+
				"Subgoal 2: Evaluate the antiderivative at unit limits: \\(%d \\left[\\frac{x^2}{2}\\right]_0^1 = %d \\left(\\frac{1}{2}\\right)\\).\n"+
				"Subgoal 3: Compute product: \\(%d \\times \\frac{1}{2} = %d\\).",
				a, b, a*b, a*b, a*b, a*b, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"answer": ansStr,
				},
			}
		default:
			// Orthogonal projection coefficient
			v1 := smallInt(2, 4+tier)
			v2 := smallInt(2, 4+tier)
			vNormSq := v1*v1 + v2*v2
			m := smallInt(2, 5+tier*2)
			u1 := m * v1
			u2 := m * v2
			dot := u1*v1 + u2*v2
			ansStr := strconv.Itoa(m)
			q := fmt.Sprintf("Find the scalar projection coefficient \\(c = \\frac{\\langle u, v \\rangle}{\\langle v, v \\rangle}\\) of \\(u = (%d,%d)\\) onto \\(v = (%d,%d)\\) under the standard dot product.", u1, u2, v1, v2)
			if m == u1 || m == u2 || m == v1 || m == v2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the projection coefficient definition: \\(c = \\frac{\\langle u, v \\rangle}{\\|v\\|^2}\\).\n"+
				"Subgoal 2: Compute inner products: \\(\\langle u, v \\rangle = (%d)(%d) + (%d)(%d) = %d\\) and \\(\\|v\\|^2 = (%d)^2 + (%d)^2 = %d\\).\n"+
				"Subgoal 3: Divide inner products: \\(c = \\frac{%d}{%d} = %d\\).",
				u1, v1, u2, v2, dot, v1, v2, vNormSq, dot, vNormSq, m)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"u1":     strconv.Itoa(u1),
					"u2":     strconv.Itoa(u2),
					"v1":     strconv.Itoa(v1),
					"v2":     strconv.Itoa(v2),
					"answer": ansStr,
				},
			}
		}
	}
}

type luGen struct{}

func (g *luGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// det(A) = det(U)
			u11 := smallInt(2, 4+tier)
			u22 := smallInt(2, 4+tier)
			u33 := smallInt(2, 4+tier)
			ans := u11 * u22 * u33
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("An \\(n \\times n\\) matrix \\(A\\) has LU factorization \\(A = LU\\) where \\(L\\) is unit lower triangular and \\(U\\) has diagonal entries %d, %d and %d. Compute \\(\\det(A)\\).", u11, u22, u33)
			if ans == u11 || ans == u22 || ans == u33 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the determinant factorization property: \\(\\det(A) = \\det(L)\\det(U)\\), with \\(\\det(L) = 1\\) for unit triangular \\(L\\).\n"+
				"Subgoal 2: Substitute upper triangular diagonal entries into the product formula: \\((%d)(%d)(%d)\\).\n"+
				"Subgoal 3: Multiply the diagonal entries: \\(\\det(A) = %d \\times %d \\times %d = %d\\).",
				u11, u22, u33, u11, u22, u33, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"u11":    strconv.Itoa(u11),
					"u22":    strconv.Itoa(u22),
					"u33":    strconv.Itoa(u33),
					"answer": ansStr,
				},
			}
		case 1:
			// Forward solve Ly = b
			l21 := smallInt(2, 4+tier)
			b1 := smallInt(2, 5+tier*2)
			y2 := smallInt(2, 5+tier*2)
			b2 := l21*b1 + y2
			ansStr := strconv.Itoa(y2)
			q := fmt.Sprintf("In solving \\(Ax = b\\) via \\(A = LU\\), solve \\(Ly = b\\) for \\(y_2\\) where \\(L = \\begin{pmatrix} 1 & 0 \\\\ %d & 1 \\end{pmatrix}\\) and \\(b = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\).", l21, b1, b2)
			if y2 == l21 || y2 == b1 || y2 == b2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the forward substitution equations: \\(y_1 = b_1 = %d\\) and \\(%d y_1 + y_2 = b_2 = %d\\).\n"+
				"Subgoal 2: Substitute \\(y_1 = %d\\) into the second row equation: \\((%d)(%d) + y_2 = %d\\), so \\(%d + y_2 = %d\\).\n"+
				"Subgoal 3: Subtract to solve for \\(y_2\\): \\(y_2 = %d - %d = %d\\).",
				b1, l21, b2, b1, l21, b1, b2, l21*b1, b2, b2, l21*b1, y2)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"l21":    strconv.Itoa(l21),
					"b1":     strconv.Itoa(b1),
					"b2":     strconv.Itoa(b2),
					"answer": ansStr,
				},
			}
		case 2:
			// Back substitution Ux = y
			u11 := smallInt(2, 4+tier)
			u12 := smallInt(1, 3+tier)
			u22 := smallInt(2, 4+tier)
			x2 := smallInt(2, 5+tier*2)
			y2 := u22 * x2
			x1 := smallInt(2, 5+tier*2)
			y1 := u11*x1 + u12*x2
			ansStr := strconv.Itoa(x2)
			q := fmt.Sprintf("In the back substitution stage of solving \\(Ax = b\\), solve \\(Ux = y\\) for \\(x_2\\) where \\(U = \\begin{pmatrix} %d & %d \\\\ 0 & %d \\end{pmatrix}\\) and \\(y = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\).", u11, u12, u22, y1, y2)
			if x2 == u11 || x2 == u12 || x2 == u22 || x2 == y1 || x2 == y2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the second row equation in upper-triangular system: \\(%d x_2 = y_2 = %d\\).\n"+
				"Subgoal 2: Set up the division by pivot entry: \\(x_2 = \\frac{%d}{%d}\\).\n"+
				"Subgoal 3: Evaluate the quotient to obtain \\(x_2 = %d\\).",
				u22, y2, y2, u22, x2)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"u11":    strconv.Itoa(u11),
					"u12":    strconv.Itoa(u12),
					"u22":    strconv.Itoa(u22),
					"y1":     strconv.Itoa(y1),
					"y2":     strconv.Itoa(y2),
					"answer": ansStr,
				},
			}
		case 3:
			// Multiplier l_{21} in Gaussian elimination
			a11 := smallInt(2, 4+tier)
			m := smallInt(2, 4+tier)
			a21 := m * a11
			a12 := smallInt(2, 5+tier)
			a22 := smallInt(2, 6+tier*2)
			ansStr := strconv.Itoa(m)
			q := fmt.Sprintf("In computing the LU factorization \\(A = LU\\) of \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\) without pivoting, what is the entry \\(l_{21}\\) of the unit lower triangular factor \\(L\\)?", a11, a12, a21, a22)
			if m == a11 || m == a12 || m == a21 || m == a22 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the elimination row multiplier formula: \\(l_{21} = \\frac{a_{21}}{a_{11}}\\).\n"+
				"Subgoal 2: Substitute first column entries \\(a_{21} = %d\\) and pivot \\(a_{11} = %d\\): \\(\\frac{%d}{%d}\\).\n"+
				"Subgoal 3: Evaluate the quotient to find \\(l_{21} = %d\\).",
				a21, a11, a21, a11, m)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a11":    strconv.Itoa(a11),
					"a12":    strconv.Itoa(a12),
					"a21":    strconv.Itoa(a21),
					"a22":    strconv.Itoa(a22),
					"answer": ansStr,
				},
			}
		default:
			// Diagonal entry u_{22} in U
			l21 := smallInt(2, 4+tier)
			u11 := smallInt(2, 4+tier)
			u12 := smallInt(2, 4+tier)
			u22 := smallInt(2, 5+tier*2)
			a11 := u11
			a12 := u12
			a21 := l21 * u11
			a22 := l21*u12 + u22
			ansStr := strconv.Itoa(u22)
			q := fmt.Sprintf("In computing the LU factorization \\(A = LU\\) of \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\) where \\(L = \\begin{pmatrix} 1 & 0 \\\\ %d & 1 \\end{pmatrix}\\), what is the diagonal entry \\(u_{22}\\)?", a11, a12, a21, a22, l21)
			if u22 == a11 || u22 == a12 || u22 == a21 || u22 == a22 || u22 == l21 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the matrix product identity for entry \\((2,2)\\): \\(a_{22} = l_{21} u_{12} + u_{22}\\).\n"+
				"Subgoal 2: Substitute known matrix values: \\(%d = (%d)(%d) + u_{22} = %d + u_{22}\\).\n"+
				"Subgoal 3: Subtract to solve for \\(u_{22}\\): \\(u_{22} = %d - %d = %d\\).",
				a22, l21, u12, l21*u12, a22, l21*u12, u22)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a11":    strconv.Itoa(a11),
					"a12":    strconv.Itoa(a12),
					"a21":    strconv.Itoa(a21),
					"a22":    strconv.Itoa(a22),
					"l21":    strconv.Itoa(l21),
					"answer": ansStr,
				},
			}
		}
	}
}

type qrGen struct{}

func (g *qrGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		m := smallInt(2, 4+tier*2)
		n := smallInt(1, m-1)
		k := smallInt(1, 2+tier)
		a := k * (m*m - n*n)
		b := k * (2 * m * n)
		r := k * (m*m + n*n)
		if rand.Intn(2) == 0 {
			a, b = b, a
		}
		u := smallInt(2, 4+tier)
		v := smallInt(2, 4+tier)
		c := u*a - v*b
		d := u*b + v*a

		mode := rand.Intn(5)
		switch mode {
		case 0:
			ansStr := strconv.Itoa(r)
			q := fmt.Sprintf("The first column of real matrix \\(A\\) is \\(a_1 = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\). In the QR decomposition \\(A = QR\\), find the diagonal entry \\(r_{11} = \\|a_1\\|\\).", a, b)
			if r == a || r == -a || r == b || r == -b || strings.Contains(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the first column norm property in QR factorization: \\(r_{11} = \\|a_1\\| = \\sqrt{a_{11}^2 + a_{21}^2}\\).\n"+
				"Subgoal 2: Substitute column entries: \\(\\sqrt{(%d)^2 + (%d)^2} = \\sqrt{%d + %d}\\).\n"+
				"Subgoal 3: Evaluate square root to obtain diagonal entry: \\(\\sqrt{%d} = %d\\).",
				a, b, a*a, b*b, a*a+b*b, r)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"answer": ansStr,
				},
			}
		case 1:
			r12 := u * r
			ansStr := strconv.Itoa(r12)
			q := fmt.Sprintf("Consider matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\) with columns \\(a_1\\) and \\(a_2\\). In the QR factorization \\(A = QR\\), compute the upper entry \\(r_{12} = q_1^T a_2\\).", a, c, b, d)
			if r12 == a || r12 == -a || r12 == b || r12 == -b || r12 == c || r12 == -c || r12 == d || r12 == -d || strings.Contains(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the Gram-Schmidt projection formula for \\(r_{12}\\): \\(r_{12} = q_1 \\cdot a_2 = \\frac{1}{r_{11}} (a_1 \\cdot a_2)\\).\n"+
				"Subgoal 2: Compute inner product with first column norm \\(r_{11} = %d\\): \\(\\frac{1}{%d} ((%d)(%d) + (%d)(%d)) = \\frac{%d}{%d}\\).\n"+
				"Subgoal 3: Evaluate the quotient to obtain \\(r_{12} = %d\\).",
				r, r, a, c, b, d, a*c+b*d, r, r12)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"c":      strconv.Itoa(c),
					"d":      strconv.Itoa(d),
					"answer": ansStr,
				},
			}
		case 2:
			r22 := v * r
			ansStr := strconv.Itoa(r22)
			q := fmt.Sprintf("Consider matrix \\(A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}\\). In the QR decomposition \\(A = QR\\) with positive diagonal entries, find \\(r_{22}\\).", a, c, b, d)
			if r22 == a || r22 == -a || r22 == b || r22 == -b || r22 == c || r22 == -c || r22 == d || r22 == -d || strings.Contains(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the second diagonal entry definition: \\(r_{22} = \\|a_2 - r_{12} q_1\\|\\).\n"+
				"Subgoal 2: Compute orthogonal component \\(v_2 = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\) and take its Euclidean length: \\(\\sqrt{(%d)^2 + (%d)^2}\\).\n"+
				"Subgoal 3: Evaluate norm to obtain \\(r_{22} = %d\\).",
				-v*b, v*a, -v*b, v*a, r22)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"c":      strconv.Itoa(c),
					"d":      strconv.Itoa(d),
					"answer": ansStr,
				},
			}
		case 3:
			d1 := smallInt(2, 4+tier)
			d2 := smallInt(2, 4+tier)
			d3 := smallInt(2, 4+tier)
			detR := d1 * d2 * d3
			ansStr := strconv.Itoa(detR)
			q := fmt.Sprintf("A \\(3 \\times 3\\) matrix \\(A\\) has QR decomposition \\(A = QR\\), where \\(Q\\) is an orthogonal matrix with \\(\\det(Q) = 1\\) and \\(R\\) has diagonal entries %d, %d and %d. Find \\(\\det(A)\\).", d1, d2, d3)
			if detR == d1 || detR == d2 || detR == d3 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the determinant product formula: \\(\\det(A) = \\det(Q) \\det(R)\\).\n"+
				"Subgoal 2: Compute triangular determinant \\(\\det(R) = (%d)(%d)(%d)\\) and substitute \\(\\det(Q) = 1\\).\n"+
				"Subgoal 3: Multiply factors: \\(\\det(A) = 1 \\times (%d \\times %d \\times %d) = %d\\).",
				d1, d2, d3, d1, d2, d3, detR)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"d1":     strconv.Itoa(d1),
					"d2":     strconv.Itoa(d2),
					"d3":     strconv.Itoa(d3),
					"answer": ansStr,
				},
			}
		default:
			r11 := smallInt(2, 4+tier)
			r12 := smallInt(2, 4+tier)
			r22 := smallInt(2, 4+tier)
			x2 := smallInt(2, 5+tier*2)
			x1 := smallInt(2, 5+tier*2)
			c2 := r22 * x2
			c1 := r11*x1 + r12*x2
			ansStr := strconv.Itoa(x1)
			q := fmt.Sprintf("To solve the least-squares problem \\(Ax = b\\) via QR decomposition \\(Rx = Q^T b\\), let \\(R = \\begin{pmatrix} %d & %d \\\\ 0 & %d \\end{pmatrix}\\) and \\(Q^T b = \\begin{pmatrix} %d \\\\ %d \\end{pmatrix}\\). Find \\(x_1\\).", r11, r12, r22, c1, c2)
			if x1 == r11 || x1 == r12 || x1 == r22 || x1 == c1 || x1 == c2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the back substitution equations: \\(%d x_2 = %d\\) and \\(%d x_1 + %d x_2 = %d\\).\n"+
				"Subgoal 2: Solve the second row for \\(x_2 = \\frac{%d}{%d} = %d\\) and substitute into row 1: \\(%d x_1 + (%d)(%d) = %d\\).\n"+
				"Subgoal 3: Isolate and solve for \\(x_1\\): \\(%d x_1 = %d - %d = %d\\), so \\(x_1 = \\frac{%d}{%d} = %d\\).",
				r22, c2, r11, r12, c1, c2, r22, x2, r11, r12, x2, c1, r11, c1, r12*x2, c1-r12*x2, c1-r12*x2, r11, x1)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"r11":    strconv.Itoa(r11),
					"r12":    strconv.Itoa(r12),
					"r22":    strconv.Itoa(r22),
					"c1":     strconv.Itoa(c1),
					"c2":     strconv.Itoa(c2),
					"answer": ansStr,
				},
			}
		}
	}
}

type svdGen struct{}

func (g *svdGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Largest singular value sigma_1 of symmetric matrix A = [[a, b], [b, a]]
			// with a > b > 0. Gramian eigenvalues are (a+b)^2 and (a-b)^2, so sigma_1 = a+b.
			b := smallInt(1, 4+tier)
			a := b + smallInt(1, 4+tier)
			sig1 := a + b
			ansStr := strconv.Itoa(sig1)
			matStr := formatMatrix2([][]int{{a, b}, {b, a}})
			q := fmt.Sprintf("Find the largest singular value \\(\\sigma_1\\) of the matrix \\(A = %s\\).", matStr)
			if containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: Form the Gramian matrix \\(A^T A = %s\\).\n"+
				"Subgoal 2: Compute its eigenvalues \\(\\lambda_1 = (%d + %d)^2 = %d\\) and \\(\\lambda_2 = (%d - %d)^2 = %d\\).\n"+
				"Subgoal 3: Evaluate the principal square root \\(\\sigma_1 = \\sqrt{%d} = %d\\).",
				formatMatrix2([][]int{{a*a + b*b, 2 * a * b}, {2 * a * b, a*a + b*b}}),
				a, b, sig1*sig1, a, b, (a-b)*(a-b),
				sig1*sig1, sig1)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"answer": ansStr,
				},
			}
		case 1:
			// Spectral norm ||A||_2 of matrix A = [[a, b], [b, a]]
			b := smallInt(1, 3+tier)
			a := b + smallInt(1, 4+tier)
			ans := a + b
			ansStr := strconv.Itoa(ans)
			matStr := formatMatrix2([][]int{{a, b}, {b, a}})
			q := fmt.Sprintf("Find the spectral norm \\(\\|A\\|_2\\) of the matrix \\(A = %s\\).", matStr)
			if containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: Recall that the matrix 2-norm equals its maximum singular value: \\(\\|A\\|_2 = \\sigma_{\\max}(A)\\).\n"+
				"Subgoal 2: For symmetric \\(A = %s\\), the eigenvalues of \\(A^T A\\) are \\((%d + %d)^2 = %d\\) and \\((%d - %d)^2 = %d\\).\n"+
				"Subgoal 3: Take the principal square root of the maximum eigenvalue: \\(\\|A\\|_2 = %d\\).",
				matStr, a, b, ans*ans, a, b, (a-b)*(a-b), ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"answer": ansStr,
				},
			}
		case 2:
			// Squared Frobenius norm ||A||_F^2 = a^2 + b^2 + c^2 + d^2
			a := smallInt(1, 3+tier)
			b := smallInt(1, 3+tier)
			c := smallInt(1, 3+tier)
			d := smallInt(1, 3+tier)
			frobSq := a*a + b*b + c*c + d*d
			ansStr := strconv.Itoa(frobSq)
			matStr := formatMatrix2([][]int{{a, b}, {c, d}})
			q := fmt.Sprintf("Compute the squared Frobenius norm \\(\\|A\\|_F^2\\) of the matrix \\(A = %s\\).", matStr)
			if containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the Frobenius norm entrywise formula: \\(\\|A\\|_F^2 = \\sum_{i,j} a_{ij}^2\\).\n"+
				"Subgoal 2: Sum the squares of all matrix entries: \\((%d)^2 + (%d)^2 + (%d)^2 + (%d)^2 = %d + %d + %d + %d\\).\n"+
				"Subgoal 3: Evaluate the sum to find \\(\\|A\\|_F^2 = %d\\).",
				a, b, c, d, a*a, b*b, c*c, d*d, frobSq)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(a),
					"b":      strconv.Itoa(b),
					"c":      strconv.Itoa(c),
					"d":      strconv.Itoa(d),
					"answer": ansStr,
				},
			}
		case 3:
			// Condition number kappa_2(A) = sigma_max / sigma_min for A = [[a, b], [b, a]]
			// We select integer pairs (a, b) such that (a+b)/(a-b) = k exactly.
			type condPair struct{ a, b, k int }
			pairs := []condPair{
				{3, 1, 2},
				{6, 2, 2},
				{2, 1, 3},
				{4, 2, 3},
				{5, 3, 4},
				{3, 2, 5},
				{6, 4, 5},
			}
			p := pairs[rand.Intn(len(pairs))]
			ansStr := strconv.Itoa(p.k)
			matStr := formatMatrix2([][]int{{p.a, p.b}, {p.b, p.a}})
			q := fmt.Sprintf("Find the 2-norm condition number \\(\\kappa_2(A) = \\frac{\\sigma_{\\max}}{\\sigma_{\\min}}\\) of the matrix \\(A = %s\\).", matStr)
			if containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			sigMax := p.a + p.b
			sigMin := p.a - p.b
			ex := fmt.Sprintf("Subgoal 1: Find singular values of \\(A\\): \\(\\sigma_{\\max} = %d + %d = %d\\) and \\(\\sigma_{\\min} = %d - %d = %d\\).\n"+
				"Subgoal 2: Substitute into condition number formula: \\(\\kappa_2(A) = \\frac{%d}{%d}\\).\n"+
				"Subgoal 3: Evaluate the quotient to find \\(\\kappa_2(A) = %d\\).",
				p.a, p.b, sigMax, p.a, p.b, sigMin, sigMax, sigMin, p.k)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"a":      strconv.Itoa(p.a),
					"b":      strconv.Itoa(p.b),
					"answer": ansStr,
				},
			}
		default:
			// Rank / non-zero singular values of A
			isFullRank := rand.Intn(2) == 1
			var a, b, c, d, rank int
			if isFullRank {
				a = smallInt(1, 4)
				b = smallInt(1, 4)
				c = smallInt(1, 4)
				d = smallInt(1, 4)
				if a*d-b*c == 0 {
					d += 2
				}
				rank = 2
			} else {
				a = smallInt(1, 3)
				b = smallInt(1, 3)
				m := smallInt(2, 3)
				c = m * a
				d = m * b
				rank = 1
			}
			ansStr := strconv.Itoa(rank)
			matStr := formatMatrix2([][]int{{a, b}, {c, d}})
			q := fmt.Sprintf("How many non-zero singular values does the matrix \\(A = %s\\) have?", matStr)
			if containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the SVD rank theorem: the number of non-zero singular values equals the matrix rank \\(\\operatorname{rank}(A)\\).\n"+
				"Subgoal 2: Determine matrix determinant: \\(\\det(A) = (%d)(%d) - (%d)(%d) = %d\\).\n"+
				"Subgoal 3: Conclude that \\(\\operatorname{rank}(A) = %d\\), so there are %d non-zero singular values.",
				a, d, b, c, a*d-b*c, rank, rank)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"answer": ansStr,
				},
			}
		}
	}
}

type jordanGen struct{}

func (g *jordanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	rng := rand.New(rand.NewSource(rand.Int63()))

	for attempt := 0; ; attempt++ {
		if ctx.Difficulty >= 0.3 {
			lambda := smallInt(2, 4+tier)
			mu := lambda + 1 + rand.Intn(3)
			singleBlock := (rand.Intn(2) == 0)
			mat3, _, ok := mathutil.JordanMatrix3x3(rng, lambda, mu, singleBlock, 2)
			if ok {
				mode := rand.Intn(3)
				switch mode {
				case 0:
					// Determinant of non-triangular matrix with Jordan canonical form
					det := lambda * lambda * mu
					if singleBlock {
						det = lambda * lambda * lambda
					}
					ansStr := strconv.Itoa(det)
					q := fmt.Sprintf("Compute the determinant \\(\\det(A)\\) of the matrix \\(A = %s\\).", mat3.LaTeX())
					if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
						if retryable(attempt) {
							continue
						}
					}
					ex := fmt.Sprintf("Subgoal 1: State that the determinant is similarity-invariant: \\(\\det(A) = \\det(J)\\), the product of the diagonal eigenvalues.\n"+
						"Subgoal 2: Multiply the diagonal entries of the Jordan canonical form: \\((%d)(%d)(%d)\\).\n"+
						"Subgoal 3: Evaluate the product to find \\(\\det(A) = %d\\).",
						lambda, lambda, mu, det)
					return generator.Problem{
						Question:    q,
						Answer:      ansStr,
						Explanation: ex,
						Facts: map[string]string{
							"answer": ansStr,
						},
					}
				case 1:
					// Power trace tr(A^2)
					trSq := 2*lambda*lambda + mu*mu
					if singleBlock {
						trSq = 3 * lambda * lambda
					}
					ansStr := strconv.Itoa(trSq)
					q := fmt.Sprintf("For the matrix \\(A = %s\\) with Jordan canonical form eigenvalues \\(\\lambda_1 = %d, \\lambda_2 = %d\\) and \\(\\lambda_3 = %d\\), compute \\(\\operatorname{tr}(A^2)\\).", mat3.LaTeX(), lambda, lambda, mu)
					if singleBlock {
						q = fmt.Sprintf("For the matrix \\(A = %s\\) with eigenvalue \\(\\lambda = %d\\) of algebraic multiplicity 3, compute \\(\\operatorname{tr}(A^2)\\).", mat3.LaTeX(), lambda)
					}
					if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
						if retryable(attempt) {
							continue
						}
					}
					ex := fmt.Sprintf("Subgoal 1: State that the trace of \\(A^2\\) equals the sum of squares of its eigenvalues: \\(\\operatorname{tr}(A^2) = \\sum \\lambda_i^2\\).\n"+
						"Subgoal 2: Substitute the eigenvalues into the sum of squares.\n"+
						"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(A^2) = %d\\).", trSq)
					return generator.Problem{
						Question:    q,
						Answer:      ansStr,
						Explanation: ex,
						Facts: map[string]string{
							"answer": ansStr,
						},
					}
				default:
					// Trace of A
					tr := mat3.Trace()
					ansStr := strconv.Itoa(tr)
					q := fmt.Sprintf("Compute the trace \\(\\operatorname{tr}(A)\\) of the matrix \\(A = %s\\).", mat3.LaTeX())
					if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
						if retryable(attempt) {
							continue
						}
					}
					ex := fmt.Sprintf("Subgoal 1: State the trace property: \\(\\operatorname{tr}(A)\\) is the sum of the main diagonal entries.\n"+
						"Subgoal 2: Sum the diagonal entries: %d + %d + %d.\n"+
						"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(A) = %d\\).",
						mat3.Entry[0][0], mat3.Entry[1][1], mat3.Entry[2][2], tr)
					return generator.Problem{
						Question:    q,
						Answer:      ansStr,
						Explanation: ex,
						Facts: map[string]string{
							"answer": ansStr,
						},
					}
				}
			}
		}

		mode := rand.Intn(3)
		switch mode {
		case 0:
			// Trace of Jordan matrix
			l1 := smallInt(2, 4+tier)
			l2 := smallInt(3, 5+tier)
			l3 := smallInt(4, 6+tier)
			tr := l1 + l2 + l3
			ansStr := strconv.Itoa(tr)
			mat := [][]int{
				{l1, 1, 0},
				{0, l2, 1},
				{0, 0, l3},
			}
			matStr := formatMatrix3(mat)
			q := fmt.Sprintf("Compute the trace \\(\\operatorname{tr}(J)\\) of the upper-triangular Jordan form matrix \\(J = %s\\).", matStr)
			if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the trace property: \\(\\operatorname{tr}(J)\\) is the sum of all diagonal entries.\n"+
				"Subgoal 2: Sum the diagonal entries: %d + %d + %d.\n"+
				"Subgoal 3: Evaluate the sum to find \\(\\operatorname{tr}(J) = %d\\).",
				l1, l2, l3, tr)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"l1":     strconv.Itoa(l1),
					"l2":     strconv.Itoa(l2),
					"l3":     strconv.Itoa(l3),
					"answer": ansStr,
				},
			}
		case 1:
			// Determinant of Jordan matrix
			l1 := smallInt(2, 3+tier)
			l2 := smallInt(2, 3+tier)
			l3 := smallInt(3, 4+tier)
			detA := l1 * l2 * l3
			ansStr := strconv.Itoa(detA)
			mat := [][]int{
				{l1, 1, 0},
				{0, l2, 0},
				{0, 0, l3},
			}
			matStr := formatMatrix3(mat)
			q := fmt.Sprintf("Find the determinant \\(\\det(J)\\) of the upper-triangular matrix in Jordan canonical form \\(J = %s\\).", matStr)
			if strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the triangular determinant property: \\(\\det(J)\\) is the product of all diagonal entries.\n"+
				"Subgoal 2: Multiply diagonal entries: (%d)(%d)(%d).\n"+
				"Subgoal 3: Evaluate the product to find \\(\\det(J) = %d\\).",
				l1, l2, l3, detA)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"l1":     strconv.Itoa(l1),
					"l2":     strconv.Itoa(l2),
					"l3":     strconv.Itoa(l3),
					"answer": ansStr,
				},
			}
		default:
			// Defect in Jordan chain basis
			alg := smallInt(4, 6+tier*2)
			geom := smallInt(1, 2)
			defect := alg - geom
			eig := smallInt(2, 5+tier)
			ansStr := strconv.Itoa(defect)
			q := fmt.Sprintf("An eigenspace for eigenvalue \\(\\lambda = %d\\) has algebraic multiplicity %d and geometric multiplicity %d. How many generalized eigenvectors that are not ordinary eigenvectors are needed to complete a Jordan chain basis for this eigenspace?", eig, alg, geom)
			if defect == eig || defect == alg || defect == geom || strings.Contains(q, ansStr) || strings.Contains(q, strings.TrimPrefix(ansStr, "-")) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the defect formula for generalized eigenvectors: \\(\\text{defect} = \\text{algebraic multiplicity} - \\text{geometric multiplicity}\\).\n"+
				"Subgoal 2: Substitute multiplicities: %d - %d.\n"+
				"Subgoal 3: Evaluate the difference to find %d non-ordinary generalized eigenvector(s) needed.",
				alg, geom, defect)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"eig":    strconv.Itoa(eig),
					"alg":    strconv.Itoa(alg),
					"geom":   strconv.Itoa(geom),
					"answer": ansStr,
				},
			}
		}
	}
}

// ----- subspace -----

type subspaceGen struct{}

func (g *subspaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tier := mathutil.DifficultyTier(ctx.Difficulty)
	for attempt := 0; ; attempt++ {
		mode := rand.Intn(5)
		switch mode {
		case 0:
			// Grassmann's dimension formula: dim(U + W) = dim(U) + dim(W) - dim(U cap W)
			d12 := smallInt(1, 3+tier)
			d1 := d12 + smallInt(1, 3+tier*2)
			d2 := d12 + smallInt(1, 3+tier*2)
			ans := d1 + d2 - d12
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("Two subspaces \\(U\\) and \\(W\\) of a vector space satisfy \\(\\dim(U) = %d\\), \\(\\dim(W) = %d\\) and \\(\\dim(U \\cap W) = %d\\). What is \\(\\dim(U + W)\\)?", d1, d2, d12)
			if ans == d1 || ans == d2 || ans == d12 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State Grassmann's dimension formula: \\(\\dim(U + W) = \\dim(U) + \\dim(W) - \\dim(U \\cap W)\\).\n"+
				"Subgoal 2: Substitute dimensions: \\((%d) + (%d) - (%d)\\).\n"+
				"Subgoal 3: Evaluate arithmetic: %d + %d - %d = %d.",
				d1, d2, d12, d1, d2, d12, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"dimU":   strconv.Itoa(d1),
					"dimW":   strconv.Itoa(d2),
					"dimCap": strconv.Itoa(d12),
					"answer": ansStr,
				},
			}
		case 1:
			// Direct sum: dim(U \oplus W) = dim(U) + dim(W)
			d1 := smallInt(2, 5+tier*2)
			d2 := smallInt(2, 5+tier*2)
			ans := d1 + d2
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("If \\(V = U \\oplus W\\) is the direct sum of subspaces \\(U\\) and \\(W\\) with \\(\\dim(U) = %d\\) and \\(\\dim(W) = %d\\), what is \\(\\dim(V)\\)?", d1, d2)
			if ans == d1 || ans == d2 || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the direct sum dimension theorem: \\(\\dim(U \\oplus W) = \\dim(U) + \\dim(W)\\).\n"+
				"Subgoal 2: Substitute subspace dimensions: \\((%d) + (%d)\\).\n"+
				"Subgoal 3: Sum dimensions to find \\(\\dim(V) = %d\\).",
				d1, d2, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"dimU":   strconv.Itoa(d1),
					"dimW":   strconv.Itoa(d2),
					"answer": ansStr,
				},
			}
		case 2:
			// Orthogonal complement dimension: dim(W^\perp) = n - dim(W)
			n := smallInt(4, 7+tier*3)
			dimW := smallInt(1, n-2)
			ans := n - dimW
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("In \\(\\mathbb{R}^{%d}\\), let \\(W\\) be a subspace of dimension %d. What is the dimension of the orthogonal complement \\(W^\\perp\\)?", n, dimW)
			if ans == n || ans == dimW || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the orthogonal complement dimension theorem: \\(\\dim(W) + \\dim(W^\\perp) = n\\).\n"+
				"Subgoal 2: Substitute ambient dimension \\(n = %d\\) and subspace dimension \\(\\dim(W) = %d\\): \\((%d) - (%d)\\).\n"+
				"Subgoal 3: Subtract to evaluate \\(\\dim(W^\\perp) = %d\\).",
				n, dimW, n, dimW, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"dimW":   strconv.Itoa(dimW),
					"answer": ansStr,
				},
			}
		case 3:
			// Hyperplane defined by a single linear constraint
			n := smallInt(4, 8+tier*3)
			ans := n - 1
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the subspace in \\(\\mathbb{R}^{%d}\\) defined by the single equation \\(x_1 + x_2 + \\dots + x_{%d} = 0\\)?", n, n)
			if ans == n || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the hyperplane codimension theorem: a single non-trivial homogeneous linear equation defines a subspace of codimension 1.\n"+
				"Subgoal 2: Subtract 1 constraint from ambient dimension \\(n = %d\\): %d - 1.\n"+
				"Subgoal 3: Evaluate dimension: %d.",
				n, n, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"answer": ansStr,
				},
			}
		default:
			// k independent linear equations
			n := smallInt(5, 8+tier*2)
			k := smallInt(2, n-2)
			ans := n - k
			ansStr := strconv.Itoa(ans)
			q := fmt.Sprintf("What is the dimension of the subspace in \\(\\mathbb{R}^{%d}\\) defined by %d linearly independent homogeneous linear equations?", n, k)
			if ans == n || ans == k || containsWord(q, ansStr) {
				if retryable(attempt) {
					continue
				}
			}
			ex := fmt.Sprintf("Subgoal 1: State the constraint codimension theorem: \\(k\\) linearly independent constraints reduce dimension by \\(k\\): \\(\\dim = n - k\\).\n"+
				"Subgoal 2: Substitute ambient dimension \\(n = %d\\) and constraint count \\(k = %d\\): %d - %d.\n"+
				"Subgoal 3: Subtract to evaluate dimension: %d.",
				n, k, n, k, ans)
			return generator.Problem{
				Question:    q,
				Answer:      ansStr,
				Explanation: ex,
				Facts: map[string]string{
					"n":      strconv.Itoa(n),
					"k":      strconv.Itoa(k),
					"answer": ansStr,
				},
			}
		}
	}
}
