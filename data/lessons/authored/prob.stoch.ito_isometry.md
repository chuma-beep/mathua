# The Ito isometry

The Ito isometry turns a second moment of a stochastic integral into an ordinary integral of a squared integrand.

## Statement

For adapted \(H\), the expectation of the square of \(\int H\,dB\) equals the expectation of \(\int H^2\,dt\).

## Use

It lets a stochastic integral be treated as an \(L^2\) object, so convergence arguments apply as they do for deterministic integrals.

## A check

For constant \(H=c\) on \([0,t]\), the integral is \(cB_t\), whose variance is \(c^2t\).

To check your understanding, square the integral, place the integral in L2 and take a constant integrand.

The isometry converts stochastic variance into a deterministic integral.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
