package machinelearning

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// mlChoice builds a multiple choice problem with an explanation. Two of these
// concepts are graded multiple_choice, and the other two are numeric, so every
// classification shape here answers in options rather than in prose.

func mlChoice(question, answer, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, answer, distractors...)
	p.Explanation = explanation
	return p
}

type confusionMatrixGen struct{}

func (g *confusionMatrixGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	tp := rand.Intn(18) + 2
	fp := rand.Intn(9) + 1
	fn := rand.Intn(9) + 1
	tn := rand.Intn(18) + 2
	type ask struct {
		question string
		answer   int
		reason   string
		key      string
	}
	asks := []ask{
		{
			question: "How many cases are there in total?",
			key:      "total",
			answer:   tp + fp + fn + tn,
			reason:   fmt.Sprintf("Every case is one of the four cells: %d + %d + %d + %d = %d.", tp, fp, fn, tn, tp+fp+fn+tn),
		},
		{
			question: "How many actual positives are there?",
			key:      "actual_pos",
			answer:   tp + fn,
			reason:   fmt.Sprintf("Actual positives are the true positives plus the false negatives: %d + %d = %d.", tp, fn, tp+fn),
		},
		{
			question: "How many cases did the classifier predict positive?",
			key:      "predicted_pos",
			answer:   tp + fp,
			reason:   fmt.Sprintf("Predicted positives are the true positives plus the false positives: %d + %d = %d.", tp, fp, tp+fp),
		},
		{
			question: "How many predictions were correct?",
			key:      "correct",
			answer:   tp + tn,
			reason:   fmt.Sprintf("Correct predictions are the true positives plus the true negatives: %d + %d = %d.", tp, tn, tp+tn),
		},
		{
			question: "How many actual negatives are there?",
			key:      "actual_neg",
			answer:   fp + tn,
			reason:   fmt.Sprintf("Actual negatives are the false positives plus the true negatives: %d + %d = %d.", fp, tn, fp+tn),
		},
	}
	a := asks[rand.Intn(len(asks))]
	return generator.Problem{
		Question: fmt.Sprintf(
			"A classifier gives TP = %d, FP = %d, FN = %d and TN = %d. %s (enter a number)",
			tp, fp, fn, tn, a.question),
		Answer:      strconv.Itoa(a.answer),
		Explanation: a.reason,
		Facts:       map[string]string{"tp": strconv.Itoa(tp), "fp": strconv.Itoa(fp), "fn": strconv.Itoa(fn), "tn": strconv.Itoa(tn), "ask": a.key},
	}
}

type scalingGen struct{}

func (g *scalingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	sd := rand.Intn(5) + 1
	mean := rand.Intn(11)
	z := rand.Intn(9) - 4
	if z == 0 {
		z = -3
	}
	x := mean + z*sd
	return generator.Problem{
		Question: fmt.Sprintf(
			"A feature value is %d, the training mean is %d and the training standard deviation is %d. What is the standardized value (x - mean) / sd? (enter a number)",
			x, mean, sd),
		Answer:      strconv.Itoa(z),
		Explanation: fmt.Sprintf("(%d - %d)/%d = %d/%d = %d.", x, mean, sd, x-mean, sd, z),
		Facts:       map[string]string{"x": strconv.Itoa(x), "mean": strconv.Itoa(mean), "sd": strconv.Itoa(sd)},
	}
}

type crossEntropyGen struct{}

func (g *crossEntropyGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		p := rand.Intn(4) + 5 // 0.5 to 0.8
		prob := float64(p) / 10
		ce := -math.Log(prob)
		answer := fmt.Sprintf("%.2f", ce)
		return mlChoice(
			fmt.Sprintf("What is the cross-entropy loss for a true label 1 and a predicted probability %.1f? (-ln p)", prob),
			answer,
			fmt.Sprintf("-ln(%.1f) = %.2f. A confident correct prediction has small loss.", prob, ce),
			fmt.Sprintf("%.2f", prob),
			fmt.Sprintf("%.2f", 1-ce),
			fmt.Sprintf("%.2f", math.Log(prob)*-1+1),
		)
	case 1:
		q := rand.Intn(4) + 1 // 0.1 to 0.4
		prob := float64(q) / 10
		ce := -math.Log(1 - prob)
		answer := fmt.Sprintf("%.2f", ce)
		return mlChoice(
			fmt.Sprintf("What is the binary cross-entropy loss for a true label 0 and a predicted probability %.1f? (-ln(1 p))", prob),
			answer,
			fmt.Sprintf("The loss is -ln(1 - %.1f) = %.2f. Predicting %.1f for a negative case is costly because the model was confident and wrong.", prob, ce, prob),
			fmt.Sprintf("%.2f", prob),
			fmt.Sprintf("%.2f", 1-ce),
			fmt.Sprintf("%.2f", ce+1),
		)
	case 2:
		return mlChoice(
			"Does cross-entropy penalise a confident wrong prediction more than mean squared error does?",
			"yes",
			"Cross-entropy grows without bound as the predicted probability for the true class approaches zero, while a squared error stays finite.",
			"no",
		)
	case 3:
		return mlChoice(
			"Which loss is the standard choice for classification with probabilities?",
			"cross-entropy",
			"Cross-entropy is the negative log likelihood of the true labels, so minimising it maximises the probability assigned to the correct class.",
			"mean squared error",
			"mean absolute error",
			"hinge loss",
		)
	default:
		return mlChoice(
			"For a true label 1, what does the cross-entropy loss approach as the predicted probability approaches 1?",
			"0",
			"-ln(p) tends to 0 as p tends to 1, which is why a confident correct prediction is cheap.",
			"1",
			"-1",
			"infinity",
		)
	}
}

type overfittingGen struct{}

func (g *overfittingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	switch rand.Intn(5) {
	case 0:
		return mlChoice(
			"A model has very low training error and much higher test error. What is this called?",
			"overfitting",
			"The model has learned the training sample, including its noise, rather than the pattern that generalises.",
			"underfitting",
			"neither, the errors are unrelated",
			"regularisation",
		)
	case 1:
		return mlChoice(
			"Does collecting more training data generally reduce overfitting?",
			"yes",
			"More samples make it harder to fit noise, so the gap between training and test error usually narrows.",
			"no",
		)
	case 2:
		return mlChoice(
			"When error is plotted against model capacity, which shape of the test-error curve indicates overfitting at high capacity?",
			"it falls then rises",
			"Test error improves while the model captures real structure, then worsens once it starts fitting noise.",
			"it falls monotonically",
			"it rises monotonically",
			"it stays flat",
		)
	case 3:
		return mlChoice(
			"Which error measures generalisation to unseen data?",
			"test error",
			"Training error is measured on the data the model fitted, so only held-out error estimates generalisation.",
			"training error",
			"training error plus test error",
			"neither",
		)
	default:
		return mlChoice(
			"A model with high training error and high test error is best described as what?",
			"underfitting",
			"Both errors are high, so the model is too simple for the pattern rather than memorising it.",
			"overfitting",
			"well fitted",
			"regularised",
		)
	}
}
