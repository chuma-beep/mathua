# The generalized eigenvalue problem

The generalized problem \(Ax=\lambda Bx\) appears when a mass matrix or a constraint sits beside the operator. It reduces to a standard problem once \(B\) is factored.

## Cholesky reduction

When \(B\) is symmetric positive definite, write \(B=LL^{T}\) and solve the standard problem for \(L^{-1}AL^{-T}\). Symmetry is preserved.

## Singular B

If B is singular the pencil has infinite eigenvalues, one for each null direction. Those directions have no finite vibration frequency.

## Applications

Generalized problems arise from vibration analysis, constrained least squares and elliptic operators with weighted inner products.

To check your understanding, factor the mass matrix, find the null directions and recognise vibration analysis.

Factoring the second matrix turns a two-matrix problem back into a familiar one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
