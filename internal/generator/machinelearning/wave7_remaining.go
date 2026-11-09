package machinelearning

// Wave 7 machine learning, remaining batch. Answers are exact integers.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerWave7Remaining(reg *generator.Registry) {
	for id, kind := range map[string]string{
		"ml.reg.elastic_net":          "elastic_net",
		"ml.reg.polynomial":           "polynomial",
		"ml.svm.soft_margin":          "soft_margin",
		"ml.cluster.dbscan":           "dbscan",
		"ml.cluster.hierarchical":     "hierarchical",
		"ml.unsupervised.autoencoder": "autoencoder",
		"ml.optim.sgd":                "sgd",
		"ml.data.leakage":             "leakage",
	} {
		reg.Register(id, &w7MLRemainingGen{kind: kind})
	}
}

type w7MLRemainingGen struct{ kind string }

func (g *w7MLRemainingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp < 1 || kp > 3 {
		kp = 1 + rand.Intn(3)
	}
	switch g.kind {
	case "elastic_net":
		return w7MRElastic(kp)
	case "polynomial":
		return w7Polynomial(kp)
	case "soft_margin":
		return w7SoftMargin(kp)
	case "dbscan":
		return w7DBSCAN(kp)
	case "hierarchical":
		return w7Hierarchical(kp)
	case "autoencoder":
		return w7Autoencoder(kp)
	case "sgd":
		return w7SGD(kp)
	case "leakage":
		return w7Leakage(kp)
	}
	return generator.Problem{Question: "What is 2 + 2? (enter a number)", Answer: "4", Explanation: "2 + 2 = 4."}
}

func w7MRElastic(kp int) generator.Problem {
	l1 := w7MR(1, 5)
	l2 := w7MR(1, 5)
	w := w7MR(-3, 3)
	f := map[string]string{"l1": strconv.Itoa(l1), "l2": strconv.Itoa(l2), "w": strconv.Itoa(w)}
	switch kp {
	case 1:
		ans := l1*w7Abs(w) + l2*w*w
		return w7M(kp, fmt.Sprintf("An elastic-net model has \\(\\lambda_1=%d\\), \\(\\lambda_2=%d\\) and weight \\(w=%d\\). What is the combined penalty \\(\\lambda_1|w|+\\lambda_2w^2\\)?", l1, l2, w), strconv.Itoa(ans),
			"The L1 term is linear in |w| and the L2 term is quadratic.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("For \\(\\lambda_1=%d\\) and \\(\\lambda_2=%d\\), what is the mixing denominator \\(\\lambda_1+\\lambda_2\\)?", l1, l2), strconv.Itoa(l1+l2),
			"The mixing parameter is the ratio of one penalty to the sum.", f)
	default:
		c := w7MR(-5, 5)
		lambda := l1 + l2
		ans := c - lambda
		if c < 0 {
			ans = c + lambda
		}
		if w7Abs(c) <= lambda {
			ans = 0
		}
		f["c"] = strconv.Itoa(c)
		f["lambda"] = strconv.Itoa(lambda)
		return w7M(kp, fmt.Sprintf("For an elastic-net correlation \\(c=%d\\) and combined penalty \\(\\lambda=%d\\), what is the soft-thresholded coefficient?", c, lambda), strconv.Itoa(ans),
			"The combined penalty is the sum of the L1 and L2 penalties.", f)
	}
}

func w7Polynomial(kp int) generator.Problem {
	x := w7MR(-3, 3)
	d := w7MR(1, 12)
	w := w7MR(-3, 3)
	pow := 1
	for i := 0; i < d; i++ {
		pow *= x
	}
	f := map[string]string{"x": strconv.Itoa(x), "d": strconv.Itoa(d), "w": strconv.Itoa(w)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("What is the polynomial feature \\(x^d\\) for \\(x=%d\\) and \\(d=%d\\)?", x, d), strconv.Itoa(pow),
			fmt.Sprintf("\\(%d^%d=%d\\).", x, d, pow), f)
	case 2:
		return w7M(kp, fmt.Sprintf("A polynomial regression uses features up to degree %d. What is the degree of the model?", d), strconv.Itoa(d),
			"The highest feature degree is the model degree.", f)
	default:
		return w7M(kp, fmt.Sprintf("A polynomial model with coefficient \\(w=%d\\) predicts \\(w x^d\\) for \\(x=%d\\) and \\(d=%d\\). What is the prediction?", w, x, d), strconv.Itoa(w*pow),
			fmt.Sprintf("\\((%d)(%d)=%d\\).", w, pow, w*pow), f)
	}
}

func w7SoftMargin(kp int) generator.Problem {
	n := w7MR(1, 4)
	slacks := make([]int, 3)
	sum := 0
	for i := range slacks {
		slacks[i] = w7MR(0, 4)
		sum += slacks[i]
	}
	norms := []int{1, 2, 4, 5, 10}
	norm := norms[rand.Intn(len(norms))]
	C := w7MR(1, 12)
	f := map[string]string{"sum": strconv.Itoa(sum), "norm": strconv.Itoa(norm), "C": strconv.Itoa(C), "n": strconv.Itoa(n)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("A soft-margin SVM has slack values %v. What is the total slack?", slacks), strconv.Itoa(sum),
			"Add the individual slack variables.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("A soft-margin SVM has weight norm %d. What is 10 times the geometric margin \\(2/\\|w\\|\\)?", norm), strconv.Itoa(20/norm),
			fmt.Sprintf("\\(2/%d=%d\\).", norm, 2/norm), f)
	default:
		return w7M(kp, fmt.Sprintf("In a soft-margin SVM with C=%d, which C gives the smaller margin: the current C or a larger C? Enter 1 for current and 2 for larger.", C), "2",
			"Larger C penalises slack more and generally produces a smaller margin.", f)
	}
}

func w7DBSCAN(kp int) generator.Problem {
	n := w7MR(4, 8)
	labels := make([]int, n)
	clusters := w7MR(1, 3)
	noise := 0
	for i := range labels {
		if rand.Intn(4) == 0 {
			labels[i] = -1
			noise++
		} else {
			labels[i] = rand.Intn(clusters)
		}
	}
	f := map[string]string{"n": strconv.Itoa(n), "clusters": strconv.Itoa(clusters), "noise": strconv.Itoa(noise)}
	f["labels"] = fmt.Sprintf("%v", labels)
	switch kp {
	case 1:
		minPts := w7MR(2, 4)
		neighbors := w7MR(0, 5)
		is := 0
		if neighbors >= minPts {
			is = 1
		}
		f["minPts"] = strconv.Itoa(minPts)
		f["neighbors"] = strconv.Itoa(neighbors)
		return w7M(kp, fmt.Sprintf("A DBSCAN point has %d neighbours and minPts=%d. Is it a core point? Enter 1 for yes and 0 for no.", neighbors, minPts), strconv.Itoa(is),
			"A point is a core point when its neighbourhood size is at least minPts.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("A DBSCAN run labels %d points as %v. How many clusters are present, ignoring noise?", n, labels), strconv.Itoa(clusters),
			"Count the distinct nonnegative labels.", f)
	default:
		return w7M(kp, fmt.Sprintf("A DBSCAN run labels %d points as %v. How many noise points are there?", n, labels), strconv.Itoa(noise),
			"Noise points are labelled -1.", f)
	}
}

func w7Hierarchical(kp int) generator.Problem {
	d1 := w7MR(1, 4)
	d2 := d1 + w7MR(1, 4)
	d3 := d2 + w7MR(1, 4)
	f := map[string]string{"d1": strconv.Itoa(d1), "d2": strconv.Itoa(d2), "d3": strconv.Itoa(d3)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("A hierarchical clustering merges at distances %d, %d and %d. What is the first merge distance?", d1, d2, d3), strconv.Itoa(d1),
			"The first merge uses the smallest distance.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("For merge distances %d, %d and %d, what is the dendrogram height?", d1, d2, d3), strconv.Itoa(d3),
			"The tallest merge sets the dendrogram height.", f)
	default:
		return w7M(kp, fmt.Sprintf("If a dendrogram is cut at distance %d, how many clusters remain after merges at %d and %d?", d2, d1, d2), "2",
			"Cutting between the second and third merge leaves the two groups formed by the first two merges.", f)
	}
}

func w7Autoencoder(kp int) generator.Problem {
	n := w7MR(3, 6)
	b := w7MR(1, n-1)
	h := w7MR(2, 5)
	f := map[string]string{"n": strconv.Itoa(n), "b": strconv.Itoa(b), "h": strconv.Itoa(h)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("An autoencoder has input dimension %d and bottleneck dimension %d. What is the bottleneck dimension?", n, b), strconv.Itoa(b),
			"The bottleneck is the smallest hidden layer.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("An autoencoder has input dimension %d and hidden width %d. If the encoder has one weight matrix, how many weights does it have?", n, h), strconv.Itoa(n*h),
			"One dense layer has input times output weights.", f)
	default:
		errs := []int{w7MR(-2, 2), w7MR(-2, 2), w7MR(-2, 2)}
		sum := 0
		for _, e := range errs {
			sum += e * e
		}
		f["errs"] = fmt.Sprintf("%v", errs)
		return w7M(kp, fmt.Sprintf("An autoencoder has reconstruction errors %v. What is the squared reconstruction loss?", errs), strconv.Itoa(sum),
			"Square each error and add.", f)
	}
}

func w7SGD(kp int) generator.Problem {
	w := w7MR(-5, 5)
	g := w7MR(1, 12)
	k := w7MR(1, 6)
	for g%k != 0 {
		g = w7MR(1, 12)
	}
	f := map[string]string{"w": strconv.Itoa(w), "g": strconv.Itoa(g), "k": strconv.Itoa(k)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("An SGD update has weight %d, gradient %d and learning rate 1/%d. What is the updated weight \\(w-g/k\\)?", w, g, k), strconv.Itoa(w-g/k),
			fmt.Sprintf("\\(%d-%d/%d=%d\\).", w, g, k, w-g/k), f)
	case 2:
		g1 := w7MR(1, 6)
		g2 := w7MR(1, 6)
		if (g1+g2)%2 != 0 {
			g2++
		}
		f["g1"] = strconv.Itoa(g1)
		f["g2"] = strconv.Itoa(g2)
		return w7M(kp, fmt.Sprintf("A minibatch has gradients %d and %d. What is the minibatch gradient?", g1, g2), strconv.Itoa((g1+g2)/2),
			"Average the per-example gradients.", f)
	default:
		return w7M(kp, fmt.Sprintf("An SGD step uses learning rate 1/%d. What is the denominator of the learning rate?", k), strconv.Itoa(k),
			"The denominator is the number of equal steps.", f)
	}
}

func w7Leakage(kp int) generator.Problem {
	scenarios := []struct {
		text string
		leak int
		bad  int
	}{
		{"Standardising the full dataset before splitting", 1, 1},
		{"Using only training rows to fit a scaler", 0, 0},
		{"Adding the target column as a feature", 1, 1},
		{"Selecting features on the test set", 1, 1},
		{"Fitting preprocessing inside cross-validation", 0, 0},
	}
	s := scenarios[rand.Intn(len(scenarios))]
	f := map[string]string{"leak": strconv.Itoa(s.leak), "bad": strconv.Itoa(s.bad)}
	switch kp {
	case 1:
		return w7M(kp, fmt.Sprintf("Does the following pipeline leak test information: %s? Enter 1 for yes and 0 for no.", s.text), strconv.Itoa(s.leak),
			"Leakage occurs when information from outside the training fold affects training.", f)
	case 2:
		return w7M(kp, fmt.Sprintf("Does the pipeline \"%s\" contaminate a feature with test information? Enter 1 for yes and 0 for no.", s.text), strconv.Itoa(s.bad),
			"Contaminated features carry information that will not be available at prediction time.", f)
	default:
		return w7M(kp, fmt.Sprintf("Does the pipeline \"%s\" keep the train and test transformations separate? Enter 1 for yes and 0 for no.", s.text), strconv.Itoa(1-s.bad),
			"A clean pipeline fits every transformation on training data only.", f)
	}
}
