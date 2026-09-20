# Paths, Cycles, and Connectivity

A graph is more than a pile of edges: what matters is how vertices link
together. This lesson defines walks, paths, and cycles, then uses them to
say when a graph hangs together and how far apart two vertices are.

## Walks, Trails, and Paths

A **walk** is a sequence of vertices \(v_0, v_1, \dots, v_k\) where each
consecutive pair is joined by an edge. Vertices and edges may repeat. The
**length** of a walk is its number of edges, \(k\).

- A **trail** is a walk with no repeated edge.
- A **path** is a walk with no repeated vertex (except that a closed walk
  returning to its start is a **circuit**, and a circuit with no other
  repetition is a **cycle**).

For example, in the graph with edges \(\{A-B, B-C, C-D\}\), the sequence
\(A-B-C-D\) is a path of length 3. The sequence \(A-B-A\) is a walk but
not a path, because \(A\) repeats.

## Connected Graphs and Components

A graph is **connected** if there is a path between every pair of vertices,
and **disconnected** otherwise. A **connected component** (or just
component) is a maximal connected piece: a vertex with no edges at all is
its own component, an **isolated vertex**.

For example, the graph with edges \(\{A-B, B-C, C-D\}\) is connected: every
vertex reaches every other along the chain. The graph with edges
\(\{A-B, C-D\}\) is disconnected, with two components \(\{A, B\}\) and
\(\{C, D\}\). There is no path from \(A\) to \(D\) in it.

## Distance and Shortest Paths

The **distance** \(d(u, v)\) between two vertices is the length of the
shortest path joining them. If no path exists, the distance is undefined:
we write **none**.

For example, with edges \(\{A-B, B-C, C-D\}\), the distance from \(A\) to
\(D\) is 3, via the unique shortest path \(A-B-C-D\). With edges
\(\{A-B, C-D\}\), there is no path from \(A\) to \(D\), so the answer is
none.
