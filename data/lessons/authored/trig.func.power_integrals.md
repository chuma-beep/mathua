# Integrals of powers of sine and cosine

An integral of a power of sine or cosine is evaluated by reducing the exponent until the remaining integral is elementary.

## The reduction formula

Integration by parts gives

\[
\int \sin^{n}x\,dx = -\frac{1}{n}\sin^{n-1}x\cos x + \frac{n-1}{n}\int \sin^{n-2}x\,dx.
\]

Applied to \(n = 3\), the result is \(-\frac{1}{3}\sin^{2}x\cos x - \frac{2}{3}\cos x\). On the
interval from 0 to \(\pi/2\) this evaluates to \(2/3\), since the cosine terms vanish at the top
limit and contribute \(1 - 1/3\) at the bottom. Each application lowers the exponent by two, so
the reduction trades a power of the sine for a cosine outside the integral.

## Odd and even powers

An odd power yields to the substitution \(u = \cos x\). An even power needs the double-angle
identity instead. For \(\sin^{2}x = (1 - \cos 2x)/2\) the integral from 0 to \(\pi/2\) is
\(\pi/4 \approx 0.7854\); the integral of \(\sin^{4}x\) over the same interval is
\(3\pi/16 \approx 0.5890\).

## Wallis products

Writing \(W_{n}\) for the integral of \(\sin^{n}x\) from 0 to \(\pi/2\), the reduction gives
\(W_{n} = \frac{n-1}{n}W_{n-2}\), with \(W_{0} = \pi/2\) and \(W_{1} = 1\). Then
\(W_{2} = \pi/4 \approx 0.7854\), \(W_{3} = 2/3\) and \(W_{4} = 3\pi/16\). The ratios
\(W_{n}/W_{n-1}\) multiply to Wallis's product for \(\pi\), which converges slowly because each
factor is close to 1.

The parity of the exponent decides which reduction to use.
