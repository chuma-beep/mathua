# The subgradient method

A subgradient method replaces the gradient with any subgradient. The step must shrink over time because a subgradient is not a descent direction.

## The step rule

The update moves against a chosen subgradient scaled by the step. Unlike gradient descent, the new point can have a larger objective value.

## Diminishing steps

A summable step sequence with a divergent sum of squares guarantees convergence. A constant step converges only to a neighbourhood of the optimum.

## The convergence rate

The best objective after k steps is within a constant over the square root of k. The rate is slower than smooth gradient descent, which is the cost of the weaker assumption.

To check your understanding, choose the step rule, pick a subgradient and take the square root of the iteration count.

The method is robust and cheap, one subgradient per iteration. It is the fallback when no smooth structure is available.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
