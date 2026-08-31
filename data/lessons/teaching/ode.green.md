# Green's Functions

**Green's function:** For linear $Ly=f$ with boundary conditions, $G(x,\xi)$ satisfies $LG=\delta(x-\xi)$ and $y(x)=\int G(x,\xi)f(\xi)d\xi$ — the inverse of $L$.

## Impulse Response

### Construction
For $y''=f$, $y(0)=y(1)=0$, $G(x,\xi)=\begin{cases}x(1-\xi)&x<\xi\\\xi(1-x)&x>\xi\end{cases}$, with jump $G'_x(\xi^+)-G'_x(\xi^-)=1$.

### Symmetry and Expansion
Self-adjoint $L$ gives $G(x,\xi)=G(\xi,x)$; eigenfunction expansion $G(x,\xi)=\sum\phi_n(x)\phi_n(\xi)/\lambda_n$.

## Example

$y''=1$, $y(0)=y(1)=0$: $y(x)=\int_0^1 G(x,\xi)d\xi = x(x-1)/2$, check $y''=1$.
