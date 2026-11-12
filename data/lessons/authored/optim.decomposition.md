# Decomposition

Decomposition splits a large problem into subproblems that are coupled by a few constraints. Prices on those constraints coordinate the subproblems.

## Blocks

The variables divide into groups that appear in separate constraints. Only a few rows link the groups, and those rows are the coupling.

## Coordinating

The master problem chooses the prices and each subproblem solves its own block. The prices are updated until the subproblem decisions agree with the coupling rows.

## Dual prices

A price is the shadow price of a coupling row. A positive price tells the subproblems to use less of the shared resource.

To check your understanding, identify independent blocks, set prices and read a shadow price.

Decomposition replaces one large solve with many small ones and a price loop.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
