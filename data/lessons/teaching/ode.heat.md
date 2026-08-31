# Heat Equation and Separation of Variables

**Heat equation:** $u_t=k u_{xx}$ on $[0,L]×[0,\infty)$. **Separation:** assume $u=X(x)T(t)$, giving $X''/X=T'/(kT)=-\lambda$.

## Fourier Series Solution

### Separation
$T'=-k\lambda T\Rightarrow T=Ce^{-k\lambda t}$, $X''+\lambda X=0$ with Dirichlet gives $\lambda_n=(n\pi/L)^{2}$, $X_n=\sin(n\pi x/L)$.

### Maximum Principle
$\max u$ on space-time boundary; initial discontinuities smooth instantly for $t>0$.

## Example

$L=\pi$, initial $u(x,0)=\sin x$: $u(x,t)=e^{-kt}\sin x$.
