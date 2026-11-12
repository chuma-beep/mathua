# Chain complexes

A chain complex is a sequence of abelian groups with boundary maps whose composition is zero. The failure of exactness at each group is measured by homology.

## The boundary condition

The maps satisfy \(d_{k} \circ d_{k+1} = 0\), so every boundary is a cycle. A complex with groups in degrees zero to \(n\) has \(n\) boundary maps.

## Euler characteristic

The alternating sum of the ranks of the chain groups equals the alternating sum of the Betti numbers. This is the Euler-Poincaré identity for a finite complex.

## Cycles and boundaries

The \(k\)-th homology group is the quotient of \(k\)-cycles by \(k\)-boundaries. Its rank is the \(k\)-th Betti number.

To check your understanding, count the boundary maps, alternate the ranks and quotient cycles by boundaries.

A chain complex packages the data that homology measures.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller values and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
