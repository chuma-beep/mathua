# Levenberg-Marquardt

Levenberg-Marquardt interpolates between Gauss-Newton and gradient descent by adding a damping term to the normal equations. It is a standard solver for nonlinear least squares.

## The system

Solve \((J^{T}J+\lambda I)s=-J^{T}r\). As \(\lambda\to 0\) the step approaches Gauss-Newton; as \(\lambda\to\infty\) it approaches a short gradient step.

## Damping update

A successful step reduces lambda, usually by a factor of ten; a failed step increases it by the same factor. The damping adapts to the local linearity.

## Cost

Each trial step needs a linear solve with the same matrix plus a diagonal shift. Factorising the normal matrix once per damping value is the dominant cost.

To check your understanding, add the damping term, reduce after success and factor once per lambda.

The damping parameter is a dial between speed and safety.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
