# Finite-difference methods for boundary-value problems

Finite differences replace derivatives by differences on a grid. The boundary-value problem becomes a linear system.

## Build the grid equations

Divide the interval into N subintervals of width h. At each interior point replace y double prime by the centred second difference. The result is one linear equation per interior point.

## Insert the boundary values

The first and last grid values are known from the boundary conditions. Move those known terms to the right-hand side. The interior unknowns remain in a tridiagonal system.

## Solve the linear system

The tridiagonal system is solved by elimination or a specialised algorithm. The solution values approximate y at the grid points. The accuracy improves as h is reduced.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
