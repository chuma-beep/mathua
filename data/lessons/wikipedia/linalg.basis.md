> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Basis_%28linear_algebra%29) — CC BY-SA 4.0

# Basis and dimension

Basis}}


In mathematics, a set  of elements of  a vector space *V* is called a **basis** (: **bases**) if every element of *V* can be written in a unique way as a finite linear combination of elements of . The coefficients of this linear combination are referred to as **components** or **coordinates** of the vector with respect to . The elements of a basis are called ****.

Equivalently, a set  is a basis if its elements are linearly independent and every element of  is a linear combination of elements of . In other words, a basis is a linearly independent spanning set.

A vector space can have several bases; however all the bases have the same number of elements, called the dimension of the vector space.

This article deals mainly with finite-dimensional vector spaces. However, many of the principles are also valid for infinite-dimensional vector spaces.

Basis vectors find applications in the study of crystal structures and frames of reference.

## Definition
A **basis** *B* of a vector space *V* over a field *F* (such as the real numbers **R** or the complex numbers **C**) is a linearly independent subset of *V* that spans *V*. This means that a subset  of *V* is a basis if it satisfies the two following conditions:
* linear independence: for every finite subset \(\{\mathbf v_1, \dotsc, \mathbf v_m\}\) of , if \(c_1 \mathbf v_1 + \cdots + c_m \mathbf v_m = \mathbf 0\) for some \(c_1,\dotsc,c_m\) in , then \(c_1 = \cdots = c_m = 0\); and
* the spanning property: for every vector \(\mathbf{v}\) in \(V\), one can choose \(a_1,\dotsc,a_n\) in \(F\) and \(\mathbf v_1, \dotsc, \mathbf v_n\) in \(B\) such that \(\mathbf v = a_1 \mathbf v_1 + \cdots + a_n \mathbf v_n\). In other words, \(\mathbb{v}\) can be represented as a linear combination of some vectors in \(B\).

The first property may be equivalently phrased as follows: if a linear combination of vectors in  is equal to the zero vector, then all the scalars coefficients of the combination must be zero.

If  is a basis for , then every vector \(\mathbf{v}\) in  can be written as a linear combination of vectors in  (by the spanning property); it follows from linear independence that this can be done in exactly one way.  Thus, the scalar coefficients \(a_i\) that appear in this combination are uniquely determined; they are called the *coordinates* of \(\mathbf{v}\) with respect to the basis \(B\).

A vector space that has a finite basis is called finite-dimensional. In this case, the finite subset can be taken as \(B\) itself to check for linear independence in the above definition.

It is often convenient or even necessary to have an ordering on the basis vectors, for example, when discussing orientation, or when one considers the scalar coefficients of a vector with respect to a basis without referring explicitly to the basis elements. In this case, the ordering is necessary for associating each coefficient with the corresponding basis element. This ordering can be done by numbering the basis elements. In order to emphasize that an order has been chosen, one speaks of an **ordered basis**, which is therefore not simply an unstructured set, but a sequence, an indexed family, or similar; see  below.

## Examples

The set **R**<sup>2</sup> of the ordered pairs of real numbers is a vector space under the operations of component-wise addition

\[
(a, b) + (c, d) = (a + c, b+d)
\]

and scalar multiplication

\[
\lambda (a,b) = (\lambda a, \lambda b),
\]

where \(\lambda\) is any real number. A simple basis of this vector space consists of the two vectors 1=**e**<sub>1</sub> = (1, 0) and 1=**e**<sub>2</sub> = (0, 1).  These vectors form a basis (called the standard basis) because any vector 1=**v** = (*a*, *b*) of **R**<sup>2</sup> may be uniquely written as
\[
\mathbf v = a \mathbf e_1 + b \mathbf e_2.
\]
 Any other pair of linearly independent vectors of **R**<sup>2</sup>, such as (1, 1) and (−1, 2), forms also a basis of **R**<sup>2</sup>.

More generally, if  is a field, the set \(F^n\) of -tuples of elements of  is a vector space for similarly defined addition and scalar multiplication. Let
\[
\mathbf e_i = (0, \ldots, 0,1,0,\ldots, 0)
\]
 be the -tuple with all components equal to 0, except the th, which is 1. Then \(\mathbf e_1, \ldots, \mathbf e_n\) is a basis of \(F^n,\) which is called the *standard basis* of \(F^n.\)

A different flavor of example is given by polynomial rings.  If  is a field, the collection *F*[*X*] of all polynomials in one indeterminate  with coefficients in  is an -vector space.  One basis for this space is the monomial basis , consisting of all monomials:
\[
B=\{1, X, X^2, \ldots\}.
\]
 Any set of polynomials such that there is exactly one polynomial of each degree (such as the Bernstein basis polynomials or Chebyshev polynomials) is also a basis. (Such a set of polynomials is called a polynomial sequence.)  But there are also many bases for *F*[*X*] that are not of this form.

## Properties
Many properties of finite bases result from the Steinitz exchange lemma, which states that, for any vector space , given a finite spanning set  and a linearly independent set  of  elements of , one may replace  well-chosen elements of  by the elements of  to get a spanning set containing , having its other elements in , and having the same number of elements as .

Most properties resulting from the Steinitz exchange lemma remain true when there is no finite spanning set, but their proofs in the infinite case generally require the axiom of choice or a weaker form of it, such as the ultrafilter lemma.

If  is a vector space over a field , then:
* If  is a linearly independent subset of a spanning set *S* ⊆ *V*, then there is a basis  such that
\[
L\subseteq B\subseteq S.
\]

*  has a basis (this is the preceding property with  being the empty set, and 1=*S* = *V*).
* All bases of  have the same cardinality, which is called the dimension of . This is the dimension theorem.
* A generating set  is a basis of  if and only if it is minimal, that is, no proper subset of  is also a generating set of .
* A linearly independent set  is a basis if and only if it is maximal, that is, it is not a proper subset of any linearly independent set.

If  is a vector space of dimension , then:
* A subset of  with  elements is a basis if and only if it is linearly independent.
* A subset of  with  elements is a basis if and only if it is a spanning set of .

## Coordinates
Let  be a vector space of finite dimension  over a field , and

\[
B = \{\mathbf b_1, \ldots, \mathbf b_n\}
\]

be a basis of . By definition of a basis, every **v** in  may be written, in a unique way, as

\[
\mathbf v = \lambda_1 \mathbf b_1 + \cdots + \lambda_n \mathbf b_n,
\]

where the coefficients \(\lambda_1, \ldots, \lambda_n\) are scalars (that is, elements of ), which are called the *coordinates* of **v** over . However, if one talks of the *set* of the coefficients, one loses the correspondence between coefficients and basis elements, and several vectors may have the same *set* of coefficients. For example, \(3 \mathbf b_1 + 2 \mathbf b_2\) and \(2 \mathbf b_1 + 3 \mathbf b_2\) have the same set of coefficients {{math|{2, 3}}}, and are different. It is therefore often convenient to work with an **ordered basis**; this is typically done by indexing the basis elements by the first natural numbers. Then, the coordinates of a vector form a sequence similarly indexed, and a vector is completely characterized by the sequence of coordinates. An ordered basis, especially when used in conjunction with an origin, is also called a ***coordinate frame*** or simply a *frame* (for example, a Cartesian frame or an affine frame).

Let, as usual, \(F^n\) be the set of the -tuples of elements of . This set is an -vector space, with addition and scalar multiplication defined component-wise. The map

\[
\varphi: (\lambda_1, \ldots, \lambda_n) \mapsto \lambda_1 \mathbf b_1 + \cdots + \lambda_n \mathbf b_n
\]

is a linear isomorphism from the vector space \(F^n\) onto . In other words, \(F^n\) is the coordinate space of , and the -tuple \(\varphi^{-1}(\mathbf v)\) is the coordinate vector of **v**.

The inverse image by \(\varphi\) of \(\mathbf b_i\) is the -tuple \(\mathbf e_i\) all of whose components are 0, except the th that is 1. The \(\mathbf e_i\) form an ordered basis of \(F^n\), which is called its standard basis or canonical basis. The ordered basis  is the image by \(\varphi\) of the canonical basis of \(F^n\).

It follows from what precedes that every ordered basis is the image by a linear isomorphism of the canonical basis of \(F^n\), and that every linear isomorphism from \(F^n\) onto  may be defined as the isomorphism that maps the canonical basis of \(F^n\) onto a given ordered basis of . In other words, it is equivalent to define an ordered basis of , or a linear isomorphism from \(F^n\) onto .

## Change of basis

Let *V* be a vector space of dimension  over a field *F*. Given two (ordered) bases \(B_\text{old} = (\mathbf v_1, \ldots, \mathbf v_n)\) and \(B_\text{new} = (\mathbf w_1, \ldots, \mathbf w_n)\) of *V*, it is often useful to express the coordinates of a vector  with respect to \(B_\mathrm{old}\) in terms of the coordinates with respect to \(B_\mathrm{new}.\) This can be done by the *change-of-basis formula*, that is described below. The subscripts "old" and "new" have been chosen because it is customary to refer to \(B_\mathrm{old}\) and \(B_\mathrm{new}\) as the *old basis* and the *new basis*, respectively. It is useful to describe the old coordinates in terms of the new ones, because, in general, one has expressions involving the old coordinates, and if one wants to obtain equivalent expressions in terms of the new coordinates; this is obtained by replacing the old coordinates by their expressions in terms of the new coordinates.

Typically, the new basis vectors are given by their coordinates over the old basis, that is,

\[
\mathbf w_j = \sum_{i=1}^n a_{i,j} \mathbf v_i.
\]

If \((x_1, \ldots, x_n)\) and \((y_1, \ldots, y_n)\) are the coordinates of a vector **x** over the old and the new basis respectively, the change-of-basis formula is

\[
x_i = \sum_{j=1}^n a_{i,j}y_j,
\]

for 1=*i* = 1, ..., *n*.

This formula may be concisely written in matrix notation. Let  be the matrix of the {{nowrap|\(a_{i,j}\),}} and

\[
X= \begin{bmatrix} x_1 \\ \vdots \\ x_n \end{bmatrix} \quad \text{and} \quad Y = \begin{bmatrix} y_1 \\ \vdots \\ y_n \end{bmatrix}
\]

be the column vectors of the coordinates of **v** in the old and the new basis respectively, then the formula for changing coordinates is

\[
X = A Y.
\]


The formula can be proven by considering the decomposition of the vector **x** on the two bases: one has

\[
\mathbf x = \sum_{i=1}^n x_i \mathbf v_i,
\]

and

\[
\mathbf x =\sum_{j=1}^n y_j \mathbf w_j
= \sum_{j=1}^n y_j\sum_{i=1}^n a_{i,j}\mathbf v_i
= \sum_{i=1}^n \biggl(\sum_{j=1}^n a_{i,j}y_j\biggr)\mathbf v_i.
\]


The change-of-basis formula results then from the uniqueness of the decomposition of a vector over a basis, here {{nowrap|\(B_\text{old}\);}} that is

\[
x_i = \sum_{j=1}^n a_{i,j} y_j,
\]

for 1=*i* = 1, ..., *n*.

## Related notions
### Free module

If one replaces the field occurring in the definition of a vector space by a ring, one gets the definition of a module. For modules, linear independence and spanning sets are defined exactly as for vector spaces, although "generating set" is more commonly used than that of "spanning set".

Like for vector spaces, a *basis* of a module is a linearly independent subset that is also a generating set. A major difference with the theory of vector spaces is that not every module has a basis. A module that has a basis is called a *free module*. Free modules play a fundamental role in module theory, as they may be used for describing the structure of non-free modules through free resolutions.

A module over the integers is exactly the same thing as an abelian group. Thus a free module over the integers is also a free abelian group. Free abelian groups have specific properties that are not shared by modules over other rings. Specifically, every subgroup of a free abelian group is a free abelian group, and, if  is a subgroup of a finitely generated free abelian group  (that is an abelian group that has a finite basis), then there is a basis \(\mathbf e_1, \ldots, \mathbf e_n\) of  and an integer 0 ≤ *k* ≤ *n* such that \(a_1 \mathbf e_1, \ldots, a_k \mathbf e_k\) is a basis of , for some nonzero integers \(a_1, \ldots, a_k\). For details, see .

### Analysis
In the context of infinite-dimensional vector spaces over the real or complex numbers, the term **** (named after Georg Hamel) or **algebraic basis** can be used to refer to a basis as defined in this article. This is to make a distinction with other notions of "basis" that exist when infinite-dimensional vector spaces are endowed with extra structure. The most important alternatives are orthogonal bases on Hilbert spaces, Schauder bases, and Markushevich bases on normed linear spaces. In the case of the real numbers **R** viewed as a vector space over the field **Q** of rational numbers, Hamel bases are uncountable, and have specifically the cardinality of the continuum, which is the cardinal number {{nowrap|\(2^{\aleph_0}\),}} where \(\aleph_0\) (aleph-nought) is the smallest infinite cardinal, the cardinal of the integers.

The common feature of the other notions is that they permit the taking of infinite linear combinations of the basis vectors in order to generate the space. This, of course, requires that infinite sums are meaningfully defined on these spaces, as is the case for topological vector spaces – a large class of vector spaces including e.g. Hilbert spaces, Banach spaces, or Fréchet spaces.

The preference of other types of bases for infinite-dimensional spaces is justified by the fact that the Hamel basis becomes "too big" in Banach spaces: If *X* is an infinite-dimensional normed vector space that is complete (i.e. *X* is a Banach space), then any Hamel basis of *X* is necessarily uncountable. This is a consequence of the Baire category theorem. The completeness as well as infinite dimension are crucial assumptions in the previous claim. Indeed, finite-dimensional spaces have by definition finite bases and there are infinite-dimensional (*non-complete*) normed spaces that have countable Hamel bases. Consider {{nowrap|\(c_{00}\),}} the space of the sequences \(x=(x_n)\) of real numbers that have only finitely many non-zero elements, with the norm \(\|x\|=\sup_n |x_n|\). Its standard basis, consisting of the sequences having only one non-zero element, which is equal to 1, is a countable Hamel basis.

#### Example
In the study of Fourier series, one learns that the functions {{math|1={1} ∪ { sin(*nx*), cos(*nx*) : *n* = 1, 2, 3, ... }}} are an "orthogonal basis" of the (real or complex) vector space of all (real or complex valued) functions on the interval [0, 2π] that are square-integrable on this interval, i.e., functions *f* satisfying

\[
\int_0^{2\pi} \left|f(x)\right|^2\,dx < \infty.
\]


The functions {{math|1={1} ∪ { sin(*nx*), cos(*nx*) : *n* = 1, 2, 3, ... }}} are linearly independent, and every function *f* that is square-integrable on [0, 2π] is an "infinite linear combination" of them, in the sense that

\[
\lim_{n\to\infty} \int_0^{2\pi} \biggl|a_0 + \sum_{k=1}^n \left(a_k\cos\left(kx\right)+b_k\sin\left(kx\right)\right)-f(x)\biggr|^2 dx = 0
\]


for suitable (real or complex) coefficients *a*<sub>*k*</sub>, *b*<sub>*k*</sub>. But many square-integrable functions cannot be represented as *finite* linear combinations of these basis functions, which therefore *do not* comprise a Hamel basis. Every Hamel basis of this space is much bigger than this merely countably infinite set of functions. Hamel bases of spaces of this kind are typically not useful, whereas orthonormal bases of these spaces are essential in Fourier analysis.

### Geometry
The geometric notions of an affine space, projective space, convex set, and cone have related notions of  *basis*. An **affine basis** for an *n*-dimensional affine space is \(n+1\) points in general linear position. A **** is \(n+2\) points in general position, in a projective space of dimension *n*. A **** of a polytope is the set of the vertices of its convex hull. A **** consists of one point by edge of a polygonal cone. See also a Hilbert basis (linear programming).

### Random basis
For a probability distribution in **R**<sup>*n*</sup> with a probability density function, such as the equidistribution in an *n*-dimensional ball with respect to Lebesgue measure, it can be shown that  randomly and independently chosen vectors will form a basis with probability one, which is due to the fact that  linearly dependent vectors **x**<sub>1</sub>, ..., **x**<sub>*n*</sub> in **R**<sup>*n*</sup> should satisfy the equation 1=det[**x**<sub>1</sub> ⋯ **x**<sub>*n*</sub>] = 0 (zero determinant of the matrix with columns **x**<sub>*i*</sub>), and the set of zeros of a non-trivial polynomial has zero measure. This observation has led to techniques for approximating random bases.


It is difficult to check numerically the linear dependence or exact orthogonality. Therefore, the notion of ε-orthogonality is used. For spaces with inner product, *x* is ε-orthogonal to *y* if \(\left|\left\langle x,y \right\rangle\right| / \left(\left\|x\right\|\left\|y\right\|\right) < \varepsilon\) (that is, cosine of the angle between  and  is less than ).

In high dimensions, two independent random vectors are with high probability almost orthogonal, and the number of independent random vectors, which all are with given high probability pairwise almost orthogonal, grows exponentially with dimension. More precisely, consider equidistribution in *n*-dimensional ball. Choose *N* independent random vectors from a ball (they are independent and identically distributed). Let *θ* be a small positive number. Then for
{{NumBlk||
\[
N\leq {\exp}\bigl(\tfrac14\varepsilon^2n\bigr)\sqrt{-\ln(1-\theta)}
\]
|Eq. 1}}

 random vectors are all pairwise ε-orthogonal with probability 1 − *θ*. This  growth exponentially with dimension  and \(N\gg n\) for sufficiently big . This property of random bases is a manifestation of the so-called .

The figure (right) illustrates distribution of lengths N of pairwise almost orthogonal chains of vectors that are independently randomly sampled from the *n*-dimensional cube [−1, 1]<sup>*n*</sup> as a function of dimension, *n*. A point is first randomly selected in the cube. The second point is randomly chosen in the same cube. If the angle between the vectors was within π/2 ± 0.037π/2 then the vector was retained. At the next step a new vector is generated in the same hypercube, and its angles with the previously generated vectors are evaluated. If these angles are within π/2 ± 0.037π/2 then the vector is retained. The process is repeated until the chain of almost orthogonality breaks, and the number of such pairwise almost orthogonal vectors (length of the chain) is recorded. For each *n*, 20 pairwise almost orthogonal chains were constructed numerically for each dimension. Distribution of the length of these chains is presented.

## Proof that every vector space has a basis
Let **V** be any vector space over some field **F**. Let **X** be the set of all linearly independent subsets of **V**.

The set **X** is nonempty since the empty set is an independent subset of **V**, and it is partially ordered by inclusion, which is denoted, as usual, by ⊆.

Let **Y** be a subset of **X** that is totally ordered by ⊆, and let L<sub>**Y**</sub> be the union of all the elements of **Y** (which are themselves certain subsets of **V**).

Since (**Y**, ⊆) is totally ordered, every finite subset of L<sub>**Y**</sub> is a subset of an element of **Y**, which is a linearly independent subset of **V**, and hence L<sub>**Y**</sub> is linearly independent. Thus L<sub>**Y**</sub> is an element of **X**. Therefore, L<sub>**Y**</sub> is an upper bound for **Y** in (**X**, ⊆): it is an element of **X**, that contains every element of **Y**.

As **X** is nonempty, and every totally ordered subset of (**X**, ⊆) has an upper bound in **X**, Zorn's lemma asserts that **X** has a maximal element. In other words, there exists some element L<sub>**max**</sub> of **X** satisfying the condition that whenever L<sub>**max**</sub> ⊆ L for some element L of **X**, then 1=L = L<sub>**max**</sub>.

It remains to prove that L<sub>**max**</sub> is a basis of **V**. Since L<sub>**max**</sub> belongs to **X**, we already know that L<sub>**max**</sub> is a linearly independent subset of **V**.

If there were some vector **w** of **V** that is not in the span of L<sub>**max**</sub>, then **w** would not be an element of L<sub>**max**</sub> either. Let {{math|1=L<sub>**w**</sub> = L<sub>**max**</sub> ∪ {**w**}}}. This set is an element of **X**, that is, it is a linearly independent subset of **V** (because **w** is not in the span of L<sub>**max**</sub>, and L<sub>**max**</sub> is independent). As L<sub>**max**</sub> ⊆ L<sub>**w**</sub>, and L<sub>**max**</sub> ≠ L<sub>**w**</sub> (because L<sub>**w**</sub> contains the vector **w** that is not contained in L<sub>**max**</sub>), this contradicts the maximality of L<sub>**max**</sub>. Thus this shows that L<sub>**max**</sub> spans **V**.

Hence L<sub>**max**</sub> is linearly independent and spans **V**. It is thus a basis of **V**, and this proves that every vector space has a basis.

This proof relies on Zorn's lemma, which is equivalent to the axiom of choice. Conversely, it has been proved that if every vector space has a basis, then the axiom of choice is true. Thus the two assertions are equivalent.

## See also
*Basis of a matroid
*Basis of a linear program
*Coordinate system
*
*
*

## Notes


## References
### General references
*
*

### Historical references
*
*
*
*
*
* , reprint:
*
*
*
*
*

## External links
* Instructional videos from Khan Academy
**Introduction to bases of subspaces
**Proof that any subspace basis has same number of elements
*
*

