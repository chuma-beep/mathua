# Continuous Functions in Topology

In metric spaces, continuity can be defined using the $\epsilon$-$\delta$ definition or equivalently using open sets.

## $\epsilon$-$\delta$ Definition

A function $f: X \to Y$ between metric spaces is **continuous at $x_0$** if for every $\epsilon > 0$, there exists $\delta > 0$ such that:

$$d_X(x, x_0) < \delta \implies d_Y(f(x), f(x_0)) < \epsilon$$

## Open Set Definition

A function $f: X \to Y$ is **continuous** if the preimage of every open set is open:
$$U \subseteq Y \text{ open} \implies f^{-1}(U) \subseteq X \text{ open}$$

This definition works for any topological space (not just metric spaces).

## Equivalent Conditions

The following are equivalent for a function $f: X \to Y$:
1. $f$ is continuous
2. For every closed set $F \subseteq Y$, $f^{-1}(F)$ is closed in $X$
3. For every set $A \subseteq X$, $f(\overline{A}) \subseteq \overline{f(A)}$
4. For every convergent sequence $x_n \to x$, $f(x_n) \to f(x)$ (in metric spaces)

**Example 1:** $f(x) = x^2$ on $\mathbb{R}$.
Given $\epsilon > 0$, choose $\delta = \min(1, \frac{\epsilon}{2|x_0| + 1})$.
If $|x - x_0| < \delta$, then $|x^2 - x_0^2| = |x - x_0||x + x_0| < \delta(2|x_0| + 1) < \epsilon$.

**Example 2:** The function $f(x) = 1/x$ is continuous at every $x \neq 0$, but is not continuous at $x = 0$ (it is not even defined there).

## Properties

- The composition of continuous functions is continuous
- The sum, product, and quotient (where defined) of continuous functions are continuous
- A function is continuous iff it is continuous at every point
- If $X$ is compact, then every continuous $f: X \to \mathbb{R}$ attains a maximum and minimum (Extreme Value Theorem)
