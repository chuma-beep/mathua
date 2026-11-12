# Sparsity in optimization

A sparse model has mostly zero coefficients. Storing only the nonzeros makes large problems tractable.

## Nonzeros

A sparse matrix stores its nonzero entries with their row and column indices. The count of nonzeros, not the row count times the column count, sets the cost.

## Storage

Coordinate storage keeps a triple per nonzero. Compressed storage keeps the values plus one index per column, which is smaller and faster to scan.

## Work

A sparse matrix-vector product touches each nonzero once. Its cost is proportional to the nonzero count, which is why a dense method becomes affordable on a sparse model.

To check your understanding, count nonzeros, size coordinate storage and cost a product.

Sparsity is the structure that makes a large problem look small.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
