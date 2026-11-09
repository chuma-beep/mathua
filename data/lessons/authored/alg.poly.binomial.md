# The binomial theorem

The binomial theorem expands a power of a sum by counting how many times each term can be chosen. It turns a long multiplication into a formula.

## Counting the terms

In \((x+a)^n\), each term has the form \(\binom{n}{k}x^k a^{n-k}\). The binomial coefficient counts the ways to choose \(k\) factors of \(x\) from \(n\) factors.

There are \(n+1\) terms, and the coefficients are symmetric because \(\binom{n}{k}=\binom{n}{n-k}\).

## The constant term and the sum

The constant term is the \(x^0\) term, so it equals \(a^n\). Nothing depends on the other coefficients.

Setting \(x=1\) makes every term equal to a coefficient, so the sum of all coefficients is \((1+a)^n\).

## Reading one coefficient

To find the coefficient of \(x^k\), use \(\binom{n}{k}a^{n-k}\). Compute the binomial coefficient first, then multiply by the power of \(a\).

A common mistake is to use \(a^k\) instead of \(a^{n-k}\). The exponents in each term must add to \(n\).

For a small power such as \((x+2)^4\), writing the five terms and checking that the coefficients are 1, 8, 24, 32 and 16 builds confidence in the formula.
