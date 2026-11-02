# Homogeneous Linear Systems

A homogeneous system has the form \(A\mathbf{x}=\mathbf{0}\). The zero vector always solves it, so the interesting question is whether there are other solutions.

## Trivial and Nontrivial Solutions

The solution \(\mathbf{x}=\mathbf{0}\) is called the trivial solution. A nontrivial solution exists exactly when the columns of \(A\) are linearly dependent. For a square matrix this happens precisely when \(\det A=0\).

For example, the system

\[
x+2y=0,\qquad 3x+6y=0
\]

has determinant \(1\cdot 6-2\cdot 3=0\), so it has infinitely many solutions. Setting \(y=t\) gives \(x=-2t\).

## Free Variables and the Null Space

The set of all solutions of \(A\mathbf{x}=\mathbf{0}\) is the null space of \(A\). Its dimension is the nullity,

\[
\operatorname{nullity}(A)=n-\operatorname{rank}(A),
\]

where \(n\) is the number of columns. Each free variable contributes one parameter to the solution set.

If the rank equals the number of columns, the null space contains only the zero vector. If the rank is smaller, the null space contains a line, a plane or a higher-dimensional subspace.

## Solving with Row Reduction

To solve a homogeneous system, reduce its augmented matrix to row echelon form. The right-hand column stays zero throughout, so it is enough to reduce \(A\) itself.

Pivot columns correspond to basic variables. Nonpivot columns correspond to free variables. Express each basic variable in terms of the free variables, then write the solution as a linear combination of vectors.

For the system above, the reduced matrix has one pivot and one free variable. The null space is the line spanned by \((-2,1)\).
