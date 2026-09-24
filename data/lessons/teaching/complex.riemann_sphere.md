# Riemann Sphere and Stereographic Projection

**Riemann sphere:** The one-point compactification $\hat{\mathbb C}=\mathbb C \cup \{\infty\}$ is the sphere $S^2$ via stereographic projection from the north pole: $(X,Y,Z)$ maps to $z=(X+iY)/(1-Z)$, and the north pole itself is $\infty$.

## Worked: one over z swaps zero and infinity

Track $w=1/z$ on the sphere:
1. As $z$ approaches 0, $|w|=1/|z|$ blows up, so $w$ tends to the north pole $\infty$.
2. As $|z|$ grows without bound, $|w|$ tends to 0, the south pole.
3. Read the swap: 0 and $\infty$ exchange places, a rotation of the sphere.

So $\infty$ is an ordinary point of the sphere: behavior there is studied in the chart $w=1/z$, where it becomes behavior at 0.

## Worked: the unit circle lifts to the equator

Lift $|z|=1$ through stereographic projection:
1. Points with $|z|=1$ satisfy $X^2+Y^2=(1-Z)^2$ on the sphere.
2. Combined with $X^2+Y^2+Z^2=1$ this forces $Z=0$.
3. Read the image: the equator $Z=0$, a circle on the sphere.

So stereographic projection sends circles on the sphere to circles or lines in the plane: the equator gives a circle, while any circle through the north pole gives a line.

## Worked: z squared has a double pole at infinity

Study $f(z)=z^2$ at $\infty$ via $w=1/z$:
1. Substitute: $g(w)=f(1/w)=w^{-2}$.
2. Read the singularity: a pole of order 2 at $w=0$.
3. Conclude $z^2$ has a double pole at $\infty$ on the sphere.

So meromorphic functions on the sphere are exactly rational functions: poles of finite order everywhere including $\infty$ leave only quotients of polynomials.
