# Arnoldi basis construction

The Arnoldi process builds an orthonormal basis of a Krylov subspace. Each new vector is orthogonalised against all previous ones.

## Step

Multiply the last basis vector by the matrix, then orthogonalise against the existing basis. The norm of the remainder gives the next basis vector.

## Hessenberg

The coefficients of the orthogonalisation form an upper Hessenberg matrix. Its eigenvalues approximate those of the original matrix, which is the basis of iterative eigensolvers.

## Cost

Step k costs k inner products and one matrix-vector product. Storage grows linearly with k unless the process is restarted.

To check your understanding, multiply by the matrix, collect the coefficients and count inner products.

Arnoldi is Gram-Schmidt applied to a Krylov sequence, and its Hessenberg matrix is the reward.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
