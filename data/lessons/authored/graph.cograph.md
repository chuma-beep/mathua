# Cographs

Cographs are built from single vertices by disjoint union and join. They are exactly the graphs with no induced path on four vertices.

## Cotrees

A cotree records how a cograph is built. Leaves are vertices, and internal nodes are labelled union or join. A binary cotree with n leaves has n minus one internal nodes.

## Complement and joins

The complement of a cograph is obtained by swapping union and join labels in its cotree. A join adds every edge between its two parts, while a union adds none. Counting edges therefore separates into the edges inside the parts and the edges added by joins.

## P4-free graphs

A P4 is an induced path on four vertices. Cographs are exactly the graphs that contain no induced P4. Checking every four-vertex subset decides P4-freeness on a small graph.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
