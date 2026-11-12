# Infeasibility certificates

An infeasibility certificate is a direction that separates the constraints from every candidate point. Farkas' lemma gives the linear-programming form.

## Farkas' lemma

Exactly one of two statements holds. Either the system with a nonnegative variable has a solution, or there is a multiplier whose transposed product with the columns is nonnegative and whose product with the right-hand side is negative.

## Reading the certificate

A negative right-hand product while every column product is nonnegative contradicts every feasible point. The inner product would be both nonnegative and negative.

## Numerical caution

Interior-point solvers approach infeasibility with a diverging certificate. Scaling decides whether the returned vector is a useful proof.

To check your understanding, compute the right-hand product, compute a column product and scale a certificate.

An infeasibility certificate is a short proof that no feasible point exists.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
