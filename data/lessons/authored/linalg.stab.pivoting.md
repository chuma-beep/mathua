# Pivoting strategies

Pivoting reorders rows or columns during elimination to keep multipliers small. Partial pivoting is the default compromise between stability and cost.

## Partial

At each column choose the largest available entry in that column as the pivot. This guarantees multipliers no larger than one in magnitude.

## Complete

Complete pivoting searches the whole remaining submatrix. It gives stronger stability at the cost of more comparisons.

## Permutations

Row interchanges are recorded as a permutation matrix. The factorisation becomes P A equals L U, and the permutation is applied when solving.

To check your understanding, scan the column, scan the submatrix and record the swaps.

Pivoting is a search for a small multiplier, and the search scope is the design choice.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
