# Eulerian Trails, Hamiltonian and Coloring

**Eulerian circuit:** A closed trail using every edge exactly once. **Hamiltonian cycle:** A cycle visiting every vertex exactly once. **Chromatic number $\chi(G)$:** Fewest colors so adjacent vertices differ.

## Traversing Graphs

### Eulerian Criteria
Connected graph has Eulerian circuit iff every vertex has even degree; has trail iff zero or two vertices have odd degree (trail between the odds).

### Coloring and Hamiltonians
Bipartite graphs $\chi=2$ (e.g. $K_{3,3}$, trees). Petersen graph has Eulerian trails but is famously non-Hamiltonian.

## Example

$K_3$ (triangle): all degrees $2$ (even) → Eulerian circuit around the triangle; it is also its own Hamiltonian cycle; $\chi(K_3)=3$.
