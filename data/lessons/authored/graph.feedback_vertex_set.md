# Minimum feedback vertex sets

## Breaking every cycle

A feedback vertex set is a set of vertices whose removal leaves a forest. A minimum feedback vertex set has the smallest possible size. Every cycle must lose at least one vertex. In a disjoint union of cycles, one vertex from each cycle is necessary and sufficient. Three disjoint cycles therefore need three vertices.

## Testing a proposed set

Remove every vertex in the proposed set and inspect the remaining graph. A cycle that contains none of the removed vertices still remains. Count those cycles rather than guessing from the size of the set. A large set can still miss a cycle, while a small set can hit every cycle. The edge list is the reliable source for this test.

## Counting the remaining forest

A cycle of length L becomes a path with L - 2 edges after one vertex is removed. Summing L - 2 over all cycles gives the number of edges in the resulting forest. This count is a useful check that the removal really broke every cycle. A forest on n vertices has at most n - 1 edges, which is another check.
