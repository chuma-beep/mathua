# Hamiltonian Cycles

## Visiting Every Vertex Once

A Hamiltonian cycle is a closed walk that visits every vertex of a graph exactly once and returns to its starting point. The definition is about vertices rather than edges, which is what separates it from an Eulerian circuit. A graph that has such a cycle is called Hamiltonian.

Deciding whether a graph is Hamiltonian is hard in general, so courses rely on sufficient conditions and on small families whose structure is known. Complete graphs are the simplest example, because every possible edge is present. Cycles themselves are Hamiltonian by definition.

## Dirac's Sufficient Condition

Dirac's theorem says that a simple graph on n at least three vertices is Hamiltonian when every vertex has degree at least n/2. The hypothesis is sufficient but not necessary, so a graph can be Hamiltonian with smaller degrees. The theorem is most useful when the graph is dense.

To apply the condition, compute the minimum degree and compare it with half the vertex count. If the bound holds, the graph is Hamiltonian without any further search. If the bound fails, the theorem says nothing and another argument is needed.

## Counting Cycles

The complete graph K_n has (n-1)!/2 distinct Hamiltonian cycles. Fixing one starting vertex removes the rotational choices, and reversing the order removes the two directions. The formula grows quickly, which is why listing cycles by hand is practical only for very small n.

For a cycle graph C_n there is exactly one Hamiltonian cycle, namely the cycle itself. Counting Hamiltonian cycles therefore depends strongly on the graph family, not only on the number of vertices.
