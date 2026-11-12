# Sparse LU factorisation

Sparse LU computes \(PAQ=LU\) while keeping the factors as sparse as the ordering allows. The symbolic phase predicts the pattern and the numeric phase fills in the values.

## Two phases

Symbolic factorisation computes the elimination tree and the nonzero pattern. Numeric factorisation then computes values without any further pattern search.

## Pivoting

Threshold partial pivoting keeps the multipliers bounded while limiting the departure from the chosen order. Stability and sparsity pull in opposite directions.

## Cost

The work is proportional to the number of nonzeros in \(L\) and \(U\), not to \(n^{3}\). For a banded matrix it is proportional to \(n\) times the bandwidth squared.

To check your understanding, build the elimination tree, bound the multipliers and count the factor nonzeros.

A sparse factorisation is planned before it is computed.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
