# Eigenspaces and Algebraic Multiplicity

An eigenvalue tells us that some nonzero vector is only stretched by a matrix. The set of all such vectors, together with zero, is a subspace. Its size measures how much of the space the eigenvalue actually controls.

## The Eigenspace of an Eigenvalue

For an eigenvalue $\lambda$ of $A$, the eigenspace is

$$E_\lambda = \ker(A - \lambda I).$$

It contains the zero vector and every eigenvector for $\lambda$. Because it is a null space, it is a subspace. Solving $(A - \lambda I)x = 0$ gives a basis, and the number of free variables is its dimension.

## Geometric Multiplicity

The geometric multiplicity of $\lambda$ is $\dim E_\lambda$, equivalently $n - \operatorname{rank}(A - \lambda I)$. It is the number of independent eigenvectors available for that eigenvalue. It is always at least one when $\lambda$ is an eigenvalue.

## Algebraic versus Geometric Multiplicity

The algebraic multiplicity is the power of $(t - \lambda)$ in the characteristic polynomial. The geometric multiplicity can never exceed it. A matrix is diagonalisable exactly when the geometric multiplicities add up to $n$, which happens exactly when every eigenvalue has geometric multiplicity equal to its algebraic multiplicity.
