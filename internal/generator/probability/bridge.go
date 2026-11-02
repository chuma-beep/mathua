package probability

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// probBridgeOrdinal returns the requested knowledge point, or a random one when
// the caller is drawing unguided practice.

func probBridgeOrdinal(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

func chooseInt(n, k int) int {
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

// ---------------------------------------------------------------------------
// Conditional probability (numeric)
// ---------------------------------------------------------------------------

type conditionalDefinitionGen struct{}

func (g *conditionalDefinitionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		nPool := []int{20, 24, 25, 30, 40, 50, 60, 80}
		bPool := []int{4, 5, 8, 10, 20, 25, 40, 50}
		n := nPool[rand.Intn(len(nPool))]
		b := bPool[rand.Intn(len(bPool))]
		if b >= n {
			b = n / 2
		}
		a := 1 + rand.Intn(b-1)
		ans := decimal6(float64(a) / float64(b))
		return generator.Problem{
			Question: fmt.Sprintf("A sample space has %d equally likely outcomes. Event B contains %d outcomes, and %d of those also lie in event A. What is P(A | B)? (enter a decimal)", n, b, a),
			Answer:   ans,
			Explanation: fmt.Sprintf("Condition on B: the denominator is the %d outcomes in B, not the full %d. The joint outcomes number %d, so P(A | B) = %d/%d = %s.",
				b, n, a, a, b, ans),
			Facts: map[string]string{"n": strconv.Itoa(n), "b": strconv.Itoa(b), "both": strconv.Itoa(a), "answer": ans},
		}
	case 2:
		pbPool := []float64{0.2, 0.25, 0.4, 0.5, 0.8}
		qPool := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.75}
		pb := pbPool[rand.Intn(len(pbPool))]
		q := qPool[rand.Intn(len(qPool))]
		pab := pb * q
		pabS := decimal6(pab)
		pbS := decimal6(pb)
		ans := decimal6(q)
		return generator.Problem{
			Question:    fmt.Sprintf("For events A and B, P(A and B) = %s and P(B) = %s. What is P(A | B)? (enter a decimal)", pabS, pbS),
			Answer:      ans,
			Explanation: fmt.Sprintf("By definition P(A | B) = P(A and B)/P(B) = %s/%s = %s.", pabS, pbS, ans),
			Facts:       map[string]string{"joint": pabS, "condition": pbS, "answer": ans},
		}
	default:
		nPool := []int{60, 80, 100, 120}
		n := nPool[rand.Intn(len(nPool))]
		a := 12 + rand.Intn(19)
		b := 12 + rand.Intn(19)
		both := 1 + rand.Intn(11)
		if both >= a || both >= b {
			both = 1
		}
		if a+b-both > n {
			n = a + b - both + 10
		}
		kind := rand.Intn(2)
		var num, den int
		var symbol string
		if kind == 0 {
			num, den, symbol = both, b, "P(A | B)"
		} else {
			num, den, symbol = both, a, "P(B | A)"
		}
		ans := decimal6(float64(num) / float64(den))
		return generator.Problem{
			Question: fmt.Sprintf("A survey has %d people. %d have property A, %d have property B, and %d have both properties. What is %s? (enter a decimal)", n, a, b, both, symbol),
			Answer:   ans,
			Explanation: fmt.Sprintf("The conditioning event supplies the denominator. Here %s = %d/%d = %s, so only the people in the conditioning group are counted.",
				symbol, num, den, ans),
			Facts: map[string]string{"n": strconv.Itoa(n), "a": strconv.Itoa(a), "b": strconv.Itoa(b), "both": strconv.Itoa(both), "answer": ans},
		}
	}
}

// ---------------------------------------------------------------------------
// Convolution of random variables (numeric)
// ---------------------------------------------------------------------------

type convolutionGen struct{}

func (g *convolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		kind := rand.Intn(5)
		switch kind {
		case 0:
			l1 := 1 + rand.Intn(5)
			l2 := 1 + rand.Intn(5)
			ans := l1 + l2
			return generator.Problem{
				Question:    fmt.Sprintf("X is Poisson with rate %d and Y is Poisson with rate %d, independently. The sum X + Y is Poisson with what rate? (enter a number)", l1, l2),
				Answer:      strconv.Itoa(ans),
				Explanation: fmt.Sprintf("Independent Poisson rates add, so the rate of X + Y is %d + %d = %d.", l1, l2, ans),
				Facts:       map[string]string{"l1": strconv.Itoa(l1), "l2": strconv.Itoa(l2), "answer": strconv.Itoa(ans)},
			}
		case 1:
			m1 := -3 + rand.Intn(7)
			m2 := -3 + rand.Intn(7)
			ans := m1 + m2
			return generator.Problem{
				Question:    fmt.Sprintf("X is normal with mean %d and Y is normal with mean %d, independently. What is the mean of X + Y? (enter a number)", m1, m2),
				Answer:      strconv.Itoa(ans),
				Explanation: fmt.Sprintf("Means add regardless of dependence; here the mean is %d + %d = %d.", m1, m2, ans),
				Facts:       map[string]string{"m1": strconv.Itoa(m1), "m2": strconv.Itoa(m2), "answer": strconv.Itoa(ans)},
			}
		case 2:
			v1 := 1 + rand.Intn(6)
			v2 := 1 + rand.Intn(6)
			ans := v1 + v2
			return generator.Problem{
				Question:    fmt.Sprintf("Independent variables X and Y have variances %d and %d. What is the variance of X + Y? (enter a number)", v1, v2),
				Answer:      strconv.Itoa(ans),
				Explanation: fmt.Sprintf("Variances add for independent variables, so Var(X + Y) = %d + %d = %d.", v1, v2, ans),
				Facts:       map[string]string{"v1": strconv.Itoa(v1), "v2": strconv.Itoa(v2), "answer": strconv.Itoa(ans)},
			}
		case 3:
			n1 := 2 + rand.Intn(8)
			n2 := 2 + rand.Intn(8)
			ans := n1 + n2
			return generator.Problem{
				Question:    fmt.Sprintf("X is binomial with %d trials and Y is binomial with %d trials, independently, and both use the same success probability. The sum X + Y is binomial with how many trials? (enter a number)", n1, n2),
				Answer:      strconv.Itoa(ans),
				Explanation: fmt.Sprintf("The trials concatenate, so the trial count of X + Y is %d + %d = %d.", n1, n2, ans),
				Facts:       map[string]string{"n1": strconv.Itoa(n1), "n2": strconv.Itoa(n2), "answer": strconv.Itoa(ans)},
			}
		default:
			m := 1 + rand.Intn(8)
			ans := m
			return generator.Problem{
				Question:    fmt.Sprintf("X and Y are independent and each is uniformly distributed on the integers 1 through %d. The sum X + Y takes the value %d in how many ordered pairs? (enter a number)", m, m+1),
				Answer:      strconv.Itoa(ans),
				Explanation: fmt.Sprintf("The pairs that sum to %d are (1,%d), (2,%d), and so on through (%d,1), which is %d pairs.", m+1, m, m-1, m, ans),
				Facts:       map[string]string{"m": strconv.Itoa(m), "answer": strconv.Itoa(ans)},
			}
		}
	case 2:
		if rand.Intn(3) == 0 {
			pPool := []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
			p := pPool[rand.Intn(len(pPool))]
			ans := decimal6(2 * p * (1 - p))
			return generator.Problem{
				Question:    fmt.Sprintf("X and Y are independent Bernoulli variables with success probability %s. What is P(X + Y = 1)? (enter a decimal)", decimal6(p)),
				Answer:      ans,
				Explanation: fmt.Sprintf("The sum is 1 in the two orders (1,0) and (0,1), so the probability is 2(%s)(%s) = %s.", decimal6(p), decimal6(1-p), ans),
				Facts:       map[string]string{"p": decimal6(p), "answer": ans},
			}
		}
		m := 4 + rand.Intn(7)
		ans := decimal6(1 / float64(m))
		return generator.Problem{
			Question:    fmt.Sprintf("X and Y are independent and each is uniform on the integers 1 through %d. What is P(X + Y = %d)? (enter a decimal)", m, m+1),
			Answer:      ans,
			Explanation: fmt.Sprintf("There are %d pairs in total and %d pairs sum to %d, so the probability is %d/%d = %s.", m*m, m, m+1, m, m*m, ans),
			Facts:       map[string]string{"m": strconv.Itoa(m), "answer": ans},
		}
	default:
		if rand.Intn(3) == 0 {
			lPool := []float64{0.5, 1, 2}
			l := lPool[rand.Intn(len(lPool))]
			ans := decimal6(l / math.E)
			return generator.Problem{
				Question:    fmt.Sprintf("X and Y are independent exponential variables with rate %s. Their sum has density f(z) = z e^{-z} scaled for the rate. What is the density at z = %s? (enter a decimal)", decimal6(l), decimal6(1/l)),
				Answer:      ans,
				Explanation: fmt.Sprintf("The sum is Gamma(2, rate %s), with density (rate^2) z e^(-rate z). At z = %s the value is %s.", decimal6(l), decimal6(1/l), ans),
				Facts:       map[string]string{"rate": decimal6(l), "z": decimal6(1 / l), "answer": ans},
			}
		}
		xPool := []float64{0.25, 0.5, 0.75}
		x := xPool[rand.Intn(len(xPool))]
		upper := rand.Intn(2) == 0
		z := x
		ans := x
		if upper {
			z = 2 - x
			ans = 2 - z
		}
		ansS := decimal6(ans)
		return generator.Problem{
			Question:    fmt.Sprintf("X and Y are independent uniform variables on the interval from 0 to 1. The density of X + Y at z = %s is what? (enter a decimal)", decimal6(z)),
			Answer:      ansS,
			Explanation: fmt.Sprintf("The sum has a triangular density: f(z) = z below 1 and f(z) = 2 - z above 1. At z = %s the density is %s.", decimal6(z), ansS),
			Facts:       map[string]string{"z": decimal6(z), "answer": ansS},
		}
	}
}

// ---------------------------------------------------------------------------
// The binomial distribution (numeric)
// ---------------------------------------------------------------------------

type binomialDistributionGen struct{}

func (g *binomialDistributionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch probBridgeOrdinal(ctx) {
	case 1:
		n := 2 + rand.Intn(5)
		k := rand.Intn(n + 1)
		ans := decimal6(float64(chooseInt(n, k)) / math.Pow(2, float64(n)))
		return generator.Problem{
			Question: fmt.Sprintf("A fair coin is flipped %d times. What is the probability of exactly %d heads? (enter a decimal)", n, k),
			Answer:   ans,
			Explanation: fmt.Sprintf("There are %d equally likely sequences, and %d of them have exactly %d heads, so the probability is %d/%d = %s.",
				int(math.Pow(2, float64(n))), chooseInt(n, k), k, chooseInt(n, k), int(math.Pow(2, float64(n))), ans),
			Facts: map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k), "answer": ans},
		}
	case 2:
		n := 3 + rand.Intn(4)
		k := rand.Intn(n + 1)
		pPool := []float64{0.2, 0.25, 0.4, 0.5, 0.6, 0.75, 0.8}
		p := pPool[rand.Intn(len(pPool))]
		prob := float64(chooseInt(n, k)) * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
		ans := decimal6(prob)
		return generator.Problem{
			Question: fmt.Sprintf("A binomial experiment has %d independent trials with success probability %s. What is the probability of exactly %d successes? (enter a decimal)", n, decimal6(p), k),
			Answer:   ans,
			Explanation: fmt.Sprintf("Use C(%d,%d) p^k (1-p)^(n-k) = %d times (%s)^%d times (%s)^%d = %s.",
				n, k, chooseInt(n, k), decimal6(p), k, decimal6(1-p), n-k, ans),
			Facts: map[string]string{"n": strconv.Itoa(n), "k": strconv.Itoa(k), "p": decimal6(p), "answer": ans},
		}
	default:
		n := 5 + rand.Intn(16)
		pPool := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
		p := pPool[rand.Intn(len(pPool))]
		kind := rand.Intn(3)
		switch kind {
		case 0:
			ans := decimal6(float64(n) * p)
			return generator.Problem{
				Question:    fmt.Sprintf("A binomial variable has %d trials and success probability %s. What is its mean? (enter a decimal)", n, decimal6(p)),
				Answer:      ans,
				Explanation: fmt.Sprintf("The mean of a binomial variable is np = %d times %s = %s.", n, decimal6(p), ans),
				Facts:       map[string]string{"n": strconv.Itoa(n), "p": decimal6(p), "answer": ans},
			}
		case 1:
			ans := decimal6(float64(n) * p * (1 - p))
			return generator.Problem{
				Question:    fmt.Sprintf("A binomial variable has %d trials and success probability %s. What is its variance? (enter a decimal)", n, decimal6(p)),
				Answer:      ans,
				Explanation: fmt.Sprintf("The variance of a binomial variable is np(1-p) = %d times %s times %s = %s.", n, decimal6(p), decimal6(1-p), ans),
				Facts:       map[string]string{"n": strconv.Itoa(n), "p": decimal6(p), "answer": ans},
			}
		default:
			ans := decimal6(math.Sqrt(float64(n) * p * (1 - p)))
			return generator.Problem{
				Question:    fmt.Sprintf("A binomial variable has %d trials and success probability %s. What is its standard deviation? (enter a decimal)", n, decimal6(p)),
				Answer:      ans,
				Explanation: fmt.Sprintf("The standard deviation is the square root of np(1-p) = %d times %s times %s, which is %s.", n, decimal6(p), decimal6(1-p), ans),
				Facts:       map[string]string{"n": strconv.Itoa(n), "p": decimal6(p), "answer": ans},
			}
		}
	}
}
