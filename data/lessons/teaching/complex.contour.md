# Contour Integration

**Contour integral:** The integral $\int_C f(z)\,dz$ follows a curve $C$ in the plane. For analytic $f$ on a simply connected domain every closed-loop integral is 0; otherwise the residue theorem converts it into $2\pi i$ times enclosed residues.

## Worked: integral of one over z around the unit circle

Evaluate $\oint_{|z|=1} dz/z$ by parametrization $z=e^{it}$:
1. Substitute: $dz=i e^{it}\,dt$, so $dz/z=i\,dt$.
2. Integrate $t$ from 0 to $2\pi$: the integral is $2\pi i$.
3. Match residues: one interior pole of residue 1 gives $2\pi i \cdot 1$.

So the direct computation agrees with the residue theorem; deforming the loop without crossing the pole keeps the same value.

## Worked: integral of one over z squared is zero

Evaluate $\oint_{|z|=1} dz/z^2$ the same way:
1. Substitute: $dz/z^2=i e^{-it}\,dt$.
2. Integrate $t$ from 0 to $2\pi$: $e^{-it}$ averages to 0.
3. Match residues: $a_{-1}=0$ at the double pole, so the theorem gives 0.

So not every singularity contributes: $1/z^2$ has antiderivative $-1/z$ away from 0, and closed loops of derivatives integrate to 0.

## Worked: ML bound for z squared on the unit circle

Bound $\oint_{|z|=1} z^2\,dz$ without computing it:
1. Take the supremum: $|z^2|=1$ on $|z|=1$, so $M=1$.
2. Measure the loop: the unit circle has length $L=2\pi$.
3. Multiply: the integral is at most $M \cdot L=2\pi$ in modulus (it is exactly 0).

So the ML estimate $|\oint_C f| \le M \cdot L$ gives a quick ceiling: supremum of $|f|$ times the length of the path.
