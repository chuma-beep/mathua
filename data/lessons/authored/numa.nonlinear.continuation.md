# Continuation methods

Continuation tracks a solution as a parameter varies, using each solution as the starting point for the next. It turns a hard global problem into a sequence of easy local ones.

## Prediction and correction

A tangent predictor estimates the next point and Newton correction restores the constraint. Adaptive step control shrinks the step near folds.

## Parameter stepping

With step h over a parameter interval of length L, at least L/h steps are needed. Each step costs one or two linear solves.

## Folds

At a turning point the solution branch bends back, so the parameter cannot serve as the continuation variable. Arc-length continuation parametrises by distance along the curve instead.

To check your understanding, take a tangent predictor, count the steps and recognise a turning point.

Continuation starts from a solution that is known and walks to one that is not.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
