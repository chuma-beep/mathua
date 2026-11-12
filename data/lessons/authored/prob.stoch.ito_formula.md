# The Ito formula

The Ito formula is the chain rule for stochastic integrals. A second-derivative term appears because quadratic variation does not vanish.

## Statement

For \(f\) twice differentiable, \(df(B_t)=f'(B_t)dB_t+\tfrac12 f''(B_t)dt\). The extra term is the Ito correction.

## The correction

The correction integrates the second derivative against the quadratic variation. It is what keeps \(B_t^2-t\) a martingale.

## Worked check

Taking \(f(x)=x^2\) gives \(d(B_t^2)=2B_t dB_t+dt\), so \(B_t^2-t\) has zero drift.

To check your understanding, write the first derivative term, integrate against time and apply f(x)=x squared.

The Ito formula is the chain rule with a drift correction.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
