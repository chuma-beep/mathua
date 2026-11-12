# Dirac's theorem

Dirac's theorem turns a minimum-degree condition into a Hamiltonian cycle. Half the vertices is enough, and the bound is tight.

## The threshold

If a graph on \(n\ge 3\) vertices has minimum degree at least \(n/2\), it has a Hamiltonian cycle. The ceiling form \(\lceil n/2\rceil\) is the useful integer statement.

## The edge count

Minimum degree \(\lceil n/2\rceil\) forces at least \(n\lceil n/2\rceil/2\) edges. For even \(n\) that is \(n^2/4\), the same order as the Turán bound for triangles.

## Tightness

Two disjoint cliques of size \(n/2\) have minimum degree \(n/2-1\) and no Hamiltonian cycle. Dirac's bound cannot be lowered, and a complete graph meets it with room to spare.

To check your understanding, state the minimum degree, multiply degree by order and build the two-clique example.

A dense enough graph cannot avoid a cycle through every vertex, and two cliques show the density is necessary.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
