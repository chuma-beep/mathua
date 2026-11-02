# Slope Fields and Solution Curves

A first-order differential equation can be read as a direction field. At every point where the right-hand side is defined, the equation gives the slope of the solution curve that passes through that point.

## Reading a slope field

For an equation written as $y' = f(x,y)$, evaluate $f$ at a point $(x_0,y_0)$. The number $f(x_0,y_0)$ is the slope of the short segment drawn at that point. A solution curve follows those segments, so a slope field lets us predict the shape of a solution before solving the equation. For example, $y' = x + y$ has slope 5 at $(2,3)$, because $2 + 3 = 5$.

## Equilibrium solutions

An equilibrium solution is constant. If $y = c$ is constant, then $y' = 0$, so every equilibrium value must satisfy $f(x,c) = 0$ for the relevant range of $x$. For the autonomous equation $y' = y(4-y)$, the equilibria are $y = 0$ and $y = 4$. Between them the slope is positive, and outside them it is negative, so the field shows solutions moving away from 0 and toward 4.

## Concavity from the differential equation

Differentiate $y' = f(x,y)$ with respect to $x$:

$$y'' = f_x(x,y) + f_y(x,y)y' = f_x(x,y) + f_y(x,y)f(x,y).$$

The sign of $y''$ tells us whether a solution curve is concave up or concave down. At $(2,3)$ for $y' = x + y$, we get $f_x = 1$ and $f_y = 1$, so $y'' = 1 + 5 = 6$. The solution curve is concave up there.
