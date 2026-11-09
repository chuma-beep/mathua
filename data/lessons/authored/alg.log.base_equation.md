# Solving for a logarithm's base

A logarithm records an exponent, and that definition turns a logarithmic equation into a power equation. When the argument and the logarithm's value are known, the base is the only unknown, and solving for it is a root problem.

## Base from a power

The equation \(\log_b(N) = y\) means \(b^{y} = N\). If \(N = 64\) and \(y = 3\), then \(b^{3} = 64\), so \(b = 4\). The base is the cube root of the argument in that example.

The base must be positive and not equal to one. Those conditions come from the definition of a logarithm, and they should be checked after the arithmetic. A negative candidate is discarded even if its power matches.

Taking the \(y\)-th root is the same as raising both sides to the power \(1/y\). For \(y = 2\) the operation is a square root, and for \(y = 3\) it is a cube root.

## Base from a negative exponent

A negative logarithm value means the argument is a reciprocal. The equation \(\log_b(1/8) = -3\) becomes \(b^{-3} = 1/8\). Taking reciprocals gives \(b^{3} = 8\), so \(b = 2\).

The sign of the exponent and the size of the argument move in opposite directions. A base larger than one with a negative exponent produces a number below one, which is why the argument is a fraction.

Rewriting the fraction as a positive power is the cleanest route. It turns the problem back into the previous case and keeps the arithmetic in whole numbers.

## Base from a fractional exponent

A fractional logarithm value produces a power equation with a fractional exponent. The equation \(\log_b(3) = 1/2\) means \(b^{1/2} = 3\), and squaring both sides gives \(b = 9\).

The operation that undoes the exponent \(1/2\) is squaring. In general, the operation that undoes \(1/n\) is the \(n\)-th power. Applying the inverse to both sides is valid because both sides are positive.

After solving, substitute the candidate back into the original logarithm. For \(b = 9\), \(\log_9(3) = 1/2\) because \(9^{1/2} = 3\). The check confirms both the value and the direction of the exponent.
