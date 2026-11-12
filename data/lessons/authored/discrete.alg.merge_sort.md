# Merge sort

Merge sort splits the array in half, sorts each half, and merges the two sorted runs. Every case costs linearithmic time and the merge is stable.

## Divide and merge

The array is halved until runs have length one. A linear merge then combines two sorted runs by repeatedly taking the smaller front element.

## The recurrence

Splitting into equal halves gives \(T(n)=2T(n/2)+O(n)\), which the master theorem resolves to \(O(n\log n)\).

## Stability

Taking the left element on ties preserves the original order of equal keys, so merge sort is stable without extra comparisons.

To check your understanding, halve the array, write the halving recurrence and break ties to the left.

Merge sort trades the extra array of a merge for a guaranteed bound. That guarantee is why it is the fallback when quicksort's worst case is unacceptable.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
