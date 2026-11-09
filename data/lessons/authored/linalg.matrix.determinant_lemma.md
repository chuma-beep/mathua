# The matrix determinant lemma

The matrix determinant lemma updates a determinant after adding a rank-one term. It reduces a large determinant to a scalar correction.

## The lemma

For an invertible A, the determinant of A plus u v transpose is det A times 1 plus v transpose A inverse u. The scalar v transpose A inverse u is the only new quantity. The lemma is the rank-one case of the Weinstein-Aronszajn identity.

## Determinant update

Compute the scalar first and then multiply by the original determinant. The result can be zero even when A is invertible. A zero result means the rank-one update destroyed invertibility.

## Change in determinant

The change in the determinant is det A times the scalar. This is often easier to read than the new determinant itself. It separates the effect of the update from the size of the original matrix.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
