# The Fokker-Planck equation

The Fokker-Planck equation describes how the density of a diffusion evolves in time. It is the forward equation dual to the backward generator.

## The equation

The density satisfies \(\partial_t p=-\partial_x(b p)+\tfrac12\partial_{xx}(\sigma^2 p)\), with drift \(b\) and diffusion \(\sigma\).

## Interpretation

The first term transports probability with the drift and the second spreads it with the diffusion.

## Stationary case

Setting the time derivative to zero gives the invariant density, which for constant coefficients is Gaussian.

To check your understanding, write the drift term, read transport and drop the time derivative.

The forward equation moves densities, not paths.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
