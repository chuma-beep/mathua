# The Hajnal-Szemeredi theorem

The Hajnal-Szemeredi theorem guarantees an equitable colouring with one more colour than the maximum degree.

## Equitable colouring

An equitable \(k\)-colouring partitions the vertices into \(k\) independent classes whose sizes differ by at most one. The theorem says every graph with maximum degree \(\Delta\) has an equitable \((\Delta+1)\)-colouring.

## Class sizes

With \(n\) vertices and \(\Delta+1\) classes, the largest class has \(\lceil n/(\Delta+1)\rceil\) vertices and the smallest has \(\lfloor n/(\Delta+1)\rfloor\). The complete graph \(K_n\) forces exactly this spread.

## Sharpness

The complete graph \(K_{\Delta+1}\) needs \(\Delta+1\) colours, and a disjoint union of cliques of order \(\Delta+1\) shows that fewer colours cannot be equitable in general. The bound is tight.

To check your understanding, state the class-size condition, divide the order by the colours and build a tight example.

One colour beyond the maximum degree is always enough, and the class sizes can always be balanced.
