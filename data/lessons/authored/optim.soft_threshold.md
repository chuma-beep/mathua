# Soft thresholding

Soft thresholding shrinks a value toward zero and sets small values to zero. It is the proximal operator of the absolute value.

## The rule

Subtract the threshold from the magnitude, keep the sign and clamp at zero. A value inside the threshold becomes exactly zero.

## The threshold

The threshold is the penalty times the step. A larger penalty removes more values and leaves fewer nonzeros.

## Two-sided

The rule is applied to positive and negative inputs alike. Its graph is flat near zero and shifts the diagonal outside the flat region.

To check your understanding, apply the rule, compute a threshold and handle a negative input.

Soft thresholding is the arithmetic that makes the lasso sparse.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
