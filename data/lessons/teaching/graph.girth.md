# Girth of a Graph

## The Shortest Cycle

The girth of a graph is the length of its shortest cycle. A tree has no cycles, so its girth is infinite by convention. In a graph with cycles, the girth measures how far the graph is from being a tree.

Every cycle contains a shortest cycle, so a finite graph with at least one cycle has a well-defined girth. The value is always at least three, because simple graphs have no cycles of length one or two.

## Girth of Common Families

A cycle graph C_n has girth n, since its only cycle is the whole graph. A complete graph K_n with n at least three has girth three, because any three vertices form a triangle. A complete bipartite graph K_{m,n} with both parts of size at least two has girth four.

These examples show that adding edges generally lowers the girth until triangles or quadrilaterals appear. Graphs with large girth and many edges are important in extremal graph theory and in the design of expanders.

## Bipartite Graphs

A graph is bipartite exactly when it has no odd cycle. Its girth is therefore even whenever the graph contains a cycle. The smallest possible value in a bipartite graph with a cycle is four.

Checking bipartiteness with a two-coloring also detects odd cycles, so it can certify that the girth is not three. For small graphs, a breadth-first search finds the shortest cycle directly.
