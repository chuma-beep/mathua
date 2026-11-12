# Bilevel programs

A bilevel program has an outer objective and an inner optimisation whose solution depends on the outer variables. The outer decision is evaluated at the inner optimum, not at any other feasible point.

## The nested structure

For each outer choice the inner problem is solved to optimality and the resulting inner variable is substituted into the outer objective. The outer feasible set may itself depend on the inner solution.

## Reformulations

When the inner problem is convex, its optimality conditions can replace it. That produces a single-level problem with complementarity constraints, which is nonconvex even though the inner problem is convex.

## Why they are hard

The value function of the inner problem is generally nonsmooth and the bilevel feasible set is not convex. Gradient methods need a derivative of the inner solution, which is where implicit differentiation enters.

To check your understanding, solve the inner problem, replace the inner problem by its conditions and face a nonsmooth value function.

A bilevel program is a nested decision, and the single-level reformulation trades an inner solver for complementarity constraints.
