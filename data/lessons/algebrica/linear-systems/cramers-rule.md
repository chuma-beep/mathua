> Content sourced from [Algebrica](https://algebrica.org/cramers-rule/) — CC BY-NC 4.0

## What is Cramer’s rule?

Cramer’s Rule provides a method for solving [systems](<../systems-of-linear-equations>) of \\( n \\) [linear equations](<../linear-equations>) in \\( n \\) unknowns, by using the determinant of the system’s coefficient [matrix](<../matrices>). This rule applies only when the coefficient matrix is square and its determinant is non-zero, ensuring that the system has a unique solution.

## Applying Cramer’s Rule

Let us consider a general system of \\( n \\) equations in \\( n \\) unknowns:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + \dots + a_{1n}x_n = b_1 \\\\[0.5em] a_{21}x_1 + a_{22}x_2 + \dots + a_{2n}x_n = b_2 \\\\[0.5em] \quad\vdots \\\\[0.5em] a_{n1}x_1 + a_{n2}x_2 + \dots + a_{nn}x_n = b_n \end{cases} \\]

We can rewrite the system using matrix notation as \\( A \cdot \mathbf{X} = \mathbf{B} \\). The system becomes:

\\[A = \begin{bmatrix} a_{11} & a_{12} & \cdots & a_{1n} \\\ a_{21} & a_{22} & \cdots & a_{2n} \\\ \vdots & \vdots & \ddots & \vdots \\\ a_{m1} & a_{n2} & \cdots & a_{nn} \end{bmatrix} \\]

The constant terms and variables can be organized into two column [vectors](<../vectors/>):

\\[X = \begin{bmatrix} x_1 \\\ x_2 \\\ \vdots \\\ x_n \end{bmatrix} \quad\quad B = \begin{bmatrix} b_1 \\\ b_2 \\\ \vdots \\\ b_m \end{bmatrix} \\]


Suppose that the coefficient matrix \\( A \\) is invertible, which means that its determinant is non-zero (\\(\det(A) \ne 0\\)). Under this condition, the system has a unique solution, and it can be expressed using the inverse of \\( A \\) as:

\\[\begin{bmatrix} x_1 \\\ x_2 \\\ \vdots \\\ x_n \end{bmatrix} = \frac{1}{\det(A)} \begin{bmatrix} A_{11} & A_{21} & \cdots & A_{n1} \\\ A_{12} & A_{22} & \cdots & A_{n2} \\\ \vdots & \vdots & \ddots & \vdots \\\ A_{1n} & A_{2n} & \cdots & A_{nn} \\\ \end{bmatrix} \begin{bmatrix} b_1 \\\ b_2 \\\ \vdots \\\ b_n \end{bmatrix} \\]

This is the general form of Cramer’s Rule, where each \\( A_{ij} \\) is the cofactor of the element \\( a_{ji} \\) in the original matrix \\( A \\).

##### In this context, \\( A_{ij} \\) refers to the cofactor of the element \\( a_{ji} \\) from the original matrix \\( A \\), not to the entry of the matrix itself. The matrix used here is the adjugate of \\( A \\), denoted as \\( \text{adj}(A) \\).


The value of the unknown in position \\( k \\) is given by a fraction. Its denominator is \\( \det(A) \\) and its numerator is the determinant of the matrix obtained by replacing the \\( k \\)-th column of \\( A \\) with the column of constants:

\\[x_k = \frac{\det(A_k)}{\det(A)} \\]

##### For a clearer understanding of this step, see the detailed explanation in Example 2.

## Solutions of homogeneous systems

A [homogeneous system](<../systems-of-linear-equations>) is a system of linear equations where all the constant terms are zero. These systems always have at least one solution: the trivial solution, where all the variables are zero. But something interesting happens when we look at the determinant of the coefficient matrix:

  * If \\( \det(A) \ne 0 \\), the system has only the trivial solution.
  * If \\( \det(A) = 0 \\), the system admits infinitely many solutions, including non-trivial ones, where at least one variable is not zero.


##### Homogeneous systems never have no solution. They’re always consistent, but the number of solutions depends entirely on the determinant.

## Example 1

Let’s consider the following homogeneous system of three equations in three unknowns:

\\[\begin{cases} x + y + z = 0 \\\\[0.5em] 2x - y + z = 0 \\\\[0.5em] 3x + y + 2z = 0 \end{cases} \\]


This system can be written in matrix form as \\( A \cdot \mathbf{x} = \mathbf{0} \\), where the coefficient matrix is:

\\[A = \begin{bmatrix} 1 & 1 & 1 \\\\[0.5em] 2 & -1 & 1 \\\\[0.5em] 3 & 1 & 2 \end{bmatrix} \\]


To understand the nature of the solutions, we compute the determinant of the matrix \\( A \\):

\\[\begin{align*} \det(A) &= 1 \cdot (-1 \cdot 2 - 1 \cdot 1) -1 \cdot (2 \cdot 2 - 1 \cdot 3) +1 \cdot (2 \cdot 1 - (-1) \cdot 3) \\\\[0.5em] &= 1(-2 - 1) - 1(4 - 3) + 1(2 + 3) \\\\[0.5em] &= -3 - 1 + 5 \\\\[0.5em] & = 1 \end{align*} \\]

Since the determinant is non-zero, the system admits only the trivial solution:

\\[x = 0 \quad y = 0 \quad z = 0\\]

##### This outcome aligns with Cramer’s Rule: when \\( \det(A) \ne 0 \\), the only possible solution to a homogeneous system is the trivial one, since all the determinants in the numerators of Cramer’s formula become zero.

## Example 2

Let’s solve the following system of two linear equations in two unknowns:

\\[\begin{cases} 2x + 3y = 8 \\\\[0.5em] 4x - y = 2 \end{cases} \\]


We identify the coefficient matrix \\( A \\), the vector of unknowns \\( \mathbf{x} \\), and the constants vector \\( \mathbf{b} \\):

\\[A = \begin{bmatrix} 2 & 3 \\\ 4 & -1 \end{bmatrix} \quad\quad \mathbf{x} = \begin{bmatrix} x \\\ y \end{bmatrix} \quad\quad \mathbf{b} = \begin{bmatrix} 8 \\\ 2 \end{bmatrix} \\]


We compute the determinant of the coefficient matrix \\( A \\):

\\[\det(A) = 2 \cdot (-1) - 3 \cdot 4 = -2 - 12 = -14 \\]

Since the determinant is non-zero, the system has exactly one solution, and we can apply Cramer’s rule.


We build the matrix \\( A_1 \\) by replacing the first column of \\( A \\) with the constants:

\\[A_1 = \begin{bmatrix} 8 & 3 \\\ 2 & -1 \end{bmatrix} \quad \rightarrow \quad \det(A_1) = 8 \cdot (-1) - 3 \cdot 2 = -8 - 6 = -14 \\]


We now build \\( A_2 \\) by replacing the second column of \\( A \\) with the constants:

\\[A_2 = \begin{bmatrix} 2 & 8 \\\ 4 & 2 \end{bmatrix} \quad \rightarrow \quad \det(A_2) = 2 \cdot 2 - 8 \cdot 4 = 4 - 32 = -28 \\]


Now we compute the values of the unknowns using the formula:

\begin{align*} x &= \frac{\det(A_1)}{\det(A)} = \frac{-14}{-14} = 1 \\\\[0.5em] y &= \frac{\det(A_2)}{\det(A)} = \frac{-28}{-14} = 2 \end{align*}

The solution to the system is:

\\[x = 1 \quad\quad y = 2 \\]

##### Remember that the solution to a system of linear equations refers to the \\(n\\)-tuple of values that satisfies all equations in the system simultaneously. In this case, the pair \\( (x, y) = (1, 2) \\) is the only combination of values that makes both equations true at the same time.
