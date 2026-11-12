# The explicit heat scheme

The explicit heat scheme advances each node from its two neighbours at the old time level. It is simple, but its time step is tied to the square of the spacing.

## The update

\(u_i^{n+1}=u_i^n+r(u_{i+1}^n-2u_i^n+u_{i-1}^n)\) with \(r=\alpha\Delta t/\Delta x^{2}\). The stencil is three points wide.

## Stability

The scheme is stable only when \(r\le 1/2\), that is \(\Delta t\le\Delta x^{2}/(2\alpha)\). Halving the spacing quarters the admissible step.

## Cost

Each step is cheap and parallel, but the step count grows like \(1/\Delta x^{2}\). The total work therefore grows like \(1/\Delta x^{4}\) for a fixed final time.

To check your understanding, write the three-point stencil, bound r by one half and count the steps.

Explicit diffusion is simple to code and expensive to refine.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
