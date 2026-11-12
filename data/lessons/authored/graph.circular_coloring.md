# Circular coloring

Circular colouring replaces the discrete palette by a circle of circumference \(k\), with adjacent vertices required to be far apart.

## The definition

A \((k,d)\)-colouring maps the vertices to points of a circle of circumference \(k\) so that adjacent vertices are at circular distance at least \(d\). The circular chromatic number \(\chi_c(G)\) is the infimum of \(k/d\) over such colourings.

## Chromatic interpolation

Every graph satisfies \(\chi(G)-1<\chi_c(G)\le\chi(G)\). The circular chromatic number therefore refines the ordinary one, and integer values occur exactly when the two agree.

## Cycles

An even cycle has \(\chi_c=2\). An odd cycle \(C_{2k+1}\) has \(\chi_c=2+1/k\), so \(C_5\) has \(5/2\) and \(C_7\) has \(7/3\). These are the standard examples of non-integer circular chromatic numbers.

To check your understanding, place the vertices on a circle, compare with the chromatic number and handle an even cycle.

A circle with a minimum gap interpolates between the integer chromatic numbers.
