# Bellman-Ford shortest paths

Bellman-Ford relaxes every edge repeatedly and tolerates negative edge weights. A further improvement after n minus 1 passes reveals a negative cycle.

## Relaxation

Start the source at distance 0 and every other vertex at infinity. For each edge, replace the distance to its head when the path through its tail is shorter. Repeat the full pass n minus 1 times.

## Negative cycles

If a distance still improves on pass n, the graph has a negative cycle reachable from the source. The cycle can be traversed repeatedly to make the path arbitrarily short. Bellman-Ford reports this instead of a finite distance.

## Path length

The predecessor pointer records the last edge used to improve a distance. Following the pointers from the target back to the source gives a shortest path. The number of pointers is the number of edges in that path.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
