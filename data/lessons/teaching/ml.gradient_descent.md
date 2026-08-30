## Definition of gradient descent

Gradient descent is an optimization algorithm that minimizes a loss by moving parameters in the direction of the negative gradient.

### Example

With gradient $2$ and learning rate $0.1$, the update is $0.2$ in the negative direction.

## The learning rate

The learning rate controls step size. Too large diverges, too small is slow. It scales the gradient update.

### Example

Learning rate $0.01$ with gradient $2$ gives update $0.02$.

## Stochastic and batch variants

Stochastic gradient descent updates after each sample; batch gradient descent uses the whole dataset. Mini-batch is a compromise.

### Example

SGD with one sample per update is noisy but fast; batch is stable but slower per step.
