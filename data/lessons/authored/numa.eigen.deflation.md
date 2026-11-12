# Deflation

Deflation removes an eigenpair that has already been found so the iteration can target the next one. Each deflation lowers the effective dimension by one.

## Symmetric deflation

For symmetric \(A\), replacing \(A\) by \(A-\lambda vv^{T}\) leaves every other eigenpair untouched when \(v\) is a unit eigenvector. The found eigenvalue becomes zero.

## Nonsymmetric deflation

A nonsymmetric matrix can be deflated by a rank-one modification that keeps the remaining spectrum. The construction needs both left and right eigenvectors.

## Use in practice

Subspace iteration deflates a whole block at once. The cost of each deflation is one rank-one update, which is cheap compared with another full solve.

To check your understanding, subtract the rank-one term, use left and right vectors and deflate a block.

Deflation lets one solve produce a sequence of eigenpairs instead of a single one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
