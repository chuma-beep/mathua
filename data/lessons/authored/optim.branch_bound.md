# Branch and bound

Branch and bound explores a tree of subproblems. Each node has a relaxation bound, and a node is discarded when its bound cannot beat the best integer solution found so far. Branching splits a fractional variable into two disjoint ranges, so the tree eventually covers the integer feasible set.

## Taking a branch

Pick a variable with a fractional value in the node relaxation and create two children, one with that variable at most the floor and one with it at least the ceiling. The two children exclude the fractional point and together cover every integer possibility. Re-solving the relaxation in each child gives new bounds.

## Computing a bound

The node bound is the value of its relaxation. For a minimization, no integer solution in that subtree can be better than the bound. The global upper bound is the best integer solution found. A node whose bound is not below the global upper bound can be pruned.

## Pruning a node

A node is pruned when its relaxation is infeasible, when its bound cannot improve the incumbent, or when its relaxation is integral and improves the incumbent. Pruning by bound is the step that makes the method practical, and the quality of the bound decides how much of the tree is visited.

Branch and bound turns a relaxation into an exact algorithm. Its correctness is exhaustive division, and its speed is the quality of its bounds and the order in which nodes are explored.
