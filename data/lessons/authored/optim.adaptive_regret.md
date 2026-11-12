# Regret bounds for adaptive methods

An adaptive method can be analysed by regret against the best fixed comparator. The bound separates the cost of the gradient scale from the cost of the number of steps.

## The regret definition

Regret is the accumulated loss of the method minus the loss of the best fixed point in hindsight. A sublinear bound means the average regret vanishes, so the method approaches the best comparator.

## The adaptive bound

The bound depends on the square root of the accumulated squared gradients and on the square root of the number of steps. Coordinates with a large history contribute more, which is why the method suits uneven scales.

## Reading the bound

A bound that grows with the gradient scale explains why adaptive methods tolerate poorly scaled data. It also shows the limit: the price is logarithmic factors and a dependence on the comparator norm.

To check your understanding, sum the loss difference, use accumulated squares and explain scale tolerance.

Regret is the right lens for online and stochastic optimisation because it compares against the best fixed choice rather than against a solution the method never sees.
