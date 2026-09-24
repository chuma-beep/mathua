# Jordan Canonical Form

**Jordan form:** Every complex matrix is similar to block-diagonal Jordan form: eigenvalues on the diagonal, ones on the superdiagonal within each block. Block sizes encode the gap between algebraic and geometric multiplicity.

## Worked: [[2,1],[0,2]] is one block

Inspect directly:
1. Eigenvalue 2 (double root of the characteristic polynomial).
2. One eigenvector direction only (rank of $A - 2I$ is 1).
3. So a single $2 \times 2$ block: one Jordan block, already in Jordan form.

So a defective matrix shows its block structure on its face: repeated eigenvalue, too few eigenvectors.

## Worked: counting blocks in a diagonalizable 3x3

If diagonalizable:
1. Every block is $1 \times 1$ (no nontrivial blocks possible).
2. Three diagonal entries means three blocks.
3. So block count equals geometric-multiplicity sum: full eigenvector supply.

So diagonalizability is the all-trivial-blocks case; anything else signals defectiveness.

## Minimal polynomial

The largest block size for $\lambda$ is its exponent in the minimal polynomial. Block sizes pin the polynomial exactly — no finer invariant exists.
