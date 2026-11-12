# Step-size control

A stiff solver chooses its step from an error estimate and a stability constraint at the same time. The controller rejects a step whose estimated error exceeds the tolerance and grows the step only while the estimate stays small.

## The error estimate

An embedded pair computes two approximations of different order from the same stages. Their difference estimates the local error, and the estimate is compared with the tolerance after scaling by the solution size.

## The controller

The new step is the old step times a power of the ratio of tolerance to error. A safety factor keeps the controller from oscillating, and lower and upper bounds prevent a collapse or an enormous jump.

## Stability ceiling

An explicit method also has a maximum stable step set by the fastest mode. A stiff solver removes that ceiling by being stable on the whole left half-plane, so the accuracy estimate alone controls the step.

To check your understanding, compare the embedded pair, apply the power law and identify the fast mode.

Step control is what makes a stiff solver efficient, because the step follows the slow solution rather than the fast transient.
