# The root test

The root test decides convergence of a series by the size of the \(n\)-th root of its terms. It is the sharpest of the elementary tests when the terms carry an \(n\)-th power.

## Roots of the terms

For the series \(\sum n/2^n\), the \(n\)-th root of the term is \(n^{1/n}/2\). Since \(n^{1/n}\) tends to 1, the roots tend to \(1/2\). Each term therefore behaves like \((1/2)^n\) in the long run.

## The limiting root

Let \(L = \limsup |a_n|^{1/n}\). If \(L < 1\) the series converges absolutely. If \(L > 1\) the terms do not tend to zero, so the series diverges. At \(L = 1\) the test says nothing. The harmonic series has roots \(n^{-1/n}\) tending to 1, and it diverges. The series \(\sum 1/n^2\) also has limiting root 1, and it converges.

## Choosing a test

Roots suit terms raised to the \(n\)-th power, such as \(\sum (n/(2n+1))^n\), whose root is \(1/(2 + 1/n)\) and tends to \(1/2\). Ratios suit factorials, because \(n!\) cancels cleanly in \(a_{n+1}/a_n\). When a root is exactly 1, compare with a \(p\)-series instead.

The root test reads the geometric decay hidden in a term, and the limiting root is the number that decay produces.
