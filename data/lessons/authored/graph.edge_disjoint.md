# Edge-disjoint paths

Menger's theorem in its edge form equates the maximum number of edge-disjoint paths with the minimum size of an edge cut.

## The edge form

The maximum number of pairwise edge-disjoint paths from \(s\) to \(t\) equals the minimum number of edges whose removal separates \(s\) from \(t\). The value is the edge connectivity between the two vertices.

## Global connectivity

The edge connectivity of a connected graph is the minimum number of edges whose removal disconnects it. It is bounded above by the minimum degree, and a complete graph on \(n\) vertices attains \(n-1\).

## Small values

A path has edge connectivity one, a cycle has two, and a tree has one because every edge is a bridge. These values are also the numbers of edge-disjoint paths between suitable endpoint pairs.

To check your understanding, pack edge-disjoint paths, bound by the minimum degree and check a path.

Paths that share no edge and edges whose removal separates two vertices are dual quantities.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
