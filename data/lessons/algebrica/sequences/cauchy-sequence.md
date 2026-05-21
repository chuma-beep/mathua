> Content sourced from [Algebrica](https://algebrica.org/cauchy-sequence/) — CC BY-NC 4.0

## What is a Cauchy sequence

A **Cauchy sequence** is a special type of [sequence](<../sequences>) where, as you move further along, the terms get closer and closer to each other. It doesn’t matter what the [limit](<../limits>) is or even if you know it: what matters is that the difference between the terms becomes smaller and smaller. This idea is important because it helps us understand whether a sequence is behaving in a stable way, even before knowing exactly what it’s approaching. In formal terms:


Let \\( {a_n} \\) be a sequence. Then \\( {a_n} \\) is a Cauchy sequence if and only if it is convergent. In other words, the following condition holds:

\\[\forall \varepsilon > 0 \quad \exists \nu \in \mathbb{N} : \quad |a_n - a_m| < \varepsilon \quad \forall n, m > \nu \\]

##### This criterion also applies to [series](<../series>) and allows us to prove [convergence](<../cauchy-convergence-criterion-series/>) without knowing the exact value of the sum. Instead of computing a limit, we check whether the terms of the sequence of partial sums get arbitrarily close to one another.


A sequence of the form \\( a_n = \frac{1}{n} \\) is a Cauchy sequence, where the values get closer and closer to each other as ( n ) increases. When plotted on a graph, it shows:

  * a curve starting from \\( a_1 = 1 \\) that decreases rapidly,
  * points that become increasingly dense near zero,
  * the distance between any two terms \\( a_n \\) and \\( a_m \\), for large \\( n \\) and \\( m \\), becomes smaller and smaller.


![Convergent sequence.](/diagrams/algebrica/sequences-conv-1.png)

As \\( n \\) increases, the terms become smaller and smaller, approaching zero. This is a classic example of a sequence that converges to 0.


Another example of a Cauchy sequence is given by the sequence defined as:

\\[a_n = 1 + \frac{1}{2} + \frac{1}{4} + \dots + \frac{1}{2^n} \\]

This is the sum of the first \\( n \\) terms of a geometric progression with ratio \\( \frac{1}{2} \\). This sequence is a Cauchy sequence because:

  * The terms get closer and closer to each other.
  * Each new term adds less and less to the total sum.
  * The distance between \\( a_n \\) and \\( a_m \\) (for \\( m > n \\)) becomes extremely small, since you are only adding values like:  
\\( \begin{align}\\\\[0.5em]\dfrac{1}{2^{n+1}}, \dfrac{1}{2^{n+2}}, \dots \end{align}\\)


![](/diagrams/algebrica/cauchy-sequence-1.png)

In the limit, this sequence converges to \\( 2 \\), reinforcing that it is both a Cauchy and convergent sequence.

##### In general, a numerical sequence is called a [geometric progression](<../sequences>) when the ratio between each term and its previous one is constant

## Theorem

Every convergent sequence \\( (x_n)_n \\) is a Cauchy sequence.

##### A Cauchy allows us to detect convergence based solely on how close the terms get to each other, without needing to know the actual limit. In complete spaces like the real numbers, this internal consistency is enough to ensure the sequence converges.


In fact, let \\( (x_n) \\) be a convergent sequence in \\( \mathbb{R} \\), and let \\( L \in \mathbb{R} \\) be its limit. By definition of convergence we have:

\\[\forall \varepsilon > 0, \ \exists N \in \mathbb{N} \ \text{such that} \ |x_n - L| < \frac{\varepsilon}{2} \quad \forall \, n \geq N \\]

Now, for any \\( n, m \geq N \\), we apply the triangle inequality:

\\[|x_n - x_m| = |x_n - L + L - x_m| \leq |x_n - L| + |x_m - L| < \frac{\varepsilon}{2} + \frac{\varepsilon}{2} = \varepsilon \\]

Thus, we have:

\\[\forall \varepsilon > 0, \ \exists N \in \mathbb{N} \ \text{such that} \ |x_n - x_m| < \varepsilon \quad \forall \, n, m \geq N \\]

This is exactly the definition of a Cauchy sequence. Therefore, every convergent sequence is a Cauchy sequence.

## Theorem

Every Cauchy sequence \\( (x_n) \\) is also a [bounded sequence](<../convergent-and-divergent-sequences/>). In fact, by definition, if \\( (x_n) \\) is a Cauchy sequence, then:

\\[\forall \varepsilon > 0, \ \exists N \in \mathbb{N} \ \text{such that} \ |x_n - x_m| < \varepsilon \quad \forall n, m \geq N \\]

Let’s choose \\( \varepsilon = 1 \\). Then there exists \\( N \in \mathbb{N} \\) such that:

\\[|x_n - x_m| < 1 \quad \forall n, m \geq N \\]

Fix \\( m = N \\), so we get:

\\[|x_n - x_N| < 1 \Rightarrow |x_n| \leq |x_N| + 1 \quad \forall n \geq N \\]

Now define:

\\[M_1 := \max{|x_0|, |x_1|, \dots, |x_{N-1}|}, \quad M_2 := |x_N| + 1 \\]

Let:

\\[M := \max{M_1, M_2} \Rightarrow |x_n| \leq M \quad \forall n \in \mathbb{N} \\]

Therefore, the sequence \\( (x_n) \\) stays entirely within a finite interval and is thus bounded.
