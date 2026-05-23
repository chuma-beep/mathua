# Fermat's Little Theorem

**Fermat's Little Theorem:** If $p$ is a prime and $a$ is an integer not divisible by $p$, then

$$a^{p-1} \equiv 1 \pmod{p}$$

Equivalently, for any integer $a$:
$$a^p \equiv a \pmod{p}$$

**Example 1:** Let $p = 7$ and $a = 2$.
$2^{6} = 64$, and $64 \div 7 = 9$ remainder $1$, so $64 \equiv 1 \pmod{7}$.

**Example 2:** Let $p = 11$ and $a = 3$.
$3^{10} = 59049$. Divide by 11: $59049 \div 11 = 5368 \times 11 + 1$, remainder $1$.
So $3^{10} \equiv 1 \pmod{11}$.

**Application: Modular inverse**
To find $a^{-1} \pmod{p}$, compute $a^{p-2} \pmod{p}$.
Since $a^{p-1} \equiv 1$, dividing by $a$ gives $a^{p-2} \equiv a^{-1}$.

**Example 3:** Find $5^{-1} \pmod{7}$.
$5^{5} = 5^4 \times 5 = 625 \times 5 = 3125$
$3125 \div 7 = 446 \times 7 + 3$, so $5^5 \equiv 3 \pmod{7}$.
Check: $5 \times 3 = 15 \equiv 1 \pmod{7}$. ✓

**Application: Primality testing**
If $a^{n-1} \not\equiv 1 \pmod{n}$ for some $a$, then $n$ is composite.
(But the converse is not true — Carmichael numbers like 561 are composite yet satisfy the condition for all $a$ coprime to $n$.)
