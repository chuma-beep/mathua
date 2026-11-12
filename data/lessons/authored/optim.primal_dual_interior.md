# Primal-dual interior-point methods

Primal-dual methods treat the primal point, the slacks and the multipliers as unknowns and drive the complementarity products to the parameter together.

## The residuals

Three residuals measure stationarity, feasibility and complementarity. A step must reduce all three, and the complementarity residual is the one the parameter controls directly.

## The Newton system

Linearising the residuals gives a system whose matrix is the KKT matrix with a diagonal complementarity block. The block changes each iteration, so the factorisation is rebuilt each time.

## The step length

A fraction-to-the-boundary rule keeps slacks and multipliers positive. The same length applies to both, which is why the iterate stays inside the cone.

To check your understanding, evaluate the complementarity residual, count the unknowns and respect positivity.

A primal-dual method is Newton's method on a perturbed KKT system, with the perturbation driven to zero.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
