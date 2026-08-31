# Normal Subgroups and Conjugacy

**Normal subgroup:** $N \le G$ is normal ($N \trianglelefteq G$) if $gNg^{-1}=N$ for all $g\in G$, equivalently $gN=Ng$ for all $g$. This means left and right cosets coincide.

## Recognizing Normality

### Index Two Is Normal
Any subgroup of index $2$ is normal: the two cosets are $N$ and $G\setminus N$ on both sides.

### Conjugacy Test
$N$ is normal iff closed under conjugation: for all $n\in N$ and $g\in G$, $gng^{-1}\in N$. Abelian groups have all subgroups normal; $S_3$ has a non-normal subgroup $\{e,(1\,2)\}$.

## Example

In $S_3$, $A_3=\{e,(1\,2\,3),(1\,3\,2)\}$ is normal (index $2$). But $\{e,(1\,2)\}$ is not: $(1\,3)(1\,2)(1\,3)^{-1}=(2\,3)\notin\{e,(1\,2)\}$.
