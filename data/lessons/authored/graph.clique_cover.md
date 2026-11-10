# Clique covers

## Cliques

A clique is a set of vertices in which every pair is adjacent. A clique cover partitions the vertices into cliques. The minimum number of cliques is the clique cover number.

A single clique covers all vertices when the graph is complete. An empty graph needs one clique per vertex.

## Covers

A clique cover is a partition, so every vertex belongs to exactly one clique. A cover with fewer cliques is better. The minimum is the clique cover number.

Check a proposed cover by verifying that each part is a clique and that the parts are disjoint and cover every vertex.

## Complement colouring

A clique in a graph is an independent set in the complement. A clique cover of the graph is therefore a colouring of the complement. The clique cover number equals the chromatic number of the complement.

This duality turns a clique problem into a colouring problem. Use the smaller or simpler graph to compute the answer.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
