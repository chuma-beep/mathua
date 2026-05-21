> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Ring_%28mathematics%29) — CC BY-SA 4.0

# Rings and fields

In mathematics, a **ring** is an algebraic structure consisting of a set with two binary operations typically called *addition* and *multiplication* and denoted like addition and multiplication of integers. They work similarly to integer addition and multiplication, except that multiplication in a ring does not need to be commutative. Ring elements may be numbers such as integers or complex numbers, but they may also be non-numerical objects such as polynomials, square matrices, functions, and power series.

More formally, a ring is a set that is endowed with two binary operations (*addition* and *multiplication*) such that the ring is an abelian group with respect to addition. The multiplication is associative, is distributive over the addition operation, and has a multiplicative identity element. Some authors apply the term *ring* to a further generalization, often called a *rng*, that omits the requirement for a multiplicative identity, and instead call the structure defined above a *ring with identity*.

A *commutative ring* is a ring with a commutative multiplication. This property has profound implications on ring properties. Commutative algebra, the theory of commutative rings, is a major branch of ring theory. Its development has been greatly influenced by problems and ideas of algebraic number theory and algebraic geometry. In turn, commutative algebra is a fundamental tool in these branches of mathematics.

Examples of commutative rings include every field (such as the real or complex numbers), the integers, the polynomials in one or several variables with coefficients in another ring, the coordinate ring of an affine algebraic variety, and the ring of integers of a number field. Examples of noncommutative rings include the ring of *n* × *n* real square matrices with *n* ≥ 2, group rings in representation theory, operator algebras in functional analysis, rings of differential operators, and cohomology rings in topology.

The conceptualization of rings spanned the 1870s to the 1920s, with key contributions by Richard Dedekind, David Hilbert, Abraham Fraenkel, and Emmy Noether. Rings were first formalized as a generalization of Dedekind domains that occur in number theory, and of polynomial rings and rings of invariants that occur in algebraic geometry and invariant theory. They later proved useful in other branches of mathematics such as geometry and analysis.

Rings appear in the following chain of class inclusions:

## Definition
A **ring** is a set equipped with two binary operations for each ordered pair of elements of. + (addition) and ⋅ (multiplication) satisfying the following three sets of axioms, called the **ring axioms**:
# is an abelian group under addition, meaning that:
#* (*a* + *b*) + *c* = *a* + (*b* + *c*) for all *a*, *b*, *c* in (that is, + is associative).
#* *a* + *b* = *b* + *a* for all *a*, *b* in (that is, + is commutative).
#* There is an element 0 in such that *a* + 0 = *a* for all in (that is, 0 is an additive identity).
#* For each in there exists −*a* in such that *a* + (−*a*) = 0 (that is, −*a* is the additive inverse of ).
# is a monoid under multiplication, meaning that:
#* (*a* · *b*) · *c* = *a* · (*b* · *c*) for all *a*, *b*, *c* in (that is, ⋅ is associative).
#* There is an element 1 in such that *a* · 1 = *a* and 1 · *a* = *a* for all in (that is, 1 is a multiplicative identity).
# Multiplication is distributive with respect to addition, meaning that:
#* *a* · (*b* + *c*) = (*a* · *b*) + (*a* · *c*) for all *a*, *b*, *c* in (left distributivity).
#* (*b* + *c*) · *a* = (*b* · *a*) + (*c* · *a*) for all *a*, *b*, *c* in (right distributivity).

In notation, the multiplication symbol · is often omitted, in which case *a* · *b* is written as *ab*.

### Variations on terminology
In the terminology of this article, a ring is defined to have a multiplicative identity, while a structure with the same axiomatic definition but without the requirement for a multiplicative identity is instead called a "" (IPA: ) with a missing "i". For example, the set of even integers with the usual + and ⋅ is a rng, but not a ring. As explained in ** below, many authors apply the term "ring" without requiring a multiplicative identity.

Although ring addition is commutative, ring multiplication is not required to be commutative: need not necessarily equal *ba*. Rings that also satisfy commutativity for multiplication (such as the ring of integers) are called *commutative rings*. Books on commutative algebra or algebraic geometry often adopt the convention that *ring* means *commutative ring*, to simplify terminology.

In a ring, multiplicative inverses are not required to exist. A nonzero ring in which every nonzero element has a multiplicative inverse is called a division ring and a commutative division ring is called a field.

The additive group of a ring is the underlying set equipped with only the operation of addition. Although the definition requires that the additive group be abelian, this can be inferred from the other ring axioms. The proof makes use of the "1", and does not work in a rng. (For a rng, omitting the axiom of commutativity of addition leaves it inferable from the remaining rng assumptions only for elements that are products: *ab* + *cd* = *cd* + *ab*.)

Some authors use the term "ring" to refer to structures in which there is no requirement for multiplication to be associative; see the nonassociative ring subsection below. For these authors, every algebra is a "ring".

## Illustration

The most familiar example of a ring is the set of all integers consisting of the numbers
\(\dots,-5,-4,-3,-2,-1,0,1,2,3,4,5,\dots\)

The axioms of a ring are modeled on familiar properties of addition and multiplication of integers.

### Some properties
Some basic properties of a ring follow immediately from the axioms:
* The additive identity is unique.
* The additive inverse of each element is unique.
* The multiplicative identity is unique.
* For any element in a ring , one has *x*0 = 0 = 0*x* (zero is an absorbing element with respect to multiplication) and (–1)*x* = –*x*.
* If 0 = 1 in a ring (or more generally, 0 is a unit element), then has only one element, and is called the zero ring.
* If a ring contains the zero ring as a subring, then itself is the zero ring.
* The binomial formula holds for any and satisfying *xy* = *yx*.

### Example: Integers modulo 4

Equip the set \(\Z /4\Z = \left\{\overline{0}, \overline{1}, \overline{2}, \overline{3}\right\}\) with the following operations:
* The sum \(\overline{x} + \overline{y}\) in is the remainder when the integer *x* + *y* is divided by 4 (as *x* + *y* is always smaller than 8, this remainder is either *x* + *y* or *x* + *y* − 4). For example, \(\overline{2} + \overline{3} = \overline{1}\) and \(\overline{3} + \overline{3} = \overline{2}.\)
* The product \(\overline{x} \cdot \overline{y}\) in is the remainder when the integer is divided by 4. For example, \(\overline{2} \cdot \overline{3} = \overline{2}\) and \(\overline{3} \cdot \overline{3} = \overline{1}.\)

Then is a ring: each axiom follows from the corresponding axiom for If is an integer, the remainder of when divided by 4 may be considered as an element of and this element is often denoted by "*x* mod 4" or \(\overline x,\) which is consistent with the notation for 0, 1, 2, 3. The additive inverse of any \(\overline x\) in is \(-\overline x=\overline{-x}.\) For example, \(-\overline{3} = \overline{-3} = \overline{1}.\)

### Example: 2-by-2 matrices
The set of 2-by-2 square matrices with entries in a field is
\(\operatorname{M}_2(F) = \left\{ \left.\begin{pmatrix} a & b \\ c & d \end{pmatrix} \right|\ a, b, c, d \in F \right\}.\)

With the operations of matrix addition and matrix multiplication, \(\operatorname{M}_2(F)\) satisfies the above ring axioms. The element \(\left( \begin{smallmatrix} 1 & 0 \\ 0 & 1 \end{smallmatrix}\right)\) is the multiplicative identity of the ring. If \(A = \left( \begin{smallmatrix} 0 & 1 \\ 1 & 0 \end{smallmatrix} \right)\) and \(B = \left( \begin{smallmatrix} 0 & 1 \\ 0 & 0 \end{smallmatrix} \right),\) then \(AB = \left( \begin{smallmatrix} 0 & 0 \\ 0 & 1 \end{smallmatrix} \right)\) while \(BA = \left( \begin{smallmatrix} 1 & 0 \\ 0 & 0 \end{smallmatrix} \right);\) this example shows that the ring is noncommutative.

More generally, for any ring , commutative or not, and any nonnegative integer , the square *n* × *n* matrices with entries in form a ring; see *Matrix ring*.

## History

### Dedekind
The study of rings originated from the theory of polynomial rings and the theory of algebraic integers. In 1871, Richard Dedekind defined the concept of the ring of integers of a number field. In this context, he introduced the terms "ideal" (inspired by Ernst Kummer's notion of ideal number) and "module" and studied their properties. Dedekind did not use the term "ring" and did not define the concept of a ring in a general setting.

### Hilbert
The term "Zahlring" (number ring) was coined by David Hilbert in 1892 and published in 1897. According to Harvey Cohn, Hilbert used the term for a ring that had the property of "circling directly back" to an element of itself (in the sense of an equivalence). Specifically, in a ring of algebraic integers, all high powers of an algebraic integer can be written as an integral combination of a fixed set of lower powers, and thus the powers "cycle back". For instance, if *a* then:
\(\begin{aligned}
a^3 &= 4a-1, \\
a^4 &= 4a^2-a, \\
a^5 &= -a^2+16a-4, \\
a^6 &= 16a^2-8a+1, \\
a^7 &= -8a^2+65a-16, \\
\vdots \ & \qquad \vdots
\end{aligned}\)
and so on; in general, *a* is going to be an integral linear combination of 1, *a*, and *a*.

### Fraenkel and Noether
The first axiomatic definition of a ring was given by Abraham Fraenkel in 1915, but his axioms were stricter than those in the modern definition. For instance, he required every non-zero-divisor to have a multiplicative inverse.|ps=In 1921, Emmy Noether gave a modern axiomatic definition of commutative rings (with and without 1) and developed the foundations of commutative ring theory in her paper *Idealtheorie in Ringbereichen*.

### Multiplicative identity and the term "ring"
Fraenkel applied the term "ring" to structures with axioms that included a multiplicative identity,|ps=whereas Noether applied it to structures that did not.

Most or all books on algebra up to around 1960 followed Noether's convention of not requiring a 1 for a "ring". Starting in the 1960s, it became increasingly common to see books including the existence of 1 in the definition of "ring", especially in advanced books by notable authors such as Artin, Bourbaki, Eisenbud, and Lang. There are also books published as late as 2022 that use the term without the requirement for a 1. Likewise, the Encyclopedia of Mathematics does not require unit elements in rings. In a research article, the authors often specify which definition of ring they use in the beginning of that article.

Gardner and Wiegandt assert that, when dealing with several objects in the category of rings (as opposed to working with a fixed ring), if one requires all rings to have a 1, then some consequences include the lack of existence of infinite direct sums of rings, and that proper direct summands of rings are not subrings. They conclude that "in many, maybe most, branches of ring theory the requirement of the existence of a unity element is not sensible, and therefore unacceptable." Poonen makes the counterargument that the natural notion for rings would be the direct product rather than the direct sum. However, his main argument is that rings without a multiplicative identity are not totally associative, in the sense that they do not contain the product of any finite sequence of ring elements, including the empty sequence.

Authors who follow either convention for the use of the term "ring" may use one of the following terms to refer to objects satisfying the other convention:
* to include a requirement for a multiplicative identity: "unital ring", "unitary ring", "unit ring", "ring with unity", "ring with identity", "ring with a unit", or "ring with 1".
* to omit a requirement for a multiplicative identity: "rng" or "pseudo-ring", although the latter may be confusing because it also has other meanings.

## Basic examples

### Commutative rings
* The prototypical example is the ring of integers with the two operations of addition and multiplication.
* The rational, real and complex numbers are commutative rings of a type called fields.
* A unital associative algebra over a commutative ring is itself a ring as well as an -module. Some examples:
** The algebra *R*[*X*] of polynomials with coefficients in. ** The algebra \(RX_1, \dots, X_n\) of formal power series with coefficients in. ** The set of all continuous real-valued functions defined on the real line forms a commutative -algebra. The operations are pointwise addition and multiplication of functions.
** Let be a set, and let be a ring. Then the set of all functions from to forms a ring, which is commutative if is commutative.

* The ring of quadratic integers, the integral closure of in a quadratic extension of It is a subring of the ring of all algebraic integers.

* The ring of profinite integers the (infinite) product of the rings of -adic integers over all prime numbers. * The Hecke ring, the ring generated by Hecke operators.
* If is a set, then the power set of becomes a ring if we define addition to be the symmetric difference of sets and multiplication to be intersection. This is an example of a Boolean ring.

### Noncommutative rings
* For any ring and any natural number , the set of all square -by- matrices with entries from , forms a ring with matrix addition and matrix multiplication as operations. For *n* = 1, this matrix ring is isomorphic to itself. For *n* > 1 (and not the zero ring), this matrix ring is noncommutative.
* If *G* is an abelian group, then the endomorphisms of *G* form a ring, the endomorphism ring End(*G*) of *G*. The operations in this ring are addition and composition of endomorphisms. More generally, if is a left module over a ring , then the set of all -linear maps forms a ring, also called the endomorphism ring and denoted by End.
*The endomorphism ring of an elliptic curve. It is a commutative ring if the elliptic curve is defined over a field of characteristic zero.
* If *G* is a group and is a ring, the group ring of *G* over is a free module over having *G* as basis. Multiplication is defined by the rules that the elements of *G* commute with the elements of and multiply together as they do in the group *G*.
* The ring of differential operators (depending on the context). In fact, many rings that appear in analysis are noncommutative. For example, most Banach algebras are noncommutative.

### Non-rings
* The set of natural numbers with the usual operations is not a ring; since is not even a group (not all the elements are invertible with respect to addition – for instance, there is no natural number which can be added to 3 to get 0 as a result). There is a natural way to enlarge it to a ring, by including negative numbers to produce the ring of integers The natural numbers (including 0) form an algebraic structure known as a semiring (which has all of the axioms of a ring excluding that of an additive inverse).
* Let be the set of all continuous functions on the real line that vanish outside a bounded interval that depends on the function, with addition as usual but with multiplication defined as convolution:
\[
(f * g)(x) = \int_{-\infty}^\infty f(y)g(x - y) \, dy.
\]
 Then is a rng, but not a ring: the Dirac delta function has the property of a multiplicative identity, but it is not a function and hence is not an element of. ## Basic concepts
### Products and powers
For each nonnegative integer , given a sequence of elements of , one can define the product \(\textstyle P_n = \prod_{i=1}^n a_i\) recursively: let *P*\(_{0}\) = 1 and let *P*\(_{*m*}\) = *P*\(_{*m*−1}\)*a*\(_{*m*}\) for 1 ≤ *m* ≤ *n*.

As a special case, one can define nonnegative integer powers of an element of a ring: *a* and *a**a*for *n* ≥ 1. Then *a*\(^{*m*+*n*}\) = *a*\(^{*m*}\)*a*\(^{*n*}\) for all *m*, *n* ≥ 0.

### Elements in a ring
A left zero divisor of a ring is an element in the ring such that there exists a nonzero element of such that *ab* = 0. A right zero divisor is defined similarly.

A nilpotent element is an element such that *a for some *n'' > 0. One example of a nilpotent element is a nilpotent matrix. A nilpotent element in a nonzero ring is necessarily a zero divisor.

An idempotent \(e\) is an element such that *e*. One example of an idempotent element is a projection in linear algebra.

A unit is an element having a multiplicative inverse; in this case the inverse is unique, and is denoted by *a*. The set of units of a ring is a group under ring multiplication; this group is denoted by *R* or *R** or *U*(*R*). For example, if is the ring of all square matrices of size over a field, then *R* consists of the set of all invertible matrices of size , and is called the general linear group.

### Subring

A subset of is called a subring if any one of the following equivalent conditions holds:
* the addition and multiplication of restrict to give operations *S* × *S* → *S* making a ring with the same multiplicative identity as. * 1 ∈ *S*; and for all in , the elements , *x* + *y*, and are in. * can be equipped with operations making it a ring such that the inclusion map *S* → *R* is a ring homomorphism.

For example, the ring of integers is a subring of the field of real numbers and also a subring of the ring of polynomials (in both cases, contains 1, which is the multiplicative identity of the larger rings). On the other hand, the subset of even integers does not contain the identity element 1 and thus does not qualify as a subring of  one could call a subrng, however.

An intersection of subrings is a subring. Given a subset of , the smallest subring of containing is the intersection of all subrings of containing , and it is called *the subring generated by E*.

For a ring , the smallest subring of is called the *characteristic subring* of. It can be generated through addition of copies of 1 and −1. It is possible that *n* · 1 = 1 + 1 +... + 1 ( times) can be zero. If is the smallest positive integer such that this occurs, then is called the *characteristic* of. In some rings, *n* · 1 is never zero for any positive integer , and those rings are said to have *characteristic zero*.

Given a ring , let Z(*R*) denote the set of all elements in such that commutes with every element in : *xy* = *yx* for any in. Then Z(*R*) is a subring of , called the center of. More generally, given a subset of , let be the set of all elements in that commute with every element in. Then is a subring of , called the centralizer (or commutant) of. The center is the centralizer of the entire ring. Elements or subsets of the center are said to be *central* in ; they (each individually) generate a subring of the center.

### Ideal

Let be a ring. A **left ideal** of is a nonempty subset of such that for any in and in , the elements *x* + *y* and are in. Let denote the -span of , that is, the set of finite sums
\(r_1 x_1 + \cdots + r_n x_n \quad \textrm{such}\;\textrm{that}\; r_i \in R \; \textrm{ and } \; x_i \in I;\)
then is a left ideal if *RI* ⊆ *I*. Similarly, a **right ideal** is a subset such that *IR* ⊆ *I*. A subset is said to be a **two-sided ideal** or simply **ideal** if it is both a left ideal and right ideal. A one-sided or two-sided ideal is then an additive subgroup of. If is a subset of , then *RE* is a left ideal, called the left ideal generated by ; it is the smallest left ideal containing. Similarly, one can consider the right ideal or the two-sided ideal generated by a subset of. If is in , then *Rx* is a left ideal, and *xR* is a right ideal; they are called the principal left ideal and right ideal generated by. The principal ideal *RxR* is written as (*x*). For example, the set of all positive and negative multiples of 2 along with 0 form an ideal of the integers, and this ideal is generated by the integer 2. In fact, every ideal of the ring of integers is principal.

A ring is simple if it is nonzero and it has no proper nonzero two-sided ideals. A commutative simple ring is precisely a field.

Rings are often studied with special conditions set upon their ideals. For example, a ring in which there is no strictly increasing infinite chain of left ideals is called a left Noetherian ring. A ring in which there is no strictly decreasing infinite chain of left ideals is called a left Artinian ring. It is a somewhat surprising fact that a left Artinian ring is left Noetherian (the Hopkins–Levitzki theorem). The integers, however, form a Noetherian ring which is not Artinian.

For commutative rings, the ideals generalize the classical notion of divisibility and decomposition of an integer into prime numbers in algebra. A proper ideal of is called a prime ideal if for any elements \(x, y\in R\) we have that \(xy \in P\) implies either \(x \in P\) or \(y\in P.\) Equivalently, is prime if for any ideals *I*, *J* we have that *IJ* ⊆ *P* implies either *I* ⊆ *P* or *J* ⊆ *P*. This latter formulation illustrates the idea of ideals as generalizations of elements.

### Homomorphism

A **homomorphism** from a ring (*R*, +, **⋅**) to a ring (*S*, ‡, ∗) is a function from to  that preserves the ring operations; namely, such that, for all *a*, *b* in the following identities hold:
\(\begin{aligned}
& f(a+b) = f(a) \ddagger f(b) \\
& f(a\cdot b) = f(a)*f(b) \\
& f(1_R) = 1_S
\end{aligned}\)

If one is working with , then the third condition is dropped.

A ring homomorphism is said to be an **isomorphism** if there exists an inverse homomorphism to (that is, a ring homomorphism that is an inverse function), or equivalently if it is bijective.

Examples:
* The function that maps each integer to its remainder modulo 4 (a number in ) is a homomorphism from the ring to the quotient ring ("quotient ring" is defined below).
* If is a unit element in a ring , then \(R \to R, x \mapsto uxu^{-1}\) is a ring homomorphism, called an inner automorphism of. * Let be a commutative ring of prime characteristic. Then *x* ↦ is a ring endomorphism of called the Frobenius homomorphism.
* The Galois group of a field extension *L* / *K* is the set of all automorphisms of whose restrictions to are the identity.
* For any ring , there is a unique ring homomorphism , and there is a unique ring homomorphism. * An algebra homomorphism from a -algebra to the endomorphism algebra of a vector space over is called a representation of the algebra.

Given a ring homomorphism *f* : *R* → *S*, the set of all elements mapped to 0 by is called the kernel of. The kernel is a two-sided ideal of. The image of , on the other hand, is not always an ideal, but it is always a subring of. To give a ring homomorphism from a commutative ring to a ring with image contained in the center of is the same as to give a structure of an algebra over to  (which in particular gives a structure of an -module).

### Quotient ring

The notion of quotient ring is analogous to the notion of a quotient group. Given a ring (*R*, +, **⋅**) and a two-sided ideal of (*R*, +, **⋅**), view as subgroup of (*R*, +); then the **quotient ring** *R* / *I* is the set of cosets of together with the operations
\(\begin{aligned}
& (a+I)+(b+I) = (a+b)+I, \\
& (a+I)(b+I) = (ab)+I.
\end{aligned}\)
for all *a*, *b* in. The ring *R* / *I* is also called a **factor ring**.

As with a quotient group, there is a canonical homomorphism *p* : *R* → *R* / *I*, given by *x* ↦ *x* + *I*. It is surjective and satisfies the following universal property:
* If *f* : *R* → *S* is a ring homomorphism such that *f*(*I*) = 0, then there is a unique homomorphism \(\overline{f} : R/I \to S\) such that \(f = \overline{f} \circ p.\)
For any ring homomorphism *f* : *R* → *S*, invoking the universal property with *I* = ker *f* produces a homomorphism \(\overline{f} : R / \ker f \to S\) that gives an isomorphism from *R* / ker *f* to the image of. ## Modules

The concept of a *module over a ring* generalizes the concept of a vector space (over a field) by generalizing from multiplication of vectors with elements of a field (scalar multiplication) to multiplication with elements of a ring. More precisely, given a ring , an -module is an abelian group equipped with an operation *R* × *M* → *M* (associating an element of to every pair of an element of and an element of ) that satisfies certain axioms. This operation is commonly denoted by juxtaposition and called multiplication. The axioms of modules are the following: for all *a*, *b* in and all *x*, *y* in ,
 is an abelian group under addition.
\(\begin{aligned}
& a(x+y) = ax+ay \\
& (a+b)x = ax+bx \\
& 1x = x \\
& (ab)x = a(bx)
\end{aligned}\)
When the ring is noncommutative these axioms define *left modules*; *right modules* are defined similarly by writing instead of. This is not only a change of notation, as the last axiom of right modules (that is *x*(*ab*) = (*xa*)*b*) becomes (*ab*)*x* = *b*(*ax*), if left multiplication (by ring elements) is used for a right module.

Basic examples of modules are ideals, including the ring itself.

Although similarly defined, the theory of modules is much more complicated than that of vector space, mainly, because, unlike vector spaces, modules are not characterized (up to an isomorphism) by a single invariant (the dimension of a vector space). In particular, not all modules have a basis.

The axioms of modules imply that (−1)*x* = −*x*, where the first minus denotes the additive inverse in the ring and the second minus the additive inverse in the module. Using this and denoting repeated addition by a multiplication by a positive integer allows identifying abelian groups with modules over the ring of integers.

Any ring homomorphism induces a structure of a module: if *f* : *R* → *S* is a ring homomorphism, then is a left module over by the multiplication: *rs* = *f*(*r*)*s*. If is commutative or if *f*(*R*) is contained in the center of , the ring is called a -algebra. In particular, every ring is an algebra over the integers.

## Constructions
### Direct product

Let and be rings. Then the product *R* × *S* can be equipped with the following natural ring structure:
\(\begin{aligned}
& (r_1,s_1) + (r_2,s_2) = (r_1+r_2,s_1+s_2) \\
& (r_1,s_1) \cdot (r_2,s_2)=(r_1\cdot r_2,s_1\cdot s_2)
\end{aligned}\)
for all *r*in and *s*in. The ring *R* × *S* with the above operations of addition and multiplication and the multiplicative identity (1, 1) is called the **direct product** of with. The same construction also works for an arbitrary family of rings: if are rings indexed by a set , then \(\prod_{i \in I} R_i\) is a ring with componentwise addition and multiplication.

Let be a commutative ring and \(\mathfrak{a}_1, \cdots, \mathfrak{a}_n\) be ideals such that \(\mathfrak{a}_i + \mathfrak{a}_j = (1)\) whenever *i* ≠ *j*. Then the Chinese remainder theorem says there is a canonical ring isomorphism:

\[
R /{\textstyle \bigcap_{i=1}^{n}{\mathfrak{a}_i\simeq \prod_{i=1}^{n}{R/ \mathfrak{a}_i}, \qquad x \bmod {\textstyle \bigcap_{i=1}^{n}\mathfrak{a}_i} \mapsto (x \bmod \mathfrak{a}_1, \ldots , x \bmod \mathfrak{a}_n).
\]

A "finite" direct product may also be viewed as a direct sum of ideals. Namely, let \(R_i, 1 \le i \le n\) be rings, \(R_i \to R = \prod R_i\) the inclusions with the images \(\mathfrak{a}_i\) (in particular \(\mathfrak{a}_i\) are rings though not subrings). Then \(\mathfrak{a}_i\) are ideals of and

\[
R = \mathfrak{a}_1 \oplus \cdots \oplus \mathfrak{a}_n, \quad \mathfrak{a}_i \mathfrak{a}_j = 0, i \ne j, \quad \mathfrak{a}_i^2 \subseteq \mathfrak{a}_i
\]

as a direct sum of abelian groups (because for abelian groups finite products are the same as direct sums). Clearly the direct sum of such ideals also defines a product of rings that is isomorphic to. Equivalently, the above can be done through central idempotents. Assume that has the above decomposition. Then we can write

\[
1 = e_1 + \cdots + e_n, \quad e_i \in \mathfrak{a}_i.
\]

By the conditions on \(\mathfrak{a}_i,\) one has that are central idempotents and *e* = 0, *i* ≠ *j* (orthogonal). Again, one can reverse the construction. Namely, if one is given a partition of 1 in orthogonal central idempotents, then let \(\mathfrak{a}_i = R e_i,\) which are two-sided ideals. If each is not a sum of orthogonal central idempotents, then their direct sum is isomorphic to. An important application of an infinite direct product is the construction of a projective limit of rings (see below). Another application is a restricted product of a family of rings (cf. adele ring).

### Polynomial ring

Given a symbol (called a variable) and a commutative ring , the set of polynomials
\(R[t] = \left\{ a_n t^n + a_{n-1} t^{n -1} + \dots + a_1 t + a_0 \mid n \ge 0, a_j \in R \right\}\)
forms a commutative ring with the usual addition and multiplication, containing as a subring. It is called the polynomial ring over. More generally, the set \(R\left[t_1, \ldots, t_n\right]\) of all polynomials in variables \(t_1, \ldots, t_n\) forms a commutative ring, containing \(R\left[t_i\right]\) as subrings.

If is an integral domain, then *R*[*t*] is also an integral domain; its field of fractions is the field of rational functions. If is a Noetherian ring, then *R*[*t*] is a Noetherian ring. If is a unique factorization domain, then *R*[*t*] is a unique factorization domain. Finally, is a field if and only if *R*[*t*] is a principal ideal domain.

Let \(R \subseteq S\) be commutative rings. Given an element of , one can consider the ring homomorphism
\(R[t] \to S, \quad f \mapsto f(x)\)

(that is, the substitution). If *S* = *R*[*t*] and *x* = *t*, then *f*(*t*) = *f*. Because of this, the polynomial is often also denoted by *f*(*t*). The image of the map is denoted by *R*[*x*]; it is the same thing as the subring of generated by and. Example: \(k\left[t^2, t^3\right]\) denotes the image of the homomorphism
\(k[x, y] \to k[t], \, f \mapsto f\left(t^2, t^3\right).\)

In other words, it is the subalgebra of *k*[*t*] generated by *t* and *t*.

Example: let be a polynomial in one variable, that is, an element in a polynomial ring. Then *f*(*x* + *h*) is an element in *R*[*h*] and *f*(*x* + *h*) – *f*(*x*) is divisible by in that ring. The result of substituting zero to in (*f*(*x* + *h*) – *f*(*x*)) / *h* is *f' *(*x*), the derivative of at. The substitution is a special case of the universal property of a polynomial ring. The property states: given a ring homomorphism \(\phi: R \to S\) and an element in there exists a unique ring homomorphism \(\overline{\phi}: R[t] \to S\) such that \(\overline{\phi}(t) = x\) and \(\overline{\phi}\) restricts to. For example, choosing a basis, a symmetric algebra satisfies the universal property and so is a polynomial ring.

To give an example, let be the ring of all functions from to itself; the addition and the multiplication are those of functions. Let be the identity function. Each in defines a constant function, giving rise to the homomorphism *R* → *S*. The universal property says that this map extends uniquely to
\(R[t] \to S, \quad f \mapsto \overline{f}\)
( maps to ) where \(\overline{f}\) is the polynomial function defined by. The resulting map is injective if and only if is infinite.

Given a non-constant monic polynomial in *R*[*t*], there exists a ring containing such that is a product of linear factors in *S*[*t*].

Let be an algebraically closed field. The Hilbert's Nullstellensatz (theorem of zeros) states that there is a natural one-to-one correspondence between the set of all prime ideals in \(k\left[t_1, \ldots, t_n\right]\) and the set of closed subvarieties of. In particular, many local problems in algebraic geometry may be attacked through the study of the generators of an ideal in a polynomial ring. (cf. Gröbner basis.)

There are some other related constructions. A formal power series ring \(R[\![t]\!]\) consists of formal power series
\(\sum_0^\infty a_i t^i, \quad a_i \in R\)

together with multiplication and addition that mimic those for convergent series. It contains *R*[*t*] as a subring. A formal power series ring does not have the universal property of a polynomial ring; a series may not converge after a substitution. The important advantage of a formal power series ring over a polynomial ring is that it is local (in fact, complete).

### Matrix ring and endomorphism ring

Let be a ring (not necessarily commutative). The set of all square matrices of size with entries in forms a ring with the entry-wise addition and the usual matrix multiplication. It is called the matrix ring and is denoted by M. Given a right -module , the set of all -linear maps from to itself forms a ring with addition that is of function and multiplication that is of composition of functions; it is called the endomorphism ring of and is denoted by End.

As in linear algebra, a matrix ring may be canonically interpreted as an endomorphism ring: \(\operatorname{End}_R(R^n) \simeq \operatorname{M}_n(R).\) This is a special case of the following fact: If \(f: \oplus_1^n U \to \oplus_1^n U\) is an -linear map, then may be written as a matrix with entries in *S* = End, resulting in the ring isomorphism:
\(\operatorname{End}_R(\oplus_1^n U) \to \operatorname{M}_n(S), \quad f \mapsto (f_{ij}).\)

Any ring homomorphism *R* → *S* induces M(*S*).

Schur's lemma says that if is a simple right -module, then End is a division ring. If \(U = \bigoplus_{i = 1}^r U_i^{\oplus m_i}\) is a direct sum of -copies of simple -modules \(U_i,\) then
\(\operatorname{End}_R(U) \simeq \prod_{i=1}^r \operatorname{M}_{m_i} (\operatorname{End}_R(U_i)).\)
The Artin–Wedderburn theorem states any semisimple ring (cf. below) is of this form.

A ring and the matrix ring M over it are Morita equivalent: the category of right modules of is equivalent to the category of right modules over M. In particular, two-sided ideals in correspond in one-to-one to two-sided ideals in M.

### Limits and colimits of rings
Let be a sequence of rings such that is a subring of *R* for all. Then the union (or filtered colimit) of is the ring \(\varinjlim R_i\) defined as follows: it is the disjoint union of all 's modulo the equivalence relation *x* ~ *y* if and only if *x* = *y* in for sufficiently large. Examples of colimits:
* A polynomial ring in infinitely many variables: \(R[t_1, t_2, \cdots] = \varinjlim R[t_1, t_2, \cdots, t_m].\)
* The algebraic closure of finite fields of the same characteristic \(\overline{\mathbf{F_p = \varinjlim \mathbf{F}_{p^m}.\)
* The function field of an algebraic variety over a field is \(\varinjlim k[U]\) where the limit runs over all the coordinate rings *k*[*U*] of nonempty open subsets (more succinctly it is the stalk of the structure sheaf at the generic point.)

Any ring is the filtered colimit (union) of its finitely generated subrings.

A projective limit (or a filtered limit) of rings is defined as follows. Suppose we are given a family of rings *R*, *i* running over positive integers, say, and ring homomorphisms *R*, *j* ≥ *i* such that *R*are all the identities and *R* → *R*is *R*whenever *k* ≥ *j* ≥ *i*. Then \(\varprojlim R_i\) is the subring of \(\textstyle \prod R_i\) consisting of (*x* such that *x* maps to *x* under *R*, *j* ≥ *i*.

For an example of a projective limit, see **.

### Localization
The localization generalizes the construction of the field of fractions of an integral domain to an arbitrary ring and modules. Given a (not necessarily commutative) ring and a subset of , there exists a ring \(R[S^{-1}]\) together with the ring homomorphism \(R \to R\left[S^{-1}\right]\) that "inverts" ; that is, the homomorphism maps elements in to unit elements in \(R\left[S^{-1}\right],\) and, moreover, any ring homomorphism from that "inverts" uniquely factors through \(R\left[S^{-1}\right].\) The ring \(R\left[S^{-1}\right]\) is called the **localization** of with respect to. For example, if is a commutative ring and an element in , then the localization \(R\left[f^{-1}\right]\) consists of elements of the form \(r/f^n, \, r \in R , \, n \ge 0\) (to be precise, \(R\left[f^{-1}\right] = R[t]/(tf - 1).\))

The localization is frequently applied to a commutative ring with respect to the complement of a prime ideal (or a union of prime ideals) in. In that case \(S = R - \mathfrak{p},\) one often writes \(R_\mathfrak{p}\) for \(R\left[S^{-1}\right].\) \(R_\mathfrak{p}\) is then a local ring with the maximal ideal \(\mathfrak{p} R_\mathfrak{p}.\) This is the reason for the terminology "localization". The field of fractions of an integral domain is the localization of at the prime ideal zero. If \(\mathfrak{p}\) is a prime ideal of a commutative ring , then the field of fractions of \(R/\mathfrak{p}\) is the same as the residue field of the local ring \(R_\mathfrak{p}\) and is denoted by \(k(\mathfrak{p}).\)

If is a left -module, then the localization of with respect to is given by a change of rings \(M\left[S^{-1}\right] = R\left[S^{-1}\right] \otimes_R M.\)

The most important properties of localization are the following: when is a commutative ring and a multiplicatively closed subset
* \(\mathfrak{p} \mapsto \mathfrak{p}\left[S^{-1}\right]\) is a bijection between the set of all prime ideals in disjoint from and the set of all prime ideals in \(R\left[S^{-1}\right].\)
* \(R\left[S^{-1}\right] = \varinjlim R\left[f^{-1}\right],\) running over elements in with partial ordering given by divisibility.
* The localization is exact:
\[
0 \to M'\left[S^{-1}\right] \to M\left[S^{-1}\right] \to M''\left[S^{-1}\right] \to 0
\]
 is exact over \(R\left[S^{-1}\right]\) whenever \(0 \to M' \to M \to M'' \to 0\) is exact over. * Conversely, if \(0 \to M'_\mathfrak{m} \to M_\mathfrak{m} \to M*_\mathfrak{m} \to 0\) is exact for any maximal ideal \(\mathfrak{m},\) then \(0 \to M' \to M \to M* \to 0\) is exact.
* A remark: localization is no help in proving a global existence. One instance of this is that if two modules are isomorphic at all prime ideals, it does not follow that they are isomorphic. (One way to explain this is that the localization allows one to view a module as a sheaf over prime ideals and a sheaf is inherently a local notion.)

In category theory, a localization of a category amounts to making some morphisms isomorphisms. An element in a commutative ring may be thought of as an endomorphism of any -module. Thus, categorically, a localization of with respect to a subset of is a functor from the category of -modules to itself that sends elements of viewed as endomorphisms to automorphisms and is universal with respect to this property. (Of course, then maps to \(R\left[S^{-1}\right]\) and -modules map to \(R\left[S^{-1}\right]\)-modules.)

### Completion
Let be a commutative ring, and let be an ideal of. The **completion** of at is the projective limit \(\hat{R} = \varprojlim R/I^n;\) it is a commutative ring. The canonical homomorphisms from to the quotients \(R/I^n\) induce a homomorphism \(R \to \hat{R}.\) The latter homomorphism is injective if is a Noetherian integral domain and is a proper ideal, or if is a Noetherian local ring with maximal ideal , by Krull's intersection theorem. The construction is especially useful when is a maximal ideal.

The basic example is the completion of at the principal ideal (*p*) generated by a prime number ; it is called the ring of -adic integers and is denoted The completion can in this case be constructed also from the -adic absolute value on The -adic absolute value on is a map \(x \mapsto |x|\) from to given by \(|n|_p=p^{-v_p(n)}\) where \(v_p(n)\) denotes the exponent of in the prime factorization of a nonzero integer into prime numbers (we also put \(|0|_p=0\) and \(|m/n|_p = |m|_p/|n|_p\)). It defines a distance function on and the completion of as a metric space is denoted by It is again a field since the field operations extend to the completion. The subring of consisting of elements with |*x*| is isomorphic to 

Similarly, the formal power series ring math
