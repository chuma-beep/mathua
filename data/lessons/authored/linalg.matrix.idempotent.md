# Idempotent matrices

An idempotent matrix satisfies A squared equal to A. It behaves like a projection, keeping one part of a vector and discarding the rest.

## Projection matrices

An idempotent matrix is a projection. Applying it twice gives the same result as applying it once. Its eigenvalues are 0 and 1, and its trace equals its rank.

## Trace and rank

The trace of an idempotent matrix is an integer between 0 and the dimension. It equals the number of eigenvalues equal to 1, which is the rank. The trace is therefore a quick check on a proposed projection.

## The idempotent test

Square the matrix and compare every entry with the original. Equality in all four positions means the matrix is idempotent. A single mismatch is enough to reject it.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
