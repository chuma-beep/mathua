# Depth-first search

Depth-first search follows one path as far as possible before backtracking. Its discovery and finish times describe the structure of the search tree.

## Discovery and finish times

A clock increments when a vertex is first reached and again when the search leaves it. The discovery time is the first increment and the finish time is the second. The interval between them contains the intervals of all descendants.

## The DFS tree

The edge used to discover an unvisited vertex becomes a tree edge. The vertex from which it was discovered is the parent. The root has no parent and is recorded as minus 1.

## Why order matters

Visiting neighbours in increasing order makes the search deterministic. A different neighbour order can produce a different tree. The question always states the order so the answer is unique.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
