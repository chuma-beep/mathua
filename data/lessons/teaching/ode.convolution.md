# Convolution and Laplace Transforms

**Convolution:** for functions $f, g$ on $[0, \infty)$, the product-like integral $(f*g)(t) = \int_0^t f(\tau) g(t-\tau) d\tau$. Its Laplace transform is an ordinary product: $L(f*g) = L(f) \cdot L(g)$.

## Worked: one convolved with one

Compute $(1*1)(t)$:
1. Insert $f = g = 1$: $(1*1)(t) = \int_0^t 1 \cdot 1 \, d\tau$.
2. The integrand is 1, so the integral is $t$.
3. Check via Laplace: $L(1) = 1/s$, product $1/s^2$, whose inverse is $t$.

So convolving the constant function with itself grows the ramp $t$.

## Worked: Laplace of a convolution

Take $L(f) = 1/s$ and $L(g) = 1/s$:
1. By the convolution theorem, $L(f*g) = (1/s)(1/s) = 1/s^2$.
2. Invert: $1/s^2$ is the transform of $t$.
3. So $(f*g)(t) = t$ without ever doing the time-domain integral.

So the theorem trades an integral for a multiplication plus a table lookup.

## Worked: one convolved with an exponential

Take $f(t) = 1$ and $g(t) = e^t$:
1. Transforms: $L(f) = 1/s$, $L(g) = 1/(s-1)$, so $L(f*g) = 1/(s(s-1))$.
2. Partial fractions: $1/(s-1) - 1/s$.
3. Invert term by term: $e^t - 1$. Direct check: $\int_0^t e^{\tau} d\tau = e^t - 1$.

So convolution is commutative ($f*g = g*f$) and associative, exactly mirroring multiplication in the transform domain.
