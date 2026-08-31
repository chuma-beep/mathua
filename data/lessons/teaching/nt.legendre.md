# Legendre Symbol

**Legendre symbol:** For odd prime $p$ and integer $a$,

$$\left(\frac{a}{p}\right)=\begin{cases}0 & p\mid a\\ 1 & a\text{ is quadratic residue mod }p\\ -1 & a\text{ is non-residue mod }p\end{cases}$$

## Quadratic Residues

### Definition
$a$ is a residue mod $p$ if $x^{2}\equiv a\pmod p$ has a solution. Exactly half the non-zero residues are residues.

### Euler's Criterion
$\left(\frac{a}{p}\right)\equiv a^{(p-1)/2}\pmod p$. Supplement: $\left(\frac{2}{p}\right)=1\iff p\equiv\pm1\pmod8$.

## Example

Mod $5$: residues are $1^{2}=1,2^{2}=4$, so $\left(\frac{1}{5}\right)=1$, $\left(\frac{4}{5}\right)=1$, $\left(\frac{2}{5}\right)=-1$.
