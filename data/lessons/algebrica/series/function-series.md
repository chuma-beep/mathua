> Content sourced from [Algebrica](https://algebrica.org/function-series/) — CC BY-NC 4.0

# Series of Functions

A **series of functions** \(\sum f_n(x)\) assigns to each \(x\) a numerical
series. It can converge at some points and diverge at others; the convergence
set may be an interval, a scattered set, or empty. The central question is
always: does the convergence behave well enough to preserve continuity,
differentiability, and integrability?

## Pointwise vs Uniform Convergence

**Pointwise** convergence fixes \(x\) first, then takes the limit in \(n\).
**Uniform** convergence controls all \(x\) at once: for every
\(\varepsilon > 0\) there is an \(N\) working simultaneously everywhere. A
uniform limit of continuous functions is continuous; a pointwise limit need
not be. The standard casualty is \(\sum x^n\) on \([0,1]\): partial sums are
continuous, the pointwise limit jumps at \(x = 1\).

## The Weierstrass M-Test

When \(|f_n(x)| \le M_n\) for all \(x\) and \(\sum M_n\) converges, the series
\(\sum f_n\) converges **absolutely and uniformly**. This **M-test** is the
workhorse: it certifies uniform convergence without ever computing the limit
function. Example: \(\sum \sin(nx)/n^2\) is squeezed by \(\sum 1/n^2\), hence
uniformly convergent on all of \(\mathbb{R}\), so its sum is continuous.

## Calculus of Uniform Limits

Uniform convergence on \([a,b]\) licenses term-by-term integration, and with
uniformly convergent derivatives it licenses term-by-term differentiation.
Power series converge uniformly on every closed subinterval of their interval
of convergence, which is why the term-by-term rules for power series need no
separate justification. Fourier series show the sharp edge: \(\sum \sin(nx)/n\)
converges pointwise but not uniformly near the jump, and term-by-term
differentiation fails there.
