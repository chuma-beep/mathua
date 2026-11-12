# Cheeger's inequality

Cheeger's inequality relates the edge expansion of a graph to the second eigenvalue of its normalized Laplacian.

## Edge expansion

The Cheeger constant \(h(G)\) is the minimum of \(|\partial S|/|S|\) over sets \(S\) of at most half the vertices. It measures the cheapest way to cut off a small piece.

## The inequality

For the normalized Laplacian, \(\frac{1-\lambda_2}{2}\le h(G)\le\sqrt{2(1-\lambda_2)}\). A large spectral gap therefore forces good expansion, and good expansion forces a gap, up to a square root.

## Examples

The complete graph \(K_n\) has \(h=n/2\) and \(\lambda_2=n/(n-1)\), so the upper bound is far from tight. A long path has small \(h\) and small \(\lambda_2\), and the two sides of the inequality are of the same order.

To check your understanding, choose a small set, lower bound the gap and compute the complete graph case.

The spectral gap and the cheapest cut bound each other, with a square-root loss in one direction.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
