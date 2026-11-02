# Termwise Differentiation and Integration

Inside the interval of convergence, a power series can be differentiated and integrated one term at a time. The new series keeps the same radius of convergence.

## Differentiating a power series

If \(f(x)=\sum a_n x^n\), then \(f'(x)=\sum n a_n x^{n-1}\) inside the interval of convergence. The constant term disappears, and the coefficient of \(x^{n-1}\) becomes \(n a_n\).

For \(f(x)=1+x+2x^2+3x^3\), the derivative is \(1+4x+9x^2\). The coefficient of \(x^2\) in the derivative is 9, which is three times the coefficient of \(x^3\) in the original series.

## Integrating a power series

Termwise integration gives \(\int f(x)\,dx=C+\sum a_n x^{n+1}/(n+1)\). The coefficient of \(x^{n+1}\) is \(a_n/(n+1)\). The constant of integration is chosen by the initial condition.

For \(f(x)=1+2x+3x^2\), the antiderivative with value zero at the origin is \(x+x^2+x^3\). The coefficient of \(x^3\) is one, which is the original coefficient 3 divided by 3.

## When termwise operations are valid

Termwise differentiation and integration are valid inside the interval of convergence because the series converges uniformly on every closed subinterval there. The radius of convergence does not change under either operation.

Endpoints require separate checking. A series may converge at an endpoint while its derivative diverges there, or the integrated series may converge at a point where the original did not. Always test endpoints with the original and the new series.
