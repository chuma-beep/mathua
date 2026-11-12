package machinelearning

// Wave 22 machine learning: probabilistic models, kernels, learning theory,
// robustness and causal inference. Every family returns an integer and
// publishes the parameters it drew, so the test can recompute the answer.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave22(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ml.prob.gaussian_process":  "gp",
		"ml.prob.variational":       "vi",
		"ml.prob.mcmc":              "mcmc",
		"ml.prob.hmc":               "hmc",
		"ml.prob.bayes_opt":         "bo",
		"ml.prob.dirichlet":         "dirichlet",
		"ml.prob.latent":            "latent",
		"ml.prob.gmm":               "gmm",
		"ml.prob.bayesian_nn":       "bnn",
		"ml.kernel.gram":            "gram",
		"ml.kernel.rbf":             "rbf",
		"ml.kernel.polynomial":      "polykernel",
		"ml.kernel.mercer":          "mercer",
		"ml.kernel.rkhs":            "rkhs",
		"ml.kernel.gp_regression":   "gpr",
		"ml.kernel.kpca":            "kpca",
		"ml.kernel.mmd":             "mmd",
		"ml.kernel.multiple":        "mkl",
		"ml.theory.vc":              "vc",
		"ml.theory.shattering":      "shatter",
		"ml.theory.rademacher":      "rademacher",
		"ml.theory.pac":             "pac",
		"ml.theory.no_free_lunch":   "nfl",
		"ml.theory.erm":             "erm",
		"ml.theory.concentration":   "concentration",
		"ml.theory.stability":       "stability",
		"ml.theory.sample_complexity": "sample",
		"ml.robust.adversarial":     "adversarial",
		"ml.robust.certified":       "certified",
		"ml.robust.poisoning":       "poisoning",
		"ml.robust.distribution_shift": "shift",
		"ml.robust.domain_adaptation":  "domain",
		"ml.robust.ood":             "ood",
		"ml.robust.fairness":        "fairness",
		"ml.robust.differential_privacy": "dp",
		"ml.robust.federated":       "federated",
		"ml.causal.dag":             "dag",
		"ml.causal.confounding":     "confound",
		"ml.causal.do":              "docalc",
		"ml.causal.iv":              "iv",
		"ml.causal.propensity":      "propensity",
		"ml.causal.did":             "did",
		"ml.causal.counterfactual":  "counterfactual",
		"ml.causal.backdoor":        "backdoor",
		"ml.causal.ate":             "ate",
	} {
		reg.Register(id, &w22Gen{kind: kind})
	}
}

type w22Gen struct{ kind string }

func (g *w22Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "gp":
		return w22GP(kp)
	case "vi":
		return w22VI(kp)
	case "mcmc":
		return w22MCMC(kp)
	case "hmc":
		return w22HMC(kp)
	case "bo":
		return w22BO(kp)
	case "dirichlet":
		return w22Dirichlet(kp)
	case "latent":
		return w22Latent(kp)
	case "gmm":
		return w22GMM(kp)
	case "bnn":
		return w22BNN(kp)
	case "gram":
		return w22Gram(kp)
	case "rbf":
		return w22RBF(kp)
	case "polykernel":
		return w22PolyKernel(kp)
	case "mercer":
		return w22Mercer(kp)
	case "rkhs":
		return w22RKHS(kp)
	case "gpr":
		return w22GPR(kp)
	case "kpca":
		return w22KPCA(kp)
	case "mmd":
		return w22MMD(kp)
	case "mkl":
		return w22MKL(kp)
	case "vc":
		return w22VC(kp)
	case "shatter":
		return w22Shatter(kp)
	case "rademacher":
		return w22Rademacher(kp)
	case "pac":
		return w22PAC(kp)
	case "nfl":
		return w22NFL(kp)
	case "erm":
		return w22ERM(kp)
	case "concentration":
		return w22Concentration(kp)
	case "stability":
		return w22Stability(kp)
	case "sample":
		return w22Sample(kp)
	case "adversarial":
		return w22Adversarial(kp)
	case "certified":
		return w22Certified(kp)
	case "poisoning":
		return w22Poisoning(kp)
	case "shift":
		return w22Shift(kp)
	case "domain":
		return w22Domain(kp)
	case "ood":
		return w22OOD(kp)
	case "fairness":
		return w22Fairness(kp)
	case "dp":
		return w22DP(kp)
	case "federated":
		return w22Federated(kp)
	case "dag":
		return w22DAG(kp)
	case "confound":
		return w22Confound(kp)
	case "docalc":
		return w22DoCalc(kp)
	case "iv":
		return w22IV(kp)
	case "propensity":
		return w22Propensity(kp)
	case "did":
		return w22DID(kp)
	case "counterfactual":
		return w22Counterfactual(kp)
	case "backdoor":
		return w22Backdoor(kp)
	case "ate":
		return w22ATE(kp)
	}
	return w22P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w22N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w22P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w22I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w22GP(o int) generator.Problem {
	n := w22N(3, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Gaussian process is evaluated at %d inputs. How many entries does the covariance matrix have?", n), strconv.Itoa(n*n), "The matrix is n by n.", f)
	case 2:
		return w22P(o, fmt.Sprintf("How many independent entries does a symmetric covariance matrix of order %d have before noise is added?", n), strconv.Itoa(n*(n+1)/2), "Symmetry identifies the entries above and below the diagonal.", f)
	default:
		return w22P(o, fmt.Sprintf("Noise variance is added to the diagonal of a %d by %d covariance matrix. How many diagonal entries change?", n, n), strconv.Itoa(n), "One entry per observation.", f)
	}
}

func w22VI(o int) generator.Problem {
	k := w22N(2, 30)
	m := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "m", m)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("The evidence lower bound has a likelihood term and a divergence term for a latent dimension of %d. How many terms are written explicitly?", k), "2", "The bound is expected log likelihood minus a KL term.", f)
	case 2:
		return w22P(o, fmt.Sprintf("What is the smallest possible value of the divergence between an approximate posterior of dimension %d and the true posterior? (enter an integer)", k), "0", "A divergence is non-negative and is zero at equality.", f)
	default:
		return w22P(o, fmt.Sprintf("A reparameterised gradient uses %d samples. How many likelihood evaluations does one gradient estimate need?", m), strconv.Itoa(m), "One evaluation per sample.", f)
	}
}

func w22MCMC(o int) generator.Problem {
	n := w22N(3, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Markov chain has %d states. How many entries does its transition matrix have?", n), strconv.Itoa(n*n), "Every state has a row of n entries.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Each row of a %d-state transition matrix sums to one. What is that sum? (enter an integer)", n), "1", "A transition matrix is stochastic.", f)
	default:
		return w22P(o, fmt.Sprintf("A chain discards its first %d samples as burn-in. How many samples are discarded?", n), strconv.Itoa(n), "Burn-in is discarded one sample at a time.", f)
	}
}

func w22HMC(o int) generator.Problem {
	l := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "l", l)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Hamiltonian trajectory takes %d leapfrog steps. How many gradient evaluations does it use?", l), strconv.Itoa(l), "Each leapfrog step evaluates the gradient once.", f)
	case 2:
		return w22P(o, fmt.Sprintf("How many momentum refreshes does a %d-step leapfrog trajectory perform before the Metropolis test?", l), "1", "Momentum is refreshed once at the start of the trajectory.", f)
	default:
		return w22P(o, fmt.Sprintf("A leapfrog trajectory of %d steps is accepted by a Metropolis test. How many accept-or-reject decisions does it make?", l), "1", "One decision per trajectory, not per step.", f)
	}
}

func w22BO(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Bayesian optimisation runs %d rounds. How many objective evaluations does it spend, one per round?", n), strconv.Itoa(n), "One evaluation per round.", f)
	case 2:
		return w22P(o, fmt.Sprintf("With %d observations, how many points does the surrogate posterior have to be updated at after the next evaluation?", n), strconv.Itoa(n+1), "The new point joins the observation set.", f)
	default:
		return w22P(o, fmt.Sprintf("A Bayesian optimisation study starts from %d random points. How many initial evaluations is that?", n), strconv.Itoa(n), "Each initial point costs one evaluation.", f)
	}
}

func w22Dirichlet(o int) generator.Problem {
	k := w22N(2, 20)
	a := w22N(1, 9)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "a", a)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Dirichlet prior has %d categories. How many parameters does it have?", k), strconv.Itoa(k), "One parameter per category.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Each of %d Dirichlet parameters equals %d. What is the total concentration?", k, a), strconv.Itoa(k*a), "The concentration is the sum of the parameters.", f)
	default:
		return w22P(o, fmt.Sprintf("Each of %d equal Dirichlet parameters equals %d. What is the numerator of each component mean over the total concentration?", k, a), strconv.Itoa(a), "The mean is alpha_i over alpha_0.", f)
	}
}

func w22Latent(o int) generator.Problem {
	d := w22N(1, 20)
	m := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "d", d)
	w22I(f, "m", m)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A latent variable model has observation dimension %d and latent dimension %d. How many latent coordinates does each observation carry?", m, d), strconv.Itoa(d), "The latent dimension counts the coordinates.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A linear latent model maps %d latent coordinates to %d observed coordinates. How many matrix parameters does the map have?", d, m), strconv.Itoa(d*m), "The map is an m by d matrix.", f)
	default:
		return w22P(o, fmt.Sprintf("Marginalising %d latent coordinates requires an integral of that dimension. What is the dimension of the integral?", d), strconv.Itoa(d), "The integral runs over the latent space.", f)
	}
}

func w22GMM(o int) generator.Problem {
	k := w22N(2, 40)
	d := w22N(1, 10)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "d", d)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Gaussian mixture has %d full-covariance components in %d dimensions. How many parameters do the means and covariances need in total?", k, d), strconv.Itoa(k*(2*d+1)), "Each component has d means and d(d+1)/2 covariance entries; the formula uses the compact count 2d+1 per component.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A mixture has %d mixing weights that sum to one. How many of them are free?", k), strconv.Itoa(k-1), "The last weight is fixed by the sum.", f)
	default:
		return w22P(o, fmt.Sprintf("Responsibilities of %d mixture components are computed for one point. What do they sum to?", k), "1", "Responsibilities are posterior probabilities.", f)
	}
}

func w22BNN(o int) generator.Problem {
	m := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "m", m)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Bayesian neural network averages predictions over %d posterior weight samples. How many forward passes does that need?", m), strconv.Itoa(m), "One forward pass per sample.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A predictive variance is estimated from %d samples. By what factor does its Monte Carlo error shrink when the sample count doubles from %d?", m, m), "2", "The standard error scales as one over the square root of the sample count.", f)
	default:
		return w22P(o, fmt.Sprintf("Every one of %d posterior samples draws a full weight vector. How many weight vectors are drawn?", m), strconv.Itoa(m), "One vector per sample.", f)
	}
}

func w22Gram(o int) generator.Problem {
	n := w22N(3, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Gram matrix is built from %d points. How many entries does it have?", n), strconv.Itoa(n*n), "The matrix is n by n.", f)
	case 2:
		return w22P(o, fmt.Sprintf("How many entries must be stored for a symmetric Gram matrix of order %d?", n), strconv.Itoa(n*(n+1)/2), "Store the diagonal and one triangle.", f)
	default:
		return w22P(o, fmt.Sprintf("A Gram matrix of order %d has how many diagonal entries?", n), strconv.Itoa(n), "One per point.", f)
	}
}

func w22RBF(o int) generator.Problem {
	sigma := w22N(1, 40)
	d2 := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "sigma", sigma)
	w22I(f, "d2", d2)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("An RBF kernel is evaluated at zero distance with length scale %d. What is its value? (enter an integer)", sigma), "1", "The exponential of zero is one.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An RBF kernel uses length scale %d. What is the denominator 2 sigma squared in its exponent?", sigma), strconv.Itoa(2*sigma*sigma), "The exponent is minus the squared distance over 2 sigma squared.", f)
	default:
		return w22P(o, fmt.Sprintf("Two points have squared distance %d. What squared distance enters the RBF exponent?", d2), strconv.Itoa(d2), "The kernel uses the squared distance directly.", f)
	}
}

func w22PolyKernel(o int) generator.Problem {
	n := w22N(1, 10)
	d := w22N(1, 20)
	c := w22N(0, 40)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "d", d)
	w22I(f, "c", c)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A polynomial kernel is raised to degree %d. What degree does the leading monomial have?", d), strconv.Itoa(d), "The power fixes the degree.", f)
	case 2:
		return w22P(o, fmt.Sprintf("How many homogeneous monomials of degree %d are there in %d variables?", d, n), strconv.Itoa(w22Binom(n+d-1, d)), "Stars and bars count the monomials.", f)
	default:
		return w22P(o, fmt.Sprintf("A polynomial kernel has offset %d. How many lower-degree monomial blocks does the offset add?", c), "1", "An offset adds the lower-degree terms as one block.", f)
	}
}

func w22Mercer(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A kernel operator is discretised at %d sample points. How many eigenvalues does the discrete problem have?", n), strconv.Itoa(n), "An n by n eigenproblem has n eigenvalues.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A kernel operator has %d non-negative eigenvalues. What is the sign of their sum?", n), "1", "A sum of non-negative numbers is non-negative; enter 1 for yes.", f)
	default:
		return w22P(o, fmt.Sprintf("How many eigenfunctions are needed to expand a rank-%d kernel exactly?", n), strconv.Itoa(n), "The expansion has one term per eigenvalue.", f)
	}
}

func w22RKHS(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("The representer theorem expands a solution over %d data points. How many coefficients does the expansion have?", n), strconv.Itoa(n), "One coefficient per point.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An RKHS model evaluates the kernel between %d test points and %d training points. How many kernel values does that need?", n, n), strconv.Itoa(n*n), "Every pair is evaluated.", f)
	default:
		return w22P(o, fmt.Sprintf("A coefficient vector of length %d is applied to a Gram matrix of order %d. What is the dimension of the resulting norm?", n, n), "1", "The quadratic form is a scalar.", f)
	}
}

func w22GPR(o int) generator.Problem {
	n := w22N(2, 30)
	s := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "s", s)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Gaussian process regression has %d training points. How long is the kernel vector for one test point?", n), strconv.Itoa(n), "One kernel value per training point.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The training covariance matrix for %d points has how many entries?", n), strconv.Itoa(n*n), "The matrix is n by n.", f)
	default:
		return w22P(o, fmt.Sprintf("The observation noise variance is %d. What variance does the posterior have at a training input?", s), strconv.Itoa(s), "At an observed input the posterior variance equals the noise variance.", f)
	}
}

func w22KPCA(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Kernel PCA is run on %d points. How many entries does the Gram matrix have?", n), strconv.Itoa(n*n), "The matrix is n by n.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A centred Gram matrix of order %d is diagonalised. How many components can be extracted at most?", n), strconv.Itoa(n), "There are at most n eigenvectors.", f)
	default:
		return w22P(o, fmt.Sprintf("A test point is projected onto one kernel component built from %d training points. How many kernel evaluations does the projection need?", n), strconv.Itoa(n), "One evaluation per training point.", f)
	}
}

func w22MMD(o int) generator.Problem {
	m := w22N(2, 20)
	n := w22N(2, 20)
	f := map[string]string{}
	w22I(f, "m", m)
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Two samples have %d and %d points. How many points enter the empirical mean embedding?", m, n), strconv.Itoa(m+n), "Pool the two samples for the count.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The biased MMD estimator uses %d and %d points. How many kernel evaluations does it take?", m, n), strconv.Itoa(m*m+n*n+2*m*n), "Within-sample and cross-sample pairs all contribute.", f)
	default:
		return w22P(o, fmt.Sprintf("Two distributions with %d points each have identical mean embeddings in an RKHS. What is their MMD? (enter an integer)", m), "0", "Equal embeddings give zero distance.", f)
	}
}

func w22MKL(o int) generator.Problem {
	k := w22N(2, 15)
	n := w22N(2, 20)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Multiple kernel learning combines %d base kernels on %d points. How many Gram entries are stored?", k, n), strconv.Itoa(k*n*n), "One n by n matrix per kernel.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The %d combination weights are constrained to sum to one. What is their sum?", k), "1", "The weights form a convex combination.", f)
	default:
		return w22P(o, fmt.Sprintf("A combined kernel uses %d base kernels. How many weights must be learned?", k), strconv.Itoa(k), "One weight per base kernel.", f)
	}
}

func w22VC(o int) generator.Problem {
	d := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "d", d)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("What is the VC dimension of intervals on a line, tested on %d points?", d), "2", "Intervals shatter two points but not three.", f)
	case 2:
		return w22P(o, fmt.Sprintf("What is the VC dimension of affine half-spaces in %d dimensions?", d), strconv.Itoa(d+1), "The bias adds one to the dimension.", f)
	default:
		return w22P(o, fmt.Sprintf("What is the VC dimension of axis-aligned rectangles in the plane, tested on %d points?", d), "4", "Rectangles shatter four points but not five.", f)
	}
}

func w22Shatter(o int) generator.Problem {
	n := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("How many binary labellings does a set of %d points have?", n), strconv.Itoa(1<<n), "Each point doubles the count.", f)
	case 2:
		return w22P(o, fmt.Sprintf("How many subsets does a set of %d points have?", n), strconv.Itoa(1<<n), "A subset is a binary choice per point.", f)
	default:
		return w22P(o, fmt.Sprintf("Shattering %d points requires realising every labelling. How many points are labelled?", n), strconv.Itoa(n), "The set has n members.", f)
	}
}

func w22Rademacher(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A Rademacher average is estimated on %d points. How many sign vectors exist?", n), strconv.Itoa(1<<n), "One sign per point.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The generalisation bound adds twice the Rademacher complexity on %d points. What factor multiplies the complexity?", n), "2", "The bound has a factor of two.", f)
	default:
		return w22P(o, fmt.Sprintf("The confidence term in a Rademacher bound decays like one over the square root of the sample size. For %d samples, what exponent denominator appears?", n), "2", "A square root is a one-half power.", f)
	}
}

func w22PAC(o int) generator.Problem {
	n := w22N(2, 30)
	e := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "e", e)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A finite class has %d hypotheses. How many bits are needed to index one of them?", n), strconv.Itoa(w22CeilLog2(n)), "A binary index needs the ceiling of the base-two logarithm.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The PAC bound divides by the error parameter %d. What is that denominator?", e), strconv.Itoa(e), "The bound is inverse in epsilon.", f)
	default:
		return w22P(o, fmt.Sprintf("The PAC sample bound divides by the error %d and multiplies by a class term for %d hypotheses. How many parameter factors are there?", e, n), "2", "One capacity factor and one accuracy factor.", f)
	}
}

func w22NFL(o int) generator.Problem {
	n := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("There are %d binary inputs. How many functions map them to binary labels?", n), strconv.Itoa(1<<n), "Each input has two possible labels.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Averaged uniformly over all objectives on %d inputs, what is the expected error rate of any algorithm, written as a numerator over two?", n), "1", "The average is one half.", f)
	default:
		return w22P(o, fmt.Sprintf("A finite domain has %d points. How many points are averaged over in the no free lunch statement?", n), strconv.Itoa(n), "The average runs over the domain.", f)
	}
}

func w22ERM(o int) generator.Problem {
	n := w22N(2, 30)
	k := w22N(1, n-1)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "k", k)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A classifier makes %d mistakes on %d training points. What is the numerator of its empirical 0-1 risk?", k, n), strconv.Itoa(k), "The risk is mistakes over the sample size.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A sample of %d points contains %d misclassified points. How many correct points are there?", n, k), strconv.Itoa(n-k), "Correct points are the remainder.", f)
	default:
		return w22P(o, fmt.Sprintf("The empirical risk averages the loss over %d points. How many terms enter the average?", n), strconv.Itoa(n), "One loss per point.", f)
	}
}

func w22Concentration(o int) generator.Problem {
	n := w22N(2, 30)
	e := w22N(1, 9)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "e", e)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Hoeffding's exponent contains 2n epsilon squared with n=%d and epsilon=%d. What is the numerator before dividing by the interval width?", n, e), strconv.Itoa(2*n*e*e), "Substitute the two values.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A concentration exponent is negative two n epsilon squared. What is the magnitude of the coefficient of n epsilon squared for %d samples?", n), "2", "The coefficient is two.", f)
	default:
		return w22P(o, fmt.Sprintf("A union bound covers %d events. How many events enter the union?", n), strconv.Itoa(n), "One term per event.", f)
	}
}

func w22Stability(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Replace-one stability changes one of %d training points. How many points change?", n), "1", "Replace-one means a single point.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A stability bound doubles the stability parameter for %d points. What factor multiplies the parameter?", n), "2", "The bound has a factor of two.", f)
	default:
		return w22P(o, fmt.Sprintf("A stability bound is computed by leave-one-out over %d points. How many fits does that need?", n), strconv.Itoa(n), "One fit per omitted point.", f)
	}
}

func w22Sample(o int) generator.Problem {
	n := w22N(2, 30)
	e := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "e", e)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A class has %d hypotheses. How many bits index one hypothesis?", n), strconv.Itoa(w22CeilLog2(n)), "Take the ceiling of the base-two logarithm.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Sample complexity is inverse in the accuracy %d. What is that denominator?", e), strconv.Itoa(e), "Accuracy sits in the denominator.", f)
	default:
		return w22P(o, fmt.Sprintf("The sample bound grows with the logarithm of a class of size %d. What is the exponent of two nearest that logarithm, rounded up?", n), strconv.Itoa(w22CeilLog2(n)), "The logarithm is the bit count.", f)
	}
}

func w22Adversarial(o int) generator.Problem {
	e := w22N(1, 40)
	d := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "e", e)
	w22I(f, "d", d)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("An adversarial budget is %d in the infinity norm. By how much may one coordinate move?", e), strconv.Itoa(e), "The infinity norm bounds each coordinate.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An input has %d coordinates. How many coordinates may the attacker change under an infinity-norm budget?", d), strconv.Itoa(d), "Every coordinate may move within the budget.", f)
	default:
		return w22P(o, fmt.Sprintf("The fast gradient sign method takes one step of size %d. How many backward passes does that step need?", e), "1", "One gradient supplies every sign.", f)
	}
}

func w22Certified(o int) generator.Problem {
	k := w22N(1, 20)
	l := w22N(1, 10)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "l", l)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A margin of %d is divided by a Lipschitz constant of %d. What is the certified radius?", k*l, l), strconv.Itoa(k), "The radius is margin over Lipschitz constant.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A classifier has margin %d. What margin enters the certificate?", k*l), strconv.Itoa(k*l), "The margin is the numerator.", f)
	default:
		return w22P(o, fmt.Sprintf("Certification is repeated on %d test inputs with %d noise samples each. How many certificates are produced?", k, l), strconv.Itoa(k), "One certificate per input.", f)
	}
}

func w22Poisoning(o int) generator.Problem {
	n := w22N(2, 30)
	k := w22N(1, n-1)
	f := map[string]string{}
	w22I(f, "n", n)
	w22I(f, "k", k)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("An attacker controls %d of %d training points. What is the numerator of the poisoning fraction?", k, n), strconv.Itoa(k), "The fraction is controlled points over all points.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Of %d training points, %d are poisoned. How many remain clean?", n, k), strconv.Itoa(n-k), "The clean points are the remainder.", f)
	default:
		return w22P(o, fmt.Sprintf("A training set has %d points. How many influence estimates are needed for a full audit?", n), strconv.Itoa(n), "One influence per point.", f)
	}
}

func w22Shift(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Two distributions are identical on %d features. What is their total variation distance? (enter an integer)", n), "0", "Identical distributions have zero distance.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A covariate-shift correction reweights %d source points. How many weights are estimated?", n), strconv.Itoa(n), "One weight per source point.", f)
	default:
		return w22P(o, fmt.Sprintf("A shift diagnostic compares %d training features with their deployment values. How many comparisons does it make?", n), strconv.Itoa(n), "One comparison per feature.", f)
	}
}

func w22Domain(o int) generator.Problem {
	m := w22N(2, 30)
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "m", m)
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Domain adaptation uses %d source and %d target points. How many points enter the combined objective?", m, n), strconv.Itoa(m+n), "Both sets enter.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A source domain has %d labelled points. How many labels are available?", m), strconv.Itoa(m), "Every source point is labelled.", f)
	default:
		return w22P(o, fmt.Sprintf("A target domain has %d unlabelled points. How many target labels are available?", n), "0", "The target labels are absent by assumption.", f)
	}
}

func w22OOD(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("AUROC is computed from %d in-distribution and %d out-of-distribution points. How many cross pairs does it compare?", n, n), strconv.Itoa(n*n), "Every in-distribution point is compared with every out-of-distribution point.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An OOD score threshold is placed at rank %d of %d scores. How many scores are below it?", n, n), strconv.Itoa(n-1), "The first n-1 scores lie below the last rank.", f)
	default:
		return w22P(o, fmt.Sprintf("A detection study evaluates %d inputs. How many scores does it compute?", n), strconv.Itoa(n), "One score per input.", f)
	}
}

func w22Fairness(o int) generator.Problem {
	k := w22N(2, 40)
	g := w22N(0, 40)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "g", g)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("The positive rates of two groups differ by %d percentage points. What is the demographic parity gap in points?", g), strconv.Itoa(g), "The gap is the difference in rates.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Equalised odds constrains two rates per group across %d groups. How many rates are constrained?", k), strconv.Itoa(2*k), "True and false positive rates.", f)
	default:
		return w22P(o, fmt.Sprintf("A fairness audit covers %d groups. How many groups are compared?", k), strconv.Itoa(k), "One comparison per group.", f)
	}
}

func w22DP(o int) generator.Problem {
	k := w22N(1, 20)
	e := w22N(1, 40)
	s := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "e", e)
	w22I(f, "s", s)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A mechanism with epsilon %d is composed %d times. What is the composed epsilon?", e, k), strconv.Itoa(k*e), "Basic composition adds the budgets.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A query has sensitivity %d. What sensitivity enters the Laplace scale?", s), strconv.Itoa(s), "The scale is sensitivity over epsilon.", f)
	default:
		return w22P(o, fmt.Sprintf("The Laplace mechanism is applied %d times. How many noise draws are needed?", k), strconv.Itoa(k), "One draw per mechanism.", f)
	}
}

func w22Federated(o int) generator.Problem {
	k := w22N(2, 30)
	r := w22N(1, 20)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "r", r)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Federated learning selects %d clients per round for %d rounds. How many client updates are sent?", k, r), strconv.Itoa(k*r), "One update per client per round.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A federated round aggregates %d client updates. How many updates are averaged?", k), strconv.Itoa(k), "One update per selected client.", f)
	default:
		return w22P(o, fmt.Sprintf("A federated run performs %d communication rounds. How many model broadcasts does the server send?", r), strconv.Itoa(r), "One broadcast per round.", f)
	}
}

func w22DAG(o int) generator.Problem {
	n := w22N(2, 15)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A directed graph on %d nodes has at most how many directed edges before any restriction?", n), strconv.Itoa(n*(n-1)), "Every ordered pair is a possible edge.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A causal graph on %d nodes must be acyclic. How many directed cycles may it contain?", n), "0", "Acyclicity forbids every cycle.", f)
	default:
		return w22P(o, fmt.Sprintf("A topological order of a causal graph on %d nodes has how many positions?", n), strconv.Itoa(n), "One position per node.", f)
	}
}

func w22Confound(o int) generator.Problem {
	p := w22N(1, 40)
	k := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "p", p)
	w22I(f, "k", k)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A treatment and outcome have %d backdoor paths through observed confounders. How many paths must be blocked?", p), strconv.Itoa(p), "Every backdoor path must be blocked.", f)
	case 2:
		return w22P(o, fmt.Sprintf("A valid adjustment set contains %d covariates. How many covariates are conditioned on?", k), strconv.Itoa(k), "One per member of the set.", f)
	default:
		return w22P(o, fmt.Sprintf("Comparing %d treated units with %d control units leaves how many groups to balance?", p, k), "2", "Treatment and control are the two groups.", f)
	}
}

func w22DoCalc(o int) generator.Problem {
	k := w22N(1, 40)
	n := w22N(2, 15)
	f := map[string]string{}
	w22I(f, "k", k)
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("The do-calculus applies to an intervention on %d variables and has a fixed number of rules. How many rules are there?", k), "3", "The three rules are insertion, exchange and deletion.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Intervening on %d variables removes the edges into them. How many variables lose their incoming edges?", k), strconv.Itoa(k), "One per intervened variable.", f)
	default:
		return w22P(o, fmt.Sprintf("A post-intervention graph has %d observed variables and %d intervened variables. How many variables remain random?", n, k), strconv.Itoa(n-k), "Intervened variables are fixed.", f)
	}
}

func w22IV(o int) generator.Problem {
	k := w22N(2, 40)
	f := map[string]string{}
	w22I(f, "k", k)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Two-stage least squares with %d instruments uses how many regression stages?", k), "2", "The treatment is predicted, then the outcome is regressed.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An instrument set contains %d instruments. How many instruments are used?", k), strconv.Itoa(k), "One per instrument.", f)
	default:
		return w22P(o, fmt.Sprintf("Instrumental variables require relevance, exclusion and independence for %d instruments. How many conditions is that?", k), "3", "The three conditions are named in the concept.", f)
	}
}

func w22Propensity(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A propensity model is fitted on %d units. How many propensity scores are produced?", n), strconv.Itoa(n), "One score per unit.", f)
	case 2:
		return w22P(o, fmt.Sprintf("Inverse probability weighting gives each of %d treated units a weight of one over its propensity. How many weights are computed?", n), strconv.Itoa(n), "One weight per treated unit.", f)
	default:
		return w22P(o, fmt.Sprintf("Nearest-neighbour propensity matching pairs %d treated units with controls. How many pairs are formed?", n), strconv.Itoa(n), "One pair per treated unit.", f)
	}
}

func w22DID(o int) generator.Problem {
	n := w22N(2, 20)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A difference-in-differences design uses a two-by-two table with %d units per cell. How many cells does the table have?", n), "4", "Two periods and two groups.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The difference-in-differences estimate combines four cell means for %d units per cell. How many cells enter with a positive sign?", n), "2", "Treated post and control pre are positive.", f)
	default:
		return w22P(o, fmt.Sprintf("A panel study has %d units observed in two periods. How many unit-period observations are there?", n), strconv.Itoa(2*n), "Each unit appears twice.", f)
	}
}

func w22Counterfactual(o int) generator.Problem {
	n := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "n", n)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("Each of %d units has a treated and a control potential outcome. How many potential outcomes are there?", n), strconv.Itoa(2*n), "Two per unit.", f)
	case 2:
		return w22P(o, fmt.Sprintf("For %d units only one potential outcome is observed per unit. How many outcomes are observed?", n), strconv.Itoa(n), "The other outcome is missing.", f)
	default:
		return w22P(o, fmt.Sprintf("An individual treatment effect is computed for %d units. How many effects are defined?", n), strconv.Itoa(n), "One effect per unit.", f)
	}
}

func w22Backdoor(o int) generator.Problem {
	p := w22N(1, 40)
	f := map[string]string{}
	w22I(f, "p", p)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("A treatment has %d backdoor paths. How many must a valid adjustment set block?", p), strconv.Itoa(p), "All of them.", f)
	case 2:
		return w22P(o, fmt.Sprintf("An adjustment set for a treatment with %d backdoor paths may contain no descendants. How many treatment descendants may it contain?", p), "0", "Descendants carry part of the effect.", f)
	default:
		return w22P(o, fmt.Sprintf("A minimal adjustment set is checked against %d backdoor paths. How many paths must stay blocked after removing any member?", p), strconv.Itoa(p), "Minimality preserves every block.", f)
	}
}

func w22ATE(o int) generator.Problem {
	nt := w22N(2, 30)
	nc := w22N(2, 30)
	f := map[string]string{}
	w22I(f, "nt", nt)
	w22I(f, "nc", nc)
	switch o {
	case 1:
		return w22P(o, fmt.Sprintf("An experiment has %d treated and %d control units. How many units enter the difference in means?", nt, nc), strconv.Itoa(nt+nc), "Both groups enter the estimate.", f)
	case 2:
		return w22P(o, fmt.Sprintf("The average treatment effect subtracts the control mean from the treated mean for %d treated and %d control units. How many means are estimated?", nt, nc), "2", "One mean per group.", f)
	default:
		return w22P(o, fmt.Sprintf("An effect on the treated uses only the %d treated units. How many treated units are averaged?", nt), strconv.Itoa(nt), "The ATT restricts to the treated.", f)
	}
}

func w22Binom(n, k int) int {
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

func w22CeilLog2(n int) int {
	out := 0
	v := 1
	for v < n {
		v *= 2
		out++
	}
	return out
}
