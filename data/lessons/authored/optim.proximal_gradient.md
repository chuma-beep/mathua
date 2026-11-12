# Proximal gradient

Proximal gradient takes a gradient step on the smooth part and applies the proximal operator of the nonsmooth part.

## The iteration

The step first moves along the negative gradient of the smooth term. The proximal operator then handles the penalty.

## The step size

A step below one over the Lipschitz constant guarantees descent. A larger step can diverge even when the proximal part is exact.

## Convergence

With a suitable step the objective decreases at every iteration and the iterates converge to a minimiser. The rate depends on the strong convexity of the smooth part.

To check your understanding, take the gradient step, bound the step and check descent.

Proximal gradient is gradient descent that leaves the differentiable world only once per step.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
