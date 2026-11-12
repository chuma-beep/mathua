# Ore's theorem

Ore's theorem relaxes Dirac's condition from every vertex to every nonadjacent pair. A degree sum of n is enough.

## The degree-sum condition

If every pair of nonadjacent vertices \(u,v\) satisfies \(d(u)+d(v)\ge n\), the graph is Hamiltonian. Adjacent pairs are exempt because the edge already joins them.

## Ore closure

Repeatedly joining nonadjacent pairs whose degree sum reaches \(n\) gives the closure of the graph. Ore's theorem says the closure is Hamiltonian exactly when the original graph is.

## Comparing the hypotheses

Dirac's condition implies Ore's, because two degrees of at least \(n/2\) sum to at least \(n\). The converse fails, so Ore's theorem covers strictly more graphs.

To check your understanding, identify a nonadjacent pair, join a qualifying pair and derive Ore from Dirac.

Pairing the degrees of nonadjacent vertices is weaker than bounding every degree, and still forces a spanning cycle.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
