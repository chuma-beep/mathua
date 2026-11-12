# Semicontinuity

Semicontinuity relaxes continuity to one side. An upper semicontinuous function can jump down but not up.

## Upper and lower

A function is upper semicontinuous at \(x\) when \(\limsup_{y \to x} f(y) \le f(x)\). The indicator of \(\{0\}\) is upper semicontinuous, since it equals 1 at the origin and 0 elsewhere. Lower semicontinuity reverses the inequality. The indicator of an open set is lower semicontinuous, because its set of ones is open.

## Extreme values

An upper semicontinuous function on a compact set attains its maximum. The function \(-t^2\) on \([-1,1]\) is upper semicontinuous and peaks at 0 with value 0. It need not attain a minimum, which is why the two one-sided notions are not the same.

## Level sets

The superlevel set \(\{x : f(x) \ge a\}\) is closed when \(f\) is upper semicontinuous. For \(f(t) = -t^2\) and \(a = -1\) the set is \([-1,1]\), closed as promised. For a lower semicontinuous \(f\) the sublevel sets \(\{f \le a\}\) are closed instead, and the two families of closed sets determine the function.

Semicontinuity keeps the half of continuity that extremal arguments need, and the level sets record it.
