# Epigraphs and sublevel sets

The epigraph of a function is the set of points above its graph, and the sublevel set is the set of points where the function is at most a value. Both turn questions about a function into questions about a set. Convexity of the function is exactly convexity of the epigraph, which is the cleanest bridge between the two languages.

## Testing an epigraph point

A point belongs to the epigraph when its height is at least the function value at its horizontal coordinate. The test is one evaluation and one comparison. The epigraph of a convex function is a convex set, so a point below the graph is outside the epigraph and a segment between two points above the graph stays above it.

## Finding a sublevel set

The sublevel set at a value is found by solving the inequality that defines it. For a convex quadratic, the sublevel set is an ellipsoid whose axes come from the eigenvalues of the curvature matrix. For a linear function it is a half-space. For the absolute value it is an interval. Solving the inequality once gives the whole set, not just one point.

## Convexity through the epigraph

A function is convex exactly when its epigraph is a convex set. This equivalence is useful because set convexity has a direct midpoint test, and because it extends convexity to functions that take infinite values outside a domain. A vertical section of the epigraph is the sublevel set, so the two objects carry the same information in different directions.

The epigraph is the geometric form of convexity. It explains why sublevel sets of convex functions are convex and gives a single picture that covers constraints, objectives and duality.
