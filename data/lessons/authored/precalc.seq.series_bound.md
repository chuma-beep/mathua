# Bounds on series remainders

A convergent series can be approximated by a partial sum, and the error is bounded by the size of the remaining tail. For geometric series the tail has a closed form.

## The tail

The remainder after n terms is the sum of all later terms. For a geometric series it is the first omitted term divided by one minus the ratio.

## Bounding the error

Taking the absolute value of the tail gives an upper bound on the error of the partial sum, provided the ratio has absolute value below one.

## Choosing n

Set the bound below a target error and solve the resulting inequality for n, then round up because the term count is an integer.

To check your understanding, identify the first omitted term, take the absolute value and set the target.

A bound is a guarantee, so a smaller computed error does not weaken it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
