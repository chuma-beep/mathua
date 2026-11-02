# Total Boundedness

## Balls of radius epsilon

A metric space is totally bounded when, for every positive epsilon, finitely many balls of radius epsilon cover the space. The number of balls may depend on epsilon, and the centres must lie in the space. Total boundedness is a uniform finiteness condition, stronger than ordinary boundedness in general metric spaces.

In Euclidean space a bounded set is totally bounded, because it fits inside a large box that can be divided into finitely many small boxes. In an infinite-dimensional space the unit ball is bounded but not totally bounded: no finite collection of small balls can cover it. This is one of the central differences between finite and infinite dimensions.

## Total boundedness and compactness

For metric spaces, compactness is equivalent to completeness together with total boundedness. Completeness prevents sequences from escaping through missing limits, while total boundedness prevents them from spreading out so far that no subsequence can cluster. Together the two conditions give the Bolzano-Weierstrass property.

The proof of the forward direction uses the open cover by balls of radius epsilon. Compactness extracts a finite subcover. The reverse direction builds a sequence with no convergent subsequence from a failure of total boundedness, or uses a finite epsilon-net to construct a Cauchy subsequence.

## Finite nets

An epsilon-net is a finite set whose epsilon-balls cover the space. Total boundedness says that an epsilon-net exists for every positive epsilon. Nets are useful computationally because they reduce an infinite space to finitely many centres while controlling the error by epsilon.

The covering number is the smallest size of an epsilon-net. It measures the complexity of the space at scale epsilon. For an interval of length L the covering number grows like L divided by epsilon, while for a high-dimensional cube it grows like the ratio raised to the dimension. That growth is the source of the curse of dimensionality.
