# The Maximum Modulus Principle

## The maximum principle

The maximum modulus principle says that a nonconstant analytic function on a domain cannot attain a maximum modulus at an interior point. If the modulus has a local maximum inside the domain, the function must be constant on that domain. The maximum can occur on the boundary, and for a bounded domain it often does.

The principle follows from the open mapping property or from Cauchy's integral formula. If the modulus were maximal at an interior point, averaging over a small circle would force every value on the circle to have the same modulus, and then the function would be constant nearby. Connectedness extends the conclusion to the whole domain.

For a continuous function on a closed bounded domain that is analytic inside, the maximum modulus is attained somewhere on the boundary. This turns a two-dimensional optimization problem into a one-dimensional one. The minimum modulus does not behave the same way, because zeros can occur inside.

## The minimum principle

If an analytic function is nonzero on a domain, then its reciprocal is analytic, so the maximum principle applied to the reciprocal gives a minimum principle for the original function. A nonzero analytic function cannot attain a positive minimum modulus at an interior point unless it is constant.

Zeros are the exception. The modulus is zero at a zero, which is certainly a minimum, but the function is not constant. The correct statement therefore excludes zeros or treats them separately. On a closed bounded domain, the minimum modulus is either zero at an interior zero or attained on the boundary.

These two principles are useful for bounding solutions of equations. If a function is known on the boundary of a disc, the maximum principle bounds it inside. If it has no zeros, the minimum principle gives a lower bound as well. Together they sandwich the modulus between boundary data.

## Applications

The maximum modulus principle proves that a nonconstant polynomial must have a root. If a polynomial had no root, its reciprocal would be entire and bounded, since the polynomial tends to infinity at infinity. Liouville's theorem would then force the reciprocal to be constant, which is impossible for a nonconstant polynomial.

The principle also gives uniqueness for boundary value problems. If two analytic functions have the same boundary values on a domain, their difference is analytic and has zero boundary values, so its maximum modulus is zero and the functions agree. This is a prototype for elliptic uniqueness theorems in analysis.

Numerically, the principle warns against searching for a maximum modulus in the interior. Evaluate on the boundary, or use the structure of the function to identify candidate boundary points. For a disc, parametrize the boundary circle and reduce the problem to a trigonometric optimization, which is often simpler than a two-variable search.

