# Aitken acceleration

Aitken acceleration estimates the limit of a slowly converging sequence from three consecutive terms. It removes the dominant geometric error.

## The delta-squared correction

For three iterates, the correction uses \(\Delta x_n=x_{n+1}-x_n\) and \(\Delta^2x_n=x_{n+2}-2x_{n+1}+x_n\).

The accelerated estimate is \(x_n-(\Delta x_n)^2/\Delta^2x_n\). It extrapolates the geometric tail to its limit.

The denominator can be small, so compute the differences before dividing.

## Why it works for geometric error

If the error shrinks by a constant factor each step, the iterates have the form \(L+Cr^n\). Aitken acceleration recovers \(L\) exactly.

A ratio close to 1 means slow convergence. The correction can then remove most of the remaining error in one step.

A ratio near 0 means the sequence is already fast, so acceleration helps less.

## When it helps

Aitken acceleration is useful for linearly convergent fixed-point iteration. It can be applied repeatedly to the accelerated sequence.

It can amplify noise when differences are tiny or when the error is not geometric. Check the sequence before trusting the estimate.

If the acceleration produces a value far outside the range of the iterates, suspect a near-zero denominator.
