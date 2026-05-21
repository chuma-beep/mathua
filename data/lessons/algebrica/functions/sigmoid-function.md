> Content sourced from [Algebrica](https://algebrica.org/sigmoid-function/) — CC BY-NC 4.0

Concept

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Advanced

2

Requires

0

Enables

The following concepts, [Exponential Function](https://algebrica.org/exponential-function/), [Functions](https://algebrica.org/functions/), are required as prerequisites for this entry.

## Definition

The sigmoid function is a real-valued [function](<../functions/>) of a real variable that takes values strictly between \\(0\\) and \\(1\\), approaching each of the two extremes [asymptotically](<../asymptotes/>). It provides a smooth mapping from the real line to the unit interval and is widely used in analysis and machine learning. Its definition is the following:

\\[\sigma(x) = \frac{1}{1 + e^{-x}} \\]

The sigmoid function can be written in equivalent forms that are sometimes more convenient for computation or for establishing certain properties. One such form is obtained by multiplying both numerator and denominator by \\(e^x\\):

\\[\sigma(x) = \frac{e^x}{e^x + 1} \\]

This expression is entirely equivalent to the original definition and can simplify certain algebraic manipulations.

![Sigmoid function.](/diagrams/algebrica/sigmoid-function.png)

  * The [domain](<../determining-the-domain-of-a-function/>) is \\(\mathbb{R}\\), while the range is the open interval \\((0,1)\\).
  * The function is strictly increasing on \\(\mathbb{R}\\), since its [derivative](<../derivatives/>) is always positive, and is therefore bijective from \\(\mathbb{R}\\) onto \\((0,1)\\).
  * The function has no local extrema and exactly one [inflection point](<../maximum-minimum-and-inflection-points/>) at \\((0, \frac{1}{2})\\).
  * The [limits](<../limits/>) at infinity are the following: \\[\lim_{x \to -\infty} \sigma(x) = 0\\] \\[\lim_{x \to +\infty} \sigma(x) = 1\\]


> The S-shaped curve reflects the behaviour of the function: slow growth for very negative values of \\(x\\), a rapid transition in a neighbourhood of the origin, and saturation for very positive values.

## Properties of the sigmoid function

The following properties characterise the sigmoid function analytically and justify its widespread use in both mathematical analysis and applications. The function satisfies the following symmetry relation with respect to the origin:

\\[\sigma(-x) = 1 - \sigma(x) \\]

This identity, which can be verified by direct substitution, implies that the graph of \\(\sigma\\) is symmetric about the point:

\\[\left(0,\, \frac{1}{2}\right)\\]

The value of the function at the origin is:

\\[\sigma(0) = \frac{1}{1 + e^{0}} = \frac{1}{2} \\]

The limits at the extremes of the real line are the following:

\\[\lim_{x \to -\infty} \sigma(x) = 0\\] \\[\lim_{x \to +\infty} \sigma(x) = 1\\]

The lines \\(y = 0\\) and \\(y = 1\\) are therefore [horizontal asymptotes](<../asymptotes/>) of the graph.

## Derivative of the sigmoid function

One of the most notable properties of the sigmoid function is that its derivative can be expressed in a remarkably compact form in terms of the function itself. The first derivative is the following:

\\[\sigma’(x) = \sigma(x)\,\bigl(1 - \sigma(x)\bigr) \\]

To verify this identity one may proceed by direct computation. Writing \\(\sigma(x) = (1 + e^{-x})^{-1}\\) and applying the chain rule gives the following:

\\[\sigma’(x) = \frac{e^{-x}}{(1 + e^{-x})^2} \\]

Observing that the numerator can be written as \\((1 + e^{-x}) - 1\\), the expression separates into the product:

\\[\begin{align} \sigma’(x) &= \frac{1}{1 + e^{-x}} \cdot \frac{e^{-x}}{1 + e^{-x}} \\\\[6pt] &= \sigma(x)\,\bigl(1 - \sigma(x)\bigr) \end{align} \\]

Since \\(\sigma(x) \in (0, 1)\\) for every \\(x \in \mathbb{R}\\), the derivative is always strictly positive, confirming that the function is [strictly increasing](<../increasing-and-decreasing-functions/>). The maximum value of the derivative is attained at \\(x = 0\\), where \\(\sigma’(0) = 1/4\\).

## Second derivative and concavity

The second derivative of the sigmoid function is obtained by differentiating the expression: \\[\sigma’(x) = \sigma(x)\,(1 - \sigma(x))\\]

Applying the product rule and substituting the expression for \\(\sigma’(x)\\) gives the following:

\\[\begin{align} \sigma’‘(x) &= \sigma’(x),(1 - \sigma(x)) - \sigma(x)\,\sigma’(x) \\\\[6pt] &= \sigma’(x)\,(1 - 2\sigma(x)) \\\\[6pt] &= \sigma(x)\,(1 - \sigma(x))\,(1 - 2\sigma(x)) \end{align} \\]

The sign of \\(\sigma’'(x)\\) is determined entirely by the factor \\(1 - 2\sigma(x)\\), since \\(\sigma(x)(1 - \sigma(x)) > 0\\) for all \\(x \in \mathbb{R}\\). Since \\(\sigma\\) is strictly increasing and \\(\sigma(0) = \tfrac{1}{2}\\), the factor \\(1 - 2\sigma(x)\\) is positive for \\(x < 0\\) and negative for \\(x > 0\\).

![](/diagrams/algebrica/sigmoid-function-second-derivative-1-1024x558.png)

It follows that the function is concave upward on \\((-\infty, 0)\\) and concave downward on \\((0, +\infty)\\). The point \\(x = 0\\) is therefore an inflection point, at which \\(\sigma’'(0) = 0\\) and the concavity changes sign.

## Relation to the logistic function

The sigmoid function coincides with the special case of the logistic function in which the growth rate equals \\(1\\) and the inflection point is located at the origin. The general form of the logistic function is the following:

\\[f(x) = \frac{L}{1 + e^{-k(x - x_0)}} \\]

In this expression \\(L\\) denotes the upper asymptotic value, \\(k\\) the growth rate, and \\(x_0\\) the inflection point. The standard sigmoid function corresponds to the choice \\(L = 1\\), \\(k = 1\\), and \\(x_0 = 0\\).

## Relation to the hyperbolic tangent

The sigmoid function is closely related to the [hyperbolic tangent](<../hyperbolic-tangent-and-cotangent/>) \\(\tanh.\\) The following identity holds:

\\[\sigma(x) = \frac{1 + \tanh\\!\left(\dfrac{x}{2}\right)}{2} \\]

An equivalent form is the following:

\\[\tanh(x) = 2\,\sigma(2x) - 1 \\]

This relation shows that the two functions differ essentially by a vertical translation and a rescaling. While the sigmoid maps \\(\mathbb{R}\\) into the interval \\((0, 1)\\), the hyperbolic tangent maps \\(\mathbb{R}\\) into the interval \\((-1, 1)\\). Both functions exhibit the same S-shaped curve and the same type of saturation at the extremes.

## Inverse of the sigmoid function

Since the sigmoid function is strictly monotone, it admits an inverse function defined on \\((0, 1)\\). This inverse is known as the logit function, and its expression is the following:

\\[\sigma^{-1}(p) = \ln\\!\left(\frac{p}{1-p}\right) \\]

The argument of the [logarithm](<../logarithms/>) is called the odds ratio. The logit function therefore maps a probability \\(p \in (0,1)\\) to the corresponding real value on the log-odds scale.

## Example

Consider the problem of computing the value of the sigmoid function at \\(x = 2\\) and verifying that its derivative at that point is consistent with the formula \\(\sigma’(x) = \sigma(x)(1 - \sigma(x))\\). The value of the function is the following:

\\[\sigma(2) = \frac{1}{1 + e^{-2}} \\]

Since \\(e^{-2} \approx 0.1353\\), one obtains:

\\[\sigma(2) \approx \frac{1}{1.1353} \approx 0.8808 \\]

Applying the derivative formula, the value of \\(\sigma’(2)\\) is the following:

\\[\sigma’(2) = \sigma(2),\bigl(1 - \sigma(2)\bigr) \approx 0.8808 \cdot 0.1192 \approx 0.1050 \\]

The value of the derivative of the sigmoid function at \\(x = 2\\) is therefore approximately \\(0.1050\\), confirming both the formula and the fact that the function grows very slowly in that region, having already approached saturation.
