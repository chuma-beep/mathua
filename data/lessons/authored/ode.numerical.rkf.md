# Runge-Kutta-Fehlberg methods

Runge-Kutta-Fehlberg methods compute two approximations of different order from the same stages. Their difference estimates the local error.

## Compute the embedded pair

The method evaluates several stages and combines them into a fourth-order and a fifth-order estimate. The two estimates share the stage values. The higher-order estimate is usually the one carried forward.

## Estimate the local error

The difference between the two estimates is the local error estimate. It is compared with the requested tolerance. A large difference calls for a smaller step.

## Choose the next step

Scale the step by the tolerance divided by the error estimate raised to one over the order difference. The usual exponent is one fifth. Clamp the factor so the step cannot jump too far in one update.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
