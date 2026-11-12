package numerical

// Wave 16, part three: stochastic numerics. The families follow the same
// integer-valued contract as the rest of the wave.

import "fmt"

var wave16FamiliesC = map[string][3]func() w16Trio{
	"numa.stochastic.standard_error": {
		func() w16Trio {
			sigma := w16Limit(1, 10)
			k := w16Limit(2, 9)
			n := k * k
			return w16Mk(fmt.Sprintf("A Monte Carlo average uses %d independent samples with standard deviation %d. What is its standard error?", n, sigma), sigma/k,
				fmt.Sprintf("The error is %d/sqrt(%d) = %d/%d = %d.", sigma, n, sigma, k, sigma/k))
		},
		func() w16Trio {
			n := w16Limit(1, 20) * 100
			return w16Mk(fmt.Sprintf("A Monte Carlo estimate uses %d samples. How many samples are needed to halve the standard error?", n), 4*n,
				fmt.Sprintf("The error falls with the square root of the count, so halving it needs four times the samples: %d.", 4*n))
		},
		func() w16Trio {
			f := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("A simulation must reduce its standard error by a factor of %d. By what factor must the sample count grow?", f), f*f,
				fmt.Sprintf("The count must grow by the square of the error factor, %d squared = %d.", f, f*f))
		},
	},
	"numa.stochastic.antithetic": {
		func() w16Trio {
			pairs := w16Limit(10, 200)
			return w16Mk(fmt.Sprintf("An antithetic estimator uses %d pairs, each costing two evaluations. How many evaluations does it perform?", pairs), 2*pairs,
				fmt.Sprintf("Two evaluations per pair times %d pairs gives %d.", pairs, 2*pairs))
		},
		func() w16Trio {
			rho := -10 * w16Limit(1, 9)
			return w16Mk(fmt.Sprintf("Two antithetic responses have correlation %d percent. The pair average has variance sigma squared (1 + rho)/2. What percentage of sigma squared is that?", rho), (100+rho)/2,
				fmt.Sprintf("The factor is (100 + (%d))/2 percent = %d percent.", rho, (100+rho)/2))
		},
		func() w16Trio {
			pairs := w16Limit(10, 300)
			return w16Mk(fmt.Sprintf("An antithetic estimator uses %d pairs, each built from one uniform draw and its complement. How many distinct uniform draws does it generate?", pairs), pairs,
				fmt.Sprintf("Each pair needs one fresh uniform, because the partner is its complement, so %d draws.", pairs))
		},
	},
	"numa.stochastic.control_variates": {
		func() w16Trio {
			rho := 10 * w16Limit(1, 9)
			return w16Mk(fmt.Sprintf("A control variate has correlation %d percent with the estimator. The variance falls to (10000 - rho squared)/10000 of its value. What percentage of the original variance remains?", rho),
				100-rho*rho/100,
				fmt.Sprintf("The remaining fraction is 100 - %d squared/100 = %d percent.", rho, 100-rho*rho/100))
		},
		func() w16Trio {
			cov := w16Limit(1, 20)
			vr := w16Limit(1, 5) * 10
			return w16Mk(fmt.Sprintf("A control variate has covariance %d with the estimator and variance %d. What is the optimal coefficient beta = cov/var, in percent?", cov, vr), 100*cov/vr,
				fmt.Sprintf("The coefficient is %d/%d = %d percent.", cov, vr, 100*cov/vr))
		},
		func() w16Trio {
			x := w16Limit(10, 50)
			beta := w16Limit(1, 4)
			y := w16Limit(1, 9)
			mu := w16Limit(1, 5)
			return w16Mk(fmt.Sprintf("A control-variate estimator has X = %d, beta = %d, Y = %d and known mean mu = %d. What is X - beta(Y - mu)?", x, beta, y, mu), x-beta*(y-mu),
				fmt.Sprintf("The correction is %d · (%d - %d) = %d, so the estimate is %d - %d = %d.", beta, y, mu, beta*(y-mu), x, beta*(y-mu), x-beta*(y-mu)))
		},
	},
	"numa.stochastic.importance": {
		func() w16Trio {
			q := w16Limit(2, 10)
			k := w16Limit(2, 10)
			p := q * k
			return w16Mk(fmt.Sprintf("A sample has target density %d and proposal density %d. What is its importance weight p/q?", p, q), k,
				fmt.Sprintf("The weight is %d/%d = %d.", p, q, k))
		},
		func() w16Trio {
			q := []int{1, 2, 4, 5, 10, 20, 25}[w16Limit(0, 6)]
			return w16Mk(fmt.Sprintf("A sample is drawn from a proposal with density %d while the target has density 1. What is the importance weight as a percentage of the target density?", q), 100/q,
				fmt.Sprintf("The weight is 1/%d, which is %d percent.", q, 100/q))
		},
		func() w16Trio {
			n := w16Limit(50, 500)
			return w16Mk(fmt.Sprintf("A importance-sampling run gives every one of %d samples the same weight. What is the effective sample size?", n), n,
				fmt.Sprintf("Equal weights leave the effective size equal to the sample count, %d.", n))
		},
	},
	"numa.stochastic.stratified": {
		func() w16Trio {
			strata := w16Limit(2, 10)
			per := w16Limit(5, 40)
			return w16Mk(fmt.Sprintf("A stratified estimator draws %d samples from each of %d strata. How many samples does it draw in total?", per, strata), per*strata,
				fmt.Sprintf("%d per stratum times %d strata gives %d.", per, strata, per*strata))
		},
		func() w16Trio {
			between := 10 * w16Limit(1, 9)
			return w16Mk(fmt.Sprintf("Between-stratum variation accounts for %d percent of the total variance and stratification removes it. What percentage of the variance remains?", between), 100-between,
				fmt.Sprintf("Removing %d percent leaves %d percent.", between, 100-between))
		},
		func() w16Trio {
			total := w16Limit(2, 10) * 100
			strata := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("A stratified run uses proportional allocation with %d samples over %d strata. How many samples does each stratum receive?", total, strata), total/strata,
				fmt.Sprintf("Proportional allocation divides evenly here: %d/%d = %d.", total, strata, total/strata))
		},
	},
	"numa.stochastic.quasi": {
		func() w16Trio {
			e := 2 * w16Limit(1, 6)
			return w16Mk(fmt.Sprintf("A quasi-Monte Carlo rule uses 10 to the power %d points, where its error is about 1/N and the random error is about 1/sqrt(N). How many decimal digits better than random sampling is it?", e), e/2,
				fmt.Sprintf("The two exponents are %d and %d, so the gain is %d digits.", e, e/2, e/2))
		},
		func() w16Trio {
			d := w16Limit(2, 8)
			return w16Mk(fmt.Sprintf("The quasi-Monte Carlo error bound carries a factor (log N) to the power d. Enter the exponent d for a %d-dimensional integral.", d), d,
				fmt.Sprintf("The dimension is the exponent of the logarithm, so d = %d.", d))
		},
		func() w16Trio {
			n := w16Limit(1, 10) * 1000
			return w16Mk(fmt.Sprintf("A low-discrepancy sequence supplies %d points. How many points does the rule evaluate?", n), n,
				fmt.Sprintf("The rule evaluates every point once, so %d evaluations.", n))
		},
	},
	"numa.stochastic.metropolis": {
		func() w16Trio {
			ratio := 10 * w16Limit(1, 10)
			accept := ratio
			if accept > 100 {
				accept = 100
			}
			return w16Mk(fmt.Sprintf("A Metropolis-Hastings proposal has acceptance ratio %d percent. What is the acceptance probability in percent?", ratio), accept,
				fmt.Sprintf("The acceptance probability is the minimum of the ratio and one, so %d percent.", accept))
		},
		func() w16Trio {
			piX := w16Limit(2, 10)
			piY := w16Limit(1, piX)
			return w16Mk(fmt.Sprintf("A symmetric proposal moves from a point with target density %d to one with density %d. What is the acceptance probability in percent?", piX, piY), 100*piY/piX,
				fmt.Sprintf("The ratio is %d/%d, so the probability is %d percent.", piY, piX, 100*piY/piX))
		},
		func() w16Trio {
			n := w16Limit(2, 10) * 100
			rate := w16Limit(1, 9) * 10
			return w16Mk(fmt.Sprintf("A chain makes %d proposals and accepts %d percent of them. How many proposals are accepted?", n, rate), n*rate/100,
				fmt.Sprintf("%d percent of %d is %d.", rate, n, n*rate/100))
		},
	},
	"numa.stochastic.gibbs": {
		func() w16Trio {
			d := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A Gibbs sampler updates every coordinate once per sweep. How many updates does one sweep of a %d-dimensional target perform?", d), d,
				fmt.Sprintf("One sweep updates each of the %d coordinates once.", d))
		},
		func() w16Trio {
			d := w16Limit(2, 12)
			k := w16Limit(2, 10)
			return w16Mk(fmt.Sprintf("A Gibbs sampler runs %d sweeps on a %d-dimensional target. How many coordinate updates does it perform?", k, d), k*d,
				fmt.Sprintf("%d updates per sweep times %d sweeps gives %d.", d, k, k*d))
		},
		func() w16Trio {
			n := w16Limit(10, 100)
			return w16Mk(fmt.Sprintf("A Gibbs sampler draws each coordinate from its exact full conditional. Out of %d proposals, how many are accepted?", n), n,
				fmt.Sprintf("Every full-conditional draw is accepted, so all %d proposals are accepted.", n))
		},
	},
	"numa.stochastic.bootstrap": {
		func() w16Trio {
			n := w16Limit(5, 200)
			return w16Mk(fmt.Sprintf("A bootstrap resamples a data set of %d observations. How many observations does each replicate contain?", n), n,
				fmt.Sprintf("Each replicate has the same size as the original sample, %d.", n))
		},
		func() w16Trio {
			b := w16Limit(2, 20) * 100
			n := w16Limit(5, 50)
			return w16Mk(fmt.Sprintf("A bootstrap draws %d replicates, each of size %d. How many observations does it draw in total?", b, n), b*n,
				fmt.Sprintf("%d replicates times %d observations gives %d.", b, n, b*n))
		},
		func() w16Trio {
			b := w16Limit(2, 20) * 100
			return w16Mk(fmt.Sprintf("A bootstrap uses %d replicates. How many statistic values does it compute?", b), b,
				fmt.Sprintf("One statistic per replicate, so %d values.", b))
		},
	},
	"numa.stochastic.rng": {
		func() w16Trio {
			a := w16Limit(2, 5)
			x := w16Limit(1, 9)
			c := w16Limit(1, 9)
			m := w16Limit(10, 30)
			return w16Mk(fmt.Sprintf("A linear congruential generator has a = %d, c = %d and m = %d. Starting from x = %d, what is the next value (a x + c) mod m?", a, c, m, x), (a*x+c)%m,
				fmt.Sprintf("%d · %d + %d = %d, and %d mod %d = %d.", a, x, c, a*x+c, a*x+c, m, (a*x+c)%m))
		},
		func() w16Trio {
			m := w16Limit(10, 100)
			return w16Mk(fmt.Sprintf("A linear congruential generator works modulo %d. What is the largest period it can have?", m), m,
				fmt.Sprintf("The state space has %d values, so no period can exceed %d.", m, m))
		},
		func() w16Trio {
			k := w16Limit(2, 12)
			return w16Mk(fmt.Sprintf("A parallel simulation runs %d workers and each needs an independent stream. How many seeds must be chosen?", k), k,
				fmt.Sprintf("One independent seed per worker, so %d seeds.", k))
		},
	},
}
