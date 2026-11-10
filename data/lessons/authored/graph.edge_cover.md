# Edge covers

## Covering every vertex

An edge cover is a set of edges that touches every vertex of a graph. The minimum edge cover is the smallest such set. A graph with an isolated vertex has no edge cover, so the concept applies to graphs in which every vertex has at least one incident edge. For a star with one centre and k leaves, every leaf needs its own edge, so the minimum has k edges.

## Testing a proposed cover

Mark the endpoints of every edge in the proposed set. Any unmarked vertex is uncovered. The number of uncovered vertices is a direct test, and zero means the proposed set is an edge cover. This test is independent of whether the set is minimum. A proposed set can be an edge cover and still be larger than necessary.

## Edge covers and matchings

For a graph with no isolated vertices, the minimum edge cover size plus the maximum matching size equals the number of vertices. If a graph has 7 vertices and a maximum matching of size 3, the minimum edge cover has size 4. The relation converts one problem into the other and gives a fast check.
