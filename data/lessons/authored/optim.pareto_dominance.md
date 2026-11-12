# Pareto dominance

Pareto dominance compares two feasible points without choosing a trade-off. It is the partial order that multiobjective optimisation works with.

## The definition

One point dominates another when no objective is worse and at least one objective is strictly better. A dominated point can be discarded without any preference between objectives.

## Weak dominance

Weak dominance allows every objective to be equal. Strict dominance requires an improvement somewhere, which is the form used to define a front.

## A test on two points

Compare the vectors coordinatewise and count the wins and losses. A loss anywhere rules out dominance, and a tie everywhere rules out strict dominance.

To check your understanding, compare two objective vectors, distinguish weak from strict and count winning coordinates.

Dominance is the only comparison that needs no weights and no units.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
