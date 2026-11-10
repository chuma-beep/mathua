# Graph thickness

## Planar layers

The thickness of a graph is the smallest number of planar subgraphs whose union is the whole graph. Each layer must be planar, and every edge belongs to one layer.

Thickness measures how far a graph is from being planar. A planar graph has thickness one.

## The edge bound

A planar graph with n at least 3 has at most 3n - 6 edges. Dividing the edge count of the graph by this bound and rounding up gives a lower bound on thickness.

The bound is not always exact, but it is a reliable first estimate. Complete graphs often meet it or come close.

## Complete graphs

The thickness of K_n grows slowly. Small complete graphs have thickness two, and larger ones need three or more planar layers.

Check a proposed layer decomposition by verifying that each layer is planar and that the layers together contain every edge exactly once.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
