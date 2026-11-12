# Crowding distance

Crowding distance measures how isolated a point is along a front. Selection uses it to keep the front spread out.

## Neighbours

Sort the front by one objective and look at the points on either side. The gap to the neighbours is the contribution of that objective.

## The distance

Sum the normalised neighbour gaps over the objectives. A large distance means the point sits in a sparse region of the front.

## Boundary points

The smallest and largest point in each objective receive an infinite distance. They are always kept, because they anchor the extent of the front.

To check your understanding, sort by an objective, normalise gaps and identify an extreme.

Crowding distance is diversity measured on the front itself.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
