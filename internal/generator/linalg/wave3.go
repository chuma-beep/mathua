package linalg

// Wave 3 linear algebra: twenty-one concepts. The generators ask for scalar
// quantities (an entry, a trace, a determinant, a rank, a norm, a sign) so the
// numeric grader is sufficient and no matrix-answer parser is needed. Facts
// publish every matrix entry used, and wave3_test.go recomputes each answer.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave3(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"linalg.matrix.powers":             "powers",
		"linalg.matrix.nilpotent":          "nilpotent",
		"linalg.matrix.idempotent":         "idempotent",
		"linalg.matrix.elementary":         "elementary",
		"linalg.matrix.block":              "block",
		"linalg.matrix.permutation":        "permutation",
		"linalg.matrix.norm":               "norm",
		"linalg.decomp.pseudoinverse":      "pseudoinverse",
		"linalg.decomp.cholesky":           "cholesky",
		"linalg.matrix.schur_complement":   "schur_complement",
		"linalg.matrix.trace_identities":   "trace_identities",
		"linalg.matrix.rank_factorization": "rank_factorization",
		"linalg.matrix.outer_product":      "outer_product",
		"linalg.matrix.hadamard":           "hadamard",
		"linalg.matrix.kronecker":          "kronecker",
		"linalg.matrix.vandermonde":        "vandermonde",
		"linalg.matrix.companion":          "companion",
		"linalg.matrix.circulant":          "circulant",
		"linalg.transform.householder":     "householder",
		"linalg.transform.givens":          "givens",
		"linalg.matrix.determinant_lemma":  "determinant_lemma",
	} {
		reg.Register(id, &w3LinalgGen{kind: kind})
	}
}

type w3LinalgGen struct{ kind string }

func (g *w3LinalgGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "powers":
		return w3MatrixPowers(kp)
	case "nilpotent":
		return w3Nilpotent(kp)
	case "idempotent":
		return w3Idempotent(kp)
	case "elementary":
		return w3Elementary(kp)
	case "block":
		return w3Block(kp)
	case "permutation":
		return w3Permutation(kp)
	case "norm":
		return w3MatrixNorm(kp)
	case "pseudoinverse":
		return w3Pseudoinverse(kp)
	case "cholesky":
		return w3Cholesky(kp)
	case "schur_complement":
		return w3SchurComplement(kp)
	case "trace_identities":
		return w3TraceIdentities(kp)
	case "rank_factorization":
		return w3RankFactorization(kp)
	case "outer_product":
		return w3OuterProduct(kp)
	case "hadamard":
		return w3Hadamard(kp)
	case "kronecker":
		return w3Kronecker(kp)
	case "vandermonde":
		return w3Vandermonde(kp)
	case "companion":
		return w3Companion(kp)
	case "circulant":
		return w3Circulant(kp)
	case "householder":
		return w3Householder(kp)
	case "givens":
		return w3Givens(kp)
	case "determinant_lemma":
		return w3DeterminantLemma(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w3L(o int, q, a, ex string, facts map[string]string) generator.Problem {
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: facts, KPOrdinal: o}
}

func w3LI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w3LR(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w3LNZ(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	for {
		v := w3LR(lo, hi)
		if v != 0 {
			return v
		}
	}
}

func w3LMat(a, b, c, d int) string { return fmt.Sprintf("[[%d,%d],[%d,%d]]", a, b, c, d) }

func w3LDet(a, b, c, d int) int { return a*d - b*c }

func w3LTrace(a, b, c, d int) int { return a + d }

func w3LMatMul2(a, b, c, d, e, f, g, h int) (int, int, int, int) {
	return a*e + b*g, a*f + b*h, c*e + d*g, c*f + d*h
}

func w3LMatPow2(a, b, c, d, k int) (int, int, int, int) {
	ra, rb, rc, rd := 1, 0, 0, 1
	for i := 0; i < k; i++ {
		ra, rb, rc, rd = ra*a+rb*c, ra*b+rb*d, rc*a+rd*c, rc*b+rd*d
	}
	return ra, rb, rc, rd
}

func w3LParse(facts map[string]string, prefix string) (a, b, c, d int) {
	a, _ = strconv.Atoi(facts[prefix+"a"])
	b, _ = strconv.Atoi(facts[prefix+"b"])
	c, _ = strconv.Atoi(facts[prefix+"c"])
	d, _ = strconv.Atoi(facts[prefix+"d"])
	return
}

func w3LSetMatrix(f map[string]string, prefix string, a, b, c, d int) {
	w3LI(f, prefix+"a", a)
	w3LI(f, prefix+"b", b)
	w3LI(f, prefix+"c", c)
	w3LI(f, prefix+"d", d)
}

func w3MatrixPowers(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4)
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	switch kp {
	case 1:
		ra, _, _, _ := w3LMatPow2(a, b, c, d, 2) // aislop-ignore-line ai-slop/swallowed-exception -- four matrix entries, not an error
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the (1,1) entry of \\(A^2\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(ra),
			fmt.Sprintf("Multiply row 1 of A by column 1 of A: \\((%d)(%d)+(%d)(%d)=%d\\).", a, a, b, c, ra), f)
	case 2:
		ra, _, _, rd := w3LMatPow2(a, b, c, d, 3)
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is \\(\\operatorname{tr}(A^3)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(ra+rd),
			fmt.Sprintf("Compute \\(A^3\\) and add its diagonal entries. The trace is %d.", ra+rd), f)
	default:
		k := w3LR(4, 6)
		_, rb, _, _ := w3LMatPow2(a, b, c, d, k) // aislop-ignore-line ai-slop/swallowed-exception -- four matrix entries, not an error
		w3LI(f, "power", k)
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the (1,2) entry of \\(A^{%d}\\)? (enter a number)", w3LMat(a, b, c, d), k), strconv.Itoa(rb),
			fmt.Sprintf("Multiply the matrix by itself %d times, or use the recurrence for the entries. The (1,2) entry is %d.", k, rb), f)
	}
}

func w3Nilpotent(kp int) generator.Problem {
	// A = [[x, y],[-x^2/y, -x]] has trace 0 and determinant 0; with y | x^2 it is nilpotent.
	x := w3LNZ(-3, 3)
	y := w3LR(1, 3)
	for x*x%y != 0 {
		y = w3LR(1, 3)
	}
	c := -(x * x) / y
	a, b, d := x, y, -x
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	ra, rb, rc, rd := w3LMatPow2(a, b, c, d, 2)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the (1,1) entry of \\(A^2\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(ra),
			fmt.Sprintf("Multiply the matrix by itself. The (1,1) entry is %d.", ra), f)
	case 2:
		idx := 1
		if !(a == 0 && b == 0 && c == 0 && d == 0) {
			idx = 2
		}
		w3LI(f, "nilpotency_index", idx)
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the nilpotency index (the least positive \\(k\\) with \\(A^k=0\\))? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(idx),
			fmt.Sprintf("The matrix has trace 0 and determinant 0, so its square is the zero matrix. The least such power is %d.", idx), f)
	default:
		is := 1
		if !(ra == 0 && rb == 0 && rc == 0 && rd == 0) {
			is = 0
		}
		w3LI(f, "nilpotent", is)
		return w3L(kp, fmt.Sprintf("Is \\(A=%s\\) nilpotent? Enter 1 for yes and 0 for no.", w3LMat(a, b, c, d)), strconv.Itoa(is),
			fmt.Sprintf("A 2x2 matrix with trace 0 and determinant 0 is nilpotent: \\(A^2=0\\). The answer is %d.", is), f)
	}
}

func w3Idempotent(kp int) generator.Problem {
	// A = [[p, p(1-p)/q],[q, 1-p]] is a projection when q divides p(1-p).
	p := w3LR(0, 2)
	q := w3LR(1, 3)
	for (p*(1-p))%q != 0 {
		q = w3LR(1, 3)
	}
	a, b, c, d := p, p*(1-p)/q, q, 1-p
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	ra, rb, rc, rd := w3LMatPow2(a, b, c, d, 2)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For the idempotent matrix \\(A=%s\\), what is the (2,2) entry of \\(A^2\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(rd),
			fmt.Sprintf("Because \\(A^2=A\\), the (2,2) entry of \\(A^2\\) is the (2,2) entry of A, which is %d.", rd), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the idempotent matrix \\(A=%s\\), what is \\(\\operatorname{tr}(A)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(a+d),
			fmt.Sprintf("Add the diagonal entries: \\(%d+%d=%d\\).", a, d, a+d), f)
	default:
		is := 1
		if !(ra == a && rb == b && rc == c && rd == d) {
			is = 0
		}
		w3LI(f, "idempotent", is)
		return w3L(kp, fmt.Sprintf("Is \\(A=%s\\) idempotent, meaning \\(A^2=A\\)? Enter 1 for yes and 0 for no.", w3LMat(a, b, c, d)), strconv.Itoa(is),
			fmt.Sprintf("Square the matrix and compare every entry with A. The answer is %d.", is), f)
	}
}

func w3Elementary(kp int) generator.Problem {
	// E is a 2x2 elementary matrix: add k times row 1 to row 2.
	k := w3LNZ(-3, 3)
	a, b, c, d := 1, 0, k, 1
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	w3LI(f, "k", k)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For \\(E=%s\\), what is \\(\\det(E)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(w3LDet(a, b, c, d)),
			fmt.Sprintf("Adding a multiple of one row to another does not change the determinant, so \\(\\det(E)=1\\)."), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(E=%s\\), what is the (2,1) entry? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(c),
			fmt.Sprintf("The elementary matrix records the row operation: its (2,1) entry is %d.", c), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(E=%s\\), what is \\(\\operatorname{tr}(E)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(a+d),
			fmt.Sprintf("Add the diagonal entries: \\(%d+%d=%d\\).", a, d, a+d), f)
	}
}

func w3Block(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-3, 3), w3LNZ(-3, 3), w3LNZ(-3, 3), w3LNZ(-3, 3)
	e, f2, g, h := w3LNZ(-3, 3), w3LNZ(-3, 3), w3LNZ(-3, 3), w3LNZ(-3, 3)
	f := map[string]string{}
	w3LSetMatrix(f, "A", a, b, c, d)
	w3LSetMatrix(f, "B", e, f2, g, h)
	detA := w3LDet(a, b, c, d)
	detB := w3LDet(e, f2, g, h)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("The block-diagonal matrix \\(M=\\begin{pmatrix}A&0\\\\0&B\\end{pmatrix}\\) has \\(A=%s\\) and \\(B=%s\\). What is \\(\\det(M)\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(detA*detB),
			fmt.Sprintf("The determinant of a block-diagonal matrix is the product of the block determinants: \\((%d)(%d)=%d\\).", detA, detB, detA*detB), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the block-diagonal matrix with \\(A=%s\\) and \\(B=%s\\), what is the trace of the whole matrix? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(a+d+e+h),
			fmt.Sprintf("The trace is the sum of the diagonal entries of both blocks: \\((%d)+(%d)+(%d)+(%d)=%d\\).", a, d, e, h, a+d+e+h), f)
	default:
		ra, _, _, _ := w3LMatMul2(a, b, c, d, a, b, c, d) // aislop-ignore-line ai-slop/swallowed-exception -- four matrix entries, not an error
		return w3L(kp, fmt.Sprintf("For the block matrix \\(M=\\begin{pmatrix}A&0\\\\0&B\\end{pmatrix}\\) with \\(A=%s\\) and \\(B=%s\\), what is the (1,1) entry of \\(M^2\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(ra),
			fmt.Sprintf("The (1,1) block of \\(M^2\\) is \\(A^2\\), and its (1,1) entry is %d.", ra), f)
	}
}

func w3Permutation(kp int) generator.Problem {
	perm := rand.Perm(3)
	// perm[i] is the column of the 1 in row i.
	inv := 0
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if perm[i] > perm[j] {
				inv++
			}
		}
	}
	sign := 1
	if inv%2 == 1 {
		sign = -1
	}
	f := map[string]string{"permutation": strings.Trim(strings.ReplaceAll(fmt.Sprint(perm), " ", ","), "[]")}
	w3LI(f, "inversions", inv)
	w3LI(f, "sign", sign)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("The permutation matrix for \\(\\sigma=%v\\) sends row \\(i\\) to column \\(\\sigma(i)\\). What is its determinant? (enter a number)", perm), strconv.Itoa(sign),
			fmt.Sprintf("The determinant is the sign of the permutation. With %d inversions, the sign is %d.", inv, sign), f)
	case 2:
		return w3L(kp, fmt.Sprintf("How many inversions does the permutation \\(\\sigma=%v\\) have? (enter a number)", perm), strconv.Itoa(inv),
			func() string {
				if inv == 1 {
					return "There is 1 inversion."
				}
				return fmt.Sprintf("There are %d inversions.", inv)
			}(), f)
	default:
		i := w3LR(0, 2)
		x := []int{w3LR(1, 5), w3LR(1, 5), w3LR(1, 5)}
		w3LI(f, "row", i)
		f["vector"] = fmt.Sprintf("%d,%d,%d", x[0], x[1], x[2])
		// (Px)_i = x_{perm[i]}.
		ans := x[perm[i]]
		return w3L(kp, fmt.Sprintf("For \\(\\sigma=%v\\) and \\(x=(%d,%d,%d)\\), what is component %d of \\(P_\\sigma x\\)? (enter a number)", perm, x[0], x[1], x[2], i+1), strconv.Itoa(ans),
			fmt.Sprintf("A permutation matrix permutes coordinates: component %d of \\(P_\\sigma x\\) is \\(x_{\\sigma(%d)}=x_{%d}=%d\\).", i+1, i+1, perm[i]+1, ans), f)
	}
}

func w3MatrixNorm(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-5, 6), w3LNZ(-5, 6), w3LNZ(-5, 6), w3LNZ(-5, 6)
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	col1 := absInt(a) + absInt(c)
	col2 := absInt(b) + absInt(d)
	row1 := absInt(a) + absInt(b)
	row2 := absInt(c) + absInt(d)
	one := col1
	if col2 > one {
		one = col2
	}
	inf := row1
	if row2 > inf {
		inf = row2
	}
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the induced 1-norm \\(\\|A\\|_1\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(one),
			fmt.Sprintf("The 1-norm is the largest absolute column sum. The column sums are %d and %d, so the norm is %d.", col1, col2, one), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the induced infinity norm \\(\\|A\\|_\\infty\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(inf),
			fmt.Sprintf("The infinity norm is the largest absolute row sum. The row sums are %d and %d, so the norm is %d.", row1, row2, inf), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is the squared Frobenius norm \\(\\|A\\|_F^2\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(a*a+b*b+c*c+d*d),
			fmt.Sprintf("Square every entry and add: \\((%d)^2+(%d)^2+(%d)^2+(%d)^2=%d\\).", a, b, c, d, a*a+b*b+c*c+d*d), f)
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func w3Pseudoinverse(kp int) generator.Problem {
	n := 3
	vals := make([]int, n)
	for i := range vals {
		if rand.Intn(3) == 0 {
			vals[i] = 0
		} else {
			vals[i] = w3LNZ(-4, 4)
		}
	}
	vals[0] = w3LNZ(-4, 4)
	rank := 0
	for _, v := range vals {
		if v != 0 {
			rank++
		}
	}
	f := map[string]string{}
	for i, v := range vals {
		w3LI(f, fmt.Sprintf("d%d", i+1), v)
	}
	w3LI(f, "rank", rank)
	switch kp {
	case 1:
		i := w3LR(0, n-1)
		ans := 0
		if vals[i] != 0 {
			if vals[i] == 1 || vals[i] == -1 {
				ans = vals[i]
			} else {
				// Answers may be fractional; choose only units for this family.
				for vals[i] != 1 && vals[i] != -1 {
					vals[i] = w3LNZ(-1, 1)
				}
				w3LI(f, fmt.Sprintf("d%d", i+1), vals[i])
				ans = vals[i]
			}
		}
		return w3L(kp, fmt.Sprintf("For the diagonal matrix \\(D=\\operatorname{diag}(%d,%d,%d)\\), what is diagonal entry %d of the pseudoinverse \\(D^+\\)? (enter a number)", vals[0], vals[1], vals[2], i+1), strconv.Itoa(ans),
			fmt.Sprintf("The pseudoinverse inverts each nonzero diagonal entry and leaves zero entries zero. The requested entry is %d.", ans), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the diagonal matrix \\(D=\\operatorname{diag}(%d,%d,%d)\\), what is the rank? (enter a number)", vals[0], vals[1], vals[2]), strconv.Itoa(rank),
			fmt.Sprintf("The rank is the number of nonzero diagonal entries: %d.", rank), f)
	default:
		ans := 0
		for _, v := range vals {
			if v != 0 {
				ans++
			}
		}
		return w3L(kp, fmt.Sprintf("For the diagonal matrix \\(D=\\operatorname{diag}(%d,%d,%d)\\), what is \\(\\operatorname{tr}(DD^+)\\)? (enter a number)", vals[0], vals[1], vals[2]), strconv.Itoa(ans),
			fmt.Sprintf("Each nonzero diagonal entry contributes 1 to \\(DD^+\\), and each zero contributes 0. The trace is %d.", ans), f)
	}
}

func w3Cholesky(kp int) generator.Problem {
	p := w3LNZ(1, 4)
	q := w3LNZ(-4, 4)
	r := w3LNZ(1, 4)
	a, b, c, d := p*p, p*q, p*q, q*q+r*r
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	w3LI(f, "p", p)
	w3LI(f, "q", q)
	w3LI(f, "r", r)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("The Cholesky factor of \\(A=%s\\) is \\(L=\\begin{pmatrix}p&0\\\\q&r\\end{pmatrix}\\). What is \\(p\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(p),
			fmt.Sprintf("\\(p^2\\) is the (1,1) entry of A, so \\(p=%d\\).", p), f)
	case 2:
		return w3L(kp, fmt.Sprintf("The Cholesky factor of \\(A=%s\\) is \\(L=\\begin{pmatrix}p&0\\\\q&r\\end{pmatrix}\\). What is \\(q\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(q),
			fmt.Sprintf("The (2,1) entry of A is \\(pq\\), and \\(p=%d\\), so \\(q=%d\\).", p, q), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is \\(\\det(A)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(a*d-b*c),
			fmt.Sprintf("For this positive-definite matrix, \\(\\det(A)=(pr)^2\\), which is %d.", a*d-b*c), f)
	}
}

func w3SchurComplement(kp int) generator.Problem {
	a := w3LNZ(1, 3)
	b := w3LNZ(-3, 3)
	c := w3LNZ(-3, 3)
	d := w3LNZ(-3, 3)
	// Make c*b divisible by a so the Schur complement is an integer.
	for (c*b)%a != 0 {
		c = w3LNZ(-3, 3)
	}
	comp := d - c*b/a
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	w3LI(f, "schur", comp)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For the block matrix \\(M=\\begin{pmatrix}a&b\\\\c&d\\end{pmatrix}\\) with \\(a=%d\\), \\(b=%d\\), \\(c=%d\\), \\(d=%d\\), what is the Schur complement \\(d-ca^{-1}b\\)? (enter a number)", a, b, c, d), strconv.Itoa(comp),
			fmt.Sprintf("The 1x1 block form gives \\(d-ca^{-1}b=%d-(%d)(%d)/(%d)=%d\\).", d, c, b, a, comp), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(M=%s\\), what is \\(\\det(M)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(a*comp),
			fmt.Sprintf("The determinant is \\(a\\) times the Schur complement: \\((%d)(%d)=%d\\).", a, comp, a*comp), f)
	default:
		is := 0
		if comp != 0 {
			is = 1
		}
		w3LI(f, "invertible", is)
		return w3L(kp, fmt.Sprintf("For \\(M=%s\\), is the matrix invertible? Enter 1 for yes and 0 for no.", w3LMat(a, b, c, d)), strconv.Itoa(is),
			fmt.Sprintf("The matrix is invertible exactly when its determinant \\(a(d-ca^{-1}b)\\) is nonzero. Here the answer is %d.", is), f)
	}
}

func w3TraceIdentities(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5)
	e, f2, g, h := w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5)
	f := map[string]string{}
	w3LSetMatrix(f, "A", a, b, c, d)
	w3LSetMatrix(f, "B", e, f2, g, h)
	switch kp {
	case 1:
		ra, _, _, rd := w3LMatMul2(a, b, c, d, e, f2, g, h)
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is \\(\\operatorname{tr}(AB)\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(ra+rd),
			fmt.Sprintf("The trace of a product is \\(\\sum_{i,j}A_{ij}B_{ji}\\). Here it is %d.", ra+rd), f)
	case 2:
		ra, _, _, rd := w3LMatMul2(a, b, c, d, a, b, c, d)
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\), what is \\(\\operatorname{tr}(A^2)\\)? (enter a number)", w3LMat(a, b, c, d)), strconv.Itoa(ra+rd),
			fmt.Sprintf("Square A and add its diagonal entries. The trace is %d.", ra+rd), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is \\(\\operatorname{tr}(A+B)\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(a+d+e+h),
			fmt.Sprintf("Trace is linear: \\(\\operatorname{tr}(A+B)=%d+%d=%d\\).", a+d, e+h, a+d+e+h), f)
	}
}

func w3RankFactorization(kp int) generator.Problem {
	u1, u2 := w3LNZ(-4, 4), w3LNZ(-4, 4)
	v1, v2 := w3LNZ(-4, 4), w3LNZ(-4, 4)
	a, b := u1*v1, u1*v2
	c, d := u2*v1, u2*v2
	f := map[string]string{}
	w3LSetMatrix(f, "", a, b, c, d)
	w3LI(f, "u1", u1)
	w3LI(f, "u2", u2)
	w3LI(f, "v1", v1)
	w3LI(f, "v2", v2)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("The matrix \\(A=%s\\) has rank-one factorisation \\(A=uv^T\\) with \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\). What is \\(\\operatorname{rank}(A)\\)? (enter a number)", w3LMat(a, b, c, d), u1, u2, v1, v2), "1",
			"The outer product of two nonzero vectors has rank 1.", f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(A=uv^T\\) with \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), what is \\(\\operatorname{tr}(A)\\)? (enter a number)", u1, u2, v1, v2), strconv.Itoa(a+d),
			fmt.Sprintf("The trace of \\(uv^T\\) is the dot product \\(u\\cdot v=%d\\).", a+d), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=uv^T\\) with \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), what is the sum of all four entries of A? (enter a number)", u1, u2, v1, v2), strconv.Itoa(a+b+c+d),
			fmt.Sprintf("The entries sum to \\((u_1+u_2)(v_1+v_2)=%d\\).", a+b+c+d), f)
	}
}

func w3OuterProduct(kp int) generator.Problem {
	u := []int{w3LNZ(-4, 4), w3LNZ(-4, 4), w3LNZ(-4, 4)}
	v := []int{w3LNZ(-4, 4), w3LNZ(-4, 4), w3LNZ(-4, 4)}
	f := map[string]string{}
	for i := 0; i < 3; i++ {
		w3LI(f, fmt.Sprintf("u%d", i+1), u[i])
		w3LI(f, fmt.Sprintf("v%d", i+1), v[i])
	}
	switch kp {
	case 1:
		i, j := w3LR(0, 2), w3LR(0, 2)
		w3LI(f, "row", i+1)
		w3LI(f, "col", j+1)
		return w3L(kp, fmt.Sprintf("For \\(u=%v\\) and \\(v=%v\\), what is the (%d,%d) entry of \\(uv^T\\)? (enter a number)", u, v, i+1, j+1), strconv.Itoa(u[i]*v[j]),
			fmt.Sprintf("The (i,j) entry of an outer product is \\(u_i v_j=(%d)(%d)=%d\\).", u[i], v[j], u[i]*v[j]), f)
	case 2:
		dot := u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
		return w3L(kp, fmt.Sprintf("For \\(u=%v\\) and \\(v=%v\\), what is \\(\\operatorname{tr}(uv^T)\\)? (enter a number)", u, v), strconv.Itoa(dot),
			fmt.Sprintf("The trace of an outer product is \\(u\\cdot v=%d\\).", dot), f)
	default:
		return w3L(kp, fmt.Sprintf("For nonzero \\(u=%v\\) and \\(v=%v\\), what is \\(\\operatorname{rank}(uv^T)\\)? (enter a number)", u, v), "1",
			"Every nonzero outer product is a rank-one matrix.", f)
	}
}

func w3Hadamard(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5)
	e, f2, g, h := w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5), w3LNZ(-4, 5)
	f := map[string]string{}
	w3LSetMatrix(f, "A", a, b, c, d)
	w3LSetMatrix(f, "B", e, f2, g, h)
	switch kp {
	case 1:
		i, j := w3LR(0, 1), w3LR(0, 1)
		w3LI(f, "row", i+1)
		w3LI(f, "col", j+1)
		valsA := [][]int{{a, b}, {c, d}}
		valsB := [][]int{{e, f2}, {g, h}}
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is the (%d,%d) entry of the Hadamard product \\(A\\circ B\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h), i+1, j+1), strconv.Itoa(valsA[i][j]*valsB[i][j]),
			fmt.Sprintf("The Hadamard product multiplies corresponding entries: \\((%d)(%d)=%d\\).", valsA[i][j], valsB[i][j], valsA[i][j]*valsB[i][j]), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is \\(\\operatorname{tr}(A\\circ B)\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(a*e+d*h),
			fmt.Sprintf("Only diagonal entries contribute: \\((%d)(%d)+(%d)(%d)=%d\\).", a, e, d, h, a*e+d*h), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is the sum of the four entries of \\(A\\circ B\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(a*e+b*f2+c*g+d*h),
			fmt.Sprintf("The entries of the Hadamard product are %d, %d, %d, %d; their sum is %d.", a*e, b*f2, c*g, d*h, a*e+b*f2+c*g+d*h), f)
	}
}

func w3Kronecker(kp int) generator.Problem {
	a, b, c, d := w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4)
	e, f2, g, h := w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4)
	f := map[string]string{}
	w3LSetMatrix(f, "A", a, b, c, d)
	w3LSetMatrix(f, "B", e, f2, g, h)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), both 2x2, what is the dimension of the square matrix \\(A\\otimes B\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), "4",
			"The Kronecker product of an \\(m\\times n\\) matrix and a \\(p\\times q\\) matrix has dimensions \\((mp)\\times(nq)\\), here 4x4.", f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is the (1,1) entry of \\(A\\otimes B\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa(a*e),
			fmt.Sprintf("The (1,1) entry is \\(a_{11}b_{11}=(%d)(%d)=%d\\).", a, e, a*e), f)
	default:
		return w3L(kp, fmt.Sprintf("For \\(A=%s\\) and \\(B=%s\\), what is \\(\\operatorname{tr}(A\\otimes B)\\)? (enter a number)", w3LMat(a, b, c, d), w3LMat(e, f2, g, h)), strconv.Itoa((a+d)*(e+h)),
			fmt.Sprintf("The trace of a Kronecker product is the product of traces: \\((%d)(%d)=%d\\).", a+d, e+h, (a+d)*(e+h)), f)
	}
}

func w3Vandermonde(kp int) generator.Problem {
	x, y, z := w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4)
	for y == x || z == x || z == y {
		y, z = w3LNZ(-3, 4), w3LNZ(-3, 4)
	}
	f := map[string]string{}
	w3LI(f, "x", x)
	w3LI(f, "y", y)
	w3LI(f, "z", z)
	switch kp {
	case 1:
		det := (y - x) * (z - x) * (z - y)
		return w3L(kp, fmt.Sprintf("For the Vandermonde matrix with nodes %d, %d and %d, what is its determinant? (enter a number)", x, y, z), strconv.Itoa(det),
			fmt.Sprintf("The Vandermonde determinant is \\((y-x)(z-x)(z-y)=%d\\).", det), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the Vandermonde matrix with nodes %d, %d and %d, what is the (2,3) entry? (enter a number)", x, y, z), strconv.Itoa(z*z),
			fmt.Sprintf("Rows use powers 0,1,2. The (2,3) entry is \\(z^2=%d\\).", z*z), f)
	default:
		return w3L(kp, fmt.Sprintf("For the Vandermonde matrix with nodes %d, %d and %d (all distinct), what is its rank? (enter a number)", x, y, z), "3",
			"Distinct nodes give a nonzero Vandermonde determinant, so the 3x3 matrix has rank 3.", f)
	}
}

func w3Companion(kp int) generator.Problem {
	b := w3LNZ(-4, 4)
	c := w3LNZ(-4, 4)
	a, bb, cc, d := 0, 1, -c, -b
	f := map[string]string{}
	w3LSetMatrix(f, "", a, bb, cc, d)
	w3LI(f, "b", b)
	w3LI(f, "c", c)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For the companion matrix \\(C=%s\\), what is \\(\\det(C)\\)? (enter a number)", w3LMat(a, bb, cc, d)), strconv.Itoa(c),
			fmt.Sprintf("The determinant of \\(\\begin{pmatrix}0&1\\\\-c&-b\\end{pmatrix}\\) is \\(c=%d\\).", c), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the companion matrix \\(C=%s\\), what is \\(\\operatorname{tr}(C)\\)? (enter a number)", w3LMat(a, bb, cc, d)), strconv.Itoa(-b),
			fmt.Sprintf("The trace is \\(0+(-b)=%d\\).", -b), f)
	default:
		return w3L(kp, fmt.Sprintf("The companion matrix \\(C=%s\\) has characteristic polynomial \\(\\lambda^2+b\\lambda+c\\). What is its constant term \\(c\\)? (enter a number)", w3LMat(a, bb, cc, d)), strconv.Itoa(c),
			fmt.Sprintf("The constant term is the determinant, which is %d.", c), f)
	}
}

func w3Circulant(kp int) generator.Problem {
	a, b, c := w3LNZ(-3, 4), w3LNZ(-3, 4), w3LNZ(-3, 4)
	// Rows are cyclic shifts of (a,b,c).
	f := map[string]string{}
	w3LI(f, "a", a)
	w3LI(f, "b", b)
	w3LI(f, "c", c)
	switch kp {
	case 1:
		det := a*a*a + b*b*b + c*c*c - 3*a*b*c
		return w3L(kp, fmt.Sprintf("For the circulant matrix with first row \\((a,b,c)=(%d,%d,%d)\\), what is its determinant? (enter a number)", a, b, c), strconv.Itoa(det),
			fmt.Sprintf("A 3x3 circulant has determinant \\(a^3+b^3+c^3-3abc=%d\\).", det), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the circulant matrix with first row \\((a,b,c)=(%d,%d,%d)\\), what is its trace? (enter a number)", a, b, c), strconv.Itoa(3*a),
			fmt.Sprintf("Every diagonal entry of a circulant with first row \\((a,b,c)\\) is \\(a\\), so the trace is \\(3a=%d\\).", 3*a), f)
	default:
		return w3L(kp, fmt.Sprintf("For the circulant matrix with first row \\((a,b,c)=(%d,%d,%d)\\), what is the (1,2) entry? (enter a number)", a, b, c), strconv.Itoa(b),
			fmt.Sprintf("The first row is \\((a,b,c)\\), so the (1,2) entry is \\(b=%d\\).", b), f)
	}
}

func w3Householder(kp int) generator.Problem {
	if kp == 2 || kp == 3 {
		// The determinant and symmetry do not need the integer matrix. Draw a
		// wider set of direction vectors so both families have enough variety.
		vv := [][2]int{
			{1, 0}, {0, 1}, {1, 1}, {1, -1},
			{2, 1}, {2, -1}, {1, 2}, {1, -2},
			{3, 1}, {3, -1}, {1, 3}, {1, -3},
			{3, 2}, {3, -2}, {2, 3}, {2, -3},
			{4, 1}, {4, -1}, {1, 4}, {1, -4},
			{4, 3}, {4, -3}, {3, 4}, {3, -4},
		}
		v := vv[rand.Intn(len(vv))]
		f := map[string]string{}
		w3LI(f, "v1", v[0])
		w3LI(f, "v2", v[1])
		if kp == 2 {
			return w3L(kp, fmt.Sprintf("For the Householder reflection defined by \\(v=(%d,%d)\\), what is \\(\\det(H)\\)? (enter a number)", v[0], v[1]), "-1",
				"A Householder reflection is an orthogonal reflection, so its determinant is -1.", f)
		}
		return w3L(kp, fmt.Sprintf("For the Householder reflection defined by \\(v=(%d,%d)\\), is \\(H\\) symmetric? Enter 1 for yes and 0 for no.", v[0], v[1]), "1",
			"A Householder matrix is symmetric: \\(H^T=H\\). The answer is 1.", f)
	}
	// For the reflected-component family, keep v among the four directions
	// whose Householder matrix has integer entries.
	vs := [][2]int{{1, 1}, {1, -1}, {1, 0}, {0, 1}}
	v := vs[rand.Intn(len(vs))]
	vv := v[0]*v[0] + v[1]*v[1]
	// H = I - 2 vv^T/(v^Tv)
	h11 := 1 - 2*v[0]*v[0]/vv
	h12 := -2 * v[0] * v[1] / vv
	h21 := h12
	h22 := 1 - 2*v[1]*v[1]/vv
	x, y := w3LNZ(-5, 5), w3LNZ(-5, 5)
	f := map[string]string{}
	w3LSetMatrix(f, "", h11, h12, h21, h22)
	w3LI(f, "x", x)
	w3LI(f, "y", y)
	switch kp {
	case 1:
		ans := h11*x + h12*y
		return w3L(kp, fmt.Sprintf("The Householder reflection for \\(v=(%d,%d)\\) is \\(H=%s\\). What is the first component of \\(H(x,y)\\) for \\((x,y)=(%d,%d)\\)? (enter a number)", v[0], v[1], w3LMat(h11, h12, h21, h22), x, y), strconv.Itoa(ans),
			fmt.Sprintf("Multiply the first row of H by the vector: \\((%d)(%d)+(%d)(%d)=%d\\).", h11, x, h12, y, ans), f)
	case 2:
		return w3L(kp, fmt.Sprintf("The Householder reflection for \\(v=(%d,%d)\\) is \\(H=%s\\). What is \\(\\det(H)\\)? (enter a number)", v[0], v[1], w3LMat(h11, h12, h21, h22)), strconv.Itoa(w3LDet(h11, h12, h21, h22)),
			"A Householder reflection is an orthogonal reflection, so its determinant is -1.", f)
	default:
		is := 1
		if h12 != h21 {
			is = 0
		}
		w3LI(f, "symmetric", is)
		return w3L(kp, fmt.Sprintf("The Householder reflection for \\(v=(%d,%d)\\) is \\(H=%s\\). Is H symmetric? Enter 1 for yes and 0 for no.", v[0], v[1], w3LMat(h11, h12, h21, h22)), strconv.Itoa(is),
			fmt.Sprintf("A Householder matrix is symmetric: \\(H^T=H\\). The answer is %d.", is), f)
	}
}

func w3Givens(kp int) generator.Problem {
	pairs := [][3]int{
		{3, 4, 5}, {4, 3, 5},
		{5, 12, 13}, {12, 5, 13},
		{8, 15, 17}, {15, 8, 17},
		{7, 24, 25}, {24, 7, 25},
		{20, 21, 29}, {21, 20, 29},
		{9, 40, 41}, {40, 9, 41},
		{28, 45, 53}, {45, 28, 53},
	}
	pair := pairs[rand.Intn(len(pairs))]
	c, s, k := pair[0], pair[1], pair[2]
	if rand.Intn(2) == 0 {
		s = -s
	}
	x, y := w3LR(-4, 4), w3LR(-4, 4)
	f := map[string]string{}
	w3LI(f, "c", c)
	w3LI(f, "s", s)
	w3LI(f, "k", k)
	w3LI(f, "x", x)
	w3LI(f, "y", y)
	g11 := fmt.Sprintf("%d/%d", c, k)
	g12 := fmt.Sprintf("%d/%d", -s, k)
	g21 := fmt.Sprintf("%d/%d", s, k)
	g22 := g11
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For the Givens matrix \\(G=\\begin{pmatrix}%s&%s\\\\%s&%s\\end{pmatrix}\\), what is \\(\\det(G)\\)? (enter a number)", g11, g12, g21, g22), "1",
			"A Givens rotation is an orthogonal rotation, so its determinant is 1.", f)
	case 2:
		return w3L(kp, fmt.Sprintf("For the Givens matrix \\(G=\\begin{pmatrix}%s&%s\\\\%s&%s\\end{pmatrix}\\), what is the (2,1) entry of \\(%dG\\)? (enter a number)", g11, g12, g21, g22, k), strconv.Itoa(s),
			fmt.Sprintf("Multiply every entry by %d: the (2,1) entry becomes \\(s=%d\\).", k, s), f)
	default:
		ans := c*x - s*y
		return w3L(kp, fmt.Sprintf("For the Givens rotation \\(G=\\begin{pmatrix}%s&%s\\\\%s&%s\\end{pmatrix}\\) and \\((x,y)=(%d,%d)\\), what is %d times the first component of \\(G(x,y)\\)? (enter a number)", g11, g12, g21, g22, x, y, k), strconv.Itoa(ans),
			fmt.Sprintf("The first component is \\((%d)(%d)-(%d)(%d)\\) divided by %d, so %d times it is %d.", c, x, s, y, k, k, ans), f)
	}
}

func w3DeterminantLemma(kp int) generator.Problem {
	a, d := w3LR(1, 4), w3LR(1, 4)
	u1, u2 := w3LR(-4, 4), w3LR(-4, 4)
	v1, v2 := w3LR(-4, 4), w3LR(-4, 4)
	// v^T A^{-1} u = v1*u1/a + v2*u2/d. Choose multiples so it is integer.
	u1 *= a
	u2 *= d
	scalar := v1*u1/a + v2*u2/d
	f := map[string]string{}
	w3LI(f, "a", a)
	w3LI(f, "d", d)
	w3LI(f, "u1", u1)
	w3LI(f, "u2", u2)
	w3LI(f, "v1", v1)
	w3LI(f, "v2", v2)
	detA := a * d
	detNew := detA * (1 + scalar)
	switch kp {
	case 1:
		return w3L(kp, fmt.Sprintf("For \\(A=\\operatorname{diag}(%d,%d)\\), \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), what is \\(v^T A^{-1}u\\)? (enter a number)", a, d, u1, u2, v1, v2), strconv.Itoa(scalar),
			fmt.Sprintf("\\(A^{-1}=\\operatorname{diag}(1/%d,1/%d)\\), so the product is \\((%d)(%d)/(%d) + (%d)(%d)/(%d) = %d\\).", a, d, v1, u1, a, v2, u2, d, scalar), f)
	case 2:
		return w3L(kp, fmt.Sprintf("For \\(A=\\operatorname{diag}(%d,%d)\\), \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), what is \\(\\det(A+uv^T)\\)? (enter a number)", a, d, u1, u2, v1, v2), strconv.Itoa(detNew),
			fmt.Sprintf("The matrix determinant lemma gives \\(\\det(A)(1+v^TA^{-1}u)=%d\\cdot(1+(%d))=%d\\).", detA, scalar, detNew), f)
	default:
		diff := detNew - detA
		w3LI(f, "difference", diff)
		return w3L(kp, fmt.Sprintf("For \\(A=\\operatorname{diag}(%d,%d)\\), \\(u=(%d,%d)\\) and \\(v=(%d,%d)\\), by how much does the determinant change when the rank-one term \\(uv^T\\) is added? (enter a number)", a, d, u1, u2, v1, v2), strconv.Itoa(diff),
			fmt.Sprintf("The new determinant is %d and the old determinant is %d, so the change is %d.", detNew, detA, diff), f)
	}
}

// Keep strings imported for the facts helpers even if a future family stops
// using it directly.
var _ = strings.TrimSpace
