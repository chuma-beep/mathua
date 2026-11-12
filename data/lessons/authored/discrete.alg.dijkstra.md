# Dijkstra's algorithm

Dijkstra's algorithm grows a set of settled vertices, always settling the unsettled vertex with the smallest tentative distance. Non-negative edge weights are what make the greedy choice safe.

## Relaxation

Settling a vertex relaxes each outgoing edge: if the tentative distance plus the edge weight is smaller, the neighbour's distance is replaced.

## The priority queue

A min-heap returns the closest unsettled vertex in logarithmic time. With a binary heap the whole algorithm costs \(O((V+E)\log V)\).

## Why weights must be non-negative

A negative edge could make a later path shorter than a settled distance, which would break the greedy invariant that settled distances are final.

To check your understanding, settle the closest vertex, extract the minimum and find the failing edge.

Dijkstra is a greedy algorithm with a proof obligation: each settled distance is the true shortest distance because no non-negative edge can improve it later.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
