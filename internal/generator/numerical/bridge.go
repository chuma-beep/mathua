package numerical

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Concepts added by the 2026 content expansion. Each one chooses its
// knowledge point from GeneratorContext.KPOrdinal so the question matches
// the section of the lesson the learner is reading.
var numaBridgeKinds = map[string]string{
	"numa.linear.cholesky":            "cholesky",
	"numa.interp.hermite":             "hermite",
	"numa.linear.sor":                 "sor",
	"numa.linear.iter_convergence":    "iter_convergence",
	"numa.linear.qr_algorithm":        "qr_algorithm",
	"numa.linear.cg":                  "cg",
	"numa.fourier.fft":                "fft",
	"numa.integration.monte_carlo":    "monte_carlo",
	"numa.least_squares.gauss_newton": "gauss_newton",
}

type numaBridgeGen struct{ kind string }

func (g *numaBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "cholesky":
		return choleskyProblem(o)
	case "hermite":
		return hermiteProblem(o)
	case "sor":
		return sorProblem(o)
	case "iter_convergence":
		return iterConvergenceProblem(o)
	case "qr_algorithm":
		return qrAlgorithmProblem(o)
	case "cg":
		return conjugateGradientProblem(o)
	case "fft":
		return fftProblem(o)
	case "monte_carlo":
		return monteCarloProblem(o)
	default:
		return gaussNewtonProblem(o)
	}
}

func numaBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		Facts:       map[string]string{"answer": answer},
		KPOrdinal:   o,
	}
}

func numaBridgeDecimal(v float64, places int) string {
	s := strconv.FormatFloat(v, 'f', places, 64)
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	if s == "-0" {
		s = "0"
	}
	return s
}

func choleskyProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 1 + rand.Intn(6)
		b := 1 + rand.Intn(5)
		return numaBridgeProblem(o,
			fmt.Sprintf("A symmetric positive definite matrix has a = %d and b = %d in its top row [[a, b], [b, c]]. What is the (1,1) entry of its Cholesky factor L? (enter a number)", k*k, b),
			strconv.Itoa(k),
			fmt.Sprintf("The first diagonal entry is the square root of a, and sqrt(%d) = %d.", k*k, k))
	case 2:
		k := 1 + rand.Intn(6)
		l := 1 + rand.Intn(5)
		return numaBridgeProblem(o,
			fmt.Sprintf("For A = [[%d, %d], [%d, c]] with a suitable c, what is the (2,1) entry of the Cholesky factor L? (enter a number)", k*k, k*l, k*l),
			strconv.Itoa(l),
			fmt.Sprintf("The (2,1) entry is b divided by sqrt(a), so %d / %d = %d.", k*l, k, l))
	default:
		k := 1 + rand.Intn(6)
		l := 1 + rand.Intn(5)
		j := 1 + rand.Intn(5)
		c := l*l + j*j
		return numaBridgeProblem(o,
			fmt.Sprintf("For A = [[%d, %d], [%d, %d]], what is the (2,2) entry of the Cholesky factor L? (enter a number)", k*k, k*l, k*l, c),
			strconv.Itoa(j),
			fmt.Sprintf("The (2,2) entry is sqrt(c - l21^2) = sqrt(%d - %d) = %d.", c, l*l, j))
	}
}

func hermiteProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := rand.Intn(10)
		b := rand.Intn(10)
		c := rand.Intn(6)
		d := rand.Intn(6)
		y0 := 2 * a
		y1 := 2 * b
		m0 := 8 * c
		m1 := 8 * d
		ans := a + b + c - d
		return numaBridgeProblem(o,
			fmt.Sprintf("A cubic Hermite interpolant has y(0) = %d, y(1) = %d, y'(0) = %d and y'(1) = %d. What is y(1/2)? (enter a number)", y0, y1, m0, m1),
			strconv.Itoa(ans),
			fmt.Sprintf("At the midpoint the value is (y0+y1)/2 + (m0-m1)/8 = %d + %d = %d.", a+b, c-d, ans))
	case 2:
		n := 2 + rand.Intn(7)
		return numaBridgeProblem(o,
			fmt.Sprintf("Hermite interpolation matches a value and a first derivative at each of %d nodes. How many conditions does that impose? (enter a number)", n),
			strconv.Itoa(2*n),
			fmt.Sprintf("Each node gives two conditions, so %d nodes give %d conditions.", n, 2*n))
	default:
		n := 2 + rand.Intn(7)
		return numaBridgeProblem(o,
			fmt.Sprintf("Hermite interpolation matches a value and a first derivative at each of %d distinct nodes. What is the degree of the interpolating polynomial? (enter a number)", n),
			strconv.Itoa(2*n-1),
			fmt.Sprintf("There are %d conditions, so the polynomial has degree at most %d - 1 = %d.", 2*n, 2*n, 2*n-1))
	}
}

func sorProblem(o int) generator.Problem {
	switch o {
	case 1:
		old := rand.Intn(11)
		gs := old + 1 + rand.Intn(5)
		ws := []float64{0.5, 1.0, 1.5, 1.8}
		w := ws[rand.Intn(len(ws))]
		ans := (1-w)*float64(old) + w*float64(gs)
		return numaBridgeProblem(o,
			fmt.Sprintf("An SOR sweep has old value %d, Gauss-Seidel value %d and relaxation parameter omega = %s. What is the updated value? (enter a decimal)", old, gs, numaBridgeDecimal(w, 2)),
			numaBridgeDecimal(ans, 2),
			fmt.Sprintf("The update is (1-omega)x_old + omega x_gs = %s.", numaBridgeDecimal(ans, 2)))
	case 2:
		rhos := []float64{0, 0.25, 0.3, 0.4, 0.6, 0.7, 0.8, 0.9}
		rho := rhos[rand.Intn(len(rhos))]
		opt := 2 / (1 + math.Sqrt(1-rho*rho))
		return numaBridgeProblem(o,
			fmt.Sprintf("The Jacobi iteration matrix has spectral radius %s. What is the optimal SOR relaxation parameter? (enter a decimal)", numaBridgeDecimal(rho, 1)),
			numaBridgeDecimal(opt, 4),
			fmt.Sprintf("omega = 2 / (1 + sqrt(1 - rho^2)) = %s.", numaBridgeDecimal(opt, 4)))
	default:
		values := []float64{0.5, 1.0, 1.5, 1.9, 2.0, 2.5}
		rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })
		chosen := values[:4]
		count := 0
		parts := make([]string, 0, len(chosen))
		for _, v := range chosen {
			if v > 0 && v < 2 {
				count++
			}
			parts = append(parts, numaBridgeDecimal(v, 1))
		}
		return numaBridgeProblem(o,
			fmt.Sprintf("For how many of these relaxation parameters does SOR converge on every symmetric positive definite system? %s (enter a number)", strings.Join(parts, ", ")),
			strconv.Itoa(count),
			fmt.Sprintf("%d of the four values lie strictly between 0 and 2.", count))
	}
}

func iterConvergenceProblem(o int) generator.Problem {
	switch o {
	case 1:
		rhos := []float64{0.5, 0.1, 0.2, 0.25}
		rho := rhos[rand.Intn(len(rhos))]
		n := 2 + rand.Intn(4)
		err := math.Pow(rho, float64(n))
		return numaBridgeProblem(o,
			fmt.Sprintf("An iteration has spectral radius %s. By what factor is the error multiplied after %d iterations? (enter a decimal)", numaBridgeDecimal(rho, 2), n),
			numaBridgeDecimal(err, 6),
			fmt.Sprintf("The factor is rho^n = %s^%d = %s.", numaBridgeDecimal(rho, 2), n, numaBridgeDecimal(err, 6)))
	case 2:
		rhos := []float64{0.5, 0.1, 0.2, 0.25}
		rho := rhos[rand.Intn(len(rhos))]
		k := 2 + rand.Intn(5)
		need := int(math.Ceil(float64(k) / math.Log10(1/rho)))
		return numaBridgeProblem(o,
			fmt.Sprintf("An iteration has spectral radius %s. How many iterations are needed to reduce the error by a factor of 10^%d? (enter a number)", numaBridgeDecimal(rho, 2), k),
			strconv.Itoa(need),
			fmt.Sprintf("Solve rho^n <= 10^-%d, so n = ceil(%d / log10(1/rho)) = %d.", k, k, need))
	default:
		values := []float64{0.1, 0.2, 0.5, 0.8, 0.9, 1.0, 1.1, 1.5}
		rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })
		chosen := values[:5]
		count := 0
		parts := make([]string, 0, len(chosen))
		for _, v := range chosen {
			if v < 1 {
				count++
			}
			parts = append(parts, numaBridgeDecimal(v, 1))
		}
		return numaBridgeProblem(o,
			fmt.Sprintf("How many of these spectral radii guarantee convergence of the iteration? %s (enter a number)", strings.Join(parts, ", ")),
			strconv.Itoa(count),
			fmt.Sprintf("%d of the five radii are strictly less than 1.", count))
	}
}

func qrAlgorithmProblem(o int) generator.Problem {
	a := 1 + rand.Intn(9)
	d := 1 + rand.Intn(9)
	b := 1 + rand.Intn(9)
	large := a
	if d > large {
		large = d
	}
	switch o {
	case 1:
		return numaBridgeProblem(o,
			fmt.Sprintf("The QR algorithm is applied to the upper triangular matrix [[%d, %d], [0, %d]]. What is its larger eigenvalue? (enter a number)", a, b, d),
			strconv.Itoa(large),
			fmt.Sprintf("The eigenvalues of an upper triangular matrix are its diagonal entries, so the larger is max(%d, %d) = %d.", a, d, large))
	case 2:
		return numaBridgeProblem(o,
			fmt.Sprintf("The QR algorithm is applied to the upper triangular matrix [[%d, %d], [0, %d]]. What is the product of its eigenvalues? (enter a number)", a, b, d),
			strconv.Itoa(a*d),
			fmt.Sprintf("The product is the determinant, and for a triangular matrix that is %d x %d = %d.", a, d, a*d))
	default:
		return numaBridgeProblem(o,
			fmt.Sprintf("The QR algorithm is applied to the upper triangular matrix [[%d, %d], [0, %d]]. What is the sum of its eigenvalues? (enter a number)", a, b, d),
			strconv.Itoa(a+d),
			fmt.Sprintf("The sum is the trace, and for a triangular matrix that is %d + %d = %d.", a, d, a+d))
	}
}

func conjugateGradientProblem(o int) generator.Problem {
	switch o {
	case 1:
		a := 1 + rand.Intn(6)
		b := 1 + rand.Intn(6)
		p := rand.Intn(7) - 3
		q := rand.Intn(7) - 3
		if p == 0 && q == 0 {
			q = 1
		}
		ans := a*p*p + b*q*q
		return numaBridgeProblem(o,
			fmt.Sprintf("For A = [[%d, 0], [0, %d]] and x = (%d, %d), what is the quadratic form x^T A x? (enter a number)", a, b, p, q),
			strconv.Itoa(ans),
			fmt.Sprintf("The form is %d x %d^2 + %d x %d^2 = %d.", a, p, b, q, ans))
	case 2:
		a := 1 + rand.Intn(9)
		b := 1 + rand.Intn(9)
		if a == b {
			b = a%9 + 1
		}
		lo, hi := a, b
		if lo > hi {
			lo, hi = hi, lo
		}
		return numaBridgeProblem(o,
			fmt.Sprintf("A diagonal matrix has diagonal entries %d and %d. What is its condition number? (enter a decimal)", lo, hi),
			numaBridgeDecimal(float64(hi)/float64(lo), 4),
			fmt.Sprintf("The condition number is the largest eigenvalue divided by the smallest, %d / %d = %s.", hi, lo, numaBridgeDecimal(float64(hi)/float64(lo), 4)))
	default:
		n := 2 + rand.Intn(9)
		return numaBridgeProblem(o,
			fmt.Sprintf("Conjugate gradient is applied to a symmetric positive definite system of order %d in exact arithmetic. In at most how many iterations does it reach the solution? (enter a number)", n),
			strconv.Itoa(n),
			fmt.Sprintf("The method is exact after at most n = %d iterations.", n))
	}
}

func fftProblem(o int) generator.Problem {
	k := 1 + rand.Intn(6)
	n := 1 << k
	switch o {
	case 1:
		return numaBridgeProblem(o,
			fmt.Sprintf("A radix-2 fast Fourier transform has length %d. How many butterflies are performed in one stage? (enter a number)", n),
			strconv.Itoa(n/2),
			fmt.Sprintf("Each stage pairs the entries, so it uses n/2 = %d butterflies.", n/2))
	case 2:
		return numaBridgeProblem(o,
			fmt.Sprintf("A radix-2 fast Fourier transform has length %d. What is the total number of stages? (enter a number)", n),
			strconv.Itoa(k),
			fmt.Sprintf("The length is 2^%d, so there are %d stages.", k, k))
	default:
		return numaBridgeProblem(o,
			fmt.Sprintf("A radix-2 fast Fourier transform has length %d. How many complex multiplications does it use, up to the leading term n log2(n)? (enter a number)", n),
			strconv.Itoa(n*k),
			fmt.Sprintf("The leading count is n log2(n) = %d x %d = %d.", n, k, n*k))
	}
}

func monteCarloProblem(o int) generator.Problem {
	k := 2 + rand.Intn(5)
	n := k * k
	s := 1 + rand.Intn(5)
	switch o {
	case 1:
		return numaBridgeProblem(o,
			fmt.Sprintf("A Monte Carlo estimate uses %d samples and the integrand has standard deviation %d. What is the standard error? (enter a number)", n, s*k),
			strconv.Itoa(s),
			fmt.Sprintf("The standard error is sigma/sqrt(n) = %d / %d = %d.", s*k, k, s))
	case 2:
		return numaBridgeProblem(o,
			fmt.Sprintf("A Monte Carlo estimate uses %d samples and the integrand has standard deviation %d. What is the variance of the estimator? (enter a number)", n, s*k),
			strconv.Itoa(s*s),
			fmt.Sprintf("The variance is sigma^2/n = %d / %d = %d.", s*s*k*k, n, s*s))
	default:
		factor := 2 + rand.Intn(5)
		mult := factor * factor
		return numaBridgeProblem(o,
			fmt.Sprintf("Using %d times as many Monte Carlo samples reduces the standard error by what factor? (enter a number)", mult),
			strconv.Itoa(factor),
			fmt.Sprintf("The error scales as 1/sqrt(n), so multiplying n by %d divides the error by sqrt(%d) = %d.", mult, mult, factor))
	}
}

func gaussNewtonProblem(o int) generator.Problem {
	m := 3 + rand.Intn(10)
	n := 2 + rand.Intn(7)
	switch o {
	case 1:
		return numaBridgeProblem(o,
			fmt.Sprintf("A nonlinear least-squares problem has %d residuals and %d parameters. How many entries does the Jacobian have? (enter a number)", m, n),
			strconv.Itoa(m*n),
			fmt.Sprintf("The Jacobian has one row per residual and one column per parameter, so it has %d x %d = %d entries.", m, n, m*n))
	case 2:
		return numaBridgeProblem(o,
			fmt.Sprintf("A nonlinear least-squares problem has %d residuals and %d parameters. The Gauss-Newton step solves a normal system built from J^T J. What is its order? (enter a number)", m, n),
			strconv.Itoa(n),
			fmt.Sprintf("J^T J is n by n, and here n = %d.", n))
	default:
		return numaBridgeProblem(o,
			fmt.Sprintf("A nonlinear least-squares problem has %d residuals and %d parameters. Forming J^T J costs about m n^2 multiplications. What is that count? (enter a number)", m, n),
			strconv.Itoa(m*n*n),
			fmt.Sprintf("The cost is m n^2 = %d x %d^2 = %d.", m, n, m*n*n))
	}
}
