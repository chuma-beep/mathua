# Characteristic roots

The characteristic roots of a linear recurrence determine its growth. Distinct roots give pure exponentials, while a repeated root of multiplicity \(m\) contributes a polynomial factor of degree \(m-1\).

## Distinct roots

Each root \(r_i\) contributes a term \(A_ir_i^n\). The constants \(A_i\) are fixed by the initial values.

## Repeated roots

A root of multiplicity \(m\) contributes \(n^j r^n\) for \(j=0,\ldots,m-1\). The polynomial factor is what distinguishes the solutions.

## Dominant root

The root of largest modulus controls the asymptotic growth. Ties in modulus produce oscillating or polynomial behaviour.

To check your understanding, form the roots, count multiplicity and compare moduli.

Characteristic roots convert a recurrence into a closed form. The multiplicity bookkeeping is the only subtle part.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
