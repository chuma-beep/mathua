# Cauchy-Riemann Equations

**Cauchy-Riemann equations:** Write $f(z)=u(x,y)+iv(x,y)$ with $z=x+iy$. If $f$ is analytic at $z$, then

$$u_{x}=v_{y}, \qquad u_{y}=-v_{x}$$

Conversely, if $u,v$ have continuous partials and satisfy CR at a point, $f$ is complex-differentiable there.

## Analyticity Criterion

### The Equations
$u_{x}=v_{y}$ ties the $x$-derivative of the real part to the $y$-derivative of the imaginary part; $u_{y}=-v_{x}$ ties the cross terms.

### Testing Examples
$f(z)=z^{2}$: $u=x^{2}-y^{2}$, $v=2xy$, $u_{x}=2x=v_{y}$, $u_{y}=-2y=-v_{x}$ everywhere → entire. $f(z)=\bar{z}$: $u=x$, $v=-y$, $u_{x}=1\neq-1=v_{y}$ → nowhere analytic.

## Example

Check $f(z)=e^{x}(\cos y+i\sin y)$: $u=e^{x}\cos y$, $v=e^{x}\sin y$, $u_{x}=e^{x}\cos y=v_{y}$, $u_{y}=-e^{x}\sin y=-v_{x}$. CR holds everywhere, so $f(z)=e^{z}$ is entire.
