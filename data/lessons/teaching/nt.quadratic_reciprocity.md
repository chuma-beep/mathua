# Quadratic Reciprocity

**Law of quadratic reciprocity:** For odd primes $p\neq q$,

$$\left(\frac{p}{q}\right)\left(\frac{q}{p}\right)=(-1)^{\frac{p-1}{2}\cdot\frac{q-1}{2}}$$

If either $p\equiv1\pmod4$ or $q\equiv1\pmod4$, the symbols are equal; if both $p\equiv q\equiv3\pmod4$, they are opposite.

## Using Reciprocity

### Swapping Primes
Reciprocity swaps numerator and denominator, reducing computation: $\left(\frac{3}{5}\right)$ via $\left(\frac{5}{3}\right)$.

### Significance
It connects residues mod $p$ and residues mod $q$, enabling fast evaluation of Legendre symbols without exhaustive search.

## Example

$\left(\frac{5}{7}\right)$: both primes, $5\equiv1\pmod4$, so $\left(\frac{5}{7}\right)=\left(\frac{7}{5}\right)=\left(\frac{2}{5}\right)=-1$ (since $2$ is non-residue mod $5$).
