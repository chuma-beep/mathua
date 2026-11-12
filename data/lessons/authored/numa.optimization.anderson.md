# Anderson acceleration

Anderson acceleration mixes a fixed-point iteration with a short history of earlier residuals. The mixing often turns a slowly converging iteration into a fast one without any derivative.

## The history

The method stores the last \(m\) residuals and solves a small least-squares problem for the mixing coefficients. A history of depth \(m\) costs \(O(m^{2})\) work per step.

## The mixing step

A damping factor \(\beta\) between zero and one interpolates between the plain iterate and the accelerated one. Setting \(\beta=1\) applies the full correction and \(\beta=0\) recovers the original iteration.

## The saving

On a contraction with rate \(\rho\), plain iteration needs about \(k\) steps to cut the error by a factor 10. Acceleration usually removes a fixed fraction of those steps rather than changing the rate itself.

To check your understanding, count stored residuals, weigh the correction and count the steps saved.

The method accelerates a fixed-point map and never needs the Jacobian.
