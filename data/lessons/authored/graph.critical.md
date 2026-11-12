# Color-critical graphs

A graph is \(k\)-critical when its chromatic number is \(k\) but every proper subgraph is \((k-1)\)-colorable. Critical graphs are the minimal obstructions to colouring.

## The minimum degree

A \(k\)-critical graph has minimum degree at least \(k-1\). A vertex of smaller degree could be removed, the rest coloured with \(k-1\) colours, and the vertex recoloured with a spare colour.

## The edge lower bound

Combining the minimum degree with the handshake lemma gives at least \(n(k-1)/2\) edges in a \(k\)-critical graph on \(n\) vertices. Odd cycles are the 3-critical case and meet the bound with exactly \(n\) edges.

## Critical examples

The complete graph \(K_k\) is \(k\)-critical with \(\binom{k}{2}\) edges. An odd cycle is 3-critical. Every \(k\)-critical graph is connected and has no cut vertex, because a separation would let the pieces be coloured independently.

To check your understanding, delete a low-degree vertex, apply the handshake lemma and count the edges of a complete graph.

Minimality with respect to colouring forces a minimum degree, and that degree forces an edge count.
