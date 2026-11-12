# Backtracking

Backtracking explores a decision tree depth first and abandons a branch as soon as its partial assignment violates a constraint. It is exhaustive search with early pruning.

## The search tree

Each level fixes one decision. A node is a partial assignment, and a leaf is either a complete solution or a dead end.

## Pruning

A partial assignment that already violates a constraint cannot be completed, so the whole subtree is skipped. Good pruning decides whether backtracking is practical.

## N queens

Placing queens row by row, the constraints are column and both diagonals. A partial board with two queens attacking is pruned immediately.

To check your understanding, fix one decision per level, detect a violated constraint and place row by row.

Backtracking is complete because it explores every branch that survives pruning. Its cost is governed by how early the constraints fail.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
