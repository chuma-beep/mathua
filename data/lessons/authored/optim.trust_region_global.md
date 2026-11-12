# Global convergence of trust region

A trust-region method is globally convergent because every accepted step lowers the objective or leaves the radius shrinking.

## The ratio test

The ratio compares the actual decrease with the decrease predicted by the model. A ratio near one accepts the step and a ratio near zero rejects it.

## The radius

An accepted step can grow the radius and a rejected step shrinks it. The radius therefore tracks the region where the model is trustworthy.

## Global convergence

If the gradient did not vanish the model would predict a decrease and a small enough radius would realise it. That contradiction forces the gradient to zero in the limit.

To check your understanding, compute the ratio, shrink the radius and argue by contradiction.

Trust region is a local model with a leash, and the leash is what makes it globally safe.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
