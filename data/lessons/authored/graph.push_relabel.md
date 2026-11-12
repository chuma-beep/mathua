# The push-relabel algorithm

Push-relabel works with a preflow that may violate conservation, and repairs it with local push and relabel operations.

## Preflow and labels

The algorithm saturates every edge out of the source, then maintains a height label on each vertex. Flow may enter a vertex faster than it leaves, so conservation is temporarily violated.

## Push and relabel

A push sends flow along an edge to a vertex one level lower. When no such edge exists, a relabel raises the vertex above its lowest neighbour, and the excess eventually returns to the source or reaches the sink.

## The complexity

Each vertex is relabelled at most \(2V-1\) times, and the naive implementation runs in \(O(V^2E)\). The highest-label selection rule and the gap heuristic make it fast in practice.

To check your understanding, saturate the source edges, push downhill and bound the relabels.

Relax conservation into a preflow, then repair it with pushes that follow a height function.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
