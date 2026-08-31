# Stability of Equilibria

**Equilibrium:** $x^*$ with $f(x^*)=0$ for $x'=f(x)$. **Stable:** nearby solutions stay nearby; **asymptotically stable:** also converge to $x^*$; **unstable:** some nearby solutions diverge.

## Linearization

### Eigenvalue Test
For $\mathbf{x}'=A\mathbf{x}$, origin is asymptotically stable if all eigenvalues have negative real part, unstable if any positive, centre (pure imaginary) stable but not asymptotically.

### Lyapunov Indirect
Jacobian at equilibrium gives linear approximation; its eigenvalues decide stability when hyperbolic (no zero/ pure imaginary).

## Example

$x'=-x$: $x(t)=x_0e^{-t}\to0$, asymptotically stable. $x'=x$: $x(t)=x_0e^{t}$ diverges, unstable.
