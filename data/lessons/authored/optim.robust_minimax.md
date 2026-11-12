# Robust minimax design

Robust design replaces a point estimate of the data with a set of plausible values. The design is chosen to minimise the worst loss over that set, which turns a single optimisation into a minimax problem over the uncertainty.

## The uncertainty set

The set encodes what the modeller is willing to assume: a box, a ball or a polytope. A larger set buys a stronger guarantee and a worse average performance, and that trade is the main design decision.

## Tractability

A linear loss over a box or a ball keeps the worst case computable in closed form, and the resulting problem stays linear or conic. An arbitrary nonlinear loss over a general set usually makes the inner maximum intractable.

## Conservatism

The worst case is often attained at a corner the data will never produce. Tempering the set, or penalising the deviation from the nominal solution, is how practitioners avoid a design that is safe and useless.

To check your understanding, choose box, ball or polytope, keep the inner maximum closed form and locate the worst-case corner.

A robust minimax design is only as meaningful as its uncertainty set, so the set deserves more scrutiny than the solver.
