# Order Modulo n

**Order:** For \(a\) coprime to \(m\), the order \(\operatorname{ord}_m(a)\) is the smallest \(k\) above 0 with \(a^k = 1\) mod \(m\). It always divides \(\varphi(m)\), and \(p - 1\) when \(m = p\) is prime.

## Worked: order of 2 mod 7

Raise 2 until 1 appears mod 7:
1. \(2^1 = 2\) and \(2^2 = 4\).
2. \(2^3 = 8 = 1\) mod 7, and no smaller positive exponent works.
3. So the order is 3, which divides \(\varphi(7) = 6\).

So the powers cycle 2, 4, 1 with period 3: the order is the cycle length.

## Worked: order of 3 mod 7

Same modulus, new base:
1. \(3^1 = 3\), \(3^2 = 9 = 2\) mod 7, \(3^3 = 6\) mod 7.
2. Continue: \(3^4 = 4\), \(3^5 = 5\), \(3^6 = 1\) mod 7, the first return to 1.
3. So the order is 6, the full \(\varphi(7)\): base 3 generates everything.

So 3 is a primitive root mod 7 while 2 is not: order measures how much of the group one element sees.

## Worked: testing divisors of phi

For \(a = 2\), \(m = 9\): \(\varphi(9) = 6\), so the order divides 6; test its divisors 1, 2, 3, 6:
1. \(2^1 = 2\), \(2^2 = 4\), \(2^3 = 8\): none is 1 mod 9.
2. \(2^6 = 64 = 1\) mod 9, since 63 is a multiple of 9.
3. So the order is 6 with no smaller candidate surviving.

So only divisors of \(\varphi(m)\) need testing: the order is the first divisor that returns 1.
