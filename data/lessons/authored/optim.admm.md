# ADMM

The alternating direction method of multipliers splits a problem into two blocks and alternates between them with a dual update.

## The split

A new variable copies the first block and an equality constraint ties the two. The split makes each block easy on its own.

## The steps

One step minimises over the first block, one over the second, and one updates the dual variable. The second step often has a closed form.

## The penalty

A quadratic penalty keeps the two blocks close. A large penalty enforces agreement quickly and can slow the first block, so the choice is a trade-off.

To check your understanding, write the copy constraint, minimise one block and choose a penalty.

ADMM is the augmented Lagrangian with the minimisation split into easy pieces.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
