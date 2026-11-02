# Estimating the sum of a series

## The alternating remainder



For an alternating series the remainder is bounded by the next term.

The bound is tight enough for numerical estimates.

Stop when the next term falls below the requested accuracy.

The next term is a convenient error bound because the alternating signs force the true sum to lie between consecutive partial sums. The estimate is therefore both simple and reliable.

## The integral remainder



The integral test gives upper and lower bounds for a sum.

The remainder lies between two adjacent integrals.

The bounds improve as the terms decrease faster.

The integral test converts a discrete tail into an area that can be computed. The upper and lower bounds differ by the first omitted term, so the interval shrinks as N grows.

## Choosing enough terms



Decide the accuracy first.

Use the remainder bound to solve for the number of terms.

Report the estimate with its error bound.

Accuracy requirements determine how many terms are needed, not the other way round. Report the bound alongside the estimate so the reader knows how much to trust the digits.
