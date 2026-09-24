# Legendre Symbol

**Legendre symbol:** For an odd prime \(p\) and integer \(a\), the symbol \((a/p)\) is 0 when \(p\) divides \(a\), 1 when \(a\) is a nonzero square mod \(p\), and -1 otherwise. It compresses the question "is \(a\) a square mod \(p\)" into one value.

## Worked: residues mod 5

Square every nonzero residue mod 5 and sort the outputs:
1. \(1^2 = 1\) and \(4^2 = 16 = 1\) mod 5, so 1 is a residue.
2. \(2^2 = 4\) and \(3^2 = 9 = 4\) mod 5, so 4 is a residue.
3. Nothing squares to 2 or 3, so both are non-residues.

So \((1/5) = 1\), \((4/5) = 1\), and \((2/5) = -1\): exactly half the nonzero residues are squares.

## Worked: Euler criterion mod 7

Euler's criterion says \((a/p) = 1\) exactly when \(a^{(p-1)/2} = 1\) mod \(p\). Test \(a = 2\), \(p = 7\):
1. The exponent is \((7-1)/2 = 3\).
2. \(2^3 = 8 = 1\) mod 7, so the criterion fires.
3. Confirm by hand: \(3^2 = 9 = 2\) mod 7, so 2 is indeed a square.

So one exponentiation replaces a full search: Euler detects the witness 3 without listing every square.

## Worked: the supplement for 2

The supplement says \((2/p) = 1\) exactly when \(p\) is 1 or 7 mod 8. Check \(p = 17\):
1. 17 mod 8 is 1, so the supplement predicts \((2/17) = 1\).
2. Confirm with \(x = 6\): \(6^2 = 36\), and \(36 - 2 = 34\) is a multiple of 17.
3. Contrast \(p = 5\): 5 mod 8 is 5, predicting -1, matching \((2/5) = -1\) above.

So the residue class of \(p\) mod 8 decides \((2/p)\) with no search at all.
