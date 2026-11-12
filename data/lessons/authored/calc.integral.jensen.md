# Jensen's inequality for integrals

## Convexity and averages

Convexity says the graph lies below its chords, so the midpoint value is at most the average of the endpoint values: \(f((a+b)/2) \le (f(a)+f(b))/2\). For \(f(x) = x^2\) with endpoints 0 and 2 the left side is 1 and the right side is 2.

## The integral form

The integral form replaces the midpoint by the mean: \(f\left(\frac{1}{b-a}\int_a^b g\right) \le \frac{1}{b-a}\int_a^b f(g)\). On the interval from 0 to 2 with \(g(x) = x\) the mean is 1, so the left side is 1 while the right side is \(4/3 \approx 1.333\).

## The gap for a square

The gap for a square is the variance. On the interval from 0 to 2 the mean of \(x^2\) is \(4/3\) and the square of the mean is 1, so the gap is \(1/3\). That matches the uniform variance \((b-a)^2/12 = 4/12 = 1/3\). The gap vanishes only when the variable is constant.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
