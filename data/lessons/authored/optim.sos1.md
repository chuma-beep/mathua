# SOS1 constraints

A special ordered set of type one allows at most one of a group of variables to be nonzero. It expresses a choice among alternatives without introducing a binary for each alternative.

## The set

The solver enforces that at most one member of the ordered set is nonzero. The variables are ordered, so a branching rule can split the set by position rather than by value.

## Against big-M

An SOS1 branch is often stronger than a binary switch because the branching decision uses the order. The model also avoids a large constant, which removes one source of numerical weakness.

## The modelling fit

Piecewise linear functions and disjunctive choices are the natural uses. A piecewise function is written as a convex combination of breakpoints with an SOS1 set on the segments, which keeps the model convex.

To check your understanding, allow one nonzero, compare bound strength and write a piecewise function.

An SOS1 set is a modelling device for exclusive choice, and its advantage over big-M is a branching rule that uses the order.
