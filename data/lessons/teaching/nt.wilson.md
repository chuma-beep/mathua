# Wilson's Theorem

**Wilson's theorem:** For a prime $p$,

$$(p-1)! \equiv -1 \pmod{p}$$

Conversely, if $(n-1)! \equiv -1 \pmod{n}$ then $n$ is prime.

## Factorial Modulo a Prime

### The Pairing Argument
Modulo a prime $p$, every $a \in \{1,\dots,p-1\}$ has a unique inverse $a^{-1}$. Only $1$ and $p-1$ are self-inverse. All other elements pair as $a\cdot a^{-1}\equiv1$, leaving $(p-1)!\equiv1\cdot(p-1)\equiv-1$.

### Converse
If $n$ is composite and $n\neq4$, then $(n-1)!\equiv0\pmod{n}$. The case $n=4$ gives $3!=6\equiv2\pmod4\neq-1$.

## Example

Let $p=5$: $(5-1)!=4!=24$. $24\div5=4$ remainder $4$, and $4\equiv-1\pmod5$. So $24\equiv-1\pmod5$.

Let $p=7$: $6!=720$. $720\div7=102$ remainder $6$, and $6\equiv-1\pmod7$.

## Application

Wilson gives a (slow) primality test and is used to prove that $a^{p-1}\equiv1$ can be refined for factorial structure.
