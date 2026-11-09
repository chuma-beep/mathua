# Schur complements

A Schur complement eliminates one block of a matrix. For a 2x2 block matrix with a scalar first block, it is d minus c a inverse b.

## The Schur complement

Partition a matrix into blocks A, B, C and D. The Schur complement of A is D minus C A inverse B. It is the block analogue of eliminating a variable. When A is the scalar a, the formula reduces to d minus c b over a.

## Determinant

The determinant of the full block matrix is the determinant of A times the determinant of its Schur complement. For a 1x1 block this is a times d minus c b over a. The product simplifies to the usual 2x2 determinant.

## Invertibility

The block matrix is invertible exactly when A is invertible and the Schur complement is nonzero. A zero Schur complement signals a dependency between the blocks. This criterion is the block version of a zero determinant.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
