# Free Groups

**Free group:** $F(S)$ on a set $S$ is the group of reduced words in $S$ with no relations beyond cancellation: every group with $|S|$ generators is a quotient of it. Free means maximally unconstrained.

## Worked: F_1 is Z

Reduce words in one generator $a$:
1. Every word collapses to $a^n$ for some integer $n$ (negative $n$ for inverses).
2. The map $a^n \mapsto n$ is a bijection respecting the operation.
3. So $F_1 \cong \mathbb Z$: rank 1 free is just the integers.

So freeness on one generator adds nothing beyond integer arithmetic.

## Worked: F_2 is non-abelian

Compare $aba^{-1}b^{-1}$ against the identity:
1. No cancellation is possible: the word is already reduced and nonempty.
2. In particular $ab \neq ba$: the generators do not commute.
3. So $F_2$ is non-abelian — two free generators already escape commutativity.

So relations are the only source of commutativity; with none imposed, even two generators stay wild.

## Abelianization

Killing all commutators collapses $F_n$ to $\mathbb Z^n$: $F_n/[F_n, F_n] \cong \mathbb Z^n$. The abelianization counts the generators and forgets everything else.
