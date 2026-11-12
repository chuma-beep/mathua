# Oscillation

The oscillation of a function measures how much it moves near a point. It turns the qualitative idea of a jump into a number.

## Oscillation at a point

On an interval \(I\) around \(x\), the oscillation is the supremum of \(f\) on \(I\) minus its infimum there. The oscillation at \(x\) is the limit of these numbers as the interval shrinks to \(x\). For \(f(t) = \sin(1/t)\) near 0 every interval meets both 1 and -1, so the oscillation is 2.

## Discontinuities as oscillation

A function is continuous at \(x\) exactly when its oscillation there is zero. For \(f(t) = t\sin(1/t)\) the values shrink with \(t\), so the oscillation at 0 is 0 once \(f(0) = 0\) is set. A jump of size 1 gives oscillation 1, and the oscillation is therefore the size of the discontinuity.

## Oscillation and integrability

A bounded function on \([a,b]\) is Riemann integrable exactly when its set of discontinuities has measure zero. The oscillation bounds the gap between the upper and lower sums on each small interval, so the points where it stays positive are exactly the obstruction. The indicator of the rationals has oscillation 1 everywhere and is not integrable.

Oscillation assigns a magnitude to each failure of continuity, and the total of those magnitudes decides integrability.
