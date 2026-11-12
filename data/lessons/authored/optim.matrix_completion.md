# Matrix completion

Matrix completion fills the missing entries of a low-rank matrix from a few observations. The nuclear norm is the convex stand-in for rank.

## Observed entries

Only some positions are known, and the model must reproduce each of them. The remaining positions are free apart from the structural penalty.

## The nuclear-norm proxy

The nuclear norm is the sum of the singular values, which is the tightest convex lower bound on rank. Penalising it pushes the solution towards few nonzero singular values.

## Recovery conditions

Recovery needs the observations to be spread across the matrix and the rank to be small. When both hold, the nuclear-norm solution equals the true low-rank matrix.

To check your understanding, count the known positions, take absolute values of the diagonal and count the nonzero singular values.

The observed entries are equality constraints and the nuclear norm is the objective. The solution is the simplest low-rank matrix that agrees with what was seen.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
