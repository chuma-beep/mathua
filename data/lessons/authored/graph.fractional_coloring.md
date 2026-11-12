# Fractional coloring

Fractional colouring allows colours to be shared, and the resulting optimum lies between the clique number and the chromatic number.

## The relaxation

Assign a nonnegative weight to each independent set so that every vertex receives total weight at least one. The fractional chromatic number is the minimum total weight, written \(\chi_f(G)\).

## Bounds

Every graph satisfies \(\omega(G)\le\chi_f(G)\le\chi(G)\). For a vertex-transitive graph, \(\chi_f(G)=n/\alpha(G)\), the order divided by the independence number.

## Values

The complete graph \(K_n\) has \(\chi_f=n\). A bipartite graph with at least one edge has \(\chi_f=2\), and the cycle \(C_5\) has \(\chi_f=5/2\). The last value is fractional even though every chromatic number is an integer.

To check your understanding, weight the independent sets, compare with the clique number and compute a complete graph.

Sharing colours across independent sets gives a linear programme whose optimum can be fractional.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
