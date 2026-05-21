> Content sourced from [Algebrica](https://algebrica.org/quadratic-equations-with-parameters/) — CC BY-NC 4.0

Concept

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Easy

0

Requires

0

Enables

## Definition

[Quadratic equations](<../quadratic-equations/>) are usually introduced in their classical form, where all coefficients are fixed real numbers. The standard expression serves as the basic model for every second-degree equation: \\[a x^{2} + b x + c = 0 \\]

The key ideas behind quadratic equations are closely linked to the geometry of the [parabola](<../parabola/>) and to the behaviour of the [discriminant](<../quadratic-formula/>). These elements determine how the graph bends, where its vertex lies, and whether the equation has two real solutions, a repeated solution, or a pair of [complex roots](<../quadratic-equations-with-complex-solutions/>).

![](/diagrams/algebrica/quadratic-equations-params-6.png)

In many situations, however, the coefficients are not constant but depend on an external quantity that we treat as a parameter. Instead of being fixed, the numbers \\( a \\), \\( b \\), and \\( c \\) become functions of a real parameter \\( k \\). A parametrised quadratic equation can therefore be written in the more general form:

\\[a(k)\,x^{2} + b(k)\,x + c(k) = 0 \\] \\[a(k) \ne 0 \\]

The condition \\( a(k) \ne 0\\) prevents the equation from collapsing into a [linear equation with parameters](<../linear-equations-with-parameters/>) when the parameter varies, ensuring that it retains the structure of a quadratic for all admissible values of \\( k \\).

###### This page focuses on the quadratic case. For a broader discussion of parametric equations across different degrees, see [equations with parameters](<../equations-with-parameters/>).

## The role of the discriminant

The discriminant allows to understand how the solutions of a quadratic equation behave. When the coefficients depend on a parameter, this idea remains the same, but the discriminant itself becomes a function of that parameter. Starting from the coefficients \\( a(k) \\), \\( b(k) \\), and \\( c(k) \\), we obtain:

\\[\Delta(k) = b(k)^{2} - 4\,a(k)\,c(k) \\]

Once written in this form, the discriminant gives a picture of what happens to the roots for each value of \\( k \\). Its sign tells us everything we need:

  * When \\( \Delta(k) > 0 \\), the equation has two distinct real solutions.
  * When \\( \Delta(k) = 0 \\), the equation has one real solution that appears as a double root.
  * When \\( \Delta(k) < 0 \\), the equation has two [complex conjugate solutions](<../quadratic-equations-with-complex-solutions/>).


## Example 1

Let us begin with a very simple case and look at the quadratic equation: \\[x^{2} + kx + 1 = 0 \\]

This example is useful because each coefficient depends on the parameter \\( k \\) in a straightforward way:

\\[a(k) = 1 \qquad b(k) = k \qquad c(k) = 1 \\]

Before analysing the roots, it is useful to recall how the discriminant is defined in the general case. For a quadratic equation in the standard form the discriminant is:

\\[\Delta = b^{2} - 4ac \\]

and its sign determines the nature of the solutions. We now examine how the discriminant changes when the coefficients depend on the parameter \\(k\\). By substituting \\( a(k) \\), \\( b(k) \\), and \\( c(k) \\) into the same formula, we obtain:

\\[\Delta(k) = b(k)^{2} - 4\,a(k)\,c(k) \\]

The expression simplifies to:

\\[\Delta(k) = k^{2} - 4 \\]

To understand when real solutions exist, we study the [quadratic inequality](<../quadratic-inequalities/>):

\\[k^{2} - 4 \ge 0 \\]

Solving it yields two separate intervals:

\\[k \le -2 \quad\text{or}\quad k \ge 2 \\]

For these values of the parameter, the discriminant is positive and the quadratic has two distinct real roots. At the boundary points \\( k = \pm 2 \\), the discriminant becomes zero and the equation collapses to a single repeated solution. When:

\\[-2 < k < 2 \\]

the discriminant is negative, and the equation no longer has real roots. Instead, it produces a pair of complex conjugate solutions. The following table summarises the [sign analysis](<../sign-analysis-in-inequalities/>) of the expressions \\(k - 2\\), \\(k + 2\\), and their product, highlighting how the sign changes across the three intervals determined by the critical points \\(k = -2\\) and (k = 2).

| | \\[-2 \\]| \\[2 \\]  
---|---|---|---  
\\( k - 2 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( k + 2 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\( (k - 2)(k + 2) > 0 \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
  

Whenever the roots turn out to be real, it is useful to have an explicit expression for them.  
By applying the quadratic formula to this particular family, the two real solutions can be written in closed form as

\\[x = \frac{-k \pm \sqrt{k^{2} - 4}}{2} \\]

This expression shows directly how each root depends on the parameter \\(k\\) and makes it possible to study how the solutions move along the real line as \\(k\\) varies.

From the analysis above, the final outcomes are:

  * If \\( k < -2 \\) or \\( k > 2 \\), the equation has two real and distinct solutions.
  * If \\( k = \pm 2 \\), there is one real double solution.
  * If \\( -2 < k < 2 \\), the equation has two complex conjugate solutions.


## Geometric interpretation of parameter changes

Making reference to the quadratic equation used in example 1, it is also helpful to look at how the graph of the corresponding parabola changes as the parameter \\(k\\) varies. To each value of \\(k\\) we can associate the [parabola](<../parabola/>):

\\[y = x^{2} + kx + 1 \\]

so analysing the roots is equivalent to studying where this curve intersects the x-axis. As the parameter moves, the entire family of parabolas shifts in a smooth and predictable way. One way to describe this evolution is to examine what happens in three key regimes of the parameter.


When \\(|k|\\) becomes very large (\\(k \to +\infty\\) or \\(k \to -\infty\\)), the linear term \\(kx\\) dominates the expression and the two roots separate widely. A simple [asymptotic](<../asymptotes/>) argument shows that one root tends to \\(0\\), while the other behaves approximately like \\(-k\\). Geometrically, this means that the intersection points move far apart: one stays close to the origin, while the other drifts further and further along the negative axis. The parabola becomes increasingly tilted in appearance.


At the critical values \\(k = \pm 2\\), the discriminant becomes zero. This corresponds to the moment when the parabola touches the x-axis at exactly one point. In these cases the graph has a single point of tangency, the equation has one repeated root, and the vertex lies precisely on the axis.


Within the interval \\(-2 < k < 2\\), the discriminant is negative, so the equation has no real solutions. From a geometric perspective, this means that the parabola does not intersect the x-axis at all. Because the leading coefficient is positive, the parabola always opens upward, and in this parameter region the vertex lies strictly above the x-axis. As a result, the entire curve stays on one side of the axis, reflecting the fact that the solutions move into the complex plane.

## Example 2

Let us now consider a slightly more elaborate example than the previous one. We look at the parameter-dependent quadratic equation:

\\[(a + 2)x^{2} + (3a - 1)x + (a - 4) = 0 a \neq -2 \\]

where the condition \\(a \neq -2\\) ensures that the coefficient of \\(x^{2}\\) does not vanish, so the expression remains a genuine quadratic equation for every admissible value of the parameter. Our goal is to determine for which values of \\(a\\) the equation admits two distinct real solutions. As usual, this depends on the sign of the [discriminant](<../quadratic-formula/>). To avoid confusion with the parameter \\(a\\), we denote the coefficients of the present quadratic by \\(A\\), \\(B\\), \\(C\\), so that the general formula reads:

\\[\Delta = B^{2} - 4AC \\]

We substitute the coefficients of the present quadratic equation and obtain:

\\[\Delta(a) = (3a - 1)^{2} - 4(a + 2)(a - 4) \\]


Expanding the two parts with only the essential intermediate steps, we first compute the square:

\\[(3a - 1)^{2} = 9a^{2} - 6a + 1 \\]

Then the product involving the quadratic and constant coefficients:

\\[(a + 2)(a - 4) = a^{2} - 2a - 8 \\]

Inserting these expressions back into the discriminant and collecting like terms leads to the simplified form

\\[\Delta(a) = 5a^{2} + 2a + 33 \\]


This is a quadratic function of \\(a\\) with positive leading coefficient. To understand whether it can ever be zero or negative, we examine its discriminant:

\\[2^{2} - 4 \cdot 5 \cdot 33 = 4 - 660 < 0 \\]

Since this value is negative, the parabola associated with (5a^{2} + 2a + 33) never intersects the horizontal axis: it remains strictly positive for all real values of \\(a\\). Therefore:

\\[\Delta(a) > 0 \quad \forall a \in \mathbb{R} \\]

The equation has two real and distinct solutions for every real value of the parameter \\(a\\), except for \\(a = -2\\), which must be excluded because it would make the coefficient of \\(x^{2}\\) vanish and reduce the expression to a first-degree equation.

## Example 3

Consider the parameter-dependent quadratic equation:

\\[x^{2} - (a - 4)x + (a^{2} - a) = 0 \\]

where \\(a\\) is a real parameter. We want to determine the values of \\(a\\) for which the product of the two solutions is strictly less than \\(2\\). For a general quadratic equation of the form:

\\[x^{2} + bx + c = 0 \\]

the product of the solutions satisfies the identity \\(x_{1}x_{2} = c\\).

###### This identity follows from the [factorised form](<../factoring-quadratic-equations/>) of a monic quadratic, where \\(x^{2} + bx + c = (x - x_{1})(x - x_{2})\\). By expanding the product, the constant term appears as \\(x_{1}x_{2}\\), which explains why the product of the solutions is equal to \\(c\\).


In this case the constant term is \\(c = a^{2} - a\\), and requiring the product of the two solutions to be less than \\(2\\) translates into the inequality:

\\[a^{2} - a < 2 \\]

Bringing all the terms to the left gives the inequality:

\\[a^{2} - a - 2 < 0 \\]

The quadratic [polynomial](<../polynomials/>) on the left factors as

\\[a^{2} - a - 2 = (a - 2)(a + 1) \\]

The critical points are therefore (a = -1) and (a = 2). To determine where the expression is negative, we analyse the signs of the two factors.

| | \\[-1\\]| \\[2\\]  
---|---|---|---  
\\( a - 2 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
\\( a + 1 > 0 \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{+} \\)  
\\( (a - 2)(a + 1) < 0 \\)| \\( \boldsymbol{+} \\)| \\( \boldsymbol{-} \\)| \\( \boldsymbol{+} \\)  
  

The sign chart shows that \\((a - 2)(a + 1) < 0\\) when the parameter lies between the two roots:

\\[-1 < a < 2 \\]

The product of the solutions satisfies \\(x_{1}x_{2} < 2\\) if and only if \\(-1 < a < 2.\\)

Equations

Equations state equalities between quantities and are fundamental for solving mathematical problems.

9.8k

[Equations](https://algebrica.org/equations/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/equations.md?plain=1)

8.2k

[Polynomial Equations](https://algebrica.org/polynomial-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/equations.md?plain=1)

9.2k

[Linear Equations](https://algebrica.org/linear-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations.md?plain=1)

17.7k

[Quadratic Equations](https://algebrica.org/quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations.md?plain=1)

13.9k

[Quadratic Formula](https://algebrica.org/quadratic-formula/)

2 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/quadratic-formula.md?plain=1)

2.1k

[Factoring Quadratic Equations](https://algebrica.org/factoring-quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/factoring-quadratic-equations.md?plain=1)

3.9k

[Incomplete Quadratic Equations](https://algebrica.org/incomplete-quadratic-equations/)

1 comment[](https://github.com/antoniolupetti/algebrica/blob/main/equations/incomplete-quadratic-equations.md?plain=1)

1.7k

[The Geometric Interpretation of Quadratic Equations](https://algebrica.org/geometrical-meaning-quadratic-equations/)

2.2k

[Loss of Roots](https://algebrica.org/loss-of-roots/)

1.4k

[Quadratic Equations with Complex Solutions](https://algebrica.org/quadratic-equations-with-complex-solutions/)

3.1k

[Binomial Equations](https://algebrica.org/binomial-equations/)

3.8k

[Trinomial Equations](https://algebrica.org/trinomial-equations/)

3.8k

[Rational Equations](https://algebrica.org/rational-equations/)

12.7k

[Irrational Equations](https://algebrica.org/irrational-equations/)

1.7k

[Absolute Value Equations](https://algebrica.org/absolute-value-equations/)

1.6k

[Exponential Equations](https://algebrica.org/exponential-equations/)

2.2k

[Logarithmic Equations](https://algebrica.org/logarithmic-equations/)

1.2k

[Homogeneous Trigonometric Equations](https://algebrica.org/homogeneous-trigonometric-equations/)

1.7k

[Trigonometric Equations: Part 1](https://algebrica.org/trigonometric-equations/)

652

[Equations with Parameters](https://algebrica.org/equations-with-parameters/)

1.3k

[Linear Equations with Parameters](https://algebrica.org/linear-equations-with-parameters/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations-with-parameters.md?plain=1)
