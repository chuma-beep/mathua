# Primitive Roots

**Primitive root:** For prime \(p\), a base \(g\) is primitive when its order is \(p - 1\): its powers \(g^0\) through \(g^{p-2}\) cover every nonzero residue mod \(p\). Every prime has \(\varphi(p-1)\) of them.

## Worked: 2 is primitive mod 5

List the powers of 2 mod 5:
1. \(2^0 = 1\), \(2^1 = 2\), \(2^2 = 4\), \(2^3 = 3\) mod 5.
2. The set 1, 2, 4, 3 is all four nonzero residues.
3. So the order is 4, which equals \(5 - 1\): base 2 is primitive mod 5.

So one orbit covers the whole group: that is what primitive means.

## Worked: 2 is not primitive mod 7

Same base, new modulus:
1. Powers: \(2^1 = 2\), \(2^2 = 4\), \(2^3 = 1\) mod 7, then the cycle repeats.
2. The orbit 2, 4, 1 has 3 elements, not 6.
3. Since 3 differs from \(7 - 1 = 6\), base 2 is not primitive mod 7.

So primitivity depends on the modulus, not the base alone: 2 works mod 5 and fails mod 7.

## Worked: counting them mod 7

The primitive roots mod 7 are exactly the elements of order 6:
1. Test \(g = 3\): powers run 3, 2, 6, 4, 5, 1, all six distinct, so the order is 6.
2. Test \(g = 5\): powers run 5, 4, 6, 2, 3, 1, all six distinct, so the order is 6.
3. Count check: \(\varphi(7-1) = \varphi(6) = 2\), matching the pair 3, 5.

So the count \(\varphi(p-1)\) predicts and listing confirms: exactly 2 primitive roots mod 7.
