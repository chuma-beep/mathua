# Banded storage

A banded matrix has all nonzeros within a fixed distance of the diagonal. Storing only that band turns an n-by-n array into a rectangle of width p+q+1.

## The band

Lower bandwidth p and upper bandwidth q bound the offsets that can be nonzero. The band holds at most (p+q+1)n entries.

## Elimination

LU without pivoting preserves the band, so the factors fit in the same rectangle. Pivoting inside the band keeps it; pivoting outside destroys it.

## Cost

Banded elimination costs O(n p squared) when p and q are comparable. That is linear in n for a fixed band.

To check your understanding, read the bandwidths, preserve the band and count the elimination work.

The band is a promise that the matrix makes and pivoting can break.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
