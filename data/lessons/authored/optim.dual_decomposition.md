# Dual decomposition

When the constraints couple blocks only through prices, the dual splits into independent subproblems. Each block then solves its own problem at the current price.

## Split the objective

The objective is a sum of block objectives and the coupling appears only in shared constraints. Relaxing those constraints with multipliers removes the coupling.

## The dual price

Each shared constraint carries one price. Given the prices, every block minimises its own Lagrangian without looking at the other blocks.

## Coordinating subproblems

A master update raises or lowers the prices from the subproblem responses. The blocks never exchange data directly, only prices.

To check your understanding, fix the shared price, multiply each shared constraint by its price and count the blocks.

Dual decomposition needs the coupling to be linear. A nonlinear coupling would leave the subproblems dependent on each other.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
