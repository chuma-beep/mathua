# Egorov's theorem

Egorov's theorem turns almost everywhere convergence into uniform convergence after a small set is discarded. It requires a finite measure space.

## Almost uniform convergence

On a finite measure space, if \(f_n \to f\) almost everywhere then for every \(\varepsilon > 0\) there is a measurable set \(A\) with \(\mu(A) < \varepsilon\) such that \(f_n \to f\) uniformly on the complement. The convergence is called almost uniform.

## The exceptional set

The set \(A\) collects the points where convergence is slow. On \([0,1]\) the functions \(t^n\) converge to 0 pointwise but not uniformly. Take \(A = (0.99, 1]\), of measure 0.01. On the complement the supremum of \(t^n\) is \(0.99^n\), which falls below 0.01 once \(n > 458\).

## Finite measure

Finiteness is essential. On the line the functions \(1_{[n, n+1]}\) converge to 0 pointwise, yet any set on which the convergence is uniform must miss almost all of each far interval. Its complement then has infinite measure, and no small exceptional set exists.

Egorov's theorem says pointwise convergence is uniform once a small set is removed, and finite measure is the price of that removal.
