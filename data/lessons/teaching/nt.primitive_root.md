# Primitive Roots

**Primitive root:** $g$ mod $p$ (prime) is a primitive root if $\operatorname{ord}_p(g)=p-1$, i.e. powers $\{g^{0},\dots,g^{p-2}\}$ cover all non-zero residues mod $p$.

## Existence and Counting

### Primes Have Primitive Roots
Every prime $p$ has $\varphi(p-1)$ primitive roots. For $p=5$, $\varphi(4)=2$ primitive roots: $2$ and $3$ (powers of $2$ mod $5$ are $2,4,3,1$).

### Order Characterization
$g$ is primitive iff no smaller $k<p-1$ gives $g^{k}\equiv1$. $2$ mod $7$ has order $3\neq6$, so not primitive; $3$ mod $7$ has order $6$, so primitive.

## Example

Mod $7$: primitive roots are $3$ and $5$ (each generates $\{1,\dots,6\}$); $2$ is not primitive (generates only $\{2,4,1\}$).
