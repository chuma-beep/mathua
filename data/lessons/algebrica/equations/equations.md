# Equations

Source: algebrica.org — CC BY-NC 4.0
https://algebrica.org/equations/

## What are equations

An equation is a mathematical statement asserting that two expressions take the same value, typically written in the form \\( F(x) = G(x) \\), where one or more variables act as unknowns ranging over some specified [set](../types-of-numbers/), such as \\( \mathbb{R} \\) or \\( \mathbb{C} \\). To solve an equation is to determine every value of the variable, within the appropriate [domain](../determining-the-domain-of-a-function/), for which the equality holds. Many equations can be reorganized into the form:

\\[
F(x) = 0
\\]

a representation that simplifies their study, highlights their structure, and provides a unified framework for approaching different methods of solution.

- - -
## The solution to an equation

The solution to an equation is any value of the variables that makes the equality true. Depending on the equation and its structure, the solution set may contain:

+ a unique solution  
+ multiple solutions  
+ infinitely many solutions  
+ no solution at all  

> In some cases, such as with trigonometric equations, the set of solutions can be infinite because the functions involved are periodic.

In other contexts, especially in polynomial, exponential, or transcendental equations, the solutions may not lie within the real numbers. When no real value satisfies the equation, the natural setting becomes the [complex plane](../complex-numbers-introduction/), where solutions can still exist and be meaningfully interpreted.

A representative instance is the quadratic equation with negative discriminant, which admits no real roots but always has two complex conjugate solutions, as discussed in detail in the page on [quadratic equations with complex solutions](../quadratic-equations-with-complex-solutions/).

- - -
## Equivalent equations and admissible operations

Two equations are considered equivalent if they possess identical solution sets. The objective in solving an equation is to apply a sequence of transformations that maintain this equivalence, thereby reducing the equation to a simpler form in which the solutions are readily apparent.

Certain operations are guaranteed to yield an equivalent equation. Adding or subtracting the same expression on both sides, or multiplying both sides by a nonzero constant, does not alter the solution set. These manipulations form the basis of most elementary solution techniques. As a simple illustration, the equation \\( 2x + 4 = 0 \\) is equivalent to \\( x + 2 = 0 \\), obtained by dividing both sides by \\( 2 \\), and both equations share the unique solution \\( x = -2 \\).

Other operations can disrupt equivalence in less obvious ways. Multiplying both sides by an expression involving the variable may introduce extraneous solutions if the expression equals zero for some value of \\( x \\). Similarly, squaring both sides, a common technique for [irrational equations](../irrational-equations/), can produce solutions that satisfy the transformed equation but not the original. Conversely, dividing both sides by a variable expression may eliminate solutions at points where the expression is zero, a phenomenon addressed in detail in the discussion on [loss of roots](../loss-of-roots/).

- - -
## Algebraic equations

Algebraic equations are equations in which both sides consist entirely of [polynomials](../polynomials/). A polynomial in one variable is a formal expression of the form

\\[
P(x) = a_{n}x^{n} + a_{n-1}x^{n-1} + \cdots + a_{1}x + a_{0}
\\]

where \\( n \\) is a non-negative integer, the coefficients \\( a_0, a_1, \ldots, a_n \\) belong to some fixed field (typically \\( \mathbb{R} \\) or \\( \mathbb{C} \\)), and \\( a_n \neq 0 \\) when \\( n \geq 1 \\). The [integer](../integers/) \\( n \\) is the degree of the polynomial. An algebraic equation then takes the form \\( P(x) = Q(x) \\), or equivalently \\( P(x) - Q(x) = 0 \\), which reduces the problem to finding the roots of a single polynomial. Algebraic equations are classified according to the degree of the polynomial involved.

[Linear equations](../linear-equations/) have degree \\( 1 \\). They involve the variable raised to no power higher than the first, and their solution is always unique when the leading coefficient is nonzero.

[Quadratic equations](../quadratic-equations/) have degree \\( 2 \\) and take the standard form \\( ax^2 + bx + c = 0 \\), with \\( a \neq 0 \\). The nature of their solutions is governed by the [discriminant](../quadratic-formula) \\( \Delta = b^2 - 4ac \\).

+ If \\( \Delta > 0 \\) there are two distinct real solutions.
+ If \\( \Delta = 0 \\) there is one repeated real solution.
+ If \\( \Delta < 0 \\) the solutions are a pair of complex conjugates.

Cubic equations have degree \\( 3 \\) and take the general form \\( ax^3 + bx^2 + cx + d = 0 \\). By the [Fundamental Theorem of Algebra](../roots-of-a-polynomial/), every polynomial of degree \\( n \\) with coefficients in \\( \mathbb{C} \\) has exactly \\( n \\) roots in \\( \mathbb{C} \\), counted with multiplicity. A cubic equation therefore has exactly three roots in \\( \mathbb{C} \\), which may be all real, or one real and two complex conjugates.

Equations of degree higher than three follow the same principle: a [polynomial equation](../polynomial-equations/) of degree \\( n \\) has exactly \\( n \\) roots in \\( \mathbb{C} \\), counted with multiplicity, though explicit formulas for the roots in terms of radicals exist only for degrees up to four, a fact made precise by the theory of Galois.

Among higher-degree equations, [binomial](../binomial-equations/) and [trinomial equations](../trinomial-equations) deserve particular attention. These are equations of degree greater than two that contain only two or three distinct terms, respectively, and can often be solved by substitution or by factoring into lower-degree polynomials.

- - -
## Rational equations

[Rational equations](../rational-equations/) are equations that contain at least one fractional expression whose numerator and denominator are polynomials. In their most general form, they involve a ratio of polynomials on both sides, and can always be reduced to the form

\\[
\frac{P(x)}{Q(x)} = 0
\\]

by transferring all terms to one side and combining them over a common denominator, with the restriction that \\( Q(x) \neq 0 \\). A standard approach is to clear denominators by multiplying both sides by the least common multiple of all denominators in the equation, thereby transforming the problem into a polynomial equation.

> This step, however, requires careful attention: any value that makes a denominator vanish must be excluded from the solution set from the outset, and candidate solutions obtained after clearing must be checked against these excluded values.

- - -
## Irrational equations

[Irrational equations](../irrational-equations/) have variables inside a radical. Typically, such equations have one radical, for example:

\\[
\sqrt[n]{\\,f(x)\\,} = g(x)
\\]

where \\( f(x) \\) and \\( g(x) \\) are polynomials with real coefficients, and the standard technique consists of isolating the radical and raising both sides to the \\( n \\)-th power to eliminate it. When the equation contains more than one radical, this process must be applied iteratively. After each step, a new radical is isolated, and the procedure is repeated until none remain. Each time both sides are raised to a power, the transformation is not equivalence-preserving and may introduce extraneous solutions.

It is therefore necessary to verify every candidate solution in the original equation, and this requirement becomes increasingly important as the number of radicals grows.

- - -
## Absolute value equations

[Absolute value equations](../absolute-value-equations/) are equations in which the unknown appears inside an [absolute value](../absolute-value/) expression. The simplest case takes the form \\( |x| = a \\), where \\( a \\) is a real constant, and its solution set depends entirely on the sign of \\( a \\).

+ If \\( a > 0 \\), the equation has exactly two solutions, \\( x = a \\) and \\( x = -a \\), since both values have the same distance from the origin on the real line.
+ If \\( a = 0 \\), the only solution is \\( x = 0 \\), as the absolute value of a real number vanishes if and only if the number itself is zero.
+ If \\( a < 0 \\), the equation has no solution in \\( \mathbb{R} \\), because the absolute value is by definition non-negative and cannot equal a negative constant.

> Absolute value equations that contain polynomial or rational expressions often require case analysis based on the sign of the inner expression. It is essential to verify that each candidate solution satisfies the original equation, since the case-splitting procedure may introduce extraneous solutions.

- - -
## Transcendental equations

Transcendental equations are equations in which one or more variables appear within transcendental functions, such as [exponential](../exponential-function/), [logarithmic](../logarithmic-function/), or [trigonometric functions](../sine-function/), that cannot be expressed as finite combinations of algebraic operations. These equations go beyond polynomial or rational forms and frequently resist closed-form solution: in many cases no explicit formula for the roots exists, and one must resort to numerical methods to approximate them.

[Logarithmic equations](../logarithmic-equations/) are equations that involve a variable inside a logarithmic expression, such as \\( 2\log_2{(2 - x)^2} = 0 \\). These equations are typically solved by applying logarithmic properties, such as the power rule, the product and quotient rules, or the change of base formula, to simplify the expression and isolate the variable.

[Exponential equations](../exponential-equations/) are equations in which the variable appears in the exponent of an exponential expression. They typically involve forms like \\( a^x = b \\), where \\( a \\) and \\( b \\) are constants and \\( x \\) is the unknown.

[Trigonometric equations](../trigonometric-equations/) are equations that involve periodic trigonometric functions such as \\( \sin(x) \\), \\( \cos(x) \\), or \\( \tan(x) \\) containing a variable. Due to the periodic nature of these functions, trigonometric equations often have infinitely many solutions.

- - -
## Systems of equations

A system of equations arises when multiple equations must be satisfied simultaneously by the same set of unknowns. Such a system consists of two or more equations involving the same variables, and a solution is any assignment of values to those variables that satisfies all equations concurrently. Systems may be classified as [linear](../systems-of-linear-equations/) or nonlinear based on the form of their equations. Analysing these systems often requires advanced techniques beyond those used for single equations, including substitution, elimination, or matrix methods for linear systems.