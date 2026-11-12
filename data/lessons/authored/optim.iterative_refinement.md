# Iterative refinement

A direct solve in floating point leaves a residual because of rounding. Iterative refinement computes that residual in higher precision and solves a correction system with the existing factorisation.

## The correction step

The residual of the computed solution is evaluated, then the same factorisation solves for a correction. Adding the correction reduces the error when the factorisation is not too ill-conditioned.

## Mixed precision

The residual is often computed in a wider format than the factorisation. That combination recovers much of the accuracy of a fully wide solve at a fraction of the cost.

## When it fails

If the condition number exceeds the reciprocal of the working precision, the correction does not converge. The remedy is a better preconditioner or a wider factorisation, not more refinement rounds.

To check your understanding, compute the residual, widen the residual and compare conditioning to precision.

Iterative refinement is a cheap way to recover accuracy from a factorisation already computed, and its limit is the condition number.
