> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Mathematical_induction) — CC BY-SA 4.0

# Mathematical induction

**Mathematical induction** is a method for proving that a statement \(P(n)\) is true for every natural number \(n\), that is, that the infinitely many cases \(P(0), P(1), P(2), P(3), \dots\)  all hold. This is done by first proving a simple case, then also showing that if we assume the claim is true for a given case, then the next case is also true. Informal metaphors help to explain this technique, such as falling dominoes or climbing a ladder:

A **proof by induction** consists of two cases. The first, the **base case**, proves the statement for \(n = 0\) without assuming any knowledge of other cases. The second case, the **induction step**, proves that *if* the statement holds for any given case \(n = k\), *then* it must also hold for the next case \(n = k + 1\). These two steps establish that the statement holds for every natural number \(n\). The base case does not necessarily begin with \(n = 0\), but often with \(n = 1\), and possibly with any fixed natural number \(n = N\), establishing the truth of the statement for all natural numbers \(n \geq N\).

The method can be extended to prove statements about more general well-founded structures, such as trees; this generalization, known as structural induction, is used in mathematical logic and computer science. Mathematical induction in this extended sense is closely related to recursion. Mathematical induction is an inference rule used in formal proofs, and is the foundation of most correctness proofs for computer programs.

Despite its name, mathematical induction differs fundamentally from inductive reasoning as used in philosophy, in which the examination of many cases results in a probable conclusion. The mathematical method examines infinitely many cases to prove a general statement, but it does so by a finite chain of deductive reasoning involving the variable \(n\), which can take infinitely many values. The result is a rigorous proof of the statement, not an assertion of its probability.

## History
Acerbi (2000) claims that Plato's *Parmenides* (370 BC) contains traces of an early implicit inductive proof. However, Negrepontis and Farmaki (2021) have challenged this. While the Pythagoreans utilized inductive proofs, they did so without the principle of mathematical induction.

The earliest implicit proof by mathematical induction was written by al-Karaji around 1000 AD, who applied it to arithmetic sequences to prove the binomial theorem and properties of Pascal's triangle. Whilst the original work was lost, it was later referenced by Al-Samawal al-Maghribi in his treatise *al-Bahir fi'l-jabr (The Brilliant in Algebra)* in around 1150 AD.

Katz says in his history of mathematics

In India, early implicit proofs by mathematical induction appear in Bhaskara's "cyclic method".

None of these ancient mathematicians, however, explicitly stated the induction hypothesis. Another similar case (contrary to what Vacca has written, as Freudenthal carefully showed) was that of Francesco Maurolico in his *Arithmeticorum libri duo* (1575), who used the technique to prove that the sum of the first odd integers 

The earliest rigorous use of induction was by Gersonides (1288–1344). The first explicit formulation of the principle of induction was given by Pascal in his *Traité du triangle arithmétique* (1665). Another Frenchman, Fermat, made ample use of a related principle: indirect proof by infinite descent.

The induction hypothesis was also employed by the Swiss Jakob Bernoulli, and from then on it became well known. The modern formal treatment of the principle came only in the 19th century, with George Boole, Augustus De Morgan, Charles Sanders Peirce, Giuseppe Peano, and Richard Dedekind.

## Description
The simplest and most common form of mathematical induction infers that a statement involving a natural number (that is, an integer *n* ≥ 0 or 1) holds for all values of. The proof consists of two steps:
# The **** (or **initial case**): prove that the statement holds for 0, or 1.
# The **** (or **inductive step**, or **step case**): prove that for every , if the statement holds for , then it holds for *n* + 1. In other words, assume that the statement holds for some arbitrary natural number , and prove that the statement holds for *n* + 1.

The hypothesis in the induction step, that the statement holds for a particular , is called the **induction hypothesis** or **inductive hypothesis**. To prove the induction step, one assumes the induction hypothesis for and then uses this assumption to prove that the statement holds for *n* + 1.

Authors who prefer to define natural numbers to begin at 0 use that value in the base case; those who define natural numbers to begin at 1 use that value.

## Examples
### Sum of consecutive natural numbers
Mathematical induction can be used to prove the following statement for all natural numbers \(n\geq0\):

\[
P(n)\!:\ \ 0 + 1 + 2 + \cdots + n = \frac{n(n + 1)}{2}.
\]

This states a general formula for the sum of the natural numbers less than or equal to a given number; in fact an infinite sequence of statements: \(0 = \tfrac{(0)(0+1)}2\), \(0+1 = \tfrac{(1)(1+1)}2\), \(0+1+2 = \tfrac{(2)(2+1)}2\), etc.

**Proposition.** For every \(n\in\mathbb{N}\), we have that \(0 + 1 + 2 + \cdots + n = \tfrac{n(n + 1)}{2}.\)

**Proof.** Let \(P(n)\) be the statement \(0 + 1 + 2 + \cdots + n = \tfrac{n(n + 1)}{2}.\) We give a proof by induction on \(n\).

*Base case:* Show that the statement holds for the smallest natural number *n* = 0.

\(P(0)\) is clearly true: \(0 = \tfrac{0(0 + 1)}{2}\,.\)

*Induction step:* Show that for every \(k \geq 0\), if \(P(k)\) holds, then \(P(k+1)\) also holds.

Assume the induction hypothesis that for a particular \(k\), the single case \(n=k\) holds, meaning \(P(k)\) is true:

\[
0 + 1 + \cdots + k = \frac{k(k+1)}2.
\]

It follows that:

\[
(0 + 1 + 2 + \cdots + k )+ (k+1) = \frac{k(k+1)}2 + (k+1).
\]

Algebraically, the right hand side simplifies as:

\[
\begin{aligned}
\frac{k(k+1)}{2} + (k+1) &= \frac{k(k+1) + 2(k+1)}{2} \\
&= \frac{(k+1)(k+2)}{2} \\
&= \frac{(k+1)((k+1) + 1)}{2}.
\end{aligned}
\]

Equating the extreme left hand and right hand sides, we deduce that:
\[
0 + 1 + 2 + \cdots + k + (k+1) = \frac{(k+1)((k+1)+1)}2.
\]
 That is, the statement \(P(k+1)\) also holds true, establishing the induction step.

*Conclusion:* Since both the base case and the induction step have been proved as true, by mathematical induction the statement \(P(n)\) holds for every natural number \(n\geq0\). Q.E.D.

### A trigonometric inequality
Induction is often used to prove inequalities. As an example, we prove that \(\left|\sin nx\right| \leq n\left|\sin x\right|\) for any real number \(x\) and natural number \(n\).

At first glance, it may appear that a more general version, \(\left|\sin nx\right| \leq n\left|\sin x\right|\) for any *real* numbers \(n,x\), could be proven without induction; but the case \(n = \frac{1}{2},\, x=\pi\) shows it may be false for non-integer values of \(n\). This suggests we examine the statement specifically for *natural* values of \(n\), and induction is the readiest tool.

**Proposition.** For any \(x \in \mathbb{R}\) and \(n \in \mathbb{N}\), \(\left|\sin nx\right| \leq n\left|\sin x\right|\).

**Proof.** Fix an arbitrary real number \(x\), and let \(P(n)\) be the statement \(\left|\sin nx\right| \leq n\left|\sin x\right|\). We induct on \(n\).

*Base case:* The calculation \(\left|\sin 0x\right| = 0 \leq 0 = 0 \left|\sin x\right|\) verifies \(P(0)\).

*Induction step:* We show the implication \(P(k) \implies P(k+1)\) for any natural number \(k\). Assume the induction hypothesis: for a given value \(n = k \geq 0\), the single case \(P(k)\) is true. Using the angle addition formula and the triangle inequality, we deduce:

\[
\begin{aligned}
\left|\sin(k+1)x\right|
&= \left|\sin kx \cos x+\sin x \cos kx\right| && \text{(angle addition)}
\\
&\leq \left|\sin kx \cos x\right| + \left|\sin x\,\cos kx\right| && \text{(triangle inequality)}
\\
&= \left|\sin kx\right|\left| \cos x\right| + \left|\sin x\right|\left|\cos kx\right|
\\
&\leq \left|\sin kx\right| + \left|\sin x\right| && (\left|\cos t\right| \leq 1)
\\
&\leq
k\left|\sin x\right|+\left|\sin x\right| && \text{(induction hypothesis})\\
&= (k+1)\left|\sin x\right|.
\end{aligned}
\]

The inequality between the extreme left-hand and right-hand quantities shows that \(P(k+1)\) is true, which completes the induction step.

*Conclusion:* The proposition \(P(n)\) holds for all natural numbers \(n.\) Q.E.D.

## Variants

In practice, proofs by induction are often structured differently, depending on the exact nature of the property to be proven.
All variants of induction are special cases of transfinite induction; see below.

### Base case other than 0 or 1
If one wishes to prove a statement, not for all natural numbers, but only for all numbers greater than or equal to a certain number , then the proof by induction consists of the following:
# Showing that the statement holds when *n* = *b*.
# Showing that if the statement holds for an arbitrary number *n* ≥ *b*, then the same statement also holds for *n* + 1.
This can be used, for example, to show that 2\(^{*n*}\) ≥ *n* + 5 for *n* ≥ 3.

In this way, one can prove that some statement *P*(*n*) holds for all *n* ≥ 1, or even for all *n* ≥ −5. This form of mathematical induction is actually a special case of the previous form, because if the statement to be proved is *P*(*n*) then proving it with these two rules is equivalent with proving *P*(*n* + *b*) for all natural numbers with an induction base case 0.

#### Example: forming dollar amounts by coins
Assume an infinite supply of 4- and 5-dollar coins. Induction can be used to prove that any whole amount of dollars greater than or equal to 12 can be formed by a combination of such coins. Let *S*(*k*) denote the statement " dollars can be formed by a combination of 4- and 5-dollar coins". The proof that *S*(*k*) is true for all *k* ≥ 12 can then be achieved by induction on as follows:

*Base case:* Showing that *S*(*k*) holds for *k* = 12 is simple: take three 4-dollar coins.

*Induction step:* Given that *S*(*k*) holds for some value of *k* ≥ 12 (*induction hypothesis*), prove that *S*(*k* + 1) holds, too. Assume *S*(*k*) is true for some arbitrary *k* ≥ 12. If there is a solution for dollars that includes at least one 4-dollar coin, replace it by a 5-dollar coin to make *k* + 1 dollars. Otherwise, if only 5-dollar coins are used, must be a multiple of 5 and so at least 15; but then we can replace three 5-dollar coins by four 4-dollar coins to make *k* + 1 dollars. In each case, *S*(*k* + 1) is true.

Therefore, by the principle of induction, *S*(*k*) holds for all *k* ≥ 12, and the proof is complete.

In this example, although *S*(*k*) also holds for \(k \in \{ 4, 5, 8, 9, 10 \}\), the above proof cannot be modified to replace the minimum amount of 12 dollar to any lower value. For *m* = 11, the base case is actually false; for *m* = 10, the second case in the induction step (replacing three 5- by four 4-dollar coins) will not work; let alone for even lower. ### Induction on more than one counter
It is sometimes desirable to prove a statement involving two natural numbers, and , by iterating the induction process. That is, one proves a base case and an induction step for , and in each of those proves a base case and an induction step for. See, for example, the proof of commutativity accompanying *addition of natural numbers*. More complicated arguments involving three or more counters are also possible.

### Infinite descent

The method of infinite descent is a variation of mathematical induction which was used by Pierre de Fermat. It is used to show that some statement *Q*(*n*) is false for all natural numbers. Its traditional form consists of showing that if *Q*(*n*) is true for some natural number , it also holds for some strictly smaller natural number. Because there are no infinite decreasing sequences of natural numbers, this situation would be impossible, thereby showing (by contradiction) that *Q*(*n*) cannot be true for any. The validity of this method can be verified from the usual principle of mathematical induction. Using mathematical induction on the statement *P*(*n*) defined as "*Q*(*m*) is false for all natural numbers less than or equal to ", it follows that *P*(*n*) holds for all , which means that *Q*(*n*) is false for every natural number. ### Limited mathematical induction
If one wishes to prove that a property *P* holds for all natural numbers less than or equal to a fixed , proving that *P* satisfies the following conditions suffices:

# *P* holds for 0,
# For any natural number less than , if *P* holds for , then *P* holds for *x* + 1

### Prefix induction
The most common form of proof by mathematical induction requires proving in the induction step that

\[
\forall k \, (P(k) \to P(k+1))
\]

whereupon the induction principle "automates" applications of this step in getting from *P*(0) to *P*(*n*). This could be called "predecessor induction" because each step proves something about a number from something about that number's predecessor.

A variant of interest in computational complexity is "prefix induction", in which one proves the following statement in the induction step:

\[
\forall k\, (P(k) \to P(2k) \land P(2k+1))
\]

or equivalently

\[
\forall k\, \left( P\!\left(\left\lfloor \frac{k}{2} \right\rfloor \right) \to P(k) \right)
\]

The induction principle then "automates" log\(_{2}\) *n* applications of this inference in getting from *P*(0) to *P*(*n*). In fact, it is called "prefix induction" because each step proves something about a number from something about the "prefix" of that number — as formed by truncating the low bit of its binary representation. It can also be viewed as an application of traditional induction on the length of that binary representation.

If traditional predecessor induction is interpreted computationally as an -step loop, then prefix induction would correspond to a log--step loop. Because of that, proofs using prefix induction are "more feasibly constructive" than proofs using predecessor induction.

Predecessor induction can trivially simulate prefix induction on the same statement. Prefix induction can simulate predecessor induction, but only at the cost of making the statement more syntactically complex (adding a bounded universal quantifier), so the interesting results relating prefix induction to polynomial-time computation depend on excluding unbounded quantifiers entirely, and limiting the alternation of bounded universal and existential quantifiers allowed in the statement.

One can take the idea a step further: one must prove

\[
\forall k \, \left( P\!\left( \left\lfloor \sqrt{k} \right\rfloor \right) \to P(k) \right)
\]

whereupon the induction principle "automates" log log *n* applications of this inference in getting from *P*(0) to *P*(*n*). This form of induction has been used, analogously, to study log-time parallel computation.

### Complete (strong) induction
Another variant, called **complete induction**, **course of values induction** or **strong induction** (in contrast to which the basic form of induction is sometimes known as **weak induction**), makes the induction step easier to prove by using a stronger hypothesis: one proves the statement \(P(m+1)\) under the assumption that \(P(n)\) holds for *all* natural numbers \(n\) less than \(m+1\); by contrast, the basic form only assumes \(P(m)\). The name "strong induction" does not mean that this method can prove more than "weak induction", but merely refers to the stronger hypothesis used in the induction step.

In fact, it can be shown that the two methods are actually equivalent, as explained below. In this form of complete induction, one still has to prove the base case, \(P(0)\), and it may even be necessary to prove extra-base cases such as \(P(1)\) before the general argument applies, as in the example below of the Fibonacci number \(F_{n}\).

Although the form just described requires one to prove the base case, this is unnecessary if one can prove \(P(m)\) (assuming \(P(n)\) for all lower \(n\)) for all \(m \geq 0\). This is a special case of transfinite induction as described below, although it is no longer equivalent to ordinary induction. In this form the base case is subsumed by the case \(m = 0\), where \(P(0)\) is proved with no other \(P(n)\) assumed; this case may need to be handled separately, but sometimes the same argument applies for \(m = 0\) and \(m > 0\), making the proof simpler and more elegant.
In this method, however, it is vital to ensure that the proof of \(P(m)\) does not implicitly assume that \(m > 0\), e.g. by saying "choose an arbitrary \(n < m\)", or by assuming that a set of elements has an element.
#### Equivalence with ordinary induction
Complete induction is equivalent to ordinary mathematical induction as described above, in the sense that a proof by one method can be transformed into a proof by the other. Suppose there is a proof of \(P(n)\) by complete induction. Then, this proof can be transformed into an ordinary induction proof by assuming a stronger inductive hypothesis. Let \(Q(n)\) be the statement "\(P(m)\) holds for all \(m\) such that \(0\leq m \leq n\)"—this becomes the inductive hypothesis for ordinary induction. We can then show \(Q(0)\) and \(Q(n + 1)\) for \(n \in \mathbb N\) assuming only \(Q(n)\) and show that \(Q(n)\) implies \(P(n)\).

If, on the other hand, \(P(n)\) had been proven by ordinary induction, the proof would already effectively be one by complete induction: \(P(0)\) is proved in the base case, using no assumptions, and \(P(n+1)\) is proved in the induction step, in which one may assume all earlier cases but need only use the case \(P(n)\).

#### Example: Fibonacci numbers
Complete induction is most useful when several instances of the inductive hypothesis are required for each induction step. For example, complete induction can be used to show that

\[
F_n = \frac{\varphi^n - \psi^n}{\varphi - \psi}
\]

where \(F_n\) is the -th Fibonacci number, and \(\varphi = \frac{1}{2}(1 + \sqrt 5)\) (the golden ratio) and \(\psi = \frac{1}{2} (1 - \sqrt 5)\) are the roots of the polynomial \(x^2-x-1\). By using the fact that \(F_{n+2} = F_{n+1} + F_{n}\) for each \(n \in \mathbb{N}\), the identity above can be verified by direct calculation for \(F_{n+2}\) if one assumes that it already holds for both \(F_{n+1}\) and \(F_n\). To complete the proof, the identity must be verified in the two base cases: \(n = 0\) and \(n = 1\).

#### Example: prime factorization
Another proof by complete induction uses the hypothesis that the statement holds for *all* smaller \(n\) more thoroughly. Consider the statement that "every natural number greater than 1 is a product of (one or more) prime numbers", which is the "existence" part of the fundamental theorem of arithmetic. For proving the induction step, the induction hypothesis is that for a given \(m>1\) the statement holds for all smaller \(n>1\). If \(m\) is prime then it is certainly a product of primes, and if not, then by definition it is a product: \(m = n_1 n_2\), where neither of the factors is equal to 1; hence neither is equal to \(m\), and so both are greater than 1 and smaller than \(m\). The induction hypothesis now applies to \(n_1\) and \(n_2\), so each one is a product of primes. Thus \(m\) is a product of products of primes, and hence by extension a product of primes itself.

#### Example: dollar amounts revisited
We shall look to prove the same example as above, this time with *strong induction*. The statement remains the same:

\[
S(n): \,\,n \geq 12 \implies \,\exists\, a,b\in\mathbb{N}. \,\, n = 4a+5b
\]

However, there will be slight differences in the structure and the assumptions of the proof, starting with the extended base case.

**Proof.**

*Base case:* Show that \(S(k)\) holds for \(k = 12,13,14,15\).

\[
\begin{aligned}
4 \cdot 3+5 \cdot 0=12\\
4 \cdot 2+5 \cdot 13\\
4 \cdot 1+5 \cdot 2=14\\
4 \cdot 0+5 \cdot 3=15
\end{aligned}
\]

The base case holds.

*Induction step:* Given some \(j>15\), assume \(S(m)\) holds for all \(m\) with \(12 \leq m< j\). Prove that \(S(j)\) holds.

Choosing \(m=j-4\), and observing that \(15 < j \implies 12 \leq j-4 < j\) shows that \(S(j-4)\) holds, by the inductive hypothesis. That is, the sum \(j-4\) can be formed by some combination of \(4\) and \(5\) dollar coins. Then, simply adding a \(4\) dollar coin to that combination yields the sum \(j\). That is, \(S(j)\) holds. Q.E.D.

### Forward-backward induction

Sometimes, it is more convenient to deduce backwards, proving the statement for \(n-1\), given its validity for \(n\). However, proving the validity of the statement for no single number suffices to establish the base case; instead, one needs to prove the statement for an infinite subset of the natural numbers. For example, Augustin Louis Cauchy first used forward (regular) induction to prove the
inequality of arithmetic and geometric means for all powers of 2, and then used backwards induction to show it for all natural numbers.

## Example of error in the induction step

The induction step must be proved for all values of. To illustrate this, Joel E. Cohen proposed the following argument, which purports to prove by mathematical induction that all horses are of the same color:

*Base case:* in a set of only *one* horse, there is only one color.

*Induction step:* assume as induction hypothesis that within any set of \(n\) horses, there is only one color. Now look at any set of \(n+1\) horses. Number them: \(1, 2, 3, \dotsc, n, n+1\). Consider the sets \(\left\{1, 2, 3, \dotsc, n\right\}\) and \(\left\{2, 3, 4, \dotsc, n+1\right\}\). Each is a set of only \(n\) horses, therefore within each there is only one color. But the two sets overlap, so there must be only one color among all \(n+1\) horses.

The base case \(n=1\) is trivial, and the induction step is correct in all cases \(n > 1\). However, the argument used in the induction step is incorrect for \(n+2\), because the statement that "the two sets overlap" is false for \(\left\{1\right\}\) and \(\left\{2\right\}\).

## Formalization
In **second-order logic**, one can write down the "axiom of induction" as follows:

\[
\forall P\,\Bigl( P(0) \land \forall k \bigl( P(k) \to P(k+1)\bigr ) \to \forall n \,\bigl(P(n)\bigr)\Bigr),
\]

where *P*(·) is a variable for predicates involving one natural number and and are variables for natural numbers.

In words, the base case *P*(0) and the induction step (namely, that the induction hypothesis *P*(*k*) implies *P*(*k* + 1)) together imply that *P*(*n*) for any natural number. The axiom of induction asserts the validity of inferring that *P*(*n*) holds for any natural number from the base case and the induction step.

The first quantifier in the axiom ranges over *predicates* rather than over individual numbers. This is a second-order quantifier, which means that this axiom is stated in second-order logic. Axiomatizing arithmetic induction in first-order logic requires an axiom schema containing a separate axiom for each possible predicate. The article Peano axioms contains further discussion of this issue.

The axiom of structural induction for the natural numbers was first formulated by Peano, who used it to specify the natural numbers together with the following four other axioms:

# 0 is a natural number.
# The successor function of every natural number yields a natural number (*s*(*x*) = *x* + 1).
# The successor function is injective.
# 0 is not in the range of. In **first-order ZFC set theory**, quantification over predicates is not allowed, but one can still express induction by quantification over sets:

\[
\forall A \Bigl( 0 \in A \land \forall k \in \N \bigl( k \in A \to (k+1) \in A \bigr) \to \N\subseteq A\Bigr)
\]

 may be read as a set representing a proposition, and containing natural numbers, for which the proposition holds. This is not an axiom, but a theorem, given that natural numbers are defined in the language of ZFC set theory by axioms, analogous to Peano's. See construction of the natural numbers using the axiom of infinity and axiom schema of specification.

## Transfinite induction

One variation of the principle of complete induction can be generalized for statements about elements of any well-founded set, that is, a set with an irreflexive relation < that contains no infinite descending chains. Every set representing an ordinal number is well-founded, the set of natural numbers is one of them.

Applied to a well-founded set, transfinite induction can be formulated as a single step. To prove that a statement *P*(*n*) holds for each ordinal number:
# Show, for each ordinal number , that if *P*(*m*) holds for all *m* < *n*, then *P*(*n*) also holds.

This form of induction, when applied to a set of ordinal numbers (which form a well-ordered and hence well-founded class), is called *transfinite induction*. It is an important proof technique in set theory, topology and other fields.

Proofs by transfinite induction typically distinguish three cases:
# when is a minimal element, i.e. there is no element smaller than ;
# when has a direct predecessor, i.e. the set of elements which are smaller than has a largest element;
# when has no direct predecessor, i.e. is a so-called limit ordinal.

Strictly speaking, it is not necessary in transfinite induction to prove a base case, because it is a vacuous special case of the proposition that if *P* is true of all *n* < *m*, then *P* is true of. It is vacuously true precisely because there are no values of *n* < *m* that could serve as counterexamples. So the special cases are special cases of the general case.

## Relationship to the well-ordering principle
The principle of mathematical induction is usually stated as an axiom of the natural numbers; see Peano axioms. It is strictly stronger than the well-ordering principle in the context of the other Peano axioms. Suppose the following:
* The trichotomy axiom: For any natural numbers and , is less than or equal to if and only if is not less than. * For any natural number , *n* + 1 is greater than *n*.
* For any natural number , no natural number is between *n* and *n* + 1.
* No natural number is less than zero.

It can then be proved that induction, given the above-listed axioms, implies the well-ordering principle. The following proof uses complete induction and the first and fourth axioms.

**Proof.** Suppose there exists a non-empty set, of natural numbers that has no least element. Let *P*(*n*) be the assertion that is not in. Then *P*(0) is true, for if it were false then 0 is the least element of. Furthermore, let be a natural number, and suppose *P*(*m*) is true for all natural numbers less than *n* + 1. Then if *P*(*n* + 1) is false *n* + 1 is in , thus being a minimal element in , a contradiction. Thus *P*(*n* + 1) is true. Therefore, by the complete induction principle, *P*(*n*) holds for all natural numbers ; so is empty, a contradiction. Q.E.D.

Dark mode invert
