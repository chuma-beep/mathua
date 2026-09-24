# Euler's Identity

**Euler's identity:** The special case of Euler's formula $e^{i\theta}=\cos\theta+i\sin\theta$ at $\theta=\pi$ is $e^{i\pi}+1=0$, linking 0, 1, $i$, $\pi$, and $e$ in one equation.

## Worked: e to the i pi is negative one

Evaluate $e^{i\pi}$ from the formula:
1. Split into parts: $e^{i\pi}=\cos\pi+i\sin\pi$.
2. Read the trig values: $\cos\pi=-1$ and $\sin\pi=0$.
3. Recombine: $e^{i\pi}=-1+0i=-1$, so $e^{i\pi}+1=0$.

So rotating by a half turn lands exactly on $-1$: the identity is the half-turn case of Euler's formula.

## Worked: e to the i pi over two is i

Evaluate $e^{i\pi/2}$ the same way:
1. Split: $e^{i\pi/2}=\cos(\pi/2)+i\sin(\pi/2)$.
2. Read the values: $\cos(\pi/2)=0$ and $\sin(\pi/2)=1$.
3. Recombine: $e^{i\pi/2}=0+1i=i$.

So a quarter turn lands on $i$: each special angle puts $e^{i\theta}$ at the matching point of the unit circle.

## Worked: e to the two i pi is one

Evaluate $e^{2i\pi}$:
1. Split: $e^{2i\pi}=\cos(2\pi)+i\sin(2\pi)$.
2. Read the values: $\cos(2\pi)=1$ and $\sin(2\pi)=0$.
3. Recombine: $e^{2i\pi}=1+0i=1$.

So a full turn returns to 1: $e^{i\theta}$ is $2\pi$-periodic, parametrizing the unit circle exactly once per period.
