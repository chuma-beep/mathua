# Strong Induction and Well-Ordering

**Strong induction:** To prove $P(n)$ for all $n\ge0$, show $P(0)$ and show that $P(k)$ for all $k<n$ implies $P(n)$. The hypothesis covers every smaller case, not just the predecessor. It is equivalent to ordinary induction and to well-ordering: every nonempty subset of $\mathbb N$ has a least element.

## Worked: every n > 1 factors into primes

Prove $P(n)$: "$n$ is a product of primes", using all smaller cases:
1. Base $P(2)$: 2 is prime, hence a (one-factor) product.
2. Fix $n>2$ and assume $P(k)$ for every $2\le k<n$; if $n$ is prime, done.
3. Else $n=ab$ with $1<a,b<n$, and both factor by hypothesis, so concatenate their factorizations for $n$.

So $P(n)$ holds for all $n$: factoring $n$ needs factors strictly below $n$, exactly what the strong hypothesis supplies.

## Worked: Fibonacci growth needs two priors

Prove $a_n<2^n$ for Fibonacci $a_1=a_2=1$, $a_n=a_{n-1}+a_{n-2}$:
1. Bases $P(1)$ and $P(2)$: $1<2$ and $1<4$ both hold.
2. Fix $n>2$ and assume $P(k)$ for all $k<n$; in particular $a_{n-1}<2^{n-1}$ and $a_{n-2}<2^{n-2}$.
3. Add: $a_n<2^{n-1}+2^{n-2}<2^{n-1}+2^{n-1}=2^n$.

So depth-2 recurrences force the strong form: ordinary induction's lone $P(n-1)$ never mentions $a_{n-2}$.

## Minimal counterexamples

Well-ordering proves the strong step by contradiction: if some $P(n)$ failed, the failing set has a least element $m$, and minimality makes $P(k)$ true for all $k<m$. The induction step then proves $P(m)$, contradicting the choice of $m$. In general, least-element arguments, strong induction, and ordinary induction are three faces of the same principle on $\mathbb N$.
