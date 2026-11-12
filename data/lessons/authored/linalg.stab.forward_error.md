# Forward error

Forward error measures the distance between a computed answer and the exact answer. It is the quantity a user cares about, but it is usually the hardest to bound.

## Absolute

The absolute forward error is the norm of the difference between the computed and exact vectors. It carries the units of the answer.

## Relative

The relative forward error divides by the norm of the exact answer. It is dimensionless and is the natural target when the scale is unknown.

## Bound

The relative error is bounded by the condition number times the relative backward error. This factorisation separates the problem from the algorithm.

To check your understanding, subtract the vectors, divide by the answer and multiply by the condition number.

Forward error is what the user sees, and the condition number explains why it can be large.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
