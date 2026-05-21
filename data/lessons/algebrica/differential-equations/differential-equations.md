> Content sourced from [Algebrica](https://algebrica.org/differential-equations/) — CC BY-NC 4.0

## What Are Differential Equations?

An ordinary differential equation is a relation involving a function \\( f(x) \\), defined on an interval \\( I \subset \mathbb{R} \\), together with a finite number of its [derivatives](<../derivatives>). The equation is required to hold for all \\( x \\) in the interval \\( I \\). The general form of an ordinary differential equation of order \\( n \\) is:

\\[F\left(x, f(x), f’(x), f’'(x), \dots, f^{(n)}(x)\right) = 0 \\]

where \\( F \\) is a given function that depends on the independent variable \\( x \\), the unknown function \\(y = f(x) \\), and its derivatives up to order \\( n \\). An example of a differential equation is:

\\[f’(x) + x f(x) = 0 \quad \text{or} \quad y’ + x y = 0 \\]

## Order of a differential equation

The order of a differential equation is defined as the order of the highest derivative that appears in the equation. For example, the following differential equation is of third order because the third derivative of \\( y \\) is the highest-order derivative appearing in the equation:

\\[y’’ - 3y’ + 2y = 0 \\]

Differential equations are often considered challenging. For this reason, we will introduce the concepts gradually, with the goal of building a solid understanding of the methods used to solve them.

## Solutions of a differential equation

In general, a differential equation may admit infinitely many solutions. A general solution is a family of functions depending on \\( n \\) arbitrary constants, where \\( n \\) corresponds to the order of the equation. This family includes every possible solution of the differential equation.

  * Any function that satisfies the differential equation is called a solution, or integral of the equation.

  * A particular solution of a differential equation is a specific function obtained by assigning concrete values to the arbitrary constants in the general solution.

  * Initial conditions are constraints that allow us to select a particular solution from the general family.


## Forms and Types of Differential Equations

A differential equation is said to be in normal form when it can be written as:

\\[y^{(n)} = F(x, y, \dots, y^{(n-1)}) \\]

That is, the equation is explicitly solved with respect to the highest-order derivative.


A differential equation is said to be autonomous if the independent variable does not appear explicitly in its expression. The equation\\(y’ = y^2 - 1\\) is autonomous, since the right-hand side depends only on \\( y \\), not on the independent variable \\( x \\).


Given a differential equation along with a set of initial conditions, the problem of finding a solution that satisfies both the equation and the specified conditions is called a Cauchy problem. A Cauchy problem for a first-order differential equation is typically written as:

\\[\begin{cases} y’(x) = f(x, y(x)) \\\\[0.5em] y(x_0) = y_0 \end{cases} \\]

where \\( y’(x) = f(x, y(x)) \\) is the differential equation, and \\( y(x_0) = y_0 \\) is the initial condition, specifying the value of the solution at \\( x = x_0 \\).

## How to solve simple differential equations

Let’s begin by solving very simple differential equations, specifically first-order equations of the form \\( y’ = f(x) \\). In this case, the solution is obtained by [integrating](<../indefinite-integrals/>) both sides of the equation, yielding:

\\[\int y’ \, dx = \int f(x) \, dx \\]

For this type of equation, the general solution is given by:

\\[y(x) = \int f(x) \,dx + c \\]

where \\( c \in \mathbb{R} \\) is an arbitrary constant.

## Example 1

Let’s solve the differential equation:

\\[y’ - 3x = 0 \\]


Let’s rewrite it in standard form and integrate both sides:

\\[y’ = 3x \\]

\\[\int y’ \,dx = \int 3x \, dx \rightarrow y(x) = \frac{3}{2}x^2 + c \\]


The solutions are represented by the integral curves given by the equation:

\\[y(x) = \frac{3}{2}x^2 + c\\]

![](/diagrams/algebrica/differential-equations-1-1.png)


##### Each curve represents a particular solution of the general form. The value of \\( C \\) determines the vertical position (a vertical translation) of the curve. All the curves share the same upward-opening [parabolic](<../parabola>) shape, but they are vertically shifted depending on the value of \\( C \\). Together, this family of curves represents the complete set of solutions to the differential equation.

The solution is:

\\[y(x) = \frac{3}{2}x^2 + c\\]

## Separable Differential Equations

A first-order differential equation is said to be separable if it can be rewritten in the form:

\\[y’ = a(x) b(y) \\]

To find the general integral of this type of equation, we begin by dividing both sides of the equation by \\( b(y) \\). We obtain:

\\[\frac{y’}{b(y)} = a(x) \\]

We then integrate both sides with respect to \\( x \\), obtaining:

\\[\int \frac{y(x)'}{b(y(x))}\, dx = \int a(x) \, dx + c \\]

By substituting \\( y = y(x) \\), we obtain:

\\[\int \frac{dy}{b(y)}\, dx = \int a(x) \, dx + c \\]

## Example 2

Let’s solve the differential equation:

\\[y’ = xy\\]


We have a separable differential equation. We rewrite all the terms involving \\( y \\) on one side and those involving \\( x \\) on the other:

\\[\frac{1}{y} \, dy = x \, dx \\]


We integrate both sides and obtain:

\\[\int \frac{1}{y} \, dy = \int x \, dx \\]

\\[\ln |y| = \frac{x^2}{2} + c \\]

##### For a complete overview of integration rules and common integrals, see the [dedicated entry](<../indefinite-integral>) on the topic.


To isolate \\( y \\), we eliminate the logarithm by exponentiating both sides:

\\[|y| = e^{\frac{x^2}{2} + c} = e^c \cdot e^{\frac{x^2}{2}} = c e^{\frac{x^2}{2}} \\]

where \\( c \\) is a constant representing \\( e^c \\).

Since \\( c \\) is an arbitrary constant in \\( \mathbb{R} \\), which may be either positive or negative, we can write the general solution of the equation as:

\\[y(x) = c e^{\frac{x^2}{2}}, \quad c \in \mathbb{R} \\]


The solutions are represented by the following integral curves:

![](/diagrams/algebrica/differential-equations-2.png)

Therefore, the general solution is:

\\[y(x) = c e^{\frac{x^2}{2}}, \quad c \in \mathbb{R} \\]

## Glossary

  * Ordinary differential equation: a relation involving a function defined on an interval and a finite number of its derivatives, which holds true at all points in the interval.

  * Order of a differential equation: the order of the highest derivative present in the equation.

  * General solution: a family of functions that satisfies the differential equation and contains arbitrary constants equal in number to the order of the equation.

  * Particular solution: a specific function obtained from the general solution by assigning concrete values to the arbitrary constants.

  * Initial conditions: conditions specifying the value of the solution and possibly its derivatives at a particular point, used to determine a unique particular solution.

  * Cauchy problem: the task of finding a solution to a differential equation that satisfies given initial conditions, typically expressed as: \\[y(x_0) = y_0,\quad y’(x_0) = y_1, \ldots \\]

  * Integral curves: the graphs of the solutions of a differential equation, representing how the function evolves along the plane.
