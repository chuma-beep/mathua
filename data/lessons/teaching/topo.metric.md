# Metric Spaces

A **metric space** is a set $X$ together with a **distance function** (metric) $d: X \times X \to \mathbb{R}$ satisfying:

1. **Non-negativity:** $d(x, y) \ge 0$, with $d(x, y) = 0$ iff $x = y$
2. **Symmetry:** $d(x, y) = d(y, x)$
3. **Triangle inequality:** $d(x, z) \le d(x, y) + d(y, z)$

## Examples

**Example 1: Euclidean space**
$X = \mathbb{R}^n$, $d(x, y) = \sqrt{\sum_{i=1}^n (x_i - y_i)^2}$
This is the usual distance in $n$-dimensional space.

**Example 2: Manhattan metric**
$X = \mathbb{R}^n$, $d(x, y) = \sum_{i=1}^n |x_i - y_i|$
Distance measured along grid lines (like city blocks).

**Example 3: Discrete metric**
Any set $X$, $d(x, y) = 0$ if $x = y$, $d(x, y) = 1$ if $x \neq y$.
Every point is at distance 1 from every other point.

**Example 4: Supremum metric**
$X = C[0,1]$ (continuous functions on $[0,1]$), $d(f, g) = \sup_{x \in [0,1]} |f(x) - g(x)|$.
The distance is the maximum vertical gap between two functions.

## Open Balls

The **open ball** centered at $x$ with radius $r > 0$:
$$B_r(x) = \{y \in X : d(x, y) < r\}$$

Open balls generalize the idea of "open intervals" to any metric space and are the building blocks for defining open sets.

**Example:** In $\mathbb{R}^2$ with the Euclidean metric, $B_r(x)$ is a disk (without its boundary). With the Manhattan metric, $B_r(x)$ is a diamond shape.
