# Multistep methods

A multistep method uses several previous solution values to advance one step. Adams-Bashforth methods are explicit and Adams-Moulton methods are implicit.

## Apply the multistep formula

The two-step Adams-Bashforth update is y next equals y plus h times three halves f sub n minus one half f sub n minus 1. The previous slope is stored from the last step. The formula is a weighted average of recent slopes.

## Determine the order

The order is the power of h in the leading local error term. Adams-Bashforth with k steps has order k. A higher order needs more history but gives a smaller error.

## Check zero stability

The roots of the first characteristic polynomial must lie inside or on the unit circle. A root outside the circle makes the method unstable. A double root on the circle also fails the strict stability test.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
