package numerical

// Wave 16 extends numerical analysis from 69 to 119 concepts. Five new
// subdomains carry eigenvalue algorithms, sparse methods, numerical PDEs,
// nonlinear systems and stochastic numerics. Every family is integer-valued so
// the numeric grader can check the answer, and every family draws at least two
// independent parameters so the variety floor holds at each ordinal.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

type w16Trio struct {
	question    string
	answer      int
	explanation string
}

func w16Limit(low, high int) int { return low + rand.Intn(high-low+1) }

func w16Mk(question string, answer int, explanation string) w16Trio {
	return w16Trio{question: question, answer: answer, explanation: explanation}
}

type w16Gen struct {
	families [3]func() w16Trio
}

func (g *w16Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	out := g.families[kp-1]()
	return generator.Problem{
		Question:    out.question,
		Answer:      strconv.Itoa(out.answer),
		Explanation: out.explanation,
		KPOrdinal:   kp,
	}
}

func registerWave16(reg *generator.Registry) {
	for id, families := range wave16FamiliesA {
		reg.Register(id, &w16Gen{families: families})
	}
	for id, families := range wave16FamiliesB {
		reg.Register(id, &w16Gen{families: families})
	}
	for id, families := range wave16FamiliesC {
		reg.Register(id, &w16Gen{families: families})
	}
}

// wave16FamiliesA covers eigenvalue algorithms and sparse methods.
var wave16FamiliesA = map[string][3]func() w16Trio{
	"numa.eigen.hessenberg": {
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A dense %d by %d matrix is reduced to upper Hessenberg form. How many entries strictly below the first subdiagonal become zero?", n, n),
				(n-1)*(n-2)/2,
				fmt.Sprintf("Rows 3 through %d have %d, %d, ..., 1 entries below the subdiagonal, which sum to %d.", n, n-2, n-3, (n-1)*(n-2)/2))
		},
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("An upper Hessenberg matrix has order %d. How many entries can be nonzero in total, counting the diagonal and the first subdiagonal?", n),
				n*(n+1)/2+(n-1),
				fmt.Sprintf("The upper triangle has %d entries and the subdiagonal adds %d, giving %d.", n*(n+1)/2, n-1, n*(n+1)/2+(n-1)))
		},
		func() w16Trio {
			n := w16Limit(4, 12)
			k := w16Limit(1, n-2)
			return w16Mk(fmt.Sprintf("A Householder reflector is applied to column %d of a %d by %d matrix during Hessenberg reduction. How many entries in that column does it zero?", k, n, n),
				n-k-1,
				fmt.Sprintf("The reflector keeps rows 1 through %d, so rows %d through %d are zeroed, which is %d entries.", k+1, k+2, n, n-k-1))
		},
	},
	"numa.eigen.rayleigh": {
		func() w16Trio {
			m := w16Limit(1, 8)
			a := w16Limit(1, 9)
			b := a + 5*m
			return w16Mk(fmt.Sprintf("A diagonal matrix has entries %d and %d. What is the Rayleigh quotient at x = (2, 1)? (enter a number)", a, b),
				a+m,
				fmt.Sprintf("The quotient is (4·%d + %d)/5 = (%d)/5 = %d.", a, b, 4*a+b, a+m))
		},
		func() w16Trio {
			m := w16Limit(1, 8)
			a := w16Limit(1, 9)
			b := a + 5*m
			return w16Mk(fmt.Sprintf("A diagonal matrix has entries %d and %d. What is the Rayleigh quotient at x = (0, 1)? (enter a number)", a, b),
				b,
				fmt.Sprintf("The second component is the only nonzero one, so the quotient picks out the second eigenvalue %d.", b))
		},
		func() w16Trio {
			r := w16Limit(1, 4)
			a := w16Limit(1, 9)
			b := a + 10*r
			return w16Mk(fmt.Sprintf("A diagonal matrix has entries %d and %d. What is the Rayleigh quotient at x = (1, 1)? (enter a number)", a, b),
				a+5*r,
				fmt.Sprintf("The quotient is the average (%d + %d)/2 = %d.", a, b, a+5*r))
		},
	},
	"numa.eigen.lanczos": {
		func() w16Trio {
			k := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("The Lanczos iteration runs for %d steps. How many matrix-vector products does it perform?", k), k,
				fmt.Sprintf("Each step needs exactly one product with the operator, so %d steps need %d products.", k, k))
		},
		func() w16Trio {
			k := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("The Lanczos iteration runs for %d steps. How many entries does the resulting tridiagonal matrix store?", k),
				2*k-1,
				fmt.Sprintf("A %d by %d tridiagonal matrix has %d diagonal entries and %d off-diagonal entries, giving %d.", k, k, k, 2*(k-1), 2*k-1))
		},
		func() w16Trio {
			k := w16Limit(2, 10)
			nnz := w16Limit(5, 40)
			return w16Mk(fmt.Sprintf("The operator has %d nonzeros and Lanczos runs for %d steps. How many multiply-add operations do the matrix-vector products cost in total?", nnz, k),
				k*nnz,
				fmt.Sprintf("One product costs %d multiply-adds, so %d products cost %d.", nnz, k, k*nnz))
		},
	},
	"numa.eigen.arnoldi": {
		func() w16Trio {
			k := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("The Arnoldi iteration completes %d steps. How many scalar products does full orthogonalisation perform in total?", k),
				k*(k+1)/2,
				fmt.Sprintf("Step j performs j scalar products, so the total is 1 + 2 + ... + %d = %d.", k, k*(k+1)/2))
		},
		func() w16Trio {
			k := w16Limit(2, 20)
			return w16Mk(fmt.Sprintf("The Arnoldi iteration completes %d steps on a matrix of order 200. How many basis vectors must be stored?", k), k,
				fmt.Sprintf("Each step adds one basis vector, so %d steps store %d vectors.", k, k))
		},
		func() w16Trio {
			k := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("Arnoldi runs for %d steps. How many entries can be nonzero in the resulting upper Hessenberg matrix?", k),
				k*(k+1)/2+(k-1),
				fmt.Sprintf("The upper triangle has %d entries and the subdiagonal adds %d, giving %d.", k*(k+1)/2, k-1, k*(k+1)/2+(k-1)))
		},
	},
	"numa.eigen.jacobi_rotation": {
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A symmetric matrix of order %d is diagonalised by Jacobi rotations. How many rotations does one full sweep perform?", n),
				n*(n-1)/2,
				fmt.Sprintf("Every unordered pair of distinct indices is visited once, and there are %d such pairs.", n*(n-1)/2))
		},
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A Jacobi rotation acts on rows p and q and on columns p and q of an %d by %d matrix. How many entries lie in those two rows and two columns?", n, n),
				4*n-4,
				fmt.Sprintf("Two full rows and two full columns contain 4n entries, and the four intersections are counted twice, so 4n - 4 = %d.", 4*n-4))
		},
		func() w16Trio {
			n := w16Limit(3, 12)
			s := w16Limit(1, 4)
			return w16Mk(fmt.Sprintf("A Jacobi sweep on a matrix of order %d performs %d rotations. How many rotations do %d sweeps perform?", n, n*(n-1)/2, s),
				s*n*(n-1)/2,
				fmt.Sprintf("Each sweep has %d rotations, so %d sweeps have %d.", n*(n-1)/2, s, s*n*(n-1)/2))
		},
	},
	"numa.eigen.deflation": {
		func() w16Trio {
			n := w16Limit(4, 12)
			d := w16Limit(1, n-2)
			return w16Mk(fmt.Sprintf("A matrix of order %d is deflated %d times. What is the order of the remaining problem?", n, d),
				n-d,
				fmt.Sprintf("Each deflation removes one dimension, so %d - %d = %d remains.", n, d, n-d))
		},
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A rank-one deflation updates a dense matrix of order %d. How many entries does the update touch?", n), n*n,
				fmt.Sprintf("The update is a full %d by %d matrix, so it touches %d entries.", n, n, n*n))
		},
		func() w16Trio {
			n := w16Limit(4, 12)
			k := w16Limit(1, n-1)
			return w16Mk(fmt.Sprintf("A matrix of order %d is to be reduced to a problem of order %d by deflation. How many eigenpairs must be deflated?", n, k),
				n-k,
				fmt.Sprintf("Deflation removes one eigenpair at a time, so %d - %d = %d pairs must be removed.", n, k, n-k))
		},
	},
	"numa.eigen.svd": {
		func() w16Trio {
			m := w16Limit(2, 9)
			n := w16Limit(2, 9)
			k := m
			if n < k {
				k = n
			}
			return w16Mk(fmt.Sprintf("A %d by %d matrix is factored by the singular value decomposition. How many singular values does it have?", m, n), k,
				fmt.Sprintf("The decomposition has min(%d, %d) = %d singular values.", m, n, k))
		},
		func() w16Trio {
			m := w16Limit(3, 9)
			n := w16Limit(3, 9)
			k := m
			if n < k {
				k = n
			}
			r := w16Limit(1, k-1)
			return w16Mk(fmt.Sprintf("A %d by %d matrix has rank %d. How many singular values are zero?", m, n, r), k-r,
				fmt.Sprintf("There are %d singular values and %d are nonzero, so %d are zero.", k, r, k-r))
		},
		func() w16Trio {
			m := w16Limit(3, 9)
			n := w16Limit(3, 9)
			k := m
			if n < k {
				k = n
			}
			return w16Mk(fmt.Sprintf("A %d by %d matrix is stored as a rank-%d truncated SVD, with U, the singular values and V. How many numbers are stored in total?", m, n, k),
				k*(m+n+1),
				fmt.Sprintf("U holds %d entries, the values hold %d and V holds %d, giving %d.", m*k, k, n*k, k*(m+n+1)))
		},
	},
	"numa.eigen.condition_svd": {
		func() w16Trio {
			smin := w16Limit(1, 9)
			c := w16Limit(2, 9)
			smax := smin * c
			return w16Mk(fmt.Sprintf("A matrix has largest singular value %d and smallest singular value %d. What is its 2-norm condition number?", smax, smin), c,
				fmt.Sprintf("The condition number is %d/%d = %d.", smax, smin, c))
		},
		func() w16Trio {
			e := w16Limit(1, 6)
			c := 1
			for i := 0; i < e; i++ {
				c *= 10
			}
			return w16Mk(fmt.Sprintf("A matrix has condition number %d. About how many decimal digits of accuracy are lost when the data carries a relative error of one unit in the last place?", c), e,
				fmt.Sprintf("Each factor of ten in the condition number costs about one decimal digit, and the ratio is 10^%d.", e))
		},
		func() w16Trio {
			smin := w16Limit(1, 6)
			c := w16Limit(2, 9)
			smax := smin * c
			return w16Mk(fmt.Sprintf("A matrix has largest singular value %d and condition number %d. What is its smallest singular value?", smax, c), smin,
				fmt.Sprintf("The smallest value is %d/%d = %d.", smax, c, smin))
		},
	},
	"numa.eigen.generalized": {
		func() w16Trio {
			n := w16Limit(3, 9)
			nullity := w16Limit(1, n-1)
			return w16Mk(fmt.Sprintf("A generalized problem has order %d and its second matrix has nullity %d. How many infinite eigenvalues does the pencil have?", n, nullity), nullity,
				fmt.Sprintf("Each null direction of the second matrix gives one infinite eigenvalue, so there are %d.", nullity))
		},
		func() w16Trio {
			n := w16Limit(2, 9)
			return w16Mk(fmt.Sprintf("A symmetric positive definite matrix of order %d is Cholesky factored while reducing a generalized problem. How many entries does the triangular factor store?", n),
				n*(n+1)/2,
				fmt.Sprintf("A lower triangular factor stores 1 + 2 + ... + %d = %d entries.", n, n*(n+1)/2))
		},
		func() w16Trio {
			n := w16Limit(2, 9)
			return w16Mk(fmt.Sprintf("A generalized problem of order %d is reduced to a standard eigenvalue problem. What is the order of the standard problem?", n), n,
				fmt.Sprintf("The Cholesky congruence preserves the order, so the standard problem is also of order %d.", n))
		},
	},
	"numa.eigen.sensitivity": {
		func() w16Trio {
			k := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("The left and right eigenvectors of an eigenvalue satisfy v^T u = 1/%d. What is the condition number of that eigenvalue?", k), k,
				fmt.Sprintf("The condition number is 1/|v^T u| = %d.", k))
		},
		func() w16Trio {
			k := w16Limit(2, 10)
			e := w16Limit(1, 6)
			return w16Mk(fmt.Sprintf("An eigenvalue has 1/|v^T u| = %d and the perturbation gives |v^T E u| = %d. What is the first-order shift |d lambda|?", k, e), e*k,
				fmt.Sprintf("The shift is |v^T E u|/|v^T u| = %d · %d = %d.", e, k, e*k))
		},
		func() w16Trio {
			eps := w16Limit(1, 5)
			k := w16Limit(2, 10)
			delta := eps * k
			return w16Mk(fmt.Sprintf("Two eigenvalues are separated by %d and the matrix perturbation has size %d. What is the largest eigenvalue condition number that still keeps the two resolvable?", delta, eps), k,
				fmt.Sprintf("Resolvability needs the separation to exceed the condition number times the perturbation, so the bound is %d/%d = %d.", delta, eps, k))
		},
	},
	"numa.sparse.csr": {
		func() w16Trio {
			nnz := w16Limit(5, 60)
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A matrix with %d nonzeros and %d rows is stored in compressed sparse row form. How many numbers are stored in total?", nnz, n),
				2*nnz+n+1,
				fmt.Sprintf("Values and column indices hold %d each and the row pointers hold %d, giving %d.", nnz, n+1, 2*nnz+n+1))
		},
		func() w16Trio {
			nnz := w16Limit(5, 60)
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A compressed sparse row matrix has %d nonzeros and %d rows. How many entries does the row-pointer array hold?", nnz, n), n+1,
				fmt.Sprintf("There is one pointer per row plus a final pointer, so %d + 1 = %d.", n, n+1))
		},
		func() w16Trio {
			m := w16Limit(3, 12)
			n := w16Limit(3, 12)
			nnz := w16Limit(2, m*n-1)
			return w16Mk(fmt.Sprintf("A %d by %d matrix has %d nonzeros. How many numbers are saved by compressed sparse row storage compared with a dense array?", m, n, nnz),
				m*n-(2*nnz+n+1),
				fmt.Sprintf("Dense storage holds %d numbers and compressed storage holds %d, so %d are saved.", m*n, 2*nnz+n+1, m*n-(2*nnz+n+1)))
		},
	},
	"numa.sparse.csc": {
		func() w16Trio {
			nnz := w16Limit(5, 60)
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A matrix with %d nonzeros and %d columns is stored in compressed sparse column form. How many numbers are stored in total?", nnz, n),
				2*nnz+n+1,
				fmt.Sprintf("Values and row indices hold %d each and the column pointers hold %d, giving %d.", nnz, n+1, 2*nnz+n+1))
		},
		func() w16Trio {
			n := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A compressed sparse column matrix has %d columns. How many entries does the column-pointer array hold?", n), n+1,
				fmt.Sprintf("One pointer per column plus a final pointer gives %d + 1 = %d.", n, n+1))
		},
		func() w16Trio {
			m := w16Limit(3, 12)
			n := w16Limit(3, 12)
			nnz := w16Limit(2, m*n-1)
			return w16Mk(fmt.Sprintf("A %d by %d matrix has %d nonzeros. How many entries does a dense array store that compressed sparse column storage omits?", m, n, nnz),
				m*n-nnz,
				fmt.Sprintf("Dense storage holds %d entries and only %d are nonzero, so %d zeros are omitted.", m*n, nnz, m*n-nnz))
		},
	},
	"numa.sparse.matvec": {
		func() w16Trio {
			nnz := w16Limit(5, 80)
			return w16Mk(fmt.Sprintf("A sparse matrix has %d nonzeros. How many multiply-add operations does one matrix-vector product cost?", nnz), nnz,
				fmt.Sprintf("Each nonzero contributes one multiplication and one addition, so the product costs %d multiply-adds.", nnz))
		},
		func() w16Trio {
			m := w16Limit(4, 20)
			return w16Mk(fmt.Sprintf("A sparse matrix-vector product multiplies a matrix with %d rows. How many entries does the output vector have?", m), m,
				fmt.Sprintf("Each row produces one output entry, so the output has %d entries.", m))
		},
		func() w16Trio {
			m := w16Limit(4, 12)
			n := w16Limit(4, 12)
			nnz := m * n / 4
			return w16Mk(fmt.Sprintf("A dense product on a %d by %d matrix costs %d multiply-adds. A sparse product on the same matrix costs %d. What percentage of the dense work does the sparse product use?", m, n, m*n, nnz),
				25,
				fmt.Sprintf("The ratio is %d/%d = 1/4, which is 25 percent.", nnz, m*n))
		},
	},
	"numa.sparse.fill_in": {
		func() w16Trio {
			original := w16Limit(10, 60)
			factor := original + w16Limit(5, 50)
			return w16Mk(fmt.Sprintf("A sparse matrix has %d nonzeros and its LU factors hold %d. How many entries are fill-in?", original, factor), factor-original,
				fmt.Sprintf("Fill is the factor count minus the original count: %d - %d = %d.", factor, original, factor-original))
		},
		func() w16Trio {
			d := w16Limit(3, 12)
			return w16Mk(fmt.Sprintf("A pivot has %d remaining neighbours. How many new couplings between pairs of neighbours does eliminating it create?", d), d*(d-1)/2,
				fmt.Sprintf("Every unordered pair of neighbours becomes coupled, and there are %d such pairs.", d*(d-1)/2))
		},
		func() w16Trio {
			original := w16Limit(2, 20) * 10
			fill := original * w16Limit(1, 4) / 2
			return w16Mk(fmt.Sprintf("A matrix has %d nonzeros and elimination creates %d fill entries. By what percentage does the factor storage grow?", original, fill), 100*fill/original,
				fmt.Sprintf("The growth is %d/%d = %d percent.", fill, original, 100*fill/original))
		},
	},
	"numa.sparse.ordering": {
		func() w16Trio {
			n := w16Limit(3, 15)
			return w16Mk(fmt.Sprintf("A fill-reducing ordering permutes a matrix of order %d. How many entries does the permutation vector store?", n), n,
				fmt.Sprintf("One new index is stored for each old index, so the vector holds %d entries.", n))
		},
		func() w16Trio {
			lo := w16Limit(1, 8)
			hi := lo + w16Limit(1, 12)
			return w16Mk(fmt.Sprintf("After reverse Cuthill-McKee ordering the smallest and largest nonzero row indices are %d and %d. What is the bandwidth?", lo, hi), hi-lo,
				fmt.Sprintf("The bandwidth is the largest offset from the diagonal, %d - %d = %d.", hi, lo, hi-lo))
		},
		func() w16Trio {
			d := w16Limit(4, 20)
			return w16Mk(fmt.Sprintf("A minimum-degree ordering inspects a graph and finds a vertex of degree %d. How many neighbours does the first pivot have?", d), d,
				fmt.Sprintf("The pivot is chosen for its degree, so it has exactly %d neighbours.", d))
		},
	},
	"numa.sparse.preconditioner": {
		func() w16Trio {
			n := w16Limit(10, 200)
			return w16Mk(fmt.Sprintf("A Jacobi preconditioner is applied to a residual vector of length %d. How many divisions does one application perform?", n), n,
				fmt.Sprintf("Each entry is divided by its diagonal entry, so the application performs %d divisions.", n))
		},
		func() w16Trio {
			n := w16Limit(10, 200)
			return w16Mk(fmt.Sprintf("A diagonal preconditioner for a problem of order %d stores only the diagonal. How many numbers does it store?", n), n,
				fmt.Sprintf("The diagonal has one entry per row, so it stores %d numbers.", n))
		},
		func() w16Trio {
			base := w16Limit(4, 20) * 4
			f := w16Limit(2, 4)
			return w16Mk(fmt.Sprintf("An iterative solver needs %d steps at condition number kappa. A preconditioner reduces the condition number by a factor of %d, and the step count scales with the square root of the condition number. How many steps are needed now?", base, f*f),
				base/f,
				fmt.Sprintf("The step count scales with the root of the condition number, so a factor %d in the condition gives a factor %d in the steps: %d/%d = %d.", f*f, f, base, f, base/f))
		},
	},
	"numa.sparse.sparse_lu": {
		func() w16Trio {
			n := w16Limit(10, 100)
			b := w16Limit(2, 8)
			return w16Mk(fmt.Sprintf("A banded matrix of order %d has bandwidth %d in both directions. Sparse LU costs about n times the bandwidth squared. How many operations is that?", n, b), n*b*b,
				fmt.Sprintf("The estimate is %d · %d squared = %d.", n, b, n*b*b))
		},
		func() w16Trio {
			nnz := w16Limit(20, 80)
			fill := w16Limit(5, 60)
			return w16Mk(fmt.Sprintf("Sparse LU starts from %d nonzeros and creates %d fill entries. How many entries do the factors hold?", nnz, fill), nnz+fill,
				fmt.Sprintf("The factors hold the original nonzeros plus the fill: %d + %d = %d.", nnz, fill, nnz+fill))
		},
		func() w16Trio {
			factor := w16Limit(30, 120)
			n := w16Limit(10, 25)
			return w16Mk(fmt.Sprintf("The LU factors of a matrix of order %d hold %d entries. The diagonal accounts for %d of them. How many multipliers are stored?", n, factor, n), factor-n,
				fmt.Sprintf("The triangular factors hold %d entries and %d sit on the diagonal, so %d multipliers remain.", factor, n, factor-n))
		},
	},
	"numa.sparse.multigrid": {
		func() w16Trio {
			levels := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("A multigrid hierarchy coarsens a grid by a factor of two until one unknown remains. How many levels does a hierarchy starting from 2^%d fine points have?", levels), levels,
				fmt.Sprintf("Halving %d times reaches one point, so the hierarchy has %d levels.", levels, levels))
		},
		func() w16Trio {
			levels := w16Limit(3, 10)
			n := 1 << levels
			return w16Mk(fmt.Sprintf("A multigrid hierarchy starts from %d fine points and halves the count at every level down to one. How many points are stored across all levels?", n), 2*n-1,
				fmt.Sprintf("The geometric sum %d + %d + ... + 1 equals %d.", n, n/2, 2*n-1))
		},
		func() w16Trio {
			n := w16Limit(8, 64) * 16
			return w16Mk(fmt.Sprintf("One multigrid V-cycle performs work proportional to the number of unknowns, with a constant factor of two across all levels. About how many work units does a cycle on %d unknowns cost?", n), 2*n,
				fmt.Sprintf("The cost is twice the finest level, so 2 · %d = %d.", n, 2*n))
		},
	},
	"numa.sparse.banded": {
		func() w16Trio {
			p := w16Limit(1, 5)
			q := w16Limit(1, 5)
			n := w16Limit(4, 20)
			return w16Mk(fmt.Sprintf("A banded matrix of order %d has lower bandwidth %d and upper bandwidth %d. How many entries does banded storage hold?", n, p, q),
				(p+q+1)*n,
				fmt.Sprintf("The band has %d + %d + 1 = %d diagonals, so %d · %d = %d entries.", p, q, p+q+1, p+q+1, n, (p+q+1)*n))
		},
		func() w16Trio {
			p := w16Limit(1, 4)
			q := w16Limit(1, 4)
			n := w16Limit(4, 15)
			return w16Mk(fmt.Sprintf("A %d by %d banded matrix has lower bandwidth %d and upper bandwidth %d. How many entries does banded storage omit compared with a dense array?", n, n, p, q),
				n*n-(p+q+1)*n,
				fmt.Sprintf("Dense storage holds %d and the band holds %d, so %d entries are omitted.", n*n, (p+q+1)*n, n*n-(p+q+1)*n))
		},
		func() w16Trio {
			n := w16Limit(10, 50)
			p := w16Limit(2, 6)
			q := w16Limit(2, 6)
			return w16Mk(fmt.Sprintf("Banded elimination on an order-%d matrix with bandwidths %d and %d costs about n times p times q. How many operations is that?", n, p, q), n*p*q,
				fmt.Sprintf("The estimate is %d · %d · %d = %d.", n, p, q, n*p*q))
		},
	},
	"numa.sparse.drop_tolerance": {
		func() w16Trio {
			nnz := w16Limit(2, 10) * 100
			drop := w16Limit(1, 5) * 10
			return w16Mk(fmt.Sprintf("An incomplete factorisation starts from %d nonzeros and drops %d percent of them. How many entries are kept?", nnz, drop),
				nnz*(100-drop)/100,
				fmt.Sprintf("Keeping %d percent of %d gives %d entries.", 100-drop, nnz, nnz*(100-drop)/100))
		},
		func() w16Trio {
			full := w16Limit(50, 200)
			kept := w16Limit(10, full-1)
			return w16Mk(fmt.Sprintf("A drop tolerance keeps %d entries out of %d. How many entries are discarded?", full, kept), full-kept,
				fmt.Sprintf("The discarded count is %d - %d = %d.", full, kept, full-kept))
		},
		func() w16Trio {
			full := w16Limit(4, 10) * 100
			kept := full / w16Limit(2, 4)
			return w16Mk(fmt.Sprintf("A full factor holds %d entries and the dropped factor holds %d. What percentage of the storage is saved?", full, kept), 100*(full-kept)/full,
				fmt.Sprintf("The saving is (%d - %d)/%d = %d percent.", full, kept, full, 100*(full-kept)/full))
		},
	},
}
