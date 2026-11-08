# Floating-point representation error

A computer stores numbers in a fixed number of binary digits. Most decimal fractions cannot be represented exactly in that format, so every stored value may already carry a small error before any arithmetic begins. Understanding that error is the first step in trusting a numerical result.

## Numbers a machine can hold

Floating-point numbers are written as a sign, a significand and an exponent. The significand has a fixed number of bits, so only finitely many values exist in any interval. The gap between neighbouring values is the machine epsilon at that scale.

A number such as 0.1 has no finite binary expansion, just as \(1/3\) has no finite decimal expansion. The stored value is the nearest representable number, and the difference is the representation error. It is small, but it is not zero.

## Absolute and relative error

Absolute error is the size of the difference: \(|\hat x-x|\). Relative error divides by the true value: \(|\hat x-x|/|x|\). Relative error is usually the more informative measure because it accounts for scale.

An absolute error of 0.01 is serious for a length of 0.05 and negligible for a length of 1000. The relative error says which case you are in. Percent error is relative error multiplied by one hundred.

## Significant digits

The number of correct significant digits is roughly the negative logarithm of the relative error. A relative error of \(10^{-3}\) leaves about three correct digits. Rounding to a stated number of significant figures keeps the leading digits and replaces the rest with zeros.

Significant figures count from the first nonzero digit, not from the decimal point. Trailing zeros after the decimal point are significant; leading zeros are not. Reporting more digits than the error justifies suggests a precision the computation does not have.
