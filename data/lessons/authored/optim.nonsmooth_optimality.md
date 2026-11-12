# Nonsmooth optimality

When a function has kinks, a zero gradient is not necessary for optimality. The subdifferential replaces the gradient.

## The condition

Zero must belong to the subdifferential at an unconstrained minimum. For a smooth function this reduces to a zero gradient.

## Kinks

At a kink the subdifferential is an interval or a set. An optimum can sit exactly at the kink with no gradient defined there.

## Constrained form

With constraints the condition combines the subdifferential with the normal cone of the feasible set. The sum must contain zero.

To check your understanding, check zero membership, compute a subdifferential and form a normal cone.

Nonsmooth optimality replaces an equation with a set-membership test.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
