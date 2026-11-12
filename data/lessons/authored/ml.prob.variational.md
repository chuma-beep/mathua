# Variational inference

Variational inference replaces an intractable posterior with the closest member of a simpler family. Closeness is measured by the Kullback-Leibler divergence, and the evidence lower bound turns minimisation into maximisation.

## The evidence lower bound

The bound is the expected log likelihood minus the divergence to the prior. Both terms are expectations under the approximate posterior.

## The gap

The difference between the log evidence and the bound is the divergence from the approximation to the true posterior. That gap is non-negative.

## Optimisation

Gradients of the bound are estimated by sampling from the approximation. The reparameterisation trick moves the randomness out of the parameters.

To check your understanding, name the two terms, state the gap and sample the approximation.

Variational inference trades exactness for a bound that can be optimised.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
