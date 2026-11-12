# Eigenvalue sensitivity

An eigenvalue moves when the matrix is perturbed, and the movement depends on the angle between its left and right eigenvectors. Near-defective eigenvalues are far more sensitive than the matrix condition number suggests.

## First-order shift

For a simple eigenvalue, \(d\lambda=v^{T}Eu/(v^{T}u)\), where \(u\) and \(v\) are right and left eigenvectors. The denominator can be tiny even when \(A\) is well conditioned.

## The amplification factor

The quantity \(1/|v^{T}u|\) is the condition number of the eigenvalue. It equals one for a normal matrix and grows without bound as the eigenvalue becomes defective.

## Practical reading

A cluster of eigenvalues is hard to resolve individually because perturbations mix the cluster. Algorithms therefore compute invariant subspaces rather than single vectors.

To check your understanding, perturb the matrix, form the reciprocal overlap and resolve a cluster.

An eigenvalue is only as stable as the angle between its two eigenvectors.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
