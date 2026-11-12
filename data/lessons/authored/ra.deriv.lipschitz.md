# Lipschitz continuity

A Lipschitz function changes no faster than a fixed multiple of the distance. The multiple is a quantitative measure of regularity.

## The Lipschitz constant

The least \(L\) with \(|f(x)-f(y)| \le L|x-y|\) is the Lipschitz constant. For \(f(t) = 3t + 1\) it is 3, and for \(\sin t\) it is 1. A differentiable function with bounded derivative is Lipschitz with \(L\) equal to the supremum of \(|f'|\).

## Contraction and absolute continuity

A Lipschitz function is absolutely continuous, so it is differentiable almost everywhere and equals the integral of its derivative. When \(L < 1\) the map is a contraction, and iteration converges to a fixed point: \(f(t) = (t+2)/3\) has \(L = 1/3\), and iterating from 0 gives 0.666, 0.888, 0.962, approaching 1.

## Functions that fail

The square root on \([0,1]\) is not Lipschitz, because \((\sqrt{x} - 0)/(x - 0) = 1/\sqrt{x}\) is unbounded near 0. The function \(t^2\) on the whole line fails as well, since the ratio of the change to the distance grows with the points. On a bounded interval \(t^2\) is Lipschitz with constant \(2b\).

Lipschitz continuity bounds the difference quotient without demanding a derivative, and the bound is what makes iteration and integration work.
