# Romberg extrapolation

Romberg integration combines trapezoid estimates at finer and finer grids. Richardson extrapolation removes the leading error term at each stage, and the result is a triangular table of increasingly accurate values.

## Richardson extrapolation

Let \(T(h)\) be the trapezoid estimate with step \(h\). For a smooth function the error has an expansion in powers of \(h^2\). The first extrapolated value is
\[
R_{1,1}=T(h/2)+\frac{T(h/2)-T(h)}{3}.
\]
The denominator 3 comes from the ratio of the two error terms. The combination cancels the \(h^2\) error and leaves a value accurate to order \(h^4\).

The same idea applies at every level. Each new column removes one more power of \(h^2\) from the error.

## The second column

The next column uses the denominator \(4^2-1=15\):
\[
R_{k,2}=R_{k,1}+\frac{R_{k,1}-R_{k-1,1}}{15}.
\]
The two entries being combined must come from grids whose step sizes differ by a factor of two. Mixing entries from other levels breaks the cancellation.

In general, column \(j\) uses the divisor \(4^j-1\). The denominators grow quickly, so the later columns need very accurate input values.

## The full table

A Romberg table is built row by row. Each row starts with a trapezoid estimate, then fills columns from left to right. The diagonal entry is the most accurate value in that row.

The order matters. Compute \(R_{1,1}\) from the first two trapezoid values, then \(R_{2,1}\) from the next two, then \(R_{2,2}\) from those two extrapolations. A mistake in an early entry propagates through every entry to its right.
