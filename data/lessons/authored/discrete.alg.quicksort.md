# Quicksort

Quicksort partitions an array around a pivot, then sorts the two sides recursively. Its expected running time is linearithmic, and its worst case is quadratic.

## Partitioning

A partition pass moves every element smaller than the pivot to its left and every larger element to its right. The pivot then sits in its final position.

## Expected cost

With random pivots the recurrence is \(T(n)=2T(n/2)+O(n)\), giving \(O(n\log n)\) expected time. Balanced splits are typical, not guaranteed.

## Worst case

An already sorted array with a first-element pivot splits off one element per pass, so the recurrence degenerates to \(T(n)=T(n-1)+O(n)\).

To check your understanding, move smaller elements left, write the recurrence and identify the degenerate split.

Quicksort is fast in expectation because random pivots balance the recursion with high probability. The deterministic worst case is why production implementations randomise or sample three pivots.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
