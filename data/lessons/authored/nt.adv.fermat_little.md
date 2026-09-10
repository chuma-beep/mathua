# Fermat's little theorem

Fermat's little theorem is the workhorse identity of modular arithmetic.

## Statement

**Theorem.** Let \\(p\\) be prime and let \\(a\\) be an integer not divisible by \\(p\\). Then

\[
a^{p-1} \equiv 1 \pmod{p}.
\]

Equivalently, with no exception on \\(a\\): \\(a^{p} \equiv a \pmod{p}\\).

## Worked Example

A quick check: take \\(p = 5\\) and \\(a=2\\). Then \\(2^4 = 16 = 3\cdot 5 + 1\\), so indeed \\(16 \equiv 1 \pmod 5\\).

## Proof Sketch

The standard proof multiplies the nonzero residue classes \\(1, 2, \dots, p-1\\) by \\(a\\). Modulo \\(p\\), this just permutes them (no two products collide, and none lands on zero because \\(p\nmid a\\)). Multiplying the two descriptions of the same product gives \\(a^{p-1}\cdot(p-1)! \equiv (p-1)!\\), and cancelling \\((p-1)!\\) — legitimate since it is coprime to \\(p\\) — leaves the theorem.

## Applications

The theorem lets huge exponents collapse: to find \\(3^{100} \bmod 7\\), write \\(100 = 6\cdot 16 + 4\\); then \\(3^6 \equiv 1\\) gives \\(3^{100} \equiv 3^4 = 81 \equiv 4 \pmod 7\\). It also underlies primality testing and, through Euler's generalization, the RSA cryptosystem.
