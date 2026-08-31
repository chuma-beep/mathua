# First-Order PDEs and Characteristics

**First-order PDE:** $a(x,y)u_x+b(x,y)u_y=c(x,y,u)$. **Characteristics:** curves $x(t),y(t)$ with $dx/dt=a$, $dy/dt=b$ along which PDE becomes ODE $du/dt=c$.

## Transport Along Curves

### Char ODEs
$a u_x+b u_y=0$ has characteristics $dx/a=dy/b$; $u$ constant along them. For $u_t+c u_x=0$, characteristics $x-ct=\text{const}$, so $u(x,t)=f(x-ct)$ shifts initial data.

### Shock Formation
When characteristics intersect, solution becomes multi-valued → shock. Quasi-linear $u$ feeds back into $a,b$, bending chars.

## Example

$u_t+2u_x=0$, $u(x,0)=f(x)$: char $x-2t=\text{const}$, solution $u(x,t)=f(x-2t)$, wave moves right at speed $2$.
