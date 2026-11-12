# Goal programming

Goal programming measures how far a solution falls from targets and minimises those deviations. It is a multiobjective method built from preferences.

## Goals and deviations

Each objective has a target and two deviation variables, one for falling short and one for overshooting. The objective is a weighted sum of the deviations.

## Weighted deviations

The weights encode the relative cost of missing each target. A weight of zero turns the target into a constraint that may be missed for free.

## Lexicographic order

A lexicographic version minimises the largest deviation first, then the next, which prevents a large miss from being hidden by many small wins.

To check your understanding, write a deviation, apply weights and rank deviations.

Goal programming is the language of targets, and the weights are the policy.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
