# Riemann-Stieltjes integrals

## Riemann-Stieltjes sums

A Riemann-Stieltjes sum replaces the width \(\Delta x\) by an increment \(\Delta\alpha\), so it reads \(\sum f(t_i)(\alpha(x_i) - \alpha(x_{i-1}))\). Taking \(\alpha(x) = x\) returns the ordinary Riemann sum. Taking \(\alpha(x) = 2x\) doubles every sum, so the integral doubles as well.

## Step integrators

A step integrator concentrates the whole integral at its jumps. If \(\alpha\) jumps by 2 at \(x = 1\) and by 3 at \(x = 2\), then the integral of \(f\) is \(2f(1) + 3f(2)\). With \(f(x) = x^2\) the value is \(2 \cdot 1 + 3 \cdot 4 = 14\).

## Summing the jumps

Summing the jumps of a staircase is the standard example. On the interval from 0 to 3 with \(\alpha(x) = \lfloor x \rfloor\) the only jumps inside are at 1 and 2, each of size 1. For \(f(x) = x\) the integral is \(f(1) + f(2) = 3\). A jump at an endpoint contributes nothing to the open interval.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
