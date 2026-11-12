# Backward error

Backward error asks which nearby problem the computed answer solves exactly. A small backward error means the algorithm is stable, whatever the forward error.

## Definition

The backward error is the smallest perturbation of the data for which the computed answer is exact. It is measured relative to the data scale.

## Stability

An algorithm is backward stable when its backward error is of the order of machine epsilon. Gaussian elimination with partial pivoting has this property in practice.

## Transfer

Forward error is bounded by the condition number times backward error. A stable algorithm on an ill-conditioned problem can still return a poor answer.

To check your understanding, perturb the data, compare with epsilon and multiply by the condition.

Backward error attributes rounding to the problem, and condition attributes the rest to the problem's sensitivity.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
