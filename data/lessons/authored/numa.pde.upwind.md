# Upwind differencing

Upwind differencing takes the spatial derivative from the side the flow comes from. It keeps the scheme stable at the cost of artificial diffusion.

## The stencil

For positive speed the backward difference is used, so the stencil leans upstream. For negative speed the forward difference is used.

## Numerical diffusion

The leading truncation term is \(\tfrac{c\Delta x}{2}(1-C)u_{xx}\), a diffusion proportional to the cell size. It smears fronts even when the physical problem has none.

## Accuracy

The scheme is first order. Higher-order upwind reconstructions reduce the smearing but can produce oscillations near discontinuities.

To check your understanding, read the flow direction, identify the truncation term and state the order.

Upwind leans into the flow and pays for stability with smearing.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
