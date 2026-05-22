> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Inequality_%28mathematics%29) — CC BY-SA 4.0

# Compound inequalities

{about|relations "greater than" and "less than"|the relation "not equal"|Inequation}
{redirect-distinguish|Less than|Less Than (song)}
{redirect-distinguish|≪|Absolute continuity of measures}


In mathematics, an **inequality** is a relation which makes a non-equal comparison between two numbers or other mathematical expressions. It is used most often to compare two numbers on the number line by their size. The main types of inequality are **less than** and **greater than** (denoted by < and >, respectively the less-than and greater-than signs).

## Notation
There are several different notations used to represent different kinds of inequalities:
* The notation *a* < *b* means that *a* is **less than** *b*.
* The notation *a* > *b* means that *a* is **greater than** *b*.
In either case, *a* is not equal to *b*. These relations are known as **strict inequalities**, meaning that *a* is strictly less than or strictly greater than *b*. Equality is excluded.

In contrast to strict inequalities, there are two types of inequality relations that are not strict:
* The notation *a* ≤ *b* or *a* ⩽ *b* or *a* ≦ *b* means that *a* is **less than or equal to** *b* (or, equivalently, at most *b*).
* The notation *a* ≥ *b* or *a* ⩾ *b* or *a* ≧ *b* means that *a* is **greater than or equal to** *b* (or, equivalently, at least *b*).

In the 17th and 18th centuries, personal notations or typewriting signs were used to signal inequalities. For example, In 1670, John Wallis used a single horizontal bar *above* rather than below the < and >.
Later in 1734, ≦ and ≧, known as "less than (greater-than) over equal to" or "less than (greater than) or equal to with double horizontal bars", first appeared in Pierre Bouguer's work . After that, mathematicians simplified Bouguer's symbol to "less than (greater than) or equal to with one horizontal bar" (≤), or "less than (greater than) or slanted equal to" (⩽).

The relation **not greater than** can also be represented by \(a \ngtr b,\) the symbol for "greater than" bisected by a slash, "not". The same is true for **not less than**, \(a \nless b.\)

The notation *a* ≠ *b* means that *a* is not equal to *b*; this *inequation* sometimes is considered a form of strict inequality. It does not say that one is greater than the other; it does not even require *a* and *b* to be a member of an ordered set.

In engineering sciences, less formal use of the notation is to state that one quantity is "much greater" than another, normally by several orders of magnitude.
* The notation *a* ≪ *b* means that *a* is **much less than** *b*.
* The notation *a* ≫ *b* means that *a* is **much greater than** *b*.
This implies that the lesser value can be neglected with little effect on the accuracy of an approximation (such as the case of ultrarelativistic limit in physics).

In all of the cases above, any two symbols mirroring each other are symmetrical; *a* < *b* and *b* > *a* are equivalent, etc.

## Properties on the number line
Inequalities are governed by the following properties. All of these properties also hold if all of the non-strict inequalities (≤ and ≥) are replaced by their corresponding strict inequalities (< and >) and — in the case of applying a function — monotonic functions are limited to *strictly* monotonic functions.

### Converse
The relations ≤ and ≥ are each other's converse, meaning that for any real numbers *a* and *b*:


### Transitivity
The transitive property of inequality states that for any real numbers *a*, *b*, *c*:

If *either* of the premises is a strict inequality, then the conclusion is a strict inequality:


### Addition and subtraction

A common constant *c* may be added to or subtracted from both sides of an inequality. So, for any real numbers *a*, *b*, *c*:


In other words, the inequality relation is preserved under addition (or subtraction) and the real numbers are an ordered group under addition.

### Multiplication and division


The properties that deal with multiplication and division state that for any real numbers, *a*, *b* and non-zero *c*:


In other words, the inequality relation is preserved under multiplication and division with positive constant, but is reversed when a negative constant is involved. More generally, this applies for an ordered field. For more information, see *Ordered fields*.

### Additive inverse
The property for the additive inverse states that for any real numbers *a* and *b*:


### Multiplicative inverse
If both numbers are positive, then the inequality relation between the multiplicative inverses is opposite of that between the original numbers. More specifically, for any non-zero real numbers *a* and *b* that are both positive (or both negative):


All of the cases for the signs of *a* and *b* can also be written in chained notation, as follows:


### Applying a function to both sides

Any monotonically increasing function, by its definition, may be applied to both sides of an inequality without breaking the inequality relation (provided that both expressions are in the domain of that function). However, applying a monotonically decreasing function to both sides of an inequality means the inequality relation would be reversed. The rules for the additive inverse, and the multiplicative inverse for positive numbers, are both examples of applying a monotonically decreasing function.

If the inequality is strict (*a* < *b*, *a* > *b*) *and* the function is strictly monotonic, then the inequality remains strict. If only one of these conditions is strict, then the resultant inequality is non-strict. In fact, the rules for additive and multiplicative inverses are both examples of applying a *strictly* monotonically decreasing function.

A few examples of this rule are:
* Raising both sides of an inequality to a power *n* > 0 (equiv., −*n* < 0), when *a* and *b* are positive real numbers:
* Taking the natural logarithm on both sides of an inequality, when *a* and *b* are positive real numbers:  0 < *a* ≤ *b* ⇔ ln(*a*) ≤ ln(*b*).  0 < *a* < *b* ⇔ ln(*a*) < ln(*b*).  (this is true because the natural logarithm is a strictly increasing function.)

## Formal definitions and generalizations
A (non-strict) **partial order** is a binary relation ≤ over a set *P* which is reflexive, antisymmetric, and transitive. That is, for all *a*, *b*, and *c* in *P*, it must satisfy the three following clauses:

* *a* ≤ *a* (reflexivity)
* if *a* ≤ *b* and *b* ≤ *a*, then *a* = *b* (antisymmetry)
* if *a* ≤ *b* and *b* ≤ *c*, then *a* ≤ *c* (transitivity)

A set with a partial order is called a **partially ordered set**. Those are the very basic axioms that every kind of order has to satisfy.

A strict partial order is a relation < that satisfies
* *a* ≮ *a* (irreflexivity),
* if *a* < *b*, then *b* ≮ *a* (asymmetry),
* if *a* < *b* and *b* < *c*, then *a* < *c* (transitivity),
where ≮ means that < does not hold.

Some types of partial orders are specified by adding further axioms, such as:

* Total order: For every *a* and *b* in *P*, *a* ≤ *b* or *b* ≤ *a* .
* Dense order: For all *a* and *b* in *P* for which *a* < *b*, there is a *c* in *P* such that *a* < *c* < *b*.
* Least-upper-bound property: Every non-empty subset of *P* with an upper bound has a *least* upper bound (supremum) in *P*.

### Ordered fields

If (*F*, +, ×) is a field and ≤ is a total order on *F*, then (*F*, +, ×, ≤) is called an **ordered field** if and only if:
* *a* ≤ *b* implies *a* + *c* ≤ *b* + *c*;
* 0 ≤ *a* and 0 ≤ *b* implies 0 ≤ *a* × *b*.

Both \((\mathbb Q, +, \times, \leq)\) and \((\mathbb R, +, \times, \leq)\) are ordered fields, but ≤ cannot be defined in order to make \((\mathbb C, +, \times, \leq)\) an ordered field, because −1 is the square of *i* and would therefore be positive.

Besides being an ordered field, **R** also has the Least-upper-bound property. In fact, **R** can be defined as the only ordered field with that quality.

## Chained notation
The notation ***a* < *b* < *c*** stands for "*a* < *b* and *b* < *c*", from which, by the transitivity property above, it also follows that *a* < *c*. By the above laws, one can add or subtract the same number to all three terms, or multiply or divide all three terms by same nonzero number and reverse all inequalities if that number is negative. Hence, for example, *a* < *b* + *e* < *c* is equivalent to *a* − *e* < *b* < *c* − *e*.

This notation can be generalized to any number of terms: for instance, ***a*\(_{1}\) ≤ *a*\(_{2}\) ≤ ... ≤ *a*\(_{*n*}\)** means that *a*\(_{*i*}\) ≤ *a*\(_{*i*+1}\) for *i* = 1, 2, ..., *n* − 1. By transitivity, this condition is equivalent to *a*\(_{*i*}\) ≤ *a*\(_{*j*}\) for any 1 ≤ *i* ≤ *j* ≤ *n*.

When solving inequalities using chained notation, it is possible and sometimes necessary to evaluate the terms independently. For instance, to solve the inequality 4*x* < 2*x* + 1 ≤ 3*x* + 2, it is not possible to isolate *x* in any one part of the inequality through addition or subtraction. Instead, the inequalities must be solved independently, yielding *x* < 1/2 and *x* ≥ −1 respectively, which can be combined into the final solution −1 ≤ *x* < 1/2.

Occasionally, chained notation is used with inequalities in different directions, in which case the meaning is the logical conjunction of the inequalities between adjacent terms. For example, the defining condition of a zigzag poset is written as *a*\(_{1}\) < *a*\(_{2}\) > *a*\(_{3}\) < *a*\(_{4}\) > *a*\(_{5}\) < *a*\(_{6}\) > ... . Mixed chained notation is used more often with compatible relations, like <, =, ≤. For instance, *a* < *b* = *c* ≤ *d* means that *a* < *b*, *b* = *c*, and *c* ≤ *d*. This notation exists in a few programming languages such as Python. In contrast, in programming languages that provide an ordering on the type of comparison results, such as C, even homogeneous chains may have a completely different meaning.

## Sharp inequalities
An inequality is said to be *sharp* if it cannot be *relaxed* and still be valid in general. Formally, a universally quantified inequality *φ* is called sharp if, for every valid universally quantified inequality *ψ*, if *ψ* material conditional holds, then *ψ* [[equivalence (logic) also holds. For instance, the inequality [[universal quantification is sharp, whereas the inequality ∀*a* ∈ **R**. *a*\(^{2}\) ≥ −1 is not sharp.

## Inequalities between means


There are many inequalities between means. For example, for any positive numbers *a*\(_{1}\), *a*\(_{2}\), ..., *a*\(_{*n*}\) we have

\(H\le G\le A\le Q,\)

where they represent the following means of the sequence:

* [[Harmonic mean : \(H = \frac{n}{\frac{1}{a_1} + \frac{1}{a_2} + \cdots + \frac{1}{a_n\)
* Geometric mean : \(G = \sqrt[n]{a_1 \cdot a_2 \cdots a_n}\)
* Arithmetic mean : \(A = \frac{a_1 + a_2 + \cdots + a_n}{n}\)
* Quadratic mean : \(Q = \sqrt{\frac{a_1^2 + a_2^2 + \cdots + a_n^2}{n\)

## Cauchy–Schwarz inequality


The Cauchy–Schwarz inequality states that for all vectors *u* and *v* of an inner product space it is true that

\[
|\langle \mathbf{u},\mathbf{v}\rangle| ^2 \leq \langle \mathbf{u},\mathbf{u}\rangle \cdot \langle \mathbf{v},\mathbf{v}\rangle,
\]

where \(\langle\cdot,\cdot\rangle\) is the inner product. Examples of inner products include the real and complex dot product; In Euclidean space *R*\(^{*n*}\) with the standard inner product, the Cauchy–Schwarz inequality is

\[
\biggl(\sum_{i=1}^n u_i v_i\biggr)^2\leq \biggl(\sum_{i=1}^n u_i^2\biggr) \biggl(\sum_{i=1}^n v_i^2\biggr).
\]


## Power inequalities
A **power inequality** is an inequality containing terms of the form *a*\(^{*b*}\), where *a* and *b* are real positive numbers or variable expressions. They often appear in mathematical olympiads exercises.

Examples:
* For any real *x*,
\[
e^x \ge 1+x.
\]

* If *x* > 0 and *p* > 0, then
\[
\frac1p\left(x^p - 1\right) \ge \ln(x) \ge \frac1p\left(1 - \frac{1}{x^p}\right).
\]
 In the limit of *p* → 0, the upper and lower bounds converge to ln(*x*).
* If *x* > 0, then
\[
x^x \ge \left( \frac{1}{e}\right)^\frac{1}{e}.
\]

* If *x* > 0, then
\[
x^{x^x} \ge x.
\]

* If *x*, *y*, *z* > 0, then
\[
\left(x+y\right)^z + \left(x+z\right)^y + \left(y+z\right)^x > 2.
\]

* For any real distinct numbers *a* and *b*,
\[
\frac{e^b-e^a}{b-a} > e^\frac{a+b}{2}.
\]

* If *x*, *y* > 0 and 0 < *p* < 1, then
\[
x^p+y^p > \left(x+y\right)^p.
\]

* If *x*, *y*, *z* > 0, then
\[
x^x y^y z^z \ge \left(xyz\right)^\frac{x+y+z}{3}.
\]

* If *a*, *b* > 0, then
\[
a^a + b^b \ge a^b + b^a.
\]

* If *a*, *b* > 0, then
\[
a^{ea} + b^{eb} \ge a^{eb} + b^{ea}.
\]

* If *a*, *b*, *c* > 0, then
\[
a^{2a} + b^{2b} + c^{2c} \ge a^{2b} + b^{2c} + c^{2a}.
\]

* If *a*, *b* > 0, then
\[
a^b + b^a > 1.
\]


## Well-known inequalities
List of inequalities

Mathematicians often use inequalities to bound quantities for which exact formulas cannot be computed easily. Some inequalities are used so often that they have names:

* Azuma's inequality
* Bernoulli's inequality
* Bell's inequality
* Boole's inequality
* Cauchy–Schwarz inequality
* Chebyshev's inequality
* Chernoff's inequality
* Cramér–Rao inequality
* Hoeffding's inequality
* Hölder's inequality
* Inequality of arithmetic and geometric means
* Jensen's inequality
* Kolmogorov's inequality
* Markov's inequality
* Minkowski inequality
* Nesbitt's inequality
* Pedoe's inequality
* Poincaré inequality
* Samuelson's inequality
* Sobolev inequality
* Triangle inequality


## Complex numbers and inequalities
The set of complex numbers \(\mathbb{C}\) with its operations of addition and multiplication is a field, but it is impossible to define any relation ≤ so that \((\Complex, +, \times, \leq)\) becomes an ordered field. To make \((\mathbb{C}, +, \times, \leq)\) an ordered field, it would have to satisfy the following two properties:
* if *a* ≤ *b*, then *a* + *c* ≤ *b* + *c*;
* if 0 ≤ *a* and 0 ≤ *b*, then 0 ≤ *ab*.

Because ≤ is a total order, for any number *a*, either 0 ≤ *a* or *a* ≤ 0 (in which case the first property above implies that 0 ≤ −*a*). In either case 0 ≤ *a*\(^{2}\); this means that *i*\(^{2}\) > 0 and 1\(^{2}\) > 0; so −1 > 0 and 1 > 0, which means (−1 + 1) > 0; contradiction.

However, an operation ≤ can be defined so as to satisfy only the first property (namely, "if *a* ≤ *b*, then *a* + *c* ≤ *b* + *c*"). Sometimes the lexicographical order definition is used:
* *a* ≤ *b*, if
** Re(*a*) < Re(*b*), or
** Re(*a*) {=} Re(*b*) and Im(*a*) ≤ Im(*b*)
It can easily be proven that for this definition *a* ≤ *b* implies *a* + *c* ≤ *b* + *c*.

## Systems of inequalities
Systems of linear inequalities can be simplified by Fourier–Motzkin elimination.

The cylindrical algebraic decomposition is an algorithm that allows testing whether a system of polynomial equations and inequalities has solutions, and, if solutions exist, describing them. The complexity of this algorithm is doubly exponential in the number of variables. It is an active research domain to design algorithms that are more efficient in specific cases.

