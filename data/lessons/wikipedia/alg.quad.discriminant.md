> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Discriminant) — CC BY-SA 4.0

# Discriminant: number of real solutions

In mathematics, the **discriminant** of a polynomial is a quantity that depends on the coefficients and allows deducing some properties of the roots without computing them. More precisely, it is a polynomial function of the coefficients of the original polynomial. The discriminant is widely used in polynomial factoring, number theory, and algebraic geometry.

The discriminant of the quadratic polynomial \(ax^2+bx+c\) is
\(b^2-4ac,\)
the quantity which appears under the square root in the quadratic formula. If \(a\ne 0,\) this discriminant is zero if and only if the polynomial has a double root. In the case of real coefficients, it is positive if the polynomial has two distinct real roots, and negative if it has two distinct complex conjugate roots. Similarly, the discriminant of a cubic polynomial is zero if and only if the polynomial has a multiple root. In the case of a cubic with real coefficients, the discriminant is positive if the polynomial has three distinct real roots, and negative if it has one real root and two distinct complex conjugate roots.

More generally, the discriminant of a univariate polynomial of positive degree is zero if and only if the polynomial has a multiple root. For real coefficients and no multiple roots, the discriminant is positive if the number of non-real roots is a multiple of 4 (including none), and negative otherwise.

Several generalizations are also called discriminant: the *discriminant of an algebraic number field*; the *discriminant of a quadratic form*; and more generally, the *discriminant* of a form, of a homogeneous polynomial, or of a projective hypersurface (these three concepts are essentially equivalent).

## Origin
The term "discriminant" was coined in 1851 by the British mathematician James Joseph Sylvester.

## Definition
Let
\(A(x) = a_nx^n+a_{n-1}x^{n-1}+\cdots+a_1x+a_0\)
be a polynomial of degree *n* (this means \(a_n\ne 0\)), such that the coefficients \(a_0, \ldots, a_n\) belong to a field, or, more generally, to a commutative ring. The resultant of *A* with its derivative,
\(A'(x) = na_nx^{n-1}+(n-1)a_{n-1}x^{n-2}+\cdots+a_1,\)
is a polynomial in \(a_0, \ldots, a_n\) with integer coefficients, which is the determinant of the Sylvester matrix of *A* and *A*. The nonzero entries of the first column of the Sylvester matrix are \(a_n\) and \(na_n,\) and the resultant is thus a multiple of \(a_n.\) Hence the discriminant—up to its sign—is defined as the quotient of the resultant of *A* and ''A' by \(a_n\):

\(\operatorname{Disc}_x(A) = \frac{(-1)^{n(n-1)/2{a_n} \operatorname{Res}_x(A,A')\)
where the resultant is computed with \(A'\) considered of degree \(n-1,\) even if it has a lower degree, which occurs when the characteristic divides \(n\).

Historically, this sign has been chosen such that, over the reals, the discriminant will be positive when all the roots of the polynomial are real. The division by \(a_n\) may not be well defined if the ring of the coefficients contains zero divisors. Such a problem may be avoided by replacing \(a_n\) by 1 in the first column of the Sylvester matrix—*before* computing the determinant. In any case, the discriminant is a polynomial in \(a_0, \ldots, a_n\) with integer coefficients.

### Expression in terms of the roots
When the above polynomial is defined over a field, it has *n* roots, \(r_1, r_2, \dots, r_n\), not necessarily all distinct, in any algebraically closed extension of the field. (If the coefficients are real numbers, the roots may be taken in the field of complex numbers, where the fundamental theorem of algebra applies.)

In terms of the roots, the discriminant is equal to

\(\operatorname{Disc}_x(A) = a_n^{2n-2}\prod_{i < j} (r_i-r_j)^2
= (-1)^{n(n-1)/2} a_n^{2n-2} \prod_{i \neq j} (r_i-r_j).\)

It is thus the square of the Vandermonde polynomial times \(a_n^{2n-2}\).

This expression for the discriminant is often taken as a definition. It makes clear that if the polynomial has a multiple root, then its discriminant is zero, and that, in the case of real coefficients, if all the roots are real and simple, then the discriminant is positive. Unlike the previous definition, this expression is not obviously a polynomial in the coefficients, but this follows either from the fundamental theorem of Galois theory, or from the fundamental theorem of symmetric polynomials and Vieta's formulas by noting that this expression is a symmetric polynomial in the roots of *A*.

## Low degrees
The discriminant of a linear polynomial (degree 1) is rarely considered. If needed, it is commonly defined to be equal to 1 (using the usual conventions for the empty product and considering that one of the two blocks of the Sylvester matrix is empty). There is no common convention for the discriminant of a constant polynomial (i.e., polynomial of degree 0).

For small degrees, the discriminant is rather simple (see below), but for higher degrees, it may become unwieldy. For example, the discriminant of a general quartic has 16 terms, that of a quintic has 59 terms, and that of a sextic has 246 terms.
This is OEIS sequence. ### Degree 2

The quadratic polynomial \(ax^2+bx+c \,\) has discriminant
\(b^2-4ac\,.\)

The square root of the discriminant appears in the quadratic formula for the roots of the quadratic polynomial:
\(x_{1,2}=\frac{-b \pm \sqrt {b^2-4ac{2a}.\)

where the discriminant is zero if and only if the two roots are equal. If *a*, *b*, *c* are real numbers, the polynomial has two distinct real roots if the discriminant is positive, and two complex conjugate roots if it is negative.

The discriminant is the product of *a* and the square of the difference of the roots.

If *a*, *b*, *c* are rational numbers, then the discriminant is the square of a rational number if and only if the two roots are rational numbers.

### Degree 3

The cubic polynomial \(ax^3+bx^2+cx+d \,\) has discriminant
\(b^2c^2-4ac^3-4b^3d-27a^2d^2+18abcd\,.\)

In the special case of a depressed cubic polynomial \(x^3+px+q\), the discriminant simplifies to
\(-4p^3-27q^2\,.\)

The discriminant is zero if and only if at least two roots are equal. If the coefficients are real numbers, and the discriminant is not zero, the discriminant is positive if the roots are three distinct real numbers, and negative if there is one real root and two complex conjugate roots.

The square root of a quantity strongly related to the discriminant appears in the formulas for the roots of a cubic polynomial. Specifically, this quantity can be −3 times the discriminant, or its product with the square of a rational number; for example, the square of 1/18 in the case of Cardano formula.

If the polynomial is irreducible and its coefficients are rational numbers (or belong to a number field), then the discriminant is a square of a rational number (or a number from the number field) if and only if the Galois group of the cubic equation is the cyclic group of order three.

### Degree 4

The quartic polynomial
\(ax^4+bx^3+cx^2+dx+e\,\)
has discriminant
\(\begin{align}
& 256a^3e^3-192a^2bde^2-128a^2c^2e^2+144a^2cd^2e \\[4pt]
&\quad -27a^2d^4+144ab^2ce^2-6ab^2d^2e-80abc^2de \\[4pt]
&\quad +18abcd^3+16ac^4e-4ac^3d^2-27b^4e^2+18b^3cde \\[4pt]
&\quad -4b^3d^3-4b^2c^3e+b^2c^2d^2.
\end{align}\)

The depressed quartic polynomial
\(x^4+cx^2+dx+e\,\)
has discriminant
\(\begin{align}
{} & 16c^4e -4c^3d^2 -128c^2e^2+144cd^2e -27d^4 + 256e^3.
\end{align}\)

The discriminant is zero if and only if at least two roots are equal. If the coefficients are real numbers and the discriminant is negative, then there are two real roots and two complex conjugate roots. Conversely, if the discriminant is positive, then the roots are either all real or all non-real.

## Properties
### Zero discriminant
The discriminant of a polynomial over a field is zero if and only if the polynomial has a multiple root in some field extension.

The discriminant of a polynomial over an integral domain is zero if and only if the polynomial and its derivative have a non-constant common divisor.

In characteristic 0, this is equivalent to saying that the polynomial is not square-free (i.e., it is divisible by the square of a non-constant polynomial).

In nonzero characteristic *p*, the discriminant is zero if and only if the polynomial is not square-free or it has an irreducible factor which is not separable (i.e., the irreducible factor is a polynomial in \(x^p\)).

### Invariance under change of the variable
The discriminant of a polynomial is, up to a scaling, invariant under any projective transformation of the variable. As a projective transformation may be decomposed into a product of translations, homotheties and inversions, this results in the following formulas for simpler transformations, where *P*(*x*) denotes a polynomial of degree *n*, with \(a_n\) as leading coefficient.

* *Invariance by translation*:
\(\operatorname{Disc}_x(P(x+\alpha)) = \operatorname{Disc}_x(P(x))\)
This results from the expression of the discriminant in terms of the roots
* *Invariance by homothety*:
\(\begin{aligned}
\operatorname{Disc}_x(P(\alpha x)) &= \alpha^{n(n-1)}\operatorname{Disc}_x(P(x)) \\
\operatorname{Disc}_x(\alpha P(x)) &= \alpha^{2n-2}\operatorname{Disc}_x(P(x))
\end{aligned}\)
This results from the expression in terms of the roots, or of the quasi-homogeneity of the discriminant.
* *Invariance by inversion*: The reciprocal polynomial of \(P(x) = a_nx^n + \cdots + a_0\) with \(a_0 \neq 0\) is:
\(P^{\mathrm{r(x) = x^nP(1/x) = a_0x^n +\cdots +a_n.\) Then:
\(\operatorname{Disc}_x(P^{\mathrm{r\!\!\,(x)) = \operatorname{Disc}_x(P(x)).\)

### Invariance under ring homomorphisms
Let \(\varphi\colon R \to S\) be a homomorphism of commutative rings. Given a polynomial
\(A = a_nx^n+a_{n-1}x^{n-1}+\cdots+a_0\)
in *R*[*x*], the homomorphism \(\varphi\) acts on *A* for producing the polynomial
\(A^\varphi = \varphi(a_n)x^n+\varphi(a_{n-1})x^{n-1}+ \cdots+\varphi(a_0)\)
in *S*[*x*].

The discriminant is invariant under \(\varphi\) in the following sense. If \(\varphi(a_n)\ne 0,\) then
\(\operatorname{Disc}_x(A^\varphi) = \varphi(\operatorname{Disc}_x(A)).\)
As the discriminant is defined in terms of a determinant, this property results immediately from the similar property of determinants.

If \(\varphi(a_n)= 0,\) then \(\varphi(\operatorname{Disc}_x(A))\) may be zero or not. One has, when \(\varphi(a_n)= 0,\)
\(\varphi(\operatorname{Disc}_x(A)) = \varphi(a_{n-1})^2\operatorname{Disc}_x(A^\varphi).\)

When one is only interested in knowing whether a discriminant is zero (as is generally the case in algebraic geometry), these properties may be summarised as:
\(\varphi(\operatorname{Disc}_x(A)) = 0\) if and only if either \(\operatorname{Disc}_x(A^\varphi)=0\) or \(\deg(A)-\deg(A^\varphi)\ge 2.\)

This is often interpreted as saying that \(\varphi(\operatorname{Disc}_x(A)) = 0\) if and only if \(A^\varphi\) has a multiple root (possibly at infinity).

### Product of polynomials
If *R* = *PQ* is a product of polynomials in *x*, then
\(\begin{align}
\operatorname{disc}_x(R) &= \operatorname{disc}_x(P)\operatorname{Res}_x(P,Q)^2\operatorname{disc}_x(Q)
\\[5pt]
{}&=(-1)^{pq}\operatorname{disc}_x(P)\operatorname{Res}_x(P,Q)\operatorname{Res}_x(Q,P)\operatorname{disc}_x(Q),
\end{align}\)
where \(\operatorname{Res}_x\) denotes the resultant with respect to the variable *x*, and *p* and *q* are the respective degrees of *P* and *Q*.

This property follows immediately by substituting the expression for the resultant, and the discriminant, in terms of the roots of the respective polynomials.

### Homogeneity
The discriminant is a homogeneous polynomial in the coefficients; it is also a homogeneous polynomial in the roots and thus quasi-homogeneous in the coefficients.

The discriminant of a polynomial of degree *n* is homogeneous of degree 2*n* − 2 in the coefficients. This can be seen in two ways. In terms of the roots-and-leading-term formula, multiplying all the coefficients by does not change the roots, but multiplies the leading term by. In terms of its expression as a determinant of a (2*n* − 1) × (2*n* − 1) matrix (the Sylvester matrix) divided by , the determinant is homogeneous of degree 2*n* − 1 in the entries, and dividing by makes the degree 2*n* − 2.

The discriminant of a polynomial of degree *n* is homogeneous of degree *n*(*n* − 1) in the roots. This follows from the expression of the discriminant in terms of the roots, which is the product of a constant and \(\binom{n}{2} = \frac{n(n-1)}{2}\) squared differences of roots.

The discriminant of a polynomial of degree *n* is quasi-homogeneous of degree *n*(*n* − 1) in the coefficients, if, for every *i*, the coefficient of \(x^i\) is given the weight *n* − *i*. It is also quasi-homogeneous of the same degree, if, for every *i*, the coefficient of \(x^i\) is given the weight *i*. This is a consequence of the general fact that every polynomial which is homogeneous and symmetric in the roots may be expressed as a quasi-homogeneous polynomial in the elementary symmetric functions of the roots.

Consider the polynomial
\(P=a_nx^n+a_{n-1}x^{n-1}+ \cdots +a_0.\)
It follows from what precedes that the exponents in every monomial \(a_0^{i_0}, \dots , a_n^{i_n}\) appearing in the discriminant satisfy the two equations
\(i_0+i_1+\cdots+i_n=2n-2\)
and
\(i_1+2i_2 + \cdots+n i_n=n(n-1),\)
and also the equation
\(ni_0 +(n-1)i_1+ \cdots+ i_{n-1}=n(n-1),\)
which is obtained by subtracting the second equation from the first one multiplied by *n*.

This restricts the possible terms in the discriminant. For the general quadratic polynomial, the discriminant \(b^2-4ac\) is a homogeneous polynomial of degree 2 which has only two terms, while the general homogeneous polynomial of degree two in three variables has 6 terms. The discriminant of the general cubic polynomial is a homogeneous polynomial of degree 4 in four variables; it has five terms, which is the maximum allowed by the above rules, while the general homogeneous polynomial of degree 4 in 4 variables has 35 terms.

For higher degrees, there may be monomials which satisfy above rules and do not appear in the discriminant. The first example is for the quartic polynomial \(ax^4 + bx^3 + cx^2 + dx + e\), in which case the monomial \(bc^4d\) satisfies the rules without appearing in the discriminant.

## Real roots
In this section, all polynomials have real coefficients.

It has been seen in that the sign of the discriminant provides useful information on the nature of the roots for polynomials of degree 2 and 3. For higher degrees, the information provided by the discriminant is less complete, but still useful. More precisely, for a polynomial of degree *n*, one has:
* The polynomial has a multiple root if and only if its discriminant is zero.
* If the discriminant is positive, the number of non-real roots is a multiple of 4. That is, there is a nonnegative integer *k* ≤ *n*/4 such that there are 2*k* pairs of complex conjugate roots and *n* − 4*k* real roots.
* If the discriminant is negative, the number of non-real roots is not a multiple of 4. That is, there is a nonnegative integer *k* ≤ (*n* − 2)/4 such that there are 2*k* + 1 pairs of complex conjugate roots and *n* − 4*k* + 2 real roots.

## Homogeneous bivariate polynomial
Let
\(A(x,y) = a_0x^n+ a_1 x^{n-1}y + \cdots + a_n y^n=\sum_{i=0}^n a_i x^{n-i}y^i\)
be a homogeneous polynomial of degree *n* in two indeterminates.

Supposing, for the moment, that \(a_0\) and \(a_n\) are both nonzero, one has
\(\operatorname{Disc}_x(A(x,1))=\operatorname{Disc}_y(A(1,y)).\)
Denoting this quantity by \(\operatorname{Disc}^h (A),\)
one has
\(\operatorname{Disc}_x (A) =y^{n(n-1)} \operatorname{Disc}^h (A),\)
and
\(\operatorname{Disc}_y (A) =x^{n(n-1)} \operatorname{Disc}^h (A).\)

Because of these properties, the quantity \(\operatorname{Disc}^h (A)\) is called the *discriminant* or the *homogeneous discriminant* of *A*.

If \(a_0\) and \(a_n\) are permitted to be zero, the polynomials *A*(*x*, 1) and *A*(1, *y*) may have a degree smaller than *n*. In this case, above formulas and definition remain valid, if the discriminants are computed as if all polynomials would have the degree. This means that the discriminants must be computed with \(a_0\) and \(a_n\) indeterminate, the substitution for them of their actual values being done *after* this computation. Equivalently, the formulas of must be used.

## Use in algebraic geometry
The typical use of discriminants in algebraic geometry is for studying plane algebraic curves, and more generally algebraic hypersurfaces. Let *V* be such a curve or hypersurface; *V* is defined as the zero set of a multivariate polynomial. This polynomial may be considered as a univariate polynomial in one of the indeterminates, with polynomials in the other indeterminates as coefficients. The discriminant with respect to the selected indeterminate defines a hypersurface *W* in the space of the other indeterminates. The points of *W* are exactly the projection of the points of *V* (including the points at infinity), which either are singular or have a tangent hyperplane that is parallel to the axis of the selected indeterminate.

For example, let be a bivariate polynomial in and with real coefficients, so that *f*  = 0 is the implicit equation of a real plane algebraic curve. Viewing as a univariate polynomial in with coefficients depending on , then the discriminant is a polynomial in whose roots are the -coordinates of the singular points, of the points with a tangent parallel to the -axis and of some of the asymptotes parallel to the -axis. In other words, the computation of the roots of the -discriminant and the -discriminant allows one to compute all of the remarkable points of the curve, except the inflection points.

## Generalizations
There are two classes of the concept of discriminant. The first class is the discriminant of an algebraic number field, which, in some cases including quadratic fields, is the discriminant of a polynomial defining the field.

Discriminants of the second class arise for problems depending on coefficients, when degenerate instances or singularities of the problem are characterized by the vanishing of a single polynomial in the coefficients. This is the case for the discriminant of a polynomial, which is zero when two roots collapse. Most of the cases, where such a generalized discriminant is defined, are instances of the following.

Let *A* be a homogeneous polynomial in *n* indeterminates over a field of characteristic 0, or of a prime characteristic that does not divide the degree of the polynomial. The polynomial *A* defines a projective hypersurface, which has singular points if and only the *n* partial derivatives of *A* have a nontrivial common zero. This is the case if and only if the multivariate resultant of these partial derivatives is zero, and this resultant may be considered as the discriminant of *A*. However, because of the integer coefficients resulting of the derivation, this multivariate resultant may be divisible by a power of *n*, and it is better to take, as a discriminant, the primitive part of the resultant, computed with generic coefficients. The restriction on the characteristic is needed because otherwise a common zero of the partial derivative is not necessarily a zero of the polynomial (see Euler's identity for homogeneous polynomials).

In the case of a homogeneous bivariate polynomial of degree *d*, this general discriminant is \(d^{d-2}\) times the discriminant defined in. Several other classical types of discriminants, that are instances of the general definition are described in next sections.

### Quadratic forms

A quadratic form is a function over a vector space, which is defined over some basis by a homogeneous polynomial of degree 2:

\(Q(x_1,\ldots,x_n) \ =\ \sum_{i=1}^n a_{ii} x_i^2+\sum_{1\le i ×)^{2}\), the quotient of the multiplicative monoid of *K* by the subgroup of the nonzero squares (that is, two elements of *K* are in the same equivalence class if one is the product of the other by a nonzero square). It follows that over the complex numbers, a discriminant is equivalent to 0 or 1. Over the real numbers, a discriminant is equivalent to −1, 0, or 1. Over the rational numbers, a discriminant is equivalent to a unique square-free integer.

By a theorem of Jacobi, a quadratic form over a field of characteristic different from 2 can be expressed, after a linear change of variables, in **diagonal form** as
\(a_1x_1^2 + \cdots + a_nx_n^2.\)
More precisely, a quadratic form may be expressed as a sum
\(\sum_{i=1}^n a_i L_i^2\)
where the *L*_{*i*} are independent linear forms and is the number of the variables (some of the *a*_{*i*} may be zero). Equivalently, for any symmetric matrix *A*, there is an elementary matrix *S* such that \(S^\mathrm T A\,S\) is a diagonal matrix.
Then the discriminant is the product of the *a*_{*i*}, which is well-defined as a class in *K*/(*K*^{×})^{2}.

Geometrically, the discriminant of a quadratic form in three variables is the equation of a quadratic projective curve. The discriminant is zero if and only if the curve is decomposed in lines (possibly over an algebraically closed extension of the field).

A quadratic form in four variables is the equation of a projective surface. The surface has a singular point if and only its discriminant is zero. In this case, either the surface may be decomposed in planes, or it has a unique singular point, and is a cone or a cylinder. Over the reals, if the discriminant is positive, then the surface either has no real point or has everywhere a negative Gaussian curvature. If the discriminant is negative, the surface has real points, and has a negative Gaussian curvature.

### Conic sections
A conic section is a plane curve defined by an implicit equation of the form
\(ax^2+ 2bxy + cy^2 + 2dx + 2ey + f = 0,\)
where *a*, *b*, *c*, *d*, *e*, *f* are real numbers.

Two quadratic forms, and thus two discriminants may be associated to a conic section.

The first quadratic form is
\(ax^2+ 2bxy + cy^2 + 2dxz + 2eyz + fz^2 = 0.\)
Its discriminant is the determinant
\(\begin{vmatrix} a & b & d\\b & c & e\\d & e & f \end{vmatrix}.\)
It is zero if the conic section degenerates into two lines, a double line or a single point.

The second discriminant, which is the only one that is considered in many elementary textbooks, is the discriminant of the homogeneous part of degree two of the equation. It is equal to

\(b^2 - ac,\)

and determines the shape of the conic section. If this discriminant is negative, the curve either has no real points, or is an ellipse or a circle, or, if degenerated, is reduced to a single point. If the discriminant is zero, the curve is a parabola, or, if degenerated, a double line or two parallel lines. If the discriminant is positive, the curve is a hyperbola, or, if degenerated, a pair of intersecting lines.

### Real quadric surfaces
A real quadric surface in the Euclidean space of dimension three is a surface that may be defined as the zeros of a polynomial of degree two in three variables. As for the conic sections there are two discriminants that may be naturally defined. Both are useful for getting information on the nature of a quadric surface.

Let \(P(x,y,z)\) be a polynomial of degree two in three variables that defines a real quadric surface. The first associated quadratic form, \(Q_4,\) depends on four variables, and is obtained by homogenizing *P*; that is
\(Q_4(x,y,z,t)=t^2P(x/t,y/t, z/t).\)
Let us denote its discriminant by \(\Delta_4.\)

The second quadratic form, \(Q_3,\) depends on three variables, and consists of the terms of degree two of *P*; that is
\(Q_3(x,y,z)=Q_4(x, y,z,0).\)
Let us denote its discriminant by \(\Delta_3.\)

If \(\Delta_4>0,\) and the surface has real points, it is either a hyperbolic paraboloid or a one-sheet hyperboloid. In both cases, this is a ruled surface that has a negative Gaussian curvature at every point.

If \(\Delta_4<0,\) the surface is either an ellipsoid or a two-sheet hyperboloid or an elliptic paraboloid. In all cases, it has a positive Gaussian curvature at every point.

If \(\Delta_4=0,\) the surface has a singular point, possibly at infinity. If there is only one singular point, the surface is a cylinder or a cone. If there are several singular points the surface consists of two planes, a double plane or a single line.

When \(\Delta_4\ne 0,\) the sign of \(\Delta_3,\) if not 0, does not provide any useful information, as changing *P* into −*P* does not change the surface, but changes the sign of \(\Delta_3.\) However, if \(\Delta_4\ne 0\) and \(\Delta_3 = 0,\) the surface is a paraboloid, which is elliptic or hyperbolic, depending on the sign of \(\Delta_4.\)

### Discriminant of an algebraic number field
The discriminant of an algebraic number field measures the size of the (ring of integers of the) algebraic number field.

More specifically, it is proportional to the squared volume of the fundamental domain of the ring of integers, and it regulates which primes are ramified.

The discriminant is one of the most basic invariants of a number field, and occurs in several important analytic formulas such as the functional equation of the Dedekind zeta function of *K*, and the analytic class number formula for *K*. A theorem of Hermite states that there are only finitely many number fields of bounded discriminant, however determining this quantity is still an open problem, and the subject of current research.

Let *K* be an algebraic number field, and let *O_{K}* be its ring of integers. Let *b*_{1},... , *b_{n}* be an integral basis of *O_{K}* (i.e. a basis as a **Z**-module), and let {σ_{1},... , σ_{*n*be the set of embeddings of *K* into the complex numbers (i.e. injective ring homomorphisms *K* → **C**). The **discriminant** of *K* is the square of the determinant of the *n* by *n* matrix *B* whose (*i*,*j*)-entry is σ_{*i*}(*b_{j}*). Symbolically,

\(\Delta_K=\det\left(\begin{array}{cccc}
\sigma_1(b_1) & \sigma_1(b_2) &\cdots & \sigma_1(b_n) \\
\sigma_2(b_1) & \ddots & & \vdots \\
\vdots & & \ddots & \vdots \\
\sigma_n(b_1) & \cdots & \cdots & \sigma_n(b_n)
\end{array}\right)^2.\)

The discriminant of *K* can be referred to as the absolute discriminant of *K* to distinguish it from an extension *K*/*L* of number fields. The latter is an ideal in the ring of integers of *L*, and like the absolute discriminant it indicates which primes are ramified in *K*/*L*. It is a generalization of the absolute discriminant allowing for *L* to be bigger than **Q**; in fact, when *L* = **Q**, the relative discriminant of *K*/**Q** is the principal ideal of **Z** generated by the absolute discriminant of *K*.
### Fundamental discriminants
A specific type of discriminant useful in the study of quadratic fields is the fundamental discriminant. It arises in the theory of integral binary quadratic forms, which are expressions of the form:
\[
Q(x, y) = ax^2 + bxy + cy^2
\]
where \(a\), \(b\), and \(c\) are integers. The discriminant of \(Q(x, y)\) is given by:
\[
D = b^2 - 4ac.
\]
Not every integer can arise as a discriminant of an integral binary quadratic form. An integer \(D\) is a fundamental discriminant if and only if it meets one of the following criteria:
* Case 1: \(D\) is congruent to 1 modulo 4 (\(D \equiv 1 \pmod{4}\)) and is square-free, meaning it is not divisible by the square of any prime number.
* Case 2: \(D\) is equal to four times an integer \(m\) (\(D = 4m\)) where \(m\) is congruent to 2 or 3 modulo 4 (\(m \equiv 2, 3 \pmod{4}\)) and is square-free.

These conditions ensure that every fundamental discriminant corresponds uniquely to a specific type of quadratic form.

The first eleven positive fundamental discriminants are:

1, 5, 8, 12, 13, 17, 21, 24, 28, 29, 33

The first eleven negative fundamental discriminants are:

−3, −4, −7, −8, −11, −15, −19, −20, −23, −24, −31. #### Quadratic number fields
A quadratic field is a field extension of the rational numbers \(\mathbb{Q}\) that has degree 2. The discriminant of a quadratic field plays a role analogous to the discriminant of a quadratic form.

There exists a fundamental connection: an integer \(D_0\) is a fundamental discriminant if and only if:

* \(D_0 = 1\), or
* \(D_0\) is the discriminant of a quadratic field.

For each fundamental discriminant \(D_0 \neq 1\), there exists a unique (up to isomorphism) quadratic field with \(D_0\) as its discriminant. This connects the theory of quadratic forms and the study of quadratic fields.

#### Prime factorization
Fundamental discriminants can also be characterized by their prime factorization. Consider the set \(S\) consisting of \(-8, 8, -4,\) the prime numbers congruent to 1 modulo 4, and the additive inverses of the prime numbers congruent to 3 modulo 4:
\[
S = \{-8, -4, 8, -3, 5, -7, -11, 13, 17, -19,... \}
\]
An integer \(D \neq 1\) is a fundamental discriminant if and only if it is a product of elements of \(S\) that are pairwise coprime.
