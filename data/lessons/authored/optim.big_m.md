# Big-M formulations

A big-M formulation switches a constraint on or off with a binary variable. The constant M must be large enough to make the constraint vacuous when the switch is off and as small as possible for the relaxation.

## The mechanism

The constraint is written with an extra term equal to M times the binary. When the binary is zero the constraint binds, and when it is one the right-hand side is pushed beyond any feasible value.

## The choice of M

A valid M comes from a bound on the constraint activity, which needs a separate analysis. An arbitrarily large M keeps the model correct but makes the relaxation weak and invites numerical trouble.

## Tightening

Per-constraint bounds, indicator constraints and perspective reformulations all reduce the effective M. The reformulation is often the difference between a model that solves and one that does not.

To check your understanding, add a binary switch, derive an activity bound and use per-row bounds.

Big-M is a modelling convenience whose entire quality lies in the tightness of the constant.
