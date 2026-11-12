# Compressed sparse row storage

Compressed sparse row storage keeps the nonzero values in row order and records where each row starts. It removes the zeros a dense array would store.

## Three arrays

The values array and the column-index array each hold one entry per nonzero. The row-pointer array holds one entry per row plus one.

## Row access

Row i occupies the half-open slice between pointers i and i+1. Reading a row is a contiguous scan of two parallel arrays.

## Cost

Storage is \(2\,\mathrm{nnz}+n+1\) numbers instead of \(mn\). The format wins whenever the matrix is sparse enough.

To check your understanding, name the three arrays, read the row pointers and count the stored numbers.

Compressed storage makes the zero structure explicit and the nonzeros contiguous.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
