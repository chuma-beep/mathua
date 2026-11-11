# Shadow prices

A shadow price is the rate at which the optimal objective changes when a constraint's right-hand side increases by one unit. It equals the corresponding dual variable at optimality, and it is valid only while the optimal basis stays the same. Shadow prices are the economic reading of a dual solution.

## Reading a shadow price

Take the multiplier of the constraint in the optimal dual solution. Its sign follows the convention of the primal form, and its magnitude is the change per unit of right-hand side. A zero price means the constraint is not binding at the current optimum.

## Estimating a constraint change

For a small change within the allowable range, the new objective is approximately the old objective plus the price times the change. The estimate is exact for linear programs inside the range and approximate for nonlinear problems. Outside the range the basis changes and the price no longer applies.

## Comparing two prices

Two constraints can be ranked by their prices, and the larger price identifies the relaxation that buys more objective per unit. A price comparison is only meaningful within the same range and the same units. Converting units before comparing avoids a common modelling misstatement.

Shadow prices turn a solved model into advice. They say which constraint to relax first and how much a relaxation is worth, provided the range is respected.
