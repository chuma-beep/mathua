# The Transpose of a Matrix

## Reading a transpose

A transpose turns each row into a column while keeping the order of the entries. The entry in row i and column j of the original matrix becomes the entry in row j and column i of the transpose, and the shape swaps from m by n to n by m.

This section develops the idea and connects it to the surrounding linear algebra.

- Identify rows of the original matrix as columns of the transpose
- Write the entry in a named position after transposing
- Recognise that the shape swaps

## Transpose rules

Transposing a sum distributes across the sum, but transposing a product reverses the order. The identity (AB)^T = B^T A^T matters because matrix multiplication is not commutative, and the double transpose returns the original matrix.

This section develops the idea and connects it to the surrounding linear algebra.

- Apply the transpose to a sum
- Apply the transpose to a product in the correct order
- Use the double transpose identity

## Symmetric matrices

A square matrix is symmetric when it equals its transpose. Symmetry forces entries mirrored across the main diagonal to match, so the condition can be checked by comparing one entry above the diagonal with its partner below it.

This section develops the idea and connects it to the surrounding linear algebra.

- State the condition for symmetry
- Test a small matrix for symmetry
- Connect symmetry to the transpose
