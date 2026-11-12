# Linear programs as cone programs

A linear program is already conic: its variables live in a cone whose only member is a ray. The reformulation matters because one cone solver can then read both models.

## A linear program in conic form

Each free variable becomes the difference of two nonnegative entries, and each slack becomes a nonnegative entry of its own. The whole vector sits in the nonnegative orthant, which is a cone.

## The scalar cone

A one-dimensional second-order cone holds every scalar t with t at least the norm of an empty vector, so every nonnegative number qualifies. This is why a linear program is a special second-order cone program.

## What the reformulation buys

The conic form exposes which constraints are conic and which are linear. A solver can then share one Newton system across both kinds and reuse the same infeasibility certificate.

To check your understanding, split each free variable into two nonnegative parts, compare the scalar with the norm of the empty tail and count original constraints.

Treat the nonnegative orthant as the base cone. A linear model is then a cone program with no nonlinear cone block, which is the cheapest case to solve.
