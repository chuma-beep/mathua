# Singular Value Decomposition

**SVD:** Any $m×n$ real $A=U\Sigma V^T$ with $U$ $m×m$ orthogonal, $V$ $n×n$ orthogonal, $\Sigma$ $m×n$ diagonal with $\sigma_i\ge0$ singular values.

## Geometry and Rank

### Singular Values
$\sigma_i=\sqrt{\lambda_i(A^TA)}$, $\text{rank}(A)=\#\{\sigma_i>0\}$, truncated SVD gives best low-rank approximation (Eckart-Young).

### Vectors
$U$'s columns left singular vectors, $V$'s right; $Av_i=\sigma_i u_i$.

## Example

$A=\begin{pmatrix}3&0\\0&2\end{pmatrix}$: $U=I$, $V=I$, $\Sigma=\text{diag}(3,2)$, singular values $3,2$, rank $2$.
