# Non-dominated sorting

Non-dominated sorting assigns each point a front number. The first front is the Pareto front, and later fronts are the best of what remains.

## Fronts

Remove the first front and sort the rest again. Repeating the removal partitions every point by how many layers separate it from the front.

## The first front

A point is in the first front when no other point dominates it. Ties are allowed, so the front can hold points that do not dominate each other.

## Later fronts

A point in the second front is dominated only by points in the first. The front number is therefore a measure of how far a point is from optimality.

To check your understanding, remove a front, identify front one and assign a front number.

Sorting into fronts turns a partial order into a ranking that a selection operator can use.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
