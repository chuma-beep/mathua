# Functions of Bounded Variation

## Partitions and variation

Given a function on a closed interval and a partition of that interval, form the sum of the absolute values of the changes over the subintervals. The total variation is the supremum of these sums over all partitions. A function has bounded variation when that supremum is finite.

Monotone functions have bounded variation equal to the absolute difference of their endpoint values. A function that oscillates more and more violently can have unbounded variation even when it is continuous. Variation measures the total up-and-down travel of the graph rather than its range.

## Monotone functions

Every function of bounded variation is the difference of two increasing functions. Split the variation into positive and negative parts: the positive part records upward movement and the negative part records downward movement. Both are increasing, and their difference recovers the original function up to a constant.

This decomposition is the main structural theorem about bounded variation. It implies that a function of bounded variation has at most countably many discontinuities, all of them jumps. It also connects variation to the theory of the Riemann-Stieltjes integral, where the integrator must have bounded variation for the integral to behave well.

## Algebra of variation

Sums and scalar multiples of functions of bounded variation again have bounded variation, and the variation of a sum is at most the sum of the variations. Products also have bounded variation on a bounded interval when both factors are bounded. These closure properties make the class convenient for integration theory.

The total variation defines a seminorm, and adding the absolute value at one point makes it a norm. The resulting space is complete. That completeness is another reason the class appears in analysis: it is a Banach space that contains all monotone functions and all Lipschitz functions on the interval.
