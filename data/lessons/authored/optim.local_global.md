# Local and global minima

A local minimum beats every nearby feasible point, while a global minimum beats every feasible point. The two coincide for convex problems, and they can differ sharply otherwise. Most algorithms can only certify local behaviour, so knowing which claim a method supports is part of reading its output.

## Comparing values

To compare two candidates, evaluate the objective at each and order the numbers. This settles which is lower, but not which is global, because a third point may be lower still. Comparison is still the basic operation behind every branch and bound and every grid search, where a bound supplies the missing guarantee.

## Identifying a global minimum

A global minimum is certified either by an argument that covers the whole feasible set or by an exhaustive bound. For a convex objective over a convex set, any stationary point is global. For a finite candidate set, the smallest evaluated value is global by exhaustion. Without one of those two structures, a computed point is a candidate and should be described as one.

## Local versus global

A point can be locally minimal and not globally minimal, and the standard picture is a landscape with two valleys of different depths. Local information, a gradient and a Hessian, cannot see the other valley. Methods that escape use restarts, random perturbations or a global bound, and each of those is a statement about how the search covers the space rather than about the local condition.

The distinction decides what a result means. Reporting a local minimum as the minimum is a modelling claim that the objective has one valley, and that claim needs its own justification.
