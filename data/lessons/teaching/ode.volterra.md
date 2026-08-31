# Volterra Integral Equations

**Volterra of second kind:** $y(x)=f(x)+\int_0^x K(x,t)y(t)dt$ with continuous kernel $K$; unique continuous solution via Picard iteration; resolvent kernel $R$ gives $y=f+\int R f$.

## Solving

### Differentiation
$K=1$: $y=f+\int_0^x y$, $y'=f'+y$, solve ODE.

### Transform
Convolution Volterra solved via Laplace.

## Example

$y(x)=1+\int_0^x y(t)dt$: $y'=y$, $y(0)=1$ ⇒ $y=e^x$.
