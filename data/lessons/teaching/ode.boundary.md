# Boundary Value Problems

**BVP:** $y''=f(x,y,y')$ with conditions at two points $y(a)=\alpha$, $y(b)=\beta$ (vs IVP at one point). May have $0,1,\infty$ solutions; Green's function and Fredholm alternative give solvability.

## Versus IVP

### Existence
Unlike IVP (unique via Picard), BVP may have none or many; shooting method iterates IVP to match boundary.

### Green's Function
Solution $y(x)=\int_a^b G(x,\xi)f(\xi)d\xi$ with $G$ satisfying $LG=\delta$ and BC.

## Example

$y''=0$, $y(0)=0$, $y(1)=1$: unique $y=x$.
