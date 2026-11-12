# Acyclic coloring

An acyclic colouring is proper and every two colour classes induce a forest.

## The definition

A proper colouring is acyclic when the union of any two colour classes is a forest. Equivalently, no cycle uses only two colours.

## The bound

Every graph of maximum degree \(\Delta\) has an acyclic colouring with \(O(\Delta^{4/3})\) colours, and the current best bounds are close to that order. Trees and complete graphs sit at the two extremes.

## Basic values

A tree has acyclic chromatic number two, because two colours already give a forest. The complete graph \(K_n\) needs \(n\) colours, since any two classes meet in at most one edge. A cycle \(C_n\) needs three colours.

To check your understanding, colour properly, state the degree bound and colour a tree.

Forbidding bichromatic cycles is much stronger than forbidding monochromatic edges.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
