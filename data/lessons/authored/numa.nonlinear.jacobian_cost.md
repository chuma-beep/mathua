# Jacobian assembly

Assembling a Jacobian can dominate the cost of a Newton step. Finite differences need one perturbation per unknown, while automatic differentiation pays a constant factor.

## Finite differences

A forward difference needs \(n+1\) evaluations for n unknowns, or \(2n\) for a central difference. The accuracy is \(O(h)\) or \(O(h^{2})\).

## Automatic differentiation

Forward mode propagates one direction per pass, so a full Jacobian costs n passes. Reverse mode gives one row per pass, which suits many residuals and few unknowns.

## Structure

A sparse Jacobian needs far fewer perturbations when columns share a sparsity pattern. Graph colouring groups compatible columns.

To check your understanding, count forward evaluations, use forward mode and colour the columns.

The cheapest Jacobian is the one that exploits the structure of the residual.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
