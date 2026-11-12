# Preconditioners

A preconditioner approximates the matrix well enough that an iterative solver converges in few steps, yet stays cheap to apply. It changes the conditioning of the iteration without changing the solution.

## Left and right

A left preconditioner acts on the residual; a right one acts on the update. Symmetric problems need a split preconditioner to preserve symmetry.

## Jacobi

The diagonal preconditioner divides each residual entry by the corresponding diagonal entry. It costs one vector operation and removes simple scaling.

## Quality measure

The condition number of \(M^{-1}A\) governs the iteration count. Reducing that number by a factor of four roughly halves the number of steps.

To check your understanding, choose the side, read the diagonal and form the preconditioned matrix.

A preconditioner is judged by the conditioning it leaves behind and the cost it adds.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
