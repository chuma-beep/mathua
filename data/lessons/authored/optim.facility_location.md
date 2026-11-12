# Facility location

Facility location chooses which facilities to open and which open facility serves each customer. Opening costs are fixed, assignment costs are linear, and every customer must be served.

## The formulation

A binary variable opens a facility and a continuous variable assigns demand. The linking constraint forbids assignment to a closed facility, which is the only place the two variable types interact.

## The strength of the relaxation

The basic relaxation is weak because a customer can split demand across many facilities. Adding a constraint that each customer is assigned at least once, even though the demand is already covered, tightens the bound considerably.

## The uncapacitated case

Without capacity the problem is still NP-hard, but the LP relaxation often comes close to the integer optimum. That is why it is a standard benchmark for cut strength.

To check your understanding, open facilities with binaries, allow split demand and keep the problem NP-hard.

Facility location is a fixed-charge model whose relaxation is weak until the assignment constraints are reinforced.
