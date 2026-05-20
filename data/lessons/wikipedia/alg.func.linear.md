> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Linear_function_%28calculus%29) — CC BY-SA 4.0

# Linear functions and their graphs

In calculus and related areas of mathematics, a **linear function** from the real numbers to the real numbers is a function whose graph (in Cartesian coordinates) is a non-vertical line in the plane.
The characteristic property of linear functions is that when the input variable is changed, the change in the output is proportional to the change in the input.

Linear functions are related to linear equations.

## Properties
A linear function is a polynomial function in which the variable  has degree at most one (a *linear polynomial*):
\(f(x)=ax+b\).
Such a function is called *linear* because its graph, the set of all points \((x,f(x))\) in the Cartesian plane, is a line. The coefficient *a* is called the *slope* of the function and of the line (see below).

If the slope is \(a=0\), this is a *constant function* \(f(x)=b\) defining a horizontal line, which some authors exclude from the class of linear functions. With this definition, the degree of a linear polynomial would be exactly one, and its graph would be a line that is neither vertical nor horizontal. However, in this article, \(a\neq 0\) is not required, so constant functions will be considered linear.

If \(b=0\) then the linear function is said to be *homogeneous*. Such function defines a line that passes through the origin of the coordinate system, that is, the point \((x,y)=(0,0)\). In advanced mathematics texts, the term *linear function* often denotes specifically homogeneous linear functions, while the term affine function is used for the general case, which includes \(b\neq0\).

The natural domain of a linear function \(f(x)\), the set of allowed input values for *x*, is the entire set of real numbers, \(x\in \mathbb R.\) One can also consider such functions with *x* in an arbitrary field, taking the coefficients *a, b* in that field.

The graph \(y=f(x)=ax+b\) is a non-vertical line having exactly one intersection with the *y*-axis, its *y*-intercept point \((x,y)=(0,b).\) The *y*-intercept value \(y=f(0)=b\) is also called the *initial value* of \(f(x).\) If \(a\neq 0,\) the graph is a non-horizontal line having exactly one intersection with the *x*-axis, the *x*-intercept point \((x,y)=(-\tfrac ba,0).\) The *x*-intercept value \(x=-\tfrac ba,\) the solution of the equation \(f(x)=0,\) is also called the *root* or *zero* of \(f(x).\)

## Slope

The slope of a nonvertical line is a number that measures how steeply the line is slanted (rise-over-run). If the line is the graph of the linear function \(f(x) = ax + b\), this slope is given by the constant .

The slope measures the constant rate of change of \(f(x)\) per unit change in *x*: whenever the input  is increased by one unit, the output changes by  units: \(f(x{+}1)=f(x)+a\), and more generally \(f(x{+}\Delta x)=f(x)+a\Delta x\) for any number \(\Delta x\). If the slope is positive, \(a > 0\), then the function \(f(x)\) is increasing; if \(a < 0\), then \(f(x)\) is decreasing

In calculus, the derivative of a general function measures its rate of change. A linear function \(f(x)=ax+b\) has a constant rate of change equal to its slope , so its derivative is the constant function \(f\,'(x)=a\).

The fundamental idea of differential calculus is that any smooth function \(f(x)\) (not necessarily linear) can be closely approximated near a given point \(x=c\) by a unique linear function. The derivative \(f\,'(c)\) is the slope of this linear function, and the approximation is: \(f(x) \approx f\,'(c)(x{-}c)+f(c)\) for \(x\approx c\). The graph of the linear approximation is the tangent line of the graph \(y=f(x)\) at the point \((c,f(c))\). The derivative slope \(f\,'(c)\) generally varies with the point *c*. Linear functions can be characterized as the only real functions whose derivative is constant: if \(f\,'(x)=a\) for all *x*, then \(f(x)=ax+b\) for \(b=f(0)\).

## Slope-intercept, point-slope, and two-point forms
A given linear function \(f(x)\) can be written in several standard formulas displaying its various properties. The simplest is the *slope-intercept form*:
\(f(x)= ax+b\),
from which one can immediately see the slope *a* and the initial value \(f(0)=b\), which is the *y*-intercept of the graph \(y=f(x)\).

Given a slope *a* and one known value \(f(x_0)=y_0\), we write the *point-slope form*:
\(f(x) = a(x{-}x_0)+y_0\).
In graphical terms, this gives the line \(y=f(x)\) with slope *a* passing through the point \((x_0,y_0)\).

The *two-point form* starts with two known values \(f(x_0)=y_0\) and \(f(x_1)=y_1\). One computes the slope \(a=\tfrac{y_1-y_0}{x_1-x_0}\) and inserts this into the point-slope form:
\(f(x) = \tfrac{y_1-y_0}{x_1-x_0}(x{-}x_0\!) + y_0\).
Its graph \(y=f(x)\) is the unique line passing through the points \((x_0,y_0\!), (x_1,y_1\!)\). The equation \(y=f(x)\) may also be written to emphasize the constant slope:
\(\frac{y-y_0}{x-x_0}=\frac{y_1-y_0}{x_1-x_0}\).

## Relationship with linear equations

Linear functions commonly arise from practical problems involving variables \(x,y\) with a linear relationship, that is, obeying a linear equation \(Ax+By=C\). If  \(B\neq 0\), one can solve this equation for *y*, obtaining
\(y = -\tfrac{A}{B}x +\tfrac{C}{B}=ax+b,\)
where we denote \(a=-\tfrac{A}{B}\) and \(b=\tfrac{C}{B}\). That is, one may consider *y* as a dependent variable (output) obtained from the independent variable (input) *x* via a linear function: \(y = f(x) = ax+b\). In the *xy*-coordinate plane, the possible values of \((x,y)\) form a line, the graph of the function \(f(x)\).  If \(B=0\) in the original equation, the resulting line \(x=\tfrac{C}{A}\) is vertical, and cannot be written as \(y=f(x)\).

The features of the graph \(y = f(x) = ax+b\) can be interpreted in terms of the variables *x* and *y*. The *y*-intercept is the initial value \(y=f(0)=b\) at \(x=0\). The slope *a* measures the rate of change of the output *y* per unit change in the input *x*. In the graph, moving one unit to the right (increasing *x* by 1) moves the *y*-value up by *a*: that is, \(f(x{+}1) = f(x) + a\). Negative slope *a* indicates a decrease in *y* for each increase in *x*.

For example, the linear function \(y = -2x + 4\) has slope \(a=-2\), *y*-intercept point \((0,b)=(0,4)\), and *x*-intercept point \((2,0)\).

### Example
Suppose salami and sausage cost €6 and €3 per kilogram, and we wish to buy €12 worth. How much of each can we purchase? If *x* kilograms of salami and *y* kilograms of sausage costs a total of €12 then, €6×*x* + €3×*y* = €12. Solving for *y* gives the point-slope form \(y = -2x + 4\), as above. That is, if we first choose the amount of salami *x*, the amount of sausage can be computed as a function \(y = f(x) = -2x + 4\). Since salami costs twice as much as sausage, adding one kilo of salami decreases the sausage by 2 kilos: \(f(x{+}1) = f(x) - 2\), and the slope is −2. The *y*-intercept point \((x,y)=(0,4)\) corresponds to buying only 4 kg of sausage; while the *x*-intercept point \((x,y)=(2,0)\) corresponds to buying only 2 kg of salami.

Note that the graph includes points with negative values of *x* or *y*, which have no meaning in terms of the original variables (unless we imagine selling meat to the butcher). Thus we should restrict our function \(f(x)\) to the domain \(0\le x\le 2\).

Also, we could choose *y* as the independent variable, and compute *x* by the inverse linear function: \(x = g(y) = -\tfrac12 y +2\) over the domain \(0\le y \le 4\).

## Relationship with other classes of functions
If the coefficient of the variable is not zero (*a* ≠ 0), then a linear function is represented by a degree 1 polynomial (also called a *linear polynomial*), otherwise it is a constant function – also a polynomial function, but of zero degree.

A straight line, when drawn in a different kind of coordinate system may represent other functions.

For example, it may represent an exponential function when its values are expressed in the logarithmic scale. It means that when log(*g*(*x*)) is a linear function of , the function  is exponential. With linear functions, increasing the input by one unit causes the output to increase by a fixed amount, which is the slope of the graph of the function. With exponential functions, increasing the input by one unit causes the output to increase by a fixed multiple, which is known as the base of the exponential function.

If *both* arguments and values of a function are in the logarithmic scale (i.e., when log(*y*) is a linear function of log(*x*)), then the straight line represents a power law:
\(\log_r y = a \log_r x + b \quad\Rightarrow\quad y = r^b\cdot x^a\)


On the other hand, the graph of a linear function in terms of polar coordinates:
\(r =f(\theta ) = a\theta  + b\)
is an Archimedean spiral if \(a \neq 0\) and a circle otherwise.

## See also
* Affine map, a generalization
* Arithmetic progression, a linear function of integer argument

## Notes


## References
*
*

## External links
* https://web.archive.org/web/20130524101825/http://www.math.okstate.edu/~noell/ebsm/linear.html
* https://web.archive.org/web/20180722042342/https://corestandards.org/assets/CCSSI_Math%20Standards.pdf

