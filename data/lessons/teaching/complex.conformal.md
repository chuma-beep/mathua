# Conformal Mappings

**Conformal:** An analytic $f$ with $f'(z_0) \ne 0$ preserves angles and orientation at $z_0$. Where the derivative vanishes the map folds angles instead; Mobius transformations are conformal everywhere on the sphere.

## Worked: e to the z is conformal everywhere

Check $f(z)=e^z$:
1. Differentiate: $f'(z)=e^z$.
2. Test for zeros: $|e^z|=e^x > 0$ for every $z$, so $f'$ never vanishes.
3. Conclude $f$ is conformal at every point: vertical lines go to circles and horizontal lines to rays, meeting at the same angles.

So nonzero derivative is the whole test: analytic plus $f' \ne 0$ preserves the angles between curves.

## Worked: z squared doubles angles at zero

Track two rays under $w=z^2$ at the origin:
1. Take rays at angles $\theta$ and $\theta+\phi$ leaving 0.
2. Square: their images leave 0 at angles $2\theta$ and $2\theta+2\phi$.
3. Compare: the opening angle $\phi$ becomes $2\phi$.

So $w=z^2$ is not conformal at 0 where $w'=0$: critical points multiply angles by the vanishing order, and conformality holds only away from them.

## Worked: Mobius maps are conformal on the sphere

Check $w=(az+b)/(cz+d)$ with $ad-bc \ne 0$:
1. Differentiate by the quotient rule: $w'=(ad-bc)/(cz+d)^2$.
2. Test for zeros: the numerator $ad-bc$ is a nonzero constant, so $w'$ never vanishes.
3. Conclude conformality at every point, including $\infty$ after the chart change.

So every Mobius transformation is a conformal automorphism of the sphere; they are exactly the bijective conformal self-maps.
