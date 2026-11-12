# The stiffness ratio

A stiff system contains rates that differ by orders of magnitude. The stiffness ratio measures that spread, and it decides how much a step is limited by stability rather than by accuracy.

## The definition

The ratio is the largest magnitude eigenvalue divided by the smallest. Equivalently it is the longest time scale divided by the shortest, so a large ratio means the solution has both slow and fast behaviour.

## Why explicit methods suffer

An explicit method is stable only when the step times the largest eigenvalue lies in a bounded region. The step is then tied to the fast mode even after that mode has decayed and only the slow mode remains.

## The symptom

An explicit solver takes far more steps than the accuracy of the answer requires. The step is chosen by stability, and reducing the tolerance changes nothing until the stability limit is reached.

To check your understanding, identify the extreme eigenvalues, write the stability bound and observe the step count.

Stiffness is a property of the spectrum rather than of the solution, and the ratio is the number that predicts the wasted work.
