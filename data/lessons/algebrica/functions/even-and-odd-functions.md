> Content sourced from [Algebrica](https://algebrica.org/even-and-odd-functions/) — CC BY-NC 4.0

## Behavior of a function

When analyzing the behavior of a [function](<../functions>), it is useful to investigate whether the function exhibits symmetry with respect to the coordinate axes. In this context, functions can be classified as **even** , showing symmetry with respect to the \\(y\\)-axis, or **odd** , exhibiting symmetry with respect to the origin. In general, a function can be:

  * even
  * odd
  * neither even nor odd


## Even function

More specifically, suppose we have a function \\( f(x): \mathbb{R} \rightarrow \mathbb{R} \\), and let \\( D \subseteq \mathbb{R} \\) be its [domain](<../determining-the-domain-of-a-function/>). The function \\( f \\) is said to be **even** if the following condition holds:

\\[f(x) = f(-x) \quad \text{for all } x \in D \\]

![](/diagrams/algebrica/even-odd-functions-1.png)

As shown in the figure, the function \\( f(x) = x^2 \\) is a [parabola](<../parabola>) symmetric with respect to the \\(y\\)-axis. In general, functions of the form \\( f(x) = x^4 \\), \\( x^6 \\), or more generally \\( x^{2n} \\), where the exponent is even, are examples of even functions.

![](/diagrams/algebrica/even-odd-functions-2-1.png)

Another example of an even function is the [cosine function](<../cosine-function>). It is a periodic function with period \\( 2\pi \\), and its graph is symmetric with respect to the \\(y\\)-axis. In fact, it is easy to verify that: \\[\cos(\pi) = \cos(-\pi) = -1 \\]

Another even function is the [absolute value function](<../absolute-value-function/>).


More generally, when considering the family of functions of the form \\(f(x) = x^{n}\\) with \\(n \in \mathbb{N},\\) the parity of the function is entirely determined by the exponent: the function behaves as an even function whenever \\(n\\) is an even [integer](<../integers/>), whereas it behaves as an odd function whenever \\(n\\) is odd.

## Definite integral of even function

One of the useful consequences of a function being even is the simplification it allows in [definite integrals](<../definite-integrals/>) over symmetric intervals. If \\( f(x) \\) is a [continuous](<../continuous-functions/>) and even function, then its graph is symmetric with respect to the \\(y\\)-axis.

This symmetry directly influences how we evaluate definite integrals over intervals of the form \\([ -a, a ]\\). Specifically, the following identity holds:

\\[\int_{-a}^{a} f(x),dx = 2\int_0^a f(x),dx \\]

![](/diagrams/algebrica/definite-integrals-5.png)

That is, the total area under the curve from \\(-a\\) to \\(a\\) is simply twice the area from 0 to (a). This works because the portion of the graph on the negative side of the \\(x\\)-axis is a mirror image of the positive side, and contributes the same value to the integral.

## Odd function

Suppose we have a function \\( f(x): \mathbb{R} \rightarrow \mathbb{R} \\), and let \\( D \subseteq \mathbb{R} \\) be its domain. The function \\( f \\) is said to be **odd** if the following condition holds:

\\[f(-x) = -f(x) \quad \text{for all } x \in D \\]

![](/diagrams/algebrica/even-odd-functions-3.png)

As shown in the figure, the function \\( f(x) = x^3 \\) is symmetric with respect to the origin. Functions of the form \\( f(x) = x^3 \\), \\( x^5 \\), or more generally \\( x^{2n+1} \\), where the exponent is odd, are examples of odd functions.

![](/diagrams/algebrica/even-odd-functions-4-1.png)

Another example of an odd function is the [sine function](<../sine-function>). It is a periodic function with period \\( 2\pi \\), and its graph is symmetric with respect to the origin. In fact, it is easy to verify that:  
\\[\sin(-\pi) = -\sin(\pi) = 0 \\]

## Definite integral of odd function

In the case of an odd function, the area between \\( [-a, 0] \\) is equal in magnitude but opposite in sign to the area between \\( [0, a] \\). Therefore, the definite integral is equal to:

\\[\int_{-a}^{a} f(x) \, dx = 0 \\]

![](/diagrams/algebrica/definite-integrals-6.png)

In both situations, the area enclosed between the graph of \\( f(x) \\) and the \\( x \\)-axis over the interval \\( [-a, a] \\) is given by:

\\[S = \int_{0}^{a} |f(x)| \, dx \\]

## The only function that is both even and odd

The function \\( f(x) = 0 \\) is the only function that is both even and odd, because it satisfies both \\( f(-x) = f(x) \\) and \\( f(-x) = -f(x) \\) for all \\( x \in \mathbb{R} \\). In fact, if a function were to be both even and odd, we would have:

  * \\(f(-x) = f(x)\\) when the function is even.
  * \\(f(-x) = -f(x)\\) when the function is odd.


Therefore, the zero function is the unique case that satisfies both properties.

## Properties

  * The sum of two even functions is even.
  * The product of an even function by a constant is even.
  * The product of two even functions is an even function.
  * The [derivative](<../derivatives>) of an even function is an odd function.
  * The sum of two odd functions is odd.
  * The product of an odd function by a constant is odd.
  * The product of two odd functions is an even function.
  * The derivative of an odd function is an even function.
  * The product of an even function and an odd function is an odd function.
