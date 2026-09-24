# Möbius Transformations

**Mobius transformation:** The map $w=(az+b)/(cz+d)$ with $ad-bc \ne 0$ is bijective on the sphere $\hat{\mathbb C}$. These maps form the group $PSL(2,\mathbb C)$: they preserve cross-ratios and send circles and lines to circles and lines.

## Worked: one over z is the swap map

Check $w=1/z$ against the definition:
1. Read the coefficients: $a=0$, $b=1$, $c=1$, $d=0$.
2. Test the determinant: $ad-bc=0-1=-1 \ne 0$.
3. Track the poles: $w$ sends 0 to $\infty$ and $\infty$ to 0.

So $1/z$ is Mobius: inversion plus reflection, swapping inside and outside the unit circle.

## Worked: one over z fixes the unit circle

Apply $w=1/z$ to $|z|=1$:
1. Take moduli: $|w|=1/|z|$.
2. Impose $|z|=1$: then $|w|=1$ as well.
3. Conclude the unit circle maps onto itself, inside exchanging with outside.

So Mobius maps preserve the circle-line family: this one sends the unit circle to itself while swapping 0 with $\infty$.

## Worked: pinning three points with z minus one over z plus one

Study $w=(z-1)/(z+1)$:
1. Plug in $z=1$: $w=0/2=0$.
2. Plug in $z=-1$: the denominator vanishes, so $w=\infty$.
3. Let $|z|$ grow: $w$ tends to 1.

So one Mobius map sends 1 to 0, $-1$ to $\infty$, and $\infty$ to 1; in general any three distinct points go to any other three, which is 3-transitivity.
