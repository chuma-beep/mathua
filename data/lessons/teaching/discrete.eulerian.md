# Eulerian Trails, Hamiltonian Cycles and Degree Tests

**Eulerian circuit:** A closed trail using every edge exactly once. A connected graph has one iff every vertex has even degree; it has an (open) Eulerian trail iff zero or two vertices have odd degree, and then the trail runs between the two odds. **Hamiltonian cycle** is a different demand: one cycle visiting every vertex exactly once, with no degree test.

## Worked: $K_3$ is Eulerian

Take the triangle with vertices $A,B,C$:
1. Each vertex has degree 2, all even, and the graph is connected.
2. Walk $A\to B\to C\to A$: it uses each of the 3 edges exactly once and returns to $A$.
3. That closed trail is the Eulerian circuit the even-degree test promises.

So a connected graph with all degrees even always closes up: entering a vertex on an unused edge leaves an unused edge to exit by until every edge is spent.

## Worked: a trail between two odds

Take the path $A-B-C$ with degrees 1, 2, 1:
1. Exactly two vertices, $A$ and $C$, have odd degree; $B$ is even.
2. Start at $A$: the only unused edge forces $A\to B$, then $B\to C$, using each edge once.
3. The trail ends at $C$, the other odd vertex, and cannot close since $A$ has no unused edge left.

So with exactly two odds the trail must start at one odd vertex and die at the other; with four or more odds no single trail covers every edge.

## Eulerian is not Hamiltonian

The Petersen graph is connected with all vertices of degree 3, yet it has no Hamiltonian cycle: no degree parity test can certify a vertex-covering cycle. Meanwhile $K_{3,3}$ is bipartite, so $\chi(K_{3,3})=2$, and every tree is bipartite too. In general, Eulerian questions are about edges and parity, while Hamiltonian and coloring questions are about vertices and need separate arguments.
