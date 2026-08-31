# Möbius Function and Inversion

**Möbius function:** For $n=p_1^{e_1}\cdots p_k^{e_k}$,

$$\mu(n)=\begin{cases}1 & n=1\\ 0 & \text{some }e_i\ge2\\ (-1)^k & \text{all }e_i=1\end{cases}$$

Möbius inversion: if $F(n)=\sum_{d\mid n}f(d)$ then $f(n)=\sum_{d\mid n}\mu(d)F(n/d)$.

## Values and Inversion

### Computing μ
$\mu(1)=1$, $\mu(p)=-1$, $\mu(6)=\mu(2\cdot3)=1$, $\mu(4)=\mu(2^{2})=0$ (squared prime divides).

### Inversion
Möbius inverts divisor sums, e.g. $\varphi(n)=n\sum_{d\mid n}\mu(d)/d$ and $\sum_{d\mid n}\mu(d)=0$ for $n>1$.

## Example

$\mu(30)=\mu(2\cdot3\cdot5)=(-1)^{3}=-1$ (three distinct primes). $\mu(12)=\mu(2^{2}\cdot3)=0$ (squared prime $2^{2}\mid12$).
