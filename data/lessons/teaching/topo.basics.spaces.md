# Topological spaces

## What a Topology Is

A topology is a chosen family of subsets declared open. It must contain the empty set and the whole space, and it must be closed under arbitrary unions and finite intersections. Those three requirements are the whole definition, and every later construction in topology respects them.

The definition is deliberately sparse. It does not mention distance, angles or order, so the same set can carry many different topologies. A coarser topology has fewer open sets and makes continuity easier to achieve; a finer topology has more open sets and makes it easier to separate points.

Comparison with a metric space is useful. Every metric space carries its metric topology, whose open sets are unions of open balls, but a topological space need not come from any metric. This generality is what lets the same arguments apply to function spaces, quotient spaces and algebraic varieties.

## Comparing Topologies

Two topologies on the same set are compared by inclusion. If every set open in one topology is open in the other, the first is coarser and the second is finer. The discrete topology is the finest possible and the indiscrete topology is the coarsest possible.

Comparison matters because topological properties are not uniform across coarseness. Compactness and connectedness become easier to satisfy in a coarse topology, while the Hausdorff property and separation become easier in a fine one. Continuity from a fixed domain is easier into a coarse codomain and harder into a fine one.

When a problem asks for the coarsest topology satisfying a condition, the answer is usually a generated topology. One starts with a subbasis, closes under finite intersections, then closes under arbitrary unions. The resulting family is the smallest topology containing the original collection, and it is unique.

## Open Sets in Practice

On a finite set, a topology can be written out completely. The discrete topology on n points has 2 to the n open sets, the indiscrete topology has two, and many intermediate examples exist. Counting open sets is therefore a concrete way to compare topologies.

The particular point topology fixes a point p and declares a set open when it is empty or contains p. The excluded point topology declares a set open when it is the whole space or omits p. Both are legitimate topologies and both are useful as small counterexamples.

Most arguments in point-set topology use only membership in open sets. To prove that a set is open, exhibit it as a union of basis elements or as a preimage of an open set under a continuous map. To prove that it is not open, find a point with no open neighbourhood contained in the set.

