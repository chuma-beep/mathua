# Green's Functions

**Green's function:** for linear $Ly = f$ with boundary conditions, the impulse response $G(x, \xi)$ satisfying $LG = \delta(x - \xi)$, so that $y(x) = \int G(x, \xi) f(\xi) d\xi$. It is the inverse of the differential operator $L$.

## Worked: the jump across the spike

Integrate $G_{xx} = \delta(x - \xi)$ from $\xi^-$ to $\xi^+$:
1. $G$ itself stays continuous: no jump in value.
2. The slope kinks: $|G_x(\xi^+) - G_x(\xi^-)| = 1$ — the spike bends $G$ by one unit of slope.
3. Away from $\xi$, $G_{xx} = 0$, so $G$ is two straight segments joined at that kink.

So a delta spike bends $G$ by exactly one unit of slope and nothing else.

## Worked: piecewise G on the unit interval

For $y'' = f$ with $y(0) = y(1) = 0$:
1. Straight segments vanishing at the ends give $G(x, \xi) = x(1-\xi)$ for $x < \xi$ and $\xi(1-x)$ for $x > \xi$.
2. Check the kink: slope from the right minus slope from the left is $-\xi - (1-\xi) = -1$, magnitude 1 as required.
3. Test with $f = -1$: $y(x) = -\int_0^1 G \, d\xi = x(x-1)/2$, and indeed $y'' = 1$ with $y(0) = y(1) = 0$.

So $G$ spreads each point load $f(\xi)$ into the tent shape the boundaries allow.

## Symmetry and spectra

Self-adjoint $L$ gives reciprocity $G(x, \xi) = G(\xi, x)$, and eigenfunctions expand it as $G = \sum \phi_n(x)\phi_n(\xi)/\lambda_n$ — the inverse operator written one mode at a time.
