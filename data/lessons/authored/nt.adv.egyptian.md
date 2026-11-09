# Egyptian fractions

An Egyptian fraction writes a rational as a sum of distinct unit fractions. The greedy algorithm chooses the largest unit fraction that does not exceed the remainder.

## The greedy algorithm

For a fraction a over b, the first denominator is the ceiling of b over a. That is the smallest d with 1 over d at most a over b. Subtract the unit fraction and repeat on the remainder.

## The next denominator

After subtracting 1 over d, reduce the remainder and take the ceiling of its reciprocal. That is the next denominator. The numerators stay positive and the process terminates for positive rationals.

## The unit-fraction identity

The identity 1 over n equals 1 over n plus 1 plus 1 over n times n plus 1 splits a unit fraction into two. The larger denominator is n times n plus 1. This identity is a useful shortcut in Egyptian-fraction exercises.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
