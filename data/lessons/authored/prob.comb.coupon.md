# The coupon collector

The coupon collector collects one of \(n\) equally likely coupons per draw. The expected time to complete the set grows like \(n\log n\).

## Expected time

The waiting time for the \(k\)-th new coupon is geometric with success probability \((n-k+1)/n\), so the expectations sum to \(nH_n\).

## Tail

Concentration around the harmonic sum follows from independence of the waiting times.

## Variants

Unequal coupon probabilities change the constants but not the logarithmic order.

To check your understanding, wait for a new coupon, use independence and unequal probabilities.

Collecting all coupons takes about n log n draws.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
