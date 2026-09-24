# Volterra Integral Equations

**Volterra of the second kind:** $y(x) = f(x) + \int_0^x K(x,t) y(t) dt$ with continuous kernel $K$ — the upper limit moves with $x$, so the past accumulates. Continuous $K$ and $f$ guarantee exactly one continuous solution.

## Worked: differentiating back to an ODE

Solve $y(x) = 1 + \int_0^x 2y(t) dt$:
1. Differentiate: $y' = 2y$ (the integral's upper limit gives the integrand at $x$).
2. Read the initial value from the equation at $x = 0$: $y(0) = 1$.
3. Solve $y' = 2y$, $y(0) = 1$: $y = e^{2x}$, so $y(1) = e^2$.

So convolution-free Volterra equations with constant kernels are ODEs in disguise.

## Worked: Picard iterates converge

Build $y(x) = 1 + \int_0^x y(t) dt$ by iteration from $y_0 = 1$:
1. $y_1 = 1 + \int_0^x 1 \, dt = 1 + x$.
2. $y_2 = 1 + \int_0^x (1 + t) dt = 1 + x + x^2/2$.
3. The pattern is the exponential series: $y_n \to e^x$.

So Picard iteration converges to the unique solution — the same $e^x$ that differentiation gives.

## Resolvent and Laplace routes

Repeated kernels sum to the resolvent $R$ with $y = f + \int R f$ (a Neumann series that always converges for Volterra), and convolution kernels $K(x - t)$ surrender to Laplace transforms. Differentiate when the kernel is simple; iterate or transform otherwise.
