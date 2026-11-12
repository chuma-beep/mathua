# The Lanczos iteration

The Lanczos iteration builds a tridiagonal matrix from a symmetric operator with a three-term recurrence. Each step needs one matrix-vector product and keeps only the last two basis vectors.

## The recurrence

\(\beta_j q_{j+1}=Aq_j-\alpha_j q_j-\beta_{j-1}q_{j-1}\). Orthogonality is inherited from symmetry rather than enforced by reorthogonalisation.

## The projected problem

After \(k\) steps the eigenvalues of the \(k\times k\) tridiagonal matrix estimate the extreme eigenvalues of \(A\). The estimates improve as \(k\) grows.

## Cost

The method costs \(k\) matrix-vector products and stores three vectors. That is why it suits very large sparse problems where a factorisation is impossible.

To check your understanding, write the three-term recurrence, form the tridiagonal matrix and count matrix-vector products.

Lanczos trades a full factorisation for a sequence of cheap projections.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
