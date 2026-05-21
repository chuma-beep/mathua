> Content sourced from [Algebrica](https://algebrica.org/asymptotes/) — CC BY-NC 4.0

## Horizontal asymptotes

Asymptotes are a fundamental concept in mathematical analysis. They are [lines](<../lines>) that a [function](<../functions>) approaches indefinitely without ever reaching, which helps to characterise the function’s behaviour, particularly near infinity or at points of discontinuity. Asymptotes play a fundamental role in the [analysis of functions](<../analyzing-the-graphs-of-functions/>), as their definition inherently relies on the concept of [limits](<../limits/>).

###### In general terms, an asymptote describes the limiting behaviour of a function through a straight line.


Consider a real-valued function \\( y = f(x) \\) defined on an interval \\( [a, +\infty[ \\), \\( ]-\infty, b] \\), or on \\( \mathbb{R} \\). The line with equation \\( y = L \\) is defined as a **horizontal asymptote** of ( f ) if:

\\[\lim_{x \to +\infty} f(x) = L \quad \text{or} \quad \lim_{x \to -\infty} f(x) = L \\]

In other words, the function approaches the horizontal line \\( y = L \\) as \\( x \\) tends to positive or negative infinity. Let us consider, for example, the function:

\\[y = \frac{x + 1}{x} \\]

![](/diagrams/algebrica/asymptotes-1.png)

By computing the limit as \\( x \to \pm\infty \\), we obtain:

\\[\lim_{x \to \pm\infty} \frac{x + 1}{x} = \lim_{x \to \pm\infty} \left(1 + \frac{1}{x}\right) = 1 \\]

Therefore, the function has a horizontal asymptote along the line \\( y = 1 \\). As we can see from the graph of the function, both branches get closer and closer to the line \\( y = 1 \\) as \\( x \\) tends to positive or negative infinity. This behavior confirms that \\( y = 1 \\) is a horizontal asymptote.

## Vertical asymptotes

Let \\( y = f(x) \\) be a real-valued function defined on an interval \\( [a, b] \setminus {x_0} \\), where \\( x_0 \in [a, b] \\). We say that the line with equation \\( x = x_0 \\) is a **vertical asymptote** of \\( f \\) if:

\\[\lim_{x \to x_0^-} f(x) = \pm\infty \quad \text{or} \quad \lim_{x \to x_0^+} f(x) = \pm\infty \\]

In other words, the function diverges as \\( x \\) approaches \\( x_0 \\) from the left or the right, getting arbitrarily large in [absolute value](<../absolute-value>). Let us consider, for example, the function:

\\[y = \frac{1}{x - 1}\\]

![](/diagrams/algebrica/asymptotes-2.png)

We observe that the [rational function](<../rational-functions>) is undefined at \\( x = 1 \\), since the denominator becomes zero. To analyze the behavior of \\( f(x) \\) near \\( x = 1 \\), we compute the one-sided limits:

\\[\lim_{x \to 1^-} \frac{1}{x - 1} = -\infty \quad \lim_{x \to 1^+} \frac{1}{x - 1} = +\infty \\]

This means that the function diverges to \\( -\infty \\) when approaching \\(1\\) from the left, and to \\( +\infty \\) when approaching from the right. Therefore, the line \\( x = 1 \\) is a vertical asymptote of the function.

###### [Rational functions](<../rational-functions>) often have vertical asymptotes at points where the denominator is zero, and the function is undefined. These points correspond to non-removable discontinuities, which are typical of this type of function.

## Oblique asymptotes

Let \\( y = f(x) \\) be a real-valued function defined on the half-line \\( ] -\infty, a] \\) or \\( [a, +\infty[ \\). We say that the line with equation \\( y = px + q \\) is an **oblique asymptote** of \\( f \\) if the following condition holds:

\\[\begin{align} \lim_{x \to -\infty} \left[f(x) - (px + q)\right] &= 0 \\\\[0.5em] \lim_{x \to +\infty} \left[f(x) - (px + q)\right] &= 0 \end{align} \\]

In other words, the difference between the function and the line \\( y = px + q \\) tends to zero as \\( x \\) tends to infinity or negative infinity. This means that the function behaves more and more like the line \\( y = px + q \\) for large values of \\( x \\).


The equation of the oblique asymptote \\( y = px + q \\) for a function \\( f(x) \\) can be determined by computing two specific limits. The slope \\( p \\) of the asymptote is found by evaluating the following limit:

\\[p = \lim_{x \to \pm\infty} \frac{f(x)}{x} \\]

Once the slope is known, we find the vertical offset \\( q \\) by computing:

\\[q = \lim_{x \to \pm\infty} \left[f(x) - px\right] \\]

If both limits exist and are finite, the line \\( y = px + q \\) is the oblique asymptote of the function.


Let us consider the function:

\\[f(x) = \frac{x^2 + 1}{x} \\]

![](/diagrams/algebrica/asymptotes-3-1.png)

To determine whether this function has an oblique asymptote as \\( x \to \pm\infty \\), we begin by analyzing its behavior for large values of \\( x \\). We start by computing the limit:

\\[\frac{f(x)}{x} = \frac{x^2 + 1}{x^2} = 1 + \frac{1}{x^2} \\]

As \\( x \to \pm\infty \\), the term \\( \dfrac{1}{x^2} \\) tends to zero, so:

\\[\lim_{x \to \pm\infty} \frac{f(x)}{x} = 1 \\]

This tells us that the slope of the asymptote is \\( p = 1 \\). Next, we compute the limit of the difference between the function and the linear term \\( px \\), to find the vertical offset:

\\[f(x) - x = \frac{x^2 + 1}{x} - x = \frac{x^2 + 1 - x^2}{x} = \frac{1}{x} \\]

And again, since \\( \frac{1}{x} \to 0 \\) as \\( x \to \pm\infty \\), we find:

\\[\lim_{x \to \pm\infty} [f(x) - x] = 0 \\]

Therefore, the function has an oblique asymptote with equation:

\\[y = x \\]

###### In this example, the oblique asymptote passes through the origin, resulting in \\( q = 0 \\), which represents a degenerate case. Generally, the vertical offset \\( q \\) does not need to be zero.

## Example 1

Consider the following function as an additional example, which extends the discussion beyond the degenerate case presented previously:

\\[f(x) = \frac{2x^2 - x + 3}{2x} \\]

To determine the oblique asymptote, first compute the slope \\( p \\):

\\[\begin{align} p &= \lim_{x \to \pm\infty} \frac{f(x)}{x} \\\\[6pt] &= \lim_{x \to \pm\infty} \frac{2x^2 - x + 3}{2x^2} \\\\[6pt] &= 1 \end{align} \\]

Next, evaluate the vertical offset \\( q \\):

\\[\begin{align} q &= \lim_{x \to \pm\infty} \bigl(f(x) - x\bigr) \\\\[6pt] &= \lim_{x \to \pm\infty} \frac{2x^2 - x + 3 - 2x^2}{2x} \\\\[6pt] &= \lim_{x \to \pm\infty} \frac{-x + 3}{2x} \\\\[6pt] &= -\frac{1}{2} \end{align} \\]

Therefore, the equation of the oblique asymptote is given by:

\\[y = x - \frac{1}{2} \\]

## Summary

|   
---|---  
Horizontal| \\[\lim_{x \to \pm\infty} f(x) = L \\]  
Vertical| \\[\lim_{x \to x_0^\pm} f(x) = \pm\infty \\]  
Oblique| \\[\lim_{x \to \pm\infty} [f(x) - (px+q)] = 0 \\]  
  
## Key properties of asymptotes

Asymptotes come in different forms and follow specific rules that are worth keeping in mind. Some of these properties are immediately intuitive, while others become clear only after working through a few examples. The following points summarize the most important facts about asymptotes and how they relate to a function’s behavior.

  * Not all functions possess asymptotes.
  * With respect to horizontal asymptotes, several configurations are possible: a function may have none, it may approach the same horizontal line as \\( x \to +\infty \\) and \\( x \to -\infty \\), or it may approach two different horizontal lines in the two directions.
  * Different types of asymptotes can also occur together. A function may simultaneously exhibit horizontal, vertical, and oblique asymptotes, depending on its behaviour near discontinuities and as the variable tends to infinity.
  * Vertical asymptotes typically occur at points where the function is undefined as a result of division by zero. These asymptotes correspond to non-removable discontinuities.
  * Horizontal asymptotes characterise the end behaviour of a function as it approaches a constant value when \\( x \\) becomes very large or very small.
  * Oblique asymptotes occur when the degree of the numerator exceeds that of the denominator by exactly one, causing the function to approach a slanted line as \\( x \\) approaches infinity.


Limits

Limits describe a function’s behavior near a point.

5.9k

[Limits](https://algebrica.org/limits/)

942

[Algebra of Limits](https://algebrica.org/algebra-of-limits/)

1.6k

[Squeeze Theorem](https://algebrica.org/squeeze-theorem/)

4.1k

[Remarkable Limits](https://algebrica.org/remarkable-limits/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/limits/remarkable-limits.md?plain=1)

1.9k

[Indeterminate Forms of Limits](https://algebrica.org/indeterminate-forms/)

5.1k

[Little-o Notation](https://algebrica.org/little-o-notation/)

1.8k

[Big O Notation](https://algebrica.org/big-o-notation/)
