# Compressed sparse column storage

Compressed sparse column storage mirrors the row format with the roles exchanged. It makes column access contiguous, which suits the column-oriented kernels of sparse factorisation.

## Layout

Values and row indices are grouped by column, and a pointer array of length n+1 marks where each column begins.

## Choosing a format

The row format suits row-wise products and triangular solves; the column format suits column updates and factorisation. Converting costs one pass over the nonzeros.

## Storage

The count is again \(2\,\mathrm{nnz}+n+1\), so the formats differ in access pattern rather than in size.

To check your understanding, group values by column, match format to kernel and count the stored numbers.

The two compressed formats hold the same numbers in a different order.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
