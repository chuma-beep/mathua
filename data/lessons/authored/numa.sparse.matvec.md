# Sparse matrix-vector products

A sparse matrix-vector product touches each nonzero once and writes one output entry per row. Its cost is proportional to the number of nonzeros rather than to the square of the order.

## Work count

Each nonzero contributes one multiplication and one addition, so the product costs nnz multiply-add operations. The output has one entry per row.

## Memory traffic

The product is memory bound: reading the values, the column indices and the input vector dominates the arithmetic. Blocking rows improves reuse of the input.

## Parallelism

Rows are independent, so the product parallelises without synchronisation. Load balance suffers when row lengths vary widely.

To check your understanding, count the multiply-adds, identify the streams and split by rows.

The sparse product is a memory problem more than an arithmetic one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
