package machinelearning

// The 2026-10-18 machine-learning batch: twelve concepts that take the domain
// from 33 to the target of 45.
//
// The survey found that the domain's missing pieces are not more named
// architectures: the CNN lesson already teaches output-size arithmetic, the
// validation lesson already teaches leakage, the gradient-descent lesson
// already teaches the learning-rate trade, the overfitting lesson already
// teaches capacity, and the transformer lesson already teaches softmax
// attention. The batch therefore fills the gaps those lessons point at but do
// not own: the sigmoid and log loss family behind logistic regression, the
// absolute-error counterpart to MSE, the confusion matrix and F1 that the
// accuracy/precision/recall lessons name but never build, feature scaling and
// one-hot encoding, class imbalance, k-means, and the decision-tree family
// (Gini impurity and bagging).
//
// Seven of the twelve are numeric, which raises the domain from 24 of 33
// multiple-choice to 29 of 45. Every numeric generator publishes the facts its
// answer follows from, so the tests recompute the answer from those facts
// instead of comparing a string with itself. The multiple-choice generators
// use generator.MultipleChoice, so every distractor is a real option the router
// has to reject.

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

func registerBatch(reg *generator.Registry) {
	reg.Register("ml.activation.sigmoid", &sigmoidGen{})
	reg.Register("ml.logistic.regression", &logisticRegressionGen{})
	reg.Register("ml.grad.vanishing", &vanishingGradientGen{})
	reg.Register("ml.loss.mae", &maeGen{})
	reg.Register("ml.eval.confusion", &confusionMatrixGen{})
	reg.Register("ml.eval.f1", &f1Gen{})
	reg.Register("ml.cluster.kmeans", &kmeansGen{})
	reg.Register("ml.data.scaling", &scalingGen{})
	reg.Register("ml.data.imbalance", &imbalanceGen{})
	reg.Register("ml.tree.gini", &giniGen{})
	reg.Register("ml.tree.ensemble", &ensembleGen{})
}

// ---------------------------------------------------------------------------
// 1. Sigmoid activation
// ---------------------------------------------------------------------------

type sigmoidGen struct{}

func (g *sigmoidGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		x        float64
		kind     string
		reason   string
	}
	claims := []claim{
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(0)? (enter a decimal)",
			answer:   "0.5", x: 0, kind: "value",
			reason: "e^0 = 1, so sigma(0) = 1/(1+1) = 0.5.",
		},
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(ln 3)? (enter a decimal)",
			answer:   "0.75", x: math.Log(3), kind: "value",
			reason: "e^(-ln 3) = 1/3, so sigma(ln 3) = 1/(1+1/3) = 3/4 = 0.75.",
		},
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(-ln 3)? (enter a decimal)",
			answer:   "0.25", x: -math.Log(3), kind: "value",
			reason: "e^(ln 3) = 3, so sigma(-ln 3) = 1/(1+3) = 1/4 = 0.25.",
		},
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(ln 9)? (enter a decimal)",
			answer:   "0.9", x: math.Log(9), kind: "value",
			reason: "e^(-ln 9) = 1/9, so sigma(ln 9) = 1/(1+1/9) = 9/10 = 0.9.",
		},
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(-ln 9)? (enter a decimal)",
			answer:   "0.1", x: -math.Log(9), kind: "value",
			reason: "e^(ln 9) = 9, so sigma(-ln 9) = 1/(1+9) = 1/10 = 0.1.",
		},
		{
			question: "The sigmoid is sigma(x) = 1/(1+e^-x). What is the limit of sigma(x) as x tends to infinity? (enter a number)",
			answer:   "1", x: 0, kind: "limit",
			reason: "As x grows, e^-x tends to 0, so the fraction tends to 1.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"kind": c.kind, "x": strconv.FormatFloat(c.x, 'g', -1, 64), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 2. Logistic regression
// ---------------------------------------------------------------------------

type logisticRegressionGen struct{}

func (g *logisticRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Logistic regression turns a linear score into...",
			correct:     "a probability through the sigmoid",
			distractors: []string{"a real number through the identity", "a class through a step function", "a distance through the absolute value"},
			reason:      "The linear score w*x + b is passed through the sigmoid, which maps every real number into the interval from 0 to 1.",
		},
		{
			question:    "Which loss is the natural training objective for logistic regression?",
			correct:     "cross-entropy",
			distractors: []string{"mean squared error", "mean absolute error", "hinge loss"},
			reason:      "Cross-entropy is the negative log-likelihood of the Bernoulli label, so it matches the probabilistic output of the model.",
		},
		{
			question:    "The decision boundary of logistic regression is...",
			correct:     "a hyperplane where the linear score is zero",
			distractors: []string{"a circle centred on the mean", "a quadratic curve", "the convex hull of the positive points"},
			reason:      "The model predicts above one half exactly when the linear score is positive, and the score is zero on a hyperplane.",
		},
		{
			question:    "If the predicted probability is 0.9, how is the point classified at a 0.5 threshold?",
			correct:     "positive",
			distractors: []string{"negative", "undecided", "it depends on the sign of the bias"},
			reason:      "A probability above 0.5 falls on the positive side of the threshold, regardless of the individual weights.",
		},
		{
			question:    "Logistic regression is called a linear classifier because...",
			correct:     "its decision boundary is linear in the inputs",
			distractors: []string{"its loss is linear in the parameters", "it can use only one feature", "its output is a linear function"},
			reason:      "The probability is nonlinear, but the boundary where the probability equals one half is the linear equation w*x + b = 0.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 3. Vanishing and exploding gradients
// ---------------------------------------------------------------------------

type vanishingGradientGen struct{}

func (g *vanishingGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Why can deep sigmoid networks suffer from vanishing gradients?",
			correct:     "the sigmoid derivative is at most 0.25, so repeated multiplication shrinks it",
			distractors: []string{"the sigmoid derivative is above 1, so gradients explode", "the sigmoid is not differentiable", "the loss is unbounded"},
			reason:      "The sigmoid derivative peaks at 0.25, so a product of many layer derivatives shrinks toward zero as depth grows.",
		},
		{
			question:    "Exploding gradients are most often caused by...",
			correct:     "repeated multiplication by factors larger than 1",
			distractors: []string{"a learning rate of zero", "too few layers", "using cross-entropy"},
			reason:      "If the layer derivatives or weights have magnitude above 1, their product grows exponentially with depth.",
		},
		{
			question:    "Which remedy targets vanishing gradients?",
			correct:     "ReLU activations and residual connections",
			distractors: []string{"tiny random weight initialization", "removing normalization", "increasing only the batch size"},
			reason:      "ReLU keeps a derivative of 1 on the positive side, and residual connections give gradients a short path back.",
		},
		{
			question:    "In a deep network, the gradient at an early layer is...",
			correct:     "a product of many layer derivatives",
			distractors: []string{"the sum of the layer outputs", "independent of depth", "always equal to the loss"},
			reason:      "The chain rule multiplies one local derivative per layer, so depth appears as a product and drives vanishing or explosion.",
		},
		{
			question:    "Gradient clipping is used to...",
			correct:     "cap the gradient norm and stop an exploding update",
			distractors: []string{"raise the learning rate", "remove the bias terms", "make the loss convex"},
			reason:      "Clipping bounds the size of the update when a rare batch produces an unusually large gradient.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 4. Mean absolute error (numeric)
// ---------------------------------------------------------------------------

type maeGen struct{}

func (g *maeGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question  string
		answer    string
		residuals []int
		reason    string
	}
	claims := []claim{
		{
			question: "A model has residuals 1, -2, 3 and 0. What is the mean absolute error? (enter a decimal)",
			answer:   "1.5", residuals: []int{1, -2, 3, 0},
			reason: "The absolute residuals are 1, 2, 3 and 0, whose sum is 6; dividing by 4 gives 1.5.",
		},
		{
			question: "A model has residuals 2, -2, 4 and -4. What is the mean absolute error? (enter a number)",
			answer:   "3", residuals: []int{2, -2, 4, -4},
			reason: "The absolute residuals are 2, 2, 4 and 4, whose sum is 12; dividing by 4 gives 3.",
		},
		{
			question: "Predictions are 3 and 5 and the targets are 1 and 5. What is the mean absolute error? (enter a number)",
			answer:   "1", residuals: []int{2, 0},
			reason: "The residuals are 2 and 0, so the mean absolute error is (2 + 0)/2 = 1.",
		},
		{
			question: "A model has residuals -1, 3 and 2. What is the mean absolute error? (enter a number)",
			answer:   "2", residuals: []int{-1, 3, 2},
			reason: "The absolute residuals are 1, 3 and 2, whose sum is 6; dividing by 3 gives 2.",
		},
		{
			question: "A model has residuals 0, 0, 0 and 8. What is the mean absolute error? (enter a number)",
			answer:   "2", residuals: []int{0, 0, 0, 8},
			reason: "The absolute residuals sum to 8, and dividing by 4 gives 2.",
		},
		{
			question: "A model has residuals 5 and -5. What is the mean absolute error? (enter a number)",
			answer:   "5", residuals: []int{5, -5},
			reason: "Both absolute residuals are 5, so their mean is 5.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	parts := make([]string, len(c.residuals))
	for i, r := range c.residuals {
		parts[i] = strconv.Itoa(r)
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"residuals": strings.Join(parts, ","), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 5. The confusion matrix (numeric)
// ---------------------------------------------------------------------------

type confusionMatrixGen struct{}

func (g *confusionMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		tp, fp   int
		fn, tn   int
		ask      string
		reason   string
	}
	claims := []claim{
		{
			question: "A classifier gives TP = 8, FP = 2, FN = 3 and TN = 7. How many cases are there in total? (enter a number)",
			answer:   "20", tp: 8, fp: 2, fn: 3, tn: 7, ask: "total",
			reason: "Every case is one of the four cells, so the total is 8 + 2 + 3 + 7 = 20.",
		},
		{
			question: "A classifier gives TP = 12, FP = 4, FN = 6 and TN = 18. How many actual positives are there? (enter a number)",
			answer:   "18", tp: 12, fp: 4, fn: 6, tn: 18, ask: "actual_pos",
			reason: "Actual positives are the true positives plus the false negatives: 12 + 6 = 18.",
		},
		{
			question: "A classifier gives TP = 5, FP = 1, FN = 5 and TN = 9. How many cases did it predict positive? (enter a number)",
			answer:   "6", tp: 5, fp: 1, fn: 5, tn: 9, ask: "predicted_pos",
			reason: "Predicted positives are the true positives plus the false positives: 5 + 1 = 6.",
		},
		{
			question: "A classifier gives TP = 20, FP = 5, FN = 0 and TN = 15. How many predictions were correct? (enter a number)",
			answer:   "35", tp: 20, fp: 5, fn: 0, tn: 15, ask: "correct",
			reason: "Correct predictions are the true positives plus the true negatives: 20 + 15 = 35.",
		},
		{
			question: "A classifier gives TP = 7, FP = 3, FN = 2 and TN = 8. How many actual negatives are there? (enter a number)",
			answer:   "11", tp: 7, fp: 3, fn: 2, tn: 8, ask: "actual_neg",
			reason: "Actual negatives are the false positives plus the true negatives: 3 + 8 = 11.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"tp": strconv.Itoa(c.tp), "fp": strconv.Itoa(c.fp),
			"fn": strconv.Itoa(c.fn), "tn": strconv.Itoa(c.tn),
			"ask": c.ask, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 6. The F1 score (numeric)
// ---------------------------------------------------------------------------

type f1Gen struct{}

func (g *f1Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question  string
		answer    string
		precision float64
		recall    float64
		reason    string
	}
	claims := []claim{
		{
			question: "A classifier has precision 1 and recall 0.6. What is its F1 score? (enter a decimal)",
			answer:   "0.75", precision: 1, recall: 0.6,
			reason: "F1 = 2*1*0.6/(1+0.6) = 1.2/1.6 = 0.75.",
		},
		{
			question: "A classifier has precision 0.6 and recall 0.4. What is its F1 score? (enter a decimal)",
			answer:   "0.48", precision: 0.6, recall: 0.4,
			reason: "F1 = 2*0.6*0.4/(0.6+0.4) = 0.48/1.0 = 0.48.",
		},
		{
			question: "A classifier has precision 0.6 and recall 0.9. What is its F1 score? (enter a decimal)",
			answer:   "0.72", precision: 0.6, recall: 0.9,
			reason: "F1 = 2*0.6*0.9/(0.6+0.9) = 1.08/1.5 = 0.72.",
		},
		{
			question: "A classifier has precision 0.5 and recall 0.75. What is its F1 score? (enter a decimal)",
			answer:   "0.6", precision: 0.5, recall: 0.75,
			reason: "F1 = 2*0.5*0.75/(0.5+0.75) = 0.75/1.25 = 0.6.",
		},
		{
			question: "A classifier has precision 0.25 and recall 1. What is its F1 score? (enter a decimal)",
			answer:   "0.4", precision: 0.25, recall: 1,
			reason: "F1 = 2*0.25*1/(0.25+1) = 0.5/1.25 = 0.4.",
		},
		{
			question: "A classifier has precision 0.3 and recall 0.7. What is its F1 score? (enter a decimal)",
			answer:   "0.42", precision: 0.3, recall: 0.7,
			reason: "F1 = 2*0.3*0.7/(0.3+0.7) = 0.42/1.0 = 0.42.",
		},
		{
			question: "A classifier has precision 0.2 and recall 0.8. What is its F1 score? (enter a decimal)",
			answer:   "0.32", precision: 0.2, recall: 0.8,
			reason: "F1 = 2*0.2*0.8/(0.2+0.8) = 0.32/1.0 = 0.32.",
		},
		{
			question: "A classifier has precision 0.1 and recall 0.9. What is its F1 score? (enter a decimal)",
			answer:   "0.18", precision: 0.1, recall: 0.9,
			reason: "F1 = 2*0.1*0.9/(0.1+0.9) = 0.18/1.0 = 0.18.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"precision": strconv.FormatFloat(c.precision, 'g', -1, 64),
			"recall":    strconv.FormatFloat(c.recall, 'g', -1, 64),
			"answer":    c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 7. k-means clustering (numeric)
// ---------------------------------------------------------------------------

type kmeansGen struct{}

func (g *kmeansGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		points   [][2]int
		axis     string
		reason   string
	}
	claims := []claim{
		{
			question: "The points (0, 0) and (2, 4) are assigned to one k-means cluster. What is the x-coordinate of the new centroid? (enter a number)",
			answer:   "1", points: [][2]int{{0, 0}, {2, 4}}, axis: "x",
			reason: "The centroid averages each coordinate: (0 + 2)/2 = 1 for x.",
		},
		{
			question: "The points (0, 0) and (2, 4) are assigned to one k-means cluster. What is the y-coordinate of the new centroid? (enter a number)",
			answer:   "2", points: [][2]int{{0, 0}, {2, 4}}, axis: "y",
			reason: "The centroid averages each coordinate: (0 + 4)/2 = 2 for y.",
		},
		{
			question: "The points (1, 1), (3, 5) and (5, 3) are assigned to one k-means cluster. What is the x-coordinate of the new centroid? (enter a number)",
			answer:   "3", points: [][2]int{{1, 1}, {3, 5}, {5, 3}}, axis: "x",
			reason: "The x-coordinates 1, 3 and 5 average to (1 + 3 + 5)/3 = 3.",
		},
		{
			question: "The points (1, 1), (3, 5) and (5, 3) are assigned to one k-means cluster. What is the y-coordinate of the new centroid? (enter a number)",
			answer:   "3", points: [][2]int{{1, 1}, {3, 5}, {5, 3}}, axis: "y",
			reason: "The y-coordinates 1, 5 and 3 average to (1 + 5 + 3)/3 = 3.",
		},
		{
			question: "The points (2, 2) and (4, 6) are assigned to one k-means cluster. What is the y-coordinate of the new centroid? (enter a number)",
			answer:   "4", points: [][2]int{{2, 2}, {4, 6}}, axis: "y",
			reason: "The y-coordinates 2 and 6 average to (2 + 6)/2 = 4.",
		},
		{
			question: "The points (0, 4), (2, 0) and (4, 2) are assigned to one k-means cluster. What is the x-coordinate of the new centroid? (enter a number)",
			answer:   "2", points: [][2]int{{0, 4}, {2, 0}, {4, 2}}, axis: "x",
			reason: "The x-coordinates 0, 2 and 4 average to (0 + 2 + 4)/3 = 2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	parts := make([]string, len(c.points))
	for i, p := range c.points {
		parts[i] = strconv.Itoa(p[0]) + "," + strconv.Itoa(p[1])
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"points": strings.Join(parts, ";"), "axis": c.axis, "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 8. Feature scaling (numeric)
// ---------------------------------------------------------------------------

type scalingGen struct{}

func (g *scalingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		x, mean  int
		sd       int
		reason   string
	}
	claims := []claim{
		{
			question: "A feature value is 10, the training mean is 4 and the training standard deviation is 2. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "3", x: 10, mean: 4, sd: 2,
			reason: "(10 - 4)/2 = 6/2 = 3.",
		},
		{
			question: "A feature value is 5, the training mean is 5 and the training standard deviation is 2. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "0", x: 5, mean: 5, sd: 2,
			reason: "A value at the mean has z-score (5 - 5)/2 = 0.",
		},
		{
			question: "A feature value is 8, the training mean is 2 and the training standard deviation is 3. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "2", x: 8, mean: 2, sd: 3,
			reason: "(8 - 2)/3 = 6/3 = 2.",
		},
		{
			question: "A feature value is 1, the training mean is 5 and the training standard deviation is 2. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "-2", x: 1, mean: 5, sd: 2,
			reason: "(1 - 5)/2 = -4/2 = -2.",
		},
		{
			question: "A feature value is 12, the training mean is 10 and the training standard deviation is 1. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "2", x: 12, mean: 10, sd: 1,
			reason: "(12 - 10)/1 = 2.",
		},
		{
			question: "A feature value is 0, the training mean is 10 and the training standard deviation is 5. What is the standardized value (x - mean) / sd? (enter a number)",
			answer:   "-2", x: 0, mean: 10, sd: 5,
			reason: "(0 - 10)/5 = -2.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"x": strconv.Itoa(c.x), "mean": strconv.Itoa(c.mean), "sd": strconv.Itoa(c.sd),
			"answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 9. One-hot encoding
// ---------------------------------------------------------------------------

type oneHotGen struct{}

func (g *oneHotGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A categorical feature has 4 categories and one-hot encoding keeps every column. How many columns result?",
			correct:     "4",
			distractors: []string{"3", "1", "2"},
			reason:      "One column per category is created, so four categories give four columns.",
		},
		{
			question:    "One-hot encoding represents a category as...",
			correct:     "a vector with a single 1 and zeros elsewhere",
			distractors: []string{"its integer rank", "the mean of the labels", "a vector of all ones"},
			reason:      "The active category gets a 1 and every other column gets a 0, so no ordering is implied.",
		},
		{
			question:    "Why drop one column in one-hot encoding before a linear model with an intercept?",
			correct:     "to avoid perfect collinearity",
			distractors: []string{"to reduce the number of categories", "to make the values continuous", "to remove missing values"},
			reason:      "The dropped column is the sum of the others, so keeping it makes the design matrix rank deficient against the intercept.",
		},
		{
			question:    "One-hot encoding is intended for which input?",
			correct:     "a categorical feature",
			distractors: []string{"a continuous feature", "a regression target", "a missing value"},
			reason:      "The encoding exists because categories have no numeric order, unlike a continuous feature.",
		},
		{
			question:    "For three categories, what does the one-hot vector [0, 1, 0] represent?",
			correct:     "the second category",
			distractors: []string{"the first category", "the third category", "all three categories"},
			reason:      "The 1 sits in the second position, so the vector selects the second category.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 10. Class imbalance
// ---------------------------------------------------------------------------

type imbalanceGen struct{}

func (g *imbalanceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "A dataset has 990 negatives and 10 positives. A model that always predicts negative has what accuracy?",
			correct:     "99%",
			distractors: []string{"50%", "10%", "1%"},
			reason:      "It is right on the 990 negatives and wrong on the 10 positives, so its accuracy is 990/1000 = 99%.",
		},
		{
			question:    "Which technique directly addresses class imbalance?",
			correct:     "resampling or class weights",
			distractors: []string{"one-hot encoding", "a larger learning rate", "removing the minority class"},
			reason:      "Resampling changes the class frequencies and class weights change their cost; neither discards the minority class.",
		},
		{
			question:    "Why is accuracy misleading on imbalanced data?",
			correct:     "the majority class dominates the count",
			distractors: []string{"the loss is always convex", "the model cannot output probabilities", "the test set is too large"},
			reason:      "A model can score high by always predicting the majority class, so the errors that matter disappear into the average.",
		},
		{
			question:    "Undersampling the majority class...",
			correct:     "discards data and can raise variance",
			distractors: []string{"increases the majority count", "guarantees better precision", "removes the need for validation"},
			reason:      "Throwing away majority examples balances the classes but spends information, so estimates become noisier.",
		},
		{
			question:    "When positives are rare, which metrics are more informative than accuracy?",
			correct:     "precision and recall",
			distractors: []string{"mean squared error", "the parameter count", "training time"},
			reason:      "Precision and recall report the two error types separately, which accuracy hides when one class dominates.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}

// ---------------------------------------------------------------------------
// 11. Gini impurity (numeric)
// ---------------------------------------------------------------------------

type giniGen struct{}

func (g *giniGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question string
		answer   string
		counts   []int
		reason   string
	}
	claims := []claim{
		{
			question: "A decision-tree node holds 1 example of one class and 1 of another. What is its Gini impurity, 1 - sum(p_i^2)? (enter a decimal)",
			answer:   "0.5", counts: []int{1, 1},
			reason: "Each class has probability 1/2, so the impurity is 1 - (1/4 + 1/4) = 0.5.",
		},
		{
			question: "A decision-tree node holds 3 examples of one class and 1 of another. What is its Gini impurity, 1 - sum(p_i^2)? (enter a decimal)",
			answer:   "0.375", counts: []int{3, 1},
			reason: "The probabilities are 3/4 and 1/4, so the impurity is 1 - (9/16 + 1/16) = 6/16 = 0.375.",
		},
		{
			question: "A decision-tree node is pure: all 4 examples have the same class. What is its Gini impurity, 1 - sum(p_i^2)? (enter a number)",
			answer:   "0", counts: []int{4, 0},
			reason: "One class has probability 1 and the rest 0, so the impurity is 1 - 1 = 0.",
		},
		{
			question: "A decision-tree node holds 2 examples of one class, 1 of a second class and 1 of a third. What is its Gini impurity, 1 - sum(p_i^2)? (enter a decimal)",
			answer:   "0.625", counts: []int{2, 1, 1},
			reason: "The probabilities are 1/2, 1/4 and 1/4, so the impurity is 1 - (1/4 + 1/16 + 1/16) = 10/16 = 0.625.",
		},
		{
			question: "A decision-tree node holds 3 examples of one class and 3 of another. What is its Gini impurity, 1 - sum(p_i^2)? (enter a decimal)",
			answer:   "0.5", counts: []int{3, 3},
			reason: "Each class has probability 1/2, so the impurity is 1 - (1/4 + 1/4) = 0.5.",
		},
		{
			question: "A decision-tree node holds 2 examples of one class and 2 of another. What is its Gini impurity, 1 - sum(p_i^2)? (enter a decimal)",
			answer:   "0.5", counts: []int{2, 2},
			reason: "Each class has probability 1/2, so the impurity is 1 - (1/4 + 1/4) = 0.5.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	parts := make([]string, len(c.counts))
	for i, n := range c.counts {
		parts[i] = strconv.Itoa(n)
	}
	return generator.Problem{
		Question:    c.question,
		Answer:      c.answer,
		Explanation: c.reason,
		Facts: map[string]string{
			"counts": strings.Join(parts, ","), "answer": c.answer,
		},
	}
}

// ---------------------------------------------------------------------------
// 12. Bagging and random forests
// ---------------------------------------------------------------------------

type ensembleGen struct{}

func (g *ensembleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type claim struct {
		question    string
		correct     string
		distractors []string
		reason      string
	}
	claims := []claim{
		{
			question:    "Bagging trains many models on...",
			correct:     "bootstrap samples of the training data",
			distractors: []string{"the same data with shuffled labels", "successively reweighted copies", "the test set"},
			reason:      "Each model sees a resample drawn with replacement, so the models differ while staying trained on real data.",
		},
		{
			question:    "Random forests add which ingredient to bagging?",
			correct:     "a random subset of features at each split",
			distractors: []string{"a larger learning rate", "one-hot encoding", "a single deep tree"},
			reason:      "Restricting the candidate features decorrelates the trees, which makes the average more reliable.",
		},
		{
			question:    "Bagging mainly reduces...",
			correct:     "variance",
			distractors: []string{"bias", "the number of features", "the training set size"},
			reason:      "Averaging many high-variance models cancels part of their independent error, but each model keeps its own bias.",
		},
		{
			question:    "Boosting mainly reduces...",
			correct:     "bias",
			distractors: []string{"variance only", "the number of trees", "the need for validation"},
			reason:      "Each new weak learner corrects the residual errors of the current ensemble, so the fit improves where it was too simple.",
		},
		{
			question:    "Averaging many high-variance trees helps because...",
			correct:     "their errors are partly independent and cancel",
			distractors: []string{"each tree has zero bias", "the average is always unbiased", "the trees see the test data"},
			reason:      "The variance of an average shrinks when the components are not perfectly correlated, which is why the feature subsampling matters.",
		},
	}
	c := claims[rand.Intn(len(claims))]
	p := generator.MultipleChoice(c.question, c.correct, c.distractors...)
	p.Explanation = c.reason
	return p
}
