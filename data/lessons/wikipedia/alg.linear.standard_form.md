> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linear_equation) — CC BY-SA 4.0

# Standard form of a linear equation Ax + By = C

## Two variables
A linear equation in two variables \(x\) and \(y\) can be written as \(ax+by+c=0,\) where \(a\) and  \(b\) are not both 0.

If \(a\) and \(b\) are real numbers, it has infinitely many solutions.

### Linear function

If *b* ≠ 0, the equation
\(ax+by+c=0\)
is a linear equation in the single variable \(y\) for every value of \(x\). It therefore has a unique solution for \(y\), which is given by
\(y=-\frac ab x-\frac cb.\)

This defines a function. The graph of this function is a line with slope \(-\frac ab\) and \(y\)-intercept \(-\frac cb.\) The functions whose graph is a line are generally called *linear functions* in the context of calculus. However, in linear algebra, a linear function is a function that maps a sum to the sum of the images of the summands. So, for this definition, the above function is linear only when *c* = 0, that is when the line passes through the origin. To avoid confusion, the functions whose graph is an arbitrary line are often called *affine functions*, and the linear functions such that *c* = 0 are often called *linear maps*.

### Geometric interpretation


Each solution (*x*, *y*) of a linear equation
\(ax+by+c=0\)
may be viewed as the Cartesian coordinates of a point in the Euclidean plane. With this interpretation, all solutions of the equation form a line, provided that \(a\) and \(b\) are not both zero. Conversely, every line is the set of all solutions of a linear equation.

The phrase "linear equation" takes its origin in this correspondence between lines and equations: a *linear equation* in two variables is an equation whose solutions form a line.

If *b* ≠ 0, the line is the graph of the function of \(x\) that has been defined in the preceding section. If *b* = 0, the line is a *vertical line* (that is a line parallel to the \(y\)-axis) of equation \(x=-\frac ca,\) which is not the graph of a function of \(x\).

Similarly, if *a* ≠ 0, the line is the graph of a function of \(y\), and, if *a* = 0, one has a horizontal line of equation \(y=-\frac cb.\)

### Equation of a line
There are various ways of defining a line. In the following subsections, a linear equation of the line is given in each case.

#### Slope–intercept form or Gradient-intercept form
A non-vertical line can be defined by its slope \(m\), and its \(y\)-intercept *y*{sub|0} (the \(y\) coordinate of its intersection with the \(y\)-axis). In this case, its *linear equation* can be written
\(y=mx+y_0.\)

If, moreover, the line is not horizontal, it can be defined by its slope and its \(x\)-intercept *x*{sub|0}. In this case, its equation can be written
\(y=m(x-x_0),\)
or, equivalently,
\(y=mx-mx_0.\)

These forms rely on the habit of considering a nonvertical line as the graph of a function. For a line given by an equation
\(ax+by+c = 0,\)
these forms can be easily deduced from the relations
\(\begin{align}
m&=-\frac ab,\\
x_0&=-\frac ca,\\
y_0&=-\frac cb.
\end{align}\)

#### Point–slope form or Point-gradient form
A non-vertical line can be defined by its slope \(m\), and the coordinates \(x_1, y_1\) of any point of the line. In this case, a linear equation of the line is
\(y=y_1 + m(x-x_1),\)
or
\(y=mx +y_1-mx_1.\)

This equation can also be written
\(y-y_1=m(x-x_1)\)
to emphasize that the slope of a line can be computed from the coordinates of any two points.

#### Intercept form
A line that is not parallel to an axis and does not pass through the origin cuts the axes into two different points. The intercept values *x*{sub|0} and *y*{sub|0} of these two points are nonzero, and an equation of the line is
\(\frac{x}{x_0} + \frac{y}{y_0} = 1.\)
(It is easy to verify that the line defined by this equation has *x*{sub|0} and *y*{sub|0} as intercept values).

#### Two-point form
Given two different points (*x*{sub|1}, *y*{sub|1}) and (*x*{sub|2}, *y*{sub|2}), there is exactly one line that passes through them. There are several ways to write a linear equation of this line.

If *x*{sub|1} ≠ *x*{sub|2}, the slope of the line is \(\frac{y_2 - y_1}{x_2 - x_1}.\) Thus, a point-slope form is
\(y - y_1 = \frac{y_2 - y_1}{x_2 - x_1} (x - x_1).\)

By clearing denominators, one gets the equation
\((x_2 - x_1)(y - y_1) - (y_2 - y_1)(x - x_1)=0,\)
which is valid also when *x*{sub|1} = *x*{sub|2} (to verify this, it suffices to verify that the two given points satisfy the equation).

This form is not symmetric in the two given points, but a symmetric form can be obtained by regrouping the constant terms:
\((y_1-y_2)x + (x_2-x_1)y + (x_1y_2 - x_2y_1) =0\)
(exchanging the two points changes the sign of the left-hand side of the equation).

#### Determinant form
The two-point form of the equation of a line can be expressed simply in terms of a determinant. There are two common ways for that.

The equation \((x_2 - x_1)(y - y_1) - (y_2 - y_1)(x - x_1)=0\) is the result of expanding the determinant in the equation
\(\begin{vmatrix}x-x_1&y-y_1\\x_2-x_1&y_2-y_1\end{vmatrix}=0.\)

The equation \((y_1-y_2)x + (x_2-x_1)y + (x_1y_2 - x_2y_1)=0\) can be obtained by expanding with respect to its first row the determinant in the equation
\(\begin{vmatrix}
x&y&1\\
x_1&y_1&1\\
x_2&y_2&1
\end{vmatrix}=0.\)

Besides being very simple and mnemonic, this form has the advantage of being a special case of the more general equation of a hyperplane passing through \(n\) points in a space of dimension *n* − 1. These equations rely on the condition of linear dependence of points in a projective space.

