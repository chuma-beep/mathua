# Polynomial Rings and Field Extensions

**Polynomial ring:** For a ring $R$, $R[x]$ is the set of polynomials with coefficients in $R$. If $R$ is a field, $R[x]$ is a PID; quotients $R[x]/(p(x))$ give field extensions when $p$ is irreducible.

## From Polynomials to Fields

### Ideals in $R[x]$
$(x)$ consists of polynomials with zero constant term; $R[x]/(x)\cong R$ (evaluate at $0$). $(x^{2}+1)$ in $\mathbb R[x]$ is maximal, so $\mathbb R[x]/(x^{2}+1)\cong\mathbb C$.

### Field Extensions
Adjoining a root: $\mathbb Q[x]/(x^{2}-2)\cong\mathbb Q(\sqrt2)$. Irreducible $p(x)$ of degree $n$ gives an $n$-dimensional extension.

## Example

$\mathbb R[x]/(x^{2}+1)$: class of $x$ satisfies $x^{2}=-1$, so it behaves as $i$. Every element is $a+bx$ with $x^{2}=-1$, matching $a+bi\in\mathbb C$.
