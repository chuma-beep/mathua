# Vandermonde matrices

A Vandermonde matrix has powers of distinct nodes in each row. Its determinant is the product of all pairwise differences.

## Determinant

For nodes x, y and z, the Vandermonde determinant is (y minus x)(z minus x)(z minus y). The order of the factors follows the row order. A repeated node makes one factor zero and the determinant zero.

## Entries

Row i contains the powers of the i-th node from 0 to n minus 1. The (2,3) entry is the square of the second node. The first column is always all ones.

## Rank

Distinct nodes give a nonzero determinant, so the square Vandermonde matrix has full rank. Repeated nodes reduce the rank. The matrix is the coefficient matrix for polynomial interpolation.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
