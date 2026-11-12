# The Grundy number

The Grundy number is the most colours a greedy colouring can be forced to use.

## Greedy colouring

A greedy colouring visits the vertices in a fixed order and gives each vertex the smallest colour not used by an earlier neighbour. Different orders can produce different numbers of colours.

## The definition

The Grundy number \(\Gamma(G)\) is the maximum number of colours over all vertex orders. Every colour class is nonempty, and colour \(i\) is used only when a neighbour already carries every colour below \(i\).

## Values

The complete graph \(K_n\) has \(\Gamma=n\) for every order. A path on at least two vertices has \(\Gamma=2\), and a complete bipartite graph has \(\Gamma=2\). Deleting one edge from \(K_n\) lowers the Grundy number to \(n-1\).

To check your understanding, order the vertices, maximise over orders and colour a complete graph.

The Grundy number measures how badly a greedy colouring can behave, not how many colours the graph needs.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
