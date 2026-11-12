# Rosenbrock methods

A Rosenbrock method linearises the implicit equation around the current state, so each step needs one Jacobian and a sequence of linear solves. It keeps the stability of an implicit method without a nonlinear iteration.

## Linearisation

The right-hand side at the new time is replaced by its first-order Taylor expansion about the current state. The unknown new value then appears linearly and can be isolated by a matrix solve.

## Stage equations

Each stage solves a linear system with the same matrix, shifted by a stage coefficient. The right-hand sides differ, so one factorisation serves every stage of the step.

## Cost and use

The method needs one Jacobian evaluation and one factorisation per step rather than an iteration to convergence. That makes it attractive for stiff problems of moderate dimension, where a dense factorisation is affordable.

To check your understanding, expand the right-hand side, form the stage matrix and count the Jacobians.

Rosenbrock methods are the linear-implicit compromise between the robustness of a fully implicit method and the cost of an explicit one.
