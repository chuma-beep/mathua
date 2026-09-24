# Class Numbers

**Class number:** \(h(D)\) counts ideal classes in the quadratic order of discriminant \(D\); \(h(D) = 1\) exactly when the order is a PID, hence a UFD. Gauss asked which \(D\) give 1; Heegner's nine fundamental cases \(-1, -2, -3, -7, -11, -19, -43, -67, -163\) close the negative side.

## Worked: 6 factors two ways

In \(\mathbb{Z}[\sqrt{-5}]\), the number 6 splits twice:
1. \(6 = 2 \cdot 3\), and neither factor splits further: \(a^2+5b^2\) never equals 2 or 3, so no element has norm 2 or 3 to divide them.
2. Also \(6 = (1+\sqrt{-5})(1-\sqrt{-5})\): each bracket has norm \(1+5 = 6\), irreducible since norms 2 and 3 are missing.
3. The factorizations are genuinely distinct: 2 divides neither bracket (the quotients would need half-integer coefficients).

So unique factorization fails: one number, two irreducible factorizations, hence the class number exceeds 1 (it is 2).

## Worked: Gaussian integers have class number 1

\(\mathbb{Z}[i]\) is Euclidean under the norm, hence a PID:
1. For any \(\alpha\) and nonzero \(\beta\), round \(\alpha/\beta\) to the nearest lattice point: the error has norm at most 1/2, below 1.
2. So division leaves a remainder of strictly smaller norm: the Euclidean algorithm runs.
3. Euclidean implies PID implies UFD: every factorization into irreducibles is unique, so \(h(-4) = 1\).

So geometry decides arithmetic: the square lattice is fine enough to divide with remainder.

## Worked: the formula at minus 163

Dirichlet's class number formula links \(h(D)\) to \(L(1,\chi_D)\):
1. Plug \(D = -163\) (two roots of unity) into \(h = (w\sqrt{|D|}/2\pi) L(1,\chi_D)\).
2. The square root and the L-value combine to exactly 1: the largest Heegner discriminant.
3. So the analytic value certifies the algebraic fact: class number 1, unique factorization in that order.

So the formula converts analysis into arithmetic: an L-value proves a factorization claim.
