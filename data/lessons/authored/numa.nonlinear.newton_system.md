# Newton's method for systems

Newton's method for a system solves a linear model at each iterate and updates by its solution. The Jacobian replaces the scalar derivative.

## The step

Solve \(J(x_k)s_k=-F(x_k)\), then set \(x_{k+1}=x_k+s_k\). The Jacobian is \(n\times n\) and changes at every step.

## Local convergence

Near a simple root the error is squared at each step. Away from a root the step can be large or the Jacobian can be singular.

## Cost

Each step needs a Jacobian evaluation and a linear solve costing \(O(n^{3})\) densely. Quasi-Newton methods trade that cost for slower convergence.

To check your understanding, evaluate the Jacobian, square the error and count the Jacobian entries.

The Jacobian is the whole difficulty of Newton's method in more than one variable.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
