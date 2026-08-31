# LU Decomposition

**LU:** $A=LU$ with $L$ unit lower triangular, $U$ upper triangular. Solve $Ax=b$ via $Ly=b$ (forward) then $Ux=y$ (back substitution). Without pivoting requires nonzero leading principal minors; $PA=LU$ handles zeros.

## Factoring Matrices

### Example
$A=\begin{pmatrix}2&1\\4&3\end{pmatrix}= \begin{pmatrix}1&0\\2&1\end{pmatrix}\begin{pmatrix}2&1\\0&1\end{pmatrix}$; $\det(A)=\det(U)=2·1=2$.

### Determinant
$\det(A)=\det(U)$ when $L$ unit, product of $U$ diagonal.

## Example

Solve $Ax=b$: $Ly=b$ gives $y$, then $Ux=y$ gives $x$.
