# Causal graphs

A causal graph encodes which variables directly influence which others. Its edges are directed, and its paths carry the assumptions that make adjustment possible.

## The structure

A graph on \(n\) nodes has at most \(n(n-1)\) directed edges. A causal graph must be acyclic.

## Paths

A directed path follows the arrows, and a backdoor path does not. Separating the two kinds is the first step of every adjustment argument.

## d-separation

Two sets are d-separated when every path between them is blocked. Conditioning on a collider blocks a path rather than opening one.

To check your understanding, count the nodes, follow the arrows and block a path.

The graph is the assumption, and every later claim is relative to it.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
