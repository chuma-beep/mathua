## Definition of a perceptron

A perceptron is the simplest neural unit. It computes a weighted sum of inputs plus a bias, then applies an activation function to produce an output.

### Example

With weights $[0.5, -0.5]$, bias $0$, and inputs $[1,1]$, the weighted sum is $0$ and the step activation gives $0$.

## The weighted sum

For inputs $x$ and weights $w$, the weighted sum is $w \cdot x + b$. This is a dot product plus bias.

### Example

$w=[1,2]$, $x=[3,4]$, $b=0$ gives $1\cdot3 + 2\cdot4 = 11$.

## Activation functions

The perceptron uses an activation to decide its output: step, sign, sigmoid, or ReLU. The choice determines the decision boundary.

### Example

Step activation outputs $1$ if the sum is positive, otherwise $0$.
