# Smoothing a dual problem

A dual problem is often nonsmooth because a maximum over discrete alternatives creates a kink. Adding a strongly convex term smooths it into a differentiable problem.

## Why the dual is rough

The dual is a minimum over multipliers of a maximum over primal variables. Swapping the order leaves a maximum of linear functions, which has kinks.

## The smoothing term

A small strongly convex term makes the smoothed dual strictly convex and differentiable. The smoothing parameter controls how far it moves.

## The temperature

The smoothing term acts like a temperature in statistical mechanics. A high temperature blurs the kinks and a low one keeps them sharp.

To check your understanding, evaluate the rough dual, compute the probabilities and compare the smoothed and rough values.

The smoothed optimum is close to the true optimum with an error of order the smoothing parameter. The gradient is Lipschitz with a constant set by the same parameter.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
