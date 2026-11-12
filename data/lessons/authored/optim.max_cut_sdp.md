# The max-cut relaxation

Max cut is written as a quadratic over signs, and its semidefinite relaxation replaces the signs with unit vectors. The relaxation supplies an upper bound that a rounding step turns into a cut.

## The cut as a quadratic

Assigning a sign to each vertex, the weight of the cut is one quarter of the sum over edges of the squared difference of the signs. Expanding it gives a quadratic form in the sign vector.

## The relaxation

The sign vector becomes a matrix constrained to be positive semidefinite with a unit diagonal. The relaxation maximises the same linear objective over that larger set.

## The rounding bound

A random hyperplane through the unit vectors cuts an edge with a probability set by the angle between its endpoints. That argument gives the approximation guarantee.

To check your understanding, sum the weights of crossing edges, compute the quadratic objective and subtract the cut value from the bound.

The relaxation bounds the cut from above. Comparing that bound with a rounded cut measures how much the relaxation gives away.
