# Completely Regular Spaces

**Completely regular:** For each closed $A$ and point $x$ outside $A$, there is a continuous $f:X\to[0, 1]$ with $f(x)=1$ and $f=0$ on $A$. With $T_1$ this is Tychonoff; normal spaces and metric spaces are completely regular.

## Worked: the metric separator

Separate $x$ from closed $A$ in a metric space with distance $d$:

1. Set $f(y)=d(y, A)/(d(y, A)+d(y, x))$; the denominator is positive since $x$ stays off $A$.
2. On $A$ the numerator vanishes, so $f=0$ there; at $x$ the second term vanishes, so $f(x)=1$.
3. Distance functions are continuous, so $f$ is continuous with values in [0, 1].

So metric spaces are completely regular by an explicit ratio formula.

## Worked: products stay Tychonoff

Check a product of Tychonoff spaces:

1. To separate $x=(x_i)$ from closed $A$, find one coordinate $i$ and a basic neighbourhood missing $A$.
2. The factor space gives a separator $g$ on $X_i$ vanishing off that neighbourhood.
3. Compose with the projection $\pi_i$: $f=g$ after $\pi_i$ separates $x$ from $A$ on the product.

So products of completely regular spaces are completely regular; $\mathbb{R}^n$ is Tychonoff.

## Worked: subspaces of compact Hausdorff

Place a Tychonoff $X$ inside a compact Hausdorff space:

1. For each pair $(x, A)$ take the separating function into [0, 1].
2. Assemble all of them into one map from $X$ into a cube $[0, 1]^J$.
3. The map is an embedding, and its closure is compact Hausdorff.

So Tychonoff spaces are exactly the subspaces of compact Hausdorff spaces: complete regularity is the embeddability condition.
