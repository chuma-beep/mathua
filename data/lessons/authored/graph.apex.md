# Apex graphs

An apex graph becomes planar after deleting one vertex, so a single vertex carries every obstruction to planarity.

## One vertex away from planar

A graph is apex when some vertex \(v\) has \(G-v\) planar. Every planar graph is apex because deleting any vertex keeps it planar, and \(K_5\) is apex because deleting one vertex leaves \(K_4\), which is planar.

## The edge budget

Deleting one vertex removes its incident edges and leaves at most \(3(n-1)-6\) edges by the planar bound. The deleted vertex contributes at most \(n-1\) edges, so an apex graph on \(n\) vertices has at most \(4n-9\) edges.

## Obstructions

The apex property is minor-closed, and the forbidden minors are exactly the graphs that stay nonplanar after any single deletion. A graph with an apex vertex is planar plus one cone, which is why apex graphs sit between planar and toroidal graphs in the minor order.

To check your understanding, delete a vertex and test planarity, apply the planar bound after deletion and close the class under minors.

A single vertex can absorb the whole obstruction, and that is what the apex bound measures.
