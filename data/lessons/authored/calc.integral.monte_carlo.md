# Monte Carlo quadrature

Monte Carlo integration estimates an integral from random samples. It is useful when the integrand is complicated, high dimensional or known only through a simulation.

## The sample-mean estimator

For \(N\) uniform samples \(x_1,\dots,x_N\) on \([a,b]\), the estimate is
\[
\int_a^b f(x)\,dx\approx\frac{b-a}{N}\sum_{i=1}^{N}f(x_i).
\]
The factor \(b-a\) converts an average height into an area. The estimate is unbiased, so its expected value is the true integral.

The error decreases like \(1/\sqrt{N}\). Four times as many samples halve the error, which is slower than a good deterministic rule in one dimension and much better than a grid in high dimensions.

## Estimating pi

A classic application estimates \(\pi\). Sample points uniformly in the unit square and count how many fall inside the quarter circle \(x^2+y^2\le1\). If \(m\) of \(N\) points fall inside, the area of the quarter circle is approximately \(m/N\), and
\[
\pi\approx\frac{4m}{N}.
\]
The estimate is a sampled number, not an exact value. Reporting it to three decimal places is honest only when the sample size is stated.

Larger samples reduce the spread, but the result still depends on the particular random draw. Two runs with the same \(N\) can differ in the last digits.

## Standard error

The sample mean of the values \(f(x_i)\) has standard error
\[
\frac{s}{\sqrt{N}},
\]
where \(s\) is the sample standard deviation. This quantity measures the typical distance between the estimate and the true integral.

The formula explains the slow convergence. To gain one more decimal digit of accuracy, the sample size must grow by a factor of about one hundred.
