# Nonhomogeneous recurrences

A nonhomogeneous linear recurrence adds a forcing term to the homogeneous part. Its general solution is one particular solution plus the general homogeneous solution.

## Particular plus homogeneous

If \(a^{(p)}_n\) solves the full recurrence, then every solution is \(a^{(p)}_n+h_n\) with \(h_n\) homogeneous. The initial values fix the homogeneous constants.

## Guessing a particular solution

A polynomial forcing term suggests a polynomial particular solution, and an exponential forcing term suggests the same exponential unless it is already a homogeneous root.

## Resonance

When the forcing matches a homogeneous root, the guess must be multiplied by \(n\). Otherwise the substitution produces a zero left side.

To check your understanding, find one particular solution, match the forcing shape and check the homogeneous roots.

Nonhomogeneous recurrences split into a transient homogeneous part and a persistent forced part. Resonance is the case where the two parts share a rate.
