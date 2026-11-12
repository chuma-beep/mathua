# Spectral methods

Spectral methods represent the solution in a global basis of smooth functions. For smooth data the error falls faster than any power of the grid spacing.

## Global basis

Fourier modes suit periodic problems; Chebyshev polynomials suit bounded intervals. The coefficients are computed by a transform.

## Convergence

The error decays exponentially when the solution is analytic. A discontinuity destroys the rate and produces Gibbs oscillations.

## Cost

A transform costs \(O(N\log N)\) with the FFT, and a derivative is a diagonal multiplication in coefficient space. The dense operator is never formed.

To check your understanding, choose Fourier modes, state exponential decay and use the FFT.

Spectral accuracy is a reward for smoothness and is withdrawn without warning.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
