> Content sourced from [Algebrica](https://algebrica.org/principle-of-mathematical-induction/) — CC BY-NC 4.0

## Inductive set

Mathematical induction is a fundamental principle used to rigorously prove statements concerning [natural numbers](<../natural-numbers>). A subset \\( A \subseteq \mathbb{R} \\) is said to be inductive if it satisfies the following two conditions:

  * \\( 0 \in A \\).
  * \\( \forall \, x \in A\\) it follows that \\(x + 1 \in A \\).


In other words, an inductive set contains the element \\( 0 \\) and is closed under the operation of adding \\( 1 \\).


Let \\( A, B \subseteq \mathbb{R} \\) be two inductive sets. Then, their intersection \\( A \cap B \\) is also an inductive set. Since \\( A \\) and \\( B \\) are inductive, we have:

  * \\( 0 \in A \\) and \\( 0 \in B \\), thus \\( 0 \in A \cap B \\).
  * If \\( x \in A \cap B \\), then \\( x \in A \\) and \\( x \in B \\). By the inductive property, \\( x + 1 \in A \\) and \\( x + 1 \in B \\), hence \\( x + 1 \in A \cap B \\).


Therefore, \\( A \cap B \\) satisfies the conditions to be an inductive set. In particular, the set of natural numbers \\( \mathbb{N} \subseteq \mathbb{R} \\) is itself an inductive set, and more precisely, it is the smallest inductive set, being the intersection of all inductive subsets of \\( \mathbb{R} \\).

## Mathematical induction

Starting from the formal definition of an inductive set, let us consider a set \\( A \subseteq \mathbb{N} \\) defined by a property \\( p(n) \\), such that \\(A = \lbrace n \in \mathbb{N} \mid p(n) \rbrace\\). Suppose the following conditions hold:

  * \\( p(0) \\) is true, that is, \\( 0 \in A. \\)
  * \\( p(n) \rightarrow p(n+1) \, \forall \ n \in \mathbb{N}, \\) that is, if \\( n \in A \\), then \\( n+1 \in A.\\)


Thus, \\( p(n) \\) is true for every \\( n \in \mathbb{N}.\\)


These two conditions correspond precisely to the structure of a proof by mathematical induction:

  * The base case consists of verifying that \\( p(0) \\) holds, meaning that \\( 0 \\) belongs to \\( A \\).
  * The inductive step consists of proving that, whenever \\( p(n) \\) is true for some \\( n \in \mathbb{N} \\), it follows that \\( p(n+1) \\) is also true, thereby ensuring that \\( n+1 \in A \\).


## Example 1

Let us see a practical application of the principle of mathematical induction by proving the following statement (the sum of the first \\( n \\) natural numbers):

\\[\sum_{i=1}^{n} i = \frac{n(n+1)}{2} \quad \forall \, n \in \mathbb{N} \\]

We first consider the base case. For \\( n = 0 \\), the left-hand side is an empty sum, which equals \\( 0 \\) by convention. The right-hand side gives:

\\[\frac{0(0+1)}{2} = 0 \\]

Thus, the formula holds when \\( n = 0 \\).


We now move to the inductive step. We assume that the formula is true for some arbitrary \\( k \in \mathbb{N} \\); that is, we suppose:

\\[\sum_{i=1}^{k} i = \frac{k(k+1)}{2}. \\]

Under this assumption, we must prove that the formula also holds for \\( k+1 \\). Starting from the left-hand side:

\\[\sum_{i=1}^{k+1} i = \frac{k(k+1)}{2} + (k+1). \\]

Simplifying, we obtain:

\\[\frac{k(k+1) + 2(k+1)}{2} = \frac{(k+1)(k+2)}{2}. \\]

This is exactly the right-hand side of the formula with \\( n = k+1 \\). Since the statement holds for \\( n = 0 \\) and its validity for \\( n = k \\) implies its validity for \\( n = k+1 \\), by the principle of mathematical induction, the formula is true for every natural number \\( n \\).

It follows that the formula holds for every \\( n \in \mathbb{N} \\).

## Applications: divisibility

We show that the expression \\( n^3 - n \\) is divisible by \\( 3 \\) for every \\( n \in \mathbb{N} \\). We start with the base case. For \\( n = 0 \\), we have \\( 0^3 - 0 = 0 \\), and \\( 3 \mid 0 \\) holds trivially since \\( 0 = 3 \cdot 0. \\)

We now turn to the inductive step. Assume that \\( 3 \mid k^3 - k \\) for some \\( k \in \mathbb{N} \\), that is, \\( k^3 - k = 3m \\) for some integer \\( m \\). Under this hypothesis, we want to show that divisibility by \\( 3 \\) carries over to the next term, meaning that \\( 3 \mid (k+1)^3 - (k+1) \\). Expanding the expression we obtain:

\\[\begin{align} (k+1)^3 - (k+1) &= k^3 + 3k^2 + 3k + 1 - k - 1 \\\\[6pt] &= (k^3 - k) + 3k^2 + 3k \end{align} \\]

By the inductive hypothesis, \\( k^3 - k = 3m \\), so the expression becomes:

\\[\begin{align} (k+1)^3 - (k+1) &= 3m + 3k^2 + 3k \\\\[6pt] &= 3(m + k^2 + k) \end{align} \\]

Since \\( m + k^2 + k \in \mathbb{Z} \\), it follows that \\( 3 \mid (k+1)^3 - (k+1) \\). By the principle of induction, the statement holds for every \\( n \in \mathbb{N} \\).

## Applications: inequality

We show that \\( 2^n > n \\) for every \\( n \in \mathbb{N} \\). This inequality expresses the fact that exponential growth eventually and permanently dominates linear growth, a pattern that induction captures precisely: once the inequality holds at a given step, the multiplicative nature of the left-hand side guarantees it continues to hold at the next.

We start with the base case. For \\( n = 0 \\), we have \\( 2^0 = 1 > 0 \\), so the inequality holds.

The inductive step proceeds as follows. Assume that \\( 2^k > k \\) for some \\( k \in \mathbb{N} \\). The goal is to demonstrate that \\( 2^{k+1} > k+1 \\). Consider the left-hand side:

\\[\begin{align} 2^{k+1} &= 2 \cdot 2^k \\\ &> 2k \\\ &= k + k \end{align} \\]

Since \\( k \geq 0 \\), we have \\( k + k \geq k + 1 \\) for every \\( k \geq 1 \\), and the case \\( k = 0 \\) is already covered by the base case. Therefore \\( 2^{k+1} > k + 1 \\), and by the principle of induction, the inequality holds for every \\( n \in \mathbb{N} \\).

## Applications: geometric series

We show that for any real number \\( r \neq 1 \\) and every \\( n \in \mathbb{N} \\):

\\[\sum_{i=0}^{n} r^i = \frac{r^{n+1} - 1}{r - 1} \\]

This identity is the closed form of the partial sum of a [geometric sequence](<../geometric-sequence/>), and induction provides a clean route to its verification. We start with the base case. For \\( n = 0 \\), the left-hand side gives \\( r^0 = 1 \\), and the right-hand side gives:

\\[\frac{r^1 - 1}{r - 1} = 1 \\]

The formula holds for \\( n = 0 \\).


The inductive step proceeds as follows. Assume the formula holds for some \\( k \in \mathbb{N} \\), that is:

\\[\sum_{i=0}^{k} r^i = \frac{r^{k+1} - 1}{r - 1} \\]

We want to show it holds for \\( k + 1 \\). Starting from the left-hand side:

\\[\begin{align} \sum_{i=0}^{k+1} r^i &= \sum_{i=0}^{k} r^i + r^{k+1} \\\\[6pt] &= \frac{r^{k+1} - 1}{r - 1} + r^{k+1} \end{align} \\]

Combining the two terms over a common denominator:

\\[\begin{align} \frac{r^{k+1} - 1}{r - 1} + r^{k+1} &= \frac{r^{k+1} - 1 + r^{k+1}(r-1)}{r - 1} \\\\[6pt] &= \frac{r^{k+2} - 1}{r - 1} \end{align} \\]

This is exactly the right-hand side of the formula with \\( n = k+1 \\). By the principle of induction, the identity holds for every \\( n \in \mathbb{N} \\).

## Selected references

  * **MIT, F. Gotti**. [Mathematical Induction](https://math.mit.edu/~fgotti/docs/Courses/C.%20Combinatorial%20Analysis/2.%20Mathematical%20Induction/Mathematical%20Induction.pdf)

  * **University of California, Berkeley, P. Vojta**. [Three Types of Induction](https://math.berkeley.edu/~vojta/115/ho2.pdf)

  * **McGill University, J. Labute**. [Axioms for Set Theory](https://www.math.mcgill.ca/labute/courses/240f05/Sets.pdf)


Sequences

Sequences are ordered collections of numbers defined by a rule, each indexed by natural numbers.

2.1k

[Sequences](https://algebrica.org/sequences/)

1.9k

[Convergent and Divergent Sequences](https://algebrica.org/convergent-and-divergent-sequences/)

670

[Monotone Sequences](https://algebrica.org/monotone-sequences/)

1k

[Arithmetic Sequence](https://algebrica.org/arithmetic-sequence/)

1.1k

[Geometric Sequence](https://algebrica.org/geometric-sequence/)

1.4k

[Cauchy Sequence](https://algebrica.org/cauchy-sequence/)

2.1k

[Euler’s Number as the Limit of a Sequence](https://algebrica.org/euler-number-limit-sequence/)

1.4k

[Sequences of Functions](https://algebrica.org/sequences-of-functions/)
