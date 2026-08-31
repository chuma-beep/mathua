# Contour Integration

**Contour integral:** $\int_C f(z)dz$ along curve $C$. For analytic $f$ on simply connected domain, $\oint_C f=0$ (Cauchy); in general $\oint_C f=2\pi i \sum \operatorname{Res}$ (residue theorem).

## Evaluating Integrals

### ML Estimate and Deformation
$|\oint_C f|\le M\cdot L$ (sup $|f|$ times length). Deforming $C$ without crossing singularities preserves integral (homotopy).

### Real Integrals
$\int_{-\infty}^{\infty} dx/(x^{2}+1)=\pi$ via semicircle contour enclosing $i$; Cauchy formula $f(a)=\frac1{2\pi i}\oint f(z)/(z-a)dz$ recovers values.

## Example

$\oint_{|z|=1} dz/z = 2\pi i$ (residue $1$ at $0$). Deform small loop around $0$ to unit circle without crossing pole — same integral.
