# Travelling salesman relaxations

The travelling salesman problem asks for the shortest tour that visits every city once. Its standard relaxation assigns two edges to every city, which permits a set of disjoint subtours rather than a single tour.

## The degree relaxation

The assignment relaxation is easy and its bound is usually weak. The gap comes from solutions made of several disjoint cycles, each satisfying the degree constraints locally.

## Subtour elimination

A cut forbids a subset of cities from forming a closed cycle. There are exponentially many subsets, so the cuts are separated on demand by a min-cut computation.

## The strength ladder

Adding subtour cuts, then comb inequalities, then general blossom inequalities tightens the bound step by step. Each level costs a harder separation routine and buys a smaller gap.

To check your understanding, assign two edges per city, separate a city subset and add subtour cuts.

The salesman relaxation is a degree problem plus cycle elimination, and the value of a cut family is measured by how much of the gap it closes.
