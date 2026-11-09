# All-pairs shortest paths

Floyd-Warshall computes the shortest distance between every pair of vertices. It adds one intermediate vertex at a time.

## The distance matrix

Initialize the diagonal to 0 and every direct edge to its weight. For each intermediate vertex k, update every pair i and j through k. After all k, the matrix holds all-pairs shortest distances.

## Diameter

The diameter is the largest finite shortest-path distance between distinct vertices. It is read from the distance matrix after the algorithm finishes. An unreachable pair is ignored because its distance is infinite.

## Finite pairs

A pair is reachable when its entry in the distance matrix is finite. Counting the finite off-diagonal entries gives the number of ordered reachable pairs. The diagonal is excluded because it is always zero.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
