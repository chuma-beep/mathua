# Star polygons

A star polygon connects every \(k\)-th vertex of a regular \(n\)-gon. When \(k\) and \(n\) share a factor, the figure splits into several cycles.

## The step

The notation \(\{n/k\}\) records \(n\) vertices joined in steps of \(k\). The step is taken modulo \(n\).

## The cycles

The number of separate cycles is the greatest common divisor of \(n\) and \(k\).

## The turning angle

Each edge turns by \(2\pi k/n\), so a single closed cycle needs \(n/\gcd(n,k)\) edges.

To check your understanding, read the notation, take the gcd and compute the turn.

A star polygon is determined by its vertex count and its step.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
