# Approximation algorithms

An approximation algorithm returns a feasible solution whose value is within a proven factor of the optimum. The factor is the price paid for polynomial running time on a hard problem.

## The ratio

For minimisation, a \(\rho\)-approximation satisfies \(\mathrm{cost}\le\rho\cdot\mathrm{OPT}\). For maximisation the inequality reverses.

## Schemes

A PTAS achieves ratio \(1+\epsilon\) in time polynomial in the input for fixed \(\epsilon\). An FPTAS is additionally polynomial in \(1/\epsilon\).

## Hardness of approximation

Some problems admit no constant-factor approximation unless P equals NP. The PCP theorem supplies many such thresholds.

To check your understanding, compare with the optimum, fix epsilon and name the threshold.

Approximation is a contract: a proven ratio in polynomial time. The ratio is a worst-case promise, not an average-case observation.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
