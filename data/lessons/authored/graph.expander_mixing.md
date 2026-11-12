# The expander mixing lemma

The expander mixing lemma says a small second eigenvalue forces every pair of vertex sets to meet in nearly the expected number of edges.

## The estimate

For a \(d\)-regular graph with second adjacency eigenvalue \(\lambda\), every two sets \(S,T\) satisfy \(\left|e(S,T)-\frac{d|S||T|}{n}\right|\le\lambda\sqrt{|S||T|}\). The first term is the edge count a random \(d\)-regular graph would give.

## Consequences

When \(\lambda\) is much smaller than \(d\), the graph behaves like a random graph for edge counts. Every set of reasonable size has many neighbours, so the graph expands.

## The parameters

The expected number of edges between \(S\) and \(T\) is \(d|S||T|/n\). The deviation is controlled by \(\lambda\), which is at most \(d\) and is zero only for a union of complete graphs.

To check your understanding, compute the expected edges, take large sets and divide by the order.

A spectral gap turns into a uniform statement about how evenly edges spread between sets.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
