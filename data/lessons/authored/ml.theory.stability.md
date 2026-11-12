# Algorithmic stability

Algorithmic stability bounds generalisation by how much the output changes when one training point is replaced. It avoids capacity measures entirely.

## Replace-one stability

Change one point and compare the two outputs. A small expected change implies a small generalisation gap.

## The bound

The gap is at most twice the stability parameter plus a confidence term. Regularisation usually supplies the stability.

## The appeal

Stability applies to any algorithm, including those with unbounded VC dimension. The price is a bound that can be loose in practice.

To check your understanding, replace a point, double the parameter and skip the capacity.

Stability measures sensitivity to a single training point.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
