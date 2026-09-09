# The Ratio Test

## Ratio test

For \\(\sum a_{n}\\) with \\(a_{n} \neq 0\\), compute

\\[L = \lim_{n\to\infty}\left|\frac{a_{n+1}}{a_{n}}\right|\\]

Then \\(L < 1\\) forces absolute convergence, \\(L > 1\\) (including \\(\infty\\)) forces divergence, and \\(L = 1\\) is inconclusive — the test says nothing and a comparison, integral, or root test is needed.

### Example

For \\(\sum 1/n!\\): \\(|a_{n+1}/a_{n}| = 1/(n+1) \to 0 < 1\\), so the series converges. For \\(\sum n!/n^{n}\\): the ratio tends to \\(1/e < 1\\), so it converges too.

The test shines on factorials and exponentials, where ratios simplify; it is useless for \\(p\\)-series and rational terms, where \\(L = 1\\) always.
