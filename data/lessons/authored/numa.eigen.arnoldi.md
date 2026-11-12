# The Arnoldi iteration

The Arnoldi iteration generalises Lanczos to nonsymmetric matrices. It builds an orthonormal basis and an upper Hessenberg matrix whose eigenvalues approximate those of the operator.

## Gram-Schmidt

Each new vector is orthogonalised against every previous basis vector, so the recurrence is long rather than three-term. Modified Gram-Schmidt is used in practice for stability.

## The Hessenberg matrix

The projected matrix is upper Hessenberg, with coefficients taken from the orthogonalisation. It is small whenever the number of steps is small.

## Cost

Step \(k\) costs one matrix-vector product and about \(2kn\) operations for orthogonalisation, so the total is \(O(k^{2}n)\).

To check your understanding, orthogonalise the new vector, read the projection and count the products.

Arnoldi projects a large nonsymmetric problem onto a small Hessenberg one.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
