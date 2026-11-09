# Disjoint-set union

A disjoint-set structure maintains a partition of vertices under union operations. Union by size keeps the trees shallow and find returns a representative.

## Union by size

Each set starts as a singleton with size 1. A union attaches the smaller root under the larger root. The size of the new root is the sum of the two sizes. This keeps the depth logarithmic.

## Roots

The representative of a set is the root of its tree. Path compression points every visited vertex directly at the root. Two vertices are in the same set exactly when their roots are equal.

## Component size

The size stored at a root is the number of vertices in its component. After all unions, the root of a target vertex gives the component that contains it. The stored size is the answer.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
