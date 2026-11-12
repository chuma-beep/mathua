# Convex functions

A convex function lies below every chord joining two of its points. The condition is local in form yet global in consequence.

## The chord condition

For \(x < y\) and \(0 < \lambda < 1\), convexity says \(f(\lambda x + (1-\lambda)y) \le \lambda f(x) + (1-\lambda) f(y)\). For \(f(t) = t^2\) with \(x = 0\) and \(y = 2\) at the midpoint, the left side is 1 and the right side is 2. The parabola sits under its chord.

## Derivatives and monotonicity

A differentiable convex function has an increasing derivative. For \(t^2\) the derivative \(2t\) increases everywhere. The tangent at any point lies below the graph, so \(f(y) \ge f(x) + f'(x)(y-x)\). At \(x = 1\) this reads \(y^2 \ge 1 + 2(y-1)\), which is \((y-1)^2 \ge 0\).

## Jensen's inequality

Averaging the chord condition gives \(f\) of a weighted mean at most the weighted mean of the values. With weights \(1/4\) and \(3/4\) on 1 and 3, the mean is 2.5, and \(6.25 \le 7\). The inequality is the working form of convexity in probability and in optimization.

Convexity is a comparison between a function and its chords, and that comparison yields both monotone derivatives and Jensen's inequality.
