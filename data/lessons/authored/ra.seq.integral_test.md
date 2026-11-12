# The integral test

The integral test compares a series with an improper integral. It applies to positive terms that decrease steadily, and it gives an error bound along with the verdict.

## Decreasing terms

Take \(f\) positive, continuous and decreasing on \([1,\infty)\), with \(a_n = f(n)\). The sum \(\sum a_n\) and the integral \(\int_1^\infty f\) converge or diverge together. For \(f(x) = 1/x^2\) the integral is 1, so \(\sum 1/n^2\) converges.

## The integral comparison

The rectangle over \([n, n+1]\) of height \(f(n)\) covers the graph, and the rectangle of height \(f(n+1)\) sits under it. Summing the inequalities gives \(\int_1^{N+1} f \le \sum_{n=1}^N f(n) \le f(1) + \int_1^N f\). For \(f(x) = 1/x\) the integral is \(\log N\), which grows without bound, so the harmonic series diverges.

## Remainder estimates

The tail satisfies \(\int_{N+1}^\infty f \le \sum_{n=N+1}^\infty f(n) \le \int_N^\infty f\). For \(\sum 1/n^2\) with \(N = 10\) the tail lies between \(1/11\) and \(1/10\), roughly. The estimate is tight enough for numerical work and it needs no extra hypothesis.

The test converts a discrete sum into an area, and the same area measures the error of stopping early.
