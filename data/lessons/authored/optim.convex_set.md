# Convex sets

A set is convex when the segment between any two of its points stays inside it. The definition is short, and it is the property that makes optimization tractable, because a convex feasible set plus a convex objective rules out the second valley. Testing convexity is therefore a routine check on a model, not an advanced topic.

## Testing a combination

Take two points in the set and a weight between zero and one, then form the weighted combination. The set is convex when the combination lies in the set for every such choice. A single pair whose segment leaves the set disproves convexity, so a counterexample is the efficient test. A proof needs an argument that covers all pairs, usually by showing the set is an intersection of half-spaces or a sublevel set of a convex function.

## Separating lines

When a point lies outside a closed convex set, a line separates the two, with the set on one side and the point on the other. The separating line is useful because it turns a membership question into a single linear inequality. For a polygon, the line can be chosen along a violated edge, and the distance from the point to that edge measures how far outside the point is.

## Convexity from a description

Several descriptions are immediately convex. A half-space is convex. An intersection of convex sets is convex. A ball, a box and a simplex are convex. A union of two disjoint disks is not, and a set with a hole is not, because a segment can cross the gap. Reading a description for unions and holes is the fastest way to decide the common cases.

Convexity is what licenses a local certificate to be a global one. Checking it early tells the modeller whether a fast method is valid or whether the problem needs a global search.
