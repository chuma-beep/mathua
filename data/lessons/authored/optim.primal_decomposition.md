# Primal decomposition

Primal decomposition assigns each shared resource a master variable. The master allocates, and the blocks optimise given the allocation.

## Blocks and coupling

Each block has local variables and consumes a share of a shared resource. The sum of the shares is the coupling constraint.

## The master problem

The master chooses the shares, knowing each block's value as a function of its share. The block value functions are the pieces the master assembles.

## Convergence

The master is a one-dimensional problem per resource when the value functions are convex. A line search or a cutting plane then converges quickly.

To check your understanding, count the independent blocks, count the shared resources and sum the allocated shares.

Primal decomposition keeps the blocks feasible at every step. Dual decomposition may violate the coupling until the prices settle.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
