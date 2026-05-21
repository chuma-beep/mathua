> Content sourced from [Algebrica](https://algebrica.org/squeeze-theorem/) — CC BY-NC 4.0

## What is the Squeeze Theorem

The **Squeeze Theorem** , also referred to as the Sandwich Theorem, provides a method for determining the [limit](<../limits/>) of a [function](<../functions/>) when direct evaluation is challenging or when the function displays complex oscillatory behaviour near a specific point. This theorem is frequently applied to functions involving [sine and cosine](<../sine-and-cosine/>), particularly when these trigonometric terms exhibit oscillatory behaviour that precludes straightforward limit evaluation, such as:

\\[\sin\left(\frac{1}{x}\right) \quad \text{or} \quad \cos\left(\frac{1}{x}\right) \\]

In these situations, the function is constrained between two other functions with known and equal limits, which facilitates the evaluation of the target limit.

## Statement

Let \\( x_0 \in \mathbb{R} \cup { \pm\infty } \\) be a limit point, meaning that every neighborhood of \\(x_0\\) contains at least one point of the [domain](<../determining-the-domain-of-a-function/>) different from \\(x_0\\). Let \\( f \\), \\( g \\), and \\( h \\) be real-valued functions defined on a neighborhood \\(I\\) of \\(x_0\\). Assume that for every \\( x \in I \\), the following [inequality](<../linear-inequalities/>) holds:

\\[g(x) \leq f(x) \leq h(x) \\]

Also assume that the limits of \\( f(x) \\) and \\( h(x) \\) as \\( x \to x_0 \\) exist and are equal to some real number \\( \ell \\): \\[\lim_{x \to x_0} g(x) = \lim_{x \to x_0} h(x) = \ell \\]

Then, under these hypotheses, the function \\( g(x) \\) also admits a limit as \\( x \to x_0 \\), and that limit is: \\[\lim_{x \to x_0} f(x) = \ell \\]


In the graph, the black curve representing \\( f(x) \\) lies entirely between the lower bound \\( g(x) \\) and the upper bound \\( h(x) \\). As both bounding functions tend to \\( \ell \\), the function \\( f(x) \\) is forced to approach the same limit.

![](/diagrams/algebrica/squeeze-theorem.png)

This demonstrates the geometric intuition underlying the theorem: if a function is bounded above and below by two functions that both converge to the same value, then it must converge to that value.

## Proof of the Squeeze Theorem

Let \\( \varepsilon > 0 \\) be arbitrary. Our objective is to prove that the function \\( f(x) \\), which is bounded between \\( g(x) \\) and \\( h(x) \\), tends to the same limit \\( \ell \\) as \\( x \to x_0 \\). By assumption, we know that \\( \lim_{x \to x_0} g(x) = \ell \\). This means that there exists a positive number \\( \delta_1 \\) such that for every \\( x \\) sufficiently close to \\( x_0 \\) (specifically, for all \\( x \\) with \\( 0 < |x - x_0| < \delta_1 \\)), we have:

\\[|g(x) - \ell| < \varepsilon \quad \to \quad \ell - \varepsilon < g(x) < \ell + \varepsilon \\]

Similarly, since \\( \lim_{x \to x_0} h(x) = \ell \\), there exists another positive number \\( \delta_2 \\) such that:

\\[|h(x) - \ell| < \varepsilon \quad \to \quad \ell - \varepsilon < h(x) < \ell + \varepsilon \\]

Now let \\( \delta = \min(\delta_1, \delta_2) \\). Then for every \\( x \\) such that \\( 0 < |x - x_0| < \delta \\), both inequalities above are satisfied. But \\( f(x) \\) is squeezed between \\( g(x) \\) and \\( h(x) \\), so:

\\[g(x) \leq f(x) \leq h(x) \\]

Combining this with the bounds on \\( g(x) \\) and \\( h(x) \\), we obtain:

\\[\ell - \varepsilon < f(x) < \ell + \varepsilon \quad \to \quad |f(x) - \ell| < \varepsilon \\]

Since this inequality holds for every \\( \varepsilon > 0 \\), we conclude that:

\\[\lim_{x \to x_0} f(x) = \ell \\]

## Example

The subsequent example demonstrates how the theorem is applied to compute the following limit:

\\[\lim_{x \to 0} x \cdot \sin\left( \frac{1}{x} \right) \\]


The term \\( \sin\left( \frac{1}{x} \right) \\) does not admit a limit as \\( x \to 0 \\), since it oscillates indefinitely between \\(-1\\) and \\(1\\). However, for every [real number](<../types-of-numbers>) \\( x \neq 0 \\), the following inequality holds:

\\[-1 \leq \sin\left( \frac{1}{x} \right) \leq 1 \\]

Multiplying the entire inequality by \\( x \\), we obtain:

\\[-|x| \leq x \cdot \sin\left( \frac{1}{x} \right) \leq |x| \\]


Indeed, when \\( x > 0 \\), the inequality is preserved, while for \\( x < 0 \\), the inequality is reversed, but the [absolute value](<../absolute-value>) ensures that the comparison remains symmetric with respect to zero. Now observe that both bounding functions \\( -|x| \\) and \\( |x| \\) tend to zero as \\( x \to 0 \\):

\\[\lim_{x \to 0} -|x| = 0 \qquad \lim_{x \to 0} |x| = 0 \\]

Since \\( f(x) \\) is squeezed between two functions that both approach zero, we can apply the squeeze theorem and conclude that:

\\[\lim_{x \to 0} x \cdot \sin\left( \frac{1}{x} \right) = 0 \\]

In many instances, when an oscillating function is multiplied by a power of \\( x \\) that approaches zero, the overall limit is zero. This result arises because the oscillation remains bounded, as demonstrated by sine and cosine functions, which are always confined between \\(-1\\) and \\(1\\). Conversely, the factor \\( x^n \\) approaches zero rapidly enough to dominate the oscillation, causing the entire product to converge to zero.

## Exercises: compute the following limits using the Squeeze Theorem

  * \\[\text{1. } \quad \lim_{x \to +\infty} \frac{\ln(3 + \sin x)}{x^3}\\] [solution](<#e1>)

  * \\[\text{2. } \quad \lim_{x \to 0} x^4 \cdot \cos\left(\frac{2}{x}\right) + 2 \\] [solution](<#e2>)


## Exercise 1

Evaluate the following limit:

\\[\lim_{x \to +\infty} \frac{\ln(3 + \sin x)}{x^3} \\]


To begin, observe that the sine function is always bounded between \\(-1\\) and \\(1\\) for all real \\( x \\), so we can write:

\\[-1 \leq \sin x \leq 1 \\]

From the previous inequality, we can write:

\\[2 \leq 3 + \sin x \leq 4 \quad \text{for all } x \in \mathbb{R} \\]


Now, since the [logarithmic function](<../logarithmic-function/>) is strictly [increasing](<../increasing-and-decreasing-functions/>), we have:

\\[\log 2 \leq \log(3 + \sin x) \leq \log 4 \\]

We now divide all parts of the inequality by \\( x^3 \\) obtaining:

\\[\frac{\log 2}{x^3} \leq \frac{\ln(3 + \sin x)}{x^3} \leq \frac{\log 4}{x^3} \quad \forall \, x > 0 \\]

Since both bounding functions tend to zero as \\( x \to +\infty \\), we apply the Squeeze Theorem and obtain:

\\[\lim_{x \to +\infty} \frac{\ln(3 + \sin x)}{x^3} = 0 \\]

## Exercise 2

Evaluate the following limit:

\\[\lim_{x \to 0} x^4 \cdot \cos\left( \frac{2}{x} \right) + 2 \\]

To do so, we start by analyzing the behavior of the function \\( x^4 \cdot \cos\left( \frac{2}{x} \right) \\). We know that the cosine function is bounded between \\(-1\\) and \\(1\\) for all real values:

\\[-1 \leq \cos\left( \frac{2}{x} \right) \leq 1 \\]

Multiplying all parts of this inequality by \\( x^4 \\), which is always non-negative, we get:

\\[-x^4 \leq x^4 \cdot \cos\left( \frac{2}{x} \right) \leq x^4 \\]


Now we take the limit of the left and right bounds as \\( x \to 0 \\):

\\[\lim_{x \to 0} (-x^4) = 0 \qquad \lim_{x \to 0} x^4 = 0 \\]

Therefore, by the Squeeze Theorem, we conclude:

\\[\lim_{x \to 0} x^4 \cdot \cos\left( \frac{2}{x} \right) = 0 \\]

Now we return to the original expression:

\\[\lim_{x \to 0} \left( x^4 \cdot \cos\left( \frac{2}{x} \right) + 2 \right) \\]

Since:

\\[\lim_{x \to 0} x^4 \cdot \cos\left( \frac{2}{x} \right) = 0 \\]

we obtain

\\[\lim_{x \to 0} x^4 \cdot \cos\left( \frac{2}{x} \right) + 2 = 0 + 2 = 2 \\]

## Selected references

  * **MIT OpenCourseWare, C. Rodriguez**. [The Squeeze Theorem](https://ocw.mit.edu/courses/18-100a-real-analysis-fall-2020/mit18_100af20_lec8.pdf)
  * **University of California, Berkeley, A. Vizeff**. [Limit Laws and the Squeeze Theorem](https://math.berkeley.edu/~avizeff/calculus-I-F22/lecture-4.pdf)
