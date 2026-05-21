> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Cramer%27s_rule) — CC BY-SA 4.0

# Cramer's rule

In linear algebra, **Cramer's rule** is an explicit formula for the solution of a system of linear equations with as many equations as unknowns, valid whenever the system has a unique solution. It expresses the solution in terms of the determinants of the (square) coefficient matrix and of matrices obtained from it by replacing one column by the column vector of right-sides of the equations. It is named after Gabriel Cramer, who published the rule for an arbitrary number of unknowns in 1750, although Colin Maclaurin also published special cases of the rule in 1748, and possibly knew of it as early as 1729.

Cramer's rule, implemented in a naive way, is computationally inefficient for systems of more than two or three equations. In the case of equations in unknowns, it requires computation of *n* + 1 determinants, while Gaussian elimination produces the result with the same (up to a constant factor independent of ) computational complexity as the computation of a single determinant. Moreover, the Bareiss algorithm is a simple modification of Gaussian elimination that produces in a single computation a matrix whose nonzero entries are the determinants involved in Cramer's rule.
In 1983, an algorithm for solving the system using Cramer's rule in \(n^3 + O(n^2)\) operations was proposed.
This algorithm can use permutations exactly like those in the Gaussian method. Therefore, if approximate calculation methods are used, the solution will be more stable than the solution using the Gaussian method; since the first steps are calculated without error, before rounding off the numbers.
In 1991, an algorithm with complexity of \(n^3/2 + O(n^2)\) operations was published.. In 1997, an algorithm for solving a system using Cramer's rule with complexity equal to that of matrix multiplication was proposed.. For example, for the Strassen's Multiplication Algorithm, this algorithm computes the solution in \((7/15)n^{\log_2(7)}+ O(n^2)\) operations.

## General case
Consider a system of linear equations for unknowns, represented in matrix multiplication form as follows:

\(A\mathbf{x} = \mathbf{b}\)

where the *n* × *n* matrix has a nonzero determinant, and the vector \(\mathbf{x} = (x_1, \ldots, x_n)^\mathsf{T}\) is the column vector of the variables. Then the theorem states that in this case the system has a unique solution, whose individual values for the unknowns are given by:

\(x_i = \frac{\det(A_i)}{\det(A)} \qquad i = 1, \ldots, n\)

where \(A_i\) is the matrix formed by replacing the -th column of by the column vector **b**.

A more general version of Cramer's rule considers the matrix equation

\(AX = B\)

where the *n* × *n* matrix has a nonzero determinant, and , are *n* × *m* matrices. Given sequences \(1 \leq i_1 < i_2 < \cdots < i_k \leq n\) and \(1 \leq j_1 < j_2 < \cdots < j_k \leq m\), let \(X_{I,J}\) be the *k* × *k* submatrix of with rows in \(I := (i_1, \ldots, i_k )\) and columns in \(J := (j_1, \ldots, j_k )\). Let \(A_{B}(I,J)\) be the *n* × *n* matrix formed by replacing the \(i_s\) column of by the \(j_s\) column of , for all \(s = 1,\ldots, k\). Then

\(\det X_{I,J} = \frac{\det(A_{B}(I,J))}{\det(A)}.\)

In the case \(k = 1\), this reduces to the normal Cramer's rule.

The rule holds for systems of equations with coefficients and unknowns in any field, not just in the real numbers.

## Proof
The proof for Cramer's rule uses the following properties of the determinants: linearity with respect to any given column and the fact that the determinant is zero whenever two columns are equal, which is implied by the property that the sign of the determinant flips if you switch two columns.

Fix the index *j* of a column, and consider that the entries of the other columns have fixed values. This makes the determinant a function of the entries of the th column. Linearity with respect to this column means that this function has the form
\(D_j(a_{1,j}, \ldots, a_{n,j})= C_{1,j}a_{1,j}+\cdots, C_{n,j}a_{n,j},\)
where the \(C_{i,j}\) are coefficients that depend on the entries of that are not in column. So, one has
\(\det(A)=D_j(a_{1,j}, \ldots, a_{n,j})=C_{1,j}a_{1,j}+\cdots, C_{n,j}a_{n,j}\)
(Laplace expansion provides a formula for computing the \(C_{i,j}\) but their expression is not important here.)

If the function \(D_j\) is applied to any *other* column *k* of , then the result is the determinant of the matrix obtained from by replacing column *j* by a copy of column *k*, so the resulting determinant is 0 (the case of two equal columns).

Now consider a system of linear equations in unknowns \(x_1, \ldots,x_n\), whose coefficient matrix with det(*A*) assumed to be nonzero:

\(\begin{matrix}
a_{11}x_1+a_{12}x_2+\cdots+a_{1n}x_n&=&b_1\\
a_{21}x_1+a_{22}x_2+\cdots+a_{2n}x_n&=&b_2\\
&\vdots&\\
a_{n1}x_1+a_{n2}x_2+\cdots+a_{nn}x_n&=&b_n.
\end{matrix}\)

If one combines these equations by taking *C*\(_{1,*j*}\) times the first equation, plus *C*\(_{2,*j*}\) times the second, and so forth until *C*\(_{*n*,*j*}\) times the last, then for every the resulting coefficient of becomes
\(D_j(a_{1,k},\ldots,a_{n,k}).\)
So, all coefficients become zero, except the coefficient of \(x_j\) that becomes \(\det(A).\) Similarly, the constant coefficient becomes \(D_j(b_1,\ldots,b_n),\) and the resulting equation is thus
\(\det(A)x_j=D_j(b_1,\ldots, b_n),\)
which gives the value of \(x_j\) as
\(x_j=\frac1{\det(A)}D_j(b_1,\ldots, b_n).\)

As, by construction, the numerator is the determinant of the matrix obtained from by replacing column *j* by **b**, we get the expression of Cramer's rule as a necessary condition for a solution.

It remains to prove that these values for the unknowns form a solution. Let be the *n* × *n* matrix that has the coefficients of \(D_j\) as th row, for \(j=1,\ldots,n\) (this is the adjugate matrix for ). Expressed in matrix terms, we have thus to prove that
\(\mathbf x = \frac1{\det(A)}M\mathbf b\)
is a solution; that is, that
\(A\left(\frac1{\det(A)}M\right)\mathbf b=\mathbf b.\)
For that, it suffices to prove that
\(A\,\left(\frac1{\det(A)}M\right)=I_n,\)
where \(I_n\) is the identity matrix.

The above properties of the functions \(D_j\) show that one has *MA* , and therefore,
\(\left(\frac1{\det(A)}M\right)\,A=I_n.\)
This completes the proof; since a left inverse of a square matrix is also a right-inverse (see Invertible matrix theorem).

For other proofs, see below.

## Finding inverse matrix

Let be an *n* × *n* matrix with entries in a field *F*. Then

\(A\,\operatorname{adj}(A) = \operatorname{adj}(A)\,A=\det(A) I\)

where adj(*A*) denotes the adjugate matrix, det(*A*) is the determinant, and *I* is the identity matrix. If det(*A*) is nonzero, then the inverse matrix of is

\(A^{-1} = \frac{1}{\det(A)} \operatorname{adj}(A).\)

This gives a formula for the inverse of , provided det(*A*) ≠ 0. In fact, this formula works whenever *F* is a commutative ring, provided that det(*A*) is a unit. If det(*A*) is not a unit, then is not invertible over the ring (it may be invertible over a larger ring in which some non-unit elements of may be invertible).

## Applications
### Explicit formulas for small systems
Consider the linear system

\(\left\{\begin{matrix}
a_1x + b_1y&= {\color{red}c_1}\\
a_2x + b_2y&= {\color{red}c_2}
\end{matrix}\right.\)

which in matrix format is

\(\begin{bmatrix} a_1 & b_1 \\ a_2 & b_2 \end{bmatrix}\begin{bmatrix} x \\ y \end{bmatrix}=\begin{bmatrix} {\color{red}c_1} \\ {\color{red}c_2} \end{bmatrix}.\)

Assume *a*\(_{1}\)*b*\(_{2}\) − *b*\(_{1}\)*a*\(_{2}\) is nonzero. Then, with the help of determinants, and can be found with Cramer's rule as

\(\begin{align}
x &= \frac{\begin{vmatrix} {\color{red}{c_1& b_1 \\ {\color{red}{c_2& b_2 \end{vmatrix{\begin{vmatrix} a_1 & b_1 \\ a_2 & b_2 \end{vmatrix= { {\color{red}c_1}b_2 - b_1{\color{red}c_2} \over a_1b_2 - b_1a_2}, \quad
y = \frac{\begin{vmatrix} a_1 & {\color{red}{c_1\\ a_2 & {\color{red}{c_2\end{vmatrix{\begin{vmatrix} a_1 & b_1 \\ a_2 & b_2 \end{vmatrix = { a_1{\color{red}c_2} - {\color{red}c_1}a_2 \over a_1b_2 - b_1a_2}
\end{align}.\)

The rules for 3 × 3 matrices are similar. Given

\(\left\{\begin{matrix}
a_1x + b_1y + c_1z&= {\color{red}d_1}\\
a_2x + b_2y + c_2z&= {\color{red}d_2}\\
a_3x + b_3y + c_3z&= {\color{red}d_3}
\end{matrix}\right.\)

which in matrix format is

\(\begin{bmatrix} a_1 & b_1 & c_1 \\ a_2 & b_2 & c_2 \\ a_3 & b_3 & c_3 \end{bmatrix}\begin{bmatrix} x \\ y \\ z \end{bmatrix}=\begin{bmatrix} {\color{red}d_1} \\ {\color{red}d_2} \\ {\color{red}d_3} \end{bmatrix}.\)

Then the values of and can be found as follows:

\(x = \frac{\begin{vmatrix} {\color{red}d_1} & b_1 & c_1 \\ {\color{red}d_2} & b_2 & c_2 \\ {\color{red}d_3} & b_3 & c_3 \end{vmatrix} } { \begin{vmatrix} a_1 & b_1 & c_1 \\ a_2 & b_2 & c_2 \\ a_3 & b_3 & c_3 \end{vmatrix, \quad
y = \frac {\begin{vmatrix} a_1 & {\color{red}d_1} & c_1 \\ a_2 & {\color{red}d_2} & c_2 \\ a_3 & {\color{red}d_3} & c_3 \end{vmatrix{\begin{vmatrix} a_1 & b_1 & c_1 \\ a_2 & b_2 & c_2 \\ a_3 & b_3 & c_3 \end{vmatrix, \quad
z = \frac { \begin{vmatrix} a_1 & b_1 & {\color{red}d_1} \\ a_2 & b_2 & {\color{red}d_2} \\ a_3 & b_3 & {\color{red}d_3} \end{vmatrix{\begin{vmatrix} a_1 & b_1 & c_1 \\ a_2 & b_2 & c_2 \\ a_3 & b_3 & c_3 \end{vmatrix} }.\)

### Differential geometry
#### Ricci calculus
Cramer's rule is used in the Ricci calculus in various calculations involving the Christoffel symbols of the first and second kind.

In particular, Cramer's rule can be used to prove that the divergence operator on a Riemannian manifold is invariant with respect to change of coordinates. We give a direct proof, suppressing the role of the Christoffel symbols.
Let \((M,g)\) be a Riemannian manifold equipped with local coordinates \((x^1, x^2, \dots, x^n)\). Let \(A=A^i \frac{\partial}{\partial x^i}\) be a vector field. We use the summation convention throughout.

**Theorem**.
*The *divergence* of \(A\),*
\(\operatorname{div} A = \frac{1}{\sqrt{\det g\frac{\partial}{\partial x^i} \left( A^i \sqrt{\det g} \right),\)
*is invariant under change of coordinates.*

Let \((x^1,x^2,\ldots,x^n)\mapsto (\bar x^1,\ldots,\bar x^n)\) be a coordinate transformation with non-singular Jacobian. Then the classical transformation laws imply that \(A=\bar A^{k}\frac{\partial}{\partial\bar x^{k\) where \(\bar A^{k}=\frac{\partial \bar x^{k{\partial x^{jA^{j}\). Similarly, if \(g=g_{mk}\,dx^{m}\otimes dx^{k}=\bar{g}_{ij}\,d\bar x^{i}\otimes d\bar x^{j}\), then \(\bar{g}_{ij}=\,\frac{\partial x^{m{\partial\bar x^{i\frac{\partial x^{k{\partial \bar x^{jg_{mk}\).
Writing this transformation law in terms of matrices yields \(\bar g=\left(\frac{\partial x}{\partial\bar{x\right)^{\text{Tg\left(\frac{\partial x}{\partial\bar{x\right)\), which implies \(\det\bar g=\left(\det\left(\frac{\partial x}{\partial\bar{x\right)\right)^{2}\det g\).

Now one computes
\(\begin{align}
\operatorname{div} A &=\frac{1}{\sqrt{\det g\frac{\partial}{\partial x^{i\left( A^{i}\sqrt{\det g}\right)\\
	&=\det\left(\frac{\partial x}{\partial\bar{x\right)\frac{1}{\sqrt{\det\bar g\frac{\partial \bar x^k}{\partial x^{i\frac{\partial}{\partial\bar x^{k\left(\frac{\partial x^{i{\partial \bar x^{\ell\bar{A}^{\ell}\det\!\left(\frac{\partial x}{\partial\bar{x\right)^{\!\!-1}\!\sqrt{\det\bar g}\right).
\end{align}\)
In order to show that this equals
\(\frac{1}{\sqrt{\det\bar g\frac{\partial}{\partial\bar x^{k\left(\bar A^{k}\sqrt{\det\bar{g\right)\),
it is necessary and sufficient to show that
\(\frac{\partial\bar x^{k{\partial x^{i\frac{\partial}{\partial\bar x^{k\left(\frac{\partial x^{i{\partial \bar x^{\ell\det\!\left(\frac{\partial x}{\partial\bar{x\right)^{\!\!\!-1}\right)=0\qquad\text{for all } \ell,\)
which is equivalent to
\(\frac{\partial}{\partial \bar x^{\ell\det\left(\frac{\partial x}{\partial\bar{x\right)
=\det\left(\frac{\partial x}{\partial\bar{x\right)\frac{\partial\bar x^{k{\partial x^{i\frac{\partial^{2}x^{i{\partial\bar x^{k}\partial\bar x^{\ell.\)
Carrying out the differentiation on the left-hand side, we get:
\(\begin{align}
	\frac{\partial}{\partial\bar x^{\ell\det\left(\frac{\partial x}{\partial\bar{x\right)
	&=(-1)^{i+j}\frac{\partial^{2}x^{i{\partial\bar x^{\ell}\partial\bar x^{j\det M(i|j)\\
	&=\frac{\partial^{2}x^{i{\partial\bar x^{\ell}\partial\bar x^{j\det\left(\frac{\partial x}{\partial\bar{x\right)\frac{(-1)^{i+j}\det\left(\frac{\partial x}{\partial\bar{x\right)\det M(i|j)=(\ast),
	\end{align}\)
where \(M(i|j)\) denotes the matrix obtained from \(\left(\frac{\partial x}{\partial\bar{x\right)\) by deleting the \(i\)th row and \(j\)th column.
But Cramer's Rule says that
\(\frac{(-1)^{i+j}\det\left(\frac{\partial x}{\partial\bar{x\right)\det M(i|j)\)
is the \((j,i)\)th entry of the matrix \(\left(\frac{\partial \bar{x{\partial x}\right)\).
Thus
\((\ast)=\det\left(\frac{\partial x}{\partial\bar{x\right)\frac{\partial^{2}x^{i{\partial\bar x^{\ell}\partial\bar x^{j\frac{\partial\bar x^{j{\partial x^{i,\)
completing the proof.

#### Computing derivatives implicitly
Consider the two equations \(F(x, y, u, v) = 0\) and \(G(x, y, u, v) = 0\). When *u* and *v* are independent variables, we can define \(x = X(u, v)\) and \(y = Y(u, v).\)

An equation for \(\dfrac{\partial x}{\partial u}\) can be found by applying Cramer's rule.

-
