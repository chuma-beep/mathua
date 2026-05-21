> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Quantifier_%28logic%29) — CC BY-SA 4.0

# Quantifiers: universal and existential

In logic, a **quantifier** is an operator that specifies how many individuals in the domain of discourse satisfy an open formula. For instance, the universal quantifier \(\forall\) in the first-order formula \(\forall x P(x)\) expresses that everything in the domain satisfies the property denoted by \(P\). On the other hand, the existential quantifier \(\exists\) in the formula \(\exists x P(x)\) expresses that there exists something in the domain which satisfies that property. A formula where a quantifier takes widest scope is called a quantified formula. A quantified formula must contain a bound variable and a subformula specifying a property of the referent of that variable.

The most commonly used quantifiers are \(\forall\) and \(\exists\). These quantifiers are standardly defined as duals; in classical logic: each can be defined in terms of the other using negation. They can also be used to define more complex quantifiers, as in the formula \(\neg \exists x P(x)\) which expresses that nothing has the property \(P\). Other quantifiers are only definable within second-order logic or higher-order logics. Quantifiers have been generalized beginning with the work of Andrzej Mostowski and Per Lindström.

In a first-order logic statement, quantifications in the same type (either universal quantifications or existential quantifications) can be exchanged without changing the meaning of the statement, while the exchange of quantifications in different types changes the meaning. As an example, the only difference in the definition of uniform continuity and (ordinary) continuity is the order of quantifications.

First-order quantifiers approximate the meanings of some natural language quantifiers such as "some" and "all". However, many natural language quantifiers can only be analyzed in terms of generalized quantifiers.

## Relations to logical conjunction and disjunction
For a finite domain of discourse \(D = \{a_1,...a_n\}\), the universally quantified formula \(\forall x \in D \; P(x)\) is equivalent to the logical conjunction \(P(a_1) \land.. . \land P(a_n)\).
Dually, the existentially quantified formula \(\exists x \in D \; P(x)\) is equivalent to the logical disjunction \(P(a_1) \lor.. . \lor P(a_n)\).
For example, if \(B = \{ 0,1 \}\) is the set of binary digits, the formula \(\forall x \in B \; x = x^2\) abbreviates \(0 = 0^2 \land 1 = 1^2\), which evaluates to *true*.

### Infinite domain of discourse
Consider the following statement (using dot notation for multiplication):

This has the appearance of an *infinite conjunction* of propositions. From the point of view of formal languages, this is immediately a problem; since syntax rules are expected to generate finite statements. A succinct equivalent formulation, which avoids these problems, uses *universal quantification*:

A similar analysis applies to the disjunction,

which can be rephrased using *existential quantification*:

## Algebraic approaches to quantification
It is possible to devise abstract algebras whose models include formal languages with quantification, but progress has been slow and interest in such algebra has been limited. Three approaches have been devised to date:
* Relation algebra, invented by Augustus De Morgan, and developed by Charles Sanders Peirce, Ernst Schröder, Alfred Tarski, and Tarski's students. Relation algebra cannot represent any formula with quantifiers nested more than three deep. Surprisingly, the models of relation algebra include the axiomatic set theory ZFC and Peano arithmetic;
* Cylindric algebra, devised by Alfred Tarski, Leon Henkin, and others;
* The polyadic algebra of Paul Halmos.

## Notation
The two most common quantifiers are the universal quantifier and the existential quantifier. The traditional symbol for the universal quantifier is "∀", a rotated letter "A", which stands for "for all" or "all". The corresponding symbol for the existential quantifier is "∃", a rotated letter "E", which stands for "there exists" or "exists".

An example of translating a quantified statement in a natural language such as English would be as follows. Given the statement, "Each of Peter's friends either likes to dance or likes to go to the beach (or both)", key aspects can be identified and rewritten using symbols including quantifiers. So, let *X* be the set of all Peter's friends, *P*(*x*) the predicate "*x* likes to dance", and *Q*(*x*) the predicate "*x* likes to go to the beach". Then the above sentence can be written in formal notation as \(\forall{x}{\in}X, (P(x) \lor Q(x))\), which is read, "for every *x* that is a member of *X*, *P* applies to *x* or *Q* applies to *x*".

Some other quantified expressions are constructed as follows,

*\(\exists{x}\, P\)
*\(\forall{x}\, P\)

for a formula *P*. These two expressions (using the definitions above) are read as "there exists a friend of Peter who likes to dance" and "all friends of Peter like to dance", respectively. Variant notations include, for set *X* and set members *x*:

*\(\bigvee_{x} P\)
*\((\exists{x}) P\)
*\((\exists x \. \ P)\)
*\(\exists x \ \cdot \ P\)
*\((\exists x : P)\)
*\(\exists{x}(P)\)
*\(\exists_{x}\, P\)
*\(\exists{x}{,}\, P\)
*\(\exists{x}{\in}X \, P\)
*\(\exists\, x{:}X \, P\)

All of these variations also apply to universal quantification.
Other variations for the universal quantifier are

*\(\bigwedge_{x} P\)
*\(\bigwedge x P\)
*\((x) \, P\)

Some versions of the notation explicitly mention the range of quantification. The range of quantification must always be specified; for a given mathematical theory, this can be done in several ways:
* Assume a fixed domain of discourse for every quantification, as is done in Zermelo–Fraenkel set theory.
* Fix several domains of discourse in advance and require that each variable have a declared domain, which is the *type* of that variable. This is analogous to the situation in statically typed computer programming languages, where variables have declared types.
* Mention explicitly the range of quantification, perhaps using a symbol for the set of all objects in that domain (or the type of the objects in that domain).

One can use any variable as a quantified variable in place of any other, under certain restrictions in which *variable capture* does not occur. Even if the notation uses typed variables, variables of that type may be used.

Informally or in natural language, the "∀*x*" or "∃*x*" might appear after or in the middle of *P*(*x*). Formally, however, the phrase that introduces the dummy variable is placed in front.

Mathematical formulas mix symbolic expressions for quantifiers with natural language quantifiers such as,

Keywords for uniqueness quantification include:

Further, *x* may be replaced by a pronoun. For example,

## Order of quantifiers (nesting)

The order of quantifiers is critical to meaning, as is illustrated by the following two propositions:

This is clearly true; it just asserts that every natural number has a square. The meaning of the assertion in which the order of quantifiers is reversed is different:

This is clearly false; it asserts that there is a single natural number *s* that is the square of *every* natural number. This is because the syntax directs that any variable cannot be a function of subsequently introduced variables.

A less trivial example from mathematical analysis regards the concepts of uniform and pointwise continuity, whose definitions differ only by an exchange in the positions of two quantifiers. A function *f* from **R** to **R** is called

* Pointwise continuous if
\[
\forall \varepsilon > 0 \; \forall x \in \R \; \exists \delta > 0 \; \forall h \in \R \; (|h| < \delta \, \Rightarrow \, |f(x) - f(x + h)| < \varepsilon )
\]

* Uniformly continuous if
\[
\forall \varepsilon > 0 \; \exists \delta > 0 \; \forall x \in \R \; \forall h \in \R \; (|h| < \delta \, \Rightarrow \, |f(x) - f(x + h)| < \varepsilon )
\]

In the former case, the particular value chosen for *δ* can be a function of both *ε* and *x*, the variables that precede it.
In the latter case, *δ* can be a function only of *ε* (i.e., it has to be chosen independent of *x*). For example, *f*(*x*) = *x*\(^{2}\) satisfies pointwise, but not uniform continuity (its slope is unbounded). In contrast, interchanging the two initial universal quantifiers in the definition of pointwise continuity does not change the meaning.

As a general rule, swapping two adjacent universal quantifiers with the same scope (or swapping two adjacent existential quantifiers with the same scope) doesn't change the meaning of the formula (see Example here), but swapping an existential quantifier and an adjacent universal quantifier may change its meaning.

The maximum depth of nesting of quantifiers in a formula is called its "quantifier rank".

## Equivalent expressions
If *D* is a domain of *x* and *P*(*x*) is a predicate dependent on object variable *x*, then the universal proposition can be expressed as

\[
\forall x\!\in\!D\; P(x).
\]

This notation is known as restricted or relativized or bounded quantification. Equivalently one can write,

\[
\forall x\;(x\!\in\!D \to P(x)).
\]

The existential proposition can be expressed with bounded quantification as

\[
\exists x\!\in\!D\; P(x),
\]

or equivalently

\[
\exists x\;(x\!\in\!\!D \land P(x)).
\]

Together with negation, only one of either the universal or existential quantifier is needed to perform both tasks:

\[
\neg (\forall x\!\in\!D\; P(x)) \equiv \exists x\!\in\!D\; \neg P(x),
\]

which shows that to disprove a "for all *x*" proposition, one needs no more than to find an *x* for which the predicate is false. Similarly,

\[
\neg (\exists x\!\in\!D\; P(x)) \equiv \forall x\!\in\!D\; \neg P(x),
\]

to disprove a "there exists an *x*" proposition, one needs to show that the predicate is false for all *x*.

In classical logic, every formula is logically equivalent to a formula in prenex normal form, that is, a string of quantifiers and bound variables followed by a quantifier-free formula.

## Range of quantification
Every quantification involves one specific variable and a domain of discourse or **range of quantification** of that variable. The range of quantification specifies the set of values that the variable takes. In the examples above, the range of quantification is the set of natural numbers. Specification of the range of quantification allows us to express the difference between, say, asserting that a predicate holds for some natural number or for some real number. Expository conventions often reserve some variable names such as "*n*" for natural numbers, and "*x*" for real numbers, although relying exclusively on naming conventions cannot work in general; since ranges of variables can change in the course of a mathematical argument.

A universally quantified formula over an empty range (like \(\forall x\!\in\!\varnothing\; x \neq x\)) is always vacuously true. Conversely, an existentially quantified formula over an empty range (like \(\exists x\!\in\!\varnothing\; x = x\)) is always false.

A more natural way to restrict the domain of discourse uses *guarded quantification*. For example, the guarded quantification

means

In some mathematical theories, a single domain of discourse fixed in advance is assumed. For example, in Zermelo–Fraenkel set theory, variables range over all sets. In this case, guarded quantifiers can be used to mimic a smaller range of quantification. Thus in the example above, to express

in Zermelo–Fraenkel set theory, one would write

where **N** is the set of all natural numbers.

## Formal semantics
Mathematical semantics is the application of mathematics to study the meaning of expressions in a formal language. It has three elements: a mathematical specification of a class of objects via syntax, a mathematical specification of various semantic domains and the relation between the two, which is usually expressed as a function from syntactic objects to semantic ones. This article only addresses the issue of how quantifier elements are interpreted.
The syntax of a formula can be given by a syntax tree. A quantifier has a scope, and an occurrence of a variable *x* is free if it is not within the scope of a quantification for that variable. Thus in

\[
\forall x (\exists y B(x,y)) \vee C(y,x)
\]

the occurrence of both *x* and *y* in *C*(*y*, *x*) is free, while the occurrence of *x* and *y* in *B*(*y*, *x*) is bound (i.e. non-free).

An interpretation for first-order predicate calculus assumes as given a domain of individuals *X*. A formula *A* whose free variables are *x*\(_{1}\),.. ., *x*\(_{n}\) is interpreted as a Boolean-valued function *F*(*v*\(_{1}\),.. ., *v*\(_{*n*}\)) of *n* arguments, where each argument ranges over the domain *X*. Boolean-valued means that the function assumes one of the values **T** (interpreted as truth) or **F** (interpreted as falsehood). The interpretation of the formula

\[
\forall x_n A(x_1, \ldots , x_n)
\]

is the function *G* of *n*-1 arguments such that *G*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\)) = **T** if and only if *F*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\), *w*) = **T** for every *w* in *X*. If *F*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\), *w*) = **F** for at least one value of *w*, then *G*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\)) = **F**. Similarly the interpretation of the formula

\[
\exists x_n A(x_1, \ldots , x_n)
\]

is the function *H* of *n*-1 arguments such that *H*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\)) = **T** if and only if *F*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\), *w*) = **T** for at least one *w* and *H*(*v*\(_{1}\),.. ., *v*\(_{*n*-1}\)) = **F** otherwise.

The semantics for uniqueness quantification requires first-order predicate calculus with equality. This means there is given a distinguished two-placed predicate "="; the semantics is also modified accordingly so that "=" is always interpreted as the two-place equality relation on *X*. The interpretation of

\[
\exists ! x_n A(x_1, \ldots , x_n)
\]

then is the function of *n*-1 arguments, which is the logical *and* of the interpretations of

\[
\begin{align}
\exists x_n & A(x_1, \ldots , x_n) \\
\forall y,z & \big( A(x_1, \ldots ,x_{n-1}, y) \wedge A(x_1, \ldots ,x_{n-1}, z) \implies y = z \big).
\end{align}
\]

Each kind of quantification defines a corresponding closure operator on the set of formulas, by adding, for each free variable *x*, a quantifier to bind *x*. For example, the *existential closure* of the open formula *n*>2 ∧ *x*\(^{*n*}\)+*y*\(^{*n*}\)=*z*\(^{*n*}\) is the closed formula ∃*n* ∃*x* ∃*y* ∃*z* (*n*>2 ∧ *x*\(^{*n*}\)+*y*\(^{*n*}\)=*z*\(^{*n*}\)); the latter formula, when interpreted over the positive integers, is known to be false by Fermat's Last Theorem. As another example, equational axioms, like *x*+*y*=*y*+*x*, are usually meant to denote their *universal closure*, like ∀*x* ∀*y* (*x*+*y*=*y*+*x*) to express commutativity.

## Paucal, multal and other degree quantifiers

None of the quantifiers previously discussed apply to a quantification such as

One possible interpretation mechanism can be obtained as follows: Suppose that in addition to a semantic domain *X*, we have given a probability measure P defined on *X* and cutoff numbers 0 < *a* ≤ *b* ≤ 1. If *A* is a formula with free variables *x*\(_{1}\),...,*x*\(_{*n*}\) whose interpretation is
the function *F* of variables *v*\(_{1}\),...,*v*\(_{*n*}\) then the interpretation of

\[
\exists^{\mathrm{manyx_n A(x_1, \ldots, x_{n-1}, x_n)
\]

is the function of *v*\(_{1}\),...,*v*\(_{*n*-1}\) which is **T** if and only if

\[
\operatorname{P} \{w: F(v_1, \ldots, v_{n-1}, w) = \mathbf{T} \} \geq b
\]

and **F** otherwise. Similarly, the interpretation of

\[
\exists^{\mathrm{few x_n A(x_1, \ldots, x_{n-1}, x_n)
\]

is the function of *v*\(_{1}\),...,*v*\(_{*n*-1}\) which is **F** if and only if

\[
0< \operatorname{P} \{w: F(v_1, \ldots, v_{n-1}, w) = \mathbf{T}\} \leq a
\]

and **T** otherwise.

## Other quantifiers
A few other quantifiers have been proposed over time. In particular, the solution quantifier, noted § (section sign) and read "those". For example,

\[
\left[ \S n \in \mathbb{N} \quad n^2 \leq 4 \right] = \{0, 1, 2\}
\]

is read "those *n* in **N** such that *n*\(^{2}\) ≤ 4 are in {0,1,2}." The same construct is expressible in set-builder notation as

\[
\{n \in \mathbb N: n^2 \le 4\} = \{0, 1, 2\}.
\]

Contrary to the other quantifiers, § yields a set rather than a formula.

Some other quantifiers sometimes used in mathematics include:
*There are infinitely many elements such that...
*For all but finitely many elements... (sometimes expressed as "for almost all elements...").
*There are uncountably many elements such that...
*For all but countably many elements...
*For all elements in a set of positive measure...
*For all elements except those in a set of measure zero...

## History
Term logic, also called Aristotelian logic, treats quantification in a manner that is closer to natural language, and also less suited to formal analysis. Term logic treated *All*, *Some* and *No* in the 4th century BC, in an account also touching on the alethic modalities.

In 1827, George Bentham published his *Outline of a New System of Logic: With a Critical Examination of Dr. Whately's Elements of Logic*, describing the principle of the quantifier, but the book was not widely circulated.

William Hamilton claimed to have coined the terms "quantify" and "quantification", most likely in his Edinburgh lectures c. 1840. Augustus De Morgan confirmed this in 1847, but modern usage began with De Morgan in 1862 where he makes statements such as "We are to take in both *all* and *some-not-all* as quantifiers".

Gottlob Frege, in his 1879 , was the first to employ a quantifier to bind a variable ranging over a domain of discourse and appearing in predicates. He would universally quantify a variable (or relation) by writing the variable over a dimple in an otherwise straight line appearing in his diagrammatic formulas. Frege did not devise an explicit notation for existential quantification, instead employing his equivalent of ~∀*x*~, or contraposition. Frege's treatment of quantification went largely unremarked until Bertrand Russell's 1903 *Principles of Mathematics*.

In work that culminated in Peirce (1885), Charles Sanders Peirce and his student Oscar Howard Mitchell independently invented universal and existential quantifiers, and bound variables. Peirce and Mitchell wrote Π\(_{x}\) and Σ\(_{x}\) where we now write ∀*x* and ∃*x*. Peirce's notation can be found in the writings of Ernst Schröder, Leopold Loewenheim, Thoralf Skolem, and Polish logicians into the 1950s. Most notably, it is the notation of Kurt Gödel's landmark 1930 paper on the completeness of first-order logic, and 1931 paper on the incompleteness of Peano arithmetic. Per Martin-Löf adopted a similar notation for dependent products and sums in his intuitionistic type theory, which are conceptually related to quantification.

Peirce's approach to quantification also influenced William Ernest Johnson and Giuseppe Peano, who invented yet another notation, namely (*x*) for the universal quantification of *x* and (in 1897) ∃*x* for the existential quantification of *x*. Hence for decades, the canonical notation in philosophy and mathematical logic was (*x*)*P* to express "all individuals in the domain of discourse have the property *P*", and "(∃*x*)*P*" for "there exists at least one individual in the domain of discourse having the property *P*". Peano, who was much better known than Peirce, in effect diffused the latter's thinking throughout Europe. Peano's notation was adopted by the *Principia Mathematica* of Whitehead and Russell, Quine, and Alonzo Church. In 1935, Gentzen introduced the ∀ symbol, by analogy with Peano's ∃ symbol. ∀ did not become canonical until the 1960s.

Around 1895, Peirce began developing his existential graphs, whose variables can be seen as tacitly quantified. Whether the shallowest instance of a variable is even or odd determines whether that variable's quantification is universal or existential. (Shallowness is the contrary of depth, which is determined by the nesting of negations.) Peirce's graphical logic has attracted some attention in recent years by those researching heterogeneous reasoning and diagrammatic inference.
