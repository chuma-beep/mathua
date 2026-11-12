# Multiway cuts

A multiway cut separates a set of terminals from each other by deleting as few edges as possible.

## The problem

Given terminals \(t_1,\dots,t_k\), a multiway cut removes edges so that no two terminals remain connected. For \(k=2\) it reduces to an ordinary minimum cut.

## The isolating cut bound

Isolating each terminal with its own minimum cut and deleting all but the largest gives a cut of cost at most \(2-2/k\) times the optimum. The bound is tight on a star whose leaves are the terminals.

## The star case

In a star with \(k\) leaves as terminals, deleting all \(k\) spokes separates every pair at cost \(k\). Any solution must delete at least \(k-1\) edges, so the star is a useful test of a heuristic.

To check your understanding, name the terminals, cut around each terminal and count the spokes.

Separating many terminals at once costs more than one cut, and the isolating-cut heuristic bounds how much more.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
