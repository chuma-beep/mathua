# Bang-bang control

When the Hamiltonian is linear in the control, the pointwise minimum sits at a bound unless the coefficient vanishes. The optimal control then switches between extremes, which is bang-bang behaviour.

## The switching function

The coefficient of the control in the Hamiltonian is the switching function. Its sign selects the lower or the upper bound, and its zeros are the candidate switching times.

## Singular arcs

If the switching function vanishes over an interval, the control is not determined by the first-order condition. A singular arc requires higher derivatives of the switching function to be analysed, and it may be optimal or not.

## Why it matters

Bang-bang controls are the rule in time-optimal problems and in actuator-limited systems. They also produce chattering when implemented on a real actuator, so the theoretical solution is often smoothed before use.

To check your understanding, read the control coefficient, detect a vanishing interval and solve time-optimal problems.

Bang-bang control is what a linear Hamiltonian implies, and the singular arc is the case where the first-order condition goes silent.
