# The Four Fundamental Subspaces

Every matrix $A$ with $m$ rows and $n$ columns determines four subspaces. Two live in $\mathbb{R}^n$ and two live in $\mathbb{R}^m$. Their dimensions are fixed by the rank, so they describe a matrix more completely than its entries do.

## Column Space and Null Space

The column space is the span of the columns of $A$. It consists of every vector $b$ for which $Ax = b$ has a solution, and its dimension is the rank $r$.

The null space is the set of vectors $x$ with $Ax = 0$. Its dimension is the nullity $n - r$. Row reduction identifies the pivot columns, which form a basis of the column space, and the free variables, which parameterise the null space.

## Row Space and Left Null Space

The row space is the span of the rows of $A$. Row operations do not change it, so the nonzero rows of an echelon form give a basis. Its dimension is also $r$.

The left null space is the set of vectors $y$ with $A^T y = 0$. It is the null space of $A^T$, so its dimension is $m - r$. It records the linear dependencies among the rows of $A$.

## Dimensions of the Four Subspaces

The four dimensions are $r$, $n - r$, $r$ and $m - r$. The column space and row space have the same dimension, the rank. The null space and left null space measure the freedom in the domain and the redundancy in the codomain.
