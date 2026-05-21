> Content sourced from [Algebrica](https://algebrica.org/geometric-sequence/) — CC BY-NC 4.0

## What is a geometric sequence

A [sequence](<../sequences>) \\( a_n \\) is called a **geometric sequence** (or geometric progression) if it consists of numbers arranged in such a way that the ratio between any term and the one before it is constant. It is characterized by terms of the form:

\\[a_1, a_2, \ldots, a_n \quad \text{with} \quad \frac{a_n}{a_{n-1}} = r \\]

  * By convention, the first term of a geometric progression is typically indexed with \\( n = 1 \\).
  * \\( r \\) represents the ratio between two consecutive terms in a geometric progression, and it is known as the common ratio.
  * If \\( r > 1 \\), the progression is increasing ([exponentially](<../exponential-function/>)).
  * If \\( 0 < r < 1 \\), the progression is decreasing toward zero.
  * If \\( r = 1 \\), the progression is constant.
  * If \\( r < 0 \\), the progression alternates in sign.


Let us consider, for example, the sequence with general term:

\\[a_n = 3 \cdot 2^{n-1}\\]

This sequence is a geometric progression starting at 3, where each term is obtained by multiplying the previous one by \\( r = 2 \\).

![](/diagrams/algebrica/geometri-sequence-1.png)


A geometric progression can also be defined recursively, meaning that each term is determined from the previous one. The recursive definition is given by:

\\[\begin{cases} a_1 = a \\\\[0.5em] a_n = a_{n-1} \cdot r \quad \text{for } n \geq 2 \end{cases} \\]

  * \\( a \in \mathbb{R} \\) is the first term,
  * \\( r \in \mathbb{R} \\) is the common ratio,
  * \\( a_n \\) is the general term of the sequence.


![](/diagrams/algebrica/geometri-sequence-2-1.png)

A geometric progression exhibits a characteristic exponential growth pattern, where the ratio between consecutive terms remains constant, leading to rapid increases (or decreases) in magnitude.

##### An [arithmetic progression](<../arithmetic-sequence/>), instead, exhibits a characteristic linear growth pattern, where the difference between consecutive terms remains constant, resulting in a steady increase (or decrease) over time.


In a geometric progression, each term \\( a_n \\) is obtained by multiplying the first term \\( a_1 \\) by the common ratio \\( r \\) raised to the power \\( (n - 1) \\). This gives the general formula for the \\( n \\)-th term:

\\[a_n = a_1 \cdot r^{n - 1} \quad \text{for } n \geq 1 \\]

This formula allows you to compute any term in the sequence directly, without knowing or listing all the previous ones.

The key difference between the explicit form and the recursive form of a sequence lies in how each term is defined:

  * In the explicit form, each term \\( a_n \\) is directly defined as a function of \\( n \\).  
You can compute any term independently, without needing the previous ones.

  * In the recursive form, each term \\( a_n \\) is defined based on one or more previous terms in the sequence. To compute a given term, you must first know the preceding ones.


**Example**

Let’s define a geometric sequence with first term \\( a_1 = 2 \\) and common ratio ( r = 3 ). We use the formula:

\\[a_n = a_1 \cdot r^{n - 1} \\]

Plug in the values:

\\[a_n = 2 \cdot 3^{n - 1} \\]

Now calculate the first few terms:

  * \\( a_1 = 2 \\)
  * \\( a_2 = 2 \cdot 3^1 = 6 \\)
  * \\( a_3 = 2 \cdot 3^2 = 18 \\)
  * \\( a_4 = 2 \cdot 3^3 = 54 \\)
  * \\( a_5 = 2 \cdot 3^4 = 162 \\)


The resulting sequence is: \\[2,\ 6,\ 18,\ 54,\ 162,\ \dots \\]

## Sum of \\( n \\) terms of a geometric progression

The sum \\( S_n \\) of the first \\( n \\) terms \\( a_1, a_2, \dots, a_n \\) of a geometric progression is given by the formula:

\\[S_n = a_1 \cdot \frac{1 - r^n}{1 - r} \quad \text{for } r \neq 1 \\]

This formula allows you to quickly compute the total sum of a finite number of terms in a geometric progression. For example, consider the geometric progression:

\\[2,\ 4,\ 8,\ 16,\ 32 \\]

We want to calculate the sum of the first 5 terms (\\( n = 5 \\)). Using the formula, we have:

\\[S_5 = 2 \cdot \frac{1 - 2^5}{1 - 2} = 2 \cdot \frac{1 - 32}{-1} = 2 \cdot 31 = 62 \\]

## Limit behavior of a geometric sequence

Regarding the [limits](<../limits>) of a geometric sequence, we observe the following behavior:

  * It diverges to \\( +\infty \\) if \\( r > 1 \\).
  * It is constant (that is, \\( a_n = a_0 \\) for every \\( n \in \mathbb{N} \\)) if \\( q = 1 \\), and thus \\(\lim_{n \to +\infty} a_n = a_0 = 1.\\)
  * It is infinitesimal if \\( |r| < 1 \\), meaning the terms approach zero.
  * It is oscillatory (irregular) if \\( r \leq -1 \\), due to alternating signs and unbounded growth.


The previously shown progression:

\\[a_n = 2 \cdot 3^{n - 1} \\]

diverges because the common ratio between its terms is \\( r = 3 \\), which satisfies \\( r > 1 \\). As a result, the terms grow exponentially and tend to \\( +\infty \\) as \\( n \\)\ increases.


Let us consider, for example, the geometric progression shown in the figure:

\\[a_n = (-2)^{n - 1} \\]

Expanding the sequence, we observe that the common ratio is \\( r = -2 \\), which satisfies \\( r \leq -1 \\). As a result, the sequence displays an irregular oscillatory behavior, alternating the sign of each term while the absolute values grow exponentially.

![](/diagrams/algebrica/geometri-sequence-3-2.png)
