# The implicit heat scheme

The implicit heat scheme evaluates the spatial operator at the new time level. It is unconditionally stable, at the price of solving a linear system at every step.

## The system

The update is \((I-rL)u^{n+1}=u^{n}\), where \(L\) is the second-difference matrix. In one dimension the matrix is tridiagonal and the Thomas algorithm costs \(O(n)\).

## Stability

The amplification factor is \(1/(1+4r\sin^{2}(\theta/2))\), whose magnitude never exceeds one. Any positive step is stable.

## Accuracy

The scheme is first order in time and second order in space, and large steps smear sharp features. Stability does not imply accuracy.

To check your understanding, form the implicit matrix, write the amplification factor and state the order in time.

Implicit diffusion removes the step limit and adds a linear solve.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
