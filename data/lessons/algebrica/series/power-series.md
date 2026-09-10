> Content sourced from [Algebrica](https://algebrica.org/power-series/) — CC BY-NC 4.0

# Power Series

A **power series** centered at \(x = a\) is an infinite polynomial:
\(\sum_{n=0}^{\infty} c_n (x-a)^n\). Near the center it behaves like a
polynomial; far away it may diverge. The set where it converges is always an
interval (possibly a point, possibly the whole line), and inside that interval
the series defines an infinitely differentiable function.

## Radius of Convergence

Every power series has a **radius of convergence** \(R \in [0, \infty]\): the
series converges absolutely for \(|x - a| < R\) and diverges for
\(|x - a| > R\). At the endpoints \(x = a \pm R\) anything can happen and each
endpoint is tested separately. The ratio test gives the standard computation:
\(R = \lim |c_n / c_{n+1}|\) when the limit exists.

Example: \(\sum x^n / n!\) has \(R = \infty\) (converges to \(e^x\) everywhere),
while \(\sum n! \, x^n\) has \(R = 0\) (converges only at \(x = 0\)).

## Differentiating Term by Term

Inside the interval of convergence, a power series can be differentiated and
integrated **term by term**, and the resulting series keeps the same radius
\(R\). So if \(f(x) = \sum c_n (x-a)^n\), then
\(f'(x) = \sum n c_n (x-a)^{n-1}\) on the same interval. Endpoint behavior may
change: differentiating \(\sum x^n/n^2\) (convergent at both endpoints) gives
\(\sum x^{n-1}/n\), which diverges at \(x = 1\).

## Taylor Series Connection

The coefficients of a convergent power series are forced:
\(c_n = f^{(n)}(a)/n!\). Every power series is the Taylor series of the
function it defines. Conversely, an infinitely differentiable function need
not equal its Taylor series (the classic \(e^{-1/x^2}\) is flat at 0 but
nonzero nearby) — equality requires remainder estimates, not just existence
of derivatives.
