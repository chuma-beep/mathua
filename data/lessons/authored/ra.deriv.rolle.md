# Rolle's theorem

Rolle's theorem guarantees a horizontal tangent between two points of equal height. It is the engine behind the mean value theorem.

## Equal endpoints

Let \(f\) be continuous on \([a,b]\) and differentiable inside, with \(f(a) = f(b)\). Then some \(c\) in \((a,b)\) has \(f'(c) = 0\). The function \(f(t) = t^2 - 4t + 3\) has \(f(0) = f(4) = 3\), and \(f'(t) = 2t - 4\) vanishes at \(c = 2\).

## The interior critical point

The proof takes a maximum or a minimum in the interior. If \(f\) is constant the derivative vanishes everywhere. Otherwise some interior point is extremal, and the derivative there is zero because the difference quotients from the two sides have opposite signs.

## Proving the mean value theorem

Subtract the chord: \(g(t) = f(t) - f(a) - \frac{f(b)-f(a)}{b-a}(t-a)\). Then \(g(a) = g(b) = 0\), so \(g'(c) = 0\) for some \(c\), which gives \(f'(c) = \frac{f(b)-f(a)}{b-a}\). For \(f(t) = t^3\) on \([0,2]\) the slope is 4, and the point is \(c = \sqrt{4/3} \approx 1.155\).

Rolle's theorem is the mean value theorem with the chord flattened, and every application of the mean value theorem runs through it.
