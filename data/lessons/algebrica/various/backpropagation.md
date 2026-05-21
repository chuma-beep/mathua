> Content sourced from [Algebrica](https://algebrica.org/backpropagation/) — CC BY-NC 4.0

Technique

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

5

Requires

0

Enables

The following concepts, [Derivative of Composite Power Functions](https://algebrica.org/derivative-of-composite-power-functions/), [Derivatives](https://algebrica.org/derivatives/), [Matrices](https://algebrica.org/matrices/), [Partial Derivatives](https://algebrica.org/partial-derivatives/), [Vectors](https://algebrica.org/vectors/), are required as prerequisites for this entry.

## What is backpropagation

Training a neural network means finding the values of the parameters, typically weight [matrices](<../matrices>), that minimize a loss function with respect to the training data. The fundamental problem is that even a moderately deep network can contain millions of parameters and computing the [gradient](<../partial-derivatives/>) of the loss with respect to each of them naively would require a number of operations proportional to the square of the network size.

Backpropagation solves this problem by exploiting the compositional structure of the network and the [chain rule](<../the-derivative-of-a-composite-function/>) of differential calculus. In this way it is possible to compute all gradients in a time proportional to that required by a single forward pass through the network.

Intuitively, backpropagation answers the question: by how much should each weight be adjusted in order to reduce the error produced by the network? The answer is contained in the gradient of the loss with respect to every parameter. The algorithm proceeds in two distinct phases.

  * In the first, called the forward pass, the input is propagated layer by layer up to the output, and the intermediate values of every computation are stored.

  * In the second, the backward pass, one starts from the error committed at the output and propagates it backwards through the network, iteratively computing the gradients with respect to each layer by means of the chain rule.


This mechanism underlies the training of virtually all modern deep learning models, from convolutional networks for computer vision to transformers employed in natural language processing.

> The backpropagation algorithm was formalized in 1986 by David E. Rumelhart, Geoffrey E. Hinton, and Ronald J. Williams in their seminal paper [Learning representations by back-propagating errors (1986)](https://www.nature.com/articles/323533a0). Their work shows that training a neural network can be reduced to a systematic application of the chain rule, enabling the computation of gradients of the loss function with respect to each parameter through local derivative updates.

## The chain rule as theoretical foundation

The computation of gradients in a neural network rests entirely on the [chain rule](<../the-derivative-of-a-composite-function/>). If a variable \\(Z\\) depends on \\(Y\\) and \\(Y\\) depends on \\(X\\), then the [derivative](<../derivatives/>) of \\(Z\\) with respect to \\(X\\) is obtained as follows:

\\[\frac{\partial Z}{\partial X} = \frac{\partial Z}{\partial Y} \cdot \frac{\partial Y}{\partial X} \\]

This identity is the mechanism that allows to connect the gradient of the loss at the network output to the gradient with respect to any internal parameter, however far from the output it may be. Each layer of the network contributes a local factor to the overall product, and backpropagation exploits this structure so that no quantity is ever recomputed twice.


For completeness, it is useful to express the chain rule in its Jacobian form. If \\( y = g(x) \\) and \\( z = f(y) \\), then holds:

\\[\nabla_x z = J_g(x)^\top \nabla_y z \\]

Here \\( J_g(x) \\) denotes the Jacobian matrix of \\( g \\) at \\( x \\), that is, the matrix of all [partial derivatives](<../partial-derivatives/>) of the components of \\( y \\) with respect to the components of \\( x \\). This formulation makes explicit how gradients are propagated backward through vector-valued transformations. For simplicity, this level of generality is not explicitly considered in the example discussed here.

## The computational graph

To make the reasoning precise, it is convenient to describe the computation performed by a neural network by means of a directed computational graph. In this graph, square nodes represent variables (parameters, activations, inputs, outputs) while circular nodes represent operations (matrix multiplication, application of a nonlinearity, computation of the loss). Directed edges indicate the dependencies between variables and operations.

![Neural network.](/diagrams/algebrica/neural-network.png)

> For simplicity, in the example presented below, the effect of bias terms along the backpropagation chain is not explicitly considered, as their contribution follows analogous derivative rules.


For example, consider a two-layer fully connected network with an activation function \\(\phi\\) applied element-wise between the two linear layers. The equations describing the forward pass are the following:

\\[\mathbf{z} = W^{(1)} \mathbf{x} \\]

\\[\mathbf{h} = \phi(\mathbf{z}) \\]

\\[\mathbf{o} = W^{(2)} \mathbf{h} \\]

\\[L = \ell(\mathbf{o},\, \mathbf{y}) \\]

In these expressions, \\(\mathbf{x}\\) is the input vector, \\(W^{(1)}\\) and \\(W^{(2)}\\) are the weight matrices of the two layers, \\(\mathbf{y}\\) is the target, and \\(L\\) is the scalar value of the loss. The goal of backpropagation is to compute the gradients \\(\frac{\partial L}{\partial W^{(1)}}\\) and \\(\frac{\partial L}{\partial W^{(2)}}\\), which will then be used by a gradient descent algorithm to update the weights.

> Typically, the problems encountered in practice are not linear in nature, and a model that applies only linear transformations to its input would be fundamentally limited in its expressive power. An activation function is a nonlinear function applied element-wise to the output of a linear layer, whose role is precisely to introduce nonlinearity into the network. Without it, the composition of any number of linear layers would itself reduce to a single linear transformation, making deeper architectures equivalent to a single-layer model. Common choices include the [sigmoid function](<../sigmoid-function/>) \\(\sigma(s) = \frac{1}{1+e^{-s}}\\), the [hyperbolic tangent](<../hyperbolic-tangent-and-cotangent/>), and the rectified linear unit \\(\text{ReLU}(s) = \max(0, s)\\).

## Reverse-mode automatic differentiation

The concrete technique by which backpropagation is implemented is called reverse-mode automatic differentiation, or reverse autodiff. The idea is that one traverses the computational graph in reverse order, from the output towards the input, and at each node computes the gradient of the loss with respect to the corresponding variable, accumulating the contributions received from the subsequent layers. The computation always starts from the derivative of the loss with respect to itself, which is trivially equal to one:

\\[\frac{\partial L}{\partial L} = 1 \\]

![Nerual network: backpropagation.](/diagrams/algebrica/neural-network-2.png)

This initial value constitutes the starting point from which the gradient is propagated backwards through all the nodes of the graph. In the following steps, we trace this propagation in detail, moving from the output layer back to the first set of parameters, and deriving at each stage the gradient that will eventually be used to update the weights.

## Step 1: gradient with respect to the output \\(\mathbf{o}\\)

The next node encountered when proceeding backwards is the one that computes the loss from the output \\(\mathbf{o}\\) and the target \\(\mathbf{y}.\\) Applying the chain rule gives:

\\[\frac{\partial L}{\partial \mathbf{o}} = \left(\frac{\partial L}{\partial L} \cdot \frac{\partial L}{\partial \mathbf{o}}\right) = \frac{\partial L}{\partial \mathbf{o}} \\]

![](/diagrams/algebrica/neural-network-3.png)

This gradient depends on the specific form of the loss function. For a quadratic loss \\(\ell(\mathbf{o}, \mathbf{y}) = |\mathbf{o} - \mathbf{y}|^2\\), for instance, one has \\(\frac{\partial L}{\partial \mathbf{o}} = 2(\mathbf{o} - \mathbf{y})\\). The concrete value of \\(\mathbf{o}\\) is available because it was computed and stored during the forward pass.

## Step 2: gradient with respect to \\(W^{(2)}\\)

The second linear layer computes \\(\mathbf{o} = W^{(2)} \mathbf{h}\\). To determine the gradient of the loss with respect to the matrix \\(W^{(2)}\\), the chain rule is applied:

\\[\frac{\partial L}{\partial W^{(2)}} = \frac{\partial L}{\partial \mathbf{o}} \cdot \frac{\partial \mathbf{o}}{\partial W^{(2)}} \\]

![](/diagrams/algebrica/neural-network-4.png)

Since \\(\mathbf{o} = W^{(2)} \mathbf{h}\\) is a linear transformation, the derivative of \\(\mathbf{o}\\) with respect to \\(W^{(2)}\\) applied to the already known gradient produces the following result:

\\[\frac{\partial L}{\partial W^{(2)}} = \frac{\partial L}{\partial \mathbf{o}} \, \mathbf{h}^{\top} \\]

The outer product between the gradient vector \\(\frac{\partial L}{\partial \mathbf{o}} \in \mathbb{R}^m\\) and the transposed activation vector \\(\mathbf{h}^{\top} \in \mathbb{R}^{1 \times d}\\) returns a matrix of the same dimensions as \\(W^{(2)}\\), as is necessary. Note that the value of \\(\mathbf{h}\\) is available because it was computed and retained during the forward pass.

## Step 3: propagating the gradient towards \\(\mathbf{h}\\)

Before being able to ascend to the first layer, it is necessary to compute the gradient of the loss with respect to the activation vector \\(\mathbf{h}\\). Once again the chain rule is applied, exploiting the fact that \\(\mathbf{o} = W^{(2)} \mathbf{h}\\):

\\[\frac{\partial L}{\partial \mathbf{h}} = \left(\frac{\partial \mathbf{o}}{\partial \mathbf{h}}\right)^{\top} \frac{\partial L}{\partial \mathbf{o}} = \left(W^{(2)}\right)^{\top} \frac{\partial L}{\partial \mathbf{o}} \\]

![](/diagrams/algebrica/neural-network-5.png)

The transpose of \\(W^{(2)}\\) arises naturally in the computation of the gradient with respect to the input of a linear transformation, and this is the mechanism that allows the error signal to propagate from the output towards the input of the network.

## Step 4: gradient through the activation function

The next node in the graph, proceeding further backwards, is the one that applies the activation function \\(\phi\\) element-wise, computing \\(\mathbf{h} = \phi(\mathbf{z})\\). Since \\(\phi\\) acts separately on each component, its derivative is diagonal and the gradient propagates as follows:

\\[\frac{\partial L}{\partial \mathbf{z}} = \frac{\partial L}{\partial \mathbf{h}} \odot \phi’(\mathbf{z}) \\]

![](/diagrams/algebrica/neural-network-6.png)

Here \\(\odot\\) denotes the element-wise product and \\(\phi’(\mathbf{z})\\) is the [vector](<../vectors/>) containing the values of the derivative of \\(\phi\\) evaluated at the points \\(z_1, z_2, \ldots, z_d\\) stored during the forward pass. This formula highlights an important aspect: if the derivative of the activation function is systematically small, the gradient attenuates at every layer until it vanishes in the deepest layers. This phenomenon is known as the vanishing gradient problem and is one of the reasons why the choice of activation function has important consequences for the training of deep networks.

## The vanishing gradient problem

The analysis above shows that every time the gradient passes through a layer with activation function \\(\phi\\), it is multiplied by \\(\phi’(\mathbf{z})\\). In a network with \\(n\\) layers, the gradient with respect to the parameters of the first layer contains a product of \\(n\\) such factors. If \\(\phi’\\) is systematically less than one in absolute value, this product decreases exponentially with \\(n\\), making the gradient negligible and rendering the update of the weights in the earliest layers effectively impossible. The sigmoid function is a particularly unfavorable example in this respect. It has the form:

\\[\sigma(s) = \frac{1}{1+e^{-s}} \\]

Its derivative satisfies \\(\sigma’(s) \leq \frac{1}{4}\\) everywhere. In a deep network, the repeated multiplication of these factors causes the gradient to decay exponentially with depth, so the earliest layers receive a signal that is a small fraction of what was computed at the output. Learning in those layers becomes extremely slow or entirely ineffective.


The ReLU function was introduced precisely to mitigate this problem: \\[\text{ReLU}(s) = \max(0, s) \\] Its derivative equals one for all positive values of the argument, which avoids the systematic attenuation of the gradient along the direction of activated neurons.

## Step 5: gradient with respect to \\(W^{(1)}\\)

The last parametric node encountered when ascending the graph is the multiplication \\(\mathbf{z} = W^{(1)} \mathbf{x}\\). The gradient with respect to \\(W^{(1)}\\) is computed with the same scheme already used for \\(W^{(2)}\\):

\\[\frac{\partial L}{\partial W^{(1)}} = \frac{\partial L}{\partial \mathbf{z}} \, \mathbf{x}^{\top} \\]

![](/diagrams/algebrica/neural-network-7.png)

Here again the gradient has the same shape as the weight matrix, and is expressed as the outer product between the error signal propagated up to this point and the input vector \\(\mathbf{x}\\), also available from the forward pass.

## Summary of the computational flow

The entire procedure can be outlined as follows. During the forward pass, \\(\mathbf{z}\\), \\(\mathbf{h}\\), \\(\mathbf{o}\\) and \\(L\\) are computed in sequence, with intermediate values stored. During the backward pass, the following gradients are computed in reverse order:

\\[\begin{array}{c} \dfrac{\partial L}{\partial L} = 1 \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial \mathbf{o}} = \dfrac{\partial \ell(\mathbf{o}, \mathbf{y})}{\partial \mathbf{o}} \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial W^{(2)}} = \dfrac{\partial L}{\partial \mathbf{o}} \, \mathbf{h}^{\top} \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial \mathbf{h}} = \left(W^{(2)}\right)^{\top} \dfrac{\partial L}{\partial \mathbf{o}} \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial \mathbf{z}} = \dfrac{\partial L}{\partial \mathbf{h}} \odot \phi’(\mathbf{z}) \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial W^{(1)}} = \dfrac{\partial L}{\partial \mathbf{z}} \, \mathbf{x}^{\top} \end{array} \\]

Each of these gradients depends exclusively on quantities already computed in previous steps, which makes the algorithm computationally efficient: every intermediate value is computed once and then reused as many times as necessary. The total number of operations is of the same order of magnitude as the forward pass, regardless of the depth of the network.

> The condition that makes backpropagation possible is the differentiability of all operations in the graph. For activation functions that are not differentiable everywhere, such as ReLU, the subgradient is adopted in practice at the point of non-differentiability, and the algorithm continues to work correctly in the context of stochastic training.

## The general case

The two-layer example developed above generalises directly to a network of arbitrary depth. Consider a network with \\(n\\) weight matrices \\(W^{(1)}, \ldots, W^{(n)}\\) and activation functions \\(\phi^{(1)}, \ldots, \phi^{(n-1)}\\) applied between successive layers. Denoting the pre-activation at layer \\(k\\) as \\(\mathbf{z}^{(k)}\\) and the post-activation as \\(\mathbf{h}^{(k)}\\), the forward pass computes, for \\(k = 1, \ldots, n\\):

\\[\mathbf{z}^{(k)} = W^{(k)} \mathbf{h}^{(k-1)}\\] \\[\mathbf{h}^{(k)} = \phi^{(k)}(\mathbf{z}^{(k)})\\]

In this case\\(\mathbf{h}^{(0)} = \mathbf{x}\\) by convention. The backward pass then propagates the gradient from layer \\(n\\) down to layer \\(1\\). Given the gradient \\(\frac{\partial L}{\partial \mathbf{h}^{(k)}}\\) arriving at layer \\(k\\), the three local updates are the following:

\\[\begin{array}{c} \dfrac{\partial L}{\partial \mathbf{z}^{(k)}} = \dfrac{\partial L}{\partial \mathbf{h}^{(k)}} \odot \left(\phi^{(k)}\right)'(\mathbf{z}^{(k)}) \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial W^{(k)}} = \dfrac{\partial L}{\partial \mathbf{z}^{(k)}} \left(\mathbf{h}^{(k-1)}\right)^{\top} \\\\[10pt] \downarrow \\\\[10pt] \dfrac{\partial L}{\partial \mathbf{h}^{(k-1)}} = \left(W^{(k)}\right)^{\top} \dfrac{\partial L}{\partial \mathbf{z}^{(k)}} \end{array} \\]

The third equation propagates the gradient to the layer below, making the recursion self-sustaining. Starting from \\(\frac{\partial L}{\partial \mathbf{h}^{(n)}} = \frac{\partial L}{\partial \mathbf{o}}\\) and iterating from \\(k = n\\) down to \\(k = 1\\), one recovers all the gradients needed to update every weight matrix in the network.

Other Topics

Topics that warrant a complete treatment but do not yet belong to a defined category.

22.8k

[Cosine Similarity](https://algebrica.org/cosine-similarity/)

1 comment

16.3k

[Propositional Logic](https://algebrica.org/propositional-logic/)
