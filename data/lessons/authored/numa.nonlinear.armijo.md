# Armijo sufficient decrease

The Armijo condition requires a sufficient decrease in the objective along the search direction. It rules out steps that are too long to help and too short to make progress.

## Statement

\(f(x+\alpha p)\le f(x)+c\alpha\nabla f(x)^{T}p\) with a small constant \(c\), typically \(10^{-4}\). The right side is a line with slight negative slope.

## Backtracking

Starting from \(\alpha=1\), the step is multiplied by a factor \(\rho\in(0,1)\) until the condition holds. Each trial costs one function evaluation.

## Why it suffices

Together with a curvature condition on the same step, Armijo guarantees convergence of the line-search method. Alone it still prevents unbounded increase.

To check your understanding, write the inequality, start at one and pair with curvature.

Armijo asks for a fraction of the decrease the gradient predicts.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
