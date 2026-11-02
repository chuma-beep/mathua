# The Cantor Set

## Constructing the Cantor set

Begin with the closed interval from 0 to 1. Remove the open middle third, leaving two closed thirds. Remove the middle third of each remaining interval, and continue indefinitely. The Cantor set is the intersection of all the stages. It is closed because it is an intersection of closed sets, and it is bounded, hence compact.

At stage n there are 2 to the n closed intervals, each of length 3 to the minus n. The total length of the stage is therefore (2/3)^n, which tends to zero. The limiting set is not empty: endpoints of removed intervals survive, and many more points survive as well.

## What is removed

The complement is a countable union of open intervals, so its measure is the sum of the removed lengths. After n stages the removed total is 1 minus (2/3)^n, and in the limit the complement has measure 1. The Cantor set therefore has measure zero even though it is uncountable.

Membership can be described in base three. A point belongs to the Cantor set exactly when it has a ternary expansion using only the digits 0 and 2. This digit description makes the uncountability transparent: every binary sequence corresponds to a ternary sequence of 0s and 2s, and different sequences give different points.

## Properties of the remainder

The Cantor set is uncountable, has empty interior and has measure zero. It is a standard example showing that uncountable does not imply positive measure and that a set can be large in cardinality while small in length. It is also totally disconnected: between any two of its points there is a gap.

The set is self-similar, and its Hausdorff dimension is log 2 divided by log 3, a number between zero and one. That fractional dimension captures the way the set sits between a point and an interval. The Cantor set reappears in real analysis as a source of counterexamples and as a building block for singular functions.
