# The uniform boundedness principle

The uniform boundedness principle turns pointwise bounds on a family of operators into a uniform bound on their norms. It needs completeness and it rests on Baire category.

## Pointwise bounds

Let \(T_\alpha\) be bounded operators from a Banach space \(X\) to a normed space \(Y\). If \(\sup_\alpha \|T_\alpha x\| < \infty\) for each \(x\), then \(\sup_\alpha \|T_\alpha\| < \infty\). The hypothesis is pointwise, and the conclusion is uniform over the whole space.

## The Baire argument

The sets \(F_n = \{x : \sup_\alpha \|T_\alpha x\| \le n\}\) are closed and cover \(X\). Baire category gives one \(F_n\) with nonempty interior, so a ball of radius \(r\) is bounded by \(n\). Scaling moves the ball to the origin and gives the norm bound \(2n/r\).

## Families of operators

The Fourier partial sums \(S_N\) on \(C[0,1]\) have norms growing like \(\log N\). If \(\sup_N |S_N f(0)|\) were finite for every continuous \(f\), the principle would bound the norms. Since the norms are unbounded, some continuous function has a divergent Fourier series at the origin. The principle supplies that existence without naming the function.

A pointwise bound is weaker than it looks, and completeness upgrades it to a bound on the norms.
