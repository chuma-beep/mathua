# Basic feasible solutions

A basic solution comes from choosing a set of columns as a basis, setting the rest to zero and solving the resulting square system. It is basic feasible when the solved values are non-negative. The vertices of a linear program's feasible region are exactly its basic feasible solutions, which is why the simplex method only visits them.

## Finding a basic solution

Choose as many columns as there are rows, require them to be linearly independent and solve the square system with the other variables at zero. The solution is a point where the chosen constraints are tight. A different basis gives a different point, and enumerating bases enumerates candidates.

## Testing feasibility

Check that every solved basic variable is non-negative and that the non-basic variables are zero. A negative value means the basis is not feasible, even though it solves the equations. Feasibility is a separate check from solving, and forgetting it is the most common error in hand computation.

## Counting possible bases

The number of ways to choose as many columns as rows from the full matrix is a binomial coefficient. For a small model this count is a concrete number, and it explains why enumeration is impractical beyond a few variables. The simplex method avoids enumeration by moving between adjacent bases.

Basic feasible solutions connect linear algebra to geometry. Each one is a vertex, and the optimal vertex is among them, which is the finite search the simplex method organizes.
