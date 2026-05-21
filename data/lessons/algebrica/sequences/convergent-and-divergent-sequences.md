> Content sourced from [Algebrica](https://algebrica.org/convergent-and-divergent-sequences/) — CC BY-NC 4.0

## Behavior of a sequence

We introduced [sequences](<../sequences>) as an ordered collection of elements, each assigned to a specific position indexed by a [natural number](<../types-of-numbers>). To every sequence \\( (a_n)_{n \in \mathbb{N}} \\), there is an associated behavior of its terms \\( a_n \\) that describes how they evolve as the index \\( n \\) increases. Analyzing this behavior helps determine whether the sequence converges to a finite [limit](<../limits>), diverges to infinity, or exhibits an oscillating pattern.

## Convergent sequence

A sequence \\( (a_n)_{n \in \mathbb{N}} \\) is said to be **convergent to the limit** \\( \ell \in \mathbb{R} \\) if for every \\( \varepsilon > 0 \\), there exists \\( n_0 \in \mathbb{N} \\) such that:

\\[|a_n - \ell| < \varepsilon \quad \text{for all } n \geq n_0. \\]

In this case, we write:

\\[\lim_{n \to +\infty} a_n = \ell \quad \text{or} \quad a_n \to \ell \quad \text{as } n \to +\infty. \\]

In other words, this means that the terms of the sequence get increasingly close to the number \\( \ell \\) as \\( n \\) grows larger. No matter how tight a margin \\( \varepsilon \\), from a certain index onward all terms will stay within that distance from \\( \ell \\). For example, consider the following sequence: \\[a_n = \left( \frac{1}{n}\right)_{n\geq} = \left(1, \frac{1}{2}, \frac{1}{3}, \ldots \right) \\]

![Convergent sequence.](/diagrams/algebrica/sequences-conv-1.png)

As \\( n \\) increases, the terms become smaller and smaller, approaching zero. This is a classic example of a sequence that converges to 0. A sequence is said to be **infinitesimal** when its terms get arbitrarily close to zero as the index grows and:

\\[\lim_{n \to +\infty} a_n = 0. \\]

The limit of a sequence \\( (a_n)_{n \in \mathbb{N}} \\), if it exists, is unique.

## Example

Let’s consider the sequence defined by:

\\[a_n = \frac{n}{n + 2} \\]

We aim to demonstrate that this sequence converges to 1 as \\( n \to +\infty \\), using the formal definition of convergence.


To prove this, we must show that for every \\( \varepsilon > 0 \\), there exists a natural number \\( n_0 \\) such that for all \\( n \geq n_0 \\):

\\[\left| \frac{n}{n + 2} - 1 \right| < \varepsilon \\]

Let’s simplify the absolute value expression:

\\[\left| \frac{n}{n + 2} - 1 \right| = \left| \frac{-2}{n + 2} \right| = \frac{2}{n + 2}. \\]

Now, we want:

\\[\frac{2}{n + 2} < \varepsilon \\]

Solving the inequality:

\\[n + 2 > \frac{2}{\varepsilon} \quad \Rightarrow \quad n > \frac{2}{\varepsilon} - 2 \\]

So we can define:

\\[n_0 = \left\lceil \frac{2}{\varepsilon} - 2 \right\rceil \\]

From this point onward, every term of the sequence stays within a distance \\( \varepsilon \\) of the limit \\(1.\\) Hence, by definition:

\\[\lim_{n \to +\infty} \frac{n}{n + 2} = 1. \\]

## Divergent sequence

A sequence \\( (a_n)_{n \in \mathbb{N}} \\) is said to be **divergent** if it does not converge to a finite limit. This can happen in the following ways.

A sequence diverges to \\( +\infty \\) if, for every \\( M > 0 \\), there exists an index \\( n_0 \in \mathbb{N} \\) such that  
\\[a_n > M \quad \text{for all } n \geq n_0 \\] In this case, we write: \\[\lim_{n \to +\infty} a_n = +\infty \quad \text{or} \quad a_n \to +\infty \text{ as } n \to +\infty \\]


A sequence diverges to \\( -\infty \\) if, for every \\( M < 0 \\), there exists an index \\( n_0 \in \mathbb{N} \\) such that  
\\[a_n < M \quad \text{for all } n \geq n_0 \\] In this case, we write: \\[\lim_{n \to +\infty} a_n = -\infty \quad \text{or} \quad a_n \to -\infty \text{ as } n \to +\infty \\]

## Bounded sequence

A bounded sequence is a sequence of numbers whose terms always stay within a fixed, finite interval, no matter how large the index becomes. In formal terms, let \\( {a_n} \\) be a sequence. We say that the sequence is bounded if there exists a constant \\( M > 0 \\) such that:

\\[|a_n| \leq M \quad \forall n \in \mathbb{N} \\]

We say that a sequence \\( {a_n} \\) is bounded above if there exists a constant \\( M \in \mathbb{R} \\) such that:

\\[a_n \leq M \quad \forall n \in \mathbb{N} \\]

We say that the sequence is bounded below if there exists a constant \\( M \in \mathbb{R} \\) such that:

\\[a_n \geq M \quad \forall n \in \mathbb{N} \\]

## Oscillating sequence

Oscillating sequences are a special type of bounded sequence. Let us consider the sequence:

\\[(a_n)_{n \in \mathbb{N}} = ((-1)^n)_{n \in \mathbb{N}} = (+1, -1, +1, -1, +1, -1, \dots) \\]

As the index \\(n\\) increases, the terms of the sequence alternate consistently between \\(+1\\) and \\(-1\\). This type of sequence does not approach any finite value and is called an **oscillating sequence**. It does not converge to a finite limit, nor does it diverge to \\( +\infty \\) or \\( -\infty \\), and its terms continue to fluctuate between different values

![Oscillating sequence.](/diagrams/algebrica/sequences-conv-2.png)

## Geometric sequence

Let us consider an example of a sequence, called a [geometric sequence](<../geometric-sequence/>), which can display different behaviors depending on the fixed real number \\( q \\). In general, a numerical sequence is called a geometric progression when the ratio between each term and its previous one is constant. More precisely, a geometric sequence is defined as follows:

\\[a_n := q^n \\]

It exhibits the following behavior:

  * It diverges to \\( +\infty \\) if \\( q > 1 \\).
  * It is constant (that is, \\( a_n = a_0 \\) for every \\( n \in \mathbb{N} \\)) if \\( q = 1 \\), and thus \\(\lim_{n \to +\infty} a_n = a_0 = 1.\\)
  * It is infinitesimal if \\( |q| < 1 \\), meaning the terms approach zero.
  * It is oscillatory (irregular) if \\( q \leq -1 \\), due to alternating signs and unbounded growth.


![](/diagrams/algebrica/sequences-conv-3.png)

As shown in the graph, when \\( q = 2 \\), the values of the geometric sequence \\( a_n = q^n \\) grow [exponentially](<../exponential-function>). As \\( n \\) increases, each term doubles the previous one, leading to a rapid escalation in magnitude.

##### Take a closer look at the difference between an [arithmetic progression](<../arithmetic-sequence>) and a geometric progression to better understand how their structures and growth patterns differ.
