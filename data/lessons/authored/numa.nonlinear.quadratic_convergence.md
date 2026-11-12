# Quadratic convergence

An iteration converges quadratically when the error at the next step is proportional to the square of the current error. The number of correct digits doubles each step.

## The definition

\(e_{k+1}\le Ce_k^{2}\) near the limit. The constant \(C\) depends on the second derivative and on the conditioning of the Jacobian.

## Digit count

If the error starts at \(10^{-1}\) and \(C\) is about one, the sequence is \(10^{-1},10^{-2},10^{-4},10^{-8}\), so six steps reach machine precision.

## Failure modes

A multiple root drops the order to linear. An inaccurate Jacobian also destroys the quadratic term and leaves superlinear convergence at best.

To check your understanding, write the error bound, double the digits and recognise a multiple root.

Doubling digits per step is the signature of a correct Newton iteration.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
