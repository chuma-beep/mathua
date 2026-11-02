package probability

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// probDecimal formats a probability or a transformed value without trailing
// zeros, so an exact integer does not read as a measurement with false
// precision.
func probDecimal(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func probDecimal4(v float64) string {
	return strconv.FormatFloat(v, 'f', 4, 64)
}

// ---------------------------------------------------------------------------
// Poisson distribution
// ---------------------------------------------------------------------------

type poissonDistributionGen struct{}

func (g *poissonDistributionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	lambda := 1 + rand.Intn(5)
	switch probBridgeOrdinal(ctx) {
	case 1:
		k := rand.Intn(5)
		p := math.Exp(-float64(lambda)) * math.Pow(float64(lambda), float64(k)) / float64(factorialInt(k))
		return generator.Problem{
			Question:    fmt.Sprintf("A Poisson random variable has rate %d. What is P(X = %d)? (enter a decimal)", lambda, k),
			Answer:      probDecimal4(p),
			Explanation: fmt.Sprintf("P(X = %d) = e^(-%d) * %d^%d / %d! = %s.", k, lambda, lambda, k, k, probDecimal4(p)),
			Facts:       map[string]string{"lambda": strconv.Itoa(lambda), "k": strconv.Itoa(k), "answer": probDecimal4(p)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		kind := "mean"
		if rand.Intn(2) == 0 {
			kind = "variance"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("A Poisson random variable has rate %d. What is its %s? (enter a number)", lambda, kind),
			Answer:      strconv.Itoa(lambda),
			Explanation: fmt.Sprintf("For a Poisson variable the mean and the variance both equal the rate, so the %s is %d.", kind, lambda),
			Facts:       map[string]string{"lambda": strconv.Itoa(lambda), "kind": kind, "answer": strconv.Itoa(lambda)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		p := 1 - math.Exp(-float64(lambda))
		return generator.Problem{
			Question:    fmt.Sprintf("A Poisson random variable has rate %d. What is P(X >= 1)? (enter a decimal)", lambda),
			Answer:      probDecimal4(p),
			Explanation: fmt.Sprintf("P(X >= 1) = 1 - P(X = 0) = 1 - e^(-%d) = %s.", lambda, probDecimal4(p)),
			Facts:       map[string]string{"lambda": strconv.Itoa(lambda), "answer": probDecimal4(p)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

func factorialInt(n int) int {
	out := 1
	for i := 2; i <= n; i++ {
		out *= i
	}
	return out
}

// ---------------------------------------------------------------------------
// Normal distribution
// ---------------------------------------------------------------------------

type normalDistributionGen struct{}

func (g *normalDistributionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mu := -2 + rand.Intn(8)
	sigma := 1 + rand.Intn(4)
	switch probBridgeOrdinal(ctx) {
	case 1:
		d := -3 + rand.Intn(7)
		if d == 0 {
			d = 1
		}
		x := mu + d*sigma
		z := float64(d)
		return generator.Problem{
			Question:    fmt.Sprintf("A normal variable has mean %d and standard deviation %d. What is the z score of x = %d? (enter a number)", mu, sigma, x),
			Answer:      probDecimal(z),
			Explanation: fmt.Sprintf("z = (x - mu) / sigma = (%d - %d) / %d = %s.", x, mu, sigma, probDecimal(z)),
			Facts:       map[string]string{"mu": strconv.Itoa(mu), "sigma": strconv.Itoa(sigma), "x": strconv.Itoa(x), "answer": probDecimal(z)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		k := 1 + rand.Intn(3)
		outside := rand.Intn(2) == 1
		percent := map[int]string{1: "68", 2: "95", 3: "99.7"}[k]
		if outside {
			percent = map[int]string{1: "32", 2: "5", 3: "0.3"}[k]
		}
		where := "within"
		if outside {
			where = "outside"
		}
		return generator.Problem{
			Question:    fmt.Sprintf("Using the empirical rule, what percent of a normal distribution lies %s %d standard deviations of the mean? (enter a number)", where, k),
			Answer:      percent,
			Explanation: fmt.Sprintf("The empirical rule gives about %s percent %s %d standard deviations.", percent, where, k),
			Facts:       map[string]string{"k": strconv.Itoa(k), "where": where, "answer": percent},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		k := 1 + rand.Intn(3)
		x := mu + k*sigma
		return generator.Problem{
			Question:    fmt.Sprintf("A normal variable has mean %d and standard deviation %d. What is the upper endpoint of the interval within %d standard deviations? (enter a number)", mu, sigma, k),
			Answer:      strconv.Itoa(x),
			Explanation: fmt.Sprintf("The upper endpoint is mu + %d sigma = %d + %d * %d = %d.", k, mu, k, sigma, x),
			Facts:       map[string]string{"mu": strconv.Itoa(mu), "sigma": strconv.Itoa(sigma), "k": strconv.Itoa(k), "answer": strconv.Itoa(x)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Skewness and kurtosis
// ---------------------------------------------------------------------------

type skewKurtosisGen struct{}

func (g *skewKurtosisGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		type shape struct {
			text string
			code int
		}
		shapes := []shape{
			{"A distribution has a long right tail and its mean exceeds its median.", 1},
			{"A distribution has a long left tail and its median exceeds its mean.", -1},
			{"A distribution is symmetric about its centre.", 0},
			{"Most observations sit near the low end with a few very large values.", 1},
			{"Most observations sit near the high end with a few very small values.", -1},
			{"The mean and median coincide and the two tails mirror each other.", 0},
		}
		s := shapes[rand.Intn(len(shapes))]
		return generator.Problem{
			Question:    fmt.Sprintf("%s What is the sign of the skewness? Enter 1 for positive, -1 for negative or 0 for zero.", s.text),
			Answer:      strconv.Itoa(s.code),
			Explanation: fmt.Sprintf("The described shape has skewness sign %d.", s.code),
			Facts:       map[string]string{"answer": strconv.Itoa(s.code)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		type dist struct {
			name string
			val  string
		}
		ds := []dist{
			{"the normal distribution", "0"},
			{"the uniform distribution on an interval", "-1.2"},
			{"the exponential distribution", "6"},
			{"the Bernoulli distribution with p = 0.5", "-2"},
			{"the Laplace distribution", "3"},
			{"the Poisson distribution with rate 2", "0.5"},
		}
		d := ds[rand.Intn(len(ds))]
		return generator.Problem{
			Question:    fmt.Sprintf("What is the excess kurtosis of %s? (enter a decimal)", d.name),
			Answer:      d.val,
			Explanation: fmt.Sprintf("The standardised kurtosis of %s gives excess kurtosis %s.", d.name, d.val),
			Facts:       map[string]string{"answer": d.val},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		skew := []string{"0.5", "-0.5", "1", "-1"}[rand.Intn(4)]
		sigma := 2 + rand.Intn(3)
		s, _ := strconv.ParseFloat(skew, 64)
		mu3 := s * math.Pow(float64(sigma), 3)
		return generator.Problem{
			Question:    fmt.Sprintf("A distribution has skewness %s and standard deviation %d. What is its third central moment? (enter a number)", skew, sigma),
			Answer:      probDecimal(mu3),
			Explanation: fmt.Sprintf("Skewness is mu3 / sigma^3, so mu3 = %s * %d^3 = %s.", skew, sigma, probDecimal(mu3)),
			Facts:       map[string]string{"skew": skew, "sigma": strconv.Itoa(sigma), "answer": probDecimal(mu3)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Conditional densities
// ---------------------------------------------------------------------------

type conditionalDensityGen struct{}

func (g *conditionalDensityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	width := 2 + rand.Intn(4)
	y := 2 + rand.Intn(8)
	switch probBridgeOrdinal(ctx) {
	case 1:
		a := rand.Intn(5)
		val := 1 / float64(width)
		return generator.Problem{
			Question:    fmt.Sprintf("X is uniform on the interval from %d to %d given Y = y. What is the conditional density f(x | y)? (enter a decimal)", a, a+width),
			Answer:      probDecimal4(val),
			Explanation: fmt.Sprintf("A uniform density on an interval of width %d has height 1/%d = %s.", width, width, probDecimal4(val)),
			Facts:       map[string]string{"width": strconv.Itoa(width), "lo": strconv.Itoa(a), "answer": probDecimal4(val)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		mean := float64(y) / 2
		return generator.Problem{
			Question:    fmt.Sprintf("Given Y = %d, X is uniform on the interval from 0 to %d. What is E[X | Y = %d]? (enter a number)", y, y, y),
			Answer:      probDecimal(mean),
			Explanation: fmt.Sprintf("The midpoint of [0, %d] is %d / 2 = %s.", y, y, probDecimal(mean)),
			Facts:       map[string]string{"y": strconv.Itoa(y), "answer": probDecimal(mean)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		c := 1 + rand.Intn(y)
		p := float64(c) / float64(y)
		return generator.Problem{
			Question:    fmt.Sprintf("Given Y = %d, X is uniform on [0, %d]. What is P(X <= %d | Y = %d)? (enter a decimal)", y, y, c, y),
			Answer:      probDecimal4(p),
			Explanation: fmt.Sprintf("The conditional probability is the interval length %d divided by %d, which is %s.", c, y, probDecimal4(p)),
			Facts:       map[string]string{"y": strconv.Itoa(y), "c": strconv.Itoa(c), "answer": probDecimal4(p)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Probability generating functions
// ---------------------------------------------------------------------------

type pgfGen struct{}

func (g *pgfGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		type model struct {
			name string
			p0   string
		}
		ms := []model{
			{"a Bernoulli variable with success probability 0.3", "0.7"},
			{"a Bernoulli variable with success probability 0.5", "0.5"},
			{"a Bernoulli variable with success probability 0.8", "0.2"},
			{"a geometric variable with success probability 0.25", "0.25"},
			{"a geometric variable with success probability 0.5", "0.5"},
			{"a Poisson variable with rate 1", "0.3679"},
			{"a Poisson variable with rate 2", "0.1353"},
			{"a Poisson variable with rate 3", "0.0498"},
		}
		m := ms[rand.Intn(len(ms))]
		return generator.Problem{
			Question:    fmt.Sprintf("For %s, what is G(0) = P(X = 0)? (enter a decimal)", m.name),
			Answer:      m.p0,
			Explanation: fmt.Sprintf("G(0) keeps only the constant term of the generating function, so it is P(X = 0) = %s.", m.p0),
			Facts:       map[string]string{"answer": m.p0},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		p := 0.2 + 0.1*float64(rand.Intn(6))
		kind := rand.Intn(3)
		var mean float64
		var name string
		switch kind {
		case 0:
			name = fmt.Sprintf("a Bernoulli variable with p = %.1f", p)
			mean = p
		case 1:
			name = fmt.Sprintf("a geometric variable with p = %.1f", p)
			mean = (1 - p) / p
		default:
			lambda := 1 + rand.Intn(5)
			name = fmt.Sprintf("a Poisson variable with rate %d", lambda)
			mean = float64(lambda)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("For %s, use G'(1) to find the mean. (enter a decimal)", name),
			Answer:      probDecimal(mean),
			Explanation: fmt.Sprintf("The derivative of the generating function at 1 is the mean, which is %s.", probDecimal(mean)),
			Facts:       map[string]string{"answer": probDecimal(mean)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		p := 0.2 + 0.1*float64(rand.Intn(6))
		kind := rand.Intn(2)
		var variance float64
		var name string
		if kind == 0 {
			name = fmt.Sprintf("a Bernoulli variable with p = %.1f", p)
			variance = p * (1 - p)
		} else {
			name = fmt.Sprintf("a geometric variable with p = %.1f", p)
			variance = (1 - p) / (p * p)
		}
		return generator.Problem{
			Question:    fmt.Sprintf("For %s, what is the variance? (enter a decimal)", name),
			Answer:      probDecimal4(variance),
			Explanation: fmt.Sprintf("Using factorial moments from the generating function gives variance %s.", probDecimal4(variance)),
			Facts:       map[string]string{"answer": probDecimal4(variance)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Recurrence and return times
// ---------------------------------------------------------------------------

type recurrenceGen struct{}

func (g *recurrenceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		self := []string{"0.25", "0.4", "0.5", "0.6", "0.75"}[rand.Intn(5)]
		return generator.Problem{
			Question:    fmt.Sprintf("A Markov chain has one-step transition probability P(i, i) = %s. What is the one-step return probability to state i? (enter a decimal)", self),
			Answer:      self,
			Explanation: fmt.Sprintf("The one-step return probability is the diagonal entry P(i, i) = %s.", self),
			Facts:       map[string]string{"answer": self},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		pi := []string{"0.05", "0.1", "0.125", "0.2", "0.25", "0.4", "0.5"}[rand.Intn(7)]
		p, _ := strconv.ParseFloat(pi, 64)
		mean := 1 / p
		return generator.Problem{
			Question:    fmt.Sprintf("An irreducible positive recurrent chain has stationary probability %s for state i. What is the expected return time to i? (enter a number)", pi),
			Answer:      probDecimal(mean),
			Explanation: fmt.Sprintf("The expected return time is 1 / pi_i = 1 / %s = %s.", pi, probDecimal(mean)),
			Facts:       map[string]string{"pi": pi, "answer": probDecimal(mean)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		a := []string{"0.2", "0.4", "0.5", "0.6", "0.8"}[rand.Intn(5)]
		b := []string{"0.3", "0.5", "0.7"}[rand.Intn(3)]
		av, _ := strconv.ParseFloat(a, 64)
		bv, _ := strconv.ParseFloat(b, 64)
		val := av*av + (1-av)*(1-bv)
		return generator.Problem{
			Question:    fmt.Sprintf("A two-state chain has P(1,1) = %s and P(2,1) = %s. What is the two-step return probability to state 1? (enter a decimal)", a, b),
			Answer:      probDecimal4(val),
			Explanation: fmt.Sprintf("(P^2)(1,1) = P(1,1)P(1,1) + P(1,2)P(2,1) = %s * %s + %s * %s = %s.", a, a, probDecimal(1-av), b, probDecimal4(val)),
			Facts:       map[string]string{"answer": probDecimal4(val)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Optional stopping
// ---------------------------------------------------------------------------

type stoppingGen struct{}

func (g *stoppingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	n := 4 + rand.Intn(8)
	k := 1 + rand.Intn(n-1)
	switch probBridgeOrdinal(ctx) {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A martingale starts at %d and is stopped at a bounded stopping time T. What is E[X_T]? (enter a number)", k),
			Answer:      strconv.Itoa(k),
			Explanation: fmt.Sprintf("Optional stopping gives E[X_T] = E[X_0] = %d when its hypotheses hold.", k),
			Facts:       map[string]string{"start": strconv.Itoa(k), "answer": strconv.Itoa(k)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		val := k * (n - k)
		return generator.Problem{
			Question:    fmt.Sprintf("A symmetric random walk starts at %d and stops at 0 or %d. What is the expected number of steps until absorption? (enter a number)", k, n),
			Answer:      strconv.Itoa(val),
			Explanation: fmt.Sprintf("The expected absorption time is k(N - k) = %d * %d = %d.", k, n-k, val),
			Facts:       map[string]string{"k": strconv.Itoa(k), "n": strconv.Itoa(n), "answer": strconv.Itoa(val)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		p := float64(k) / float64(n)
		return generator.Problem{
			Question:    fmt.Sprintf("A symmetric random walk starts at %d and stops at 0 or %d. What is the probability it reaches %d before 0? (enter a decimal)", k, n, n),
			Answer:      probDecimal4(p),
			Explanation: fmt.Sprintf("The hitting probability is k / N = %d / %d = %s.", k, n, probDecimal4(p)),
			Facts:       map[string]string{"k": strconv.Itoa(k), "n": strconv.Itoa(n), "answer": probDecimal4(p)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

// ---------------------------------------------------------------------------
// Transformations of random vectors
// ---------------------------------------------------------------------------

type multivariateTransformGen struct{}

func (g *multivariateTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	a := 1 + rand.Intn(4)
	b := 1 + rand.Intn(4)
	c := 1 + rand.Intn(3)
	d := 1 + rand.Intn(4)
	det := a*d - b*c
	if det == 0 {
		det = a*d + 1
	}
	switch probBridgeOrdinal(ctx) {
	case 1:
		return generator.Problem{
			Question:    fmt.Sprintf("A linear map has matrix [[%d, %d], [%d, %d]]. What is the absolute value of its Jacobian determinant? (enter a number)", a, b, c, d),
			Answer:      strconv.Itoa(absInt(det)),
			Explanation: fmt.Sprintf("The determinant is %d * %d - %d * %d = %d, whose absolute value is %d.", a, d, b, c, det, absInt(det)),
			Facts:       map[string]string{"a": strconv.Itoa(a), "b": strconv.Itoa(b), "c": strconv.Itoa(c), "d": strconv.Itoa(d), "answer": strconv.Itoa(absInt(det))},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		val := 1 / float64(absInt(det))
		return generator.Problem{
			Question:    fmt.Sprintf("A uniform density of 1 on the unit square is mapped by a matrix with determinant %d. What is the density on the image region? (enter a decimal)", absInt(det)),
			Answer:      probDecimal4(val),
			Explanation: fmt.Sprintf("Densities scale by 1 / |det|, so the new density is 1 / %d = %s.", absInt(det), probDecimal4(val)),
			Facts:       map[string]string{"det": strconv.Itoa(absInt(det)), "answer": probDecimal4(val)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		p := 1 + rand.Intn(4)
		q := 1 + rand.Intn(4)
		j := absInt(p * q)
		return generator.Problem{
			Question:    fmt.Sprintf("The map u = %d x and v = %d y has Jacobian determinant %d. What is the absolute determinant used in the density formula? (enter a number)", p, q, p*q),
			Answer:      strconv.Itoa(j),
			Explanation: fmt.Sprintf("The determinant of the diagonal map is %d * %d = %d.", p, q, p*q),
			Facts:       map[string]string{"answer": strconv.Itoa(j)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ---------------------------------------------------------------------------
// Probability integral transform
// ---------------------------------------------------------------------------

type integralTransformGen struct{}

func (g *integralTransformGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	lambda := 1 + rand.Intn(3)
	x := 1 + rand.Intn(6)
	switch probBridgeOrdinal(ctx) {
	case 1:
		f := 1 - math.Exp(-float64(lambda*x))
		return generator.Problem{
			Question:    fmt.Sprintf("An exponential variable has rate %d. What is F(%d)? (enter a decimal)", lambda, x),
			Answer:      probDecimal4(f),
			Explanation: fmt.Sprintf("F(x) = 1 - e^(-lambda x) = 1 - e^(-%d * %d) = %s.", lambda, x, probDecimal4(f)),
			Facts:       map[string]string{"lambda": strconv.Itoa(lambda), "x": strconv.Itoa(x), "answer": probDecimal4(f)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	case 2:
		u := []string{"0.2", "0.35", "0.5", "0.7", "0.9"}[rand.Intn(5)]
		return generator.Problem{
			Question:    fmt.Sprintf("U = F(X) is the probability integral transform of a continuous variable. What is P(U <= %s)? (enter a decimal)", u),
			Answer:      u,
			Explanation: fmt.Sprintf("The transform makes U uniform on [0, 1], so P(U <= %s) = %s.", u, u),
			Facts:       map[string]string{"u": u, "answer": u},
			KPOrdinal:   ctx.KPOrdinal,
		}
	default:
		u := []string{"0.25", "0.5", "0.75"}[rand.Intn(3)]
		uv, _ := strconv.ParseFloat(u, 64)
		v := -math.Log(1-uv) / float64(lambda)
		return generator.Problem{
			Question:    fmt.Sprintf("An exponential variable has rate %d. Use the inverse transform with U = %s. What is the simulated value? (enter a decimal)", lambda, u),
			Answer:      probDecimal4(v),
			Explanation: fmt.Sprintf("x = -ln(1 - U) / lambda = -ln(%s) / %d = %s.", probDecimal(1-uv), lambda, probDecimal4(v)),
			Facts:       map[string]string{"lambda": strconv.Itoa(lambda), "u": u, "answer": probDecimal4(v)},
			KPOrdinal:   ctx.KPOrdinal,
		}
	}
}
