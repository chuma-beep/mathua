# Fill-reducing orderings

A fill-reducing ordering permutes rows and columns before factorisation so elimination creates fewer new nonzeros. The permutation changes the cost and never the solution.

## Local heuristics

Minimum degree eliminates the vertex of smallest degree first. Reverse Cuthill-McKee instead reduces bandwidth, which helps banded solvers.

## Nested dissection

Recursive separator ordering suits meshes from finite elements. It gives asymptotically optimal fill for two- and three-dimensional grids.

## Trade-offs

Finding the optimal order is NP-hard, so every practical method is a heuristic. The symbolic factorisation measures the result before the numeric phase begins.

To check your understanding, pick the minimum degree, find a separator and recognise the hard problem.

Every sparse direct solver starts by choosing an order.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
