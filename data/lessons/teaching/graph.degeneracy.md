# Graph Degeneracy

## The Degeneracy Number

A graph is k-degenerate if every subgraph has a vertex of degree at most k. The degeneracy is the smallest such k. Equivalently, repeatedly remove a vertex of minimum degree; the largest minimum degree seen is the degeneracy.

Degeneracy measures how sparse a graph is locally. Trees are 1-degenerate, cycles are 2-degenerate, and complete graphs have the largest possible degeneracy for their order.

## Small Families

The complete graph K_n has degeneracy n-1, because every vertex has degree n-1 and every subgraph is complete. A tree on at least two vertices has degeneracy 1, because it always has a leaf. A cycle has degeneracy 2, because removing one vertex leaves a path.

These values give quick sanity checks. If a proposed degeneracy is smaller than the minimum degree of the whole graph, the proposal is immediately wrong.

## Colouring by Degeneracy

A k-degenerate graph can be colored with k+1 colors by removing vertices in degeneracy order and coloring them back in reverse. This is a constructive proof, not merely an existence statement. It also explains why planar graphs, which are 5-degenerate, are 6-colorable by this method before the sharper four color theorem is applied.

Degeneracy therefore links sparsity, ordering and coloring. The same removal order appears in algorithms for list coloring and for bounded-treewidth computation.
