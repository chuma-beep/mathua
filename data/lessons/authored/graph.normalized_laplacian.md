# The normalized Laplacian

The normalized Laplacian divides the Laplacian by the degrees, which makes its spectrum comparable across graphs of different size.

## The definition

The normalized Laplacian is \(\mathcal{L}=I-D^{-1/2}AD^{-1/2}\) for a graph with no isolated vertices. Its entries are one on the diagonal, \(-1/\sqrt{d_i d_j}\) across an edge and zero elsewhere.

## The spectrum

The eigenvalues lie in \([0,2]\), and the number of zero eigenvalues equals the number of connected components. The trace is \(n\), the number of vertices.

## Comparison

For a \(d\)-regular graph, \(\mathcal{L}=L/d\), so the normalized spectrum is the ordinary Laplacian spectrum scaled by \(d\). The complete graph \(K_n\) has eigenvalues zero and \(n/(n-1)\) with multiplicity \(n-1\).

To check your understanding, invert the square-root degrees, bound the eigenvalues and scale a regular graph.

Normalizing by degree puts every graph on the same spectral scale without changing the component count.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
