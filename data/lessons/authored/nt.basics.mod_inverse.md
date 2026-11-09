# Modular inverses

A modular inverse undoes multiplication modulo n. It exists exactly when the multiplier is coprime to the modulus, and the extended Euclidean algorithm finds it.

## The inverse congruence

The inverse of a modulo n is the least residue x with a x congruent to 1 modulo n. The extended Euclidean algorithm writes 1 as a combination of a and n, and the coefficient of a is the inverse. Reduce that coefficient into the range 0 to n minus 1.

## Solving a linear congruence

To solve a x congruent to b, multiply both sides by the inverse of a. The answer is the least nonnegative residue of the product. Coprimality of a and n is what makes the inverse exist and the solution unique modulo n.

## Why coprimality matters

If a and n share a factor greater than 1, no inverse exists. The congruence may still have solutions, but they are not unique modulo n. Checking the gcd before inverting prevents a false answer.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
