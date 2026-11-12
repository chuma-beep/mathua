# The Hadwiger number

The Hadwiger number is the order of the largest complete graph that appears as a minor.

## The definition

The Hadwiger number \(h(G)\) is the largest \(t\) such that \(K_t\) is a minor of \(G\). Contracting edges and deleting vertices and edges produces the minor.

## Relation to colouring

Hadwiger's conjecture states that \(\chi(G)\le h(G)\). The conjecture holds for \(t\le 6\), and it would imply the four colour theorem at \(t=5\).

## Values

The complete graph \(K_n\) has \(h=n\). A path has \(h=2\), and a cycle has \(h=3\) because contracting a path leaves a triangle. A tree with at least one edge also has \(h=2\).

To check your understanding, contract an edge, state the conjecture and compute a complete graph.

Minors can only make a graph denser, and the largest complete minor measures how dense it is.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
