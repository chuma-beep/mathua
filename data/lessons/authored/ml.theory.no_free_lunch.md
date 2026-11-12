# The no free lunch theorem

The no free lunch theorem says that averaged over all possible objectives, every optimisation algorithm performs the same. Any advantage on one problem is paid for on another.

## The averaging

Uniformly average over all functions from a finite domain to a finite range. Every algorithm has the same expected performance.

## The count

With \(n\) binary inputs there are \(2^n\) objectives. The symmetry of the average is what forces the result.

## The escape

Real problems are not uniform over all objectives. A prior or an inductive bias is what makes learning possible.

To check your understanding, average over functions, count the objectives and name the prior.

No free lunch is a statement about averages, not about every problem.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
