# Dual certificates

A dual certificate is a dual feasible point whose value proves a bound on the primal optimum. When the bound meets a primal value, the certificate proves optimality.

## The bound

Weak duality says every dual feasible value bounds a minimisation from below. A certificate is therefore a checkable proof of a bound rather than a claim about an algorithm.

## The gap

If a primal feasible point and a dual certificate have the same value, both are optimal. The gap is the only quantity that needs to be checked.

## Certificates in practice

Solvers return a certificate with the solution. A certificate that is slightly infeasible is a numerical warning rather than a proof.

To check your understanding, evaluate the dual objective, compute the gap and test a certificate.

A certificate converts an answer into a check: plug in, evaluate and compare.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
