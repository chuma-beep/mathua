# Blocks and Biconnected Components

## Cut Vertices and Blocks

A cut vertex is a vertex whose removal disconnects the graph. Blocks are the maximal connected subgraphs that contain no cut vertex of the subgraph itself. Every edge belongs to exactly one block, and blocks meet only at cut vertices.

The block structure turns a graph into a tree of pieces, which is useful because many problems can be solved inside each block and then combined. Bridges are blocks with two vertices, while a cycle is a single block with no cut vertex.

## Block Structure of a Tree

Every edge of a tree is a bridge, so each edge is its own block. A tree on n vertices therefore has exactly n-1 blocks. Removing any internal vertex disconnects the tree, so every internal vertex is a cut vertex.

A path on n vertices has n-2 cut vertices, because the two endpoints are the only vertices whose removal leaves a connected path. This gives a clean way to test the definitions on a familiar family.

## Biconnected Components

A graph is biconnected when it remains connected after any single vertex is removed. A cycle is the simplest biconnected graph, and a complete graph is biconnected for three or more vertices. A single edge is not biconnected in the strict sense but is still a block.

Finding blocks is a standard depth-first search task: low-link values identify the edges that cannot be bypassed. The resulting decomposition supports network reliability questions and separates independent cycles from connecting bridges.
