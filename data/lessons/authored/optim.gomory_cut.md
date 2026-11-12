# Gomory cuts

A Gomory cut is derived from a row of the optimal simplex tableau. It removes the current fractional solution without cutting off any integer point, so it strengthens the relaxation exactly where it is weak.

## The derivation

The tableau expresses a basic variable as an integer part plus a fractional part of the nonbasic variables. A fractional right-hand side forces a rounding argument, and the resulting inequality holds for every integer point.

## The separation

Each fractional basic variable yields one candidate cut, so a node can generate several. Choosing the deepest cut is a heuristic that correlates with bound improvement.

## Numerical care

The coefficients of repeated cuts grow and the tableau becomes ill-conditioned. Modern solvers limit the cut rounds and remove cuts that are no longer useful, which keeps the basis stable.

To check your understanding, read a tableau row, take each fractional basic and limit cut rounds.

A Gomory cut is a rounding argument read off the tableau, and its weakness is numerical growth rather than a lack of validity.
