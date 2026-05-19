> Content sourced from [Algebrica](https://algebrica.org/quadratic-equations-with-complex-solutions/) — CC BY-NC 4.0

## Real and complex solutions

The nature of the solutions of a [quadratic equation](<../quadratic-equations/>) is determined entirely by the sign of the discriminant \\(\Delta = b^2 - 4ac\\), which appears under the radical in the [quadratic formula](<../quadratic-formula/>):

\\[x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}\\]

When \\(\Delta > 0\\), the equation has two distinct real solutions. When \\(\Delta = 0\\), the two solutions coincide and the equation has a unique real root, counted with multiplicity two. When \\(\Delta < 0\\), the square root of a negative number is involved, and the equation has no real solutions. It does, however, admit exactly two solutions in the field of [complex numbers](<../complex-numbers-introduction/>), and these always appear as a conjugate pair.


By the [Fundamental Theorem of Algebra](<../roots-of-a-polynomial/>), a [polynomial](<../polynomials/>) equation of degree \\(n\\) has exactly \\(n\\) roots in \\(\mathbb{C}\\), counted with multiplicity. Applied to the quadratic case, this guarantees that every equation of the form \\(ax^2 + bx + c = 0\\) has exactly two roots in \\(\mathbb{C}\\), regardless of the sign of the discriminant.


Recall that given a complex number \\(z = a + bi\\), its [conjugate](<../complex-numbers-introduction/>) is defined as the following.

\\[\overline{z} = a - bi\\]

The number \\(\overline{z}\\) is represented in the Gaussian plane by the point symmetric to \\(z\\) with respect to the \\(x\\)-axis.

![](https://algebrica.org/wp-content/uploads/resources/images/quad-eq-complex-roots.png)

##### This geometric interpretation highlights how conjugate roots naturally arise when solving quadratic equations with negative discriminant, providing a visual link between algebraic expressions and their representation in the complex plane.

Every quadratic equation can be solved, regardless of whether its discriminant is positive, zero, or negative. When the discriminant happens to be negative, the equation simply has no real roots. However, it still possesses two legitimate solutions within the [complex number](<../complex-numbers-introduction/>) system. These solutions always arise as a pair of complex conjugates, reflecting the fundamental structure of quadratic equations in the complex plane.

## Complex conjugate roots

If the coefficients \\(a\\), \\(b\\), and \\(c\\) are [real](<../properties-of-real-numbers/>), the complex roots of the quadratic equation \\(ax^2 + bx + c = 0\\) always occur in conjugate pairs. Specifically, if \\(z \in \mathbb{C}\\) is a solution, then \\(\overline{z}\\) is also a solution. This result follows directly from the properties of complex conjugation with respect to arithmetic operations.

Suppose \\(az^2 + bz + c = 0\\). Taking the conjugate of both sides yields the following:

\\[\overline{az^2 + bz + c} = \overline{0}\\]

Since conjugation distributes over sums and products, and since every real number coincides with its own conjugate, it follows that \\(\overline{a} = a\\), \\(\overline{b} = b\\), and \\(\overline{c} = c\\). The left-hand side therefore reduces to the following.

\\[a\overline{z}^2 + b\overline{z} + c = 0\\]

This is the original equation evaluated at \\(\overline{z}\\), which is therefore a root of the same equation.

## Example 1

Solve the quadratic equation:

\\[x^2 + 4x + 5 = 0\\]

The equation is already in standard form \\(ax^2 + bx + c = 0\\), with coefficients \\(a = 1\\), \\(b = 4\\), and \\(c = 5\\). Substituting into the [quadratic formula](<../quadratic-formula/>) yields the following.

\\[x_{1,2} = \frac{-4 \pm \sqrt{4^2 - 4 \cdot 1 \cdot 5}}{2 \cdot 1}\\]

The discriminant is negative:

\\[\Delta = 16 - 20 = -4\\]

Since \\(\Delta < 0\\), the equation admits no real solutions. Proceeding with the quadratic formula and writing \\(\sqrt{-4} = 2i\\):

\\[x_{1,2} = \frac{-4 \pm 2i}{2}\\]

Simplifying, the two complex conjugate solutions are the following.

\\[x_1 = -2 + i \qquad x_2 = -2 - i\\]

## Example 2

The following example illustrates the case in which the discriminant is negative and the coefficients are not all integers, requiring an additional simplification step before the solutions can be expressed in standard complex form.

\\[3x^2 - 2x + 4 = 0\\]

The equation is in standard form with coefficients \\(a = 3\\), \\(b = -2\\), and \\(c = 4\\). Substituting into the [quadratic formula](<../quadratic-formula/>) yields the following.

\\[x_{1,2} = \frac{2 \pm \sqrt{(-2)^2 - 4 \cdot 3 \cdot 4}}{2 \cdot 3}\\]

Computing the discriminant:

\\[\Delta = 4 - 48 = -44\\]

Since \\(\Delta < 0\\), the equation admits no real solutions. Writing \\(\sqrt{-44} = 2\sqrt{11}\, i\\):

\\[x_{1,2} = \frac{2 \pm 2\sqrt{11}\, i}{6} = \frac{1 \pm \sqrt{11}\, i}{3}\\]

The solution is:

\\[x_1 = \frac{1 + \sqrt{11}\, i}{3} \qquad x_2 = \frac{1 - \sqrt{11}\, i}{3}\\]

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

651

[Equations with Parameters](https://algebrica.org/equations-with-parameters/)

1.3k

[Linear Equations with Parameters](https://algebrica.org/linear-equations-with-parameters/)

0 comments[](https://github.com/antoniolupetti/algebrica/blob/main/equations/linear-equations-with-parameters.md?plain=1)

1.6k

[Quadratic Equations with Parameters](https://algebrica.org/quadratic-equations-with-parameters/)
