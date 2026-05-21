> Content sourced from [Algebrica](https://algebrica.org/weierstrass-theorem/) — CC BY-NC 4.0

## Statement

Let \\( f : [a,b] \to \mathbb{R} \\) be a continuous [function](<../functions/>) on a closed and bounded interval \\( [a,b] \\). Then there exist points \\( x_{\min}, x_{\max} \in [a,b] \\) such that:

\\[f(x_{\min}) \le f(x) \le f(x_{\max}) \quad \forall \; x \in [a,b] \\]

In other words, a [continuous function](<../continuous-functions/>) on a closed interval attains both its minimum and its maximum. This result is often referred to as the Extreme Value Theorem. The following graph illustrates the theorem: the function reaches a maximum and a minimum at two interior points of the interval \\( [a, b] \\).

![](/diagrams/algebrica/weierstrass-theorem-1.png)

It is helpful to relate the notation to the figure. The quantity \\( f(x_{\max}) \\) denotes the value of the function at the point marked Max in the graph, that is, the highest point reached on the interval. Likewise, \\( f(x_{\min}) \\) represents the value of the function at the point marked min, where the graph attains its lowest value.


Each assumption plays a precise role.

  * The interval must be closed and bounded. Closed means that the endpoints are included. Bounded means that it has finite length. Together, these conditions express compactness on the real line.

  * Continuity is equally essential. The theorem does not require differentiability. It only requires that the function have no jumps or breaks on the interval.

  * If either condition is removed, the conclusion may fail.


A full proof lies beyond this page, but the main idea is straightforward. On a closed and bounded interval, a continuous function cannot diverge or skip values. Because the interval includes its endpoints and continuity prevents jumps, the function does not merely approach its largest and smallest values; it actually attains them.

## Why closed and bounded matters

Consider the function \\( f(x) = x \\) on the open interval \\( (0,1) \\).

![](/diagrams/algebrica/weierstrass-theorem-2.png)

The function is continuous, yet it has no maximum and no minimum on that interval. The infimum is \\( 0 \\) and the supremum is \\( 1 \\), but neither value is attained because the endpoints are not included.


Now consider \\( f(x) = \dfrac{1}{x} \\) on \\( (0,1] \\). The interval is bounded but not closed. The function is continuous on its [domain](<../determining-the-domain-of-a-function/>), yet it does not attain a maximum. As \\( x \to 0^+ \\), the function grows without bound.

###### These examples show that compactness of the domain is not a technical detail but the structural reason the theorem holds.

## Example 1

Let us now see a concrete application of the theorem. Consider the function:

\\[f(x) = x^3 - 3x \\]

on the interval \\( [-2,2] \\). This is a [polynomial](<../polynomials/>), therefore continuous on the whole real line. In particular, it is continuous on the closed and bounded interval \\( [-2,2] \\). By Weierstrass’ theorem, we already know that the function must attain both a maximum and a minimum somewhere in this interval. To determine where these extreme values occur, we proceed using differential calculus. First compute the [derivative](<../derivatives/>):

\\[f’(x) = 3x^2 - 3 = 3(x^2 - 1) \\]

The critical points are obtained by solving \\( f’(x)=0 \\):

\\[3(x^2 - 1)=0 \quad \rightarrow \quad x^2=1 \\]

so we have:

\\[x=-1 \quad x=1 \\]

These are the only interior points where the slope of the tangent line vanishes. However, the absolute extrema on a closed interval may also occur at the endpoints. For this reason, we evaluate the function at all candidates: the critical points and the endpoints.

\\[\begin{array}{ll} f(-2) = -2 & f(-1) = 2 \\\\[6pt] f(1) = -2 & f(2) = 2 \end{array} \\]

Comparing these values, we observe that the maximum value is \\( 2 \\), attained at \\( x=-1 \\) and \\( x=2 \\), while the minimum value is \\( -2 \\), attained at \\( x=-2 \\) and \\( x=1 \\).

###### To summarize, one important point of this example is that Weierstrass’ theorem does not tell us where the extreme values are located, nor how many there are. It simply guarantees that they exist. It is the derivative that allows us to find them explicitly.

## The range of a continuous function

Let \\( f \\) be continuous on a closed and bounded interval \\( [a,b] \\). By the Weierstrass theorem, \\( f \\) attains a minimum value \\( m \\) and a maximum value \\( M \\) on \\( [a,b] \\). This means there exist points \\( x_{\min}, x_{\max} \in [a,b] \\) such that

\\[f(x_{\min}) = m \qquad f(x_{\max}) = M \\]

By the
