# QR Decomposition

**QR:** $A=QR$ with $Q$ orthogonal ($Q^TQ=I$) and $R$ upper triangular. Via Gram-Schmidt or Householder reflections; $R$'s diagonal positive by convention.

## Least Squares

### Solving
$Ax=b$ least squares via $Rx=Q^Tb$ — stabler than normal equations.

### Stability
Householder reflections give numerically stable QR; $Q$'s columns orthonormal.

## Example

$A=\begin{pmatrix}1&1\\1&0\end{pmatrix}$: $q_1=(1,1)/\sqrt2$, $q_2=(1,-1)/\sqrt2$, $R=\begin{pmatrix}\sqrt2&1/\sqrt2\\0&1/\sqrt2\end{pmatrix}$.
