# The transportation problem

The transportation problem moves goods from sources with supplies to destinations with demands at minimum cost. Its constraint matrix has a special structure, one row per source and one column per destination, and a basic feasible solution has one fewer positive route than the number of rows plus columns.

## Balancing a table

Add the supplies and the demands. If they differ, add a dummy source or destination with the difference and zero costs. A balanced instance has equal totals and a feasible flow. Without balancing, the constraints are inconsistent or a destination is left unmet.

## Finding a basic solution

The northwest corner rule fills the table from the top left, moving right or down as a supply or demand is exhausted. It produces a basic feasible solution with the required number of positive entries. Other rules, such as minimum cost, produce a cheaper starting point but the same structure.

## Computing a cost

Multiply each shipped quantity by its unit cost and sum. The cost is what the objective measures, and a proposed improvement is accepted only when the recomputed total is lower. The arithmetic is a check on the routing, not a substitute for the optimality test.

Transportation is the canonical structured linear program. Its table is a compact model, and its basic solutions have a counting rule that makes hand computation feasible.
