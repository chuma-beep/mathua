# Moore-Penrose pseudoinverse

The Moore-Penrose pseudoinverse extends the inverse to matrices that are not square or not invertible. On a diagonal matrix it simply inverts the nonzero diagonal entries.

## Diagonal pseudoinverse

For a diagonal matrix D, the pseudoinverse D plus has the reciprocal of each nonzero diagonal entry and zero elsewhere. A zero diagonal entry stays zero. This is the simplest case of the singular-value construction.

## Rank

The rank of a diagonal matrix is the number of nonzero diagonal entries. The pseudoinverse has the same rank. Zero rows and columns remain zero in both matrices.

## The product D D plus

Multiplying a diagonal matrix by its pseudoinverse gives a diagonal matrix with 1 where the original entry was nonzero and 0 where it was zero. The trace therefore counts the nonzero diagonal entries. It equals the rank.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
