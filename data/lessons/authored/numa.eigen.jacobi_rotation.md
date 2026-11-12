# Jacobi rotations

The Jacobi method drives a symmetric matrix towards diagonal form with plane rotations. Each rotation annihilates one off-diagonal pair and leaves the spectrum unchanged.

## One rotation

The angle is chosen so the chosen off-diagonal entry becomes zero. The rotation acts on rows and columns p and q only, so the rest of the matrix is untouched.

## A sweep

A sweep applies one rotation to every off-diagonal pair, so it performs \(n(n-1)/2\) rotations. Repeated sweeps converge quadratically once the matrix is nearly diagonal.

## Guarantees

Every rotation is orthogonal, so the eigenvalues are preserved exactly in exact arithmetic. The method is slower than QR but remarkably robust.

To check your understanding, choose the rotation angle, count the pairs and preserve the spectrum.

Jacobi buys certainty with arithmetic that is easy to verify.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
