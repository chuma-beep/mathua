# Adam

Adam combines momentum with a per-coordinate second-moment estimate and bias correction.

## The two averages

The first moment averages gradients and the second averages their squares. Both use exponential weights, so recent values dominate.

## Bias correction

Both averages start at zero, which biases them low. Dividing by one minus the decay raised to the step count removes the bias in the first steps.

## The update

The step is the base step times the corrected first moment divided by the square root of the corrected second moment. The ratio is free of the gradient scale.

To check your understanding, update a moment, compute the correction and form the ratio.

Adam is momentum with a normaliser, and the normaliser is what lets one step size work across coordinates.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
