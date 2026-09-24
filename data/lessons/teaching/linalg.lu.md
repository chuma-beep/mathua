# LU Decomposition

**LU decomposition:** $A = LU$ with $L$ unit lower-triangular and $U$ upper-triangular. It solves $Ax = b$ as two triangular systems: $Ly = b$ forward, then $Ux = y$ backward. Needs nonzero leading pivots (else permute: $PA = LU$).

## Worked: solve via Ly=b, Ux=y

Take $L = [[1,0],[2,1]]$, $b = [1,3]$:
1. Forward: $y_1 = 1$, then $2\cdot 1 + y_2 = 3$ gives $y_2 = 1$.
2. So $y = [1, 1]$: forward substitution reads top-down.
3. Then $Ux = y$ solves bottom-up the same way.

So triangular systems never need elimination — substitution walks straight through.

## Worked: det from U

With unit $L$ and $U$ diagonal $(2, 3)$:
1. $\det(L) = 1$ (unit diagonal).
2. $\det(U) = 2 \cdot 3 = 6$.
3. So $\det(A) = 6$: the determinant is the $U$-diagonal product, no expansion needed.

So LU factors the determinant computation into a diagonal read-off.

## Pivoting

Zero leading pivots need row swaps: $PA = LU$ with a permutation $P$. Partial pivoting (largest available pivot) is the stable default.
