# Jensen's inequality

Jensen's inequality says a convex function of an average is at most the average of the values. It is the most used consequence of convexity.

## The finite form

For weights \(\lambda_i \ge 0\) summing to 1 and a convex \(f\), \(\sum \lambda_i f(x_i) \ge f\left(\sum \lambda_i x_i\right)\). Take \(f(t) = t^2\) with weights \(1/4\) on 1 and \(3/4\) on 3. The left side is 7 and the right side is 6.25.

## Weights and expectations

The weights are a probability distribution on finitely many points, so the inequality reads \(E[f(X)] \ge f(E[X])\). For \(f(t) = e^t\) with equal weights on 0 and \(\ln 4\), the mean of the values is 2.5 and \(e\) of the mean is 2. The exponential of an average is smaller than the average of exponentials.

## Equality cases

Equality holds when all the \(x_i\) with positive weight are equal, or when \(f\) is affine on the interval they span. For the strict case \(f(t) = t^2\) and points 1 and 3, equality would need \(7 = 6.25\), which fails. For the affine \(f(t) = 2t + 1\), the two sides agree at 6 on the same data.

Jensen's inequality converts an average inside a convex function into a bound outside it, and the equality case says when nothing is lost.
