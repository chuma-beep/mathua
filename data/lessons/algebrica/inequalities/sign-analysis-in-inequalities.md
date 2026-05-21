> Content sourced from [Algebrica](https://algebrica.org/sign-analysis-in-inequalities/) — CC BY-NC 4.0

## What is sign analysis

Sign analysis of inequalities is a method for determining the intervals in which a given expression is positive, negative, or zero. This approach is especially valuable for solving [polynomial](<../polynomials>) and [rational inequalities](https://algebrica.org/rational-inequalities/) and for analysing the behaviour of functions within specified[domains](<../determining-the-domain-of-a-function/>).


The sign of a product is determined by whether the number of negative factors is even or odd. A product is positive if the number of negative factors is even, and negative if the number of negative factors is odd. We have:

\\[\begin{align} +\times + &\quad = + \\\ +\times - &\quad = - \\\ -\times - &\quad = + \\\ -\times + &\quad = - \end{align} \\]

## Example 1

Let’s consider a simple [quadratic inequality](<../quadratic-inequalities/>):

\\[x^2 + x - 2 > 0 \\]

We need to determine the range of x values that satisfy the inequality. First, we can [factor the polynomial](<../factoring-ac-method/>) in the form:

\\[(x - 1)(x + 2) > 0 \\]

According to the sign product rule, this polynomial factored into the product of two factors is positive when the following condition is met:

\\[\begin{align} &x - 1 > 0 \to x > 1 \\\\[6pt] &x + 2 > 0 \to x > -2 \end{align} \\]

The values \\(x = 1\\) and \\(x = -2\\) divide the real line into disjoint intervals. Since the expression is [continuous](<../continuous-functions/>) and can change sign only at its zeros, its sign remains constant within each interval. These values are then marked on a number line, with \\(+\\) and \\(-\\) indicating the sign on each interval.

| | \\[-2 \\]| \\[1 \\]  
---|---|---|---  
\\( x - 1 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( x +2 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\( (x - 1)(x+2) > 0 \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
  

In the last row, we insert the result of the product of the signs between the sign from row 1 and that from row 2, for each interval. The intervals that satisfy our initial inequality \\( (x - 1)(x + 2) > 0 \\) are:

\\[x < -2 \quad \text{and} \quad x > 1 \\]

In interval notation, the solution set is:

\\((-\infty, -2) \cup (1, +\infty) \\)

###### If the inequality is non-strict, the zeros of the expression belong to the solution. For the inequality: \\(x^2 + x - 2 \geq 0\\) the procedure is identical, but since the inequality is satisfied also when \\(f(x) = 0\\), the zeros \\(x = -2\\) and \\(x = 1\\) are included in the solution set \\((-\infty, -2] \cup [1, +\infty)\\).

## Systematic procedure for sign analysis

In general, sign analysis of inequalities follows a systematic procedure that applies consistently to any polynomial or rational expression.

  * Rewrite the inequality in the form \\(f(x) > 0\\) (or \\(< 0\\), \\(\geq 0\\), \\(\leq 0\\)), with zero on the right-hand side.
  * Factor \\(f(x)\\) into a product or quotient of linear or irreducible factors. For rational inequalities, analyse the numerator and denominator independently.
  * Find the zeros of each factor. Zeros of the denominator never belong to the solution and must be marked on the number line with an open circle.
  * Arrange the zeros in ascending order along the real line, thereby partitioning \\(\mathbb{R}\\) into disjoint intervals. For each interval, determine the sign of each factor.
  * Combine the signs for each interval using the product rule: the overall sign is positive if the number of negative factors is even and negative if it is odd. Factors with even multiplicity do not produce a sign change at their corresponding zero.
  * Identify the intervals where the overall sign fulfils the original inequality. For strict inequalities, exclude zeros from the solution; for non-strict inequalities, include them. Always exclude zeros of the denominator.
  * Present the solution using set notation, combining disjoint intervals with \\(\cup\\).


## Geometric representation

Plotting the curve on the axes, we obtain:

![](/diagrams/algebrica/sign-analysis-1.png)

In this way, we have solved the inequality using the sign table without resorting to solving the associated quadratic equation using the [quadratic formula](<../quadratic-formula>), which would have led to the same result.

In fact, when the inequality is of the form \\( ax^2 + bx + c \ge 0 \\) or \\( ax^2 + bx + c > 0 \\), and the corresponding [quadratic equation](<../quadratic-equations/>) \\( ax^2 + bx + c = 0 \\) has two distinct real solutions \\( x_1 < x_2 \\), we have:

\\[\begin{align} &x \leq x_1 \lor x \geq x_2 &&\text{for} \quad ax^2+bx+c \geq 0 \\\ &x < x_1 \lor x > x_2 &&\text{for} \quad ax^2+bx+c > 0 \end{align} \\]

If a function \\( f(x) \\) is [continuous](<../continuous-functions/>) on an interval \\([a,b]\\) and \\( f(a) \\) and \\( f(b) \\) have opposite signs, then by the
