# The alternating direction step

The alternating direction method splits a variable into two copies and updates one at a time. Each update is a simpler problem than the joint one.

## The split variable

The original variable is replaced by two copies tied by an equality. The objective keeps one copy and the constraint keeps the other.

## The two updates

The first copy minimises the objective plus the penalty, and the second enforces the constraint plus the penalty. Alternating them costs one solve each.

## Stopping

The primal residual measures the disagreement between the copies and the dual residual measures the change in the multiplier. Both must be small.

To check your understanding, fix the second copy, fix the first copy and subtract the copies.

Each subproblem is often closed form. That is why the method scales even when the joint problem does not.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
