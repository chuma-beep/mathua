# The quadratic assignment problem

The quadratic assignment problem places facilities at locations so that the flow between facilities times the distance between locations is minimised. The objective couples every pair of assignments.

## The model

A permutation matrix assigns facilities to locations. The objective is a double sum over facility pairs and location pairs, and it is quadratic in the permutation entries.

## Linearisation

Introducing a variable for each facility pair and location pair turns the objective linear. The linearisation is exact but the model grows with the fourth power of the instance size.

## Why it is hard

The problem contains the travelling salesman and graph partitioning as special cases. Even small instances resist exact solution, so heuristics and semidefinite relaxations carry the practical load.

To check your understanding, assign by a permutation, create pair variables and contain the salesman.

Quadratic assignment is the archetype of a compact model with a hard structure, and every exact approach pays for the linearisation.
