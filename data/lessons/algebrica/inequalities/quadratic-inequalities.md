> Content sourced from [Algebrica](https://algebrica.org/quadratic-inequalities/) — CC BY-NC 4.0

## What are quadratic inequalities?

A quadratic inequality or [inequality](<../inequalities/>) of degree two is a second-degree [polynomial](<../polynomials/>) inequality in one variable. The standard form of a quadratic inequality presents the same organization of terms as a [quadratic equation](<../quadratic-equations/>), but instead of the equal sign, it presents the signs \\( > \\), \\( < \\), \\( \geq \\), or \\( \leq \\).

A quadratic inequality is of the form:

\\[ax^2 + bx + c \gt 0 \quad \text{where} \quad a \neq 0\\]

  * The coefficients \\(a\\), \\(b\\), and \\(c\\) are constants and \\(x\\) represents the variable.
  * \\(a\\) is the coefficient of the quadratic term \\(x^2\\), \\(b\\) the coefficient of the linear term \\(x\\) and \\(c\\) the constant term.
  * When \\(a\\) is equal to zero, the equation degenerates into a linear equation \\(bx + c = 0\\), or a constant equation, depending on the values of \\(b\\) and \\(c\\).


Quadratic inequalities not presented in standard form can be converted by transferring all terms to one side, resulting in zero on the right-hand side. For example, the inequality:

\\[x^2 + 3 < 2x + 1\\]

is rewritten as \\(x^2 - 2x + 2 < 0\\) prior to applying the methods outlined in this section.

## The nature of quadratic inequality solutions

The solution to a quadratic inequality is typically expressed as a range of values that satisfy the inequality, rather than a single value. Depending on the sign of the quadratic expression and the nature of its roots, the solution may consist of:

  * an open or closed interval (a continuous set of values),
  * two separate intervals,
  * a single value, in cases where the inequality is satisfied only at one point,
  * or no solution, when the inequality is never true for any real number.


These outcomes depend on the discriminant of the associated quadratic equation and on the direction of the inequality symbol \\(\gt, \lt \\). Understanding this structure is essential for interpreting the graphical behavior of the parabola and identifying which regions of the number line satisfy the condition.

## Resolution method

The first step in solving a quadratic inequality is to find solutions to the corresponding second-degree equation. Then, based on the sign of the inequality, the range of values that satisfy it can be determined. Two widely used methods exist to obtain the solutions of the quadratic equation: the [quadratic formula](<../quadratic-formula/>) and the [factorization method](<../factoring-quadratic-equations/>).

##### The quadratic formula provides solutions to any quadratic equation in a straightforward and systematic way, while the factorization method involves factoring the quadratic expression to obtain its solutions. Both methods are practical and widely used for solving quadratic equations.


Given an inequality in standard form \\(ax^2 +bx +c \gt 0\\) or \\(ax^2 +bx +c \lt 0\\), let’s translate it into the associated equation by setting the polynomial equal to zero. Thus, we obtain a second-degree equation in the form:  
\\[ax^2 +bx +c = 0 \\]

The real [roots](<../roots-of-a-polynomial/>) of the associated equation partition the real line into subintervals on which the quadratic polynomial preserves a constant sign. By determining the sign of the expression on each of these intervals, we can identify the values that satisfy the inequality.

###### It is often useful to understand how the behaviour of a quadratic expression changes when the coefficients depend on a parameter. In these cases, the sign of the quadratic expression is governed by the discriminant as a function of the parameter, and the intervals of positivity or negativity may shift accordingly. A complete discussion of this topic can be found in the dedicated entry on [quadratic equations with parameters](<../quadratic-equations-with-parameters/>).

## Solutions when \\(\Delta > 0\\)

As mentioned earlier, quadratic inequalities, unlike equations, typically yield a solution set consisting of a range of values, determined by the sign of the inequality. Given \\(\Delta > 0\\) (the discriminant of the [quadratic formula](<../quadratic-formula/>)), the associated equation has two distinct real roots \\(x_1\\) and \\(x_2\\). Assuming for simplicity that \\(x_1 < x_2\\), we have the following cases:

When the inequality is of the form \\(ax^2+bx+c \geq 0\\) or \\(ax^2+bx+c > 0\\):

\\[\begin{align} x \leq x_1 \lor x \geq x_2 &\quad \text{for} \quad ax^2+bx+c \geq 0 \\\\[0.5em] x < x_1 \lor x > x_2 &\quad \text{for} \quad ax^2+bx+c > 0 \end{align} \\]

When the inequality is of the form \\(ax^2+bx+c \leq 0\\) or \\(ax^2+bx+c < 0\\):

\\[\begin{align} x_1 \leq x \leq x_2 &\quad \text{for} \quad ax^2+bx+c \leq 0 \\\\[0.5em] x_1 < x < x_2 &\quad \text{for} \quad ax^2+bx+c < 0 \end{align} \\]

From a geometric perspective, recalling that a second-degree polynomial represents a [parabola](<../parabola/>), we obtain the following intervals:

![](https://algebrica.org/wp-content/uploads/resources/images/quadratic-inequalities-1-1.png)

###### Unless otherwise specified, the assumption \\(a > 0\\) is maintained throughout this page. If \\(a < 0\\), the parabola opens downward, and the solution sets in the subsequent sections are reversed: expressions positive for \\(a > 0\\) become negative for \\(a < 0\\), and the converse holds.

## Solutions when \\(\Delta = 0\\)

When \\(\Delta = 0\\), the associated equation has a single real root of multiplicity two, denoted \\(x_1 = x_2\\). We have the following cases:

When the inequality is of the form \\(ax^2+bx+c \geq 0\\) or \\(ax^2+bx+c > 0\\):

\\[\begin{align} \forall , x &\quad \text{for} \quad ax^2+bx+c \geq 0 \\\\[0.5em] \forall , x \neq x_1 &\quad \text{for} \quad ax^2+bx+c > 0 \end{align} \\]

When the inequality is of the form \\(ax^2+bx+c \leq 0\\) or \\(ax^2+bx+c < 0\\):

\\[\begin{align} x = x_1 &\quad \text{for} \quad ax^2+bx+c \leq 0 \\\\[0.5em] \not\exists , x &\quad \text{for} \quad ax^2+bx+c < 0 \end{align} \\]

From a geometric point of view, when \\(\Delta = 0\\), the curve always takes positive values (assuming the coefficient \\(a\\) is greater than zero). The only point where the curve touches the x-axis is at the double root, where the value of the expression is zero \\(x_1 = x_2\\).

![](https://algebrica.org/wp-content/uploads/resources/images/quadratic-inequalities-2.png)

## Solutions when \\(\Delta \lt 0\\)

Given \\( \Delta < 0 \\), the associated quadratic equation has [complex roots](<../quadratic-equations-with-complex-solutions/>). Therefore, the parabola does not intersect the x-axis. Based on the inequality sign and the sign of the leading coefficient ( a ), we distinguish the following cases:

  * When the inequality is of the form \\(ax^2+bx+c \geq 0 \\) or \\(ax^2+bx+c \gt 0 \\) we have for both cases: \\(\forall \, x\\)

  * When the inequality is of the form \\(ax^2+bx+c \leq 0 \\) or \\(ax^2+bx+c \lt 0 \\) we have for both cases: \\(\not\exists \; x\\)


![](https://algebrica.org/wp-content/uploads/resources/images/quadratic-inequalities-3.png)

When the coefficient \\( a \\) is positive, the concavity of the parabola is upward, while when \\( a \\) is negative, the concavity is downward. In the examples, we assumed the expressions were rewritten in the form \\( ax^2 + bx + c > 0 \\) with \\( a \\) positive, resulting in a concave-up parabola.

## Example 1

Solve the quadratic inequality: \\[2x^2 +5x - 3 \gt 0 \\]

As illustrated above, the first fundamental step is to move to the associated quadratic equation by setting the polynomial equal to zero:

\\[2x^2+5x-3 = 0 \\]

Now we use the [quadratic formula](<../quadratic-formula/>) to find the solutions of the equation.

\\[\begin{align} x_{1,2} &= \frac{-5 \pm \sqrt{5^2 -4\cdot(2)\cdot(-3)}}{2\cdot 2} \\\\[1em] &= \frac{-5 \pm \sqrt{25 +24}}{4} \\\\[1em] &= \frac{-5 \pm 7}{4} \\\\[1em] \end{align}\\]


We obtain: \\[x_1 = \frac{2}{4} \to \frac{1}{2} \\] \\[x_2 = - \frac{12}{4} \to -3 \\]


Now determine the range of solutions to the inequality. We can use the graphical method to determine it visually. The inequality is of the form \\(ax^2 +bx +c \gt 0\\). We have:

| \\[-3\\]| \\[\frac{1}{2}\\]|   
---|---|---|---  
| | |   
| | |   
  
The solution of the inequality is: \\[x \lt -3 \;\lor\; x \gt \frac{1}{2} \quad\text{or}\quad x \in \left(-\infty, -3\right) \cup \left(\frac{1}{2}, +\infty\right) \\]

If the inequality had been of the form \\(2x^2 +5x - 3 \lt 0\\) we would have had:

| \\[-3\\]| \\[\frac{1}{2}\\]|   
---|---|---|---  
| | |   
| | |   
| | |   
| | |   
  
In this case, the range of values to be considered is \\(-3 \lt x \lt \frac{1}{2}\\) and the solution to the equation is \\(x \in \left(-3,\frac{1}{2}\right)\\).

## Sign analysis

Another useful method for solving inequalities is [sign analysis](<../sign-analysis-in-inequalities/>). This method is used to determine the intervals where a given expression is positive, negative, or zero. For example, consider the following inequality:

\\[x^2 - 2x - 3 > 0 \\]

The polynomial can be [factored](<../factoring-ac-method/>) as follows:

\\[(x - 3)(x + 1) > 0 \\]

According to the sign product rule, this polynomial factored into the product of two factors is positive when the following condition is met:

\\[\begin{align} &x - 3 > 0 \rightarrow x > 3 \\\\[0.5em] &x + 1 >0 \rightarrow x > -1 \end{align} \\]

Now, we represent the values on a number line and mark the corresponding intervals of positivity and negativity. We have:

| | \\[-1 \\]| \\[3 \\]  
---|---|---|---  
\\[x + 1 > 0 \\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\[x -3 > 0 \\]| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\[(x + 1)(x - 3) > 0 \\]| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
  
###### A product of two factors is positive when both factors share the same sign, and negative when they have opposite signs; a complete treatment of this method is available in the dedicated entry on [sign analysis](<../sign-analysis-in-inequalities/>).

In the last row, we insert the result of the product of the signs between the sign from row 1 and that from row 2, for each interval. The intervals that satisfy our initial inequality are:

\\[x < -1 \quad \text{and} \quad x > 3 \\]

In interval notation, the solution set is:

\\[(-\infty, -1) \cup (3, +\infty)\\]

Plotting the curve on the axes, we obtain:

![](https://algebrica.org/wp-content/uploads/resources/images/quadratic-inequalities-4-1.png)

In this way, we have solved the inequality using the sign table without resorting to solving the associated quadratic equation using the [quadratic formula](<../quadratic-formula>), which would have led to the same result.

## Example 2

Consider the following inequality, where the discriminant will play a decisive role in determining the solution: \\[3x^2 - 2x + 5 > 0\\]

The discriminant of the associated equation \\(3x^2 - 2x + 5 = 0\\) is computed as follows: \\[\Delta = (-2)^2 - 4 \cdot 3 \cdot 5 = 4 - 60 = -56\\]

Because \\(\Delta < 0\\) and the leading coefficient \\(a = 3 > 0\\), the parabola opens upward and does not intersect the x-axis.

![Quadratic inequalities.](https://algebrica.org/wp-content/uploads/resources/images/quadratic-inequalities-8.png)

As a result, the quadratic expression remains strictly positive for all real values of \\(x\\).

The solution is: \\[\forall \, x \in \mathbb{R}\\]

## Selected references

  * **Stony Brook University**. [Quadratic Inequalities](https://www.math.stonybrook.edu/Videos/MAP103Online/Handouts/Lecture-31-Handout.pdf)
  * **University of Illinois at Chicago, J. Baldwin**. [Solving Inequalities and Quadratics](http://homepages.math.uic.edu/~jbaldwin/alginit/solquadsh.pdf)
  * **University of Chicago M. Hidegkuti**. [Quadratic Inequalities](https://teaching.martahidegkuti.com/shared/lnotes/4_collegealgebra/inequalities/inequality2.pdf)


Inequalities

Inequalities express order relations between quantities, defining ranges and constraints.

2.6k

[Linear Inequalities](https://algebrica.org/linear-inequalities/)

2.7k

[Sign Analysis in Inequalities](https://algebrica.org/sign-analysis-in-inequalities/)

1.2k

[Rational Inequalities](https://algebrica.org/rational-inequalities/)

1.1k

[Irrational Inequalities](https://algebrica.org/irrational-inequalities/)

1.2k

[Inequalities with Absolute Value](https://algebrica.org/inequalities-with-absolute-value/)

3.6k

[Logarithmic Inequalities](https://algebrica.org/logarithmic-inequalities/)

869

[Trigonometric Inequalities](https://algebrica.org/trigonometric-inequalities/)

1.3k

[Systems of Inequalities](https://algebrica.org/systems-of-inequalities/)
