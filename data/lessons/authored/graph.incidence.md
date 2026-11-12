# The incidence matrix

The incidence matrix records which vertex meets which edge, and its Gram matrix recovers the Laplacian.

## The matrix

For a graph with \(n\) vertices and \(m\) edges, the incidence matrix \(B\) has size \(n\times m\). An entry is one when the vertex is an endpoint of the edge and zero otherwise.

## The Gram matrix

The product \(BB^{T}\) equals \(D+A\), where \(D\) is the degree matrix and \(A\) the adjacency matrix. This identity is the reason the Laplacian \(D-A\) has a signed incidence version.

## Rank and components

For a connected graph the incidence matrix has rank \(n-1\) over the reals. In general the rank is \(n-c\), where \(c\) is the number of connected components, and the null space records component indicators.

To check your understanding, size the matrix, multiply the incidence matrix and count the components.

Incidence ties edges to vertices, and multiplying by its transpose rebuilds the degree and adjacency data.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
