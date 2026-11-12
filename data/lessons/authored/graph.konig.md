# Konig's theorem

Konig's theorem equates the maximum matching and the minimum vertex cover in a bipartite graph.

## The equality

In a bipartite graph the size of a maximum matching equals the size of a minimum vertex cover. One inequality is immediate, since a cover must meet every edge of a matching, and the other comes from an augmenting-path argument.

## Duality

The theorem is the integral case of linear-programming duality for the matching polytope. The matching and cover problems are a primal-dual pair, and bipartiteness makes both optima integral.

## Small examples

A path on \(n\) vertices has a maximum matching of \(\lfloor n/2\rfloor\) and a minimum cover of the same size. A cycle on \(n\) vertices has matching \(\lfloor n/2\rfloor\) and cover \(\lfloor n/2\rfloor\). The complete bipartite graph \(K_{m,n}\) has both equal to \(\min(m,n)\).

To check your understanding, bound a cover below by a matching, write the matching LP and compute a path matching.

Matchings and vertex covers are dual in bipartite graphs, and their optimal values coincide.
