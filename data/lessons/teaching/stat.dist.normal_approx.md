# Normal approximation to the binomial

A binomial count can be approximated by a normal distribution when the expected number of successes and failures are both large enough. The approximation turns a sum of many trials into a familiar normal calculation.

## When a binomial looks normal

If $X$ is binomial with $n$ trials and success probability $p$, then

$$\mu = np, \qquad \sigma^2 = np(1-p).$$

The normal approximation is reasonable when $np$ and $n(1-p)$ are both at least about 10. Small expected counts produce a skewed distribution that a normal curve fits poorly.

## The continuity correction

A binomial count takes whole-number values, while a normal variable is continuous. To approximate $P(X \le k)$, use the boundary $k + 0.5$:

$$z = \frac{k + 0.5 - np}{\sqrt{np(1-p)}}.$$

For $P(X \ge k)$, use $k - 0.5$. The half-unit shift is the continuity correction.

## Accuracy of the approximation

The approximation improves as $n$ grows and as $p$ moves away from 0 and 1. It is weakest in the tails, where exact binomial probabilities can differ noticeably.

When an exact answer is required, compute the binomial sum directly. Use the normal approximation for speed and for the connection to the central limit theorem.
