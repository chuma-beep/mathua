# The Kauffman bracket

The Kauffman bracket is computed by smoothing every crossing in two ways. A diagram with \(n\) crossings has \(2^n\) states, and the bracket is the weighted sum over all of them.

## States

Each crossing offers two smoothings, so a diagram with \(n\) crossings has \(2^n\) states. The trefoil has eight states.

## The loop count

Each state produces some number of disjoint loops. The bracket weights a state by \(A\) raised to the difference of the two smoothing counts and by the loop count.

## The unknot

The bracket of the unknot is one by convention. The Jones polynomial then follows after the writhe correction.

To check your understanding, raise two to the crossings, count the loops and read the normalisation.

The bracket is a state sum over smoothings.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
