# Least-squares objectives

A least-squares problem minimises the squared distance between a model and data. Its normal equations are linear.

## The objective

The residual is the data minus the model, and the objective is the sum of squared residuals. Squaring makes the objective smooth and convex.

## The normal equations

Setting the gradient to zero gives a linear system with the transpose times the matrix on the left. The right-hand side is the transpose times the data.

## Residual at the solution

At the optimum the residual is orthogonal to the columns of the matrix. That orthogonality is the geometric content of the normal equations.

To check your understanding, evaluate a residual, form the left side and compute the residual.

Least squares is the meeting point of linear algebra and optimisation.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
