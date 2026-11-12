# Block coordinate descent

Block coordinate descent moves one group of variables at a time and holds the rest fixed. Each step solves a smaller problem exactly or approximately.

## One block at a time

The variables are partitioned into blocks. A step optimises one block with the others frozen, which reduces the dimension of each solve.

## The block gradient

With the other blocks fixed the gradient of the block is the matching slice of the full gradient. Its size measures how far the block is from its own optimum.

## Sweep order

A cyclic sweep visits every block in order and starts again. Randomised orders can be faster when the blocks are strongly coupled.

To check your understanding, fix the other blocks, partition the variables and count the blocks per sweep.

Convergence needs each block subproblem to be solved well enough. A loose block solve can stall the sweep.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
