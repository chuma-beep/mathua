# Componentwise stability

Componentwise analysis bounds each entry of the error separately. It is sharper than a norm bound when the data have widely different scales.

## Perturbation

Componentwise backward error perturbs each data entry by a relative amount. The bound is expressed with absolute values entry by entry.

## Condition

The componentwise condition number divides the solution sensitivity by the data perturbation, again entry by entry. It can differ greatly from the norm condition number.

## Scaling

Row and column scaling change the componentwise condition number. Choosing a good scaling is part of formulating a stable problem.

To check your understanding, perturb entries, divide sensitivities and scale rows.

Componentwise analysis respects the units of each entry, which a single norm discards.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
