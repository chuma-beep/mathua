# Lebesgue Covering Dimension

**Covering dimension:** The smallest $n$ such that every open cover has a refinement with no point in more than $n+1$ sets. Then $\dim\mathbb{R}^n=n$, the Cantor set has dimension 0, and subspaces never exceed the ambient dimension.

## Worked: the line has dimension 1

Show $\dim\mathbb{R}=1$:

1. Cover $\mathbb{R}$ by small overlapping intervals; refine so no point lies in more than 2 of them.
2. Order 1 always suffices: chains of intervals overlap pairwise only.
3. Order 0 fails: a connected cover by disjoint opens is impossible.

So the line needs multiplicity 2, giving $n=1$: dimension counts overlap multiplicity minus one.

## Worked: the Cantor set has dimension 0

Refine any cover of the Cantor set $C$:

1. $C$ is totally disconnected: its points separate by clopen pieces at each finite stage.
2. Refine the given cover by stage-$k$ clopen pieces for large $k$; they are pairwise disjoint.
3. Multiplicity 1 suffices, so order 0 works.

So totally disconnected compacta drop to dimension 0: no overlaps needed at all.

## Worked: dimension is topological

Conclude invariance from the definition:

1. A homeomorphism carries open covers to open covers preserving refinement multiplicity.
2. Hence both spaces admit exactly the same minimal $n$.
3. Brouwer's invariance of domain is the deep case: $\mathbb{R}^n$ and $\mathbb{R}^m$ differ for $n$ differing from $m$.

So dimension is a topological invariant; for separable metric spaces all three dimensions (covering, small and large inductive) coincide, and Menger-Nobeling embeds $n$-dimensional compact metric spaces in $\mathbb{R}^{2n+1}$.
