# Graph Automorphisms

## Symmetries of a Graph

An automorphism is a bijection from a graph to itself that preserves adjacency. Automorphisms are exactly the isomorphisms whose domain and codomain are the same graph. They describe the symmetries that survive when the labels are ignored.

The identity map is always an automorphism, and the set of all automorphisms forms a group under composition. That group records how many distinct labelings of the graph represent the same unlabeled structure.

## Counting Automorphisms

The complete graph K_n has n! automorphisms because every permutation of its vertices preserves adjacency. The cycle graph C_n has 2n automorphisms, coming from n rotations and n reflections. The path graph P_n has only two automorphisms, the identity and the reversal.

For a complete bipartite graph with parts of different sizes, the count is m! times n!, because each part must map onto itself. When the parts have equal size, the two parts can also be swapped, doubling the count.

## Automorphism Groups

Automorphism groups connect graph theory with group theory: every finite group occurs as the automorphism group of some graph, although the construction may need many vertices. This makes graph symmetry a rich source of examples.

Counting automorphisms is often done by fixing a small set of vertices that determines the rest of the map. The orbit-stabilizer relation then turns the count into a product of choices, which is why highly symmetric graphs have large automorphism groups.
