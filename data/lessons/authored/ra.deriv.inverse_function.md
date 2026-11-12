# The inverse function theorem

The inverse function theorem says a differentiable map with invertible derivative is locally invertible, with a differentiable inverse.

## Local invertibility

If \(f\) is continuously differentiable near \(a\) and \(f'(a) \ne 0\), then \(f\) maps a small interval around \(a\) one-to-one onto an interval around \(f(a)\). For \(f(t) = t^3 + t\) at \(a = 1\) the derivative is 4, and near 2 the map is invertible on a small interval.

## The derivative of the inverse

The inverse satisfies \((f^{-1})'(b) = 1/f'(a)\) where \(b = f(a)\). Here \(f(1) = 2\) and \(f'(1) = 4\), so the inverse has derivative \(1/4\) at 2. Solving \(x^3 + x = 2.1\) gives \(x \approx 1.0245\), and \((1.0245 - 1)/0.1 \approx 0.245\), which matches \(1/4\).

## Global consequences

If \(f'\) never vanishes and \(f\) is proper, the local inverses patch into a global one. The exponential map has derivative \(e^t > 0\) everywhere and is a bijection onto the positive reals, with inverse \(\ln\). The theorem is the reason the logarithm is differentiable with derivative \(1/b\).

The theorem turns a nonzero derivative into a local coordinate change, and a patching argument carries the local statement to the whole domain.
