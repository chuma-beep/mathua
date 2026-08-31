# Prime Number Theorem

**Prime number theorem (PNT):** Let $\pi(x)$ count primes $\le x$. Then

$$\pi(x) \sim \frac{x}{\log x}\quad (x\to\infty)$$

Equivalently $p_n\sim n\log n$ for the $n$-th prime, and density near $x$ is $\approx1/\log x$.

## Counting Primes

### Asymptotic Density
$\pi(10)=4$ ($2,3,5,7$); $10/\log10\approx4.3$. Larger $x$ matches more closely; error $O(x\exp(-c\sqrt{\log x}))$ unconditionally, $O(\sqrt{x}\log x)$ under RH.

### Use
$n$-th prime heuristics, cryptography key-size estimates, and Möbius sums $\sum_{n\le x}\mu(n)=o(x)$ equivalent to PNT.

## Example

Estimate $\pi(100)$: $100/\log100\approx100/4.605\approx21.7$; actual $\pi(100)=25$.
