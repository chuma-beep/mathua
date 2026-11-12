# The Rayleigh quotient

The Rayleigh quotient is a weighted average of the eigenvalues with weights from the vector components. It is stationary exactly at the eigenvectors.

## Definition

For a symmetric matrix, \(R(x)=x^{T}Ax/x^{T}x\) lies between the smallest and largest eigenvalue. Scaling \(x\) leaves the quotient unchanged.

## Stationarity

A vanishing gradient means \(Ax=R(x)x\), so every stationary point is an eigenvector. The quotient there equals the corresponding eigenvalue.

## Convergence rate

Power iteration converges at a rate set by the ratio of the two largest eigenvalues in modulus. The quotient converges twice as fast as the vector.

To check your understanding, compute the quotient, set the gradient to zero and compare the eigenvalue ratio.

The quotient turns a vector into the eigenvalue it is closest to representing.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
