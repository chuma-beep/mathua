# Zeckendorf representation

Zeckendorf's theorem writes every positive integer uniquely as a sum of non-consecutive Fibonacci numbers. The greedy choice of the largest Fibonacci number produces the representation.

## The theorem

Use \(F_2=1,F_3=2,\ldots\) and forbid consecutive indices. Every positive integer has exactly one such sum.

## The greedy algorithm

Take the largest Fibonacci number not exceeding the remainder, subtract it, and repeat. The gap left is always smaller than the previous Fibonacci number, so indices cannot be consecutive.

## Uniqueness

Two representations would differ at the largest index used by one and not the other. The remaining sum is then too small to reach that Fibonacci number, which is a contradiction.

To check your understanding, list the Fibonacci numbers, take the largest term and compare largest terms.

Zeckendorf representation is a numeral system with Fibonacci place values and a no-consecutive-digits rule. Its greedy algorithm is optimal because Fibonacci numbers grow fast enough.
