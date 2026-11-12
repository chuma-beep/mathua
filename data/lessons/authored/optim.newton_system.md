# The interior-point Newton system

The Newton system of an interior-point method is the linear system solved each iteration. Its matrix is symmetric indefinite and its right-hand side holds the residuals.

## The block form

With equality constraints the matrix is a Hessian block, a constraint block and a zero corner. The Hessian block carries a diagonal term from the barrier that keeps it positive definite on the feasible subspace.

## The Schur complement

Eliminating the primal block leaves the constraint block times the inverse Hessian times its transpose. The result is a smaller system for the multipliers, often the cheaper solve.

## The conditioning

As the parameter falls the diagonal term grows and the system becomes ill-conditioned. Regularisation or a different elimination order keeps the factorisation usable.

To check your understanding, identify the blocks, eliminate the primal block and compute a diagonal term.

The Newton system is where the cost of an interior-point method lives.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
