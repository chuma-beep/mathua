> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_trees.html) by Oscar Levin — CC BY-SA 4.0

# Trees

## Section 4.2 Trees

###### Investigate!

1. Find a subgraph with the smallest number of edges that is still connected and contains all the vertices.ð

ð
1. Find a subgraph with the largest number of edges that doesnât contain any cycles.ð

ð
1. What do you notice about the number of edges in your examples above?  Is this a coincidence?ð

ð
### Definition of a Tree.

### Subsection  Properties of Trees

#### Proposition 4.2.1.

#### Proof.

#### Corollary 4.2.2.

#### Proposition 4.2.3.

#### Proof.

#### Proposition 4.2.4.

#### Proof.

### Subsection  Rooted Trees

#### Example 4.2.5.

#### Example 4.2.6.

### Subsection  Spanning Trees

#### Spanning tree.

#### Example 4.2.7.

### Exercises  Exercises

#### 1.

1. \(G = (V, E)\) with \(V = \{a, b, c, d, e\}\) and \(E = \{\{a, b\}, \{a,e\}, \{b, c\}, \{c,d\}, \{d,e\} \}\)ð

ð
1. \(G = (V, E)\) with \(V = \{a, b, c, d, e\}\) and \(E = \{\{a, b\}, \{b, c\}, \{c,d\}, \{d,e\}\}\)ð

ð
1. \(G = (V, E)\) with \(V = \{a, b, c, d, e\}\) and \(E = \{\{a, b\}, \{a, c\}, \{a,d\}, \{a,e\}\}\)ð

ð
1. \(G = (V, E)\) with \(V = \{a, b, c, d, e\}\) and \(E = \{\{a, b\}, \{a, c\}, \{d,e\}\}\)ð

ð
1. This is not a tree since it contains a cycle. Note also that there are too many edges to be a tree, since we know that all trees with \(v\) vertices have \(v-1\) edges.ð

ð
1. This is a tree since it is connected and contains no cycles (which you can see by drawing the graph). All paths are trees.ð

ð
1. This is a tree since it is connected and contains no cycles (draw the graph). All stars are trees.ð

ð
1. This is a not a tree since it is not connected. Note that there are not enough edges to be a tree.ð

ð
#### 2.

1. \(\displaystyle (4,1,1,1,1)\)ð

ð
1. \(\displaystyle (3,3,2,1,1)\)ð

ð
1. \(\displaystyle (2,2,2,1,1)\)ð

ð
1. \(\displaystyle (4, 4, 3, 3, 3, 2, 2, 1, 1, 1, 1, 1, 1, 1)\)ð

ð
1. This must be the degree sequence for a tree. This is because the vertex of degree 4 must be adjacent to the four vertices of degree 1 (there are no other vertices for it to be adjacent to), and thus we get a star.ð

ð
1. This cannot be a tree. Each degree 3 vertex is adjacent to all but one of the vertices in the graph. Thus each must be adjacent to one of the degree 1 vertices (and not the other). That means both degree 3 vertices are adjacent to the degree 2 vertex, and to each other, so that means there is a cycle.ð

Alternatively, count how many edges there are!ð

ð
1. This might or might not be a tree. The length 4 path has this degree sequence (this is a tree), but so does the union of a 3-cycle and a length 1 path (which is not connected, so not a tree).ð

ð
1. This cannot be a tree. The sum of the degrees is 28, so there are 14 edges. But there are 14 vertices as well, so we donât have \(v = e+1\text{,}\) meaning this cannot be a tree.ð

ð
#### 3.

1. \(\displaystyle (3, 3, 2, 2, 2)\)ð

ð
1. \(\displaystyle (3, 2, 2, 1, 1, 1)\)ð

ð
1. \(\displaystyle (3, 3, 3, 1, 1, 1)\)ð

ð
1. \(\displaystyle (4, 4, 1, 1, 1, 1, 1, 1)\)ð

ð
#### 4.

#### 5.

#### 6.

#### 7.

1. Explain why this is a good name. That is, explain why a forest is a union of trees.ð

ð
1. Suppose \(F\) is a forest consisting of \(m\) trees and \(v\) vertices. How many edges does \(F\) have? Explain.ð

ð
1. Prove that any graph \(G\) with \(v\) vertices and \(e\) edges that satisfies \(v \lt e+1\) must contain a cycle (i.e., not be a forest).ð

ð
#### 8.

#### 9.

#### 10.

1. Suppose we designate vertex \(e\) as the root.  List the children, parents and siblings of each vertex.  Does any vertex other than \(e\) have grandchildren?ð

ð
1. Suppose \(e\) is not chosen as the root.  Does our choice of root vertex change the number of children \(e\) has?  The number of grandchildren?  How many are there of each?ð

ð
1. In fact, pick any vertex in the tree and suppose it is not the root.  Explain why the number of children of that vertex does not depend on which other vertex is the root.ð

ð
1. Does the previous part work for other trees?  Give an example of a different tree for which it holds.  Then either prove that it always holds or give an example of a tree for which it doesnât.ð

ð
#### 11.

#### 12.

1. Must all spanning trees of a given graph be isomorphic to each other? Explain why or give a counterexample.ð

ð
1. Must all spanning trees of a given graph have the same number of edges? Explain why or give a counterexample.ð

ð
1. Must all spanning trees of a graph have the same number of leaves (vertices of degree 1)? Explain why or give a counterexample.ð

ð
1. No, although there are graphs for which this is true. For example, \(K_4\) has a spanning tree that is a path (of three edges) and also a spanning tree that is a star (with center vertex of degree 3).ð

ð
1. Yes. For a fixed graph, we have a fixed number \(v\) of vertices. Any spanning tree of the graph will also have \(v\) vertices, and since it is a tree, must have \(v-1\) edges.ð

ð
1. No, although there are graph for which this is true (note that if all spanning trees are isomorphic, then all spanning trees will have the same number of leaves). Again, \(K_4\) is a counterexample. One spanning tree is a path, with only two leaves, another spanning tree is a star with 3 leaves.ð

ð
#### 13.

#### 14.

#### 15.

#### 16.