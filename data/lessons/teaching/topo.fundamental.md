# Fundamental Group Calculations

**Fundamental group:** For a based space $(X, x_0)$, $\pi_1(X)$ collects homotopy classes of loops at $x_0$ under concatenation. $\pi_1(S^1)$ is $\mathbb{Z}$ by winding number, $\pi_1(S^2)$ is trivial, $\pi_1(T^2)$ is $\mathbb{Z}\times\mathbb{Z}$.

## Worked: the torus

Compute $\pi_1$ of $T^2=S^1\times S^1$:

1. A loop in the product is a pair of loops, one in each circle factor.
2. Each factor contributes its winding number, an independent integer.
3. Concatenation adds the pairs componentwise.

So $\pi_1(T^2)$ is $\mathbb{Z}\times\mathbb{Z}$: loops carry two winding numbers, one per factor.

## Worked: the twice-punctured plane

Compute $\pi_1$ of $X=\mathbb{R}^2$ minus $\{0, 1\}$:

1. Loop $a$ around 0 and loop $b$ around 1; neither contracts since its puncture blocks it.
2. Van Kampen on the union of the plane-minus-0 and the plane-minus-1 gives no relation between $a$ and $b$.
3. Hence every loop is a unique word in $a$ and $b$.

So $\pi_1(X)$ is free on 2 generators: punctures contribute independent loops with no relations.

## Worked: detecting simple connectivity

Test whether $S^2$ is simply connected:

1. Take any loop $w$ on $S^2$ and perturb it to miss some point $p$.
2. $S^2$ minus $p$ is a plane, where $w$ contracts to a point.
3. Thus every loop is null-homotopic and $\pi_1(S^2)$ is trivial.

So trivial $\pi_1$ means every loop contracts; $S^1$ fails this while $S^2$ passes.
