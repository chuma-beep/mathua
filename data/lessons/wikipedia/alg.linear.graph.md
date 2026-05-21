> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Line_graph) — CC BY-SA 4.0

# Graph a linear equation

In the mathematical discipline of graph theory, the **line graph** of an undirected graph is another graph L(*G*) that represents the adjacencies between edges of. L(*G*) is constructed in the following way: for each edge in , make a vertex in L(*G*); for every two edges in that have a vertex in common, make an edge between their corresponding vertices in L(*G*).

The name *line graph* comes from a paper by although both and used the construction before this. Other terms used for the line graph include the **covering graph**, the **derivative**, the **edge-to-vertex dual**, the **conjugate**, the **representative graph**, and the **θ-obrazom**, as well as the **edge graph**, the **interchange graph**, the **adjoint graph**, and the **derived graph**.

 proved that with one exceptional case the structure of a connected graph can be recovered completely from its line graph. Many other properties of line graphs follow by translating the properties of the underlying graph from vertices into edges, and by Whitney's theorem the same translation can also be done in the other direction. Line graphs are claw-free, and the line graphs of bipartite graphs are perfect. Line graphs are characterized by nine forbidden subgraphs and can be recognized in linear time.

Various extensions of the concept of a line graph have been studied, including line graphs of line graphs, line graphs of multigraphs, line graphs of hypergraphs, and line graphs of weighted graphs.

## Formal definition
Given a graph , its line graph *L*(*G*) is a graph such that
* each vertex of *L*(*G*) represents an edge of ; and
* two vertices of *L*(*G*) are adjacent if and only if their corresponding edges share a common endpoint ("are incident") in. That is, it is the intersection graph of the edges of , representing each edge by the set of its two endpoints.

## Example
The following figures show a graph (left, with blue vertices) and its line graph (right, with green vertices). Each vertex of the line graph is shown labeled with the pair of endpoints of the corresponding edge in the original graph. For instance, the green vertex on the right labeled 1,3 corresponds to the edge on the left between the blue vertices 1 and 3. Green vertex 1,3 is adjacent to three other green vertices: 1,4 and 1,2 (corresponding to edges sharing the endpoint 1 in the blue graph) and 4,3 (corresponding to an edge sharing the endpoint 3 in the blue graph).

File:Line graph construction 1.svg|Graph *G*
File:Line graph construction 2.svg|Vertices in L(*G*) constructed from edges in *G*
File:Line graph construction 3.svg|Added edges in L(*G*)
File:Line graph construction 4.svg|The line graph L(*G*)

## Properties
### Translated properties of the underlying graph
Properties of a graph that depend only on adjacency between edges may be translated into equivalent properties in *L*(*G*) that depend on adjacency between vertices. For instance, a matching in is a set of edges no two of which are adjacent, and corresponds to a set of vertices in *L*(*G*) no two of which are adjacent, that is, an independent set.

Thus,
* The line graph of a connected graph is connected. If is connected, it contains a path connecting any two of its edges, which translates into a path in *L*(*G*) containing any two of the vertices of *L*(*G*). However, a graph that has some isolated vertices, and is therefore disconnected, may nevertheless have a connected line graph.
* A line graph has an articulation point if and only if the underlying graph has a bridge for which neither endpoint has degree one.
* For a graph with vertices and edges, the number of vertices of the line graph *L*(*G*) and the number of edges of *L*(*G*) is half the sum of the squares of the degrees of the vertices in , minus. * An independent set in *L*(*G*) corresponds to a matching in. In particular, a maximum independent set in *L*(*G*) corresponds to maximum matching in. Since maximum matchings may be found in polynomial time, so may the maximum independent sets of line graphs, despite the hardness of the maximum independent set problem for more general families of graphs. Similarly, a rainbow-independent set in *L*(*G*) corresponds to a rainbow matching in. * The edge chromatic number of a graph is equal to the vertex chromatic number of its line graph *L*(*G*).
* The line graph of an edge-transitive graph is vertex-transitive. This property can be used to generate families of graphs that (like the Petersen graph) are vertex-transitive but are not Cayley graphs: if is an edge-transitive graph that has at least five vertices, is not bipartite, and has odd vertex degrees, then *L*(*G*) is a vertex-transitive non-Cayley graph.
* If a graph has an Euler cycle, that is, if is connected and has an even number of edges at each vertex, then the line graph of is Hamiltonian. However, not all Hamiltonian cycles in line graphs come from Euler cycles in this way; for instance, the line graph of a Hamiltonian graph is itself Hamiltonian, regardless of whether is also Eulerian.
* If two simple graphs are isomorphic then their line graphs are also isomorphic. The Whitney graph isomorphism theorem provides a converse to this for all but one pair of connected graphs.
* In the context of complex network theory, the line graph of a random network preserves many of the properties of the network such as the small-world property (the existence of short paths between all pairs of vertices) and the shape of its degree distribution. observe that any method for finding vertex clusters in a complex network can be applied to the line graph and used to cluster its edges instead.

### Whitney isomorphism theorem

If the line graphs of two connected graphs are isomorphic, then the underlying graphs are isomorphic, except in the case of the triangle graph *K* and the claw *K*, which have isomorphic line graphs but are not themselves isomorphic.

As well as *K* and *K*, there are some other exceptional small graphs with the property that their line graph has a higher degree of symmetry than the graph itself. For instance, the diamond graph *K* (two triangles sharing an edge) has four graph automorphisms but its line graph *K* has eight. In the illustration of the diamond graph shown, rotating the graph by 90 degrees is not a symmetry of the graph, but is a symmetry of its line graph. However, all such exceptional cases have at most four vertices. A strengthened version of the Whitney isomorphism theorem states that, for connected graphs with more than four vertices, there is a one-to-one correspondence between isomorphisms of the graphs and isomorphisms of their line graphs.

Analogues of the Whitney isomorphism theorem have been proven for the line graphs of multigraphs, but are more complicated in this case.

### Strongly regular and perfect line graphs

The line graph of the complete graph is also known as the **triangular graph**, the Johnson graph *J*(*n*, 2), or the complement of the Kneser graph *KG*. Triangular graphs are characterized by their spectra, except for *n* = 8. They may also be characterized (again with the exception of *K*) as the strongly regular graphs with parameters srg(*n*(*n* − 1)/2, 2(*n* − 2), *n* − 2, 4). The three strongly regular graphs with the same parameters and spectrum as *L*(*K* are the Chang graphs, which may be obtained by graph switching from *L*(*K*.

The line graph of a bipartite graph is perfect (see Kőnig's theorem), but need not be bipartite as the example of the claw graph shows. The line graphs of bipartite graphs form one of the key building blocks of perfect graphs, used in the proof of the strong perfect graph theorem. A special case of these graphs are the rook's graphs, line graphs of complete bipartite graphs. Like the line graphs of complete graphs, they can be characterized with one exception by their numbers of vertices, numbers of edges, and number of shared neighbors for adjacent and non-adjacent points. The one exceptional case is *L*(*K*, which shares its parameters with the Shrikhande graph. When both sides of the bipartition have the same number of vertices, these graphs are again strongly regular. It has been shown that, except for *C* , *C* , and *C* , all connected strongly regular graphs can be made non-strongly regular within two line graph transformations. The extension to disconnected graphs would require that the graph is not a disjoint union of *C*.

More generally, a graph is said to be a line perfect graph if *L*(*G*) is a perfect graph. The line perfect graphs are exactly the graphs that do not contain a simple cycle of odd length greater than three. Equivalently, a graph is line perfect if and only if each of its biconnected components is either bipartite or of the form *K* (the tetrahedron) or *K* (a book of one or more triangles all sharing a common edge). Every line perfect graph is itself perfect.

### Other related graph families
All line graphs are claw-free graphs, graphs without an induced subgraph in the form of a three-leaf tree. As with claw-free graphs more generally, every connected line graph *L*(*G*) with an even number of edges has a perfect matching; equivalently, this means that if the underlying graph has an even number of edges, its edges can be partitioned into two-edge paths.

The line graphs of trees are exactly the claw-free block graphs. These graphs have been used to solve a problem in extremal graph theory, of constructing a graph with a given number of edges and vertices whose largest tree induced as a subgraph is as small as possible.

All eigenvalues of the adjacency matrix of a line graph are at least −2. The reason for this is that can be written as \(A = J^\mathsf{T}J - 2I\), where is the signless incidence matrix of the pre-line graph and is the identity. In particular, *A* + 2*I* is the Gramian matrix of a system of vectors: all graphs with this property have been called generalized line graphs.

## Characterization and recognition
### Clique partition

For an arbitrary graph , and an arbitrary vertex in , the set of edges incident to corresponds to a clique in the line graph *L*(*G*). The cliques formed in this way partition the edges of *L*(*G*). Each vertex of *L*(*G*) belongs to exactly two of them (the two cliques corresponding to the two endpoints of the corresponding edge in ).

The existence of such a partition into cliques can be used to characterize the line graphs: A graph is the line graph of some other graph or multigraph if and only if it is possible to find a collection of cliques in (allowing some of the cliques to be single vertices) that partition the edges of , such that each vertex of belongs to exactly two of the cliques. It is the line graph of a graph (rather than a multigraph) if this set of cliques satisfies the additional condition that no two vertices of are both in the same two cliques. Given such a family of cliques, the underlying graph for which is the line graph can be recovered by making one vertex in for each clique, and an edge in for each vertex in with its endpoints being the two cliques containing the vertex in. By the strong version of Whitney's isomorphism theorem, if the underlying graph has more than four vertices, there can be only one partition of this type.

For example, this characterization can be used to show that the following graph is not a line graph:

In this example, the edges going upward, to the left, and to the right from the central degree-four vertex do not have any cliques in common. Therefore, any partition of the graph's edges into cliques would have to have at least one clique for each of these three edges, and these three cliques would all intersect in that central vertex, violating the requirement that each vertex appear in exactly two cliques. Thus, the graph shown is not a line graph.

### Forbidden subgraphs

Another characterization of line graphs was proven in (and reported earlier without proof by ). He showed that there are nine minimal graphs that are not line graphs, such that any graph that is not a line graph has one of these nine graphs as an induced subgraph. That is, a graph is a line graph if and only if no subset of its vertices induces one of these nine graphs. In the example above, the four topmost vertices induce a claw (that is, a complete bipartite graph *K*), shown on the top left of the illustration of forbidden subgraphs. Therefore, by Beineke's characterization, this example cannot be a line graph. For graphs with minimum degree at least 5, only the six subgraphs in the left and right columns of the figure are needed in the characterization.

### Algorithms
 and described linear time algorithms for recognizing line graphs and reconstructing their original graphs. generalized these methods to directed graphs. described an efficient data structure for maintaining a dynamic graph, subject to vertex insertions and deletions, and maintaining a representation of the input as a line graph (when it exists) in time proportional to the number of changed edges at each step.

The algorithms of and are based on characterizations of line graphs involving odd triangles (triangles in the line graph with the property that there exists another vertex adjacent to an odd number of triangle vertices). However, the algorithm of uses only Whitney's isomorphism theorem. It is complicated by the need to recognize deletions that cause the remaining graph to become a line graph, but when specialized to the static recognition problem only insertions need to be performed, and the algorithm performs the following steps:
*Construct the input graph by adding vertices one at a time, at each step choosing a vertex to add that is adjacent to at least one previously-added vertex. While adding vertices to , maintain a graph for which *L* = *L*(*G*); if the algorithm ever fails to find an appropriate graph , then the input is not a line graph and the algorithm terminates.
*When adding a vertex to a graph *L*(*G*) for which has four or fewer vertices, it might be the case that the line graph representation is not unique. But in this case, the augmented graph is small enough that a representation of it as a line graph can be found by a brute force search in constant time.
*When adding a vertex to a larger graph that equals the line graph of another graph , let be the subgraph of formed by the edges that correspond to the neighbors of in. Check that has a vertex cover consisting of one vertex or two non-adjacent vertices. If there are two vertices in the cover, augment by adding an edge (corresponding to ) that connects these two vertices. If there is only one vertex in the cover, then add a new vertex to , adjacent to this vertex.
Each step either takes constant time, or involves finding a vertex cover of constant size within a graph whose size is proportional to the number of neighbors of. Thus, the total time for the whole algorithm is proportional to the sum of the numbers of neighbors of all vertices, which (by the handshaking lemma) is proportional to the number of input edges.

## Iterating the line graph operator
 consider the sequence of graphs
\(G, L(G), L(L(G)), L(L(L(G))), \dots.\\)
They show that, when is a finite connected graph, only four behaviors are possible for this sequence:
*If is a cycle graph then *L*(*G*) and each subsequent graph in this sequence are isomorphic to itself. These are the only connected graphs for which *L*(*G*) is isomorphic to. *If is a claw *K*, then *L*(*G*) and all subsequent graphs in the sequence are triangles.
*If is a path graph then each subsequent graph in the sequence is a shorter path until eventually the sequence terminates with an empty graph.
*In all remaining cases, the sizes of the graphs in this sequence eventually increase without bound.
If is not connected, this classification applies separately to each component of. For connected graphs that are not paths, all sufficiently high numbers of iteration of the line graph operation produce graphs that are Hamiltonian.

## Generalizations
### Medial graphs and convex polyhedra

When a planar graph has maximum vertex degree three, its line graph is planar, and every planar embedding of can be extended to an embedding of *L*(*G*). However, there exist planar graphs with higher degree whose line graphs are nonplanar. These include, for example, the 5-star *K*, the gem graph formed by adding two non-crossing diagonals within a regular pentagon, and all convex polyhedra with a vertex of degree four or more.

An alternative construction, the medial graph, coincides with the line graph for planar graphs with maximum degree three, but is always planar. It has the same vertices as the line graph, but potentially fewer edges: two vertices of the medial graph are adjacent if and only if the corresponding two edges are consecutive on some face of the planar embedding. The medial graph of the dual graph of a plane graph is the same as the medial graph of the original plane graph.

For regular polyhedra or simple polyhedra, the medial graph operation can be represented geometrically by the operation of cutting off each vertex of the polyhedron by a plane through the midpoints of all its incident edges. This operation is known variously as the second truncation, degenerate truncation, or rectification.

### Total graphs
The total graph *T*(*G*) of a graph has as its vertices the elements (vertices or edges) of , and has an edge between two elements whenever they are either incident or adjacent. The total graph may also be obtained by subdividing each edge of and then taking the square of the subdivided graph.

### Multigraphs
The concept of the line graph of may naturally be extended to the case where is a multigraph. In this case, the characterizations of these graphs can be simplified: the characterization in terms of clique partitions no longer needs to prevent two vertices from belonging to the same to cliques, and the characterization by forbidden graphs has seven forbidden graphs instead of nine.

However, for multigraphs, there are larger numbers of pairs of non-isomorphic graphs that have the same line graphs. For instance a complete bipartite graph *K* has the same line graph as the dipole graph and Shannon multigraph with the same number of edges. Nevertheless, analogues to Whitney's isomorphism theorem can still be derived in this case.

### Line digraphs

It is also possible to generalize line graphs to directed graphs. If is a directed graph, its **directed line graph** or **line digraph** has one vertex for each edge of. Two vertices representing directed edges from to and from to in are connected by an edge from to in the line digraph when *v* = *w*. That is, each edge in the line digraph of represents a length-two directed path in. The de Bruijn graphs may be formed by repeating this process of forming directed line graphs, starting from a complete directed graph.

### Weighted line graphs
In a line graph *L*(*G*), each vertex of degree in the original graph creates *k*(*k* − 1)/2 edges in the line graph. For many types of analysis this means high-degree nodes in are over-represented in the line graph *L*(*G*). For instance, consider a random walk on the vertices of the original graph. This will pass along some edge with some frequency. On the other hand, this edge is mapped to a unique vertex, say , in the line graph *L*(*G*). If we now perform the same type of random walk on the vertices of the line graph, the frequency with which is visited can be completely different from *f*. If our edge in was connected to nodes of degree *O*(*k*), it will be traversed *O*(*k* more frequently in the line graph *L*(*G*). Put another way, the Whitney graph isomorphism theorem guarantees that the line graph almost always encodes the topology of the original graph faithfully but it does not guarantee that dynamics on these two graphs have a simple relationship. One solution is to construct a weighted line graph, that is, a line graph with weighted edges. There are several natural ways to do this. For instance if edges and in the graph are incident at a vertex with degree , then in the line graph *L*(*G*) the edge connecting the two vertices and can be given weight 1/(*k* − 1). In this way every edge in (provided neither end is connected to a vertex of degree 1) will have strength 2 in the line graph *L*(*G*) corresponding to the two ends that the edge has in. It is straightforward to extend this definition of a weighted line graph to cases where the original graph was directed or even weighted. The principle in all cases is to ensure the line graph *L*(*G*) reflects the dynamics as well as the topology of the original graph. ### Line graphs of hypergraphs

The edges of a hypergraph may form an arbitrary family of sets, so the line graph of a hypergraph is the same as the intersection graph of the sets from the family.

### Disjointness graph
The **disjointness graph** of , denoted *D*(*G*), is constructed in the following way: for each edge in , make a vertex in *D*(*G*); for every two edges in that do *not* have a vertex in common, make an edge between their corresponding vertices in *D*(*G*). In other words, *D*(*G*) is the complement graph of *L*(*G*). A clique in *D*(*G*) corresponds to an independent set in *L*(*G*), and vice versa.
