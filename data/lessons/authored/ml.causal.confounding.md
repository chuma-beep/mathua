# Confounding

A confounder influences both the treatment and the outcome. Ignoring it mixes the treatment effect with the effect of the confounder.

## The backdoor path

A confounder opens a path from treatment to outcome that does not pass through the treatment effect. That path must be blocked.

## Adjustment

Conditioning on a valid adjustment set blocks every backdoor path. The set must not contain descendants of the treatment.

## The bias

The unadjusted difference equals the true effect plus the confounding bias. The bias vanishes only when the confounder is balanced across groups.

To check your understanding, draw the path, choose a set and compare the difference.

Confounding is a path problem before it is a modelling problem.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
