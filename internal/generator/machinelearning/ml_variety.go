package machinelearning

import (
	"fmt"
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
	return crossEntropyProblem()
}

type overfittingGen struct{}

func (g *overfittingGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return overfittingProblem()
}
