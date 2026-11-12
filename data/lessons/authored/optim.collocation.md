# Direct collocation

Direct collocation discretises both the state and the control, and imposes the dynamics as equality constraints at collocation points. The optimiser then handles states and controls together under a sparse constraint set.

## The transcription

The horizon is split into intervals, and within each interval the state is approximated by a polynomial. The dynamics are enforced at chosen nodes, which become equality constraints linking consecutive intervals.

## Sparsity

Each constraint involves only the variables of one or two intervals. The Jacobian is therefore banded, and a sparse linear solver keeps the cost close to linear in the number of intervals.

## Accuracy against robustness

More nodes improve the dynamics approximation but enlarge the problem. Collocation tolerates unstable dynamics far better than shooting because the state is a decision variable rather than a simulation output.

To check your understanding, split the horizon, couple neighbouring intervals and add nodes for accuracy.

Collocation trades a larger decision space for a sparse structure and a dynamics error that does not compound along the horizon.
