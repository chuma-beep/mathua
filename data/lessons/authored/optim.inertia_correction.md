# Inertia correction

A trust-region step needs the reduced Hessian to be positive definite. Inertia correction shifts the Hessian until the KKT factorisation reports the right number of negative eigenvalues.

## The inertia test

A symmetric factorisation yields the number of positive, negative and zero eigenvalues of the shifted matrix. The step is a descent step only when the reduced Hessian has the expected inertia.

## The shift

If the inertia is wrong the primal block is shifted by a positive multiple of the identity and the factorisation is repeated. The shift is increased until the inertia matches, so the step approaches a gradient direction in the worst case.

## The trade

A large shift produces a short, safe step and slow progress. An adaptive update that increases the shift on failure and decreases it on success keeps the method close to a pure Newton step when the problem allows it.

To check your understanding, factor the KKT matrix, add a multiple of identity and shorten the step.

Inertia correction is the mechanism that keeps a constrained Newton method moving when the local model is not convex.
