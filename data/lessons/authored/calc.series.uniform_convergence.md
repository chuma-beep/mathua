# Uniform convergence on a set

Uniform convergence controls a whole family of functions at once. Pointwise convergence fixes one \(x\) at a time, while uniform convergence asks for one bound that works for every \(x\) in the set.

## Sequences of functions

A sequence \(f_n\) converges uniformly to \(f\) on a set \(S\) when
\[
\sup_{x\in S}|f_n(x)-f(x)|\to0.
\]
The supremum is taken before the limit. If the supremum stays away from zero, the convergence is only pointwise.

For \(f_n(x)=x^n\) on \([0,1]\), the pointwise limit is zero for \(x<1\) and 1 at \(x=1\). The supremum of \(|f_n|\) is 1 for every \(n\), so the convergence is not uniform. On \([0,1/2]\) the supremum is \((1/2)^n\), which tends to zero.

## Series of functions

A series of functions converges uniformly when its sequence of partial sums does. The Weierstrass M-test gives a direct sufficient condition. If \(|f_n(x)|\le M_n\) for every \(x\) and \(\sum M_n\) converges, then \(\sum f_n\) converges uniformly.

The test turns a function question into a numerical one. A bound such as \(1/n^2\) is summable, so the series is uniform. A bound such as \(1/n\) is not summable, and the test gives no conclusion.

## The supremum test

The definition and the M-test are two views of the same idea. The definition asks for the actual supremum of the remainder. The M-test replaces that supremum by a simpler numerical bound.

When a series converges on a closed interval by the alternating-series test, the remainder is bounded by the first omitted term. If that bound is independent of \(x\) and tends to zero, the convergence is uniform. On a half-open interval the same bound may fail at the endpoint, which is why the endpoint often decides the answer.
