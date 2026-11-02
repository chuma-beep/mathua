# The Arzela-Ascoli Theorem

## Equicontinuity

A family of functions is equicontinuous when one modulus of continuity works for every member of the family. For every positive epsilon there is a positive delta such that inputs within delta produce outputs within epsilon, uniformly over the family. This is stronger than continuity of each individual function.

On a compact interval, differentiability with a common bound on the derivative implies equicontinuity by the mean value theorem. A common Lipschitz constant is another sufficient condition. Uniform equicontinuity is the key hypothesis that prevents the family from developing ever steeper oscillations.

## Uniform boundedness

A family is uniformly bounded when there is one constant that bounds the absolute value of every member at every point. Boundedness at a single point together with equicontinuity already gives uniform boundedness on a compact interval, because equicontinuity spreads the bound from a finite net of points to the whole interval.

Uniform boundedness alone is not enough for compactness in the space of continuous functions. The family of functions x to the n-th power on the unit interval is uniformly bounded but not equicontinuous, and its limit is discontinuous. Equicontinuity is what rules out that kind of degeneration.

## The compactness conclusion

Arzela-Ascoli says that a uniformly bounded and equicontinuous sequence of continuous functions on a compact interval has a uniformly convergent subsequence. The limit is continuous. In the language of metric spaces, the family is relatively compact in the uniform metric.

The proof uses a countable dense set of points and a diagonal argument to extract convergence there, then equicontinuity upgrades pointwise convergence to uniform convergence. The theorem is the function-space analogue of Bolzano-Weierstrass and is a standard tool for proving existence of solutions to differential equations.
