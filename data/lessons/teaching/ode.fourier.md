# Fourier Transforms and ODEs

**Fourier transform:** $\hat{y}(\omega) = \int_{-\infty}^{\infty} y(t) e^{-i\omega t} dt$, with $F(y') = i\omega \hat{y}$ and $F(y'') = -\omega^2 \hat{y}$. Differentiation becomes multiplication, so constant-coefficient ODEs become algebra in frequency.

## Worked: derivatives become factors

Check $F(y') = i\omega \hat{y}$:
1. Write $F(y') = \int_{-\infty}^{\infty} y'(t) e^{-i\omega t} dt$.
2. Integrate by parts: boundary terms vanish for decaying $y$, leaving $i\omega \int y e^{-i\omega t} dt$.
3. That integral is $\hat{y}$, so $F(y') = i\omega \hat{y}$; applying twice gives $F(y'') = -\omega^2 \hat{y}$.

So each derivative just multiplies the transform by $i\omega$.

## Worked: an ODE turned into division

Solve $y'' + 4y = f(t)$ in frequency:
1. Transform both sides: $(-\omega^2 + 4)\hat{y} = \hat{f}$.
2. Divide: $\hat{y} = \hat{f} / (4 - \omega^2)$.
3. Invert the result by the Fourier inversion integral to recover $y(t)$.

So the differential operator $d^2/dt^2 + 4$ acts as the number $4 - \omega^2$ on each frequency.

## Convolution and inversion

Fourier also turns convolution into a product, $F(f*g) = F(f)F(g)$, which is how Green's functions enter: the impulse response $\hat{y} = 1/(1 - \omega^2)$ for $y'' + y = \delta$ inverts to ringing $\sin t$ switched on at zero. Transform, divide, invert — that is the whole method.
