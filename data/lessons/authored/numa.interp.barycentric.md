# Barycentric interpolation

Barycentric interpolation writes the interpolating polynomial as a weighted average of the data values. The weights depend only on the nodes.

## Barycentric weights

For nodes \(x_i\), the barycentric weight is \(w_i=1/\prod_{j\ne i}(x_i-x_j)\). For equally spaced nodes, its absolute denominator is \(i!(n-i)!\).

The weights encode the geometry of the nodes. Changing the data values does not change the weights.

Only the ratios of the weights matter, so a common scale factor can be discarded.

## Evaluating at a node

At a node \(x_k\), the interpolant equals the data value \(y_k\). The barycentric formula avoids a division by zero by taking the limit.

Between nodes the interpolant is a weighted average of all data values. Near a node, that node’s weight dominates.

The barycentric form is stable to evaluate because it does not form the polynomial coefficients.

## Degree and node count

A polynomial of degree \(n\) has \(n+1\) coefficients, so it needs \(n+1\) distinct nodes to be determined uniquely.

More nodes do not always mean a better approximation. High-degree interpolation on equally spaced nodes can oscillate.

Barycentric interpolation is a way to evaluate the polynomial, not a cure for a bad node distribution.
