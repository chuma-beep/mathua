> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Metric_space) — CC BY-SA 4.0

# Metric spaces

In mathematics, a **metric space** is a set together with a notion of distance between its points. The distance is measured by a function called a **metric** or **distance function**. Metric spaces are a general setting for studying many of the concepts of mathematical analysis and geometry.

The most familiar example of a metric space is 3-dimensional Euclidean space with its usual notion of distance. Other well-known examples are a sphere equipped with the angular distance and the hyperbolic plane. A metric may correspond to a metaphorical, rather than physical, notion of distance. For example, the set of 100-character Unicode strings can be equipped with the Hamming distance, which measures the number of characters that need to be changed to get from one string to another.

Metric spaces appear in many different branches of mathematics. For example, Riemannian manifolds, normed vector spaces, and graphs may be viewed as metric spaces. In abstract algebra, the field of *p*-adic numbers is the completion of the field of rational numbers with respect to a certain metric. Metric spaces are also studied in their own right in **metric geometry** and **analysis on metric spaces**.

Many notions of analysis, including balls, completeness, as well as uniform, Lipschitz, and Hölder continuity can be defined for metric spaces. Other notions, such as continuity, compactness, and open and closed sets can be defined for metric spaces, but also in the even more general setting of topological spaces.

## Definition and illustration
### Motivation

To see the utility of different notions of distance, consider the surface of the Earth as a set of points. We can measure the distance between two such points by the length of the shortest path along the surface, "as the crow flies"; this is particularly useful for shipping and aviation. We can also measure the straight-line distance between two points through the Earth's interior; this notion is, for example, natural in seismology; since it roughly corresponds to the length of time it takes for seismic waves to travel between those two points.

The notion of distance encoded by the metric space axioms has relatively few requirements. This generality gives metric spaces a lot of flexibility. At the same time, the notion is strong enough to encode many intuitive facts about what distance means. This means that general results about metric spaces can be applied in many different contexts.

Like many fundamental mathematical concepts, the metric on a metric space can be interpreted in many different ways. A particular metric may not be best thought of as measuring physical distance, but, instead, as the cost of changing from one state to another (as with Wasserstein metrics on spaces of measures) or the degree of difference between two objects (for example, the Hamming distance between two strings of characters, or the Gromov–Hausdorff distance between metric spaces themselves).

### Definition
Formally, a **metric space** is an ordered pair (*M*, *d*) where is a set (whose elements are called points), and is a **metric** on , i.e., a function
\[
d\,\colon M \times M \to \mathbb{R}
\]
satisfying the following axioms for all points \(x,y,z \in M\):
# The distance from a point to itself is zero:
\[
d(x, x) = 0
\]

# (Positivity) The distance between two distinct points is always positive:
\[
\text{If }x \neq y\text{, then }d(x, y)>0
\]

# (Symmetry) The distance from to is always the same as the distance from to :
\[
d(x, y) = d(y, x)
\]

# The triangle inequality holds:
\[
d(x, z) \leq d(x, y) + d(y, z)
\]
This is a natural property of both physical and metaphorical notions of distance: you can arrive at from by taking a detour through , but this will not make your journey any shorter than the direct path.

It is traditional to write as an abbreviation for (*M*, *d*) when it is understood what is.

### Simple examples
#### The real numbers
The real numbers with the distance function \(d(x,y) = | y - x |\) given by the absolute difference form a metric space. Many properties of metric spaces and functions between them are generalizations of concepts in real analysis and coincide with those concepts when applied to the real line.

#### Metrics on Euclidean spaces

The Euclidean plane \(\R^2\) can be equipped with many different metrics. The Euclidean distance familiar from school mathematics can be defined by

\[
d_2((x_1,y_1),(x_2,y_2))=\sqrt{(x_2-x_1)^2+(y_2-y_1)^2}.
\]

The *taxicab* or *Manhattan* distance is defined by

\[
d_1((x_1,y_1),(x_2,y_2))=|x_2-x_1|+|y_2-y_1|
\]

and can be thought of as the distance you need to travel along horizontal and vertical lines to get from one point to the other, as illustrated at the top of the article.

The *maximum*, \(L^\infty\), or *Chebyshev distance* is defined by

\[
d_\infty((x_1,y_1),(x_2,y_2))=\max\{|x_2-x_1|,|y_2-y_1|\}.
\]

This distance does not have an easy explanation in terms of paths in the plane, but it still satisfies the metric space axioms. It can be thought of similarly to the number of moves a king would have to make on a chess board to travel from one point to another on the given space.

In fact, these three distances, while they have distinct properties, are similar in some ways. Informally, points that are close in one are close in the others, too. This observation can be quantified with the formula

\[
d_\infty(p,q) \leq d_2(p,q) \leq d_1(p,q) \leq 2d_\infty(p,q),
\]

which holds for every pair of points \(p, q \in \R^2\). Moreover, there exists a continuous family of distances \(d_\alpha\) for \(\alpha \in (1,\infty)\) that interpolates between these three distances. It is defined by
\[
d_\alpha((x_1,y_1),(x_2,y_2))=\left(\vert x_2-x_1\vert^\alpha+\vert y_2-y_1\vert^\alpha\right)^{1/\alpha}.
\]

A radically different distance can be defined by setting

\[
d(p,q)=\begin{cases}0, & \text{if }p=q, \\ 1, & \text{otherwise.}\end{cases}
\]

In this *discrete metric*, all distinct points are 1 unit apart: none of them are close to each other, and none of them are very far away from each other either. Intuitively, the discrete metric no longer remembers that the set is a plane, but treats it just as an undifferentiated set of points.

All of these metrics can be easily extended to make sense on \(\R^n\) as well as \(\R^2\).

#### Subspaces
Given a metric space (*M*, *d*) and a subset \(A \subseteq M\), we can consider to be a metric space by measuring distances the same way we would in. Formally, the *induced metric* on is a function \(d_A:A \times A \to \R\) defined by

\[
d_A(x,y)=d(x,y).
\]

For example, if we take the two-dimensional sphere S\(^{2}\) as a subset of \(\R^3\), the Euclidean metric on \(\R^3\) induces the straight-line metric on S\(^{2}\) described above. Two more useful examples are the open interval and the closed interval thought of as subspaces of the real line.

## History
Arthur Cayley, in his article "On Distance", extended metric concepts beyond Euclidean geometry into domains bounded by a conic in a projective space. His distance was given by logarithm of a cross ratio. Any projectivity leaving the conic stable also leaves the cross ratio constant, so isometries are implicit. This method provides models for elliptic geometry and hyperbolic geometry, and Felix Klein, in several publications, established the field of non-euclidean geometry through the use of the Cayley-Klein metric.

The idea of an abstract space with metric properties was addressed in 1906 by René Maurice Fréchet and the term *metric space* was coined by Felix Hausdorff in 1914.

Fréchet's work laid the foundation for understanding convergence, continuity, and other key concepts in non-geometric spaces. This allowed mathematicians to study functions and sequences in a broader and more flexible way. This was important for the growing field of functional analysis. Mathematicians like Hausdorff and Stefan Banach further refined and expanded the framework of metric spaces. Hausdorff introduced topological spaces as a generalization of metric spaces. Banach's work in functional analysis heavily relied on the metric structure. Over time, metric spaces became a central part of modern mathematics. They have influenced various fields including topology, geometry, and applied mathematics. Metric spaces continue to play a crucial role in the study of abstract mathematical concepts.

## Basic notions
A distance function is enough to define notions of closeness and convergence that were first developed in real analysis. Properties that depend on the structure of a metric space are referred to as *metric properties*. Every metric space is also a topological space, and some metric properties can also be rephrased without reference to distance in the language of topology; that is, they are really topological properties.

### The topology of a metric space
For any point in a metric space and any real number *r* > 0, the *open ball* of radius around is defined to be the set of points that are strictly less than distance from :

\[
B_r(x)=\{y \in M : d(x,y) < r\}.
\]

This is a natural way to define a set of points that are relatively close to. Therefore, a set \(N \subseteq M\) is a *neighborhood* of (informally, it contains all points "close enough" to ) if it contains an open ball of radius around for some *r* > 0.

An *open set* is a set which is a neighborhood of all its points. It follows that the open balls form a base for a topology on. In other words, the open sets of are exactly the unions of open balls. As in any topology, closed sets are the complements of open sets. Sets may be both open and closed as well as neither open nor closed.

This topology does not carry all the information about the metric space. For example, the distances *d*\(_{1}\), *d*\(_{2}\), and *d*\(_{∞}\) defined above all induce the same topology on \(\R^2\), although they behave differently in many respects. Similarly, \(\R\) with the Euclidean metric and its subspace the interval with the induced metric are homeomorphic but have very different metric properties.

Conversely, not every topological space can be given a metric. Topological spaces which are compatible with a metric are called *metrizable* and are particularly well-behaved in many ways: in particular, they are paracompact Hausdorff spaces (hence normal) and first-countable. form a neighborhood basis for that point.The Nagata–Smirnov metrization theorem gives a characterization of metrizability in terms of other topological properties, without reference to metrics.

### Convergence
Convergence of sequences in Euclidean space is defined as follows:
A sequence (*x\(_{n}\)*) converges to a point if for every ε > 0 there is an integer such that for all *n* > *N*, *d*(*x\(_{n}\)*, *x*) < ε.
Convergence of sequences in a topological space is defined as follows:
A sequence (*x\(_{n}\)*) converges to a point if for every open set containing there is an integer such that for all *n* > *N*, \(x_n \in U\).
In metric spaces, both of these definitions make sense and they are equivalent. This is a general pattern for topological properties of metric spaces: while they can be defined in a purely topological way, there is often a way that uses the metric which is easier to state or more familiar from real analysis.

### Completeness

Informally, a metric space is *complete* if it has no "missing points": every sequence that looks like it should converge to something actually converges.

To make this precise: a sequence (*x\(_{n}\)*) in a metric space is *Cauchy* if for every ε > 0 there is an integer such that for all *m*, *n* > *N*, *d*(*x\(_{m}\)*, *x\(_{n}\)*) < ε. By the triangle inequality, any convergent sequence is Cauchy: if and are both less than ε away from the limit, then they are less than 2ε away from each other. If the converse is true—every Cauchy sequence in converges—then is complete.

Euclidean spaces are complete, as is \(\R^2\) with the other metrics described above. Two examples of spaces which are not complete are and the rationals, each with the metric induced from \(\R\). One can think of as "missing" its endpoints 0 and 1. The rationals are missing all the irrationals; since any irrational has a sequence of rationals converging to it in \(\R\) (for example, its successive decimal approximations). These examples show that completeness is *not* a topological property; since \(\R\) is complete but the homeomorphic space is not.

This notion of "missing points" can be made precise. In fact, every metric space has a unique *completion*, which is a complete space that contains the given space as a dense subset. For example, is the completion of , and the real numbers are the completion of the rationals.

Since complete spaces are generally easier to work with, completions are important throughout mathematics. For example, in abstract algebra, the *p*-adic numbers are defined as the completion of the rationals under a different metric. Completion is particularly common as a tool in functional analysis. Often one has a set of nice functions and a way of measuring distances between them. Taking the completion of this metric space gives a new set of functions which may be less nice, but nevertheless useful because they behave similarly to the original nice functions in important ways. For example, weak solutions to differential equations typically live in a completion (a Sobolev space) rather than the original space of nice functions for which the differential equation actually makes sense.

### Bounded and totally bounded spaces

A metric space is *bounded* if there is an such that no pair of points in is more than distance apart. The least such is called the diameter of. The space is called *precompact* or *totally bounded* if for every *r* > 0 there is a finite cover of by open balls of radius. Every totally bounded space is bounded. To see this, start with a finite cover by -balls for some arbitrary. Since the subset of consisting of the centers of these balls is finite, it has finite diameter, say. By the triangle inequality, the diameter of the whole space is at most *D* + 2*r*. The converse does not hold: an example of a metric space that is bounded but not totally bounded is \(\R^2\) (or any other infinite set) with the discrete metric.

### Compactness

Compactness is a topological property which generalizes the properties of a closed and bounded subset of Euclidean space. There are several equivalent definitions of compactness in metric spaces:
# A metric space is compact if every open cover has a finite subcover (the usual topological definition).
# A metric space is compact if every sequence has a convergent subsequence. (For general topological spaces this is called sequential compactness and is not equivalent to compactness.)
# A metric space is compact if it is complete and totally bounded. (This definition is written in terms of metric properties and does not make sense for a general topological space, but it is nevertheless topologically invariant since it is equivalent to compactness.)
One example of a compact space is the closed interval. Compactness is important for similar reasons to completeness: it makes it easy to find limits. Another important tool is Lebesgue's number lemma, which shows that for any open cover of a compact space, every point is relatively deep inside one of the sets of the cover.

## Functions between metric spaces

Unlike in the case of topological spaces or algebraic structures such as groups or rings, there is no single "right" type of structure-preserving function between metric spaces. Instead, one works with different types of functions depending on one's goals. Throughout this section, suppose that \((M_1,d_1)\) and \((M_2,d_2)\) are two metric spaces. The words "function" and "map" are used interchangeably.

### Isometries

One interpretation of a "structure-preserving" map is one that fully preserves the distance function:
A function \(f:M_1 \to M_2\) is *distance-preserving* if for every pair of points and in *M*\(_{1}\),
\[
d_2(f(x),f(y))=d_1(x,y).
\]

It follows from the metric space axioms that a distance-preserving function is injective. A bijective distance-preserving function is called an *isometry*. One perhaps non-obvious example of an isometry between spaces described in this article is the map \(f:(\R^2,d_1) \to (\R^2,d_\infty)\) defined by

\[
f(x,y)=(x+y,x-y).
\]

If there is an isometry between the spaces *M*\(_{1}\) and *M*\(_{2}\), they are said to be *isometric*. Metric spaces that are isometric are essentially identical.

### Continuous maps

On the other end of the spectrum, one can forget entirely about the metric structure and study continuous maps, which only preserve topological structure. There are several equivalent definitions of continuity for metric spaces. The most important are:
* **Topological definition.** A function \(f\,\colon M_1\to M_2\) is continuous if for every open set in *M*\(_{2}\), the preimage \(f^{-1}(U)\) is open.
* **Sequential continuity.** A function \(f\,\colon M_1\to M_2\) is continuous if whenever a sequence (*x\(_{n}\)*) converges to a point in *M*\(_{1}\), the sequence \(f(x_1),f(x_2),\ldots\) converges to the point *f*(*x*) in *M*\(_{2}\).
(These first two definitions are *not* equivalent for all topological spaces.)
* **ε–δ definition.** A function \(f\,\colon M_1\to M_2\) is continuous if for every point in *M*\(_{1}\) and every ε > 0 there exists δ > 0 such that for all in *M*\(_{1}\) we have
\[
d_1(x,y) < \delta \implies d_2(f(x),f(y)) < \varepsilon.
\]

A *homeomorphism* is a continuous bijection whose inverse is also continuous; if there is a homeomorphism between *M*\(_{1}\) and *M*\(_{2}\), they are said to be *homeomorphic*. Homeomorphic spaces are the same from the point of view of topology, but may have very different metric properties. For example, \(\R\) is unbounded and complete, while is bounded but not complete.

### Uniformly continuous maps

A function \(f\,\colon M_1\to M_2\) is *uniformly continuous* if for every real number ε > 0 there exists δ > 0 such that for all points and in *M*\(_{1}\) such that \(d(x,y)<\delta\), we have
\[
d_2(f(x),f(y)) < \varepsilon.
\]

The only difference between this definition and the ε–δ definition of continuity is the order of quantifiers: the choice of δ must depend only on ε and not on the point. However, this subtle change makes a big difference. For example, uniformly continuous maps take Cauchy sequences in *M*\(_{1}\) to Cauchy sequences in *M*\(_{2}\). In other words, uniform continuity preserves some metric properties which are not purely topological.

On the other hand, the Heine–Cantor theorem states that if *M*\(_{1}\) is compact, then every continuous map is uniformly continuous. In other words, uniform continuity cannot distinguish any non-topological features of compact metric spaces.

### Lipschitz maps and contractions

A Lipschitz map is one that stretches distances by at most a bounded factor. Formally, given a real number *K* > 0, the map \(f\,\colon M_1\to M_2\) is -*Lipschitz* if

\[
d_2(f(x),f(y))\leq K d_1(x,y)\quad\text{for all}\quad x,y\in M_1.
\]

Lipschitz maps are particularly important in metric geometry; since they provide more flexibility than distance-preserving maps, but still make essential use of the metric. For example, a curve in a metric space is rectifiable (has finite length) if and only if it has a Lipschitz reparametrization.

A 1-Lipschitz map is sometimes called a *nonexpanding* or *metric map*. Metric maps are commonly taken to be the morphisms of the category of metric spaces.

A -Lipschitz map for *K* < 1 is called a *contraction*. The Banach fixed-point theorem states that if is a complete metric space, then every contraction \(f:M \to M\) admits a unique fixed point. If the metric space is compact, the result holds for a slightly weaker condition on : a map \(f:M \to M\) admits a unique fixed point if

\[
d(f(x), f(y)) < d(x, y) \quad \mbox{for all} \quad x \ne y \in M_1.
\]

### Quasi-isometries

A quasi-isometry is a map that preserves the "large-scale structure" of a metric space. Quasi-isometries need not be continuous. For example, \(\R^2\) and its subspace \(\Z^2\) are quasi-isometric, even though one is connected and the other is discrete. The equivalence relation of quasi-isometry is important in geometric group theory: the Švarc–Milnor lemma states that all spaces on which a group acts geometrically are quasi-isometric.

Formally, the map \(f\,\colon M_1\to M_2\) is a *quasi-isometric embedding* if there exist constants *A* ≥ 1 and *B* ≥ 0 such that

\[
\frac{1}{A} d_2(f(x),f(y))-B\leq d_1(x,y)\leq A d_2(f(x),f(y))+B \quad\text{ for all }\quad x,y\in M_1.
\]

It is a *quasi-isometry* if in addition it is *quasi-surjective*, i.e. there is a constant *C* ≥ 0 such that every point in \(M_2\) is at distance at most from some point in the image \(f(M_1)\).

### Notions of metric space equivalence

Given two metric spaces \((M_1, d_1)\) and \((M_2, d_2)\):
*They are called **homeomorphic** (topologically isomorphic) if there is a homeomorphism between them (i.e., a continuous bijection with a continuous inverse). If \(M_1=M_2\) and the identity map is a homeomorphism, then \(d_1\) and \(d_2\) are said to be **topologically equivalent**.
*They are called **uniformic** (uniformly isomorphic) if there is a uniform isomorphism between them (i.e., a uniformly continuous bijection with a uniformly continuous inverse).
*They are called **bilipschitz homeomorphic** if there is a bilipschitz bijection between them (i.e., a Lipschitz bijection with a Lipschitz inverse).
*They are called **isometric** if there is a (bijective) isometry between them. In this case, the two metric spaces are essentially identical.
*They are called **quasi-isometric** if there is a quasi-isometry between them.

## Metric spaces with additional structure
### Normed vector spaces

A normed vector space is a vector space equipped with a *norm*, which is a function that measures the length of vectors. The norm of a vector is typically denoted by \(\lVert v \rVert\). Any normed vector space can be equipped with a metric in which the distance between two vectors and is given by

\[
d(x,y):=\lVert x-y \rVert.
\]

The metric is said to be *induced* by the norm \(\lVert{\cdot}\rVert\).

Conversely, if a metric on a vector space is
* translation invariant: \(d(x,y) = d(x+a,y+a)\) for every , and in ; and
* : \(d(\alpha x, \alpha y) = |\alpha| d(x,y)\) for every and in and real number α;
then \(\lVert x \rVert := d(x,0)\) is a norm induced by the metric.
A similar relationship holds between seminorms and pseudometrics.

Among examples of metrics induced by a norm are the metrics *d*\(_{1}\), *d*\(_{2}\), and *d*\(_{∞}\) on \(\R^2\), which are induced by the Manhattan norm, the Euclidean norm, and the maximum norm, respectively. More generally, the Kuratowski embedding allows one to see any metric space as a subspace of a normed vector space.

Infinite-dimensional normed vector spaces, particularly spaces of functions, are studied in functional analysis. Completeness is particularly important in this context: a complete normed vector space is known as a Banach space. An unusual property of normed vector spaces is that linear transformations between them are continuous if and only if they are Lipschitz. Such transformations are known as bounded operators.

### Length spaces

A curve in a metric space (*M*, *d*) is a continuous function \(\gamma:[0,T] \to M\). The length of γ is measured by

\[
L(\gamma)=\sup_{0=x_0intrinsic on by setting the distance between points and to be the infimum of the -lengths of paths between them. For instance, if is the straight-line distance on the sphere, then *d*_{intrinsic} is the great-circle distance. However, in some cases *d*_{intrinsic} may have infinite values. For example, if is the Koch snowflake with the subspace metric induced from \(\R^2\), then the resulting intrinsic distance is infinite for any pair of distinct points.

### Riemannian manifolds

A Riemannian manifold is a space equipped with a Riemannian metric tensor, which determines lengths of tangent vectors at every point. This can be thought of defining a notion of distance infinitesimally. In particular, a differentiable path \(\gamma:[0, T] \to M\) in a Riemannian manifold has length defined as the integral of the length of the tangent vector to the path:

\[
L(\gamma)=\int_0^T |\dot\gamma(t)|dt.
\]

On a connected Riemannian manifold, one then defines the distance between two points as the infimum of lengths of smooth paths between them. This construction generalizes to other kinds of infinitesimal metrics on manifolds, such as sub-Riemannian and Finsler metrics.

The Riemannian metric is uniquely determined by the distance function; this means that in principle, all information about a Riemannian manifold can be recovered from its distance function. One direction in metric geometry is finding purely metric ("synthetic") formulations of properties of Riemannian manifolds. For example, a Riemannian manifold is a CAT(*k*) space (a synthetic condition which depends purely on the metric) if and only if its sectional curvature is bounded above by. Thus CAT(*k*) spaces generalize upper curvature bounds to general metric spaces.

### Metric measure spaces
Real analysis makes use of both the metric on \(\R^n\) and the Lebesgue measure. Therefore, generalizations of many ideas from analysis naturally reside in metric measure spaces: spaces that have both a measure and a metric which are compatible with each other. Formally, a *metric measure space* is a metric space equipped with a Borel regular measure such that every ball has positive measure. For example Euclidean spaces of dimension , and more generally -dimensional Riemannian manifolds, naturally have the structure of a metric measure space, equipped with the Lebesgue measure. Certain fractal metric spaces such as the Sierpiński gasket can be equipped with the α-dimensional Hausdorff measure where α is the Hausdorff dimension. In general, however, a metric space may not have an "obvious" choice of measure.

One application of metric measure spaces is generalizing the notion of Ricci curvature beyond Riemannian manifolds. Just as CAT(*k*) and Alexandrov spaces generalize sectional curvature bounds, RCD spaces are a class of metric measure spaces which generalize lower bounds on Ricci curvature.

## Further examples and applications
### Graphs and finite metric spaces
A if its induced topology is the discrete topology. Although many concepts, such as completeness and compactness, are not interesting for such spaces, they are nevertheless an object of study in several branches of mathematics. In particular, (those having a finite number of points) are studied in combinatorics and theoretical computer science. Embeddings in other metric spaces are particularly well-studied. For example, not every finite metric space can be isometrically embedded in a Euclidean space or in Hilbert space. On the other hand, in the worst case the required distortion (bilipschitz constant) is only logarithmic in the number of points.

For any undirected connected graph , the set of vertices of can be turned into a metric space by defining the distance between vertices and to be the length of the shortest edge path connecting them. In graph theory, this is also called the shortest-path distance or geodesic distance. In geometric group theory this construction is applied to the Cayley graph of a (typically infinite) finitely-generated group, yielding the word metric. Up to quasi-isometry, the word metric depends only on the group and not on the chosen finite generating set.

### Metric embeddings and approximations
An important area of study in finite metric spaces is the embedding of complex metric spaces into simpler ones while controlling the distortion of distances. This is particularly useful in computer science and discrete mathematics, where algorithms often perform more efficiently on simpler structures like tree metrics.

A significant result in this area is that any finite metric space can be probabilistically embedded into a *tree metric* with an expected distortion of \(O(\log n)\), where \(n\) is the number of points in the metric space.

This embedding is notable because it achieves the best possible asymptotic bound on distortion, matching the lower bound of \(\Omega(\log n)\). The tree metrics produced in this embedding *dominate* the original metrics, meaning that distances in the tree are greater than or equal to those in the original space. This property is particularly useful for designing approximation algorithms, as it allows for the preservation of distance-related properties while simplifying the underlying structure.

The result has significant implications for various computational problems:

* **Network design**: Improves approximation algorithms for problems like the *Group Steiner tree problem* (a generalization of the *Steiner tree problem*) and *Buy-at-bulk network design* (a problem in *Network planning and design*) by simplifying the metric space to a tree metric.
* **Clustering**: Enhances algorithms for clustering problems where hierarchical clustering can be performed more efficiently on tree metrics.
* **Online algorithms**: Benefits problems like the *k-server problem* and *metrical task system* by providing better competitive ratios through simplified metrics.

The technique involves constructing a hierarchical decomposition of the original metric space and converting it into a tree metric via a randomized algorithm. The \(O(\log n)\) distortion bound has led to improved approximation ratios in several algorithmic problems, demonstrating the practical significance of this theoretical result.

### Distances between mathematical objects
In modern mathematics, one often studies spaces whose points are themselves mathematical objects. A distance function on such a space generally aims to measure the dissimilarity between two objects. Here are some examples:
* **Functions to a metric space.** If is any set and is a metric space, then the set of all bounded functions \(f \colon X \to M\) (i.e. those functions whose image is a bounded subset of \(M\)) can be turned into a metric space by defining the distance between two bounded functions and to be
\[
d(f,g) = \sup_{x \in X} d(f(x),g(x)).
\]
 This metric is called the uniform metric or supremum metric. If is complete, then this function space is complete as well; moreover, if is also a topological space, then the subspace consisting of all bounded continuous functions from to is also complete. When is a subspace of \(\R^n\), this function space is known as a classical Wiener space.
* **String metrics and edit distances.** There are many ways of measuring distances between strings of characters, which may represent sentences in computational linguistics or code words in coding theory. *Edit distances* attempt to measure the number of changes necessary to get from one string to another. For example, the Hamming distance measures the minimal number of substitutions needed, while the Levenshtein distance measures the minimal number of deletions, insertions, and substitutions; both of these can be thought of as distances in an appropriate graph.
* Graph edit distance is a measure of dissimilarity between two graphs, defined as the minimal number of graph edit operations required to transform one graph into another.
* Wasserstein metrics measure the distance between two measures on the same metric space. The Wasserstein distance between two measures is, roughly speaking, the cost of transporting one to the other.
* The set of all by matrices over some field is a metric space with respect to the rank distance \(d(A,B) = \mathrm{rank}(B - A)\).
* The Helly metric in game theory measures the difference between strategies in a game.

### Hausdorff and Gromov–Hausdorff distance
The idea of spaces of mathematical objects can also be applied to subsets of a metric space, as well as metric spaces themselves. Hausdorff and Gromov–Hausdorff distance define metrics on the set of compact subsets of a metric space and the set of compact metric spaces, respectively.

Suppose (*M*, *d*) is a metric space, and let be a subset of. The *distance from to a point of * is, informally, the distance from to the closest point of. However; since there may not be a single closest point, it is defined via an infimum:

\[
d(x,S) = \inf\{d(x,s) : s \in S \}.
\]

In particular, \(d(x, S)=0\) if and only if belongs to the closure of. Furthermore, distances between points and sets satisfy a version of the triangle inequality:

\[
d(x,S) \leq d(x,y) + d(y,S),
\]

and therefore the map \(d_S:M \to \R\) defined by \(d_S(x)=d(x,S)\) is continuous. Incidentally, this shows that metric spaces are completely regular.

Given two subsets and of , their *Hausdorff distance* is

\[
d_H(S,T) = \max \{ \sup\{d(s,T) : s \in S \} , \sup\{ d(t,S) : t \in T \} \}.
\]

Informally, two sets and are close to each other in the Hausdorff distance if no element of is too far from and vice versa. For example, if is an open set in Euclidean space is an ε-net inside , then \(d_H(S,T)<\varepsilon\). In general, the Hausdorff distance \(d_H(S,T)\) can be infinite or zero. However, the Hausdorff distance between two distinct compact sets is always positive and finite. Thus the Hausdorff distance defines a metric on the set of compact subsets of. The Gromov–Hausdorff metric defines a distance between (isometry classes of) compact metric spaces. The *Gromov–Hausdorff distance* between compact spaces and is the infimum of the Hausdorff distance over all metric spaces that contain and as subspaces. While the exact value of the Gromov–Hausdorff distance is rarely useful to know, the resulting topology has found many applications.

=== Miscellaneous examples ===
* Given a metric space (*X*, *d*) and an increasing concave function \(f \colon [0,\infty) \to [0,\infty)\) such that *f*(*t*) if and only if *t* , then \(d_f(x,y)=f(d(x,y))\) is also a metric on. If *f*(*t*) for some real number α < 1, such a metric is known as a **snowflake** of. * The tight span of a metric space is another metric space which can be thought of as an abstract version of the convex hull.
* The *knight's move metric*, the minimal number of knight's moves to reach one point in \(\mathbb{Z}^2\) from another, is a metric on \(\mathbb{Z}^2\).
* The British Rail metric (also called the "post office metric" or the "French railway metric") on a normed vector space is given by \(d(x,y) = \lVert x \rVert + \lVert y \rVert\) for distinct points \(x\) and \(y\), and \(d(x,x) = 0\). More generally \(\lVert \cdot \rVert\) can be replaced with a function \(f\) taking an arbitrary set \(S\) to non-negative reals and taking the value \(0\) at most once: then the metric is defined on \(S\) by \(d(x,y) = f(x) + f(y)\) for distinct points \(x\) and \(y\), and \(d(x,x) = 0\). The name alludes to the tendency of railway journeys to proceed via London (or Paris) irrespective of their final destination.
* The Robinson–Foulds metric used for calculating the distances between Phylogenetic trees in Phylogenetics

## Constructions
### Product metric spaces

If \((M_1,d_1),\ldots,(M_n,d_n)\) are metric spaces, and is the Euclidean norm on \(\mathbb R^n\), then \(\bigl(M_1 \times \cdots \times M_n, d_\times\bigr)\) is a metric space, where the product metric is defined by

\[
d_\times\bigl((x_1,\ldots,x_n),(y_1,\ldots,y_n)\bigr) = N\bigl(d_1(x_1,y_1),\ldots,d_n(x_n,y_n)\bigr),
\]

and the induced topology agrees with the product topology. By the equivalence of norms in finite dimensions, a topologically equivalent metric is obtained if is the taxicab norm, a p-norm, the maximum norm, or any other norm which is non-decreasing as the coordinates of a positive -tuple increase (yielding the triangle inequality).

Similarly, a metric on the topological product of countably many metric spaces can be obtained using the metric

\[
d(x,y)=\sum_{i=1}^\infty \frac1{2^i}\frac{d_i(x_i,y_i)}{1+d_i(x_i,y_i)}.
\]

The topological product of uncountably many metric spaces need not be metrizable. For example, an uncountable product of copies of \(\mathbb{R}\) is not first-countable and thus is not metrizable.

### Quotient metric spaces
If is a metric space with metric , and \(\sim\) is an equivalence relation on , then we can endow the quotient set \(M/\!\sim\) with a pseudometric. The distance between two equivalence classes \([x]\) and \([y]\) is defined as

\[
d'([x],[y]) = \inf\{d(p_1,q_1)+d(p_2,q_2)+\dotsb+d(p_{n},q_{n})\},
\]

where the infimum is taken over all finite sequences \((p_1, p_2, \dots, p_n)\) and \((q_1, q_2, \dots, q_n)\) with \(p_1 \sim x\), \(q_n \sim y\), \(q_i \sim p_{i+1}, i=1,2,\dots, n-1\). In general this will only define a pseudometric, i.e. \(d'([x],[y])=0\) does not necessarily imply that \([x]=[y]\). However, for some equivalence relations (e.g., those given by gluing together polyhedra along faces), \(d'\) is a metric.

The quotient metric \(d'\) is characterized by the following universal property. If \(f\,\colon(M,d)\to(X,\delta)\) is a metric (i.e. 1-Lipschitz) map between metric spaces satisfying *f*(*x*) whenever \(x \sim y\), then the induced function \(\overline{f}\,\colon {M/\sim}\to X\), given by \(\overline{f}([x])=f(x)\), is a metric map \(\overline{f}\,\colon (M/\sim,d')\to (X,\delta).\)

The quotient metric does not always induce the quotient topology. For example, the topological quotient of the metric space \(\N \times [0,1]\) identifying all points of the form \((n, 0)\) is not metrizable since it is not first-countable, but the quotient metric is a well-defined metric on the same set which induces a coarser topology. Moreover, different metrics on the original topological space (a disjoint union of countably many intervals) lead to different topologies on the quotient.

A topological space is sequential if and only if it is a (topological) quotient of a metric space.

## Generalizations of metric spaces
There are several notions of spaces which have less structure than a metric space, but more than a topological space.
* Uniform spaces are spaces in which distances are not defined, but uniform continuity is.
* Approach spaces are spaces in which point-to-set distances are defined, instead of point-to-point distances. They have particularly good properties from the point of view of category theory.
* Continuity spaces are a generalization of metric spaces and posets that can be used to unify the notions of metric spaces and domains.

There are also numerous ways of relaxing the axioms for a metric, giving rise to various notions of generalized metric spaces. These generalizations can also be combined. The terminology used to describe them is not completely standardized. Most notably, in functional analysis pseudometrics often come from seminorms on vector spaces, and so it is natural to call them "semimetrics". This conflicts with the use of the term in topology.

### Extended metrics
Some authors define metrics so as to allow the distance function to attain the value ∞, i.e. distances are non-negative numbers on the extended real number line. Such a function is also called an *extended metric* or "∞-metric". Every extended metric can be replaced by a real-valued metric that is topologically equivalent. This can be done using a subadditive monotonically increasing bounded function which is zero at zero, e.g. \(d'(x, y) = d(x, y) / (1 + d(x, y))\) or \(d''(x, y) = \min(1, d(x, y))\).

### Metrics valued in structures other than the real numbers
The requirement that the metric take values in \([0,\infty)\) can be relaxed to consider metrics with values in other structures, including:
* Ordered fields, yielding the notion of a generalised metric.
* More general directed sets. In the absence of an addition operation, the triangle inequality does not make sense and is replaced with an ultrametric inequality. This leads to the notion of a *generalized ultrametric*.
These generalizations still induce a uniform structure on the space.

### Pseudometrics

A *pseudometric* on \(X\) is a function \(d: X \times X \to \R\) which satisfies the axioms for a metric, except that instead of the second (identity of indiscernibles) only \(d(x,x)=0\) for all *\(x\)* is required. In other words, the axioms for a pseudometric are:

# \(d(x, y) \geq 0\)
# \(d(x,x)=0\)
# \(d(x,y)=d(y,x)\)
# \(d(x,z)\leq d(x,y) + d(y,z)\).

In some contexts, pseudometrics are referred to as *semimetrics* because of their relation to seminorms.

### Quasimetrics
Occasionally, a **quasimetric** is defined as a function that satisfies all axioms for a metric with the possible exception of symmetry. The name of this generalisation is not entirely standardized.

# \(d(x, y) \geq 0\)
# \(d(x,y)=0 \iff x=y\)
# \(d(x,z) \leq d(x,y) + d(y,z)\)

Quasimetrics are common in real life. For example, given a set of mountain villages, the typical walking times between elements of form a quasimetric because travel uphill takes longer than travel downhill. Another example is the length of car rides in a city with one-way streets: here, a shortest path from point to point goes along a different set of streets than a shortest path from to and may have a different length.

A quasimetric on the reals can be defined by setting

\[
d(x,y)=\begin{cases}
x-y & \text{if }x\geq y,\\
1 & \text{otherwise.}
\end{cases}
\]

The 1 may be replaced, for example, by infinity or by \(1 + \sqrt{y-x}\) or any other subadditive function of *y*-*x*. This quasimetric describes the cost of modifying a metal stick: it is easy to reduce its size by filing it down, but it is difficult or impossible to grow it.

Given a quasimetric on , one can define an -ball around to be the set \(\{y \in X | d(x,y) \leq R\}\). As in the case of a metric, such balls form a basis for a topology on , but this topology need not be metrizable. For example, the topology induced by the quasimetric on the reals described above is the (reversed) Sorgenfrey line.

### Metametrics or partial metrics
In a *metametric*, all the axioms of a metric are satisfied except that the distance between identical points is not necessarily zero. In other words, the axioms for a metametric are:

# \(d(x,y)\geq 0\)
# \(d(x,y)=0 \implies x=y\)
# \(d(x,y)=d(y,x)\)
# \(d(x,z)\leq d(x,y)+d(y,z).\)

Metametrics appear in the study of Gromov hyperbolic metric spaces and their boundaries. The *visual metametric* on such a space satisfies \(d(x,x)=0\) for points \(x\) on the boundary, but otherwise \(d(x,x)\) is approximately the distance from *\(x\)* to the boundary. Metametrics were first defined by Jussi Väisälä. In other work, a function satisfying these axioms is called a *partial metric* or a *dislocated metric*.

### Semimetrics
A **semimetric** on \(X\) is a function \(d: X \times X \to \R\) that satisfies the first three axioms, but not necessarily the triangle inequality:

# \(d(x,y)\geq 0\)
# \(d(x,y)=0 \iff x=y\)
# \(d(x,y)=d(y,x)\)

Some authors work with a weaker form of the triangle inequality, such as:

<dl>
  <dt>\(d(x,z)\leq \rho\,(d(x,y)+d(y,z))\)</dt>
  <dd>ρ-relaxed triangle inequality</dd>
  <dt>\(d(x,z)\leq \rho\,\max\{d(x,y),d(y,z)\}\)</dt>
  <dd>ρ-inframetric inequality</dd>
</dl>

The ρ-inframetric inequality implies the ρ-relaxed triangle inequality (assuming the first axiom), and the ρ-relaxed triangle inequality implies the 2ρ-inframetric inequality. Semimetrics satisfying these equivalent conditions have sometimes been referred to as *quasimetrics*, *nearmetrics* or **inframetrics**.

The ρ-inframetric inequalities were introduced to model round-trip delay times in the internet. The triangle inequality implies the 2-inframetric inequality, and the ultrametric inequality is exactly the 1-inframetric inequality.

### Premetrics
Relaxing the last three axioms leads to the notion of a **premetric**, i.e. a function satisfying the following conditions:

# \(d(x,y)\geq 0\)
# \(d(x,x)=0\)

This is not a standard term. Sometimes it is used to refer to other generalizations of metrics such as pseudosemimetrics or pseudometrics; in translations of Russian books it sometimes appears as "prametric". A premetric that satisfies symmetry, i.e. a pseudosemimetric, is also called a distance.

Any premetric gives rise to a topology as follows. For a positive real \(r\), the centered at a point \(p\) is defined as
\(B_r(p)=\{ x | d(x,p) < r \}.\)
A set is called *open* if for any point *\(p\)* in the set there is an centered at *\(p\)* which is contained in the set. Every premetric space is a topological space, and in fact a sequential space.
In general, the themselves need not be open sets with respect to this topology.
As for metrics, the distance between two sets \(A\) and *\(B\)*, is defined as
\(d(A,B)=\underset{x\in A, y\in B}\inf d(x,y).\)
This defines a premetric on the power set of a premetric space. If we start with a (pseudosemi-)metric space, we get a pseudosemimetric, i.e. a symmetric premetric.
Any premetric gives rise to a preclosure operator \(cl\) as follows:
\(cl(A)=\{ x | d(x,A) = 0 \}.\)

### Pseudoquasimetrics
The prefixes *pseudo-*, *quasi-* and *semi-* can also be combined, e.g., a **pseudoquasimetric** (sometimes called **hemimetric**) relaxes both the indiscernibility axiom and the symmetry axiom and is simply a premetric satisfying the triangle inequality. For pseudoquasimetric spaces the open form a basis of open sets. A very basic example of a pseudoquasimetric space is the set \(\{0,1\}\) with the premetric given by \(d(0,1) = 1\) and \(d(1,0) = 0.\) The associated topological space is the Sierpiński space.

Sets equipped with an extended pseudoquasimetric were studied by William Lawvere as "generalized metric spaces". From a categorical point of view, the extended pseudometric spaces and the extended pseudoquasimetric spaces, along with their corresponding nonexpansive maps, are the best behaved of the metric space categories. One can take arbitrary products and coproducts and form quotient objects within the given category. If one drops "extended", one can only take finite products and coproducts. If one drops "pseudo", one cannot take quotients.

Lawvere also gave an alternate definition of such spaces as enriched categories. The ordered set \((\mathbb{R},\geq)\) can be seen as a category with one morphism \(a\to b\) if \(a\geq b\) and none otherwise. Using + as the tensor product and 0 as the identity makes this category into a monoidal category \(R^*\).
Every (extended pseudoquasi-)metric space \((M,d)\) can now be viewed as a category \(M^*\) enriched over \(R^*\):
* The objects of the category are the points of. * For every pair of points and such that \(d(x,y)<\infty\), there is a single morphism which is assigned the object \(d(x,y)\) of \(R^*\).
* The triangle inequality and the fact that \(d(x,x)=0\) for all points derive from the properties of composition and identity in an enriched category.
* Since \(R^*\) is a poset, all diagrams that are required for an enriched category commute automatically.

### Metrics on multisets
The notion of a metric can be generalized from a distance between two elements to a number assigned to a multiset of elements. A multiset is a generalization of the notion of a set in which an element can occur more than once. Define the multiset union \(U=XY\) as follows: if an element occurs times in and times in then it occurs *m* + *n* times in. A function on the set of nonempty finite multisets of elements of a set is a metric if
# \(d(X)=0\) if all elements of are equal and \(d(X) > 0\) otherwise (positive definiteness)
# \(d(X)\) depends only on the (unordered) multiset (symmetry)
# \(d(XY) \leq d(XZ)+d(ZY)\) (triangle inequality)
By considering the cases of axioms 1 and 2 in which the multiset has two elements and the case of axiom 3 in which the multisets , and have one element each, one recovers the usual axioms for a metric. That is, every multiset metric yields an ordinary metric when restricted to sets of two elements.

A simple example is the set of all nonempty finite multisets \(X\) of integers with \(d(X)=\max (X)- \min (X)\). More complex examples are information distance in multisets; and normalized compression distance (NCD) in multisets.
