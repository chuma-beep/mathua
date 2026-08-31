# Residues and the Residue Theorem

**Residue:** For $f$ with isolated singularity at $z_{0}$, $\operatorname{Res}(f,z_{0})$ is the coefficient $a_{-1}$ of $(z-z_{0})^{-1}$ in its Laurent expansion.

**Residue theorem:** If $f$ is analytic inside and on a positively oriented simple closed contour $C$ except at finitely many interior poles $z_{k}$,

$$\oint_{C} f(z)\,dz = 2\pi i \sum_{k}\operatorname{Res}(f,z_{k})$$

## Computing and Using Residues

### Simple Poles
If $f=g/h$ with $g(z_{0})\neq0$ and $h$ has a simple zero at $z_{0}$, then $\operatorname{Res}(f,z_{0})=g(z_{0})/h'(z_{0})$. For $f(z)=1/z$, $\operatorname{Res}(f,0)=1$.

### Closing Contours
The theorem turns contour integrals into algebra: sum the residues inside $C$ and multiply by $2\pi i$. For $f(z)=1/z^{2}$, the residue at $0$ is $0$, so $\oint 1/z^{2}dz=0$.

## Example

Evaluate $\oint_{|z|=1} \frac{dz}{z}$. One interior pole at $0$ with residue $1$, so the integral is $2\pi i \cdot 1 = 2\pi i$.
