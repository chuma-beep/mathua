# Kernel principal components

Kernel principal components perform principal component analysis in a feature space without computing its coordinates. The eigenvectors of the centred Gram matrix supply the projections.

## Centring

The Gram matrix is centred by subtracting row means, column means and the grand mean. Without centring the components do not have zero mean in feature space.

## The eigenproblem

An \(n\)-point data set yields an \(n \times n\) eigenproblem. The number of non-zero eigenvalues is at most the feature dimension.

## Projection

A test point projects onto a component through kernel values against the training points. Each projection costs \(n\) kernel evaluations.

To check your understanding, subtract the row means, size the matrix and evaluate the kernels.

Kernel PCA finds nonlinear directions with a linear eigenproblem on the Gram matrix.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
