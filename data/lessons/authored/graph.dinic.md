# Dinic's algorithm

Dinic's algorithm augments along shortest paths in a layered network and finishes in a bounded number of phases.

## The level graph

Breadth-first search labels each vertex with its distance from the source. The level graph keeps only edges that advance one level, and a blocking flow saturates every source-sink path in it.

## The phase count

The distance from source to sink strictly increases after each blocking flow, so there are at most \(n-1\) phases. On a unit-capacity network the bound improves to \(O(\sqrt{E})\) phases.

## The running time

Each phase costs \(O(VE)\) with a simple blocking-flow routine, giving \(O(V^2E)\) overall. The unit-capacity case runs in \(O(E\sqrt{E})\), which is why Dinic is the default choice for matching as flow.

To check your understanding, run breadth-first search, see the distance grow and multiply phase cost by phases.

Layering turns a flow problem into a sequence of acyclic blocking-flow problems, and each phase raises the shortest path length.
