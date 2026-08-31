# Fourier Transforms and ODEs

**Fourier transform:** $\hat{y}(\omega)=\int_{-\infty}^{\infty}y(t)e^{-i\omega t}dt$, with $F\{y'\}=i\omega\hat{y}$, $F\{y''\}=-\omega^{2}\hat{y}$.

## Solving ODEs in Frequency Domain

### Algebraic Equation
$y''+y=f(t)$ becomes $(-\omega^{2}+1)\hat{y}=\hat{f}$, so $\hat{y}=\hat{f}/(1-\omega^{2})$, then invert via integral.

### Convolution Like Laplace
Fourier also turns convolution to product: $F\{f*g\}=F\{f\}F\{g\}$, useful for Green's functions.

## Example

$y''+y=\delta(t)$: Fourier gives $(1-\omega^{2})\hat{y}=1$, so $\hat{y}=1/(1-\omega^{2})$, inverse yields $y(t)=\sin t\cdot H(t)$ (up to constants).
