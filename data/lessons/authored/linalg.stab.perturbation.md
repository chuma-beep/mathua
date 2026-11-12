# Perturbation theory

Perturbation theory describes how eigenvalues and eigenvectors move when a matrix is perturbed. Eigenvalues are stable for normal matrices and can be far more sensitive otherwise.

## Eigenvalues

A perturbation of norm delta moves each eigenvalue by at most delta when the matrix is normal. For non-normal matrices the bound involves the condition of the eigenvector basis.

## Eigenvectors

Eigenvector sensitivity is governed by the gap between eigenvalues. A small gap makes the direction turn quickly.

## Pseudospectra

For non-normal matrices the pseudospectrum can be much larger than the union of epsilon discs around eigenvalues. It is the right tool for transient behaviour.

To check your understanding, bound by delta, measure the spectral gap and compare with discs.

Eigenvalue perturbation is bounded by the perturbation for normal matrices and by eigenvector conditioning otherwise.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
