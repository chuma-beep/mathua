# The simplex pivot

The simplex method moves from one basic feasible solution to an adjacent one by exchanging a basic variable for a non-basic variable. The exchange is a pivot, and the rules for choosing the entering and leaving variable guarantee that the objective does not increase. Repeating until no improving exchange exists gives an optimal vertex.

## Choosing an entering variable

In the objective row, a negative reduced cost for a non-basic variable means increasing that variable lowers the objective. Any such column can enter, and the most negative is a common choice. If every reduced cost is non-negative, the current basis is optimal and the method stops.

## Performing a pivot

Divide the pivot row by the pivot entry to make it one, then eliminate the entering column from every other row, including the objective row. The result is a new tableau with a new basis. The right-hand side stays non-negative if the leaving row was chosen by the minimum ratio test.

## Testing optimality

Optimality is the absence of an improving entering column. In the tableau this means every entry of the objective row is at least zero for a minimization problem. The test is local to the tableau, yet it certifies global optimality because of linear programming duality.

The pivot is the elementary step of the simplex method. One pivot, one ratio test and one optimality check are enough to follow any small instance by hand.
