> Content sourced from [Algebrica](https://algebrica.org/backpropagation/) — CC BY-NC 4.0

# Backpropagation

Training a neural network means choosing weights that minimize a loss, and the
gradient of the loss with respect to every weight is needed for each update.
**Backpropagation** computes all of these gradients in one backward sweep by
reusing intermediate results — it is the chain rule, memoized across the
network's computation graph.

## Forward Pass and Loss

Fix the weights and run the network forward: each layer applies its weights,
adds biases, and passes the result through an activation, producing a
prediction \(\hat{y}\). A **loss** \(\ell(\hat{y}, y)\) scores the prediction
against the target — mean squared error for regression, cross-entropy for
classification. The forward pass caches every intermediate value, because the
backward pass will need them.

## Backward Sweep

Starting from \(d\ell\), walk the graph in reverse. At each node, multiply the
incoming gradient by the node's local Jacobian and pass the result to its
predecessors. A weight matrix \(W\) with input \(h\) and output \(o = Wh\)
receives gradient \(dW = (do) h^T\) and passes \(W^T (do)\) backward. One sweep
yields every \(dW\) at cost proportional to a single forward evaluation —
versus one full evaluation per weight for naive finite differences.

## Gradient Descent Update

With gradients in hand, each weight steps downhill:
\(W \leftarrow W - \eta \, dW\) for learning rate \(\eta\). Repeating
forward-backward-update over minibatches is **stochastic gradient descent**.
Too large an \(\eta\) diverges; too small crawls. Momentum and adaptive
variants (Adam) smooth the trajectory, but every one of them consumes the
gradients backpropagation supplies.
