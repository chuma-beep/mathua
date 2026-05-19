> Content sourced from [Algebrica](https://algebrica.org/solving-linear-systems-using-gaussian-elimination/) — CC BY-NC 4.0

## What is Gaussian elimination

The Gauss method, or **Gaussian elimination** , is a technique used to solve [systems](<../systems-of-linear-equations/>) of \\( n \\) linear equations in \\( n \\) unknowns. The process involves applying a sequence of operations iteratively to eliminate one variable at a time, transforming the system into a form that is easy to solve. Let us consider the following system of equations:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + a_{13}x_3 = b_1 \\\\[0.5em] a_{21}x_1 + a_{22}x_2 + a_{23}x_3 = b_2 \\\\[0.5em] a_{31}x_1 + a_{32}x_2 + a_{33}x_3 = b_3 \end{cases} \\]

##### To apply the method, it is essential that the system is square, that is, it must have the same number of equations and unknowns, like the \\( 3 \times 3 \\) system shown above.

## The process involves the following steps:

  * Eliminate the variable \\( x_1 \\) from all equations except the first one.
  * Eliminate the variable \\( x_2 \\) from the third equation.
  * Solve the third equation to find the value of \\( x_3 \\).
  * Work backwards to find the values of the remaining variables.


##### When the algorithm produces an inconsistent or indeterminate equation, the system cannot proceed further: this signals either no solution or an infinite set of solutions.

## Connection with Matrices

The Gaussian elimination method transforms a [matrix](<../matrices>) into its row-echelon form. The objective is to rewrite the matrix so that each pivot position (on the main diagonal) contains a 1, and all entries below each pivot are 0. For example, starting from a generic matrix:

\\[A = \begin{bmatrix} a_{11} & a_{12} & a_{13} \\\ a_{21} & a_{22} & a_{23} \\\ a_{31} & a_{32} & a_{33} \end{bmatrix} \quad \xrightarrow{\text{to}} \quad \begin{bmatrix} 1 & b_{12} & b_{13} \\\ 0 & 1 & b_{23} \\\ 0 & 0 & 1 \end{bmatrix} \\]

## Let’s start by eliminating the variable \\( x_1 \\) from all equations except the first one.

To eliminate \\( x_1 \\) from the second equation, we multiply the first equation by \\( a_{21} \\) and the second by \\( a_{11} \\). We then subtract the two resulting equations and replace the second equation with the result. We obtain:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + a_{13}x_3 = b_1 \\\\[0.5em] \left(a_{11}a_{22} - a_{21}a_{12}\right)x_2 + \left(a_{11}a_{23} - a_{21}a_{13}\right)x_3 = a_{11}b_2 - a_{21}b_1 \\\\[0.5em] a_{31}x_1 + a_{32}x_2 + a_{33}x_3 = b_3 \end{cases} \\]

To simplify the calculations, we rewrite the second equation as:

\\[a^\prime_{22} x_2 + a^\prime_{23} x_3 = b^\prime_2 \\]

The system becomes:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + a_{13}x_3 = b_1 \\\\[0.5em] a^\prime_{22} x_2 + a^\prime_{23} x_3 = b^\prime_2\\\\[0.5em] a_{31}x_1 + a_{32}x_2 + a_{33}x_3 = b_3 \end{cases} \\]


Now, we multiply the first equation by \\( a_{31} \\) and the third equation by \\( a_{11} \\). We then subtract the two resulting equations and replace the third equation with the result.

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + a_{13}x_3 = b_1 \\\\[0.5em] a^\prime_{22} x_2 + a^\prime_{23} x_3 = b^\prime_2\\\\[0.5em] \left(a_{11}a_{32} - a_{31}a_{12}\right)x_2 + \left(a_{11}a_{33} - a_{31}a_{13}\right)x_3 = a_{11}b_3 - a_{31}b_1 \end{cases} \\]

We rewrite the third equation as:

\\[a^\prime_{32}x_2 + a^\prime_{33}x_3 = b’_3 \\]

We obtain:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + a_{13}x_3 = b_1 \\\\[0.5em] a^\prime_{22} x_2 + a^\prime_{23} x_3 = b^\prime_2\\\\[0.5em] a^\prime_{32}x_2 + a^\prime_{33}x_3 = b’_3\end{cases} \\]

## We proceed with the second step and eliminate the variable \\( x_2 \\) from the third equation.

We multiply the second equation by \\( a^\prime_{32} \\) and the third equation by \\( a^\prime_{22} \\). We then subtract the two resulting equations and replace the third equation with the result. By completing the calculations as in the first step, the system reduces to the following:

\\[\begin{cases} a_{11}x_1 + a_{12} x_2 + a_{13} x_3 = b_1 \\\\[0.5em] a^\prime_{22} x_2 + a^\prime_{23} x_3 = b^\prime_2 \\\\[0.5em] a^{\prime\prime}_{33} x_3 = b^{\prime\prime}_3 \end{cases} \\]

We have obtained a triangular system. At this point, we can solve for \\( x_3 \\) from the third equation, obtaining: \\[x_3 = \frac{b^{\prime\prime_3}}{a^{\prime\prime_{33}}} \\]

We can now determine each variable, starting from the known value of \\(x_3\\).

## Example

Let’s walk through a concrete example of Gaussian elimination applied to a 3×3 system. Consider the following system of linear equations:

\\[\begin{cases} x + y + z = 6 \\\\[0.5em] 2x + 3y + z = 14 \\\\[0.5em] x + 2y + 3z = 14 \end{cases} \\]

To eliminate \\( x \\) from the second equation, subtract \\( 2 \times \\) equation 1 from equation 2:

\\[(2x + 3y + z) - 2(x + y + z) = 0x + y - z = 2 \\]


Now eliminate \\( x \\) from the third equation by subtracting equation 1 from equation 3:

\\[(x + 2y + 3z) - (x + y + z) = 0x + y + 2z = 8 \\]

The system becomes:

\\[\begin{cases} x + y + z = 6 \\\\[0.5em] y - z = 2 \\\\[0.5em] y + 2z = 8 \end{cases} \\]


Now use the second equation as the new pivot. To eliminate \\( y \\) from the third equation, subtract equation 2 from equation 3:

\\[(y + 2z) - (y - z) = 3z = 6 \Rightarrow z = 2 \\]

Substituting back to find the remaining variables. From equation 2: \\[y - z = 2 \rightarrow y = 2 + z = 2 + 2 = 4 \\]

From equation 1: \\[x + y + z = 6 \Rightarrow x = 6 - y - z = 6 - 4 - 2 = 0 \\]

The solution is

\\[x = 0 \quad y = 4 \quad z = 2 \\]

Linear Systems

A linear system is a set of equations solved simultaneously by values that satisfy all of them.

2k

[Systems of Linear Equations](https://algebrica.org/systems-of-linear-equations/)

3.4k

[Cramer’s Rule](https://algebrica.org/cramers-rule/)

1.2k

[Rouché-Capelli Theorem](https://algebrica.org/rouche-capelli-theorem/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/linear-systems/rouche-capelli-theorem.md?plain=1)
