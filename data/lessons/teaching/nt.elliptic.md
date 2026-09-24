# Elliptic Curves and BSD

**Elliptic curve:** A nonsingular cubic \(E: y^2 = x^3+ax+b\) whose discriminant is nonzero, i.e. \(4a^3+27b^2\) differs from 0. Its points form an abelian group under chord-tangent addition, and over the rationals that group is finitely generated (Mordell-Weil).

## Worked: doubling on y squared equals x cubed minus 2

Take \(E: y^2 = x^3-2\) and \(P = (3,5)\); double it via the tangent slope \(m = 3x^2/2y\):
1. Slope at \(P\): \(m = 27/10\).
2. New \(x\): \(m^2-2x = 729/100-6 = 129/100\).
3. Reflect for the new \(y\): \(m(x-x_{new})-y = (27/10)(171/100)-5 = -383/1000\).

So \(2P = (129/100, -383/1000)\) lies on \(E\): doubling manufactures genuine new rationals from old.

## Worked: Hasse count mod 5

Take \(E: y^2 = x^3+1\) over the field of 5 elements and count points:
1. \(x = 0\) gives \(y^2 = 1\), so \(y = 1, 4\): two points.
2. \(x = 1\) gives \(y^2 = 2\): no solution (squares mod 5 are 0, 1, 4); \(x = 2\) gives \(y^2 = 9 = 4\), so \(y = 2, 3\): two points.
3. \(x = 3\) gives \(y^2 = 28 = 3\): none; \(x = 4\) gives \(y^2 = 65 = 0\), so \(y = 0\): one point; plus the point at infinity: 6 total.

So the count is 6 against \(p+1 = 6\): error 0, inside Hasse's \(2\sqrt{5}\) bound.

## Worked: rank zero and BSD

Take \(E: y^2 = x^3-x\) over the rationals:
1. Rational points: \((-1,0)\), \((0,0)\), \((1,0)\), plus infinity, each doubling to infinity.
2. So every rational point is torsion: the rank is 0 with no point of infinite order.
3. BSD predicts \(L(E,1)\) differs from 0 for rank 0, and indeed it is about 0.66.

So torsion-only forces rank 0 while the L-value stays nonzero: BSD holding in the simplest case.
