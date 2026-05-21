> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Tree_%28graph_theory%29) — CC BY-SA 4.0

# Trees and spanning trees

In graph theory, a **tree** is an undirected graph in which every pair of distinct vertices is connected by path, or equivalently, a connected acyclic undirected graph. A **forest** is an undirected graph in which any two vertices are connected by path, or equivalently an acyclic undirected graph, or equivalently a disjoint union of trees.

A directed tree, oriented tree, polytree, or singly connected network is a directed acyclic graph (DAG) whose underlying undirected graph is a tree. A polyforest (or directed forest or oriented forest) is a directed acyclic graph whose underlying undirected graph is a forest.

The various kinds of data structures referred to as trees in computer science have underlying graphs that are trees in graph theory, although such data structures are generally rooted trees. A rooted tree may be directed, called a directed rooted tree, either making all its edges point away from the root—in which case it is called an arborescence or out-tree—or making all its edges point towards the root—in which case it is called an anti-arborescence or in-tree. A rooted tree itself has been defined by some authors as a directed graph. A rooted forest is a disjoint union of rooted trees. A rooted forest may be directed, called a directed rooted forest, either making all its edges point away from the root in each rooted tree—in which case it is called a branching or out-forest—or making all its edges point towards the root in each rooted tree—in which case it is called an anti-branching or in-forest.

The term was coined in 1857 by the British mathematician Arthur Cayley.

## Definitions
### Tree
A *tree* is an undirected graph that satisfies any of the following equivalent conditions:
* is connected and acyclic (contains no cycles).
* is acyclic, and a simple cycle is formed if any edge is added to. * is connected, but would become disconnected if any single edge is removed from. * is connected and the complete graph *K* is not a minor of. * Any two vertices in can be connected by a unique simple path.
If has finitely many vertices, say of them, then the above statements are also equivalent to any of the following conditions:
* is connected and has *n* − 1 edges.
* is connected, and every subgraph of includes at least one vertex with zero or one incident edges. (That is, is connected and 1-degenerate.)
* has no simple cycles and has *n* − 1 edges.
As elsewhere in graph theory, the order-zero graph (graph with no vertices) is generally not considered to be a tree: while it is vacuously connected as a graph (any two vertices can be connected by a path), it is not 0-connected (or even (−1)-connected) in algebraic topology, unlike non-empty trees, and violates the "one more vertex than edges" relation. It may, however, be considered as a forest consisting of zero trees.

An (or inner vertex) is a vertex of degree at least 2. Similarly, an (or outer vertex, terminal vertex or leaf) is a vertex of degree 1. A branch vertex in a tree is a vertex of degree at least 3.

An (or series-reduced tree) is a tree in which there is no vertex of degree 2 (enumerated at sequence in the OEIS).

### Forest
A is an undirected acyclic graph or equivalently a disjoint union of trees. Trivially so, each connected component of a forest is a tree. As special cases, the order-zero graph (a forest consisting of zero trees), a single tree, and an edgeless graph, are examples of forests.
Since for every tree 1=*V* − *E* = 1, we can easily count the number of trees that are within a forest by subtracting the difference between total vertices and total edges. 1=*V* − *E* = number of trees in a forest.

### Polytree

A (or *directed tree* or *oriented tree* or *singly connected network*) is a directed acyclic graph (DAG) whose underlying undirected graph is a tree. In other words, if we replace its directed edges with undirected edges, we obtain an undirected graph that is both connected and acyclic.

Some authors restrict the phrase "directed tree" to the case where the edges are all directed towards a particular vertex, or all directed away from a particular vertex (see arborescence).

### Polyforest
A (or directed forest or oriented forest) is a directed acyclic graph whose underlying undirected graph is a forest. In other words, if we replace its directed edges with undirected edges, we obtain an undirected graph that is acyclic.

As with directed trees, some authors restrict the phrase "directed forest" to the case where the edges of each connected component are all directed towards a particular vertex, or all directed away from a particular vertex (see branching).

### Rooted tree
A is a tree in which one vertex has been designated the root. The edges of a rooted tree can be assigned a natural orientation, either away from or towards the root, in which case the structure becomes a directed rooted tree. When a directed rooted tree has an orientation away from the root, it is called an *arborescence* or *out-tree*; when it has an orientation towards the root, it is called an *anti-arborescence* or *in-tree*. The tree-order is the partial ordering on the vertices of a tree with *u* < *v* if and only if the unique path from the root to passes through. A rooted tree that is a subgraph of some graph is a normal tree if the ends of every -path in are comparable in this tree-order. Rooted trees, often with an additional structure such as an ordering of the neighbors at each vertex, are a key data structure in computer science; see tree data structure.

In a context where trees typically have a root, a tree without any designated root is called a *free tree*.

A *labeled tree* is a tree in which each vertex is given a unique label. The vertices of a labeled tree on vertices (for nonnegative integers ) are typically given the labels 1, 2, …, *n*. A *recursive tree* is a labeled rooted tree where the vertex labels respect the tree order (i.e., if *u* < *v* for two vertices and , then the label of is smaller than the label of ).

In a rooted tree, the **parent** of a vertex is the vertex connected to on the path to the root; every vertex has a unique parent, except the root has no parent. A *child* of a vertex is a vertex of which is the parent. An *ascendant* of a vertex is any vertex that is either the parent of or is (recursively) an ascendant of a parent of. A *descendant* of a vertex is any vertex that is either a child of or is (recursively) a descendant of a child of. A *sibling* to a vertex is any other vertex on the tree that shares a parent with. A *leaf* is a vertex with no children. An *internal vertex* is a vertex that is not a leaf.

The *height* of a vertex in a rooted tree is the length of the longest downward path to a leaf from that vertex. The *height* of the tree is the height of the root. The *depth* of a vertex is the length of the path to its root (*root path*). The depth of a tree is the maximum depth of any vertex. Depth is commonly needed in the manipulation of the various self-balancing trees, AVL trees in particular. The root has depth zero, leaves have height zero, and a tree with only a single vertex (hence both a root and leaf) has depth and height zero. Conventionally, an empty tree (a tree with no vertices, if such are allowed) has depth and height −1.

A *-ary tree* (for nonnegative integers ) is a rooted tree in which each vertex has at most children. 2-ary trees are often called *binary trees*, while 3-ary trees are sometimes called *ternary trees*.

### Ordered tree
An *ordered tree* (alternatively, *plane tree* or *positional tree*) is a rooted tree in which an ordering is specified for the children of each vertex. This is called a "plane tree" because an ordering of the children is equivalent to an embedding of the tree in the plane, with the root at the top and the children of each vertex lower than that vertex. Given an embedding of a rooted tree in the plane, if one fixes a direction of children, say left to right, then an embedding gives an ordering of the children. Conversely, given an ordered tree, and conventionally drawing the root at the top, then the child vertices in an ordered tree can be drawn left-to-right, yielding an essentially unique planar embedding.

## Properties
* Every tree is a bipartite graph. A graph is bipartite if and only if it contains no cycles of odd length. Since a tree contains no cycles at all, it is bipartite.
* Every tree with only countably many vertices is a planar graph.
* Every connected graph *G* admits a spanning tree, which is a tree that contains every vertex of *G* and whose edges are edges of *G*. More specific types spanning trees, existing in every connected finite graph, include depth-first search trees and breadth-first search trees. Generalizing the existence of depth-first-search trees, every connected graph with only countably many vertices has a Trémaux tree. However, some uncountable-order graphs do not have such a tree.
* Every finite tree with *n* vertices, with *n* > 1, has at least two terminal vertices (leaves). This minimal number of leaves is characteristic of path graphs; the maximal number, *n* − 1, is attained only by star graphs. The number of leaves is at least the maximum vertex degree.
* For any three vertices in a tree, the three paths between them have exactly one vertex in common. More generally, a vertex in a graph that belongs to three shortest paths among three vertices is called a median of these vertices. Because every three vertices in a tree have a unique median, every tree is a median graph.
* Every tree has a center consisting of one vertex or two adjacent vertices. The center is the middle vertex or middle two vertices in every longest path. Similarly, every *n*-vertex tree has a centroid consisting of one vertex or two adjacent vertices. In the first case removal of the vertex splits the tree into subtrees of fewer than *n*/2 vertices. In the second case, removal of the edge between the two centroidal vertices splits the tree into two subtrees of exactly *n*/2 vertices.
* The maximal cliques of a tree are precisely its edges, implying that the class of trees has few cliques.

## Enumeration
### Labeled trees
Cayley's formula states that there are *n* trees on labeled vertices. A classic proof uses Prüfer sequences, which naturally show a stronger result: the number of trees with vertices 1, 2, …, *n* of degrees *d*, …, *d*respectively, is the multinomial coefficient
\({n - 2 \choose d_1 - 1, d_2 - 1, \ldots, d_n - 1}.\)

A more general problem is to count spanning trees in an undirected graph, which is addressed by the matrix tree theorem. (Cayley's formula is the special case of spanning trees in a complete graph.) The similar problem of counting all the subtrees regardless of size is #P-complete in the general case ().

### Unlabeled trees
Counting the number of unlabeled free trees is a harder problem. No closed formula for the number *t*(*n*) of trees with vertices up to graph isomorphism is known. The first few values of *t*(*n*) are
1, 1, 1, 1, 2, 3, 6, 11, 23, 47, 106, 235, 551, 1301, 3159, …. proved the asymptotic estimate
\(t(n) \sim C \alpha^n n^{-5/2} \quad\text{as } n\to\infty,\)
with *C* ≈ 0.534949606... and *α* ≈ 2.95576528565.... Here, the ~ symbol means that
\(\lim_{n \to \infty} \frac{t(n)}{C \alpha^n n^{-5/2} } = 1.\)
This is a consequence of his asymptotic estimate for the number *r*(*n*) of unlabeled rooted trees with vertices:
\(r(n) \sim D\alpha^n n^{-3/2} \quad\text{as } n\to\infty,\)
with D ≈ 0.43992401257... and the same as above (cf. , chap. 2.3.4.4 and , chap. VII.5, p. 475).

The first few values of *r*(*n*) are
1, 1, 2, 4, 9, 20, 48, 115, 286, 719, 1842, 4766, 12486, 32973, …. ## Types of trees
* A *path graph* (or *linear graph*) consists of vertices arranged in a line, so that vertices and *i* + 1 are connected by an edge for 1=*i* = 1, …, *n* – 1.
* A *starlike tree* consists of a central vertex called *root* and several path graphs attached to it. More formally, a tree is starlike if it has exactly one vertex of degree greater than 2.
* A *star tree* is a tree which consists of a single internal vertex (and *n* – 1 leaves). In other words, a star tree of order is a tree of order with as many leaves as possible.
* A *caterpillar tree* is a tree in which all vertices are within distance 1 of a central path subgraph.
* A *lobster tree* is a tree in which all vertices are within distance 2 of a central path subgraph.
* A *regular tree* of degree is the infinite tree with edges at each vertex. These arise as the Cayley graphs of free groups, and in the theory of Tits buildings. In statistical mechanics they are known as *Bethe lattices*.
