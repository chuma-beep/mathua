> Content sourced from [Algebrica](https://algebrica.org/limits/) — CC BY-NC 4.0

## Understanding limits in calculus

The concept of a **limit** is fundamental in mathematics. Intuitively, the limit of a [function](<../functions>) \\( f(x) \\) as \\( x \\) approaches a point \\( x_0 \\) allows us to analyse the behaviour of the function as the values of \\( x \\) get arbitrarily close to \\( x_0 \\).

A neighbourhood of \\( x \\) refers to an interval consisting of all points sufficiently close to \\( x \\). More formally, a neighbourhood of \\( x \\) is any open interval \\( (x - \delta, x + \delta) \\) where \\( \delta > 0 \\). This concept is essential for defining limits and understanding the behaviour of functions as they approach a given point.

![](https://algebrica.org/wp-content/uploads/resources/images/limits-1.png)

The smaller the neighbourhood, the closer the points are to \\( x \\). In other words, as the interval \\( (x - \delta, x + \delta) \\) becomes narrower (with \\( \delta \\) approaching zero), the distance between the points within the neighbourhood and \\( x \\) decreases.

## Definition

Suppose we have a function \\( f(x) \\) whose behaviour we wish to study as \\( x \\) approaches the point \\( x_0 \\). We say that as \\( x \\) tends to \\( x_0 \\), the function \\( f(x) \\) has a limit \\( \ell \\), and we write:

\\[\lim_{x \to x_0} f(x) = \ell \\]

Formally, this statement asserts that for every tolerance \\( \varepsilon > 0 \\), there exists a corresponding distance \\( \delta > 0 \\) such that whenever:

\\[0 < |x - x_0| < \delta \\]

it follows that:

\\[|f(x) - \ell| < \varepsilon \\]

In other words, for every neighbourhood of \\( \ell \\), there exists a sufficiently small neighbourhood of \\( x_0 \\) such that all corresponding function values remain within this neighbourhood. This formal definition precisely articulates the intuitive concept that the limit represents the value approached by \\( f(x) \\) as \\( x \\) becomes arbitrarily close to \\( x_0 \\).


When the definition of a limit applies only to a right neighbourhood or a left neighbourhood of \\( x_0 \\), we refer to these as the right-hand limit and left-hand limit, respectively. They are represented as follows:

\\[\lim_{x \to x_0^+} f(x) \quad \text{and} \quad \lim_{x \to x_0^-} f(x) \\]

## Asymptotes and infinite limits

In general, the value of \\( x \\) in a limit can approach a real number \\( x_0 \\) or \\( \pm \infty \\).

\\[\lim_{x \to x_0} f(x) = \ell \quad \text{or} \quad \lim_{x \to x_0} f(x) = \pm \infty \\]

Additionally, the value of the limit itself can be either a finite number or \\( \pm \infty \\).

\\[\lim_{x \to \pm \infty} f(x) = \ell \quad \text{or} \quad \lim_{x \to \pm \infty} f(x) = \pm \infty \\]


When the limit of \\( f(x) \\) exists and tends to \\(\pm \infty\\) as \\( x \\) approaches a finite real number, the behaviour of the function can resemble the simplified pattern shown in the figure. The line \\( x = k \\) is called a **vertical asymptote** :

![](https://algebrica.org/wp-content/uploads/resources/images/limits-2.png)

In the example, we have the case where the right-hand and left-hand limits of \\( f(x) \\) are, respectively:

\\[\lim_{x \to x_0^+} f(x) = - \infty \quad \text{and} \quad \lim_{x \to x_0^-} f(x) = + \infty \\]


When the limit of \\( f(x) \\) exists and approaches a finite value \\( L \\) as \\( x \\) tends to \\( \pm \infty \\), the behaviour of the function can resemble the simplified pattern shown in the figure. The line \\( y = L \\) is called a **horizontal asymptote** :

![](https://algebrica.org/wp-content/uploads/resources/images/limits-3.png)

In the example, we have the case where the right-hand and left-hand limits of \\( f(x) \\) are, respectively:

\\[\lim_{x \to +\infty} f(x) = L \quad \text{and} \quad \lim_{x \to -\infty} f(x) = L \\]

##### An asymptote is defined as a line that the graph of a function approaches arbitrarily closely as either the \\( x \\)-value or \\( y \\)-value increases or decreases without bound. Consequently, the distance between the curve and the asymptote approaches zero as the graph extends toward the extremes of the coordinate plane.

## Conditions for limit existence and continuity

When the left-hand and right-hand limits of a function both exist and are finite, but have different values with \\( \ell_1 \neq \ell_2 \\), we have:

\\[\begin{cases} \lim\limits_{x \to x_0^-} f(x) = \ell_1 \in \mathbb{R} \\\\[0.5em] \lim\limits_{x \to x_0^+} f(x) = \ell_2 \in \mathbb{R} \end{cases} \implies \nexists \lim\limits_{x \to x_0} f(x) \\]

In this scenario, the limit of \\( f(x) \\) as \\( x \\) approaches \\( x_0 \\) does not exist because the function approaches two distinct values depending on the direction of approach. However, the left-hand limit and the right-hand limit are well-defined and finite when considered separately.


According to the Uniqueness Theorem of Limits, if the limit of a function \\( f(x) \\) as \\( x \\) approaches \\( x_0 \\) exists, whether finite or infinite, such a limit is unique. This statement can be formally represented as:

\\[\lim_{x \to x_0} f(x) = \ell \in \overline{\mathbb{R}} \implies \ell \text{ is unique} \\]

The theorem ensures that if a limit exists, there cannot be two different values satisfying the definition of a limit for the same function and point.


The concept of a limit is fundamental for introducing and defining the concept of a [continuous](<../continuous-functions/>) function. A function \\( y = f(x) \\) is said to be continuous at a point \\( x_0 \\) if the limit of the function as \\( x \\) approaches \\( x_0 \\) exists and is equal to the value of the function at that point. Formally, this is expressed as:

\\[\lim_{x \to x_0} f(x) = f(x_0) \\]

## Properties

The following properties of limits are particularly useful for performing calculations and simplifying complex expressions. They establish the foundational rules for manipulating limits and are essential for working with more advanced mathematical problems.


The limit of the product of a **constant** and a function is equal to the product of the constant and the limit of the function, provided the limit exists.

\\[\lim_{x \to x_0} \left( c \, f(x) \right) = c \, \lim_{x \to x_0} f(x) = c \cdot \ell \\]

Multiplying a function by a constant does not affect the process of taking the limit, other than scaling the result by that constant.


The limit of the algebraic **sum** of two functions is equal to the sum of their individual limits, provided both limits exist.

\\[\lim_{x \to x_0} \left( f(x) + g(x) \right) = \lim_{x \to x_0} f(x) + \lim_{x \to x_0} g(x) = \ell_1 + \ell_2 \\]

The limits of each function can therefore be evaluated separately and then added. This is particularly useful when working with [polynomial functions](<../polynomial-function/>), [trigonometric functions](<../sine-and-cosine>), and other common mathematical expressions.


The limit of the **product** of two functions is equal to the product of their individual limits, provided both limits exist.

\\[\lim\limits_{x \to x_0} \left( f(x) \, g(x) \right) = \lim\limits_{x \to x_0} f(x) \cdot \lim\limits_{x \to x_0} g(x) = \ell_1 \cdot \ell_2 \\]


The limit of the **quotient** of two functions is equal to the quotient of their individual limits, provided both limits exist and the limit of the denominator is not zero.

\\[\lim\limits_{x \to x_0} \left( \frac{f(x)}{g(x)} \right) = \frac{\lim\limits_{x \to x_0} f(x)}{\lim\limits_{x \to x_0} g(x)} = \frac{\ell_1}{\ell_2} \\]

## When standard properties do not apply

The previously outlined properties are valid only when all relevant limits exist and are finite, and the denominator remains nonzero. However, in practical applications, it is common to encounter expressions where direct substitution produces an undefined result, such as:

\\[\dfrac{0}{0} \quad \dfrac{\infty}{\infty} \quad \infty - \infty \\]

Such expressions are classified as [indeterminate forms](<../indeterminate-forms/>). Resolving them requires specialised techniques that extend beyond standard algebraic manipulation of limits. A classic example illustrates this point:

\\[\lim_{x \to 0} \frac{\sin x}{x} \\]

Direct substitution of \\( x = 0 \\) yields \\( \frac{0}{0} \\), which is undefined. The quotient property of limits does not apply in this case because the limit of the denominator is zero. The resolution of such expressions requires a specialised analytical approach, as detailed in the referenced page.

## Fundamental limits of elementary functions

The following limits characterise the asymptotic behaviour of common elementary functions. These limits provide foundational tools for evaluating more complex limits and are frequently encountered in mathematical analysis.


For the constant function \\(f(x) = k\\) with \\(k \in \mathbb{R}\\), we have:

\\[\lim_{x \to -\infty} k = k \quad \text{and} \quad \lim_{x \to +\infty} k = k\\]


For the function \\( f(x) = x \\), we have:

\\[\lim_{x \to -\infty} x = -\infty \quad \text{and} \quad \lim_{x \to +\infty} x = +\infty\\]


For the [exponential function](<../exponential-function/>) with base \\( a > 1 \\), we have:

\\[\lim_{x \to -\infty} a^x = 0 \quad \text{and} \quad \lim_{x \to +\infty} a^x = +\infty\\]


For the exponential function with base \\( 0 < a < 1 \\), we have:

\\[\lim_{x \to -\infty} a^x = +\infty \quad \text{and} \quad \lim_{x \to +\infty} a^x = 0\\]

###### If the base is greater than \\(1\\), the exponential function increases without bound in one direction and approaches zero in the other. If the base is strictly between \\(0\\) and \\(1\\), this behaviour is reversed.


For the power function with an even exponent, we have:

\\[\lim_{x \to -\infty} x^n = +\infty \quad \text{and} \quad \lim_{x \to +\infty} x^n = +\infty\\]


For the power function with an odd exponent, we have:

\\[\lim_{x \to -\infty} x^n = -\infty \quad \text{and} \quad \lim_{x \to +\infty} x^n = +\infty\\]


For root functions with an even index, we have:

\\[\lim_{x \to +\infty} \sqrt[n]{x} = +\infty\\]

###### For even indices, the root function is defined only for \\( x \geq 0 \\). Therefore, the limit as \\( x \to -\infty \\) is not applicable.


For root functions with an odd index, we have:

\\[\lim_{x \to +\infty} \sqrt[n]{x} = +\infty \quad \text{and} \quad \lim_{x \to -\infty} \sqrt[n]{x} = -\infty\\]


For the [logarithmic function](<../logarithmic-function>) with base \\( a > 1 \\), we have:

\\[\lim_{x \to 0^+} \log_a x = -\infty \quad \text{and} \quad \lim_{x \to +\infty} \log_a x = +\infty\\]


For the [logarithmic function](<../logarithmic-function>) with base \\( 0 < a < 1 \\), we have:

\\[\lim_{x \to 0^+} \log_a x = +\infty \quad \text{and} \quad \lim_{x \to +\infty} \log_a x = -\infty\\]


For the [absolute value function](<../absolute-value-function/>) \\( f(x) = |x| \\), we have:

\\[\lim_{x \to -\infty} |x| = +\infty \quad \text{and} \quad \lim_{x \to +\infty} |x| = +\infty\\]


For the sign function \\( \text{sgn}(x) \\), we have:

\\[\lim_{x \to -\infty} \text{sgn}(x) = -1 \quad \text{and} \quad \lim_{x \to +\infty} \text{sgn}(x) = 1\\]


When working with limits, you use algebraic operations to combine, break apart, and work with functions in a systematic way. These rules, such as limits of sums, products, quotients, powers, and compositions, are explained in more detail on the page about the [algebra of limits](<../algebra-of-limits/>).

## Selected references

  * **Harvard University O. Knill**. [Unit 3: Limits](https://people.math.harvard.edu/~knill/teaching/math1a2021/handouts/lecture03.pdf)

  * **University of California, Berkeley A. Vizeff**. [Continuity and Discontinuities](https://math.berkeley.edu/~avizeff/calculus-I-F22/lecture-5.pdf)

  * **University of California, Los Angeles N. Hu**. [Limits of Functions](https://www.math.ucla.edu/~njhu/teaching/math-131a-2024u/notes/10-compact.pdf)

  * **University of California, Berkeley R. Wang**. [Limits](https://math.berkeley.edu/~ruiwang/pdf/104.pdf)

  * **MIT D. Jerison**. [Limits and Continuity](https://ocw.mit.edu/courses/18-01sc-single-variable-calculus-fall-2010/pages/1.-differentiation/part-a-definition-and-basic-rules/session-4-limits-and-continuity/)


Limits

Limits describe a function’s behavior near a point.

941

[Algebra of Limits](https://algebrica.org/algebra-of-limits/)

1.6k

[Squeeze Theorem](https://algebrica.org/squeeze-theorem/)

4.1k

[Remarkable Limits](https://algebrica.org/remarkable-limits/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/limits/remarkable-limits.md?plain=1)

3.8k

[Asymptotes](https://algebrica.org/asymptotes/)

1.9k

[Indeterminate Forms of Limits](https://algebrica.org/indeterminate-forms/)

5.1k

[Little-o Notation](https://algebrica.org/little-o-notation/)

1.8k

[Big O Notation](https://algebrica.org/big-o-notation/)
