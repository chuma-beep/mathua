# Per-coordinate step sizes

A single global rate must accommodate the largest gradient in any coordinate. Per-coordinate rates let each coordinate move at its own scale.

## The diagonal scaling

Each coordinate divides its gradient by a coordinate-specific statistic, usually a root-mean-square. The resulting step is invariant to a rescaling of that coordinate's features.

## The condition number

Diagonal scaling improves the effective condition number when the Hessian is dominated by its diagonal. It cannot fix rotation, where the large curvature lies off the diagonal.

## The cost

The method stores one statistic per coordinate, which is the same order as the parameters. The arithmetic per step is linear, so the overhead is small compared with the gradient itself.

To check your understanding, scale each coordinate, improve the diagonal case and store one statistic.

Per-coordinate scaling is the reason adaptive methods work well on sparse features, where most coordinates are updated rarely and a global rate would be set by the few dense ones.
