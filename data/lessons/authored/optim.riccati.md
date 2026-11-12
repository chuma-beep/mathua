# The Riccati equation

The Riccati equation is the matrix equation whose solution gives the regulator gain. For a finite horizon it is integrated backward in time, and for an infinite horizon it becomes an algebraic equation.

## The differential form

The matrix evolves backward from the terminal weight under a quadratic differential equation. Each step adds the state weight and subtracts a term that rewards control authority, which keeps the matrix positive semidefinite.

## The algebraic form

At steady state the derivative vanishes and the equation becomes algebraic. Its stabilising solution is the one that makes the closed loop stable, and it is unique under detectability and stabilisability.

## Solving it

The algebraic equation is solved by iteration, by a Schur decomposition of the associated Hamiltonian matrix, or by a structured eigensolver. The Hamiltonian approach is preferred when the dimensions are moderate because it is numerically stable.

To check your understanding, integrate backward, set the derivative to zero and iterate to convergence.

The Riccati equation is where the regulator's gain comes from, and its stabilising solution is the one an implementation must select.
