# Inverse power iteration

Inverse power iteration applies the power method to the inverse of a shifted matrix. It converges to the eigenvalue closest to the shift.

## The shift

The method computes with \((A-\sigma I)^{-1}\). Its dominant eigenvalue corresponds to the eigenvalue of \(A\) closest to \(\sigma\).

With shift zero, it finds the eigenvalue of smallest magnitude. A good shift can target a chosen eigenvalue.

If the shift equals an eigenvalue exactly, the shifted matrix is singular and the iteration needs a modified step.

## Convergence

The convergence rate depends on how much closer the target eigenvalue is to the shift than the next nearest one.

If two eigenvalues are equally close to the shift, convergence can be slow. The shift must avoid that tie.

A shift very close to the target makes the first factor large and the convergence fast.

## The power method connection

The power method on \(A\) converges to the eigenvalue of largest magnitude. Inverse power reverses the ordering through the inverse.

The two methods are the same iteration applied to different matrices. Both need a starting vector with a component in the target direction.

In practice each inverse-power step solves a linear system rather than forming an inverse.
