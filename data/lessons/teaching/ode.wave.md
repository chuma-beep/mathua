# Wave Equation and d'Alembert Solution

**Wave equation:** $u_{tt}=c^{2}u_{xx}$ on $\mathbb R×[0,\infty)$. **d'Alembert:** $u(x,t)=f(x-ct)+g(x+ct)$ — superposition of right- and left-travelling waves.

## Characteristics and Energy

### Travelling Waves
Characteristics $x\pm ct=\text{const}$ carry data. Initial $u(x,0)=f(x)$, $u_t(x,0)=0$ gives $u(x,t)=[f(x-ct)+f(x+ct)]/2$.

### Domain of Dependence
$u(x,t)$ depends on $f$ on $[x-ct,x+ct]$; disturbances propagate at speed $c$ with conserved energy $\int(u_t^{2}+c^{2}u_x^{2})dx$.

## Example

$f(x)=\sin x$: $u(x,t)=\sin(x-ct)+\sin(x+ct)$-type superposition; check via d'Alembert produces standing wave $\sin x\cos ct$.
