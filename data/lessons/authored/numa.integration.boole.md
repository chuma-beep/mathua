# Boole’s rule

Boole’s rule integrates over four subintervals using five function values and a fixed set of weights.

## The five-point weights

Boole’s rule uses the weights 7, 32, 12, 32 and 7. The weighted sum is \(7f_0+32f_1+12f_2+32f_3+7f_4\).

The weights sum to 90. The symmetry around the middle value is part of the rule.

The pattern is not the same as Simpson’s rule, which uses three points.

## The integral estimate

With step \(h\), the estimate is \((2h/45)\) times the weighted sum. The factor 2/45 comes from integrating the interpolating quartic.

When \(h=45\), the factor is 2, which makes hand calculations easy. The rule is exact for polynomials up to degree 5.

The rule is not exact for a general degree-6 polynomial, so the error is of order \(h^6\).

## Composite Boole

A composite rule applies Boole’s rule to panels of four subintervals. Adjacent panels share endpoints.

For \(m\) panels, the number of distinct function values is \(4m+1\). Reusing shared endpoints avoids extra evaluations.

The composite error is the sum of the panel errors, which is why smooth panels are preferred.
