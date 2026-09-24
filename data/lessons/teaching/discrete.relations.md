# Relations, Equivalence and Partitions

**Relation:** A subset $R\subseteq A\times A$. It is an **equivalence relation** when reflexive ($(a,a)\in R$ for all $a$), symmetric ($(a,b)\in R\Rightarrow(b,a)\in R$), and transitive ($(a,b),(b,c)\in R\Rightarrow(a,c)\in R$). Each $a$ then defines a class $[a]=\{b:(a,b)\in R\}$, and distinct classes are disjoint and cover $A$: a partition.

## Worked: one property can fail

Test $R=\{(1,1),(2,2),(1,2)\}$ on $\{1,2\}$:
1. Reflexive: both loops $(1,1)$ and $(2,2)$ are present, so yes.
2. Symmetric: $(1,2)$ is present but $(2,1)$ is missing, so no.
3. One failure is enough: $R$ is not an equivalence relation and defines no partition.

So all three properties must be checked independently; two out of three buys nothing.

## Worked: congruence mod 3

Test $R=\{(a,b):a\equiv b\pmod 3\}$ on $\mathbb Z$:
1. Reflexive: $a-a=0$ is divisible by 3 for every $a$.
2. Symmetric: if $a-b=3k$ then $b-a=3(-k)$, still divisible by 3.
3. Transitive: $a-b=3k$ and $b-c=3\ell$ give $a-c=3(k+\ell)$.

So congruence mod 3 is an equivalence relation with classes $[0],[1],[2]$, and mod $n$ partitions $\mathbb Z$ into $n$ residue classes generally.

## Counting partitions with Bell numbers

Equivalent elements clump into blocks: each block is one class, blocks never overlap, and every element sits in exactly one. On $\{1,2,3\}$ the partitions are $\{123\}$, three of type $\{12|3\}$, and $\{1|2|3\}$: 5 total. So the number of equivalence relations on a 3-element set is the Bell number $B_3=5$, and $B_n$ counts them for $n$ generally.
