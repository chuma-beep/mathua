# The Singleton bound

The Singleton bound says a code of length \(n\), size \(M\) and minimum distance \(d\) over an alphabet of size \(q\) satisfies \(M\le q^{n-d+1}\). Codes meeting it are called MDS.

## The bound

Puncturing the last \(d-1\) coordinates must leave distinct words, otherwise two codewords would differ in at most \(d-1\) places. Hence \(M\le q^{n-d+1}\).

## MDS codes

Reed-Solomon codes meet the bound with equality. Their parameters are optimal for the given length and distance.

## Consequences

The bound limits how much redundancy is needed. It is weaker than the Hamming bound for some parameter ranges and stronger for others.

To check your understanding, puncture the last coordinates, check equality and compare with Hamming.

The Singleton bound is a puncturing argument, not a sphere-packing one. It is the right tool when the alphabet is large and the distance is the design target.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
