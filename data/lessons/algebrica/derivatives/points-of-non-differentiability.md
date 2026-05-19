> Content sourced from [Algebrica](https://algebrica.org/points-of-non-differentiability/) — CC BY-NC 4.0

## What are non differentiable points

In the entry on [derivatives](<../derivatives>), we saw that if a function \\( f(x) \\) is differentiable at a point \\( c \\), then the function is [continuous](<../continuous-functions/>) at that point. However, there are cases where a function is continuous at \\( c \\) but not differentiable. More generally, the **non-differentiable points** of a function \\( f(x) \\) occur when:

  * The right-hand and left-hand [limits](<../limits/>) of the [difference quotient](<../difference-quotient>) exist and are finite but are not equal. \\[f_{-}’ \left( c \right) \neq f_{+}’ \left( c \right)\\]

  * The limit of the difference quotient is infinite.


These points are categorized into three main types, which we will discuss below.

## Inflection point with vertical tangent

An **inflection point** is a point where the concavity of a function changes. In this case, we have a point of non-differentiability \\( c \\) of the function, which results in an [inflection point](https://algebrica.org/maximum-minimum-and-inflection-points/) with a tangent parallel to the \\( y \\)-axis (a vertical tangent). At such a point, the following occurs:

![](https://algebrica.org/wp-content/uploads/resources/images/non-differentiable-points-1.png)

This behavior indicates that the slope of the tangent becomes vertical at \\( x = c \\) while the function may change concavity around this point. In the case shown in the figure, we have \\[f_{-}’ \left (c \right) = f_{+}’ \left(c \right) = +\infty \\]

If the curve were reflected across the y-axis, we would have \\[f_{-}’ \left (c \right) = f_{+}’ \left(c \right) = -\infty \\]

## Cusps

In the case of **cusps** , the right-hand and left-hand limits are infinite and have opposite signs.

![](https://algebrica.org/wp-content/uploads/resources/images/non-differentiable-points-2.png)

In the case shown in the figure, we have: \\[f_{-}’ \left( c \right) = -\infty \quad \text{and} \quad f_{+}’ \left( c \right) = +\infty \\]

If the cusp were facing upwards instead of downwards, we would have: \\[f_{-}’ \left( c \right) = +\infty \quad \text{and} \quad f_{+}’ \left( c \right) = -\infty \\]

## Corners

A **corner** occurs when the left-hand derivative and the right-hand derivative exist but are not equal. In the case of corners points, there are two tangents to the graph at the same point, and they are different from each other.

![](https://algebrica.org/wp-content/uploads/resources/images/non-differentiable-points-3.png)

In this case we have:

\\[f_{-}’ \left (c \right) \neq f_{+}’ \left(c \right) \\]

How can we verify the differentiability of a function without relying on the limit of its difference quotient?

In general, let \\( f(x) \\) be a function continuous on an interval ([a,b]) and differentiable on that interval, except possibly at the point \\( x_0 \in [a,b] \\). If the limits \\(\lim_{x \to x_0^-} f’(x) \\) and \\( \lim_{x \to x_0^+} f’(x)\\) exist, then:

\\[f_{-}’ (x_o) = \lim_{x \to x_0^-} f’(x) \quad \text{and} \quad f_{+}’ (x_o) = \lim_{x \to x_0^+} f’(x) \\]

if \\( \underset{x \to x_0^-}{\lim} f{\prime}(x) = \underset{x \to x_0^+}{\lim} f{\prime}(x) = \ell\\), with \\(\ell \in \mathbb{R}\\) then the function is differentiable at \\(x_0\\), and it follows that \\(f’(x_0) = \ell\\).

Derivatives

The derivative describes a function’s rate of change.

9.1k

[Difference Quotient](https://algebrica.org/difference-quotient/)

7.8k

[Derivatives](https://algebrica.org/derivatives/)

1.6k

[Derivative of a Composite Function](https://algebrica.org/the-derivative-of-a-composite-function/)

1k

[Differential of a Function](https://algebrica.org/differential-of-a-function/)

1.4k

[Derivative of Composite Power Functions](https://algebrica.org/derivative-of-composite-power-functions/)

7.4k

[Maximum, Minimum, and Inflection Points](https://algebrica.org/maximum-minimum-and-inflection-points/)

1.4k

[Partial Derivatives](https://algebrica.org/partial-derivatives/)
