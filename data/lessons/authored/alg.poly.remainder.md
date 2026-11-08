# The remainder theorem

Dividing a polynomial by \(x-a\) looks like long division, but the remainder can be found with one substitution. The remainder theorem says the remainder is \(f(a)\), the value of the polynomial at the number that makes the divisor zero.

## Evaluate instead of divide

The divisor \(x-a\) is zero when \(x=a\). Substituting that value into the polynomial gives the remainder directly. For \(f(x)=x^3-2x+5\) divided by \(x-2\), the remainder is \(f(2)=8-4+5=9\).

This shortcut matters because long division is slow and easy to misalign. The theorem replaces a column of arithmetic with one evaluation. It works for any degree and any divisor of the form \(x-a\).

## Why the theorem works

Write the division as \(f(x)=(x-a)q(x)+r\), where \(q\) is the quotient and \(r\) is a constant remainder. The remainder has degree lower than the divisor, so for a linear divisor it is a number.

Set \(x=a\). The factor \(x-a\) becomes zero, the quotient term disappears, and \(f(a)=r\). The identity holds for every \(x\), so it must hold at this particular one. That is the whole proof.

## Remainders and roots

A remainder of zero means \(x-a\) divides the polynomial exactly, so \(a\) is a root. The remainder theorem therefore turns root testing into evaluation: if \(f(a)=0\), the factor belongs.

A nonzero remainder is also useful. It says how far the polynomial is from being divisible at that point. In practice, the theorem is the fastest way to check a suspected factor before committing to a full factorisation.
