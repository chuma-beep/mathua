package linalg

// Wave 24 linear algebra: multilinear and functional structure, iterative
// solvers and numerical stability. Every family returns an integer and
// publishes the parameters it drew, so the test can recompute the answer.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave24(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"linalg.ml.tensor_space":         "tensor",
		"linalg.ml.symmetric_square":     "symsquare",
		"linalg.ml.exterior_square":      "extsquare",
		"linalg.ml.wedge":                "wedge",
		"linalg.ml.alternating":          "alternating",
		"linalg.ml.pfaffian":             "pfaffian",
		"linalg.ml.bilinear_rank":        "bilrank",
		"linalg.ml.quadratic_rank":       "quadrank",
		"linalg.ml.inertia":              "inertia",
		"linalg.ml.polarization":         "polarization",
		"linalg.func.dual":               "dual",
		"linalg.func.annihilator":        "annihilator",
		"linalg.func.quotient":           "quotient",
		"linalg.func.exact":              "exact",
		"linalg.func.adjoint":            "adjoint",
		"linalg.func.riesz":              "riesz",
		"linalg.func.invariant":          "invariant",
		"linalg.func.cyclic":             "cyclic",
		"linalg.func.minimal_polynomial": "minimal",
		"linalg.func.direct_sum":         "directsum",
		"linalg.itr.jacobi":              "jacobi",
		"linalg.itr.gauss_seidel":        "gauss",
		"linalg.itr.sor":                 "sor",
		"linalg.itr.residual":            "residual",
		"linalg.itr.convergence":         "convergence",
		"linalg.itr.krylov":              "krylov",
		"linalg.itr.arnoldi":             "arnoldi",
		"linalg.itr.gmres":               "gmres",
		"linalg.itr.cg":                  "cg",
		"linalg.itr.preconditioner":      "precond",
		"linalg.stab.forward_error":      "forward",
		"linalg.stab.backward_error":     "backward",
		"linalg.stab.machine_epsilon":    "epsilon",
		"linalg.stab.growth":             "growth",
		"linalg.stab.pivoting":           "pivoting",
		"linalg.stab.condition_error":    "conderr",
		"linalg.stab.residual_bound":     "resbound",
		"linalg.stab.componentwise":      "component",
		"linalg.stab.perturbation":       "perturb",
		"linalg.stab.rank_revealing":     "rankreveal",
	} {
		reg.Register(id, &w24Gen{kind: kind})
	}
}

type w24Gen struct{ kind string }

func (g *w24Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "tensor":
		return w24Tensor(kp)
	case "symsquare":
		return w24SymSquare(kp)
	case "extsquare":
		return w24ExtSquare(kp)
	case "wedge":
		return w24Wedge(kp)
	case "alternating":
		return w24Alternating(kp)
	case "pfaffian":
		return w24Pfaffian(kp)
	case "bilrank":
		return w24BilRank(kp)
	case "quadrank":
		return w24QuadRank(kp)
	case "inertia":
		return w24Inertia(kp)
	case "polarization":
		return w24Polarization(kp)
	case "dual":
		return w24Dual(kp)
	case "annihilator":
		return w24Annihilator(kp)
	case "quotient":
		return w24Quotient(kp)
	case "exact":
		return w24Exact(kp)
	case "adjoint":
		return w24Adjoint(kp)
	case "riesz":
		return w24Riesz(kp)
	case "invariant":
		return w24Invariant(kp)
	case "cyclic":
		return w24Cyclic(kp)
	case "minimal":
		return w24Minimal(kp)
	case "directsum":
		return w24DirectSum(kp)
	case "jacobi":
		return w24Jacobi(kp)
	case "gauss":
		return w24Gauss(kp)
	case "sor":
		return w24SOR(kp)
	case "residual":
		return w24Residual(kp)
	case "convergence":
		return w24Convergence(kp)
	case "krylov":
		return w24Krylov(kp)
	case "arnoldi":
		return w24Arnoldi(kp)
	case "gmres":
		return w24GMRES(kp)
	case "cg":
		return w24CG(kp)
	case "precond":
		return w24Precond(kp)
	case "forward":
		return w24Forward(kp)
	case "backward":
		return w24Backward(kp)
	case "epsilon":
		return w24Epsilon(kp)
	case "growth":
		return w24Growth(kp)
	case "pivoting":
		return w24Pivoting(kp)
	case "conderr":
		return w24CondErr(kp)
	case "resbound":
		return w24ResBound(kp)
	case "component":
		return w24Component(kp)
	case "perturb":
		return w24Perturb(kp)
	case "rankreveal":
		return w24RankReveal(kp)
	}
	return w24P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w24N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w24P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w24I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w24Binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	out := 1
	for i := 1; i <= k; i++ {
		out = out * (n - k + i) / i
	}
	return out
}

func w24Fact(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func w24Pow(a, b int) int {
	out := 1
	for i := 0; i < b; i++ {
		out *= a
	}
	return out
}

func w24Tensor(kp int) generator.Problem {
	m, n := w24N(2, 60), w24N(2, 60)
	f := map[string]string{}
	w24I(f, "m", m)
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A space of dimension %d is tensored with a space of dimension %d. What is the dimension of the tensor product?", m, n), strconv.Itoa(m*n), "Tensor dimensions multiply.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("How many basis elements does the tensor product of spaces of dimensions %d and %d have?", m, n), strconv.Itoa(m*n), "Basis pairs give the basis.", f)
	default:
		return w24P(kp, fmt.Sprintf("How many bilinear maps on a %d-dimensional space and a %d-dimensional space does the universal property convert into linear maps?", m, n), "1", "The universal property gives a unique linear map.", f)
	}
}

func w24SymSquare(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("What is the dimension of the symmetric square of a %d-dimensional space?", n), strconv.Itoa(n*(n+1)/2), "Count unordered pairs including repeats.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("How many basis elements of the symmetric square of a %d-dimensional space come from repeated basis vectors?", n), strconv.Itoa(n), "One repeated pair per basis vector.", f)
	default:
		return w24P(kp, fmt.Sprintf("What is the dimension of the symmetric square of a %d-dimensional space, counted as unordered pairs of basis vectors?", n), strconv.Itoa(n*(n+1)/2), "The formula is n(n+1)/2.", f)
	}
}

func w24ExtSquare(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("What is the dimension of the exterior square of a %d-dimensional space?", n), strconv.Itoa(n*(n-1)/2), "Count unordered pairs of distinct basis vectors.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("How many basis elements does the exterior square of a %d-dimensional space have, one per distinct pair?", n), strconv.Itoa(n*(n-1)/2), "The formula is n(n-1)/2.", f)
	default:
		return w24P(kp, fmt.Sprintf("What is the wedge product of a %d-dimensional vector with itself?", n), "0", "Alternation forces the square to vanish.", f)
	}
}

func w24Wedge(kp int) generator.Problem {
	n, k := w24N(2, 60), 0
	k = w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("How many basis 2-vectors does the exterior square of a %d-dimensional space have?", n), strconv.Itoa(w24Binom(n, 2)), "Count pairs of distinct indices.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("What is the dimension of the degree-%d exterior power of a %d-dimensional space?", k, n), strconv.Itoa(w24Binom(n, k)), "The dimension is C(n,k).", f)
	default:
		return w24P(kp, fmt.Sprintf("A wedge of %d vectors in a %d-dimensional space is tested for dependence. What value does the wedge take when two arguments coincide?", k, n), "0", "Repeated arguments make an alternating form vanish.", f)
	}
}

func w24Alternating(kp int) generator.Problem {
	n, k := w24N(2, 60), 0
	k = w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("How many independent alternating %d-forms exist on a %d-dimensional space?", k, n), strconv.Itoa(w24Binom(n, k)), "The dimension is C(n,k).", f)
	case 2:
		k = w24N(2, 12)
		q := w24N(2, 40)
		w24I(f, "k", k)
		w24I(f, "q", q)
		return w24P(kp, fmt.Sprintf("An alternating %d-form over a field of %d elements has how many orderings of its %d arguments?", k, q, k), strconv.Itoa(w24Fact(k)), "The arguments may be permuted in k factorial ways.", f)
	default:
		return w24P(kp, fmt.Sprintf("What is the dimension of the space of alternating %d-forms on a %d-dimensional space, when the form has top degree?", n, n), "1", "The top exterior power is one-dimensional.", f)
	}
}

func w24Pfaffian(kp int) generator.Problem {
	m := w24N(1, 12)
	f := map[string]string{}
	w24I(f, "m", m)
	double := 1
	for i := 1; i <= 2*m-1; i += 2 {
		double *= i
	}
	switch kp {
	case 1:
		k := w24N(1, 40)
		w24I(f, "k", k)
		return w24P(kp, fmt.Sprintf("A Pfaffian is defined for an alternating matrix of order 2k with k=%d. What is the order?", k), strconv.Itoa(2*k), "The order is twice k.", f)
	case 2:
		q := w24N(2, 40)
		w24I(f, "q", q)
		return w24P(kp, fmt.Sprintf("How many perfect matchings pair the %d indices of an alternating matrix of order %d over a field of %d elements?", 2*m, 2*m, q), strconv.Itoa(double), "The count is the double factorial (2m-1)!!.", f)
	default:
		q := w24N(2, 40)
		w24I(f, "q", q)
		return w24P(kp, fmt.Sprintf("The Pfaffian of an alternating matrix of order %d over a field of %d elements is squared. What power of the Pfaffian equals the determinant?", 2*m, q), "2", "The square of the Pfaffian is the determinant.", f)
	}
}

func w24BilRank(kp int) generator.Problem {
	n, r := w24N(2, 60), 0
	r = w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "r", r)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A bilinear form on a %d-dimensional space has rank %d. What is its rank?", n, r), strconv.Itoa(r), "The rank is given by the representing matrix.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A bilinear form of rank %d acts on a %d-dimensional space. What is its nullity?", r, n), strconv.Itoa(n-r), "Rank plus nullity equals the dimension.", f)
	default:
		return w24P(kp, fmt.Sprintf("A non-degenerate bilinear form on a %d-dimensional space has rank %d. By how much does a change of basis alter the rank?", n, r), "1", "Rank is invariant; enter 1 for unchanged.", f)
	}
}

func w24QuadRank(kp int) generator.Problem {
	n, r := w24N(2, 60), 0
	r = w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "r", r)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A quadratic form on a %d-dimensional space has rank %d. What is its rank?", n, r), strconv.Itoa(r), "The rank is the number of nonzero squares.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A quadratic form of rank %d is polarised into a symmetric bilinear form. What factor relates the two ranks?", r), "2", "Polarisation carries a factor of two.", f)
	default:
		return w24P(kp, fmt.Sprintf("A quadratic form of rank %d on a %d-dimensional space has how many zero diagonal directions?", r, n), strconv.Itoa(n-r), "The nullity is the dimension minus the rank.", f)
	}
}

func w24Inertia(kp int) generator.Problem {
	n := w24N(3, 20)
	p := w24N(1, n-1)
	q := w24N(1, n-p)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "p", p)
	w24I(f, "q", q)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A real quadratic form has %d positive and %d negative squares. What is its rank?", p, q), strconv.Itoa(p+q), "The rank is the sum of the nonzero counts.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A real quadratic form has %d positive and %d negative squares. What is its signature?", p, q), strconv.Itoa(p-q), "The signature subtracts the negative count.", f)
	default:
		return w24P(kp, fmt.Sprintf("Sylvester's law fixes %d positive squares for a form of rank %d. How many positive squares does every diagonal form have?", p, p+q), strconv.Itoa(p), "The positive index is invariant.", f)
	}
}

func w24Polarization(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("The polarization identity in dimension %d combines how many squared norms in the complex case?", n), "4", "Four norms with signs.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("The real polarization identity in dimension %d uses how many squared norms?", n), "2", "Two norms suffice over the reals.", f)
	default:
		return w24P(kp, fmt.Sprintf("How many symmetric bilinear forms on a %d-dimensional space correspond to one quadratic form?", n), "1", "Polarisation is a bijection away from characteristic two.", f)
	}
}

func w24Dual(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A vector space has dimension %d. What is the dimension of its dual?", n), strconv.Itoa(n), "The dual has the same finite dimension.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A basis of a %d-dimensional space has how many dual basis functionals?", n), strconv.Itoa(n), "One functional per basis vector.", f)
	default:
		return w24P(kp, fmt.Sprintf("A space of dimension %d is mapped canonically into its double dual. What is the dimension of the image?", n), strconv.Itoa(n), "The canonical map is an isomorphism in finite dimension.", f)
	}
}

func w24Annihilator(kp int) generator.Problem {
	n := w24N(2, 60)
	k := w24N(1, n-1)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A subspace of dimension %d sits inside a space of dimension %d. What is the dimension of its annihilator?", k, n), strconv.Itoa(n-k), "The annihilator has the codimension.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("The annihilator of the whole %d-dimensional space has what dimension?", n), "0", "Only the zero functional vanishes everywhere.", f)
	default:
		return w24P(kp, fmt.Sprintf("The annihilator of a %d-dimensional subspace of a %d-dimensional space is taken twice. What is the dimension of the double annihilator?", k, n), strconv.Itoa(k), "The double annihilator recovers the subspace.", f)
	}
}

func w24Quotient(kp int) generator.Problem {
	n := w24N(2, 60)
	k := w24N(1, n-1)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A %d-dimensional space is quotiented by a %d-dimensional subspace. What is the dimension of the quotient?", n, k), strconv.Itoa(n-k), "The quotient dimension is the codimension.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A quotient of a %d-dimensional space has dimension %d. What is the dimension of the subspace quotiented out?", n, k), strconv.Itoa(k), "The codimension is the quotient dimension.", f)
	default:
		return w24P(kp, fmt.Sprintf("A %d-dimensional space is quotiented by a %d-dimensional subspace. How many dimensions does the quotient map forget?", n, k), strconv.Itoa(k), "The kernel has dimension k.", f)
	}
}

func w24Exact(kp int) generator.Problem {
	a, b, c := w24N(1, 50), w24N(1, 50), w24N(1, 50)
	f := map[string]string{}
	w24I(f, "a", a)
	w24I(f, "b", b)
	w24I(f, "c", c)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A short exact sequence has middle dimension %d and end dimensions %d and %d. What is the alternating sum of dimensions?", b, a, c), "0", "Exactness makes the alternating sum vanish.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A short exact sequence has end dimensions %d and %d. What is the dimension of the middle space?", a, c), strconv.Itoa(a+c), "The middle dimension is the sum of the ends.", f)
	default:
		return w24P(kp, fmt.Sprintf("A short exact sequence has %d maps between its three nonzero spaces. How many maps does it contain?", a+b+c), "3", "Two inclusions and one quotient map.", f)
	}
}

func w24Adjoint(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A linear map on a %d-dimensional space has rank %d. What is the rank of its adjoint?", n, n), strconv.Itoa(n), "The adjoint has the same rank.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A map on a %d-dimensional space is represented by a matrix. How many entries does the adjoint matrix have?", n), strconv.Itoa(n*n), "The adjoint is represented by the conjugate transpose.", f)
	default:
		return w24P(kp, fmt.Sprintf("The kernel of the adjoint of a map on a %d-dimensional space is orthogonal to the image. What is the dimension of the kernel when the map is invertible?", n), "0", "An invertible map has full rank and zero kernel.", f)
	}
}

func w24Riesz(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A linear functional on a %d-dimensional inner product space is represented by how many vectors?", n), "1", "Riesz representation gives a unique vector.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("The Riesz map on a %d-dimensional inner product space has what dimension as its domain?", n), strconv.Itoa(n), "The domain is the space itself.", f)
	default:
		return w24P(kp, fmt.Sprintf("How many inner products on a %d-dimensional space correspond to the same functional through the Riesz map?", n), "1", "The representing vector is unique.", f)
	}
}

func w24Invariant(kp int) generator.Problem {
	n := w24N(2, 60)
	k := w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A linear map on a %d-dimensional space preserves a subspace of dimension %d. What is the dimension of the restriction's domain?", k, n), strconv.Itoa(k), "The restriction acts on the invariant subspace.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("An invariant subspace of dimension %d is spanned by how many eigenvectors when it is an eigenspace?", k), "1", "An eigenspace is the span of one eigenvector per distinct eigenvalue; for a single eigenvalue the line count is one.", f)
	default:
		return w24P(kp, fmt.Sprintf("A %d-dimensional invariant subspace sits in a %d-dimensional space. What is the dimension of a complementary invariant subspace in the semisimple case?", n, k), strconv.Itoa(n-k), "The complement has the remaining dimension.", f)
	}
}

func w24Cyclic(kp int) generator.Problem {
	n := w24N(2, 60)
	d := w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "d", d)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A cyclic subspace of a %d-dimensional space has dimension %d. What is its dimension?", n, d), strconv.Itoa(d), "The dimension is the degree of the first relation.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A cyclic subspace of dimension %d is represented by a companion block. What is the order of the block?", d), strconv.Itoa(d), "The block has the dimension of the subspace.", f)
	default:
		return w24P(kp, fmt.Sprintf("A cyclic vector in a %d-dimensional space generates a subspace of dimension %d. What is the dimension of the whole space?", n, d), strconv.Itoa(n), "The ambient dimension is n.", f)
	}
}

func w24Minimal(kp int) generator.Problem {
	n := w24N(2, 60)
	k := w24N(1, 50)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A matrix of order %d has a minimal polynomial. What is the largest possible degree?", n), strconv.Itoa(n), "The degree is at most the order.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A minimal polynomial has %d distinct roots. How many eigenvalues does the matrix have?", k), strconv.Itoa(k), "The roots are the eigenvalues.", f)
	default:
		return w24P(kp, fmt.Sprintf("A matrix of order %d has minimal polynomial degree %d. How many factors does the characteristic polynomial share with it?", n, k), "1", "The minimal polynomial divides the characteristic polynomial.", f)
	}
}

func w24DirectSum(kp int) generator.Problem {
	n, m := w24N(2, 60), w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "m", m)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A space of dimension %d is the direct sum of a space of dimension %d and a complement. What is the dimension of the complement?", n, m), strconv.Itoa(n-m), "Direct sum dimensions subtract.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("Two subspaces of dimensions %d and %d form a direct sum. What is the dimension of their intersection?", m, n), "0", "A direct sum has trivial intersection.", f)
	default:
		return w24P(kp, fmt.Sprintf("A direct sum of %d summands gives how many projection maps onto the summands?", m), strconv.Itoa(m), "One projection per summand.", f)
	}
}

func w24Jacobi(kp int) generator.Problem {
	n := w24N(2, 90)
	k := w24N(1, 90)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A Jacobi sweep on a %d by %d system updates how many coordinates?", n, n), strconv.Itoa(n), "One update per coordinate.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A Jacobi iteration performs %d sweeps. How many sweeps is that?", k), strconv.Itoa(k), "One matrix-vector product per sweep.", f)
	default:
		return w24P(kp, fmt.Sprintf("A Jacobi sweep on a %d by %d system uses how many matrix-vector products?", n, n), "1", "The whole sweep applies the matrix once.", f)
	}
}

func w24Gauss(kp int) generator.Problem {
	n := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A Gauss-Seidel sweep on a %d by %d system updates how many coordinates?", n, n), strconv.Itoa(n), "One update per coordinate.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A Gauss-Seidel sweep uses how many triangular solves on a %d by %d system?", n, n), "2", "One forward and one backward substitution are implicit in the sweep.", f)
	default:
		return w24P(kp, fmt.Sprintf("A Gauss-Seidel sweep on a %d by %d system uses how many matrix-vector products?", n, n), "1", "The triangular structure is applied once.", f)
	}
}

func w24SOR(kp int) generator.Problem {
	n := w24N(2, 90)
	k := w24N(1, 90)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("An SOR parameter is written as a fraction with denominator %d. What is the numerator of the optimal value formula for spectral radius zero?", n), "1", "With zero spectral radius the optimal omega is one.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("An SOR iteration performs %d sweeps. How many sweeps is that?", k), strconv.Itoa(k), "One relaxation sweep per step.", f)
	default:
		return w24P(kp, fmt.Sprintf("An SOR sweep on a %d by %d system updates how many coordinates?", n, n), strconv.Itoa(n), "One update per coordinate.", f)
	}
}

func w24Residual(kp int) generator.Problem {
	n := w24N(2, 90)
	m := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "m", m)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A residual is computed for a system with %d unknowns. How many entries does the residual vector have?", n), strconv.Itoa(n), "The residual matches the right-hand side.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A relative residual divides by the norm of a right-hand side with %d entries. How many entries does that norm take?", m), strconv.Itoa(m), "The norm runs over all entries.", f)
	default:
		return w24P(kp, fmt.Sprintf("A residual recursion is run for %d iterations on a %d by %d system. How many recurrences are updated per iteration?", n, n, n), "1", "One residual vector per iteration.", f)
	}
}

func w24Convergence(kp int) generator.Problem {
	k := w24N(2, 90)
	d := w24N(1, 6)
	f := map[string]string{}
	w24I(f, "k", k)
	w24I(f, "d", d)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A stationary iteration runs %d steps. How many applications of the iteration matrix is that?", k), strconv.Itoa(k), "One per step.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("Each iteration gains %d decimal digits and the method runs %d steps. How many digits are gained?", d, k), strconv.Itoa(d*k), "Multiply digits per step by steps.", f)
	default:
		return w24P(kp, fmt.Sprintf("A method has convergence factor below one and runs %d steps. How many error reductions occur?", k), strconv.Itoa(k), "One reduction per step.", f)
	}
}

func w24Krylov(kp int) generator.Problem {
	k := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A Krylov subspace of order %d is spanned by how many vectors before truncation?", k), strconv.Itoa(k), "One vector per power.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("The Krylov subspace of order %d is built from a matrix and a vector. What is the exponent of the last matrix power used?", k), strconv.Itoa(k-1), "The powers run from zero to k-1.", f)
	default:
		return w24P(kp, fmt.Sprintf("A Krylov method of order %d performs how many matrix-vector products to build its basis?", k), strconv.Itoa(k), "One product per basis vector.", f)
	}
}

func w24Arnoldi(kp int) generator.Problem {
	k := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("The Arnoldi process runs for %d steps. How many orthonormal basis vectors does it produce?", k), strconv.Itoa(k), "One per step.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("An Arnoldi matrix after %d steps is upper Hessenberg. How many nonzero subdiagonal entries does it have?", k), strconv.Itoa(k-1), "All but the last step produce a subdiagonal entry.", f)
	default:
		return w24P(kp, fmt.Sprintf("Step %d of Arnoldi orthogonalises against how many earlier basis vectors?", k), strconv.Itoa(k), "Every earlier vector is used.", f)
	}
}

func w24GMRES(kp int) generator.Problem {
	k := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("GMRES improves its residual at every step. How many residual minimisations does it solve in %d steps?", k), strconv.Itoa(k), "One per step.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("Full GMRES after %d steps stores how many basis vectors?", k), strconv.Itoa(k), "One per step.", f)
	default:
		return w24P(kp, fmt.Sprintf("A restarted GMRES keeps %d vectors per cycle. How many vectors does it store?", k), strconv.Itoa(k), "The restart length is the storage.", f)
	}
}

func w24CG(kp int) generator.Problem {
	k := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("Conjugate gradients takes %d steps. How many matrix-vector products does it use?", k), strconv.Itoa(k), "One per step.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("In exact arithmetic conjugate gradients terminates in at most n steps. How many steps for a system of order %d?", k), strconv.Itoa(k), "Termination is bounded by the order.", f)
	default:
		return w24P(kp, fmt.Sprintf("A conjugate-gradient step of order %d keeps how many search directions?", k), "1", "The short recurrence keeps one direction.", f)
	}
}

func w24Precond(kp int) generator.Problem {
	k := w24N(2, 90)
	f := map[string]string{}
	w24I(f, "k", k)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A preconditioner reduces the condition number from %d to one. By what factor is it reduced?", k), strconv.Itoa(k), "The factor is the original condition number.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A preconditioner is applied once per iteration for %d iterations. How many applications is that?", k), strconv.Itoa(k), "One per iteration.", f)
	default:
		return w24P(kp, fmt.Sprintf("An identity preconditioner on a %d by %d system changes the spectrum by what factor?", k, k), "1", "The identity leaves the spectrum unchanged.", f)
	}
}

func w24Forward(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A computed vector differs from the exact vector in %d entries. How many entries carry the absolute forward error?", n), strconv.Itoa(n), "Every entry can differ.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A relative forward error divides by the norm of a solution with %d entries.", n), strconv.Itoa(n), "The norm runs over all entries.", f)
	default:
		return w24P(kp, fmt.Sprintf("A forward error bound multiplies the condition number by the backward error for a problem of order %d. How many factors does the bound contain?", n), "2", "Condition and backward error.", f)
	}
}

func w24Backward(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A backward error perturbs %d data entries. How many entries are perturbed?", n), strconv.Itoa(n), "Every entry may move.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("For a problem of order %d, first-order backward analysis bounds the error by a constant times machine epsilon. What is the exponent of machine epsilon?", n), "1", "Enter 1 for first order in epsilon.", f)
	default:
		return w24P(kp, fmt.Sprintf("A backward error bound for a %d by %d system measures the perturbation in how many matrix entries?", n, n), strconv.Itoa(n*n), "Every entry can be perturbed.", f)
	}
}

func w24Epsilon(kp int) generator.Problem {
	b := w24N(30, 90)
	f := map[string]string{}
	w24I(f, "b", b)
	digits := int(float64(b) * 0.30103)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A floating-point format has %d significand bits. What is the exponent of two for the gap at one?", b), strconv.Itoa(b-1), "The gap at one is two to the minus (b-1).", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A format has %d significand bits. About how many decimal digits does it carry?", b), strconv.Itoa(digits), "Multiply bits by log base ten of two.", f)
	default:
		return w24P(kp, fmt.Sprintf("A format has %d significand bits. What is the exponent of two for the unit roundoff?", b), strconv.Itoa(b), "Unit roundoff is half the gap, so the exponent rises by one.", f)
	}
}

func w24Growth(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("Gaussian elimination without pivoting on an order-%d matrix has a worst-case growth factor of two to what exponent?", n), strconv.Itoa(n-1), "The worst case is two to the n minus one.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A growth factor of two to the power %d is reported. What is the exponent?", n-1), strconv.Itoa(n-1), "The exponent is n minus one.", f)
	default:
		return w24P(kp, fmt.Sprintf("Partial pivoting on an order-%d matrix keeps multipliers bounded by what integer?", n), "1", "Partial pivoting bounds multipliers by one.", f)
	}
}

func w24Pivoting(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("Partial pivoting on an order-%d matrix performs how many row swaps in the worst case?", n), strconv.Itoa(n-1), "One swap per elimination step.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("Complete pivoting on an order-%d matrix keeps the multipliers bounded by what integer?", n), "1", "Complete pivoting also bounds multipliers by one.", f)
	default:
		return w24P(kp, fmt.Sprintf("A permutation of %d rows can be chosen in how many ways?", n), strconv.Itoa(w24Fact(n)), "Count the permutations.", f)
	}
}

func w24CondErr(kp int) generator.Problem {
	k := w24N(2, 60)
	d := w24N(1, 12)
	f := map[string]string{}
	w24I(f, "k", k)
	w24I(f, "d", d)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A matrix has condition number %d. By what factor can a relative data perturbation be amplified?", k), strconv.Itoa(k), "The condition number is the amplification factor.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A condition number has %d decimal digits and the working precision has %d more digits. How many digits can the answer lose?", d, k), strconv.Itoa(d), "The condition number's digits are the loss.", f)
	default:
		return w24P(kp, fmt.Sprintf("A problem of order %d has condition number %d. How many digits of accuracy can be lost in the worst case?", k, d), strconv.Itoa(d), "The condition number controls the loss.", f)
	}
}

func w24ResBound(kp int) generator.Problem {
	n := w24N(2, 60)
	m := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "m", m)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A residual bound for a system with %d unknowns bounds how many entries of the error?", n), strconv.Itoa(n), "The bound covers the whole vector.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A backward stable solver controls the residual to first order in epsilon for a system with %d equations. What order is that?", m), "1", "Enter 1 for first order.", f)
	default:
		return w24P(kp, fmt.Sprintf("A residual bound multiplies the inverse norm by the residual for a %d by %d system. How many matrix norms enter?", n, n), "2", "The inverse norm and the residual norm.", f)
	}
}

func w24Component(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A componentwise bound tracks %d entries separately. How many entrywise bounds are computed?", n), strconv.Itoa(n), "One per entry.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A componentwise backward error measures relative perturbations of %d data entries. How many relative errors are measured?", n), strconv.Itoa(n), "One per entry.", f)
	default:
		return w24P(kp, fmt.Sprintf("A componentwise condition number respects the units of %d entries. How many scaling choices does row scaling add?", n), strconv.Itoa(n), "One scale per row.", f)
	}
}

func w24Perturb(kp int) generator.Problem {
	n := w24N(2, 60)
	f := map[string]string{}
	w24I(f, "n", n)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A normal matrix of order %d is perturbed by a matrix of norm one. By how much can each eigenvalue move at most?", n), "1", "Normal matrices move eigenvalues by at most the perturbation norm.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A matrix of order %d has how many eigenvalues counted with multiplicity?", n), strconv.Itoa(n), "The characteristic polynomial has degree n.", f)
	default:
		return w24P(kp, fmt.Sprintf("An eigenvector bound for a matrix of order %d divides by the spectral gap. How many eigenvalue gaps are considered?", n), strconv.Itoa(n), "One gap per eigenvalue.", f)
	}
}

func w24RankReveal(kp int) generator.Problem {
	n, r := w24N(2, 60), 0
	r = w24N(1, n)
	f := map[string]string{}
	w24I(f, "n", n)
	w24I(f, "r", r)
	switch kp {
	case 1:
		return w24P(kp, fmt.Sprintf("A rank-revealing factorisation of an order-%d matrix keeps %d large singular values. What is the numerical rank?", n, r), strconv.Itoa(r), "The numerical rank counts the large values.", f)
	case 2:
		return w24P(kp, fmt.Sprintf("A singular value of an order-%d matrix is compared with a tolerance. How many thresholds are applied?", n), "1", "One threshold separates rank from nullity.", f)
	default:
		return w24P(kp, fmt.Sprintf("An order-%d matrix of numerical rank %d is truncated. How many small singular values are dropped?", n, r), strconv.Itoa(n-r), "The dropped values are the remainder.", f)
	}
}
