# Decomposing a vector into orthogonal components

Any vector can be split into a component parallel to a chosen direction and a component perpendicular to it. The two pieces add back to the original vector.

## The parallel component

The projection of v onto u is the scalar k u, where k is the dot product of v and u divided by the squared length of u. This component lies along u. Its length is the absolute value of k times the length of u.

## The perpendicular component

Subtract the parallel component from v to get the perpendicular component. The dot product of the perpendicular component with u is zero. That zero is the check that the split is orthogonal.

## Checking the decomposition

Add the two components and recover v. Then dot the perpendicular component with u. A nonzero dot product means the decomposition has an arithmetic error.

Practice connects the three sections: state the relation, carry out the computation and interpret the result. A second method is a useful check when one is available.
