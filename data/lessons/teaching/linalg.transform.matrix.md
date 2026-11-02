# The Matrix of a Linear Transformation

A linear transformation is determined by what it does to a basis. Once the images of the basis vectors are known, every other input is fixed by linearity. A matrix is simply the record of those images in a chosen pair of bases.

## Images of Basis Vectors

Let $T$ be linear and let $e_1, \ldots, e_n$ be the standard basis of $\mathbb{R}^n$. Every vector is a combination $x = x_1 e_1 + \cdots + x_n e_n$. Linearity gives

$$T(x) = x_1 T(e_1) + \cdots + x_n T(e_n).$$

So the whole transformation is known from the $n$ output vectors $T(e_j)$.

## Building the Standard Matrix

Place $T(e_j)$ in column $j$ of a matrix $A$. Then $T(x) = Ax$ for every $x$. For example, a rotation through angle $\theta$ sends $e_1$ to $(\cos\theta, \sin\theta)$ and $e_2$ to $(-\sin\theta, \cos\theta)$, so those two vectors are the columns of its matrix.

The same rule works for any basis. If $\mathcal{B}$ is a basis, the matrix of $T$ relative to $\mathcal{B}$ has the coordinate vectors of $T(b_j)$ as its columns.

## Matrices under a Change of Basis

Changing the basis changes the matrix but not the transformation. If $P$ converts $\mathcal{B}$-coordinates to standard coordinates, the matrix in the new basis is $P^{-1} A P$. Similar matrices represent the same linear map in different bases, so trace, determinant and eigenvalues are unchanged.
