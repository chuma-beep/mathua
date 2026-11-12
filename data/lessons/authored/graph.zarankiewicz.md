# The Zarankiewicz problem

The Zarankiewicz problem asks for the most edges in a bipartite graph with no complete bipartite subgraph of fixed size. The \(K_{2,2}\) case forbids a four-cycle.

## The forbidden subgraph

A \(K_{s,t}\)-free bipartite graph may contain no \(s\) vertices on one side and \(t\) on the other that are all adjacent. For \(s=t=2\) the forbidden pattern is a cycle of length four.

## The Reiman bound

The number of edges in a \(K_{2,2}\)-free graph on \(n\) vertices is at most \(\frac{n}{4}\left(1+\sqrt{4n-3}\right)\). The bound comes from counting length-two paths and applying convexity.

## Counting four-cycles

The complete bipartite graph \(K_{n,n}\) contains \(\binom{n}{2}^2\) four-cycles, one for each pair of vertices on each side. Removing edges destroys cycles, and the extremal count sits well below the complete graph.

To check your understanding, draw the forbidden \(K_{s,t}\), count length-two paths and choose two vertices on each side.

The Zarankiewicz problem measures how many edges a bipartite graph can carry before a fixed complete bipartite pattern appears.
