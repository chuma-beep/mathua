# Choosability

Choosability asks for the smallest list size that guarantees a choice of distinct colours from arbitrary lists.

## The definition

A graph is \(k\)-choosable when every assignment of lists of size \(k\) admits a proper colouring that picks one colour from each list. The choosability \(\mathrm{ch}(G)\) is the least such \(k\).

## Relation to chromatic number

Every graph satisfies \(\chi(G)\le\mathrm{ch}(G)\), because a common list of size \(\chi(G)\) reduces to ordinary colouring. Equality can fail, and the smallest failures are complete bipartite graphs.

## Basic values

A path is 2-choosable, an even cycle is 2-choosable and an odd cycle is 3-choosable. The complete graph \(K_n\) has choosability \(n\), since all lists can be identical.

To check your understanding, assign a list to each vertex, use identical lists and colour a path from lists.

Lists are harder than a shared palette, and the gap between the two numbers is real.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
