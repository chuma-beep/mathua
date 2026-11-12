# Conditioning from singular values

The 2-norm condition number is the ratio of the largest to the smallest singular value. It measures how much a relative error in the data can grow in the solution.

## Definition

\(\kappa_2(A)=\sigma_{\max}/\sigma_{\min}\) for an invertible matrix. A singular matrix has an infinite condition number.

## Error bound

A relative perturbation of size \(\epsilon\) in \(A\) can change the solution of \(Ax=b\) by about \(\kappa_2(A)\epsilon\). The bound is sharp for some right-hand sides.

## Reading the ratio

A condition number near one means well conditioned. Each factor of ten adds roughly one lost decimal digit in the answer.

To check your understanding, form the ratio, perturb the matrix and interpret a small ratio.

The singular value ratio is the honest measure of how much an answer can drift.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
