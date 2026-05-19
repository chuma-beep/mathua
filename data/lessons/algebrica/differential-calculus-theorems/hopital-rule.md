> Content sourced from [Algebrica](https://algebrica.org/hopital-rule/) — CC BY-NC 4.0

## Indeterminate forms

L’Hôpital’s rule is a method for evaluating certain [limits](<../limits>) that result in [indeterminate forms](<../indeterminate-forms/>). The theorem establishes a criterion for resolving the indeterminate form of the limit of one or more functions by utilizing their [derivatives](<../derivatives>). By indeterminate forms, we refer to expressions of uncertainty such as: \\[\frac{0}{0} \quad \text{and} \quad \frac{\infty}{\infty} \\]

These forms prevent the direct evaluation of a limit, as they signify cases where standard limit theorems alone are insufficient to determine the result and require additional analytical techniques.

## Statement

According to l’Hôpital’s rule, if two [functions](<../functions/>), \\( f(x) \\) and \\( g(x) \\), are defined on a punctured neighbourhood \\( I \\) of a point \\( x_0 \\), the following conditions must be satisfied:

  * \\( f(x) \\) and \\( g(x) \\) are differentiable in \\( I \\), except possibly at \\( x_0 \\).
  * \\( g’(x) \neq 0 \\) for all \\( x \in I \\), with \\( x \neq x_0 \\).
  * \\( \displaystyle \lim_{x \to x_0} f(x) = \lim_{x \to x_0} g(x) = 0 \\).
  * The following limit exists (finite or infinite): \\[\lim\limits_{x \to x_0} \frac{f’(x)}{g’(x)} \\]


If the above conditions hold, then the following limit exists: \\[\lim\limits_{x \to x_0} \frac{f(x)}{g(x)} = \lim\limits_{x \to x_0} \frac{f’(x)}{g’(x)} \\]

##### In simpler terms, if the quotient of two functions results in an indeterminate form, its limit can be determined by evaluating the limit of their derivatives as \\( x \to x_0 \\), provided that this latter limit exists.


The rule similarly applies to the indeterminate form \\( \infty/\infty \\). Assume that \\( f(x) \to \infty \\) and \\( g(x) \to \infty \\) as \\( x \to x_0 \\), and that both functions are differentiable in a neighbourhood of \\( x_0 \\) with \\( g’(x) \neq 0 \\). If the limit of the quotient of the derivatives exists, then the original limit is equal to this value, as follows:

\\[\lim_{x \to x_0} \frac{f(x)}{g(x)} = \lim_{x \to x_0} \frac{f’(x)}{g’(x)} \\]

## Proof

To prove the theorem, let us consider an arbitrary point \\( x \in I \\) with \\( x > x_0 \\). We apply [Cauchy’s theorem](<../cauchy-theorem/>) to \\( f(x) \\) and \\( g(x) \\). The theorem states that under certain conditions, there exists a point \\( c \in ]x_0, x[ \\) such that:  
\\[\frac{f’\left( c \right )}{g’\left( c \right)} = \frac{f(x) - f(x_0)}{g(x) - g(x_0)} \\]


By hypothesis, we have \\( f(x_0) = g(x_0) = 0 \\). Therefore, \\((1)\\) becomes:  
\\[\frac{f’\left( c \right)}{g’\left( c \right)} = \frac{f(x)}{g(x)} \\]


At this stage, the point \\( c \\) is not fixed: it depends on \\( x \\). More precisely, for each \\( x \neq x_0 \\), Cauchy’s Theorem guarantees the existence of a point \\( c = c(x) \\) such that: \\[x_0 < c(x) < x \quad \text{if } x > x_0\\] \\[x < c(x) < x_0 \quad \text{if } x < x_0\\] Assume for instance that \\( x \to x_0^+ \\). Then, we have \\(0 < c(x) - x_0 < x - x_0.\\) Since \\( x - x_0 \to 0 \\), by the [Squeeze Theorem](<../squeeze-theorem/>) it follows that \\(c(x) - x_0 \to 0\\) and therefore \\(c(x) \to x_0\\). An analogous argument holds when \\( x \to x_0^- \\). Hence, we conclude that: \\[\lim_{x \to x_0} \frac{f’(c(x))}{g’(c(x))} = \lim_{x \to x_0} \frac{f(x)}{g(x)} \\]


If \\( f’(x) \\) and \\( g’(x) \\) are continuous at \\( x_0 \\), then their limits at the points \\( c \\) and \\( x \\) coincide. Therefore, the following equality holds: \\[\lim\limits_{x \to x_0} \frac{f’ \left( c \right)}{g’ \left(c \right)} = \lim\limits_{x \to x_0} \frac{f’(x)}{g’(x)} \\]


Then we obtain what we wanted to prove: \\[\lim\limits_{x \to x_0} \frac{f(x)}{g(x)} = \lim\limits_{x \to x_0} \frac{f’(x)}{g’(x)} \\]

## Example 1

Let’s compute the following limit involving the [sine function](<../sine-function>).

\\[\lim_{x \to 0} \frac{\sin x}{x}\\]


This is a fairly simple limit, but at first glance, it leads to an indeterminate form. Indeed, substituting \\(0\\) for \\(x\\), we get:

\\[\frac{\sin 0}{0} = \frac{0}{0} \\]

Because the functions meet the necessary conditions for l’Hôpital’s Rule, the limit of the quotient may be replaced with the limit of the ratio of their derivatives:

\\[\lim_{x \to 0} \frac{\sin x}{x} = \lim_{x \to 0} \frac{(\sin x)'}{(x)'} = \lim_{x \to 0} \frac{\cos x}{1} \\]

##### The conditions of the theorem are satisfied. Indeed, we have that \\(\sin(x)\\) and \\(x\\) are [continuous functions](<../continuous-functions/>) at \\(x_0 = 0\\), and \\(\sin(0) = 0,\, x\big|_{x=0} = 0\\). Moreover, both functions are differentiable in an interval \\(I\\) containing \\(0\\), and the derivative of the denominator \\(g’(x)\\), which in this case is simply \\(1\\), is different from zero.


In this case, by evaluating the limit and computing \\(\cos(x)\\), we find that the limit is equal to (1):

\\[\lim_{x \to 0} \frac{\cos x}{1} = \frac{\cos(0)}{1} = \frac{1}{1} = 1 \\]

We can therefore conclude that:

\\[\lim_{x \to 0} \frac{\sin x}{x} = \lim_{x \to 0} \frac{\cos x}{1} = 1\\]

## Example 2

Let us now consider a more complex scenario in which the expression results in an indeterminate form of the type \\( -\infty + \infty \\). In these situations, the recommended approach is to rewrite the difference between the two functions as either a product or a quotient. This transformation allows the expression to be converted into one of the standard indeterminate forms to which L’Hôpital’s Rule applies, namely:

\\[\frac{0}{0} \quad \text{or} \quad \frac{\infty}{\infty} \\]


Let’s consider the following limit:

\\[\lim_{x \to 0} \left(\frac{1}{\sin x} - \frac{2}{x}\right)\\]

This limit leads to an indeterminate form of type \\(-\infty + \infty\\). To apply L’Hôpital’s Rule\\(^*\\), we must first rewrite it as a single fraction, thus obtaining an indeterminate form of type \\(\frac{0}{0}\\) or \\(\frac{\infty}{\infty}\\). We have:

\\[\lim_{x \to 0}\left(\frac{1}{\sin x} - \frac{2}{x}\right) = \lim_{x \to 0}\frac{x - 2\sin x}{x\sin x} \\]

##### Before applying the theorem, it is always necessary to verify that the initial conditions are satisfied.


Computing the derivative, the limit becomes: \\[\lim_{x \to 0} \frac{1 - 2\cos x}{\sin x + x\cos x}\\]

Substituting \\( x = 0 \\) in the resulting expression yields \\(-1/0\\), which means the limit diverges.

The result is:

\\[\lim_{x \to 0} \frac{1 – 2\cos x}{\sin x + x\cos x} = -\infty\\]

Generally, if the application of L’Hôpital’s Rule yields a limit that remains an [indeterminate form](<../indeterminate-forms/>), the rule may be applied repeatedly, provided the necessary conditions are met at each stage. At every iteration, it is essential to confirm that both the new numerator and denominator approach either \\(0\\) or \\(\infty\\).

## Indeterminate products

The same principle illustrated in Example 2 can be applied when encountering indeterminate forms of the type \\(0 \cdot \infty\\), arising from the product of two functions \\(f(x)\cdot g(x)\\). In this case, to rewrite the expression in a form suitable for applying L’Hôpital’s Rule, it is sufficient to express the product as follows:

\\[f(x)\cdot g(x) = \frac{f(x)}{\dfrac{1}{g(x)}} \quad\text{or}\quad f(x)\cdot g(x) = \frac{g(x)}{\dfrac{1}{f(x)}} \\]


For example, consider the following limit:

\\[\lim_{x \to 0^+} x \ln x \\]

This expression represents an indeterminate form of type \\(0 \cdot \infty\\). The product can be rewritten as a quotient:

\\[\lim_{x \to 0^+} x \ln x = \lim_{x \to 0^+} \frac{\ln x}{\dfrac{1}{x}} \\]

The resulting expression is now an indeterminate form of type \\(\dfrac{\infty}{\infty}\\), which is suitable for the application of L’Hôpital’s Rule:

\\[\lim_{x \to 0^+} \frac{\ln x}{\dfrac{1}{x}} = \lim_{x \to 0^+} \frac{(\ln x)'}{\left(\dfrac{1}{x}\right)'} = \lim_{x \to 0^+} \frac{\dfrac{1}{x}}{-\dfrac{1}{x^2}} = \lim_{x \to 0^+} (-x) = 0 \\]

## Selected references

  * **MIT, E. Herman**. [Calculus Vol. 1 — 4.8 L’Hôpital’s Rule](https://openstax.org/books/calculus-volume-1/pages/4-8-lhopitals-rule-introduction-to-improper-integrals)

  * **University of British Columbia (Feldman, Rechnitzer, Yeager)**. [L’Hôpital’s Rule](https://personal.math.ubc.ca/~CLP/CLP1/clp_1_dc/sec_4_7.html)

  * **MAA Convergence, (D. E. Otero, Xavier University)**. [L’Hôpital’s Rule](https://old.maa.org/press/periodicals/convergence/l-h-pital-s-rule-a-mini-primary-source-project-for-calculus-1-students)


Differential Calculus Theorems

Core theorems describe the behavior of differentiable functions, including limits, tangents, and mean value properties.

2k

[Weierstrass Theorem](https://algebrica.org/weierstrass-theorem/)

4.7k

[Fermat’s Theorem](https://algebrica.org/fermat-theorem/)

2.7k

[Rolle’s Theorem](https://algebrica.org/rolles-theorem/)

2.1k

[Lagrange’s Theorem](https://algebrica.org/lagrange-theorem/)

3k

[Cauchy’s Theorem](https://algebrica.org/cauchy-theorem/)
