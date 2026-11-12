# Broyden's method

Broyden's method updates an approximate Jacobian by a rank-one correction instead of recomputing it. It keeps superlinear convergence at a fraction of the cost per step.

## The update

\(B_{k+1}=B_k+\frac{(y_k-B_ks_k)s_k^{T}}{s_k^{T}s_k}\), where \(s_k\) is the step and \(y_k\) the change in the residual. The new matrix satisfies the secant condition.

## Cost

One matrix-vector product and one rank-one update per step, so the step costs \(O(n^{2})\). A dense Newton step costs \(O(n^{3})\).

## Trade-off

Convergence is superlinear rather than quadratic. The method is attractive when the Jacobian is expensive and n is moderate.

To check your understanding, form the residual change, count the products and state the convergence order.

Broyden replaces a derivative computation with a rank-one guess.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
