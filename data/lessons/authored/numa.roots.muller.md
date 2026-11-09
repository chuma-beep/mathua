# Muller’s method

Muller’s method fits a quadratic through three points and uses one of its roots as the next estimate. It can converge to complex roots.

## Quadratic interpolation

Three points determine a quadratic. Muller’s method writes that quadratic and solves it with the quadratic formula.

The leading coefficient is \((y_0+y_2-2y_1)/(2h^2)\) for equally spaced points. It measures the curvature of the interpolant.

The linear coefficient is \((y_2-y_0)/(2h)\). Together with the constant term it determines the quadratic completely.

## Three starting points

Muller’s method needs three initial guesses, unlike the secant method, which needs two. The three points define the quadratic.

The method chooses the root closest to the current iterate to keep the step stable. The quadratic formula supplies both roots.

A poor triangle of points can make the interpolating quadratic nearly degenerate.

## Complex roots

Because the quadratic may have complex roots, Muller’s method can find complex solutions of real equations. Newton’s method with real arithmetic cannot leave the real line.

The discriminant decides whether the roots are real or complex. A negative discriminant means a complex conjugate pair.

Complex arithmetic makes each iteration more expensive, but it widens the set of reachable roots.
