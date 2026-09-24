# Prime Number Theorem

**Prime number theorem (PNT):** With \(\pi(x)\) counting primes up to \(x\), the ratio of \(\pi(x)\) to \(x/\log x\) tends to 1 as \(x\) grows. Equivalently the \(n\)-th prime sits near \(n \log n\), and a random integer near \(x\) is prime with probability about \(1/\log x\).

## Worked: estimating pi of 100

Compare the one-line estimate against the true count:
1. \(\log 100\) is about 4.6, so \(100/\log 100\) is about 21.7.
2. The actual count is 25: there are 25 primes up to 100.
3. The estimate lands within 4 of the truth already, and the relative error shrinks as \(x\) grows.

So \(x/\log x\) tracks \(\pi(x)\): a one-line estimate with honest accuracy.

## Worked: density near a million

Apply the density form \(1/\log x\) at \(x = 10^6\):
1. \(\log 10^6\) is about 13.8.
2. So about 1 in 14 integers near a million is prime.
3. Equivalently a block of 140 consecutive integers there holds about 10 primes.

So the theorem predicts local frequency, not just global totals: primality thins out logarithmically.

## Worked: the nth prime heuristic

Invert the estimate to guess the \(n\)-th prime via \(p_n \approx n \log n\):
1. For \(n = 1000\): \(\log 1000\) is about 6.9, so the guess is about 6900.
2. The true 1000th prime is 7919.
3. The guess runs low by about 13 percent, typical first-order accuracy.

So the same law reads backwards: counts give sizes, and sizes give counts.
