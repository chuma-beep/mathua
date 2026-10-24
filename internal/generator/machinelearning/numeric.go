package machinelearning

// The numeric machine-learning concepts: least-squares slope and prediction,
// accuracy, precision and recall from classification counts, the ridge penalty,
// explained variance percentage in principal component analysis and nearest-neighbour distance.
//
// Converted to parametric generators scaling with ctx.Difficulty using mathutil primitives
// and structured three-subgoal cognitive derivations in explanations.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mathutil"
)

func registerNumeric(reg *generator.Registry) {
	reg.Register("ml.linreg.least_squares", &leastSquaresGen{})
	reg.Register("ml.linreg.predict", &predictGen{})
	reg.Register("ml.eval.accuracy", &accuracyGen{})
	reg.Register("ml.eval.precision", &precisionGen{})
	reg.Register("ml.eval.recall", &recallGen{})
	reg.Register("ml.reg.ridge_penalty", &ridgePenaltyGen{})
	reg.Register("ml.pca.variance_percent", &variancePercentGen{})
	reg.Register("ml.knn.nearest_distance", &nearestDistanceGen{})
}

func smallInt(lo, hi int) int {
	if hi < lo {
		hi = lo
	}
	return lo + rand.Intn(hi-lo+1)
}

// leastSquaresGen fits a line to three non-collinear points using ordinary least squares.
type leastSquaresGen struct{}

func (g *leastSquaresGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxM := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	m := 1 + rand.Intn(maxM)
	maxB := mathutil.ScaleInt(ctx.Difficulty, 5, 20)
	b := rand.Intn(2*maxB+1) - maxB

	deltas := []int{-2, -1, 1, 2}
	delta := deltas[rand.Intn(len(deltas))]

	y0 := b
	y1 := m + b + delta
	y2 := 2*m + b

	// Ensure the slope is not read off from any coordinate. A printed coordinate
	// equal to the slope, or to its negative, repeats the answer in the stem: the
	// answer 10 next to a point (1, -10) is a leak even though no coordinate is
	// the slope itself.
	for b == 0 || y0 == m || y1 == m || y2 == m || y0 == -m || y1 == -m || y2 == -m {
		b = rand.Intn(2*maxB+1) - maxB
		delta = deltas[rand.Intn(len(deltas))]
		y0 = b
		y1 = m + b + delta
		y2 = 2*m + b
	}
	ans := strconv.Itoa(m)

	q := fmt.Sprintf("Fit a line to the points \\((0,%d)\\), \\((1,%d)\\) and \\((2,%d)\\). What is the least-squares slope?", y0, y1, y2)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The ordinary least squares slope for points (x_i, y_i) is m = Cov(x, y) / Var(x) = [Σ(x_i - x̄)(y_i - ȳ)] / [Σ(x_i - x̄)^2]. For x-coordinates (0, 1, 2), the mean is x̄ = 1 and variance sum is Σ(x_i - 1)^2 = (-1)^2 + 0^2 + 1^2 = 2. The covariance reduces to (y_2 - y_0) / 2.\n\nSubgoal 2: Substitute parameters. For coordinates (0, %d), (1, %d) and (2, %d), substitute endpoints y_0 = %d and y_2 = %d into the least-squares formula: m = (%d - (%d)) / 2.\n\nSubgoal 3: Evaluate and simplify. Compute difference %d - (%d) = %d. Dividing by 2 yields least-squares slope m = %d / 2 = %d.",
		y0, y1, y2, y0, y2, y2, y0, y2, y0, y2-y0, y2-y0, m)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"b":      strconv.Itoa(y0),
			"y1":     strconv.Itoa(y1),
			"y2":     strconv.Itoa(y2),
			"answer": ans,
		},
	}
}

// predictGen evaluates the fitted regression line at a given query point.
type predictGen struct{}

func (g *predictGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxM := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	m := 1 + rand.Intn(maxM)
	maxB := mathutil.ScaleInt(ctx.Difficulty, 5, 20)
	b := rand.Intn(2*maxB+1) - maxB
	maxX := mathutil.ScaleInt(ctx.Difficulty, 4, 15)
	x := 2 + rand.Intn(maxX)
	pred := m*x + b
	// A prediction equal to the slope or to the input can be read off the
	// question instead of computed, which happens for intercepts like b = -3 with
	// m = 2 and x = 3. Redraw the intercept until the prediction is its own
	// number.
	for pred == m || pred == x {
		b = rand.Intn(2*maxB+1) - maxB
		pred = m*x + b
	}
	ans := strconv.Itoa(pred)

	q := fmt.Sprintf("A model has been fitted as \\(\\hat{y} = %dx %+d\\). Find the predicted value of \\(\\hat{y}\\) when \\(x = %d\\).", m, b, x)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Prediction in a fitted univariate linear model ŷ = mx + b is evaluated by directly substituting query point x into the equation.\n\nSubgoal 2: Substitute parameters. Insert slope m = %d, intercept b = %d and query value x = %d: ŷ = %d(%d) %+d.\n\nSubgoal 3: Evaluate and simplify. Multiply and add: ŷ = %d %+d = %d.",
		m, b, x, m, x, b, m*x, b, pred)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"m":      strconv.Itoa(m),
			"b":      strconv.Itoa(b),
			"x":      strconv.Itoa(x),
			"answer": ans,
		},
	}
}

// accuracyGen evaluates classification accuracy as an exact integer percentage.
type accuracyGen struct{}

func (g *accuracyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	denoms := []int{10, 20, 25, 50}
	total := denoms[rand.Intn(len(denoms))]
	correct := 2 + rand.Intn(total-3)
	pct := 100 * correct / total

	tp := 1 + rand.Intn(correct-1)
	tn := correct - tp
	incorrect := total - correct
	fp := 1
	if incorrect > 1 {
		fp = 1 + rand.Intn(incorrect-1)
	}
	fn := incorrect - fp
	ans := strconv.Itoa(pct)

	q := fmt.Sprintf("A binary classifier has %d true positives, %d true negatives, %d false positives and %d false negatives. What is its accuracy, as a whole percentage?", tp, tn, fp, fn)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Classification accuracy equals the proportion of correct predictions over total instances: Accuracy = [(TP + TN) / (TP + TN + FP + FN)] · 100%%.\n\nSubgoal 2: Substitute parameters. Substitute TP = %d, TN = %d, FP = %d and FN = %d: Accuracy = [(%d + %d) / (%d + %d + %d + %d)] · 100%%.\n\nSubgoal 3: Evaluate and simplify. Total correct predictions equal %d + %d = %d across %d total samples. Divide to find accuracy: (%d / %d) · 100%% = %d%%.",
		tp, tn, fp, fn, tp, tn, tp, tn, fp, fn, tp, tn, correct, total, correct, total, pct)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"tp":     strconv.Itoa(tp),
			"tn":     strconv.Itoa(tn),
			"fp":     strconv.Itoa(fp),
			"fn":     strconv.Itoa(fn),
			"answer": ans,
		},
	}
}

// precisionGen evaluates precision as an exact whole percentage.
type precisionGen struct{}

func (g *precisionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	denoms := []int{10, 20, 25, 50}
	total := denoms[rand.Intn(len(denoms))]
	tp := 1 + rand.Intn(total-1)
	fp := total - tp
	pct := 100 * tp / total
	ans := strconv.Itoa(pct)

	q := fmt.Sprintf("A classifier produces %d true positives and %d false positives. What is its precision, as a whole percentage?", tp, fp)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Precision measures the proportion of positive predictions that are truly positive: Precision = [TP / (TP + FP)] · 100%%.\n\nSubgoal 2: Substitute parameters. Substitute true positives TP = %d and false positives FP = %d: Precision = [%d / (%d + %d)] · 100%%.\n\nSubgoal 3: Evaluate and simplify. Total predicted positives equal %d + %d = %d. Compute precision: (%d / %d) · 100%% = %d%%.",
		tp, fp, tp, tp, fp, tp, fp, total, tp, total, pct)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"tp":     strconv.Itoa(tp),
			"fp":     strconv.Itoa(fp),
			"answer": ans,
		},
	}
}

// recallGen evaluates recall as an exact whole percentage.
type recallGen struct{}

func (g *recallGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	denoms := []int{10, 20, 25, 50}
	total := denoms[rand.Intn(len(denoms))]
	tp := 1 + rand.Intn(total-1)
	fn := total - tp
	pct := 100 * tp / total
	ans := strconv.Itoa(pct)

	q := fmt.Sprintf("Of the positive cases in a test set, %d were found and %d were missed. What is the recall, as a whole percentage?", tp, fn)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. Recall measures the proportion of actual positive cases successfully detected: Recall = [TP / (TP + FN)] · 100%%.\n\nSubgoal 2: Substitute parameters. Substitute detected positives TP = %d and missed positives FN = %d: Recall = [%d / (%d + %d)] · 100%%.\n\nSubgoal 3: Evaluate and simplify. Total actual positive cases equal %d + %d = %d. Compute recall: (%d / %d) · 100%% = %d%%.",
		tp, fn, tp, tp, fn, tp, fn, total, tp, total, pct)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"tp":     strconv.Itoa(tp),
			"fn":     strconv.Itoa(fn),
			"answer": ans,
		},
	}
}

// ridgePenaltyGen computes the L2 penalty term lambda * (w1^2 + w2^2).
type ridgePenaltyGen struct{}

func (g *ridgePenaltyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	maxLambda := mathutil.ScaleInt(ctx.Difficulty, 3, 10)
	lambda := 1 + rand.Intn(maxLambda)
	maxW := mathutil.ScaleInt(ctx.Difficulty, 3, 8)
	w1 := 1 + rand.Intn(maxW)
	w2 := 1 + rand.Intn(maxW)
	penalty := lambda * (w1*w1 + w2*w2)
	ans := strconv.Itoa(penalty)

	q := fmt.Sprintf("A ridge regression has weights %d and %d and penalty strength %d. What is the value of the L2 penalty term \\(\\lambda(w_1^2 + w_2^2)\\)?", w1, w2, lambda)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The L2 ridge regularization penalty term is defined as R(w) = λ · (w_1^2 + w_2^2).\n\nSubgoal 2: Substitute parameters. Substitute penalty strength λ = %d and weights w_1 = %d and w_2 = %d: R(w) = %d · ((%d)^2 + (%d)^2).\n\nSubgoal 3: Evaluate and simplify. Compute squared weights (%d)^2 + (%d)^2 = %d + %d = %d. Multiply by penalty strength: R(w) = %d · %d = %d.",
		lambda, w1, w2, lambda, w1, w2, w1, w2, w1*w1, w2*w2, w1*w1+w2*w2, lambda, w1*w1+w2*w2, penalty)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"w1":     strconv.Itoa(w1),
			"w2":     strconv.Itoa(w2),
			"lambda": strconv.Itoa(lambda),
			"answer": ans,
		},
	}
}

// variancePercentGen computes PCA explained variance percentage.
type variancePercentGen struct{}

func (g *variancePercentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// Draw the answer first, then build eigenvalues that produce it exactly. A
	// fixed pool of seven ratios made this concept serve seven answers however
	// often a learner met it.
	pct := 10 + rand.Intn(81)
	k := 2 + rand.Intn(mathutil.ScaleInt(ctx.Difficulty, 1, 4))
	// λ2 = k(100 - pct); it equals pct at pct = 80, k = 4, which would print the
	// answer as an eigenvalue.
	for k*(100-pct) == pct {
		k++
	}
	lam1 := k * pct
	lam2 := k * (100 - pct)
	ans := strconv.Itoa(pct)

	q := fmt.Sprintf("A principal component analysis finds two eigenvalues, %d and %d. What percentage of the total variance is accounted for by the first component?", lam1, lam2)
	exp := fmt.Sprintf("Subgoal 1: State governing formula. In principal component analysis, the percentage of total variance explained by the first component equals Explained Variance = [λ_1 / (λ_1 + λ_2)] · 100%%.\n\nSubgoal 2: Substitute parameters. Insert eigenvalues λ_1 = %d and λ_2 = %d: Explained Variance = [%d / (%d + %d)] · 100%%.\n\nSubgoal 3: Evaluate and simplify. Total variance is %d + %d = %d. Compute percentage: (%d / %d) · 100%% = %d%%.",
		lam1, lam2, lam1, lam1, lam2, lam1, lam2, lam1+lam2, lam1, lam1+lam2, pct)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"lam1":   strconv.Itoa(lam1),
			"lam2":   strconv.Itoa(lam2),
			"answer": ans,
		},
	}
}

// nearestDistanceGen evaluates squared Euclidean distance to the nearest training point.
type nearestDistanceGen struct{}

func (g *nearestDistanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	best := 1 << 30
	var bx, by int
	pts := [][2]int{}
	dists := []int{}
	for i := 0; i < 3; i++ {
		x := smallInt(-6, 6)
		y := smallInt(-6, 6)
		if x == 0 && y == 0 {
			x = 1
		}
		pts = append(pts, [2]int{x, y})
		d := x*x + y*y
		dists = append(dists, d)
		if d < best {
			best, bx, by = d, x, y
		}
	}
	ans := strconv.Itoa(best)

	q := fmt.Sprintf("A nearest-neighbour classifier holds the training points \\((%d,%d)\\), \\((%d,%d)\\) and \\((%d,%d)\\). What is the squared distance from the origin to the nearest of them?",
		pts[0][0], pts[0][1], pts[1][0], pts[1][1], pts[2][0], pts[2][1])
	exp := fmt.Sprintf("Subgoal 1: State governing formula. The squared Euclidean distance from the origin (0, 0) to point (x, y) is d^2 = x^2 + y^2. A nearest-neighbour classifier selects the point minimizing this squared distance.\n\nSubgoal 2: Substitute parameters. Compute squared distance to each candidate point: d_1^2 = (%d)^2 + (%d)^2 = %d, d_2^2 = (%d)^2 + (%d)^2 = %d and d_3^2 = (%d)^2 + (%d)^2 = %d.\n\nSubgoal 3: Evaluate and simplify. Comparing squared distances min(%d, %d, %d) yields the minimum value %d at nearest point (%d, %d).",
		pts[0][0], pts[0][1], dists[0], pts[1][0], pts[1][1], dists[1], pts[2][0], pts[2][1], dists[2], dists[0], dists[1], dists[2], best, bx, by)

	return generator.Problem{
		Question:    q,
		Answer:      ans,
		Explanation: exp,
		Facts: map[string]string{
			"answer": ans,
		},
	}
}
