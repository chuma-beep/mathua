# Set cover

Set cover chooses the fewest sets whose union contains every element. It is an integer program with a binary variable per set and one covering constraint per element. Vertex cover is the special case where each set has two elements, namely the endpoints of an edge.

## Testing a cover

For each element, check whether at least one chosen set contains it. A single uncovered element disqualifies the selection. The test is a union computation, and it is the feasibility check for both the integer model and any heuristic.

## Finding a minimum cover

For a small instance, enumerate subsets of sets in increasing size and stop at the first cover. That is exact and finite. For larger instances, the LP relaxation gives a lower bound and a greedy choice by coverage gives an upper bound, and the two together bound the optimum.

## Comparing with a vertex cover

A vertex cover is a set cover whose elements are edges and whose sets are the stars at each vertex. The reduction shows why set cover is at least as hard as vertex cover. A cover of a graph by vertices is therefore a special case, and a general set cover instance need not come from a graph.

Set cover is the general covering model behind many selection problems. Its LP relaxation and its greedy bound are the two facts most often used in practice.
