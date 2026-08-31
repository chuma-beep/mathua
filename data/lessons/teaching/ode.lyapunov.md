# Lyapunov Functions

**Lyapunov function:** $V(x)>0$ for $x\neq0$, $V(0)=0$, $dV/dt = \nabla V·f <0$ along trajectories proves stability of $x^*=0$ without solving ODE (direct method); strict $V$ gives asymptotic.

## Direct Method

### LaSalle
LaSalle's invariance principle extends to $dV/dt\le0$ with invariant set.

### Converse
Asymptotically stable ⇒ exists Lyapunov $V$.

## Example

$x'=-x$: $V=x^2$, $dV/dt=2x·(-x)=-2x^2<0$, so $0$ asymptotically stable.
