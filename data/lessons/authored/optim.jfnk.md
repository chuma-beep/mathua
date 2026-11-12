# Jacobian-free Newton-Krylov

A Newton step needs the action of the Jacobian on a vector, not the Jacobian itself. A Krylov method needs only that action, and a finite difference of the residual supplies it without forming the matrix.

## The product

The Jacobian-vector product is approximated by a directional difference of the residual. The step in that difference is small and scaled to the iterate, and no matrix is ever stored.

## The forcing term

The Krylov solve is run to a tolerance that tightens as the outer iteration converges. An inexact Newton step is sufficient near the solution, and the forcing term prevents oversolving far from it.

## Preconditioning

Without a preconditioner the Krylov iteration may stall. A cheap approximate inverse is built from a simplified Jacobian, and its quality decides whether the method is practical.

To check your understanding, difference the residual, tighten with convergence and build a cheap inverse.

Jacobian-free Newton-Krylov replaces a factorisation with a sequence of residual evaluations, and a preconditioner is what keeps the sequence short.
