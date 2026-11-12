# Complex-step differentiation

The complex step evaluates a real function at a complex argument and reads the derivative from the imaginary part. It has no subtractive cancellation, so the step can be extremely small.

## The formula

For an analytic function, the imaginary part of the value at a point shifted by an imaginary step divided by the step equals the derivative to second order. No difference of two nearly equal real numbers appears.

## The step size

The step can be as small as a fraction of machine precision, because cancellation is absent. The accuracy is then limited by the second-order truncation term rather than by rounding.

## The requirement

The function must be analytic and the code must accept complex arguments. Branch cuts, absolute values and comparisons break analyticity, so the technique applies to a restricted but useful class.

To check your understanding, shift by an imaginary step, avoid cancellation and require analyticity.

Complex-step differentiation buys near machine-precision derivatives at the price of a complex-aware implementation.
