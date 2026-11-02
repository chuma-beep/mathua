# Least-Squares Solutions

## When a system has no solution

An overdetermined system usually has no exact solution, because the right-hand side does not lie in the column space of the coefficient matrix. Least squares replaces the impossible exact solution with the vector in the column space that is closest to the data.

This section develops the idea and connects it to the surrounding linear algebra.

- Recognise an inconsistent overdetermined system
- Describe the least-squares goal
- Distinguish exact and approximate solutions

## The normal equations

The normal equations express the optimality condition that the residual is orthogonal to the column space. Solving them gives the projection of the data onto the column space and the coefficient vector of the best-fitting model.

This section develops the idea and connects it to the surrounding linear algebra.

- Write the normal equations
- Solve a small normal system
- Interpret the projection onto the column space

## Residuals and fit

Residuals measure the vertical distance between observed and predicted values. The sum of squared residuals is minimised by the least-squares solution, and comparing this sum across candidate models shows which one fits the data better.

This section develops the idea and connects it to the surrounding linear algebra.

- Compute a residual
- Sum squared residuals
- Compare candidate lines by residual size
