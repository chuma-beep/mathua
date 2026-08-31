# Order Modulo n

**Order:** For $\gcd(a,m)=1$, $\operatorname{ord}_m(a)$ is the smallest $k>0$ with $a^{k}\equiv1\pmod m$. It always divides $\varphi(m)$ (and $p-1$ when $m=p$ prime).

## Properties

### Divisors of Phi
$\operatorname{ord}_m(a)\mid\varphi(m)$. For $a=2$, $m=5$, powers $2,4,3,1$ cycle length $4=\varphi(5)$; for $m=7$, $2\to4\to1$ length $3$ dividing $6$.

### Computing the Order
Test divisors of $\varphi(m)$: if $a^{d}\equiv1$ for $d\mid\varphi(m)$, the minimal such $d$ is the order.

## Example

Order of $2$ mod $7$: $2^{1}=2$, $2^{2}=4$, $2^{3}=8\equiv1$, so $\operatorname{ord}_7(2)=3$.
