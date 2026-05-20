> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Open_set) — CC BY-SA 4.0

# Open and closed sets

In general topology and mathematical analysis, an **open set** is a generalization of an open interval in the real line.

In a metric space (a set with a distance defined between every two points), an open set is a set that, with every point  in it, contains all points of the metric space that are sufficiently near to  (that is, all points whose distance to  is less than some value depending on ).

More generally, an open set is a member of a given collection of subsets of a given set, a collection that has the property of containing every union of its members, every finite intersection of its members, the empty set, and the whole set itself. A set in which such a collection is given is called a topological space, and the collection is called a topology. These conditions are very loose, and allow enormous flexibility in the choice of open sets. For example, *every* subset can be open (the discrete topology), or *no* subset can be open except the space itself and the empty set (the indiscrete topology).

In practice, however, open sets are usually chosen to provide a notion of nearness that is similar to that of metric spaces, without having a notion of distance defined. In particular, a topology allows defining properties such as continuity, connectedness, and compactness, which were originally defined by means of a distance.

The most common case of a topology without any distance is given by manifolds, which are topological spaces that, *near* each point, resemble an open set of a Euclidean space, but on which no distance is defined in general. Less intuitive topologies are used in other branches of mathematics; for example, the Zariski topology, which is fundamental in algebraic geometry and scheme theory.

## Motivation
Intuitively, an open set provides a method to distinguish two points. For example, if about one of two points in a topological space, there exists an open set not containing the other (distinct) point, the two points are referred to as topologically distinguishable. In this manner, one may speak of whether two points, or more generally two subsets, of a topological space are "near" without concretely defining a distance. Therefore, topological spaces may be seen as a generalization of spaces equipped with a notion of distance, which are called metric spaces.

In the set of all real numbers, one has the natural Euclidean metric; that is, a function which measures the distance between two real numbers: 1=*d*(*x*, *y*) = . Therefore, given a real number *x*, one can speak of the set of all points close to that real number; that is, within *ε* of *x*. In essence, points within ε of *x* approximate *x* to an accuracy of degree *ε*. Note that *ε* > 0 always but as *ε* becomes smaller and smaller, one obtains points that approximate *x* to a higher and higher degree of accuracy. For example, if *x* = 0 and *ε* = 1, the points within *ε* of *x* are precisely the points of the interval (−1, 1); that is, the set of all real numbers between −1 and 1. However, with *ε* = 0.5, the points within *ε* of *x* are precisely the points of (−0.5, 0.5). Clearly, these points approximate *x* to a greater degree of accuracy than when *ε* = 1.

The previous discussion shows, for the case *x* = 0, that one may approximate *x* to higher and higher degrees of accuracy by defining *ε* to be smaller and smaller. In particular, sets of the form (−*ε*, *ε*) give us a lot of information about points close to *x* = 0. Thus, rather than speaking of a concrete Euclidean metric, one may use sets to describe points close to *x*. This innovative idea has far-reaching consequences; in particular, by defining different collections of sets containing 0 (distinct from the sets (−*ε*, *ε*)), one may find different results regarding the distance between 0 and other real numbers. For example, if we were to define **R** as the only such set for "measuring distance", all points would be close to 0 since there is only one possible degree of accuracy one may achieve in approximating 0: being a member of **R**. Thus, we find that in some sense, every real number is distance 0 away from 0. It may help in this case to think of the measure as being a binary condition: all things in **R** are equally close to 0, while any item that is not in **R** is not close to 0.

In general, one refers to the family of sets containing 0, used to approximate 0, as a ***neighborhood basis***; a member of this neighborhood basis is referred to as an open set. In fact, one may generalize these notions to an arbitrary set (*X*); rather than just the real numbers. In this case, given a point (*x*) of that set, one may define a collection of sets "around" (that is, containing) *x*, used to approximate *x*. Of course, this collection would have to satisfy certain properties (known as **axioms**) for otherwise we may not have a well-defined method to measure distance. For example, every point in *X* should approximate *x* to *some* degree of accuracy. Thus *X* should be in this family. Once we begin to define "smaller" sets containing *x*, we tend to approximate *x* to a greater degree of accuracy. Bearing this in mind, one may define the remaining axioms that the family of sets about *x* is required to satisfy.

## Definitions
Several definitions are given here, in an increasing order of technicality. Each one is a special case of the next one.

### Euclidean space
A subset \(U\) of the Euclidean *n*-space **R**<sup>*n*</sup> is *open* if, for every point  in \(U\), there exists a positive real number  (depending on ) such that any point in **R**<sup>*n*</sup> whose Euclidean distance from  is smaller than  belongs to \(U\). Equivalently, a subset \(U\) of **R**<sup>*n*</sup> is open if every point in \(U\) is the center of an open ball contained in \(U.\)

An example of a subset of **R** that is not open is the closed interval , since neither 0 - *ε* nor 1 + *ε* belongs to  for any *ε* > 0, no matter how small.

### Metric space
A subset *U* of a metric space (*M*, *d*) is called *open* if, for any point *x* in *U*, there exists a real number *ε* > 0 such that any point \(y \in M\) satisfying *d*(*x*, *y*) < *ε* belongs to *U*. Equivalently, *U* is open if every point in *U* has a neighborhood contained in *U*.

This generalizes the Euclidean space example, since Euclidean space with the Euclidean distance is a metric space.

### Topological space
A *topology* \(\tau\) on a set  is a set of subsets of  with the following properties:
*\(X \in \tau\) and \(\varnothing \in \tau\).
*Any union of sets in \(\tau\) belong to \(\tau\): if \(\left\{ U_i : i \in I \right\} \subseteq \tau\) then
\[
\bigcup_{i \in I} U_i \in \tau  \,.
\]

*Any finite intersection of sets in \(\tau\) belong to \(\tau\): for a positive integer \(n\), if \(U_1, \ldots, U_n \in \tau\) then
\[
U_1 \cap \cdots \cap U_n \in \tau \,.
\]

Each member of \(\tau\) is called an *open set*.  The set  together with \(\tau\) is called a *topological space*.

Infinite intersections of open sets need not be open. For example, the intersection of all intervals of the form \(\left( -1/n, 1/n \right),\) where \(n\) is a positive integer, is the set \(\{ 0 \}\), which is not open in the real line.

A metric space is a topological space, whose topology consists of the collection of all subsets that are unions of open balls. There are, however, topological spaces that are not metric spaces.

## Properties
The union of any number of open sets, or infinitely many open sets, is open. The intersection of a finite number of open sets is open.

A complement of an open set (relative to the space that the topology is defined on) is called a closed set. A set may be both open and closed (a clopen set). The empty set and the full space are examples of sets that are both open and closed.

A set can never be considered as open by itself. This notion is relative to a containing set and a specific topology on it.

Whether a set is open depends on the topology under consideration. Having opted for greater brevity over greater clarity, we refer to a set *X* endowed with a topology \(\tau\) as "the topological space *X*" rather than "the topological space \((X, \tau)\)", despite the fact that all the topological data is contained in \(\tau.\) If there are two topologies on the same set, a set *U* that is open in the first topology might fail to be open in the second topology. For example, if *X* is any topological space and *Y* is any subset of *X*, the set *Y* can be given its own topology (called the 'subspace topology') defined by "a set *U* is open in the subspace topology on *Y* if and only if *U* is the intersection of *Y* with an open set from the original topology on *X*." This potentially introduces new open sets: if *V* is open in the original topology on *X*, but \(V \cap Y\) isn't open in the original topology on *X*, then \(V \cap Y\) is open in the subspace topology on *Y*.

As a concrete example of this, if *U* is defined as the set of rational numbers in the interval \((0, 1),\) then *U* is an open subset of the rational numbers, but not of the real numbers. This is because when the surrounding space is the rational numbers, for every point *x* in *U*, there exists a positive number *ε* such that all  points within distance *ε* of *x* are also in *U*. On the other hand, when the surrounding space is the reals, then for every point *x* in *U* there is  positive *ε* such that all  points within distance *ε* of *x* are in *U* (because *U* contains no non-rational numbers).

## Uses
Open sets have a fundamental importance in topology. The concept is required to define and make sense of topological space and other topological structures that deal with the notions of closeness and convergence for spaces such as metric spaces and uniform spaces.

Every subset *A* of a topological space *X* contains a (possibly empty) open set; the maximum (ordered under inclusion) such open set is called the interior of *A*.
It can be constructed by taking the union of all the open sets contained in *A*.

A function \(f : X \to Y\) between two topological spaces \(X\) and \(Y\) is  if the preimage of every open set in \(Y\) is open in \(X.\)
The function \(f : X \to Y\) is called  if the image of every open set in \(X\) is open in \(Y.\)

An open set on the real line has the characteristic property that it is a countable union of disjoint open intervals.

## Special types of open sets
### Clopen sets and non-open and/or non-closed sets
A set might be open, closed, both, or neither. In particular, open and closed sets are not mutually exclusive, meaning that it is in general possible for a subset of a topological space to simultaneously be both an open subset  a closed subset. Such subsets are known as ****. Explicitly, a subset \(S\) of a topological space \((X, \tau)\) is called  if both \(S\) and its complement \(X \setminus S\) are open subsets of \((X, \tau)\); or equivalently, if \(S \in \tau\) and \(X \setminus S \in \tau.\)

In  topological space \((X, \tau),\) the empty set \(\varnothing\) and the set \(X\) itself are always clopen. These two sets are the most well-known examples of clopen subsets and they show that clopen subsets exist in  topological space. To see, it suffices to remark that, by definition of a topology, \(X\) and \(\varnothing\) are both open, and that they are also closed, since each is the complement of the other.

The open sets of the usual Euclidean topology of the real line \(\R\) are the empty set, the open intervals and every union of open intervals.

* The interval \(I = (0, 1)\) is open in \(\R\) by definition of the Euclidean topology. It is not closed since its complement in \(\R\) is \(I^\complement = (-\infty, 0] \cup [1, \infty),\) which is not open; indeed, an open interval contained in \(I^\complement\) cannot contain 1, and it follows that \(I^\complement\) cannot be a union of open intervals. Hence, \(I\) is an example of a set that is open but not closed.
* By a similar argument, the interval \(J = [0, 1]\) is a closed subset but not an open subset.
* Finally, neither \(K = [0, 1)\) nor its complement \(\R \setminus K = (-\infty, 0) \cup [1, \infty)\) are open (because they cannot be written as a union of open intervals); this means that \(K\) is neither open nor closed.

If a topological space \(X\) is endowed with the discrete topology (so that by definition, every subset of \(X\) is open) then every subset of \(X\) is a clopen subset.
For a more advanced example reminiscent of the discrete topology, suppose that \(\mathcal{U}\) is an ultrafilter on a non-empty set \(X.\) Then the union \(\tau := \mathcal{U} \cup \{ \varnothing \}\) is a topology on \(X\) with the property that  non-empty proper subset \(S\) of \(X\) is  an open subset or else a closed subset, but never both; that is, if \(\varnothing \neq S \subsetneq X\) (where \(S \neq X\)) then  of the following two statements is true: either (1) \(S \in \tau\) or else, (2) \(X \setminus S \in \tau.\) Said differently,  subset is open or closed but the   subsets that are both (i.e. that are clopen) are \(\varnothing\) and \(X.\)

### Regular open sets
A subset \(S\) of a topological space \(X\) is called a **** if \(\operatorname{Int} \left( \overline{S} \right) = S\) or equivalently, if \(\operatorname{Bd} \left( \overline{S} \right) = \operatorname{Bd} S\), where \(\operatorname{Bd} S\), \(\operatorname{Int} S\), and \(\overline{S}\) denote, respectively, the topological boundary, interior, and closure of \(S\) in \(X\).  A topological space for which there exists a base consisting of regular open sets is called a ****.
A subset of \(X\) is a regular open set if and only if its complement in \(X\) is a regular closed set, where by definition a subset \(S\) of \(X\) is called a **** if \(\overline{\operatorname{Int} S} = S\) or equivalently, if \(\operatorname{Bd} \left( \operatorname{Int} S \right) = \operatorname{Bd} S.\)
Every regular open set (resp. regular closed set) is an open subset (resp. is a closed subset) although in general, the converses are  true.

## Generalizations of open sets


Throughout, \((X, \tau)\) will be a topological space.

A subset \(A \subseteq X\) of a topological space \(X\) is called:


**** if \(A ~\subseteq~ \operatorname{int}_X \left( \operatorname{cl}_X \left( \operatorname{int}_X A \right) \right)\), and the complement of such a set is called ****.
****, ****, or **** if it satisfies any of the following equivalent conditions:

\(A ~\subseteq~ \operatorname{int}_X \left( \operatorname{cl}_X A \right).\)
There exists subsets \(D, U \subseteq X\) such that \(U\) is open in \(X,\) \(D\) is a dense subset of \(X,\) and \(A = U \cap D.\)
There exists an open (in \(X\)) subset \(U \subseteq X\) such that \(A\) is a dense subset of \(U.\)

The complement of a preopen set is called ****.

**** if \(A ~\subseteq~ \operatorname{int}_X \left( \operatorname{cl}_X A \right) ~\cup~ \operatorname{cl}_X \left( \operatorname{int}_X A \right)\). The complement of a b-open set is called ****.
**** or **** if it satisfies any of the following equivalent conditions:

\(A ~\subseteq~ \operatorname{cl}_X \left( \operatorname{int}_X \left( \operatorname{cl}_X A \right) \right)\)
\(\operatorname{cl}_X A\) is a regular closed subset of \(X.\)
There exists a preopen subset \(U\) of \(X\) such that \(U \subseteq A \subseteq \operatorname{cl}_X U.\)

The complement of a β-open set is called ****.

**** if it satisfies any of the following equivalent conditions:

Whenever a sequence in \(X\) converges to some point of \(A,\) then that sequence is eventually in \(A.\) Explicitly, this means that if \(x_{\bull} = \left( x_i \right)_{i=1}^{\infty}\) is a sequence in \(X\) and if there exists some \(a \in A\) is such that \(x_{\bull} \to x\) in \((X, \tau),\) then \(x_{\bull}\) is eventually in \(A\) (that is, there exists some integer \(i\) such that if \(j \geq i,\) then \(x_j \in A\)).
\(A\) is equal to its **** in \(X,\) which by definition is the set
\(\begin{alignat}{4}
\operatorname{SeqInt}_X A
&= \{ a \in A ~:~ \text{ whenever a sequence in } X \text{ converges to } a \text{ in } (X, \tau), \text{ then that sequence is eventually in } A \} \\
&= \{ a \in A ~:~ \text{ there does NOT exist a sequence in } X \setminus A \text{ that converges in } (X, \tau) \text{ to a point in } A \} \\
\end{alignat}\)


The complement of a sequentially open set is called ****. A subset \(S \subseteq X\) is sequentially closed in \(X\) if and only if \(S\) is equal to its ****, which by definition is the set \(\operatorname{SeqCl}_X S\) consisting of all \(x \in X\) for which there exists a sequence in \(S\) that converges to \(x\) (in \(X\)).

**** and is said to have **** if there exists an open subset \(U \subseteq X\) such that \(A \bigtriangleup U\) is a meager subset, where \(\bigtriangleup\) denotes the symmetric difference.
* The subset \(A \subseteq X\) is said to have **the Baire property in the restricted sense** if for every subset \(E\) of \(X\) the intersection \(A\cap E\) has the Baire property relative to \(E\).
**** if \(A ~\subseteq~ \operatorname{cl}_X \left( \operatorname{int}_X A \right)\) or, equivalently, \(\operatorname{cl}_X A = \operatorname{cl}_X \left( \operatorname{int}_X A \right)\). The complement in \(X\) of a semi-open set is called a ** set**.
* The **** (in \(X\)) of a subset \(A \subseteq X,\) denoted by \(\operatorname{sCl}_X A,\) is the intersection of all semi-closed subsets of \(X\) that contain \(A\) as a subset.
**** if for each \(x \in A\) there exists some semiopen subset \(U\) of \(X\) such that \(x \in U \subseteq \operatorname{sCl}_X U \subseteq A.\)
**** (resp. ****) if its complement in \(X\) is a θ-closed (resp. ) set, where by definition, a subset of \(X\) is called **** (resp. ****) if it is equal to the set of all of its θ-cluster points (resp. δ-cluster points). A point \(x \in X\) is called a **** (resp. a ****) of a subset \(B \subseteq X\) if for every open neighborhood \(U\) of \(x\) in \(X,\) the intersection \(B \cap \operatorname{cl}_X U\) is not empty (resp. \(B \cap \operatorname{int}_X\left( \operatorname{cl}_X U \right)\) is not empty).


Using the fact that
\(A ~\subseteq~ \operatorname{cl}_X A ~\subseteq~ \operatorname{cl}_X B\) and \(\operatorname{int}_X A ~\subseteq~ \operatorname{int}_X B ~\subseteq~ B\)

whenever two subsets \(A, B \subseteq X\) satisfy \(A \subseteq B,\) the following may be deduced:

* Every α-open subset is semi-open, semi-preopen, preopen, and b-open.
* Every b-open set is semi-preopen (i.e. β-open).
* Every preopen set is b-open and semi-preopen.
* Every semi-open set is b-open and semi-preopen.

Moreover, a subset is a regular open set if and only if it is preopen and semi-closed. The intersection of an α-open set and a semi-preopen (resp. semi-open, preopen, b-open) set is a semi-preopen (resp. semi-open, preopen, b-open) set.  Preopen sets need not be semi-open and semi-open sets need not be preopen.

Arbitrary unions of preopen (resp. α-open, b-open, semi-preopen) sets are once again preopen (resp. α-open, b-open, semi-preopen). However, finite intersections of preopen sets need not be preopen. The set of all α-open subsets of a space \((X, \tau)\) forms a topology on \(X\) that is finer than \(\tau.\)

A topological space \(X\) is Hausdorff if and only if every compact subspace of \(X\) is θ-closed.
A space \(X\) is totally disconnected if and only if every regular closed subset is preopen or equivalently, if every semi-open subset is preopen. Moreover, the space is totally disconnected if and only if the **** of every preopen subset is open.

## See also
*
*
*
*
*
*
*
*

## Notes


## References


## Bibliography
*
*
*

## External links
*

