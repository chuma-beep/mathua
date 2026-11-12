# Bernoulli numbers

## The generating function

The Bernoulli numbers are the coefficients in \(x/(e^x-1) = \sum B_n x^n/n!\). Expanding both factors around zero fixes them one at a time: \(B_0 = 1\) from the constant term, then \(B_1 = -1/2\) from the coefficient of \(x\). The series converges for \(|x| < 2\pi\), the distance to the nearest pole.

## The first Bernoulli numbers

The first values are \(B_0 = 1\), \(B_1 = -1/2\), \(B_2 = 1/6\), \(B_4 = -1/30\) and \(B_6 = 1/42\). The expansion gives \(x/(e^x-1) = 1 - x/2 + x^2/12 - x^4/720 + \cdots\), so \(B_2 = 2!/12 = 1/6\) and \(B_4 = -4!/720 = -1/30\), matching the list.

## Odd terms vanish

Every odd Bernoulli number past the first is zero. The reason is that \(x/(e^x-1) + x/2 = (x/2)\coth(x/2)\) is an even function, so no odd powers survive. The recurrence \(\sum_{k=0}^{n}\binom{n+1}{k}B_k = 0\) confirms it: at \(n = 3\) the equation \(1 - 2 + 1 + 4B_3 = 0\) forces \(B_3 = 0\).

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
