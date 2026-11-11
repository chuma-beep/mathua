# Strict convexity and uniqueness

Strict convexity strengthens the midpoint condition to a strict inequality for distinct points. Its practical value is uniqueness: a strictly convex function on a convex set has at most one minimizer. Many objectives are convex but not strictly convex, and the difference shows up as a flat valley or a whole segment of minimizers.

## Testing strict convexity

For a twice differentiable function of one variable, strict convexity follows from a positive second derivative everywhere, and a zero second derivative at isolated points is allowed. For a quadratic with matrix, the test is positive definiteness rather than positive semidefiniteness. A single pair of distinct points with equality in the midpoint inequality is enough to show the function is not strictly convex.

## Counting minimizers

If the function is strictly convex and a minimizer exists, there is exactly one. If it is merely convex, the minimizer set is convex, so it can be a segment or a face. A constant function on an interval is the extreme case, with every point minimal. Counting minimizers therefore starts with the definiteness of the curvature, not with the solver output.

## The midpoint inequality

The midpoint form says the value at the average is below the average of the values, with equality allowed only for convexity and forbidden for distinct points under strict convexity. Rearranging it gives a numerical test: evaluate the function at two points and at their midpoint, then compare. The size of the gap measures how curved the function is over that interval and is the quantity a line search implicitly exploits.

Uniqueness is what lets a solver stop with confidence and what lets a proof say the solution rather than a solution. Strict convexity is the usual hypothesis behind that word.
