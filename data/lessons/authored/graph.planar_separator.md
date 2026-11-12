# Planar separators

Every planar graph has a small set of vertices whose removal splits it into two parts of comparable size.

## The separator theorem

A planar graph on \(n\) vertices has a separator of at most \(2\sqrt{2n}\) vertices whose removal leaves components of at most \(2n/3\) vertices each. The bound is the Lipton-Tarjan theorem.

## Consequences

The separator theorem gives divide-and-conquer algorithms for planar graphs with \(O(\sqrt n)\) overhead per level. It also implies that planar graphs have treewidth \(O(\sqrt n)\).

## Grids

The square grid has a separator of size about \(\sqrt n\), so the square-root order is tight. A long path has a separator of size one, which shows the bound is not always attained.

To check your understanding, state the separator size, divide and recurse and separate a grid.

Planarity limits how much a graph can resist being cut in half.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
