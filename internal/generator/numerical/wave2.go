package numerical

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func numKP(ctx generator.GeneratorContext) int {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = rand.Intn(3) + 1
	}
	return kp
}

// ---------------------------------------------------------------------------
// Aitken acceleration
// ---------------------------------------------------------------------------

type aitkenGen struct{}

func (g *aitkenGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	a := rand.Intn(6) + 1
	b := 4 * (rand.Intn(3) + 1)
	x0, x1, x2 := a+b, a+b/2, a+b/4

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = a
		question = fmt.Sprintf("Three iterates are \\(x_0 = %d\\), \\(x_1 = %d\\) and \\(x_2 = %d\\). What limit does Aitken's Δ² method estimate? (enter a number)", x0, x1, x2)
		explanation = fmt.Sprintf("The errors shrink by a factor of 1/2, so Aitken acceleration removes the geometric tail and gives \\(x_0 - (Δx_0)^2/Δ^2x_0 = %d\\).", a)
	case 2:
		answer = b
		question = fmt.Sprintf("Three iterates are \\(x_0 = %d\\), \\(x_1 = %d\\) and \\(x_2 = %d\\). Aitken's method estimates the limit %d. By how much does it correct \\(x_0\\)? (enter a number)", x0, x1, x2, a)
		explanation = fmt.Sprintf("The estimate is %d, so the correction from \\(x_0 = %d\\) is %d.", a, x0, b)
	default:
		answer = b / 4
		question = fmt.Sprintf("For the iterates \\(x_0 = %d\\), \\(x_1 = %d\\) and \\(x_2 = %d\\), what is \\(Δ^2x_0 = x_2 - 2x_1 + x_0\\)? (enter a number)", x0, x1, x2)
		explanation = fmt.Sprintf("\\(Δ^2x_0 = %d - 2(%d) + %d = %d\\).", x2, x1, x0, b/4)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"x0": strconv.Itoa(x0), "x1": strconv.Itoa(x1), "x2": strconv.Itoa(x2), "limit": strconv.Itoa(a)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Muller's method
// ---------------------------------------------------------------------------

type mullerGen struct{}

func (g *mullerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	y1 := rand.Intn(7) - 3
	u := rand.Intn(4) + 1
	w := rand.Intn(3) - 1
	v := u + 2*w
	if u == v {
		v = u + 2
	}
	y0, y2 := y1+u, y1+v
	aa := (y0 + y2 - 2*y1) / 2
	bb := (y2 - y0) / 2
	cc := y1

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = aa
		question = fmt.Sprintf("Muller's method fits a quadratic through \\((-1, %d)\\), \\((0, %d)\\) and \\((1, %d)\\). What is its leading coefficient? (enter a number)", y0, y1, y2)
		explanation = fmt.Sprintf("Write the quadratic as \\(ax^2 + bx + c\\). From the three values, \\(a = (%d + %d - 2·%d)/2 = %d\\).", y0, y2, y1, aa)
	case 2:
		answer = 4*aa + 2*bb + cc
		question = fmt.Sprintf("Muller's method fits a quadratic through \\((-1, %d)\\), \\((0, %d)\\) and \\((1, %d)\\). What is the value of that quadratic at x = 2? (enter a number)", y0, y1, y2)
		explanation = fmt.Sprintf("The quadratic is \\(%dx^2%s + %d\\). At \\(x = 2\\) it equals %d.", aa, signedXTerm(bb), cc, answer)
	default:
		answer = bb*bb - 4*aa*cc
		question = fmt.Sprintf("Muller's method fits a quadratic through \\((-1, %d)\\), \\((0, %d)\\) and \\((1, %d)\\). What is its discriminant? (enter a number)", y0, y1, y2)
		explanation = fmt.Sprintf("The quadratic is \\(%dx^2%s + %d\\), so the discriminant is \\(%d^2 - 4(%d)(%d) = %d\\).", aa, signedXTerm(bb), cc, bb, aa, cc, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"a": strconv.Itoa(aa), "b": strconv.Itoa(bb), "c": strconv.Itoa(cc)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Barycentric interpolation
// ---------------------------------------------------------------------------

type barycentricGen struct{}

func (g *barycentricGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	n := rand.Intn(7) + 2 // 2..8
	k := rand.Intn(n + 1)
	ys := make([]int, n+1)
	for i := range ys {
		ys[i] = rand.Intn(9) + 1
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = factorial(k) * factorial(n-k)
		question = fmt.Sprintf("Interpolation uses the equally spaced nodes \\(x_i = i\\) for i = 0, 1, ..., %d. What is the absolute value of the denominator in the barycentric weight for node %d? (enter a number)", n, k)
		explanation = fmt.Sprintf("The weight is \\(1/∏_{j≠%d}(x_%d - x_j)\\). For equally spaced nodes its absolute denominator is \\(%d!·%d! = %d\\).", k, k, k, n-k, answer)
	case 2:
		answer = ys[k]
		question = fmt.Sprintf("Barycentric interpolation uses nodes 0, 1, ..., %d with values %s. What value does the interpolant take at node %d? (enter a number)", n, intList(ys), k)
		explanation = fmt.Sprintf("An interpolant passes through every data point, so at \\(x_%d = %d\\) its value is exactly %d.", k, k, ys[k])
	default:
		answer = n + 1
		question = fmt.Sprintf("How many interpolation nodes are needed to determine a polynomial of degree %d uniquely? (enter a number)", n)
		explanation = fmt.Sprintf("A polynomial of degree %d has %d coefficients, so it needs %d nodes.", n, n+1, n+1)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k), "value": strconv.Itoa(ys[k])},
		KPOrdinal:   kp,
	}
}

func factorial(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func intList(nums []int) string {
	out := ""
	for i, n := range nums {
		if i > 0 {
			out += ", "
		}
		out += strconv.Itoa(n)
	}
	return out
}

// ---------------------------------------------------------------------------
// Runge phenomenon
// ---------------------------------------------------------------------------

type rungeGen struct{}

func (g *rungeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	n := rand.Intn(9) + 4 // 4..12

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = n + 1
		question = fmt.Sprintf("Runge's example uses %d equally spaced subintervals on [-1, 1]. How many interpolation nodes are there? (enter a number)", n)
		explanation = fmt.Sprintf("Each subinterval contributes one new node after the first, so %d subintervals give %d nodes.", n, n+1)
	case 2:
		answer = n - 1
		question = fmt.Sprintf("An interpolating polynomial uses %d equally spaced nodes on [-1, 1], including both endpoints. How many nodes lie strictly inside the interval? (enter a number)", n+1)
		explanation = fmt.Sprintf("The two endpoints are excluded, leaving \\(%d - 2 = %d\\) interior nodes.", n+1, n-1)
	default:
		answer = 1
		question = fmt.Sprintf("For Runge's function on [-1, 1] with %d equally spaced nodes, the interpolation error oscillates most strongly near which location? Enter 1 for the endpoints and 0 for the center.", n+1)
		explanation = "Equally spaced nodes cluster near the center and leave the endpoints poorly resolved, so the largest oscillations appear near the endpoints."
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"nodes": strconv.Itoa(n + 1), "intervals": strconv.Itoa(n)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Piecewise linear interpolation
// ---------------------------------------------------------------------------

type piecewiseLinearGen struct{}

func (g *piecewiseLinearGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	x0 := rand.Intn(5)
	step := 2 * (rand.Intn(3) + 1)
	x1 := x0 + step
	y0 := rand.Intn(9) - 4
	j := rand.Intn(7) - 3
	y1 := y0 + step*j

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = y0 + step*j/2
		question = fmt.Sprintf("A piecewise linear interpolant passes through \\((%d, %d)\\) and \\((%d, %d)\\). What is its value at the midpoint x = %d? (enter a number)", x0, y0, x1, y1, x0+step/2)
		explanation = fmt.Sprintf("The midpoint value is the average of the endpoint values: \\((%d + %d)/2 = %d\\).", y0, y1, answer)
	case 2:
		answer = j
		question = fmt.Sprintf("A piecewise linear interpolant passes through \\((%d, %d)\\) and \\((%d, %d)\\). What is its slope on that interval? (enter a number)", x0, y0, x1, y1)
		explanation = fmt.Sprintf("The slope is \\((%d - (%d))/(%d - %d) = %d\\).", y1, y0, x1, x0, j)
	default:
		answer = x0
		question = fmt.Sprintf("A piecewise linear interpolant has nodes at x = %d and x = %d. A query point x = %d lies between them. What is the left endpoint of its interval? (enter a number)", x0, x1, x0+1)
		explanation = fmt.Sprintf("The point %d lies in \\([%d, %d]\\), so the left endpoint is %d.", x0+1, x0, x1, x0)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"x0": strconv.Itoa(x0), "x1": strconv.Itoa(x1), "y0": strconv.Itoa(y0), "y1": strconv.Itoa(y1)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Boole's rule
// ---------------------------------------------------------------------------

type booleGen struct{}

func (g *booleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	f := make([]int, 5)
	for i := range f {
		f[i] = rand.Intn(10)
	}
	weighted := 7*f[0] + 32*f[1] + 12*f[2] + 32*f[3] + 7*f[4]

	var question, explanation string
	var answer int
	panels := 0
	switch kp {
	case 1:
		answer = weighted
		question = fmt.Sprintf("Boole's rule uses the weights 7, 32, 12, 32, 7. For the values \\(f_0 = %d\\), \\(f_1 = %d\\), \\(f_2 = %d\\), \\(f_3 = %d\\) and \\(f_4 = %d\\), what is the weighted sum? (enter a number)", f[0], f[1], f[2], f[3], f[4])
		explanation = fmt.Sprintf("The weighted sum is 7(%d) + 32(%d) + 12(%d) + 32(%d) + 7(%d) = %d.", f[0], f[1], f[2], f[3], f[4], weighted)
	case 2:
		answer = 2 * weighted
		question = fmt.Sprintf("Apply Boole's rule with step h = 45 to the values \\(f_0 = %d\\), \\(f_1 = %d\\), \\(f_2 = %d\\), \\(f_3 = %d\\) and \\(f_4 = %d\\). What is the integral estimate? (enter a number)", f[0], f[1], f[2], f[3], f[4])
		explanation = fmt.Sprintf("The rule is \\((2h/45)(7f_0 + 32f_1 + 12f_2 + 32f_3 + 7f_4)\\). With h = 45 the factor is 2, and the weighted sum is %d, so the estimate is %d.", weighted, 2*weighted)
	default:
		m := rand.Intn(5) + 1
		panels = m
		answer = 4*m + 1
		question = fmt.Sprintf("Boole's rule is applied to %d panels, each made of four subintervals. How many function values are needed? (enter a number)", m)
		explanation = fmt.Sprintf("Neighbouring panels share endpoints, so %d panels need \\(4(%d) + 1 = %d\\) values.", m, m, 4*m+1)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"weighted": strconv.Itoa(weighted), "panels": strconv.Itoa(panels)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Gauss-Kronrod quadrature
// ---------------------------------------------------------------------------

type gaussKronrodGen struct{}

func (g *gaussKronrodGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	left := rand.Intn(5) - 2
	right := left + rand.Intn(4) + 2

	var question, explanation string
	var answer int
	panels := 0
	switch kp {
	case 1:
		answer = 15
		question = fmt.Sprintf("A 7-15 Gauss-Kronrod rule is applied on \\([%d, %d]\\). How many nodes does the Kronrod rule use? (enter a number)", left, right)
		explanation = "The Kronrod extension of the 7-point Gauss rule has 15 nodes, arranged symmetrically about the midpoint."
	case 2:
		answer = 7
		question = fmt.Sprintf("A 7-15 Gauss-Kronrod pair is applied on \\([%d, %d]\\). How many of the 15 nodes belong to the embedded Gauss rule? (enter a number)", left, right)
		explanation = "The embedded rule is the original 7-point Gauss rule. Its nodes are a subset of the 15 Kronrod nodes, so 8 nodes are new."
	default:
		n := rand.Intn(6) + 1
		panels = n
		answer = 15 * n
		question = fmt.Sprintf("A composite 7-15 Gauss-Kronrod rule uses %d panels on \\([%d, %d]\\). How many function evaluations are needed in total? (enter a number)", n, left, right)
		explanation = fmt.Sprintf("Each panel evaluates the function at its 15 nodes, so %d panels need \\(15 × %d = %d\\) evaluations.", n, n, 15*n)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"nodes": "15", "gauss_nodes": "7", "panels": strconv.Itoa(panels)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Gauss-Jordan elimination
// ---------------------------------------------------------------------------

type gaussJordanGen struct{}

func (g *gaussJordanGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	x := rand.Intn(9) - 4
	y := rand.Intn(9) - 4
	var a, b, c, d int
	for {
		a = rand.Intn(7) - 3
		b = rand.Intn(7) - 3
		c = rand.Intn(7) - 3
		d = rand.Intn(7) - 3
		if a*d-b*c != 0 && a != 0 && c != 0 {
			break
		}
	}
	e := a*x + b*y
	f := c*x + d*y

	var answer int
	var asked, question, explanation string
	switch kp {
	case 1:
		answer = x
		asked = "x"
	case 2:
		answer = y
		asked = "y"
	default:
		answer = a*d - b*c
		asked = "the determinant of the coefficient matrix"
	}
	question = fmt.Sprintf("Use Gauss-Jordan elimination on \\(%s\\) and \\(%s\\). What is %s? (enter a number)",
		formatSystemRow(a, b, e), formatSystemRow(c, d, f), asked)
	explanation = fmt.Sprintf("Elimination gives the reduced system \\(x = %d\\), \\(y = %d\\). The coefficient determinant is \\((%d)(%d) - (%d)(%d) = %d\\).",
		x, y, a, d, b, c, a*d-b*c)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d),
			"e": strconv.Itoa(e), "f": strconv.Itoa(f), "x": strconv.Itoa(x), "y": strconv.Itoa(y),
			"det": strconv.Itoa(a*d - b*c),
		},
		KPOrdinal: kp,
	}
}

func signedXTerm(coef int) string {
	switch {
	case coef == 0:
		return ""
	case coef == 1:
		return " + x"
	case coef == -1:
		return " - x"
	case coef > 0:
		return fmt.Sprintf(" + %dx", coef)
	default:
		return fmt.Sprintf(" - %dx", -coef)
	}
}

func formatSystemRow(a, b, c int) string {
	out := fmt.Sprintf("%dx", a)
	switch {
	case a == 1:
		out = "x"
	case a == -1:
		out = "-x"
	}
	switch {
	case b == 1:
		out += " + y"
	case b == -1:
		out += " - y"
	case b > 0:
		out += fmt.Sprintf(" + %dy", b)
	case b < 0:
		out += fmt.Sprintf(" - %dy", -b)
	}
	return fmt.Sprintf("%s = %d", out, c)
}

// ---------------------------------------------------------------------------
// Thomas algorithm
// ---------------------------------------------------------------------------

type thomasGen struct{}

func (g *thomasGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	off := rand.Intn(2) + 2
	diag := rand.Intn(4) + 4
	x1 := rand.Intn(7) - 3
	x2 := rand.Intn(7) - 3
	x3 := rand.Intn(7) - 3
	d1 := diag*x1 + off*x2
	d2 := off*x1 + diag*x2 + off*x3
	d3 := off*x2 + diag*x3

	var answer int
	var asked string
	switch kp {
	case 1:
		answer = x1
		asked = "x1"
	case 2:
		answer = x2
		asked = "x2"
	default:
		answer = x3
		asked = "x3"
	}

	question := fmt.Sprintf("The tridiagonal system has equations \\(%dx_1 + %dx_2 = %d\\), \\(%dx_1 + %dx_2 + %dx_3 = %d\\) and \\(%dx_2 + %dx_3 = %d\\). Use the Thomas algorithm. What is %s? (enter a number)",
		diag, off, d1, off, diag, off, d2, off, diag, d3, asked)
	explanation := fmt.Sprintf("Forward elimination followed by back substitution gives \\(x_1 = %d\\), \\(x_2 = %d\\) and \\(x_3 = %d\\).", x1, x2, x3)

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts: map[string]string{
			"diag": strconv.Itoa(diag), "off": strconv.Itoa(off),
			"x1": strconv.Itoa(x1), "x2": strconv.Itoa(x2), "x3": strconv.Itoa(x3),
			"d1": strconv.Itoa(d1), "d2": strconv.Itoa(d2), "d3": strconv.Itoa(d3),
		},
		KPOrdinal: kp,
	}
}

// ---------------------------------------------------------------------------
// Matrix condition number
// ---------------------------------------------------------------------------

type matrixConditionGen struct{}

func (g *matrixConditionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	var vals []int
	var cond int
	for {
		s := rand.Intn(3) + 1 // smallest entry 1..3
		q := rand.Intn(3) + 2 // condition number 2..4
		largest := s * q
		if largest-s < 2 {
			continue // no room for a distinct middle entry
		}
		middle := s + rand.Intn(largest-s-1) + 1
		vals = []int{s, middle, largest}
		cond = q
		break
	}
	cond2 := cond
	cond3 := cond

	var question, explanation string
	var answer int
	condExp := 0
	switch kp {
	case 1:
		answer = cond2
		question = fmt.Sprintf("For the diagonal matrix \\(A = \\operatorname{diag}(%d, %d)\\), what is the 2-norm condition number? (enter a number)", vals[0], vals[2])
		explanation = fmt.Sprintf("For a diagonal matrix the condition number is the largest diagonal entry divided by the smallest: \\(%d/%d = %d\\).", vals[2], vals[0], cond2)
	case 2:
		answer = cond3
		question = fmt.Sprintf("For the diagonal matrix \\(A = \\operatorname{diag}(%d, %d, %d)\\), what is the 2-norm condition number? (enter a number)", vals[0], vals[1], vals[2])
		explanation = fmt.Sprintf("The largest and smallest diagonal entries are %d and %d, so the condition number is \\(%d/%d = %d\\).", vals[2], vals[0], vals[2], vals[0], cond3)
	default:
		k := rand.Intn(7) + 1
		condExp = k
		answer = k
		question = fmt.Sprintf("A matrix has 2-norm condition number \\(10^{%d}\\). About how many decimal digits of accuracy can be lost when solving a system with it? (enter a number)", k)
		explanation = fmt.Sprintf("A condition number of \\(10^{%d}\\) can amplify relative error by that factor, which corresponds to about %d lost decimal digits.", k, k)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"condition": strconv.Itoa(cond3), "exponent": strconv.Itoa(condExp)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Inverse power iteration
// ---------------------------------------------------------------------------

type inversePowerGen struct{}

func (g *inversePowerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	vals := make([]int, 0, 3)
	for len(vals) < 3 {
		v := rand.Intn(9) + 1
		dup := false
		for _, old := range vals {
			if old == v {
				dup = true
			}
		}
		if !dup {
			vals = append(vals, v)
		}
	}
	sort.Ints(vals)
	shift := vals[0]
	closest := vals[0]
	for {
		shift = vals[0] + rand.Intn(vals[2]-vals[0]+1)
		closest = vals[0]
		tie := false
		for _, v := range vals[1:] {
			d := absInt(v - shift)
			best := absInt(closest - shift)
			switch {
			case d < best:
				closest = v
			case d == best:
				tie = true
			}
		}
		if !tie {
			break
		}
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = vals[0]
		question = fmt.Sprintf("Inverse power iteration with shift 0 is applied to \\(A = \\operatorname{diag}(%d, %d, %d)\\). Which eigenvalue does it converge to? (enter a number)", vals[0], vals[1], vals[2])
		explanation = fmt.Sprintf("Inverse iteration with shift 0 converges to the eigenvalue of smallest magnitude. The entries are %d, %d and %d, so the answer is %d.", vals[0], vals[1], vals[2], vals[0])
	case 2:
		answer = closest
		question = fmt.Sprintf("Inverse power iteration with shift \\(σ = %d\\) is applied to \\(A = \\operatorname{diag}(%d, %d, %d)\\). Which eigenvalue does it converge to? (enter a number)", shift, vals[0], vals[1], vals[2])
		explanation = fmt.Sprintf("The method converges to the eigenvalue closest to %d. The distances are %d, %d and %d, so the closest eigenvalue is %d.", shift, absInt(vals[0]-shift), absInt(vals[1]-shift), absInt(vals[2]-shift), closest)
	default:
		answer = vals[2]
		question = fmt.Sprintf("The power method is applied to \\(A = \\operatorname{diag}(%d, %d, %d)\\). Which eigenvalue does it converge to? (enter a number)", vals[0], vals[1], vals[2])
		explanation = fmt.Sprintf("The power method converges to the eigenvalue of largest magnitude, which is %d.", vals[2])
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"smallest": strconv.Itoa(vals[0]), "middle": strconv.Itoa(vals[1]), "largest": strconv.Itoa(vals[2]), "shift": strconv.Itoa(shift), "closest": strconv.Itoa(closest)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Discrete Fourier transform
// ---------------------------------------------------------------------------

type dftGen struct{}

func (g *dftGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	n := []int{4, 6, 8, 10, 12}[rand.Intn(5)]
	x := make([]int, n)
	for i := range x {
		x[i] = rand.Intn(6)
	}
	sum := 0
	alt := 0
	for i, v := range x {
		sum += v
		if i%2 == 0 {
			alt += v
		} else {
			alt -= v
		}
	}

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = sum
		question = fmt.Sprintf("A length-%d signal is \\(x = (%s)\\). What is the DFT coefficient \\(X_0\\)? (enter a number)", n, intList(x))
		explanation = fmt.Sprintf("\\(X_0\\) is the sum of the samples: %d.", sum)
	case 2:
		answer = alt
		question = fmt.Sprintf("A length-%d signal is \\(x = (%s)\\). What is the DFT coefficient \\(X_{%d}\\)? (enter a number)", n, intList(x), n/2)
		explanation = fmt.Sprintf("At \\(k = %d\\) the exponential alternates sign, so \\(X_{%d}\\) is the alternating sum %d.", n/2, n/2, alt)
	default:
		answer = n * n
		question = fmt.Sprintf("A direct DFT of length %d is computed without the FFT. How many complex multiplications does the direct sum require? (enter a number)", n)
		explanation = fmt.Sprintf("Each of the %d outputs sums %d terms, so the direct computation uses \\(%d × %d = %d\\) complex multiplications.", n, n, n, n, n*n)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"n": strconv.Itoa(n), "sum": strconv.Itoa(sum), "alternating": strconv.Itoa(alt)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Aliasing
// ---------------------------------------------------------------------------

type aliasingGen struct{}

func (g *aliasingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	fs := (rand.Intn(7) + 4) * 2 // 8..20, even
	alias := rand.Intn(fs/2-1) + 1
	f := fs + alias

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = fs / 2
		question = fmt.Sprintf("A signal is sampled at %d Hz. What is the Nyquist frequency? (enter a number)", fs)
		explanation = fmt.Sprintf("The Nyquist frequency is half the sampling rate: \\(%d/2 = %d\\) Hz.", fs, fs/2)
	case 2:
		answer = alias
		question = fmt.Sprintf("A %d Hz sinusoid is sampled at %d Hz. What alias frequency is observed? (enter a number)", f, fs)
		explanation = fmt.Sprintf("The sampled frequency is \\(|%d - %d| = %d\\) Hz, which lies below Nyquist.", f, fs, alias)
	default:
		answer = 2 * f
		question = fmt.Sprintf("A signal contains a %d Hz component. What is the minimum sampling rate that avoids aliasing it? (enter a number)", f)
		explanation = fmt.Sprintf("The sampling rate must exceed twice the highest frequency: \\(2 × %d = %d\\) Hz.", f, 2*f)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"fs": strconv.Itoa(fs), "f": strconv.Itoa(f), "alias": strconv.Itoa(alias)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Pade approximants
// ---------------------------------------------------------------------------

type padeGen struct{}

func (g *padeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	m := rand.Intn(3) + 1
	n := rand.Intn(3) + 1

	var question, explanation string
	var answer int
	padeX := 0
	switch kp {
	case 1:
		answer = m + n + 1
		question = fmt.Sprintf("A Padé approximant of type \\([%d/%d]\\) is built. How many free coefficients does it have in total? (enter a number)", m, n)
		explanation = fmt.Sprintf("The numerator has %d coefficients and the denominator has %d, with one normalization, so \\(%d + %d + 1 = %d\\) are free.", m+1, n+1, m, n, m+n+1)
	case 2:
		answer = m + n
		question = fmt.Sprintf("A Padé approximant of type \\([%d/%d]\\) is matched to a Taylor series at 0. Through what degree does it match? (enter a number)", m, n)
		explanation = fmt.Sprintf("The \\([%d/%d]\\) Padé approximant matches the Taylor series through degree \\(%d + %d = %d\\).", m, n, m, n, m+n)
	default:
		x := 2 * (rand.Intn(9) + 1)
		padeX = x
		answer = 1 + x/2
		question = fmt.Sprintf("The \\([1/1]\\) Padé approximant of \\(e^x\\) is \\((1 + x/2)/(1 - x/2)\\). What is its numerator at x = %d? (enter a number)", x)
		explanation = fmt.Sprintf("The numerator is \\(1 + %d/2 = %d\\).", x, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"m": strconv.Itoa(m), "n": strconv.Itoa(n), "x": strconv.Itoa(padeX)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Finite differences
// ---------------------------------------------------------------------------

type numDiffGen struct{}

func (g *numDiffGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	x := rand.Intn(6) + 1
	h := rand.Intn(3) + 1

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = 2*x + h
		question = fmt.Sprintf("For \\(f(x) = x^2\\), what is the forward difference \\((f(%d + %d) - f(%d))/%d\\)? (enter a number)", x, h, x, h)
		explanation = fmt.Sprintf("\\(f(%d) = %d\\) and \\(f(%d) = %d\\), so the quotient is \\((%d - %d)/%d = %d\\).", x, x*x, x+h, (x+h)*(x+h), (x+h)*(x+h), x*x, h, answer)
	case 2:
		answer = 2 * x
		question = fmt.Sprintf("For \\(f(x) = x^2\\), what is the central difference \\((f(%d + %d) - f(%d - %d))/(2·%d)\\)? (enter a number)", x, h, x, h, h)
		explanation = fmt.Sprintf("The central difference of \\(x^2\\) is exact: \\(((%d)^2 - (%d)^2)/(2·%d) = %d\\).", x+h, x-h, h, answer)
	default:
		answer = 2 * h * h
		question = fmt.Sprintf("For \\(f(x) = x^2\\), what is the second difference \\(f(%d + %d) - 2f(%d) + f(%d - %d)\\)? (enter a number)", x, h, x, x, h)
		explanation = fmt.Sprintf("The second difference of \\(x^2\\) with step %d is \\(2·%d^2 = %d\\).", h, h, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"x": strconv.Itoa(x), "h": strconv.Itoa(h)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Richardson extrapolation
// ---------------------------------------------------------------------------

type richardsonGen struct{}

func (g *richardsonGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	t := rand.Intn(9) + 1
	k := rand.Intn(3) + 1
	d1 := t + 4*k
	d2 := t + k

	var question, explanation string
	var answer int
	p := 0
	switch kp {
	case 1:
		answer = t
		question = fmt.Sprintf("Two approximations with error \\(O(h^2)\\) are \\(D(h) = %d\\) and \\(D(h/2) = %d\\). What is the Richardson extrapolated value? (enter a number)", d1, d2)
		explanation = fmt.Sprintf("For second-order error, \\(D_{ext} = (4D(h/2) - D(h))/3 = (4·%d - %d)/3 = %d\\).", d2, d1, t)
	case 2:
		answer = d2 - d1
		question = fmt.Sprintf("Two approximations with error \\(O(h^2)\\) are \\(D(h) = %d\\) and \\(D(h/2) = %d\\). What is \\(D(h/2) - D(h)\\)? (enter a number)", d1, d2)
		explanation = fmt.Sprintf("Subtract the approximations: \\(%d - %d = %d\\). The error term has shrunk by a factor of four.", d2, d1, answer)
	default:
		p = rand.Intn(7) + 2
		answer = 1 << p
		question = fmt.Sprintf("Richardson extrapolation combines two approximations whose error is \\(O(h^{%d})\\). By what factor does the error shrink when h is halved? (enter a number)", p)
		explanation = fmt.Sprintf("The error scales as \\((h/2)^{%d}/h^{%d} = 1/2^{%d}\\), so the correction uses the factor \\(2^{%d} = %d\\).", p, p, p, p, answer)
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"extrapolated": strconv.Itoa(t), "d1": strconv.Itoa(d1), "d2": strconv.Itoa(d2), "p": strconv.Itoa(p)},
		KPOrdinal:   kp,
	}
}

// ---------------------------------------------------------------------------
// Line search
// ---------------------------------------------------------------------------

type lineSearchGen struct{}

func (g *lineSearchGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := numKP(ctx)
	center := rand.Intn(8) + 1
	x0 := rand.Intn(7) - 3
	dir := 1
	if rand.Intn(2) == 0 {
		dir = -1
	}
	alpha := (center - x0) * dir

	var question, explanation string
	var answer int
	switch kp {
	case 1:
		answer = alpha
		question = fmt.Sprintf("An exact line search minimizes \\(f(x) = (x - %d)^2\\) from \\(x_0 = %d\\) along the direction \\(d = %d\\). What step length α is optimal? (enter a number)", center, x0, dir)
		explanation = fmt.Sprintf("The derivative at \\(x_0\\) is \\(2(%d - %d)\\). Setting \\(x_0 + αd = %d\\) gives \\(α = %d\\).", x0, center, center, alpha)
	case 2:
		answer = center
		question = fmt.Sprintf("An exact line search minimizes \\(f(x) = (x - %d)^2\\) from \\(x_0 = %d\\) along \\(d = %d\\). At what x does it stop? (enter a number)", center, x0, dir)
		explanation = fmt.Sprintf("The minimum of \\(f\\) is at \\(x = %d\\), and the step reaches it exactly.", center)
	default:
		answer = 2
		question = fmt.Sprintf("A line search for \\(f(x) = (x - %d)^2\\) starts at \\(x_0 = %d\\). The strong Wolfe conditions combine an Armijo sufficient-decrease condition and a curvature condition. How many inequalities is that? (enter a number)", center, x0)
		explanation = "The strong Wolfe conditions are two inequalities: sufficient decrease, and a bound on the magnitude of the directional derivative."
	}

	return generator.Problem{
		Question:    question,
		Answer:      strconv.Itoa(answer),
		Explanation: explanation,
		Facts:       map[string]string{"center": strconv.Itoa(center), "x0": strconv.Itoa(x0), "alpha": strconv.Itoa(alpha)},
		KPOrdinal:   kp,
	}
}
