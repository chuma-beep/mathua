# Defective coloring

Defective colouring allows a bounded number of same-colour neighbours instead of forbidding them.

## The definition

A \((k,d)\)-colouring uses \(k\) colours and gives every vertex at most \(d\) neighbours of its own colour. Ordinary colouring is the case \(d=0\).

## The degree bound

A graph of maximum degree \(\Delta\) always has a \((k,d)\)-colouring when \(k(d+1)\ge\Delta+1\). Greedily colouring in decreasing degree order keeps at most \(d\) same-colour neighbours.

## Edge counts

In a \((k,d)\)-colouring, each colour class induces a graph of maximum degree \(d\), so it has at most \(nd/2\) edges. Summing over classes bounds the total edges that can stay within colour classes.

To check your understanding, allow \(d\) same-colour neighbours, order by degree and bound one colour class.

Tolerating a few monochromatic edges buys a large reduction in the number of colours.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
