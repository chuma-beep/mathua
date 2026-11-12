# Distribution shift

Distribution shift is a change between the training distribution and the deployment distribution. A model that fits the training law can fail when that law moves.

## Covariate shift

The conditional distribution of the label given the input stays fixed while the input distribution changes. Reweighting by the density ratio corrects for it.

## Label shift

The class prior changes while the class-conditional input distribution stays fixed. The correction adjusts the predicted priors.

## The measurement

Total variation distance bounds how much any event's probability can change. A distance of zero means the distributions agree.

To check your understanding, fix the conditional, fix the class conditionals and compute the distance.

Shift names which part of the joint distribution moved.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
