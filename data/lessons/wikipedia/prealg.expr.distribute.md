> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Distributive_property) — CC BY-SA 4.0

# Distributive property

| statement =
| symbolic statement =
}}

In mathematics, the **distributive property** of binary operations is a generalization of the **distributive law**, which asserts that the equality

\[
x \cdot (y + z) = x \cdot y + x \cdot z
\]

is always true in elementary algebra.
For example, in elementary arithmetic, one has

\[
2 \cdot (1 + 3) = (2 \cdot 1) + (2 \cdot 3).
\]

Therefore, one would say that multiplication *distributes* over addition.

This basic property of numbers is part of the definition of most algebraic structures that have two operations called addition and multiplication, such as complex numbers, polynomials, matrices, rings, and fields. It is also encountered in Boolean algebra and mathematical logic, where each of the logical and (denoted \(\,\land\,\)) and the logical or (denoted \(\,\lor\,\)) distributes over the other.

## Definition
Given a set \(S\) and two binary operators \(\,*\,\) and \(\,+\,\) on \(S,\)
*the operation \(\,*\,\) is  over (or with respect to) \(\,+\,\) if, given any elements \(x, y, \text{ and } z\) of \(S,\)
\(x * (y + z) = (x * y) + (x * z);\)
*the operation \(\,*\,\) is  over \(\,+\,\) if, given any elements \(x, y, \text{ and } z\) of \(S,\)
\((y + z) * x = (y * x) + (z * x);\)
*and the operation \(\,*\,\) is  over \(\,+\,\) if it is left- and right-distributive.

When \(\,*\,\) is commutative, the three conditions above are logically equivalent.

## Meaning
The operators used for examples in this section are those of the usual addition \(\,+\,\) and multiplication \(\,\cdot.\,\)

If the operation denoted \(\cdot\) is not commutative, there is a distinction between left-distributivity and right-distributivity:


\[
a \cdot \left( b \pm c \right) = a \cdot b \pm a \cdot c \qquad \text{ (left-distributive) }
\]


\[
(a \pm b) \cdot c = a \cdot c \pm b \cdot c \qquad \text{ (right-distributive) }.
\]


In either case, the distributive property can be described in words as:

To multiply a sum (or difference) by a factor, each summand (or minuend and subtrahend) is multiplied by this factor and the resulting products are added (or subtracted).

If the operation outside the parentheses (in this case, the multiplication) is commutative, then left-distributivity implies right-distributivity and vice versa, and one talks simply of .

One example of an operation that is "only" right-distributive is division, which is not commutative:

\[
(a \pm b) \div c = a \div c \pm b \div c.
\]

In this case, left-distributivity does not apply:

\[
a \div(b \pm c) \neq a \div b \pm a \div c
\]


The distributive laws are among the axioms for rings (like the ring of integers) and fields (like the field of rational numbers). Here multiplication is distributive over addition, but addition is not distributive over multiplication. Examples of structures with two operations that are each distributive over the other are Boolean algebras such as the algebra of sets or the switching algebra.

Multiplying sums can be put into words as follows: When a sum is multiplied by a sum, multiply each summand of a sum with each summand of the other sum (keeping track of signs) then add up all of the resulting products.

## Examples
### Real numbers
In the following examples, the use of the distributive law on the set of real numbers \(\R\) is illustrated. When multiplication is mentioned in elementary mathematics, it usually refers to this kind of multiplication. From the point of view of algebra, the real numbers form a field, which ensures the validity of the distributive law.


{{defn|

\[
\begin{align}
(a + b) \cdot (a - b) & = a \cdot (a - b) + b \cdot (a - b) = a^2 - ab + ba - b^2 = a^2 - b^2 \\
                      & = (a + b) \cdot a - (a + b) \cdot b = a^2 + ba - ab - b^2 = a^2 - b^2 \\
\end{align}
\]


Here the distributive law was applied twice, and it does not matter which bracket is first multiplied out.
}}


### Matrices
The distributive law is valid for matrix multiplication. More precisely,

\[
(A + B) \cdot C = A \cdot C + B \cdot C
\]

for all  \(l \times m\)-matrices \(A, B\) and \(m \times n\)-matrices \(C,\) as well as

\[
A \cdot (B + C) = A \cdot B + A \cdot C
\]

for all \(l \times m\)-matrices \(A\) and \(m \times n\)-matrices \(B, C.\)
Because the commutative property does not hold for matrix multiplication, the second law does not follow from the first law. In this case, they are two different laws.

### Other examples
* Multiplication of ordinal numbers, in contrast, is only left-distributive, not right-distributive.
* The cross product is left- and right-distributive over vector addition, though not commutative.
* For sets, the union is distributive over intersection, and intersection is distributive over union.
* Logical disjunction ("or") is distributive over logical conjunction ("and"), and vice versa.
* For real numbers (and for any totally ordered set), the maximum operation is distributive over the minimum operation, and vice versa:
\[
\max(a, \min(b, c)) = \min(\max(a, b), \max(a, c)) \quad \text{ and } \quad \min(a, \max(b, c)) = \max(\min(a, b), \min(a, c)).
\]

* For integers, the greatest common divisor is distributive over the least common multiple, and vice versa:
\[
\gcd(a, \operatorname{lcm}(b, c)) = \operatorname{lcm}(\gcd(a, b), \gcd(a, c)) \quad \text{ and } \quad \operatorname{lcm}(a, \gcd(b, c)) = \gcd(\operatorname{lcm}(a, b), \operatorname{lcm}(a, c)).
\]

* For real numbers, addition distributes over the maximum operation, and also over the minimum operation:
\[
a + \max(b, c) = \max(a + b, a + c) \quad \text{ and } \quad a + \min(b, c) = \min(a + b, a + c).
\]

* For binomial multiplication, distribution is sometimes referred to as the FOIL Method (First terms \(a c,\) Outer \(a d,\) Inner \(b c,\) and Last \(b d\)) such as: \((a + b) \cdot (c + d) = a c + a d + b c + b d.\)
* In all semirings, including the complex numbers, the quaternions, polynomials, and matrices,  multiplication distributes over addition: \(u (v + w) = u v + u w, (u + v)w = u w + v w.\)
* In all algebras over a field, including the octonions and other non-associative algebras, multiplication distributes over addition.

## Propositional logic


### Rule of replacement
In standard truth-functional propositional logic,  in logical proofs uses two valid rules of replacement to expand individual occurrences of certain logical connectives, within some formula, into separate applications of those connectives across subformulas of the given formula. The rules are

\[
(P \land (Q \lor R)) \Leftrightarrow ((P \land Q) \lor (P \land R)) \qquad \text{ and } \qquad (P \lor (Q \land R)) \Leftrightarrow ((P \lor Q) \land (P \lor R))
\]

where "\(\Leftrightarrow\)", also written \(\,\equiv,\,\) is a metalogical symbol representing "can be replaced in a proof with" or "is logically equivalent to".

### Truth functional connectives
 is a property of some logical connectives of truth-functional propositional logic. The following logical equivalences demonstrate that distributivity is a property of particular connectives. The following are truth-functional tautologies.

\[
\begin{alignat}{13}
&(P &&\;\land &&(Q \lor R)) &&\;\Leftrightarrow\;&& ((P \land Q) &&\;\lor (P \land R)) && \quad\text{ Distribution of } && \text{ conjunction } && \text{ over } && \text{ disjunction } \\
&(P &&\;\lor &&(Q \land R)) &&\;\Leftrightarrow\;&& ((P \lor Q) &&\;\land (P \lor R)) && \quad\text{ Distribution of } && \text{ disjunction } && \text{ over } && \text{ conjunction } \\
&(P &&\;\land &&(Q \land R)) &&\;\Leftrightarrow\;&& ((P \land Q) &&\;\land (P \land R)) && \quad\text{ Distribution of } && \text{ conjunction } && \text{ over } && \text{ conjunction } \\
&(P &&\;\lor &&(Q \lor R)) &&\;\Leftrightarrow\;&& ((P \lor Q) &&\;\lor (P \lor R)) && \quad\text{ Distribution of } && \text{ disjunction } && \text{ over } && \text{ disjunction } \\
&(P &&\to &&(Q \to R)) &&\;\Leftrightarrow\;&& ((P \to Q) &&\to (P \to R)) && \quad\text{ Distribution of } && \text{ implication } && \text{  } && \text{  } \\
&(P &&\to &&(Q \leftrightarrow R)) &&\;\Leftrightarrow\;&& ((P \to Q) &&\leftrightarrow (P \to R)) && \quad\text{ Distribution of } && \text{ implication } && \text{ over } && \text{ equivalence } \\
&(P &&\to &&(Q \land R)) &&\;\Leftrightarrow\;&& ((P \to Q) &&\;\land (P \to R)) && \quad\text{ Distribution of } && \text{ implication } && \text{ over } && \text{ conjunction } \\
&(P &&\;\lor &&(Q \leftrightarrow R)) &&\;\Leftrightarrow\;&& ((P \lor Q) &&\leftrightarrow (P \lor R)) && \quad\text{ Distribution of } && \text{ disjunction } && \text{ over } && \text{ equivalence } \\
\end{alignat}
\]


;Double distribution:

\[
\begin{alignat}{13}
&((P \land Q) &&\;\lor (R \land S)) &&\;\Leftrightarrow\;&& (((P \lor R) \land (P \lor S)) &&\;\land ((Q \lor R) \land (Q \lor S))) && \\
&((P \lor Q) &&\;\land (R \lor S)) &&\;\Leftrightarrow\;&& (((P \land R) \lor (P \land S)) &&\;\lor ((Q \land R) \lor (Q \land S))) && \\
\end{alignat}
\]


## Distributivity and rounding
In approximate arithmetic, such as floating-point arithmetic, the distributive property of multiplication (and division) over addition may fail because of the limitations of arithmetic precision.  For example, the identity \(1/3 + 1/3 + 1/3 = (1 + 1 + 1) / 3\) fails in decimal arithmetic, regardless of the number of significant digits. Methods such as banker's rounding may help in some cases, as may increasing the precision used, but ultimately some calculation errors are inevitable.

## In rings and other structures
Distributivity is most commonly found in semirings, notably the particular cases of rings and distributive lattices.

A semiring has two binary operations, commonly denoted \(\,+\,\) and \(\,*,\) and requires that \(\,*\,\) must distribute over \(\,+.\)

A ring is a semiring with additive inverses.

A lattice is another kind of algebraic structure with two binary operations, \(\,\land \text{ and } \lor.\)
If either of these operations distributes over the other (say \(\,\land\,\) distributes over \(\,\lor\)), then the reverse also holds (\(\,\lor\,\) distributes over \(\,\land\,\)), and the lattice is called distributive. See also .

A Boolean algebra can be interpreted either as a special kind of ring (a Boolean ring) or a special kind of distributive lattice (a Boolean lattice). Each interpretation is responsible for different distributive laws in the Boolean algebra.

In any semiring, distributivity can be used to show that any product of sums is a sum of products (though not every sum of products is necessarily a product of sums). The general formula reads:
\[
\prod_{i=1}^m \left(\sum_{j=1}^{n_i}  a_{i,j}\right) = \sum_{j_1=1}^{n_1} \sum_{j_2=1}^{n_2} \cdots \sum_{j_m=1}^{n_m} \;\prod_{i=1}^m a_{i,j_i}
\]
Structures without two-sided distributive laws are near-rings and near-fields. The operations are usually defined to be distributive on the right but not on the left.

## Generalizations


In several mathematical areas, generalized distributivity laws are considered. This may involve the weakening of the above conditions or the extension to infinitary operations. Especially in order theory one finds numerous important variants of distributivity, some of which include infinitary operations, such as the infinite distributive law; others being defined in the presence of only  binary operation, such as the according definitions and their relations are given in the article distributivity (order theory). This also includes the notion of a completely distributive lattice.

In the presence of an ordering relation, one can also weaken the above equalities by replacing \(\,=\,\) by either \(\,\leq\,\) or \(\,\geq.\) Naturally, this will lead to meaningful concepts only in some situations. An application of this principle are the notions of *sub-distributivity*, where equality is replaced by "less than or equal"; and *super-distributivity*, where equality is replaced by "greater than or equal".

In category theory, if \((S, \mu, \nu)\) and \(\left(S^{\prime}, \mu^{\prime}, \nu^{\prime}\right)\) are monads on a category \(C,\) a *distributive law* \(S . S^{\prime} \to S^{\prime} . S\) is a natural transformation \(\lambda : S . S^{\prime} \to S^{\prime} . S\) such that \(\left(S^{\prime}, \lambda\right)\) is a lax map of monads \(S \to S\) and \((S, \lambda)\) is a colax map of monads \(S^{\prime} \to S^{\prime}.\) This is exactly the data needed to define a monad structure on \(S^{\prime} . S\): the multiplication map is \(S^{\prime} \mu . \mu^{\prime} S^2 . S^{\prime} \lambda S\) and the unit map is \(\eta^{\prime} S . \eta.\)
A generalized distributive law  has also been proposed in the area of information theory.

### Antidistributivity
The ubiquitous identity that relates inverses to the binary operation in any group, namely \((x y)^{-1} = y^{-1} x^{-1},\) which is taken as an axiom in the more general context of a semigroup with involution, has sometimes been called an *antidistributive property* (of inversion as a unary operation).

In the context of a near-ring, which removes the commutativity of the additively written group and assumes only one-sided distributivity, one can speak of (two-sided) *distributive elements* but also of *antidistributive elements*. The latter reverse the order of (the non-commutative) addition; assuming a left-nearring (i.e. one which all elements distribute when multiplied on the left), then an antidistributive element \(a\) reverses the order of addition when multiplied to the right: \((x + y) a = y a + x a.\)

In the study of propositional logic and Boolean algebra, the term *antidistributive law* is sometimes used to denote the interchange between conjunction and disjunction when implication factors over them:

\[
(a \lor b) \Rightarrow c \equiv (a \Rightarrow c) \land (b \Rightarrow c)
\]


\[
(a \land b) \Rightarrow c \equiv (a \Rightarrow c) \lor (b \Rightarrow c).
\]


These two tautologies are a direct consequence of the duality in De Morgan's laws.

## Notes


## External links

* A demonstration of the Distributive Law for integer arithmetic (from cut-the-knot)

