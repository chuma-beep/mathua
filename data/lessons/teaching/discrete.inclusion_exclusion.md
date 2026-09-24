# Inclusion-Exclusion Principle

**Inclusion-exclusion:** $|\cup A_i|=\sum|A_i|-\sum|A_i\cap A_j|+\sum|A_i\cap A_j\cap A_k|-\cdots$, alternating over intersections. Adding the sets double-counts overlaps; subtracting pairwise intersections over-corrects triple overlaps, and so on.

## Worked: two sets

Take $|A|=3$, $|B|=4$, $|A\cap B|=1$:
1. Adding counts gives $3+4=7$, but the shared element was counted twice.
2. Subtract the overlap once: $7-1=6$.
3. Check: 2 elements only in $A$, 3 only in $B$, 1 shared — $2+3+1=6$.

So $|A\cup B|=|A|+|B|-|A\cap B|$ corrects exactly the double count.

## Worked: integers up to 10 coprime to 6

Count $n\le10$ divisible by neither 2 nor 3:
1. Multiples of 2: $\lfloor10/2\rfloor=5$; of 3: $\lfloor10/3\rfloor=3$; of both, i.e. of 6: $\lfloor10/6\rfloor=1$.
2. Sieve: $10-5-3+1=3$.
3. Survivors are 1, 5, 7 — indeed 3 numbers.

So each divisibility condition removes its multiples and the overlap is added back once.

## Three sets and derangements

For three sets there are 3 singles, 3 pairs, and 1 triple: 7 terms with signs $+,-,+$. Derangements count permutations of $n$ with no fixed point by sieving out each $A_i=\{\text{$i$ fixed}\}$: $!n=n!\sum_{k=0}^n(-1)^k/k!$. In general, unions become alternating sums over all nonempty intersections.
