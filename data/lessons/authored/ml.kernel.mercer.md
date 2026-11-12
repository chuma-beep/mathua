# Mercer's theorem

Mercer's theorem turns a continuous positive semidefinite kernel into an inner product in a feature space. The eigenfunctions of the integral operator supply the coordinates.

## The integral operator

The kernel defines an operator that maps a function to a weighted integral against the kernel. Positive semidefiniteness makes its eigenvalues non-negative.

## The expansion

The kernel equals a sum of eigenvalue-weighted eigenfunction products. Truncating the sum gives a finite approximation.

## The consequence

Every valid kernel is an inner product in some space. That is the licence behind the kernel trick.

To check your understanding, define the operator, write the sum and state the inner product.

Mercer's theorem is the bridge from a kernel function to a feature space.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
