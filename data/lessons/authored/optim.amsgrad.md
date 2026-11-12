# AMSGrad and long-term memory

Adam's second moment is a decaying average, so a single small gradient can shrink the denominator and produce a large step. AMSGrad keeps the maximum of the second moment instead.

## The maximum

The denominator uses the largest decayed second moment seen so far. It therefore never decreases, which prevents the step from growing after a quiet period.

## The consequence

Steps become non-increasing for a fixed gradient magnitude. Convergence analyses that fail for Adam in some settings hold for the maximum form, at the cost of a slower response to genuinely easier regions.

## When it matters

The difference is largest in adversarial or sparse settings where the gradient scale changes sharply. In smooth, well-scaled problems the two methods behave almost identically.

To check your understanding, keep the largest second moment, make steps non-increasing and expect sparse settings.

AMSGrad is a reminder that a decaying average can hide the largest scale a coordinate has seen, and that the choice of memory shape is part of the method.
