# The handshake lemma

The handshake lemma counts each edge twice when the degrees are summed. It is the first constraint every degree sequence must satisfy.

## Summing the degrees

For a graph with \(m\) edges and degrees \(d_1,\dots,d_n\), \(\sum_i d_i = 2m\). Each edge contributes one to the degree of each endpoint, so the total is even.

## Counting odd vertices

The number of vertices of odd degree is even, because the sum of the even degrees is even and the total is even. A degree sequence with an odd number of odd entries cannot be drawn.

## Average degree

The average degree is \(2m/n\). A graph with average degree \(d\) therefore has \(dn/2\) edges, and a graph with minimum degree \(\delta\) has at least \(\delta n/2\) edges.

To check your understanding, double the edge count, count odd degrees and divide the degree sum by the order.

Every edge is counted at both ends, and that single observation constrains every degree list.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
