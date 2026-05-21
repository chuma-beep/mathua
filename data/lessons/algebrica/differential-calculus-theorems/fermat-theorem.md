> Content sourced from [Algebrica](https://algebrica.org/fermat-theorem/) — CC BY-NC 4.0

## Introduction

Fermat’s Theorem states that any relative [maximum or minimum](<../maximum-minimum-and-inflection-points/>) of a differentiable function within its [domain](<../determining-the-domain-of-a-function/>) must occur at a stationary point, that is, a point where the first [derivative](<../derivatives/>) is equal to zero, and the tangent line is horizontal (parallel to the \\(x\\)-axis).

Fermat’s Theorem states a necessary, but not sufficient, condition for identifying local extrema. Specifically, if a function has a local maximum or minimum at a point and is differentiable there, then its derivative at that point must be zero. However, the converse does not hold: a zero derivative does not necessarily imply the presence of an extremum.

###### Fermat’s Theorem plays a key role in the proof of [Rolle’s Theorem](<../rolles-theorem/>).

## Statement

Given a [function](<../functions/>) \\( y = f(x) \\) defined on a closed and bounded interval \\([a, b]\\), and differentiable on the open interval \\((a, b)\\), if the function attains a local maximum or minimum at a point \\( x_0 \in (a, b) \\), then the derivative at that point must be zero:

\\[f’(x_0) = 0 \\]

The following graph illustrates theorem. At a local maximum point \\( \mu \\) within the interval \\( I \\), the derivative of the function is equal to zero.

![](/diagrams/algebrica/fermat-1.png)

###### In certain cases, the derivative of a function becomes zero at a point that is neither a maximum nor a minimum. Such points are referred to as stationary points without an extremum, including stationary [inflection points](<../maximum-minimum-and-inflection-points/>), where the derivative is zero but the function maintains its direction.

## Proof

To prove the theorem, let us assume that \\( x_0 \\) is a point of local maximum. Then, in a neighborhood \\( I \\) around \\( x_0 \\), the following inequality must hold:

\\[f(x) \leq f(x_0) \quad \forall \, x \in I \\] From this, it follows that the [difference quotient](<../difference-quotient>) satisfies the following:

For \\( h > 0 \\): \\[\frac{f(x_0 + h) - f(x_0)}{h} \leq 0 \\]

For \\( h < 0 \\): \\[\frac{f(x_0 + h) - f(x_0)}{h} \geq 0 \\]

From these inequalities, and by the definition of the derivative as the limit of the difference quotient, it follows that the respective limits satisfy:

\\[\lim_{h \to 0^+} \frac{f(x_0 + h) - f(x_0)}{h} \leq 0 \\] \\[\lim_{h \to 0^-} \frac{f(x_0 + h) - f(x_0)}{h} \geq 0 \\]

If the function is differentiable at \\( x_0\\), then both the left-hand and right-hand limits exist and are equal to the derivative. The only way these two inequalities can be true simultaneously is if:

\\[f’(x_0) = 0 \\]

##### The theorem is thus proven, as we have shown that if a differentiable function attains a local extremum at an interior point, the derivative at that point must necessarily be zero.

Fermat’s Theorem states that if a function attains a local maximum or minimum at a point and is differentiable at that point, then its derivative must be zero. However, the converse does not hold: a zero derivative at a point does not necessarily imply that the point is a local maximum or minimum.

## Example 1

Let’s explore an example that applies Fermat’s Theorem. Suppose we study the real-valued function:

\\[f(x) = x^{3} - 3x^{2} + 2 \\]

defined for every real number. Being a [polynomial](<../polynomials/>), the function is [continuous](<../continuous-functions/>) and [differentiable](<../derivatives/>) on the entire real line. This ensures that if it attains a local maximum or minimum at some interior point of its domain, Fermat’s Theorem guarantees that the derivative at that point must be equal to zero. To explore where the function might have extrema, we compute its derivative:

\\[f’(x) = 3x^{2} - 6x = 3x(x - 2) \\]

The derivative vanishes precisely when \\(3x(x - 2) = 0\\) which occurs at:

\\[x = 0 \quad \text{and} \quad x = 2 \\]

These two values are therefore the only candidates for interior extrema, because Fermat’s Theorem states that any differentiable function reaching a local extremum must have a horizontal tangent line at that point.


To determine the nature of these points, we examine how the derivative behaves around them. For \\( x < 0 \\), the derivative is positive and the function increases. Between \\( 0 \\) and \\( 2 \\), the derivative becomes negative, so the function decreases. For \\( x > 2 \\), the derivative returns to positive, meaning the function increases again. This change in [monotonicity](<../increasing-and-decreasing-functions/>) reveals that:

  * for \\( x = 0 \\), the function transitions from increasing to decreasing, indicating a local maximum;
  * for \\( x = 2 \\), the function transitions from decreasing to increasing, indicating a local minimum.


| | \\[0\\]| \\[2\\]  
---|---|---|---  
\\( f'(x) \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( f(x) \\)| \\( \boldsymbol{\nearrow} \\)| \\( \boldsymbol{\searrow} \\)| \\( \boldsymbol{\nearrow} \\)  
  
Evaluating the function confirms this classification: \\(f(0) = 2 \\) is a local maximum, and \\(f(2)= -2 \\) is a local minimum.

###### This example highlights the essential role of Fermat’s Theorem: it does not ensure that every point with a zero derivative is an extremum, but it does guarantee that every interior extremum of a differentiable function must occur at such a point. It therefore provides a necessary condition that guides and narrows the search for maxima and minima.

## Not all stationary points are extrema

The function \\( f(x) = x^3 \\) provides a clear example showing that a zero derivative does not necessarily imply the presence of a local extremum.

![](/diagrams/algebrica/fermat-2.png)

The derivative is:

\\[f’(x) = 3x^2 \\]

At \\( x = 0 \\), we have \\(f’(0) = 0\\). Thus, \\( x = 0 \\) is a stationary point. Although the derivative at \\( x = 0 \\) is zero, the point is neither a local maximum nor a local minimum. Instead, x = 0 represents a stationary inflection point, where the derivative equals zero. However, the function increases on both sides of this point without changing direction.

## Stationary points and boundary behavior

The presence of a stationary point—namely, a point at which the derivative vanishes—does not by itself determine whether the function achieves a local maximum or minimum there. The classification of such points requires a more refined analysis of the function in a neighborhood of the candidate point, as well as an understanding of the structure of the domain. Consider, for instance, the [absolute value function](https://algebrica.org/absolute-value-function/)

\\[y = f(x) = |x| = \begin{cases} +x & \text{if } x \ge 0\\\\[4pt] -x & \text{if } x < 0 \end{cases} \\]

which attains a global minimum at \\(x = 0\\) even though the derivative does not exist at that point.

![The function | x | has an absolute minimum at the origin, but it is not differentiable there, so Fermat’s theorem does not apply.](/diagrams/algebrica/fermath-theorem-4.png)

This example illustrates that extrema may occur not only where the derivative is zero, but also where differentiability fails or at boundary points of the domain. When dealing with closed intervals, a function may reach its extrema at the endpoints regardless of the behavior of the derivative in the interior.

Conversely, on open intervals or unbounded [domains](<../determining-the-domain-of-a-function/>), extrema may be absent even if stationary points exist. A complete understanding of maxima and minima therefore requires the combined use of Fermat’s condition, the study of differentiability, and the careful analysis of the domain on which the function is defined.

## Selected references

  * **Harvard University O. Knill**. [The Mean Value Theorem and Rolle’s Theorem](https://people.math.harvard.edu/~knill/teaching/math1a_2014/handouts/26-rolle.pdf)

  * **UC Davis J. Hunter**. [Differentiable Functions](https://www.math.ucdavis.edu/~hunter/m125a/intro_analysis_ch4.pdf)

  * **MIT J. Starr**. [Critical Points and Extrema](https://ocw.mit.edu/courses/18-01-single-variable-calculus-fall-2005/e623940ed32c0dc3f942b64b9db07485_lecture_9.pdf)
