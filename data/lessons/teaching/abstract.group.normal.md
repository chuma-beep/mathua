# Normal Subgroups

**Normal subgroup:** $N \le G$ is normal if $gNg^{-1} = N$ for every $g$: conjugation never leaves the set. Normality is exactly what quotients require — only normal subgroups admit a well-defined coset multiplication.

## Worked: A_3 is normal in S_3

Check conjugation stays inside $A_3 = \{e, (1\,2\,3), (1\,3\,2)\}$:
1. Conjugating a 3-cycle by any permutation gives a 3-cycle (conjugation preserves cycle type).
2. Both 3-cycles already lie in $A_3$, and $geg^{-1} = e$.
3. So $gA_3g^{-1} = A_3$ for every $g$: normal. (Index 2 always forces this.)

So normality can follow from a counting argument rather than checking every conjugate.

## Worked: {e, (1 2)} is not normal in S_3

Conjugate $(1\,2)$ by $(1\,3)$:
1. Track 2 through $(1\,3)(1\,2)(1\,3)$, right to left: $(1\,3)$ fixes 2, $(1\,2)$ sends $2 \to 1$, $(1\,3)$ sends $1 \to 3$. So $2 \to 3$; similarly $3 \to 2$. The conjugate is $(2\,3)$.
2. $(2\,3)$ is not in $\{e, (1\,2)\}$: conjugation escaped the set.
3. So the subgroup is not normal — and indeed $S_3/\{e,(1\,2)\}$ has no group structure.

So one escaped conjugate kills normality, and with it any hope of quotienting.

## Abelian shortcut

Every subgroup of an abelian group is normal, since $ghg^{-1} = h$ always. Non-abelian groups are where normality becomes a real question.
