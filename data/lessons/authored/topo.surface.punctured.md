# Punctured surfaces

Removing points from a surface changes its homotopy type and its Euler characteristic. A punctured surface deformation retracts onto a graph, so its homology is computable by counting.

## Euler characteristic

Removing \(p\) points from a closed orientable surface of genus \(g\) gives \(\chi = 2 - 2g - p\). Each puncture is an open disc removed, and the Euler characteristic falls by one.

## First Betti number

A surface of genus \(g\) with \(p \ge 1\) punctures has \(b_1 = 2g + p - 1\). The punctures add independent loops, which is why a plane with \(p\) punctures has first Betti number \(p - 1\).

## The plane

The plane has \(\chi = 1\). Removing \(p\) points gives \(\chi = 1 - p\), and the complement retracts onto a wedge of \(p\) circles when \(p \ge 1\).

To check your understanding, subtract one per puncture, add twice the genus and the punctures and start from one.

Punctures trade two-dimensional cells for one-dimensional loops.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
