# Surface triangulations

A triangulation of a surface counts vertices, edges and faces. The Euler characteristic constrains the counts, and the edge-face incidence gives a second linear relation.

## Euler characteristic

For a triangulated surface, \(\chi = V - E + F\). The count is independent of the triangulation, which is what makes it an invariant.

## Counting edges

Every face has three edges and every edge borders two faces, so \(3F = 2E\). Combined with the Euler characteristic, a closed triangulated surface satisfies \(E = 3(V - \chi)\).

## Counting faces

The same relations give \(F = 2(V - \chi)\) for a closed triangulated surface. A triangulated torus with seven vertices has fourteen faces.

To check your understanding, add vertices and faces, use three faces per edge count and double the vertex surplus.

Triangulations turn the Euler characteristic into bookkeeping.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
