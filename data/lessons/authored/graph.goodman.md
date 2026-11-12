# Goodman's formula

Goodman's formula counts the monochromatic triangles in a two-colouring of the edges from the degree sequence alone.

## The formula

For a graph with degrees \(d_i\), the number of monochromatic triangles in the two-colouring where one colour is the graph and the other its complement is \(\binom{n}{3}-\frac12\sum_i d_i(n-1-d_i)\).

## The two-colour split

At a vertex \(i\), the \(d_i\) edges of the graph and the \(n-1-d_i\) edges of the complement form \(d_i(n-1-d_i)\) mixed pairs. Each non-monochromatic triangle is counted twice this way.

## The lower bound

Every two-colouring of \(K_6\) contains a monochromatic triangle, and Goodman's formula gives the exact count once the degrees are known. A regular graph minimises the mixed pairs and therefore maximises the monochromatic ones.

To check your understanding, count all triangles, split the edges at a vertex and evaluate the formula on a degree list.

Counting mixed pairs at each vertex and subtracting from all triangles gives the monochromatic count exactly.
