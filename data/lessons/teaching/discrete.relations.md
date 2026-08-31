# Relations, Functions and Partitions

**Relation:** A subset $R\subseteq A\times A$. **Equivalence relation:** reflexive, symmetric, transitive — it partitions $A$ into disjoint equivalence classes $[a]=\{b:(a,b)\in R\}$.

## Properties and Partitions

### Reflexive, Symmetric, Transitive
$R$ is reflexive if $(a,a)\in R$ for all $a$; symmetric if $(a,b)\in R\Rightarrow(b,a)\in R$; transitive if $(a,b),(b,c)\in R\Rightarrow(a,c)\in R$.

### Partitions from Equivalence
Each class $[a]$ collects mutually related elements; distinct classes are disjoint and cover $A$. Number of equivalence relations on an $n$-element set equals Bell number $B_n$ ($B_3=5$).

## Example

$R=\{(a,b):a\equiv b\pmod 3\}$ on $\mathbb Z$: $a-a=0$ divisible by $3$ (reflexive), $a-b$ divisible implies $b-a$ divisible (symmetric), and transitivity holds via sums. Classes are $[0],[1],[2]$.
