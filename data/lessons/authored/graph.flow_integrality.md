# The integral flow theorem

When every capacity is an integer, a maximum flow can be chosen integral. The theorem is what makes flow algorithms combinatorial rather than continuous.

## Integral optima

If all capacities are integers, the maximum flow value is an integer and some maximum flow assigns an integer to every edge. The result follows from augmenting along paths with integral residual capacity.

## The cut value

The max-flow min-cut theorem identifies the maximum flow value with the minimum cut capacity. An integral cut capacity therefore forces an integral optimum, and a cut is a finite sum of integers.

## Total unimodularity

The incidence matrix of a directed graph is totally unimodular. Integer right-hand sides therefore give integral vertices of the flow polytope, which is the linear-programming explanation of the theorem.

To check your understanding, augment along a path, sum the cut capacities and write the incidence matrix.

Integer capacities guarantee an integer optimum, so no rounding step is ever needed.
