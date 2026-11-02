# Shortest Path Algorithms

## Weighted Paths

A weighted graph assigns a cost to each edge, and the length of a path is the sum of its edge weights. A shortest path minimizes that sum among all paths between two vertices. When all weights are one, the shortest path is simply the path with the fewest edges.

Shortest paths are not always unique, so a problem may ask for the distance or for the number of routes that achieve it. Distinguishing the two questions prevents counting a single path several times.

## Dijkstra's Algorithm

Dijkstra's algorithm grows a set of vertices whose shortest distance is settled. At each step it selects the unsettled vertex with the smallest tentative distance and relaxes its outgoing edges. The algorithm requires nonnegative weights.

On a grid with unit steps, the distance from one corner to another is the sum of the horizontal and vertical separations. On a path graph with listed weights, the distance is the sum of all weights between the endpoints.

## Counting Shortest Routes

Monotone lattice paths from (0,0) to (a,b) use only right and up steps, so every route has a+b steps with a of one kind. The number of routes is the binomial coefficient C(a+b,a). This count is a shortest-path count because every monotone route has the same length.

Counting routes is different from finding a distance: the distance is a+b, while the number of routes grows combinatorially. Both questions are common in applications such as robot motion and network routing.
