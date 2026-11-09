# Gauss-Kronrod quadrature

A Gauss-Kronrod pair adds nodes to a Gauss rule so the same evaluations also give an error estimate.

## The embedded pair

A 7-15 Gauss-Kronrod rule uses 15 nodes. Seven of them are the original Gauss nodes.

The 15-node Kronrod rule gives a high-order estimate, while the 7-node subset gives a lower-order estimate.

The two estimates bracket the true integral closely when the function is smooth.

## The error estimate

The difference between the Kronrod and Gauss estimates estimates the quadrature error. A large difference signals that the interval needs refinement.

The two rules share function values, so the error estimate costs only the eight new Kronrod evaluations.

The estimate is not a rigorous bound, but it is a reliable practical indicator.

## Composite cost

With \(n\) panels, each panel uses 15 evaluations, so the total is \(15n\). Panels do not share nodes unless they meet at an endpoint.

Adaptive quadrature uses the error estimate to split only the panels that need it.

The cost per panel is higher than for Simpson’s rule, but the error estimate comes for free.
