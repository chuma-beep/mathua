# Singular Value Decomposition

**SVD:** Any real $A = U\Sigma V^T$ with $U, V$ orthogonal and $\Sigma$ diagonal of singular values $\sigma_i \ge 0$. Rank equals the nonzero count; truncating gives the best low-rank approximation (Eckart-Young).

## Worked: rank from singular values 5, 2, 0

Count the nonzero ones:
1. $\sigma = (5, 2, 0)$: two nonzero.
2. So rank is 2: the zero singular value marks a dead direction.
3. Best rank-1 approximation keeps just $\sigma_1 = 5$.

So the singular spectrum reads off rank and compressibility at a glance.

## Worked: largest singular value from A^TA

Given $A^TA$ eigenvalues $9, 4, 0$:
1. $\sigma_i = \sqrt{\lambda_i}$: singular values are root-eigenvalues.
2. Largest: $\sqrt{9} = 3$.
3. So $\sigma_{\max} = 3$: the top stretch factor of $A$.

So $A^TA$ converts SVD questions into symmetric eigenvalue questions.

## Vectors

$U$'s columns are left singular vectors, $V$'s are right: $Av_i = \sigma_i u_i$. The decomposition aligns domain and codomain along matched orthonormal bases.
