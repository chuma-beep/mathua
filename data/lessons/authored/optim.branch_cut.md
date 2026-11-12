# Branch and cut

Branch and cut combines branch and bound with cutting planes. Cuts tighten the relaxation at each node, and branching splits the remaining region when the relaxation is still fractional.

## The loop at a node

The node relaxation is solved, then valid inequalities that the fractional solution violates are added and the relaxation is re-solved. The loop stops when no violated cut is found.

## The bound

The relaxed value is a bound for every integer point in the node's region. Pruning happens when that bound is worse than the best known integer solution, so tighter cuts prune sooner.

## Separation

A cut is useful only if it can be found without enumerating the integer points. Separation routines are heuristics for finding a violated inequality, and their strength determines the practical speed of the solver.

To check your understanding, solve the relaxation, bound the node region and find a violated inequality.

Branch and cut is branch and bound with a relaxation that keeps improving at each node, and separation is the engine that makes the improvement cheap.
