# Backward differentiation formulas

Backward differentiation formulas fit a polynomial through past values and impose the derivative at the new time. The family is implicit and includes methods that remain stable on stiff problems.

## Construction

A polynomial is interpolated through the new value and several past values. Its derivative at the new time is set equal to the right-hand side, which produces one implicit equation per step.

## First and second order

The first-order formula is backward Euler. The second-order formula uses the last three values and is still A-stable, which is why it is the workhorse of stiff solvers.

## Stability limits

Higher-order formulas lose A-stability, and the family is zero-stable only up to order six. Beyond that order the method is unstable in principle, so a practical solver stops at five.

To check your understanding, interpolate the past values, identify BDF1 and state the order bound.

BDF trades the large stability region of one-step methods for higher order at the cost of a history of past values.
