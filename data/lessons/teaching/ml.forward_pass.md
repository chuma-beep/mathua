# The Forward Pass

**Forward pass:** For layer inputs $a$, compute $z = Wa + b$ then $a' = \sigma(z)$, layer after layer from input to output. It turns data into predictions — and into the loss that backpropagation later differentiates.

## Worked: an identity layer

Take $W=[[1, 0], [0, 1]]$ (identity), $b=[0, 0]$, input $x=[2, 3]$, activation $\sigma$ the identity:
1. Multiply: $Wx = [1\cdot 2 + 0\cdot 3,\; 0\cdot 2 + 1\cdot 3] = [2, 3]$.
2. Add the bias: $[2, 3] + [0, 0] = [2, 3]$.
3. Apply $\sigma$: identity leaves $[2, 3]$ unchanged.

So a layer with identity weights and zero bias passes its input through untouched — the base case every deeper computation builds on.

## Worked: a ReLU layer

Take $W=[[1, -1]]$, $b=[0]$, $x=[2, 5]$, $\sigma = \text{ReLU}$:
1. Multiply and add: $z = 1\cdot 2 + (-1)\cdot 5 + 0 = -3$.
2. Apply ReLU: $\max(0, -3) = 0$.
3. The neuron is silent: its output 0 carries no signal to the next layer.

So each layer is just affine map plus pointwise nonlinearity, repeated: the output of one layer is the next layer's input.

## Before backpropagation

The forward pass must finish before gradients exist: predictions feed the loss, the loss feeds the backward sweep. The order is always input, then forward pass, then loss, then backpropagation, then the gradient descent update — never reversed.
