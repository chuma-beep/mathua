# Pigeonhole Principle

**Pigeonhole principle:** If $n$ pigeons go into $k$ holes and $n>k$, some hole has at least $\lceil n/k\rceil$ pigeons. In particular, $k+1$ pigeons in $k$ holes guarantees a collision.

## Counting Guarantees

### Simple Collisions
$13$ people, $12$ months → two share a birth month ($\lceil13/12\rceil=2$). $n+1$ integers → two share remainder mod $n$ ($n$ remainders).

### Generalized Version
$n$ pigeons, $k$ holes → one hole has $\lceil n/k\rceil$. For $30$ students and $4$ grades, $\lceil30/4\rceil=8$ share a grade.

## Example

$5$ points in a unit square: partition into $4$ squares of side $1/2$. One small square contains $\ge2$ points; its diagonal $\sqrt2/2$ bounds their distance.
