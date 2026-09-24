# Continued Fractions

**Continued fraction:** \([a_0; a_1, a_2, \dots] = a_0+1/(a_1+1/(a_2+\cdots))\) with positive integer tails. Its convergents \(p_k/q_k\) are best rational approximations, and periodicity is equivalent to quadratic irrationality (Lagrange).

## Worked: 22 over 7 from 3 semicolon 7

Expand \([3; 7]\) bottom-up:
1. Start at the tail: the last entry is 7.
2. Invert and add: \(3+1/7 = 22/7\).
3. So \([3; 7] = 22/7\), matching \(\pi\) within 0.04 percent.

So truncation turns an infinite object into a fraction: convergents are just stopped expansions.

## Worked: root 2 repeats 2 forever

Show \(\sqrt{2} = [1; 2, 2, 2, \dots]\) via the regenerating remainder:
1. Split off the integer part: \(\sqrt{2} = 1+(\sqrt{2}-1)\) with remainder below 1.
2. Invert the remainder: \(1/(\sqrt{2}-1) = \sqrt{2}+1 = 2+(\sqrt{2}-1)\).
3. The same remainder recurs, so the digit 2 repeats forever.

So self-similarity forces periodicity: quadratic irrationals repeat because the remainder regenerates.

## Worked: best approximation quality

Convergents satisfy \(|x-p_k/q_k|\) below \(1/(q_k q_{k+1})\). Test \(x = \pi\) at \(22/7\) (next denominator 106):
1. Error of \(22/7\): about 0.00126.
2. Bound: \(1/(7 \cdot 106) = 1/742\), about 0.00135.
3. The error sits under the bound, as claimed.

So convergents carry certificates: each approximation ships with its own error bar.
