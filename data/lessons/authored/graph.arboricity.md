# Arboricity

## Forests

A forest is an acyclic graph. The arboricity of a graph is the smallest number of forests needed to cover all of its edges. Every graph has an arboricity.

A forest on n vertices has at most n - 1 edges. This edge bound is the key to estimating arboricity.

## The edge bound

For any subgraph with n vertices and m edges, the arboricity is at least the ceiling of m divided by n - 1. The maximum of this bound over all subgraphs is the exact arboricity by the Nash-Williams theorem.

Use the whole graph as one subgraph to get a quick lower bound. Dense subgraphs can raise the bound.

## Complete graphs

The complete graph K_n has arboricity equal to the ceiling of n/2. This follows from its edge count and the forest edge bound, and the bound is attainable.

A forest decomposition partitions the edges, not the vertices. The same vertex can appear in several forests.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
