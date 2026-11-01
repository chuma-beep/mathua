package machinelearning

import (
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("ml.backpropagation", &backpropGen{})
	reg.Register("ml.loss.cross_entropy", &crossEntropyGen{})
	reg.Register("ml.overfitting", &overfittingGen{})
	reg.Register("ml.regularization.dropout", &dropoutGen{})
	reg.Register("ml.arch.rnn", &rnnGen{})
	reg.Register("ml.optim.batchnorm", &batchNormGen{})
	reg.Register("ml.eval.roc", &rocGen{})
	reg.Register("ml.arch.transformer", &transformerGen{})
	reg.Register("ml.arch.gan", &ganGen{})
	reg.Register("ml.optim.adam", &adamDetailsGen{})
	reg.Register("ml.data.augmentation", &augmentationGen{})
	reg.Register("ml.nlp.transformer", &nlpTransformerGen{})
	reg.Register("ml.reinforcement.qlearning", &qlearningGen{})
	reg.Register("ml.reinforcement.policy_gradient", &policyGradientGen{})
	reg.Register("ml.nlp.bert", &bertGen{})
	reg.Register("ml.reinforcement.actor_critic", &actorCriticGen{})
	registerNumeric(reg)
	registerBatch(reg)
	registerConversions(reg)
}

type backpropGen struct{}

func (g *backpropGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{
			"In backpropagation, which rule from calculus is used to compute gradients through multiple layers?",
			"chain rule",
			"Backpropagation applies the chain rule to compute the gradient of the loss with respect to each weight.",
		},
		{
			"The forward pass computes the ____ given input data.",
			"output",
			"The forward pass propagates input through each layer to compute the network's output or prediction.",
		},
		{
			"The backward pass computes the ____ of the loss with respect to each weight.",
			"gradient",
			"The backward pass uses the chain rule to compute how much each weight contributed to the error.",
		},
		{
			"What optimization algorithm uses gradients computed by backpropagation to update weights?",
			"gradient descent",
			"Gradient descent updates weights in the opposite direction of the gradient to minimize the loss.",
		},
		{
			"In a neural network with sigmoid activation, the derivative at x=0 is σ(0)(1-σ(0)). What is it? (enter a number)",
			"0.25",
			"The derivative of the sigmoid function is \\(\\sigma(x)(1-\\sigma(x))\\); at x=0: 0.5·0.5=0.25, which makes backpropagation efficient.",
		},
		{
			"What does the learning rate control in gradient descent with backpropagation?",
			"step size",
			"The learning rate determines how large a step is taken in the direction of the negative gradient.",
		},
		{
			"Backpropagation is an application of the ____ rule from multivariable calculus.",
			"chain",
			"The multivariable chain rule allows computing how changes in each weight affect the final loss.",
		},
		{
			"What is the name of the problem where gradients become very small in deep networks, making learning slow?",
			"vanishing gradient",
			"Vanishing gradients occur when gradients get progressively smaller as they are backpropagated through many layers.",
		},
		{
			"Which activation function helps mitigate the vanishing gradient problem compared to sigmoid?",
			"ReLU",
			"ReLU (Rectified Linear Unit) has derivative 1 for positive inputs, helping gradients flow through deep networks.",
		},
		{
			"The chain rule means each layer's gradient depends on the gradient from the layer above it. Which layer: above or below? (type one word)",
			"above",
			"The chain rule means each layer's gradient depends on the gradient from the layer above it.",
		},
	}
	e := entries[rand.Intn(len(entries))]

	return generator.Problem{
		Question:    e.question,
		Answer:      e.answer,
		Explanation: e.exp,
	}
}

type perceptronGen struct{}

func (g *perceptronGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"What does a perceptron compute? (weighted sum + activation / convolution)", "weighted sum + activation", "A perceptron computes a weighted sum of inputs plus bias, then applies an activation function."},
		{"A perceptron with weights [0.5, -0.5] and bias 0, inputs [1,1] gives weighted sum?", "0", "0.5*1 + (-0.5)*1 + 0 = 0."},
		{"Perceptron activation is typically a ____ function.", "step (or sign)", "Original perceptron uses a step function to output 0 or 1."},
		{"Is a single perceptron able to solve XOR? (yes/no)", "no", "XOR is not linearly separable, so a single perceptron cannot solve it; needs a multilayer network."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type mseGen struct{}

func (g *mseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Mean squared error between prediction 2 and target 3 is?", "1", "MSE = (2-3)^2 = 1. Mean over one sample is 1."},
		{"MSE between predictions [0,1] and targets [1,1] is?", "0.5", "((0-1)^2 + (1-1)^2)/2 = (1+0)/2 = 0.5."},
		{"MSE is appropriate for ____ tasks.", "regression", "Mean squared error measures distance for continuous regression targets."},
		{"Derivative of MSE (1/2)(y-t)^2 w.r.t y is?", "y - t", "d/dy 1/2(y-t)^2 = y-t."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type forwardPassGen struct{}

func (g *forwardPassGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"In a forward pass, what is computed layer by layer?", "activations (or outputs)", "Forward pass propagates input through each layer's weights and activation to compute the output."},
		{"Forward pass of a linear layer with W=[[1,0],[0,1]] and x=[2,3] gives?", "[2,3]", "Identity matrix leaves the input unchanged: Wx = x."},
		{"What operation follows the weighted sum in each neuron?", "activation", "After computing the weighted sum plus bias, an activation function is applied."},
		{"Is forward pass done before or after backpropagation? (before/after)", "before", "Forward pass computes predictions; backpropagation then computes gradients."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type gradientDescentGen struct{}

func (g *gradientDescentGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct {
		question string
		answer   string
		exp      string
	}
	entries := []entry{
		{"Gradient descent updates weights in the direction of the ____ gradient.", "negative (or opposite)", "Weights are moved opposite to the gradient to minimize loss."},
		{"What does learning rate control in gradient descent?", "step size", "Learning rate scales the gradient step; too large diverges, too small is slow."},
		{"With gradient 2 and learning rate 0.1, what is the weight update magnitude?", "0.2", "Update = learning_rate * gradient = 0.1*2 = 0.2."},
		{"What variant of gradient descent uses a single sample per update?", "stochastic (or SGD)", "Stochastic gradient descent updates after each sample, while batch uses the whole dataset."},
	}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type regularizationGen struct{}

func (g *regularizationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"L2 regularization adds λ||w||² to the loss. What does λ control? (regularization strength)", "regularization strength", "λ trades off."}, {"Does L2 regularization shrink weights toward zero? (yes/no)", "yes", "Gradient 2λw pulls."}, {"What is the gradient of λ||w||² w.r.t w? (2λw)", "2λw", "d/dw λw² = 2λw."}, {"Is regularization used to reduce overfitting? (yes/no)", "yes", "Penalizing large weights."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type validationGen struct{}

func (g *validationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"What is the purpose of a validation set?", "tune hyperparameters (or model selection)", "Validation for hyperparams."}, {"k-fold cross-validation with k=5 splits data into how many folds?", "5", "k folds."}, {"Should the test set be used for hyperparameter tuning? (yes/no)", "no", "Leakage."}, {"Does cross-validation give a more reliable estimate than a single split? (yes/no)", "yes", "Averaging."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type momentumGen struct{}

func (g *momentumGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type entry struct{ question, answer, exp string }
	entries := []entry{{"Momentum update: v_{t+1}=β v_t + g_t. What does β control? (momentum / decay)", "momentum", "β controls persistence."}, {"Adam combines momentum and ____ scaling.", "adaptive (or RMSProp)", "Adaptive."}, {"Does momentum help escape shallow local minima? (yes/no)", "yes", "Velocity carries."}, {"With β=0.9, the effective memory horizon is roughly how many steps? (10)", "10", "Horizon ≈1/(1-β)."}}
	e := entries[rand.Intn(len(entries))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type dropoutGen struct{}

func (g *dropoutGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return dropoutProblem()
}

type cnnGen struct{}

func (g *cnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ question, answer, exp string }
	easy := []entry{{"CNNs use ____ layers to detect local patterns.", "convolutional", "Convolutions."}, {"Does weight sharing in CNNs reduce parameters vs fully connected? (yes/no)", "yes", "Same filter reused."}, {"What operation follows convolution in a CNN block? (activation + pooling)", "activation", "After convolution activation."}}
	hard := []entry{{"For a 5x5 input, 3x3 filter, stride 1, no padding, what is output size? (3x3)", "3x3", "(5-3)/1+1=3."}, {"Does max-pooling provide translation invariance? (yes/no)", "yes", "Robust to shifts."}, {"Are CNNs especially suited for image data? (yes/no)", "yes", "Local connectivity."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.question, Answer: e.answer, Explanation: e.exp}
}

type rnnGen struct{}

func (g *rnnGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return rnnProblem()
}

type batchNormGen struct{}

func (g *batchNormGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return batchNormProblem()
}

type rocGen struct{}

func (g *rocGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return rocProblem()
}

type transformerGen struct{}

func (g *transformerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return transformerProblem()
}

type ganGen struct{}

func (g *ganGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return ganProblem()
}

type adamDetailsGen struct{}

func (g *adamDetailsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return adamProblem()
}

type augmentationGen struct{}

func (g *augmentationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return augmentationProblem()
}

type nlpTransformerGen struct{}

func (g *nlpTransformerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return nlpTransformerProblem()
}

type qlearningGen struct{}

func (g *qlearningGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return qlearningProblem()
}

type policyGradientGen struct{}

func (g *policyGradientGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return policyGradientProblem()
}

type bertGen struct{}

func (g *bertGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return bertProblem()
}

type actorCriticGen struct{}

func (g *actorCriticGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	return actorCriticProblem()
}
