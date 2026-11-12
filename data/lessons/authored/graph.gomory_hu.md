# Gomory-Hu trees

A Gomory-Hu tree compresses all pairwise minimum cuts of an undirected graph into a single weighted tree.

## The tree

A Gomory-Hu tree has the same vertices as the graph and \(n-1\) weighted edges. The minimum cut between any two vertices equals the minimum edge weight on the tree path joining them.

## Construction cost

Building the tree needs \(n-1\) maximum-flow computations, one per added tree edge. Each computation splits a current component, so the total is polynomial and independent of the number of vertex pairs.

## Cuts on the tree

Removing a tree edge splits the vertices into the two sides of the corresponding minimum cut. The \(n-1\) tree edges therefore encode \(n-1\) distinct cuts that answer all \(\binom{n}{2}\) pairwise queries.

To check your understanding, count the tree edges, count the flow computations and remove one tree edge.

One weighted tree answers every pairwise minimum-cut question in the original graph.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
