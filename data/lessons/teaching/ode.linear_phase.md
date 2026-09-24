# Linear Systems and Phase Portraits

**Linear system:** $x' = Ax$ with $A$ a constant 2-by-2 matrix. The eigenvalues of $A$ classify the equilibrium at the origin — read them, and the whole portrait follows.

## Worked: a saddle from diag(2, -3)

Take $A$ with eigenvalues 2 and $-3$:
1. Opposite signs: one direction grows, one shrinks — saddle.
2. Along the $-3$ direction trajectories approach the origin; along the 2 direction they depart.
3. Off-axis starts ride hyperbolic curves that hug the shrinking axis inbound and swing out along the growing axis.

So a saddle is unstable: almost every start eventually leaves along the growing direction.

## Worked: a node and a spiral

Take eigenvalues $-1, -2$:
1. Both real and negative: every direction shrinks — a stable node, all trajectories sink to the origin.
2. Now take $0.5 \pm 2i$: the imaginary part means rotation, the positive real part 0.5 means growth.
3. So that pair is an unstable spiral — circles winding outward.

So the recipe is: real pair, check signs (node vs saddle); complex pair, check the real part (spiral in, out, or centre when it is zero).

## Eigenvalues decide the portrait

Both negative: stable node. Both positive: unstable node. Split signs: saddle. Complex $\alpha \pm i\beta$: spiral, stable exactly when $\alpha < 0$, centre when $\alpha = 0$. A zero eigenvalue gives a whole line of equilibria instead of an isolated point.
