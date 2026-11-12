# The distance matrix

The distance matrix records shortest-path distances, and its spectrum carries global information about the graph.

## The matrix

The distance matrix \(D\) has entry \(d(i,j)\), the length of a shortest path. It is symmetric with zero diagonal, and it has size \(n\times n\).

## Path sums

For a path on \(n\) vertices, the sum of all distances is \(n(n^2-1)/6\). The average distance is that sum divided by \(\binom{n}{2}\), which grows linearly with \(n\).

## Diameter and radius

The largest entry is the diameter and the smallest positive row maximum is the radius. The distance matrix determines both, while the adjacency matrix does not determine the distance matrix in general.

To check your understanding, index the rows by vertices, sum the distances on a path and take the largest entry.

Distances summarise the whole graph, and their sums and extremes are the first invariants to read.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
