# Proximal operators

A proximal operator returns the point that balances a small step on a function against closeness to a given point.

## The definition

The operator minimises the function plus one over twice a step times the squared distance to the input. It exists even when the function is not differentiable.

## A quadratic example

For a quadratic the operator is a scaled average of the input and the minimiser. The step decides how far the result moves.

## Why it helps

The operator handles the nonsmooth part exactly, so the smooth part can be graded separately. This is the basis of proximal gradient methods.

To check your understanding, write the proximal objective, apply the formula and split smooth and nonsmooth.

A proximal operator is a small optimisation solved exactly inside a larger iteration.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
