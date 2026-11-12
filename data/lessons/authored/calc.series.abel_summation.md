# Abel summation by parts

## The summation identity

Abel summation rewrites a sum of products using partial sums. With \(A_n = a_1 + \cdots + a_n\), the identity is \(\sum_{n=1}^{N} a_n b_n = A_N b_N + \sum_{n=1}^{N-1} A_n (b_n - b_{n+1})\). For \(a_n = 1\) and \(b_n = 1/n\) at \(N = 3\) both sides give \(1 + 0.5 + 0.3333 = 1.8333\).

## Applying it to a product

The identity converts a product into a sum weighted by the change in \(b\). For \(a_n = (-1)^{n-1}\) the partial sums stay in the pair 0 and 1, while \(b_n = 1/n\) decreases, so the tail is bounded by \(b_{N+1}\). At \(N = 4\) the tail \(1/5 - 1/6 + 1/7 - 1/8\) equals 0.0512, below the bound 0.2.

## Dirichlet's test

Dirichlet's test follows: bounded partial sums of \(a_n\) plus a decreasing \(b_n\) with limit zero give convergence. The series \(\sum \sin(n)/n\) qualifies, since the partial sums of \(\sin(n)\) never exceed \(1/\sin(1/2) \approx 2.086\). The first four terms total 1.1539. The bound keeps every later partial sum finite.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
