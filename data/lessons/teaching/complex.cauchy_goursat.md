# Cauchy-Goursat Theorem

**Cauchy-Goursat:** If $f$ is analytic inside and on a simple closed contour $C$, then $\oint_C f(z)\,dz=0$. Goursat's proof needs only analyticity, not continuity of $f'$; locally $f$ then has a path-independent antiderivative.

## Worked: integral of z squared over a triangle is zero

Integrate $f(z)=z^2$ around any triangle $T$:
1. Take the antiderivative $F(z)=z^3/3$, valid everywhere since $z^2$ is entire.
2. Walk the three edges: consecutive endpoint values telescope and cancel.
3. Total 0: the starting and ending values of $F$ coincide.

So entire functions integrate to zero over closed loops; analyticity on and inside the loop supplies the antiderivative locally.

## Worked: one over z breaks the conclusion

Integrate $f(z)=1/z$ around $|z|=1$:
1. Check the hypothesis: $f$ is not analytic at 0, which lies inside the circle.
2. Compute directly: the integral is $2\pi i \ne 0$.
3. Diagnose: the interior singularity voids the theorem.

So the hypothesis "analytic everywhere inside" is sharp: one interior pole changes the integral from 0 to $2\pi i$.

## Worked: Goursat bisection on a triangle

Sketch why only analyticity is needed, for a triangle $T$:
1. Split $T$ into 4 congruent subtriangles; one of them, call it $T_1$, carries at least a quarter of the integral.
2. Repeat: nested triangles $T \supset T_1 \supset T_2 \supset \cdots$ shrink to a point $z_0$ while concentrating the integral.
3. Linearize at $z_0$: analyticity makes $f$ nearly linear there, and linear functions integrate to 0 over triangles, squeezing the integral to 0.

So Goursat's subdivision removes the continuity-of-$f'$ assumption: analyticity plus shrinking geometry already forces the integral to vanish, which then yields local antiderivatives and infinite differentiability.
