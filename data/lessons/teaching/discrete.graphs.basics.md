> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed](https://discrete.openmathbooks.org/dmoi3/sec_gt-intro.html) by Oscar Levin — CC BY-SA 4.0

# Graph Theory Basics

## Section 4.1 Definitions

###### Investigate!

### Graph Definition.

### Example 4.1.1.

### Example 4.1.2.

### Isomorphic Graphs.

### Example 4.1.3.

### Subgraphs.

### Example 4.1.4.

### Lemma 4.1.5. Handshake Lemma.

### Example 4.1.6.

### Example 4.1.7.

### Proposition 4.1.8.

### Proof.

### Named Graphs.

### Graph Theory Definitions.

### Exercises  Exercises

#### 1.

#### 2.

#### 3.

#### 4.

#### 5.

1. Let \(f:G_1 \rightarrow G_2\) be a function that takes the vertices of Graph 1 to vertices of Graph 2. The function is given by the following table:ð



\(x\)
\(a\)
\(b\)
\(c\)
\(d\)
\(e\)
\(f\)
\(g\)


\(f(x)\)
\(v_4\)
\(v_5\)
\(v_1\)
\(v_6\)
\(v_2\)
\(v_3\)
\(v_7\)


Does \(f\) define an isomorphism between Graph 1 and Graph 2?ð

ð
| \(x\) | \(a\) | \(b\) | \(c\) | \(d\) | \(e\) | \(f\) | \(g\) |
| \(f(x)\) | \(v_4\) | \(v_5\) | \(v_1\) | \(v_6\) | \(v_2\) | \(v_3\) | \(v_7\) |

1. Define a new function \(g\) (with \(g \ne f\)) that defines an isomorphism between Graph 1 and Graph 2.ð

ð
1. Is the graph pictured below isomorphic to Graph 1 and Graph 2? Explain.ð


ð
#### 6.

#### 7.

#### 8.

#### 9.

1. Two different trees with the same number of vertices and the same number of edges. A tree is a connected graph with no cycles.ð

ð
1. Two different graphs with 8 vertices all of degree 2.ð

ð
1. Two different graphs with 5 vertices all of degree 4.ð

ð
1. Two different graphs with 5 vertices all of degree 3.ð

ð
1. For example:ð





ð
1. This is not possible if we require the graphs to be connected. If not, we could take \(C_8\) as one graph and two copies of \(C_4\) as the other.ð

ð
1. Not possible. If you have a graph with 5 vertices all of degree 4, then every vertex must be adjacent to every other vertex. This is the graph \(K_5\text{.}\)ð

ð
1. This is not possible. In fact, there is not even one graph with this property (such a graph would have \(5\cdot 3/2 = 7.5\) edges).ð

ð
#### 10.

1. Any subgraph of a complete graph is also complete.ð

ð
1. Any induced subgraph of a complete graph is also complete.ð

ð
1. Any subgraph of a bipartite graph is bipartite.ð

ð
1. Any subgraph of a tree is a tree.ð

ð
1. False.ð

ð
1. True.ð

ð
1. True.ð

ð
1. False.ð

ð
#### 11.

#### 12.

1. Let \(G\) be the graph with \(V = \{a,b,c,d,e,f\}\) and \(E = \{\{a,b\}, \{a,e\},\{b, c\}, \{b,e\}, \{c,d\}, \{c, f\}, \{d, f\}, \{e,f\}\}\text{.}\)  Find \(N(a)\text{,}\) \(N[a]\text{,}\) \(N(c)\text{,}\) and \(N[c]\text{.}\)ð

ð
1. What is the largest and smallest possible values for \(|N(v)|\) and \(|N[v]|\) for the graph in part (a)?  Explain.ð

ð
1. Give an example of a graph \(G = (V, E)\) (probably different than the one above) for which \(N[v] = V\) for some vertex \(v \in V\text{.}\)  Is there a graph for which \(N[v] = V\) for all \(v \in V\text{?}\)  Explain.ð

ð
1. Give an example of a graph \(G = (V,E)\) for which \(N(v) = \emptyset\) for some \(v \in V\text{.}\)  Is there an example of such a graph for which \(N[u] = V\) for some other \(u \in V\) as well?  Explain.ð

ð
1. Describe in words what \(N(v)\) and \(N[v]\) mean in general.ð

ð
#### 13.

1. The set \(V = \{1,2, \ldots, 9\}\) and the relationship \(x \sim y\) when \(x-y\) is a non-zero multiple of 3.ð

ð
1. The set \(V = \{1,2, \ldots, 9\}\) and the relationship \(x \sim y\) when \(y\) is a multiple of \(x\text{.}\)ð

ð
1. The set \(V = \{1,2,\ldots, 9\}\) and the relationship \(x \sim y\) when \(0 \lt |x-y| \lt 3\text{.}\)ð

ð
#### 14.

1. How many edges must the graph have to guarantee at least one vertex has degree two or more?  Prove your answer.ð

ð
1. How many edges must the graph have to guarantee all vertices have degree two or more?  Prove your answer.ð

ð
#### 15.

#### 16.