# Solving systems with inverse matrices

A linear system can be written as \(A\mathbf{x} = \mathbf{b}\). When the coefficient matrix \(A\) is invertible, multiplying both sides by \(A^{-1}\) gives \(\mathbf{x} = A^{-1}\mathbf{b}\). The inverse method solves the system in one matrix product.

## The determinant

For a \(2 \times 2\) matrix with entries \(a, b, c, d\), the determinant is \(ad - bc\). The matrix has an inverse only when the determinant is nonzero.

The determinant also measures how the matrix scales area. A zero determinant collapses the plane onto a line or a point, and that collapse is why the system loses a unique solution.

For \(A = \begin{pmatrix}2 & 1\\ 3 & 4\end{pmatrix}\), the determinant is \(2 \times 4 - 1 \times 3 = 5\). The matrix is invertible, so the system has one solution.

## The x-coordinate

The inverse of a \(2 \times 2\) matrix is \(\frac{1}{ad-bc}\begin{pmatrix}d & -b\\ -c & a\end{pmatrix}\). Applying it to \(\mathbf{b} = (e, f)\) gives \(x = (de - bf)/(ad-bc)\).

The numerator is a determinant of its own. Replacing the first column of \(A\) by \(\mathbf{b}\) and taking the determinant produces \(de - bf\), which is Cramer's rule in disguise.

For the system \(2x + y = 8\) and \(3x + 4y = 20\), the determinant is 5 and \(x = (4 \times 8 - 1 \times 20)/5 = 12/5\). The inverse method and substitution agree.

## The y-coordinate

The second coordinate is \(y = (af - ce)/(ad-bc)\). It is the determinant with the second column replaced by \(\mathbf{b}\), divided by the determinant of \(A\).

The two coordinates are obtained from the same inverse matrix. Once the determinant and the four inverse entries are known, each coordinate is one numerator and one division.

For the same system, \(y = (2 \times 20 - 3 \times 8)/5 = 16/5\). Substituting both values into either original equation confirms the pair.

The inverse method is most useful when several systems share the same coefficient matrix. The inverse is computed once, and each new right-hand side costs only a matrix product.
