# Complete coloring

A complete colouring is proper and every pair of colour classes is joined by an edge.

## The definition

A proper colouring is complete when each pair of colour classes has at least one edge between them. Complete colourings are the colourings counted by the achromatic number.

## Edge requirement

With \(k\) classes the colouring needs at least \(\binom{k}{2}\) edges, one for each pair. A graph with \(m\) edges therefore satisfies \(\binom{\psi(G)}{2}\le m\).

## Refinement

Every colouring can be refined to a complete one by splitting classes that have no edge between them. Refinement never decreases the number of colours, so a maximum complete colouring is also a maximal proper colouring.

To check your understanding, colour properly, count the class pairs and split a disconnected pair.

Completeness is a condition on pairs of classes, and it caps the number of colours by the edge count.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
