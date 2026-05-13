# Definite integrals

Source: algebrica.org — CC BY-NC 4.0
https://algebrica.org/definite-integrals/

## Area under a function: from curve to integral

Consider a [function](../functions/) \\( f(x) \\) defined on a [closed interval](../intervals/) \\( [a, b] \\). The definite integral of \\( f(x) \\) over this interval represents the oriented area of the region bounded by the graph of the function, the x-axis, and the vertical lines \\( x = a \\) and \\( x = b \\):

\\[\int_{a}^{b} f(x) \\, dx\\]

When \\( f(x) \\) is [continuous](../continuous-functions/) on \\( [a, b] \\), this region is called a curvilinear trapezoid: a planar figure bounded above by the graph of \\( f(x) \\), below by the x-axis, and laterally by the vertical lines \\( x = a \\) and \\( x = b \\).

![Img. 1](svg/definite-integrals-1.svg)

Its area cannot be determined by the standard formulas of elementary geometry, since one of its sides is a curve rather than a straight segment.

- - -
## Approximating area with rectangular sums

The area of the curvilinear trapezoid can be approximated by dividing the interval \\( [a, b] \\) into \\( n \\) subintervals of equal width:

\\[\Delta x = \frac{b - a}{n}\\]

Over each subinterval, the region is approximated by a rectangle, and the sum of the areas of these rectangles provides an estimate of the total area.

![Img. 2](svg/definite-integrals-2.svg)

Denoting by \\( m_i \\) and \\( M_i \\) respectively the minimum and maximum values of \\( f(x) \\) on the \\( i \\)-th subinterval, the lower and upper sums are defined as:

\\[s_n^{-} = \sum_{i=1}^{n} m_i \\,\Delta x_i \qquad s_n^{+} = \sum_{i=1}^{n} M_i \\,\Delta x_i\\]

The lower sum \\( s_n^{-} \\) approximates the area from below, the upper sum \\( s_n^{+} \\) from above.

![Img. 3](svg/definite-integrals-3.svg)

As the width \\( \Delta x \\) tends to zero, both sums converge to the same [limit](../limits/). A bounded function \\( f(x) \\) on \\( [a, b] \\) is said to be integrable if:

\\[\lim_{\Delta x \to 0} s_n^{-} = \lim_{\Delta x \to 0} s_n^{+} = I\\]

This common limit is the definite integral of \\( f(x) \\) over \\( [a, b] \\):

\\[I = \int_{a}^{b} f(x) \\, dx\\]

This approach, based on the convergence of lower and upper sums, is known as the [Riemann definition](../riemann-integrability-criteria/) of the integral. The values \\( a \\) and \\( b \\) are called the lower and upper limits of integration, and \\( f(x) \\) is the integrand. Geometrically, \\( f(x) \\) and \\( dx \\) represent the height and the base of the infinitesimal rectangles whose areas are summed in the limit.

- - -
## Computing definite integrals

If \\( f(x) \\) is a continuous function on \\( [a, b] \\) and \\( F(x) \\) is any [antiderivative](../indefinite-integrals/) of \\( f(x) \\), then:

\\[\int_{a}^{b} f(x) \\, dx = F(b) - F(a)\\]

- \\( F(x) \\) satisfies \\( F'(x) = f(x) \\) for all \\( x \in [a, b] \\).
- \\( F(b) \\) and \\( F(a) \\) are the values of the antiderivative evaluated at the upper and lower limits of integration, respectively.

This formula is the conclusion of the Second Fundamental Theorem of Calculus. The First Fundamental Theorem establishes that the function defined by accumulating area from a fixed point is differentiable, with derivative equal to the original integrand. Specifically, if:

\\[F(x) = \int_a^x f(t) \\, dt\\]

then \\( F'(x) = f(x) \\), making differentiation and integration inverse operations in a precise sense. Both results are covered in the dedicated page on the [Fundamental Theorem of Calculus](../fundamental-theorem-of-calculus/).

- - -
## Properties

When the two extremes of integration coincide, the integral is zero:

\\[\int_{a}^{a} f(x) \\, dx = 0\\]

This follows directly from the definition: if the interval has no width, there is no area to accumulate. Reversing the limits of integration changes the sign of the integral:

\\[\int_{a}^{b} f(x) \\, dx = -\int_{b}^{a} f(x) \\, dx\\]

This reflects the oriented nature of the definite integral: traversing the interval in the opposite direction reverses the sign of the accumulated area.

If \\( f(x) = k \\) is constant on \\( [a, b] \\), the region under the graph is a rectangle of height \\( k \\) and base \\( b - a \\):

\\[\int_{a}^{b} k \\, dx = k(b - a)\\]

If \\( k \\) is a constant, a constant factor can be moved outside the integral:

\\[\int_{a}^{b} k \cdot f(x) \\, dx = k \int_{a}^{b} f(x) \\, dx\\]

- - -

The integral is additive over sums of functions:

\\[\int_{a}^{b} \left( f(x) + g(x) \right) \\, dx = \int_{a}^{b} f(x) \\, dx + \int_{a}^{b} g(x) \\, dx\\]

Together, these two properties make the definite integral a linear operator. The integral is also additive over adjacent intervals: if \\( a \\), \\( b \\), and \\( c \\) are points in the domain of \\( f \\), then:

\\[\int_{a}^{c} f(x) \\, dx = \int_{a}^{b} f(x) \\, dx + \int_{b}^{c} f(x) \\, dx\\]

If \\( f(x) \leq g(x) \\) for every \\( x \in [a, b] \\), then:

\\[\int_{a}^{b} f(x) \\, dx \leq \int_{a}^{b} g(x) \\, dx\\]

> This is known as the comparison property of integrals. Geometrically, it states that the area under \\( f \\) does not exceed the area under \\( g \\) over the same interval.

- - -
## Mean Value Theorem for Integrals

If \\( f(x) \\) is continuous on \\( [a, b] \\), then there exists at least one point \\( c \in (a, b) \\) such that:

\\[\int_{a}^{b} f(x) \\, dx = f\(c\)(b - a)\\]

The value \\( f\(c\) \\) is the average value of the function over the interval. Geometrically, there exists a rectangle with base \\( b - a \\) and height \\( f\(c\) \\) whose area equals exactly the area under the curve. The theorem guarantees the existence of such a point without providing a method to locate it. The average value of \\( f \\) over \\( [a, b] \\) can be written explicitly as:

\\[f\(c\) = \frac{1}{b - a} \int_{a}^{b} f(x) \\, dx\\]

> The Mean Value Theorem for Integrals is the integral counterpart of [Lagrange's Mean Value Theorem](../lagrange-theorem/). While the latter guarantees a point where the instantaneous rate of change equals the average rate of change, this theorem guarantees a point where the function value equals the average value over the interval.

- - -
## Example 1

Compute the following definite integral:

\\[\int_{0}^{3} (3x - x^2) \\, dx\\]

Applying linearity and moving the constant factor outside the first integral gives:

\\[3\int_{0}^{3} x \\, dx - \int_{0}^{3} x^2 \\, dx\\]

Computing the antiderivative of each term:

\\[F(x) = \frac{3x^2}{2} - \frac{x^3}{3}\\]

Evaluating \\( F(3) - F(0) \\):

\\[\begin{align}
F(3) - F(0) &= \left( \frac{3 \cdot 9}{2} - \frac{27}{3} \right) - \left( \frac{3 \cdot 0}{2} - \frac{0}{3} \right) \\\\[6pt]
&= \frac{27}{2} - 9 \\\\[6pt]
&= \frac{27 - 18}{2} \\\\[6pt]
&= \frac{9}{2}
\end{align}\\]

> For a systematic treatment of antiderivatives and their properties, see the page on [indefinite integrals](../indefinite-integrals/).

The area of the region bounded by the graph of \\( f(x) = 3x - x^2 \\) and the x-axis over \\( [0, 3] \\) is equal to: 

\\[ \dfrac{9}{2} \\]

> This is just a simple example that generally shows the procedure for calculating definite integrals. Very often, integrals are not so straightforward to compute, and it is necessary to resort to other solving methods such as [substitution](../integration-by-substitution/) and [integration by parts](../integration-by-parts/).

- - -
## Example 2

Compute the following definite integral:

\\[\int_{0}^{\pi} (x + \sin x) \\, dx\\]

> This example combines a polynomial term with a trigonometric function. For a review of the relevant antiderivatives, see the page on [integrals of trigonometric functions](../integral-of-trigonometric-functions/).

- - -

Applying linearity:

\\[\int_{0}^{\pi} x \\, dx + \int_{0}^{\pi} \sin x \\, dx\\]

Computing the antiderivative of each term:

\\[F(x) = \frac{x^2}{2} - \cos x\\]

Evaluating \\( F(\pi) - F(0) \\):

\\[\begin{align}
F(\pi) - F(0) &= \left( \frac{\pi^2}{2} - \cos\pi \right) - \left( \frac{0}{2} - \cos 0 \right) \\\\[6pt]
&= \left( \frac{\pi^2}{2} + 1 \right) - (0 - 1) \\\\[6pt]
&= \frac{\pi^2}{2} + 2
\end{align}\\]

The area of the region bounded by the graph of \\( f(x) = x + \sin x \\) and the x-axis over \\( [0, \pi] \\) is equal to:

\\[
\frac{\pi^2}{2} + 2
\\]

- - -
## Handling definite integrals with positive and negative areas

The interpretation of the definite integral as an area holds when \\( f(x) \geq 0 \\) throughout \\( [a, b] \\). When \\( f(x) \\) changes sign within the interval, the integral assigns a negative value to the portions of the region lying below the x-axis, so the result is an oriented area rather than a geometric one.

![Img. 4](svg/definite-integrals-4.svg)

The interval \\( [a, b] \\) must be divided into subintervals over which \\( f(x) \\) maintains constant sign. If \\( f(x) > 0 \\) on \\( [a, c] \\) and \\( f(x) < 0 \\) on \\( [c, b] \\), the integral over each subinterval is computed separately and the results combined algebraically:

\\[\int_{a}^{b} f(x) \\, dx = \int_{a}^{c} f(x) \\, dx + \int_{c}^{b} f(x) \\, dx\\]

For an [even function](../even-and-odd-functions/), symmetry about the y-axis implies that the contributions from \\( [-a, 0] \\) and \\( [0, a] \\) are equal, so:

\\[\int_{-a}^{a} f(x) \\, dx = 2\int_{0}^{a} f(x) \\, dx\\]

![Img. 5](svg/definite-integrals-5.svg)

For an [odd function](../even-and-odd-functions/), symmetry about the origin implies that the contributions from \\( [-a, 0] \\) and \\( [0, a] \\) are equal in magnitude but opposite in sign, so:

\\[\int_{-a}^{a} f(x) \\, dx = 0\\]

![Img. 6](svg/definite-integrals-6.svg)

In both cases, the geometric area enclosed between the graph of \\( f(x) \\) and the x-axis over \\( [-a, a] \\) is obtained by integrating the absolute value of the function. For an even function this coincides with the oriented integral:

\\[S = 2\int_{0}^{a} f(x) \\, dx\\]

For an odd function, where the oriented integral vanishes, the geometric area is:

\\[S = 2\int_{0}^{a} |f(x)| \\, dx\\]

- - -
## Improper integrals

Everything covered on this page assumes that the interval \\( [a, b] \\) is finite and that \\( f(x) \\) remains bounded throughout. These conditions are not always satisfied: it is common to encounter integrals over unbounded intervals, or functions that diverge at some point in the domain.

The standard [Riemann integral](../riemann-integrability-criteria/) cannot handle these cases directly. The approach is to replace the problematic bound with a parameter and take a limit. An integral over an unbounded interval is defined as:

\\[\int_{a}^{+\infty} f(x) \\, dx = \lim_{t \to +\infty} \int_{a}^{t} f(x) \\, dx\\]

When the limit exists and is finite, the integral converges; otherwise it diverges. These cases are treated in the dedicated page on [improper integrals](../improper-integrals/).

> The definite integral, as developed on this page, measures signed area and provides a precise way to compute accumulated quantities. One of the most direct applications is finding the area of regions bounded by curves, a problem that reduces entirely to setting up and evaluating definite integrals of the kind studied here. This is covered in detail in the dedicated page on [finding areas by integration](../finding-areas-by-integration/).