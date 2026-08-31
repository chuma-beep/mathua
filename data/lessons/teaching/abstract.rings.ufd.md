# Unique Factorization Domains

**UFD:** An integral domain where every non-zero non-unit factors into irreducibles uniquely up to order and units; irreducible $\Rightarrow$ prime.

## Factorization Hierarchy

### PID ⇒ UFD
Every PID is a UFD: $ \mathbb Z$ and $F[x]$ are UFDs. The converse fails: $\mathbb Z[\sqrt{-5}]$ has $6=2\cdot3=(1+\sqrt{-5})(1-\sqrt{-5})$ so not a UFD.

### Polynomial Rings over UFDs
Gauss's lemma: if $R$ is a UFD then $R[x]$ is a UFD. Hence $\mathbb Z[x]$ is a UFD even though $\mathbb Z$ has non-principal ideals.

## Example

$\mathbb Z$: $12=2^{2}\cdot3$ uniquely. $\mathbb Z[\sqrt{-5}]$: $6$ has two factorizations, so not a UFD.
