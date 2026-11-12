# The trapezoidal rule

The trapezoidal rule averages the derivative at both ends of a step. It is second order and A-stable, but it is not L-stable, so it can leave a visible oscillation on a very stiff mode.

## The formula

The new value equals the old value plus half the step times the sum of the old and new derivatives. The new derivative makes the formula implicit, so each step requires a solve.

## Stability function

On the scalar test problem the function is a ratio of two linear factors. Its modulus is at most one on the whole left half-plane, which is the A-stability proof.

## Behaviour at infinity

The ratio tends to minus one rather than zero. A very stiff mode therefore flips sign each step and decays only slowly, which is the ringing that L-stable methods avoid.

To check your understanding, write the average of the two derivatives, apply the test problem and take the limit.

The trapezoidal rule is the standard example of a method that is stable without being sufficiently damped.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
