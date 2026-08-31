# Convolution and Laplace Transforms

**Convolution:** For functions $f,g$ on $[0,\infty)$,

$$(f*g)(t)=\int_0^t f(\tau)g(t-\tau)\,d\tau$$

The Laplace transform turns convolution into multiplication: $L\{f*g\}=L\{f\}\cdot L\{g\}$.

## The Convolution Theorem

### Definition
The integral defining $f*g$ is commutative ($f*g=g*f$) and associative, matching the algebraic role of multiplication in the transform domain.

### Transform of a Convolution
Instead of integrating a convolution directly, multiply the transforms: if $L\{f\}=F(s)$ and $L\{g\}=G(s)$, then $L\{f*g\}=F(s)G(s)$. Inversion then gives $(f*g)(t)$ without a time-domain integral.

## Example

Let $f(t)=1$ and $g(t)=e^{t}$. $L\{f\}=1/s$, $L\{g\}=1/(s-1)$, so $L\{f*g\}=1/(s(s-1))$. Partial fractions give $1/(s-1)-1/s$, whose inverse is $e^{t}-1$. Indeed $(f*g)(t)=\int_0^t e^{\tau}d\tau=e^{t}-1$.
