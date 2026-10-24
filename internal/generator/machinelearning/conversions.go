package machinelearning

// The second batch of multiple-choice concepts turned into computed questions.
//
// Machine learning asked a learner to recognise a definition twenty-nine times
// out of forty-five. These ten have an answer that is a number, so they keep
// their content and lose only the option list: a mean squared error, a
// perceptron's predictions, a rectified forward pass, one gradient-descent
// step, an L2 penalty, a fold size, a one-hot width, a convolutional
// parameter count, an expected calibration error and one momentum update.
//
// Every generator publishes the numbers its answer follows from, so the test
// recomputes the answer rather than trusting the arithmetic here.

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
)

// ---------------------------------------------------------------------------
// Mean squared error
// ---------------------------------------------------------------------------

type mseConvertedGen struct{}

func (g *mseConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; attempt < 100; attempt++ {
		t0 := rand.Intn(10) + 2
		t1 := rand.Intn(10) + 2
		t2 := rand.Intn(10) + 2

		d0 := rand.Intn(13) - 6
		d1 := rand.Intn(13) - 6
		d2 := rand.Intn(13) - 6
		sumSq := d0*d0 + d1*d1 + d2*d2
		if sumSq == 0 || sumSq%3 != 0 {
			continue
		}

		p0 := t0 + d0
		p1 := t1 + d1
		p2 := t2 + d2
		if p0 < 0 || p1 < 0 || p2 < 0 {
			continue
		}

		mse := sumSq / 3
		if p0 == mse || p1 == mse || p2 == mse || t0 == mse || t1 == mse || t2 == mse {
			continue
		}

		ansStr := strconv.Itoa(mse)
		return generator.Problem{
			Question: fmt.Sprintf("A model predicts %d, %d and %d for three targets that are %d, %d and %d. What is the mean squared error?",
				p0, p1, p2, t0, t1, t2),
			Answer: ansStr,
			Explanation: fmt.Sprintf("The errors are %d, %d and %d. Their squares add to %d, and dividing by the three samples gives a mean squared error of %d.",
				d0, d1, d2, sumSq, mse),
			Facts: map[string]string{
				"predictions": fmt.Sprintf("%d, %d, %d", p0, p1, p2),
				"targets":     fmt.Sprintf("%d, %d, %d", t0, t1, t2),
				"answer":      ansStr,
			},
		}
	}

	return generator.Problem{
		Question:    "A model predicts 1, 6 and 0 for three targets that are 3, 5 and 2. What is the mean squared error?",
		Answer:      "3",
		Explanation: "The errors are -2, 1 and -2. Their squares add to 9, and dividing by three gives a mean squared error of 3.",
		Facts: map[string]string{
			"predictions": "1, 6, 0",
			"targets":     "3, 5, 2",
			"answer":      "3",
		},
	}
}

// ---------------------------------------------------------------------------
// Perceptron predictions
// ---------------------------------------------------------------------------

type perceptronConvertedGen struct{}

func (g *perceptronConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	w1 := smallInt(-3, 3)
	w2 := smallInt(-3, 3)
	b := smallInt(-2, 2)
	type point struct{ x1, x2 int }
	points := []point{{x1: 1, x2: 1}, {x1: 2, x2: 1}, {x1: 1, x2: 2}}
	positives := 0
	for _, p := range points {
		if w1*p.x1+w2*p.x2+b > 0 {
			positives++
		}
	}
	return generator.Problem{
		Question: fmt.Sprintf("A perceptron has weights %d and %d and bias %d, and predicts +1 when the weighted sum plus the bias is positive. For how many of the points (1, 1), (2, 1) and (1, 2) does it predict +1?",
			w1, w2, b),
		Answer: fmt.Sprintf("%d", positives),
		Explanation: fmt.Sprintf("The weighted sums plus the bias at (1, 1), (2, 1) and (1, 2) are %d, %d and %d, so %d of the three are positive.",
			w1+w2+b, 2*w1+w2+b, w1+2*w2+b, positives),
		Facts: map[string]string{
			"weights": fmt.Sprintf("%d, %d", w1, w2),
			"bias":    fmt.Sprintf("%d", b),
			"answer":  fmt.Sprintf("%d", positives),
		},
	}
}

// ---------------------------------------------------------------------------
// A rectified forward pass
// ---------------------------------------------------------------------------

type forwardPassConvertedGen struct{}

func (g *forwardPassConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	w1 := smallInt(-2, 3)
	w2 := smallInt(-2, 3)
	x1 := smallInt(1, 3)
	x2 := smallInt(1, 3)
	b := smallInt(-3, 2)
	raw := w1*x1 + w2*x2 + b
	out := raw
	if out < 0 {
		out = 0
	}
	return generator.Problem{
		Question: fmt.Sprintf("A neuron computes \\(%d \\cdot %d + %d \\cdot %d + %d\\) and then applies the rectifier, which replaces a negative value with zero. What is its output?",
			w1, x1, w2, x2, b),
		Answer: fmt.Sprintf("%d", out),
		Explanation: fmt.Sprintf("The weighted sum plus the bias is %d. The rectifier leaves it unchanged when it is positive and returns zero otherwise, so the output is %d.",
			raw, out),
		Facts: map[string]string{
			"weights": fmt.Sprintf("%d, %d", w1, w2),
			"inputs":  fmt.Sprintf("%d, %d", x1, x2),
			"bias":    fmt.Sprintf("%d", b),
			"answer":  fmt.Sprintf("%d", out),
		},
	}
}

// ---------------------------------------------------------------------------
// One gradient-descent step
// ---------------------------------------------------------------------------

type gradientDescentConvertedGen struct{}

func (g *gradientDescentConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	weight := smallInt(4, 20)
	gradient := smallInt(1, 5)
	rate := smallInt(1, 3)
	next := weight - rate*gradient
	return generator.Problem{
		Question: fmt.Sprintf("A weight is %d, the gradient of the loss with respect to it is %d, and the learning rate is %d. What is the weight after one gradient-descent step?",
			weight, gradient, rate),
		Answer: fmt.Sprintf("%d", next),
		Explanation: fmt.Sprintf("Gradient descent subtracts the learning rate times the gradient, so the new weight is %d minus %d times %d, which is %d.",
			weight, rate, gradient, next),
		Facts: map[string]string{
			"weight":        fmt.Sprintf("%d", weight),
			"gradient":      fmt.Sprintf("%d", gradient),
			"learning_rate": fmt.Sprintf("%d", rate),
			"answer":        fmt.Sprintf("%d", next),
		},
	}
}

// ---------------------------------------------------------------------------
// The L2 penalty
// ---------------------------------------------------------------------------

type regularizationConvertedGen struct{}

func (g *regularizationConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	lambda := smallInt(1, 4)
	w1 := smallInt(-4, 4)
	w2 := smallInt(-4, 4)
	penalty := lambda * (w1*w1 + w2*w2)
	return generator.Problem{
		Question: fmt.Sprintf("L2 regularization with strength %d adds the penalty \\(\\lambda (w_1^2 + w_2^2)\\) to the loss. For weights %d and %d, what is the added penalty?",
			lambda, w1, w2),
		Answer: fmt.Sprintf("%d", penalty),
		Explanation: fmt.Sprintf("The squared weights are %d and %d. Their sum is %d, and multiplying by the strength %d gives a penalty of %d.",
			w1*w1, w2*w2, w1*w1+w2*w2, lambda, penalty),
		Facts: map[string]string{
			"lambda":  fmt.Sprintf("%d", lambda),
			"weights": fmt.Sprintf("%d, %d", w1, w2),
			"answer":  fmt.Sprintf("%d", penalty),
		},
	}
}

// ---------------------------------------------------------------------------
// Fold size in cross-validation
// ---------------------------------------------------------------------------

type validationConvertedGen struct{}

func (g *validationConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	folds := smallInt(3, 10)
	perFold := smallInt(8, 40)
	samples := folds * perFold
	return generator.Problem{
		Question: fmt.Sprintf("A dataset has %d examples and is split into %d folds of equal size for cross-validation. How many examples are in each fold?",
			samples, folds),
		Answer: fmt.Sprintf("%d", perFold),
		Explanation: fmt.Sprintf("Equal folds divide the %d examples into %d groups, so each fold holds %d examples.",
			samples, folds, perFold),
		Facts: map[string]string{
			"samples":   fmt.Sprintf("%d", samples),
			"folds":     fmt.Sprintf("%d", folds),
			"fold_size": fmt.Sprintf("%d", perFold),
			"answer":    fmt.Sprintf("%d", perFold),
		},
	}
}

// ---------------------------------------------------------------------------
// One-hot width with drop-first
// ---------------------------------------------------------------------------

type oneHotConvertedGen struct{}

func (g *oneHotConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	categories := smallInt(4, 12)
	columns := categories - 1
	return generator.Problem{
		Question: fmt.Sprintf("A categorical feature has %d categories. One-hot encoding with drop-first keeps one column fewer than the number of categories so that the columns are not collinear. How many columns does it use?",
			categories),
		Answer: fmt.Sprintf("%d", columns),
		Explanation: fmt.Sprintf("Without drop-first there would be %d columns, one per category. Dropping the first category leaves %d columns.",
			categories, columns),
		Facts: map[string]string{
			"categories": fmt.Sprintf("%d", categories),
			"columns":    fmt.Sprintf("%d", columns),
			"answer":     fmt.Sprintf("%d", columns),
		},
	}
}

// ---------------------------------------------------------------------------
// Parameters in a convolutional layer
// ---------------------------------------------------------------------------

type cnnConvertedGen struct{}

func (g *cnnConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	kernel := 3
	if rand.Intn(2) == 1 {
		kernel = 5
	}
	inChannels := smallInt(1, 4)
	outChannels := smallInt(4, 32)
	parameters := (kernel*kernel*inChannels + 1) * outChannels
	return generator.Problem{
		Question: fmt.Sprintf("A convolutional layer has a %d by %d kernel, %d input channels, %d output channels and one bias per output channel. How many trainable parameters does it have?",
			kernel, kernel, inChannels, outChannels),
		Answer: fmt.Sprintf("%d", parameters),
		Explanation: fmt.Sprintf("Each output channel owns one kernel per input channel, so it has %d times %d weights plus one bias, which is %d. Across %d output channels that is %d parameters.",
			kernel*kernel, inChannels, kernel*kernel*inChannels+1, outChannels, parameters),
		Facts: map[string]string{
			"kernel":       fmt.Sprintf("%d", kernel),
			"in_channels":  fmt.Sprintf("%d", inChannels),
			"out_channels": fmt.Sprintf("%d", outChannels),
			"answer":       fmt.Sprintf("%d", parameters),
		},
	}
}

// ---------------------------------------------------------------------------
// Expected calibration error
// ---------------------------------------------------------------------------

type calibrationConvertedGen struct{}

func (g *calibrationConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	// The gap and the bin weights are drawn, and the answer is the weighted average of the
	// gaps, so the arithmetic is exact: the weights are multiples of 20 out of 100 and the
	// gaps are multiples of 5, which keeps the division integral.
	gaps := []int{5, 10, 15, 20, 25}
	confs := []int{70, 75, 80, 85, 90, 95}
	k := rand.Intn(4) + 1 // bin one holds 20k of the 100 examples
	n1, n2 := 20*k, 100-20*k
	g1 := gaps[rand.Intn(len(gaps))]
	g2 := gaps[rand.Intn(len(gaps))]
	conf1 := confs[rand.Intn(len(confs))]
	conf2 := confs[rand.Intn(len(confs))]
	acc1 := conf1 - g1
	acc2 := conf2 - g2
	if acc1 <= 0 || acc2 <= 0 {
		acc1, acc2 = conf1, conf2
		g1, g2 = 0, 0
	}
	weighted := n1*g1 + n2*g2
	ece := weighted / 100
	return generator.Problem{
		Question: fmt.Sprintf("A model is checked in two confidence bins. The first holds %d examples with %d percent accuracy at %d percent confidence; the second holds %d examples with %d percent accuracy at %d percent confidence. What is the expected calibration error, in percentage points?",
			n1, acc1, conf1, n2, acc2, conf2),
		Answer: fmt.Sprintf("%d", ece),
		Explanation: fmt.Sprintf("The gaps between accuracy and confidence are %d and %d points. Weighting them by the %d and %d examples and dividing by %d gives an expected calibration error of %d points.",
			g1, g2, n1, n2, n1+n2, ece),
		Facts: map[string]string{
			"bin_one": fmt.Sprintf("%d examples, %d accuracy, %d confidence", n1, acc1, conf1),
			"bin_two": fmt.Sprintf("%d examples, %d accuracy, %d confidence", n2, acc2, conf2),
			"samples": fmt.Sprintf("%d", n1+n2),
			"answer":  fmt.Sprintf("%d", ece),
		},
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// ---------------------------------------------------------------------------
// One momentum update
// ---------------------------------------------------------------------------

type momentumConvertedGen struct{}

func (g *momentumConvertedGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	velocity := 2 * smallInt(1, 5)
	gradient := smallInt(1, 4)
	weight := smallInt(6, 20)
	rate := smallInt(1, 2)
	nextVelocity := velocity/2 + gradient
	next := weight - rate*nextVelocity
	return generator.Problem{
		Question: fmt.Sprintf("Momentum with factor 0.5 keeps a velocity of %d, the new gradient is %d, the learning rate is %d and the current weight is %d. With the update \\(v \\leftarrow 0.5v + g\\) and \\(w \\leftarrow w - \\alpha v\\), what is the new weight?",
			velocity, gradient, rate, weight),
		Answer: fmt.Sprintf("%d", next),
		Explanation: fmt.Sprintf("The velocity becomes half of %d plus %d, which is %d. The weight then drops by %d times %d, so it becomes %d.",
			velocity, gradient, nextVelocity, rate, nextVelocity, next),
		Facts: map[string]string{
			"velocity":      fmt.Sprintf("%d", velocity),
			"gradient":      fmt.Sprintf("%d", gradient),
			"momentum":      "0.5",
			"learning_rate": fmt.Sprintf("%d", rate),
			"weight":        fmt.Sprintf("%d", weight),
			"answer":        fmt.Sprintf("%d", next),
		},
	}
}

// registerConversions registers the concepts of this batch that used to be
// multiple choice. Their old registrations are removed from Register, because
// a concept id may be registered only once.
func registerConversions(reg *generator.Registry) {
	reg.Register("ml.loss.mse", &mseConvertedGen{})
	reg.Register("ml.perceptron", &perceptronConvertedGen{})
	reg.Register("ml.forward_pass", &forwardPassConvertedGen{})
	reg.Register("ml.gradient_descent", &gradientDescentConvertedGen{})
	reg.Register("ml.regularization", &regularizationConvertedGen{})
	reg.Register("ml.validation", &validationConvertedGen{})
	reg.Register("ml.data.onehot", &oneHotConvertedGen{})
	reg.Register("ml.arch.cnn", &cnnConvertedGen{})
	reg.Register("ml.eval.calibration", &calibrationConvertedGen{})
	reg.Register("ml.optim.momentum", &momentumConvertedGen{})
}
