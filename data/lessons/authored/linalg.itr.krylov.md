# Krylov subspaces

The Krylov subspace of order k is spanned by the first k images of a vector under a matrix. Krylov methods search for approximate solutions inside it.

## Definition

The subspace is spanned by v, Av, A squared v and so on up to A to the k minus one v. Its dimension is at most k and can be smaller when an invariant subspace appears.

## Nestedness

Each Krylov subspace contains the previous one. The dimensions increase until the sequence becomes dependent, which happens at the degree of the minimal polynomial of v.

## Methods

Conjugate gradients, GMRES and Arnoldi all build bases of these subspaces. They differ in the inner product and the optimality condition they impose.

To check your understanding, list the images, nest the subspaces and build a basis.

Krylov subspaces are the natural search spaces for matrix-free iterative methods.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
