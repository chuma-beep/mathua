# The backdoor criterion

The backdoor criterion names the covariate sets that identify a causal effect by adjustment. A valid set blocks every backdoor path and contains no descendant of the treatment.

## Blocking

Every path from treatment to outcome that starts with an arrow into the treatment must be blocked. Conditioning on a non-collider on the path blocks it.

## Descendants

A descendant of the treatment carries part of the effect. Conditioning on it removes part of the effect as well as the confounding.

## Minimal sets

A minimal sufficient set has no redundant member. Removing any member leaves a backdoor path open.

To check your understanding, find a backdoor path, identify descendants and remove a member.

The backdoor criterion is the graph-theoretic test for a valid adjustment set.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
