## Definition of the forward pass

The forward pass propagates input data through each layer of a neural network, computing activations from weighted sums plus biases and activation functions.

### Example

With $W=[[1,0],[0,1]]$ and $x=[2,3]$, the forward pass gives $Wx = [2,3]$.

## Layer by layer computation

Each layer computes $a^{(l)} = \sigma(W^{(l)} a^{(l-1)} + b^{(l)})$. The output of one layer becomes the input to the next.

### Example

First layer: $z = Wx + b$, then $a = \sigma(z)$.

## Before backpropagation

The forward pass must complete before backpropagation can compute gradients. Predictions are needed to evaluate the loss.

### Example

Input → forward pass → loss → backpropagation → gradient descent update.
