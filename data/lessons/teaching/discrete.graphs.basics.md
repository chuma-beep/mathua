> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_gt-intro.html) by Oscar Levin, CC BY-SA 4.0

# Graph Theory Basics

##  Section 4.1 Definitions

######  Investigate!

Which (if any) of the graphs below are the same?

![Five vertices in two rows, with every top vertex connected to every bottom vertex](/diagrams/graphs/sec_gt-intro-2-2-1.svg)
![Five vertices in a pentagon, each connected to its two neighbours](/diagrams/graphs/sec_gt-intro-2-2-2.svg)
![Five vertices in a diamond with one in the middle](/diagrams/graphs/sec_gt-intro-2-2-3.svg)
![Five vertices in a pentagon, joined as a five-pointed star](/diagrams/graphs/sec_gt-intro-2-2-4.svg)
![Five vertices, with the centre vertex joined to all four outer vertices](/diagrams/graphs/sec_gt-intro-2-2-5.svg)

The graphs above are unlabeled. Usually we think of a graph as having a specific set of vertices. Which (if any) of the graphs below are the same?

![Six labeled vertices in two rows, joined like two crossing squares](/diagrams/graphs/sec_gt-intro-2-4-1.svg)
![Six labeled vertices in two rows, joined like two adjacent squares](/diagrams/graphs/sec_gt-intro-2-4-2.svg)
![Six labeled vertices in two rows, with the middle labels swapped](/diagrams/graphs/sec_gt-intro-2-4-3.svg)
![Six labeled vertices in two rows, using v labels](/diagrams/graphs/sec_gt-intro-2-4-4.svg)

Actually, all the graphs we have seen above are just _drawings_ of graphs. A graph is really an abstract mathematical object consisting of two sets \(V\) and \(E\) where \(E\) is a set of 2-element subsets of \(V\text{.}\) Are the graphs below the same or different?

Graph 1:

\(V = \{a, b, c, d, e\}\text{,}\)

\(E = \{\{a,b\}, \{a, c\}, \{a,d\}, \{a,e\}, \{b,c\}, \{d,e\}\}\text{.}\)

Graph 2:

\(V = \{v_1, v_2, v_3, v_4, v_5\}\text{,}\)

\(E = \{\{v_1, v_3\}, \{v_1, v_5\}, \{v_2, v_4\}, \{v_2, v_5\}, \{v_3, v_5\}, \{v_4, v_5\}\}\text{.}\)

Before we start studying graphs, we need to agree upon what a graph is. While we almost always think of graphs as pictures (dots connected by lines) this is fairly ambiguous. Do the lines need to be straight? Does it matter how long the lines are or how large the dots are? Can there be two lines connecting the same pair of dots? Can one line connect three dots?

The way we avoid ambiguities in mathematics is to provide concrete and rigorous _definitions_. Crafting good definitions is not easy, but it is incredibly important. The definition is the agreed upon starting point from which all truths in mathematics proceed. Is there a graph with no edges? We have to look at the definition to see if this is possible.

We want our definition to be precise and unambiguous, but it also must agree with our intuition for the objects we are studying. It needs to be useful: we _could_ define a graph to be a six legged mammal, but that would not let us solve any problems about bridges. Instead, here is the (now) standard definition of a graph.

### Graph Definition.

A graph is an ordered pair \(G = (V, E)\) consisting of a nonempty set \(V\) (called the vertices) and a set \(E\) (called the edges) of two-element subsets of \(V\text{.}\)

Strange. Nowhere in the definition is there talk of dots or lines. From the definition, a graph could be

$$
(\{a,b,c,d\}, \{\{a,b\}, \{a,c\}, \{b,c\}, \{b,d\}, \{c,d\}\})\text{.}
$$

Here we have a graph with four vertices (the letters \(a, b, c, d\)) and five edges (the pairs \(\{a,b\}, \{a,c\}, \{b,c\}, \{b,d\}, \{c,d\})\)).

Looking at sets and sets of 2-element sets is difficult to process. That is why we often draw a representation of these sets. We put a dot down for each vertex, and connect two dots with a line precisely when those two vertices are one of the 2-element subsets in our set of edges. Thus one way to draw the graph described above is this:

![Four vertices in a square, with edges on the border and one diagonal](/diagrams/graphs/sec_gt-intro-9.svg)

However we could also have drawn the graph differently. For example either of these:

![The same graph drawn with a different placement of the vertices](/diagrams/graphs/sec_gt-intro-11-1.svg)
![The same graph drawn as a row with two crossing curved edges](/diagrams/graphs/sec_gt-intro-11-2.svg)

We should be careful about what it means for two graphs to be “the same.” Actually, given our definition, this is easy: Are the vertex sets equal? Are the edge sets equal? We know what it means for sets to be equal, and graphs are nothing but a pair of two special sorts of sets.

###  Example 4.1.1.

Are the graphs below equal?

$$
G_1 = (\{a,b,c\}, \{\{a,b\}, \{b,c\}\}); \qquad G_2 = (\{a,b,c\}, \{\{a,c\}, \{c, b\}\})\text{.}
$$

Even if two graphs are not _equal_ , they might be _basically_ the same. The graphs in the previous example could be drawn like this:

![Two three-vertex graphs drawn as rows, showing the same path with different labels](/diagrams/graphs/sec_gt-intro-15.svg)

Graphs that are basically the same (but perhaps not equal) are called isomorphic. We will give a precise definition of this term after a quick example:

###  Example 4.1.2.

Consider the graphs:

$$
G_1 = (V_1, E_1) \text{ where }
V_1 = \{a, b, c\} \text{ and } E_1 = \{\{a,b\}, \{a,c\}, \{b,c\}\};
$$ $$
G_2 = (V_2, E_2) \text{ where }
V_2 = \{u,v,w\} \text{ and }E_2 = \{\{u,v\}, \{u,w\}, \{v,w\}\}.
$$

![A triangle with vertices a, b and c](/diagrams/graphs/sec_gt-intro-17-2-2-1.svg)
![A triangle with vertices u, v and w](/diagrams/graphs/sec_gt-intro-17-2-2-2.svg)

Are these graphs the same?

Intuitively, graphs are isomorphic if they are basically the same, or better yet, if they are the same except for the names of the vertices. To make the concept of renaming vertices precise, we give the following definitions:

### Isomorphic Graphs.

An isomorphism between two graphs \(G_1\) and \(G_2\) is a bijection \(f:V_1 \to V_2\) between the vertices of the graphs such that \(\{a,b\}\) is an edge in \(G_1\) if and only if \(\{f(a), f(b)\}\) is an edge in \(G_2\text{.}\)

Two graphs are isomorphic if there is an isomorphism between them. In this case we write \(G_1 \cong G_2\text{.}\)

An isomorphism is simply a function which renames the vertices. It must be a bijection so every vertex gets a new name. These newly named vertices must be connected by edges precisely when they were connected by edges with their old names.

###  Example 4.1.3.

Decide whether the graphs \(G_1 = (V_1, E_1)\) and \(G_2 = (V_2, E_2)\) are equal or isomorphic.

\(V_1 = \{a,b,c,d\}\text{,}\) \(E_1 = \{\{a,b\}, \{a,c\}, \{a,d\}, \{c,d\}\}\)

\(V_2 = \{a,b,c,d\}\text{,}\) \(E_2 = \{\{a,b\}, \{a,c\}, \{b,c\}, \{c,d\}\}\)

![The first four-vertex graph drawn as a diamond](/diagrams/graphs/sec_gt-intro-21-3-4-1.svg)
![The second four-vertex graph drawn as a diamond](/diagrams/graphs/sec_gt-intro-21-3-4-2.svg)

Sometimes we will talk about a graph with a special name (like \(K_n\) or the _Petersen graph_) or perhaps draw a graph without any labels. In this case we are really referring to _all_ graphs isomorphic to any copy of that particular graph. A collection of isomorphic graphs is often called an isomorphism class.  1

This is not unlike geometry, where we might have more than one copy of a particular triangle. There instead of _isomorphic_ we say _congruent_.

There are other relationships between graphs that we care about, other than equality and being isomorphic. For example, compare the following pair of graphs:

![A complete graph on six vertices drawn as a hexagon](/diagrams/graphs/sec_gt-intro-24-1.svg)
![A complete graph on four vertices drawn as a diamond](/diagrams/graphs/sec_gt-intro-24-2.svg)

These are definitely not isomorphic, but notice that the graph on the right looks like it might be part of the graph on the left, especially if we draw it like this:

![The six-vertex graph with a four-vertex subgraph highlighted](/diagrams/graphs/sec_gt-intro-26.svg)

We would like to say that the smaller graph is a _subgraph_ of the larger.

We should give a careful definition of this. In fact, there are two reasonable notions for what a subgraph should mean.

### Subgraphs.

We say that \(G' = (V', E')\) is a subgraph of \(G = (V, E)\text{,}\) and write \(G' \subseteq G\text{,}\) provided \(V' \subseteq V\) and \(E' \subseteq E\text{.}\)

We say that \(G' = (V', E')\) is an induced subgraph of \(G = (V, E)\) provided \(V' \subseteq V\) and every edge in \(E\) whose vertices are still in \(V'\) is also an edge in \(E'\text{.}\)

Notice that every induced subgraph is also an ordinary subgraph, but not conversely. Think of a subgraph as the result of deleting some vertices and edges from the larger graph. For the subgraph to be an induced subgraph, we can still delete vertices, but now we only delete those edges that included the deleted vertices.

###  Example 4.1.4.

Consider the graphs:

![The large graph G1 on six vertices](/diagrams/graphs/sec_gt-intro-31-1-2-1.svg)
![The induced subgraph G2](/diagrams/graphs/sec_gt-intro-31-1-2-2.svg)
![A subgraph G3 that deletes an edge as well as vertices](/diagrams/graphs/sec_gt-intro-31-1-2-3.svg)
![A graph G4 that is not a subgraph of G1](/diagrams/graphs/sec_gt-intro-31-1-2-4.svg)

Here both \(G_2\) and \(G_3\) are subgraphs of \(G_1\text{.}\) But only \(G_2\) is an _induced_ subgraph. Every edge in \(G_1\) that connects vertices in \(G_2\) is also an edge in \(G_2\text{.}\) In \(G_3\text{,}\) the edge \(\{a,b\}\) is in \(E_1\) but not \(E_3\text{,}\) even though vertices \(a\) and \(b\) are in \(V_3\text{.}\)

The graph \(G_4\) is NOT a subgraph of \(G_1\text{,}\) even though it looks like all we did is remove vertex \(e\text{.}\) The reason is that in \(E_4\) we have the edge \(\{c,f\}\) but this is not an element of \(E_1\text{,}\) so we don’t have the required \(E_4 \subseteq E_1\text{.}\)

Back to some basic graph theory definitions. Notice that all the graphs we have drawn above have the property that no pair of vertices is connected more than once, and no vertex is connected to itself. Graphs like these are sometimes called simple, although we will just call them _graphs_. This is because our definition for a graph says that the edges form a set of 2-element subsets of the vertices. Remember that it doesn’t make sense to say a set contains an element more than once. So no pair of vertices can be connected by an edge more than once. Also, since each edge must be a set containing two vertices, we cannot have a single vertex connected to itself by an edge.

That said, there are times we want to consider double (or more) edges and single edge loops. For example, the “graph” we drew for the Bridges of Königsberg problem had double edges because there really are two bridges connecting a particular island to the near shore. We will call these objects multigraphs. This is a good name: a _multiset_ is a set in which we are allowed to include a single element multiple times.

The graphs above are also connected: you can get from any vertex to any other vertex by following some path of edges. A graph that is not connected can be thought of as two separate graphs drawn close together. For example, the following graph is NOT connected because there is no path from \(a\) to \(b\text{:}\)

![A graph made of two separate diamond-shaped components, labeled a on the left and b on the right](/diagrams/graphs/ex-gt-non-connected.svg)

Vertices in a graph do not always have edges between them. If we add all possible edges, then the resulting graph is called complete. That is, a graph is complete if every pair of vertices is connected by an edge. Since a graph is determined completely by which vertices are adjacent to which other vertices, there is only one complete graph with a given number of vertices. We give these a special name: \(K_n\) is the complete graph on \(n\) vertices.

Each vertex in \(K_n\) is adjacent to \(n-1\) other vertices. We call the number of edges emanating from a given vertex the degree of that vertex. So every vertex in \(K_n\) has degree \(n-1\text{.}\) How many edges does \(K_n\) have? One might think the answer should be \(n(n-1)\text{,}\) since we count \(n-1\) edges \(n\) times (once for each vertex). However, each edge is incident to 2 vertices, so we counted every edge exactly twice. Thus there are \(n(n-1)/2\) edges in \(K_n\text{.}\) Alternatively, we can say there are \({n \choose 2}\) edges, since to draw an edge we must choose 2 of the \(n\) vertices.

In general, if we know the degrees of all the vertices in a graph, we can find the number of edges. The sum of the degrees of all vertices will always be _twice_ the number of edges, since each edge adds to the degree of two vertices. Notice this means that the sum of the degrees of all vertices in any graph must be even!

This is our first example of a general result about all graphs. It seems innocent enough, but we will use it to prove all sorts of other statements. So let’s give it a name and state it formally.

###  Lemma 4.1.5. Handshake Lemma.

In any graph, the sum of the degrees of vertices in the graph is always twice the number of edges.

The handshake lemma 2

A _lemma_ is a mathematical statement that is primarily of importance in that it is used to establish other results.

is sometimes called the _degree sum formula_ , and can be written symbolically as

$$
\sum_{v\in V} d(v) = 2e\text{.}
$$

Here we are using the notation \(d(v)\) for the degree of the vertex \(v\text{.}\)

One use for the lemma is to actually find the number of edges in a graph. To do this, you must be given the degree sequence for the graph (or be able to find it from other information). This is a list of every degree of every vertex in the graph, generally written in non-increasing order.

###  Example 4.1.6.

How many vertices and edges must a graph have if its degree sequence is

$$
(4, 4, 3, 3, 3, 2, 1)\text{?}
$$

The handshake lemma also tells us what is not possible.

###  Example 4.1.7.

At a recent math seminar, 9 mathematicians greeted each other by shaking hands. Is it possible that each mathematician shook hands with exactly 7 people at the seminar?

We can generalize the previous example to get the following proposition. 3

A proposition is a general statement in mathematics, similar to a theorem, although generally of lesser importance.

###  Proposition 4.1.8.

In any graph, the number of vertices with odd degree must be even.

### Proof.

Suppose there were a graph with an odd number of vertices with odd degree. Then the sum of the degrees in the graph would be odd, which is impossible, by the handshake lemma.

We will consider further applications of the handshake lemma in the exercises.

One final definition: we say a graph is bipartite if the vertices can be divided into two sets, \(A\) and \(B\text{,}\) with no two vertices in \(A\) adjacent and no two vertices in \(B\) adjacent. The vertices in \(A\) can be adjacent to some or all of the vertices in \(B\text{.}\) If each vertex in \(A\) is adjacent to all the vertices in \(B\text{,}\) then the graph is a complete bipartite graph, and gets a special name: \(K_{m,n}\text{,}\) where \(|A| = m\) and \(|B| = n\text{.}\) The graph in the houses and utilities puzzle is \(K_{3,3}\text{.}\)

### Named Graphs.

Some graphs are used more than others, and get special names.

\(K_n\)

The complete graph on \(n\) vertices.

![The complete graph K5](/diagrams/graphs/sec_gt-intro-50-11-1.svg)

\(K_{m,n}\)

The complete bipartite graph with sets of \(m\) and \(n\) vertices.

![The complete bipartite graph K2,3](/diagrams/graphs/sec_gt-intro-50-11-2.svg)

\(C_n\)

The cycle on \(n\) vertices, just one big loop.

![The cycle C6](/diagrams/graphs/sec_gt-intro-50-11-3.svg)

\(P_n\)

The path on \(n+1\) vertices (so \(n\) edges), just one long path.

![The path P5 on six vertices](/diagrams/graphs/sec_gt-intro-50-11-4.svg)

### Graph Theory Definitions.

There are a lot of definitions to keep track of in graph theory. Here is a glossary of the terms we have already used and will soon encounter.

Graph

A collection of vertices, some of which are connected by edges. More precisely, a pair of sets \(V\) and \(E\) where \(V\) is a set of vertices and \(E\) is a set of 2-element subsets of \(V\text{.}\)

Adjacent

Two vertices are adjacent if they are connected by an edge. Two edges are adjacent if they share a vertex.

Bipartite graph

A graph for which it is possible to divide the vertices into two disjoint sets such that there are no edges between any two vertices in the same set.

Complete bipartite graph

A bipartite graph for which every vertex in the first set is adjacent to every vertex in the second set.

Complete graph

A graph in which every pair of vertices is adjacent.

Connected

A graph is connected if there is a path from any vertex to any other vertex.

Chromatic number

The minimum number of colors required in a proper vertex coloring of the graph.

Cycle

A path (see below) that starts and stops at the same vertex, but contains no other repeated vertices.

Degree of a vertex

The number of edges incident to a vertex.

Euler path

A walk which uses each edge exactly once.

Euler circuit

An Euler path which starts and stops at the same vertex.

Multigraph

A multigraph is just like a graph but can contain multiple edges between two vertices as well as single edge loops (that is an edge from a vertex to itself).

Path

A path is a walk that doesn’t repeat any vertices (or edges) except perhaps the first and last. If a path starts and ends at the same vertex, it is called a cycle.

Planar

A graph which can be drawn (in the plane) without any edges crossing.

Subgraph

We say that \(H\) is a subgraph of \(G\) if every vertex and edge of \(H\) is also a vertex or edge of \(G\text{.}\) We say \(H\) is an induced subgraph of \(G\) if every vertex of \(H\) is a vertex of \(G\) and each pair of vertices in \(H\) are adjacent in \(H\) if and only if they are adjacent in \(G\text{.}\)

Tree

A connected graph with no cycles. (If we remove the requirement that the graph is connected, the graph is called a forest.) The vertices in a tree with degree 1 are called leaves.

Vertex coloring

An assignment of colors to each of the vertices of a graph. A vertex coloring is proper if adjacent vertices are always colored differently.

Walk

A sequence of vertices such that consecutive vertices (in the sequence) are adjacent (in the graph). A walk in which no edge is repeated is called a trail, and a trail in which no vertex is repeated (except possibly the first and last) is called a path.
