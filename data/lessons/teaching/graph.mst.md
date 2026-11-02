# Minimum Spanning Trees

## Spanning Trees and Weight

A spanning tree of a connected graph is a tree that uses every vertex. Its weight is the sum of its edge weights, and a minimum spanning tree has the smallest possible weight. Every connected graph has at least one spanning tree.

A spanning tree on n vertices always has exactly n-1 edges. This is why the edge count of an MST depends only on the number of vertices, while the total weight depends on the chosen edges.

## Kruskal's Algorithm

Kruskal's algorithm sorts the edges by weight and adds each edge that does not create a cycle. The process stops when n-1 edges have been added. The greedy choice is safe because of the cut property of minimum spanning trees.

The algorithm is easy to simulate on small graphs. Sorting the weights first and then checking connectivity with a disjoint-set structure gives an efficient implementation.

## Minimum Weight Examples

In a star with unit weights, every spanning tree must use all n-1 edges, so the minimum weight is n-1. In a path with listed weights, the graph is already a tree, so the only spanning tree is the path itself and its weight is the sum of the listed weights.

These examples separate the two quantities that students often confuse: the number of edges in an MST and the total weight of an MST. The first is always n-1, while the second varies with the weights.
