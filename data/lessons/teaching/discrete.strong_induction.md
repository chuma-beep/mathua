# Strong Induction and Well-Ordering

**Strong induction:** To prove $P(n)$ for all $n\ge0$, prove $P(0)$ and show that if $P(k)$ holds for all $k<n$ then $P(n)$ holds. Equivalent to ordinary induction and to well-ordering of $\mathbb N$.

## Stronger Hypotheses

### Why Strong Helps
Recurrences like $a_n=a_{n-1}+a_{n-2}$ need two prior cases; ordinary induction's single $P(n-1)$ is insufficient, strong gives $P(n-1)$ and $P(n-2)$.

### Well-Ordering
Every non-empty subset of $\mathbb N$ has a least element. Minimal counterexample proves strong induction: if some $P(n)$ failed, the least failure contradicts the induction step.

## Example

Prove every $n>1$ is product of primes: assume every $k<n$ factors; if $n$ is prime done, else $n=ab$ with $1<a,b<n$, each factors by hypothesis, so $n$ factors. Hence $P(n)$ holds for all $n$.
