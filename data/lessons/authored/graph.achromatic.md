# The achromatic number

The achromatic number is the largest number of colours in a proper colouring where every pair of colours appears on an edge.

## Complete colourings

A complete colouring uses every pair of colour classes on at least one edge. The achromatic number \(\psi(G)\) is the maximum number of colours in a complete colouring.

## The edge requirement

With \(k\) colour classes there are \(\binom{k}{2}\) pairs, and each pair needs a distinct edge. A complete colouring with \(k\) colours therefore needs at least \(\binom{k}{2}\) edges.

## Values

The complete graph \(K_n\) has \(\psi=n\). A path on \(n\) vertices has \(\psi=\lfloor (n+1)/2\rfloor\), and a star has \(\psi=2\). A cycle on \(n\) vertices has \(\psi=\lfloor n/2\rfloor\) for \(n\ge 4\).

To check your understanding, colour properly, count the colour pairs and colour a complete graph.

Complete colourings balance two demands, many classes and an edge for every pair of them.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
