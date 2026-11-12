# The singular value decomposition

Every matrix factors as \(A=U\Sigma V^{T}\) with orthogonal factors and a diagonal \(\Sigma\) of nonnegative singular values. The factorisation exposes rank, range and conditioning at once.

## Singular values

The singular values are the square roots of the eigenvalues of \(A^{T}A\). There are \(\min(m,n)\) of them, counted with multiplicity.

## Low rank

Zero singular values mark directions the matrix annihilates. A truncated SVD keeps the largest k values and gives the best rank-k approximation in the Frobenius norm.

## Computation

Golub-Kahan bidiagonalisation followed by implicit QR is the standard dense algorithm. Its cost is \(O(mn^{2})\) when \(m\ge n\).

To check your understanding, relate values to eigenvalues, identify the null directions and bidiagonalise the matrix.

The SVD separates what a matrix stretches from how much it stretches.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
