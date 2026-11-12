# Total coloring

Total colouring colours vertices and edges together so that every two incident or adjacent elements differ.

## The definition

A total colouring assigns colours to \(V\cup E\). Two vertices joined by an edge, two edges sharing a vertex and a vertex with an incident edge must all receive different colours.

## The lower bound

At a vertex of maximum degree \(\Delta\), the vertex and its \(\Delta\) incident edges form a clique of size \(\Delta+1\) in the conflict graph. Hence the total chromatic number is at least \(\Delta+1\).

## Values

The complete graph \(K_n\) has total chromatic number \(n\) for odd \(n\) and \(n+1\) for even \(n\ge 4\). A path on at least three vertices needs three colours, and a cycle needs three or four according to its length.

To check your understanding, colour vertices and edges, build the conflict clique and colour a complete graph.

Vertices and edges compete for colours, and the maximum degree sets the floor.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
