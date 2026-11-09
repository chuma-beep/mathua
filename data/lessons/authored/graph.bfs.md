# Breadth-first search

Breadth-first search visits vertices in layers by distance. It finds shortest paths in an unweighted graph.

## Layers

The source is at distance 0. Its neighbours are at distance 1. Every vertex is reached by a shortest path, so its distance is the number of layers between it and the source.

## The queue order

A queue processes vertices in the order they were discovered. All vertices at one distance are processed before any vertex at the next distance. The visit order is therefore nondecreasing in distance.

## Counting a layer

The number of vertices at a given distance can be counted from the distance array. It is not the same as the number of vertices within that distance. The layer is exactly the set of vertices whose distance equals the target distance.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
