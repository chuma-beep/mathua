# Quadratic Reciprocity

**Law of quadratic reciprocity:** For distinct odd primes \(p\) and \(q\), the symbols \((p/q)\) and \((q/p)\) agree unless both primes are 3 mod 4, in which case they disagree. The product is \((-1)^{((p-1)/2)((q-1)/2)}\).

## Worked: 5 and 7 agree

Since 5 is 1 mod 4, the sign is positive and \((5/7) = (7/5)\):
1. Reduce: 7 is 2 mod 5, so \((7/5) = (2/5)\).
2. 2 is a non-residue mod 5 (the squares are 1 and 4 only), so \((2/5) = -1\).
3. Hence \((5/7) = -1\) with no search mod 7.

So reciprocity trades a mod 7 question for a mod 5 question, then finishes by inspection.

## Worked: 3 and 11 disagree

Both 3 and 11 are 3 mod 4, so the sign is negative: \((3/11) = -(11/3)\):
1. Reduce: 11 is 2 mod 3, so \((11/3) = (2/3)\).
2. Squares mod 3 are just 1, so \((2/3) = -1\).
3. Hence \((3/11) = 1\): indeed \(5^2 = 25\), and \(25 - 3 = 22\) is a multiple of 11.

So the minus sign flips the answer: \((3/11) = 1\) while \((11/3) = -1\).

## Worked: computing 3 over 5 fast

Evaluate \((3/5)\) two ways and compare:
1. Direct: residues mod 5 are 1 and 4, and 3 is missing, so \((3/5) = -1\).
2. Reciprocity: 5 is 1 mod 4, so \((3/5) = (5/3) = (2/3) = -1\).
3. Both routes agree at -1.

So reciprocity is a shortcut, not a new rule: small cases verify by hand.
