# Euler's Method

**Euler's method:** For $y'=f(x,y)$, $y(x_0)=y_0$, step size $h$,

$$y_{n+1}=y_n + h\,f(x_n,y_n), \quad x_{n+1}=x_n+h$$

It follows the tangent line at each step to approximate the solution curve.

## One Step of Euler

### The Tangent Approximation
At $(x_n,y_n)$ the ODE gives slope $f(x_n,y_n)$. Euler moves $h$ in $x$ and $h\cdot f$ in $y$, as if the slope stayed constant over the interval.

### Local Error
Error per step is $O(h^2)$; global error after $1/h$ steps is $O(h)$. Smaller $h$ improves accuracy but costs more steps.

## Example

Use $h=0.1$ for $y'=y$, $y(0)=1$ to estimate $y(0.1)$. $f(0,1)=1$, so $y_1=1+0.1\cdot1=1.1$. The exact $y(0.1)=e^{0.1}\approx1.1052$, error $0.0052$.
