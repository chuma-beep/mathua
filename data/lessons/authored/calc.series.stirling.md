# Stirling's approximation

## The approximation

Stirling's approximation is \(n! \approx \sqrt{2\pi n}\,(n/e)^n\). At \(n = 10\) the factors are 7.9267 and 453999, giving about 3,598,700 against the exact 3,628,800. The relative error is under one percent and falls as \(n\) grows, so the formula is a tool for large \(n\).

## Where the constants come from

The constant \(\sqrt{2\pi}\) comes from the central binomial coefficient, where the Wallis product supplies the factor \(\pi\). The \(1/e\) comes from comparing consecutive factorials: the ratio \((n+1)!/n! = n+1\) must match the ratio of the approximation, which forces the base \(e\). The first correction is a factor \(1 + 1/(12n)\), about 0.0083 at \(n = 10\).

## Counting digits

Digit counts follow from logarithms. The number of digits of \(n!\) is one more than the floor of \(\log_{10}(n!)\). Stirling gives \(n\log_{10}(n/e) + \frac{1}{2}\log_{10}(2\pi n)\). At \(n = 100\) this is \(156.57 + 1.40 = 157.97\), so \(100!\) has 158 digits. The exact count agrees.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
