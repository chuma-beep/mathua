# QR Decomposition

**QR decomposition:** $A = QR$ with $Q$ orthogonal ($Q^TQ = I$) and $R$ upper triangular. Via Gram-Schmidt or Householder reflections; $R$'s diagonal positive by convention. It solves least squares as $Rx = Q^Tb$ — stabler than the normal equations.

## Worked: normalize (1,1)

First Gram-Schmidt step on $A = [[1,1],[1,0]]$:
1. Take $a_1 = (1, 1)$.
2. Norm: $\|(1,1)\| = \sqrt{2}$.
3. So $q_1 = (1,1)/\sqrt{2}$: divide by the length just computed.

So orthonormalization is normalize-as-you-go; each $q$ is a unit vector by construction.

## Worked: count R's entries

For $2 \times 2$ upper-triangular $R$:
1. Diagonal: 2 entries.
2. Above diagonal: 1 entry.
3. So 3 possibly-nonzero entries — the zeros below the diagonal are structural.

So triangularity is a sparsity promise: half the matrix is known zeros.

## Back substitution

$Rx = Q^Tb$ solves bottom-up (back substitution), since $R$ is upper triangular. Householder reflections compute the same factorization more stably than classical Gram-Schmidt.
