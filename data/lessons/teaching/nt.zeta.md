# Riemann Zeta Function

**Riemann zeta:** \(\zeta(s) = \sum 1/n^s\) over \(n\) from 1 upward converges for \(\Re(s)\) above 1 and continues analytically to every \(s\) except 1, where it has a simple pole. Its Euler product ties it to the primes: \(\zeta(s) = \prod_p (1-p^{-s})^{-1}\).

## Worked: Basel at 2

Evaluate at \(s = 2\) two ways:
1. Series: 1 + 1/4 + 1/9 + 1/16 already reaches about 1.42, climbing toward the limit.
2. Closed form: \(\zeta(2) = \pi^2/6\), about 1.64.
3. The partial sum sits below the limit, as a positive-terms series must.

So \(\zeta(2)\) is the settled value of the reciprocal squares: Basel's answer in zeta language.

## Worked: Euler product at 2

Factor the same value over primes:
1. The \(p = 2\) factor is \((1-1/4)^{-1} = 4/3\).
2. Times the \(p = 3\) factor \((1-1/9)^{-1} = 9/8\): the running product is \(4/3 \cdot 9/8 = 3/2\).
3. More primes push 3/2 upward toward \(\pi^2/6\): the product over all primes equals the sum over all integers.

So the prime factorization of integers becomes a product formula for \(\zeta\): analysis meets arithmetic.

## Worked: pole and zeros

Locate the pole and the two zero families:
1. At \(s = 1\) the harmonic series diverges, so \(\zeta\) has a pole there with residue 1.
2. Trivial zeros sit at the negative even integers -2, -4, and so on.
3. The remaining zeros are conjectured onto the line \(\Re(s) = 1/2\): the Riemann hypothesis.

So one pole, predictable trivial zeros, and one famous open line: the anatomy of \(\zeta\) in three facts.
