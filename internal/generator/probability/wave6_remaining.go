package probability

// Wave 6 probability theory, remaining batch. The generator uses small exact
// models: renewal counts, two-state queues, Ito moments, finite Gaussians,
// HMMs, Kalman updates, concentration bounds, Bayes risks and entropy on
// powers of two. Every answer is an integer and every instance is in Facts.

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave6Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"prob.process.renewal":       "renewal",
		"prob.process.queue":         "queue",
		"prob.process.ito":           "ito",
		"prob.process.sde":           "sde",
		"prob.process.gaussian":      "gaussian",
		"prob.process.hidden_markov": "hmm",
		"prob.process.kalman":        "kalman",
		"prob.conc.hoeffding":        "hoeffding",
		"prob.conc.bernstein":        "bernstein",
		"prob.conc.mcdiarmid":        "mcdiarmid",
		"prob.limit.lindeberg":       "lindeberg",
		"prob.limit.berry_esseen":    "berry",
		"prob.limit.local_clt":       "local_clt",
		"prob.bayes.conjugate":       "conjugate",
		"prob.bayes.predictive":      "predictive",
		"prob.decision.bayes_risk":   "bayes_risk",
		"prob.decision.minimax":      "minimax",
		"prob.info.entropy":          "entropy",
		"prob.info.kl":               "kl",
	} {
		reg.Register(id, &w6ProbRemainingGen{kind: kind})
	}
}

type w6ProbRemainingGen struct{ kind string }

func (g *w6ProbRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "renewal":
		return w6Renewal(kp)
	case "queue":
		return w6Queue(kp)
	case "ito":
		return w6Ito(kp)
	case "sde":
		return w6SDE(kp)
	case "gaussian":
		return w6Gaussian(kp)
	case "hmm":
		return w6HMM(kp)
	case "kalman":
		return w6Kalman(kp)
	case "hoeffding":
		return w6Hoeffding(kp)
	case "bernstein":
		return w6Bernstein(kp)
	case "mcdiarmid":
		return w6McDiarmid(kp)
	case "lindeberg":
		return w6Lindeberg(kp)
	case "berry":
		return w6Berry(kp)
	case "local_clt":
		return w6LocalCLT(kp)
	case "conjugate":
		return w6Conjugate(kp)
	case "predictive":
		return w6Predictive(kp)
	case "bayes_risk":
		return w6BayesRisk(kp)
	case "minimax":
		return w6Minimax(kp)
	case "entropy":
		return w6Entropy(kp)
	case "kl":
		return w6KL(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w6PM(o int, q, a, ex string, f map[string]string) generator.Problem {
	if f == nil {
		f = map[string]string{}
	}
	f["answer"] = a
	return generator.Problem{Question: q, Answer: a, Explanation: ex, Facts: f, KPOrdinal: o}
}

func w6PRI(m map[string]string, k string, v int) { m[k] = strconv.Itoa(v) }

func w6PRem(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

func w6Renewal(kp int) generator.Problem {
	mu := w6PRem(2, 12)
	t := mu * w6PRem(3, 9)
	f := map[string]string{"mu": strconv.Itoa(mu), "t": strconv.Itoa(t)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A renewal process has deterministic interarrival time %d. How many renewals occur by time %d?", mu, t), strconv.Itoa(t/mu),
			"Divide the horizon by the interarrival time.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A renewal process has deterministic interarrival time %d. What is the mean interarrival time?", mu), strconv.Itoa(mu),
			"The mean is the deterministic interarrival time.", f)
	default:
		return w6PM(kp, fmt.Sprintf("A renewal process has deterministic interarrival time %d. What is the expected number of renewals by time %d?", mu, 2*t), strconv.Itoa(2*t/mu),
			"Twice the earlier horizon gives twice the renewal count.", f)
	}
}

type w6QueuePair struct{ lambda, mu int }

func w6Queue(kp int) generator.Problem {
	pairs := []w6QueuePair{{1, 2}, {2, 4}, {3, 6}, {2, 3}, {3, 4}, {4, 5}, {1, 3}, {2, 5}, {3, 5}, {5, 6}, {1, 4}, {4, 6}}
	p := pairs[rand.Intn(len(pairs))]
	rho100 := 100 * p.lambda / p.mu
	empty100 := 100 - rho100
	L := p.lambda / (p.mu - p.lambda)
	f := map[string]string{"lambda": strconv.Itoa(p.lambda), "mu": strconv.Itoa(p.mu)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("An M/M/1 queue has arrival rate %d and service rate %d. What is 100 times its traffic intensity?", p.lambda, p.mu), strconv.Itoa(rho100),
			"Traffic intensity is lambda/mu.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("An M/M/1 queue has arrival rate %d and service rate %d. What is 100 times the stationary probability of an empty system?", p.lambda, p.mu), strconv.Itoa(empty100),
			"The empty probability is 1-rho.", f)
	default:
		return w6PM(kp, fmt.Sprintf("An M/M/1 queue has arrival rate %d and service rate %d. What is the expected number of customers in the system?", p.lambda, p.mu), strconv.Itoa(L),
			"The mean number in the system is rho/(1-rho).", f)
	}
}

func w6Ito(kp int) generator.Problem {
	t := w6PRem(1, 6)
	c := w6PRem(1, 12)
	f := map[string]string{"t": strconv.Itoa(t), "c": strconv.Itoa(c)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("For the Ito integral %d times the integral from 0 to %d of dB_s, what is the mean?", c, t), "0",
			"Ito integrals of deterministic integrands have mean zero.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For the Ito integral %d times the integral from 0 to %d of dB_s, what is the variance?", c, t), strconv.Itoa(c*c*t),
			"The variance is the integral of the squared integrand.", f)
	default:
		return w6PM(kp, fmt.Sprintf("Ito's formula gives B_t^2 = t + C times the integral of B_s dB_s. What is C for the process %d B_t?", c), strconv.Itoa(2*c),
			"For c B_t, squaring introduces the factor c^2 and the cross term has coefficient 2c.", f)
	}
}

func w6SDE(kp int) generator.Problem {
	a := w6PRem(1, 5)
	b := w6PRem(1, 4)
	h := w6PRem(1, 4)
	x0 := w6PRem(0, 5)
	f := map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "h": strconv.Itoa(h), "x0": strconv.Itoa(x0)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("The SDE dX_t = %d dt + %d dW_t has what drift coefficient?", a, b), strconv.Itoa(a),
			"The drift is the dt coefficient.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("The SDE dX_t = %d dt + %d dW_t has what diffusion coefficient?", a, b), strconv.Itoa(b),
			"The diffusion is the dW coefficient.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For dX_t = %d dt + %d dW_t with X_0=%d, what is the Euler mean after step h=%d?", a, b, x0, h), strconv.Itoa(x0+a*h),
			"The Brownian increment has mean zero, so the mean advances by a*h.", f)
	}
}

func w6Gaussian(kp int) generator.Problem {
	s := w6PRem(1, 4)
	t := s + w6PRem(1, 5)
	n := w6PRem(2, 12)
	f := map[string]string{"s": strconv.Itoa(s), "t": strconv.Itoa(t), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("For Brownian motion, what is Cov(B_%d,B_%d)?", s, t), strconv.Itoa(s),
			"The covariance is min(s,t).", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A Gaussian process is observed at %d time points. What is the dimension of its finite-dimensional distribution?", n), strconv.Itoa(n),
			"One observation per time point gives an n-dimensional Gaussian vector.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For Brownian motion with s=%d and t=%d, E[B_t|B_s] equals (s/t) B_s. What is t times the coefficient s/t?", s, t), strconv.Itoa(s),
			"Multiply the coefficient by t.", f)
	}
}

func w6HMM(kp int) generator.Problem {
	k := w6PRem(2, 12)
	e00 := w6PRem(1, 9)
	e01 := w6PRem(1, 9)
	f := map[string]string{"states": strconv.Itoa(k), "emit00": strconv.Itoa(e00), "emit01": strconv.Itoa(e01)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A hidden Markov model has %d states. How many free transition parameters are needed when each row of the transition matrix sums to one?", k), strconv.Itoa(k*(k-1)),
			"Each of the k rows has k-1 free entries.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For a two-state HMM, the emission probability of symbol 0 from state 0 is %d/10. What is 10 times that probability?", e00), strconv.Itoa(e00),
			"The probability is already stated in tenths.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For a two-state HMM, the emission probability of symbol 1 from state 0 is %d/10. What is 10 times that probability?", e01), strconv.Itoa(e01),
			"The probability is already stated in tenths.", f)
	}
}

func w6Kalman(kp int) generator.Problem {
	P := w6PRem(2, 6)
	Q := w6PRem(1, 4)
	R := w6PRem(1, 4)
	prior := w6PRem(0, 9)
	z := w6PRem(0, 9)
	f := map[string]string{"P": strconv.Itoa(P), "Q": strconv.Itoa(Q), "R": strconv.Itoa(R), "prior": strconv.Itoa(prior), "z": strconv.Itoa(z)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A Kalman filter has prior variance %d and process noise variance %d. What is the predicted variance?", P, Q), strconv.Itoa(P+Q),
			"The predicted variance is P+Q.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A Kalman update has predicted variance %d and measurement noise variance %d. What is (P+R) times the gain P/(P+R)?", P, R), strconv.Itoa(P),
			"Multiply the gain by P+R.", f)
	default:
		return w6PM(kp, fmt.Sprintf("A Kalman update has P=%d, R=%d, prior mean %d and measurement %d. What is (P+R) times the posterior mean?", P, R, prior, z), strconv.Itoa(R*prior+P*z),
			"The posterior mean is a weighted average with weights R and P.", f)
	}
}

func w6Hoeffding(kp int) generator.Problem {
	n := 2 * w6PRem(1, 10)
	f := map[string]string{"n": strconv.Itoa(n), "eps": "1/2"}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("For n=%d and eps=1/2, Hoeffding's bound is 2 exp(-2 n eps^2). What is the exponent 2 n eps^2?", n), strconv.Itoa(n/2),
			"Substitute eps=1/2.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For n=%d and eps=1/2, what is the Hoeffding exponent after the sample size is doubled?", n), strconv.Itoa(2*n),
			"The exponent is n/2 at sample size n, so doubling n gives exponent n.", f)
	default:
		t := w6PRem(1, 12)
		f["target"] = strconv.Itoa(t)
		return w6PM(kp, fmt.Sprintf("For eps=1/2, what sample size makes the Hoeffding exponent 2 n eps^2 equal to %d?", t), strconv.Itoa(2*t),
			"Solve n/2=t.", f)
	}
}

func w6Bernstein(kp int) generator.Problem {
	v := w6PRem(1, 12) // variance numerator over 4
	n := w6PRem(2, 6)
	den := 3*v + 2
	f := map[string]string{"v_num": strconv.Itoa(v), "v_den": "4", "n": strconv.Itoa(n), "den": strconv.Itoa(den)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("For variance %d/4 and eps=1/2, the Bernstein denominator is (3v+2)/6. What is 12 times that denominator?", v), strconv.Itoa(6*v+4),
			"Multiply (3v+2)/6 by 12: 2(3v+2) = 6v+4.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For variance %d/4 and eps=1/2, what is 3v+2 in the Bernstein exponent denominator?", v), strconv.Itoa(den),
			"Use the scaled denominator.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For variance %d/4, after the variance is doubled, what is the numerator over 6 of the Bernstein denominator?", v), strconv.Itoa(6*v+2),
			"Double the variance term: 3(2v)+2 = 6v+2.", f)
	}
}

func w6McDiarmid(kp int) generator.Problem {
	n := 2 * w6PRem(2, 10)
	t := n / 2
	f := map[string]string{"n": strconv.Itoa(n), "t": strconv.Itoa(t)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A function of %d independent variables has bounded differences all equal to 1. What is the sum of squared bounds?", n), strconv.Itoa(n),
			"Sum n ones.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For n=%d and t=%d, McDiarmid's exponent 2t^2/n equals what integer?", n, t), strconv.Itoa(n/2),
			"Substitute t=n/2.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For n=%d and t=%d, by what factor does the McDiarmid exponent change when t is doubled?", n, t), "4",
			"The exponent is 2t^2/n, so doubling t multiplies the exponent by 4.", f)
	}
}

func w6Lindeberg(kp int) generator.Problem {
	total := w6PRem(10, 20)
	max := w6PRem(1, total/4)
	f := map[string]string{"total": strconv.Itoa(total), "max": strconv.Itoa(max)}
	switch kp {
	case 1:
		ok := 0
		if max*10 < total {
			ok = 1
		}
		return w6PM(kp, fmt.Sprintf("A triangular array has total variance %d and largest individual variance %d. Does the small-check Lindeberg condition max*10 < total hold? Enter 1 for yes and 0 for no.", total, max), strconv.Itoa(ok),
			"Compare the largest variance with one tenth of the total.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A triangular array has total variance %d and largest individual variance %d. What is the integer ratio total/max?", total, max), strconv.Itoa(total/max),
			"Divide the total variance by the maximum variance.", f)
	default:
		fail := 0
		if max*10 >= total {
			fail = 1
		}
		return w6PM(kp, fmt.Sprintf("For total variance %d and largest variance %d, does the small-check Lindeberg condition fail? Enter 1 for yes and 0 for no.", total, max), strconv.Itoa(fail),
			"The condition fails when max*10 is at least the total.", f)
	}
}

func w6Berry(kp int) generator.Problem {
	sigma := w6PRem(1, 10)
	n := w6PRem(4, 15)
	f := map[string]string{"sigma": strconv.Itoa(sigma), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		k := w6PRem(1, 8)
		rho := sigma * sigma * sigma * k
		f["rho"] = strconv.Itoa(rho)
		f["k"] = strconv.Itoa(k)
		return w6PM(kp, fmt.Sprintf("A Berry-Esseen bound has sigma=%d and third absolute moment %d. The moment is sigma cubed times an integer k. What is k?", sigma, rho), strconv.Itoa(k),
			"Divide the third absolute moment by sigma cubed.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A Berry-Esseen bound has sigma=%d. What is sigma cubed?", sigma), strconv.Itoa(sigma*sigma*sigma),
			"Cube sigma.", f)
	default:
		return w6PM(kp, fmt.Sprintf("A Berry-Esseen bound uses n=%d and sigma=%d. What is n times sigma cubed?", n, sigma), strconv.Itoa(n*sigma*sigma*sigma),
			"Multiply the sample size by sigma cubed.", f)
	}
}

func w6LocalCLT(kp int) generator.Problem {
	v := w6PRem(1, 10)
	n := w6PRem(1, 10)
	f := map[string]string{"v": strconv.Itoa(v), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("For a sample of size n=%d, the local central limit normalising factor is proportional to n to what exponent denominator? Enter the denominator.", n), "2",
			"The factor is proportional to n^(-1/2), so the exponent denominator is 2.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("For %d lattice variables each of variance %d, what is the variance of their sum?", n, v), strconv.Itoa(n*v),
			"Variances add for independent variables.", f)
	default:
		return w6PM(kp, fmt.Sprintf("For a sample of size n=%d, the local central limit normalising factor is proportional to n^(-1/2). What is the positive exponent numerator?", n), "1",
			"The factor is proportional to n^(-1/2), so the positive numerator is 1.", f)
	}
}

func w6Conjugate(kp int) generator.Problem {
	a := w6PRem(1, 5)
	b := w6PRem(1, 5)
	s := w6PRem(1, 5)
	fail := w6PRem(1, 5)
	f := map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "successes": strconv.Itoa(s), "failures": strconv.Itoa(fail)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) prior sees %d successes. What is the posterior alpha?", a, b, s), strconv.Itoa(a+s),
			"Add successes to alpha.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) prior sees %d failures. What is the posterior beta?", a, b, fail), strconv.Itoa(b+fail),
			"Add failures to beta.", f)
	default:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) prior sees %d successes and %d failures. What is the posterior total strength?", a, b, s, fail), strconv.Itoa(a+b+s+fail),
			"Add the four counts.", f)
	}
}

func w6Predictive(kp int) generator.Problem {
	a := w6PRem(1, 5)
	b := w6PRem(1, 5)
	s := w6PRem(1, 5)
	fail := w6PRem(1, 5)
	total := a + b + s + fail
	f := map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "successes": strconv.Itoa(s), "failures": strconv.Itoa(fail), "total": strconv.Itoa(total)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) posterior after %d successes and %d failures is used to predict the next trial. What is the numerator of the predictive success probability?", a, b, s, fail), strconv.Itoa(a+s),
			"The predictive numerator is alpha plus successes.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) prior sees %d successes and %d failures. What is the denominator of the predictive success probability?", a, b, s, fail), strconv.Itoa(total),
			"The denominator is the posterior total strength.", f)
	default:
		return w6PM(kp, fmt.Sprintf("A Beta(%d,%d) prior sees %d successes and %d failures. What is the numerator of the predictive failure probability?", a, b, s, fail), strconv.Itoa(b+fail),
			"The failure numerator is beta plus failures.", f)
	}
}

func w6BayesRisk(kp int) generator.Problem {
	p := w6PRem(2, 8)
	l1 := w6PRem(1, 6)
	l2 := w6PRem(1, 6)
	r1 := p * l1
	r2 := (10 - p) * l2
	f := map[string]string{"prior": strconv.Itoa(p), "loss1": strconv.Itoa(l1), "loss2": strconv.Itoa(l2)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("Hypothesis 1 has prior %d/10. Action 1 loses %d under hypothesis 1 and 0 otherwise. What is 10 times its Bayes risk?", p, l1), strconv.Itoa(r1),
			"Risk is prior times loss.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("Hypothesis 2 has prior (10-%d)/10. Action 2 loses %d under hypothesis 2 and 0 otherwise. What is 10 times its Bayes risk?", p, l2), strconv.Itoa(r2),
			"Risk is prior times loss.", f)
	default:
		action := 1
		if r2 < r1 {
			action = 2
		}
		return w6PM(kp, fmt.Sprintf("For the two actions with scaled risks %d and %d, which action is the Bayes action? Enter 1 or 2.", r1, r2), strconv.Itoa(action),
			"Choose the action with the smaller risk.", f)
	}
}

func w6Minimax(kp int) generator.Problem {
	a, b := w6PRem(1, 6), w6PRem(1, 6)
	c, d := w6PRem(1, 6), w6PRem(1, 6)
	r1, r2 := maxInt(a, b), maxInt(c, d)
	f := map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d)}
	switch kp {
	case 1:
		return w6PM(kp, fmt.Sprintf("A 2x2 loss matrix has row 1 entries %d and %d. What is the maximum loss in row 1?", a, b), strconv.Itoa(r1),
			"Take the larger row entry.", f)
	case 2:
		return w6PM(kp, fmt.Sprintf("A 2x2 loss matrix has row 2 entries %d and %d. What is the maximum loss in row 2?", c, d), strconv.Itoa(r2),
			"Take the larger row entry.", f)
	default:
		return w6PM(kp, fmt.Sprintf("The two row maxima are %d and %d. What is the minimax value?", r1, r2), strconv.Itoa(minInt(r1, r2)),
			"The minimax value is the smaller row maximum.", f)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func w6Entropy(kp int) generator.Problem {
	probs := [][]int{
		{1, 1}, {2, 2}, {4, 4},
		{1, 1, 2}, {2, 2, 4}, {4, 4, 8},
		{1, 1, 1, 1}, {2, 2, 2, 2}, {4, 4, 4, 4},
		{1, 1, 1, 1, 4},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 2, 2, 2},
	}
	p := probs[rand.Intn(len(probs))]
	f := map[string]string{"weights": joinInts(p)}
	switch kp {
	case 1:
		total := 0
		for _, w := range p {
			total += w
		}
		weighted := 0
		for _, w := range p {
			weighted += w * (log2Int(total) - log2Int(w))
		}
		twoH := 2 * weighted / total
		return w6PM(kp, fmt.Sprintf("A distribution has weights %s over outcomes. What is 2 times its entropy in bits?", joinInts(p)), strconv.Itoa(twoH),
			"Entropy is the sum of p log2(1/p); multiply the result by 2.", f)
	case 2:
		n := w6PRem(2, 10)
		f["n"] = strconv.Itoa(n)
		return w6PM(kp, fmt.Sprintf("Two independent variables have %d equally likely outcomes each. What is 2 times the entropy of the pair in bits?", n), strconv.Itoa(2*log2Int(n*n)),
			"Independent entropy adds.", f)
	default:
		return w6PM(kp, fmt.Sprintf("What is the maximum entropy in bits of a distribution on %d equally likely outcomes?", len(p)), strconv.Itoa(log2Int(len(p))),
			"The uniform distribution maximises entropy.", f)
	}
}

func log2Int(n int) int {
	c := 0
	for n > 1 {
		n /= 2
		c++
	}
	return c
}

func joinInts(xs []int) string {
	var s strings.Builder
	for i, x := range xs {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(strconv.Itoa(x))
	}
	return s.String()
}

func w6KL(kp int) generator.Problem {
	f := map[string]string{}
	switch kp {
	case 1:
		n := w6PRem(2, 16)
		f["n"] = strconv.Itoa(n)
		return w6PM(kp, fmt.Sprintf("What is the KL divergence in bits from a point mass to the uniform distribution on %d outcomes?", n), strconv.Itoa(log2Int(n)),
			"The point mass assigns probability 1 to an event of probability 1/n, giving log2 n.", f)
	case 2:
		n := w6PRem(2, 12)
		f["n"] = strconv.Itoa(n)
		return w6PM(kp, fmt.Sprintf("Two independent fair variables with %d outcomes each have what mutual information in bits?", n), "0",
			"Independent variables have zero mutual information.", f)
	default:
		a := w6PRem(1, 9)
		b := w6PRem(1, 9)
		f["a"] = strconv.Itoa(a)
		f["b"] = strconv.Itoa(b)
		return w6PM(kp, fmt.Sprintf("For independent variables, H(X,Y)=H(X)+H(Y). If H(X)=%d and H(Y)=%d, what is H(X,Y)?", a, b), strconv.Itoa(a+b),
			"Add the entropies.", f)
	}
}
