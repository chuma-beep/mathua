# Coercivity and existence

A function is coercive when it tends to infinity as its argument leaves every bounded region. Coercivity plus continuity guarantees that a minimum exists, because a minimizing sequence cannot escape to infinity. Without it, an infimum can be approached but never attained, and a solver can run forever toward a value it will not reach.

## Testing coercivity

For a quadratic with matrix, coercivity is positive definiteness, since the quadratic term dominates the linear term in every direction. For a polynomial in one variable, the leading degree must be even and its coefficient positive. A function that is constant along a direction, such as a linear function, is not coercive, because it does not grow in that direction.

## Deciding existence

Continuity and coercivity together imply that a global minimizer exists. Continuity alone gives an infimum. On a closed bounded feasible set, continuity alone is enough, because the set is compact. Existence is therefore a statement about the interaction between the function and the set, and each hypothesis covers a different way a minimizing sequence could fail to converge.

## Bounding a minimizing sequence

Coercivity gives an explicit radius beyond which the function exceeds any chosen level. If a point with value below that level is known, the minimizer must lie inside the corresponding ball. That bound turns an infinite search into a bounded one, and it is the first step in many global methods. The bound is computable from the growth rate, not merely existential.

Existence is the hypothesis behind writing the minimum rather than the infimum. Coercivity is the usual way to obtain it on an unbounded domain, and it also supplies the radius that makes a search finite.
