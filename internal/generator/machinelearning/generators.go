package machinelearning

import (
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("ml.perceptron", &perceptronGen{})
	reg.Register("ml.loss.mse", &mseGen{})
	reg.Register("ml.forward_pass", &forwardPassGen{})
	reg.Register("ml.gradient_descent", &gradientDescentGen{})
	reg.Register("ml.backpropagation", &backpropGen{})
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
			"output (or prediction)",
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
			"In a neural network with sigmoid activation, what is the derivative of \\(\\sigma(x) = 1/(1+e^{-x})\\) in terms of \\(\\sigma(x)\\)?",
			"σ(x)(1-σ(x))",
			"The derivative of the sigmoid function is \\(\\sigma(x)(1-\\sigma(x))\\), which makes backpropagation efficient.",
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
			"During backpropagation, the gradient of the loss with respect to a weight depends on the gradient from the ____ layer.",
			"next (or subsequent)",
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
