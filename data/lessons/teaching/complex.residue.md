# Residues and the Residue Theorem

**Residue:** For $f$ with an isolated singularity at $z_0$, $\operatorname{Res}(f,z_0)$ is the coefficient $a_{-1}$ of $(z-z_0)^{-1}$ in its Laurent expansion. The residue theorem says $\oint_C f(z)\,dz=2\pi i$ times the sum of the residues inside $C$.

## Worked: residue of one over z

Find $\operatorname{Res}(1/z,0)$:
1. Write the Laurent series at 0: $f(z)=z^{-1}$ is already a Laurent expansion.
2. Read the coefficient of $z^{-1}$: it is $a_{-1}=1$.
3. Conclude $\operatorname{Res}(f,0)=1$.

So a simple pole shows its residue directly: for $f=g/h$ with a simple zero of $h$ at $z_0$, $\operatorname{Res}(f,z_0)=g(z_0)/h'(z_0)$.

## Worked: residue of one over z squared

Find $\operatorname{Res}(1/z^2,0)$:
1. Write the Laurent series at 0: $f(z)=z^{-2}+0 \cdot z^{-1}$.
2. Read the coefficient of $z^{-1}$: it is $a_{-1}=0$.
3. Conclude $\operatorname{Res}(f,0)=0$, so $\oint_{|z|=1} dz/z^2=0$.

So higher poles can hide a zero residue: only the $z^{-1}$ coefficient survives integration, since every other power has an antiderivative.

## Worked: closing the contour on one over z

Evaluate $\oint_{|z|=1} dz/z$ by residues:
1. List interior singularities: one simple pole at 0 with residue 1.
2. Sum the residues: the total is 1.
3. Multiply by $2\pi i$: the integral is $2\pi i \cdot 1=2\pi i$.

So contour integration becomes algebra: sum the enclosed residues and multiply by $2\pi i$.
