# Checking a dual certificate

A dual certificate is a short proof that a proposed point is optimal. Checking it costs one residual evaluation instead of a full solve.

## What a certificate proves

A feasible dual point with the same value as the primal point proves both are optimal. The certificate is the dual point together with that equality.

## The residual test

Stationarity is a residual that must vanish and a dual feasibility condition that must hold. A nonzero residual or a violated sign rejects the certificate.

## Reading a failure

A small but nonzero residual points at the stationarity row that is wrong. A violated sign points at a constraint whose multiplier has the wrong direction.

To check your understanding, form the stationarity expression, check every multiplier and compare the primal and dual values.

The check is cheap because it evaluates a candidate rather than searching for one. It is the last step of many conic solvers.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
