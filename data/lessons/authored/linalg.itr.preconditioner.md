# Preconditioned iterations

A preconditioner is an approximate inverse applied to a system before iteration. A good one clusters the spectrum and reduces the iteration count.

## Construction

The preconditioner should be cheap to apply and close to the inverse of the matrix. Diagonal scaling, incomplete factorisations and multigrid are common choices.

## Spectrum

Preconditioning changes the eigenvalue distribution of the iteration operator. Clustering eigenvalues near one is the goal.

## Trade-off

A stronger preconditioner costs more per step but needs fewer steps. The best choice minimises total work, not either count alone.

To check your understanding, approximate the inverse, cluster eigenvalues and count work per step.

Preconditioning is the practical difference between a textbook Krylov method and a solver.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
