# Counting polygon diagonals

## Counting diagonals

A diagonal joins two nonadjacent vertices of a polygon. From each vertex there are n - 3 diagonals, and each diagonal is counted twice, so the total is n(n - 3)/2.

The formula applies to convex polygons and to simple polygons counted combinatorially. It counts segments, not crossings.

## Recovering the side count

If the number of diagonals is known, solve n(n - 3)/2 = D. The positive integer solution is the number of sides. A quick estimate places n between the square root of 2D and that value plus 3.

Check the solution by substituting it back into the formula. A non-integer solution means the given diagonal count does not belong to a polygon.

## All segments

The total number of segments joining pairs of vertices is n(n - 1)/2. This includes the n sides and all diagonals. Subtracting n sides leaves the diagonal count.

The two formulas are consistent. Use the total-segment formula as a check on the diagonal formula.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
