package machinelearning

import (
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("ml.backpropagation", &backpropGen{})
}

type backpropGen struct{}

func (g *backpropGen) Generate(difficulty float64) generator.Problem {
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
