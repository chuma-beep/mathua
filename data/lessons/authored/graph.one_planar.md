# 1-planar graphs

A 1-planar graph can be drawn with at most one crossing on every edge, which buys a larger edge budget than planarity.

## The definition

A graph is 1-planar when it has a drawing in which every edge is crossed at most once. Planar graphs are the 1-planar graphs with no crossings.

## The edge bound

A simple 1-planar graph on \(n\ge 3\) vertices has at most \(4n-8\) edges. Every crossing pairs two edges, and the resulting planarisation has a vertex of degree four at each crossing.

## Examples

The complete graph \(K_6\) is 1-planar with 15 edges, below the bound 16 for six vertices. The complete graph \(K_7\) has 21 edges and exceeds \(4\cdot 7-8=20\), so it is not 1-planar.

To check your understanding, draw the graph, planarise the crossings and check \(K_6\).

One crossing per edge is a modest allowance that still leaves a linear edge bound.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
