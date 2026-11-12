# Hall's marriage theorem

Hall's theorem characterises when a bipartite graph has a matching that covers one whole side.

## The condition

A bipartite graph with sides \(A\) and \(B\) has a matching covering \(A\) exactly when every subset \(S\subseteq A\) has at least \(|S|\) neighbours in \(B\). The condition is necessary because a matching sends distinct vertices to distinct neighbours.

## Deficiency

When the condition fails, the largest matching leaves \(|A|-\max_S(|S|-|N(S)|)\) vertices of \(A\) unmatched. The deficit is the worst shortfall of neighbourhoods.

## Regular bipartite graphs

A \(k\)-regular bipartite graph with equal sides satisfies Hall's condition for every \(k\ge 1\), so it has a perfect matching. Repeating the argument after removing a perfect matching decomposes the edges into \(k\) perfect matchings.

To check your understanding, pick a subset of one side, find the worst subset and verify Hall on a regular graph.

A matching covering one side exists exactly when no subset of that side is short of neighbours.
