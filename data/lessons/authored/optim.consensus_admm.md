# Consensus ADMM

Consensus ADMM gives each block a local copy of the shared variable. The copies are pulled together by a penalty, and the average is the consensus value.

## Local copies

Every block keeps its own copy of the global variable. Its local objective sees only its copy, so the blocks stay independent.

## The consensus penalty

A quadratic penalty measures each copy against the global average. The penalty weight controls how strictly the copies must agree.

## The update order

Each block updates its copy, then the global variable becomes the average, then the dual variables absorb the disagreement. One pass is one ADMM iteration.

To check your understanding, sum the local copies, subtract the average from a copy and compute each deviation.

The average is what makes the method distributed. No block ever sees another block's data, only the shared average.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
