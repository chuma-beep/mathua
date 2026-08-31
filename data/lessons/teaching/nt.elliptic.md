# Elliptic Curves and BSD

**Elliptic curve:** $E: y^2=x^3+ax+b$, $\Delta=-16(4a^3+27b^2)\neq0$; points form abelian group via chord-tangent; Mordell-Weil $E(\mathbb Q)\cong\mathbb Z^r\times\text{torsion}$; Hasse bound $|#E(\mathbb F_p)-(p+1)|\le2\sqrt p$; BSD relates rank $r$ to $L(E,1)$.

## Arithmetic of Curves

### Group Law
Chord through $P,Q$ meets $E$ at $R$, reflect to $P+Q$; used in ECM factorization and ECC cryptography.

### L-Function
$L(E,s)$ at $s=1$ conjecturally rank $0$ iff $L(E,1)\neq0$.

## Example

$E: y^2=x^3-x$, $E(\mathbb Q)$ rank $0$, $L(E,1)\neq0$ fits BSD.
