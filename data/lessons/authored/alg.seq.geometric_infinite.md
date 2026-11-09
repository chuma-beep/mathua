# Infinite geometric series

An infinite geometric series adds the terms of a geometric sequence forever. It can have a finite sum when each term is smaller than the one before.

## When the sum exists

The series \(a+ar+ar^2+\cdots\) converges exactly when \(|r|<1\). In that case its sum is \(S=a/(1-r)\).

If \(|r|\ge 1\), the terms do not shrink to zero, so the series diverges and has no finite sum.

## Finding the ratio

If the first term and the sum are known, solve \(S=a/(1-r)\) for \(r\). Rearranging gives \(r=1-a/S\).

The ratio must satisfy \(|r|<1\). A computed ratio outside that range means the data cannot describe a convergent series.

## Finding the first term

Rearrange the sum formula as \(a=S(1-r)\). This is useful when a model gives the long-run total and the decay rate.

For example, a series with ratio \(1/3\) and sum 12 has first term \(12(2/3)=8\).

A repeating decimal is a geometric series in disguise. For instance, \(0.333\ldots\) has first term \(3/10\) and ratio \(1/10\), so its sum is \(1/3\).
