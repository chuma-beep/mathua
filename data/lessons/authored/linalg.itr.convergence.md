# Convergence rates

A stationary iteration converges linearly when its error shrinks by a fixed factor per step. The factor is the spectral radius of the iteration matrix.

## Rate

After k steps the error is bounded by the spectral radius to the power k times the initial error. The logarithm of the reciprocal gives the asymptotic rate.

## Digits

The number of correct decimal digits gained per step is the negative logarithm base ten of the spectral radius. A radius near one therefore means slow progress.

## Acceleration

Chebyshev acceleration and Krylov methods improve on the plain rate by using information from earlier iterates. They do not change the underlying spectral radius.

To check your understanding, raise the radius to k, take a base-ten logarithm and use earlier iterates.

Linear convergence is an exponential decay, and the spectral radius is its decay constant.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
