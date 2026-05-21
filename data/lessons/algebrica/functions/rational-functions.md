> Content sourced from [Algebrica](https://algebrica.org/rational-functions/) — CC BY-NC 4.0

## What are rational functions

Rational functions are [functions](<../functions>) in which both the numerator and the denominator are [polynomials](<../polynomials>), typically of degrees \\(n\\) and \\(m\\). They are usually written in the form

\\[y = \frac{P(x)}{Q(x)} \\]

where \\(P(x)\\) and \\(Q(x)\\) are polynomials and \\(Q(x) \neq 0\\). In their most explicit form, a rational function can be expressed as

\\[y = \frac{a_n x^n + a_{n-1} x^{n-1} + \dots + a_1 x + a_0}{b_m x^m + b_{m-1} x^{m-1} + \dots + b_1 x + b_0} \\]

To make this structure more concrete, consider the simple function:

\\[R(x) = \frac{2x + 3}{x - 1} \\]

In this example both the numerator and the denominator are first–degree polynomials. The expression is well defined for every real number except the value that makes the denominator vanish. Here this happens at \\(x = 1\\), which must therefore be excluded. The [domain](<../determining-the-domain-of-a-function/>) of the function is:

\\[{\, x \in \mathbb{R} : x \neq 1 \,} \\]

Even in such an elementary case, we see the typical features of a rational function: a clearly recognisable algebraic form, a domain determined by the points where the denominator becomes zero, and a behaviour that can be further explored through limits, asymptotes, and algebraic simplification.

## Properties

  * [Domain](<../determining-the-domain-of-a-function/>): it consists of all real numbers except those that make the denominator \\( Q(x) \\) equal to zero. \\(D = \mathbb{R} \setminus { x \in \mathbb{R} \mid Q(x) = 0 .}\\)
  * Range: it represents the set of all [real values](<types-of-numbers>) that the function is capable of attaining, except for those that are inherently unreachable based on the function’s algebraic structure. In essence, it depends on the specific form and behavior of the function.


##### As for the [evenness or oddness](<../even-and-odd-functions/>) of the function, as well as its boundedness, [monotonicity](<../increasing-and-decreasing-functions/>), concavity, and convexity, these properties cannot be determined in advance. They depend entirely on the specific form of the rational function under consideration.

## Continuity of rational functions

Rational functions are among the most well-behaved objects in real analysis when it comes to continuity. Since they are defined as ratios of polynomials and polynomials are continuous for all real numbers, a rational function is automatically continuous at any point \\(x_0\\) where \\(Q(x_0) \neq 0\\). Possible discontinuities occur only at the solutions of:

\\[Q(x) = 0 \\]

where the expression is not defined. Even in these cases, the type of discontinuity is highly structured. If numerator and denominator both vanish:

\\[P(x_0) = 0 \quad Q(x_0) = 0 \\]

the discontinuity may be removable, since a common factor can often be cancelled. If instead

\\[P(x_0) \neq 0 \quad Q(x_0) = 0 \\]

the function diverges and a vertical asymptote appears.

###### Because their behaviour is dictated by algebraic properties, rational functions are continuous everywhere on their domain and exhibit only a small, predictable set of discontinuities.

## Limits at infinity for rational functions

To understand how a rational function behaves for very large values of \\(x\\), it is enough to examine how fast the numerator and the denominator grow. Far from the origin, the lower–degree terms become negligible, and the function is essentially governed by the degrees and leading coefficients of the two polynomials. Consider: \\[R(x) = \frac{P(x)}{Q(x)} \\] where \\(P(x)\\) and \\(Q(x)\\) have degrees \\(n\\) and \\(m\\). The comparison between \\(n\\) and \\(m\\) determines the behaviour of the function as \\(x \to \pm\infty\\):

If \\(n < m\\), the denominator dominates and the limit is: \\[\lim_{x \to \pm\infty} R(x) = 0 \\]

If \\(n = m\\), the highest–degree terms balance. The limit becomes the ratio of the leading coefficients:  
\\[\lim_{x \to \pm\infty} R(x) = \frac{a_n}{b_m} \\]

If \\(n = m + 1\\), the function behaves like a line for large \\(|x|\\), giving rise to an oblique asymptote. Although the limit does not exist as a finite number, the difference: \\[R(x) - (ax + b) \\] tends to zero for a suitable line \\(ax + b\\).

If \\(n > m + 1\\), the numerator grows too fast, and the function diverges, depending on the signs of the leading terms:  
\\[\lim_{x \to \pm\infty} R(x) = \pm\infty \\]

##### These cases offer a structured way to predict the long-range behaviour of any rational function simply by looking at the degrees of the polynomials involved, without the need for detailed algebraic manipulation.

## Limits at points where the denominator becomes zero

When analysing a rational function the most delicate situations arise at the values of \\(x\\) for which the denominator vanishes. Around these points the behaviour of the function can change drastically, and understanding what happens requires distinguishing two fundamentally different cases. In the first case, the denominator becomes zero at a point \\(x_0\\) while the numerator remains nonzero. This corresponds to the situation: \\[Q(x_0) = 0 \qquad P(x_0) \neq 0 \\] Near \\(x_0\\), the expression looks like: \\[\frac{P(x_0)}{0} \\] which indicates that the function grows without bound as \\(x\\) approaches that value. What emerges is a true vertical asymptote, and the limit takes the form: \\[\lim_{x \to x_0^\pm} R(x) = \pm\infty \\] with the sign determined by how the denominator changes on either side of \\(x_0\\).


A different situation occurs when both numerator and denominator become zero at the same point which produces the indeterminate form  
\\[\frac{0}{0} \\] In this scenario, the value of the limit cannot be inferred directly from the polynomials evaluated at \\(x_0\\). The function might simplify to something well-behaved, it might diverge, or it might lead to a finite limit after suitable manipulation. To uncover what truly happens, one usually proceeds by factoring both numerator and denominator to cancel any common factors, or by applying [L’Hôpital’s rule](<../hopital-rule/>) when its conditions are satisfied.

## Asymptotes

Rational functions always exhibit at least one type of [asymptote](<../asymptotes/>), the nature of which depends on the specific structure of the function. In general, the following cases can be identified:

  * If \\( Q(x_0) = 0 \\) and \\( P(x_0) \neq 0 \\), then \\( x = x_0 \\), that is the point where the denominator becomes zero, is a vertical asymptote.
  * If the degree of \\( P(x) \\) is equal to the degree of \\( Q(x) \\), that is, \\( n = m \\), the function has a horizontal asymptote at the ratio of the leading coefficients \\(a_n/b_m\\). If \\( n < m \\), the horizontal asymptote is the x-axis, represented by the line \\( y = 0 \\).
  * Finally, an oblique asymptote may occur when the degree of the polynomial \\( P(x) \\) is greater than the degree of \\( Q(x) \\), specifically when \\( n = m + 1 \\).


## Derivatives and integrals

Rational functions are [continuous](<../continuous-functions/>) and differentiable throughout their entire domain. It is not possible to provide a single general form for the [derivative](<../derivatives>), but the expression used to compute it is given by the following formula:

\\[\frac{d}{dx} \left[ \frac{N(x)}{D(x)} \right] = \frac{N’(x)D(x) - N(x)D’(x)}{[D(x)]^2} \\]


Similarly, there is no single standard method for integrating rational functions; the integration process depends entirely on the specific form of the function. However, in general, the [integral of a rational function](<../integral-of-rational-functions/>) can be computed by applying the following formula:

\\[\int \frac{N(x)}{D(x)} \, dx = \int Q(x) \, dx + \int \frac{R(x)}{D(x)} \, dx \\]

  * \\( Q(x) \\) is the quotient obtained from the division of \\( N(x) \\) by \\( D(x) \\).
  * \\( R(x) \\) is the remainder of the division.
