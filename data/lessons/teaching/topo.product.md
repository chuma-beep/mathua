# Product Topology

**Product topology:** For spaces $X_i$, the product carries the coarsest topology making each projection $\pi_j$ continuous. A basic open constrains only finitely many coordinates; the rest stay whole spaces.

## Worked: finite products

Describe a basic open in $[0, 1]\times[0, 1]$:

1. Pick opens $U_1$ in the first factor and $U_2$ in the second, e.g. $U_1=(0, 1/2)$ and $U_2=(1/4, 3/4)$.
2. Form $U_1\times U_2$: a rectangle constraining both coordinates.
3. Every product open is a union of such rectangles, and each projection sends them to opens.

So finite products are simple: basic opens are products of opens, and finite products of compact spaces stay compact.

## Worked: infinite products vs the box

Compare the two topologies on $\mathbb{R}^{\mathbb{N}}$:

1. A product-basic open looks like $U_1\times\cdots\times U_k\times\mathbb{R}\times\cdots$: all but finitely many factors are whole.
2. A box-basic open allows every factor to be proper, e.g. $(-1/n, 1/n)$ in coordinate $n$ for all $n$.
3. That box open contains no product-basic open around the zero sequence, so the box is strictly finer.

So on infinite products the box topology has genuinely more opens; only finitely many coordinates may be constrained in the product topology.

## Worked: projections are continuous

Check $\pi_3:\prod_i X_i\to X_3$:

1. Take an open $V$ in $X_3$.
2. Its preimage is $X_1\times X_2\times V\times X_4\times\cdots$, constraining one coordinate only.
3. That preimage is basic-open by the finite-constraint rule, so $\pi_3$ is continuous.

So each projection is continuous by design, and $\mathbb{R}^{\mathbb{N}}$ with this topology is even metrizable.
