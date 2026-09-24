# Pigeonhole Principle

**Pigeonhole principle:** If $n$ pigeons go into $k$ holes and $n>k$, some hole holds at least $\lceil n/k\rceil$ pigeons. The basic case $k+1$ pigeons in $k$ holes forces a collision of at least 2.

## Worked: 13 people, 12 months

Put 13 birthdays into 12 month-holes:
1. With 12 or fewer people each month could hold at most one birthday.
2. The 13th person lands in an occupied month, so some month holds at least 2.
3. The ceiling confirms it: $\lceil13/12\rceil=2$.

So any 13 people contain two sharing a birth month, and likewise $n+1$ integers contain two with the same remainder mod $n$.

## Worked: 30 students, 4 grades

Bound the largest grade group among 30 students and 4 possible grades:
1. If every grade held at most 7 students, the class would hold at most $4\cdot7=28<30$.
2. So some grade holds at least 8, and $\lceil30/4\rceil=\lceil7.5\rceil=8$ is the exact bound.
3. The bound is sharp: 29 students can split 8, 7, 7, 7 with no grade above 8.

So $\lceil n/k\rceil$ is the best guarantee obtainable from counts alone.

## Worked: 5 points in a unit square

Force two close points among 5 in a $1\times1$ square:
1. Cut the square into 4 quarter-squares of side $1/2$.
2. With 5 points and 4 quarters, some quarter holds at least 2 points.
3. Two points in one quarter are at most a diagonal apart: $\sqrt2/2$.

So geometry plus counting guarantees a pair within $\sqrt2/2$, and the partition is the whole trick.
