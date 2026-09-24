# Power Series and Radius of Convergence

**Power series:** An $f$ analytic at $z_0$ expands as $\sum a_n(z-z_0)^n$ converging in the disk out to the nearest singularity. The radius comes from Hadamard's formula $1/R=\limsup |a_n|^{1/n}$; inside the disk the series differentiates term by term.

## Worked: radius of one over one minus z

Find the radius of $\sum z^n$ at 0:
1. Sum the geometric series: it equals $1/(1-z)$ for $|z| < 1$.
2. Locate the nearest singularity: a pole at $z=1$, distance 1 from 0.
3. Read the radius: $R=1$, and the series diverges for $|z| > 1$.

So the radius is the distance to the nearest singularity: the series converges exactly until it hits the first obstruction.

## Worked: e to the z has infinite radius

Find the radius of $\sum z^n/n!$ at 0:
1. Apply the ratio test: $|a_{n+1}/a_n|=1/(n+1)$, tending to 0.
2. Read the radius: $R=\infty$.
3. Match the function: $e^z$ is entire, with no finite singularity to stop at.

So entire functions have everywhere-convergent Taylor series; finite radius always signals a hidden singularity at exactly that distance.

## Worked: differentiating the geometric series

Differentiate $\sum z^n=1/(1-z)$ term by term for $|z| < 1$:
1. Differentiate inside: $\sum n z^{n-1}$ for $n \ge 1$.
2. Differentiate the sum: the derivative of $1/(1-z)$ is $1/(1-z)^2$.
3. Equate: $\sum n z^{n-1}=1/(1-z)^2$ on the same disk.

So uniform convergence on smaller disks justifies termwise differentiation; analytic continuation then extends identities like this one beyond the original disk through overlapping disks.
