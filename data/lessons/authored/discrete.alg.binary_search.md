# Binary search

Binary search halves a sorted interval after each comparison. It finds a key in logarithmic time, provided the invariant that the key can only lie inside the interval is maintained.

## The halving

Compare the key with the middle element. One comparison discards half of the remaining interval, so the interval length follows a geometric sequence.

## The invariant

At every step the key, if present, lies between the low and high bounds inclusive. The loop stops when the bounds cross, which is the not-found case.

## Logarithmic steps

After \(k\) comparisons at most \(n/2^k\) candidates remain, so the loop runs at most \(\lceil\log_2 n\rceil\) times.

To check your understanding, compare with the middle, state the bound invariant and halve the candidate count.

Binary search is the canonical example of an algorithm whose correctness rests on an invariant rather than on the code alone. Off-by-one errors in the bounds are the usual failure.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
