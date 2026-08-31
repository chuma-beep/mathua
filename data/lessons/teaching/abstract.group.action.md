# Group Actions and Orbit-Stabilizer

**Group action:** A homomorphism $G\to\operatorname{Sym}(X)$, written $g\cdot x$. The orbit $\operatorname{Orb}(x)=\{g\cdot x:g\in G\}$ partitions $X$; the stabilizer $\operatorname{Stab}(x)=\{g:g\cdot x=x\}$ is a subgroup.

## Structure of Actions

### Orbits as Equivalence Classes
$x\sim y$ iff $y=g\cdot x$ for some $g$ is an equivalence relation. Orbits are the classes.

### Orbit-Stabilizer
Bijection $G/\operatorname{Stab}(x)\cong \operatorname{Orb}(x)$ gives $|G|=|\operatorname{Orb}(x)|\cdot|\operatorname{Stab}(x)|$ for finite $G$. Burnside's lemma: $|X/G|=\frac1{|G|}\sum_g|\operatorname{Fix}(g)|$.

## Example

$S_3$ acting on $\{1,2,3\}$: orbit of $1$ is $\{1,2,3\}$ (transitive), stabilizer of $1$ is $\{e,(2\,3)\}$ of order $2$, and $6=3\cdot2$.
