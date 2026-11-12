# Mini-batch gradients

A mini-batch average reduces gradient variance by the batch size while multiplying the cost by the same factor.

## The average

The batch gradient is the mean of the per-sample gradients. Its expectation is unchanged and its variance falls as one over the batch size.

## The cost

One batch step costs the batch size times one sample step. The trade is variance against arithmetic, and the best batch size depends on the hardware.

## Choosing a batch

Beyond a point, doubling the batch no longer halves the wall-clock time because the computation stops being the bottleneck. That point is the practical limit.

To check your understanding, average two gradients, count the arithmetic and measure the speedup.

A mini-batch is variance reduction paid for with arithmetic.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
