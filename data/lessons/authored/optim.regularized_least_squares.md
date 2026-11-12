# Regularised least squares

A regulariser adds a penalty on the coefficients to a least-squares objective. The penalty trades fit against size.

## The penalty

A quadratic penalty adds a positive multiple of the squared norm of the coefficients. The objective stays smooth and the solution is a linear system with a shifted diagonal.

## The effect

The penalty shrinks the coefficients toward zero. A larger penalty gives a smaller norm and a larger residual on the training data.

## Choosing the penalty

The penalty is chosen by cross-validation or by a validation split. Too small overfits and too large underfits.

To check your understanding, write the penalised objective, compare two penalties and try a grid.

A regulariser is a prior written as a penalty, and the penalty strength is the decision.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
