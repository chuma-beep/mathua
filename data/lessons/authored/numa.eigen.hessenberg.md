# Hessenberg reduction

A Hessenberg reduction clears every entry below the first subdiagonal with similarity transformations. The spectrum does not change, and the QR algorithm becomes far cheaper.

## Similarity transformations

Applying \(Q^{-1}AQ\) with an orthogonal \(Q\) preserves eigenvalues and symmetry. A single Householder reflector zeroes a whole column segment at once.

## The cost

Direct reduction costs about \(\tfrac{10}{3}n^{3}\) operations, the same order as a dense LU factorisation. That cost is paid once, before any iteration begins.

## Why Hessenberg

Each QR step on a Hessenberg matrix costs \(O(n^{2})\), because only the first subdiagonal is nonzero. An unreduced matrix would cost \(O(n^{3})\) per step.

To check your understanding, apply a similarity transform, count the reduction flops and count the QR step cost.

The reduction changes the representation and never the eigenvalues.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
