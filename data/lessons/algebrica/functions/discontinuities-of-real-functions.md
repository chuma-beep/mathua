> Content sourced from [Algebrica](https://algebrica.org/discontinuities-of-real-functions/) — CC BY-NC 4.0

## Introduction

Continuity is a property of a [function](<../functions/>) in which small variations in the input result in correspondingly small variations in the output within the neighbourhood of a given point. If this local stability does not hold, the function is considered discontinuous. Discontinuities are typically classified into three distinct types:

  * A removable discontinuity occurs when the limit exists and is finite, but the function is either undefined at the point or its value does not equal the limit.
  * A jump discontinuity is present when both the left-hand and right-hand limits exist and are finite, but these limits are not equal.
  * An infinite discontinuity occurs when at least one of the one-sided limits is infinite, causing the function to diverge near the point rather than approach a finite value.


A discontinuity at \\(x_0\\) can occur in exactly one of the three mutually exclusive ways described above. A point cannot simultaneously exhibit more than one type of discontinuity.

###### Each of these types will be examined in detail in the following sections. In this discussion, \\( f \\) denotes a real-valued function, and \\( x_0 \\) represents a point in its [domain](<../determining-the-domain-of-a-function/>) or a point at which the function may fail to be defined.

## Recall of continuity

A function \\( f \\) is [continuous](<../continuous-functions/>) at \\( x_0 \\) if the limit as \\( x \\) approaches \\( x_0 \\) exists, is finite, and coincides with the value of the function at that point. This condition is expressed by the following [limit](<../limits/>):

\\[\lim_{x \to x_0} f(x) = f(x_0) \\]

[Polynomials](<../polynomials/>) constitute a fundamental class of elementary continuous functions. These functions represent smooth curves in the plane and exhibit no points of discontinuity. Below is the graph of the quadratic function \\( x^2 + 2x + 1 \\), which represents a [parabola](<../parabola/>):

![The graph of a second-degree polynomial is a continuous parabola, with no jumps or interruptions.](https://algebrica.org/wp-content/uploads/resources/images/discontinuity.png)

A discontinuity at \\( x_0 \\) arises whenever this equality does not hold, and the specific way in which the condition fails determines the type of discontinuity.

###### Intuitively, a function is continuous if its graph can be drawn in the plane without any interruptions, breaks, or sudden jumps.

## Removable discontinuity

A removable discontinuity arises when a function possesses a well-defined finite limit at \\( x_0 \\), yet the function’s value at that point is either undefined or does not coincide with the limit. Formally, a function \\( f \\) has a removable discontinuity at \\( x_0 \\) if the following limit exists and is finite:

\\[\lim_{x \to x_0} f(x) = \ell \in \mathbb{R}\\]

Moreover, at least one of the following conditions is satisfied:

  * \\( f(x_0) \\) is undefined.
  * \\( f(x_0) \neq \ell \\).


In such cases, the discontinuity may be removed by redefining the function at a single point as follows: \\[g(x) = \begin{cases} \ell & \text{if } x = x_0 \\\\[6pt] f(x) & \text{if } x \ne x_0 \end{cases} \\]

With this definition, the function \\( g \\) becomes continuous at \\( x_0 \\). The term “removable” refers to the fact that the discontinuity can be resolved in this manner.

###### Removable discontinuities typically occur in rational functions containing cancellable factors, resulting in a hole in the graph. They can also be present in piecewise-defined functions or in functions where the value at a single point has been modified, provided the limit at that point exists and is finite.

## Example 1

Consider the function defined by the following [rational](<../rational-functions/>) expression, which is undefined at \\( x = 1 \\):

\\[f(x) = \frac{x^2 - 1}{x - 1} \\]

[Factoring](<../factoring-ac-method/>) the numerator demonstrates that the expression simplifies for all values of \\( x \\) except \\(1\\) since \\(x=1\\) would cancel the denominator and make the function undefined.

\\[x^2 - 1 = (x - 1)(x + 1) \\]

For all \\( x \neq 1 \\), the function is equivalent to a linear function: \\[f(x) = x + 1 \\]

Although the function is undefined at \\( x = 1 \\), the limit as \\( x \\) approaches \\(1\\) exists and is finite:

\\[\lim_{x \to 1} \frac{x^2 - 1}{x - 1} = 2 \\]

This demonstrates that \\( x = 1 \\) is a removable discontinuity, as the graph corresponds to the straight line \\( y = x + 1 \\) with a single missing point at\\( (1,2) .\\)

Redefining the function at that point by assigning it the value of the limit eliminates the discontinuity: \\[g(x) = \begin{cases} 2 & \text{if } x = 1 \\\\[6pt] f(x) & \text{if } x \ne 1 \end{cases} \\]

###### With this modification, the function is continuous at \\( x = 1 \\).

## Jump discontinuity

A jump discontinuity arises when both the left-hand and right-hand limits at \\( x_0 \\) exist and are finite, yet these limits are not equal. Formally, \\( f \\) has a jump discontinuity at \\( x_0 \\) if:

\\[\begin{align} \lim_{x \to x_0^-} f(x) &= \ell_1 \in \mathbb{R} \\\\[6pt] \lim_{x \to x_0^+} f(x) &= \ell_2 \in \mathbb{R} \\\\[6pt] \ell_1 &\neq \ell_2 \end{align} \\]

In this case, the limit \\( \lim_{x \to x_0} f(x) \\) does not exist. The function approaches two distinct finite values depending on the direction of approach. Unlike a removable discontinuity, this type cannot be resolved by redefining the function at a single point, as the discrepancy is inherent to the local behavior.

## Example 2

To analyse the jump discontinuity, consider the following simple function, which exhibits a discontinuity at the point \\(x = 1.\\)

\\[f(x) = \begin{cases} 0 & \text{if } x < 1 \\\\[6pt] 2 & \text{if } x \ge 1 \end{cases} \\]

![](https://algebrica.org/wp-content/uploads/resources/images/discontinuity-2.png)

For values of \\(x\\) approaching \\(1\\) from the left, the function remains constant at \\(0\\). Therefore:

\\[\lim_{x \to 1^-} f(x) = 0 \\]

For values of \\( x \\) approaching 1 from the right, the function remains constantly equal to \\(2\\), and therefore the limit is: \\[\lim_{x \to 1^+} f(x) = 2 \\]

Both one-sided limits exist and are finite but they are not equal. Since \\( 0 \neq 2 \\), it follows that the two one-sided limits do not coincide, and consequently, the limit \\(\lim_{x \to 1} f(x)\\) does not exist. The graph of the function shows a vertical jump at \\(x = 1\\), transitioning from \\(0\\) to \\(2\\).

This discontinuity cannot be removed by redefining the function at \\(x = 1\\), as the difference between the two limiting values indicates a break in the local behaviour of the function.

## Infinite discontinuity

An infinite discontinuity occurs when a function diverges as \\( x \\) approaches \\( x_0 \\), with at least one of the one-sided limits being infinite. Formally, a function \\( f \\) exhibits an infinite discontinuity at \\( x_0 \\) if at least one of the following conditions is satisfied:

\\[\begin{align} \lim_{x \to x_0^-} f(x) &= \pm \infty \\\\[6pt] \lim_{x \to x_0^+} f(x) &= \pm \infty \end{align} \\]

In such cases, the function does not approach any finite value as \\( x \\) nears \\( x_0 \\). The graph typically displays a vertical [asymptote](<../asymptotes/>). This discontinuity reflects unbounded growth rather than a finite discontinuity.

## Example 3

For example, consider the following function:

\\[f(x) = \frac{1}{x - 2} \\]

![](https://algebrica.org/wp-content/uploads/resources/images/discontinuity-3.png)

The behaviour of this function near \\( x_0 = 2 \\) is analysed as follows. As \\( x \to 2^- \\), the denominator \\( x - 2 \\) becomes negative and approaches zero, causing the function to decrease without bound. Therefore we have:

\\[\lim_{x \to 2^-} f(x) = -\infty \\]

As \\( x \to 2^+ \\), the denominator is positive and approaches zero, which causes the function to increase without bound. The limit is:

\\[\lim_{x \to 2^+} f(x) = +\infty \\]

At least one of the one-sided limits is infinite, and they diverge with opposite signs. Consequently, the function exhibits an infinite discontinuity at \\( x = 2 \\). The graph displays a [vertical asymptote](<../asymptotes/>) at the line \\( x = 2 \\), and this divergence indicates unbounded growth rather than a finite jump or a removable discontinuity.

## Discontinuity, continuity and differentiability

It is instructive to establish a precise link between the notions of discontinuity and [differentiability](<../derivatives/>). We know that if a function f is differentiable at a point \\(x_0\\), it must also be [continuous](<../continuous-functions/>) at that point. The existence of the derivative ensures that the function satisfies the condition of continuity:

\\[f’(x_0) = \lim_{x \to x_0} \frac{f(x) - f(x_0)}{x - x_0} \\] \\[\lim_{x \to x_0} f(x) = f(x_0) \\]

Therefore, if a function exhibits a discontinuity of the type just described at \\( x_0 \\), meaning the limit does not exist or does not equal the function’s value the derivative at that point does not exist.

However, the converse is not true. A function can be continuous at \\(x_0\\) yet not differentiable there. This situation arises when the one-sided derivatives exist but differ, when at least one is infinite, or when one of the limits diverges.

\\[\lim_{x \to x_0^-} \frac{f(x) - f(x_0)}{x - x_0} \neq \lim_{x \to x_0^+} \frac{f(x) - f(x_0)}{x - x_0} \\]

A common example is the [absolute value function](<../absolute-value-function/>), which is continuous at \\(x = 0\\) but not differentiable there, resulting in a corner on its graph. In summary, every discontinuity implies non-differentiability, whereas not every point of non-differentiability is associated with a discontinuity.

## A particular case: essential discontinuity

An additional category, known as essential discontinuity, is sometimes recognised but not universally adopted as a formal classification. This type arises when the limit does not exist and cannot be described as infinite. Unlike a jump discontinuity, where both one-sided limits exist but are unequal, or an infinite discontinuity, where the function diverges in a particular direction, an essential discontinuity reflects fundamentally irregular behaviour that cannot be reduced to simpler forms.

A classic example is the following function, which exhibits an essential discontinuity at \\(x = 0\\):

\\[f(x) = \sin\left(\frac{1}{x}\right) \\]

As \\(x\\) approaches \\(0\\), the argument \\(1/x\\) grows without bound, causing the function to oscillate between \\(-1\\) and \\(1\\) with increasing frequency. Neither one-sided limit exists, and no value can be assigned to \\(f(0)\\) that would restore any form of continuity.

## Selected references

  * **Harvard University N. Masson**. [Discontinuities and Monotonic Functions](https://people.math.harvard.edu/~nate/teaching/UPenn/2009/spring/math_360/lectures/week_8/lecture_14/lecture_14.pdf)

  * **Harvard University O. Knill**. [Introduction to Functions and Calculus – Continuity](https://people.math.harvard.edu/~knill/teaching/math1a_2012/handouts/14-continuity.pdf)

  * **UC Berkeley A. Vizeff**. [Lecture 5: Continuity and Discontinuities](https://math.berkeley.edu/~avizeff/calculus-I-F22/lecture-5.pdf)


Functions

A function maps each input to a unique output.

17.6k

[Functions](https://algebrica.org/functions/)

3 comments

10.2k

[Determining the Domain of a Function](https://algebrica.org/determining-the-domain-of-a-function/)

1.6k

[Even and Odd Functions](https://algebrica.org/even-and-odd-functions/)

2.7k

[Increasing, Decreasing and Monotonic Functions](https://algebrica.org/increasing-and-decreasing-functions/)

1.9k

[Convexity and Concavity of Functions](https://algebrica.org/convexity-and-concavity-of-functions/)

1.3k

[Composite Functions](https://algebrica.org/composite-functions/)

1.2k

[Inverse Function](https://algebrica.org/inverse-function/)

1.9k

[Continuous Functions](https://algebrica.org/continuous-functions/)

2 comments

1.1k

[Uniform Continuity](https://algebrica.org/uniform-continuity/)

1.9k

[Analyzing the Graphs of Functions](https://algebrica.org/analyzing-the-graphs-of-functions/)

1.1k

[Polynomial Function](https://algebrica.org/polynomial-function/)

2.6k

[Rational Functions](https://algebrica.org/rational-functions/)

2k

[Logarithmic Function](https://algebrica.org/logarithmic-function/)

3.1k

[Exponential Function](https://algebrica.org/exponential-function/)

1.5k

[Absolute Value Function](https://algebrica.org/absolute-value-function/)

847

[Sign Function](https://algebrica.org/sign-function/)

2.2k

[Sine Function](https://algebrica.org/sine-function/)

2k

[Cosine Function](https://algebrica.org/cosine-function/)

1.7k

[Tangent Function](https://algebrica.org/tangent-function/)

2k

[Cotangent Function](https://algebrica.org/cotangent-function/)

1.8k

[Secant Function](https://algebrica.org/secant-function/)

1.4k

[Cosecant Function](https://algebrica.org/cosecant-function/)

1.2k

[Dirichlet Function](https://algebrica.org/dirichlet-function/)

1.4k

[Sigmoid Function](https://algebrica.org/sigmoid-function/)
