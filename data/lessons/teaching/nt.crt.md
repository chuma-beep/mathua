# Chinese Remainder Theorem

**Chinese remainder theorem (CRT):** If $m_1$ and $m_2$ are coprime, the system

$$x \equiv r_1 \pmod{m_1}, \quad x \equiv r_2 \pmod{m_2}$$

has a unique solution modulo $M=m_1m_2$.

## Constructing the Solution

### Coprime Moduli
The hypothesis $\gcd(m_1,m_2)=1$ is essential. It guarantees Bézout coefficients $u,v$ with $um_1+vm_2=1$, which build the solution.

### Combining Congruences
Write $M_1=M/m_1$, $M_2=M/m_2$. Let $y_1$ be the inverse of $M_1$ mod $m_1$ and $y_2$ the inverse of $M_2$ mod $m_2$. Then

$$x = r_1M_1y_1 + r_2M_2y_2 \pmod{M}$$

## Example

Solve $x \equiv 2 \pmod{3}$ and $x \equiv 3 \pmod{5}$.
$M=15$, $M_1=5$, $M_2=3$. Inverses: $5^{-1}\equiv 2 \pmod{3}$, $3^{-1}\equiv 2 \pmod{5}$.
$x = 2\cdot5\cdot2 + 3\cdot3\cdot2 = 20+18=38 \equiv 8 \pmod{15}$. Check: $8\equiv2\pmod3$, $8\equiv3\pmod5$.
