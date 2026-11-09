# The Runge phenomenon

The Runge phenomenon shows that high-degree interpolation on equally spaced nodes can oscillate violently near the endpoints.

## Equally spaced nodes

With \(n+1\) equally spaced nodes there are \(n\) subintervals. The nodes cluster near the center when measured against the endpoints.

The interpolation error depends on the node distribution, not only on the number of nodes.

The two endpoints are always included, and the interior spacing is uniform.

## Endpoint oscillations

For Runge’s function \(1/(1+25x^2)\), the interpolant develops large oscillations near \(x=\pm1\) as the degree grows.

The oscillations can grow even though the function is smooth. Adding more equally spaced nodes can make the maximum error worse.

The effect is not a roundoff error. It comes from the interpolation formula itself.

## Choosing better nodes

Chebyshev nodes cluster near the endpoints and control the oscillation. They make high-degree interpolation stable for smooth functions.

Piecewise interpolation with low-degree pieces is another remedy. It avoids a single high-degree polynomial.

The choice of nodes is a design decision, not an implementation detail.
