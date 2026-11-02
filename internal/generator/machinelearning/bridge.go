package machinelearning

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// The 2026 content expansion adds ten concepts whose questions are computed
// rather than recalled. Each concept carries one generator with three
// section-aware shapes, selected by GeneratorContext.KPOrdinal.
var mlBridgeKinds = map[string]string{
	"ml.svm.margin":           "svm_margin",
	"ml.svm.kernel":           "svm_kernel",
	"ml.naive_bayes":          "naive_bayes",
	"ml.bias_variance":        "bias_variance",
	"ml.mdp.bellman":          "bellman",
	"ml.cluster.em":           "cluster_em",
	"ml.feature.selection":    "feature_selection",
	"ml.multiclass":           "multiclass",
	"ml.nlp.embedding":        "nlp_embedding",
	"ml.optim.early_stopping": "early_stopping",
}

type mlBridgeGen struct{ kind string }

func (g *mlBridgeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	o := ctx.KPOrdinal
	if o < 1 || o > 3 {
		o = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "svm_margin":
		return mlSVMarginProblem(o)
	case "svm_kernel":
		return mlSVMKernelProblem(o)
	case "naive_bayes":
		return mlNaiveBayesProblem(o)
	case "bias_variance":
		return mlBiasVarianceProblem(o)
	case "bellman":
		return mlBellmanProblem(o)
	case "cluster_em":
		return mlClusterEMProblem(o)
	case "feature_selection":
		return mlFeatureSelectionProblem(o)
	case "multiclass":
		return mlMulticlassProblem(o)
	case "nlp_embedding":
		return mlNLPEmbeddingProblem(o)
	case "early_stopping":
		return mlEarlyStoppingProblem(o)
	}
	return mlBridgeProblem(o, "What is 2 + 2? (enter a number)", "4", "Two plus two is four.")
}

func mlBridgeProblem(o int, question, answer, explanation string) generator.Problem {
	return generator.Problem{
		Question:    question,
		Answer:      answer,
		Explanation: explanation,
		KPOrdinal:   o,
		Facts:       map[string]string{"answer": answer},
	}
}

func mlBridgeDecimal(v float64, places int) string {
	s := strconv.FormatFloat(v, 'f', places, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "-0" || s == "" {
		s = "0"
	}
	return s
}

func mlBridgePow(base, exp int) int {
	r := 1
	for i := 0; i < exp; i++ {
		r *= base
	}
	return r
}

func mlSVMarginProblem(o int) generator.Problem {
	switch o {
	case 1:
		triples := [][2]int{{3, 4}, {5, 12}, {8, 15}, {7, 24}, {20, 21}, {9, 40}}
		t := triples[rand.Intn(len(triples))]
		norm := math.Sqrt(float64(t[0]*t[0] + t[1]*t[1]))
		ans := 2 / norm
		return mlBridgeProblem(o,
			fmt.Sprintf("A support vector machine has weight vector w = (%d, %d). What is the width of its margin, 2 divided by the norm of w? (enter a decimal)", t[0], t[1]),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("The norm of w is the square root of %d squared plus %d squared, so the margin width is 2 divided by that norm, about %s.", t[0], t[1], mlBridgeDecimal(ans, 4)))
	case 2:
		a := 1 + rand.Intn(8)
		b := 1 + rand.Intn(8)
		ans := a*a + b*b
		return mlBridgeProblem(o,
			fmt.Sprintf("For the weight vector w = (%d, %d), what is the squared norm of w? (enter a number)", a, b),
			strconv.Itoa(ans),
			fmt.Sprintf("The squared norm is the sum of squares: %d squared plus %d squared equals %d.", a, b, ans))
	default:
		widths := []string{"0.25", "0.4", "0.5", "1", "2", "4"}
		w := widths[rand.Intn(len(widths))]
		width, _ := strconv.ParseFloat(w, 64)
		ans := 2 / width
		return mlBridgeProblem(o,
			fmt.Sprintf("A maximum-margin classifier has margin width %s. What is the norm of its weight vector? (enter a decimal)", w),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("Because the margin width is 2 divided by the norm, the norm is 2 divided by %s, which is %s.", w, mlBridgeDecimal(ans, 4)))
	}
}

func mlSVMKernelProblem(o int) generator.Problem {
	switch o {
	case 1:
		x1, x2 := 1+rand.Intn(3), 1+rand.Intn(3)
		y1, y2 := 1+rand.Intn(3), 1+rand.Intn(3)
		c := rand.Intn(3)
		d := 2 + rand.Intn(2)
		dot := x1*y1 + x2*y2
		ans := mlBridgePow(dot+c, d)
		return mlBridgeProblem(o,
			fmt.Sprintf("For the polynomial kernel (x dot y + %d)^%d, evaluate it on x = (%d, %d) and y = (%d, %d). (enter a number)", c, d, x1, x2, y1, y2),
			strconv.Itoa(ans),
			fmt.Sprintf("The dot product is %d. Adding %d gives %d, and raising it to the power %d gives %d.", dot, c, dot+c, d, ans))
	case 2:
		x1, x2 := 1+rand.Intn(6), 1+rand.Intn(6)
		y1, y2 := 1+rand.Intn(6), 1+rand.Intn(6)
		ans := (x1-y1)*(x1-y1) + (x2-y2)*(x2-y2)
		return mlBridgeProblem(o,
			fmt.Sprintf("An RBF kernel uses the squared distance between x = (%d, %d) and y = (%d, %d). What is that squared distance? (enter a number)", x1, x2, y1, y2),
			strconv.Itoa(ans),
			fmt.Sprintf("Subtract coordinates to get (%d, %d), then add the squares: %d squared plus %d squared equals %d.", x1-y1, x2-y2, x1-y1, x2-y2, ans))
	default:
		n := 3 + rand.Intn(10)
		ans := n * (n + 1) / 2
		return mlBridgeProblem(o,
			fmt.Sprintf("How many distinct entries are in the kernel matrix for %d training points? (enter a number)", n),
			strconv.Itoa(ans),
			fmt.Sprintf("The matrix is symmetric, so it has n squared entries but only n times n plus one divided by two distinct ones, which is %d.", ans))
	}
}

func mlNaiveBayesProblem(o int) generator.Problem {
	switch o {
	case 1:
		n := 5 + rand.Intn(16)
		count := rand.Intn(n + 1)
		k := 2 + rand.Intn(4)
		ans := float64(count+1) / float64(n+k)
		return mlBridgeProblem(o,
			fmt.Sprintf("A feature value occurs %d times in %d trials and has %d possible values. Using Laplace smoothing, what is the estimated probability? (enter a decimal)", count, n, k),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("Laplace smoothing adds one to the count and %d to the denominator, giving %d divided by %d, about %s.", k, count+1, n+k, mlBridgeDecimal(ans, 4)))
	case 2:
		k := 2 + rand.Intn(3)
		d := 2 + rand.Intn(5)
		v := 2 + rand.Intn(4)
		ans := k * d * v
		return mlBridgeProblem(o,
			fmt.Sprintf("A multinomial naive Bayes model has %d classes, %d features and %d values per feature. How many conditional probabilities does it store? (enter a number)", k, d, v),
			strconv.Itoa(ans),
			fmt.Sprintf("Every class and feature pair needs one probability per value, so the count is %d times %d times %d equals %d.", k, d, v, ans))
	default:
		k := 2 + rand.Intn(3)
		d := 2 + rand.Intn(6)
		ans := k * (2*d + 1)
		return mlBridgeProblem(o,
			fmt.Sprintf("A Bernoulli naive Bayes model has %d classes and %d binary features. How many parameters does it store, including the class priors? (enter a number)", k, d),
			strconv.Itoa(ans),
			fmt.Sprintf("Each class stores two probabilities per feature plus one prior, so the total is %d times the quantity 2 times %d plus 1, which is %d.", k, d, ans))
	}
}

func mlBiasVarianceProblem(o int) generator.Problem {
	switch o {
	case 1:
		b := float64(1+rand.Intn(5)) / 10
		v := float64(1+rand.Intn(5)) / 10
		noise := float64(1+rand.Intn(3)) / 10
		ans := b*b + v + noise
		return mlBridgeProblem(o,
			fmt.Sprintf("A model has bias %s, variance %s and irreducible noise %s. What is its expected squared error? (enter a decimal)", mlBridgeDecimal(b, 1), mlBridgeDecimal(v, 1), mlBridgeDecimal(noise, 1)),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("The error is squared bias plus variance plus noise: %s squared plus %s plus %s equals %s.", mlBridgeDecimal(b, 1), mlBridgeDecimal(v, 1), mlBridgeDecimal(noise, 1), mlBridgeDecimal(ans, 4)))
	case 2:
		b := float64(1+rand.Intn(5)) / 10
		ans := b * b
		return mlBridgeProblem(o,
			fmt.Sprintf("A model has bias %s. What is its squared bias? (enter a decimal)", mlBridgeDecimal(b, 1)),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("Squaring the bias gives %s squared, which is %s.", mlBridgeDecimal(b, 1), mlBridgeDecimal(ans, 4)))
	default:
		v := 1 + rand.Intn(9)
		n := 2 + rand.Intn(4)
		ans := float64(v) / float64(n)
		return mlBridgeProblem(o,
			fmt.Sprintf("Each of %d independent predictions has variance %d. What is the variance of their average? (enter a decimal)", n, v),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("Averaging divides the variance by the number of predictions, so the result is %d divided by %d, which is %s.", v, n, mlBridgeDecimal(ans, 4)))
	}
}

func mlBellmanProblem(o int) generator.Problem {
	switch o {
	case 1:
		r1 := 1 + rand.Intn(20)
		v1 := 2 * (1 + rand.Intn(10))
		r2 := 1 + rand.Intn(20)
		v2 := 2 * (1 + rand.Intn(10))
		a1 := float64(r1) + 0.5*float64(v1)
		a2 := float64(r2) + 0.5*float64(v2)
		ans := math.Max(a1, a2)
		return mlBridgeProblem(o,
			fmt.Sprintf("A state has two actions. Action one gives reward %d then value %d. Action two gives reward %d then value %d. With discount factor 0.5, what is the optimal value? (enter a decimal)", r1, v1, r2, v2),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("The two action values are %s and %s, and the Bellman optimum is the larger one, %s.", mlBridgeDecimal(a1, 4), mlBridgeDecimal(a2, 4), mlBridgeDecimal(ans, 4)))
	case 2:
		r0 := 1 + rand.Intn(10)
		r1 := 1 + rand.Intn(10)
		r2 := 1 + rand.Intn(10)
		ans := float64(r0) + 0.5*float64(r1) + 0.25*float64(r2)
		return mlBridgeProblem(o,
			fmt.Sprintf("A deterministic chain gives rewards %d, then %d, then %d. With discount factor 0.5, what is the value of the first state? (enter a decimal)", r0, r1, r2),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("Sum the rewards with discount factors one, one half and one quarter to get %s.", mlBridgeDecimal(ans, 4)))
	default:
		s := 2 + rand.Intn(9)
		a := 2 + rand.Intn(8)
		ans := s * a
		return mlBridgeProblem(o,
			fmt.Sprintf("A tabular Markov decision process has %d states and %d actions. How many state-action values does it store? (enter a number)", s, a),
			strconv.Itoa(ans),
			fmt.Sprintf("Every state and action pair needs one value, so the table has %d times %d equals %d entries.", s, a, ans))
	}
}

func mlClusterEMProblem(o int) generator.Problem {
	switch o {
	case 1:
		vals := []int{1, 2, 4, 5, 10}
		p := vals[rand.Intn(len(vals))]
		q := vals[rand.Intn(len(vals))]
		for q == p {
			q = vals[rand.Intn(len(vals))]
		}
		ans := float64(p) / float64(p+q)
		return mlBridgeProblem(o,
			fmt.Sprintf("Two equally weighted mixture components have densities %d and %d at a point. What is the responsibility of the first component? (enter a decimal)", p, q),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("The responsibility is the first density divided by the sum, %d divided by %d, which is %s.", p, p+q, mlBridgeDecimal(ans, 4)))
	case 2:
		k := 2 + rand.Intn(3)
		d := 1 + rand.Intn(4)
		cov := d * (d + 1) / 2
		ans := k*(d+cov) + (k - 1)
		return mlBridgeProblem(o,
			fmt.Sprintf("A Gaussian mixture has %d components in %d dimensions with full covariances. How many free parameters does it have? (enter a number)", k, d),
			strconv.Itoa(ans),
			fmt.Sprintf("Each component needs %d means, %d covariance entries and a weight. With one weight fixed, the total is %d times %d plus %d, equals %d.", d, cov, k, d+cov, k-1, ans))
	default:
		k := 2 + rand.Intn(4)
		d := 2 + rand.Intn(4)
		ans := k * d * (d + 1) / 2
		return mlBridgeProblem(o,
			fmt.Sprintf("A Gaussian mixture has %d components in %d dimensions with full covariance matrices. How many covariance entries does it store? (enter a number)", k, d),
			strconv.Itoa(ans),
			fmt.Sprintf("Each full covariance has d times d plus one divided by two entries, so the total is %d times %d times %d divided by two, which is %d.", k, d, d+1, ans))
	}
}

func mlFeatureSelectionProblem(o int) generator.Problem {
	switch o {
	case 1:
		d := 2 + rand.Intn(11)
		ans := 1 << d
		return mlBridgeProblem(o,
			fmt.Sprintf("How many feature subsets can be formed from %d features? (enter a number)", d),
			strconv.Itoa(ans),
			fmt.Sprintf("Each feature is either included or excluded, so there are 2 to the power %d equals %d subsets.", d, ans))
	case 2:
		d := 3 + rand.Intn(13)
		ans := d * (d - 1) / 2
		return mlBridgeProblem(o,
			fmt.Sprintf("How many distinct pairs of features are there among %d features? (enter a number)", d),
			strconv.Itoa(ans),
			fmt.Sprintf("Choose two of %d without regard to order: %d times %d divided by two equals %d.", d, d, d-1, ans))
	default:
		rs := []string{"0", "0.5", "0.75", "0.8", "0.9", "0.95"}
		r := rs[rand.Intn(len(rs))]
		rf, _ := strconv.ParseFloat(r, 64)
		ans := 1 / (1 - rf)
		return mlBridgeProblem(o,
			fmt.Sprintf("A feature has squared multiple correlation %s with the other features. What is its variance inflation factor? (enter a decimal)", r),
			mlBridgeDecimal(ans, 4),
			fmt.Sprintf("The variance inflation factor is one divided by one minus %s, which is %s.", r, mlBridgeDecimal(ans, 4)))
	}
}

func mlMulticlassProblem(o int) generator.Problem {
	switch o {
	case 1:
		k := 3 + rand.Intn(8)
		return mlBridgeProblem(o,
			fmt.Sprintf("A one-versus-rest strategy handles %d classes. How many binary classifiers does it train? (enter a number)", k),
			strconv.Itoa(k),
			fmt.Sprintf("One classifier separates each class from the rest, so %d classes need %d classifiers.", k, k))
	case 2:
		k := 3 + rand.Intn(8)
		ans := k * (k - 1) / 2
		return mlBridgeProblem(o,
			fmt.Sprintf("A one-versus-one strategy handles %d classes. How many binary classifiers does it train? (enter a number)", k),
			strconv.Itoa(ans),
			fmt.Sprintf("There is one classifier per unordered pair of classes, so the count is %d times %d divided by two equals %d.", k, k-1, ans))
	default:
		n := 4 + rand.Intn(20)
		k := 3 + rand.Intn(8)
		ans := n * k
		return mlBridgeProblem(o,
			fmt.Sprintf("A softmax model scores %d examples over %d classes. How many scores does it produce? (enter a number)", n, k),
			strconv.Itoa(ans),
			fmt.Sprintf("Each example receives one score per class, so the total is %d times %d equals %d.", n, k, ans))
	}
}

func mlNLPEmbeddingProblem(o int) generator.Problem {
	switch o {
	case 1:
		v := 100 + rand.Intn(4901)
		dims := []int{50, 100, 200, 300}
		d := dims[rand.Intn(len(dims))]
		ans := v * d
		return mlBridgeProblem(o,
			fmt.Sprintf("An embedding table covers %d tokens with vectors of dimension %d. How many numbers does it store? (enter a number)", v, d),
			strconv.Itoa(ans),
			fmt.Sprintf("Every token has one vector, so the table has %d times %d equals %d entries.", v, d, ans))
	case 2:
		w := 1 + rand.Intn(5)
		n := 10 + rand.Intn(41)
		ans := 2 * w * n
		return mlBridgeProblem(o,
			fmt.Sprintf("A skip-gram model uses window size %d over %d tokens. About how many centre-context pairs does it form, ignoring edge effects? (enter a number)", w, n),
			strconv.Itoa(ans),
			fmt.Sprintf("Each token pairs with up to %d tokens on each side, so the count is 2 times %d times %d equals %d.", w, w, n, ans))
	default:
		d1 := 50 + 50*rand.Intn(5)
		d2 := 50 + 50*rand.Intn(5)
		ans := d1 + d2
		return mlBridgeProblem(o,
			fmt.Sprintf("Two embedding spaces have dimensions %d and %d. After concatenating them, what is the vector dimension? (enter a number)", d1, d2),
			strconv.Itoa(ans),
			fmt.Sprintf("Concatenation places the vectors side by side, so the dimensions add: %d plus %d equals %d.", d1, d2, ans))
	}
}

func mlEarlyStoppingProblem(o int) generator.Problem {
	switch o {
	case 1:
		b := 5 + rand.Intn(26)
		p := 2 + rand.Intn(9)
		ans := b + p
		return mlBridgeProblem(o,
			fmt.Sprintf("The best validation epoch is %d and patience is %d. At which epoch does early stopping trigger? (enter a number)", b, p),
			strconv.Itoa(ans),
			fmt.Sprintf("Training continues for %d epochs without improvement after epoch %d, so it stops at epoch %d.", p, b, ans))
	case 2:
		total := 40 + rand.Intn(61)
		b := 5 + rand.Intn(26)
		p := 2 + rand.Intn(9)
		if b+p >= total {
			total = b + p + 10
		}
		ans := total - b - p
		return mlBridgeProblem(o,
			fmt.Sprintf("Training is budgeted for %d epochs, the best epoch is %d and patience is %d. How many budgeted epochs are saved? (enter a number)", total, b, p),
			strconv.Itoa(ans),
			fmt.Sprintf("Stopping happens at epoch %d, so %d minus %d equals %d epochs of the budget are not run.", b+p, total, b+p, ans))
	default:
		b := 5 + rand.Intn(31)
		return mlBridgeProblem(o,
			fmt.Sprintf("Early stopping selects the epoch with the best validation metric. If that epoch is %d, which checkpoint is restored? (enter the epoch number)", b),
			strconv.Itoa(b),
			fmt.Sprintf("The checkpoint from epoch %d is restored because it had the best validation performance.", b))
	}
}
