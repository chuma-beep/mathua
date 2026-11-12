# The Moreau envelope

The Moreau envelope smooths a function by minimising it plus a distance penalty. The smoothed function shares the minimisers of the original.

## The envelope

The envelope at a point is the smallest value of the function plus the scaled squared distance. It is always at least as small as the function itself.

## The gradient

The gradient of the envelope is the scaled displacement from the point to its proximal image. That identity is what makes the envelope useful for optimisation.

## The smoothing effect

A larger smoothing parameter gives a flatter envelope. The original function is recovered in the limit as the parameter goes to zero.

To check your understanding, compute the proximal image, subtract the proximal image from the point and evaluate the original at the point.

The envelope is convex whenever the original function is convex. Its gradient need not exist for the original, which is the point of the construction.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
