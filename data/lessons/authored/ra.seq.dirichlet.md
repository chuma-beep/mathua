# Dirichlet's test

Dirichlet's test proves convergence of a series whose factors play different roles. One factor has bounded partial sums and the other decreases to zero.

## Partial summation

Summation by parts rewrites \(\sum_{n=1}^N a_n b_n\) as \(A_N b_N + \sum_{n=1}^{N-1} A_n (b_n - b_{n+1})\), where \(A_n\) is the \(n\)-th partial sum of the \(a_n\). The identity is the discrete analogue of integration by parts. It moves the variation from the \(a_n\) onto the \(b_n\).

## Bounded partial sums

Suppose \(|A_n| \le M\) for every \(n\), and \(b_n\) decreases to zero. Then \(\left|\sum a_n b_n\right| \le 2M b_1\). The series converges, though not necessarily absolutely. The factor of two comes from the total drop of the decreasing sequence.

## Conditional examples

Take \(a_n = \sin(nx)\) with \(x\) not a multiple of \(2\pi\). The partial sums are bounded by \(1/|\sin(x/2)|\), because they telescope through the product-to-sum identity. With \(b_n = 1/n\) the series \(\sum \sin(nx)/n\) converges. At \(x = \pi/2\) the sum is the Leibniz series \(1 - 1/3 + 1/5 - \dots\), whose value is \(\pi/4\). Absolute convergence fails at every such \(x\).

Dirichlet's test separates oscillation from decay, and a bounded oscillation multiplied by a decay to zero is still summable.
