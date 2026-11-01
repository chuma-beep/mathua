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
	"fmt"
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
	if rand.Intn(5) == 0 {
		return generator.Problem{
			Question:    "The sigmoid is sigma(x) = 1/(1+e^-x). What is the limit of sigma(x) as x tends to infinity? (enter a number)",
			Answer:      "1",
			Explanation: "As x grows, e^-x tends to 0, so the fraction tends to 1.",
			Facts: map[string]string{
				"kind":   "limit",
				"x":      "0",
				"answer": "1",
			},
		}
	}

	cleanK := []int{1, 3, 4, 7, 9, 15, 19, 24, 39, 49, 99}
	k := cleanK[rand.Intn(len(cleanK))]

	if k == 1 {
		return generator.Problem{
			Question:    "The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(0)? (enter a decimal)",
			Answer:      "0.5",
			Explanation: "e^0 = 1, so sigma(0) = 1/(1+1) = 0.5.",
			Facts: map[string]string{
				"kind":   "value",
				"x":      "0",
				"answer": "0.5",
			},
		}
	}

	isPositive := rand.Intn(2) == 1
	var xVal, val float64
	var q, ex string
	if isPositive {
		xVal = math.Log(float64(k))
		val = float64(k) / float64(k+1)
		ansStr := strconv.FormatFloat(val, 'g', -1, 64)
		q = fmt.Sprintf("The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(ln %d)? (enter a decimal)", k)
		ex = fmt.Sprintf("e^(-ln %d) = 1/%d, so sigma(ln %d) = 1/(1+1/%d) = %d/%d = %s.", k, k, k, k, k, k+1, ansStr)
		return generator.Problem{
			Question:    q,
			Answer:      ansStr,
			Explanation: ex,
			Facts: map[string]string{
				"kind":   "value",
				"x":      strconv.FormatFloat(xVal, 'g', -1, 64),
				"answer": ansStr,
			},
		}
	}

	xVal = -math.Log(float64(k))
	val = 1.0 / float64(1+k)
	ansStr := strconv.FormatFloat(val, 'g', -1, 64)
	q = fmt.Sprintf("The sigmoid is sigma(x) = 1/(1+e^-x). What is sigma(-ln %d)? (enter a decimal)", k)
	ex = fmt.Sprintf("e^(ln %d) = %d, so sigma(-ln %d) = 1/(1+%d) = 1/%d = %s.", k, k, k, k, k+1, ansStr)
	return generator.Problem{
		Question:    q,
		Answer:      ansStr,
		Explanation: ex,
		Facts: map[string]string{
			"kind":   "value",
			"x":      strconv.FormatFloat(xVal, 'g', -1, 64),
			"answer": ansStr,
		},
	}
}

// ---------------------------------------------------------------------------
// 2. Logistic regression
// ---------------------------------------------------------------------------

type logisticRegressionGen struct{}

func (g *logisticRegressionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return logisticRegressionProblem()
}

// ---------------------------------------------------------------------------
// 3. Vanishing and exploding gradients
// ---------------------------------------------------------------------------

type vanishingGradientGen struct{}

func (g *vanishingGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return vanishingGradientProblem()
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
	return imbalanceProblem()
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
	return ensembleProblem()
}
