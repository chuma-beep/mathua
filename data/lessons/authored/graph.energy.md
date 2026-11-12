# Graph energy

Graph energy is the sum of the absolute values of the adjacency eigenvalues, a spectral measure of how far the graph is from being empty.

## The definition

The energy of a graph is \(E(G)=\sum_i|\lambda_i|\), where the sum runs over the adjacency eigenvalues. It is always nonnegative and equals zero only for the empty graph.

## Complete graphs

The complete graph \(K_n\) has eigenvalues \(n-1\) and \(-1\) with multiplicity \(n-1\), so its energy is \((n-1)+(n-1)=2(n-1)\).

## Bipartite graphs

A bipartite graph has a symmetric spectrum, so its energy is twice the sum of the positive eigenvalues. The complete bipartite graph \(K_{m,n}\) has energy \(2\sqrt{mn}\), which follows from its two nonzero eigenvalues \(\pm\sqrt{mn}\).

To check your understanding, list the eigenvalues, find the eigenvalues and use spectral symmetry.

Energy adds the eigenvalue magnitudes, so it counts spectral size without regard to sign.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
