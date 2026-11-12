# Branch-and-bound search

Branch and bound searches for an optimum by branching on decisions and bounding the best value any completion of a node can achieve. A node is discarded when its bound cannot beat the incumbent.

## Bounding

Each node carries an optimistic bound: the best value reachable if every remaining choice went perfectly. A relaxation such as ignoring a constraint supplies one.

## Pruning by bound

Keep an incumbent, the best known feasible solution. A node whose bound is no better than the incumbent cannot contain a better solution and is discarded.

## Branching order

Exploring the most promising child first raises the incumbent early, which prunes more of the remaining tree.

To check your understanding, build an optimistic bound, keep the best feasible value and order children by bound.

Branch and bound is exact, unlike a heuristic, because the bound is optimistic. If it is wrong in the unsafe direction the search can discard the true optimum.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
