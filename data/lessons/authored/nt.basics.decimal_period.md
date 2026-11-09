# Repeating decimal periods

The repeating part of a rational decimal is governed by the order of 10 modulo the denominator after the factors 2 and 5 are removed.

## The period as an order

For a denominator n coprime to 10, the decimal of 1 over n repeats with the smallest k such that 10 to the k is congruent to 1 modulo n. That k is the multiplicative order of 10 modulo n. The repetend then has exactly k digits.

## Reducing the fraction

A numerator that shares a factor with the denominator changes the digits but not the period. Reduce first and then compute the order of 10 modulo the reduced denominator. The numerator does not enter the period calculation.

## Removing terminating factors

Factors of 2 and 5 produce a terminating prefix and do not affect the repeating tail. Remove them, compute the order modulo the remaining core and attach the prefix afterward. If the core is 1, the decimal terminates and the period is zero.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.
