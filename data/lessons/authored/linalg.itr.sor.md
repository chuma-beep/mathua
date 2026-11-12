# Over-relaxation parameters

Successive over-relaxation extrapolates the Gauss-Seidel correction by a parameter omega. A good omega can cut the iteration count substantially.

## Parameter

The update is a convex combination of the old value and the Gauss-Seidel value when omega lies between zero and one. Values above one extrapolate beyond the correction.

## Optimal value

For consistently ordered matrices the optimal omega is two divided by one plus the square root of one minus the Jacobi spectral radius squared.

## Range

Convergence requires omega strictly between zero and two. Outside that interval the iteration diverges even for well-conditioned systems.

To check your understanding, form the combination, use the spectral radius and bound omega.

SOR is Gauss-Seidel with a tunable extrapolation, and the tuning is a one-dimensional optimisation.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
