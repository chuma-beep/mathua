# Path-following methods

A path-following method takes damped Newton steps to stay near the central path while the barrier parameter falls. A neighbourhood condition bounds how far a point may stray from the path.

## The neighbourhood

A point lies in the neighbourhood when every slack-multiplier product is between one minus theta and one plus theta times mu. The threshold theta is small, often a tenth.

## The damped step

A full Newton step can leave the neighbourhood. The method scales the step by the largest factor that keeps every product inside the bounds.

## The reduction

After a step the parameter becomes sigma times mu with sigma below one. The gap falls by roughly the same factor, so the iteration count grows with the logarithm of the required accuracy.

To check your understanding, compute a lower product bound, compare a product with the bound and apply the reduction factor.

Path following is the practical face of the central path: a neighbourhood, a damped Newton step and a fixed reduction.
