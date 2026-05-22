> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Graph_%28discrete_mathematics%29) — CC BY-SA 4.0

# Graphs: vertices, edges, and degrees

Sets of vertices connected by edgesgraphs of mathematical functionsGraph of a functionother usesGraph (disambiguation)


In discrete mathematics, particularly in graph theory, a **graph** is a structure consisting of a set of objects where some pairs of the objects are in some sense "related". The objects are represented by abstractions called *vertices* (also called *nodes* or *points*) and each of the related pairs of vertices is called an *edge* (also called *link* or *line*). Typically, a graph is depicted in diagrammatic form as a set of dots or circles for the vertices, joined by lines or curves for the edges.

The edges may be directed or undirected. For example, if the vertices represent people at a party, and there is an edge between two people if they shake hands, then this graph is undirected because any person *A* can shake hands with a person *B* only if *B* also shakes hands with *A*. In contrast, if an edge from a person *A* to a person *B* means that *A* owes money to *B*, then this graph is directed, because owing money is not necessarily reciprocated.

Graphs are the basic subject studied by graph theory. The word "graph" was first used in this sense by J. J. Sylvester in 1878 due to a direct relation between mathematics and chemical structure (what he called a chemico-graphical image).

## Definitions

Definitions in graph theory vary. The following are some of the more basic ways of defining graphs and related mathematical structures.

###  Graph


A **graph** (sometimes called an *undirected graph* to distinguish it from a directed graph, or a *simple graph* to distinguish it from a multigraph) is a pair *G* = (*V*, *E*), where \(V\) is a set whose elements are called *vertices* (singular: vertex), and \(E\) is a set of unordered pairs \(\{v_1, v_2\}\) of vertices, whose elements are called *edges* (sometimes *links* or *lines*).

An empty graph is a graph that has an empty set of vertices (and thus an empty set of edges). The *order* of a graph is its number {abs|*V*} of vertices, usually denoted by \(n\). The *size* of a graph is its number {abs|*E*} of edges, typically denoted by \(m\). However, in some contexts, such as for expressing the computational complexity of algorithms, the term *size* is used for the quantity {abs|*V*} + {abs|*E*} (otherwise, a non-empty graph could have size 0). The *degree* or *valency* of a vertex is the number of edges that are incident to it; for graphs with loops, a loop is counted twice.

The vertices \(u\) and \(v\) of an edge {*u*, *v*} are called the edge's *endpoints*. The edge is said to *join* \(u\) and \(v\) and to be *incident* on them. A vertex may belong to no edge, in which case it is not joined to any other vertex and is called *isolated*. When an edge \(\{u,v\}\) exists, the vertices \(u\) and \(v\) are called *adjacent*.

A multigraph is a generalization that allows multiple edges to have the same pair of endpoints. In some texts, multigraphs are simply called graphs.

Sometimes, graphs are allowed to contain *loops*, which are edges that join a vertex to itself. To allow loops, the pairs of vertices in \(E\) must be allowed to have the same node twice.  Such generalized graphs are called *graphs with loops* or simply *graphs* when it is clear from the context that loops are allowed.

Generally, the vertex set \(V\) is taken to be finite (which implies that the edge set \(E\) is also finite). Sometimes infinite graphs are considered, but they are usually viewed as a special kind of binary relation, because most results on finite graphs either do not extend to the infinite case or need a rather different proof.

In a graph of order *n*, the maximum degree of each vertex is *n* − 1 (or *n* + 1 if loops are allowed, because a loop contributes 2 to the degree), and the maximum number of edges is *n*(*n* − 1)/2 (or *n*(*n* + 1)/2 if loops are allowed).

The edges of a graph define a symmetric relation on the vertices, called the *adjacency relation*. Specifically, two vertices \(x\) and \(y\) are *adjacent* if {*x*, *y*} is an edge. A graph is fully determined by its adjacency matrix \(A\), which is an *n* × *n* square matrix, with \(A{sub|ij}\) specifying the number of connections from vertex \(i\) to vertex \(j\). For a simple graph, *A{sub|ij}* is either 0, indicating disconnection, or 1, indicating connection; moreover *A{sub|ii}* = 0 because an edge in a simple graph cannot start and end at the same vertex. Graphs with self-loops will be characterized by some or all \(A{sub|ii}\) being equal to a positive integer, and multigraphs (with multiple edges between vertices) will be characterized by some or all \(A{sub|ij}\) being equal to a positive integer. Undirected graphs will have a symmetric adjacency matrix (meaning *A{sub|ij}* = *A{sub|ji}*).

### Directed graph


A **directed graph** or **digraph** is a graph in which edges have orientations.

In one restricted but very common sense of the term, a **directed graph** is a pair *G* = (*V*, *E*) comprising:
* \(V\), a set of *vertices* (also called *nodes* or *points*);
* \(E\), a set of *edges* (also called *directed edges*, *directed links*, *directed lines*, *arrows*, or *arcs*), which are ordered pairs of distinct vertices: \(E \subseteq \{(x,y) \mid (x,y) \in V^2 \;\textrm{ and }\; x \neq y \}\).
To avoid ambiguity, this type of object may be called precisely a **directed simple graph**.

In the edge (*x*, *y*) directed from \(x\) to \(y\), the vertices \(x\) and \(y\) are called the *endpoints* of the edge, \(x\) the *tail* of the edge and \(y\) the *head* of the edge. The edge is said to *join* \(x\) and \(y\) and to be *incident* on \(x\) and on \(y\). A vertex may exist in a graph and not belong to an edge. The edge (*y*, *x*) is called the *inverted edge* of (*x*, *y*). *Multiple edges*, not allowed under the definition above, are two or more edges with both the same tail and the same head.

In one more general sense of the term allowing multiple edges, a directed graph is sometimes defined to be an ordered triple *G* = (*V*, *E*, *ϕ*) comprising:
* \(V\), a set of *vertices* (also called *nodes* or *points*);
* \(E\), a set of *edges* (also called *directed edges*, *directed links*, *directed lines*, *arrows* or *arcs*);
* \(ϕ\), an *incidence function* mapping every edge to an ordered pair of vertices (that is, an edge is associated with two distinct vertices): \(\phi : E \to \{(x,y) \mid (x,y) \in V^2 \;\textrm{ and }\; x \neq y \}\).

To avoid ambiguity, this type of object may be called precisely a **directed multigraph**.

A *loop* is an edge that joins a vertex to itself. Directed graphs as defined in the two definitions above cannot have loops, because a loop joining a vertex \(x\) to itself is the edge (for a directed simple graph) or is incident on (for a directed multigraph) \((x,x)\) which is not in \(\{(x,y) \mid (x,y) \in V^2 \;\textrm{ and }\; x \neq y \}\). So to allow loops the definitions must be expanded. For directed simple graphs, the definition of \(E\) should be modified to \(E \subseteq V^2\). For directed multigraphs, the definition of \(\phi\) should be modified to \(\phi : E \to V^2\). To avoid ambiguity, these types of objects may be called precisely a **directed simple graph permitting loops** and a **directed multigraph permitting loops** (or a *quiver*) respectively.

The edges of a directed simple graph permitting loops \(G\) is a homogeneous relation ~ on the vertices of \(G\) that is called the *adjacency relation* of \(G\). Specifically, for each edge (*x*, *y*), its endpoints \(x\) and \(y\) are said to be *adjacent* to one another, which is denoted *x* ~ *y*.

### Mixed graph


A *mixed graph* is a graph in which some edges may be directed and some may be undirected. It is an ordered triple *G* = (*V*, *E*, *A*) for a *mixed simple graph* and *G* = (*V*, *E*, *A*, *ϕ{sub|E}*, *ϕ{sub|A}*) for a *mixed multigraph* with \(V\), \(E\) (the undirected edges), \(A\) (the directed edges), \(ϕ{sub|E}\) and \(ϕ{sub|A}\) defined as above. Directed and undirected graphs are special cases.

### Weighted graph


A *weighted graph* or a *network* is a graph in which a number (the weight) is assigned to each edge. Such weights might represent for example costs, lengths or capacities, depending on the problem at hand. Such graphs arise in many contexts, for example in shortest path problems such as the traveling salesman problem.

## Types of graphs
### Oriented graph
One definition of an *oriented graph* is that it is a directed graph in which at most one of (*x*, *y*) and (*y*, *x*) may be edges of the graph. That is, it is a directed graph that can be formed as an orientation of an undirected (simple) graph.

Some authors use "oriented graph" to mean the same as "directed graph".  Some authors use "oriented graph" to mean any orientation of a given undirected graph or multigraph.

### Regular graph

A *regular graph* is a graph in which each vertex has the same number of neighbours, i.e., every vertex has the same degree. A regular graph with vertices of degree *k* is called a *k*‑regular graph or regular graph of degree *k*.

### Complete graph


A *complete graph* is a graph in which each pair of vertices is joined by an edge. A complete graph contains all possible edges.

### Finite graph
A *finite graph* is a graph in which the vertex set and the edge set are finite sets. Otherwise, it is called an *infinite graph*.

Most commonly in graph theory it is implied that the graphs discussed are finite. If the graphs are infinite, that is usually specifically stated.

### Connected graph

In an undirected graph, an unordered pair of vertices {mset|*x*, *y*} is called *connected* if a path leads from *x* to *y*. Otherwise, the unordered pair is called *disconnected*.

A *connected graph* is an undirected graph in which every unordered pair of vertices in the graph is connected. Otherwise, it is called a *disconnected graph*.

In a directed graph, an ordered pair of vertices (*x*, *y*) is called *strongly connected* if a directed path leads from *x* to *y*. Otherwise, the ordered pair is called *weakly connected* if an undirected path leads from *x* to *y* after replacing all of its directed edges with undirected edges. Otherwise, the ordered pair is called *disconnected*.

A *strongly connected graph* is a directed graph in which every ordered pair of vertices in the graph is strongly connected. Otherwise, it is called a *weakly connected graph* if every ordered pair of vertices in the graph is weakly connected. Otherwise it is called a *disconnected graph*.

A *k-vertex-connected graph* or *k-edge-connected graph* is a graph in which no set of *k* − 1 vertices (respectively, edges) exists that, when removed, disconnects the graph. A *k*-vertex-connected graph is often called simply a *k-connected graph*.

### Bipartite graph

A *bipartite graph* is a simple graph in which the vertex set can be partitioned into two sets, *W* and *X*, so that no two vertices in *W* share a common edge and no two vertices in *X* share a common edge. Alternatively, it is a graph with a chromatic number of 2.

In a complete bipartite graph, the vertex set is the union of two disjoint sets, *W* and *X*, so that every vertex in *W* is adjacent to every vertex in *X* but there are no edges within *W* or *X*.

### Path graph

A *path graph* or *linear graph* of order *n* ≥ 2 is a graph in which the vertices can be listed in an order *v*\(_{1}\), *v*\(_{2}\), …, *v*\(_{*n*}\) such that the edges are the {mset|*v*\(_{*i*}\), *v*\(_{*i*+1}\)} where *i* = 1, 2, …, *n* − 1. Path graphs can be characterized as connected graphs in which the degree of all but two vertices is 2 and the degree of the two remaining vertices is 1. If a path graph occurs as a subgraph of another graph, it is a path in that graph.

### Planar graph

A *planar graph* is a graph whose vertices and edges can be drawn in a plane such that no two of the edges intersect.

### Cycle graph

A *cycle graph* or *circular graph* of order *n* ≥ 3 is a graph in which the vertices can be listed in an order *v*\(_{1}\), *v*\(_{2}\), …, *v*\(_{*n*}\) such that the edges are the {mset|*v*\(_{*i*}\), *v*\(_{*i*+1}\)} where *i* = 1, 2, …, *n* − 1, plus the edge {mset|*v*\(_{*n*}\), *v*\(_{1}\)}. Cycle graphs can be characterized as connected graphs in which the degree of all vertices is 2. If a cycle graph occurs as a subgraph of another graph, it is a cycle or circuit in that graph.

### Tree

A *tree* is an undirected graph in which any two vertices are connected by *exactly one* path, or equivalently a connected acyclic undirected graph.

A *forest* is an undirected graph in which any two vertices are connected by *at most one* path, or equivalently an acyclic undirected graph, or equivalently a disjoint union of trees.

### Polytree


A *polytree* (or *directed tree* or *oriented tree* or *singly connected network*) is a directed acyclic graph (DAG) whose underlying undirected graph is a tree.

A *polyforest* (or *directed forest* or *oriented forest*) is a directed acyclic graph whose underlying undirected graph is a forest.

### Advanced classes
More advanced kinds of graphs are:
* Petersen graph and its generalizations;
* perfect graphs;
* cographs;
* chordal graphs;
* other graphs with large automorphism groups: vertex-transitive, arc-transitive, and distance-transitive graphs;
* strongly regular graphs and their generalizations distance-regular graphs.

## Properties of graphs
Glossary of graph theoryGraph property
Two vertices of a graph are called *adjacent* if they share a common edge. Two vertices of a directed graph are called *consecutive* if the head of the first one is the tail of the second one. Similarly, two vertices are called *adjacent* if they share a common edge (*consecutive* if the first one is the tail and the second one is the head of an edge), in which case the common edge is said to *join* the two vertices. An edge and a vertex on that edge are called *incident*.

The graph with only one vertex and no edges is called the *trivial graph*. A graph with only vertices and no edges is known as an *edgeless graph*. The graph with no vertices and no edges is sometimes called the *null graph* or *empty graph*, but the terminology is not consistent and not all mathematicians allow this object.

Normally, the vertices of a graph, by their nature as elements of a set, are distinguishable. This kind of graph may be called *vertex-labeled*. However, for many questions it is better to treat vertices as indistinguishable. (Of course, the vertices may be still distinguishable by the properties of the graph itself, e.g., by the numbers of incident edges.) The same remarks apply to edges, so graphs with labeled edges are called *edge-labeled*. Graphs with labels attached to edges or vertices are more generally designated as *labeled*. Consequently, graphs in which vertices are indistinguishable and edges are indistinguishable are called *unlabeled*. (In the literature, the term *labeled* may apply to other kinds of labeling, besides that which serves only to distinguish different vertices or edges.)

The category of directed multigraphs permitting loops is the comma category Set ↓ *D* where *D*: Set → Set is the functor taking a set *s* to *s* × *s*.

## Examples


* The diagram is a schematic representation of the graph with vertices \(V = \{1, 2, 3, 4, 5, 6\}\) and edges \(E = \{\{1, 2\}, \{1, 5\}, \{2, 3\}, \{2, 5\}, \{3, 4\}, \{4, 5\}, \{4, 6\}\}.\)
* In computer science, directed graphs are used to represent knowledge (e.g., conceptual graph), finite-state machines, and many other discrete structures.
* A binary relation *R* on a set *X* defines a directed graph. An element *x* of *X* is a direct predecessor of an element *y* of *X* if and only if *xRy*.
* A directed graph can model information networks such as Twitter, with one user following another.
*Particularly regular examples of directed graphs are given by the Cayley graphs of finitely-generated groups, as well as Schreier coset graphs
*In category theory, every small category has an underlying directed multigraph whose vertices are the objects of the category, and whose edges are the arrows of the category.  In the language of category theory, one says that there is a forgetful functor from the category of small categories to the category of quivers.

## Graph operations

There are several operations that produce new graphs from initial ones, which might be classified into the following categories:
* *unary operations*, which create a new graph from an initial one, such as:
** edge contraction,
** line graph,
** dual graph,
** complement graph,
** graph rewriting;
* *binary operations*, which create a new graph from two initial ones, such as:
** disjoint union of graphs,
** cartesian product of graphs,
** tensor product of graphs,
** strong product of graphs,
** lexicographic product of graphs,
** series–parallel graphs.

## Generalizations
In a hypergraph, an edge can join any positive number of vertices.

An undirected graph can be seen as a simplicial complex consisting of 1-simplices (the edges) and 0-simplices (the vertices). As such, complexes are generalizations of graphs since they allow for higher-dimensional simplices.

Every graph gives rise to a matroid.

In model theory, a graph is just a structure. But in that case, there is no limitation on the number of edges: it can be any cardinal number, see continuous graph.

In computational biology, power graph analysis introduces power graphs as an alternative representation of undirected graphs.

In geographic information systems, geometric networks are closely modeled after graphs, and borrow many concepts from graph theory to perform spatial analysis on road networks or utility grids.

