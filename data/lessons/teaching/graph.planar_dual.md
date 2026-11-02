# Planar Dual Graphs

## Faces and Euler's Formula

A planar drawing divides the plane into faces, including the unbounded outside face. Euler's formula for a connected planar graph says V - E + F = 2. Rearranging gives F = E - V + 2, which is the quickest way to count faces once the drawing is fixed.

The formula assumes that the graph is connected and that the drawing has no unnecessary crossings. For disconnected graphs the count changes by the number of components, so the connected case is the one used in most exercises.

## Constructing the Dual

The dual graph places one vertex inside each face of the original drawing and joins two dual vertices whenever the corresponding faces share an edge. Bridges become loops in the dual, because the same face lies on both sides of a bridge.

The dual depends on the drawing, not only on the abstract graph, although different planar drawings of the same connected graph give duals that are related by planar isomorphism. The number of dual vertices always equals the number of faces.

## Dual Edge Counts

Each edge of the original graph corresponds to one edge of the dual, so the two graphs have the same number of edges. The dual of the dual returns the original connected planar graph. These two facts let you move counts back and forth without redrawing.

Face counts, dual vertex counts and dual edge counts are therefore three views of the same Euler-formula calculation. Checking one against another is a useful way to catch arithmetic slips in a hand computation.
