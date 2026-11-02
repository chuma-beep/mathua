# Matroids

## Independence and Rank

A matroid is a set with a notion of independence that satisfies the hereditary and exchange axioms. The rank of a subset is the size of the largest independent subset it contains. Graphic matroids take the edges of a graph as the ground set and forests as the independent sets.

Matroids abstract the parts of linear algebra and graph theory that behave alike. The same rank language describes linearly independent vectors and acyclic edge sets, which is why the theory appears in both subjects.

## Bases and Circuits

A basis is a maximal independent set, and every basis of a matroid has the same size, called the rank. A circuit is a minimal dependent set. In a graphic matroid, bases are spanning forests and circuits are simple cycles.

The uniform matroid U_{r,n} declares every set of at most r elements independent. Its bases are the r-element subsets of an n-element ground set, so their number is the binomial coefficient C(n,r). Its circuits all have size r+1.

## The Graphic Matroid

For a connected graph on n vertices, the graphic matroid has rank n-1, the number of edges in a spanning tree. Every spanning tree is a basis, and every cycle is a circuit. Disconnected graphs have rank n minus the number of components.

This dictionary lets graph algorithms be read as matroid algorithms. Kruskal's algorithm, for example, grows a maximum independent set by adding edges that do not create a circuit.
