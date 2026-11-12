# Asymptotic expansions

An asymptotic expansion approximates a function by a series that need not converge, judged by its remainder.

## The meaning of an asymptotic series

A series is asymptotic to a function when the remainder after finitely many terms is smaller than the last term kept.

The approximation improves for a while as terms are added, then worsens, so the best truncation is near the smallest term. The series may diverge for every value of the argument.

At \(x=10\) the first two terms of the exponential integral expansion are 4.54e-6 and 4.09e-6.

## Integrating by parts repeatedly

Repeated integration by parts produces the series, each step lowering the power of the denominator and putting a factorial in the numerator.

For the exponential integral the series is \(e^{-x}/x\) times the alternating sum of \(n!/x^{n}\). The factorials grow faster than the powers.

At \(x=10\) the tenth term is 0.000363 of the prefactor and the eleventh is 0.000399, so the terms stop shrinking.

## A worked expansion

At \(x=10\) the prefactor is about \(4.54\times 10^{-6}\). Keeping one, two and three corrections gives 4.086e-6, 4.177e-6 and 4.150e-6.

The true value is 4.157e-6, so the error falls below 1 percent before the factorials take over.

An expansion is a tool for a fixed large argument, not a convergent representation. The error is controlled by the first omitted term.
