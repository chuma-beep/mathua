# Metric spaces

To talk about *closeness* in full generality, mathematicians strip the idea of distance down to its essential properties.

## Definition and Axioms

A **metric space** is a set \\(X\\) together with a function \\(d: X \times X \to \mathbb{R}\\), called the **metric**, that assigns a distance \\(d(x,y)\\) to every pair of points and satisfies four rules for all \\(x, y, z \in X\\):

1. **Non-negativity:** \\(d(x,y) \ge 0\\).
2. **Identity:** \\(d(x,y) = 0\\) if and only if \\(x = y\\).
3. **Symmetry:** \\(d(x,y) = d(y,x)\\).
4. **Triangle inequality:** \\(d(x,z) \le d(x,y) + d(y,z)\\).

## Examples

The real line with ordinary distance \\(d(x,y)=|x-y|\\) is a metric space, and so is the plane with Euclidean distance. But the definition is deliberately flexible: the same set can carry many different metrics (discrete metric, taxicab metric), and the points need not be numbers at all — strings, functions, or images can live in a metric space whenever distances behave sensibly.

## Why Metrics Matter

Every statement you will meet about limits, open sets, or continuity in this domain is phrased purely in terms of the metric, which is why these three axioms power so much of analysis and topology.
