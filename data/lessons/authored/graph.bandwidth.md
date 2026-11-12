# Graph bandwidth

Bandwidth is the smallest stretch of any edge under a linear layout of the vertices.

## The definition

Bandwidth is the minimum over vertex orderings of the largest difference between the positions of an edge's endpoints. Equivalently, every window of \(b+1\) consecutive positions is a vertex separator of the graph.

## Values

A path has bandwidth one, a cycle has bandwidth two and the complete graph \(K_n\) has bandwidth \(n-1\). A star with \(n\) leaves has bandwidth \(\lceil n/2\rceil\).

## Relation to pathwidth

Bandwidth is at least the pathwidth and at most a constant times it. Both measure how close the graph is to a path, but bandwidth fixes the vertex order more rigidly.

To check your understanding, order the vertices, lay out a path and compare the two parameters.

Bandwidth measures the longest edge in the best linear layout, and that single number controls many layout costs.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
