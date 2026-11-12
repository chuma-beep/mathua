# Drop tolerances

An incomplete factorisation keeps only entries that pass a drop tolerance, trading exactness for sparsity. The result is applied as a preconditioner rather than solved exactly.

## Incomplete LU

ILU with no fill keeps exactly the pattern of the original matrix. ILU with a drop tolerance keeps entries larger than a threshold and discards the rest.

## Effect on convergence

A looser tolerance gives a sparser factor and a worse preconditioner. The iteration count rises while the cost per step falls.

## Stability

Dropping entries can make the factorisation unstable, so the diagonal is shifted when a pivot becomes too small. The shift is reduced as the iteration proceeds.

To check your understanding, keep the original pattern, trade fill for steps and detect a small pivot.

Every dropped entry is a memory saving paid for in iteration count.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
