package probability

import "github.com/chuma-beep/mathua/internal/generator"

// The limit-theorem concepts are statements rather than numbers. Asking how
// many of eight drawn statements are true keeps the answer typeable while the
// explanation still states why each one holds or fails.

func rvIndependenceProblem() generator.Problem {
	pool := []probVerdict{
		{"X and Y are independent exactly when P(X in A, Y in B) = P(X in A) P(Y in B) for all Borel A and B", true, "this is the definition"},
		{"Independent random variables with a joint density have a density that factorises", true, "f_{X,Y}(x,y) = f_X(x) f_Y(y)"},
		{"Independent random variables are uncorrelated", true, "independence makes E[XY] = E[X]E[Y]"},
		{"Zero covariance implies independence", false, "uncorrelated is weaker than independent"},
		{"Functions of independent random variables are independent", true, "measurable maps preserve independence"},
		{"A sum of independent normal variables is normal", true, "means and variances add"},
		{"The indicators of independent events are independent random variables", true, "their joint law is the product law"},
		{"Independent Bernoulli variables have P(X = 1, Y = 1) = p q", true, "the joint probability factorises"},
		{"Pairwise independence of three variables implies mutual independence", false, "the triple product need not hold"},
		{"Independence is preserved by taking limits in distribution", false, "the limiting joint law need not factorise"},
		{"P(X = Y) = 0 for every pair of independent continuous variables", false, "it fails when both have atoms at the same value"},
		{"Independent variables must have zero means", false, "independence says nothing about location"},
		{"Independent variables must be identically distributed", false, "only the joint law factorises"},
		{"A constant is independent of every random variable", true, "its sigma-algebra is trivial"},
		{"E[XY] = E[X]E[Y] always implies independence", false, "this is only lack of correlation"},
		{"The joint CDF of independent variables is the product of the marginal CDFs", true, "F(x,y) = F_X(x) F_Y(y)"},
	}
	return countProbVerdicts("How many of these statements about independence are true?", pool)
}

func convergenceInProbabilityProblem() generator.Problem {
	pool := []probVerdict{
		{"X_n converges to X in probability when P(|X_n - X| > eps) tends to 0 for every eps > 0", true, "this is the definition"},
		{"Convergence in probability implies convergence in distribution", true, "the portmanteau theorem gives it"},
		{"Convergence in distribution to a constant implies convergence in probability to that constant", true, "the limit is deterministic"},
		{"Almost-sure convergence implies convergence in probability", true, "pointwise convergence is stronger"},
		{"Convergence in probability implies almost-sure convergence", false, "only a subsequence converges almost surely"},
		{"Convergence in mean square implies convergence in probability", true, "Chebyshev's inequality"},
		{"Convergence in probability implies convergence in mean", false, "the variables need not be integrable"},
		{"The limit in probability is unique up to a null set", true, "two limits must agree almost surely"},
		{"The continuous mapping theorem holds for convergence in probability", true, "continuous g preserves it"},
		{"Convergence in distribution alone implies convergence in probability", false, "a constant limit is needed"},
		{"Every subsequence of a sequence converging in probability also converges in probability", true, "to the same limit"},
		{"Boundedness in probability is enough for convergence", false, "bounded sequences need not converge"},
		{"Convergence in probability is preserved by addition and multiplication", true, "the algebra of limits holds"},
		{"A sequence can converge in probability to two different constants", false, "the constants would have to be equal"},
		{"If X_n converges in probability to X then E[X_n] converges to E[X]", false, "integrability is required"},
		{"Convergence in probability is metrised by d(X,Y) = E[min(1, |X-Y|)]", true, "this metric induces the topology"},
	}
	return countProbVerdicts("How many of these statements about convergence in probability are true?", pool)
}

func convergenceInDistributionProblem() generator.Problem {
	pool := []probVerdict{
		{"X_n converges in distribution to X when F_n(x) tends to F(x) at every continuity point of F", true, "this is the definition"},
		{"The portmanteau theorem characterises convergence in distribution by weak convergence of laws", true, "integrals of bounded continuous functions converge"},
		{"Convergence in distribution to a constant implies convergence in probability", true, "the limit law is a point mass"},
		{"Convergence in distribution is preserved by continuous maps", true, "the continuous mapping theorem"},
		{"Slutsky's theorem combines distribution and probability limits", true, "sums, products and ratios can be handled"},
		{"Convergence in distribution always implies convergence in probability", false, "that needs a constant limit"},
		{"Convergence in distribution implies convergence of densities", false, "densities need not exist or converge"},
		{"Convergence of characteristic functions pointwise implies convergence in distribution", true, "Levy's continuity theorem"},
		{"The central limit theorem is a statement about convergence in distribution", true, "the standardised sums approach a normal law"},
		{"Convergence in distribution implies convergence of means", false, "the variables need not be uniformly integrable"},
		{"A sequence of point masses at 1/n converges in distribution to a point mass at 0", true, "the CDFs converge at continuity points"},
		{"Convergence in distribution is unaffected by changing finitely many terms", true, "the limit depends on the tail"},
		{"Every tight sequence has a distributionally convergent subsequence", true, "Prokhorov's theorem"},
		{"Convergence in distribution requires all moments to converge", false, "only the laws converge"},
		{"Convergence in distribution is equivalent to convergence of all bounded continuous expectations", true, "the portmanteau theorem"},
		{"A sequence converging in distribution must converge almost surely", false, "Skorokhod gives a coupling, not the original sequence"},
	}
	return countProbVerdicts("How many of these statements about convergence in distribution are true?", pool)
}

func almostSureProblem() generator.Problem {
	pool := []probVerdict{
		{"X_n converges to X almost surely when P(X_n tends to X) = 1", true, "this is the definition"},
		{"Almost-sure convergence implies convergence in probability", true, "convergence on a full-measure set"},
		{"Convergence in probability implies almost-sure convergence along a subsequence", true, "a standard subsequence argument"},
		{"If the sum of P(|X_n - X| > eps) is finite then X_n converges almost surely", true, "the first Borel-Cantelli lemma"},
		{"Almost-sure limits are unique up to null sets", true, "two limits agree almost surely"},
		{"Almost-sure convergence is preserved by continuous maps", true, "evaluate the map along the convergent sequence"},
		{"Convergence in probability implies almost-sure convergence", false, "the full sequence need not converge"},
		{"Almost-sure convergence implies convergence in mean", false, "integrability and domination are needed"},
		{"Almost-sure convergence implies convergence of expectations", false, "the dominated convergence theorem needs a dominating variable"},
		{"Almost-sure convergence is equivalent to convergence in distribution", false, "distributional convergence is much weaker"},
		{"If X_n converges almost surely then every subsequence does", true, "restrict the full-measure set"},
		{"Almost-sure convergence requires independence of the variables", false, "no independence is needed"},
		{"A sequence that converges almost surely must be bounded", false, "the limit may be unbounded"},
		{"The set where a sequence fails to converge can have positive probability", true, "almost sure allows a null exceptional set only"},
		{"Almost-sure convergence is metrised by convergence in probability", false, "almost sure is strictly stronger"},
		{"If X_n converges almost surely to X then X_n - X converges almost surely to 0", true, "subtract the same limit"},
	}
	return countProbVerdicts("How many of these statements about almost-sure convergence are true?", pool)
}

func weakLawProblem() generator.Problem {
	pool := []probVerdict{
		{"The weak law says the sample mean converges in probability to the population mean", true, "this is the conclusion"},
		{"Finite variance and Chebyshev's inequality prove the weak law for iid samples", true, "Var(mean) = sigma^2/n"},
		{"The weak law implies almost-sure convergence of the sample mean", false, "that is the strong law"},
		{"The weak law requires the variables to be normally distributed", false, "only iid and a finite mean are needed"},
		{"Khinchin's weak law needs only a finite first moment for iid variables", true, "finite variance is not required"},
		{"The weak law is a statement about convergence in probability", true, "the probability of a large deviation vanishes"},
		{"The sample mean converges to the mean at rate 1/sqrt(n) in standard deviation", true, "Var(mean) = sigma^2/n"},
		{"The weak law guarantees the sample mean equals the mean eventually", false, "it only makes large deviations unlikely"},
		{"The weak law applies to Bernoulli samples with success probability p", true, "the sample proportion converges to p"},
		{"The weak law requires the observations to be independent but not identically distributed in every version", true, "there are versions with weaker hypotheses"},
		{"A sequence with infinite mean still satisfies the weak law", false, "the mean must exist"},
		{"The weak law is equivalent to the central limit theorem", false, "the CLT gives a rate and a limit law"},
		{"The weak law can be proved by truncating the variables", true, "truncation is a standard proof technique"},
		{"The weak law says the sample mean converges in distribution to a normal law", false, "that is the CLT"},
		{"Chebyshev's inequality gives a quantitative bound in the weak law", true, "P(|mean - mu| > eps) <= sigma^2/(n eps^2)"},
		{"The weak law fails for bounded random variables", false, "bounded variables have finite variance"},
	}
	return countProbVerdicts("How many of these statements about the weak law of large numbers are true?", pool)
}

func strongLawProblem() generator.Problem {
	pool := []probVerdict{
		{"The strong law says the sample mean converges almost surely to the population mean", true, "this is the conclusion"},
		{"The strong law implies the weak law", true, "almost-sure convergence is stronger"},
		{"The weak law implies the strong law", false, "the implication only goes one way"},
		{"Kolmogorov's strong law needs only a finite first moment for iid variables", true, "finite variance is not required"},
		{"The strong law requires finite variance", false, "that is a sufficient condition in simple proofs, not necessary"},
		{"The strong law applies to iid Bernoulli variables", true, "the empirical frequency converges almost surely"},
		{"The strong law says the sample mean converges in distribution", false, "it asserts almost-sure convergence"},
		{"The strong law can fail when the mean is infinite", true, "the sample mean then has no finite limit"},
		{"The strong law holds for every sequence of independent variables", false, "identically distributed or moment conditions are needed"},
		{"The strong law is proved with the Borel-Cantelli lemmas", true, "summable tail probabilities give almost-sure convergence"},
		{"The strong law guarantees the sample mean is eventually constant", false, "it converges but need not become constant"},
		{"The strong law is a limit theorem about averages, not sums", true, "sums grow like n times the mean"},
		{"The strong law holds for iid variables with finite fourth moment", true, "a standard sufficient condition"},
		{"The strong law says the empirical distribution converges weakly", true, "the Glivenko-Cantelli theorem extends this idea"},
		{"The strong law requires the variables to be bounded", false, "unbounded variables with a finite mean are allowed"},
		{"The strong law and the central limit theorem concern the same fluctuations", false, "the CLT gives 1/sqrt(n) fluctuations"},
	}
	return countProbVerdicts("How many of these statements about the strong law of large numbers are true?", pool)
}

func borelCantelliProblem() generator.Problem {
	pool := []probVerdict{
		{"If the sum of P(A_n) is finite then P(A_n infinitely often) = 0", true, "the first Borel-Cantelli lemma"},
		{"If the sum of P(A_n) diverges and the events are independent then P(A_n infinitely often) = 1", true, "the second Borel-Cantelli lemma"},
		{"The second lemma needs some independence or mixing hypothesis", true, "divergence alone is not enough"},
		{"The first lemma requires independence", false, "it holds for every sequence of events"},
		{"P(A_n infinitely often) is the probability of the limsup of the events", true, "this is the definition"},
		{"A summable sequence of probabilities can still have infinitely many occurrences almost surely", false, "the first lemma forbids that"},
		{"The lemmas are used to prove the strong law of large numbers", true, "summable tails give almost-sure convergence"},
		{"If P(A_n) = 1/n then the sum diverges", true, "the harmonic series diverges"},
		{"For independent events with P(A_n) = 1/n, infinitely many occur almost surely", true, "the second lemma applies"},
		{"The first lemma gives an upper bound on the expected number of occurrences", true, "the sum of probabilities bounds the expectation"},
		{"The lemmas require the events to be disjoint", false, "disjointness is irrelevant"},
		{"If the sum of P(A_n) converges then only finitely many A_n occur", true, "almost surely"},
		{"The second lemma applies to pairwise independent events in general", false, "the standard version needs full independence"},
		{"Borel-Cantelli can show that a sequence converges almost surely", true, "control the tail events"},
		{"The first lemma says the events cannot occur infinitely often with positive probability", true, "the limsup has probability zero"},
		{"Divergence of the sum of probabilities always implies infinitely many occurrences", false, "independence or a mixing condition is needed"},
	}
	return countProbVerdicts("How many of these statements about the Borel-Cantelli lemmas are true?", pool)
}

func slutskyProblem() generator.Problem {
	pool := []probVerdict{
		{"If X_n converges in distribution to X and Y_n converges in probability to c then X_n + Y_n converges in distribution to X + c", true, "Slutsky's theorem"},
		{"Under the same hypotheses X_n Y_n converges in distribution to cX", true, "the product form of Slutsky"},
		{"If c is nonzero then X_n / Y_n converges in distribution to X/c", true, "the ratio form of Slutsky"},
		{"Slutsky's theorem does not require the two sequences to be independent", true, "the probability limit of Y_n is constant"},
		{"Slutsky combines convergence in distribution with convergence in probability", true, "that is its role"},
		{"Slutsky's theorem requires X_n and Y_n to be independent", false, "independence is not assumed"},
		{"Slutsky's theorem allows Y_n to converge to a nonconstant random variable", false, "the limit must be a constant for the standard form"},
		{"If Y_n converges in probability to 0 then X_n / Y_n converges in distribution to X/0", false, "the ratio need not exist"},
		{"Slutsky's theorem is used to justify the t statistic", true, "the sample variance converges in probability to sigma^2"},
		{"Slutsky's theorem says X_n + Y_n converges in probability", false, "the conclusion is distributional"},
		{"If X_n converges in distribution and Y_n converges in probability to c, then X_n Y_n converges in probability to cX", false, "distributional convergence is the conclusion"},
		{"Slutsky's theorem is a consequence of the continuous mapping theorem", true, "apply the map (x,y) -> x + y"},
		{"Slutsky's theorem requires finite means", false, "no moments are needed"},
		{"Slutsky's theorem can be applied when Y_n is a consistent estimator of a parameter", true, "consistency is convergence in probability"},
		{"The theorem fails if Y_n converges in distribution but not in probability to a constant", true, "the standard hypotheses must hold"},
		{"Slutsky's theorem preserves the limiting distribution under addition of a vanishing term", true, "a term tending to 0 in probability does not change the limit"},
	}
	return countProbVerdicts("How many of these statements about Slutsky's theorem are true?", pool)
}

func characteristicFunctionProblem() generator.Problem {
	pool := []probVerdict{
		{"The characteristic function of X is phi(t) = E[e^{itX}]", true, "this is the definition"},
		{"Every random variable has a characteristic function", true, "the complex exponential is bounded"},
		{"The characteristic function is uniformly continuous", true, "it is a Fourier transform of a probability measure"},
		{"phi(0) = 1 for every characteristic function", true, "e^0 = 1"},
		{"|phi(t)| <= 1 for every t", true, "the modulus of an expectation of a unit-modulus variable"},
		{"The characteristic function determines the distribution", true, "inversion recovers the law"},
		{"Independent variables have characteristic functions that multiply", true, "phi_{X+Y} = phi_X phi_Y"},
		{"A real characteristic function corresponds to a symmetric distribution", true, "the imaginary part vanishes"},
		{"The characteristic function requires a density", false, "every distribution has one"},
		{"A characteristic function can exceed 1 in modulus", false, "bounded by 1"},
		{"The characteristic function is always differentiable", false, "differentiability depends on moments"},
		{"If the k-th moment is finite then phi has a k-th derivative at 0", true, "differentiate under the expectation"},
		{"The characteristic function is the Fourier transform of the distribution", true, "up to a sign convention"},
		{"Two variables with the same characteristic function have the same distribution", true, "uniqueness follows by inversion"},
		{"The characteristic function of a point mass at a is e^{iat}", true, "direct computation"},
		{"The characteristic function is always real-valued", false, "it is complex in general"},
	}
	return countProbVerdicts("How many of these statements about characteristic functions are true?", pool)
}

func characteristicSumsProblem() generator.Problem {
	pool := []probVerdict{
		{"For independent X and Y, phi_{X+Y} = phi_X phi_Y", true, "the expectation factorises"},
		{"For iid variables, the characteristic function of the sum is phi^n", true, "multiply n copies"},
		{"The characteristic function of a sum is always the product of the characteristic functions", false, "independence is required"},
		{"The characteristic function of a normal variable is e^{i mu t - sigma^2 t^2/2}", true, "standard formula"},
		{"The characteristic function of a Poisson variable is exp(lambda (e^{it} - 1))", true, "standard formula"},
		{"The characteristic function of a sum of independent normals is normal", true, "multiply the exponentials"},
		{"Convolution of densities corresponds to multiplication of characteristic functions", true, "Fourier transform turns convolution into a product"},
		{"The characteristic function of X + c is e^{ict} phi_X(t)", true, "a shift multiplies by a phase"},
		{"The characteristic function of -X is the complex conjugate of phi_X", true, "conjugate the exponential"},
		{"The product of two characteristic functions need not be a characteristic function", false, "it is the characteristic function of the sum of independent copies"},
		{"The characteristic function of a Cauchy sum is the product of the factors", true, "independence gives the product"},
		{"The characteristic function of a sum is the sum of the characteristic functions", false, "that is not how Fourier transforms behave"},
		{"If X and Y are independent then phi_{X-Y} = phi_X conjugate(phi_Y)", true, "apply the shift and conjugate rules"},
		{"The characteristic function of an average is phi(t/n)^n", true, "scale the sum"},
		{"Independence is irrelevant for characteristic functions of sums", false, "the product formula needs it"},
		{"The characteristic function of a sum of independent indicators is a product of Bernoulli factors", true, "each factor is 1 - p + p e^{it}"},
	}
	return countProbVerdicts("How many of these statements about characteristic functions of sums are true?", pool)
}

func continuityTheoremProblem() generator.Problem {
	pool := []probVerdict{
		{"If characteristic functions converge pointwise to a function continuous at 0 then the laws converge weakly", true, "Levy's continuity theorem"},
		{"Weak convergence of laws implies pointwise convergence of characteristic functions", true, "bounded continuous test functions"},
		{"The limit function must be continuous at 0 for the theorem to apply", true, "otherwise mass can escape"},
		{"A pointwise limit of characteristic functions that is continuous at 0 is itself a characteristic function", true, "the limit is the transform of the limiting law"},
		{"The continuity theorem is used to prove the central limit theorem", true, "converge the characteristic functions to e^{-t^2/2}"},
		{"Pointwise convergence of characteristic functions is enough even if the limit is discontinuous at 0", false, "mass can escape to infinity"},
		{"The continuity theorem requires the distributions to have densities", false, "no densities are needed"},
		{"The theorem gives a one-to-one correspondence between laws and characteristic functions", true, "uniqueness plus inversion"},
		{"Convergence of characteristic functions at a single point implies weak convergence", false, "pointwise convergence on all of R is needed"},
		{"The theorem is a special case of Prokhorov's theorem", true, "tightness follows from continuity at 0"},
		{"The continuity theorem says convergence in probability is equivalent to convergence of characteristic functions", false, "it characterises weak convergence"},
		{"If the limit characteristic function is continuous at 0 then the sequence is tight", true, "Levy's tightness criterion"},
		{"The theorem can fail if the limiting function is not continuous at 0", true, "point masses can drift to infinity"},
		{"The continuity theorem requires finite means", false, "no moments are needed"},
		{"The characteristic function of the limit is the pointwise limit of the characteristic functions", true, "this is the conclusion"},
		{"Weak convergence and pointwise convergence of characteristic functions are equivalent", true, "Levy's theorem in both directions"},
	}
	return countProbVerdicts("How many of these statements about the continuity theorem are true?", pool)
}

func characteristicCLTProblem() generator.Problem {
	pool := []probVerdict{
		{"The characteristic function of the standardised sum tends to e^{-t^2/2}", true, "this is the Fourier proof of the CLT"},
		{"The limiting characteristic function e^{-t^2/2} belongs to a standard normal variable", true, "standard normal transform"},
		{"The CLT requires iid variables with finite variance", true, "the classical hypotheses"},
		{"The characteristic-function proof expands log phi near 0", true, "a Taylor expansion gives the quadratic term"},
		{"The CLT says the standardised sum converges in distribution to N(0,1)", true, "this is the conclusion"},
		{"The CLT requires the variables to be normal", false, "only the limit is normal"},
		{"The CLT gives almost-sure convergence of the standardised sum", false, "the convergence is in distribution"},
		{"The Lindeberg condition gives a triangular-array CLT", true, "it replaces identical distribution"},
		{"The characteristic function of the average tends to a degenerate limit without scaling", true, "the law of large numbers"},
		{"The CLT needs finite fourth moments", false, "finite variance suffices"},
		{"The characteristic-function proof uses the continuity theorem", true, "pointwise convergence plus continuity at 0"},
		{"The CLT implies the weak law of large numbers", true, "distributional convergence to a constant gives probability convergence"},
		{"The CLT applies to sums of independent but not identically distributed variables under conditions", true, "Lindeberg-Feller"},
		{"The standardised sum has variance 1 in the limit", true, "the normal limit has unit variance"},
		{"The CLT says the characteristic function of the sum converges pointwise to 1", false, "after standardisation the limit is e^{-t^2/2}"},
		{"The CLT can fail when the variance is infinite", true, "stable laws can replace the normal limit"},
	}
	return countProbVerdicts("How many of these statements about the characteristic-function proof of the CLT are true?", pool)
}

func martingaleProblem() generator.Problem {
	pool := []probVerdict{
		{"A martingale satisfies E[X_{n+1} | F_n] = X_n", true, "this is the definition"},
		{"A martingale models a fair game", true, "the expected next value equals the current value"},
		{"A submartingale satisfies E[X_{n+1} | F_n] >= X_n", true, "the definition of a submartingale"},
		{"A supermartingale satisfies E[X_{n+1} | F_n] <= X_n", true, "the definition of a supermartingale"},
		{"Every martingale is a submartingale", false, "equality is a special case, not the general one"},
		{"Optional stopping holds for every martingale and every stopping time", false, "integrability or boundedness conditions are needed"},
		{"A martingale has constant expectation", true, "E[X_n] = E[X_0] when integrable"},
		{"A martingale must be a sequence of independent variables", false, "only the conditional expectations are constrained"},
		{"Sums of independent mean-zero variables form a martingale", true, "the partial sums have the martingale property"},
		{"A martingale can be unbounded", true, "integrability does not give boundedness"},
		{"Doob's martingale inequality bounds tail probabilities", true, "it controls maxima by the final moment"},
		{"Every martingale converges almost surely", false, "convergence needs additional hypotheses"},
		{"A nonnegative martingale converges almost surely", true, "the martingale convergence theorem"},
		{"The martingale property is preserved under stopping", true, "the stopped process is a martingale"},
		{"A martingale difference sequence has conditional mean zero", true, "this is the increment form"},
		{"A martingale must have finite state space", false, "state spaces can be general"},
	}
	return countProbVerdicts("How many of these statements about martingales are true?", pool)
}
