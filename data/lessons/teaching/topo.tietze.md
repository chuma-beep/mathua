# Tietze Extension Theorem

**Tietze extension:** If $X$ is normal and $A\subseteq X$ closed, every continuous $f:A\to[0,1]$ (or $\mathbb R$) extends to continuous $F:X\to[0,1]$ with $F|_A=f$. This characterizes normality.

## Extending from Closed Sets

### Closedness Matters
$A$ must be closed; e.g. $A=\mathbb Q\subset\mathbb R$ not closed cannot extend all bounded functions.

### Urysohn as Corollary
Tietze with $A=A_0\cup A_1$ ($A_0,A_1$ disjoint closed) and $f=0$ on $A_0$, $1$ on $A_1$ gives Urysohn's lemma.

## Example

$X=\mathbb R$, $A=\{0,1\}$ closed: $f(0)=0,f(1)=1$ extends to $F(x)=x$ on $[0,1]$ and constant outside, continuous on $\mathbb R$.
