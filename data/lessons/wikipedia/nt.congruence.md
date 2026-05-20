> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Congruence_relation) — CC BY-SA 4.0

# Congruence relations

In abstract algebra, a **congruence relation** (or simply **congruence**) is an equivalence relation on an algebraic structure (such as a group, ring, or vector space) that is compatible with the structure in the sense that algebraic operations done with equivalent elements will yield equivalent elements.  Every congruence relation has a corresponding quotient structure, whose elements are the equivalence classes (or **congruence classes**) for the relation.

## Definition
The definition of a congruence depends on the type of algebraic structure under consideration.  Particular definitions of congruence can be made for groups, rings, vector spaces, modules, semigroups, lattices, and so forth.  The common theme is that a congruence is an equivalence relation on an algebraic object that is compatible with the algebraic structure, in the sense that the operations are well-defined on the equivalence classes.

### General
The general notion of a congruence relation can be formally defined in the context of universal algebra, a field which studies ideas common to all algebraic structures.  In this setting, a relation \(R\) on a given algebraic structure is called **compatible** if for each \(n\) and each \(n\)-ary operation \(\mu\) defined on the structure: whenever \(a_1 \mathrel{R} a'_1\) and ... and \(a_n \mathrel{R} a'_n\), then \(\mu(a_1,\ldots,a_n) \mathrel{R} \mu(a'_1,\ldots,a'_n)\).

A congruence relation on the structure is then defined as an equivalence relation that is also compatible.

## Examples
### Basic example


The prototypical example of a congruence relation is congruence modulo \(n\) on the set of integers.  For a given positive integer \(n\), two integers \(a\) and \(b\) are called **congruent modulo \(n\)**, written
\(a \equiv b \pmod{n}\)
if \(a - b\) is divisible by \(n\) (or equivalently if \(a\) and \(b\) have the same remainder when divided by \(n\)).

For example, \(37\) and \(57\) are congruent modulo \(10\),
\(37 \equiv 57 \pmod{10}\)
since \(37 - 57 = -20\) is a multiple of 10, or equivalently since both \(37\) and \(57\) have a remainder of \(7\) when divided by \(10\).

Congruence modulo \(n\) (for a fixed \(n\)) is compatible with both addition and multiplication on the integers.  That is,

if
\(a_1 \equiv a_2 \pmod{n}\) and \(b_1 \equiv b_2 \pmod{n}\)
then
\(a_1 + b_1 \equiv a_2 + b_2 \pmod{n}\)  and  \(a_1 b_1 \equiv a_2b_2 \pmod{n}\)

The corresponding addition and multiplication of equivalence classes is known as modular arithmetic.  From the point of view of abstract algebra, congruence modulo \(n\) is a congruence relation on the ring of integers, and arithmetic modulo \(n\) occurs on the corresponding quotient ring.

### Example: Groups
For example, a group is an algebraic object consisting of a set together with a single binary operation, satisfying certain axioms.  If \(G\) is a group with operation \(\ast\), a **congruence relation** on \(G\) is an equivalence relation \(\equiv\) on the elements of \(G\) satisfying
\(g_1 \equiv g_2 \ \ \,\) and \(\ \ \, h_1 \equiv h_2 \implies g_1 \ast h_1 \equiv g_2 \ast h_2\)
for all \(g_1, g_2, h_1, h_2 \in G\). For a congruence on a group, the equivalence class containing the identity element is always a normal subgroup, and the other equivalence classes are the other cosets of this subgroup.  Together, these equivalence classes are the elements of a quotient group.

### Example: Rings
When an algebraic structure includes more than one operation, congruence relations are required to be compatible with each operation.  For example, a ring possesses both addition and multiplication, and a congruence relation on a ring must satisfy
\(r_1 + s_1 \equiv r_2 + s_2\) and \(r_1 s_1 \equiv r_2 s_2\)
whenever \(r_1 \equiv r_2\) and \(s_1 \equiv s_2\). For a congruence on a ring, the equivalence class containing 0 is always a two-sided ideal, and the two operations on the set of equivalence classes define the corresponding quotient ring.

## Relation with homomorphisms
If \(f:A\, \rightarrow B\) is a homomorphism between two algebraic structures (such as homomorphism of groups, or a linear map between vector spaces), then the relation \(R\) defined by
\(a_1\, R\, a_2\) if and only if \(f(a_1) = f(a_2)\)
is a congruence relation on \(A\).  By the first isomorphism theorem, the image of *A* under \(f\) is a substructure of *B* isomorphic to the quotient of *A* by this congruence.

On the other hand, the congruence relation \(R\) induces a unique homomorphism \(f: A \rightarrow A/R\) given by
\(f(x) = \{y \mid x \, R \, y\}\).

Thus, there is a natural correspondence between the congruences and the homomorphisms of any given algebraic structure.

## Congruences of groups, and normal subgroups and ideals
In the particular case of groups, congruence relations can be described in elementary terms as follows:
If *G* is a group (with identity element *e* and operation *) and ~ is a binary relation on *G*, then ~ is a congruence whenever:
# Given any element *a* of *G*, *a* ~ *a* (**reflexivity**);
# Given any elements *a* and *b* of *G*, if *a* ~ *b*, then *b* ~ *a* (**symmetry**);
# Given any elements *a*, *b*, and *c* of *G*, if *a* ~ *b* and *b* ~ *c*, then *a* ~ *c* (**transitivity**);
# Given any elements *a*, *a*′, *b*, and *b*′ of *G*, if *a* ~ *a*′ and *b* ~ *b*′, then *a* * *b* ~ *a*′ * *b*′;
# Given any elements *a* and *a*′ of *G*, if *a* ~ *a*′, then *a*<sup>−1</sup> ~ *a*′<sup>−1</sup> (this is implied by the other four, so is strictly redundant).

Conditions 1, 2, and 3 say that ~ is an equivalence relation.

A congruence ~ is determined entirely by the set  of those elements of *G* that are congruent to the identity element, and this set is a normal subgroup.
Specifically, *a* ~ *b* if and only if *b*<sup>−1</sup> * *a* ~ *e*.
So instead of talking about congruences on groups, people usually speak in terms of normal subgroups of them; in fact, every congruence corresponds uniquely to some normal subgroup of *G*.

### Ideals of rings and the general case
A similar trick allows one to speak of kernels in ring theory as ideals instead of congruence relations, and in module theory as submodules instead of congruence relations.

A more general situation where this trick is possible is with Omega-groups (in the general sense allowing operators with multiple arity). But this cannot be done with, for example, monoids, so the study of congruence relations plays a more central role in monoid theory.

## Universal algebra
The general notion of a congruence is particularly useful in universal algebra. An equivalent formulation in this context is the following:

A congruence relation on an algebra *A* is a subset of the direct product *A* × *A* that is both an equivalence relation on *A* and a subalgebra of *A* × *A*.

The kernel of a homomorphism is always a congruence. Indeed, every congruence arises as a kernel.
For a given congruence ~ on *A*, the set *A* / ~ of equivalence classes can be given the structure of an algebra in a natural fashion, the quotient algebra.
The function that maps every element of *A* to its equivalence class is a homomorphism, and the kernel of this homomorphism is ~.

The lattice **Con**(*A*) of all congruence relations on an algebra *A* is algebraic.

John M. Howie described how semigroup theory illustrates congruence relations in universal algebra:
In a group a congruence is determined if we know a single congruence class, in particular if we know the normal subgroup which is the class containing the identity. Similarly, in a ring a congruence is determined if we know the ideal which is the congruence class containing the zero. In semigroups there is no such fortunate occurrence, and we are therefore faced with the necessity of studying congruences as such. More than anything else, it is this necessity that gives semigroup theory its characteristic flavour. Semigroups are in fact the first and simplest type of algebra to which the methods of universal algebra must be applied ...

## Category theory
In category theory, a congruence relation *R* on a category *C* is given by: for each pair of objects *X*, *Y* in *C*, an equivalence relation *R*<sub>*X*,*Y*</sub> on Hom(*X*,*Y*), such that the equivalence relations respect composition of morphisms. See  for details.

## See also
* Chinese remainder theorem
* Congruence lattice problem
* Table of congruences

## Explanatory notes


## Notes


## References
*
*
*  (Section 4.5 discusses congruency of matrices.)
*
*
*

