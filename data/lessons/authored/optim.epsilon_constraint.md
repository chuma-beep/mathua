# The epsilon-constraint method

The epsilon-constraint method keeps one objective and turns the others into bounds. Sweeping the bounds traces the front, including nonconvex parts.

## The idea

One objective is minimised while every other objective is constrained below a bound. The result is Pareto optimal when the bound is tight.

## Choosing epsilon

A tight bound can make the problem infeasible. A loose bound returns a dominated point, so the useful range lies between the extremes.

## Feasibility

Each sweep solves a constrained problem, and feasibility is checked before optimality. An infeasible bound is information about the front rather than a failure.

To check your understanding, fix one objective, test a tight bound and check the bounds.

The epsilon-constraint method trades a weight sweep for a sequence of constrained solves.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
