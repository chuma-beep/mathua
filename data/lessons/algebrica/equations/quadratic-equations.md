# Quadratic Equations

Source: algebrica.org — CC BY-NC 4.0
https://algebrica.org/quadratic-equations/

## Introduction

A quadratic equation is a second-degree [polynomial equation](../polynomial-equations/)
in one variable. Its standard form is the following:

\\[ ax^2 + bx + c = 0 \\]

where \\(a\\), \\(b\\), and \\(c\\) are real coefficients, \\(x\\) is the unknown, and \\(a \neq 0\\).

- The coefficients \\(a\\), \\(b\\), and \\(c\\) are constants.
- \\(x\\) represents the variable.
- \\(a\\) is the coefficient of the quadratic term \\(x^2\\), \\(b\\) the coefficient of the linear term \\(x\\) and \\(c\\) the constant term.
- When \\(a = 0\\), the equation reduces to the linear equation \\(bx + c = 0\\). If \\(b = 0\\) as well, the equation becomes constant and may have no solution or infinitely many solutions, depending on whether \\(c \neq 0\\) or \\(c = 0\\).

Quadratic equations are the simplest case of [trinomial equation](../trinomial-equations/), which has the general form:

\\[ ax^{2n} + bx^{n} + c = 0 \\]

Setting \\(n = 1\\) recovers the standard quadratic form \\(ax^2 + bx + c = 0.\\) For \\(n \geq 2\\) the equation can be reduced to a quadratic in the auxiliary variable \\(y = x^n\\) and solved with the same techniques.

- - -
## Geometrical interpretation

The equation \\(y = ax^2 + bx + c\\), with \\(a \neq 0\\), represents a [parabola](../parabola/) in the plane defined by the variables \\(x\\) and \\(y\\).

![Img. 1](quadratic-equations.svg)

When the coefficient \\(a > 0\\) the parabola opens upward and the vertex is a minimum of the function. When \\(a < 0\\) it opens downward and the vertex is a maximum.

The real solutions of the equation \\(ax^2 + bx + c = 0\\) correspond to the points at which the parabola intersects the \\(x\\)-axis. The sign of the [discriminant](../quadratic-formula/) \\(\Delta = b^2 - 4ac\\) controls this configuration:

+ For \\(\Delta > 0\\) the parabola crosses the axis at two distinct points.
+ For \\(\Delta = 0\\) it is tangent to the axis.
+ For \\(\Delta < 0\\) it does not intersect the axis and the equation has no real solutions.

> The condition \\(a \neq 0\\) ensures that the equation describes a parabolic curve rather than a [linear equation](../linear-equation/).

- - -
## Resolution methods

A quadratic equation is [incomplete](../incomplete-quadratic-equations/) when either the
coefficient \\(b\\) or \\(c\\) is equal to zero. In this case the equation takes a simpler form and can be solved directly, without applying the general formula. The first step in solving a quadratic equation is to rewrite it in standard form:

\\[ax^2 + bx + c = 0\\]

This form allows the coefficients to be identified and the [discriminant](../quadratic-formula/) \\(\Delta = b^2 - 4ac\\) to be computed. The discriminant determines the nature of the solutions:

+ Two distinct real roots when \\(\Delta > 0\\).
+ One real root of multiplicity two when \\(\Delta = 0\\).
+ A pair of complex conjugate roots when \\(\Delta < 0\\).

The most general method of resolution is the [quadratic formula](../quadratic-formula/).
In some cases, however, [factoring](../factoring-quadratic-equations/) or [completing the square](../completing-the-square) can offer a more direct route to the solution.

The [fundamental theorem of algebra](../roots-of-a-polynomial/) guarantees that a quadratic equation has exactly two roots in \\(\mathbb{C}\\), counted with multiplicity. The roots are both real when \\(\Delta \geq 0\\), and form a pair of [complex conjugates](../quadratic-equations-with-complex-solutions/)when \\(\Delta < 0\\).

- - -
## Quadratic formula

 Given a quadratic equation in the standard form \\(ax^2+bx+c = 0\\), the [quadratic formula](../quadratic-equations/quadratic-formula/) is:

\\[ x_{1,2} = \frac{{-b \pm \sqrt{{b^2 - 4ac}}}}{{2a}} \\]

+ \\(a\\), \\(b\\), \\(c\\) are real coefficients and \\(a \neq 0\\).
+ The \\(\pm\\) symbol reflects the existence of two solutions, corresponding to the two signs.
+ A quadratic equation has exactly two roots in \\(\mathbb{C}\\), counted with multiplicity.
+ By Vieta's formulas, the roots satisfy \\(x_1 + x_2 = -b/a\\) and \\(x_1 \cdot x_2 = c/a\\).

A further property of the discriminant is the following:

\\[ \Delta = a^2(x_1 - x_2)^2 \\]

This identity shows directly that \\(\Delta \geq 0\\) when the roots are real, and that
\\(\Delta = 0\\) if and only if the two roots coincide.

>When the discriminant is negative, the solutions are complex. The dedicated entry on [quadratic equations with complex solutions](../quadratic-equations-with-complex-solutions/) covers this case in full.

- - -
## Factoring

A quadratic equation can be [factored](../factoring-quadratic-equations/) into the following form:

\\[ ax^2 + bx + c = 0 \quad \Longleftrightarrow \quad a(x - x_1)(x - x_2) = 0 \\]

where \\(x_1\\) and \\(x_2\\) are the roots of the equation. By [Vieta's formulas](../trinomials/), the roots satisfy the following relations:

\\[ x_1 + x_2 = -\frac{b}{a}\\]
\\[ x_1 \cdot x_2 = \frac{c}{a} \\]

This method is effective when the roots can be identified by inspection or by simple
trial, but becomes impractical for equations with irrational or complex roots, where the
[quadratic formula](../quadratic-formula/) is preferable.

- - -
## How to solve a quadratic equation

+ Rewrite the equation in standard form: \\( ax^2 + bx + c = 0 \\).  
+ Calculate the discriminant: \\( \Delta = b^2 - 4ac \\).  
+ Use the quadratic formula:  
   \\[
   x = \frac{-b \pm \sqrt{\Delta}}{2a}
   \\]  
+ Simplify the result.  
+ If \\( \Delta \geq 0 \\), the solutions are real; if \\( \Delta < 0 \\), the solutions are complex conjugates.

> The quadratic formula is the universal method, but not always the most efficient. When the equation is incomplete or admits an obvious factorization, the roots can often be obtained more quickly by direct inspection.

- - -
## Quadratic equations with parameters

A natural extension of the study of quadratic equations is to consider the case in which the coefficients are not fixed numbers but depend on an external parameter. In this setting we speak of [quadratic equations with a parameter](../quadratic-equations-with-parameters/),
also called literal quadratic equations, which take the form:

\\[ a(k)\\,x^2 + b(k)\\,x + c(k) = 0, \quad a(k) \neq 0 \\]

Varying the parameter \\(k\\) alters the equation and, consequently, the nature of its solutions. The analysis relies on the discriminant:

\\[ \Delta(k) = b(k)^2 - 4\\,a(k)\\,c(k) \\]

which, exactly as in the classical case, determines whether the equation admits two
distinct real solutions, a repeated solution, or a pair of complex conjugate solutions.

- - -
## Flowchart

- `Quadratic equation to solve`
  - **IF** \\( a = 0 \\)
    - _reduce to a linear equation_
      the leading coefficient vanishes, so the equation is no longer quadratic: \\( bx + c = 0 \\)
    - _solve directly_
      \\[x = -\frac{c}{b}\\]
  - `ELSE IF` \\( b = 0 \\) and \\( c = 0 \\)
    - _the equation reduces to \\( ax^2 = 0 \\)_
      both the linear and constant terms vanish, leaving a perfect square
    - _conclude_ \\( x = 0 \\) is a double root
  - `ELSE IF` \\( b = 0 \\) and \\( c \neq 0 \\)
    - _the equation reduces to \\( ax^2 + c = 0 \\)_
      the linear term vanishes; the equation is a pure quadratic in \\( x^2 \\)
    - _solve for \\( x^2 \\)_
      \\[x^2 = -\frac{c}{a}\\]
    - **IF** \\( -c/a < 0 \\)
      - _no real solutions_
        the right-hand side is negative, so no real number satisfies the equation
    - `ELSE`
      - _two real solutions_
        \\[x = \pm\sqrt{-\frac{c}{a}}\\]
  - `ELSE IF` \\( c = 0 \\)
    - _factor out \\( x \\)_
      the constant term vanishes, so \\( x \\) is an immediate common factor: \\( x(ax + b) = 0 \\)
    - _two solutions_
      \\[x = 0 \qquad x = -\frac{b}{a}\\]
  - `ELSE IF` the equation is easily factorable
    - _look for two numbers \\( p \\) and \\( q \\) such that_
      \\[p + q = b \qquad p \cdot q = ac\\]
    - _rewrite and factor_
      if such numbers exist, the trinomial splits cleanly: \\( ax^2 + bx + c = a(x - r_1)(x - r_2) \\)
    - _read off the roots_ \\( x = r_1 \\) and \\( x = r_2 \\)
  - `ELSE`
    - _compute the discriminant_
      when no shortcut applies, the quadratic formula is the general method: \\( \Delta = b^2 - 4ac \\)
    - **IF** \\( \Delta < 0 \\)
      - _no real solutions_
        the square root of a negative number is not real; the equation has two complex conjugate roots
    - **IF** \\( \Delta = 0 \\)
      - _one real solution_ (double root)
        the two roots coincide; the parabola is tangent to the \\( x \\)-axis
        \\[x = -\frac{b}{2a}\\]
    - **IF** \\( \Delta > 0 \\)
      - _two distinct real solutions_
        the parabola crosses the \\( x \\)-axis at two separate points
        \\[x = \frac{-b \pm \sqrt{\Delta}}{2a}\\]