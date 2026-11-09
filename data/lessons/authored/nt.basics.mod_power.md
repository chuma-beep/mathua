# Modular exponentiation

Modular exponentiation computes a power without ever forming the full integer. The method is repeated squaring, and it is the arithmetic behind every public-key cryptosystem in this domain.

## Repeated squaring

Write the exponent in binary and square the base at each step. A one bit multiplies the running result by the current square, while a zero bit only advances to the next square. Every intermediate value is reduced modulo the modulus, so the numbers stay small.

## Reducing the exponent

When the base is coprime to the modulus, Euler's theorem lets the exponent be reduced modulo phi of the modulus. The reduced exponent gives the same residue and keeps the repeated-squaring loop short. The modulus and the base are unchanged by the reduction.

## Multiplying residues

A product of two powers is the power of the sum of the exponents. Reduce each power first, multiply the residues and reduce once more. This is the move that turns a long product into a short sequence of modular multiplications.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
