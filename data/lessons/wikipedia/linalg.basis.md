> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Basis_%28linear_algebra%29) — CC BY-SA 4.0

# Basis and dimension

Basis

In mathematics, a set of elements of a vector space *V* is called a **basis** (: **bases**) if every element of *V* can be written in a unique way as a finite linear combination of elements of. The coefficients of this linear combination are referred to as **components** or **coordinates** of the vector with respect to. The elements of a basis are called ****.

Equivalently, a set is a basis if its elements are linearly independent and every element of is a linear combination of elements of. In other words, a basis is a linearly independent spanning set.

A vector space can have several bases; however all the bases have the same number of elements, called the dimension of the vector space.

This article deals mainly with finite-dimensional vector spaces. However, many of the principles are also valid for infinite-dimensional vector spaces.

Basis vectors find applications in the study of crystal structures and frames of reference.

## Definition
A **basis** *B* of a vector space *V* over a field *F* (such as the real numbers **R** or the complex numbers **C**) is a linearly independent subset of *V* that spans *V*. This means that a subset of *V* is a basis if it satisfies the two following conditions:
* linear independence: for every finite subset \(\{\mathbf v_1, \dotsc, \mathbf v_m\}\) of , if \(c_1 \mathbf v_1 + \cdots + c_m \mathbf v_m = \mathbf 0\) for some \(c_1,\dotsc,c_m\) in , then \(c_1 = \cdots = c_m = 0\); and
* the spanning property: for every vector \(\mathbf{v}\) in \(V\), one can choose \(a_1,\dotsc,a_n\) in \(F\) and \(\mathbf v_1, \dotsc, \mathbf v_n\) in \(B\) such that \(\mathbf v = a_1 \mathbf v_1 + \cdots + a_n \mathbf v_n\). In other words, \(\mathbb{v}\) can be represented as a linear combination of some vectors in \(B\).

The first property may be equivalently phrased as follows: if a linear combination of vectors in is equal to the zero vector, then all the scalars coefficients of the combination must be zero.

If is a basis for , then every vector \(\mathbf{v}\) in can be written as a linear combination of vectors in (by the spanning property); it follows from linear independence that this can be done in exactly one way. Thus, the scalar coefficients \(a_i\) that appear in this combination are uniquely determined; they are called the *coordinates* of \(\mathbf{v}\) with respect to the basis \(B\).

A vector space that has a finite basis is called finite-dimensional. In this case, the finite subset can be taken as \(B\) itself to check for linear independence in the above definition.

It is often convenient or even necessary to have an ordering on the basis vectors, for example, when discussing orientation, or when one considers the scalar coefficients of a vector with respect to a basis without referring explicitly to the basis elements. In this case, the ordering is necessary for associating each coefficient with the corresponding basis element. This ordering can be done by numbering the basis elements. In order to emphasize that an order has been chosen, one speaks of an **ordered basis**, which is therefore not simply an unstructured set, but a sequence, an indexed family, or similar; see below.

## Examples

The set **R**\(^{2}\) of the ordered pairs of real numbers is a vector space under the operations of component-wise addition

\[
(a, b) + (c, d) = (a + c, b+d)
\]

and scalar multiplication

\[
\lambda (a,b) = (\lambda a, \lambda b),
\]

where \(\lambda\) is any real number. A simple basis of this vector space consists of the two vectors **e**\(_{1}\) = (1, 0) and **e**\(_{2}\) = (0, 1). These vectors form a basis (called the standard basis) because any vector **v** = (*a*, *b*) of **R**\(^{2}\) may be uniquely written as
\[
\mathbf v = a \mathbf e_1 + b \mathbf e_2.
\]
 Any other pair of linearly independent vectors of **R**\(^{2}\), such as (1, 1) and (−1, 2), forms also a basis of **R**\(^{2}\).

More generally, if is a field, the set \(F^n\) of -tuples of elements of is a vector space for similarly defined addition and scalar multiplication. Let
\[
\mathbf e_i = (0, \ldots, 0,1,0,\ldots, 0)
\]
 be the -tuple with all components equal to 0, except the th, which is 1. Then \(\mathbf e_1, \ldots, \mathbf e_n\) is a basis of \(F^n,\) which is called the *standard basis* of \(F^n.\)

A different flavor of example is given by polynomial rings. If is a field, the collection *F*[*X*] of all polynomials in one indeterminate with coefficients in is an -vector space. One basis for this space is the monomial basis , consisting of all monomials:
\[
B=\{1, X, X^2, \ldots\}.
\]
 Any set of polynomials such that there is exactly one polynomial of each degree (such as the Bernstein basis polynomials or Chebyshev polynomials) is also a basis. (Such a set of polynomials is called a polynomial sequence.) But there are also many bases for *F*[*X*] that are not of this form.

## Properties
Many properties of finite bases result from the Steinitz exchange lemma, which states that, for any vector space , given a finite spanning set and a linearly independent set of elements of , one may replace well-chosen elements of by the elements of to get a spanning set containing , having its other elements in , and having the same number of elements as. Most properties resulting from the Steinitz exchange lemma remain true when there is no finite spanning set, but their proofs in the infinite case generally require the axiom of choice or a weaker form of it, such as the ultrafilter lemma.

If is a vector space over a field , then:
* If is a linearly independent subset of a spanning set *S* ⊆ *V*, then there is a basis such that
\[
L\subseteq B\subseteq S.
\]

* has a basis (this is the preceding property with being the empty set, and *S* = *V*).
* All bases of have the same cardinality, which is called the dimension of. This is the dimension theorem.
* A generating set is a basis of if and only if it is minimal, that is, no proper subset of is also a generating set of. * A linearly independent set is a basis if and only if it is maximal, that is, it is not a proper subset of any linearly independent set.

If is a vector space of dimension , then:
* A subset of with elements is a basis if and only if it is linearly independent.
* A subset of with elements is a basis if and only if it is a spanning set of. ## Coordinates
Let be a vector space of finite dimension over a field , and

\[
B = \{\mathbf b_1, \ldots, \mathbf b_n\}
\]

be a basis of. By definition of a basis, every **v** in may be written, in a unique way, as

\[
\mathbf v = \lambda_1 \mathbf b_1 + \cdots + \lambda_n \mathbf b_n,
\]

where the coefficients \(\lambda_1, \ldots, \lambda_n\) are scalars (that is, elements of ), which are called the *coordinates* of **v** over. However, if one talks of the *set* of the coefficients, one loses the correspondence between coefficients and basis elements, and several vectors may have the same *set* of coefficients. For example, \(3 \mathbf b_1 + 2 \mathbf b_2\) and \(2 \mathbf b_1 + 3 \mathbf b_2\) have the same set of coefficients.
