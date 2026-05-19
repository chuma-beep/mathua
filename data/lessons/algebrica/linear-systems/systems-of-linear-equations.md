> Content sourced from [Algebrica](https://algebrica.org/systems-of-linear-equations/) — CC BY-NC 4.0

## What is a linear system?

Linear systems model problems where multiple conditions must be satisfied at the same time. They form the basis of many solution methods in algebra and applied mathematics. Given \\( n \\) variables \\( x_1, x_2, \dots, x_n \\), a system is called linear if all the equations are [linear equations](<../linear-equations>), meaning each variable appears to the first power, with no products between variables. The standard form of a linear system with \\( m \\) equations and \\( n \\) unknowns is written as:

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + \dots + a_{1n}x_n = b_1 \\\\[0.5em] a_{21}x_1 + a_{22}x_2 + \dots + a_{2n}x_n = b_2 \\\\[0.5em] \quad\vdots \\\\[0.5em] a_{m1}x_1 + a_{m2}x_2 + \dots + a_{mn}x_n = b_m \end{cases} \\]

In a linear system written in standard form, the coefficients \\( a_{ij} \\) represent the value multiplying the variable \\( x_j \\) in the \\( i \\)-th equation.

## Omogeneous systems

Each \\( b_i \\) denotes the constant term (also called the known term) on the right-hand side of the \\( i \\)-th equation. If all constant terms are zero, that is, \\( \forall i,\ b_i = 0 \\), the system is called **homogeneous**.

\\[\begin{cases} a_{11}x_1 + a_{12}x_2 + \dots + a_{1n}x_n = 0 \\\\[0.5em] a_{21}x_1 + a_{22}x_2 + \dots + a_{2n}x_n = 0 \\\\[0.5em] \quad\vdots \\\\[0.5em] a_{m1}x_1 + a_{m2}x_2 + \dots + a_{mn}x_n = 0 \end{cases} \\]

## Solutions

Solving a linear system means finding an ordered \\( n \\)-tuple of values denoted as \\( (s_1, s_2, \dots, s_n) \\), that satisfies all the equations in the system.

  * A system is said to be **consistent** (or possible) if at least one solution exists, and its equations are compatible.
  * If no solution exists, the system is called **inconsistent** (or impossible), and its equations are incompatible.
  * If there is only one solution, the system is called **determined**.
  * If there are infinitely many solutions, it is called **undetermined**.


Linear systems are defined by their coefficients and constants, which can be naturally organized into [matrices](<../matrices>). This connection allows us to represent the system compactly and apply matrix methods to analyze and solve it.

Any linear system in standard form can be represented by an \\( m \times n \\) matrix, where \\( m \\) is the number of equations and \\( n \\) is the number of variables. The matrix formed by the coefficients of the variables is called the **coefficient matrix** :

\\[A = \begin{bmatrix} a_{11} & a_{12} & \cdots & a_{1n} \\\\[0.5em] a_{21} & a_{22} & \cdots & a_{2n} \\\\[0.5em] \vdots & \vdots & \ddots & \vdots \\\\[0.5em] a_{m1} & a_{m2} & \cdots & a_{mn} \end{bmatrix} \\]

The constant terms and variables can be organized into two column [vectors](<../vectors/>):

\\[X = \begin{bmatrix} x_1 \\\ x_2 \\\ \vdots \\\ x_n \end{bmatrix} \quad\quad B = \begin{bmatrix} b_1 \\\ b_2 \\\ \vdots \\\ b_m \end{bmatrix} \\]

Therefore, the system can be rewritten in matrix form as:

\\[A \cdot X = B \\]

##### This compact form summarizes the entire system, where \\( A \\) is the coefficient matrix, \\( X \\) is the vector of variables, and \\( B \\) is the vector of constants.

## Why is the matrix form important?

Solving systems of linear equations can become challenging, especially as the number of equations and variables increases. While we’ll explore various solution methods, it’s worth noting that rewriting a system in matrix form is particularly useful. It provides a more compact representation and often simplifies the process of finding solutions.

## How to solve a linear system with the same number of equations and unknowns (\\( n = m \\))

A linear system with \\( n \\) equations and \\( n \\) unknowns can be solved using the [inverse matrix](<../inverse-matrix/>) method. If the coefficient matrix \\( A \\) is non-singular, that is, if its determinant is nonzero the system \\(A \cdot X = B\\) has a unique solution, given by:

\\[X = A^{-1}B \\]

## Example

Let’s solve a relatively simple example: a linear system with 3 equations in 3 unknowns, where \\( n = m \\).

\\[\begin{cases} x_1 + x_2 + x_3 = 3 \\\\[0.5em] 2x_1 + x_2 + x_3 = 4 \\\\[0.5em] 2x_1 + x_2 + 3x_3 = 8 \end{cases} \\]


First, let’s determine the coefficient matrix \\( A \\) and compute the determinant:

\\[A = \begin{bmatrix} 1 & 1 & 1 \\\ 2 & 1 & 1 \\\ 2 & 1 & 3 \end{bmatrix} \quad \det(A) = -2 \\]

##### Refer to the section on [computing the determinant](<../determinant/>) of a square matrix for a clearer understanding of the method.


Since the determinant of the matrix is nonzero, the matrix is non-singular and its inverse can be computed.

\\[A^{-1} = \begin{bmatrix} -\dfrac{1}{2} & -\dfrac{1}{2} & \dfrac{1}{2} \\\ 0 & -\dfrac{1}{2} & \dfrac{1}{2} \\\ \dfrac{1}{2} & 1 & -\dfrac{1}{2} \end{bmatrix} \\]

##### Refer to the section on computing the [inverse matrix](<../inverse-matrix>) for all the steps.


We now express the system as \\( X = A^{-1}B \\) in order to determine the unknown variables:

\\[\begin{bmatrix} x_1 \\\ x_2 \\\ x_3 \end{bmatrix} = \begin{bmatrix} -\dfrac{1}{2} & -\dfrac{1}{2} & \dfrac{1}{2} \\\ 0 & -\dfrac{1}{2} & \dfrac{1}{2} \\\ \dfrac{1}{2} & 1 & -\dfrac{1}{2} \end{bmatrix} \cdot \begin{bmatrix} 3 \\\ 4 \\\ 8 \end{bmatrix} \\]


We now compute the values of \\( x_1 \\), \\( x_2 \\), and \\( x_3 \\) from the matrix equation.

\\[\begin{aligned} x_1 &= -\dfrac{1}{2} \cdot 3 - \dfrac{1}{2} \cdot 4 + \dfrac{1}{2} \cdot 8 = \dfrac{1}{2} \\\\[0.5em] x_2 &= 0 \cdot 3 - \dfrac{1}{2} \cdot 4 + \dfrac{1}{2} \cdot 8 = 2 \\\\[0.5em] x_3 &= \dfrac{1}{2} \cdot 3 + 1 \cdot 4 - \dfrac{1}{2} \cdot 8 = \dfrac{3}{2} \end{aligned} \\]

Therefore, the solutions of the system are:

\\[x_1 = \dfrac{1}{2} \quad\quad x_2 = 2 \quad\quad x_3 = \dfrac{3}{2} \\]

Linear Systems

A linear system is a set of equations solved simultaneously by values that satisfy all of them.

3.4k

[Cramer’s Rule](https://algebrica.org/cramers-rule/)

1.1k

[Gaussian Elimination](https://algebrica.org/solving-linear-systems-using-gaussian-elimination/)

1.2k

[Rouché-Capelli Theorem](https://algebrica.org/rouche-capelli-theorem/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/linear-systems/rouche-capelli-theorem.md?plain=1)
