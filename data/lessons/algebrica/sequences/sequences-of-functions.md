> Content sourced from [Algebrica](https://algebrica.org/sequences-of-functions/) — CC BY-NC 4.0

## Introduction

Imagine you have a list of different [functions](<../functions>), where each function in the list is linked to a number \\(n = 1, 2, 3… \in \mathbb{N} \\). So, for each \\(n\\), you get a different function, and that ordered list of functions is essentially what a [sequence](<../sequences>) of functions is. In formal terms, let \\( A \subseteq \mathbb{R} \\) be a non-empty subset and suppose that for each \\( n \in \mathbb{N} \\) we have a function \\( f_n: A \rightarrow \mathbb{R} \\). We then say that \\( (f_n) = (f_1, f_2, f_3, \dots) \\) is a **sequence of functions** on \\( A. \\)


Let’s consider a simple practical example. Let \\( (f_n) \\) be a sequence of functions, with \\(n \in \mathbb{N}_0\\) and \\( x \in \mathbb{R} \\), defined by:

\\[f_n(x) = \frac{x}{n+1} \\]

This is a family of functions where each function is linear, and the slope decreases as \\(n\\) increases. For \\(n = 0, 1, 2, 3, …\\), we have:

\\[\begin{aligned} f_0(x) &= \frac{x}{0+1} = x \\\\[0.5em] f_1(x) &= \frac{x}{1+1} = \frac{x}{2} \\\\[0.5em] f_2(x) &= \frac{x}{2+1} = \frac{x}{3} \\\\[0.5em] f_3(x) &= \frac{x}{3+1} = \frac{x}{4} \\\\[0.5em] \vdots \end{aligned} \\]

Thus, the sequence of functions is:

\\[(f_n) = \left( x, \frac{x}{2}, \frac{x}{3}, \frac{x}{4}, \dots \right) \\]

Graphically, this situation can be observed as follows:

![](/diagrams/algebrica/sequence-functions-1.png)

The graph shows how, as the index \\(n\\) increases, the slope of the line \\(f_n(x)\\) progressively decreases. This reflects the fact that the function flattens toward the zero function \\(f(x) = 0\\) for every \\(x\\). In other words, the sequence of functions \\(f_n(x)\\) converges pointwise to the zero function as \\(n\\) approaches infinity.

## Pointwise convergence

Let \\( \lbrace f_n(x) \rbrace \\) be a sequence of functions defined on a common [domain](<../determining-the-domain-of-a-function/>) \\( A \subseteq \mathbb{R} \\), with \\( n \in \mathbb{N} \\). We say that the sequence \\( \lbrace f_n(x) \rbrace \\) converges pointwise on a set \\( C \subseteq A \\) if, for every \\( x \in C \\), the numerical sequence \\( \lbrace f_n(x) \rbrace \\) converges. In this case, the [limit](<../limits>) function \\( f(x) \\) is defined as:

\\[f(x) = \lim_{n \to +\infty} f_n(x) \quad \forall \, x \in C \\]

The set \\( C \\) is called the **pointwise convergence** set of the sequence \\( {f_n(x)} \\).


Pointwise convergence can also be expressed as follows. Let \\((f_n)\\) be a sequence of functions defined on a set \\(A\\). Then \\((f_n)\\) converges pointwise to \\(f : A \to \mathbb{R}\\) if and only if \\(\forall \, x \in A\\) and \\(\forall \varepsilon > 0\\) exists \\(K \in \mathbb{N} \\) such that:

\\[|f_n(x) - f(x)| < \varepsilon \quad \forall \; n \geq K\\]

##### In other words, for each fixed point \\( x \\), we can make \\( f_n(x) \\) as close as we like to \\( f(x) \\) by choosing \\( n \\) large enough. The index \\( K \\) required to achieve the desired accuracy may vary depending on \\( x \\) and \\( \varepsilon \\).


Let us consider the sequence of functions as previously discussed:

\\[f_n(x) = \frac{x}{n+1}, \quad x \in \mathbb{R}, \quad n \in \mathbb{N}_0 \\]

Let us examine what happens for each \\(x\\) as \\(n \to \infty\\). If we fix a generic \\(x\\), for example \\(x = 2\\), the associated sequence is:

\\[\begin{aligned} f_0(2) &= 2 \\\\[0.5em] f_1(2) &= 1 \\\\[0.5em] f_2(2) &= \frac{2}{3} \\\\[0.5em] f_3(2) &= \frac{2}{4} \\\\[0.5em] &\vdots \end{aligned} \\]

This sequence of numbers tends to zero as \\(n \to \infty\\). In general, for every \\(x \in \mathbb{R}\\):

\\[\lim_{n \to \infty} f_n(x) = \lim_{n \to \infty} \frac{x}{n+1} = 0 \\]

Therefore, the limit function is:

\\[f(x) = 0 \quad \forall \, x \in \mathbb{R} \\]

By the uniqueness of limits of sequences of real numbers, the pointwise limit of a sequence of functions \\( (f_n) \\) is unique.

## Example

Let’s study the behavior of the following sequence of functions in the interval \\(-1 < x < 1\\):

\\[f_n(x) = x^n \\]

For a fixed value of \\(x\\) in this interval, we know that the [absolute value](<../absolute-value>) of \\(x\\) is less than \\(1\\), that is, \\(|x| < 1\\). This means we are considering [powers](<../powers>) of a number smaller than \\(1\\) in absolute value. By properties of [exponents](<../exponential-function>), when the base has an absolute value less than \\(1\\), the sequence \\(x^n\\) tends to zero as \\(n\\) tends to infinity:

\\[\lim_{n \to \infty} x^n = 0 \\]

The sequence of powers of a real number \\(x\\) with \\(|x| < 1\\) converges to zero.

Therefore, for every \\(x\\) in the interval \\(-1 < x < 1\\), the sequence of functions \\(f_n(x) = x^n\\) converges pointwise to the zero function.

\\[f(x) = 0 \quad \forall \, x \in (-1,1) \\]

## Consequences of pointwise convergence

Let \\( {f_n} \\) be a sequence of functions \\( f_n : A \to \mathbb{R} \\) that converges pointwise to a function \\( f : A \to \mathbb{R} \\). The following properties hold:

  * If each \\( f_n(x) \geq 0 \\) for all \\( x \in A \\), then \\( f(x) \geq 0 \\) for all \\( x \in A \\). In practice, if each function \\(f_n(x)\\) is non-negative on \\(A\\), then the limit function \\(f(x)\\) will also be non-negative on \\(A\\). This reflects the fact that the limit of a sequence of non-negative real numbers cannot be negative.

  * If each \\( f_n \\) is non-decreasing on \\( A \\), then \\( f \\) is non-decreasing on \\( A \\). Consequently, if each function \\(f_n\\) is non-decreasing on \\(A\\), then the limit function \\(f\\) will also be non-decreasing on \\(A\\). In other words, the property of monotonicity is preserved under pointwise convergence.


## Uniform convergence

Let \\( (f_n) \\) be a sequence of functions defined on a set \\( A \subseteq \mathbb{R} \\). We say that \\( (f_n) \\) **converges uniformly** on \\( A \\) to the function \\( f : A \to \mathbb{R} \\) if, for any \\( \varepsilon > 0 \\), there exists a [natural number](<../natural-numbers/>) \\( K \\) such that, for all \\( n \geq K \\) and all \\( x \in A \\), the following inequality holds:

\\[|f_n(x) - f(x)| < \varepsilon \\]

If \\( (f_n) \\) converges uniformly to \\( f \\), then \\( (f_n) \\) also converges pointwise to \\( f. \\)


Let us consider the sequence of functions:

\\[f_n(x) = \frac{x}{n}, \quad x \in [0,1], \quad n \in \mathbb{N}. \\]

For each fixed \\(x\\) in the interval \\([0,1]\\), we have:

\\[\lim_{n \to \infty} f_n(x) = 0 \\]

This means that the sequence \\(f_n(x)\\) converges pointwise to the function (f(x) = 0). Let us now check whether the convergence is uniform on \\([0,1]\\). We compute the difference between \\(f_n(x)\\) and the limit function \\(f(x)\\):

\\[|f_n(x) - f(x)| = \left| \frac{x}{n} - 0 \right| = \frac{x}{n} \\]

The maximum value of this difference on the interval \\([0,1]\\) is:

\\[\sup_{x \in [0,1]} |f_n(x) - f(x)| = \frac{1}{n} \\]

Given any \\(\varepsilon > 0\\), we can choose \\(N\\) such that:

\\[\frac{1}{N} < \varepsilon \\]

Therefore, for all \\(n \geq N\\) and for all \\(x \in [0,1]\\), we have:

\\[|f_n(x) - f(x)| < \varepsilon \\]

This confirms that the sequence of functions \\(f_n(x) = \frac{x}{n}\\) converges uniformly to the limit function \\(f(x) = 0\\) on the interval \\([0,1]\\).
