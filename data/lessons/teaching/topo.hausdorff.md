# Hausdorff Spaces

**Hausdorff space:** A space $X$ is Hausdorff (or $T_2$) if any two distinct points $x$ and $y$ have disjoint open neighbourhoods $U$ of $x$ and $V$ of $y$. Hausdorff separation makes limits unique.

## Definition

Test the indiscrete topology on the set $\{a, b\}$:

1. The only opens are the empty set and the whole $\{a, b\}$.
2. Any neighbourhood of $a$ is the whole set, and any neighbourhood of $b$ is the whole set.
3. Those neighbourhoods intersect, so $a$ and $b$ admit no disjoint pair: not Hausdorff.

So the axiom has teeth: with too few opens, distinct points cannot be pulled apart.

## Examples and non-examples

Separate two points in the cofinite topology on $\mathbb{R}$:

1. Take distinct $x$ and $y$; any opens $U$ around $x$ and $V$ around $y$ miss only finitely many points each.
2. The union of the two finite missed sets is finite, so $U$ and $V$ share all but finitely many reals.
3. Hence $U$ meets $V$: separation fails, and the cofinite topology is not Hausdorff.

So metric and discrete spaces pass easily, while indiscrete and cofinite topologies fail for lack of disjoint opens.

## Metric spaces are Hausdorff

Separate distinct $x, y$ in a metric space with distance $d$:

1. Set $r = d(x, y)/2$, a positive radius since $x$ and $y$ differ.
2. Take the balls $B$ around $x$ and $C$ around $y$, both of radius $r$.
3. By the triangle inequality no point sits within $r$ of both, so $B$ and $C$ are disjoint.

So every metric space is Hausdorff: the metric itself builds the separating neighbourhoods.
