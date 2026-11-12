# The SAG estimator

SAG keeps the same table as SAGA but uses the table average as the direction, with one entry refreshed per step. It is the older and simpler member of the family.

## The averaging rule

Each step samples an index, replaces the stored gradient at that index, and uses the mean of the table as the direction. The mean changes by one term per step, so it is maintained with a running sum.

## Bias at the start

The table begins empty or at zero, so early averages are biased towards the initial values. A common remedy is a short full-gradient phase, after which the table holds real gradients for every sample.

## Convergence behaviour

The method converges for a decreasing step size under smoothness and strong convexity. With a constant step it reaches a neighbourhood whose radius depends on the step and the gradient spread.

To check your understanding, refresh one entry, expect a biased table and decrease the step.

SAG shows the central idea of the family: remember past sample gradients, reuse them, and let the average carry the signal while the noise averages out.
