# Analytic Functions

**Analytic:** $f$ is analytic (holomorphic) on a domain $D$ if it is complex-differentiable at every point of $D$. Equivalently it satisfies Cauchy-Riemann with continuous partials, or is locally given by a convergent power series.

## Worked: z squared is entire

Check $f(z)=z^2$ with $u=x^2-y^2$ and $v=2xy$:
1. Compute the partials: $u_x=2x$, $v_y=2x$, $u_y=-2y$, $v_x=2y$.
2. Verify CR: $u_x=v_y$ and $u_y=-v_x$ at every point.
3. Differentiate directly: $f'(z)=2z$ exists everywhere.

So $z^2$ is analytic on all of $\mathbb C$; in general every polynomial is entire, with derivative from the usual power rule.

## Worked: the conjugate is analytic nowhere

Check $f(z)=\bar{z}$ with $u=x$ and $v=-y$:
1. Compute the partials: $u_x=1$ and $v_y=-1$.
2. Compare: $u_x \ne v_y$, so CR fails at every point.
3. Conclude no complex derivative exists anywhere.

So continuity is far weaker than analyticity: $\bar{z}$ is smooth as a real map but holomorphic nowhere.

## Worked: modulus squared is analytic only at zero

Check $f(z)=|z|^2$ with $u=x^2+y^2$ and $v=0$:
1. Compute the partials: $u_x=2x$, $v_y=0$, $u_y=2y$, $v_x=0$.
2. Impose CR: $2x=0$ and $2y=0$, forcing $x=y=0$.
3. Note the equations hold only at the origin, never on a neighborhood.

So CR at an isolated point does not make $f$ analytic there: analyticity needs differentiability on a whole neighborhood, and indeed a function analytic at every point is infinitely differentiable.
