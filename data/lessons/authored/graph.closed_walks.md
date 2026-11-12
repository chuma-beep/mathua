# Closed walks and eigenvalues

The trace of a power of the adjacency matrix counts closed walks, which links walk counts to the spectrum.

## The identity

The number of closed walks of length \(k\) equals \(\sum_i \lambda_i^k\), the trace of \(A^k\). Every closed walk contributes to a diagonal entry of the power.

## Small lengths

For \(k=2\) the count is \(2m\), twice the number of edges. For \(k=3\) the count is six times the number of triangles, because each triangle supports six starting points and directions.

## Open walks

The number of walks of length \(k\) from vertex \(i\) to vertex \(j\) is the \((i,j)\) entry of \(A^k\). A single step reaches the neighbours, and two steps count the length-two paths through common neighbours.

To check your understanding, raise the adjacency matrix, count length-two closed walks and read a matrix entry.

Powers of the adjacency matrix count walks, and their traces count the closed ones through the eigenvalues.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
