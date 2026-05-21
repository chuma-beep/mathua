> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Path_%28graph_theory%29) — CC BY-SA 4.0

# Paths, cycles, and connectivity

In graph theory, a **path** in a graph is a finite or infinite sequence of edges which joins a sequence of vertices which, by most definitions, are all distinct (and since the vertices are distinct, so are the edges). A **directed path** (sometimes called **dipath**) in a directed graph is a finite or infinite sequence of edges which joins a sequence of distinct vertices, but with the added restriction that the edges be all directed in the same direction.

Paths are fundamental concepts of graph theory, described in the introductory sections of most graph theory texts. See e.g. , or. cover more advanced algorithmic topics concerning paths in graphs.

## Definitions
### Walk, trail, and path

* A **walk** is a finite or infinite sequence of edges which joins a sequence of vertices.
Let *G* = (*V*, *E*, *Φ*) be a graph. A finite walk is a sequence of edges (*e*\(_{1}\), *e*\(_{2}\),... , *e*\(_{*n* − 1}\)) for which there is a sequence of vertices (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) such that *Φ*(*e*\(_{*i*}\)) = {*v*\(_{*i*}\), *v*\(_{*i* + 1}\)} for *i* = 1, 2,... , *n* − 1. (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) is the *vertex sequence* of the walk. The walk is *closed* if *v*\(_{1}\) = *v*\(_{*n*}\), and it is *open* otherwise. An infinite walk is a sequence of edges of the same type described here, but with no first or last vertex, and a semi-infinite walk (or ray) has a first vertex but no last vertex.
* A **trail** is a walk in which all edges are distinct.
* A **path** is a trail in which all vertices (and therefore also all edges) are distinct.

If *w* = (*e*\(_{1}\), *e*\(_{2}\),... , *e*\(_{*n* − 1}\)) is a finite walk with vertex sequence (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) then *w* is said to be a *walk from* *v*\(_{1}\) *to* *v*\(_{*n*}\). Similarly for a trail or a path. If there is a finite walk between two *distinct* vertices then there is also a finite trail and a finite path between them.

Some authors do not require that all vertices of a path be distinct and instead use the term **simple path** to refer to such a path where all vertices are distinct.

A weighted graph associates a value (*weight*) with every edge in the graph. The *weight of a walk* (or trail or path) in a weighted graph is the sum of the weights of the traversed edges. Sometimes the words *cost* or *length* are used instead of weight.

### Directed walk, directed trail, and directed path
* A **directed walk** is a finite or infinite sequence of edges directed in the same direction which joins a sequence of vertices.
Let *G* = (*V*, *E*, *Φ*) be a directed graph. A finite directed walk is a sequence of edges (*e*\(_{1}\), *e*\(_{2}\),... , *e*\(_{*n* − 1}\)) for which there is a sequence of vertices (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) such that *Φ*(*e*\(_{*i*}\)) = (*v*\(_{*i*}\), *v*\(_{*i* + 1}\)) for *i* = 1, 2,... , *n* − 1. (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) is the *vertex sequence* of the directed walk. The directed walk is *closed* if *v*\(_{1}\) = *v*\(_{*n*}\), and it is *open* otherwise. An infinite directed walk is a sequence of edges of the same type described here, but with no first or last vertex, and a semi-infinite directed walk (or ray) has a first vertex but no last vertex.
* A **directed trail** is a directed walk in which all edges are distinct.
* A **directed path** is a directed trail in which all vertices are distinct.

If *w* = (*e*\(_{1}\), *e*\(_{2}\),... , *e*\(_{*n* − 1}\)) is a finite directed walk with vertex sequence (*v*\(_{1}\), *v*\(_{2}\),... , *v*\(_{*n*}\)) then *w* is said to be a *walk from* *v*\(_{1}\) *to* *v*\(_{*n*}\). Similarly for a directed trail or a path. If there is a finite directed walk between two *distinct* vertices then there is also a finite directed trail and a finite directed path between them.

A "simple directed path" is a path where all vertices are distinct.

A weighted directed graph associates a value (*weight*) with every edge in the directed graph. The *weight of a directed walk* (or trail or path) in a weighted directed graph is the sum of the weights of the traversed edges. Sometimes the words *cost* or *length* are used instead of weight.

## Examples
* A graph is connected if there are paths containing each pair of vertices.
* A directed graph is strongly connected if there are oppositely oriented directed paths containing each pair of vertices.
* A path such that no graph edges connect two nonconsecutive path vertices is called an induced path.
* A path that includes every vertex of the graph without repeats is known as a Hamiltonian path.
* Two paths are *vertex-independent* (alternatively, *internally disjoint* or *internally vertex-disjoint*) if they do not have any internal vertex or edge in common. Similarly, two paths are *edge-independent* (or *edge-disjoint*) if they do not have any edge in common. Two internally disjoint paths are edge-disjoint, but the converse is not necessarily true.
* The distance between two vertices in a graph is the length of a shortest path between them, if one exists, and otherwise the distance is infinity.
* The diameter of a connected graph is the largest distance (defined above) between pairs of vertices of the graph.

## Finding paths
Several algorithms exist to find shortest and longest paths in graphs, with the important distinction that the former problem is computationally much easier than the latter.

Dijkstra's algorithm produces a list of shortest paths from a source vertex to every other vertex in directed and undirected graphs with non-negative edge weights (or no edge weights), whilst the Bellman–Ford algorithm can be applied to directed graphs with negative edge weights. The Floyd–Warshall algorithm can be used to find the shortest paths between all pairs of vertices in weighted directed graphs.

## The path partition problem
The **k-path partition** problem is the problem of partitioning a given graph to a smallest collection of vertex-disjoint paths of length at most *k*.
