# Multigrid

Multigrid attacks smooth error on a coarse grid and oscillatory error on the fine grid. One cycle reduces error at every frequency, so the iteration count barely grows with the grid size.

## Smoothing

A few Jacobi or Gauss-Seidel sweeps remove high-frequency error quickly. The remaining error is smooth and can be represented on a coarser grid.

## The cycle

Restriction transfers the residual down, a coarse solve corrects it and prolongation interpolates the correction back. The V-cycle visits each level once.

## Complexity

With geometric coarsening the number of levels is \(\log_2 n\) and the work per cycle is proportional to the number of unknowns. That is optimal for a linear solver.

To check your understanding, damp the high frequencies, restrict the residual and count the levels.

Multigrid is the rare method whose cost per unknown does not grow with the grid.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
