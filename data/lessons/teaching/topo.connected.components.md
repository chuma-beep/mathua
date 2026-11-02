# Connected components

## Components as Equivalence Classes

Declare two points equivalent when some connected subset contains both. The relation is reflexive because a singleton is connected, symmetric by definition, and transitive because the union of two connected sets that meet is connected. The equivalence classes are the connected components.

Each component is connected, and distinct components are separated. In fact every component is closed, because the closure of a connected set is connected. A space is connected exactly when it has one component.

Components give a canonical decomposition of a space into connected pieces. The decomposition is not usually a quotient map in the naive sense, but it is well behaved for locally connected spaces, where components are also open.

## Counting Components

Counting components is often easier than proving connectedness directly. A disjoint union of k nonempty connected spaces has k components. A circle has one component, two disjoint circles have two, and a space with an isolated point has at least one component consisting of that point alone.

Removing points can change the count. The real line is connected, but removing n distinct points leaves n plus one open intervals, hence n plus one components. Removing a closed interval from the line leaves two unbounded components.

Products behave predictably when the factors are nonempty. The product of a space with k components and a space with m components has k times m components, because a product is connected exactly when both factors are connected.

## Components and Continuity

A continuous map sends connected sets to connected sets, so it sends each component into a single component of the target. This gives a necessary condition for continuity and a quick way to rule out maps that would split a component.

The converse fails: a map that preserves components need not be continuous. Components are a coarse invariant, and many different maps preserve them. Stronger invariants such as path components or homology are needed when the component count is not enough.

In a locally connected space, components and path components coincide and are open. In the topologist's sine curve they do not coincide: the space is connected but has two path components. That example is the standard warning against treating the two notions as interchangeable.

