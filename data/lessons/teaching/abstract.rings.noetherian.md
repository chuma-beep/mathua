# Noetherian Rings

**Noetherian ring:** Every ascending chain of ideals stabilizes — equivalently, every ideal is finitely generated. Noetherian means ideal theory stays finite: no infinite strictly-growing towers.

## Worked: the failing chain in k[x_1, x_2, ...]

Build $(x_1) \subset (x_1, x_2) \subset (x_1, x_2, x_3) \subset \cdots$:
1. Each step adds a fresh variable no earlier ideal contains: $x_{n+1} \notin (x_1, \dots, x_n)$.
2. So every inclusion is strict, forever.
3. Hence $k[x_1, x_2, \dots]$ is not Noetherian — infinitely many variables defeat finiteness.

So non-Noetherian means a concrete infinite tower exists, not an abstract failure.

## Worked: (6, 10) = (2) in Z

Find one generator for the ideal $(6, 10)$:
1. Any common divisor story: $d$ generates iff $d = \gcd(6, 10)$.
2. $\gcd(6, 10) = 2$: indeed $6 = 2 \cdot 3$ and $10 = 2 \cdot 5$, and $2 = 10 - 6$ lies in the ideal.
3. So $(6, 10) = (2)$: principal, as Noetherian (here PID) structure demands.

So in a PID every multi-generator ideal collapses to a gcd computation.

## Hilbert basis

$R$ Noetherian implies $R[x]$ Noetherian: adjoining one variable preserves finiteness. By induction, $k[x_1, \dots, x_n]$ is Noetherian for finite $n$ — only infinitely many variables break it.
