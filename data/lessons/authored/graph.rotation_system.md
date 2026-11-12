# Rotation systems

A rotation system records the cyclic order of edges around each vertex, and that order determines the faces of an embedding.

## Darts

Each edge has two darts, one in each direction. A graph with \(m\) edges has \(2m\) darts, and a rotation system orders the darts out of each vertex cyclically.

## Faces from the rotation

Following the rule that the next dart is the successor of the reverse dart traces the faces. For a connected planar embedding, Euler's formula gives \(f=m-n+2\) faces.

## Genus

The rotation system determines the embedding surface through \(n-m+f=2-2g\) for an orientable surface of genus \(g\). A planar rotation system has genus zero, and every other rotation system reveals a handle.

To check your understanding, double the edge count, follow the successor rule and count vertices, edges and faces.

A cyclic order at each vertex is enough to reconstruct the faces and the surface.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
