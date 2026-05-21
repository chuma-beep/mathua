> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linear_transformation) — CC BY-SA 4.0

# Linear transformations

A **linear transformation** (also called a linear map) is a function \(T: V \to W\) between vector spaces that preserves vector addition and scalar multiplication:

\[
T(u + v) = T(u) + T(v)
\]
\[
T(cv) = c\,T(v)
\]

### Matrix representation

Every linear transformation between finite-dimensional vector spaces can be represented by a matrix. If \(T: \mathbb{R}^n \to \mathbb{R}^m\) and \(\{e_1, \dots, e_n\}\) is the standard basis, then the matrix of \(T\) has columns \(T(e_1), \dots, T(e_n)\).

### Key examples

- **Rotation:** rotates vectors by an angle \(\theta\)
- **Reflection:** mirrors vectors across a line or plane
- **Scaling:** stretches or compresses vectors by a factor
- **Projection:** maps vectors onto a subspace
- **Shear:** shifts coordinates proportionally to another coordinate

The **kernel** (null space) of \(T\) is \(\{v \in V : T(v) = 0\}\). The **image** (range) is \(\{T(v) : v \in V\}\). The rank-nullity theorem states \(\dim(\ker T) + \dim(\operatorname{im} T) = \dim V\).
