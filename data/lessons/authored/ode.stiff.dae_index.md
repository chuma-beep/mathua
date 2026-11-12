# Differential-algebraic equations

A differential-algebraic equation carries algebraic constraints beside differential equations. The index counts how many differentiations turn the constraint into an explicit differential equation, and it controls which integrators apply.

## Semi-explicit form

The differential variables have their own equations and the algebraic variables appear only through a constraint. The constraint must hold at every time, so the initial values cannot be chosen freely.

## The index

Index one means the constraint can be solved for the algebraic variable. Index two means the constraint must be differentiated once before that is possible, and each extra differentiation raises the index by one.

## Consequences

Consistent initial values must satisfy the constraint and its derivatives. BDF and implicit Runge-Kutta methods handle index one directly, while higher index problems need an index reduction or a specialised formulation.

To check your understanding, separate differential and algebraic variables, try to solve for the algebraic variable and state the consistency requirement.

The index is the number that decides whether a general-purpose stiff solver can be applied without rewriting the model.
