# Weak solutions

A weak solution moves the derivatives onto a smooth test function, so a solution need not have two classical derivatives. The definition extends the heat and wave equations to data that no classical solution could accept.

## Test functions

A test function is smooth and vanishes outside a compact set inside the space-time domain. Multiplying the equation by the test function and integrating leaves every term meaningful for a solution with only one derivative.

## Integration by parts

Every derivative on the solution is transferred to the test function, with a sign for each transfer. The boundary terms disappear because the test function vanishes near the boundary, which is why no boundary condition has to be imposed on the test space.

## Existence

A Galerkin approximation in finite-dimensional subspaces produces weak solutions, and the Lax-Milgram theorem gives existence and uniqueness for a coercive bilinear form. Elliptic regularity then recovers classical smoothness when the data allow it.

To check your understanding, choose a compactly supported test function, transfer each derivative and build the Galerkin subspace.

The weak formulation is the version a finite-element method actually discretises, so it is the definition with numerical consequences.
