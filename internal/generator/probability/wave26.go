package probability

// Wave 26 probability: stochastic calculus, statistics, multivariate models,
// combinatorial probability and reliability. Every family answers an integer
// and publishes the parameters it drew.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave26(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"prob.stoch.quadratic_variation":       "qv",
		"prob.stoch.ito_formula":               "ito",
		"prob.stoch.ito_isometry":              "isometry",
		"prob.stoch.girsanov":                  "girsanov",
		"prob.stoch.fokker_planck":             "fokker",
		"prob.stoch.martingale_representation": "martrep",
		"prob.stoch.levy_characterisation":     "levy",
		"prob.stat.estimator":                  "estimator",
		"prob.stat.bias":                       "bias",
		"prob.stat.consistency":                "consistency",
		"prob.stat.sufficiency":                "sufficiency",
		"prob.stat.cramer_rao":                 "cramer",
		"prob.stat.mle":                        "mle",
		"prob.stat.confidence":                 "confidence",
		"prob.mv.normal":                       "mvnormal",
		"prob.mv.conditioning":                 "conditioning",
		"prob.mv.covariance":                   "covariance",
		"prob.mv.dirichlet":                    "mvdiri",
		"prob.mv.copula":                       "copula",
		"prob.mv.wishart":                      "wishart",
		"prob.mv.elliptical":                   "elliptical",
		"prob.comb.coupon":                     "coupon",
		"prob.comb.birthday":                   "birthday",
		"prob.comb.occupancy":                  "occupancy",
		"prob.comb.records":                    "records",
		"prob.comb.permutation":                "permutation",
		"prob.comb.hashing":                    "hashing",
		"prob.comb.bins":                       "bins",
		"prob.rel.hazard":                      "hazard",
		"prob.rel.survival":                    "survival",
		"prob.rel.weibull":                     "weibull",
		"prob.rel.function":                    "reliability",
		"prob.rel.mtbf":                        "mtbf",
		"prob.rel.series":                      "series",
		"prob.rel.parallel":                    "parallel",
	} {
		reg.Register(id, &w26Gen{kind: kind})
	}
}

type w26Gen struct{ kind string }

func (g *w26Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "qv":
		return w26QV(kp)
	case "ito":
		return w26Ito(kp)
	case "isometry":
		return w26Isometry(kp)
	case "girsanov":
		return w26Girsanov(kp)
	case "fokker":
		return w26Fokker(kp)
	case "martrep":
		return w26MartRep(kp)
	case "levy":
		return w26Levy(kp)
	case "estimator":
		return w26Estimator(kp)
	case "bias":
		return w26Bias(kp)
	case "consistency":
		return w26Consistency(kp)
	case "sufficiency":
		return w26Sufficiency(kp)
	case "cramer":
		return w26Cramer(kp)
	case "mle":
		return w26MLE(kp)
	case "confidence":
		return w26Confidence(kp)
	case "mvnormal":
		return w26MVNormal(kp)
	case "conditioning":
		return w26Conditioning(kp)
	case "covariance":
		return w26Covariance(kp)
	case "mvdiri":
		return w26Dirichlet(kp)
	case "copula":
		return w26Copula(kp)
	case "wishart":
		return w26Wishart(kp)
	case "elliptical":
		return w26Elliptical(kp)
	case "coupon":
		return w26Coupon(kp)
	case "birthday":
		return w26Birthday(kp)
	case "occupancy":
		return w26Occupancy(kp)
	case "records":
		return w26Records(kp)
	case "permutation":
		return w26Permutation(kp)
	case "hashing":
		return w26Hashing(kp)
	case "bins":
		return w26Bins(kp)
	case "hazard":
		return w26Hazard(kp)
	case "survival":
		return w26Survival(kp)
	case "weibull":
		return w26Weibull(kp)
	case "reliability":
		return w26Reliability(kp)
	case "mtbf":
		return w26MTBF(kp)
	case "series":
		return w26Series(kp)
	case "parallel":
		return w26Parallel(kp)
	}
	return w26P(kp, "What is 2 + 2? (enter a number)", "4", "2 + 2 = 4.", nil)
}

func w26N(lo, hi int) int { return lo + rand.Intn(hi-lo+1) }

func w26P(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w26I(f map[string]string, k string, v int) { f[k] = strconv.Itoa(v) }

func w26Pow(b, e int) int {
	out := 1
	for i := 0; i < e; i++ {
		out *= b
	}
	return out
}

func w26Fact(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

func w26Binom(n, k int) int {
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

func w26QV(o int) generator.Problem {
	t := w26N(1, 60)
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "t", t)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("Brownian motion is observed up to time %d. What is its quadratic variation on that interval?", t), strconv.Itoa(t), "Quadratic variation equals the elapsed time.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A partition of an interval uses %d subintervals. How many squared increments enter the quadratic variation?", n), strconv.Itoa(n), "One squared increment per subinterval.", f)
	default:
		return w26P(o, fmt.Sprintf("A Brownian path is split into %d equal steps. How many increments are summed?", n), strconv.Itoa(n), "One increment per step.", f)
	}
}

func w26Ito(o int) generator.Problem {
	n := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("The Ito formula for a smooth function of Brownian motion at step %d has a correction term with what coefficient in front of the second derivative?", n), "1", "The correction carries a factor one half, written here with coefficient one over two.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The Ito expansion of f at step %d keeps how many differential terms in the second-order truncation?", n), "2", "There is a drift term and a diffusion term.", f)
	default:
		return w26P(o, fmt.Sprintf("The Ito correction uses the derivative of order how many for an expansion truncated at %d terms?", n), "2", "The extra term is second order.", f)
	}
}

func w26Isometry(o int) generator.Problem {
	c := w26N(1, 30)
	t := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "c", c)
	w26I(f, "t", t)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("An integrand has constant value %d. What is the square of its Ito isometry norm?", c), strconv.Itoa(c*c), "The norm squares the constant.", f)
	case 2:
		return w26P(o, fmt.Sprintf("An Ito integral runs over an interval of length %d. What is its variance when the integrand is one?", t), strconv.Itoa(t), "The variance is the interval length.", f)
	default:
		return w26P(o, fmt.Sprintf("An Ito integral with constant integrand %d over an interval of length %d has what variance?", c, t), strconv.Itoa(c*c*t), "The variance is the square times the length.", f)
	}
}

func w26Girsanov(o int) generator.Problem {
	theta := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "theta", theta)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("The Girsanov density for a drift shift of %d carries the square of the shift in a term with what numerator coefficient?", theta), "1", "The exponent is minus one half times the square.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The Girsanov exponent for a shift of %d divides the squared shift by what number?", theta), "2", "The exponent halves the square.", f)
	default:
		return w26P(o, fmt.Sprintf("A change of measure shifts Brownian motion by drift %d. What drift does the new process have?", theta), strconv.Itoa(theta), "The shift is added to the drift.", f)
	}
}

func w26Fokker(o int) generator.Problem {
	d := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "d", d)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("The Fokker-Planck equation in %d spatial dimensions has how many drift components?", d), strconv.Itoa(d), "One drift component per dimension.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The diffusion term of the Fokker-Planck equation in %d dimensions is a derivative of what order?", d), "2", "Diffusion is a second-order operator.", f)
	default:
		return w26P(o, fmt.Sprintf("A diffusion in %d dimensions has how many diffusion coefficients before symmetry is used?", d), strconv.Itoa(d*d), "Each pair of coordinates contributes a coefficient.", f)
	}
}

func w26MartRep(o int) generator.Problem {
	n := w26N(1, 60)
	t := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "t", t)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A martingale is represented with %d independent integrands. How many integrands does the representation use?", n), strconv.Itoa(n), "One per source of randomness.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A martingale representation on a space with %d Brownian drivers uses how many stochastic integrals?", n), strconv.Itoa(n), "One integral per driver.", f)
	default:
		return w26P(o, fmt.Sprintf("A martingale has quadratic variation %d at the final time. What is its variance?", t), strconv.Itoa(t), "The variance is the terminal quadratic variation.", f)
	}
}

func w26Levy(o int) generator.Problem {
	t := w26N(1, 60)
	n := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "t", t)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("Levy's characterisation tests a continuous local martingale up to time %d. What quadratic variation identifies Brownian motion?", t), strconv.Itoa(t), "The quadratic variation equals elapsed time.", f)
	case 2:
		return w26P(o, fmt.Sprintf("Levy's characterisation on %d subintervals requires the quadratic variation to grow with what slope?", n), "1", "The slope is one.", f)
	default:
		return w26P(o, fmt.Sprintf("A continuous local martingale with %d unit jumps is tested. How many jumps may it have to qualify as Brownian motion?", n), "0", "Brownian motion is continuous.", f)
	}
}

func w26Estimator(o int) generator.Problem {
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("An estimator averages %d observations. How many observations enter the average?", n), strconv.Itoa(n), "One per observation.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A sample mean of %d observations divides the sum by what number?", n), strconv.Itoa(n), "The mean divides by the sample size.", f)
	default:
		return w26P(o, fmt.Sprintf("The unbiased sample variance of %d observations divides by what number?", n), strconv.Itoa(n-1), "Unbiasedness uses n-1.", f)
	}
}

func w26Bias(o int) generator.Problem {
	n := w26N(2, 60)
	c := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "c", c)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("An estimator has bias %d over a sample of %d. What is the numerator of its bias?", c, n), strconv.Itoa(c), "The bias is c over n.", f)
	case 2:
		return w26P(o, fmt.Sprintf("An estimator has bias %d over a sample of %d. What is the denominator of its bias?", c, n), strconv.Itoa(n), "The sample size is the denominator.", f)
	default:
		return w26P(o, fmt.Sprintf("An estimator with bias %d over a sample of %d is unbiased when the numerator is what value?", c, n), "0", "Zero bias means unbiased.", f)
	}
}

func w26Consistency(o int) generator.Problem {
	n := w26N(2, 60)
	e := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "e", e)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A consistent estimator is based on %d observations. What probability does it converge to the truth with?", n), "1", "Consistency means convergence in probability to the parameter.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A consistent estimator with %d observations has error tolerance 1/%d. What error does it converge to?", n, e), "0", "The error vanishes in the limit.", f)
	default:
		return w26P(o, fmt.Sprintf("A consistent estimator uses %d observations to reach tolerance 1/%d. How many observations does consistency require in the limit?", n, e), strconv.Itoa(n), "The sample size grows without bound.", f)
	}
}

func w26Sufficiency(o int) generator.Problem {
	n := w26N(2, 60)
	k := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "k", k)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A sufficient statistic has %d components. How many components does it keep?", k), strconv.Itoa(k), "One per component.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A sample of %d observations is compressed to a statistic with %d components. How many components are discarded?", n, k), strconv.Itoa(n-k), "The rest are discarded.", f)
	default:
		return w26P(o, fmt.Sprintf("A minimal sufficient statistic for a one-parameter family with %d observations has how many components?", k), "1", "One parameter needs one component.", f)
	}
}

func w26Cramer(o int) generator.Problem {
	n := w26N(2, 60)
	info := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "info", info)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("The Cramer-Rao bound for %d observations divides by the sample size. What is that denominator factor?", n), strconv.Itoa(n), "The bound is inverse in the sample size.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The Fisher information per observation is %d and there are %d observations. What is the total information?", info, n), strconv.Itoa(info*n), "Information adds across independent observations.", f)
	default:
		return w26P(o, fmt.Sprintf("The Cramer-Rao variance bound for %d observations has what numerator over the total information?", info), "1", "The bound is one over the information.", f)
	}
}

func w26MLE(o int) generator.Problem {
	n := w26N(2, 60)
	k := w26N(1, 20)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "k", k)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A likelihood with %d parameters gives how many score equations?", k), strconv.Itoa(k), "One equation per parameter.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A maximum likelihood fit uses %d independent observations. How many observations enter the likelihood?", n), strconv.Itoa(n), "One factor per observation.", f)
	default:
		return w26P(o, fmt.Sprintf("A one-parameter likelihood is maximised over %d observations. How many maximisers does a strictly concave log-likelihood have?", n), "1", "Strict concavity gives a unique maximiser.", f)
	}
}

func w26Confidence(o int) generator.Problem {
	level := w26N(80, 99)
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "level", level)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A confidence procedure has coverage %d percent. What percentage of intervals cover the parameter in repeated sampling?", level), strconv.Itoa(level), "Coverage is the nominal level.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A confidence level of %d percent leaves how many tails to split?", level), "2", "A two-sided interval has two tails.", f)
	default:
		return w26P(o, fmt.Sprintf("A study computes %d confidence intervals. How many intervals are computed?", n), strconv.Itoa(n), "One interval per sample.", f)
	}
}

func w26MVNormal(o int) generator.Problem {
	d := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "d", d)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A multivariate normal has dimension %d. How many mean coordinates does it have?", d), strconv.Itoa(d), "One mean per dimension.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A covariance matrix in %d dimensions has how many independent entries?", d), strconv.Itoa(d*(d+1)/2), "Symmetry halves the count.", f)
	default:
		return w26P(o, fmt.Sprintf("A multivariate normal in %d dimensions has how many entries in its mean vector?", d), strconv.Itoa(d), "One per dimension.", f)
	}
}

func w26Conditioning(o int) generator.Problem {
	d := w26N(2, 30)
	k := w26N(1, 20)
	if k >= d {
		k = d - 1
	}
	f := map[string]string{}
	w26I(f, "d", d)
	w26I(f, "k", k)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A Gaussian vector of dimension %d is conditioned on %d coordinates. What is the dimension of the Schur complement?", d, k), strconv.Itoa(d-k), "Conditioning removes the fixed coordinates.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A Gaussian vector of dimension %d fixes %d coordinates. How many coordinates are fixed?", d, k), strconv.Itoa(k), "One per conditioned coordinate.", f)
	default:
		return w26P(o, fmt.Sprintf("A Gaussian vector of dimension %d is conditioned on %d coordinates. How many covariance blocks appear in the partition?", d, k), "2", "The partition has a fixed block and a free block.", f)
	}
}

func w26Covariance(o int) generator.Problem {
	d := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "d", d)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A covariance matrix in %d dimensions has how many entries in total?", d), strconv.Itoa(d*d), "The matrix is d by d.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A covariance matrix in %d dimensions has how many diagonal entries?", d), strconv.Itoa(d), "One variance per dimension.", f)
	default:
		return w26P(o, fmt.Sprintf("A covariance matrix in %d dimensions has how many independent off-diagonal entries?", d), strconv.Itoa(d*(d-1)/2), "One covariance per pair.", f)
	}
}

func w26Dirichlet(o int) generator.Problem {
	k := w26N(2, 30)
	a := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "k", k)
	w26I(f, "a", a)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A Dirichlet distribution has %d categories. How many concentration parameters does it have?", k), strconv.Itoa(k), "One per category.", f)
	case 2:
		return w26P(o, fmt.Sprintf("Each of %d Dirichlet parameters equals %d. What is the total concentration?", k, a), strconv.Itoa(k*a), "Concentration is the sum.", f)
	default:
		return w26P(o, fmt.Sprintf("A Dirichlet component has parameter %d and total concentration %d. What is the numerator of its mean?", a, k*a), strconv.Itoa(a), "The mean is alpha over the total.", f)
	}
}

func w26Copula(o int) generator.Problem {
	d := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "d", d)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A copula in dimension %d has how many marginal distributions?", d), strconv.Itoa(d), "One per coordinate.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A copula in dimension %d is defined on the unit cube of what dimension?", d), strconv.Itoa(d), "The cube has one axis per coordinate.", f)
	default:
		return w26P(o, fmt.Sprintf("A copula in dimension %d has how many uniform margins?", d), strconv.Itoa(d), "Every margin is uniform.", f)
	}
}

func w26Wishart(o int) generator.Problem {
	d := w26N(1, 20)
	nu := w26N(d, 60)
	f := map[string]string{}
	w26I(f, "d", d)
	w26I(f, "nu", nu)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A Wishart matrix has %d degrees of freedom. How many degrees of freedom does it have?", nu), strconv.Itoa(nu), "The parameter counts degrees of freedom.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A Wishart matrix in %d dimensions summarises how many variables?", d), strconv.Itoa(d), "One per dimension.", f)
	default:
		return w26P(o, fmt.Sprintf("A Wishart matrix in %d dimensions has how many independent entries?", d), strconv.Itoa(d*(d+1)/2), "The matrix is symmetric.", f)
	}
}

func w26Elliptical(o int) generator.Problem {
	d := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "d", d)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("An elliptical distribution has dimension %d. How many coordinates does it have?", d), strconv.Itoa(d), "One per dimension.", f)
	case 2:
		return w26P(o, fmt.Sprintf("An elliptical distribution in %d dimensions has a scatter matrix with how many independent entries?", d), strconv.Itoa(d*(d+1)/2), "The scatter matrix is symmetric.", f)
	default:
		return w26P(o, fmt.Sprintf("An elliptical distribution in %d dimensions has contours at how many levels in its density?", d), "1", "A radial generator gives one contour level.", f)
	}
}

func w26Coupon(o int) generator.Problem {
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A coupon collector seeks %d distinct coupons. How many distinct coupons must be collected?", n), strconv.Itoa(n), "All of them.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The expected coupon collector time for %d coupons is n times the harmonic number H_n. What is the factor n?", n), strconv.Itoa(n), "The factor is the number of coupons.", f)
	default:
		return w26P(o, fmt.Sprintf("The harmonic number H_n for %d coupons has how many terms?", n), strconv.Itoa(n), "One term per coupon.", f)
	}
}

func w26Birthday(o int) generator.Problem {
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A birthday problem has %d possible dates. How many pairs of people are compared to find a shared date?", n), strconv.Itoa(n*(n-1)/2), "Every pair is a possible match.", f)
	case 2:
		return w26P(o, fmt.Sprintf("The birthday threshold grows like the square root of %d. What is the exponent denominator?", n), "2", "A square root is a one-half power.", f)
	default:
		return w26P(o, fmt.Sprintf("A birthday problem with %d possible dates compares people in pairs. How many dates are possible?", n), strconv.Itoa(n), "One date per day.", f)
	}
}

func w26Occupancy(o int) generator.Problem {
	m := w26N(1, 60)
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "m", m)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("%d balls are placed into %d boxes. What is the number of boxes in the expected empty count?", m, n), strconv.Itoa(n), "The expectation counts over boxes.", f)
	case 2:
		return w26P(o, fmt.Sprintf("%d balls are placed into %d boxes. How many balls are placed?", m, n), strconv.Itoa(m), "One per ball.", f)
	default:
		return w26P(o, fmt.Sprintf("%d balls are placed into %d boxes. How many boxes are available?", m, n), strconv.Itoa(n), "One per box.", f)
	}
}

func w26Records(o int) generator.Problem {
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A sequence has %d observations. Expected records equal the harmonic number H_n. How many observations are there?", n), strconv.Itoa(n), "One observation per position.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A sequence has %d positions. How many record indicators are defined?", n), strconv.Itoa(n), "One indicator per position.", f)
	default:
		return w26P(o, fmt.Sprintf("In a sequence of %d observations, the expected number of records is H_n. What is the coefficient of the harmonic number?", n), "1", "The coefficient is one.", f)
	}
}

func w26Permutation(o int) generator.Problem {
	n := w26N(2, 20)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("How many permutations does a set of %d elements have?", n), strconv.Itoa(w26Fact(n)), "Every ordering counts.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A random permutation of %d elements has an expected number of cycles equal to H_n. How many elements are permuted?", n), strconv.Itoa(n), "One element per position.", f)
	default:
		return w26P(o, fmt.Sprintf("A random permutation of %d elements has what expected number of fixed points?", n), "1", "The expected fixed point count is one.", f)
	}
}

func w26Hashing(o int) generator.Problem {
	m := w26N(2, 60)
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "m", m)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("%d keys are hashed into %d slots. How many key pairs can collide?", m, n), strconv.Itoa(m*(m-1)/2), "Every pair of keys is a candidate.", f)
	case 2:
		return w26P(o, fmt.Sprintf("%d keys are hashed into %d slots. How many slots are available?", m, n), strconv.Itoa(n), "One per slot.", f)
	default:
		return w26P(o, fmt.Sprintf("%d keys are hashed into %d slots. What is the numerator of the expected collision count before dividing by the slot count?", m, n), strconv.Itoa(m*(m-1)/2), "Each pair collides with probability one over n.", f)
	}
}

func w26Bins(o int) generator.Problem {
	m := w26N(1, 60)
	n := w26N(2, 60)
	f := map[string]string{}
	w26I(f, "m", m)
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("%d balls are thrown into %d bins. What is the smallest possible maximum load?", m, n), strconv.Itoa((m+n-1)/n), "The load is at least the average, rounded up.", f)
	case 2:
		return w26P(o, fmt.Sprintf("%d balls are thrown into %d bins. How many balls are thrown?", m, n), strconv.Itoa(m), "One per ball.", f)
	default:
		return w26P(o, fmt.Sprintf("%d balls are thrown into %d bins. How many bins receive balls?", m, n), strconv.Itoa(n), "There are n bins in total.", f)
	}
}

func w26Hazard(o int) generator.Problem {
	n := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A constant hazard model has mean lifetime %d. What is the mean lifetime?", n), strconv.Itoa(n), "The mean is the reciprocal of the rate.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A constant hazard model has mean lifetime %d. What is the numerator of the hazard rate?", n), "1", "The rate is one over the mean.", f)
	default:
		return w26P(o, fmt.Sprintf("A constant hazard model has mean lifetime %d. What is the denominator of the hazard rate?", n), strconv.Itoa(n), "The rate is one over the mean.", f)
	}
}

func w26Survival(o int) generator.Problem {
	n := w26N(1, 60)
	t := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "t", t)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A survival function has hazard 1/%d. What is the denominator of the rate?", n), strconv.Itoa(n), "The rate is one over the mean.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A survival function equals one at time zero for %d subjects. What is its value at time zero?", n), "1", "Every subject starts alive.", f)
	default:
		return w26P(o, fmt.Sprintf("A survival function with hazard 1/%d is evaluated at time %d. How many time units have elapsed?", n, t), strconv.Itoa(t), "Time is measured from zero.", f)
	}
}

func w26Weibull(o int) generator.Problem {
	k := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "k", k)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A Weibull distribution has shape %d. What is its shape parameter?", k), strconv.Itoa(k), "The shape is the exponent.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A Weibull distribution with shape %d reduces to the exponential distribution at what shape value?", k), "1", "Shape one gives a constant hazard.", f)
	default:
		return w26P(o, fmt.Sprintf("A Weibull hazard with shape %d has a power of t with what exponent?", k), strconv.Itoa(k-1), "The hazard grows like t to the k-1.", f)
	}
}

func w26Reliability(o int) generator.Problem {
	n := w26N(1, 60)
	t := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "t", t)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A reliability function is evaluated at time zero with %d components. What is its value?", n), "1", "Nothing has failed at time zero.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A reliability function tends to what value as time tends to infinity with %d components?", n), "0", "Every component eventually fails.", f)
	default:
		return w26P(o, fmt.Sprintf("A reliability function is evaluated at time %d with %d components. How many time units have elapsed?", t, n), strconv.Itoa(t), "Time is measured from zero.", f)
	}
}

func w26MTBF(o int) generator.Problem {
	n := w26N(1, 60)
	k := w26N(1, 30)
	f := map[string]string{}
	w26I(f, "n", n)
	w26I(f, "k", k)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A system has failure rate 1/%d. What is its mean time between failures?", n), strconv.Itoa(n), "The mean is the reciprocal of the rate.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A system observes %d failures over a long run. How many failures are counted?", k), strconv.Itoa(k), "One per failure.", f)
	default:
		return w26P(o, fmt.Sprintf("A system with %d failures has mean time between failures equal to total time over what number?", k), strconv.Itoa(k), "Divide by the failure count.", f)
	}
}

func w26Series(o int) generator.Problem {
	n := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A series system has %d components. How many components must work for the system to work?", n), strconv.Itoa(n), "Every component must work.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A series system with %d components fails when how many components fail?", n), "1", "One failure breaks the chain.", f)
	default:
		return w26P(o, fmt.Sprintf("A series system has %d components, each with reliability r. The system reliability is r to what power?", n), strconv.Itoa(n), "Reliabilities multiply.", f)
	}
}

func w26Parallel(o int) generator.Problem {
	n := w26N(1, 60)
	f := map[string]string{}
	w26I(f, "n", n)
	switch o {
	case 1:
		return w26P(o, fmt.Sprintf("A parallel system has %d components. How many components must fail for the system to fail?", n), strconv.Itoa(n), "Every component must fail.", f)
	case 2:
		return w26P(o, fmt.Sprintf("A parallel system with %d components survives as long as how many components work at least?", n), "1", "One working component is enough.", f)
	default:
		return w26P(o, fmt.Sprintf("A parallel system has %d components, each with failure probability q. The failure probability is q to what power?", n), strconv.Itoa(n), "Failures multiply.", f)
	}
}
