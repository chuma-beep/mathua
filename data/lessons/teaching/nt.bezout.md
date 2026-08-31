# Bézout's Identity

**Bézout's identity:** For integers $a$ and $b$ with $\gcd(a,b)=d$, there exist integers $x$ and $y$ such that

$$ax + by = d$$

Moreover, the set $\{ax+by : x,y \in \mathbb{Z}\}$ is exactly the set of multiples of $d$.

## The GCD as a Linear Combination

### Finding the Coefficients
The extended Euclidean algorithm finds $x$ and $y$ alongside $d$.
For $a=12$ and $b=8$: $\gcd(12,8)=4$. Indeed $12(1)+8(-1)=4$. Also $12(-1)+8(2)=4$.

### Smallest Positive Combination
The smallest positive value of $ax+by$ is $d$. Any larger combination is a multiple of $d$.

## Example

Find integers $x,y$ with $27x+18y=\gcd(27,18)=9$.
Since $27=18\cdot1+9$ and $18=9\cdot2$, working backwards: $9=27-18\cdot1$, so $x=1$, $y=-1$.

## Application: Inverses

If $\gcd(a,m)=1$, Bézout gives $ax+my=1$, so $ax \equiv 1 \pmod{m}$ and $x$ is the modular inverse of $a$.
