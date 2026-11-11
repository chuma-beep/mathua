# The knapsack problem

The knapsack problem chooses items with weights and values to maximize total value under a weight capacity. It is the standard integer program with a single constraint, and its fractional relaxation is solved by taking items in value-per-weight order. The difference between the fractional and integer answers is the integrality gap of the instance.

## Solving a small instance

For a few items, enumerate the subsets that fit the capacity and take the largest value. The enumeration is finite and exact, and it exposes the trade-off between value and weight. For larger instances, dynamic programming over capacities replaces the enumeration.

## Reading a DP table

The table has one row per item and one column per capacity. Each entry is the best value using the items so far within that capacity. The recurrence either skips the item or takes it and adds the value from the row above at the reduced capacity. The final entry is the optimum.

## Testing a fractional bound

Sort items by value per unit weight, take whole items until the capacity is nearly full and then take a fraction of the next. The result is the LP relaxation optimum and an upper bound on the integer optimum. The gap to the best integer solution measures how much the indivisibility costs.

The knapsack is small enough to solve by hand and rich enough to show relaxation, dynamic programming and the integrality gap in one model.
