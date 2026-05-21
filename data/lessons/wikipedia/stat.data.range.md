> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Range_%28statistics%29) — CC BY-SA 4.0

# Find the range of a data set

In descriptive statistics, the **range** of a set of data is the size or width of the narrowest interval which contains all the data.
It is calculated as the difference between the largest and smallest values (also known as the sample maximum and minimum).
It is expressed in the same units as the data.

The range provides an indication of statistical dispersion. Robust measures of range include the interdecile range and the interquartile range.

## Range of continuous IID random variables
For *n* independent and identically distributed continuous random variables *X*\(_{1}\), *X*\(_{2}\),.. ., *X*\(_{*n*}\) with the cumulative distribution function G(*x*) and a probability density function g(*x*), let T denote the range of them, that is, T= max(*X*\(_{1}\), *X*\(_{2}\),.. ., *X*\(_{*n*}\))-
min(*X*\(_{1}\), *X*\(_{2}\),.. ., *X*\(_{*n*}\)).

### Distribution
The range, T, has the cumulative distribution function
\(F(t)= n \int_{-\infty}^\infty g(x)[G(x+t)-G(x)]^{n-1} \, \text{d}x.\)
Gumbel notes that the "beauty of this formula is completely marred by the facts that, in general, we cannot express *G*(*x* + *t*) by *G*(*x*), and that the numerical integration is lengthy and tiresome."

If the distribution of each *X*\(_{*i*}\) is limited to the right (or left) then the asymptotic distribution of the range is equal to the asymptotic distribution of the largest (smallest) value. For more general distributions the asymptotic distribution can be expressed as a Bessel function.

### Moments
The mean range is given by
\(n \int_0^1 x(G)[G^{n-1}-(1-G)^{n-1}] \,\text{d}G\)
where *x*(*G*) is the inverse function. In the case where each of the *X*\(_{*i*}\) has a standard normal distribution, the mean range is given by
\[
\int_{-\infty}^\infty (1-(1-\Phi(x))^n-\Phi(x)^n ) \,\text{d}x.
\]

### Derivation of the distribution
Please note that the following is an informal derivation of the result. It is a bit loose with the calculation of the probabilities.

Let \(m, M\) denote respectively the min and max of the random variables \(X_1 \dots X_n\).

The event that the range is smaller than \(T\) can be decomposed into smaller events according to:

* the index of the minimum value
* and the value \(x\) of the minimum.

For a given index \(i\) and minimum value \(x\), the probability of the joint event:

# \(X_i\) is the minimum,
# and \(X_i=x\),
# and the range is smaller than \(T\),

is:
\[
g(x) \left[ G(x+T) - G(x) \right]^{n-1}
\]
Summing over the indices and integrating over \(x\) yields the total probability of the event: "the range is smaller than \(T\)" which is exactly the cumulative density function of the range:
\[
F(t) = n \int_{-\infty}^{\infty} g(x) \left[G(t+x)-G(x) \right]^{n-1} \, \text{d}x
\]
which concludes the proof.

## The range in other models
Outside of the IID case with continuous random variables, other cases have explicit formulas. These cases are of marginal interest.

* non-IID continuous random variables.
* Discrete variables supported on \(\mathbb N\). A key difficulty for discrete variables is that the range is discrete. This makes the derivation of the formula require combinatorics.

## Related quantities
The range is a specific example of order statistics. In particular, the range is a linear function of order statistics, which brings it into the scope of L-estimation.
