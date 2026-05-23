# Systems of Differential Equations

A system of first-order linear ODEs has the form:

$$x_1' = a_{11}x_1 + a_{12}x_2 + \cdots + a_{1n}x_n + f_1(t)$$
$$x_2' = a_{21}x_1 + a_{22}x_2 + \cdots + a_{2n}x_n + f_2(t)$$
$$\vdots$$
$$x_n' = a_{n1}x_1 + a_{n2}x_2 + \cdots + a_{nn}x_n + f_n(t)$$

In matrix form: $\mathbf{x}' = A\mathbf{x} + \mathbf{f}(t)$

## Homogeneous Systems ($\mathbf{f} = 0$)

Assume solutions of the form $\mathbf{x} = \mathbf{v}e^{\lambda t}$.

Substitute: $A\mathbf{v} = \lambda\mathbf{v}$, so $\lambda$ is an eigenvalue and $\mathbf{v}$ an eigenvector of $A$.

**Case 1: Real distinct eigenvalues**
General solution: $\mathbf{x} = C_1\mathbf{v}_1 e^{\lambda_1 t} + C_2\mathbf{v}_2 e^{\lambda_2 t}$

**Example 1:** Solve $\mathbf{x}' = \begin{pmatrix} 3 & 1 \\ 1 & 3 \end{pmatrix} \mathbf{x}$.

Eigenvalues: $\det(A - \lambda I) = (3-\lambda)^2 - 1 = \lambda^2 - 6\lambda + 8 = 0$
$\lambda_1 = 4$, $\lambda_2 = 2$

Eigenvectors:
$\lambda = 4$: $(A - 4I)\mathbf{v} = 0$, $\begin{pmatrix} -1 & 1 \\ 1 & -1 \end{pmatrix}\mathbf{v} = 0$, $\mathbf{v}_1 = \begin{pmatrix} 1 \\ 1 \end{pmatrix}$
$\lambda = 2$: $(A - 2I)\mathbf{v} = 0$, $\begin{pmatrix} 1 & 1 \\ 1 & 1 \end{pmatrix}\mathbf{v} = 0$, $\mathbf{v}_2 = \begin{pmatrix} 1 \\ -1 \end{pmatrix}$

Solution: $\mathbf{x} = C_1\begin{pmatrix} 1 \\ 1 \end{pmatrix}e^{4t} + C_2\begin{pmatrix} 1 \\ -1 \end{pmatrix}e^{2t}$

**Case 2: Complex eigenvalues** $\lambda = a \pm bi$
The real solution uses $e^{at}(\cos bt + i\sin bt)$.

**Case 3: Repeated eigenvalues**
Need generalized eigenvectors for a second independent solution.

## Nonhomogeneous Systems

Use variation of parameters:
$$\mathbf{x}(t) = X(t)X^{-1}(0)\mathbf{x}_0 + X(t)\int_0^t X^{-1}(s)\mathbf{f}(s) ds$$
where $X(t)$ is a fundamental matrix whose columns are solutions of the homogeneous system.
