# Schwarz Lemma

**Schwarz lemma:** If $f$ is analytic on the unit disk $\mathbb D$ with $f(0)=0$ and $|f(z)| \le 1$, then $|f(z)| \le |z|$ and $|f'(0)| \le 1$. Equality at any nonzero point forces the rotation $f(z)=e^{i\theta}z$.

## Worked: z squared obeys the bound strictly

Test $f(z)=z^2$ on $\mathbb D$:
1. Check hypotheses: $f$ is analytic, $f(0)=0$, and $|z^2| < 1$ on the disk.
2. Compare: $|f(z)|=|z|^2 \le |z|$ since $|z| < 1$.
3. Check the derivative: $f'(0)=0 \le 1$, strict on both counts.

So typical maps sit strictly inside the bound; the lemma caps how fast a normalized bounded map can grow.

## Worked: rotations attain the bound

Test $f(z)=e^{i\theta}z$ on $\mathbb D$:
1. Check hypotheses: analytic, $f(0)=0$, $|f(z)|=|z| \le 1$.
2. Compare: $|f(z)|=|z|$ exactly, at every point.
3. Check the derivative: $|f'(0)|=|e^{i\theta}|=1$.

So rotations are the extremal case: equality anywhere nonzero forces the map to be a rotation, with no other possibilities.

## Worked: bounding a shifted map via Schwarz-Pick

Bound $f(z)=(z-1/2)/(1-z/2)$ at $z=0$ using the invariant form:
1. Note $f$ is a disk automorphism with $f(1/2)=0$.
2. Apply Schwarz-Pick with $z_1=0$ and $z_2=1/2$: both sides reduce to $|f(0)| \le 1/2$.
3. Verify: $f(0)=-1/2$, saturating the bound.

So the invariant form $|f(z_1)-f(z_2)|/|1-\overline{f(z_2)}f(z_1)| \le |z_1-z_2|/|1-\bar{z_2}z_1|$ controls every analytic self-map of the disk, with automorphisms $(z-a)/(1-\bar{a}z)$ as the extremals.
